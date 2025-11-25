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

package cbus

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

//go:generate go tool plc4xGenerator -type=AlphaGenerator
type AlphaGenerator struct {
	currentAlpha byte `hasLocker:"lock"`
	lock         sync.Mutex
}

func (t *AlphaGenerator) getAndIncrement() byte {
	t.lock.Lock()
	defer t.lock.Unlock()
	// If we've reached the max value 'z', reset back to 'g'
	if t.currentAlpha > 'z' {
		t.currentAlpha = 'g'
	}
	result := t.currentAlpha
	t.currentAlpha += 1
	return result
}

//go:generate go tool plc4xGenerator -type=Connection
type Connection struct {
	_default.DefaultConnection

	alphaGenerator AlphaGenerator `stringer:"true"`
	messageCodec   *MessageCodec
	subscribers    []*Subscriber
	tm             transactions.RequestTransactionManager

	configuration Configuration `stringer:"true"`
	driverContext DriverContext `stringer:"true"`

	handlerWaitGroup sync.WaitGroup

	connectionId string
	tracer       tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(messageCodec *MessageCodec, configuration Configuration, driverContext DriverContext, tagHandler spi.PlcTagHandler, tm transactions.RequestTransactionManager, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		alphaGenerator: AlphaGenerator{currentAlpha: 'g'},
		messageCodec:   messageCodec,
		configuration:  configuration,
		driverContext:  driverContext,
		tm:             tm,

		log:      customLogger,
		_options: _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(
		connection,
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

func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "Error connecting codec")
	}

	// For testing purposes we can skip the waiting for a complete connection
	if !c.driverContext.awaitSetupComplete {
		c.wg.Go(func() {
			if err := c.setupConnection(ctx); err != nil {
				c.log.Error().Err(err).Msg("Error during setup")
			}
		})
		c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		// Here we write directly and don't wait till the connection is "really" connected
		// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
		c.SetConnected(true)
		return nil
	}

	return c.setupConnection(ctx)
}

func (c *Connection) Close() error {
	err := c.DefaultConnection.Close()
	c.log.Trace().Msg("Waiting for handlers to stop")
	c.handlerWaitGroup.Wait()
	c.log.Trace().Msg("handlers stopped, dispatching result")
	return err
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    true,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			&c.alphaGenerator,
			c.messageCodec,
			c.tm,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), NewWriter(&c.alphaGenerator, c.messageCodec, c.tm))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewSubscriber(
			c.addSubscriber,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return spiModel.NewDefaultPlcBrowseRequestBuilder(
		c.GetPlcTagHandler(),
		NewBrowser(
			c,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) addSubscriber(subscriber *Subscriber) {
	for _, sub := range c.subscribers {
		if sub == subscriber {
			c.log.Debug().
				Stringer("subscriber", subscriber).
				Msg("Subscriber already added")
			return
		}
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) setupConnection(ctx context.Context) error {
	cbusOptions := &c.messageCodec.cbusOptions
	requestContext := &c.messageCodec.requestContext

	c.log.Trace().Msg("Starting connection setup")
	c.log.Trace().Msg("Sending reset")
	if err := c.sendReset(ctx); err != nil {
		c.log.Warn().Err(err).Msg("First reset failed")
		// We try a second reset in case we get a power up
		if err := c.sendReset(ctx); err != nil {
			return errors.Wrap(err, "Error during sending of Reset Request")
		}
	}
	c.log.Trace().Msg("setting application filter")
	if err := c.setApplicationFilter(ctx, requestContext); err != nil {
		c.log.Trace().Msg("Set application filter failed")
		return errors.Wrap(err, "Error during sending of Application Filter")
	}
	c.log.Trace().Msg("Setting interface options 3")
	if err := c.setInterfaceOptions3(ctx, requestContext, cbusOptions); err != nil {
		c.log.Trace().Msg("Set interface options 3 failed")
		return errors.Wrap(err, "Error during sending of Interface Options 3")
	}
	c.log.Trace().Msg("Setting interface options 1 power up settings")
	if err := c.setInterface1PowerUpSettings(ctx, requestContext, cbusOptions); err != nil {
		c.log.Trace().Msg("Set interface options 1 power up settings failed")
		return errors.Wrap(err, "Error during sending of Interface Options 1 Power Up Settings")
	}
	c.log.Trace().Msg("Setting interface options 1")
	if err := c.setInterfaceOptions1(ctx, requestContext, cbusOptions); err != nil {
		c.log.Trace().Msg("Set interface options 1 failed")
		return errors.Wrap(err, "Error during sending of Interface Options 1")
	}
	c.log.Trace().Msg("Connection setup done")
	c.SetConnected(true)
	c.startSubscriptionHandler()
	c.log.Trace().Msg("subscription handler started")
	return nil
}

func (c *Connection) startSubscriptionHandler() {
	c.log.Debug().Msg("Starting SAL handler")
	c.handlerWaitGroup.Go(c.handleSAL)
	c.log.Debug().Msg("Starting MMI handler")
	c.handlerWaitGroup.Go(c.handleMMI)
}

func (c *Connection) handleSAL() {
	salLogger := c.log.With().Str("handlerType", "SAL").Logger()
	defer func() {
		if err := recover(); err != nil {
			salLogger.Error().
				Str("stack", string(debug.Stack())).
				Interface("err", err).
				Msg("panic-ed")
		}
	}()
	salLogger.Debug().Msg("SAL handler started")
	for c.IsConnected() && c.messageCodec.IsRunning() {
		for monitoredSal := range c.messageCodec.monitoredSALs {
			if monitoredSal == nil {
				salLogger.Trace().Msg("monitoredSal chan closed")
				break
			}
			salLogger.Trace().
				Stringer("monitoredSal", monitoredSal).
				Msg("got a SAL")
			handled := false
			for _, subscriber := range c.subscribers {
				if ok := subscriber.handleMonitoredSAL(monitoredSal); ok {
					salLogger.Debug().
						Stringer("monitoredSal", monitoredSal).
						Stringer("subscriber", subscriber).
						Msg("handled")
					handled = true
				}
			}
			if !handled {
				salLogger.Debug().
					Stringer("monitoredSal", monitoredSal).
					Msg("SAL was not handled")
			}
		}
	}
	salLogger.Info().Msg("handler ended")
}

func (c *Connection) handleMMI() {
	mmiLogger := c.log.With().Str("handlerType", "MMI").Logger()
	defer func() {
		if err := recover(); err != nil {
			mmiLogger.Error().
				Str("stack", string(debug.Stack())).
				Interface("err", err).
				Msg("panic-ed")
		}
	}()
	mmiLogger.Debug().Msg("default MMI started")
	for c.IsConnected() && c.messageCodec.IsRunning() {
		for calReply := range c.messageCodec.monitoredMMIs {
			if calReply == nil {
				mmiLogger.Trace().Msg("channel closed")
				break
			}
			mmiLogger.Trace().Msg("got a MMI")
			handled := false
			for _, subscriber := range c.subscribers {
				if ok := subscriber.handleMonitoredMMI(calReply); ok {
					mmiLogger.Debug().
						Stringer("subscriber", subscriber).
						Msg("handled")
					handled = true
				}
			}
			if !handled {
				mmiLogger.Debug().Msg("MMI was not handled")
			}
		}
	}
	mmiLogger.Info().Msg("handler ended")
}

func (c *Connection) sendReset(ctx context.Context) error {
	c.log.Debug().Msg("Send a reset")
	requestTypeReset := readWriteModel.RequestType_RESET
	requestReset := readWriteModel.NewRequestReset(
		requestTypeReset,
		nil,
		&requestTypeReset,
		requestTypeReset,
		readWriteModel.NewRequestTermination(),
		requestTypeReset,
		&requestTypeReset,
		requestTypeReset,
		&requestTypeReset,
	)
	cBusMessage := readWriteModel.NewCBusMessageToServer(requestReset)

	receivedResetEchoChan := make(chan bool, 1)
	receivedResetEchoErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, "send_reset", cBusMessage, func(message spi.Message) bool {
		c.log.Trace().Msg("Checking message")
		switch message := message.(type) {
		case readWriteModel.CBusMessageToClient:
			switch reply := message.GetReply().(type) {
			case readWriteModel.ReplyOrConfirmationReply:
				switch reply.GetReply().(type) {
				case readWriteModel.PowerUpReply:
					c.log.Debug().Msg("Received a PUN reply")
					return true
				default:
					c.log.Trace().Type("reply", reply).Msg("not relevant")
					return false
				}
			default:
				c.log.Trace().Type("reply", reply).Msg("not relevant")
				return false
			}
		case readWriteModel.CBusMessageToServer:
			switch request := message.GetRequest().(type) {
			case readWriteModel.RequestReset:
				c.log.Debug().Msg("Received a Reset reply")
				return true
			default:
				c.log.Trace().Type("request", request).Msg("not relevant")
				return false
			}
		default:
			c.log.Trace().Type("message", message).Msg("not relevant")
			return false
		}
	}, func(message spi.Message) error {
		c.log.Trace().Msg("Handling message")
		switch message.(type) {
		case readWriteModel.CBusMessageToClient:
			// This is the powerup notification
			select {
			case receivedResetEchoChan <- false:
				c.log.Trace().Msg("notified reset chan from message to client")
			default:
			}
		case readWriteModel.CBusMessageToServer:
			// This is the echo
			select {
			case receivedResetEchoChan <- true:
				c.log.Trace().Msg("notified reset chan from message to server")
			default:
			}
		default:
			return errors.Errorf("Unmapped type %T", message)
		}
		return nil
	}, func(err error) error {
		select {
		case receivedResetEchoErrorChan <- errors.Wrap(err, "got error processing request"):
			c.log.Trace().Msg("notified error chan")
		default:
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "Error during sending of Reset Request")
	}

	startTime := time.Now()
	select {
	case <-receivedResetEchoChan:
		c.log.Debug().Msg("We received the echo")
	case err := <-receivedResetEchoErrorChan:
		return errors.Wrap(err, "Error receiving of Reset")
	case <-ctx.Done():
		return errors.Errorf("Timeout after %v", time.Since(startTime))
	}
	c.log.Debug().Msg("Reset done")
	return nil
}

func (c *Connection) setApplicationFilter(ctx context.Context, requestContext *readWriteModel.RequestContext) error {
	c.log.Debug().Msg("Set application filter to all")
	applicationAddress1 := readWriteModel.NewParameterValueApplicationAddress1(readWriteModel.NewApplicationAddress1(c.configuration.MonitoredApplication1), nil)
	if err := c.sendCalDataWrite(ctx, readWriteModel.Parameter_APPLICATION_ADDRESS_1, applicationAddress1, requestContext); err != nil {
		return errors.Wrap(err, "Error during sending of Application Address 1")
	}
	applicationAddress2 := readWriteModel.NewParameterValueApplicationAddress2(readWriteModel.NewApplicationAddress2(c.configuration.MonitoredApplication2), nil)
	if err := c.sendCalDataWrite(ctx, readWriteModel.Parameter_APPLICATION_ADDRESS_2, applicationAddress2, requestContext); err != nil {
		return errors.New("Error during sending of Application Address 2")
	}
	c.log.Debug().Msg("Application filter set")
	return nil
}

func (c *Connection) setInterfaceOptions3(ctx context.Context, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) error {
	c.log.Debug().Msg("Set interface options 3")
	interfaceOptions3 := readWriteModel.NewParameterValueInterfaceOptions3(readWriteModel.NewInterfaceOptions3(c.configuration.Exstat, c.configuration.Pun, c.configuration.LocalSal, c.configuration.Pcn), nil)
	if err := c.sendCalDataWrite(ctx, readWriteModel.Parameter_INTERFACE_OPTIONS_3, interfaceOptions3, requestContext); err != nil {
		return errors.Wrap(err, "Error during sending of Interface Options 3")
	}
	// TODO: add localsal to the options
	*cbusOptions = readWriteModel.NewCBusOptions(false, false, false, c.configuration.Exstat, false, false, c.configuration.Pun, c.configuration.Pcn, false)
	c.log.Debug().Msg("Interface options 3 set")
	return nil
}

func (c *Connection) setInterface1PowerUpSettings(ctx context.Context, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) error {
	c.log.Debug().Msg("Set interface options 1 power up settings")
	interfaceOptions1PowerUpSettings := readWriteModel.NewParameterValueInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1(c.configuration.Idmon, c.configuration.Monitor, c.configuration.Smart, c.configuration.Srchk, c.configuration.XonXoff, c.configuration.Connect)))
	if err := c.sendCalDataWrite(ctx, readWriteModel.Parameter_INTERFACE_OPTIONS_1_POWER_UP_SETTINGS, interfaceOptions1PowerUpSettings, requestContext); err != nil {
		return errors.Wrap(err, "Error during sending of Interface Options 1 Power Up Settings")
	}
	// TODO: what is with monall
	*cbusOptions = readWriteModel.NewCBusOptions(c.configuration.Connect, c.configuration.Smart, c.configuration.Idmon, c.configuration.Exstat, c.configuration.Monitor, false, c.configuration.Pun, c.configuration.Pcn, c.configuration.Srchk)
	c.log.Debug().Msg("Interface options 1 power up settings set")
	return nil
}

func (c *Connection) setInterfaceOptions1(ctx context.Context, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) error {
	c.log.Debug().Msg("Set interface options 1")
	interfaceOptions1 := readWriteModel.NewParameterValueInterfaceOptions1(readWriteModel.NewInterfaceOptions1(c.configuration.Idmon, c.configuration.Monitor, c.configuration.Smart, c.configuration.Srchk, c.configuration.XonXoff, c.configuration.Connect), nil)
	if err := c.sendCalDataWrite(ctx, readWriteModel.Parameter_INTERFACE_OPTIONS_1, interfaceOptions1, requestContext); err != nil {
		return errors.Wrap(err, "Error during sending of Interface Options 1")
	}
	// TODO: what is with monall
	*cbusOptions = readWriteModel.NewCBusOptions(c.configuration.Connect, c.configuration.Smart, c.configuration.Idmon, c.configuration.Exstat, c.configuration.Monitor, false, c.configuration.Pun, c.configuration.Pcn, c.configuration.Srchk)
	c.log.Debug().Msg("Interface options 1 set")
	return nil
}

// This is used for connection setup
func (c *Connection) sendCalDataWrite(ctx context.Context, paramNo readWriteModel.Parameter, parameterValue readWriteModel.ParameterValue, requestContext *readWriteModel.RequestContext) error {
	calCommandTypeContainer := readWriteModel.CALCommandTypeContainer_CALCommandWrite_2Bytes + readWriteModel.CALCommandTypeContainer(parameterValue.GetLengthInBytes(ctx))
	calData := readWriteModel.NewCALDataWrite(
		*requestContext,
		calCommandTypeContainer,
		nil,
		paramNo,
		0x0,
		parameterValue,
	)
	directCommand := readWriteModel.NewRequestDirectCommandAccess(
		0x40,
		nil,
		nil,
		0x0,
		readWriteModel.NewRequestTermination(),
		calData,
		/*we don't want an alpha otherwise the PCI will auto-switch*/ nil,
	)
	cBusMessage := readWriteModel.NewCBusMessageToServer(directCommand)

	directCommandAckChan := make(chan bool, 1)
	directCommandAckErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, "send_cal_data_write", cBusMessage, func(message spi.Message) bool {
		switch message := message.(type) {
		case readWriteModel.CBusMessageToClient:
			switch reply := message.GetReply().(type) {
			case readWriteModel.ReplyOrConfirmationReply:
				switch reply := reply.GetReply().(type) {
				case readWriteModel.ReplyEncodedReply:
					switch encodedReply := reply.GetEncodedReply().(type) {
					case readWriteModel.EncodedReplyCALReply:
						switch data := encodedReply.GetCalReply().GetCalData().(type) {
						case readWriteModel.CALDataAcknowledge:
							if data.GetParamNo() == paramNo {
								return true
							}
						}
					}
				}
			}
		}
		return false
	}, func(message spi.Message) error {
		switch message := message.(type) {
		case readWriteModel.CBusMessageToClient:
			switch reply := message.GetReply().(type) {
			case readWriteModel.ReplyOrConfirmationReply:
				switch reply := reply.GetReply().(type) {
				case readWriteModel.ReplyEncodedReply:
					switch encodedReply := reply.GetEncodedReply().(type) {
					case readWriteModel.EncodedReplyCALReply:
						switch data := encodedReply.GetCalReply().GetCalData().(type) {
						case readWriteModel.CALDataAcknowledge:
							if data.GetParamNo() == paramNo {
								select {
								case directCommandAckChan <- true:
								default:
								}
							}
						}
					}
				}
			}
		}
		return nil
	}, func(err error) error {
		c.log.Trace().Err(err).Msg("got error processing request")
		select {
		case directCommandAckErrorChan <- errors.Wrap(err, "got error processing request"):
			c.log.Trace().Msg("error redirected")
		default:
			c.log.Trace().Err(err).Msg("error discarded")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "Error during sending of write request")
	}

	startTime := time.Now()
	select {
	case <-directCommandAckChan:
		c.log.Trace().Msg("We received the ack")
	case err := <-directCommandAckErrorChan:
		return errors.Wrap(err, "Error receiving of ack")
	case <-ctx.Done():
		return errors.Wrapf(ctx.Err(), "Timeout after %v", time.Since(startTime))
	}
	return nil
}
