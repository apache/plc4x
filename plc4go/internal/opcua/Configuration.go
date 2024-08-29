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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=Configuration
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
	Password          string
	SecurityPolicy    string
	KeyStoreFile      string
	CertDirectory     string
	KeyStorePassword  string
	Ckp               *CertificateKeyPair

	log zerolog.Logger `ignore:"true"`
}

func ParseFromOptions(log zerolog.Logger, options map[string][]string) (Configuration, error) {
	titleOptions(options)
	configuration := createDefaultConfiguration()
	reflectConfiguration := reflect.ValueOf(&configuration).Elem()
	for i := 0; i < reflectConfiguration.NumField(); i++ {
		field := reflectConfiguration.Type().Field(i)
		key := field.Name
		if optionValue := getFromOptions(log, options, key); optionValue != "" {
			switch field.Type.Kind() {
			case reflect.Uint8:
				parseUint, err := strconv.ParseUint(optionValue, 0, 8)
				if err != nil {
					return Configuration{}, errors.Wrapf(err, "Error parsing %s", key)
				}
				reflectConfiguration.FieldByName(key).SetUint(parseUint)
			case reflect.Bool:
				parseBool, err := strconv.ParseBool(optionValue)
				if err != nil {
					return Configuration{}, errors.Wrapf(err, "Error parsing %s", key)
				}
				reflectConfiguration.FieldByName(key).SetBool(parseBool)
			default:
				return configuration, errors.Errorf("%s not yet supported", field.Type.Kind())
			}
		}
	}
	configuration.log = log
	return configuration, nil
}

func titleOptions(options map[string][]string) {
	caser := cases.Title(language.AmericanEnglish)
	for key, value := range options {
		options[caser.String(key)] = value
	}
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
