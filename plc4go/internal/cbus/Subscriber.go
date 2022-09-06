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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type Subscriber struct {
	connection *Connection
	consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
}

func NewSubscriber(connection *Connection) *Subscriber {
	return &Subscriber{
		connection: connection,
		consumers:  make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),
	}
}

func (m *Subscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult)
	go func() {
		internalPlcSubscriptionRequest := subscriptionRequest.(spiModel.DefaultPlcSubscriptionRequest)

		// Add this subscriber to the connection.
		m.connection.addSubscriber(m)

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
		for _, fieldName := range internalPlcSubscriptionRequest.GetFieldNames() {
			responseCodes[fieldName] = apiModel.PlcResponseCode_OK
			subscriptionValues[fieldName] = NewSubscriptionHandle(m, fieldName, internalPlcSubscriptionRequest.GetField(fieldName), internalPlcSubscriptionRequest.GetType(fieldName), internalPlcSubscriptionRequest.GetInterval(fieldName))
		}

		result <- &spiModel.DefaultPlcSubscriptionRequestResult{
			Request:  subscriptionRequest,
			Response: spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes, subscriptionValues),
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
	for registration, consumer := range m.consumers {
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			subscriptionHandle := subscriptionHandle.(*SubscriptionHandle)
			field, ok := subscriptionHandle.field.(*mmiMonitorField)
			if !ok {
				log.Debug().Msgf("Unusable field for mmi subscription %s", subscriptionHandle.field)
				continue
			}

			fields := map[string]apiModel.PlcField{}
			types := map[string]spiModel.SubscriptionType{}
			intervals := map[string]time.Duration{}
			responseCodes := map[string]apiModel.PlcResponseCode{}
			address := map[string]string{}
			sources := map[string]string{}
			plcValues := map[string]apiValues.PlcValue{}
			fieldName := subscriptionHandle.fieldName

			if unitAddress := field.GetUnitAddress(); unitAddress != nil {
				unitSuffix := fmt.Sprintf("u%d", unitAddress.GetAddress())
				if !strings.HasSuffix(unitAddressString, unitSuffix) {
					log.Debug().Msgf("Current address string %s has not the suffix %s", unitAddressString, unitSuffix)
					continue
				}
			}
			sources[fieldName] = unitAddressString

			subscriptionType := subscriptionHandle.fieldType
			// TODO: handle subscriptionType
			_ = subscriptionType

			fields[fieldName] = field
			types[fieldName] = subscriptionHandle.fieldType
			intervals[fieldName] = subscriptionHandle.interval

			var applicationString string

			isLevel := true
			blockStart := byte(0x0)
			//	var application readWriteModel.ApplicationIdContainer
			switch calData := calData.(type) {
			case readWriteModel.CALDataStatusExactly:
				application := calData.GetApplication()
				applicationString = application.ApplicationId().String()
				blockStart = calData.GetBlockStart()

				statusBytes := calData.GetStatusBytes()
				responseCodes[fieldName] = apiModel.PlcResponseCode_OK
				plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
				for i, statusByte := range statusBytes {
					plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
					plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
					plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
					plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
				}
				plcValues[fieldName] = spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
					"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
					"blockStart":  spiValues.NewPlcBYTE(blockStart),
					"values":      spiValues.NewPlcList(plcListValues),
				})
			case readWriteModel.CALDataStatusExtendedExactly:
				application := calData.GetApplication()
				applicationString = application.ApplicationId().String()
				isLevel = calData.GetCoding() == readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE || calData.GetCoding() == readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE
				blockStart = calData.GetBlockStart()
				coding := calData.GetCoding()
				switch coding {
				case readWriteModel.StatusCoding_BINARY_BY_THIS_SERIAL_INTERFACE:
					fallthrough
				case readWriteModel.StatusCoding_BINARY_BY_ELSEWHERE:
					statusBytes := calData.GetStatusBytes()
					responseCodes[fieldName] = apiModel.PlcResponseCode_OK
					plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
					for i, statusByte := range statusBytes {
						plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
						plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
						plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
						plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
					}
					plcValues[fieldName] = spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
						"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
						"blockStart":  spiValues.NewPlcBYTE(blockStart),
						"values":      spiValues.NewPlcList(plcListValues),
					})
				case readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE:
					fallthrough
				case readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE:
					levelInformation := calData.GetLevelInformation()
					responseCodes[fieldName] = apiModel.PlcResponseCode_OK
					plcListValues := make([]apiValues.PlcValue, len(levelInformation))
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
				event := NewSubscriptionEvent(fields, types, intervals, responseCodes, address, sources, plcValues)
				consumer(event)
			}
		}
	}
	return true
}

func (m *Subscriber) handleMonitoredSal(sal readWriteModel.MonitoredSAL) bool {
	for registration, consumer := range m.consumers {
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			subscriptionHandle := subscriptionHandle.(*SubscriptionHandle)
			field, ok := subscriptionHandle.field.(*salMonitorField)
			if !ok {
				log.Debug().Msgf("Unusable field for mmi subscription %s", subscriptionHandle.field)
				continue
			}
			fields := map[string]apiModel.PlcField{}
			types := map[string]spiModel.SubscriptionType{}
			intervals := map[string]time.Duration{}
			responseCodes := map[string]apiModel.PlcResponseCode{}
			address := map[string]string{}
			sources := map[string]string{}
			plcValues := map[string]apiValues.PlcValue{}
			fieldName := subscriptionHandle.fieldName

			subscriptionType := subscriptionHandle.fieldType
			// TODO: handle subscriptionType
			_ = subscriptionType

			fields[fieldName] = field
			types[fieldName] = subscriptionType
			intervals[fieldName] = subscriptionHandle.interval

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
			sources[fieldName] = unitAddressString

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
				event := NewSubscriptionEvent(fields, types, intervals, responseCodes, address, sources, plcValues)
				consumer(event)
			}
		}
	}
	return true
}

func (m *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(m, consumer, handles...)
	m.consumers[consumerRegistration] = consumer
	return consumerRegistration
}

func (m *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	delete(m.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
