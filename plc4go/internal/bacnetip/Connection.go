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
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/rs/zerolog/log"
)

type Connection struct {
	_default.DefaultConnection
	messageCodec spi.MessageCodec
	subscribers  []*Subscriber
	connectionId string
	tracer       *spi.Tracer
}

func NewConnection(messageCodec spi.MessageCodec, fieldHandler spi.PlcFieldHandler, options map[string][]string) *Connection {
	connection := &Connection{
		messageCodec: messageCodec,
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

func (c *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		connectionConnectResult := <-c.DefaultConnection.Connect()
		go func() {
			for c.IsConnected() {
				log.Debug().Msg("Polling data")
				incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
				select {
				case message := <-incomingMessageChannel:
					// TODO: implement mapping to subscribers
					log.Info().Msgf("Received \n%v", message)
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

func (c *Connection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	return internalModel.NewDefaultPlcSubscriptionRequestBuilder(c.GetPlcFieldHandler(), c.GetPlcValueHandler(), NewSubscriber(c))
}

func (c *Connection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	panic("Not implemented yet. (at least as a default)")
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
