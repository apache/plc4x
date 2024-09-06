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
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

//go:generate plc4xGenerator -type=Subscriber
type Subscriber struct {
	connection *Connection
	consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer

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
		passLogToModel: passLoggerToModel,
		log:            customLogger,
		_options:       _options,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
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
			subscriptionValues[tagName] = NewSubscriptionHandle(s, tagName, internalPlcSubscriptionRequest.GetTag(tagName), tagType, internalPlcSubscriptionRequest.GetInterval(tagName))
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
	}()
	return result
}

func (s *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("Not Implemented"))
	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more an internal handling of which values to pass where.

	return result
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
	for registration, consumer := range s.consumers {
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			subscriptionHandle := subscriptionHandle.(*SubscriptionHandle)
			groupAddressTag, ok := subscriptionHandle.tag.(GroupAddressTag)
			if !ok || !groupAddressTag.matches(groupAddress) {
				continue
			}
			if subscriptionHandle.tagType != apiModel.SubscriptionChangeOfState || !changed {
				continue
			}
			tags := map[string]apiModel.PlcTag{}
			types := map[string]apiModel.PlcSubscriptionType{}
			intervals := map[string]time.Duration{}
			responseCodes := map[string]apiModel.PlcResponseCode{}
			addresses := map[string][]byte{}
			plcValues := map[string]values.PlcValue{}
			tagName := subscriptionHandle.tagName
			rb := utils.NewReadBufferByteBased(payload)
			if groupAddressTag.GetTagType() == nil {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_DATATYPE
				plcValues[tagName] = nil
				continue
			}
			// If the size of the tag is greater than 6, we have to skip the first byte
			if groupAddressTag.GetTagType().GetLengthInBits(context.Background()) > 6 {
				_, _ = rb.ReadUint8("groupAddress", 8)
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
					plcValue, err2 := driverModel.KnxDatapointParseWithBuffer(context.Background(), rb, elementType)
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
	s.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (s *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	delete(s.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
