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

package model

import (
	"github.com/rs/zerolog"
	"strconv"
	"strings"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/pkg/errors"
)

type Configuration struct {
	SourceAmsNetId readWriteModel.AmsNetId
	SourceAmsPort  uint16
	TargetAmsNetId readWriteModel.AmsNetId
	TargetAmsPort  uint16
}

func ParseFromOptions(localLogger zerolog.Logger, options map[string][]string) (Configuration, error) {
	configuration := Configuration{}

	sourceAmsNetId := getFromOptions(localLogger, options, "sourceAmsNetId")
	if sourceAmsNetId == "" {
		return Configuration{}, errors.New("Required parameter sourceAmsNetId missing")
	}
	split := strings.Split(sourceAmsNetId, ".")
	octet1, err := strconv.ParseUint(split[0], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing SourceAmsNetId")
	}
	octet2, err := strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing SourceAmsNetId")
	}
	octet3, err := strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing SourceAmsNetId")
	}
	octet4, err := strconv.ParseUint(split[3], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing SourceAmsNetId")
	}
	octet5, err := strconv.ParseUint(split[4], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing SourceAmsNetId")
	}
	octet6, err := strconv.ParseUint(split[5], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing SourceAmsNetId")
	}
	configuration.SourceAmsNetId = readWriteModel.NewAmsNetId(
		uint8(octet1),
		uint8(octet2),
		uint8(octet3),
		uint8(octet4),
		uint8(octet5),
		uint8(octet6),
	)
	sourceAmsPort := getFromOptions(localLogger, options, "sourceAmsPort")
	if sourceAmsPort == "" {
		return Configuration{}, errors.New("Required parameter sourceAmsPort missing")
	}
	parsedUint, err := strconv.ParseUint(sourceAmsPort, 10, 16)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing sourceAmsPort")
	}
	configuration.SourceAmsPort = uint16(parsedUint)
	targetAmsNetId := getFromOptions(localLogger, options, "targetAmsNetId")
	if sourceAmsNetId == "" {
		return Configuration{}, errors.New("Required parameter targetAmsNetId missing")
	}
	split = strings.Split(targetAmsNetId, ".")
	octet1, err = strconv.ParseUint(split[0], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing TargetAmsNetId")
	}
	octet2, err = strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing TargetAmsNetId")
	}
	octet3, err = strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing TargetAmsNetId")
	}
	octet4, err = strconv.ParseUint(split[3], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing TargetAmsNetId")
	}
	octet5, err = strconv.ParseUint(split[4], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing TargetAmsNetId")
	}
	octet6, err = strconv.ParseUint(split[5], 10, 8)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing TargetAmsNetId")
	}
	configuration.TargetAmsNetId = readWriteModel.NewAmsNetId(
		uint8(octet1),
		uint8(octet2),
		uint8(octet3),
		uint8(octet4),
		uint8(octet5),
		uint8(octet6),
	)
	targetAmsPort := getFromOptions(localLogger, options, "targetAmsPort")
	if targetAmsPort == "" {
		return Configuration{}, errors.New("Required parameter targetAmsPort missing")
	}
	parsedUint, err = strconv.ParseUint(targetAmsPort, 10, 16)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error parsing targetAmsPort")
	}
	configuration.TargetAmsPort = uint16(parsedUint)

	return configuration, nil
}

func getFromOptions(localLogger zerolog.Logger, options map[string][]string, key string) string {
	if optionValues, ok := options[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLogger.Warn().Msgf("Options %s must be unique", key)
		}
		return optionValues[0]
	}
	return ""
}
