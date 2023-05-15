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
	"context"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	plc4goModel "github.com/apache/plc4x/plc4go/spi/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type Subscriber struct {
	connection *Connection
	consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
}

func NewSubscriber(connection *Connection) *Subscriber {
	return &Subscriber{
		connection: connection,
		consumers:  make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),
	}
}

func (m *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go func() {
		internalPlcSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

		// Add this subscriber to the connection.
		m.connection.addSubscriber(m)

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
		for _, tagName := range internalPlcSubscriptionRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			subscriptionValues[tagName] = spiModel.NewDefaultPlcSubscriptionHandle(m)
		}

		result <- &plc4goModel.DefaultPlcSubscriptionRequestResult{
			Request:  subscriptionRequest,
			Response: spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes, subscriptionValues),
			Err:      nil,
		}
	}()
	return result
}

func (m *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle ctx
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("not implemented"))

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more an internal handling of which values to pass where.

	return result
}

func (m *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(m, consumer, handles...)
	m.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (m *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	delete(m.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
