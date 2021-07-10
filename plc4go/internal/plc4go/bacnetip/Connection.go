//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package bacnetip

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/rs/zerolog/log"
)

type Connection struct {
	_default.DefaultConnection
	messageCodec spi.MessageCodec
	subscribers  []*Subscriber
}

func NewConnection(messageCodec spi.MessageCodec, fieldHandler spi.PlcFieldHandler) *Connection {
	connection := &Connection{
		messageCodec: messageCodec,
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcFieldHandler(fieldHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
	)
	return connection
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
	panic("Not implementec yet. (at least as a default)")
}

func (m *Connection) addSubscriber(subscriber *Subscriber) {
	for _, sub := range m.subscribers {
		if sub == subscriber {
			log.Debug().Msgf("Subscriber %v already added", subscriber)
			return
		}
	}
	m.subscribers = append(m.subscribers, subscriber)
}

func (c *Connection) String() string {
	return fmt.Sprintf("bacnetip.Connection")
}
