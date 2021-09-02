/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Configuration struct {
	backplane int8
	slot      int8
}

func ParseFromOptions(options map[string][]string) (Configuration, error) {
	configuration := Configuration{
		backplane: 0,
		slot:      0,
	}
	if localRackString := getFromOptions(options, "backplane"); localRackString != "" {
		parsedBackplane, err := strconv.ParseInt(localRackString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing backplane")
		}
		configuration.backplane = int8(parsedBackplane)
	}
	if localSlotString := getFromOptions(options, "slot"); localSlotString != "" {
		parsedSlot, err := strconv.ParseInt(localSlotString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing slot")
		}
		configuration.slot = int8(parsedSlot)
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
