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
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

//go:generate go tool plc4xGenerator -type=Subscriber
type Subscriber struct {
	consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer `hasLocker:"consumersMutex"`
	addSubscriber func(subscriber *Subscriber)

	consumersMutex sync.RWMutex

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

func NewSubscriber(addSubscriber func(subscriber *Subscriber), _options ...options.WithOption) *Subscriber {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Subscriber{
		addSubscriber: addSubscriber,
		consumers:     make(map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer),

		log:      customLogger,
		_options: _options,
	}
}

func (s *Subscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	result := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		internalPlcSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

		// Add this subscriber to the connection.
		s.addSubscriber(s)

		// Just populate all requests with an OK
		responseCodes := map[string]apiModel.PlcResponseCode{}
		subscriptionValues := make(map[string]apiModel.PlcSubscriptionHandle)
		for _, tagName := range internalPlcSubscriptionRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			handle := NewSubscriptionHandle(
				s,
				tagName,
				internalPlcSubscriptionRequest.GetTag(tagName),
				internalPlcSubscriptionRequest.GetType(tagName),
				internalPlcSubscriptionRequest.GetInterval(tagName),
			)
			preRegisteredConsumers := internalPlcSubscriptionRequest.GetPreRegisteredConsumers(tagName)
			for _, consumer := range preRegisteredConsumers {
				_ = handle.Register(consumer)
			}
			subscriptionValues[tagName] = handle
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
	// subscriptions are more a internal handling of which values to pass where.
	_ = ctx
	_ = unsubscriptionRequest

	return result
}

func (s *Subscriber) handleMonitoredMMI(calReply readWriteModel.CALReply) bool {
	s.log.Debug().Stringer("calReply", calReply).Msg("handling")
	var unitAddressString string
	switch calReply := calReply.(type) {
	case readWriteModel.CALReplyLong:
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
	s.log.Debug().Str("unitAddressString", unitAddressString).Msg("Unit address string")
	calData := calReply.GetCalData()
	handled := false
	s.consumersMutex.RLock()
	defer s.consumersMutex.RUnlock()
	for registration, consumer := range s.consumers {
		s.log.Debug().
			Stringer("registration", registration).
			Interface("consumer", consumer).
			Msg("Checking with registration and consumer")
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			s.log.Debug().Stringer("subscriptionHandle", subscriptionHandle).Msg("offering to")
			handleHandled := s.offerMMI(unitAddressString, calData, subscriptionHandle.(*SubscriptionHandle), consumer)
			s.log.Debug().Bool("handleHandled", handleHandled).Msg("handle handled")
			handled = handled || handleHandled
		}
	}
	s.log.Debug().Bool("handled", handled).Msg("final handled")
	return handled
}

func (s *Subscriber) offerMMI(unitAddressString string, calData readWriteModel.CALData, subscriptionHandle *SubscriptionHandle, consumer apiModel.PlcSubscriptionEventConsumer) bool {
	tag, ok := subscriptionHandle.tag.(*mmiMonitorTag)
	if !ok {
		s.log.Debug().
			Interface("tag", subscriptionHandle.tag).
			Msg("Unusable tag for mmi subscription")
		return false
	}

	tags := map[string]apiModel.PlcTag{}
	types := map[string]apiModel.PlcSubscriptionType{}
	intervals := map[string]time.Duration{}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	address := map[string]string{}
	sources := map[string]string{}
	plcValues := map[string]apiValues.PlcValue{}
	tagName := subscriptionHandle.tagName

	if unitAddress := tag.GetUnitAddress(); unitAddress != nil {
		unitSuffix := fmt.Sprintf("u%d", unitAddress.GetAddress())
		if !strings.HasSuffix(unitAddressString, unitSuffix) {
			s.log.Debug().
				Str("unitAddressString", unitAddressString).
				Str("unitSuffix", unitSuffix).
				Msg("Current address string unitAddressString has not the suffix unitSuffix")
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
	case readWriteModel.CALDataStatus:
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
	case readWriteModel.CALDataStatusExtended:
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
				case readWriteModel.LevelInformationAbsent:
					plcListValues[i] = spiValues.NewPlcSTRING("is absent")
				case readWriteModel.LevelInformationCorrupted:
					plcListValues[i] = spiValues.NewPlcSTRING("corrupted")
				case readWriteModel.LevelInformationNormal:
					plcListValues[i] = spiValues.NewPlcUSINT(levelInformation.GetActualLevel())
				}
			}
			plcValues[tagName] = spiValues.NewPlcList(plcListValues)
		}
	default:
		s.log.Error().Type("calData", calData).Msg("Unmapped type")
		return false
	}
	if application := tag.GetApplication(); application != nil {
		if actualApplicationIdString := application.ApplicationId().String(); applicationString != actualApplicationIdString {
			s.log.Debug().
				Str("unitAddressString", unitAddressString).
				Str("actualApplicationIdString", actualApplicationIdString).
				Msg("Current application id unitAddressString doesn't match actual id actualApplicationIdString")
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

func (s *Subscriber) handleMonitoredSAL(sal readWriteModel.MonitoredSAL) bool {
	handled := false
	s.consumersMutex.RLock()
	defer s.consumersMutex.RUnlock()
	for registration, consumer := range s.consumers {
		for _, subscriptionHandle := range registration.GetSubscriptionHandles() {
			handled = handled || s.offerSAL(sal, subscriptionHandle.(*SubscriptionHandle), consumer)
		}
	}
	return handled
}

func (s *Subscriber) offerSAL(sal readWriteModel.MonitoredSAL, subscriptionHandle *SubscriptionHandle, consumer apiModel.PlcSubscriptionEventConsumer) bool {
	tag, ok := subscriptionHandle.tag.(*salMonitorTag)
	if !ok {
		s.log.Debug().Interface("tag", subscriptionHandle.tag).Msg("Unusable tag for mmi subscription")
		return false
	}
	tags := map[string]apiModel.PlcTag{}
	types := map[string]apiModel.PlcSubscriptionType{}
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
	case readWriteModel.MonitoredSALShortFormBasicMode:
		unitAddressString = "u0" // On short form it should be always unit 0 TODO: double check that
		applicationString = sal.GetApplication().ApplicationId().String()
		salData = sal.GetSalData()
	case readWriteModel.MonitoredSALLongFormSmartMode:
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
			s.log.Debug().
				Str("unitAddressString", unitAddressString).
				Str("unitSuffix", unitSuffix).
				Msg("Current address string unitAddressString has not the suffix unitSuffix")
			return false
		}
	}
	sources[tagName] = unitAddressString

	if application := tag.GetApplication(); application != nil {
		if actualApplicationIdString := application.ApplicationId().String(); applicationString != actualApplicationIdString {
			s.log.Debug().
				Str("unitAddressString", unitAddressString).
				Str("actualApplicationIdString", actualApplicationIdString).
				Msg("Current application id unitAddressString  doesn't match actual id actualApplicationIdString")
			return false
		}
	}

	var commandTypeGetter interface {
		PLC4XEnumName() string
	}
	switch salData := salData.(type) {
	case readWriteModel.SALDataAccessControl:
		commandTypeGetter = salData.GetAccessControlData().GetCommandType()
	case readWriteModel.SALDataAirConditioning:
		commandTypeGetter = salData.GetAirConditioningData().GetCommandType()
	case readWriteModel.SALDataAudioAndVideo:
		commandTypeGetter = salData.GetAudioVideoData().GetCommandType()
	case readWriteModel.SALDataClockAndTimekeeping:
		commandTypeGetter = salData.GetClockAndTimekeepingData().GetCommandType()
	case readWriteModel.SALDataEnableControl:
		commandTypeGetter = salData.GetEnableControlData().GetCommandType()
	case readWriteModel.SALDataErrorReporting:
		commandTypeGetter = salData.GetErrorReportingData().GetCommandType()
	case readWriteModel.SALDataFreeUsage:
		s.log.Info().Msg("Unknown command type")
	case readWriteModel.SALDataHeating:
		commandTypeGetter = salData.GetHeatingData().GetCommandType()
	case readWriteModel.SALDataHvacActuator:
		commandTypeGetter = salData.GetHvacActuatorData().GetCommandType()
	case readWriteModel.SALDataIrrigationControl:
		commandTypeGetter = salData.GetIrrigationControlData().GetCommandType()
	case readWriteModel.SALDataLighting:
		commandTypeGetter = salData.GetLightingData().GetCommandType()
	case readWriteModel.SALDataMeasurement:
		commandTypeGetter = salData.GetMeasurementData().GetCommandType()
	case readWriteModel.SALDataMediaTransport:
		commandTypeGetter = salData.GetMediaTransportControlData().GetCommandType()
	case readWriteModel.SALDataMetering:
		commandTypeGetter = salData.GetMeteringData().GetCommandType()
	case readWriteModel.SALDataPoolsSpasPondsFountainsControl:
		commandTypeGetter = salData.GetPoolsSpaPondsFountainsData().GetCommandType()
	case readWriteModel.SALDataReserved:
		s.log.Info().Msg("Unknown command type")
	case readWriteModel.SALDataRoomControlSystem:
		s.log.Info().Msg("Unknown command type not implemented yet") // TODO: implement once there
	case readWriteModel.SALDataSecurity:
		commandTypeGetter = salData.GetSecurityData().GetCommandType()
	case readWriteModel.SALDataTelephonyStatusAndControl:
		commandTypeGetter = salData.GetTelephonyData().GetCommandType()
	case readWriteModel.SALDataTemperatureBroadcast:
		commandTypeGetter = salData.GetTemperatureBroadcastData().GetCommandType()
	case readWriteModel.SALDataTesting:
		s.log.Info().Msg("Unknown command type not implemented yet") // TODO: implement once there
	case readWriteModel.SALDataTriggerControl:
		commandTypeGetter = salData.GetTriggerControlData().GetCommandType()
	case readWriteModel.SALDataVentilation:
		commandTypeGetter = salData.GetVentilationData().GetCommandType()
	default:
		s.log.Error().Type("salData", salData).Msg("Unmapped type")
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
		s.log.Error().Err(err).Msg("Error serializing to plc value... just returning it as string")
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

func (s *Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	consumerRegistration := spiModel.NewDefaultPlcConsumerRegistration(s, consumer, handles...)
	s.consumers[consumerRegistration.(*spiModel.DefaultPlcConsumerRegistration)] = consumer
	return consumerRegistration
}

func (s *Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	s.log.Trace().Msg("unregister")
	s.consumersMutex.Lock()
	defer s.consumersMutex.Unlock()
	delete(s.consumers, registration.(*spiModel.DefaultPlcConsumerRegistration))
	s.log.Trace().Msg("registration removed")
}
