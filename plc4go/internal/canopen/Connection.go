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

package canopen

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// Connection is a CANopen connection, the port of plc4j's CANOpenConnection.
//
// CANopen has no session: connecting opens the bus, and from then on the connection is one more
// node on it. It reads and writes object dictionary entries with SDO transfers (as an SDO client),
// sends PDOs and NMT commands, optionally produces a heartbeat, and turns the PDO, NMT and heartbeat
// traffic it sees into subscription events.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	messageCodec  *MessageCodec
	subscriber    *Subscriber

	// sdoMutex serialises SDO transfers. A node serves one SDO transfer at a time per SDO channel
	// and its answers carry nothing that ties them to a request, so two transfers in flight on the
	// default channel would take each other's answers.
	sdoMutex sync.Mutex

	lifecycleMutex  sync.Mutex
	heartbeatCancel context.CancelFunc

	connectionId string
	tracer       tracer.Tracer

	wg sync.WaitGroup

	log      zerolog.Logger
	_options []options.WithOption
}

var _ spi.TransportInstanceExposer = (*Connection)(nil)

func NewConnection(messageCodec *MessageCodec, configuration Configuration, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	_options = slices.Clip(_options)
	connection := &Connection{
		configuration: configuration,
		messageCodec:  messageCodec,
		log:           customLogger,
		_options:      _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok && len(traceEnabledOption) == 1 {
		connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		append(_options,
			_default.WithPlcTagHandler(NewTagHandler()),
			_default.WithPlcValueHandler(spiValues.NewDefaultValueHandler(_options...)),
		)...,
	)
	connection.subscriber = NewSubscriber(connection, _options...)
	messageCodec.handler = connection.handleFrame
	messageCodec.onFatal = func(err error) {
		connection.log.Warn().Err(err).Msg("CANopen transport failed, invalidating the connection")
		connection.Invalidate()
	}
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

func (c *Connection) Connect(ctx context.Context) error {
	if err := c.DefaultConnection.Connect(ctx); err != nil {
		return errors.Wrap(err, "error opening the CAN bus")
	}
	if c.configuration.heartbeat {
		if err := c.messageCodec.sendFrame(ctx, c.heartbeatFrame(readWriteModel.NMTState_BOOTED_UP)); err != nil {
			_ = c.Close()
			return errors.Wrap(err, "error sending the boot-up message")
		}
		c.startHeartbeat()
	}
	c.log.Info().Stringer("configuration", c.configuration).Msg("CANopen connection established")
	return nil
}

func (c *Connection) Close() error {
	c.lifecycleMutex.Lock()
	heartbeatCancel := c.heartbeatCancel
	c.heartbeatCancel = nil
	c.lifecycleMutex.Unlock()
	if heartbeatCancel != nil {
		heartbeatCancel()
	}
	c.subscriber.Close()
	err := c.DefaultConnection.Close()
	c.wg.Wait()
	return err
}

func (c *Connection) heartbeatFrame(state readWriteModel.NMTState) readWriteModel.CANOpenFrame {
	return readWriteModel.NewCANOpenFrame(c.configuration.nodeID, readWriteModel.CANOpenService_HEARTBEAT, readWriteModel.NewCANOpenHeartbeatPayload(state))
}

// startHeartbeat produces OPERATIONAL heartbeats, the way plc4j does after its boot-up message.
func (c *Connection) startHeartbeat() {
	ctx, cancel := context.WithCancel(context.Background())
	c.lifecycleMutex.Lock()
	c.heartbeatCancel = cancel
	c.lifecycleMutex.Unlock()
	c.wg.Go(func() {
		ticker := time.NewTicker(c.configuration.heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sendCtx, cancelSend := context.WithTimeout(ctx, c.configuration.requestTimeout)
				if err := c.messageCodec.sendFrame(sendCtx, c.heartbeatFrame(readWriteModel.NMTState_OPERATIONAL)); err != nil {
					c.log.Warn().Err(err).Msg("error sending the CANopen heartbeat")
				}
				cancelSend()
			}
		}
	})
}

// handleFrame gets every frame no SDO transfer was waiting for.
func (c *Connection) handleFrame(frame readWriteModel.CANOpenFrame) {
	if c.configuration.nodeID != 0 && frame.GetNodeId() == c.configuration.nodeID && frame.GetService() != readWriteModel.CANOpenService_NMT {
		// Our own traffic echoed back by a transport that loops frames to the sender.
		return
	}
	c.subscriber.publish(frame)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Read: SDO upload
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	result := make(chan apiModel.PlcReadRequestResult, 1)
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil,
					errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		responseCodes := map[string]apiModel.PlcResponseCode{}
		plcValues := map[string]apiValues.PlcValue{}
		// plc4j refuses requests with more than one tag; here each tag is its own transfer.
		for _, tagName := range readRequest.GetTagNames() {
			tag, ok := castToTag(readRequest.GetTag(tagName))
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			if tag.Kind != TagKindSDO {
				// PDO, NMT and heartbeat traffic is pushed; it is subscribed to, not read.
				responseCodes[tagName] = apiModel.PlcResponseCode_UNSUPPORTED
				continue
			}
			code, value := c.readSDO(ctx, tagName, tag)
			responseCodes[tagName] = code
			if value != nil {
				plcValues[tagName] = value
			}
		}
		utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest,
			spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil))
	})
	return result
}

func (c *Connection) readSDO(ctx context.Context, tagName string, tag Tag) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	c.sdoMutex.Lock()
	data, err := newSDOTransfer(c.messageCodec, tag, c.configuration.requestTimeout).upload(ctx)
	c.sdoMutex.Unlock()
	if err != nil {
		c.log.Warn().Err(err).Str("tagName", tagName).Msg("SDO upload failed")
		var abort *AbortError
		if errors.As(err, &abort) {
			// The abort code is what says why, so it goes back as the value; the protocol's driver
			// test suite expects it there.
			return responseCodeFor(ctx, err), spiValues.NewPlcLINT(int64(abort.Code))
		}
		return responseCodeFor(ctx, err), nil
	}
	value, err := decodeValue(ctx, tag.DataType, data)
	if err != nil {
		c.log.Warn().Err(err).Str("tagName", tagName).Hex("data", data).Msg("SDO upload data doesn't decode as the tag's type")
		return apiModel.PlcResponseCode_INVALID_DATA, nil
	}
	return apiModel.PlcResponseCode_OK, value
}

// responseCodeFor maps a failed transfer onto a response code.
func responseCodeFor(ctx context.Context, err error) apiModel.PlcResponseCode {
	var abort *AbortError
	switch {
	case errors.As(err, &abort):
		switch abort.Code {
		case 0x06020000, 0x06090011:
			return apiModel.PlcResponseCode_NOT_FOUND
		case 0x06010000, 0x06010001, 0x06010002:
			return apiModel.PlcResponseCode_ACCESS_DENIED
		case 0x06070010, 0x06070012, 0x06070013, 0x06090030, 0x06090031, 0x06090032:
			return apiModel.PlcResponseCode_INVALID_DATA
		}
		return apiModel.PlcResponseCode_REMOTE_ERROR
	case errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil:
		return apiModel.PlcResponseCode_REQUEST_TIMEOUT
	}
	return apiModel.PlcResponseCode_INTERNAL_ERROR
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Write: SDO download, PDO, NMT command
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil,
					errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		responseCodes := map[string]apiModel.PlcResponseCode{}
		for _, tagName := range writeRequest.GetTagNames() {
			tag, ok := castToTag(writeRequest.GetTag(tagName))
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			responseCodes[tagName] = c.writeTag(ctx, tagName, tag, writeRequest.GetValue(tagName))
		}
		utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest,
			spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil))
	})
	return result
}

func (c *Connection) writeTag(ctx context.Context, tagName string, tag Tag, value apiValues.PlcValue) apiModel.PlcResponseCode {
	if value == nil {
		return apiModel.PlcResponseCode_INVALID_DATA
	}
	switch tag.Kind {
	case TagKindSDO:
		data, err := encodeValue(tag.DataType, value)
		if err != nil {
			c.log.Warn().Err(err).Str("tagName", tagName).Msg("value doesn't encode as the tag's type")
			return apiModel.PlcResponseCode_INVALID_DATA
		}
		c.sdoMutex.Lock()
		err = newSDOTransfer(c.messageCodec, tag, c.configuration.requestTimeout).download(ctx, data)
		c.sdoMutex.Unlock()
		if err != nil {
			c.log.Warn().Err(err).Str("tagName", tagName).Msg("SDO download failed")
			return responseCodeFor(ctx, err)
		}
		return apiModel.PlcResponseCode_OK
	case TagKindPDO:
		data, err := encodeValue(tag.DataType, value)
		if err != nil || len(data) > 8 {
			c.log.Warn().Err(err).Str("tagName", tagName).Msg("value doesn't encode as a PDO of the tag's type")
			return apiModel.PlcResponseCode_INVALID_DATA
		}
		frame := readWriteModel.NewCANOpenFrame(tag.Node, tag.Service, readWriteModel.NewCANOpenPDOPayload(readWriteModel.NewCANOpenPDO(data)))
		return c.send(ctx, tagName, frame)
	case TagKindNMT:
		// plc4j has no NMT write; the protocol's NMT driver test suite expects one. The value is the
		// NMT command specifier: 1 start, 2 stop, 128 pre-operational, 129 reset node, 130 reset
		// communication. NMT always goes out as COB-ID 0, the target is in the payload, 0 for all.
		var request readWriteModel.NMTStateRequest
		ok := value.IsUint8()
		if ok {
			request, ok = readWriteModel.NMTStateRequestByValue(value.GetUint8())
		}
		if !ok {
			c.log.Warn().Str("tagName", tagName).Stringer("value", value).Msg("not an NMT command (1, 2, 128, 129 or 130)")
			return apiModel.PlcResponseCode_INVALID_DATA
		}
		frame := readWriteModel.NewCANOpenFrame(0, readWriteModel.CANOpenService_NMT, readWriteModel.NewCANOpenNetworkPayload(request, tag.Node))
		return c.send(ctx, tagName, frame)
	}
	return apiModel.PlcResponseCode_UNSUPPORTED
}

func (c *Connection) send(ctx context.Context, tagName string, frame readWriteModel.CANOpenFrame) apiModel.PlcResponseCode {
	sendCtx, cancel := context.WithTimeout(ctx, c.configuration.requestTimeout)
	defer cancel()
	if err := c.messageCodec.sendFrame(sendCtx, frame); err != nil {
		c.log.Warn().Err(err).Str("tagName", tagName).Msg("error sending the frame")
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
	return apiModel.PlcResponseCode_OK
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Metadata and builders
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(c.GetPlcTagHandler(), readerFunc(c.read))
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), writerFunc(c.write))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), c.subscriber)
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) String() string {
	return fmt.Sprintf("canopen.Connection{%s}", c.configuration)
}

type readerFunc func(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult

func (f readerFunc) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	return f(ctx, readRequest)
}

type writerFunc func(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult

func (f writerFunc) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	return f(ctx, writeRequest)
}
