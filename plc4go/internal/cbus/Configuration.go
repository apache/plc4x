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

package cbus

import (
	"github.com/rs/zerolog"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=Configuration
type Configuration struct {
	Srchk    bool
	Exstat   bool
	Pun      bool
	LocalSal bool
	Pcn      bool
	Idmon    bool
	Monitor  bool
	Smart    bool
	XonXoff  bool
	Connect  bool

	MonitoredApplication1 byte
	MonitoredApplication2 byte
}

func ParseFromOptions(log zerolog.Logger, options map[string][]string) (Configuration, error) {
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
	return configuration, nil
}

func createDefaultConfiguration() Configuration {
	return Configuration{
		Exstat:   true,
		LocalSal: true,
		Idmon:    true,
		Monitor:  true,
		Smart:    true,
		Srchk:    true,
		Connect:  true,

		MonitoredApplication1: 0xFF,
		MonitoredApplication2: 0xFF,
	}
}

func getFromOptions(localLog zerolog.Logger, options map[string][]string, key string) string {
	if optionValues, ok := options[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLog.Warn().Msgf("Options %s must be unique", key)
		}
		return optionValues[0]
	}
	return ""
}
