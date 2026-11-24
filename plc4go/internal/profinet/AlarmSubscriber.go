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

package profinet

import (
	"context"
	"fmt"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// AlarmTag subscribes to the device's alarms: ALARM for all of them, ALARM:<slot> for one slot's,
// ALARM:<slot>.<subslot> for one submodule's. Each event's value is a struct (see Alarm.PlcValue).
// Alarms only arrive by subscription, in Event or ChangeOfState mode; they cannot be read.
type AlarmTag struct {
	// Slot and Subslot filter the alarms; -1 is any.
	Slot, Subslot int
}

var alarmTagPattern = regexp.MustCompile(`^ALARM(?::(\d{1,5})(?:\.(\d{1,5}))?)?$`)

func parseAlarmTag(address string) (AlarmTag, bool, error) {
	m := alarmTagPattern.FindStringSubmatch(address)
	if m == nil {
		return AlarmTag{}, false, nil
	}
	tag := AlarmTag{Slot: -1, Subslot: -1}
	for i, part := range []string{m[1], m[2]} {
		if part == "" {
			continue
		}
		n, err := strconv.ParseUint(part, 10, 16)
		if err != nil {
			return AlarmTag{}, true, errors.Errorf("%s is out of range", part)
		}
		if i == 0 {
			tag.Slot = int(n)
		} else {
			tag.Subslot = int(n)
		}
	}
	return tag, true, nil
}

func (t AlarmTag) matches(a Alarm) bool {
	return (t.Slot < 0 || int(a.Slot) == t.Slot) && (t.Subslot < 0 || int(a.Subslot) == t.Subslot)
}

func (t AlarmTag) GetAddressString() string {
	switch {
	case t.Slot < 0:
		return "ALARM"
	case t.Subslot < 0:
		return fmt.Sprintf("ALARM:%d", t.Slot)
	default:
		return fmt.Sprintf("ALARM:%d.%d", t.Slot, t.Subslot)
	}
}

func (t AlarmTag) GetValueType() apiValues.PlcValueType { return apiValues.Struct }
func (t AlarmTag) GetArrayInfo() []apiModel.ArrayInfo   { return []apiModel.ArrayInfo{} }
func (t AlarmTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionEvent
}
func (t AlarmTag) GetDuration() time.Duration { return 0 }
func (t AlarmTag) String() string             { return "profinet.AlarmTag{" + t.GetAddressString() + "}" }

func (t AlarmTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := t.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (t AlarmTag) SerializeWithWriteBuffer(_ context.Context, wb utils.WriteBuffer) error {
	for _, c := range []string{"PlcTagItem", "tag", "ProfinetAlarmTag"} {
		if err := wb.PushContext(c); err != nil {
			return err
		}
	}
	address := t.GetAddressString()
	if err := wb.WriteString("address", uint32(len(address)*8), address, utils.WithEncoding("UTF-8")); err != nil {
		return err
	}
	for _, c := range []string{"ProfinetAlarmTag", "tag", "PlcTagItem"} {
		if err := wb.PopContext(c); err != nil {
			return err
		}
	}
	return nil
}

func castToAlarmTag(tag apiModel.PlcTag) (AlarmTag, bool) {
	if tag == nil {
		return AlarmTag{}, false
	}
	if t, ok := tag.(AlarmTag); ok {
		return t, true
	}
	t, ok, err := parseAlarmTag(tag.GetAddressString())
	return t, ok && err == nil
}

// alarmHandle is one subscribed alarm tag.
type alarmHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle
	tagName string
	tag     AlarmTag
	tagType apiModel.PlcSubscriptionType
}

// alarmSubscriber turns alarms into subscription events. Subscribing sends nothing to the device:
// it sends its alarms whatever happens, and the driver acknowledges them whether or not anybody
// subscribed.
type alarmSubscriber struct {
	mutex     sync.RWMutex
	consumers map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	handles   map[*alarmHandle]struct{}

	wg       sync.WaitGroup
	log      zerolog.Logger
	_options []options.WithOption
}

var _ spi.PlcSubscriber = (*alarmSubscriber)(nil)

func newAlarmSubscriber(_options ...options.WithOption) *alarmSubscriber {
	return &alarmSubscriber{
		consumers: map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{},
		handles:   map[*alarmHandle]struct{}{},
		log:       options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
		_options:  slices.Clip(_options),
	}
}

func (s *alarmSubscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	s.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		if err := ctx.Err(); err != nil {
			utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, err))
			return
		}
		codes := map[string]apiModel.PlcResponseCode{}
		handles := map[string]apiModel.PlcSubscriptionHandle{}
		for _, name := range subscriptionRequest.GetTagNames() {
			tagType := subscriptionRequest.(interface {
				GetType(string) apiModel.PlcSubscriptionType
			}).GetType(name)
			if tagType == apiModel.SubscriptionCyclic {
				// An alarm is an event; there is nothing to sample.
				codes[name] = apiModel.PlcResponseCode_UNSUPPORTED
				continue
			}
			tag, ok := castToAlarmTag(subscriptionRequest.GetTag(name))
			if !ok {
				codes[name] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			handle := &alarmHandle{tagName: name, tag: tag, tagType: tagType}
			handle.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(s, handle)
			s.mutex.Lock()
			s.handles[handle] = struct{}{}
			s.mutex.Unlock()
			codes[name] = apiModel.PlcResponseCode_OK
			handles[name] = handle
		}
		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, codes, handles, append(s._options, options.WithCustomLogger(s.log))...), nil))
	})
	return result
}

func (s *alarmSubscriber) Unsubscribe(_ context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	if request, ok := unsubscriptionRequest.(*spiModel.DefaultPlcUnsubscriptionRequest); ok {
		s.mutex.Lock()
		for _, h := range request.GetSubscriptionHandles() {
			if handle, ok := h.(*alarmHandle); ok {
				delete(s.handles, handle)
			}
		}
		s.mutex.Unlock()
	}
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, spiModel.NewDefaultPlcUnsubscriptionResponse(unsubscriptionRequest), nil)
	return result
}

func (s *alarmSubscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	registration := spiModel.NewDefaultPlcConsumerRegistration(s, consumer, handles...)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.consumers[registration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return registration
}

func (s *alarmSubscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	if r, ok := registration.(*spiModel.DefaultPlcConsumerRegistration); ok {
		s.mutex.Lock()
		delete(s.consumers, r)
		s.mutex.Unlock()
	}
}

func (s *alarmSubscriber) Close() {
	s.mutex.Lock()
	s.consumers = map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{}
	s.handles = map[*alarmHandle]struct{}{}
	s.mutex.Unlock()
	s.wg.Wait()
}

// publish hands an alarm to the subscriptions it matches. It runs on the cyclic receive loop, so
// consumers should hand slow work off.
func (s *alarmSubscriber) publish(alarm Alarm) {
	type delivery struct {
		handle    *alarmHandle
		consumers []apiModel.PlcSubscriptionEventConsumer
	}
	s.mutex.RLock()
	var deliveries []delivery
	for registration, consumer := range s.consumers {
		for _, h := range registration.GetSubscriptionHandles() {
			handle, ok := h.(*alarmHandle)
			if !ok {
				continue
			}
			if _, active := s.handles[handle]; !active || !handle.tag.matches(alarm) {
				continue
			}
			deliveries = append(deliveries, delivery{handle: handle, consumers: []apiModel.PlcSubscriptionEventConsumer{consumer}})
		}
	}
	s.mutex.RUnlock()
	value := alarm.PlcValue()
	for _, d := range deliveries {
		name := d.handle.tagName
		event := spiModel.NewDefaultPlcSubscriptionEvent(nil,
			map[string]apiModel.PlcTag{name: d.handle.tag},
			map[string]apiModel.PlcSubscriptionType{name: d.handle.tagType},
			map[string]time.Duration{name: 0},
			map[string]apiModel.PlcResponseCode{name: apiModel.PlcResponseCode_OK},
			map[string]apiValues.PlcValue{name: value},
			append(s._options, options.WithCustomLogger(s.log))...)
		for _, consumer := range d.consumers {
			consumer(event)
		}
	}
}

// subscriberRouter sends subscriptions of ALARM tags to the alarm subscriber and the rest to the
// polling one. A request mixing the two is refused: subscribe to alarms in a request of their own.
type subscriberRouter struct {
	polling spi.PlcSubscriber
	alarms  *alarmSubscriber
}

var _ spi.PlcSubscriber = subscriberRouter{}

func (r subscriberRouter) Subscribe(ctx context.Context, request apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	alarms, others := 0, 0
	for _, name := range request.GetTagNames() {
		if _, ok := request.GetTag(name).(AlarmTag); ok {
			alarms++
		} else {
			others++
		}
	}
	switch {
	case alarms > 0 && others > 0:
		result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
		result <- spiModel.NewDefaultPlcSubscriptionRequestResult(request, nil, errors.New("subscribe to ALARM tags in a request of their own"))
		return result
	case alarms > 0:
		return r.alarms.Subscribe(ctx, request)
	default:
		return r.polling.Subscribe(ctx, request)
	}
}

func (r subscriberRouter) Unsubscribe(ctx context.Context, request apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	if u, ok := request.(*spiModel.DefaultPlcUnsubscriptionRequest); ok {
		for _, h := range u.GetSubscriptionHandles() {
			if _, isAlarm := h.(*alarmHandle); isAlarm {
				return r.alarms.Unsubscribe(ctx, request)
			}
		}
	}
	return r.polling.Unsubscribe(ctx, request)
}

func (r subscriberRouter) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	for _, h := range handles {
		if _, isAlarm := h.(*alarmHandle); isAlarm {
			return r.alarms.Register(consumer, handles)
		}
	}
	return r.polling.Register(consumer, handles)
}

func (r subscriberRouter) Unregister(registration apiModel.PlcConsumerRegistration) {
	r.alarms.Unregister(registration)
	r.polling.Unregister(registration)
}
