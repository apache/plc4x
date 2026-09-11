/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package configparity checks that a connection string means the same thing here as it does in
// PLC4J.
//
// **Where the expectation comes from.** The two bindings share a specification, not a runtime:
// PLC4J resolves a connection string through ConfigurationFactory inside a JVM, PLC4Go parses it
// by hand in a Go binary, and no test process holds both. But PLC4J's half of the specification is
// declarative - every option it accepts is a @ConfigurationParameter annotation on a driver's
// configuration class - so it can be read out of the Java sources without building them. That is
// what this package does: it extracts the names PLC4J declares and asks each PLC4Go driver whether
// it reads them.
//
// Taking the expectation from PLC4J rather than listing it here is the whole point. A hand-written
// list of names on the Go side can only catch "PLC4Go stopped reading a name we listed"; it cannot
// catch "PLC4J gained an option PLC4Go never grew", because nobody adds the new name to a Go-side
// list. That second drift is the one that actually happened: max-coils-per-request and
// max-registers-per-request were declared, defaulted and documented in PLC4J from 0.13.0, and
// PLC4Go ignored them in silence for years.
//
// **What the extraction covers.** Every configuration class under plc4j/drivers - the driver's own
// Configuration, a transport configuration it declares (whose names carry the transport's code as
// a prefix, which is how a connection string writes them), and a class pulled in by
// @ComplexConfigurationParameter (whose names carry that parameter's prefix). It does not cover
// the parameters a driver inherits from a transport configuration living outside the driver tree -
// CotpTransportConfiguration's cotp.local-tsap and cotp.remote-tsap reach an s7 connection string
// too, and are not checked here. Widening to plc4j/transports would pull every transport's options
// into every driver's comparison, which is a different question from the one this package asks.
//
// This exists because the drift is real. PLC4Go's modbus driver accepted "unit-identifier" while
// PLC4J declared only "default-unit-identifier", so one connection string set the unit here and
// was silently ignored there - for as long as anyone had been reading the Go getting-started
// page, which documented exactly that string.
package configparity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/abeth"
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/eip"
	"github.com/apache/plc4x/plc4go/internal/firmata"
	"github.com/apache/plc4x/plc4go/internal/iec608705104"
	"github.com/apache/plc4x/plc4go/internal/modbus"
	"github.com/apache/plc4x/plc4go/internal/opcua"
	"github.com/apache/plc4x/plc4go/internal/s7"
	"github.com/apache/plc4x/plc4go/internal/slmp"
	"github.com/apache/plc4x/plc4go/internal/umas"
)

// plc4jDrivers is PLC4J's driver tree, from this package. Both bindings live in one repository, so
// the Java sources are always beside the Go ones; a checkout without them cannot answer the
// question this package asks, and says so rather than passing.
const plc4jDrivers = "../../../plc4j/drivers"

// configurationParameter matches one option declaration and the field it is attached to. PLC4J's
// ConfigurationFactory.getConfigurationName reads exactly this: the annotation's value, or - when
// it carries none, as EIPConfiguration's backplane and slot do - the name of the field itself. So
// the field name is captured too, and used when the annotation is bare.
var configurationParameter = regexp.MustCompile(
	`@ConfigurationParameter(?:\("([^"]*)"\))?(?s:.*?)(?:private|protected|public)\s+[\w.<>\[\], ]+?\s+(\w+)\s*[=;]`)

// complexConfigurationParameter matches a parameter whose options live on another class, and
// captures the prefix they are addressed under together with that class's name. OPC UA's encoding
// limits are declared this way: Limits.java says "receive-buffer-size", and a connection string
// says "encoding.receive-buffer-size".
var complexConfigurationParameter = regexp.MustCompile(
	`@ComplexConfigurationParameter\(prefix\s*=\s*"([^"]+)"(?s:.*?)\)(?s:.*?)(?:private|protected|public)\s+([\w.]+)(?:<[^>]*>)?\s+\w+\s*[=;]`)

// transportPackage recovers the transport a *TransportConfiguration belongs to from the import of
// the class it extends. PLC4J resolves a transport's configuration under that transport's code as
// a prefix, so S7CotpTransportConfiguration's "local-rack" is written "cotp.local-rack" in a
// connection string - which is the spelling PLC4Go's s7 driver reads, and the one the
// documentation uses. Without this the prefixed names would be reported as gaps that are not real.
var transportPackage = regexp.MustCompile(`import org\.apache\.plc4x\.java\.transport\.([a-zA-Z0-9]+)\.config\.`)

// canaryOption is a name no driver can read. Probing it first proves the driver reports what it
// did not read; without that, "nothing was reported" would not distinguish a driver that reads
// every option from a driver that reports nothing at all, and the whole comparison would pass
// vacuously.
const canaryOption = "plc4x-config-parity-canary"

// probeValue is what every probe sets the option to. Its value does not matter: an option is read
// if the driver asked for it, and a driver that asks and then rejects the value has still asked.
const probeValue = "1"

type parseOptions func(zerolog.Logger, map[string][]string) error

// comparedDriver is one driver both bindings implement.
type comparedDriver struct {
	// parse is PLC4Go's entry point for this driver's connection options. The driver's key in
	// comparedDrivers is its directory under plc4j/drivers, which is where the declarations are.
	parse parseOptions
	// notReadHere are the options PLC4J declares that this binding does not report as read, each
	// with the reason. An entry here is a divergence someone looked at and decided to keep; an
	// option that is neither read nor listed is a divergence nobody decided to have, and fails.
	notReadHere map[string]string
}

// comparedDrivers are the drivers whose option vocabularies are compared. Every PLC4J driver that
// declares an option is either here or in notCompared, and TestEveryPlc4jDriverIsAccountedFor
// insists on it, so a new Java driver cannot slip past by being in neither.
var comparedDrivers = map[string]comparedDriver{
	"ab-eth": {
		parse:       wrap(abeth.ParseFromOptions),
		notReadHere: map[string]string{},
	},
	"bacnet": {
		parse: wrap(bacnetip.ParseFromOptions),
		notReadHere: map[string]string{
			"ede-file-path":      "plc4j names objects from an EDE (Engineering Data Exchange) file; plc4go's bacnet-ip has no EDE support",
			"ede-directory-path": "same: no EDE support here, so there is no directory of EDE files to point at",
		},
	},
	"c-bus": {
		parse: wrap(cbus.ParseFromOptions),
		notReadHere: map[string]string{
			"srchk": "honoured here but reported anyway - the driver title-cases the supplied keys and then consumes only the title-cased copy, so the lower-case spelling plc4j declares is applied AND warned about. A plc4go reporting bug rather than a vocabulary gap; fixing it removes this entry",
		},
	},
	"eip": {
		parse: wrap(eip.ParseFromOptions),
		notReadHere: map[string]string{
			"request-timeout-ms": "plc4go's eip connection has no configurable request timeout to bind it to",
		},
	},
	"firmata": {
		parse:       wrap(firmata.ParseFromOptions),
		notReadHere: map[string]string{},
	},
	"iec-60870": {
		parse:       wrap(iec608705104.ParseFromOptions),
		notReadHere: map[string]string{},
	},
	"modbus": {
		parse:       wrap(modbus.ParseFromOptions),
		notReadHere: map[string]string{},
	},
	"opcua": {
		parse: wrap(opcua.ParseFromOptions),
		notReadHere: map[string]string{
			// Plumbing rather than device settings.
			"protocol-code":    "not a device setting: plc4j lets the string restate the protocol code it was parsed from; plc4go takes it from the URL",
			"transport-code":   "not a device setting: as protocol-code, for the transport",
			"transport-config": "not a device setting: as protocol-code, for the transport's own configuration string",
			// Certificate handling. plc4go's secure channel verifies no server certificate, which
			// is why AllowUnverifiedSecurityPolicies exists; everything a trust store would feed
			// has nothing to feed.
			"message-security":         "plc4go derives the message security mode from whether the policy encrypts (SignAndEncrypt or None) and cannot be told a different one",
			"tls.trust-store":          "plc4go's opcua does not verify server certificates, so it has no trust store",
			"tls.trust-store-type":     "no trust store here, so no type for one",
			"tls.trust-store-password": "no trust store here, so no password for one",
			"tls.verify":               "plc4go's opcua does not verify server certificates at all; the opt-in it does have is allow-unverified-security-policies",
			"server-certificate-file":  "same: nothing here reads a server certificate to check against",
			"tls.keystore-type":        "plc4go reads a PEM key store only, so there is no type to choose",
			"generated-key-size":       "plc4go's generated client key is a fixed 4096-bit RSA key (CertificateGenerator.go)",
			// Features plc4go's opcua does not have.
			"allow-insecure-credentials":     "plc4go has no guard against sending credentials over an unsecured channel, so there is nothing to allow",
			"browse-max-references-per-node": "plc4go's opcua has no browser",
			"browse-max-total-nodes":         "plc4go's opcua has no browser",
			"browse-max-depth":               "plc4go's opcua has no browser",
			"subscription-queue-size":        "plc4go's opcua subscriber does not set a server-side queue depth",
			// Timings plc4go fixes rather than configures.
			"channel-lifetime-ms":     "plc4go's secure channel uses DEFAULT_CONNECTION_LIFETIME and whatever the server revises it to",
			"min-channel-lifetime-ms": "no configurable lifetime here, so no floor to put under it",
			"session-timeout-ms":      "plc4go's session timeout is not configurable",
			"handshake-timeout-ms":    "plc4go does not bound the open-channel/create-session/close steps separately",
			"request-timeout-ms":      "plc4go's opcua read/write/subscribe calls have no configurable timeout",
			// Endpoint overrides for a server that advertises something other than the address it
			// was reached on.
			"endpoint-host": "plc4go builds the endpoint from the connection URL and cannot be told a different host",
			"endpoint-port": "plc4go builds the endpoint from the connection URL and cannot be told a different port",
			// The TCP encoding limits, declared on Limits.java under the "encoding" prefix.
			"encoding.receive-buffer-size": "plc4go's opcua does not negotiate the TCP encoding limits",
			"encoding.send-buffer-size":    "plc4go's opcua does not negotiate the TCP encoding limits",
			"encoding.max-message-size":    "plc4go's opcua does not negotiate the TCP encoding limits",
			"encoding.max-chunk-count":     "plc4go's opcua does not negotiate the TCP encoding limits",
		},
	},
	"s7": {
		parse: wrap(s7.ParseFromOptions),
		notReadHere: map[string]string{
			"cotp.local-device-group":  "plc4go hardcodes OTHERS when it encodes the calling TSAP (DriverContext.go)",
			"cotp.remote-device-group": "plc4go hardcodes PG_OR_PC when it encodes the called TSAP (DriverContext.go)",
			"ha-heartbeat-interval-ms": "plc4go's s7 has no S7H dual-path, so there is no heartbeat to pace",
			"ha-failover-timeout-ms":   "plc4go's s7 has no S7H dual-path, so there is nothing to fail over to",
			"read-timeout-ms":          "plc4go's s7 has no configurable per-exchange timeout",
		},
	},
	"slmp": {
		parse:       wrap(slmp.ParseFromOptions),
		notReadHere: map[string]string{},
	},
	"umas": {
		parse: wrap(umas.ParseFromOptions),
		notReadHere: map[string]string{
			"browser-generate-array-nodes": "plc4go's umas browser always reports an array as one tag carrying ArrayInfo and never generates a node per element, so there is nothing to switch off",
		},
	},
}

// notCompared are the PLC4J drivers that declare options and are still not compared, with the
// reason. Being listed here says the comparison cannot be made - not that the two agree - so the
// reasons say which of the two it is.
var notCompared = map[string]string{
	"ads":           "plc4go's ads driver parses its options without an OptionReader, so it reports nothing as unknown and the probe below could not tell a read option from an ignored one",
	"knxnetip":      "plc4go's knxnetip driver has no connection-option parsing at all, so it reads none of the five options plc4j declares",
	"can":           "no plc4go driver",
	"canopen":       "no plc4go driver",
	"ctrlx":         "no plc4go driver",
	"open-protocol": "no plc4go driver",
	"plc4x":         "no plc4go driver",
	"profinet":      "no plc4go driver",
	"profinet-ng":   "no plc4go driver",
}

func wrap[T any](parse func(zerolog.Logger, map[string][]string) (T, error)) parseOptions {
	return func(log zerolog.Logger, options map[string][]string) error {
		_, err := parse(log, options)
		return err
	}
}

// parseReporting runs a driver's option parsing and returns whatever it logged, which is where an
// unrecognised name is reported.
func parseReporting(t *testing.T, parse parseOptions, options map[string][]string) string {
	t.Helper()
	var logged bytes.Buffer
	require.NoError(t, parse(zerolog.New(&logged), options))
	return logged.String()
}

func modbusParse(log zerolog.Logger, options map[string][]string) error {
	_, err := modbus.ParseFromOptions(log, options)
	return err
}

func s7Parse(log zerolog.Logger, options map[string][]string) error {
	_, err := s7.ParseFromOptions(log, options)
	return err
}

// isReadBy says whether a driver reads the given option name.
//
// It asks the driver rather than reading its source: the option is offered on its own, and the
// driver's own report of what it did not read is the answer. Offering it alone is what makes the
// answer trustworthy - no other option can make the parse give up early and leave a later one
// unread - and the report survives a parse error, because every driver defers it.
func isReadBy(parse parseOptions, option string) bool {
	var logged bytes.Buffer
	// The error is deliberately dropped: a driver that rejects probeValue has still asked for the
	// option, which is the question being asked.
	_ = parse(zerolog.New(&logged), map[string][]string{option: {probeValue}})
	for _, line := range strings.Split(logged.String(), "\n") {
		var event map[string]any
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		message, _ := event["message"].(string)
		if event["option"] == option && strings.Contains(message, "not known") {
			return false
		}
	}
	return true
}

// plc4jOptionsOf is every option name PLC4J declares for one driver, mapped to the file declaring
// it so a failure can point at the source rather than at a bare name.
func plc4jOptionsOf(t *testing.T, javaDir string) map[string]string {
	t.Helper()
	root := filepath.Join(plc4jDrivers, javaDir, "src", "main", "java")
	require.DirExists(t, root,
		"the plc4j sources are what this package compares against; without them it cannot answer")

	sources := configurationSourcesIn(t, root)

	// A class named by a @ComplexConfigurationParameter declares its options under that
	// parameter's prefix, wherever it lives in the same package.
	prefixOfClass := map[string]string{}
	for _, source := range sources {
		for _, match := range complexConfigurationParameter.FindAllSubmatch(source.content, -1) {
			prefixOfClass[string(match[2])] = string(match[1]) + "."
		}
	}

	declared := map[string]string{}
	for _, source := range sources {
		class := strings.TrimSuffix(source.name, ".java")
		prefix, complex := prefixOfClass[class]
		switch {
		case complex:
		case !strings.HasSuffix(source.name, "Configuration.java"):
			// Not a configuration and not pulled in as one.
			continue
		case strings.HasSuffix(source.name, "TransportConfiguration.java"):
			transport := transportPackage.FindSubmatch(source.content)
			require.NotNil(t, transport, "%s declares options under a transport's prefix, but names "+
				"no org.apache.plc4x.java.transport.<code>.config parent to take that prefix from",
				source.name)
			prefix = string(transport[1]) + "."
		}
		for _, match := range configurationParameter.FindAllSubmatch(source.content, -1) {
			// PLC4J falls back to the field name when the annotation carries no name of its own.
			name := string(match[1])
			if name == "" {
				name = string(match[2])
			}
			declared[prefix+name] = source.name
		}
	}
	require.NotEmpty(t, declared, "plc4j's %s driver declares no option at all, which is odd enough "+
		"to be a bug in this extraction rather than in the driver", javaDir)
	return declared
}

type javaSource struct {
	name    string
	content []byte
}

// configurationSourcesIn reads every Java file in a "config" or "configuration" package under the
// given root, which is where PLC4J keeps a driver's configuration classes.
func configurationSourcesIn(t *testing.T, root string) []javaSource {
	t.Helper()
	var sources []javaSource
	require.NoError(t, filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".java") {
			return err
		}
		switch filepath.Base(filepath.Dir(path)) {
		case "config", "configuration":
		default:
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sources = append(sources, javaSource{name: entry.Name(), content: content})
		return nil
	}))
	return sources
}

// Every option PLC4J declares is either read by the PLC4Go driver of the same name or written down
// as a deliberate divergence. This is the direction a hand-written list cannot check: the names
// come from PLC4J, so an option PLC4J grows and PLC4Go never grows shows up here by itself.
func TestEveryPlc4jOptionIsReadHereOrDeliberatelyNot(t *testing.T) {
	for driverName, driver := range comparedDrivers {
		t.Run(driverName, func(t *testing.T) {
			require.False(t, isReadBy(driver.parse, canaryOption),
				"%s does not report an option it never read, so this probe cannot tell a read "+
					"option from an ignored one and every comparison below would pass vacuously. "+
					"Either give the driver an OptionReader, or move it to notCompared with the "+
					"reason", driverName)

			declared := plc4jOptionsOf(t, driverName)
			names := make([]string, 0, len(declared))
			for name := range declared {
				names = append(names, name)
			}
			sort.Strings(names)

			var gaps []string
			for _, option := range names {
				if isReadBy(driver.parse, option) {
					assert.NotContains(t, driver.notReadHere, option,
						"%s reads %q now, so it is no longer a divergence - drop it from notReadHere",
						driverName, option)
					continue
				}
				if reason, deliberate := driver.notReadHere[option]; deliberate {
					assert.NotEmpty(t, reason, "the notReadHere entry for %q must say why", option)
					continue
				}
				gaps = append(gaps, fmt.Sprintf("  %-38s declared in %s", option, declared[option]))
			}
			assert.Empty(t, gaps, "plc4j's %s driver declares options this binding does not read, "+
				"and nobody wrote down that it should not:\n%s\n"+
				"Setting one of these in a connection string works in plc4j and does nothing here. "+
				"Either read it in plc4go/internal/... (the fix), or add it to the %q entry of "+
				"comparedDrivers with a one-line reason (the decision). Do not delete the name from "+
				"this test - it is not maintained here, it is read out of the Java sources.",
				driverName, strings.Join(gaps, "\n"), driverName)

			for option := range driver.notReadHere {
				assert.Contains(t, declared, option,
					"notReadHere lists %q, which plc4j's %s driver no longer declares - drop it",
					option, driverName)
			}
		})
	}
}

// Every PLC4J driver that declares an option is classified: compared, or explicitly not. A driver
// in neither list would be invisible here, which is the same silence this package exists to break,
// so a new one fails until someone says which it is.
func TestEveryPlc4jDriverIsAccountedFor(t *testing.T) {
	entries, err := os.ReadDir(plc4jDrivers)
	require.NoError(t, err, "the plc4j sources are what this package compares against")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		javaDir := entry.Name()
		root := filepath.Join(plc4jDrivers, javaDir, "src", "main", "java")
		if _, err := os.Stat(root); err != nil {
			continue
		}
		declaresAnything := false
		for _, source := range configurationSourcesIn(t, root) {
			if bytes.Contains(source.content, []byte("@ConfigurationParameter")) {
				declaresAnything = true
				break
			}
		}
		if !declaresAnything {
			continue
		}
		_, compared := comparedDrivers[javaDir]
		reason, excused := notCompared[javaDir]
		assert.True(t, compared != excused,
			"plc4j's %s driver declares connection-string options, so it belongs in exactly one of "+
				"comparedDrivers (with the plc4go entry point that reads them) or notCompared (with "+
				"the reason it cannot be compared)", javaDir)
		if excused {
			assert.NotEmpty(t, reason, "the notCompared entry for %s must say why", javaDir)
		}
	}
}

// The rack and slot carry the "cotp." prefix in PLC4J, which declares them on the COTP transport's
// configuration, and every s7 example in the documentation spells them that way. This binding read
// them unprefixed, so the documented connection string set nothing here and said nothing about it.
// This is the divergence this package was written to find; it is fixed, and pinned here.
func TestTheUnprefixedS7RackAndSlotAreNotAccepted(t *testing.T) {
	logged := parseReporting(t, s7Parse, map[string][]string{
		"remote-rack": {"0"}, "remote-slot": {"1"},
	})
	assert.Contains(t, logged, "remote-rack")
	assert.Contains(t, logged, "remote-slot")
}

// The other half of that: a PLC4J parameter this binding does not implement is reported, so an
// operator finds out the setting does nothing here rather than believing it applied.
func TestAParameterThisBindingDoesNotImplementIsReported(t *testing.T) {
	logged := parseReporting(t, s7Parse, map[string][]string{"ha-heartbeat-interval-ms": {"4000"}})
	assert.Contains(t, logged, "ha-heartbeat-interval-ms")
	assert.Contains(t, logged, "not known",
		"PLC4Go's s7 has no S7H dual-path, and says so rather than accepting the setting")
}

// A spelling only this binding ever had must not be silently accepted, or the same string sets a
// value here and is ignored in PLC4J. This is the divergence that prompted the test: modbus read
// "unit-identifier", which PLC4J never declared, while UMAS uses that name for a different thing.
func TestAGoOnlySpellingIsReported(t *testing.T) {
	logged := parseReporting(t, modbusParse, map[string][]string{"unit-identifier": {"9"}})
	assert.Contains(t, logged, "unit-identifier")
	assert.Contains(t, logged, "not known", "it must be reported rather than quietly applied")
}

// The pre-migration names are unknown in both bindings, so neither accepts what the other
// rejects. PLC4J asserts the same list in DriverBaseUnknownParameterTest.
func TestPreMigrationNamesAreReported(t *testing.T) {
	for _, old := range []string{"request-timeout", "read-timeout", "connect-timeout"} {
		t.Run(old, func(t *testing.T) {
			logged := parseReporting(t, modbusParse, map[string][]string{old: {"1234"}})
			assert.Contains(t, logged, old, "the old spelling must be named")
		})
	}
}
