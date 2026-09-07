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

package umas

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
)

// Configuration is what a umas connection string can say about the connection. Ported from plc4j's
// UmasConfiguration, minus its browser-generate-array-nodes flag: that option has a getter and a
// setter in plc4j and nothing reads it, so accepting it here would only promise a behaviour this
// driver doesn't have.
type Configuration struct {
	// unitIdentifier is the Modbus unit/slave address every ADU carries. UMAS PLCs answer on 0.
	unitIdentifier uint8
	// requestTimeout bounds every single request/response exchange.
	requestTimeout time.Duration
	// maxFrameSize is the frame size the driver assumes until the PLC reports its own limit in the
	// InitComms response. It only matters for the echo request of the handshake, whose payload is
	// sized from it.
	maxFrameSize uint16
}

const (
	// defaultUnitIdentifier is plc4j's @IntDefaultValue(0) for the unit-identifier option.
	defaultUnitIdentifier = uint8(0)
	// defaultRequestTimeout is plc4j's @IntDefaultValue(4000) for the request-timeout-ms option.
	defaultRequestTimeout = 4000 * time.Millisecond
	// defaultMaxFrameSize is plc4j's @IntDefaultValue(65535) for the max-frame-size option.
	defaultMaxFrameSize = uint16(65535)
	// minMaxFrameSize is the smallest frame size the echo request of the handshake can be built
	// from: it sends maxFrameSize - 3 bytes of payload.
	minMaxFrameSize = uint16(4)
)

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		unitIdentifier: defaultUnitIdentifier,
		requestTimeout: defaultRequestTimeout,
		maxFrameSize:   defaultMaxFrameSize,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string. The timeout is
// spelled in milliseconds, the way plc4j's @IntDefaultValue(4000) does.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	// Every option this driver reads goes through the reader, so the ones nothing read can be
	// reported rather than silently discarded. Deferred, so no return path can skip it.
	reader := spiOptions.NewOptionReader(localLog, connectionOptions)
	defer reader.ReportUnknown("umas")

	configuration := DefaultConfiguration()

	if unitIdentifierString := reader.Get("unit-identifier"); unitIdentifierString != "" {
		parsedInt, err := strconv.ParseUint(unitIdentifierString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing unit-identifier %s (has to fit into a single byte)", unitIdentifierString)
		}
		configuration.unitIdentifier = uint8(parsedInt)
	}

	if requestTimeoutString := reader.Get("request-timeout-ms"); requestTimeoutString != "" {
		parsedInt, err := strconv.ParseUint(requestTimeoutString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing request-timeout-ms %s", requestTimeoutString)
		}
		if parsedInt == 0 {
			return Configuration{}, errors.New("request-timeout-ms must be greater than zero")
		}
		configuration.requestTimeout = time.Duration(parsedInt) * time.Millisecond
	}

	if maxFrameSizeString := reader.Get("max-frame-size"); maxFrameSizeString != "" {
		parsedInt, err := strconv.ParseUint(maxFrameSizeString, 10, 16)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing max-frame-size %s (has to fit into two bytes)", maxFrameSizeString)
		}
		if parsedInt < uint64(minMaxFrameSize) {
			return Configuration{}, errors.Errorf("max-frame-size must be at least %d", minMaxFrameSize)
		}
		configuration.maxFrameSize = uint16(parsedInt)
	}

	return configuration, nil
}

func (c Configuration) String() string {
	return fmt.Sprintf("umas.Configuration{unitIdentifier: %d, requestTimeout: %s, maxFrameSize: %d}",
		c.unitIdentifier, c.requestTimeout, c.maxFrameSize)
}
