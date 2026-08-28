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
// **What this can and cannot do.** The two bindings share a specification, not a runtime: PLC4J
// resolves a connection string through ConfigurationFactory inside a JVM, PLC4Go parses it by
// hand in a Go binary, and no test process holds both. So parity is asserted by *shared cases*
// rather than by execution - the expectations below are the ones PLC4J's own tests assert for the
// same string, and the Java counterpart of this file names it in a comment so the pair is
// findable. Two copies of one table can drift; if that ever bites, the stronger form is a single
// checked-in fixture both bindings read as test data.
//
// This exists because the drift is real. PLC4Go's modbus driver accepted "unit-identifier" while
// PLC4J declared only "default-unit-identifier", so one connection string set the unit here and
// was silently ignored there - for as long as anyone had been reading the Go getting-started
// page, which documented exactly that string.
package configparity

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/modbus"
	"github.com/apache/plc4x/plc4go/internal/s7"
)

type parseOptions func(zerolog.Logger, map[string][]string) (interface{}, error)

// parseReporting runs a driver's option parsing and returns whatever it logged, which is where an
// unrecognised name is reported.
func parseReporting(t *testing.T, parse func(zerolog.Logger, map[string][]string) error,
	options map[string][]string) string {
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

// A canonical name is read by the driver, so nothing reports it. This is the property that makes
// one connection string mean one thing: PLC4J declares these names, and PLC4Go reads them.
func TestTheCanonicalNamesAreRecognisedHere(t *testing.T) {
	logged := parseReporting(t, modbusParse, map[string][]string{
		"default-unit-identifier":    {"3"},
		"request-timeout-ms":         {"5000"},
		"default-payload-byte-order": {"LITTLE_ENDIAN"},
	})
	assert.NotContains(t, logged, "not known", "every canonical modbus name must be read here")

	// s7, for a setting both bindings implement. Not every PLC4J parameter has a counterpart
	// here - PLC4Go's s7 has no S7H dual-path, so it reads no "ha-*" name, and reports one as
	// unknown, which is the truth: setting it here would do nothing. Parity is that a capability
	// *both* have is spelled the same, not that the capability sets match.
	logged = parseReporting(t, s7Parse, map[string][]string{
		"controller-type":  {"S7_1200"},
		"cotp.remote-rack": {"0"},
		"cotp.remote-slot": {"1"},
		"pdu-size":         {"1024"},
	})
	assert.NotContains(t, logged, "not known", "the s7 settings both bindings implement")
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
