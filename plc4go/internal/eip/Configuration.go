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
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
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
	// Every option this driver reads goes through the reader, so the ones nothing read can be
	// reported rather than silently discarded. Deferred, so no return path can skip it.
	reader := spiOptions.NewOptionReader(localLogger, options)
	defer reader.ReportUnknown("eip")

	configuration := Configuration{
		backplane: 1,
		slot:      0,
		bigEndian: true,
	}
	if localRackString := reader.Get("backplane"); localRackString != "" {
		parsedBackplane, err := strconv.ParseInt(localRackString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing backplane")
		}
		configuration.backplane = int8(parsedBackplane)
	}
	if localSlotString := reader.Get("slot"); localSlotString != "" {
		parsedSlot, err := strconv.ParseInt(localSlotString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing slot")
		}
		configuration.slot = int8(parsedSlot)
	}
	if bigEndianString := getFromOptionsAliases(reader, "bigEndian", "big-endian"); bigEndianString != "" {
		parsedBigEndian, err := strconv.ParseBool(bigEndianString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing bigEndian")
		}
		configuration.bigEndian = parsedBigEndian
	}
	if forceUnconnectedString := getFromOptionsAliases(reader, "forceUnconnectedOperation", "force-unconnected-operation"); forceUnconnectedString != "" {
		parsedForceUnconnected, err := strconv.ParseBool(forceUnconnectedString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing forceUnconnectedOperation")
		}
		configuration.forceUnconnectedOperation = parsedForceUnconnected
	}
	configuration.communicationPath = getFromOptionsAliases(reader, "communicationPath", "communication-path")
	if serialNumberString := getFromOptionsAliases(reader, "connectionSerialNumber", "connection-serial-number"); serialNumberString != "" {
		parsedSerialNumber, err := strconv.ParseUint(serialNumberString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing connectionSerialNumber")
		}
		configuration.connectionSerialNumber = uint16(parsedSerialNumber)
	}
	return configuration, nil
}

func getFromOptionsAliases(reader *spiOptions.OptionReader, keys ...string) string {
	for _, key := range keys {
		if value := reader.Get(key); value != "" {
			return value
		}
	}
	return ""
}
