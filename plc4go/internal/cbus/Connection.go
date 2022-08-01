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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
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
	if t.currentAlpha >= 'z' {
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
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		err := c.messageCodec.Connect()
		if err != nil {
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !c.driverContext.awaitSetupComplete {
			go c.setupConnection(ch)
			log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
			c.SetConnected(true)
			return
		}

		c.setupConnection(ch)
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
	// TODO: where do we get the browser from
	return internalModel.NewDefaultPlcBrowseRequestBuilder(nil)
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

func (c *Connection) setupConnection(ch chan plc4go.PlcConnectionConnectResult) {
	cbusOptions := &c.messageCodec.(*MessageCodec).cbusOptions
	requestContext := &c.messageCodec.(*MessageCodec).requestContext

	{
		log.Debug().Msg("Send a reset Request")
		requestTypeReset := readWriteModel.RequestType_RESET
		requestTypeResetByte := byte(readWriteModel.RequestType_RESET)
		requestReset := readWriteModel.NewRequestReset(requestTypeReset, &requestTypeResetByte, requestTypeReset, &requestTypeResetByte, requestTypeReset, nil, &requestTypeReset, requestTypeReset, readWriteModel.NewRequestTermination(), *cbusOptions)
		if err := c.messageCodec.Send(readWriteModel.NewCBusMessageToServer(requestReset, *requestContext, *cbusOptions)); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error writing reset"), ch)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	{
		log.Debug().Msg("Set application filter to all")
		applicationAddress1 := readWriteModel.NewParameterValueApplicationAddress1(readWriteModel.NewApplicationAddress1(0xFF), 1)
		calData := readWriteModel.NewCALDataWrite(readWriteModel.Parameter_APPLICATION_ADDRESS_1, 0x0, applicationAddress1, readWriteModel.CALCommandTypeContainer_CALCommandWrite_3Bytes, nil, *requestContext)
		directCommand := readWriteModel.NewRequestDirectCommandAccess(calData, 0x40, nil, nil, 0x0, readWriteModel.NewRequestTermination(), *cbusOptions)
		if err := c.messageCodec.Send(readWriteModel.NewCBusMessageToServer(directCommand, *requestContext, *cbusOptions)); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error writing reset"), ch)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	{
		log.Debug().Msg("Set interface options 3")
		interfaceOptions3 := readWriteModel.NewParameterValueInterfaceOptions3(readWriteModel.NewInterfaceOptions3(true, false, true, false), 1)
		calData := readWriteModel.NewCALDataWrite(readWriteModel.Parameter_INTERFACE_OPTIONS_3, 0x0, interfaceOptions3, readWriteModel.CALCommandTypeContainer_CALCommandWrite_3Bytes, nil, *requestContext)
		directCommand := readWriteModel.NewRequestDirectCommandAccess(calData, 0x40, nil, nil, 0x0, readWriteModel.NewRequestTermination(), *cbusOptions)
		if err := c.messageCodec.Send(readWriteModel.NewCBusMessageToServer(directCommand, *requestContext, *cbusOptions)); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error writing reset"), ch)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	{
		log.Debug().Msg("Set interface options 1 power up settings")
		interfaceOptions1PowerUpSettings := readWriteModel.NewParameterValueInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1PowerUpSettings(readWriteModel.NewInterfaceOptions1(true, true, true, true, false, true)), 1)
		calData := readWriteModel.NewCALDataWrite(readWriteModel.Parameter_INTERFACE_OPTIONS_1_POWER_UP_SETTINGS, 0x0, interfaceOptions1PowerUpSettings, readWriteModel.CALCommandTypeContainer_CALCommandWrite_3Bytes, nil, *requestContext)
		directCommand := readWriteModel.NewRequestDirectCommandAccess(calData, 0x40, nil, nil, 0x0, readWriteModel.NewRequestTermination(), *cbusOptions)
		if err := c.messageCodec.Send(readWriteModel.NewCBusMessageToServer(directCommand, *requestContext, *cbusOptions)); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error writing reset"), ch)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	{
		log.Debug().Msg("Set interface options 1")
		interfaceOptions1 := readWriteModel.NewParameterValueInterfaceOptions1(readWriteModel.NewInterfaceOptions1(true, true, true, true, false, true), 1)
		calData := readWriteModel.NewCALDataWrite(readWriteModel.Parameter_INTERFACE_OPTIONS_1, 0x0, interfaceOptions1, readWriteModel.CALCommandTypeContainer_CALCommandWrite_3Bytes, nil, *requestContext)
		directCommand := readWriteModel.NewRequestDirectCommandAccess(calData, 0x40, nil, nil, 0x0, readWriteModel.NewRequestTermination(), *cbusOptions)
		if err := c.messageCodec.Send(readWriteModel.NewCBusMessageToServer(directCommand, *requestContext, *cbusOptions)); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error writing reset"), ch)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	c.fireConnected(ch)

	log.Debug().Msg("Starting subscription handler")
	go func() {
		for c.IsConnected() {
			log.Debug().Msg("Handling incoming message")
			for monitoredSal := range c.messageCodec.(*MessageCodec).monitoredSALs {
				for _, subscriber := range c.subscribers {
					subscriber.handleMonitoredSal(monitoredSal)
				}
			}
		}
	}()
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
