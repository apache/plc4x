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
	"sync"
	"time"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
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
	sessionState              *SessionState
	cipEncapsulationAvailable bool
	tracer                    tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

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
		sessionState:  NewSessionState(customLogger, configuration),
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

func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting message codec")
	}

	// For testing purposes we can skip the waiting for a complete connection
	if !c.driverContext.awaitSetupComplete {
		// The caller's ctx is canceled as soon as GetConnection returns - the
		// background handshake must not inherit that cancellation (GH-954).
		setupCtx := context.WithoutCancel(ctx)
		c.wg.Go(func() {
			setupCtx, cancel := utils.WithNamedTimeout(setupCtx, "eip setup timeout", 10*time.Second)
			defer cancel()
			if err := c.setupConnection(setupCtx); err != nil {
				c.log.Error().Err(err).Msg("error during setup connection")
			}
		})
		c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		// Here we write directly and don't wait till the connection is "really" connected
		c.SetConnected(true)
		return nil
	}

	if err := c.setupConnection(ctx); err != nil {
		if disconnectErr := c.messageCodec.Disconnect(); disconnectErr != nil {
			c.log.Debug().Err(disconnectErr).Msg("error disconnecting after failed setup")
		}
		return errors.Wrap(err, "error during setup connection")
	}
	return nil
}

func (c *Connection) Close() error {
	ctx := context.TODO()
	ctx, cancelFunc := utils.WithNamedTimeout(ctx, "connection close timeout", 5*time.Second)
	defer cancelFunc()
	if c.sessionState.connectionId != 0 {
		c.log.Debug().Msg("Sending ForwardClose request")
		forwardClose := readWriteModel.NewCipRRData(
			c.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), c.sessionState.senderContext, 0,
			EmptyInterfaceHandle, 0,
			[]readWriteModel.TypeId{
				readWriteModel.NewNullAddressItem(),
				readWriteModel.NewUnConnectedDataItem(readWriteModel.NewCipConnectionManagerCloseRequest(
					2,
					readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6)),
					readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1)),
					0, 10, 14,
					c.sessionState.connectionSerialNumber, 4919, 42,
					c.sessionState.connectionPathSize, c.sessionState.routingAddress,
				)),
			},
		)
		// Many devices close the socket right after this - errors are expected and ignored.
		if err := c.messageCodec.Send(ctx, "forward_close", forwardClose); err != nil {
			c.log.Debug().Err(err).Msg("error sending forward close")
		}
		c.sessionState.connectionId = 0
	}
	c.log.Debug().Msg("Sending UnregisterSession EIP Packet")
	if err := c.messageCodec.Send(ctx, "unregister_session",
		readWriteModel.NewEipDisconnectRequest(c.sessionState.sessionHandle, 0, []byte(DefaultSenderContext), 0),
	); err != nil {
		c.log.Debug().Err(err).Msg("error sending unregister session request")
	}
	// Unregister gets no response
	time.Sleep(100 * time.Millisecond) // Just to make sure it gets out
	if err := c.messageCodec.Disconnect(); err != nil {
		c.log.Warn().Err(err).Msg("error disconnecting message codec")
	}
	c.log.Debug().
		Uint32("sessionHandle", c.sessionState.sessionHandle).
		Msg("Unregistered session")

	// Wait for background goroutines (e.g. the async setup handshake spawned in
	// Connect) to finish now that the codec is disconnected - pending setup
	// requests will fail fast against the disconnected codec. Bound the wait so
	// a stuck goroutine can't hang Close() forever.
	waitDone := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(waitDone)
	}()
	select {
	case <-waitDone:
	case <-time.After(15 * time.Second):
		c.log.Warn().Msg("timed out waiting for background goroutines to finish during close")
	}
	return nil
}

func (c *Connection) setupConnection(ctx context.Context) error {
	if err := c.listServices(ctx); err != nil {
		return errors.Wrap(err, "error listing services")
	}
	if err := c.registerSession(ctx); err != nil {
		return errors.Wrap(err, "error registering session")
	}
	if c.configuration.forceUnconnectedOperation {
		c.log.Debug().Msg("Unconnected operation forced, skipping the capability probe")
		c.SetConnected(true)
		return nil
	}
	if err := c.probeAttributes(ctx); err != nil {
		// Any probe failure (timeout, error status, malformed reply) means "device has
		// no message router / connection manager": fall back to unconnected operation
		// like plc4j instead of failing the connect.
		c.log.Debug().Err(err).Msg("GetAttributeAll probe failed, falling back to unconnected mode")
		c.sessionState.useMessageRouter = false
		c.sessionState.useConnectionManager = false
	}
	if c.sessionState.useConnectionManager {
		if err := c.openConnectionManager(ctx); err != nil {
			// Deliberate deviation from plc4j (which fails the connect here): a device
			// that advertises a connection manager but rejects the ForwardOpen is still
			// usable through unconnected messaging.
			c.log.Debug().Err(err).Msg("ForwardOpen failed, falling back to unconnected mode")
			c.sessionState.useConnectionManager = false
		}
	}
	c.SetConnected(true)
	return nil
}

func (c *Connection) listServices(ctx context.Context) error {
	c.log.Debug().Msg("Sending ListServices request")
	message, err := sendRequestAndWait(ctx, c.log, c.messageCodec, "list_services",
		readWriteModel.NewListServicesRequest(
			EmptySessionHandle, uint32(readWriteModel.CIPStatus_Success), []byte(DefaultSenderContext), 0,
		), func(message spi.Message) bool {
			_, ok := message.(readWriteModel.EipPacket)
			return ok
		})
	if err != nil {
		return err
	}
	listServicesResponse, ok := message.(readWriteModel.ListServicesResponse)
	if !ok {
		// Like plc4j, tolerate a device that replies with something other than a
		// well-formed ListServicesResponse and just proceed to RegisterSession.
		c.log.Debug().Type("responseType", message).Msg("Device did not reply with a ListServicesResponse, proceeding without it")
		return nil
	}
	if len(listServicesResponse.GetTypeIds()) == 0 {
		c.log.Debug().Msg("ListServices response contains no services, proceeding without it")
		return nil
	}
	servicesResponse, ok := listServicesResponse.GetTypeIds()[0].(readWriteModel.ServicesResponse)
	if !ok {
		c.log.Debug().Type("typeId", listServicesResponse.GetTypeIds()[0]).Msg("Unexpected type id in ListServices response, proceeding without it")
		return nil
	}
	if !servicesResponse.GetSupportsCIPEncapsulation() {
		return errors.New("device does not support CIP encapsulation")
	}
	c.cipEncapsulationAvailable = true
	c.log.Debug().Msg("Device supports CIP over EIP encapsulation")
	return nil
}

func (c *Connection) registerSession(ctx context.Context) error {
	c.log.Debug().Msg("Sending RegisterSession request")
	message, err := sendRequestAndWait(ctx, c.log, c.messageCodec, "register_session",
		readWriteModel.NewEipConnectionRequest(
			EmptySessionHandle, uint32(readWriteModel.CIPStatus_Success), []byte(DefaultSenderContext), 0,
		), func(message spi.Message) bool {
			_, ok := message.(readWriteModel.EipPacket)
			return ok
		})
	if err != nil {
		return err
	}
	connectionResponse, ok := message.(readWriteModel.EipConnectionResponse)
	if !ok {
		// Some devices skip ahead to the connection manager - proceed without a
		// registered session, like plc4j.
		c.log.Debug().Type("responseType", message).Msg("Device skipped session registration")
		return nil
	}
	if connectionResponse.GetStatus() != uint32(readWriteModel.CIPStatus_Success) {
		return errors.Errorf("got status code while registering session [%d]", connectionResponse.GetStatus())
	}
	c.sessionState.sessionHandle = connectionResponse.GetSessionHandle()
	c.sessionState.senderContext = connectionResponse.GetSenderContext()
	c.log.Debug().Uint32("sessionHandle", c.sessionState.sessionHandle).Msg("Got assigned with session handle")
	return nil
}

func (c *Connection) probeAttributes(ctx context.Context) error {
	c.log.Debug().Msg("Sending GetAttributeAll probe")
	probe := readWriteModel.NewCipRRData(
		c.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), c.sessionState.senderContext, 0,
		EmptyInterfaceHandle, 0,
		[]readWriteModel.TypeId{
			readWriteModel.NewNullAddressItem(),
			readWriteModel.NewUnConnectedDataItem(readWriteModel.NewGetAttributeAllRequest(
				readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 2)),
				readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1)),
			)),
		},
	)
	message, err := sendRequestAndWait(ctx, c.log, c.messageCodec, "get_attribute_all", probe,
		func(message spi.Message) bool {
			_, ok := message.(readWriteModel.CipRRData)
			return ok
		})
	if err != nil {
		return err
	}
	cipRRData := message.(readWriteModel.CipRRData)
	if cipRRData.GetStatus() != uint32(readWriteModel.CIPStatus_Success) {
		return errors.Errorf("got status code on GetAttributeAll probe [%d]", cipRRData.GetStatus())
	}
	if len(cipRRData.GetTypeIds()) < 2 {
		return errors.New("GetAttributeAll probe response contains no data item")
	}
	dataItem, ok := cipRRData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
	if !ok {
		return errors.Errorf("unexpected type id in GetAttributeAll probe response: %T", cipRRData.GetTypeIds()[1])
	}
	response, ok := dataItem.GetService().(readWriteModel.GetAttributeAllResponse)
	if !ok {
		return errors.Errorf("unexpected service in GetAttributeAll probe response: %T", dataItem.GetService())
	}
	if response.GetStatus() != uint8(readWriteModel.CIPStatus_Success) {
		return errors.Errorf("got status code on GetAttributeAll response [%d]", response.GetStatus())
	}
	if attributes := response.GetAttributes(); attributes != nil {
		for _, classId := range attributes.GetClassId() {
			if cipClassId, ok := readWriteModel.CIPClassIDByValue(classId); ok {
				switch cipClassId {
				case readWriteModel.CIPClassID_MessageRouter:
					c.sessionState.useMessageRouter = true
				case readWriteModel.CIPClassID_ConnectionManager:
					c.sessionState.useConnectionManager = true
				default:
				}
			}
		}
	}
	c.log.Debug().
		Bool("useMessageRouter", c.sessionState.useMessageRouter).
		Bool("useConnectionManager", c.sessionState.useConnectionManager).
		Msg("Probed device capabilities")
	return nil
}

func (c *Connection) openConnectionManager(ctx context.Context) error {
	c.log.Debug().Msg("Sending ForwardOpen request")
	forwardOpen := readWriteModel.NewCipRRData(
		c.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), c.sessionState.senderContext, 0,
		EmptyInterfaceHandle, 0,
		[]readWriteModel.TypeId{
			readWriteModel.NewNullAddressItem(),
			readWriteModel.NewUnConnectedDataItem(readWriteModel.NewCipConnectionManagerRequest(
				readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6)),
				readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1)),
				0, 10, 14, 536870914, 33944,
				c.sessionState.connectionSerialNumber, 4919, 42, 3, 2101812,
				readWriteModel.NewNetworkConnectionParameters(4002, false, 2, 0, true),
				2113537,
				readWriteModel.NewNetworkConnectionParameters(4002, false, 2, 0, true),
				readWriteModel.NewTransportType(true, 2, 3),
				c.sessionState.connectionPathSize, c.sessionState.routingAddress,
			)),
		},
	)
	message, err := sendRequestAndWait(ctx, c.log, c.messageCodec, "forward_open", forwardOpen,
		func(message spi.Message) bool {
			_, ok := message.(readWriteModel.CipRRData)
			return ok
		})
	if err != nil {
		return err
	}
	cipRRData := message.(readWriteModel.CipRRData)
	if cipRRData.GetStatus() != uint32(readWriteModel.CIPStatus_Success) {
		return errors.Errorf("got status code while opening connection manager [%d]", cipRRData.GetStatus())
	}
	if len(cipRRData.GetTypeIds()) < 2 {
		return errors.New("ForwardOpen response contains no data item")
	}
	dataItem, ok := cipRRData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
	if !ok {
		return errors.Errorf("unexpected type id in ForwardOpen response: %T", cipRRData.GetTypeIds()[1])
	}
	connectionManagerResponse, ok := dataItem.GetService().(readWriteModel.CipConnectionManagerResponse)
	if !ok {
		return errors.Errorf("unexpected service in ForwardOpen response: %T", dataItem.GetService())
	}
	c.sessionState.connectionId = connectionManagerResponse.GetOtConnectionId()
	c.log.Debug().Uint32("connectionId", c.sessionState.connectionId).Msg("Got assigned with connection id")
	return nil
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
		// Stated explicitly rather than left to the zero value: the eip driver implements
		// neither subscribing nor browsing, so both builders fall through to
		// _default.DefaultConnection and panic. Flip these the moment that changes.
		ProvidesSubscribing: false,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			c.messageCodec,
			c.tm,
			c.configuration,
			c.sessionState,
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
			c.sessionState,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) String() string {
	return fmt.Sprintf("eip.Connection")
}
