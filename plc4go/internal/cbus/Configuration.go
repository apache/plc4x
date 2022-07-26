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
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	srchk bool
}

func ParseFromOptions(options map[string][]string) (Configuration, error) {
	configuration := Configuration{
		srchk: true,
	}
	if srchk := getFromOptions(options, "srchk"); srchk != "" {
		parseBool, err := strconv.ParseBool(srchk)
		if err != nil {
			return Configuration{}, errors.Wrap(err, "Error parsing srchk")
		}
		configuration.srchk = parseBool
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
