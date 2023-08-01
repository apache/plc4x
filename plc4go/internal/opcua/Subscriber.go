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
	"encoding/binary"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=Subscriber
type Subscriber struct {
	consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	addSubscriber func(subscriber *Subscriber)
	messageCodec  *MessageCodec
	subscriptions map[uint32]*SubscriptionHandle `ignore:"true"` // TODO: we don't have support for non string key maps yet

	consumersMutex sync.RWMutex

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

func NewSubscriber(addSubscriber func(subscriber *Subscriber), messageCodec *MessageCodec, _options ...options.WithOption) *Subscriber {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Subscriber{
		consumers:     make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),
		addSubscriber: addSubscriber,
		messageCodec:  messageCodec,
		subscriptions: map[uint32]*SubscriptionHandle{},

		log:      customLogger,
		_options: _options,
	}
}

func (s *Subscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go s.subscribeSync(result, subscriptionRequest)
	return result
}

func (s *Subscriber) subscribeSync(result chan apiModel.PlcSubscriptionRequestResult, subscriptionRequest apiModel.PlcSubscriptionRequest) {
	defer func() {
		if err := recover(); err != nil {
			result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
		}
	}()
	internalPlcSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

	cycleTime := subscriptionRequest.GetTag(subscriptionRequest.GetTagNames()[0]).GetDuration()
	if cycleTime == 0 {
		cycleTime = 1 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_TIMEOUT)
	defer cancel()
	subscription, err := s.onSubscribeCreateSubscription(ctx, cycleTime)
	if err != nil {
		result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Wrap(err, "error create subscription"))
		return
	}
	subscriptionId := subscription.GetSubscriptionId()
	handle := NewSubscriptionHandle(s.log, s, s.messageCodec, subscriptionRequest, subscriptionId, cycleTime)
	s.subscriptions[subscriptionId] = handle

	// Add this subscriber to the connection.
	s.addSubscriber(s)

	// Just populate all requests with an OK
	responseCodes := map[string]apiModel.PlcResponseCode{}
	subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
	for _, tagName := range internalPlcSubscriptionRequest.GetTagNames() {
		responseCodes[tagName] = apiModel.PlcResponseCode_OK
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
}

func (s *Subscriber) onSubscribeCreateSubscription(ctx context.Context, cycleTime time.Duration) (readWriteModel.CreateSubscriptionResponse, error) {
	s.log.Trace().Msg("Entering creating subscription request")

	channel := s.messageCodec.channel

	requestHeader := readWriteModel.NewRequestHeader(channel.getAuthenticationToken(),
		channel.getCurrentDateTime(),
		channel.getRequestHandle(),
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	createSubscriptionRequest := readWriteModel.NewCreateSubscriptionRequest(
		requestHeader,
		float64(cycleTime),
		12000,
		5,
		65536,
		true,
		0,
	)

	identifier, err := strconv.ParseUint(createSubscriptionRequest.GetIdentifier(), 10, 16)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing identifier")
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		createSubscriptionRequest,
		false)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		return nil, errors.Wrap(err, "error serializing")
	}

	responseChan := make(chan readWriteModel.CreateSubscriptionResponse, 100) // TODO: bit oversized to not block anything. Discards errors
	errorChan := make(chan error, 100)                                        // TODO: bit oversized to not block anything. Discards errors
	/* Functional Consumer example using inner class */
	consumer := func(opcuaResponse []byte) {
		extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			errorChan <- errors.Wrap(err, "error Parsing")
			return
		}
		responseMessage := extensionObject.GetBody().(readWriteModel.CreateSubscriptionResponse)

		// Pass the response back to the application.
		responseChan <- responseMessage
	}

	errorDispatcher := func(err error) {
		errorChan <- errors.Wrap(err, "error received")
	}
	channel.submit(ctx, s.messageCodec, errorDispatcher, consumer, buffer)

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errorChan:
		return nil, errors.Wrap(err, "error received")
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "context ended")
	}
}

func (s *Subscriber) onDisconnect(codec *MessageCodec) {
	s.log.Trace().Msg("disconnecting")
	for _, handle := range s.subscriptions {
		handle.stopSubscriber()
	}
	codec.channel.onDisconnect(context.Background(), codec)
}

func (s *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("Not Implemented"))

	for _, handle := range unsubscriptionRequest.(*spiModel.DefaultPlcUnsubscriptionRequest).GetSubscriptionHandles() {
		handle.(*SubscriptionHandle).stopSubscriber()
	}

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
