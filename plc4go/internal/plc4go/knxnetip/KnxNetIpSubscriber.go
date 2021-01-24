//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package knxnetip

import (
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"time"
)

type KnxNetIpSubscriber struct {
	connection           *KnxNetIpConnection
	subscriptionRequests []internalModel.DefaultPlcSubscriptionRequest
	spi.PlcWriter
}

func NewKnxNetIpSubscriber(connection *KnxNetIpConnection) *KnxNetIpSubscriber {
	return &KnxNetIpSubscriber{
		connection:           connection,
		subscriptionRequests: []internalModel.DefaultPlcSubscriptionRequest{},
	}
}

func (m *KnxNetIpSubscriber) Subscribe(subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
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

		result <- apiModel.PlcSubscriptionRequestResult{
			Request:  subscriptionRequest,
			Response: internalModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes),
			Err:      nil,
		}
	}()
	return result
}

func (m *KnxNetIpSubscriber) Unsubscribe(unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult)

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more an internal handling of which values to pass where.

	return result
}

/*
 * Callback for incoming value change events from the KNX bus
 */
func (m *KnxNetIpSubscriber) handleValueChange(lDataFrame *driverModel.LDataFrame, changed bool) {
	var destinationAddress []int8
    var apdu *driverModel.Apdu
	switch lDataFrame.Child.(type) {
	case *driverModel.LDataFrameData:
		dataFrame := driverModel.CastLDataFrameData(lDataFrame)
		destinationAddress = dataFrame.DestinationAddress
        apdu = dataFrame.Apdu
	case *driverModel.LDataFrameDataExt:
		dataFrame := driverModel.CastLDataFrameDataExt(lDataFrame)
		destinationAddress = dataFrame.DestinationAddress
        apdu = dataFrame.Apdu
	}
    container := driverModel.CastApduDataContainer(apdu)
    if container == nil {
        return
    }
    groupValueWrite := driverModel.CastApduDataGroupValueWrite(container.DataApdu)
    if groupValueWrite == nil {
        return
    }

	if destinationAddress != nil {
		// Decode the group-address according to the settings in the driver
		// Group addresses can be 1, 2 or 3 levels (3 being the default)
		garb := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(destinationAddress))
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
			addresses := map[string][]int8{}
			plcValues := map[string]values.PlcValue{}

			// Check if this datagram matches any address in this subscription request
			// As depending on the address used for fields, the decoding is different, we need to decode on-demand here.
			for _, fieldName := range subscriptionRequest.GetFieldNames() {
				field, err := CastToKnxNetIpFieldFromPlcField(subscriptionRequest.GetField(fieldName))
				if err != nil {
					continue
				}
				subscriptionType := subscriptionRequest.GetType(fieldName)
				// If it matches, take the datatype of each matching field and try to decode the payload
				if field.matches(groupAddress) {
                    // If this is a CHANGE_OF_STATE field, filter out the events where the value actually hasn't changed.
					if subscriptionType == internalModel.SUBSCRIPTION_CHANGE_OF_STATE && changed {
                        var payload []int8
                        payload = append(payload, groupValueWrite.DataFirstByte)
                        payload = append(payload, groupValueWrite.Data...)
                        rb := utils.NewReadBuffer(utils.Int8ArrayToByteArray(payload))
                        plcValue, err := driverModel.KnxDatapointParse(rb, field.GetTypeName())
                        fields[fieldName] = field
                        types[fieldName] = subscriptionRequest.GetType(fieldName)
                        intervals[fieldName] = subscriptionRequest.GetInterval(fieldName)
                        addresses[fieldName] = destinationAddress
                        if err == nil {
                            responseCodes[fieldName] = apiModel.PlcResponseCode_OK
                            plcValues[fieldName] = plcValue
                        } else {
                            // TODO: Do a little more here ...
                            responseCodes[fieldName] = apiModel.PlcResponseCode_INTERNAL_ERROR
                            plcValues[fieldName] = nil
                        }
					}
				}
			}

			// Assemble a PlcSubscription event
			if len(plcValues) > 0 {
				event := NewKnxNetIpSubscriptionEvent(
					fields, types, intervals, responseCodes, addresses, plcValues)
				eventHandler := subscriptionRequest.GetEventHandler()
				eventHandler(event)
			}
		}
	}
}
