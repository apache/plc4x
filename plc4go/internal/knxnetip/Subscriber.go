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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/spi/values"
	"time"
)

type Subscriber struct {
	connection           *Connection
	subscriptionRequests []internalModel.DefaultPlcSubscriptionRequest
}

func NewSubscriber(connection *Connection) *Subscriber {
	return &Subscriber{
		connection:           connection,
		subscriptionRequests: []internalModel.DefaultPlcSubscriptionRequest{},
	}
}

func (m *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcSubscriptionRequestResult)
	go func() {
		// Add this subscriber to the connection.
		m.connection.addSubscriber(m)

		// Save the subscription request
		m.subscriptionRequests = append(m.subscriptionRequests, subscriptionRequest.(internalModel.DefaultPlcSubscriptionRequest))

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		for _, fieldName := range subscriptionRequest.GetFieldNames() {
			responseCodes[fieldName] = apiModel.PlcResponseCode_OK
		}

		result <- &internalModel.DefaultPlcSubscriptionRequestResult{
			Request:  subscriptionRequest,
			Response: internalModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes),
			Err:      nil,
		}
	}()
	return result
}

func (m *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcUnsubscriptionRequestResult)

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more an internal handling of which values to pass where.

	return result
}

/*
 * Callback for incoming value change events from the KNX bus
 */
func (m *Subscriber) handleValueChange(destinationAddress []byte, payload []byte, changed bool) {
	// Decode the group-address according to the settings in the driver
	// Group addresses can be 1, 2 or 3 levels (3 being the default)
	garb := utils.NewReadBufferByteBased(destinationAddress)
	groupAddress, err := driverModel.KnxGroupAddressParse(garb, m.connection.getGroupAddressNumLevels())
	if err != nil {
		return
	}

	// Go through all subscription-requests and process each separately
	for _, subscriptionRequest := range m.subscriptionRequests {
		fields := map[string]apiModel.PlcField{}
		types := map[string]internalModel.SubscriptionType{}
		intervals := map[string]time.Duration{}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		addresses := map[string][]byte{}
		plcValues := map[string]values.PlcValue{}

		// Check if this datagram matches any address in this subscription request
		// As depending on the address used for fields, the decoding is different, we need to decode on-demand here.
		for _, fieldName := range subscriptionRequest.GetFieldNames() {
			field, err := CastToFieldFromPlcField(subscriptionRequest.GetField(fieldName))
			if err != nil {
				continue
			}
			switch field.(type) {
			case GroupAddressField:
				subscriptionType := subscriptionRequest.GetType(fieldName)
				groupAddressField := field.(GroupAddressField)
				// If it matches, take the datatype of each matching field and try to decode the payload
				if groupAddressField.matches(groupAddress) {
					// If this is a CHANGE_OF_STATE field, filter out the events where the value actually hasn't changed.
					if subscriptionType == internalModel.SubscriptionChangeOfState && changed {
						rb := utils.NewReadBufferByteBased(payload)
						if groupAddressField.GetFieldType() == nil {
							responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_DATATYPE
							plcValues[fieldName] = nil
							continue
						}
						// If the size of the field is greater than 6, we have to skip the first byte
						if groupAddressField.GetFieldType().GetLengthInBits() > 6 {
							_, _ = rb.ReadUint8("groupAddress", 8)
						}
						elementType := *groupAddressField.GetFieldType()
						numElements := groupAddressField.GetQuantity()

						fields[fieldName] = groupAddressField
						types[fieldName] = subscriptionRequest.GetType(fieldName)
						intervals[fieldName] = subscriptionRequest.GetInterval(fieldName)
						addresses[fieldName] = destinationAddress

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
								plcValue := values2.NewRawPlcValue(rb, NewValueDecoder(rb))
								plcValueList = append(plcValueList, plcValue)
							} else {
								plcValue, err2 := driverModel.KnxDatapointParse(rb, elementType)
								if err2 == nil {
									plcValueList = append(plcValueList, plcValue)
								} else {
									// TODO: Do a little more here ...
									responseCode = apiModel.PlcResponseCode_INTERNAL_ERROR
									break
								}
							}
						}
						responseCodes[fieldName] = responseCode
						if responseCode == apiModel.PlcResponseCode_OK {
							if len(plcValueList) == 1 {
								plcValues[fieldName] = plcValueList[0]
							} else {
								plcValues[fieldName] = values2.NewPlcList(plcValueList)
							}
						}
					}
				}
			default:
				responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[fieldName] = nil
			}
		}

		// Assemble a PlcSubscription event
		if len(plcValues) > 0 {
			event := NewSubscriptionEvent(
				fields, types, intervals, responseCodes, addresses, plcValues)
			eventHandler := subscriptionRequest.GetEventHandler()
			eventHandler(event)
		}
	}
}
