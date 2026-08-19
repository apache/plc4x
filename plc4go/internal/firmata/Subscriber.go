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

package firmata

import (
	"context"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// Subscriber turns the stream of pin updates a firmata board pushes into subscription events.
// Ported from the subscription half of plc4j's FirmataConnection.
//
// Subscribing is the one thing which does go out on the wire: a pin has to be put into input mode
// and reporting has to be switched on for it before the board says anything about it. Everything
// after that is local filtering - the board reports a pin to everybody or to nobody.
//
//go:generate go tool plc4xGenerator -type=Subscriber
type Subscriber struct {
	connection *Connection

	consumers map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer `hasLocker:"consumersMutex"`
	// handles are the subscription handles which are still active. Firmata has no wire-level
	// unsubscribe worth sending, so a handle is a local filter and unsubscribing means dropping it.
	handles map[*SubscriptionHandle]struct{} `ignore:"true"`
	// consumersMutex guards consumers and handles, which Subscribe, Unsubscribe, Register and
	// Unregister write while the connection's incoming-message worker reads them.
	consumersMutex sync.RWMutex

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Subscribe and deliver both hand their own logger option downstream with append, from
	// different goroutines - clipping makes each of those appends allocate instead of writing into
	// a backing array they share.
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
		for _, tagName := range request.GetTagNames() {
			handle, responseCode := s.subscribeTag(ctx, request, tagName)
			responseCodes[tagName] = responseCode
			if handle != nil {
				subscriptionHandles[tagName] = handle
			}
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

// subscribeTag switches reporting on for one tag and hands back its handle.
func (s *Subscriber) subscribeTag(ctx context.Context, request *spiModel.DefaultPlcSubscriptionRequest, tagName string) (apiModel.PlcSubscriptionHandle, apiModel.PlcResponseCode) {
	tagType := request.GetType(tagName)
	if tagType != apiModel.SubscriptionChangeOfState {
		// A firmata board reports a pin when it changes and nothing else, so change-of-state is the
		// only subscription type this driver can serve.
		//
		// Cyclic: the global sampling interval a board does have is a property of the board, not of
		// a single tag, so there is no honest way to answer a request for one particular interval.
		// Event: the board emits no discrete events at all - and since only changes are published
		// (plc4j's FirmataConnection dedupes the same way), an event subscription would silently be
		// served as a change-of-state one.
		//
		// plc4j accepts both and quietly turns them into change-of-state subscriptions; saying
		// UNSUPPORTED is the honest answer, and it is what a caller can act on.
		s.log.Debug().
			Str("tagName", tagName).
			Stringer("tagType", tagType).
			Msg("Firmata only has change-of-state subscriptions")
		return nil, apiModel.PlcResponseCode_UNSUPPORTED
	}

	tag := s.resolveTag(request.GetTag(tagName))
	if tag == nil {
		s.log.Debug().Str("tagName", tagName).Msg("Not a firmata tag")
		return nil, apiModel.PlcResponseCode_INVALID_ADDRESS
	}

	messages, claim, responseCode := s.subscribeMessagesFor(tag)
	if responseCode != apiModel.PlcResponseCode_OK {
		return nil, responseCode
	}
	if err := s.connection.sendAll(ctx, "subscribe", messages); err != nil {
		// Reporting is only really switched on once these messages are out, so a claim whose
		// messages didn't make it has to be given back - otherwise a retry would find the pins
		// claimed, send nothing and hand out a subscription which never delivers.
		claim.rollback()
		s.log.Debug().Err(err).Str("tagName", tagName).Msg("Error sending the subscribe messages")
		return nil, apiModel.PlcResponseCode_INTERNAL_ERROR
	}

	// Only start being fed once there is something to deliver.
	s.connection.addSubscriber(s)
	handle := NewSubscriptionHandle(s, tagName, tag, tagType, request.GetInterval(tagName))
	s.consumersMutex.Lock()
	s.handles[handle] = struct{}{}
	s.consumersMutex.Unlock()
	return handle, apiModel.PlcResponseCode_OK
}

// subscribeMessagesFor claims the pins of a tag and reports the messages which switch reporting on
// for them.
func (s *Subscriber) subscribeMessagesFor(tag Tag) ([]readWriteModel.FirmataMessage, pinClaim, apiModel.PlcResponseCode) {
	switch typed := tag.(type) {
	case digitalTag:
		mode := readWriteModel.PinMode_PinModeInput
		if typed.pinMode != nil {
			mode = *typed.pinMode
		}
		if mode != readWriteModel.PinMode_PinModeInput && mode != readWriteModel.PinMode_PinModePullup {
			// The address syntax can only spell PULLUP, so this is unreachable through parsing and
			// only guards a hand-built tag.
			s.log.Debug().Stringer("mode", mode).Msg("A subscribed digital pin must be an input or a pullup input")
			return nil, nil, apiModel.PlcResponseCode_INVALID_ADDRESS
		}
		messages, claim, err := s.connection.claimDigitalInputPins(typed.address, typed.quantity, mode)
		if err != nil {
			s.log.Debug().Err(err).Msg("Unable to subscribe to these digital pins")
			return nil, nil, apiModel.PlcResponseCode_INVALID_ADDRESS
		}
		return messages, claim, apiModel.PlcResponseCode_OK
	case analogTag:
		messages, claim, err := s.connection.claimAnalogInputPins(typed.address, typed.quantity)
		if err != nil {
			s.log.Debug().Err(err).Msg("Unable to subscribe to these analog pins")
			return nil, nil, apiModel.PlcResponseCode_INVALID_ADDRESS
		}
		return messages, claim, apiModel.PlcResponseCode_OK
	default:
		return nil, nil, apiModel.PlcResponseCode_INVALID_ADDRESS
	}
}

// resolveTag maps a subscription tag back to a driver tag. A tag added by address already is one; a
// tag wrapped in a DefaultPlcSubscriptionTag doesn't expose its inner tag, so it is re-parsed from
// its address string.
func (s *Subscriber) resolveTag(tag apiModel.PlcSubscriptionTag) Tag {
	if tag == nil {
		return nil
	}
	if firmataTag, ok := tag.(Tag); ok {
		return firmataTag
	}
	parsed, err := NewTagHandler(s._options...).ParseTag(tag.GetAddressString())
	if err != nil {
		s.log.Debug().Err(err).Str("address", tag.GetAddressString()).Msg("Unable to re-parse subscription tag")
		return nil
	}
	firmataTag, _ := parsed.(Tag)
	return firmataTag
}

// Unsubscribe drops the handles of the given request. Reporting stays switched on at the board:
// firmata has no per-tag unsubscribe (reporting is per port for digital pins) and another
// subscription may well still want the same pins, so the pins keep their mode and the updates are
// simply no longer delivered - which is what plc4j's onUnsubscribe amounts to as well.
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

// unsubscribeHandles drops the given handles and every consumer registration which is left without
// a single active handle. Handles which don't belong to this subscriber are silently ignored, as an
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

// activeRegistration is a snapshot of one consumer registration reduced to the handles which are
// still subscribed.
type activeRegistration struct {
	consumer apiModel.PlcSubscriptionEventConsumer
	handles  []*SubscriptionHandle
}

// activeRegistrations snapshots the current registrations so events can be delivered without
// holding the lock while calling into the (possibly re-entrant) consumers.
func (s *Subscriber) activeRegistrations() []activeRegistration {
	s.consumersMutex.RLock()
	defer s.consumersMutex.RUnlock()
	registrations := make([]activeRegistration, 0, len(s.consumers))
	for registration, consumer := range s.consumers {
		var handles []*SubscriptionHandle
		for _, handle := range registration.GetSubscriptionHandles() {
			subscriptionHandle, ok := handle.(*SubscriptionHandle)
			if !ok {
				continue
			}
			if _, active := s.handles[subscriptionHandle]; !active {
				continue
			}
			handles = append(handles, subscriptionHandle)
		}
		if len(handles) == 0 {
			continue
		}
		registrations = append(registrations, activeRegistration{consumer: consumer, handles: handles})
	}
	return registrations
}

// handleAnalogUpdate delivers a changed analog pin to every handle whose run of pins covers it. A
// handle covering several pins gets all of their current values, not just the one which changed,
// the way plc4j's publishAnalogEvents does.
func (s *Subscriber) handleAnalogUpdate(pin uint8) {
	for _, registration := range s.activeRegistrations() {
		for _, handle := range registration.handles {
			tag, ok := handle.tag.(analogTag)
			if !ok || !covers(tag.address, tag.quantity, pin) {
				continue
			}
			plcValues := make([]apiValues.PlcValue, 0, tag.quantity)
			for offset := uint8(0); offset < tag.quantity; offset++ {
				plcValues = append(plcValues, spiValues.NewPlcINT(s.connection.analogValue(tag.address+offset)))
			}
			s.deliver(registration.consumer, handle, plcValues)
		}
	}
}

// handleDigitalUpdate delivers a batch of changed digital pins to every handle whose run of pins
// covers at least one of them.
func (s *Subscriber) handleDigitalUpdate(changedPins []uint8) {
	for _, registration := range s.activeRegistrations() {
		for _, handle := range registration.handles {
			tag, ok := handle.tag.(digitalTag)
			if !ok || !coversAny(tag.address, tag.quantity, changedPins) {
				continue
			}
			plcValues := make([]apiValues.PlcValue, 0, tag.quantity)
			for offset := uint8(0); offset < tag.quantity; offset++ {
				plcValues = append(plcValues, spiValues.NewPlcBOOL(s.connection.digitalValue(tag.address+offset)))
			}
			s.deliver(registration.consumer, handle, plcValues)
		}
	}
}

// deliver assembles and hands out one event for one handle.
func (s *Subscriber) deliver(consumer apiModel.PlcSubscriptionEventConsumer, handle *SubscriptionHandle, plcValues []apiValues.PlcValue) {
	tagName := handle.tagName
	var value apiValues.PlcValue
	if len(plcValues) == 1 {
		value = plcValues[0]
	} else {
		value = spiValues.NewPlcList(plcValues)
	}
	event := NewSubscriptionEvent(
		map[string]apiModel.PlcTag{tagName: handle.tag},
		map[string]apiModel.PlcSubscriptionType{tagName: handle.tagType},
		map[string]time.Duration{tagName: handle.interval},
		map[string]apiModel.PlcResponseCode{tagName: apiModel.PlcResponseCode_OK},
		map[string]string{tagName: handle.tag.GetAddressString()},
		map[string]apiValues.PlcValue{tagName: value},
		append(s._options, options.WithCustomLogger(s.log))...,
	)
	consumer(&event)
}

// covers says whether a run of pins starting at address contains pin.
func covers(address uint8, quantity uint8, pin uint8) bool {
	return int(pin) >= int(address) && int(pin) < int(address)+int(quantity)
}

// coversAny says whether a run of pins starting at address contains any of the given pins.
func coversAny(address uint8, quantity uint8, pins []uint8) bool {
	for _, pin := range pins {
		if covers(address, quantity, pin) {
			return true
		}
	}
	return false
}
