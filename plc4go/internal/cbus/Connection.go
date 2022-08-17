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
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/plcerrors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type AlphaGenerator struct {
	currentAlpha byte
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

type Connection struct {
	_default.DefaultConnection
	alphaGenerator AlphaGenerator
	messageCodec   spi.MessageCodec
	subscribers    []*Subscriber
	tm             *spi.RequestTransactionManager

	configuration Configuration
	driverContext DriverContext

	connectionId string
	tracer       *spi.Tracer
}

func NewConnection(messageCodec spi.MessageCodec, configuration Configuration, driverContext DriverContext, fieldHandler spi.PlcFieldHandler, tm *spi.RequestTransactionManager, options map[string][]string) *Connection {
	connection := &Connection{
		alphaGenerator: AlphaGenerator{currentAlpha: 'g'},
		messageCodec:   messageCodec,
		configuration:  configuration,
		driverContext:  driverContext,
		tm:             tm,
	}
	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcFieldHandler(fieldHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
	)
	return connection
}

func (c *Connection) GetConnectionId() string {
	return c.connectionId
}

func (c *Connection) IsTraceEnabled() bool {
	return c.tracer != nil
}

func (c *Connection) GetTracer() *spi.Tracer {
	return c.tracer
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
}

func (c *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	// TODO: use proper context
	ctx := context.TODO()
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		err := c.messageCodec.Connect()
		if err != nil {
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !c.driverContext.awaitSetupComplete {
			go c.setupConnection(ctx, ch)
			log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
			c.SetConnected(true)
			return
		}

		c.setupConnection(ctx, ch)
	}()
	return ch
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    true,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(c.GetPlcFieldHandler(), NewReader(&c.alphaGenerator, c.messageCodec, c.tm))
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(c.GetPlcFieldHandler(), c.GetPlcValueHandler(), NewWriter(&c.alphaGenerator, c.messageCodec, c.tm))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return internalModel.NewDefaultPlcSubscriptionRequestBuilder(c.GetPlcFieldHandler(), c.GetPlcValueHandler(), NewSubscriber(c))
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	// TODO: where do we get the unsubscriber from
	return nil
}

func (c *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return internalModel.NewDefaultPlcBrowseRequestBuilder(c.GetPlcFieldHandler(), NewBrowser(c, c.messageCodec))
}

func (c *Connection) addSubscriber(subscriber *Subscriber) {
	for _, sub := range c.subscribers {
		if sub == subscriber {
			log.Debug().Msgf("Subscriber %v already added", subscriber)
			return
		}
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) String() string {
	return fmt.Sprintf("cbus.Connection")
}

func (c *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	cbusOptions := &c.messageCodec.(*MessageCodec).cbusOptions
	requestContext := &c.messageCodec.(*MessageCodec).requestContext

	if !c.sendReset(ctx, ch, cbusOptions, requestContext, false) {
		log.Warn().Msg("First reset failed")
		// We try a second reset in case we get a power up
		if !c.sendReset(ctx, ch, cbusOptions, requestContext, true) {
			log.Trace().Msg("Reset failed")
			return
		}
	}
	if !c.setApplicationFilter(ctx, ch, requestContext, cbusOptions) {
		log.Trace().Msg("Set application filter failed")
		return
	}
	if !c.setInterfaceOptions3(ctx, ch, requestContext, cbusOptions) {
		log.Trace().Msg("Set interface options 3 failed")
		return
	}
	if !c.setInterface1PowerUpSettings(ctx, ch, requestContext, cbusOptions) {
		log.Trace().Msg("Set interface options 1 power up settings failed")
		return
	}
	if !c.setInterfaceOptions1(ctx, ch, requestContext, cbusOptions) {
		log.Trace().Msg("Set interface options 1 failed")
		return
	}
	c.fireConnected(ch)

	log.Debug().Msg("Starting subscription handler")
	go func() {
		log.Debug().Msg("Subscription handler stated")
		for c.IsConnected() {
			for monitoredSal := range c.messageCodec.(*MessageCodec).monitoredSALs {
				for _, subscriber := range c.subscribers {
					if ok := subscriber.handleMonitoredSal(monitoredSal); ok {
						log.Debug().Msgf("%v handled\n%s", subscriber, monitoredSal)
					}
				}
			}
		}
		log.Info().Msg("Ending subscription handler")
	}()

	log.Debug().Msg("Starting default incoming message handler")
	go func() {
		log.Debug().Msg("default incoming message handler started")
		for c.IsConnected() {
			for message := range c.messageCodec.GetDefaultIncomingMessageChannel() {
				switch message := message.(type) {
				case readWriteModel.CBusMessageToClientExactly:
					switch reply := message.GetReply().(type) {
					case readWriteModel.ReplyOrConfirmationReplyExactly:
						switch reply := reply.GetReply().(type) {
						case readWriteModel.ReplyEncodedReplyExactly:
							switch encodedReply := reply.GetEncodedReply().(type) {
							case readWriteModel.EncodedReplyCALReplyExactly:
								for _, subscriber := range c.subscribers {
									calReply := encodedReply.GetCalReply()
									if ok := subscriber.handleMonitoredMMI(calReply); ok {
										log.Debug().Msgf("%v handled\n%s", subscriber, calReply)
										continue
									}
								}
							}
						}
					}
				}
				log.Debug().Msgf("Received unhandled \n%v", message)
			}
		}
		log.Info().Msg("Ending default incoming message handler")
	}()
}

func (c *Connection) sendReset(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, cbusOptions *readWriteModel.CBusOptions, requestContext *readWriteModel.RequestContext, sendOutErrorNotification bool) (ok bool) {
	log.Debug().Msgf("Send a reset (sendOutErrorNotification: %t)", sendOutErrorNotification)
	requestTypeReset := readWriteModel.RequestType_RESET
	requestReset := readWriteModel.NewRequestReset(requestTypeReset, &requestTypeReset, requestTypeReset, &requestTypeReset, requestTypeReset, nil, &requestTypeReset, requestTypeReset, readWriteModel.NewRequestTermination(), *cbusOptions)
	cBusMessage := readWriteModel.NewCBusMessageToServer(requestReset, *requestContext, *cbusOptions)

	receivedResetEchoChan := make(chan bool)
	receivedResetEchoErrorChan := make(chan error)
	if err := c.messageCodec.SendRequest(ctx, cBusMessage, func(message spi.Message) bool {
		switch message := message.(type) {
		case readWriteModel.CBusMessageToClientExactly:
			if reply, ok := message.GetReply().(readWriteModel.ReplyOrConfirmationReplyExactly); ok {
				_, ok := reply.GetReply().(readWriteModel.PowerUpReplyExactly)
				return ok
			}
		case readWriteModel.CBusMessageToServerExactly:
			_, ok = message.GetRequest().(readWriteModel.RequestResetExactly)
			return ok
		}
		return false
	}, func(message spi.Message) error {
		switch message.(type) {
		case readWriteModel.CBusMessageToClientExactly:
			// This is the powerup notification
			go func() { receivedResetEchoChan <- false }()
		case readWriteModel.CBusMessageToServerExactly:
			// This is the echo
			go func() { receivedResetEchoChan <- true }()
		default:
			return errors.Errorf("Unmapped type %T", message)
		}
		return nil
	}, func(err error) error {
		// If this is a timeout, do a check if the connection requires a reconnection
		if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
			log.Warn().Msg("Timeout during Connection establishing, closing channel...")
			c.Close()
		}
		receivedResetEchoErrorChan <- errors.Wrap(err, "got error processing request")
		return nil
	}, c.GetTtl()); err != nil {
		if sendOutErrorNotification {
			c.fireConnectionError(errors.Wrap(err, "Error during sending of Reset Request"), ch)
		} else {
			log.Warn().Err(err).Msg("connect failed")
		}
		return false
	}

	startTime := time.Now()
	select {
	case <-receivedResetEchoChan:
		log.Debug().Msgf("We received the echo")
	case err := <-receivedResetEchoErrorChan:
		if sendOutErrorNotification {
			c.fireConnectionError(errors.Wrap(err, "Error receiving of Reset"), ch)
		} else {
			log.Trace().Err(err).Msg("connect failed")
		}
		return false
	case timeout := <-time.After(time.Millisecond * 500):
		if sendOutErrorNotification {
			c.fireConnectionError(errors.Errorf("Timeout after %v", timeout.Sub(startTime)), ch)
		} else {
			log.Trace().Msg("timeout")
		}
		return false
	}
	log.Debug().Msg("Reset done")
	return true
}

func (c *Connection) setApplicationFilter(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) (ok bool) {
	log.Debug().Msg("Set application filter to all")
	applicationAddress1 := readWriteModel.NewParameterValueApplicationAddress1(readWriteModel.NewApplicationAddress1(c.configuration.MonitoredApplication1), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_APPLICATION_ADDRESS_1, applicationAddress1, requestContext, cbusOptions) {
		return false
	}
	applicationAddress2 := readWriteModel.NewParameterValueApplicationAddress2(readWriteModel.NewApplicationAddress2(c.configuration.MonitoredApplication2), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_APPLICATION_ADDRESS_2, applicationAddress2, requestContext, cbusOptions) {
		return false
	}
	log.Debug().Msg("Application filter set")
	return true
}

func (c *Connection) setInterfaceOptions3(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) (ok bool) {
	log.Debug().Msg("Set interface options 3")
	interfaceOptions3 := readWriteModel.NewParameterValueInterfaceOptions3(readWriteModel.NewInterfaceOptions3(c.configuration.Exstat, c.configuration.Pun, c.configuration.LocalSal, c.configuration.Pcn), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_INTERFACE_OPTIONS_3, interfaceOptions3, requestContext, cbusOptions) {
		return false
	}
	// TODO: add localsal to the options
	*cbusOptions = readWriteModel.NewCBusOptions(false, false, false, c.configuration.Exstat, false, false, c.configuration.Pun, c.configuration.Pcn, false)
	log.Debug().Msg("Interface options 3 set")
	return true
}

func (c *Connection) setInterface1PowerUpSettings(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) (ok bool) {
	log.Debug().Msg("Set interface options 1 power up settings")
	interfaceOptions1PowerUpSettings := readWriteModel.NewParameterValueInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1(c.configuration.Idmon, c.configuration.Monitor, c.configuration.Smart, c.configuration.Srchk, c.configuration.XonXoff, c.configuration.Connect)), 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_INTERFACE_OPTIONS_1_POWER_UP_SETTINGS, interfaceOptions1PowerUpSettings, requestContext, cbusOptions) {
		return false
	}
	// TODO: what is with monall
	*cbusOptions = readWriteModel.NewCBusOptions(c.configuration.Connect, c.configuration.Smart, c.configuration.Idmon, c.configuration.Exstat, c.configuration.Monitor, false, c.configuration.Pun, c.configuration.Pcn, c.configuration.Srchk)
	log.Debug().Msg("Interface options 1 power up settings set")
	return true
}

func (c *Connection) setInterfaceOptions1(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) bool {
	log.Debug().Msg("Set interface options 1")
	interfaceOptions1 := readWriteModel.NewParameterValueInterfaceOptions1(readWriteModel.NewInterfaceOptions1(c.configuration.Idmon, c.configuration.Monitor, c.configuration.Smart, c.configuration.Srchk, c.configuration.XonXoff, c.configuration.Connect), nil, 1)
	if !c.sendCalDataWrite(ctx, ch, readWriteModel.Parameter_INTERFACE_OPTIONS_1, interfaceOptions1, requestContext, cbusOptions) {
		return false
	}
	// TODO: what is with monall
	*cbusOptions = readWriteModel.NewCBusOptions(c.configuration.Connect, c.configuration.Smart, c.configuration.Idmon, c.configuration.Exstat, c.configuration.Monitor, false, c.configuration.Pun, c.configuration.Pcn, c.configuration.Srchk)
	log.Debug().Msg("Interface options 1 set")
	return true
}

// This is used for connection setup
func (c *Connection) sendCalDataWrite(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult, paramNo readWriteModel.Parameter, parameterValue readWriteModel.ParameterValue, requestContext *readWriteModel.RequestContext, cbusOptions *readWriteModel.CBusOptions) bool {
	// TODO: we assume that is always a one byte request otherwise we need to map the length here
	calData := readWriteModel.NewCALDataWrite(paramNo, 0x0, parameterValue, readWriteModel.CALCommandTypeContainer_CALCommandWrite_3Bytes, nil, *requestContext)
	directCommand := readWriteModel.NewRequestDirectCommandAccess(calData /*we don't want a alpha otherwise the PCI will auto-switch*/, nil, 0x40, nil, nil, 0x0, readWriteModel.NewRequestTermination(), *cbusOptions)
	cBusMessage := readWriteModel.NewCBusMessageToServer(directCommand, *requestContext, *cbusOptions)

	directCommandAckChan := make(chan bool)
	directCommandAckErrorChan := make(chan error)
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
								directCommandAckChan <- true
							}
						}
					}
				}
			}
		}
		return nil
	}, func(err error) error {
		// If this is a timeout, do a check if the connection requires a reconnection
		if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
			log.Warn().Msg("Timeout during Connection establishing, closing channel...")
			c.Close()
		}
		directCommandAckErrorChan <- errors.Wrap(err, "got error processing request")
		return nil
	}, c.GetTtl()); err != nil {
		c.fireConnectionError(errors.Wrap(err, "Error during sending of write request"), ch)
		return false
	}

	startTime := time.Now()
	select {
	case <-directCommandAckChan:
		log.Debug().Msgf("We received the ack")
	case err := <-directCommandAckErrorChan:
		c.fireConnectionError(errors.Wrap(err, "Error receiving of ack"), ch)
		return false
	case timeout := <-time.After(time.Second * 2):
		c.fireConnectionError(errors.Errorf("Timeout after %v", timeout.Sub(startTime)), ch)
		return false
	}
	return true
}

func (c *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
	} else {
		log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
}

func (c *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
	} else {
		log.Info().Msg("Successfully connected")
	}
	c.SetConnected(true)
}
