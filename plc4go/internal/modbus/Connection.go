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
	"fmt"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

type Connection struct {
	_default.DefaultConnection
	unitIdentifier     uint8
	messageCodec       spi.MessageCodec
	options            map[string][]string
	requestInterceptor interceptors.RequestInterceptor

	connectionId string
	tracer       tracer.Tracer

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewConnection(unitIdentifier uint8, messageCodec spi.MessageCodec, connectionOptions map[string][]string, tagHandler spi.PlcTagHandler, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		unitIdentifier: unitIdentifier,
		messageCodec:   messageCodec,
		options:        connectionOptions,
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(
			spiModel.NewDefaultPlcReadRequest,
			spiModel.NewDefaultPlcWriteRequest,
			spiModel.NewDefaultPlcReadResponse,
			spiModel.NewDefaultPlcWriteResponse,
			_options...,
		),
		log:      customLogger,
		_options: _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithDefaultTtl(5*time.Second),
		_default.WithPlcTagHandler(tagHandler),
		_default.WithPlcValueHandler(NewValueHandler(_options...)),
	)
	return connection
}

func (c *Connection) GetConnectionId() string {
	return c.connectionId
}

func (c *Connection) IsTraceEnabled() bool {
	return c.tracer != nil
}

func (c *Connection) GetTracer() tracer.Tracer {
	return c.tracer
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
}

func (c *Connection) Ping() <-chan plc4go.PlcConnectionPingResult {
	// TODO: use proper context
	ctx := context.TODO()
	c.log.Trace().Msg("Pinging")
	result := make(chan plc4go.PlcConnectionPingResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- _default.NewDefaultPlcConnectionPingResult(errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		diagnosticRequestPdu := readWriteModel.NewModbusPDUDiagnosticRequest(0, 0x42)
		pingRequest := readWriteModel.NewModbusTcpADU(1, c.unitIdentifier, diagnosticRequestPdu, false)
		if err := c.messageCodec.SendRequest(ctx, pingRequest,
			func(message spi.Message) bool {
				responseAdu, ok := message.(readWriteModel.ModbusTcpADU)
				if !ok {
					return false
				}
				return responseAdu.GetTransactionIdentifier() == 1 && responseAdu.GetUnitIdentifier() == c.unitIdentifier
			},
			func(message spi.Message) error {
				c.log.Trace().Msg("Received Message")
				if message != nil {
					// If we got a valid response (even if it will probably contain an error, we know the remote is available)
					c.log.Trace().Msg("got valid response")
					result <- _default.NewDefaultPlcConnectionPingResult(nil)
				} else {
					c.log.Trace().Msg("got no response")
					result <- _default.NewDefaultPlcConnectionPingResult(errors.New("no response"))
				}
				return nil
			},
			func(err error) error {
				c.log.Trace().Msg("Received Error")
				result <- _default.NewDefaultPlcConnectionPingResult(errors.Wrap(err, "got error processing request"))
				return nil
			},
			time.Second*1,
		); err != nil {
			result <- _default.NewDefaultPlcConnectionPingResult(err)
		}
	}()
	return result
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilderWithInterceptor(
		c.GetPlcTagHandler(),
		NewReader(
			c.unitIdentifier,
			c.messageCodec,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
		c.requestInterceptor,
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilderWithInterceptor(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewWriter(
			c.unitIdentifier,
			c.messageCodec,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
		c.requestInterceptor,
	)
}

func (c *Connection) String() string {
	return fmt.Sprintf("modbus.Connection{unitIdentifier: %d}", c.unitIdentifier)
}
