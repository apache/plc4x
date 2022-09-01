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
	"context"
	"fmt"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/rs/zerolog/log"
	"strings"
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

func (m *Subscriber) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	// TODO: handle context
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

func (m *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcUnsubscriptionRequestResult)

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more an internal handling of which values to pass where.

	return result
}

func (m *Subscriber) handleMonitoredMMI(calReply readWriteModel.CALReply) bool {
	var unitAddressString string
	switch calReply := calReply.(type) {
	case readWriteModel.CALReplyLongExactly:
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
	for _, subscriptionRequest := range m.subscriptionRequests {
		fields := map[string]apiModel.PlcField{}
		types := map[string]spiModel.SubscriptionType{}
		intervals := map[string]time.Duration{}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		address := map[string]string{}
		plcValues := map[string]values.PlcValue{}

		for _, fieldName := range subscriptionRequest.GetFieldNames() {
			field, ok := subscriptionRequest.GetField(fieldName).(*mmiMonitorField)
			if !ok {
				log.Debug().Msgf("Unusable field for mmi subscription %s", field)
				continue
			}
			if unitAddress := field.GetUnitAddress(); unitAddress != nil {
				unitSuffix := fmt.Sprintf("u%d", unitAddress.GetAddress())
				if !strings.HasSuffix(unitAddressString, unitSuffix) {
					log.Debug().Msgf("Current address string %s has not the suffix %s", unitAddressString, unitSuffix)
					continue
				}
			}

			subscriptionType := subscriptionRequest.GetType(fieldName)
			// TODO: handle subscriptionType
			_ = subscriptionType

			fields[fieldName] = field
			types[fieldName] = subscriptionRequest.GetType(fieldName)
			intervals[fieldName] = subscriptionRequest.GetInterval(fieldName)

			var applicationString string

			isLevel := true
			blockStart := byte(0x0)
			switch calData := calData.(type) {
			case readWriteModel.CALDataStatusExactly:
				applicationString = calData.GetApplication().ApplicationId().String()
				blockStart = calData.GetBlockStart()

				statusBytes := calData.GetStatusBytes()
				responseCodes[fieldName] = apiModel.PlcResponseCode_OK
				plcListValues := make([]values.PlcValue, len(statusBytes)*4)
				for i, statusByte := range statusBytes {
					plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
					plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
					plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
					plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
				}
				plcValues[fieldName] = spiValues.NewPlcList(plcListValues)
			case readWriteModel.CALDataStatusExtendedExactly:
				applicationString = calData.GetApplication().ApplicationId().String()
				isLevel = calData.GetCoding() == readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE || calData.GetCoding() == readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE
				blockStart = calData.GetBlockStart()
				coding := calData.GetCoding()
				switch coding {
				case readWriteModel.StatusCoding_BINARY_BY_THIS_SERIAL_INTERFACE:
					fallthrough
				case readWriteModel.StatusCoding_BINARY_BY_ELSEWHERE:
					statusBytes := calData.GetStatusBytes()
					responseCodes[fieldName] = apiModel.PlcResponseCode_OK
					plcListValues := make([]values.PlcValue, len(statusBytes)*4)
					for i, statusByte := range statusBytes {
						plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
						plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
						plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
						plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
					}
					plcValues[fieldName] = spiValues.NewPlcList(plcListValues)
				case readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE:
					fallthrough
				case readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE:
					levelInformation := calData.GetLevelInformation()
					responseCodes[fieldName] = apiModel.PlcResponseCode_OK
					plcListValues := make([]values.PlcValue, len(levelInformation))
					for i, levelInformation := range levelInformation {
						switch levelInformation := levelInformation.(type) {
						case readWriteModel.LevelInformationAbsentExactly:
							plcListValues[i] = spiValues.NewPlcSTRING("is absent")
						case readWriteModel.LevelInformationCorruptedExactly:
							plcListValues[i] = spiValues.NewPlcSTRING("corrupted")
						case readWriteModel.LevelInformationNormalExactly:
							plcListValues[i] = spiValues.NewPlcUSINT(levelInformation.GetActualLevel())
						default:
							panic("Impossible case")
						}
					}
					plcValues[fieldName] = spiValues.NewPlcList(plcListValues)
				}
			default:
				return false
			}
			if application := field.GetApplication(); application != nil {
				if actualApplicationIdString := application.ApplicationId().String(); applicationString != actualApplicationIdString {
					log.Debug().Msgf("Current application id %s  doesn't match actual id %s", unitAddressString, actualApplicationIdString)
					continue
				}
			}
			statusType := "binary"
			if isLevel {
				statusType = fmt.Sprintf("level=0x%X", blockStart)
			}
			address[fieldName] = fmt.Sprintf("status/%s/%s", statusType, applicationString)

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

func (m *Subscriber) handleMonitoredSal(sal readWriteModel.MonitoredSAL) bool {
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
				log.Debug().Msgf("Unusable field for sal subscription %s", field)
				continue
			}

			subscriptionType := subscriptionRequest.GetType(fieldName)
			// TODO: handle subscriptionType
			_ = subscriptionType

			fields[fieldName] = field
			types[fieldName] = subscriptionRequest.GetType(fieldName)
			intervals[fieldName] = subscriptionRequest.GetInterval(fieldName)

			var salData readWriteModel.SALData
			var unitAddressString, applicationString string
			switch sal := sal.(type) {
			case readWriteModel.MonitoredSALLongFormSmartModeExactly:
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
			case readWriteModel.MonitoredSALShortFormBasicModeExactly:
				unitAddressString = "u0" // On short form it should be always unit 0 TODO: double check that
				applicationString = sal.GetApplication().ApplicationId().String()
				salData = sal.GetSalData()
			}
			if unitAddress := field.GetUnitAddress(); unitAddress != nil {
				unitSuffix := fmt.Sprintf("u%d", unitAddress.GetAddress())
				if !strings.HasSuffix(unitAddressString, unitSuffix) {
					log.Debug().Msgf("Current address string %s has not the suffix %s", unitAddressString, unitSuffix)
					continue
				}
			}

			if application := field.GetApplication(); application != nil {
				if actualApplicationIdString := application.ApplicationId().String(); applicationString != actualApplicationIdString {
					log.Debug().Msgf("Current application id %s  doesn't match actual id %s", unitAddressString, actualApplicationIdString)
					continue
				}
			}

			// TODO: we need to map commands e.g. if we get a MeteringDataElectricityConsumption we can map that to MeteringDataMeasureElectricity
			address[fieldName] = fmt.Sprintf("sal/%s/%s", applicationString, "TODO")

			rbvb := spiValues.NewWriteBufferPlcValueBased()
			err := salData.Serialize(rbvb)
			if err != nil {
				log.Error().Err(err).Msg("Error serializing to plc value... just returning it as string")
				plcValues[fieldName] = spiValues.NewPlcSTRING(fmt.Sprintf("%s", salData))
			} else {
				plcValues[fieldName] = rbvb.GetPlcValue()
			}

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
