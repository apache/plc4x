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

package slmp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Configuration is what an slmp connection string can say about the connection. Ported from plc4j's
// SlmpConfiguration, which knows exactly these two options.
type Configuration struct {
	// monitoringTimer is the SLMP monitoring timer written into every 3E request frame. Zero means
	// "wait infinitely" as far as the device is concerned; the client-side bound is requestTimeout.
	// It is an unsigned 16-bit field on the wire, which is why it is typed as one here - plc4j
	// carries it as an int and has to range-check it by hand in SlmpConnection.validateConfiguration.
	monitoringTimer uint16
	// requestTimeout bounds every single read and write.
	requestTimeout time.Duration
}

const (
	// defaultMonitoringTimer is plc4j's @IntDefaultValue(0x0000) for the monitoring-timer option.
	defaultMonitoringTimer = uint16(0x0000)
	// defaultRequestTimeout is plc4j's @IntDefaultValue(5_000) for the request-timeout option.
	defaultRequestTimeout = 5 * time.Second
)

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		monitoringTimer: defaultMonitoringTimer,
		requestTimeout:  defaultRequestTimeout,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string. The timeout is
// spelled in milliseconds, the way plc4j's @IntDefaultValue(5_000) does.
//
// plc4j defers both range checks to SlmpConnection.validateConfiguration, because its configuration
// object is populated by field injection and its setters are bypassed. Here the parse is the only
// way in, so both checks live at the parse.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	configuration := DefaultConfiguration()

	if monitoringTimerString := getFromOptions(localLog, connectionOptions, "monitoring-timer"); monitoringTimerString != "" {
		parsedInt, err := strconv.ParseUint(monitoringTimerString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing monitoring-timer %s (it is an unsigned 16-bit field in the 3E frame, so it has to be in [0, 65535])", monitoringTimerString)
		}
		configuration.monitoringTimer = uint16(parsedInt)
	}

	if requestTimeoutString := getFromOptions(localLog, connectionOptions, "request-timeout"); requestTimeoutString != "" {
		parsedInt, err := strconv.ParseUint(requestTimeoutString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing request-timeout %s", requestTimeoutString)
		}
		if parsedInt == 0 {
			// plc4j rejects this at connect time with the same reasoning: a non-positive timeout
			// would time out every request immediately.
			return Configuration{}, errors.New("request-timeout must be greater than zero")
		}
		configuration.requestTimeout = time.Duration(parsedInt) * time.Millisecond
	}

	return configuration, nil
}

func (c Configuration) String() string {
	return fmt.Sprintf("slmp.Configuration{monitoringTimer: 0x%04X, requestTimeout: %s}", c.monitoringTimer, c.requestTimeout)
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
