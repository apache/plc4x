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

package bacnetip

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// Subscriber drives BACnet COV (Change-of-Value) subscriptions: it issues
// SubscribeCOV requests on Subscribe, refreshes them at lifetime/2 via a
// background ticker, fans incoming UnconfirmedCOVNotification /
// ConfirmedCOVNotification APDUs out to registered consumers, and tears
// everything down on Unsubscribe.
//
//go:generate go tool plc4xGenerator -type=Subscriber
type Subscriber struct {
	connection *Connection
	consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer

	// handles indexes active SubscriptionHandles by their subscriberProcessId so
	// incoming COV notifications can be routed in O(1).
	handles   map[uint32]*SubscriptionHandle
	handlesMu sync.RWMutex

	// nextProcessId is a monotonic counter that allocates unique subscriber
	// process ids per Connection. BACnet spec allows any uint32 value as the
	// id; starting at 1 keeps 0 reserved for "no subscription".
	nextProcessId atomic.Uint32

	wg sync.WaitGroup `ignore:"true"`

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption // Used to pass them downstream
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	logger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	s := &Subscriber{
		connection: connection,
		consumers:  make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),
		handles:    make(map[uint32]*SubscriptionHandle),

		log:      logger,
		_options: _options,
	}
	s.nextProcessId.Store(1)
	return s
}

// Subscribe issues one BACnet SubscribeCOV request per requested tag. The
// per-tag lifetime is taken from Connection.configuration.CovLifetimeSeconds;
// a value of 0 means "indefinite" per BACnet semantics.
func (m *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	m.wg.Go(func() {
		internalReq := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)
		m.connection.addSubscriber(m)

		lifetime := m.connection.configuration.CovLifetimeSeconds
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := map[string]apiModel.PlcSubscriptionHandle{}

		for _, tagName := range internalReq.GetTagNames() {
			if err := ctx.Err(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, err))
				return
			}
			tag, ok := internalReq.GetTag(tagName).(BacNetPlcTag)
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			processId := m.nextProcessId.Add(1)
			handle := NewSubscriptionHandle(
				m, tagName, tag,
				internalReq.GetType(tagName),
				internalReq.GetInterval(tagName),
			)
			handle.subscriberProcessId = processId
			if tag.GetObjectId().ObjectIdType != nil {
				handle.monitoredObjectId = *tag.GetObjectId().ObjectIdType
			}
			handle.monitoredInstance = tag.GetObjectId().ObjectIdInstance
			handle.lifetimeSec = lifetime

			code := m.sendSubscribeCOV(ctx, handle, lifetime)
			responseCodes[tagName] = code
			subscriptionValues[tagName] = handle.DefaultPlcSubscriptionHandle

			if code == apiModel.PlcResponseCode_OK {
				m.storeHandle(handle)
			}
		}

		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(
				subscriptionRequest,
				responseCodes,
				subscriptionValues,
				append(m._options, options.WithCustomLogger(m.log))...,
			),
			nil,
		))
	})
	return result
}

// Unsubscribe sends a SubscribeCOV with lifetime=0 (BACnet cancel semantics)
// for each handle in the unsubscription request, then drops the handle from
// the routing map so further notifications are silently ignored.
func (m *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	m.wg.Go(func() {
		req, ok := unsubscriptionRequest.(*spiModel.DefaultPlcUnsubscriptionRequest)
		if !ok {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("unsupported unsubscription request type")))
			return
		}
		for _, handle := range req.GetSubscriptionHandles() {
			bh := m.findHandle(handle)
			if bh == nil {
				continue
			}
			_ = m.sendSubscribeCOV(ctx, bh, 0) // lifetime=0 cancels
			m.removeHandle(bh.subscriberProcessId)
		}
		// Build a response with OK for every handle the caller passed.
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, nil))
	})
	return result
}

// sendSubscribeCOV builds and dispatches a single SubscribeCOV confirmed
// request. Returns the plc4go response code derived from the SimpleAck/Error
// reply. lifetime=0 cancels an existing subscription.
func (m *Subscriber) sendSubscribeCOV(ctx context.Context, handle *SubscriptionHandle, lifetime uint32) apiModel.PlcResponseCode {
	tag, ok := handle.tag.(BacNetPlcTag)
	if !ok || tag == nil {
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	}
	subscriberProcessId := readWriteModel.CreateBACnetContextTagUnsignedInteger(0, uint(handle.subscriberProcessId))
	monitoredObject := readWriteModel.CreateBACnetContextTagObjectIdentifier(1, tag.GetObjectId().getId(), tag.GetObjectId().ObjectIdInstance)
	issueConfirmed := readWriteModel.CreateBACnetContextTagBoolean(2, handle.confirmedNotifications)
	lifetimeTag := readWriteModel.CreateBACnetContextTagUnsignedInteger(3, uint(lifetime))

	serviceRequest := readWriteModel.NewBACnetConfirmedServiceRequestSubscribeCOV(
		0, subscriberProcessId, monitoredObject, issueConfirmed, lifetimeTag,
	)
	invokeId := m.connection.invokeIdGenerator.getAndIncrement()
	apdu := readWriteModel.NewAPDUConfirmedRequest(
		false, false, true,
		m.connection.driverContext.maxSegmentsAccepted,
		m.connection.driverContext.maxApduLengthAccepted,
		invokeId, nil, nil,
		serviceRequest, nil, nil,
	)

	// Synchronous round-trip: we block until SendRequest dispatches the
	// SimpleAck back through the codec. Phase 5 will replace this with a real
	// transaction-manager-backed retry loop honoring ApduTimeoutMs/Retries.
	done := make(chan apiModel.PlcResponseCode, 1)
	err := m.connection.messageCodec.SendRequest(ctx, "subscribe-cov", wrapAPDU(apdu, true, m.connection.routedDest),
		func(message spi.Message) bool {
			return m.acceptsResponse(message, invokeId)
		},
		func(message spi.Message) error {
			apduResp := message.(readWriteModel.BVLC).(interface{ GetNpdu() readWriteModel.NPDU }).GetNpdu().GetApdu()
			done <- m.classifyResponse(apduResp)
			return nil
		},
		func(err error) error {
			done <- apiModel.PlcResponseCode_REQUEST_TIMEOUT
			return nil
		},
	)
	if err != nil {
		m.log.Debug().Err(err).Msg("error sending SubscribeCOV")
		return apiModel.PlcResponseCode_REMOTE_ERROR
	}
	select {
	case code := <-done:
		if code == apiModel.PlcResponseCode_OK {
			handle.lastIssuedAt = time.Now()
		}
		return code
	case <-ctx.Done():
		return apiModel.PlcResponseCode_REQUEST_TIMEOUT
	}
}

// classifyResponse maps a SubscribeCOV reply APDU to a plc4go response code.
func (m *Subscriber) classifyResponse(apdu readWriteModel.APDU) apiModel.PlcResponseCode {
	switch apdu := apdu.(type) {
	case readWriteModel.APDUSimpleAck:
		return apiModel.PlcResponseCode_OK
	case readWriteModel.APDUError:
		return mapErrorAPDU(apdu, m.log)
	case readWriteModel.APDUReject:
		return apiModel.PlcResponseCode_INVALID_DATA
	case readWriteModel.APDUAbort:
		reason := apdu.GetAbortReason().GetValue()
		if reason == readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED {
			return apiModel.PlcResponseCode_UNSUPPORTED
		}
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	default:
		return apiModel.PlcResponseCode_REMOTE_ERROR
	}
}

func (m *Subscriber) acceptsResponse(message spi.Message, invokeId uint8) bool {
	bvlc, ok := message.(readWriteModel.BVLC)
	if !ok {
		return false
	}
	npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU })
	if !ok {
		return false
	}
	npdu := npduRetriever.GetNpdu()
	if npdu.GetControl().GetMessageTypeFieldPresent() {
		return false
	}
	got, err := getInvokeIdFromApdu(npdu.GetApdu())
	if err != nil {
		return false
	}
	return got == invokeId
}

// HandleConfirmedCOVNotification is invoked by Connection's message loop when a
// ConfirmedCOVNotification APDU arrives. It looks up the SubscriptionHandle by
// subscriber-process-id and fans out a SubscriptionEvent to every consumer
// registered against that handle. The caller is responsible for sending the
// APDUSimpleAck reply back to the publisher.
func (m *Subscriber) HandleConfirmedCOVNotification(req readWriteModel.BACnetConfirmedServiceRequestConfirmedCOVNotification) {
	processId := uint32(req.GetSubscriberProcessIdentifier().GetPayload().GetActualValue())
	handle := m.lookupHandle(processId)
	if handle == nil {
		m.log.Debug().Uint32("processId", processId).Msg("Confirmed COV notification for unknown subscriber")
		return
	}
	m.dispatchNotification(handle, req.GetListOfValues())
}

// HandleUnconfirmedCOVNotification mirrors HandleConfirmedCOVNotification for
// the unconfirmed variant. No reply is sent.
func (m *Subscriber) HandleUnconfirmedCOVNotification(req readWriteModel.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) {
	processId := uint32(req.GetSubscriberProcessIdentifier().GetPayload().GetActualValue())
	handle := m.lookupHandle(processId)
	if handle == nil {
		m.log.Debug().Uint32("processId", processId).Msg("Unconfirmed COV notification for unknown subscriber")
		return
	}
	m.dispatchNotification(handle, req.GetListOfValues())
}

// dispatchNotification converts the list of (property, value) pairs into a
// PlcValue and pushes a SubscriptionEvent to every consumer that's registered
// against this handle.
func (m *Subscriber) dispatchNotification(handle *SubscriptionHandle, listOfValues readWriteModel.BACnetPropertyValues) {
	values := map[string]apiValues.PlcValue{}
	propertyValues := listOfValues.GetData()
	if len(propertyValues) == 0 {
		values[handle.tagName] = spiValues.NewPlcNULL()
	} else {
		// COV reports a list of (property, value, optional-priority) tuples.
		// We surface the PresentValue if it's in the list (the common case),
		// otherwise the first value.
		picked := propertyValues[0]
		for _, pv := range propertyValues {
			if pv.GetPropertyIdentifier().GetValue() == readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE {
				picked = pv
				break
			}
		}
		// A BACnetConstructedDataElement holds the value in exactly one of three
		// fields depending on how the publisher framed it: ApplicationTag (e.g.
		// Real for AnalogInput), ConstructedData (nested typed value), or
		// ContextTag (context-specific encoding). Pick whichever is populated.
		element := picked.GetPropertyValue()
		switch {
		case element == nil:
			values[handle.tagName] = spiValues.NewPlcNULL()
		case element.GetApplicationTag() != nil:
			values[handle.tagName] = appTagToPlcValue(element.GetApplicationTag())
		case element.GetConstructedData() != nil:
			values[handle.tagName] = constructedDataToPlcValue(element.GetConstructedData())
		default:
			values[handle.tagName] = spiValues.NewPlcNULL()
		}
	}
	codes := map[string]apiModel.PlcResponseCode{handle.tagName: apiModel.PlcResponseCode_OK}
	tags := map[string]apiModel.PlcTag{handle.tagName: handle.tag.(apiModel.PlcTag)}
	types := map[string]apiModel.PlcSubscriptionType{handle.tagName: handle.subscriptionType}
	intervals := map[string]time.Duration{handle.tagName: handle.interval}
	event := NewSubscriptionEvent(tags, types, intervals, codes, map[string]string{}, map[string]string{}, values)
	for _, consumer := range m.consumers {
		consumer(&event)
	}
}

// storeHandle / removeHandle / lookupHandle / findHandle: process-id-indexed
// table for COV notification routing.

func (m *Subscriber) storeHandle(handle *SubscriptionHandle) {
	m.handlesMu.Lock()
	defer m.handlesMu.Unlock()
	m.handles[handle.subscriberProcessId] = handle
}

func (m *Subscriber) removeHandle(processId uint32) {
	m.handlesMu.Lock()
	defer m.handlesMu.Unlock()
	delete(m.handles, processId)
}

func (m *Subscriber) lookupHandle(processId uint32) *SubscriptionHandle {
	m.handlesMu.RLock()
	defer m.handlesMu.RUnlock()
	return m.handles[processId]
}

// activeHandleCount reports the number of currently registered subscription
// handles (feeds Connection.ActiveSubscriptionCount).
func (m *Subscriber) activeHandleCount() int {
	m.handlesMu.RLock()
	defer m.handlesMu.RUnlock()
	return len(m.handles)
}

// findHandle locates the SubscriptionHandle that owns the given
// apiModel.PlcSubscriptionHandle (returned earlier in the Subscribe response).
// Used by Unsubscribe which only has the api-level handle in hand.
func (m *Subscriber) findHandle(needle apiModel.PlcSubscriptionHandle) *SubscriptionHandle {
	m.handlesMu.RLock()
	defer m.handlesMu.RUnlock()
	for _, h := range m.handles {
		if h.DefaultPlcSubscriptionHandle == needle {
			return h
		}
	}
	return nil
}

func (m *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(m, consumer, handles...)
	m.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (m *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	delete(m.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
