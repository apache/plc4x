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

package transports

import (
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options/converter"
	spiTransports "github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/serial"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

// TransportError exposes the SPI transport error wrapper on the public API so callers can perform
// errors.As checks and inspect transport severity without importing SPI internals.
type TransportError = spiTransports.TransportError

// TransportErrorKind mirrors the SPI severity enumeration so API users can make decisions based on
// transport error classification.
type TransportErrorKind = spiTransports.TransportErrorKind

const (
	// TransportErrorUnknown represents an error the transport could not classify; treat as fatal when unsure.
	TransportErrorUnknown TransportErrorKind = spiTransports.TransportErrorUnknown

	// TransportErrorTransient signals a short-lived transport hiccup that usually succeeds if re-tried immediately.
	TransportErrorTransient TransportErrorKind = spiTransports.TransportErrorTransient

	// TransportErrorRetryable indicates the caller should retry the operation after resetting or reconnecting the transport.
	TransportErrorRetryable TransportErrorKind = spiTransports.TransportErrorRetryable

	// TransportErrorFatal marks the transport as unusable for further work; callers must tear down and rebuild the connection.
	TransportErrorFatal TransportErrorKind = spiTransports.TransportErrorFatal
)

// RegisterTcpTransport registers the TCP transport implementation with the supplied driver manager using the provided options.
func RegisterTcpTransport(driverManager plc4go.PlcDriverManager, _options ...config.WithOption) {
	driverManager.(spi.TransportAware).RegisterTransport(tcp.NewTransport(converter.WithOptionToInternal(_options...)...))
}

// RegisterUdpTransport registers the UDP transport implementation with the supplied driver manager using the provided options.
func RegisterUdpTransport(driverManager plc4go.PlcDriverManager, _options ...config.WithOption) {
	driverManager.(spi.TransportAware).RegisterTransport(udp.NewTransport(converter.WithOptionToInternal(_options...)...))
}

// RegisterSerialTransport registers the serial transport implementation with the supplied driver manager using the provided options.
func RegisterSerialTransport(driverManager plc4go.PlcDriverManager, _options ...config.WithOption) {
	driverManager.(spi.TransportAware).RegisterTransport(serial.NewTransport(converter.WithOptionToInternal(_options...)...))
}
