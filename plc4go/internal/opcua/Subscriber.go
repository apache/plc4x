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

package opcua

import (
	"context"
	"runtime/debug"
	"sync"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=Subscriber
type Subscriber struct {
	consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer `ignore:"true"`
	addSubscriber func(subscriber *Subscriber)

	consumersMutex sync.RWMutex

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

func NewSubscriber(addSubscriber func(subscriber *Subscriber), _options ...options.WithOption) *Subscriber {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Subscriber{
		addSubscriber: addSubscriber,
		consumers:     make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),

		log:      customLogger,
		_options: _options,
	}
}

func (s *Subscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		internalPlcSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

		// Add this subscriber to the connection.
		s.addSubscriber(s)

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
		for _, tagName := range internalPlcSubscriptionRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			handle := NewSubscriptionHandle(
				s,
				tagName,
				internalPlcSubscriptionRequest.GetTag(tagName),
				internalPlcSubscriptionRequest.GetType(tagName),
				internalPlcSubscriptionRequest.GetInterval(tagName),
			)
			preRegisteredConsumers := internalPlcSubscriptionRequest.GetPreRegisteredConsumers(tagName)
			for _, consumer := range preRegisteredConsumers {
				_ = handle.Register(consumer)
			}
			subscriptionValues[tagName] = handle
		}

		result <- spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(
				subscriptionRequest,
				responseCodes,
				subscriptionValues,
				append(s._options, options.WithCustomLogger(s.log))...,
			),
			nil,
		)
	}()
	return result
}

func (s *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("Not Implemented"))

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more a internal handling of which values to pass where.
	_ = ctx
	_ = unsubscriptionRequest

	return result
}

func (s *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(s, consumer, handles...)
	s.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (s *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	delete(s.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
