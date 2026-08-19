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

package eip

import (
	"strconv"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

type Configuration struct {
	backplane                 int8
	slot                      int8
	bigEndian                 bool
	forceUnconnectedOperation bool
	communicationPath         string
	connectionSerialNumber    uint16 // 0 = pick a random one per connection (test support option)
}

func ParseFromOptions(localLogger zerolog.Logger, options map[string][]string) (Configuration, error) {
	configuration := Configuration{
		backplane: 1,
		slot:      0,
		bigEndian: true,
	}
	if localRackString := getFromOptions(localLogger, options, "backplane"); localRackString != "" {
		parsedBackplane, err := strconv.ParseInt(localRackString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing backplane")
		}
		configuration.backplane = int8(parsedBackplane)
	}
	if localSlotString := getFromOptions(localLogger, options, "slot"); localSlotString != "" {
		parsedSlot, err := strconv.ParseInt(localSlotString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing slot")
		}
		configuration.slot = int8(parsedSlot)
	}
	if bigEndianString := getFromOptionsAliases(localLogger, options, "bigEndian", "big-endian"); bigEndianString != "" {
		parsedBigEndian, err := strconv.ParseBool(bigEndianString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing bigEndian")
		}
		configuration.bigEndian = parsedBigEndian
	}
	if forceUnconnectedString := getFromOptionsAliases(localLogger, options, "forceUnconnectedOperation", "force-unconnected-operation"); forceUnconnectedString != "" {
		parsedForceUnconnected, err := strconv.ParseBool(forceUnconnectedString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing forceUnconnectedOperation")
		}
		configuration.forceUnconnectedOperation = parsedForceUnconnected
	}
	configuration.communicationPath = getFromOptionsAliases(localLogger, options, "communicationPath", "communication-path")
	if serialNumberString := getFromOptionsAliases(localLogger, options, "connectionSerialNumber", "connection-serial-number"); serialNumberString != "" {
		parsedSerialNumber, err := strconv.ParseUint(serialNumberString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing connectionSerialNumber")
		}
		configuration.connectionSerialNumber = uint16(parsedSerialNumber)
	}
	return configuration, nil
}

func getFromOptions(localLogger zerolog.Logger, options map[string][]string, key string) string {
	if optionValues, ok := options[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLogger.Warn().Str("key", key).Msg("Options %s must be unique")
		}
		return optionValues[0]
	}
	return ""
}

func getFromOptionsAliases(localLogger zerolog.Logger, options map[string][]string, keys ...string) string {
	for _, key := range keys {
		if value := getFromOptions(localLogger, options, key); value != "" {
			return value
		}
	}
	return ""
}
