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
	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=AlphaGenerator
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

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=Connection
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

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

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

func (c *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	c.log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.fireConnectionError(errors.Errorf("panic-ed %v. Stack:\n%s", err, debug.Stack()), ch)
			}
		}()
		if err := c.messageCodec.ConnectWithContext(ctx); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error connecting codec"), ch)
			return
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !c.driverContext.awaitSetupComplete {
			go c.setupConnection(ctx, ch)
			c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
			c.SetConnected(true)
			return
		}

		c.setupConnection(ctx, ch)
	}()
	return ch
}

func (c *Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	results := make(chan plc4go.PlcConnectionCloseResult, 1)
	go func() {
		result := <-c.DefaultConnection.Close()
		c.log.Trace().Msg("Waiting for handlers to stop")
		c.handlerWaitGroup.Wait()
		c.log.Trace().Msg("handlers stopped, dispatching result")
		results <- result
	}()
	return results
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

func (c *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	cbusOptions := &c.messageCodec.cbusOptions
	requestContext := &c.messageCodec.requestContext

	if !c.sendReset(ctx, ch, cbusOptions, requestContext, false) {
		c.log.Warn().Msg("First reset failed")
		// We try a second reset in case we get a power up
		if !c.sendReset(ctx, ch, cbusOptions, requestContext, true) {
			c.log.Trace().Msg("Reset failed")
			return
		}
	}
	if !c.setApplicationFilter(ctx, ch, requestContext, cbusOptions) {
		c.log.Trace().Msg("Set application filter failed")
		return
	}
	if !c.setInterfaceOptions3(ctx, ch, requestContext, cbusOptions) {
		c.log.Trace().Msg("Set interface options 3 failed")
		return
	}
	if !c.setInterface1PowerUpSettings(ctx, ch, requestContext, cbusOptions) {
		c.log.Trace().Msg("Set interface options 1 power up settings failed")
		return
	}
	if !c.setInterfaceOptions1(ctx, ch, requestContext, cbusOptions) {
		c.log.Trace().Msg("Set interface options 1 failed")
		return
	}
	c.log.Trace().Msg("Connection setup done")
	c.fireConnected(ch)
	c.log.Trace().Msg("Connect fired")
	c.startSubscriptionHandler()
	c.log.Trace().Msg("subscription handler started")
}

func (c *Connection) startSubscriptionHandler() {
	c.log.Debug().Msg("Starting SAL handler")
	c.handlerWaitGroup.Add(1)
	go func() {
		salLogger := c.log.With().Str("handlerType", "SAL").Logger()
		defer c.handlerWaitGroup.Done()
		defer func() {
			if err := recover(); err != nil {
				salLogger.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		salLogger.Debug().Msg("SAL handler started")
		for c.IsConnected() {
			for monitoredSal := range c.messageCodec.monitoredSALs {
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
	}()
	c.log.Debug().Msg("Starting MMI handler")
	c.handlerWaitGroup.Add(1)
	go func() {
		mmiLogger := c.log.With().Str("handlerType", "MMI").Logger()
		defer c.handlerWaitGroup.Done()
		defer func() {
			if err := recover(); err != nil {
				mmiLogger.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		mmiLogger.Debug().Msg("default MMI started")
		for c.IsConnected() {
			for calReply := range c.messageCodec.monitoredMMIs {
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
	}()
}

func (c *Connection) sendReset(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, cbusOptions *readWriteModel.CBusOptions, requestContext *readWriteModel.RequestContext, sendOutErrorNotification bool) (ok bool) {
	c.log.Debug().Bool("sendOutErrorNotification", sendOutErrorNotification).Msg("Send a reset")
	requestTypeReset := readWriteModel.RequestType_RESET
	requestReset := readWriteModel.NewRequestReset(requestTypeReset, &requestTypeReset, requestTypeReset, &requestTypeReset, requestTypeReset, nil, &requestTypeReset, requestTypeReset, readWriteModel.NewRequestTermination(), *cbusOptions)
	cBusMessage := readWriteModel.NewCBusMessageToServer(requestReset, *requestContext, *cbusOptions)

	receivedResetEchoChan := make(chan bool, 1)
	receivedResetEchoErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(
		ctx,
		cBusMessage,
		func(message spi.Message) bool {
			c.log.Trace().Msg("Checking message")
			switch message := message.(type) {
			case readWriteModel.CBusMessageToClientExactly:
				switch reply := message.GetReply().(type) {
				case readWriteModel.ReplyOrConfirmationReplyExactly:
					switch reply.GetReply().(type) {
					case readWriteModel.PowerUpReplyExactly:
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
			case readWriteModel.CBusMessageToServerExactly:
				switch request := message.GetRequest().(type) {
				case readWriteModel.RequestResetExactly:
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
		},
		func(message spi.Message) error {
			c.log.Trace().Msg("Handling message")
			switch message.(type) {
			case readWriteModel.CBusMessageToClientExactly:
				// This is the powerup notification
				select {
				case receivedResetEchoChan <- false:
					c.log.Trace().Msg("notified reset chan from message to client")
				default:
				}
			case readWriteModel.CBusMessageToServerExactly:
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
		},
		func(err error) error {
			select {
			case receivedResetEchoErrorChan <- errors.Wrap(err, "got error processing request"):
				c.log.Trace().Msg("notified error chan")
			default:
			}
			return nil
		}, c.GetTtl()); err != nil {
		if sendOutErrorNotification {
			c.fireConnectionError(errors.Wrap(err, "Error during sending of Reset Request"), ch)
		} else {
			c.log.Warn().Err(err).Msg("connect failed")
		}
		return false
	}

	startTime := time.Now()
	timeout := time.NewTimer(time.Millisecond * 500)
	defer utils.CleanupTimer(timeout)
	select {
	case <-receivedResetEchoChan:
		c.log.Debug().Msg("We received the echo")
	case err := <-receivedResetEchoErrorChan:
		if sendOutErrorNotification {
			c.fireConnectionError(errors.Wrap(err, "Error receiving of Reset"), ch)
		} else {
			c.log.Trace().Err(err).Msg("connect failed")
		}
		return false
	case timeout := <-timeout.C:
		if sendOutErrorNotification {
			c.fireConnectionError(errors.Errorf("Timeout after %v", timeout.Sub(startTime)), ch)
		} else {
			c.log.Trace().Dur("timeout", timeout.Sub(startTime)).Msg("Timeout")
		}
		return false
	}
	c.log.Debug().Msg("Reset done")
	return true
}

func (c *Connection) setApplicationFilter(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) (ok bool) {
	c.log.Debug().Msg("Set application filter to all")
	applicationAddress1 := readWriteModel.NewParameterValueApplicationAddress1(readWriteModel.NewApplicationAddress1(c.configuration.MonitoredApplication1), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_APPLICATION_ADDRESS_1, applicationAddress1, requestContext, cbusOptions) {
		return false
	}
	applicationAddress2 := readWriteModel.NewParameterValueApplicationAddress2(readWriteModel.NewApplicationAddress2(c.configuration.MonitoredApplication2), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_APPLICATION_ADDRESS_2, applicationAddress2, requestContext, cbusOptions) {
		return false
	}
	c.log.Debug().Msg("Application filter set")
	return true
}

func (c *Connection) setInterfaceOptions3(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) (ok bool) {
	c.log.Debug().Msg("Set interface options 3")
	interfaceOptions3 := readWriteModel.NewParameterValueInterfaceOptions3(readWriteModel.NewInterfaceOptions3(c.configuration.Exstat, c.configuration.Pun, c.configuration.LocalSal, c.configuration.Pcn), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_INTERFACE_OPTIONS_3, interfaceOptions3, requestContext, cbusOptions) {
		return false
	}
	// TODO: add localsal to the options
	*cbusOptions = readWriteModel.NewCBusOptions(false, false, false, c.configuration.Exstat, false, false, c.configuration.Pun, c.configuration.Pcn, false)
	c.log.Debug().Msg("Interface options 3 set")
	return true
}

func (c *Connection) setInterface1PowerUpSettings(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) (ok bool) {
	c.log.Debug().Msg("Set interface options 1 power up settings")
	interfaceOptions1PowerUpSettings := readWriteModel.NewParameterValueInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1(c.configuration.Idmon, c.configuration.Monitor, c.configuration.Smart, c.configuration.Srchk, c.configuration.XonXoff, c.configuration.Connect)), 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_INTERFACE_OPTIONS_1_POWER_UP_SETTINGS, interfaceOptions1PowerUpSettings, requestContext, cbusOptions) {
		return false
	}
	// TODO: what is with monall
	*cbusOptions = readWriteModel.NewCBusOptions(c.configuration.Connect, c.configuration.Smart, c.configuration.Idmon, c.configuration.Exstat, c.configuration.Monitor, false, c.configuration.Pun, c.configuration.Pcn, c.configuration.Srchk)
	c.log.Debug().Msg("Interface options 1 power up settings set")
	return true
}

func (c *Connection) setInterfaceOptions1(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) bool {
	c.log.Debug().Msg("Set interface options 1")
	interfaceOptions1 := readWriteModel.NewParameterValueInterfaceOptions1(readWriteModel.NewInterfaceOptions1(c.configuration.Idmon, c.configuration.Monitor, c.configuration.Smart, c.configuration.Srchk, c.configuration.XonXoff, c.configuration.Connect), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_INTERFACE_OPTIONS_1, interfaceOptions1, requestContext, cbusOptions) {
		return false
	}
	// TODO: what is with monall
	*cbusOptions = readWriteModel.NewCBusOptions(c.configuration.Connect, c.configuration.Smart, c.configuration.Idmon, c.configuration.Exstat, c.configuration.Monitor, false, c.configuration.Pun, c.configuration.Pcn, c.configuration.Srchk)
	c.log.Debug().Msg("Interface options 1 set")
	return true
}

// This is used for connection setup
func (c *Connection) sendCalDataWrite(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, paramNo readWriteModel.Parameter, parameterValue readWriteModel.ParameterValue, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) bool {
	calCommandTypeContainer := readWriteModel.CALCommandTypeContainer_CALCommandWrite_2Bytes + readWriteModel.CALCommandTypeContainer(parameterValue.GetLengthInBytes(ctx))
	calData := readWriteModel.NewCALDataWrite(paramNo, 0x0, parameterValue, calCommandTypeContainer, nil, *requestContext)
	directCommand := readWriteModel.NewRequestDirectCommandAccess(calData /*we don't want an alpha otherwise the PCI will auto-switch*/, nil, 0x40, nil, nil, 0x0, readWriteModel.NewRequestTermination(), *cbusOptions)
	cBusMessage := readWriteModel.NewCBusMessageToServer(directCommand, *requestContext, *cbusOptions)

	directCommandAckChan := make(chan bool, 1)
	directCommandAckErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, cBusMessage, func(message spi.Message) bool {
		switch message := message.(type) {
		case readWriteModel.CBusMessageToClientExactly:
			switch reply := message.GetReply().(type) {
			case readWriteModel.ReplyOrConfirmationReplyExactly:
				switch reply := reply.GetReply().(type) {
				case readWriteModel.ReplyEncodedReplyExactly:
					switch encodedReply := reply.GetEncodedReply().(type) {
					case readWriteModel.EncodedReplyCALReplyExactly:
						switch data := encodedReply.GetCalReply().GetCalData().(type) {
						case readWriteModel.CALDataAcknowledgeExactly:
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
		case readWriteModel.CBusMessageToClientExactly:
			switch reply := message.GetReply().(type) {
			case readWriteModel.ReplyOrConfirmationReplyExactly:
				switch reply := reply.GetReply().(type) {
				case readWriteModel.ReplyEncodedReplyExactly:
					switch encodedReply := reply.GetEncodedReply().(type) {
					case readWriteModel.EncodedReplyCALReplyExactly:
						switch data := encodedReply.GetCalReply().GetCalData().(type) {
						case readWriteModel.CALDataAcknowledgeExactly:
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
		select {
		case directCommandAckErrorChan <- errors.Wrap(err, "got error processing request"):
		default:
		}
		return nil
	}, c.GetTtl()); err != nil {
		c.fireConnectionError(errors.Wrap(err, "Error during sending of write request"), ch)
		return false
	}

	startTime := time.Now()
	timeout := time.NewTimer(2 * time.Second)
	defer utils.CleanupTimer(timeout)
	select {
	case <-directCommandAckChan:
		c.log.Debug().Msg("We received the ack")
	case err := <-directCommandAckErrorChan:
		c.fireConnectionError(errors.Wrap(err, "Error receiving of ack"), ch)
		return false
	case timeout := <-timeout.C:
		c.fireConnectionError(errors.Errorf("Timeout after %v", timeout.Sub(startTime)), ch)
		return false
	}
	return true
}

func (c *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
	} else {
		c.log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
	if err := c.messageCodec.Disconnect(); err != nil {
		c.log.Debug().Err(err).Msg("Error disconnecting message codec on connection error")
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
