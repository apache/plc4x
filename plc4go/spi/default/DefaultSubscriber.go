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
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
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

// DefaultPollingSubscriberRequirements adds the methods a driver has to provide in order to use a
// DefaultPollingSubscriber. It is the Go counterpart of the abstract parts of plc4j's
// PollingSubscriptionConnectionBase.
type DefaultPollingSubscriberRequirements interface {
	// PollingReadRequestBuilder has to return a *fresh* apiModel.PlcReadRequestBuilder on every call.
	// Usually this is simply the connection's ReadRequestBuilder(). The subscriber assembles one read
	// request per poll cycle from it, using apiModel.PlcReadRequestBuilder.AddTagAddress with the
	// address string of the subscribed tag.
	PollingReadRequestBuilder() apiModel.PlcReadRequestBuilder
}

// DefaultPollingSubscriber emulates subscriptions for drivers whose protocol has no native
// subscription support (Modbus, ADS-less setups, ...) by periodically issuing read requests.
//
// It supports:
//   - apiModel.SubscriptionCyclic: an event is emitted on every poll.
//   - apiModel.SubscriptionChangeOfState: an event is emitted only when a value differs from the
//     value seen on the previous poll (the very first poll always emits).
//   - apiModel.SubscriptionEvent: not emulatable, such tags are answered with
//     apiModel.PlcResponseCode_UNSUPPORTED.
//
// Adoption is opt-in and by construction, in the same style as DefaultConnection or DefaultCodec:
//
//	type Connection struct {
//		_default.DefaultConnection
//		subscriber _default.DefaultPollingSubscriber
//	}
//
//	connection.subscriber = _default.NewDefaultPollingSubscriber(connection)
//
// The driver then delegates spi.PlcSubscriber to the embedded subscriber and has to call Close on it
// while closing the connection so no poller outlives the connection.
type DefaultPollingSubscriber interface {
	spi.PlcSubscriber
	// Close stops every running poller, drops all subscriptions and consumer registrations and waits
	// for the spawned go routines to finish - which includes the consumer callbacks currently being
	// delivered (see fanOut). It is safe to call more than once and safe to call concurrently with
	// Subscribe.
	Close()
}

// PollTicker abstracts a time.Ticker so tests (and drivers bringing their own scheduler) can drive
// the poll cycle from a fake clock.
type PollTicker interface {
	// Chan returns the channel a tick is delivered on.
	Chan() <-chan time.Time
	// Stop releases the resources of this ticker. It doesn't have to close the channel.
	Stop()
}

// PollTickerFactory creates a PollTicker for the given poll interval.
type PollTickerFactory func(interval time.Duration) PollTicker

// PlcValueComparator decides if two apiValues.PlcValue are to be considered equal. It is used to
// suppress unchanged values on apiModel.SubscriptionChangeOfState subscriptions.
type PlcValueComparator func(previous, current apiValues.PlcValue) bool

// DefaultPollingInterval is used for tags which don't come with an interval of their own (which is
// the case for all apiModel.SubscriptionChangeOfState tags). Mirrors plc4j's
// PollingSubscriptionConnectionBase.getDefaultPollingInterval.
const DefaultPollingInterval = time.Second

// NewDefaultPollingSubscriber is the factory for a DefaultPollingSubscriber. Supported options are
// WithDefaultPollingInterval, WithPlcValueComparator, WithPollTickerFactory and the generic
// options.WithCustomLogger.
func NewDefaultPollingSubscriber(requirements DefaultPollingSubscriberRequirements, _options ...options.WithOption) DefaultPollingSubscriber {
	return buildDefaultPollingSubscriber(requirements, _options...)
}

// WithDefaultPollingInterval overrides the interval used for tags without an explicit one.
func WithDefaultPollingInterval(interval time.Duration) options.WithOption {
	return withDefaultPollingInterval{defaultPollingInterval: interval}
}

// WithPlcValueComparator overrides the change detection used for
// apiModel.SubscriptionChangeOfState subscriptions.
func WithPlcValueComparator(comparator PlcValueComparator) options.WithOption {
	return withPlcValueComparator{plcValueComparator: comparator}
}

// WithPollTickerFactory replaces the time.Ticker backed clock driving the poll cycles.
func WithPollTickerFactory(factory PollTickerFactory) options.WithOption {
	return withPollTickerFactory{pollTickerFactory: factory}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withDefaultPollingInterval struct {
	options.Option
	defaultPollingInterval time.Duration
}

type withPlcValueComparator struct {
	options.Option
	plcValueComparator PlcValueComparator
}

type withPollTickerFactory struct {
	options.Option
	pollTickerFactory PollTickerFactory
}

// realPollTicker is the production PollTicker, backed by a time.Ticker.
type realPollTicker struct {
	ticker *time.Ticker
}

func (r *realPollTicker) Chan() <-chan time.Time {
	return r.ticker.C
}

func (r *realPollTicker) Stop() {
	r.ticker.Stop()
}

func newRealPollTicker(interval time.Duration) PollTicker {
	return &realPollTicker{ticker: time.NewTicker(interval)}
}

// subscriptionKey groups the tags of one subscription request: all tags sharing a subscription type
// and a poll interval are served by a single poller (and therefore a single read request).
type subscriptionKey struct {
	subscriptionType apiModel.PlcSubscriptionType
	interval         time.Duration
}

// pollingSubscription is one running poller together with the tags it polls.
type pollingSubscription struct {
	id               uint64
	subscriptionType apiModel.PlcSubscriptionType
	interval         time.Duration

	// mutex guards tagNames, tags and previousValues: the poll go routine reads them while
	// Unsubscribe may concurrently drop tags from the group.
	mutex          sync.Mutex
	tagNames       []string
	tags           map[string]apiModel.PlcSubscriptionTag
	previousValues map[string]apiValues.PlcValue

	done     chan struct{}
	stopOnce sync.Once
}

// snapshot returns the currently subscribed tags in request order.
func (p *pollingSubscription) snapshot() ([]string, map[string]apiModel.PlcSubscriptionTag) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	tagNames := make([]string, len(p.tagNames))
	copy(tagNames, p.tagNames)
	tags := make(map[string]apiModel.PlcSubscriptionTag, len(p.tags))
	for tagName, tag := range p.tags {
		tags[tagName] = tag
	}
	return tagNames, tags
}

// tagCount returns how many tags are currently subscribed in this group.
func (p *pollingSubscription) tagCount() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return len(p.tagNames)
}

// hasTag tells if the given tag is still subscribed in this group.
func (p *pollingSubscription) hasTag(tagName string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.tags[tagName]
	return ok
}

// previousValue returns the value seen on the previous poll, if any.
func (p *pollingSubscription) previousValue(tagName string) (apiValues.PlcValue, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	value, ok := p.previousValues[tagName]
	return value, ok
}

// rememberValue stores the value of the current poll as the baseline for the next one.
func (p *pollingSubscription) rememberValue(tagName string, value apiValues.PlcValue) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.previousValues[tagName] = value
}

// removeTag drops one tag from this group and returns how many tags are left.
func (p *pollingSubscription) removeTag(tagName string) int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.tags[tagName]; !ok {
		return len(p.tagNames)
	}
	delete(p.tags, tagName)
	delete(p.previousValues, tagName)
	remaining := make([]string, 0, len(p.tagNames))
	for _, name := range p.tagNames {
		if name != tagName {
			remaining = append(remaining, name)
		}
	}
	p.tagNames = remaining
	return len(p.tagNames)
}

// stop signals the poll go routine to terminate. It is safe to call more than once.
func (p *pollingSubscription) stop() {
	p.stopOnce.Do(func() {
		close(p.done)
	})
}

// pollingSubscriptionHandle is handed out per subscribed tag. Embedding
// spiModel.DefaultPlcSubscriptionHandle makes it both registerable and unsubscribable through
// spiModel.DefaultPlcUnsubscriptionRequest.
type pollingSubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle
	subscriptionId uint64
	tagName        string
}

func newPollingSubscriptionHandle(subscriber spi.PlcSubscriber, subscriptionId uint64, tagName string) *pollingSubscriptionHandle {
	handle := &pollingSubscriptionHandle{
		subscriptionId: subscriptionId,
		tagName:        tagName,
	}
	handle.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(subscriber, handle)
	return handle
}

// pollingSubscriptionEvent adds the address information the generic
// spiModel.DefaultPlcSubscriptionEvent requires from its embedder.
type pollingSubscriptionEvent struct {
	*spiModel.DefaultPlcSubscriptionEvent
	addresses map[string]string
}

func newPollingSubscriptionEvent(
	tags map[string]apiModel.PlcTag,
	types map[string]apiModel.PlcSubscriptionType,
	intervals map[string]time.Duration,
	responseCodes map[string]apiModel.PlcResponseCode,
	addresses map[string]string,
	values map[string]apiValues.PlcValue,
	_options ...options.WithOption,
) apiModel.PlcSubscriptionEvent {
	event := &pollingSubscriptionEvent{addresses: addresses}
	event.DefaultPlcSubscriptionEvent = spiModel.NewDefaultPlcSubscriptionEvent(
		event, tags, types, intervals, responseCodes, values, _options...,
	).(*spiModel.DefaultPlcSubscriptionEvent)
	return event
}

// GetAddress returns the address the value was polled from.
func (p *pollingSubscriptionEvent) GetAddress(name string) string {
	return p.addresses[name]
}

type defaultPollingSubscriber struct {
	DefaultPollingSubscriberRequirements

	defaultPollingInterval time.Duration
	plcValueComparator     PlcValueComparator
	pollTickerFactory      PollTickerFactory

	subscriptionIdGenerator atomic.Uint64
	subscriptionsMutex      sync.Mutex
	subscriptions           map[uint64]*pollingSubscription

	consumersMutex sync.RWMutex
	consumers      map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer

	// ctx is cancelled by Close, which unblocks in-flight polling reads.
	ctx    context.Context
	cancel context.CancelFunc
	closed atomic.Bool

	// spawnMutex guards the "not closed yet, so spawn" decision of Subscribe against the closed flag
	// Close sets: adding to a sync.WaitGroup whose counter dropped to zero while somebody is inside
	// its Wait panics, and that panic happens on the caller's stack, outside of every recover this
	// file installs. Close only passes through this mutex as a fence (after it set closed), so it
	// never holds it while waiting - a consumer callback re-entering Subscribe gets the closed error
	// instead of deadlocking.
	spawnMutex sync.Mutex
	wg         sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func buildDefaultPollingSubscriber(requirements DefaultPollingSubscriberRequirements, _options ...options.WithOption) DefaultPollingSubscriber {
	defaultPollingInterval := DefaultPollingInterval
	var plcValueComparator PlcValueComparator
	var pollTickerFactory PollTickerFactory
	for _, option := range _options {
		switch option := option.(type) {
		case withDefaultPollingInterval:
			if option.defaultPollingInterval > 0 {
				defaultPollingInterval = option.defaultPollingInterval
			}
		case withPlcValueComparator:
			plcValueComparator = option.plcValueComparator
		case withPollTickerFactory:
			pollTickerFactory = option.pollTickerFactory
		}
	}
	if plcValueComparator == nil {
		plcValueComparator = PlcValuesEqual
	}
	if pollTickerFactory == nil {
		pollTickerFactory = newRealPollTicker
	}

	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	ctx, cancel := context.WithCancel(context.Background())
	return &defaultPollingSubscriber{
		DefaultPollingSubscriberRequirements: requirements,

		defaultPollingInterval: defaultPollingInterval,
		plcValueComparator:     plcValueComparator,
		pollTickerFactory:      pollTickerFactory,

		subscriptions: map[uint64]*pollingSubscription{},
		consumers:     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{},

		ctx:    ctx,
		cancel: cancel,

		log:      customLogger,
		_options: _options,
	}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

// PlcValuesEqual is the default PlcValueComparator: two values are equal when they are deeply equal.
// Everything reflect.DeepEqual can't tell apart is reported as *changed*, because for a
// apiModel.SubscriptionChangeOfState subscription a superfluous event is harmless while a swallowed
// one is not. Especially note that there is no String() based fallback: the String() of
// spi/values.PlcValueAdapter renders "not implemented" for every value type which doesn't override
// it (values.PlcRawByteArray, as produced by s7, knxnetip, bacnetip, cbus and opcua, is one of them),
// so such a fallback would consider any two of those values equal and silence the subscription for
// good.
func PlcValuesEqual(previous, current apiValues.PlcValue) bool {
	if previous == nil || current == nil {
		return previous == nil && current == nil
	}
	return reflect.DeepEqual(previous, current)
}

// Subscribe emulates the requested subscriptions by starting one poller per subscription type and
// interval combination. Note that the passed context only scopes the setup: the pollers outlive it
// and are bound to the lifetime of this subscriber (see Close).
func (d *defaultPollingSubscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	// The closed check and the spawning have to happen as one atomic step, otherwise a concurrent
	// Close could slip in between and the wg.Go would add to a wait group which is already being
	// drained - which panics the process.
	d.spawnMutex.Lock()
	if d.closed.Load() {
		d.spawnMutex.Unlock()
		utils.DeliverResult(d.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest, nil, errors.New("subscriber is closed"),
		))
		return result
	}
	d.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(d.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
					subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()),
				))
			}
		}()
		request, ok := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)
		if !ok {
			utils.DeliverResult(d.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
				subscriptionRequest, nil, errors.Errorf("%T is not a supported subscription request", subscriptionRequest),
			))
			return
		}

		responseCodes := map[string]apiModel.PlcResponseCode{}
		handles := map[string]apiModel.PlcSubscriptionHandle{}

		// Group the tags so tags sharing type and interval end up in one read request.
		var keyOrder []subscriptionKey
		groups := map[subscriptionKey][]string{}
		for _, tagName := range request.GetTagNames() {
			switch subscriptionType := request.GetType(tagName); subscriptionType {
			case apiModel.SubscriptionCyclic, apiModel.SubscriptionChangeOfState:
				tag := request.GetTag(tagName)
				if tag == nil {
					responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
					continue
				}
				interval := request.GetInterval(tagName)
				if interval <= 0 {
					interval = d.defaultPollingInterval
				}
				key := subscriptionKey{subscriptionType: subscriptionType, interval: interval}
				if _, alreadySeen := groups[key]; !alreadySeen {
					keyOrder = append(keyOrder, key)
				}
				groups[key] = append(groups[key], tagName)
			default:
				// Spontaneous events can't be emulated by polling.
				d.log.Debug().
					Str("tagName", tagName).
					Stringer("subscriptionType", subscriptionType).
					Msg("polling subscriptions can only emulate CYCLIC and CHANGE_OF_STATE")
				responseCodes[tagName] = apiModel.PlcResponseCode_UNSUPPORTED
			}
		}

		var startedSubscriptions []*pollingSubscription
		for _, key := range keyOrder {
			tagNames := groups[key]
			subscription := &pollingSubscription{
				id:               d.subscriptionIdGenerator.Add(1),
				subscriptionType: key.subscriptionType,
				interval:         key.interval,
				tagNames:         tagNames,
				tags:             make(map[string]apiModel.PlcSubscriptionTag, len(tagNames)),
				previousValues:   map[string]apiValues.PlcValue{},
				done:             make(chan struct{}),
			}
			for _, tagName := range tagNames {
				subscription.tags[tagName] = request.GetTag(tagName)
			}
			if !d.addSubscription(subscription) {
				for _, tagName := range tagNames {
					responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
				}
				continue
			}
			for _, tagName := range tagNames {
				handles[tagName] = newPollingSubscriptionHandle(d, subscription.id, tagName)
				responseCodes[tagName] = apiModel.PlcResponseCode_OK
			}
			startedSubscriptions = append(startedSubscriptions, subscription)
		}

		// The response has to be assembled *before* the first poll can happen: building it is what
		// hooks the consumers pre-registered on the request up to their handles, and because the very
		// first poll of a apiModel.SubscriptionChangeOfState subscription establishes the baseline
		// (whether anybody listens or not), an event lost to a consumer which isn't wired up yet is
		// lost until the value changes again.
		response := spiModel.NewDefaultPlcSubscriptionResponse(
			subscriptionRequest,
			responseCodes,
			handles,
			append(d._options, options.WithCustomLogger(d.log))...,
		)

		for _, subscription := range startedSubscriptions {
			d.log.Debug().
				Uint64("subscriptionId", subscription.id).
				Stringer("subscriptionType", subscription.subscriptionType).
				Dur("interval", subscription.interval).
				Int("numTags", subscription.tagCount()).
				Msg("starting polling subscription")
			d.startPoller(subscription)
		}

		utils.DeliverResult(d.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest, response, nil,
		))
	})
	d.spawnMutex.Unlock()
	return result
}

// Unsubscribe drops the passed handles' tags from their pollers, stopping a poller as soon as none
// of its tags are subscribed anymore. Handles which don't belong to this subscriber are skipped:
// spiModel.DefaultPlcUnsubscriptionRequest hands the complete request to every involved subscriber,
// so a request may well carry foreign handles.
func (d *defaultPollingSubscriber) Unsubscribe(_ context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	handleProvider, ok := unsubscriptionRequest.(interface {
		GetSubscriptionHandles() []apiModel.PlcSubscriptionHandle
	})
	if !ok {
		utils.DeliverResult(d.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(
			unsubscriptionRequest, nil, errors.Errorf("%T is not a supported unsubscription request", unsubscriptionRequest),
		))
		return result
	}
	touchedOurs := false
	for _, handle := range handleProvider.GetSubscriptionHandles() {
		pollingHandle, ok := handle.(*pollingSubscriptionHandle)
		if !ok {
			d.log.Trace().Type("handle", handle).Msg("not a handle of this subscriber")
			continue
		}
		d.unsubscribeTag(pollingHandle)
		touchedOurs = true
	}
	if touchedOurs {
		d.purgeStaleRegistrations()
	}
	utils.DeliverResult(d.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(
		unsubscriptionRequest,
		spiModel.NewDefaultPlcUnsubscriptionResponse(unsubscriptionRequest),
		nil,
	))
	return result
}

// Register registers a consumer for the given handles. Registering on a closed subscriber hands out
// an inert registration: there is no poller left which could feed it.
func (d *defaultPollingSubscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(d, consumer, handles...)
	d.consumersMutex.Lock()
	defer d.consumersMutex.Unlock()
	if d.closed.Load() {
		// Don't repopulate the map Close cleared (or is about to clear).
		d.log.Debug().Msg("subscriber is closed, consumer not registered")
		return consumerRegistration
	}
	d.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

// Unregister removes a previously registered consumer.
func (d *defaultPollingSubscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	defaultRegistration, ok := registration.(*spiModel.DefaultPlcConsumerRegistration)
	if !ok {
		d.log.Debug().Type("registration", registration).Msg("not a registration of this subscriber")
		return
	}
	d.consumersMutex.Lock()
	defer d.consumersMutex.Unlock()
	delete(d.consumers, defaultRegistration)
}

// Close stops all pollers and waits for the spawned go routines to finish.
func (d *defaultPollingSubscriber) Close() {
	if !d.closed.CompareAndSwap(false, true) {
		return
	}
	// Unblocks polling reads which are currently in flight.
	d.cancel()
	d.subscriptionsMutex.Lock()
	for subscriptionId, subscription := range d.subscriptions {
		subscription.stop()
		delete(d.subscriptions, subscriptionId)
	}
	d.subscriptionsMutex.Unlock()
	// Passing through the spawn mutex once is enough to fence off Subscribe: closed is already set,
	// and every Subscribe checks it while holding this very mutex. So a Subscribe which is past the
	// check has finished spawning by the time we get the mutex, and one arriving later sees the
	// subscriber as closed and doesn't spawn at all. The go routines we are about to wait for are
	// free to spawn pollers of their own: their own wait group entry keeps the counter above zero,
	// which is what makes such an Add legal.
	d.spawnMutex.Lock()
	//nolint:staticcheck // SA2001: the empty critical section is the point, see above.
	d.spawnMutex.Unlock()
	d.wg.Wait()
	d.consumersMutex.Lock()
	d.consumers = map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{}
	d.consumersMutex.Unlock()
	d.log.Trace().Msg("polling subscriber closed")
}

// addSubscription registers a subscription unless the subscriber is already closing.
func (d *defaultPollingSubscriber) addSubscription(subscription *pollingSubscription) bool {
	d.subscriptionsMutex.Lock()
	defer d.subscriptionsMutex.Unlock()
	if d.closed.Load() {
		return false
	}
	d.subscriptions[subscription.id] = subscription
	return true
}

// unsubscribeTag removes one tag from its poller and stops the poller when it runs out of tags.
func (d *defaultPollingSubscriber) unsubscribeTag(handle *pollingSubscriptionHandle) {
	d.subscriptionsMutex.Lock()
	subscription, ok := d.subscriptions[handle.subscriptionId]
	if !ok {
		d.subscriptionsMutex.Unlock()
		// Either already unsubscribed or a sibling handle of the same group did it.
		d.log.Trace().Uint64("subscriptionId", handle.subscriptionId).Msg("subscription already gone")
		return
	}
	if subscription.removeTag(handle.tagName) == 0 {
		delete(d.subscriptions, handle.subscriptionId)
		d.subscriptionsMutex.Unlock()
		subscription.stop()
		d.log.Debug().Uint64("subscriptionId", handle.subscriptionId).Msg("stopped polling subscription")
		return
	}
	d.subscriptionsMutex.Unlock()
}

// purgeStaleRegistrations drops the consumer registrations which reference nothing but handles of
// tags this subscriber doesn't poll anymore. Without it a long living connection which repeatedly
// subscribes and unsubscribes piles up dead registrations, which fanOut then walks on every poll.
// Registrations carrying at least one live or foreign handle are left alone.
func (d *defaultPollingSubscriber) purgeStaleRegistrations() {
	d.consumersMutex.Lock()
	defer d.consumersMutex.Unlock()
	for registration := range d.consumers {
		handles := registration.GetSubscriptionHandles()
		if len(handles) == 0 {
			continue
		}
		stale := true
		for _, handle := range handles {
			pollingHandle, ok := handle.(*pollingSubscriptionHandle)
			if !ok || d.isSubscribed(pollingHandle) {
				stale = false
				break
			}
		}
		if stale {
			delete(d.consumers, registration)
			d.log.Trace().Msg("dropped a consumer registration of an unsubscribed tag")
		}
	}
}

// isSubscribed tells if the tag a handle points to is still being polled.
func (d *defaultPollingSubscriber) isSubscribed(handle *pollingSubscriptionHandle) bool {
	d.subscriptionsMutex.Lock()
	subscription, ok := d.subscriptions[handle.subscriptionId]
	d.subscriptionsMutex.Unlock()
	return ok && subscription.hasTag(handle.tagName)
}

// startPoller spawns the go routine driving one subscription.
func (d *defaultPollingSubscriber) startPoller(subscription *pollingSubscription) {
	ticker := d.pollTickerFactory(subscription.interval)
	d.wg.Go(func() {
		defer ticker.Stop()
		tickChan := ticker.Chan()
		for {
			select {
			case <-d.ctx.Done():
				return
			case <-subscription.done:
				return
			case _, ok := <-tickChan:
				if !ok {
					return
				}
				d.poll(subscription)
			}
		}
	})
}

// poll issues one read request for all tags of a subscription and turns the response into an event.
func (d *defaultPollingSubscriber) poll(subscription *pollingSubscription) {
	defer func() {
		if err := recover(); err != nil {
			d.log.Error().
				Str("stack", string(debug.Stack())).
				Interface("err", err).
				Uint64("subscriptionId", subscription.id).
				Msg("panic-ed while polling")
		}
	}()
	tagNames, tags := subscription.snapshot()
	if len(tagNames) == 0 {
		return
	}
	builder := d.PollingReadRequestBuilder()
	if builder == nil {
		d.log.Error().Msg("driver didn't supply a read request builder")
		return
	}
	for _, tagName := range tagNames {
		builder.AddTagAddress(tagName, tags[tagName].GetAddressString())
	}
	readRequest, err := builder.Build()
	if err != nil {
		d.log.Warn().Err(err).Uint64("subscriptionId", subscription.id).Msg("error building the polling read request")
		return
	}

	var readResult apiModel.PlcReadRequestResult
	select {
	case readResult = <-readRequest.Execute(d.ctx):
	case <-d.ctx.Done():
		return
	case <-subscription.done:
		return
	}
	if readResult == nil {
		d.log.Warn().Uint64("subscriptionId", subscription.id).Msg("no result from the polling read request")
		return
	}
	if err := readResult.GetErr(); err != nil {
		d.log.Warn().Err(err).Uint64("subscriptionId", subscription.id).Msg("error executing the polling read request")
		return
	}
	readResponse := readResult.GetResponse()
	if readResponse == nil {
		d.log.Warn().Uint64("subscriptionId", subscription.id).Msg("no response from the polling read request")
		return
	}
	d.emit(subscription, tagNames, tags, readResponse)
}

// emit assembles the subscription event for one poll and hands it to the interested consumers.
// For apiModel.SubscriptionChangeOfState only tags whose value actually changed are part of the
// event, and an event consisting solely of failed tags is not emitted at all (matching plc4j).
func (d *defaultPollingSubscriber) emit(
	subscription *pollingSubscription,
	tagNames []string,
	tags map[string]apiModel.PlcSubscriptionTag,
	readResponse apiModel.PlcReadResponse,
) {
	eventTags := map[string]apiModel.PlcTag{}
	eventTypes := map[string]apiModel.PlcSubscriptionType{}
	eventIntervals := map[string]time.Duration{}
	eventResponseCodes := map[string]apiModel.PlcResponseCode{}
	eventAddresses := map[string]string{}
	eventValues := map[string]apiValues.PlcValue{}

	addTag := func(tagName string, code apiModel.PlcResponseCode) {
		eventTags[tagName] = tags[tagName]
		eventTypes[tagName] = subscription.subscriptionType
		eventIntervals[tagName] = subscription.interval
		eventResponseCodes[tagName] = code
		eventAddresses[tagName] = tags[tagName].GetAddressString()
	}

	shouldFire := false
	for _, tagName := range tagNames {
		code := readResponse.GetResponseCode(tagName)
		if code != apiModel.PlcResponseCode_OK {
			// Failures are reported, but on their own they don't make an event.
			addTag(tagName, code)
			continue
		}
		currentValue := readResponse.GetValue(tagName)
		if subscription.subscriptionType == apiModel.SubscriptionChangeOfState {
			if previousValue, seen := subscription.previousValue(tagName); seen && d.plcValueComparator(previousValue, currentValue) {
				continue
			}
		}
		subscription.rememberValue(tagName, currentValue)
		addTag(tagName, code)
		eventValues[tagName] = currentValue
		shouldFire = true
	}
	if !shouldFire {
		return
	}

	event := newPollingSubscriptionEvent(
		eventTags,
		eventTypes,
		eventIntervals,
		eventResponseCodes,
		eventAddresses,
		eventValues,
		append(d._options, options.WithCustomLogger(d.log))...,
	)
	d.fanOut(subscription.id, event)
}

// fanOut delivers an event to every consumer registered for a handle of the given subscription. The
// consumers are snapshotted so a consumer is free to (un)register from within its own callback.
//
// Delivery is synchronous on the poll go routine, which keeps the events of one subscription ordered
// (a consumer sees them in the order they were polled). The price is that a consumer which blocks
// stalls the polling of its own subscription - and makes Close block until the callback returns,
// since Close waits for the poll go routines. Consumers doing more than a hand-full of work should
// therefore queue the event and process it elsewhere.
func (d *defaultPollingSubscriber) fanOut(subscriptionId uint64, event apiModel.PlcSubscriptionEvent) {
	d.consumersMutex.RLock()
	interested := make([]apiModel.PlcSubscriptionEventConsumer, 0, len(d.consumers))
	for registration, consumer := range d.consumers {
		for _, handle := range registration.GetSubscriptionHandles() {
			if pollingHandle, ok := handle.(*pollingSubscriptionHandle); ok && pollingHandle.subscriptionId == subscriptionId {
				interested = append(interested, consumer)
				break
			}
		}
	}
	d.consumersMutex.RUnlock()
	for _, consumer := range interested {
		if consumer == nil {
			continue
		}
		consumer(event)
	}
}
