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

package ads

import (
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type Configuration struct {
	sourceAmsNetId readWriteModel.AmsNetId
	sourceAmsPort  uint16
	targetAmsNetId readWriteModel.AmsNetId
	targetAmsPort  uint16
}

func ParseFromOptions(options map[string][]string) (Configuration, error) {
	configuration := Configuration{}

	sourceAmsNetId := getFromOptions(options, "sourceAmsNetId")
	if sourceAmsNetId == "" {
		return Configuration{}, errors.New("Required parameter sourceAmsNetId missing")
	}
	split := strings.Split(sourceAmsNetId, ".")
	octet1, err := strconv.ParseUint(split[0], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet2, err := strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet3, err := strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet4, err := strconv.ParseUint(split[3], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet5, err := strconv.ParseUint(split[4], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet6, err := strconv.ParseUint(split[5], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	configuration.sourceAmsNetId = readWriteModel.AmsNetId{
		Octet1: uint8(octet1),
		Octet2: uint8(octet2),
		Octet3: uint8(octet3),
		Octet4: uint8(octet4),
		Octet5: uint8(octet5),
		Octet6: uint8(octet6),
	}
	sourceAmsPort := getFromOptions(options, "sourceAmsPort")
	if sourceAmsPort == "" {
		return Configuration{}, errors.New("Required parameter sourceAmsPort missing")
	}
	parsedUint, err := strconv.ParseUint(sourceAmsPort, 10, 16)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsPort")
	}
	configuration.sourceAmsPort = uint16(parsedUint)
	targetAmsNetId := getFromOptions(options, "targetAmsNetId")
	if sourceAmsNetId == "" {
		return Configuration{}, errors.New("Required parameter targetAmsNetId missing")
	}
	split = strings.Split(targetAmsNetId, ".")
	octet1, err = strconv.ParseUint(split[0], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet2, err = strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet3, err = strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet4, err = strconv.ParseUint(split[3], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet5, err = strconv.ParseUint(split[4], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet6, err = strconv.ParseUint(split[5], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	configuration.targetAmsNetId = readWriteModel.AmsNetId{
		Octet1: uint8(octet1),
		Octet2: uint8(octet2),
		Octet3: uint8(octet3),
		Octet4: uint8(octet4),
		Octet5: uint8(octet5),
		Octet6: uint8(octet6),
	}
	targetAmsPort := getFromOptions(options, "targetAmsPort")
	if targetAmsPort == "" {
		return Configuration{}, errors.New("Required parameter targetAmsPort missing")
	}
	parsedUint, err = strconv.ParseUint(targetAmsPort, 10, 16)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsPort")
	}
	configuration.targetAmsPort = uint16(parsedUint)

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
