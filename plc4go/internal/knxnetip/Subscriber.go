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

package knxnetip

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

//go:generate go tool plc4xGenerator -type=Subscriber
type Subscriber struct {
	connection *Connection
	consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	// handles contains all subscription handles which are still active. KNX has no
	// wire-level subscribe, the handles are a purely local filter, so unsubscribing
	// simply means dropping the handle here.
	handles map[*SubscriptionHandle]struct{} `ignore:"true"`
	// consumersMutex guards consumers and handles, which are written by Subscribe,
	// Unsubscribe, Register and Unregister while being read from the codec worker
	// (via Connection.handleValueCacheUpdate).
	consumersMutex sync.RWMutex

	wg sync.WaitGroup // use to track spawned go routines

	passLogToModel bool
	log            zerolog.Logger       `ignore:"true"`
	_options       []options.WithOption // Used to pass them downstream
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Subscriber{
		connection:     connection,
		consumers:      make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),
		handles:        make(map[*SubscriptionHandle]struct{}),
		passLogToModel: passLoggerToModel,
		log:            customLogger,
		_options:       _options,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	s.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		internalPlcSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

		// Add this subscriber to the connection.
		s.connection.addSubscriber(s)

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
		for _, tagName := range internalPlcSubscriptionRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			tagType := internalPlcSubscriptionRequest.GetType(tagName)
			subscriptionHandle := NewSubscriptionHandle(s, tagName, internalPlcSubscriptionRequest.GetTag(tagName), tagType, internalPlcSubscriptionRequest.GetInterval(tagName))
			subscriptionValues[tagName] = subscriptionHandle
			s.consumersMutex.Lock()
			s.handles[subscriptionHandle] = struct{}{}
			s.consumersMutex.Unlock()
		}

		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(
				subscriptionRequest,
				responseCodes,
				subscriptionValues,
				append(s._options, options.WithCustomLogger(s.log))...,
			),
			nil,
		))
	})
	return result
}

// Unsubscribe deregisters the handles of the given request. As soon as we establish a
// connection we start getting data, subscriptions are just an internal handling of which
// values to pass where, so unsubscribing means dropping the handles (and the consumer
// registrations which only referenced them) again.
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
		err := s.unsubscribeHandles(request.GetSubscriptionHandles())
		utils.DeliverResult(s.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(
			unsubscriptionRequest,
			spiModel.NewDefaultPlcUnsubscriptionResponse(unsubscriptionRequest),
			err,
		))
	})
	return result
}

// unsubscribeHandles drops the given handles and all consumer registrations which are
// left without a single active handle. Handles which don't belong to this subscriber are
// silently ignored as an unsubscription request can span multiple subscribers.
func (s *Subscriber) unsubscribeHandles(subscriptionHandles []apiModel.PlcSubscriptionHandle) error {
	var collectedErrors []error
	s.consumersMutex.Lock()
	for _, handle := range subscriptionHandles {
		subscriptionHandle, ok := handle.(*SubscriptionHandle)
		if !ok {
			collectedErrors = append(collectedErrors, errors.Errorf("%T is not a knx subscription handle", handle))
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
	return errors.Join(collectedErrors...)
}

// activeRegistration is a snapshot of one consumer registration reduced to the handles
// which are still subscribed.
type activeRegistration struct {
	consumer apiModel.PlcSubscriptionEventConsumer
	handles  []*SubscriptionHandle
}

// activeRegistrations snapshots the current registrations so events can be delivered
// without holding the lock while calling into the (possibly re-entrant) consumers.
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

// deliversFor decides if a handle of the given subscription-type wants to see the current
// group-value write. KNX is event driven: change-of-state handles only get fed if the
// value actually changed, event handles get every write.
// (Java: KnxNetIpConnection advertises PlcSubscriptionType.EVENT)
func deliversFor(tagType apiModel.PlcSubscriptionType, changed bool) bool {
	switch tagType {
	case apiModel.SubscriptionChangeOfState:
		return changed
	case apiModel.SubscriptionEvent:
		return true
	default:
		// Cyclic subscriptions are not supported by the knx driver.
		return false
	}
}

/*
 * Callback for incoming value change events from the KNX bus
 */
func (s *Subscriber) handleValueChange(ctx context.Context, destinationAddress []byte, payload []byte, changed bool) {
	// Decode the group-address according to the settings in the driver
	// Group addresses can be 1, 2 or 3 levels (3 being the default)
	ctxForModel := options.GetLoggerContextForModel(ctx, s.log, options.WithPassLoggerToModel(s.passLogToModel))
	groupAddress, err := driverModel.KnxGroupAddressParse[driverModel.KnxGroupAddress](ctxForModel, destinationAddress, s.connection.getGroupAddressNumLevels())
	if err != nil {
		return
	}

	// TODO: aggregate tags and send it to a consumer which want's all of them
	for _, registration := range s.activeRegistrations() {
		consumer := registration.consumer
		for _, subscriptionHandle := range registration.handles {
			groupAddressTag, ok := subscriptionHandle.tag.(GroupAddressTag)
			if !ok || !groupAddressTag.matches(groupAddress) {
				continue
			}
			if !deliversFor(subscriptionHandle.tagType, changed) {
				continue
			}
			tags := map[string]apiModel.PlcTag{}
			types := map[string]apiModel.PlcSubscriptionType{}
			intervals := map[string]time.Duration{}
			responseCodes := map[string]apiModel.PlcResponseCode{}
			addresses := map[string][]byte{}
			plcValues := map[string]values.PlcValue{}
			tagName := subscriptionHandle.tagName
			// The payload is the raw group-value payload: the byte carrying the 6 embedded
			// data bits followed by whatever data bytes came after it. The generated
			// datapoint parser already accounts for that layout (it reads the reserved
			// bits/byte itself), so it gets the payload as-is - skipping anything here
			// would eat the value of a small datapoint-type and shift a bigger one.
			// (Java: KnxNetIpConnection hands the payload to KnxDatapoint.staticParse unchanged)
			rb := utils.NewReadBufferByteBased(payload)
			if groupAddressTag.GetTagType() == nil {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_DATATYPE
				plcValues[tagName] = nil
				continue
			}
			elementType := *groupAddressTag.GetTagType()
			numElements := uint16(1)
			if len(groupAddressTag.GetArrayInfo()) > 0 {
				numElements = uint16(groupAddressTag.GetArrayInfo()[0].GetUpperBound() - groupAddressTag.GetArrayInfo()[0].GetLowerBound())
			}

			tags[tagName] = groupAddressTag
			types[tagName] = subscriptionHandle.tagType
			intervals[tagName] = subscriptionHandle.interval
			addresses[tagName] = destinationAddress

			var plcValueList []values.PlcValue
			responseCode := apiModel.PlcResponseCode_OK
			for i := uint16(0); i < numElements; i++ {
				// If we don't know the datatype, we'll create a RawPlcValue instead
				// so the application can decode the content later on.
				if elementType == driverModel.KnxDatapointType_DPT_UNKNOWN {
					// If this is an unknown 1 byte payload, we need the first byte.
					if !rb.HasMore(1) {
						rb.Reset(0)
					}
					plcValue := spiValues.NewPlcRawByteArray(rb.GetBytes())
					plcValueList = append(plcValueList, plcValue)
				} else {
					plcValue, err2 := driverModel.KnxDatapointParseWithBuffer(ctx, rb, elementType)
					if err2 == nil {
						plcValueList = append(plcValueList, plcValue)
					} else {
						// TODO: Do a little more here ...
						responseCode = apiModel.PlcResponseCode_INTERNAL_ERROR
						break
					}
				}
			}
			responseCodes[tagName] = responseCode
			if responseCode == apiModel.PlcResponseCode_OK {
				if len(plcValueList) == 1 {
					plcValues[tagName] = plcValueList[0]
				} else {
					plcValues[tagName] = spiValues.NewPlcList(plcValueList)
				}
			}
			event := NewSubscriptionEvent(
				tags,
				types,
				intervals,
				responseCodes,
				addresses,
				plcValues,
				append(s._options, options.WithCustomLogger(s.log))...,
			)
			consumer(&event)
		}
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
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	delete(s.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
