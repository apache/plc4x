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
	"math"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

type Connection struct {
	_default.DefaultConnection

	configuration      Configuration
	messageCodec       spi.MessageCodec
	options            map[string][]string
	requestInterceptor interceptors.RequestInterceptor

	// transactionIdentifier numbers the requests Ping sends; reads and writes have counters of
	// their own in Reader and Writer.
	transactionIdentifier atomic.Int32

	connectionId string
	tracer       tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(configuration Configuration, messageCodec spi.MessageCodec, connectionOptions map[string][]string, tagHandler spi.PlcTagHandler, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		configuration: configuration,
		messageCodec:  messageCodec,
		options:       connectionOptions,
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
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
			_default.WithPlcValueHandler(NewValueHandler(_options...)),
		)...,
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

// nextTransactionIdentifier hands out the identifier of the next request this connection sends
// itself. The field on the wire is 16 bits wide, and zero is left out so that it never collides
// with an uninitialized one.
func (c *Connection) nextTransactionIdentifier() uint16 {
	next := c.transactionIdentifier.Add(1)
	if next > math.MaxUint16 {
		c.transactionIdentifier.Store(1)
		next = 1
	}
	return uint16(next)
}

// Ping checks whether the device is still there by reading the configured ping-address, the way
// plc4j's ModbusTcpConnection.onPing does. The diagnostic function code this used to send (FC 0x08)
// is optional and a good many devices answer it with an exception or not at all, while a read of a
// register every device has is a request the device is built to answer.
func (c *Connection) Ping(ctx context.Context) error {
	if c.DefaultConnection.IsInvalidated() {
		return errors.New("connection has been invalidated")
	}
	c.log.Trace().Msg("Pinging")
	tag, err := c.GetPlcTagHandler().ParseTag(c.configuration.pingAddress)
	if err != nil {
		return errors.Wrapf(err, "error parsing ping address '%s'", c.configuration.pingAddress)
	}
	pingTag, err := castToModbusTagFromPlcTag(tag)
	if err != nil {
		return errors.Wrapf(err, "error parsing ping address '%s'", c.configuration.pingAddress)
	}
	pdu, err := readRequestPdu(pingTag)
	if err != nil {
		return errors.Wrapf(err, "can't ping by reading '%s'", c.configuration.pingAddress)
	}
	transactionIdentifier := c.nextTransactionIdentifier()
	unitIdentifier := pingTag.resolveUnitId(c.configuration.unitIdentifier)
	adus := c.configuration.adus()
	pingRequest := adus.buildRequest(transactionIdentifier, unitIdentifier, pdu)

	requestCtx, cancelRequest := withRequestTimeout(ctx, c.configuration.requestTimeout)
	defer cancelRequest()

	errChan := make(chan error, 1)
	successChan := make(chan struct{}, 1)
	if err := c.messageCodec.SendRequest(requestCtx, "ping", pingRequest, func(message spi.Message) bool {
		return adus.acceptsResponse(pingRequest, message)
	}, func(message spi.Message) error {
		c.log.Trace().Msg("Received Message")
		if message == nil {
			c.log.Trace().Msg("got no response")
			select {
			case errChan <- errors.New("no response"):
			default:
				c.log.Warn().Msg("failed to send error signal")
			}
			return nil
		}
		// A modbus exception still means the device answered, so it is reachable - which is all a
		// ping asks. plc4j answers a ping the same way; the ping address is a guess that a given
		// device need not have, and a device that rejects it is still very much alive.
		if responsePdu, err := adus.extractPdu(message); err == nil {
			if errorPdu, isError := responsePdu.(readWriteModel.ModbusPDUError); isError {
				c.log.Warn().
					Stringer("exceptionCode", errorPdu.GetExceptionCode()).
					Str("pingAddress", c.configuration.pingAddress).
					Msg("The device rejected the ping address, but answering at all means it is reachable")
			}
		}
		c.log.Trace().Msg("got valid response")
		select {
		case successChan <- struct{}{}:
		default:
			c.log.Warn().Msg("failed to send success signal")
		}
		return nil
	}, func(err error) error {
		c.log.Trace().Msg("Received Error")
		select {
		case errChan <- errors.Wrap(err, "got error processing request"):
		default:
			c.log.Warn().Msg("failed to send error signal")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "error sending ping request")
	}
	select {
	case err := <-errChan:
		return errors.Wrap(err, "got error while waiting for response")
	case <-successChan:
		return nil
	case <-requestCtx.Done():
		return requestCtx.Err()
	}
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
		// Stated explicitly rather than left to the zero value: the modbus driver implements
		// neither subscribing nor browsing, so both builders fall through to
		// _default.DefaultConnection and panic. Flip these the moment that changes.
		ProvidesSubscribing: false,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilderWithInterceptor(
		c.GetPlcTagHandler(),
		NewReader(
			c.configuration,
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
			c.configuration,
			c.messageCodec,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
		c.requestInterceptor,
	)
}

func (c *Connection) String() string {
	return fmt.Sprintf("modbus.Connection{flavor: %s, unitIdentifier: %d, defaultPayloadByteOrder: %s, pingAddress: %s, requestTimeout: %s}",
		c.configuration.flavor, c.configuration.unitIdentifier, c.configuration.defaultPayloadByteOrder, c.configuration.pingAddress, c.configuration.requestTimeout)
}
