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

package firmata

import (
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Configuration is what a firmata connection string can say about the connection as a whole.
// Ported from plc4j's FirmataConfiguration, which knows exactly one option.
type Configuration struct {
	// requestTimeout bounds how long the connect handshake waits for the board to answer the
	// system-reset with its report-firmware message. Firmata has no other request/response
	// interaction - writes are unacknowledged and subscriptions are push - so this is the only
	// place a timeout applies.
	requestTimeout time.Duration
}

const (
	// defaultRequestTimeout is plc4j's FirmataConfiguration request-timeout default of ten
	// seconds. A board that has just been reset needs a moment before it reports its firmware.
	defaultRequestTimeout = 10 * time.Second
)

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		requestTimeout: defaultRequestTimeout,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string. The timeout is
// spelled in milliseconds, the way plc4j's @IntDefaultValue(10_000) does.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	configuration := DefaultConfiguration()

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
