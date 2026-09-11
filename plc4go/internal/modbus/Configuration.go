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
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
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
	// maxCoilsPerRequest is the widest run of coils or discrete inputs the read optimizer merges
	// into one request.
	maxCoilsPerRequest uint16
	// maxRegistersPerRequest is the widest run of registers the read optimizer merges into one
	// request. It applies to all three register areas; the extended registers are cut down further
	// where FC 0x14's framing demands it (see maxPerRequestFor).
	maxRegistersPerRequest uint16
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
	// defaultRequestTimeout is plc4j's request-timeout-ms default of five seconds.
	defaultRequestTimeout = 5 * time.Second
	// defaultMaxCoilsPerRequest and defaultMaxRegistersPerRequest are plc4j's
	// max-coils-per-request and max-registers-per-request defaults, which are the ceilings the
	// modbus specification itself imposes - merging up to them asks for the largest request that
	// can be answered at all.
	defaultMaxCoilsPerRequest     = uint16(maxCoilQuantity)
	defaultMaxRegistersPerRequest = uint16(maxRegisterQuantity)
)

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		unitIdentifier:          defaultUnitIdentifier,
		defaultPayloadByteOrder: BigEndianOrder,
		pingAddress:             defaultPingAddress,
		requestTimeout:          defaultRequestTimeout,
		maxCoilsPerRequest:      defaultMaxCoilsPerRequest,
		maxRegistersPerRequest:  defaultMaxRegistersPerRequest,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	// Every option this driver reads goes through the reader, so the ones nothing read can be
	// reported rather than silently discarded. Deferred, so no return path can skip it.
	reader := spiOptions.NewOptionReader(localLog, connectionOptions)
	defer reader.ReportUnknown("modbus")

	configuration := DefaultConfiguration()

	// One name for one concept: "default-unit-identifier", the same as plc4j. This driver also
	// accepted "unit-identifier", which plc4j never declared - so one connection string set the
	// unit here and was ignored there. Worse, "unit-identifier" *is* the name UMAS uses, where it
	// means something subtly different: modbus has a per-tag override ({unit-id: 3}), so this is
	// a default, while UMAS has none, so its is absolute. Two spellings meaning two things is
	// exactly what this vocabulary exists to stop. Supplying the old name is now reported.
	unitIdentifierString := reader.Get("default-unit-identifier")
	if unitIdentifierString != "" {
		parsedUint, err := strconv.ParseUint(unitIdentifierString, 10, 8)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing unit-identifier %s", unitIdentifierString)
		}
		configuration.unitIdentifier = uint8(parsedUint)
	}

	if byteOrderString := reader.Get("default-payload-byte-order"); byteOrderString != "" {
		byteOrder, ok := ByteOrderByName(byteOrderString)
		if !ok {
			return Configuration{}, errors.Errorf("Unknown default-payload-byte-order %s", byteOrderString)
		}
		configuration.defaultPayloadByteOrder = byteOrder
	}

	if pingAddress := reader.Get("ping-address"); pingAddress != "" {
		if _, err := NewTagHandler().ParseTag(pingAddress); err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing ping-address %s", pingAddress)
		}
		configuration.pingAddress = pingAddress
	}

	// plc4j states the request timeout in milliseconds.
	if requestTimeoutString := reader.Get("request-timeout-ms"); requestTimeoutString != "" {
		parsedUint, err := strconv.ParseUint(requestTimeoutString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing request-timeout-ms %s", requestTimeoutString)
		}
		if parsedUint == 0 {
			return Configuration{}, errors.Errorf("request-timeout-ms must be greater than zero. Was %s", requestTimeoutString)
		}
		configuration.requestTimeout = time.Duration(parsedUint) * time.Millisecond
	}

	// The two ceilings the read optimizer merges within. plc4j declares them on all three modbus
	// configurations - ModbusTcpConfiguration and its RTU and ASCII twins - and the website
	// documents them (modbus-tcp.adoc and friends), so a connection string that sets them is one
	// users have been told to write. This driver never read them, which made merging silently
	// ignore a device that can't answer a full-width request.
	maxCoilsPerRequest, err := parseMaxPerRequest(localLog, reader, "max-coils-per-request",
		configuration.maxCoilsPerRequest, maxCoilQuantity)
	if err != nil {
		return Configuration{}, err
	}
	configuration.maxCoilsPerRequest = maxCoilsPerRequest

	maxRegistersPerRequest, err := parseMaxPerRequest(localLog, reader, "max-registers-per-request",
		configuration.maxRegistersPerRequest, maxRegisterQuantity)
	if err != nil {
		return Configuration{}, err
	}
	configuration.maxRegistersPerRequest = maxRegistersPerRequest

	return configuration, nil
}

// parseMaxPerRequest reads one of the two per-request ceilings, falling back to the given default
// when the connection string doesn't name it.
//
// Both bounds are refused rather than clamped. Zero would ask the device for nothing, which it
// answers happily, so a read would report success and no data; anything above what the
// specification allows produces a response that doesn't fit into a modbus PDU, which no device can
// send. Either way a request built from such a number cannot work, and saying so beats sending it.
func parseMaxPerRequest(localLog zerolog.Logger, reader *spiOptions.OptionReader, name string, defaultValue uint16, ceiling uint16) (uint16, error) {
	value := reader.Get(name)
	if value == "" {
		return defaultValue, nil
	}
	parsedUint, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, errors.Wrapf(err, "Error parsing %s %s", name, value)
	}
	// Clamp rather than reject. plc4j declares these as plain ints with no validation at all
	// (ModbusTcpConfiguration and its RTU/ASCII siblings), so a connection string that opens a
	// connection there must open one here too -- refusing it would make the option LESS portable
	// than it was when plc4go ignored it entirely, which is the opposite of the point.
	//
	// Zero would mean "no tag may ever be read" and the ceiling is what one request can physically
	// carry, so both ends are pinned to something usable and the reason is logged.
	if parsedUint == 0 {
		localLog.Warn().Str("option", name).Str("value", value).Uint16("using", defaultValue).
			Msg("a per-request ceiling of zero would read nothing, so the default is used")
		return defaultValue, nil
	}
	if parsedUint > uint64(ceiling) {
		localLog.Warn().Str("option", name).Str("value", value).Uint16("using", ceiling).
			Msg("a per-request ceiling beyond what one modbus request can carry is clamped")
		return ceiling, nil
	}
	return uint16(parsedUint), nil
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
