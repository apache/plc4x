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

package opcua

import (
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
)

//go:generate go tool plc4xGenerator -type=Configuration
type Configuration struct {
	Code              string
	Host              string
	Port              string
	Endpoint          string
	TransportEndpoint string
	Params            string
	IsEncrypted       bool
	Thumbprint        readWriteModel.PascalByteString
	SenderCertificate []byte
	Discovery         bool
	Username          string
	// The credentials below render as <redacted>, never as their values - see the generator's
	// secret tag. Username is not marked: it says who is connecting, which is a diagnostic.
	Password         string `secret:"true"`
	SecurityPolicy   string
	KeyStoreFile     string
	CertDirectory    string
	KeyStorePassword string `secret:"true"`
	// Ckp holds the client's private key. Rendering it prints pointer addresses rather than key
	// material, so nothing leaks today - but a struct that holds a key by value later would, and
	// a field whose contents must never be printed should say so where it is declared.
	Ckp *CertificateKeyPair `secret:"true"`
	// AllowUnverifiedSecurityPolicies is an explicit opt-in for security policies other than "None".
	// The plc4go OPC UA secure-channel implementation does not verify server certificates or message
	// signatures yet, so any other policy is refused unless this is set to true.
	AllowUnverifiedSecurityPolicies bool

	log zerolog.Logger
}

// fieldForOption maps a connection-string option to the Configuration field it sets.
//
// The names are the ones PLC4J declares and the documentation lists, not this struct's Go field
// names. Deriving them from the field names is what made this driver read "keyStoreFile" and
// "keyStorePassword": the documented "tls.keystore" and "tls.keystore-password" reached it as
// unknown options and were ignored, so the two bindings disagreed about the same string.
//
// Keys are lower-case; lookup lower-cases the option, which keeps the case-insensitive matching
// this driver has always had.
var fieldForOption = map[string]string{
	"discovery":             "Discovery",
	"username":              "Username",
	"password":              "Password",
	"security-policy":       "SecurityPolicy",
	"tls.keystore":          "KeyStoreFile",
	"tls.keystore-password": "KeyStorePassword",
	// No PLC4J counterpart: this binding generates its certificate into a directory rather than
	// taking a key store, so the name is this binding's own - in the shared spelling.
	"cert-directory":                     "CertDirectory",
	"allow-unverified-security-policies": "AllowUnverifiedSecurityPolicies",
}

func ParseFromOptions(log zerolog.Logger, options map[string][]string) (Configuration, error) {
	configuration := createDefaultConfiguration()
	reflectConfiguration := reflect.ValueOf(&configuration).Elem()
	for optionKey := range options {
		fieldName, ok := fieldForOption[strings.ToLower(optionKey)]
		if !ok {
			// Read by the transport rather than by this driver. The names come from the
			// transports themselves, which register what they read, so this driver does not
			// keep its own copy of a list that would drift from them.
			if spiOptions.IsTransportOption(optionKey) {
				continue
			}
			// Warn rather than fail, matching plc4j, which reports an unknown parameter and
			// carries on for every driver. This driver used to be the only one anywhere in
			// PLC4X that refused the connection, which meant one connection string was
			// accepted by plc4j and rejected here.
			//
			// The reason it refused is still real and is now carried by the warning instead:
			// a typo in a security-relevant option (securityPolicy, allowUnverified...) falls
			// back to a default, and the operator has to see that it did. An unread option is
			// reported by name, so it is visible - but it no longer stops the connection.
			log.Warn().
				Str("option", optionKey).
				Msg("Connection string option is not known to the opcua driver and is ignored")
			continue
		}
		optionValue := getFromOptions(log, options, optionKey)
		if optionValue == "" {
			continue
		}
		field := reflectConfiguration.FieldByName(fieldName)
		switch field.Kind() {
		case reflect.String:
			field.SetString(optionValue)
		case reflect.Uint8:
			parseUint, err := strconv.ParseUint(optionValue, 0, 8)
			if err != nil {
				return Configuration{}, errors.Wrapf(err, "Error parsing %s", optionKey)
			}
			field.SetUint(parseUint)
		case reflect.Bool:
			parseBool, err := strconv.ParseBool(optionValue)
			if err != nil {
				return Configuration{}, errors.Wrapf(err, "Error parsing %s", optionKey)
			}
			field.SetBool(parseBool)
		default:
			return Configuration{}, errors.Errorf("%s not yet supported", field.Kind())
		}
	}
	configuration.log = log
	if err := configuration.validateSecurityPolicy(); err != nil {
		return Configuration{}, err
	}
	return configuration, nil
}

// validateSecurityPolicy refuses security policies the driver cannot actually enforce: the
// secure-channel implementation performs no certificate or message-signature verification yet,
// so requesting anything but "None" must fail closed unless the user explicitly opted in.
func (c *Configuration) validateSecurityPolicy() error {
	if c.SecurityPolicy == "" || c.SecurityPolicy == "None" {
		return nil
	}
	if !c.AllowUnverifiedSecurityPolicies {
		return errors.Errorf("security-policy %s is not supported: the plc4go OPC UA driver does not verify server certificates or message signatures yet. "+
			"Set allow-unverified-security-policies=true to connect anyway (NOT recommended for production use)", c.SecurityPolicy)
	}
	c.log.Warn().
		Str("securityPolicy", c.SecurityPolicy).
		Msg("allowUnverifiedSecurityPolicies is set: the requested security policy is used WITHOUT certificate or message-signature verification")
	return nil
}

func (c *Configuration) openKeyStore() error {
	c.IsEncrypted = true
	securityTempDir := path.Join(c.CertDirectory, "security")
	if _, err := os.Stat(securityTempDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(securityTempDir, 700); err != nil {
			return errors.New("Unable to create directory please confirm folder permissions on " + securityTempDir)
		}
	}

	serverKeyStore := path.Join(securityTempDir, c.KeyStoreFile)
	if _, err := os.Stat(securityTempDir); errors.Is(err, os.ErrNotExist) {
		var err error
		c.Ckp, err = generateCertificate()
		if err != nil {
			return errors.Wrap(err, "error generating certificate")
		}
		c.log.Info().Str("serverKeyStore", serverKeyStore).Msg("Creating keystore")
		// TODO: not sure how to do that in golang. Seems pkc12 can only decode for now
		_ = os.WriteFile(serverKeyStore, []byte{0xA}, 0700)
	} else {
		c.log.Info().Str("serverKeyStore", serverKeyStore).Msg("Loading keystore")
		serverKeyStoreContent, err := os.ReadFile(serverKeyStore)
		if err != nil {
			return errors.Wrap(err, "error reading "+serverKeyStore)
		}
		// TODO: here we can parse with "golang.org/x/crypto/pkcs12" Decode
		_ = serverKeyStoreContent
	}

	return nil
}

func createDefaultConfiguration() Configuration {
	return Configuration{
		SecurityPolicy: "None",
	}
}

func getFromOptions(localLog zerolog.Logger, options map[string][]string, key string) string {
	if optionValues, ok := options[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLog.Warn().Str("key", key).Msg("Options must be unique")
		}
		return optionValues[0]
	}
	return ""
}
