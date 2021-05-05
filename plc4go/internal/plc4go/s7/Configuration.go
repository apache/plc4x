//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package s7

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Configuration struct {
	localRack      int32
	localSlot      int32
	remoteRack     int32
	remoteSlot     int32
	pduSize        uint16
	maxAmqCaller   uint16
	maxAmqCallee   uint16
	controllerType ControllerType
}

func ParseFromOptions(options map[string][]string) (Configuration, error) {
	configuration := Configuration{
		localRack:      1,
		localSlot:      1,
		remoteRack:     0,
		remoteSlot:     0,
		pduSize:        1024,
		maxAmqCaller:   8,
		maxAmqCallee:   8,
		controllerType: ControllerType_UNKNOWN,
	}
	if localRackString := getFromOptions(options, "local-rack"); localRackString != "" {
		atoi, err := strconv.Atoi(localRackString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing local-rack")
		}
		configuration.localRack = int32(atoi)
	}
	if localSlotString := getFromOptions(options, "local-rack"); localSlotString != "" {
		atoi, err := strconv.Atoi(localSlotString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing local-slot")
		}
		configuration.localSlot = int32(atoi)
	}
	if remoteRackString := getFromOptions(options, "remote-rack"); remoteRackString != "" {
		atoi, err := strconv.Atoi(remoteRackString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing remote-rack")
		}
		configuration.remoteRack = int32(atoi)
	}
	if remoteSlotString := getFromOptions(options, "remote-rack"); remoteSlotString != "" {
		atoi, err := strconv.Atoi(remoteSlotString)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing remote-slot")
		}
		configuration.remoteSlot = int32(atoi)
	}
	if controllerTypeString := getFromOptions(options, "controller-type"); controllerTypeString != "" {
		switch controllerTypeString {
		case "ANY":
			configuration.controllerType = ControllerType_ANY
		case "S7_300":
			configuration.controllerType = ControllerType_S7_300
		case "S7_400":
			configuration.controllerType = ControllerType_S7_400
		case "S7_1200":
			configuration.controllerType = ControllerType_S7_1200
		case "S7_1500":
			configuration.controllerType = ControllerType_S7_1500
		case "LOGO":
			configuration.controllerType = ControllerType_LOGO
		default:
			return Configuration{}, errors.Errorf("Unknown controller type %s", controllerTypeString)
		}
	}

	pduSizeString := getFromOptions(options, "pdu-size")
	if pduSizeString != "" {
		atoi, err := strconv.Atoi(pduSizeString)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing pdu-size %s", pduSizeString)
		}
		configuration.pduSize = uint16(atoi)
	}

	if maxAmqCallerString := getFromOptions(options, "max-amq-caller"); maxAmqCallerString != "" {
		atoi, err := strconv.Atoi(maxAmqCallerString)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing max-amq-caller %s", maxAmqCallerString)
		}
		configuration.maxAmqCaller = uint16(atoi)
	}

	if maxAmqCalleeString := getFromOptions(options, "max-amq-callee"); maxAmqCalleeString != "" {
		atoi, err := strconv.Atoi(maxAmqCalleeString)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing max-amq-callee %s", maxAmqCalleeString)
		}
		configuration.maxAmqCallee = uint16(atoi)
	}
	return configuration, nil
}

func getFromOptions(options map[string][]string, key string) string {
	if optionValues, ok := options[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			log.Warn().Msgf("Options %s must be unique", key)
		}
		return optionValues[0]
	}
	return ""
}
