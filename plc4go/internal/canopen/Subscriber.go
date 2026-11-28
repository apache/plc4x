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

package canopen

import (
	"context"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// Subscriber turns PDO, NMT and heartbeat frames into subscription events.
//
// Nothing about subscribing goes on the bus: nodes send their PDOs and heartbeats whether or not
// anybody listens, so a subscription is a local filter. Event subscriptions get every matching
// frame; change-of-state subscriptions only those whose value differs from the last one delivered
// for the same node. Cyclic subscriptions are refused: how often a node transmits is the node's
// configuration (its PDO transmission type and event timer, its heartbeat producer time), not
// something a listener can ask for. plc4j accepts event subscriptions only.
type Subscriber struct {
	connection *Connection

	mutex     sync.RWMutex
	consumers map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	handles   map[*SubscriptionHandle]struct{}

	wg sync.WaitGroup

	log      zerolog.Logger
	_options []options.WithOption
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	return &Subscriber{
		connection: connection,
		consumers:  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{},
		handles:    map[*SubscriptionHandle]struct{}{},
		log:        options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
		_options:   slices.Clip(_options),
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	s.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		request, ok := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)
		if !ok {
			utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("unsupported subscription request type %T", subscriptionRequest)))
			return
		}
		if err := ctx.Err(); err != nil {
			utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, err))
			return
		}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionHandles := map[string]apiModel.PlcSubscriptionHandle{}
		for _, tagName := range request.GetTagNames() {
			tagType := request.GetType(tagName)
			if tagType == apiModel.SubscriptionCyclic {
				responseCodes[tagName] = apiModel.PlcResponseCode_UNSUPPORTED
				continue
			}
			tag, ok := castToTag(request.GetTag(tagName))
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			if tag.Kind == TagKindSDO {
				// An object dictionary entry is only ever read on request.
				responseCodes[tagName] = apiModel.PlcResponseCode_UNSUPPORTED
				continue
			}
			handle := NewSubscriptionHandle(s, tagName, tag, tagType)
			s.mutex.Lock()
			s.handles[handle] = struct{}{}
			s.mutex.Unlock()
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			subscriptionHandles[tagName] = handle
		}
		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes, subscriptionHandles, append(s._options, options.WithCustomLogger(s.log))...),
			nil,
		))
	})
	return result
}

func (s *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	s.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		request, ok := unsubscriptionRequest.(*spiModel.DefaultPlcUnsubscriptionRequest)
		if !ok {
			utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.Errorf("unsupported unsubscription request type %T", unsubscriptionRequest)))
			return
		}
		s.mutex.Lock()
		for _, handle := range request.GetSubscriptionHandles() {
			if subscriptionHandle, ok := handle.(*SubscriptionHandle); ok {
				delete(s.handles, subscriptionHandle)
			}
		}
		s.mutex.Unlock()
		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, spiModel.NewDefaultPlcUnsubscriptionResponse(unsubscriptionRequest), nil))
	})
	return result
}

func (s *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	registration := spiModel.NewDefaultPlcConsumerRegistration(s, consumer, handles...)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.consumers[registration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return registration
}

func (s *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	defaultRegistration, ok := registration.(*spiModel.DefaultPlcConsumerRegistration)
	if !ok {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.consumers, defaultRegistration)
}

// Close drops every subscription and waits for in-flight requests.
func (s *Subscriber) Close() {
	s.mutex.Lock()
	s.consumers = map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{}
	s.handles = map[*SubscriptionHandle]struct{}{}
	s.mutex.Unlock()
	s.wg.Wait()
}

type delivery struct {
	handle    *SubscriptionHandle
	consumers []apiModel.PlcSubscriptionEventConsumer
}

func (s *Subscriber) deliveries() []delivery {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	byHandle := map[*SubscriptionHandle][]apiModel.PlcSubscriptionEventConsumer{}
	var order []*SubscriptionHandle
	for registration, consumer := range s.consumers {
		for _, handle := range registration.GetSubscriptionHandles() {
			subscriptionHandle, ok := handle.(*SubscriptionHandle)
			if !ok {
				continue
			}
			if _, active := s.handles[subscriptionHandle]; !active {
				continue
			}
			if _, seen := byHandle[subscriptionHandle]; !seen {
				order = append(order, subscriptionHandle)
			}
			byHandle[subscriptionHandle] = append(byHandle[subscriptionHandle], consumer)
		}
	}
	result := make([]delivery, 0, len(order))
	for _, handle := range order {
		result = append(result, delivery{handle: handle, consumers: byHandle[handle]})
	}
	return result
}

// publish hands one frame to the handles it matches. It runs on the codec's receive loop, so
// consumers should hand slow work off.
func (s *Subscriber) publish(frame readWriteModel.CANOpenFrame) {
	deliveries := s.deliveries()
	if len(deliveries) == 0 {
		return
	}
	kind, node, ok := classify(frame)
	if !ok {
		return
	}
	ctx := context.Background()
	for _, d := range deliveries {
		if !d.handle.matches(kind, frame.GetService(), node) {
			continue
		}
		code, value := eventValue(ctx, d.handle.tag, frame, node)
		if !d.handle.shouldPublish(node, value) {
			continue
		}
		event := s.buildEvent(d.handle, code, value)
		for _, consumer := range d.consumers {
			consumer(&event)
		}
	}
}

// classify says which tag kind a frame feeds and which node it concerns. For NMT that is the node
// the command is addressed to, since every NMT frame goes out under COB-ID 0 (plc4j matches the COB
// node, so only its wildcard NMT subscriptions ever fire).
func classify(frame readWriteModel.CANOpenFrame) (TagKind, uint8, bool) {
	switch payload := frame.GetPayload().(type) {
	case readWriteModel.CANOpenPDOPayload:
		return TagKindPDO, frame.GetNodeId(), true
	case readWriteModel.CANOpenHeartbeatPayload:
		return TagKindHeartbeat, frame.GetNodeId(), true
	case readWriteModel.CANOpenNetworkPayload:
		return TagKindNMT, payload.GetNode(), true
	}
	return 0, 0, false
}

// eventValue is what an event carries: the decoded PDO value, or {state, node} for NMT and heartbeat
// frames, the shape plc4j uses.
func eventValue(ctx context.Context, tag Tag, frame readWriteModel.CANOpenFrame, node uint8) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	switch payload := frame.GetPayload().(type) {
	case readWriteModel.CANOpenPDOPayload:
		value, err := decodeValue(ctx, tag.DataType, payload.GetPdo().GetData())
		if err != nil {
			return apiModel.PlcResponseCode_INVALID_DATA, spiValues.NewPlcNULL()
		}
		return apiModel.PlcResponseCode_OK, value
	case readWriteModel.CANOpenHeartbeatPayload:
		return apiModel.PlcResponseCode_OK, stateStruct(uint8(payload.GetState()), node)
	case readWriteModel.CANOpenNetworkPayload:
		return apiModel.PlcResponseCode_OK, stateStruct(uint8(payload.GetRequest()), node)
	}
	return apiModel.PlcResponseCode_INTERNAL_ERROR, spiValues.NewPlcNULL()
}

func stateStruct(state uint8, node uint8) apiValues.PlcValue {
	return spiValues.NewPlcStructOrdered(map[string]apiValues.PlcValue{
		"state": spiValues.NewPlcUSINT(state),
		"node":  spiValues.NewPlcUSINT(node),
	}, []string{"state", "node"})
}

func (s *Subscriber) buildEvent(handle *SubscriptionHandle, code apiModel.PlcResponseCode, value apiValues.PlcValue) spiModel.DefaultPlcSubscriptionEvent {
	tagName := handle.tagName
	event := spiModel.NewDefaultPlcSubscriptionEvent(
		nil,
		map[string]apiModel.PlcTag{tagName: handle.tag},
		map[string]apiModel.PlcSubscriptionType{tagName: handle.tagType},
		map[string]time.Duration{tagName: 0},
		map[string]apiModel.PlcResponseCode{tagName: code},
		map[string]apiValues.PlcValue{tagName: value},
		append(s._options, options.WithCustomLogger(s.log))...,
	)
	return *event.(*spiModel.DefaultPlcSubscriptionEvent)
}
