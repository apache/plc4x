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

package s7

import (
	"strconv"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
)

type Configuration struct {
	localRack      int32
	localSlot      int32
	remoteRack     int32
	remoteSlot     int32
	pduSize        uint16
	maxAmqCaller   uint16
	maxAmqCallee   uint16
	controllerType readWriteModel.ControllerType
}

func ParseFromOptions(localLog zerolog.Logger, options map[string][]string) (Configuration, error) {
	// Every option this driver reads goes through the reader, so the ones nothing read can be
	// reported rather than silently discarded. Deferred, so no return path can skip it.
	reader := spiOptions.NewOptionReader(localLog, options)
	defer reader.ReportUnknown("s7")

	configuration := Configuration{
		localRack:      1,
		localSlot:      1,
		remoteRack:     0,
		remoteSlot:     0,
		pduSize:        1024,
		maxAmqCaller:   8,
		maxAmqCallee:   8,
		controllerType: readWriteModel.ControllerType_ANY,
	}
	if localRackString := reader.Get("local-rack"); localRackString != "" {
		parsedInt, err := strconv.ParseInt(localRackString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing local-rack")
		}
		configuration.localRack = int32(parsedInt)
	}
	if localSlotString := reader.Get("local-slot"); localSlotString != "" {
		parsedInt, err := strconv.ParseInt(localSlotString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing local-slot")
		}
		configuration.localSlot = int32(parsedInt)
	}
	if remoteRackString := reader.Get("remote-rack"); remoteRackString != "" {
		parsedInt, err := strconv.ParseInt(remoteRackString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing remote-rack")
		}
		configuration.remoteRack = int32(parsedInt)
	}
	if remoteSlotString := reader.Get("remote-slot"); remoteSlotString != "" {
		parsedInt, err := strconv.ParseInt(remoteSlotString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing remote-slot")
		}
		configuration.remoteSlot = int32(parsedInt)
	}
	if controllerTypeString := reader.Get("controller-type"); controllerTypeString != "" {
		controllerType, ok := readWriteModel.ControllerTypeByName(controllerTypeString)
		if !ok {
			return Configuration{}, errors.Errorf("Unknown controller type %s", controllerTypeString)
		}
		configuration.controllerType = controllerType
	}

	pduSizeString := reader.Get("pdu-size")
	if pduSizeString != "" {
		parsedUint, err := strconv.ParseUint(pduSizeString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing pdu-size %s", pduSizeString)
		}
		configuration.pduSize = uint16(parsedUint)
	}

	if maxAmqCallerString := reader.Get("max-amq-caller"); maxAmqCallerString != "" {
		parsedUint, err := strconv.ParseUint(maxAmqCallerString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing max-amq-caller %s", maxAmqCallerString)
		}
		configuration.maxAmqCaller = uint16(parsedUint)
	}

	if maxAmqCalleeString := reader.Get("max-amq-callee"); maxAmqCalleeString != "" {
		parsedUint, err := strconv.ParseUint(maxAmqCalleeString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing max-amq-callee %s", maxAmqCalleeString)
		}
		configuration.maxAmqCallee = uint16(parsedUint)
	}
	return configuration, nil
}
