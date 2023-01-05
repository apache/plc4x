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

package bacnetip

import (
	"fmt"
	"sync"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog/log"
)

type Connection struct {
	_default.DefaultConnection
	invokeIdGenerator InvokeIdGenerator
	messageCodec      spi.MessageCodec
	subscribers       []*Subscriber
	tm                *spi.RequestTransactionManager

	connectionId string
	tracer       *spi.Tracer
}

func NewConnection(messageCodec spi.MessageCodec, tagHandler spi.PlcTagHandler, tm *spi.RequestTransactionManager, options map[string][]string) *Connection {
	connection := &Connection{
		invokeIdGenerator: InvokeIdGenerator{currentInvokeId: 0},
		messageCodec:      messageCodec,
		tm:                tm,
	}
	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcTagHandler(tagHandler),
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

func (c *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		connectionConnectResult := <-c.DefaultConnection.Connect()
		go func() {
			for c.IsConnected() {
				log.Trace().Msg("Polling data")
				incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
				timeout := time.NewTimer(20 * time.Millisecond)
				select {
				case message := <-incomingMessageChannel:
					// TODO: implement mapping to subscribers
					log.Info().Msgf("Received \n%v", message)
					utils.CleanupTimer(timeout)
				case <-timeout.C:
				}
			}
			log.Info().Msg("Ending incoming message transfer")
		}()
		ch <- connectionConnectResult
	}()
	return ch
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(c.GetPlcTagHandler(), NewReader(&c.invokeIdGenerator, c.messageCodec, c.tm))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return internalModel.NewDefaultPlcSubscriptionRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), NewSubscriber(c))
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
	return fmt.Sprintf("bacnetip.Connection")
}

type InvokeIdGenerator struct {
	currentInvokeId uint8
	lock            sync.Mutex
}

func (t *InvokeIdGenerator) getAndIncrement() uint8 {
	t.lock.Lock()
	defer t.lock.Unlock()
	result := t.currentInvokeId
	t.currentInvokeId += 1
	return result
}
