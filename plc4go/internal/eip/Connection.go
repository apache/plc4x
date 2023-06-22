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

package eip

import (
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/rs/zerolog"
	"runtime/debug"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

const (
	DefaultSenderContext = "PLC4X   "
	EmptySessionHandle   = uint32(0)
	EmptyInterfaceHandle = uint32(0)
)

type Connection struct {
	_default.DefaultConnection
	messageCodec              spi.MessageCodec
	configuration             Configuration
	driverContext             DriverContext
	tm                        transactions.RequestTransactionManager
	sessionHandle             uint32
	senderContext             []uint8
	connectionId              uint32
	cipEncapsulationAvailable bool
	connectionSerialNumber    uint16
	connectionPathSize        uint8
	useMessageRouter          bool
	useConnectionManager      bool
	routingAddress            []readWriteModel.PathSegment
	tracer                    tracer.Tracer

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewConnection(
	messageCodec spi.MessageCodec,
	configuration Configuration,
	driverContext DriverContext,
	tagHandler spi.PlcTagHandler,
	tm transactions.RequestTransactionManager,
	connectionOptions map[string][]string,
	_options ...options.WithOption,
) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		messageCodec:  messageCodec,
		configuration: configuration,
		driverContext: driverContext,
		tm:            tm,
		log:           customLogger,
		_options:      _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			// TODO: Fix this.
			//			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
			_default.WithPlcValueHandler(NewValueHandler(_options...)),
		)...,
	)

	// TODO: connectionPathSize
	// TODO: routingAddress
	return connection
}

func (c *Connection) GetConnectionId() string {
	// TODO: Fix this
	return "" //c.connectionId
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

func (c *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	c.log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		err := c.messageCodec.ConnectWithContext(ctx)
		if err != nil {
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !c.driverContext.awaitSetupComplete {
			go c.setupConnection(ctx, ch)
			c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against c.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
			c.SetConnected(true)
			return
		}

		c.setupConnection(ctx, ch)
	}()
	return ch
}

func (c *Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	// TODO: use proper context
	ctx := context.TODO()
	result := make(chan plc4go.PlcConnectionCloseResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		c.log.Debug().Msg("Sending UnregisterSession EIP Packet")
		_ = c.messageCodec.SendRequest(
			ctx,
			readWriteModel.NewEipDisconnectRequest(c.sessionHandle, 0, []byte(DefaultSenderContext), 0), func(message spi.Message) bool {
				return true
			},
			func(message spi.Message) error {
				return nil
			},
			func(err error) error {
				return nil
			},
			c.GetTtl(),
		) //Unregister gets no response
		c.log.Debug().Msgf("Unregistred Session %d", c.sessionHandle)
	}()
	return result
}

func (c *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	if err := c.listServiceRequest(ctx, ch); err != nil {
		c.fireConnectionError(errors.Wrap(err, "error listing service request"), ch)
		return
	}

	if err := c.connectRegisterSession(ctx, ch); err != nil {
		c.fireConnectionError(errors.Wrap(err, "error connect register session"), ch)
		return
	}

	if err := c.listAllAttributes(ctx, ch); err != nil {
		c.fireConnectionError(errors.Wrap(err, "error list all attributes"), ch)
		return
	}

	if c.useConnectionManager {
		// TODO: Continue here ....
	} else {
		// Send an event that connection setup is complete.
		c.fireConnected(ch)
	}
}

func (c *Connection) listServiceRequest(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) error {
	c.log.Debug().Msg("Sending ListServices Request")
	listServicesResultChan := make(chan readWriteModel.ListServicesResponse, 1)
	listServicesResultErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(
		ctx,
		readWriteModel.NewListServicesRequest(
			EmptySessionHandle,
			uint32(readWriteModel.CIPStatus_Success),
			[]byte(DefaultSenderContext),
			uint32(0),
		),
		func(message spi.Message) bool {
			eipPacket := message.(readWriteModel.EipPacketExactly)
			if eipPacket == nil {
				return false
			}
			eipPacketListServicesResponse := eipPacket.(readWriteModel.ListServicesResponseExactly)
			return eipPacketListServicesResponse != nil
		},
		func(message spi.Message) error {
			listServicesResponse := message.(readWriteModel.ListServicesResponse)
			serviceResponse := listServicesResponse.GetTypeIds()[0].(readWriteModel.ServicesResponse)
			if serviceResponse.GetSupportsCIPEncapsulation() {
				c.log.Debug().Msg("Device is capable of CIP over EIP encapsulation")
			}
			c.cipEncapsulationAvailable = serviceResponse.GetSupportsCIPEncapsulation()
			listServicesResultChan <- listServicesResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				c.Close()
			}
			listServicesResultErrorChan <- errors.Wrap(err, "got error processing request")
			return nil
		},
		c.GetTtl()); err != nil {
		c.fireConnectionError(errors.Wrap(err, "Error during sending of EIP ListServices Request"), ch)
	}

	timeout := time.NewTimer(1 * time.Second)
	defer utils.CleanupTimer(timeout)
	select {
	case <-timeout.C:
		return errors.New("timeout")
	case err := <-listServicesResultErrorChan:
		return errors.Wrap(err, "Error receiving of ListServices response")
	case _ = <-listServicesResultChan:
		return nil
	}
}

func (c *Connection) connectRegisterSession(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) error {
	c.log.Debug().Msg("Sending EipConnectionRequest")
	connectionResponseChan := make(chan readWriteModel.EipConnectionResponse, 1)
	connectionResponseErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(
		ctx,
		readWriteModel.NewEipConnectionRequest(
			EmptySessionHandle,
			uint32(readWriteModel.CIPStatus_Success),
			[]byte(DefaultSenderContext),
			uint32(0),
		),
		func(message spi.Message) bool {
			eipPacket := message.(readWriteModel.EipPacketExactly)
			return eipPacket != nil
		},
		func(message spi.Message) error {
			eipPacket := message.(readWriteModel.EipPacket)
			connectionResponse := eipPacket.(readWriteModel.EipConnectionResponse)
			if connectionResponse != nil {
				if connectionResponse.GetStatus() == 0 {
					c.sessionHandle = connectionResponse.GetSessionHandle()
					c.senderContext = connectionResponse.GetSenderContext()
					c.log.Debug().Msgf("Got assigned with Session %d", c.sessionHandle)
					connectionResponseChan <- connectionResponse
				} else {
					c.log.Error().Msgf("Got unsuccessful status for connection request: %d", connectionResponse.GetStatus())
					connectionResponseErrorChan <- errors.New("got unsuccessful connection response")
				}
			} else {
				// TODO: This seems pretty hard-coded ... possibly find out if we can't simplify this.
				classSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6))
				instanceSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 1))
				exchange := readWriteModel.NewUnConnectedDataItem(
					readWriteModel.NewCipConnectionManagerRequest(classSegment, instanceSegment, 0, 10,
						14, 536870914, 33944, c.connectionSerialNumber,
						4919, 42, 3, 2101812,
						readWriteModel.NewNetworkConnectionParameters(4002, false, 2, 0, true),
						2113537,
						readWriteModel.NewNetworkConnectionParameters(4002, false, 2, 0, true),
						readWriteModel.NewTransportType(true, 2, 3),
						c.connectionPathSize, c.routingAddress, 1))
				typeIds := []readWriteModel.TypeId{readWriteModel.NewNullAddressItem(), exchange}
				eipWrapper := readWriteModel.NewCipRRData(c.sessionHandle, 0, typeIds,
					c.sessionHandle, uint32(readWriteModel.CIPStatus_Success), c.senderContext, 0)
				if err := c.messageCodec.SendRequest(
					ctx,
					eipWrapper,
					func(message spi.Message) bool {
						eipPacket := message.(readWriteModel.EipPacketExactly)
						if eipPacket == nil {
							return false
						}
						cipRRData := eipPacket.(readWriteModel.CipRRDataExactly)
						return cipRRData != nil
					},
					func(message spi.Message) error {
						cipRRData := message.(readWriteModel.CipRRData)
						if cipRRData.GetStatus() == 0 {
							unconnectedDataItem := cipRRData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
							connectionManagerResponse := unconnectedDataItem.GetService().(readWriteModel.CipConnectionManagerResponse)
							c.connectionId = connectionManagerResponse.GetOtConnectionId()
							c.log.Debug().Msgf("Got assigned with connection if %d", c.connectionId)
							connectionResponseChan <- connectionResponse
						} else {
							connectionResponseErrorChan <- fmt.Errorf("got status code while opening Connection manager: %d", cipRRData.GetStatus())
						}
						return nil
					},
					func(err error) error {
						// If this is a timeout, do a check if the connection requires a reconnection
						if _, isTimeout := err.(utils.TimeoutError); isTimeout {
							c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
							c.Close()
						}
						connectionResponseErrorChan <- errors.Wrap(err, "got error processing request")
						return nil
					},
					c.GetTtl(),
				); err != nil {
					c.fireConnectionError(errors.Wrap(err, "Error during sending of EIP ListServices Request"), ch)
				}
			}
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				c.Close()
			}
			connectionResponseErrorChan <- errors.Wrap(err, "got error processing request")
			return nil
		},
		c.GetTtl(),
	); err != nil {
		c.fireConnectionError(errors.Wrap(err, "Error during sending of EIP ListServices Request"), ch)
	}
	timeout := time.NewTimer(1 * time.Second)
	defer utils.CleanupTimer(timeout)
	select {
	case <-timeout.C:
		return errors.New("timeout")
	case err := <-connectionResponseErrorChan:
		return errors.Wrap(err, "Error receiving of ListServices response")
	case _ = <-connectionResponseChan:
		return nil
	}
}

func (c *Connection) listAllAttributes(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) error {
	c.log.Debug().Msg("Sending ListAllAttributes Request")
	listAllAttributesResponseChan := make(chan readWriteModel.GetAttributeAllResponse, 1)
	listAllAttributesErrorChan := make(chan error, 1)
	classSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(uint8(0), uint8(2)))
	instanceSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(uint8(0), uint8(1)))
	if err := c.messageCodec.SendRequest(
		ctx,
		readWriteModel.NewCipRRData(
			EmptyInterfaceHandle,
			0,
			[]readWriteModel.TypeId{
				readWriteModel.NewNullAddressItem(),
				readWriteModel.NewUnConnectedDataItem(
					readWriteModel.NewGetAttributeAllRequest(
						classSegment, instanceSegment, uint16(0))),
			},
			c.sessionHandle,
			uint32(readWriteModel.CIPStatus_Success),
			c.senderContext,
			0,
		),
		func(message spi.Message) bool {
			eipPacket := message.(readWriteModel.CipRRDataExactly)
			return eipPacket != nil
		},
		func(message spi.Message) error {
			cipRrData := message.(readWriteModel.CipRRData)
			if cipRrData.GetStatus() == uint32(readWriteModel.CIPStatus_Success) {
				dataItem := cipRrData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
				response := dataItem.GetService().(readWriteModel.GetAttributeAllResponse)
				if response.GetStatus() != uint8(readWriteModel.CIPStatus_Success) {
					// TODO: Return an error ...
				} else if response.GetAttributes() != nil {
					for _, classId := range response.GetAttributes().GetClassId() {
						if curCipClassId, ok := readWriteModel.CIPClassIDByValue(classId); ok {
							switch curCipClassId {
							case readWriteModel.CIPClassID_MessageRouter:
								c.useMessageRouter = true
							case readWriteModel.CIPClassID_ConnectionManager:
								c.useConnectionManager = true
							}
						}
					}
				}
				c.log.Debug().Msgf("Connection using message router %t, using connection manager %t", c.useMessageRouter, c.useConnectionManager)
				listAllAttributesResponseChan <- response
			}
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				c.Close()
			}
			c.fireConnectionError(errors.Wrap(err, "got error processing request"), ch)
			return nil
		},
		c.GetTtl(),
	); err != nil {
		c.fireConnectionError(errors.Wrap(err, "Error during sending of EIP ListServices Request"), ch)
	}

	timeout := time.NewTimer(1 * time.Second)
	defer utils.CleanupTimer(timeout)
	select {
	case <-timeout.C:
		return errors.New("timeout")
	case err := <-listAllAttributesErrorChan:
		return errors.Wrap(err, "Error receiving of ListServices response")
	case _ = <-listAllAttributesResponseChan:
		return nil
	}
}

func (c *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
	} else {
		c.log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
}

func (c *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
	} else {
		c.log.Info().Msg("Successfully connected")
	}
	c.SetConnected(true)
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			c.messageCodec,
			c.tm,
			c.configuration,
			&c.sessionHandle,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewWriter(
			c.messageCodec,
			c.tm,
			c.configuration,
			&c.sessionHandle,
			&c.senderContext,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) String() string {
	return fmt.Sprintf("eip.Connection")
}
