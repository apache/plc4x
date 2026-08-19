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

package modbus

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Configuration is what a modbus connection string can say about the connection as a whole. Ported
// from plc4j's ModbusTcpConfiguration, which the RTU and ASCII configurations repeat.
type Configuration struct {
	// unitIdentifier is the unit identifier (slave id) a tag that doesn't name one of its own is
	// addressed at. On RS485 several modbus devices share one line and are told apart by it.
	unitIdentifier uint8
	// defaultPayloadByteOrder is the byte order a tag that doesn't name one of its own is read and
	// written with.
	defaultPayloadByteOrder ByteOrder
	// pingAddress is the address Ping reads to find out whether the device is still there.
	pingAddress string
	// requestTimeout bounds how long a single request waits for its response.
	requestTimeout time.Duration
	// flavor is the framing this connection speaks. It is not something the connection string can
	// say - each driver speaks exactly one flavor and sets it on the configuration it hands to the
	// connection. The zero value is flavorTcp, so a configuration nobody told otherwise behaves
	// exactly as it did before the RTU flavor got a framing of its own.
	flavor modbusFlavor
}

// withFlavor returns a copy of the configuration that speaks the given flavor.
func (c Configuration) withFlavor(flavor modbusFlavor) Configuration {
	c.flavor = flavor
	return c
}

// adus is the factory that builds and reads the ADUs of the configured flavor.
func (c Configuration) adus() aduFactory {
	return c.flavor.adus()
}

const (
	defaultUnitIdentifier = uint8(1)
	// defaultPingAddress reads the first holding register, as plc4j's ModbusTcpConfiguration does.
	defaultPingAddress = "4x00001:BOOL"
	// defaultRequestTimeout is plc4j's request-timeout default of five seconds.
	defaultRequestTimeout = 5 * time.Second
)

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		unitIdentifier:          defaultUnitIdentifier,
		defaultPayloadByteOrder: BigEndianOrder,
		pingAddress:             defaultPingAddress,
		requestTimeout:          defaultRequestTimeout,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	configuration := DefaultConfiguration()

	// plc4j spells the option default-unit-identifier; the Go driver has always called it
	// unit-identifier, so both are accepted and the plc4j spelling wins when somebody sets both.
	unitIdentifierString := getFromOptions(localLog, connectionOptions, "unit-identifier")
	if defaultUnitIdentifierString := getFromOptions(localLog, connectionOptions, "default-unit-identifier"); defaultUnitIdentifierString != "" {
		unitIdentifierString = defaultUnitIdentifierString
	}
	if unitIdentifierString != "" {
		parsedUint, err := strconv.ParseUint(unitIdentifierString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing unit-identifier %s", unitIdentifierString)
		}
		configuration.unitIdentifier = uint8(parsedUint)
	}

	if byteOrderString := getFromOptions(localLog, connectionOptions, "default-payload-byte-order"); byteOrderString != "" {
		byteOrder, ok := ByteOrderByName(byteOrderString)
		if !ok {
			return Configuration{}, errors.Errorf("Unknown default-payload-byte-order %s", byteOrderString)
		}
		configuration.defaultPayloadByteOrder = byteOrder
	}

	if pingAddress := getFromOptions(localLog, connectionOptions, "ping-address"); pingAddress != "" {
		if _, err := NewTagHandler().ParseTag(pingAddress); err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing ping-address %s", pingAddress)
		}
		configuration.pingAddress = pingAddress
	}

	// plc4j states the request timeout in milliseconds.
	if requestTimeoutString := getFromOptions(localLog, connectionOptions, "request-timeout"); requestTimeoutString != "" {
		parsedUint, err := strconv.ParseUint(requestTimeoutString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing request-timeout %s", requestTimeoutString)
		}
		if parsedUint == 0 {
			return Configuration{}, errors.Errorf("request-timeout must be greater than zero. Was %s", requestTimeoutString)
		}
		configuration.requestTimeout = time.Duration(parsedUint) * time.Millisecond
	}

	return configuration, nil
}

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

// withRequestTimeout bounds a single request. The codec turns the deadline of the context it is
// handed into the lifetime of its expectation (spi/default.defaultCodec.expect), so a deadline is
// all it takes to time a request out; a caller who brought a deadline of their own keeps it.
// The returned cancel function has to be called once a result has been delivered - it releases the
// timer and tells the codec to drop the expectation if one is still registered.
func withRequestTimeout(ctx context.Context, requestTimeout time.Duration) (context.Context, context.CancelFunc) {
	if requestTimeout <= 0 {
		return context.WithCancel(ctx)
	}
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, requestTimeout)
}
