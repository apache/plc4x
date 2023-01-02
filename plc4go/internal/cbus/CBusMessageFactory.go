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
	"strings"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/pkg/errors"
)

func TagToCBusMessage(tag model.PlcTag, value values.PlcValue, alphaGenerator *AlphaGenerator, messageCodec *MessageCodec) (cBusMessage readWriteModel.CBusMessage, supportsRead, supportsWrite, supportsSubscribe bool, err error) {
	cbusOptions := messageCodec.cbusOptions
	requestContext := messageCodec.requestContext
	switch tagType := tag.(type) {
	case *statusTag:
		var statusRequest readWriteModel.StatusRequest
		switch tagType.statusRequestType {
		case StatusRequestTypeBinaryState:
			statusRequest = readWriteModel.NewStatusRequestBinaryState(tagType.application, 0x7A)
		case StatusRequestTypeLevel:
			statusRequest = readWriteModel.NewStatusRequestLevel(tagType.application, *tagType.startingGroupAddressLabel, 0x73)
		}
		command := readWriteModel.NewCBusPointToMultiPointCommandStatus(statusRequest, byte(tagType.application), cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToMultiPoint(command, header, cbusOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead, supportsSubscribe = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true, true
		return
	case *calRecallTag:
		calData := readWriteModel.NewCALDataRecall(tagType.parameter, tagType.count, readWriteModel.CALCommandTypeContainer_CALCommandRecall, nil, requestContext)
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToPointCommandDirect(tagType.unitAddress, 0x0000, calData, cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToPoint(command, header, cbusOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true
		return
	case *calIdentifyTag:
		calData := readWriteModel.NewCALDataIdentify(tagType.attribute, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, requestContext)
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToPointCommandDirect(tagType.unitAddress, 0x0000, calData, cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToPoint(command, header, cbusOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true
		return
	case *calGetstatusTag:
		calData := readWriteModel.NewCALDataGetStatus(tagType.parameter, tagType.count, readWriteModel.CALCommandTypeContainer_CALCommandGetStatus, nil, requestContext)
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToPointCommandDirect(tagType.unitAddress, 0x0000, calData, cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToPoint(command, header, cbusOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true
		return
	case *salTag:
		var salCommand = tagType.salCommand
		if salCommand == "" {
			return nil, false, false, false, errors.New("Empty sal command not supported")
		}
		var salData readWriteModel.SALData
		switch tagType.application.ApplicationId() {
		case readWriteModel.ApplicationId_FREE_USAGE:
			panic("Not yet implemented") // TODO: implement
		case readWriteModel.ApplicationId_TEMPERATURE_BROADCAST:
			var temperatureBroadcastData readWriteModel.TemperatureBroadcastData
			switch salCommand {
			case readWriteModel.TemperatureBroadcastCommandType_BROADCAST_EVENT.PLC4XEnumName():
				if value == nil || !value.IsList() || len(value.GetList()) != 2 || !value.GetList()[0].IsByte() || !value.GetList()[1].IsByte() {
					return nil, false, false, false, errors.Errorf("%s requires exactly 2 arguments [temperatureGroup,temperatureByte]", salCommand)
				}
				commandTypeContainer := readWriteModel.TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes
				temperatureGroup := value.GetList()[0].GetByte()
				temperatureByte := value.GetList()[1].GetByte()
				temperatureBroadcastData = readWriteModel.NewTemperatureBroadcastData(commandTypeContainer, temperatureGroup, temperatureByte)
				supportsWrite = true
			default:
				return nil, false, false, false, errors.Errorf("Unsupported command %s for %s", salCommand, tagType.application.ApplicationId())
			}
			salData = readWriteModel.NewSALDataTemperatureBroadcast(temperatureBroadcastData, nil)
		case readWriteModel.ApplicationId_ROOM_CONTROL_SYSTEM:
			panic("Implement me")
		case
			readWriteModel.ApplicationId_LIGHTING,
			readWriteModel.ApplicationId_VENTILATION,
			readWriteModel.ApplicationId_IRRIGATION_CONTROL,
			readWriteModel.ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL,
			readWriteModel.ApplicationId_HEATING,
			readWriteModel.ApplicationId_AUDIO_AND_VIDEO,
			readWriteModel.ApplicationId_HVAC_ACTUATOR:
			var lightingData readWriteModel.LightingData
			switch salCommand {
			case readWriteModel.LightingCommandType_OFF.PLC4XEnumName():
				commandTypeContainer := readWriteModel.LightingCommandTypeContainer_LightingCommandOff
				if value == nil || !value.IsByte() {
					return nil, false, false, false, errors.Errorf("%s requires exactly 1 arguments [group]", salCommand)
				}
				group := value.GetByte()
				lightingData = readWriteModel.NewLightingDataOff(group, commandTypeContainer)
				supportsWrite = true
			case readWriteModel.LightingCommandType_ON.PLC4XEnumName():
				commandTypeContainer := readWriteModel.LightingCommandTypeContainer_LightingCommandOn
				if value == nil || (!value.IsByte() && (!value.IsList() || len(value.GetList()) != 1 || value.GetList()[0].IsByte())) {
					return nil, false, false, false, errors.Errorf("%s requires exactly 1 arguments [group]", salCommand)
				}
				group := value.GetByte()
				lightingData = readWriteModel.NewLightingDataOn(group, commandTypeContainer)
				supportsWrite = true
			case readWriteModel.LightingCommandType_RAMP_TO_LEVEL.PLC4XEnumName():
				if value == nil || !value.IsList() || len(value.GetList()) != 3 || !value.GetList()[0].IsString() || !value.GetList()[1].IsByte() || !value.GetList()[2].IsByte() {
					return nil, false, false, false, errors.Errorf("%s requires exactly 2 arguments [delay,group,level]", salCommand)
				}
				commandTypeContainer, ok := readWriteModel.LightingCommandTypeContainerByName(fmt.Sprintf("LightingCommandRampToLevel_%s", value.GetList()[0].GetString()))
				if !ok {
					var possibleValues []string
					for _, v := range readWriteModel.LightingCommandTypeContainerValues {
						possibleValues = append(possibleValues, strings.TrimPrefix(v.String(), "LightingCommandRampToLevel_"))
					}
					return nil, false, false, false, errors.Errorf("No level found for %s. Possible values %s", value.GetList()[0].GetString(), possibleValues)
				}
				group := value.GetList()[1].GetByte()
				level := value.GetList()[2].GetByte()
				lightingData = readWriteModel.NewLightingDataRampToLevel(group, level, commandTypeContainer)
				supportsWrite = true
			case readWriteModel.LightingCommandType_TERMINATE_RAMP.PLC4XEnumName():
				commandTypeContainer := readWriteModel.LightingCommandTypeContainer_LightingCommandTerminateRamp
				if value == nil || !value.IsByte() {
					return nil, false, false, false, errors.Errorf("%s requires exactly 1 arguments [group]", salCommand)
				}
				group := value.GetByte()
				lightingData = readWriteModel.NewLightingDataTerminateRamp(group, commandTypeContainer)
				supportsWrite = true
			case readWriteModel.LightingCommandType_LABEL.PLC4XEnumName():
				panic("Implement me")
			default:
				return nil, false, false, false, errors.Errorf("Unsupported command %s for %s", salCommand, tagType.application.ApplicationId())
			}
			salData = readWriteModel.NewSALDataLighting(lightingData, nil)
		case readWriteModel.ApplicationId_AIR_CONDITIONING:
			panic("Implement me")
		case readWriteModel.ApplicationId_TRIGGER_CONTROL:
			panic("Implement me")
		case readWriteModel.ApplicationId_ENABLE_CONTROL:
			panic("Implement me")
		case readWriteModel.ApplicationId_SECURITY:
			panic("Implement me")
		case readWriteModel.ApplicationId_METERING:
			panic("Implement me")
		case readWriteModel.ApplicationId_ACCESS_CONTROL:
			panic("Implement me")
		case readWriteModel.ApplicationId_CLOCK_AND_TIMEKEEPING:
			panic("Implement me")
		case readWriteModel.ApplicationId_TELEPHONY_STATUS_AND_CONTROL:
			panic("Implement me")
		case readWriteModel.ApplicationId_MEASUREMENT:
			panic("Implement me")
		case readWriteModel.ApplicationId_TESTING:
			panic("Implement me")
		case readWriteModel.ApplicationId_MEDIA_TRANSPORT_CONTROL:
			panic("Implement me")
		case readWriteModel.ApplicationId_ERROR_REPORTING:
			panic("Implement me")
		default:
			return nil, false, false, false, errors.Errorf("No support for %s", tagType.application)
		}
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToMultiPointCommandNormal(tagType.application, salData, 0x00, cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToMultiPoint(command, header, cbusOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)
		cBusMessage = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions)
		return
	default:
		return nil, false, false, false, errors.Errorf("Unsupported type %T", tagType)
	}
}
