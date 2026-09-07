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

package _default

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const (
	pollTestTimeout  = 5 * time.Second
	pollTestQuietFor = 150 * time.Millisecond
)

///////////////////////////////////////
// Test doubles
//

// fakePollTag is a minimal tag which is both an apiModel.PlcTag (so it can be read) and an
// apiModel.PlcSubscriptionTag (so it can be subscribed).
type fakePollTag struct {
	address string
}

var (
	_ apiModel.PlcTag             = fakePollTag{}
	_ apiModel.PlcSubscriptionTag = fakePollTag{}
)

func (f fakePollTag) String() string                       { return f.address }
func (f fakePollTag) GetAddressString() string             { return f.address }
func (f fakePollTag) GetValueType() apiValues.PlcValueType { return apiValues.DINT }
func (f fakePollTag) GetArrayInfo() []apiModel.ArrayInfo   { return nil }
func (f fakePollTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionCyclic
}
func (f fakePollTag) GetDuration() time.Duration { return 0 }

type fakePollTagHandler struct{}

func (fakePollTagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	return fakePollTag{address: tagAddress}, nil
}

func (fakePollTagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("not implemented")
}

// fakePollTarget plays both roles a driver plays for the polling subscriber: it hands out read
// request builders (DefaultPollingSubscriberRequirements) and answers the reads (spi.PlcReader).
type fakePollTarget struct {
	mutex  sync.Mutex
	values map[string]apiValues.PlcValue
	codes  map[string]apiModel.PlcResponseCode

	reads       atomic.Int32
	builderFail atomic.Bool
	// block, if set, makes every read wait until it is closed (or the context is done). Used to
	// pin down the shutdown behaviour with a read in flight.
	block chan struct{}
	// reading is signalled every time a read is entered.
	reading chan struct{}
}

func newFakePollTarget() *fakePollTarget {
	return &fakePollTarget{
		values:  map[string]apiValues.PlcValue{},
		codes:   map[string]apiModel.PlcResponseCode{},
		reading: make(chan struct{}, 128),
	}
}

func (f *fakePollTarget) setValue(address string, value apiValues.PlcValue) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.values[address] = value
}

func (f *fakePollTarget) setCode(address string, code apiModel.PlcResponseCode) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.codes[address] = code
}

func (f *fakePollTarget) PollingReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	if f.builderFail.Load() {
		return nil
	}
	return spiModel.NewDefaultPlcReadRequestBuilder(fakePollTagHandler{}, f)
}

func (f *fakePollTarget) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	result := make(chan apiModel.PlcReadRequestResult, 1)
	f.reads.Add(1)
	select {
	case f.reading <- struct{}{}:
	default:
	}
	if f.block != nil {
		select {
		case <-f.block:
		case <-ctx.Done():
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, ctx.Err())
			return result
		}
	}
	codes := map[string]apiModel.PlcResponseCode{}
	values := map[string]apiValues.PlcValue{}
	f.mutex.Lock()
	for _, tagName := range readRequest.GetTagNames() {
		address := readRequest.GetTag(tagName).GetAddressString()
		if code, ok := f.codes[address]; ok {
			codes[tagName] = code
		} else {
			codes[tagName] = apiModel.PlcResponseCode_OK
		}
		values[tagName] = f.values[address]
	}
	f.mutex.Unlock()
	result <- spiModel.NewDefaultPlcReadRequestResult(
		readRequest,
		spiModel.NewDefaultPlcReadResponse(readRequest, codes, values),
		nil,
	)
	return result
}

// manualTicker is a PollTicker driven by the test instead of the wall clock.
type manualTicker struct {
	interval time.Duration
	c        chan time.Time
	stopped  atomic.Bool
}

func (m *manualTicker) Chan() <-chan time.Time { return m.c }

func (m *manualTicker) Stop() { m.stopped.Store(true) }

// tick delivers a single tick and blocks until the poller picked it up.
func (m *manualTicker) tick(t *testing.T) {
	t.Helper()
	select {
	case m.c <- time.Now():
	case <-time.After(pollTestTimeout):
		t.Fatalf("poller did not pick up a tick within %s", pollTestTimeout)
	}
}

// requireStopped waits until the poll go routine driven by this ticker really returned - it stops
// its ticker on the way out. Without this wait a tickExpectingNoPoller right after an unsubscribe
// races the poller's last loop iteration: the poller may still be on its way back to the select and
// then picks up the tick even though it is about to shut down.
func (m *manualTicker) requireStopped(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(pollTestTimeout)
	for !m.stopped.Load() {
		if time.Now().After(deadline) {
			t.Fatalf("the poller didn't stop its ticker within %s", pollTestTimeout)
		}
		time.Sleep(time.Millisecond)
	}
}

// tickExpectingNoPoller tries to deliver a tick and reports whether anybody was listening.
func (m *manualTicker) tickExpectingNoPoller() bool {
	select {
	case m.c <- time.Now():
		return true
	case <-time.After(pollTestQuietFor):
		return false
	}
}

// slowStopTicker never ticks but takes its time to stop, which keeps the poll go routine using it
// alive for a moment after it was told to shut down. Used to widen the shutdown window in
// TestDefaultPollingSubscriberCloseRacingSubscribe.
type slowStopTicker struct {
	c chan time.Time
}

func newSlowStopTicker(_ time.Duration) PollTicker {
	return &slowStopTicker{c: make(chan time.Time)}
}

func (s *slowStopTicker) Chan() <-chan time.Time { return s.c }

func (s *slowStopTicker) Stop() { time.Sleep(20 * time.Millisecond) }

// manualClock hands out manualTickers and records them in creation order.
type manualClock struct {
	mutex   sync.Mutex
	tickers []*manualTicker
}

func (m *manualClock) factory(interval time.Duration) PollTicker {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ticker := &manualTicker{interval: interval, c: make(chan time.Time)}
	m.tickers = append(m.tickers, ticker)
	return ticker
}

func (m *manualClock) all() []*manualTicker {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tickers := make([]*manualTicker, len(m.tickers))
	copy(tickers, m.tickers)
	return tickers
}

func (m *manualClock) only(t *testing.T) *manualTicker {
	t.Helper()
	tickers := m.all()
	require.Len(t, tickers, 1, "expected exactly one poller")
	return tickers[0]
}

///////////////////////////////////////
// Test helpers
//

// eventSink collects the events delivered to a consumer.
type eventSink struct {
	events chan apiModel.PlcSubscriptionEvent
}

func newEventSink() *eventSink {
	return &eventSink{events: make(chan apiModel.PlcSubscriptionEvent, 64)}
}

func (e *eventSink) consume(event apiModel.PlcSubscriptionEvent) {
	e.events <- event
}

func (e *eventSink) requireEvent(t *testing.T) apiModel.PlcSubscriptionEvent {
	t.Helper()
	select {
	case event := <-e.events:
		return event
	case <-time.After(pollTestTimeout):
		t.Fatalf("no subscription event within %s", pollTestTimeout)
		return nil
	}
}

func (e *eventSink) requireNoEvent(t *testing.T) {
	t.Helper()
	select {
	case event := <-e.events:
		t.Fatalf("unexpected subscription event: %v", event)
	case <-time.After(pollTestQuietFor):
	}
}

// subscribe executes a subscription request and returns the response.
func subscribe(t *testing.T, subscriber DefaultPollingSubscriber, build func(builder apiModel.PlcSubscriptionRequestBuilder)) apiModel.PlcSubscriptionResponse {
	t.Helper()
	builder := spiModel.NewDefaultPlcSubscriptionRequestBuilder(fakePollTagHandler{}, nil, subscriber)
	build(builder)
	request, err := builder.Build()
	require.NoError(t, err)
	select {
	case result := <-request.Execute(context.Background()):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		return result.GetResponse()
	case <-time.After(pollTestTimeout):
		t.Fatalf("subscribe did not return within %s", pollTestTimeout)
		return nil
	}
}

func registerConsumer(t *testing.T, subscriber DefaultPollingSubscriber, response apiModel.PlcSubscriptionResponse, sink *eventSink, tagNames ...string) apiModel.PlcConsumerRegistration {
	t.Helper()
	var handles []apiModel.PlcSubscriptionHandle
	for _, tagName := range tagNames {
		handle, err := response.GetSubscriptionHandle(tagName)
		require.NoError(t, err)
		handles = append(handles, handle)
	}
	require.NotEmpty(t, handles)
	// A single tag goes through the api-level flow (handle.Register) to keep that path covered.
	if len(handles) == 1 {
		return handles[0].Register(sink.consume)
	}
	return subscriber.Register(sink.consume, handles)
}

func unsubscribe(t *testing.T, handles ...apiModel.PlcSubscriptionHandle) {
	t.Helper()
	builder := spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
	builder.AddHandles(handles...)
	request, err := builder.Build()
	require.NoError(t, err)
	select {
	case result := <-request.Execute(context.Background()):
		require.NoError(t, result.GetErr())
	case <-time.After(pollTestTimeout):
		t.Fatalf("unsubscribe did not return within %s", pollTestTimeout)
	}
}

func handleFor(t *testing.T, response apiModel.PlcSubscriptionResponse, tagName string) apiModel.PlcSubscriptionHandle {
	t.Helper()
	handle, err := response.GetSubscriptionHandle(tagName)
	require.NoError(t, err)
	return handle
}

// registrationCount reports how many consumer registrations the subscriber currently keeps.
func registrationCount(subscriber *defaultPollingSubscriber) int {
	subscriber.consumersMutex.RLock()
	defer subscriber.consumersMutex.RUnlock()
	return len(subscriber.consumers)
}

// assertNoPollerLeak checks that the go routine count settles back at (or below) the baseline.
// Go routines the runtime tears down lazily make a single sample flaky, hence the retries.
func assertNoPollerLeak(t *testing.T, baseline int) {
	t.Helper()
	current := runtime.NumGoroutine()
	for i := 0; i < 100 && current > baseline; i++ {
		time.Sleep(20 * time.Millisecond)
		current = runtime.NumGoroutine()
	}
	if current > baseline {
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		t.Fatalf("goroutine leak detected: baseline %d, now %d\n%s", baseline, current, buf[:n])
	}
}

func closeWithin(t *testing.T, subscriber DefaultPollingSubscriber) {
	t.Helper()
	closed := make(chan struct{})
	go func() {
		defer close(closed)
		subscriber.Close()
	}()
	select {
	case <-closed:
	case <-time.After(pollTestTimeout):
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		t.Fatalf("Close did not return within %s\n%s", pollTestTimeout, buf[:n])
	}
}

///////////////////////////////////////
// Tests
//

// TestDefaultPollingSubscriberCyclicEmitsEveryPoll pins the CYCLIC contract: every poll produces an
// event, even when the value didn't move.
func TestDefaultPollingSubscriberCyclicEmitsEveryPoll(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("holding-register:1", spiValues.NewPlcDINT(42))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "holding-register:1", 10*time.Millisecond)
	})
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("tag"))

	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	assert.Equal(t, 10*time.Millisecond, ticker.interval, "the poller has to run at the configured interval")

	ticker.tick(t)
	event := sink.requireEvent(t)
	assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("tag"))
	assert.Equal(t, int32(42), event.GetValue("tag").GetInt32())
	assert.Equal(t, "holding-register:1", event.GetAddress("tag"))

	// Unchanged value, but CYCLIC fires anyway.
	ticker.tick(t)
	event = sink.requireEvent(t)
	assert.Equal(t, int32(42), event.GetValue("tag").GetInt32())

	target.setValue("holding-register:1", spiValues.NewPlcDINT(43))
	ticker.tick(t)
	event = sink.requireEvent(t)
	assert.Equal(t, int32(43), event.GetValue("tag").GetInt32())

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberCyclicUsesRealClock runs the very same flow through the production
// time.Ticker with a short interval, so the wiring of the default ticker factory is covered too.
func TestDefaultPollingSubscriberCyclicUsesRealClock(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("holding-register:1", spiValues.NewPlcDINT(7))
	interval := 20 * time.Millisecond
	subscriber := NewDefaultPollingSubscriber(target)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "holding-register:1", interval)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	start := time.Now()
	for range 3 {
		event := sink.requireEvent(t)
		assert.Equal(t, int32(7), event.GetValue("tag").GetInt32())
	}
	// Three events can't have arrived faster than two intervals after the first one.
	assert.GreaterOrEqual(t, time.Since(start), 2*interval, "events must not be emitted faster than the interval")

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberChangeOfStateSuppressesUnchangedValues pins the CHANGE_OF_STATE
// contract: the first poll fires, further polls only fire on an actual change.
func TestDefaultPollingSubscriberChangeOfStateSuppressesUnchangedValues(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("coil:3", spiValues.NewPlcBOOL(false))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target,
		WithPollTickerFactory(clock.factory),
		WithDefaultPollingInterval(30*time.Millisecond),
	)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddChangeOfStateTagAddress("tag", "coil:3")
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	assert.Equal(t, 30*time.Millisecond, ticker.interval, "change-of-state tags have to fall back to the default interval")

	// The very first poll has no baseline yet, so it fires.
	ticker.tick(t)
	event := sink.requireEvent(t)
	assert.False(t, event.GetValue("tag").GetBool())

	// Same value twice: no event.
	ticker.tick(t)
	sink.requireNoEvent(t)
	ticker.tick(t)
	sink.requireNoEvent(t)

	// Changed value: event.
	target.setValue("coil:3", spiValues.NewPlcBOOL(true))
	ticker.tick(t)
	event = sink.requireEvent(t)
	assert.True(t, event.GetValue("tag").GetBool())

	// And it stays quiet again on the new baseline.
	ticker.tick(t)
	sink.requireNoEvent(t)

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberGroupsTagsAndRejectsEventTags pins the grouping: tags sharing type and
// interval are polled by a single read, and EVENT tags are answered as unsupported.
func TestDefaultPollingSubscriberGroupsTagsAndRejectsEventTags(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	target.setValue("b", spiValues.NewPlcDINT(2))
	target.setValue("c", spiValues.NewPlcDINT(3))
	target.setValue("d", spiValues.NewPlcDINT(4))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("first", "a", 10*time.Millisecond)
		builder.AddCyclicTagAddress("second", "b", 10*time.Millisecond)
		builder.AddCyclicTagAddress("third", "c", 20*time.Millisecond)
		builder.AddChangeOfStateTagAddress("fourth", "d")
		builder.AddEventTagAddress("fifth", "e")
	})

	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("first"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("second"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("third"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("fourth"))
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, response.GetResponseCode("fifth"),
		"a poller can't emulate spontaneous events")

	tickers := clock.all()
	require.Len(t, tickers, 3, "one poller per subscription type and interval")
	assert.Equal(t, 10*time.Millisecond, tickers[0].interval)
	assert.Equal(t, 20*time.Millisecond, tickers[1].interval)
	assert.Equal(t, DefaultPollingInterval, tickers[2].interval)

	// "first" and "second" share a poller, so they share an event and a single read.
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "first", "second")
	tickers[0].tick(t)
	event := sink.requireEvent(t)
	assert.ElementsMatch(t, []string{"first", "second"}, event.GetTagNames())
	assert.Equal(t, int32(1), event.GetValue("first").GetInt32())
	assert.Equal(t, int32(2), event.GetValue("second").GetInt32())
	assert.Equal(t, int32(1), target.reads.Load(), "grouped tags have to be polled with one read")

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberUnsubscribeStopsEmission pins that unsubscribing tears the poller down.
func TestDefaultPollingSubscriberUnsubscribeStopsEmission(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	ticker.tick(t)
	sink.requireEvent(t)

	unsubscribe(t, handleFor(t, response, "tag"))

	// The poll go routine is gone, so nobody picks up a tick anymore and no event is emitted.
	ticker.requireStopped(t)
	assert.False(t, ticker.tickExpectingNoPoller(), "the poller has to be gone after unsubscribing")
	assert.True(t, ticker.stopped.Load(), "the ticker has to be stopped")
	sink.requireNoEvent(t)

	readsAfterUnsubscribe := target.reads.Load()
	assert.Equal(t, int32(1), readsAfterUnsubscribe, "no further polls after unsubscribing")

	// Unsubscribing twice is a no-op, not an error.
	unsubscribe(t, handleFor(t, response, "tag"))

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberUnsubscribeSingleTagOfGroup pins that unsubscribing one tag of a group
// keeps the remaining tags of that group polled.
func TestDefaultPollingSubscriberUnsubscribeSingleTagOfGroup(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	target.setValue("b", spiValues.NewPlcDINT(2))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("first", "a", 10*time.Millisecond)
		builder.AddCyclicTagAddress("second", "b", 10*time.Millisecond)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "first", "second")

	ticker := clock.only(t)
	ticker.tick(t)
	event := sink.requireEvent(t)
	assert.ElementsMatch(t, []string{"first", "second"}, event.GetTagNames())

	unsubscribe(t, handleFor(t, response, "first"))

	ticker.tick(t)
	event = sink.requireEvent(t)
	assert.Equal(t, []string{"second"}, event.GetTagNames(), "only the remaining tag stays polled")

	unsubscribe(t, handleFor(t, response, "second"))
	ticker.requireStopped(t)
	assert.False(t, ticker.tickExpectingNoPoller(), "the poller has to be gone once the group ran empty")

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberFailedReadsDoNotFire pins that failing tags are reported but never
// make an event on their own (mirroring plc4j's PollingSubscriptionConnectionBase).
func TestDefaultPollingSubscriberFailedReadsDoNotFire(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setCode("a", apiModel.PlcResponseCode_NOT_FOUND)
	target.setValue("b", spiValues.NewPlcDINT(2))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("broken", "a", 10*time.Millisecond)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "broken")

	ticker := clock.only(t)
	ticker.tick(t)
	sink.requireNoEvent(t)

	// Add a healthy tag on the same interval; now the event carries both.
	secondResponse := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("broken", "a", 10*time.Millisecond)
		builder.AddCyclicTagAddress("healthy", "b", 10*time.Millisecond)
	})
	secondSink := newEventSink()
	registerConsumer(t, subscriber, secondResponse, secondSink, "broken", "healthy")

	tickers := clock.all()
	require.Len(t, tickers, 2, "a second subscribe has to get a poller of its own")
	tickers[1].tick(t)
	event := secondSink.requireEvent(t)
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, event.GetResponseCode("broken"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("healthy"))
	assert.Equal(t, int32(2), event.GetValue("healthy").GetInt32())

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberConcurrentRegisterUnregister hammers the consumer bookkeeping while
// the poller is running. Meant to be run with -race.
func TestDefaultPollingSubscriberConcurrentRegisterUnregister(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	subscriber := NewDefaultPollingSubscriber(target)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", time.Millisecond)
	})
	handle := handleFor(t, response, "tag")

	var deliveries atomic.Int64
	var churn sync.WaitGroup
	stop := make(chan struct{})
	for range 8 {
		churn.Go(func() {
			for {
				select {
				case <-stop:
					return
				default:
				}
				registration := handle.Register(func(_ apiModel.PlcSubscriptionEvent) {
					deliveries.Add(1)
				})
				require.NotNil(t, registration)
				registration.Unregister()
			}
		})
	}
	// Keep a stable consumer around so we can be sure polls really happened while we churned.
	sink := newEventSink()
	stableRegistration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		select {
		case sink.events <- event:
		default:
		}
	})
	sink.requireEvent(t)
	time.Sleep(200 * time.Millisecond)
	close(stop)
	churn.Wait()
	stableRegistration.Unregister()

	assert.Positive(t, target.reads.Load(), "the poller has to have polled while we churned")

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberCloseWithReadInFlight pins that Close doesn't block on a read which is
// still waiting for the PLC and doesn't leak the poll go routine.
func TestDefaultPollingSubscriberCloseWithReadInFlight(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.block = make(chan struct{})
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	ticker.tick(t)
	// Wait until the read is actually in flight (and stuck).
	select {
	case <-target.reading:
	case <-time.After(pollTestTimeout):
		t.Fatalf("no read started within %s", pollTestTimeout)
	}

	closeWithin(t, subscriber)
	assert.True(t, ticker.stopped.Load(), "the ticker has to be stopped on close")
	sink.requireNoEvent(t)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberCloseIsIdempotent pins that a second Close is a no-op.
func TestDefaultPollingSubscriberCloseIsIdempotent(t *testing.T) {
	target := newFakePollTarget()
	subscriber := NewDefaultPollingSubscriber(target)
	closeWithin(t, subscriber)
	closeWithin(t, subscriber)
}

// TestDefaultPollingSubscriberSubscribeAfterCloseFails pins that a closed subscriber refuses new
// subscriptions instead of silently starting a poller nobody stops.
func TestDefaultPollingSubscriberSubscribeAfterCloseFails(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))
	closeWithin(t, subscriber)

	builder := spiModel.NewDefaultPlcSubscriptionRequestBuilder(fakePollTagHandler{}, nil, subscriber)
	builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	request, err := builder.Build()
	require.NoError(t, err)
	select {
	case result := <-request.Execute(context.Background()):
		assert.Error(t, result.GetErr())
	case <-time.After(pollTestTimeout):
		t.Fatalf("subscribe did not return within %s", pollTestTimeout)
	}
	assert.Empty(t, clock.all(), "no poller may be started after close")
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberUnsupportedRequestTypes pins the guards against foreign message types.
func TestDefaultPollingSubscriberUnsupportedRequestTypes(t *testing.T) {
	target := newFakePollTarget()
	subscriber := NewDefaultPollingSubscriber(target)
	defer closeWithin(t, subscriber)

	select {
	case result := <-subscriber.Subscribe(context.Background(), &foreignSubscriptionRequest{}):
		assert.Error(t, result.GetErr())
	case <-time.After(pollTestTimeout):
		t.Fatal("Subscribe did not return")
	}

	select {
	case result := <-subscriber.Unsubscribe(context.Background(), &foreignUnsubscriptionRequest{}):
		assert.Error(t, result.GetErr())
	case <-time.After(pollTestTimeout):
		t.Fatal("Unsubscribe did not return")
	}
}

// TestDefaultPollingSubscriberUnsubscribeIgnoresForeignHandles pins that a mixed unsubscription
// request (spiModel.DefaultPlcUnsubscriptionRequest passes the whole request to every subscriber)
// doesn't produce spurious errors.
func TestDefaultPollingSubscriberUnsubscribeIgnoresForeignHandles(t *testing.T) {
	target := newFakePollTarget()
	subscriber := NewDefaultPollingSubscriber(target)
	defer closeWithin(t, subscriber)

	foreign := spiModel.NewDefaultPlcSubscriptionHandle(subscriber)
	request := spiModel.NewDefaultPlcUnsubscriptionRequest([]apiModel.PlcSubscriptionHandle{foreign})
	select {
	case result := <-subscriber.Unsubscribe(context.Background(), request):
		assert.NoError(t, result.GetErr())
	case <-time.After(pollTestTimeout):
		t.Fatal("Unsubscribe did not return")
	}
}

// TestDefaultPollingSubscriberSurvivesBrokenBuilder pins that a driver failing to hand out a read
// request builder neither panics nor kills the poller.
func TestDefaultPollingSubscriberSurvivesBrokenBuilder(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	target.builderFail.Store(true)
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	ticker.tick(t)
	sink.requireNoEvent(t)

	// The poller is still alive and recovers as soon as the driver hands out builders again.
	target.builderFail.Store(false)
	ticker.tick(t)
	event := sink.requireEvent(t)
	assert.Equal(t, int32(1), event.GetValue("tag").GetInt32())

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberPreRegisteredConsumers pins that consumers registered on the request
// receive events without an explicit Register call.
func TestDefaultPollingSubscriberPreRegisteredConsumers(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(5))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	sink := newEventSink()
	subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
		builder.AddPreRegisteredConsumer("tag", sink.consume)
	})

	clock.only(t).tick(t)
	event := sink.requireEvent(t)
	assert.Equal(t, int32(5), event.GetValue("tag").GetInt32())

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberCustomValueComparator pins that a driver can override the change
// detection.
func TestDefaultPollingSubscriberCustomValueComparator(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	// This comparator considers everything equal, so nothing but the first poll ever fires.
	subscriber := NewDefaultPollingSubscriber(target,
		WithPollTickerFactory(clock.factory),
		WithPlcValueComparator(func(_, _ apiValues.PlcValue) bool { return true }),
	)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddChangeOfStateTagAddress("tag", "a")
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	ticker.tick(t)
	sink.requireEvent(t)

	target.setValue("a", spiValues.NewPlcDINT(2))
	ticker.tick(t)
	sink.requireNoEvent(t)

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberChangeOfStateOnValuesWithoutString pins that the change detection also
// works for value types which don't implement String() (PlcValueAdapter renders those as
// "not implemented"), as produced by s7, knxnetip, bacnetip, cbus and opcua.
func TestDefaultPollingSubscriberChangeOfStateOnValuesWithoutString(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("raw", spiValues.NewPlcRawByteArray([]byte{0x01, 0x02}))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddChangeOfStateTagAddress("tag", "raw")
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "tag")

	ticker := clock.only(t)
	ticker.tick(t)
	event := sink.requireEvent(t)
	assert.Equal(t, []byte{0x01, 0x02}, event.GetValue("tag").GetRaw())

	// Unchanged, so nothing fires.
	ticker.tick(t)
	sink.requireNoEvent(t)

	target.setValue("raw", spiValues.NewPlcRawByteArray([]byte{0xFF, 0xEE}))
	ticker.tick(t)
	event = sink.requireEvent(t)
	assert.Equal(t, []byte{0xFF, 0xEE}, event.GetValue("tag").GetRaw())

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberCloseRacingSubscribe pins that a Close overlapping an in-flight
// Subscribe stays a well behaved no-op for both sides. Adding to a sync.WaitGroup whose counter drops
// to zero while somebody sits in its Wait panics on the caller's stack, i.e. outside of every recover
// of this file, and would take the whole process down.
//
// The slowStopTicker is what makes the window big enough to be hit reliably: it keeps the poll go
// routine (and with it the wait group counter) alive for a moment after Close already registered as a
// waiter, so the Subscribe calls hammering away below land in exactly that window.
func TestDefaultPollingSubscriberCloseRacingSubscribe(t *testing.T) {
	// No assertNoPollerLeak here: the hundreds of thousands of request executions below make the
	// shared spi/pool executor scale up a worker of its own, which has nothing to do with this
	// subscriber. The state assertions after every Close cover what this test is about.
	for range 20 {
		target := newFakePollTarget()
		target.setValue("a", spiValues.NewPlcDINT(1))
		subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(newSlowStopTicker))
		// A first subscription so a poll go routine holds the wait group counter above zero.
		subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
			builder.AddCyclicTagAddress("first", "a", time.Hour)
		})

		builder := spiModel.NewDefaultPlcSubscriptionRequestBuilder(fakePollTagHandler{}, nil, subscriber)
		builder.AddCyclicTagAddress("second", "a", time.Hour)
		request, err := builder.Build()
		require.NoError(t, err)

		closed := make(chan struct{})
		go func() {
			defer close(closed)
			subscriber.Close()
		}()

		var subscribes atomic.Int64
		var hammers sync.WaitGroup
		for range 8 {
			hammers.Go(func() {
				for {
					select {
					case <-closed:
						return
					default:
					}
					// A response or a "subscriber is closed" error, both are fine - it just must not
					// blow up.
					select {
					case <-request.Execute(context.Background()):
						subscribes.Add(1)
					case <-time.After(pollTestTimeout):
						t.Error("subscribe didn't return in time")
						return
					}
				}
			})
		}
		hammers.Wait()
		require.Positive(t, subscribes.Load(), "the racing subscribes have to have run")

		// Nothing the race let through may outlive the Close.
		subscriber.Close()
		internal, ok := subscriber.(*defaultPollingSubscriber)
		require.True(t, ok)
		internal.subscriptionsMutex.Lock()
		assert.Empty(t, internal.subscriptions, "no subscription may survive the Close")
		internal.subscriptionsMutex.Unlock()
		assert.Zero(t, registrationCount(internal), "no consumer registration may survive the Close")
	}
}

// TestDefaultPollingSubscriberPreRegisteredConsumerGetsTheFirstPoll pins that the pollers are only
// let loose after the response - and with it the consumers pre-registered on the request - is wired
// up. The first poll of a CHANGE_OF_STATE subscription is the only one reporting the current value,
// so a consumer missing it never learns the state until it happens to change. The real clock with a
// minimal interval is used on purpose: a manual ticker only ticks when the test says so and can
// therefore never observe this race.
func TestDefaultPollingSubscriberPreRegisteredConsumerGetsTheFirstPoll(t *testing.T) {
	baseline := runtime.NumGoroutine()
	for iteration := range 50 {
		target := newFakePollTarget()
		target.setValue("a", spiValues.NewPlcDINT(int32(iteration)))
		subscriber := NewDefaultPollingSubscriber(target, WithDefaultPollingInterval(time.Nanosecond))

		sink := newEventSink()
		subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
			builder.AddChangeOfStateTagAddress("tag", "a")
			builder.AddPreRegisteredConsumer("tag", sink.consume)
		})

		event := sink.requireEvent(t)
		assert.Equal(t, int32(iteration), event.GetValue("tag").GetInt32())
		closeWithin(t, subscriber)
	}
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberBlockingConsumerDoesNotBlockClose pins the bound on the wait Close does
// for the callbacks in flight: a consumer which never returns must not be able to keep Close from
// returning. A callback in flight can't be cancelled, so the only way out is for Close to give up on
// it once the grace period passed - which also has to hold for a callback which re-enters Subscribe
// while Close is running (that re-entrant subscribe has to be refused instead of deadlocking, which
// is why Close only fences on the spawn mutex instead of holding it).
func TestDefaultPollingSubscriberBlockingConsumerDoesNotBlockClose(t *testing.T) {
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target,
		WithPollTickerFactory(clock.factory),
		// Keeps the grace period Close spends on the stuck callback below short.
		WithConsumerCallbackTimeout(50*time.Millisecond),
	)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	// Built up front, used from within the callback below.
	reentrantBuilder := spiModel.NewDefaultPlcSubscriptionRequestBuilder(fakePollTagHandler{}, nil, subscriber)
	reentrantBuilder.AddCyclicTagAddress("reentrant", "a", 10*time.Millisecond)
	reentrantRequest, err := reentrantBuilder.Build()
	require.NoError(t, err)

	entered := make(chan struct{})
	release := make(chan struct{})
	reentrantErrs := make(chan error, 1)
	registration := handleFor(t, response, "tag").Register(func(_ apiModel.PlcSubscriptionEvent) {
		close(entered)
		<-release
		result := <-subscriber.Subscribe(context.Background(), reentrantRequest)
		reentrantErrs <- result.GetErr()
	})
	require.NotNil(t, registration)

	ticker := clock.only(t)
	ticker.tick(t)
	select {
	case <-entered:
	case <-time.After(pollTestTimeout):
		t.Fatalf("the consumer wasn't called within %s", pollTestTimeout)
	}

	// The callback is still sitting in its channel receive and will stay there until we let it go -
	// Close has to come back anyway (after its grace period), and it has to have stopped the poller
	// while doing so.
	closeWithin(t, subscriber)
	assert.True(t, ticker.stopped.Load(), "the ticker has to be stopped once Close returned")

	close(release)
	select {
	case reentrantErr := <-reentrantErrs:
		assert.Error(t, reentrantErr, "the closed subscriber has to refuse the re-entrant subscribe")
	case <-time.After(pollTestTimeout):
		t.Fatalf("the re-entrant subscribe didn't return within %s", pollTestTimeout)
	}
	// No assertNoPollerLeak here: the re-entrant subscribe drags process wide, lazily started go
	// routines in which have nothing to do with this subscriber.
}

// TestDefaultPollingSubscriberBlockingConsumerDoesNotStallPolling pins the other half of that bound:
// a consumer which never returns must not stall the poller of its subscription, must not keep the
// other consumers of that subscription from being served, must not be handed a second event while the
// first one is still in flight, and must not collect one abandoned go routine per poll. Once its
// event queue ran full the poller waits for it for at most the consumer callback timeout and then
// drops the event.
func TestDefaultPollingSubscriberBlockingConsumerDoesNotStallPolling(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target,
		WithPollTickerFactory(clock.factory),
		WithConsumerCallbackTimeout(50*time.Millisecond),
	)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	handle := handleFor(t, response, "tag")

	// The blocking consumer deliberately doesn't use eventSink: its channel is buffered 64 deep, so
	// it could never block in the first place.
	entered := make(chan struct{}, 8)
	release := make(chan struct{})
	require.NotNil(t, handle.Register(func(_ apiModel.PlcSubscriptionEvent) {
		entered <- struct{}{}
		<-release
	}))
	sink := newEventSink()
	require.NotNil(t, handle.Register(sink.consume))

	ticker := clock.only(t)
	ticker.tick(t)
	select {
	case <-entered:
	case <-time.After(pollTestTimeout):
		t.Fatalf("the blocking consumer wasn't called within %s", pollTestTimeout)
	}

	// The poll go routine has to keep picking up ticks although the callback is still stuck. The first
	// consumerEventQueueDepth events go into the stuck consumer's queue, the ones after that hit the
	// full queue and are dropped once the callback timeout expires.
	for range consumerEventQueueDepth + 2 {
		ticker.tick(t)
	}
	sink.requireEvent(t)
	sink.requireEvent(t)
	assert.Greater(t, target.reads.Load(), int32(1), "the poller has to have kept polling")
	select {
	case <-entered:
		t.Fatal("the stuck consumer must not be handed a second event while the first one is in flight")
	default:
	}

	closeWithin(t, subscriber)
	close(release)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberConcurrentPollersFeedOneRegistration pins that a registration which
// spans several subscriptions gets *all* their events: Register takes a slice of handles and one
// Subscribe is split into one poller per type and interval, so two poll go routines legitimately fan
// out to the same registration at the same time. Serializing them must not turn into dropping one of
// them, which is what the second tick below would run into if the in flight guard were a "skip while
// busy" instead of a queue.
func TestDefaultPollingSubscriberConcurrentPollersFeedOneRegistration(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	target.setValue("b", spiValues.NewPlcDINT(2))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	// Two intervals, hence two pollers, both feeding the single registration below.
	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("first", "a", 10*time.Millisecond)
		builder.AddCyclicTagAddress("second", "b", 20*time.Millisecond)
	})
	entered := make(chan struct{}, 4)
	gate := make(chan struct{})
	delivered := make(chan apiModel.PlcSubscriptionEvent, 4)
	registration := subscriber.Register(func(event apiModel.PlcSubscriptionEvent) {
		entered <- struct{}{}
		<-gate
		delivered <- event
	}, []apiModel.PlcSubscriptionHandle{handleFor(t, response, "first"), handleFor(t, response, "second")})
	require.NotNil(t, registration)

	tickers := clock.all()
	require.Len(t, tickers, 2)

	// The first poller's event puts the callback in flight and keeps it there ...
	tickers[0].tick(t)
	select {
	case <-entered:
	case <-time.After(pollTestTimeout):
		t.Fatalf("the consumer wasn't called within %s", pollTestTimeout)
	}
	// ... while the second poller hands over its event.
	tickers[1].tick(t)

	close(gate)
	var tagNames []string
	for i := range 2 {
		select {
		case event := <-delivered:
			tagNames = append(tagNames, event.GetTagNames()...)
		case <-time.After(pollTestTimeout):
			t.Fatalf("only %d of the 2 events were delivered within %s", i, pollTestTimeout)
		}
	}
	assert.ElementsMatch(t, []string{"first", "second"}, tagNames,
		"both pollers' events have to reach the registration")

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberCloseWaitsForARunningCallback pins the other side of the Close bound:
// Close does wait for a callback which is merely slow. A driver closes its connection right after
// closing the subscriber, and the event a callback is holding points at that connection's state, so
// returning from Close while a well behaved consumer is still running would hand it a torn down
// connection.
func TestDefaultPollingSubscriberCloseWaitsForARunningCallback(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	entered := make(chan struct{})
	var returned atomic.Bool
	require.NotNil(t, handleFor(t, response, "tag").Register(func(_ apiModel.PlcSubscriptionEvent) {
		close(entered)
		time.Sleep(2 * pollTestQuietFor)
		returned.Store(true)
	}))

	ticker := clock.only(t)
	ticker.tick(t)
	select {
	case <-entered:
	case <-time.After(pollTestTimeout):
		t.Fatalf("the consumer wasn't called within %s", pollTestTimeout)
	}

	closeWithin(t, subscriber)
	assert.True(t, returned.Load(), "Close returned while a consumer callback was still running")
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberUnsubscribeDropsStaleRegistrations pins that unsubscribing also cleans
// up the consumer registrations which only referenced the unsubscribed tags. A connection which
// subscribes and unsubscribes over and over would otherwise pile up dead registrations, which the
// fan out walks on every single poll.
func TestDefaultPollingSubscriberUnsubscribeDropsStaleRegistrations(t *testing.T) {
	baseline := runtime.NumGoroutine()
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	target.setValue("b", spiValues.NewPlcDINT(2))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))
	internal, ok := subscriber.(*defaultPollingSubscriber)
	require.True(t, ok)

	for range 5 {
		response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
			builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
		})
		sink := newEventSink()
		registerConsumer(t, subscriber, response, sink, "tag")
		require.Equal(t, 1, registrationCount(internal))

		unsubscribe(t, handleFor(t, response, "tag"))
		assert.Zero(t, registrationCount(internal), "the registration of an unsubscribed tag has to be gone")
	}

	// A registration still referencing a subscribed tag survives.
	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("first", "a", 10*time.Millisecond)
		builder.AddCyclicTagAddress("second", "b", 10*time.Millisecond)
	})
	sink := newEventSink()
	registerConsumer(t, subscriber, response, sink, "first", "second")
	unsubscribe(t, handleFor(t, response, "first"))
	assert.Equal(t, 1, registrationCount(internal), "a registration with a live handle has to stay")
	unsubscribe(t, handleFor(t, response, "second"))
	assert.Zero(t, registrationCount(internal))

	closeWithin(t, subscriber)
	assertNoPollerLeak(t, baseline)
}

// TestDefaultPollingSubscriberRegisterAfterCloseIsInert pins that a late Register doesn't repopulate
// the consumer map Close cleared.
func TestDefaultPollingSubscriberRegisterAfterCloseIsInert(t *testing.T) {
	target := newFakePollTarget()
	target.setValue("a", spiValues.NewPlcDINT(1))
	clock := &manualClock{}
	subscriber := NewDefaultPollingSubscriber(target, WithPollTickerFactory(clock.factory))
	internal, ok := subscriber.(*defaultPollingSubscriber)
	require.True(t, ok)

	response := subscribe(t, subscriber, func(builder apiModel.PlcSubscriptionRequestBuilder) {
		builder.AddCyclicTagAddress("tag", "a", 10*time.Millisecond)
	})
	handle := handleFor(t, response, "tag")
	closeWithin(t, subscriber)
	require.Zero(t, registrationCount(internal))

	sink := newEventSink()
	registration := subscriber.Register(sink.consume, []apiModel.PlcSubscriptionHandle{handle})
	assert.NotNil(t, registration, "the caller still has to get something it can Unregister")
	assert.Zero(t, registrationCount(internal), "a closed subscriber must not take new consumers")
	// Unregistering the inert registration stays a no-op.
	registration.Unregister()
	assert.Zero(t, registrationCount(internal))
}

func TestPlcValuesEqual(t *testing.T) {
	tests := []struct {
		name     string
		previous apiValues.PlcValue
		current  apiValues.PlcValue
		want     bool
	}{
		{name: "both nil", previous: nil, current: nil, want: true},
		{name: "previous nil", previous: nil, current: spiValues.NewPlcDINT(1), want: false},
		{name: "current nil", previous: spiValues.NewPlcDINT(1), current: nil, want: false},
		{name: "same value", previous: spiValues.NewPlcDINT(1), current: spiValues.NewPlcDINT(1), want: true},
		{name: "different value", previous: spiValues.NewPlcDINT(1), current: spiValues.NewPlcDINT(2), want: false},
		{name: "different type", previous: spiValues.NewPlcDINT(1), current: spiValues.NewPlcBOOL(true), want: false},
		{name: "same bool", previous: spiValues.NewPlcBOOL(true), current: spiValues.NewPlcBOOL(true), want: true},
		{name: "different bool", previous: spiValues.NewPlcBOOL(true), current: spiValues.NewPlcBOOL(false), want: false},
		{
			name:     "same list",
			previous: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcDINT(1), spiValues.NewPlcDINT(2)}),
			current:  spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcDINT(1), spiValues.NewPlcDINT(2)}),
			want:     true,
		},
		{
			name:     "different list",
			previous: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcDINT(1)}),
			current:  spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcDINT(2)}),
			want:     false,
		},
		// The value types below don't implement String() of their own, so PlcValueAdapter renders
		// them all as "not implemented". A String() based comparison would call every one of them
		// equal and silence a CHANGE_OF_STATE subscription for good.
		{
			name:     "same raw byte array",
			previous: spiValues.NewPlcRawByteArray([]byte{0x01, 0x02}),
			current:  spiValues.NewPlcRawByteArray([]byte{0x01, 0x02}),
			want:     true,
		},
		{
			name:     "different raw byte array",
			previous: spiValues.NewPlcRawByteArray([]byte{0x01, 0x02}),
			current:  spiValues.NewPlcRawByteArray([]byte{0xFF, 0xEE}),
			want:     false,
		},
		{
			name:     "raw byte array of a different length",
			previous: spiValues.NewPlcRawByteArray([]byte{0x01, 0x02}),
			current:  spiValues.NewPlcRawByteArray([]byte{0x01, 0x02, 0x03}),
			want:     false,
		},
		{
			name:     "different list of raw byte arrays",
			previous: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcRawByteArray([]byte{0x01})}),
			current:  spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcRawByteArray([]byte{0x02})}),
			want:     false,
		},
		{
			name:     "same struct of raw byte arrays",
			previous: spiValues.NewPlcStruct(map[string]apiValues.PlcValue{"raw": spiValues.NewPlcRawByteArray([]byte{0x01})}),
			current:  spiValues.NewPlcStruct(map[string]apiValues.PlcValue{"raw": spiValues.NewPlcRawByteArray([]byte{0x01})}),
			want:     true,
		},
		{
			name:     "different struct of raw byte arrays",
			previous: spiValues.NewPlcStruct(map[string]apiValues.PlcValue{"raw": spiValues.NewPlcRawByteArray([]byte{0x01})}),
			current:  spiValues.NewPlcStruct(map[string]apiValues.PlcValue{"raw": spiValues.NewPlcRawByteArray([]byte{0x02})}),
			want:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, PlcValuesEqual(test.previous, test.current))
		})
	}
}

// TestDefaultPollingSubscriberOptionPassing pins the option plumbing.
func TestDefaultPollingSubscriberOptionPassing(t *testing.T) {
	target := newFakePollTarget()
	subscriber := NewDefaultPollingSubscriber(target,
		options.WithCustomLogger(options.ExtractCustomLoggerOrDefaultToGlobal()),
		WithDefaultPollingInterval(0),  // ignored, has to keep the default
		WithConsumerCallbackTimeout(0), // ignored, has to keep the default
	)
	internal, ok := subscriber.(*defaultPollingSubscriber)
	require.True(t, ok)
	assert.Equal(t, DefaultPollingInterval, internal.defaultPollingInterval)
	assert.Equal(t, DefaultConsumerCallbackTimeout, internal.consumerCallbackTimeout)
	assert.NotNil(t, internal.plcValueComparator)
	assert.NotNil(t, internal.pollTickerFactory)
	closeWithin(t, subscriber)
}

///////////////////////////////////////
// Foreign message types used by the guard tests
//

type foreignSubscriptionRequest struct{}

func (f *foreignSubscriptionRequest) String() string      { return "foreignSubscriptionRequest" }
func (f *foreignSubscriptionRequest) IsAPlcMessage() bool { return true }
func (f *foreignSubscriptionRequest) GetTagNames() []string {
	return nil
}
func (f *foreignSubscriptionRequest) GetTag(_ string) apiModel.PlcSubscriptionTag { return nil }
func (f *foreignSubscriptionRequest) Execute(_ context.Context) <-chan apiModel.PlcSubscriptionRequestResult {
	return nil
}

type foreignUnsubscriptionRequest struct{}

func (f *foreignUnsubscriptionRequest) String() string      { return "foreignUnsubscriptionRequest" }
func (f *foreignUnsubscriptionRequest) IsAPlcMessage() bool { return true }
func (f *foreignUnsubscriptionRequest) Execute(_ context.Context) <-chan apiModel.PlcUnsubscriptionRequestResult {
	return nil
}
