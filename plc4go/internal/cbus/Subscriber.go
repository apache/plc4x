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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"strings"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type Subscriber struct {
	connection *Connection
	consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer

	log zerolog.Logger
}

func NewSubscriber(connection *Connection, _options ...options.WithOption) *Subscriber {
	return &Subscriber{
		connection: connection,
		consumers:  make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),

		log: options.ExtractCustomLogger(_options...),
	}
}

func (m *Subscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v", err))
			}
		}()
		internalPlcSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

		// Add this subscriber to the connection.
		m.connection.addSubscriber(m)

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
		for _, tagName := range internalPlcSubscriptionRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			subscriptionValues[tagName] = NewSubscriptionHandle(m, tagName, internalPlcSubscriptionRequest.GetTag(tagName), internalPlcSubscriptionRequest.GetType(tagName), internalPlcSubscriptionRequest.GetInterval(tagName))
		}

		result <- spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(subscriptionRequest, responseCodes, subscriptionValues),
			nil,
		)
	}()
	return result
}

func (m *Subscriber) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("Not Implemented"))

	// TODO: As soon as we establish a connection, we start getting data...
	// subscriptions are more a internal handling of which values to pass where.
	_ = ctx
	_ = unsubscriptionRequest

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
	handled := false
	for registration, consumer := range m.consumers {
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			handled = handled || m.offerMMI(unitAddressString, calData, subscriptionHandle.(*SubscriptionHandle), consumer)
		}
	}
	return handled
}

func (m *Subscriber) offerMMI(unitAddressString string, calData readWriteModel.CALData, subscriptionHandle *SubscriptionHandle, consumer apiModel.PlcSubscriptionEventConsumer) bool {
	tag, ok := subscriptionHandle.tag.(*mmiMonitorTag)
	if !ok {
		m.log.Debug().Msgf("Unusable tag for mmi subscription %s", subscriptionHandle.tag)
		return false
	}

	tags := map[string]apiModel.PlcTag{}
	types := map[string]spiModel.SubscriptionType{}
	intervals := map[string]time.Duration{}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	address := map[string]string{}
	sources := map[string]string{}
	plcValues := map[string]apiValues.PlcValue{}
	tagName := subscriptionHandle.tagName

	if unitAddress := tag.GetUnitAddress(); unitAddress != nil {
		unitSuffix := fmt.Sprintf("u%d", unitAddress.GetAddress())
		if !strings.HasSuffix(unitAddressString, unitSuffix) {
			m.log.Debug().Msgf("Current address string %s has not the suffix %s", unitAddressString, unitSuffix)
			return false
		}
	}
	sources[tagName] = unitAddressString

	subscriptionType := subscriptionHandle.subscriptionType
	// TODO: handle subscriptionType
	_ = subscriptionType

	tags[tagName] = tag
	types[tagName] = subscriptionHandle.subscriptionType
	intervals[tagName] = subscriptionHandle.interval

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
		responseCodes[tagName] = apiModel.PlcResponseCode_OK
		plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
		for i, statusByte := range statusBytes {
			plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
			plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
			plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
			plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
		}
		plcValues[tagName] = spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
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
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
			for i, statusByte := range statusBytes {
				plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
				plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
				plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
				plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
			}
			plcValues[tagName] = spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
				"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
				"blockStart":  spiValues.NewPlcBYTE(blockStart),
				"values":      spiValues.NewPlcList(plcListValues),
			})
		case readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE:
			fallthrough
		case readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE:
			levelInformation := calData.GetLevelInformation()
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			plcListValues := make([]apiValues.PlcValue, len(levelInformation))
			for i, levelInformation := range levelInformation {
				switch levelInformation := levelInformation.(type) {
				case readWriteModel.LevelInformationAbsentExactly:
					plcListValues[i] = spiValues.NewPlcSTRING("is absent")
				case readWriteModel.LevelInformationCorruptedExactly:
					plcListValues[i] = spiValues.NewPlcSTRING("corrupted")
				case readWriteModel.LevelInformationNormalExactly:
					plcListValues[i] = spiValues.NewPlcUSINT(levelInformation.GetActualLevel())
				}
			}
			plcValues[tagName] = spiValues.NewPlcList(plcListValues)
		}
	default:
		m.log.Error().Msgf("Unmapped type %T", calData)
		return false
	}
	if application := tag.GetApplication(); application != nil {
		if actualApplicationIdString := application.ApplicationId().String(); applicationString != actualApplicationIdString {
			m.log.Debug().Msgf("Current application id %s  doesn't match actual id %s", unitAddressString, actualApplicationIdString)
			return false
		}
	}
	statusType := "binary"
	if isLevel {
		statusType = fmt.Sprintf("level=%#02X", blockStart)
	}
	address[tagName] = fmt.Sprintf("status/%s/%s", statusType, applicationString)

	// Assemble a PlcSubscription event
	event := NewSubscriptionEvent(tags, types, intervals, responseCodes, address, sources, plcValues)
	consumer(&event)
	return true
}

func (m *Subscriber) handleMonitoredSAL(sal readWriteModel.MonitoredSAL) bool {
	handled := false
	for registration, consumer := range m.consumers {
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			handled = handled || m.offerSAL(sal, subscriptionHandle.(*SubscriptionHandle), consumer)
		}
	}
	return handled
}

func (m *Subscriber) offerSAL(sal readWriteModel.MonitoredSAL, subscriptionHandle *SubscriptionHandle, consumer apiModel.PlcSubscriptionEventConsumer) bool {
	tag, ok := subscriptionHandle.tag.(*salMonitorTag)
	if !ok {
		m.log.Debug().Msgf("Unusable tag for mmi subscription %s", subscriptionHandle.tag)
		return false
	}
	tags := map[string]apiModel.PlcTag{}
	types := map[string]spiModel.SubscriptionType{}
	intervals := map[string]time.Duration{}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	address := map[string]string{}
	sources := map[string]string{}
	plcValues := map[string]apiValues.PlcValue{}
	tagName := subscriptionHandle.tagName

	subscriptionType := subscriptionHandle.subscriptionType
	// TODO: handle subscriptionType
	_ = subscriptionType

	tags[tagName] = tag
	types[tagName] = subscriptionType
	intervals[tagName] = subscriptionHandle.interval

	var salData readWriteModel.SALData
	var unitAddressString, applicationString string
	switch sal := sal.(type) {
	case readWriteModel.MonitoredSALShortFormBasicModeExactly:
		unitAddressString = "u0" // On short form it should be always unit 0 TODO: double check that
		applicationString = sal.GetApplication().ApplicationId().String()
		salData = sal.GetSalData()
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
	}
	if unitAddress := tag.GetUnitAddress(); unitAddress != nil {
		unitSuffix := fmt.Sprintf("u%d", unitAddress.GetAddress())
		if !strings.HasSuffix(unitAddressString, unitSuffix) {
			m.log.Debug().Msgf("Current address string %s has not the suffix %s", unitAddressString, unitSuffix)
			return false
		}
	}
	sources[tagName] = unitAddressString

	if application := tag.GetApplication(); application != nil {
		if actualApplicationIdString := application.ApplicationId().String(); applicationString != actualApplicationIdString {
			m.log.Debug().Msgf("Current application id %s  doesn't match actual id %s", unitAddressString, actualApplicationIdString)
			return false
		}
	}

	var commandTypeGetter interface {
		PLC4XEnumName() string
	}
	switch salData := salData.(type) {
	case readWriteModel.SALDataAccessControlExactly:
		commandTypeGetter = salData.GetAccessControlData().GetCommandType()
	case readWriteModel.SALDataAirConditioningExactly:
		commandTypeGetter = salData.GetAirConditioningData().GetCommandType()
	case readWriteModel.SALDataAudioAndVideoExactly:
		commandTypeGetter = salData.GetAudioVideoData().GetCommandType()
	case readWriteModel.SALDataClockAndTimekeepingExactly:
		commandTypeGetter = salData.GetClockAndTimekeepingData().GetCommandType()
	case readWriteModel.SALDataEnableControlExactly:
		commandTypeGetter = salData.GetEnableControlData().GetCommandType()
	case readWriteModel.SALDataErrorReportingExactly:
		commandTypeGetter = salData.GetErrorReportingData().GetCommandType()
	case readWriteModel.SALDataFreeUsageExactly:
		m.log.Info().Msg("Unknown command type")
	case readWriteModel.SALDataHeatingExactly:
		commandTypeGetter = salData.GetHeatingData().GetCommandType()
	case readWriteModel.SALDataHvacActuatorExactly:
		commandTypeGetter = salData.GetHvacActuatorData().GetCommandType()
	case readWriteModel.SALDataIrrigationControlExactly:
		commandTypeGetter = salData.GetIrrigationControlData().GetCommandType()
	case readWriteModel.SALDataLightingExactly:
		commandTypeGetter = salData.GetLightingData().GetCommandType()
	case readWriteModel.SALDataMeasurementExactly:
		commandTypeGetter = salData.GetMeasurementData().GetCommandType()
	case readWriteModel.SALDataMediaTransportExactly:
		commandTypeGetter = salData.GetMediaTransportControlData().GetCommandType()
	case readWriteModel.SALDataMeteringExactly:
		commandTypeGetter = salData.GetMeteringData().GetCommandType()
	case readWriteModel.SALDataPoolsSpasPondsFountainsControlExactly:
		commandTypeGetter = salData.GetPoolsSpaPondsFountainsData().GetCommandType()
	case readWriteModel.SALDataReservedExactly:
		m.log.Info().Msg("Unknown command type")
	case readWriteModel.SALDataRoomControlSystemExactly:
		m.log.Info().Msg("Unknown command type not implemented yet") // TODO: implement once there
	case readWriteModel.SALDataSecurityExactly:
		commandTypeGetter = salData.GetSecurityData().GetCommandType()
	case readWriteModel.SALDataTelephonyStatusAndControlExactly:
		commandTypeGetter = salData.GetTelephonyData().GetCommandType()
	case readWriteModel.SALDataTemperatureBroadcastExactly:
		commandTypeGetter = salData.GetTemperatureBroadcastData().GetCommandType()
	case readWriteModel.SALDataTestingExactly:
		m.log.Info().Msg("Unknown command type not implemented yet") // TODO: implement once there
	case readWriteModel.SALDataTriggerControlExactly:
		commandTypeGetter = salData.GetTriggerControlData().GetCommandType()
	case readWriteModel.SALDataVentilationExactly:
		commandTypeGetter = salData.GetVentilationData().GetCommandType()
	default:
		m.log.Error().Msgf("Unmapped type %T", salData)
	}
	commandType := "Unknown"
	if commandTypeGetter != nil {
		commandType = commandTypeGetter.PLC4XEnumName()
	}

	// TODO: we need to map commands e.g. if we get a MeteringDataElectricityConsumption we can map that to MeteringDataMeasureElectricity
	address[tagName] = fmt.Sprintf("sal/%s/%s", applicationString, commandType)

	rbvb := spiValues.NewWriteBufferPlcValueBased()
	err := salData.SerializeWithWriteBuffer(context.Background(), rbvb)
	if err != nil {
		m.log.Error().Err(err).Msg("Error serializing to plc value... just returning it as string")
		plcValues[tagName] = spiValues.NewPlcSTRING(fmt.Sprintf("%s", salData))
	} else {
		plcValues[tagName] = rbvb.GetPlcValue()
	}

	responseCodes[tagName] = apiModel.PlcResponseCode_OK

	// Assemble a PlcSubscription event
	event := NewSubscriptionEvent(tags, types, intervals, responseCodes, address, sources, plcValues)
	consumer(&event)
	return true
}

func (m *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(m, consumer, handles...)
	m.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (m *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	delete(m.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
}
