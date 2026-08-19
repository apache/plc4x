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

package abeth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Configuration is what an ab-eth connection string can say about the connection. Ported from
// plc4j's AbEthConfiguration, which knows exactly these two options.
type Configuration struct {
	// station is the DF1 destination address every request is addressed to. It ends up in the
	// destinationAddress field of the DF1 request message.
	station uint8
	// requestTimeout bounds the connect handshake and every single read.
	requestTimeout time.Duration
}

const (
	// defaultStation is plc4j's @IntDefaultValue(0) for the station option.
	defaultStation = uint8(0)
	// defaultRequestTimeout is plc4j's @IntDefaultValue(10_000) for the request-timeout option.
	defaultRequestTimeout = 10 * time.Second
)

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		station:        defaultStation,
		requestTimeout: defaultRequestTimeout,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string. The timeout is
// spelled in milliseconds, the way plc4j's @IntDefaultValue(10_000) does.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	configuration := DefaultConfiguration()

	if stationString := getFromOptions(localLog, connectionOptions, "station"); stationString != "" {
		// The station is the DF1 destination address, which is a single byte on the wire.
		parsedInt, err := strconv.ParseUint(stationString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing station %s (has to fit into a single byte)", stationString)
		}
		configuration.station = uint8(parsedInt)
	}

	if requestTimeoutString := getFromOptions(localLog, connectionOptions, "request-timeout"); requestTimeoutString != "" {
		parsedInt, err := strconv.ParseUint(requestTimeoutString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing request-timeout %s", requestTimeoutString)
		}
		if parsedInt == 0 {
			return Configuration{}, errors.New("request-timeout must be greater than zero")
		}
		configuration.requestTimeout = time.Duration(parsedInt) * time.Millisecond
	}

	return configuration, nil
}

func (c Configuration) String() string {
	return fmt.Sprintf("abeth.Configuration{station: %d, requestTimeout: %s}", c.station, c.requestTimeout)
}

// getFromOptions plucks a single-valued option out of the parsed connection string.
func getFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string, key string) string {
	if optionValues, ok := connectionOptions[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLog.Warn().Str("key", key).Msg("Option must be unique")
		}
		return optionValues[0]
	}
	return ""
}
