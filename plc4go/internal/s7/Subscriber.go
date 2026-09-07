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

package s7

import (
	"context"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// Subscriber implements the three S7 subscription flavours (matching plc4j):
//   - AlarmTag "ALM"          + EVENT  -> alarm indication push (UserData group 0x04)
//   - AlarmTag "QUERY:..."    + EVENT  -> one-shot alarm query, delivered on Register
//   - address tag             + CYCLIC -> cyclic value push (UserData group 0x02)
//
// There is one Subscriber per connection: the alarm push subscription and the cyclic jobs
// are PLC-side state bound to the TCP connection.
type Subscriber struct {
	connection *Connection

	consumers      map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	consumersMutex sync.RWMutex

	alarmHandles            []*SubscriptionHandle
	alarmSubscriptionActive bool
	alarmMutex              sync.Mutex

	cyclicHandlesByJob map[uint8][]*SubscriptionHandle
	cyclicMutex        sync.RWMutex

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Subscriber{
		connection:         connection,
		consumers:          map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{},
		cyclicHandlesByJob: map[uint8][]*SubscriptionHandle{},
		log:                customLogger,
		_options:           _options,
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
		request := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)
		responseCodes := map[string]apiModel.PlcResponseCode{}
		handles := map[string]apiModel.PlcSubscriptionHandle{}
		deliver := func() {
			utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
				subscriptionRequest,
				spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes, handles, append(s._options, options.WithCustomLogger(s.log))...),
				nil,
			))
		}

		if !s.connection.driverContext.UserDataServicesSupported {
			for _, tagName := range request.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_UNSUPPORTED
			}
			deliver()
			return
		}

		var alarmNames, queryNames, cyclicNames []string
		alarmTags := map[string]*AlarmTag{}
		queryTags := map[string]*AlarmTag{}
		cyclicTags := map[string]PlcTag{}
		for _, tagName := range request.GetTagNames() {
			subscriptionType := request.GetType(tagName)
			switch tag := s.resolveTag(request.GetTag(tagName)).(type) {
			case *AlarmTag:
				if subscriptionType != apiModel.SubscriptionEvent {
					responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
					continue
				}
				if tag.GetKind() == AlarmTagQuery {
					queryNames = append(queryNames, tagName)
					queryTags[tagName] = tag
				} else {
					alarmNames = append(alarmNames, tagName)
					alarmTags[tagName] = tag
				}
			case PlcTag:
				if subscriptionType != apiModel.SubscriptionCyclic {
					responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
					continue
				}
				cyclicNames = append(cyclicNames, tagName)
				cyclicTags[tagName] = tag
			default:
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
			}
		}

		s.subscribeAlarms(ctx, alarmNames, alarmTags, responseCodes, handles)
		s.subscribeCyclic(ctx, request, cyclicNames, cyclicTags, responseCodes, handles)
		s.runAlarmQueries(ctx, queryNames, queryTags, responseCodes, handles)
		deliver()
	})
	return result
}

// resolveTag maps a subscription tag back to a driver tag. Tags added by address are already
// driver types; tags wrapped in a DefaultPlcSubscriptionTag don't expose their inner tag, so
// they are re-parsed from the address string.
func (s *Subscriber) resolveTag(tag apiModel.PlcSubscriptionTag) apiModel.PlcTag {
	switch typed := tag.(type) {
	case *AlarmTag:
		return typed
	case PlcTag:
		return typed
	default:
		parsed, err := NewTagHandler(s._options...).ParseTag(tag.GetAddressString())
		if err != nil {
			s.log.Debug().Err(err).Str("address", tag.GetAddressString()).Msg("Unable to re-parse subscription tag")
			return nil
		}
		return parsed
	}
}

func (s *Subscriber) subscribeAlarms(ctx context.Context, tagNames []string, tags map[string]*AlarmTag, responseCodes map[string]apiModel.PlcResponseCode, handles map[string]apiModel.PlcSubscriptionHandle) {
	if len(tagNames) == 0 {
		return
	}
	if err := s.ensureAlarmSubscriptionActive(ctx); err != nil {
		s.log.Debug().Err(err).Msg("Alarm subscription failed")
		for _, tagName := range tagNames {
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		}
		return
	}
	s.alarmMutex.Lock()
	defer s.alarmMutex.Unlock()
	for _, tagName := range tagNames {
		handle := NewSubscriptionHandle(s, tagName, tags[tagName], subscriptionKindAlarm)
		s.alarmHandles = append(s.alarmHandles, handle)
		handles[tagName] = handle
		responseCodes[tagName] = apiModel.PlcResponseCode_OK
	}
}

func (s *Subscriber) ensureAlarmSubscriptionActive(ctx context.Context) error {
	s.alarmMutex.Lock()
	alreadyActive := s.alarmSubscriptionActive
	s.alarmMutex.Unlock()
	if alreadyActive {
		return nil
	}
	tpduId := s.connection.tpduGenerator.getAndIncrement()
	alarmState := alarmStateForController(s.connection.driverContext.ControllerType, false)
	response, err := s.connection.sendUserData(ctx, tpduId, buildMsgSubscriptionRequest(tpduId, alarmState), "alarm_subscribe")
	if err != nil {
		return errors.Wrap(err, "error sending alarm subscription request")
	}
	accepted, err := parseMsgSubscriptionResponse(response)
	if err != nil {
		return err
	}
	if !accepted {
		return errors.New("alarm subscription rejected by PLC")
	}
	s.alarmMutex.Lock()
	s.alarmSubscriptionActive = true
	s.alarmMutex.Unlock()
	return nil
}

func (s *Subscriber) subscribeCyclic(ctx context.Context, request *spiModel.DefaultPlcSubscriptionRequest, tagNames []string, tags map[string]PlcTag, responseCodes map[string]apiModel.PlcResponseCode, handles map[string]apiModel.PlcSubscriptionHandle) {
	if len(tagNames) == 0 {
		return
	}
	// Group tags by their effective cadence so callers requesting the same interval ride
	// one PLC job instead of one per tag. Iterate names in request order for stable indices.
	groups := map[cyclicInterval][]string{}
	var groupOrder []cyclicInterval
	for _, tagName := range tagNames {
		interval := pickCyclicInterval(request.GetInterval(tagName))
		if _, ok := groups[interval]; !ok {
			groupOrder = append(groupOrder, interval)
		}
		groups[interval] = append(groups[interval], tagName)
	}
	for _, interval := range groupOrder {
		groupNames := groups[interval]
		groupTags := make([]PlcTag, len(groupNames))
		for i, tagName := range groupNames {
			groupTags[i] = tags[tagName]
		}
		tpduId := s.connection.tpduGenerator.getAndIncrement()
		response, err := s.connection.sendUserData(ctx, tpduId, buildCyclicSubscribeRequest(tpduId, groupTags, interval), "cyclic_subscribe")
		var jobId uint8
		if err == nil {
			jobId, err = parseCyclicSubscribeResponse(response)
		}
		if err != nil {
			s.log.Debug().Err(err).Msg("Cyclic subscribe failed")
			for _, tagName := range groupNames {
				responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
			}
			continue
		}
		jobHandles := make([]*SubscriptionHandle, 0, len(groupNames))
		for i, tagName := range groupNames {
			handle := NewSubscriptionHandle(s, tagName, groupTags[i], subscriptionKindCyclic)
			handle.jobId = jobId
			handle.itemIndex = i
			jobHandles = append(jobHandles, handle)
			handles[tagName] = handle
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		}
		s.cyclicMutex.Lock()
		s.cyclicHandlesByJob[jobId] = jobHandles
		s.cyclicMutex.Unlock()
		s.log.Debug().
			Uint8("jobId", jobId).
			Strs("tags", groupNames).
			Dur("interval", interval.toDuration()).
			Msg("Cyclic subscription active")
	}
}

func (s *Subscriber) runAlarmQueries(ctx context.Context, tagNames []string, tags map[string]*AlarmTag, responseCodes map[string]apiModel.PlcResponseCode, handles map[string]apiModel.PlcSubscriptionHandle) {
	for _, tagName := range tagNames {
		tpduId := s.connection.tpduGenerator.getAndIncrement()
		response, err := s.connection.sendUserData(ctx, tpduId, buildAlarmQueryRequest(tpduId, tags[tagName].GetQueryType()), "alarm_query")
		var payload []byte
		if err == nil {
			payload, err = parseAlarmQueryResponse(response)
		}
		if err != nil {
			s.log.Debug().Err(err).Str("tagName", tagName).Msg("Alarm query failed")
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
			continue
		}
		handle := NewSubscriptionHandle(s, tagName, tags[tagName], subscriptionKindQuery)
		handle.queryPayload = payload
		handles[tagName] = handle
		responseCodes[tagName] = apiModel.PlcResponseCode_OK
	}
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
		err := s.removeHandles(ctx, request.GetSubscriptionHandles())
		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(
			unsubscriptionRequest,
			spiModel.NewDefaultPlcUnsubscriptionResponse(unsubscriptionRequest),
			err,
		))
	})
	return result
}

// removeHandles drops the given handles from the local bookkeeping and cancels the PLC-side
// subscriptions that no longer have any handle attached.
func (s *Subscriber) removeHandles(ctx context.Context, subscriptionHandles []apiModel.PlcSubscriptionHandle) error {
	jobsToCancel := map[uint8]struct{}{}
	cancelAlarms := false
	for _, subscriptionHandle := range subscriptionHandles {
		handle, ok := subscriptionHandle.(*SubscriptionHandle)
		if !ok {
			continue
		}
		switch handle.kind {
		case subscriptionKindAlarm:
			s.alarmMutex.Lock()
			s.alarmHandles = slices.DeleteFunc(s.alarmHandles, func(other *SubscriptionHandle) bool { return other == handle })
			if len(s.alarmHandles) == 0 && s.alarmSubscriptionActive {
				cancelAlarms = true
			}
			s.alarmMutex.Unlock()
		case subscriptionKindCyclic:
			s.cyclicMutex.Lock()
			remaining := slices.DeleteFunc(s.cyclicHandlesByJob[handle.jobId], func(other *SubscriptionHandle) bool { return other == handle })
			if len(remaining) == 0 {
				delete(s.cyclicHandlesByJob, handle.jobId)
				jobsToCancel[handle.jobId] = struct{}{}
			} else {
				s.cyclicHandlesByJob[handle.jobId] = remaining
			}
			s.cyclicMutex.Unlock()
		}
	}
	var firstErr error
	for jobId := range jobsToCancel {
		tpduId := s.connection.tpduGenerator.getAndIncrement()
		if _, err := s.connection.sendUserData(ctx, tpduId, buildCyclicUnsubscribeRequest(tpduId, jobId), "cyclic_unsubscribe"); err != nil {
			s.log.Debug().Err(err).Uint8("jobId", jobId).Msg("Cyclic cancel failed (ignoring)")
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if cancelAlarms {
		tpduId := s.connection.tpduGenerator.getAndIncrement()
		alarmState := alarmStateForController(s.connection.driverContext.ControllerType, true)
		if _, err := s.connection.sendUserData(ctx, tpduId, buildMsgSubscriptionRequest(tpduId, alarmState), "alarm_unsubscribe"); err != nil {
			s.log.Debug().Err(err).Msg("Alarm cancel failed (ignoring)")
			if firstErr == nil {
				firstErr = err
			}
		}
		s.alarmMutex.Lock()
		s.alarmSubscriptionActive = false
		s.alarmMutex.Unlock()
	}
	return firstErr
}

// drain cancels every active PLC-side subscription; used on connection close.
func (s *Subscriber) drain(ctx context.Context) {
	s.alarmMutex.Lock()
	handles := make([]apiModel.PlcSubscriptionHandle, 0, len(s.alarmHandles))
	for _, handle := range s.alarmHandles {
		handles = append(handles, handle)
	}
	s.alarmMutex.Unlock()
	s.cyclicMutex.Lock()
	for _, jobHandles := range s.cyclicHandlesByJob {
		for _, handle := range jobHandles {
			handles = append(handles, handle)
		}
	}
	s.cyclicMutex.Unlock()
	if len(handles) == 0 {
		return
	}
	if err := s.removeHandles(ctx, handles); err != nil {
		s.log.Debug().Err(err).Msg("Error draining subscriptions on close")
	}
}

func (s *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	s.consumersMutex.Lock()
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(s, consumer, handles...)
	s.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	s.consumersMutex.Unlock()

	// Query handles carry their result inline (fetched at subscribe time). Deliver it as a
	// single event now; the handle goes inert thereafter.
	tags := map[string]apiModel.PlcTag{}
	types := map[string]apiModel.PlcSubscriptionType{}
	intervals := map[string]time.Duration{}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	addresses := map[string]string{}
	values := map[string]apiValues.PlcValue{}
	for _, subscriptionHandle := range handles {
		handle, ok := subscriptionHandle.(*SubscriptionHandle)
		if !ok || handle.kind != subscriptionKindQuery {
			continue
		}
		tags[handle.tagName] = handle.tag
		types[handle.tagName] = apiModel.SubscriptionEvent
		intervals[handle.tagName] = 0
		responseCodes[handle.tagName] = apiModel.PlcResponseCode_OK
		addresses[handle.tagName] = handle.tag.GetAddressString()
		values[handle.tagName] = spiValues.NewPlcRawByteArray(handle.queryPayload)
	}
	if len(values) > 0 {
		event := NewSubscriptionEvent(tags, types, intervals, responseCodes, addresses, values, append(s._options, options.WithCustomLogger(s.log))...)
		consumer(&event)
	}
	return consumerRegistration
}

func (s *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	delete(s.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}

// handleUserDataPush routes an unsolicited UserData push to the matching subscriptions.
// Returns true if the message was consumed.
func (s *Subscriber) handleUserDataPush(message readWriteModel.S7MessageUserData) bool {
	group, functionType, subfunction, ok := userDataPushKey(message)
	if !ok || functionType != 0x00 {
		return false
	}
	switch {
	case group == 0x02 && (subfunction == 0x01 || subfunction == 0x05):
		return s.handleCyclicPush(message)
	case group == 0x04:
		if _, isAlarmIndication := alarmIndicationSubfunctions[subfunction]; isAlarmIndication {
			return s.handleAlarmPush(message)
		}
	}
	return false
}

func (s *Subscriber) handleCyclicPush(message readWriteModel.S7MessageUserData) bool {
	jobId, ok := userDataSequenceNumber(message)
	if !ok {
		return false
	}
	s.cyclicMutex.RLock()
	jobHandles := slices.Clone(s.cyclicHandlesByJob[jobId])
	s.cyclicMutex.RUnlock()
	if len(jobHandles) == 0 {
		return false
	}
	items := extractCyclicPushItems(message)
	if items == nil {
		return false
	}
	tags := map[string]apiModel.PlcTag{}
	types := map[string]apiModel.PlcSubscriptionType{}
	intervals := map[string]time.Duration{}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	addresses := map[string]string{}
	values := map[string]apiValues.PlcValue{}
	for _, handle := range jobHandles {
		tags[handle.tagName] = handle.tag
		types[handle.tagName] = apiModel.SubscriptionCyclic
		intervals[handle.tagName] = 0
		addresses[handle.tagName] = handle.tag.GetAddressString()
		if handle.itemIndex >= len(items) {
			responseCodes[handle.tagName] = apiModel.PlcResponseCode_NOT_FOUND
			continue
		}
		value, err := parsePlcValue(context.Background(), handle.tag.(PlcTag), items[handle.itemIndex], s.connection.driverContext.ControllerType)
		if err != nil {
			s.log.Debug().Err(err).Str("tagName", handle.tagName).Msg("Cyclic push decode failed")
			responseCodes[handle.tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		responseCodes[handle.tagName] = apiModel.PlcResponseCode_OK
		values[handle.tagName] = value
	}
	event := NewSubscriptionEvent(tags, types, intervals, responseCodes, addresses, values, append(s._options, options.WithCustomLogger(s.log))...)

	s.consumersMutex.RLock()
	defer s.consumersMutex.RUnlock()
	for registration, consumer := range s.consumers {
		for _, registeredHandle := range registration.GetSubscriptionHandles() {
			if handle, isS7Handle := registeredHandle.(*SubscriptionHandle); isS7Handle && handle.kind == subscriptionKindCyclic && handle.jobId == jobId {
				// One event per registration, even if multiple handles match.
				consumer(&event)
				break
			}
		}
	}
	return true
}

func (s *Subscriber) handleAlarmPush(message readWriteModel.S7MessageUserData) bool {
	payload := parseAlarmIndication(message)
	if payload == nil {
		return false
	}
	s.alarmMutex.Lock()
	anyHandles := len(s.alarmHandles) > 0
	s.alarmMutex.Unlock()
	if !anyHandles {
		return false
	}
	s.consumersMutex.RLock()
	defer s.consumersMutex.RUnlock()
	for registration, consumer := range s.consumers {
		for _, registeredHandle := range registration.GetSubscriptionHandles() {
			handle, isS7Handle := registeredHandle.(*SubscriptionHandle)
			if !isS7Handle || handle.kind != subscriptionKindAlarm {
				continue
			}
			event := NewSubscriptionEvent(
				map[string]apiModel.PlcTag{handle.tagName: handle.tag},
				map[string]apiModel.PlcSubscriptionType{handle.tagName: apiModel.SubscriptionEvent},
				map[string]time.Duration{handle.tagName: 0},
				map[string]apiModel.PlcResponseCode{handle.tagName: apiModel.PlcResponseCode_OK},
				map[string]string{handle.tagName: handle.tag.GetAddressString()},
				map[string]apiValues.PlcValue{handle.tagName: payload},
				append(s._options, options.WithCustomLogger(s.log))...,
			)
			consumer(&event)
		}
	}
	return true
}
