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

package cbus

import (
	"fmt"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/rs/zerolog/log"
	"time"
)

type Subscriber struct {
	connection           *Connection
	subscriptionRequests []spiModel.DefaultPlcSubscriptionRequest
}

func NewSubscriber(connection *Connection) *Subscriber {
	return &Subscriber{
		connection:           connection,
		subscriptionRequests: []spiModel.DefaultPlcSubscriptionRequest{},
	}
}

func (m *Subscriber) Subscribe(subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult)
	go func() {
		// Add this subscriber to the connection.
		m.connection.addSubscriber(m)

		// Save the subscription request
		m.subscriptionRequests = append(m.subscriptionRequests, subscriptionRequest.(spiModel.DefaultPlcSubscriptionRequest))

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		for _, fieldName := range subscriptionRequest.GetFieldNames() {
			responseCodes[fieldName] = apiModel.PlcResponseCode_OK
		}

		result <- &spiModel.DefaultPlcSubscriptionRequestResult{
			Request:  subscriptionRequest,
			Response: spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes),
			Err:      nil,
		}
	}()
	return result
}

func (m *Subscriber) Unsubscribe(unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	result := make(chan apiModel.PlcUnsubscriptionRequestResult)

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more an internal handling of which values to pass where.

	return result
}

func (m *Subscriber) handleMonitoredMMI(calReply model.CALReply) bool {
	var unitAddressString string
	switch calReply := calReply.(type) {
	case model.CALReplyLongExactly:
		if calReply.GetIsUnitAddress() {
			unitAddressString = fmt.Sprintf("u%d", calReply.GetUnitAddress().GetAddress())
		} else {
			unitAddressString = fmt.Sprintf("b%d", calReply.GetBridgeAddress().GetAddress())
			replyNetwork := calReply.GetReplyNetwork()
			for _, bridgeAddress := range replyNetwork.GetNetworkRoute().GetAdditionalBridgeAddresses() {
				unitAddressString += fmt.Sprintf("-b%d", bridgeAddress.GetAddress())
			}
			unitAddressString += fmt.Sprintf("-u%d", replyNetwork.GetUnitAddress().GetAddress())
		}
	default:
		unitAddressString = "u0" // On short form it should be always unit 0 TODO: double check that
	}
	calData := calReply.GetCalData()
	// TODO: filter
	for _, subscriptionRequest := range m.subscriptionRequests {
		fields := map[string]apiModel.PlcField{}
		types := map[string]spiModel.SubscriptionType{}
		intervals := map[string]time.Duration{}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		address := map[string]string{}
		plcValues := map[string]values.PlcValue{}

		for _, fieldName := range subscriptionRequest.GetFieldNames() {
			field, ok := subscriptionRequest.GetField(fieldName).(MMIMonitorField)
			if !ok {
				log.Warn().Msgf("Unusable field for subscription %s", field)
				responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[fieldName] = nil
				continue
			}
			if unitAddress := field.GetUnitAddress(); unitAddress != nil {
				// TODO: filter in unit address
			}
			application := field.GetApplication()
			// TODO: filter in unit address
			_ = application

			subscriptionType := subscriptionRequest.GetType(fieldName)
			// TODO: handle subscriptionType
			_ = subscriptionType

			fields[fieldName] = field
			types[fieldName] = subscriptionRequest.GetType(fieldName)
			intervals[fieldName] = subscriptionRequest.GetInterval(fieldName)

			var applicationString string

			switch calData := calData.(type) {
			case model.CALDataStatusExactly:
				applicationString = calData.GetApplication().ApplicationId().String()
			case model.CALDataStatusExtendedExactly:
				applicationString = calData.GetApplication().ApplicationId().String()
			default:
				return false
			}
			// TODO: we might need to encode more data into the address from sal data
			address[fieldName] = fmt.Sprintf("/%s/%s", unitAddressString, applicationString)

			// TODO: map values properly
			plcValues[fieldName] = spiValues.NewPlcSTRING(fmt.Sprintf("%s", calData))
			responseCodes[fieldName] = apiModel.PlcResponseCode_OK

			// Assemble a PlcSubscription event
			if len(plcValues) > 0 {
				event := NewSubscriptionEvent(fields, types, intervals, responseCodes, address, plcValues)
				eventHandler := subscriptionRequest.GetEventHandler()
				eventHandler(event)
			}
		}
	}
	return true
}

func (m *Subscriber) handleMonitoredSal(sal model.MonitoredSAL) bool {
	// TODO: filter
	for _, subscriptionRequest := range m.subscriptionRequests {
		fields := map[string]apiModel.PlcField{}
		types := map[string]spiModel.SubscriptionType{}
		intervals := map[string]time.Duration{}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		address := map[string]string{}
		plcValues := map[string]values.PlcValue{}

		for _, fieldName := range subscriptionRequest.GetFieldNames() {
			field, ok := subscriptionRequest.GetField(fieldName).(SALMonitorField)
			if !ok {
				log.Warn().Msgf("Unusable field for subscription %s", field)
				responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[fieldName] = nil
				continue
			}
			if unitAddress := field.GetUnitAddress(); unitAddress != nil {
				// TODO: filter in unit address
			}
			application := field.GetApplication()
			// TODO: filter in unit address
			_ = application

			subscriptionType := subscriptionRequest.GetType(fieldName)
			// TODO: handle subscriptionType
			_ = subscriptionType

			fields[fieldName] = field
			types[fieldName] = subscriptionRequest.GetType(fieldName)
			intervals[fieldName] = subscriptionRequest.GetInterval(fieldName)

			var salData model.SALData
			var unitAddressString, applicationString string
			switch sal := sal.(type) {
			case model.MonitoredSALLongFormSmartModeExactly:
				if sal.GetIsUnitAddress() {
					unitAddressString = fmt.Sprintf("u%d", sal.GetUnitAddress().GetAddress())
				} else {
					unitAddressString = fmt.Sprintf("b%d", sal.GetBridgeAddress().GetAddress())
					replyNetwork := sal.GetReplyNetwork()
					for _, bridgeAddress := range replyNetwork.GetNetworkRoute().GetAdditionalBridgeAddresses() {
						unitAddressString += fmt.Sprintf("-b%d", bridgeAddress.GetAddress())
					}
					unitAddressString += fmt.Sprintf("-u%d", replyNetwork.GetUnitAddress().GetAddress())
				}
				applicationString = sal.GetApplication().ApplicationId().String()
				salData = sal.GetSalData()
			case model.MonitoredSALShortFormBasicModeExactly:
				unitAddressString = "u0" // On short form it should be always unit 0 TODO: double check that
				applicationString = sal.GetApplication().ApplicationId().String()
				salData = sal.GetSalData()
			}
			// TODO: we might need to encode more data into the address from sal data
			address[fieldName] = fmt.Sprintf("/%s/%s", unitAddressString, applicationString)

			// TODO: map values properly
			plcValues[fieldName] = spiValues.NewPlcSTRING(fmt.Sprintf("%s", salData))
			responseCodes[fieldName] = apiModel.PlcResponseCode_OK

			// Assemble a PlcSubscription event
			if len(plcValues) > 0 {
				event := NewSubscriptionEvent(fields, types, intervals, responseCodes, address, plcValues)
				eventHandler := subscriptionRequest.GetEventHandler()
				eventHandler(event)
			}
		}
	}
	return true
}
