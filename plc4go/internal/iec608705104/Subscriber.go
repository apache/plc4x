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

package iec608705104

import (
	"context"
	"runtime/debug"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Subscriber turns the ASDUs a controlled station pushes into subscription events.
//
// Nothing about subscribing goes out on the wire: after the STARTDT handshake the station reports
// what its own configuration tells it to, and a subscription is a local filter over that stream. What
// makes the filter worth having is the tag syntax - a tag may name one point or, through wildcards,
// a whole station - which is precisely what plc4j never implemented: its Iec608705104TagHandler
// resolves every address to ASDU 0 / IOA 0 and its connection then hands every event to every
// consumer regardless of what they subscribed to.
type Subscriber struct {
	connection *Connection

	// consumersMutex guards consumers and handles, which Subscribe, Unsubscribe, Register and
	// Unregister write while the connection's incoming-message worker reads them.
	consumersMutex sync.RWMutex
	consumers      map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	// handles are the subscription handles which are still active. There is no wire-level
	// unsubscribe in IEC 60870-5-104, so unsubscribing means dropping the handle.
	handles map[*SubscriptionHandle]struct{}

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Subscribe and deliver both hand their own logger option downstream with append, from
	// different goroutines - clipping makes each of those appends allocate instead of writing into a
	// backing array they share.
	_options = slices.Clip(_options)
	return &Subscriber{
		connection: connection,
		consumers:  make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),
		handles:    make(map[*SubscriptionHandle]struct{}),
		log:        customLogger,
		_options:   _options,
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
		anyHandle := false
		for _, tagName := range request.GetTagNames() {
			handle, responseCode := s.subscribeTag(request, tagName)
			responseCodes[tagName] = responseCode
			if handle != nil {
				subscriptionHandles[tagName] = handle
				anyHandle = true
			}
		}
		// Only start being fed once there is something to deliver.
		if anyHandle && s.connection != nil {
			s.connection.addSubscriber(s)
		}

		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(
				subscriptionRequest,
				responseCodes,
				subscriptionHandles,
				append(s._options, options.WithCustomLogger(s.log))...,
			),
			nil,
		))
	})
	return result
}

// subscribeTag hands back the handle for one tag of a request, or the reason it can't have one.
func (s *Subscriber) subscribeTag(request *spiModel.DefaultPlcSubscriptionRequest, tagName string) (apiModel.PlcSubscriptionHandle, apiModel.PlcResponseCode) {
	tagType := request.GetType(tagName)
	if tagType == apiModel.SubscriptionCyclic {
		// How often a point is reported is a property of the station's own configuration: the
		// controlling station can ask for everything once (a general interrogation) but it cannot ask
		// for a particular point every n milliseconds. Answering a cyclic request with whatever the
		// station happens to send would be a subscription whose interval is silently ignored, so it
		// is refused instead.
		s.log.Debug().
			Str("tagName", tagName).
			Stringer("tagType", tagType).
			Msg("IEC 60870-5-104 has no per-tag reporting interval")
		return nil, apiModel.PlcResponseCode_UNSUPPORTED
	}

	tag, ok := s.resolveTag(request.GetTag(tagName))
	if !ok {
		s.log.Debug().Str("tagName", tagName).Msg("Not an IEC 60870-5-104 tag")
		return nil, apiModel.PlcResponseCode_INVALID_ADDRESS
	}

	handle := NewSubscriptionHandle(s, tagName, tag, tagType, request.GetInterval(tagName))
	s.consumersMutex.Lock()
	s.handles[handle] = struct{}{}
	s.consumersMutex.Unlock()
	return handle, apiModel.PlcResponseCode_OK
}

// resolveTag maps a subscription tag back to a driver tag. A tag added by address already is one; a
// tag wrapped in a DefaultPlcSubscriptionTag doesn't expose its inner tag, so it is re-parsed from
// its address string - which only works because Tag.GetAddressString really spells the address (the
// plc4j tag returns null there).
func (s *Subscriber) resolveTag(tag apiModel.PlcSubscriptionTag) (Tag, bool) {
	if tag == nil {
		return Tag{}, false
	}
	if iecTag, ok := tag.(Tag); ok {
		return iecTag, true
	}
	parsed, err := NewTagHandler(s._options...).ParseTag(tag.GetAddressString())
	if err != nil {
		s.log.Debug().Err(err).Str("address", tag.GetAddressString()).Msg("Unable to re-parse subscription tag")
		return Tag{}, false
	}
	iecTag, ok := parsed.(Tag)
	return iecTag, ok
}

// Unsubscribe drops the handles of the given request. Nothing goes out on the wire: reporting is a
// property of the session, not of a tag, and STOPDT would silence every other subscription of the
// same connection.
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
		if err := ctx.Err(); err != nil {
			utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, err))
			return
		}
		s.unsubscribeHandles(request.GetSubscriptionHandles())
		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(
			unsubscriptionRequest,
			spiModel.NewDefaultPlcUnsubscriptionResponse(unsubscriptionRequest),
			nil,
		))
	})
	return result
}

// unsubscribeHandles drops the given handles and every consumer registration which is left without a
// single active handle. Handles which don't belong to this subscriber are silently ignored, as an
// unsubscription request is handed to every subscriber it names a handle of.
func (s *Subscriber) unsubscribeHandles(subscriptionHandles []apiModel.PlcSubscriptionHandle) {
	s.consumersMutex.Lock()
	for _, handle := range subscriptionHandles {
		subscriptionHandle, ok := handle.(*SubscriptionHandle)
		if !ok {
			continue
		}
		delete(s.handles, subscriptionHandle)
	}
	for registration := range s.consumers {
		registeredHandles := registration.GetSubscriptionHandles()
		if len(registeredHandles) == 0 {
			continue
		}
		stillActive := false
		for _, handle := range registeredHandles {
			subscriptionHandle, ok := handle.(*SubscriptionHandle)
			if !ok {
				continue
			}
			if _, active := s.handles[subscriptionHandle]; active {
				stillActive = true
				break
			}
		}
		if !stillActive {
			delete(s.consumers, registration)
		}
	}
	nothingLeft := len(s.handles) == 0
	s.consumersMutex.Unlock()

	// Once there is nothing left to deliver, stop being fed by the connection.
	if nothingLeft && s.connection != nil {
		s.connection.removeSubscriber(s)
	}
}

func (s *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(s, consumer, handles...)
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	s.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (s *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	defaultRegistration, ok := registration.(*spiModel.DefaultPlcConsumerRegistration)
	if !ok {
		return
	}
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	delete(s.consumers, defaultRegistration)
}

// Close drops everything this subscriber holds and waits for its in-flight requests. It is called
// from the connection's close path; a subscriber which is closed delivers nothing further.
func (s *Subscriber) Close() {
	s.consumersMutex.Lock()
	s.consumers = make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer)
	s.handles = make(map[*SubscriptionHandle]struct{})
	s.consumersMutex.Unlock()
	s.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Event delivery
////////////////////////////////////////////////////////////////////////////////////////////////////

// handleDelivery is one active handle together with every consumer registered on it. Handles are
// collapsed rather than walked per registration, because change detection is per handle: were the
// same handle visited once per registration, the first visit would swallow the report for the rest.
type handleDelivery struct {
	handle    *SubscriptionHandle
	consumers []apiModel.PlcSubscriptionEventConsumer
}

// activeDeliveries snapshots the current registrations so events can be delivered without holding
// the lock while calling into the (possibly re-entrant) consumers.
func (s *Subscriber) activeDeliveries() []handleDelivery {
	s.consumersMutex.RLock()
	defer s.consumersMutex.RUnlock()
	byHandle := make(map[*SubscriptionHandle][]apiModel.PlcSubscriptionEventConsumer, len(s.handles))
	order := make([]*SubscriptionHandle, 0, len(s.handles))
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
	deliveries := make([]handleDelivery, 0, len(order))
	for _, handle := range order {
		deliveries = append(deliveries, handleDelivery{handle: handle, consumers: byHandle[handle]})
	}
	return deliveries
}

// publish hands one incoming ASDU to the handles which cover its points. An ASDU can carry several
// information objects, and each of them is a point of its own.
func (s *Subscriber) publish(asdu readWriteModel.ASDU) {
	deliveries := s.activeDeliveries()
	if len(deliveries) == 0 {
		return
	}
	asduAddress := asdu.GetAsduAddressField()
	for _, informationObject := range asdu.GetInformationObjects() {
		if informationObject == nil {
			continue
		}
		informationObjectAddress := informationObject.GetAddress()

		// Decoding is only worth doing once something is actually listening for this point.
		var value apiValues.PlcValue
		var fingerprint string
		var code apiModel.PlcResponseCode
		decoded := false

		for _, delivery := range deliveries {
			if !delivery.handle.tag.Matches(asduAddress, informationObjectAddress) {
				continue
			}
			if !decoded {
				value, fingerprint, code = decodePoint(asdu, informationObject)
				decoded = true
			}
			if !delivery.handle.shouldPublish(asduAddress, informationObjectAddress, fingerprint) {
				continue
			}
			event := s.buildEvent(delivery.handle, asduAddress, informationObjectAddress, value, code)
			for _, consumer := range delivery.consumers {
				consumer(&event)
			}
		}
	}
}

// buildEvent assembles the event for one point of one handle.
func (s *Subscriber) buildEvent(handle *SubscriptionHandle, asduAddress uint16, informationObjectAddress uint32, value apiValues.PlcValue, code apiModel.PlcResponseCode) SubscriptionEvent {
	tagName := handle.tagName
	return NewSubscriptionEvent(
		map[string]apiModel.PlcTag{tagName: handle.tag},
		map[string]apiModel.PlcSubscriptionType{tagName: handle.tagType},
		map[string]time.Duration{tagName: handle.interval},
		map[string]apiModel.PlcResponseCode{tagName: code},
		map[string]string{tagName: pointAddressString(asduAddress, informationObjectAddress)},
		map[string]apiValues.PlcValue{tagName: value},
		append(s._options, options.WithCustomLogger(s.log))...,
	)
}

// pointAddressString spells the concrete address of a point the way the tag handler parses it back.
func pointAddressString(asduAddress uint16, informationObjectAddress uint32) string {
	return strconv.FormatUint(uint64(asduAddress), 10) + "/" +
		strconv.FormatUint(uint64(informationObjectAddress&maxInformationObjectAddress), 10)
}
