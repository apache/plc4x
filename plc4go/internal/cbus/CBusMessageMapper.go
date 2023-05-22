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
	"github.com/apache/plc4x/plc4go/spi"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/pkg/errors"
)

func TagToCBusMessage(tag apiModel.PlcTag, value apiValues.PlcValue, alphaGenerator *AlphaGenerator, messageCodec *MessageCodec) (cBusMessage readWriteModel.CBusMessage, supportsRead, supportsWrite, supportsSubscribe bool, err error) {
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
		var cbusCommand readWriteModel.CBusCommand
		cbusCommand, err = producePointToMultiPointCommandStatus(tagType.bridgeAddresses, tagType.application, statusRequest, cbusOptions)
		if err != nil {
			return nil, false, false, false, errors.Wrap(err, "error producing point to multipoint command")
		}
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead, supportsSubscribe = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true, true
		return
	case *calRecallTag:
		calData := readWriteModel.NewCALDataRecall(tagType.parameter, tagType.count, readWriteModel.CALCommandTypeContainer_CALCommandRecall, nil, requestContext)
		var command readWriteModel.CBusCommand
		command, err = producePointToPointCommand(tagType.unitAddress, tagType.bridgeAddresses, calData, cbusOptions)
		if err != nil {
			return nil, false, false, false, errors.Wrap(err, "error producing cal command")
		}
		request := readWriteModel.NewRequestCommand(command, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true
		return
	case *calIdentifyTag:
		calData := readWriteModel.NewCALDataIdentify(tagType.attribute, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, requestContext)
		var command readWriteModel.CBusCommand
		command, err = producePointToPointCommand(tagType.unitAddress, tagType.bridgeAddresses, calData, cbusOptions)
		if err != nil {
			return nil, false, false, false, errors.Wrap(err, "error producing cal command")
		}
		request := readWriteModel.NewRequestCommand(command, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

		cBusMessage, supportsRead = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions), true
		return
	case *calGetStatusTag:
		calData := readWriteModel.NewCALDataGetStatus(tagType.parameter, tagType.count, readWriteModel.CALCommandTypeContainer_CALCommandGetStatus, nil, requestContext)
		var command readWriteModel.CBusCommand
		command, err = producePointToPointCommand(tagType.unitAddress, tagType.bridgeAddresses, calData, cbusOptions)
		if err != nil {
			return nil, false, false, false, errors.Wrap(err, "error producing cal command")
		}
		request := readWriteModel.NewRequestCommand(command, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)

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
			err = errors.New("Not yet implemented") // TODO: implement
			return
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
			err = errors.New("Not yet implemented") // TODO: implement
			return
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
					return nil, false, false, false, errors.Errorf("%s requires exactly 3 arguments [delay,group,level]", salCommand)
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
				err = errors.New("Not yet implemented") // TODO: implement
				return
			default:
				return nil, false, false, false, errors.Errorf("Unsupported command %s for %s", salCommand, tagType.application.ApplicationId())
			}
			salData = readWriteModel.NewSALDataLighting(lightingData, nil)
		case readWriteModel.ApplicationId_AIR_CONDITIONING:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_TRIGGER_CONTROL:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_ENABLE_CONTROL:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_SECURITY:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_METERING:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_ACCESS_CONTROL:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_CLOCK_AND_TIMEKEEPING:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_TELEPHONY_STATUS_AND_CONTROL:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_MEASUREMENT:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_TESTING:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_MEDIA_TRANSPORT_CONTROL:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		case readWriteModel.ApplicationId_ERROR_REPORTING:
			err = errors.New("Not yet implemented") // TODO: implement
			return
		default:
			return nil, false, false, false, errors.Errorf("No support for %s", tagType.application)
		}
		var cbusCommand readWriteModel.CBusCommand
		cbusCommand, err = producePointToMultiPointCommandNormal(tagType.bridgeAddresses, tagType.application, salData, cbusOptions)
		if err != nil {
			return nil, false, false, false, errors.Wrap(err, "error producing point to multipoint command")
		}
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), cbusOptions)
		cBusMessage = readWriteModel.NewCBusMessageToServer(request, requestContext, cbusOptions)
		return
	default:
		return nil, false, false, false, errors.Errorf("Unsupported type %T", tagType)
	}
}

func producePointToPointCommand(unitAddress readWriteModel.UnitAddress, bridgeAddresses []readWriteModel.BridgeAddress, calData readWriteModel.CALData, cbusOptions readWriteModel.CBusOptions) (readWriteModel.CBusCommand, error) {
	var command readWriteModel.CBusPointToPointCommand
	numberOfBridgeAddresses := len(bridgeAddresses)
	if numberOfBridgeAddresses > 0 {
		if numberOfBridgeAddresses > 6 {
			return nil, errors.Errorf("Can't have a path longer than 6. Actuall path length = %d", numberOfBridgeAddresses)
		}
		networkRoute := readWriteModel.NewNetworkRoute(readWriteModel.NewNetworkProtocolControlInformation(uint8(numberOfBridgeAddresses), uint8(numberOfBridgeAddresses)), bridgeAddresses[1:])

		command = readWriteModel.NewCBusPointToPointCommandIndirect(bridgeAddresses[0], networkRoute, unitAddress, 0x0000, calData, cbusOptions)
	} else {
		command = readWriteModel.NewCBusPointToPointCommandDirect(unitAddress, 0x0000, calData, cbusOptions)
	}

	header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
	return readWriteModel.NewCBusCommandPointToPoint(command, header, cbusOptions), nil
}

func producePointToMultiPointCommandStatus(bridgeAddresses []readWriteModel.BridgeAddress, application readWriteModel.ApplicationIdContainer, statusRequest readWriteModel.StatusRequest, cbusOptions readWriteModel.CBusOptions) (readWriteModel.CBusCommand, error) {
	numberOfBridgeAddresses := len(bridgeAddresses)
	if numberOfBridgeAddresses > 0 {
		if numberOfBridgeAddresses > 6 {
			return nil, errors.Errorf("Can't have a path longer than 6. Actuall path length = %d", numberOfBridgeAddresses)
		}
		networkRoute := readWriteModel.NewNetworkRoute(readWriteModel.NewNetworkProtocolControlInformation(uint8(numberOfBridgeAddresses), uint8(numberOfBridgeAddresses)), bridgeAddresses[1:])
		command := readWriteModel.NewCBusPointToPointToMultiPointCommandStatus(statusRequest, bridgeAddresses[0], networkRoute, byte(application), cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint)
		return readWriteModel.NewCBusCommandPointToPointToMultiPoint(command, header, cbusOptions), nil
	}
	command := readWriteModel.NewCBusPointToMultiPointCommandStatus(statusRequest, byte(application), cbusOptions)
	header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint)
	return readWriteModel.NewCBusCommandPointToMultiPoint(command, header, cbusOptions), nil
}

func producePointToMultiPointCommandNormal(bridgeAddresses []readWriteModel.BridgeAddress, application readWriteModel.ApplicationIdContainer, salData readWriteModel.SALData, cbusOptions readWriteModel.CBusOptions) (readWriteModel.CBusCommand, error) {
	numberOfBridgeAddresses := len(bridgeAddresses)
	if numberOfBridgeAddresses > 0 {
		if numberOfBridgeAddresses > 6 {
			return nil, errors.Errorf("Can't have a path longer than 6. Actuall path length = %d", numberOfBridgeAddresses)
		}
		networkRoute := readWriteModel.NewNetworkRoute(readWriteModel.NewNetworkProtocolControlInformation(uint8(numberOfBridgeAddresses), uint8(numberOfBridgeAddresses)), bridgeAddresses[1:])
		command := readWriteModel.NewCBusPointToPointToMultiPointCommandNormal(application, salData, bridgeAddresses[0], networkRoute, byte(application), cbusOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint)
		return readWriteModel.NewCBusCommandPointToPointToMultiPoint(command, header, cbusOptions), nil
	}

	command := readWriteModel.NewCBusPointToMultiPointCommandNormal(application, salData, 0x00, cbusOptions)
	header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
	return readWriteModel.NewCBusCommandPointToMultiPoint(command, header, cbusOptions), nil
}

func MapEncodedReply(transaction spi.RequestTransaction, encodedReply readWriteModel.EncodedReply, tagName string, addResponseCode func(name string, responseCode apiModel.PlcResponseCode), addPlcValue func(name string, plcValue apiValues.PlcValue)) error {
	switch reply := encodedReply.(type) {
	case readWriteModel.EncodedReplyCALReplyExactly:
		calData := reply.GetCalReply().GetCalData()
		addResponseCode(tagName, apiModel.PlcResponseCode_OK)
		switch calData := calData.(type) {
		case readWriteModel.CALDataStatusExactly:
			application := calData.GetApplication()
			// TODO: verify application... this should be the same
			_ = application
			blockStart := calData.GetBlockStart()
			// TODO: verify application... this should be the same
			_ = blockStart
			statusBytes := calData.GetStatusBytes()
			addResponseCode(tagName, apiModel.PlcResponseCode_OK)
			plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
			for i, statusByte := range statusBytes {
				plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
				plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
				plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
				plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
			}
			addPlcValue(tagName, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
				"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
				"blockStart":  spiValues.NewPlcBYTE(blockStart),
				"values":      spiValues.NewPlcList(plcListValues),
			}))
		case readWriteModel.CALDataStatusExtendedExactly:
			coding := calData.GetCoding()
			// TODO: verify coding... this should be the same
			_ = coding
			application := calData.GetApplication()
			// TODO: verify application... this should be the same
			_ = application
			blockStart := calData.GetBlockStart()
			// TODO: verify application... this should be the same
			_ = blockStart
			switch coding {
			case readWriteModel.StatusCoding_BINARY_BY_THIS_SERIAL_INTERFACE:
				fallthrough
			case readWriteModel.StatusCoding_BINARY_BY_ELSEWHERE:
				statusBytes := calData.GetStatusBytes()
				addResponseCode(tagName, apiModel.PlcResponseCode_OK)
				plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
				for i, statusByte := range statusBytes {
					plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
					plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
					plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
					plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
				}
				addPlcValue(tagName, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
					"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
					"blockStart":  spiValues.NewPlcBYTE(blockStart),
					"values":      spiValues.NewPlcList(plcListValues),
				}))
			case readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE:
				fallthrough
			case readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE:
				levelInformation := calData.GetLevelInformation()
				addResponseCode(tagName, apiModel.PlcResponseCode_OK)
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
						return transaction.FailRequest(errors.Errorf("Impossible case %v", levelInformation))
					}
				}
				addPlcValue(tagName, spiValues.NewPlcList(plcListValues))
			}
		case readWriteModel.CALDataIdentifyReplyExactly:
			switch identifyReplyCommand := calData.GetIdentifyReplyCommand().(type) {
			case readWriteModel.IdentifyReplyCommandCurrentSenseLevelsExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetCurrentSenseLevels()))
			case readWriteModel.IdentifyReplyCommandDelaysExactly:
				addPlcValue(tagName, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
					"reStrikeDelay": spiValues.NewPlcUSINT(identifyReplyCommand.GetReStrikeDelay()),
					"terminalLevel": spiValues.NewPlcRawByteArray(identifyReplyCommand.GetTerminalLevels()),
				}))
			case readWriteModel.IdentifyReplyCommandDSIStatusExactly:
				addPlcValue(tagName, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
					"channelStatus1":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus1().String()),
					"channelStatus2":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus2().String()),
					"channelStatus3":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus3().String()),
					"channelStatus4":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus4().String()),
					"channelStatus5":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus5().String()),
					"channelStatus6":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus6().String()),
					"channelStatus7":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus7().String()),
					"channelStatus8":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus8().String()),
					"unitStatus":              spiValues.NewPlcSTRING(identifyReplyCommand.GetUnitStatus().String()),
					"dimmingUCRevisionNumber": spiValues.NewPlcUSINT(identifyReplyCommand.GetDimmingUCRevisionNumber()),
				}))
			case readWriteModel.IdentifyReplyCommandExtendedDiagnosticSummaryExactly:
				addPlcValue(tagName, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
					"lowApplication":         spiValues.NewPlcSTRING(identifyReplyCommand.GetLowApplication().String()),
					"highApplication":        spiValues.NewPlcSTRING(identifyReplyCommand.GetHighApplication().String()),
					"area":                   spiValues.NewPlcUSINT(identifyReplyCommand.GetArea()),
					"crc":                    spiValues.NewPlcUINT(identifyReplyCommand.GetCrc()),
					"serialNumber":           spiValues.NewPlcUDINT(identifyReplyCommand.GetSerialNumber()),
					"networkVoltage":         spiValues.NewPlcUSINT(identifyReplyCommand.GetNetworkVoltage()),
					"unitInLearnMode":        spiValues.NewPlcBOOL(identifyReplyCommand.GetUnitInLearnMode()),
					"networkVoltageLow":      spiValues.NewPlcBOOL(identifyReplyCommand.GetNetworkVoltageLow()),
					"networkVoltageMarginal": spiValues.NewPlcBOOL(identifyReplyCommand.GetNetworkVoltageMarginal()),
					"enableChecksumAlarm":    spiValues.NewPlcBOOL(identifyReplyCommand.GetEnableChecksumAlarm()),
					"outputUnit":             spiValues.NewPlcBOOL(identifyReplyCommand.GetOutputUnit()),
					"installationMMIError":   spiValues.NewPlcBOOL(identifyReplyCommand.GetInstallationMMIError()),
					"EEWriteError":           spiValues.NewPlcBOOL(identifyReplyCommand.GetEEWriteError()),
					"EEChecksumError":        spiValues.NewPlcBOOL(identifyReplyCommand.GetEEChecksumError()),
					"EEDataError":            spiValues.NewPlcBOOL(identifyReplyCommand.GetEEDataError()),
					"microReset":             spiValues.NewPlcBOOL(identifyReplyCommand.GetMicroReset()),
					"commsTxError":           spiValues.NewPlcBOOL(identifyReplyCommand.GetCommsTxError()),
					"internalStackOverflow":  spiValues.NewPlcBOOL(identifyReplyCommand.GetInternalStackOverflow()),
					"microPowerReset":        spiValues.NewPlcBOOL(identifyReplyCommand.GetMicroPowerReset()),
				}))
			case readWriteModel.IdentifyReplyCommandSummaryExactly:
				addPlcValue(tagName, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
					"partName":        spiValues.NewPlcSTRING(identifyReplyCommand.GetPartName()),
					"unitServiceType": spiValues.NewPlcUSINT(identifyReplyCommand.GetUnitServiceType()),
					"version":         spiValues.NewPlcSTRING(identifyReplyCommand.GetVersion()),
				}))
			case readWriteModel.IdentifyReplyCommandFirmwareVersionExactly:
				addPlcValue(tagName, spiValues.NewPlcSTRING(identifyReplyCommand.GetFirmwareVersion()))
			case readWriteModel.IdentifyReplyCommandGAVPhysicalAddressesExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetValues()))
			case readWriteModel.IdentifyReplyCommandGAVValuesCurrentExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetValues()))
			case readWriteModel.IdentifyReplyCommandGAVValuesStoredExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetValues()))
			case readWriteModel.IdentifyReplyCommandLogicalAssignmentExactly:
				var plcValues []apiValues.PlcValue
				for _, logicAssigment := range identifyReplyCommand.GetLogicAssigment() {
					plcValues = append(plcValues, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
						"greaterOfOrLogic": spiValues.NewPlcBOOL(logicAssigment.GetGreaterOfOrLogic()),
						"reStrikeDelay":    spiValues.NewPlcBOOL(logicAssigment.GetReStrikeDelay()),
						"assignedToGav16":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav16()),
						"assignedToGav15":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav15()),
						"assignedToGav14":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav14()),
						"assignedToGav13":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav13()),
					}))
				}
				addPlcValue(tagName, spiValues.NewPlcList(plcValues))
			case readWriteModel.IdentifyReplyCommandManufacturerExactly:
				addPlcValue(tagName, spiValues.NewPlcSTRING(identifyReplyCommand.GetManufacturerName()))
			case readWriteModel.IdentifyReplyCommandMaximumLevelsExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetMaximumLevels()))
			case readWriteModel.IdentifyReplyCommandMinimumLevelsExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetMinimumLevels()))
			case readWriteModel.IdentifyReplyCommandNetworkTerminalLevelsExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetNetworkTerminalLevels()))
			case readWriteModel.IdentifyReplyCommandNetworkVoltageExactly:
				volts := identifyReplyCommand.GetVolts()
				voltsFloat, err := strconv.ParseFloat(volts, 0)
				if err != nil {
					addResponseCode(tagName, apiModel.PlcResponseCode_INTERNAL_ERROR)
					return transaction.FailRequest(errors.Wrap(err, "Error parsing volts"))
				}
				voltsDecimalPlace := identifyReplyCommand.GetVoltsDecimalPlace()
				voltsDecimalPlaceFloat, err := strconv.ParseFloat(voltsDecimalPlace, 0)
				if err != nil {
					addResponseCode(tagName, apiModel.PlcResponseCode_INTERNAL_ERROR)
					return transaction.FailRequest(errors.Wrap(err, "Error parsing volts decimal place"))
				}
				voltsFloat += voltsDecimalPlaceFloat / 10
				addPlcValue(tagName, spiValues.NewPlcLREAL(voltsFloat))
			case readWriteModel.IdentifyReplyCommandOutputUnitSummaryExactly:
				unitFlags := identifyReplyCommand.GetUnitFlags()
				structContent := map[string]apiValues.PlcValue{
					"unitFlags": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
						"assertingNetworkBurden": spiValues.NewPlcBOOL(unitFlags.GetAssertingNetworkBurden()),
						"restrikeTimingActive":   spiValues.NewPlcBOOL(unitFlags.GetRestrikeTimingActive()),
						"remoteOFFInputAsserted": spiValues.NewPlcBOOL(unitFlags.GetRemoteOFFInputAsserted()),
						"remoteONInputAsserted":  spiValues.NewPlcBOOL(unitFlags.GetRemoteONInputAsserted()),
						"localToggleEnabled":     spiValues.NewPlcBOOL(unitFlags.GetLocalToggleEnabled()),
						"localToggleActiveState": spiValues.NewPlcBOOL(unitFlags.GetLocalToggleActiveState()),
						"clockGenerationEnabled": spiValues.NewPlcBOOL(unitFlags.GetClockGenerationEnabled()),
						"unitGeneratingClock":    spiValues.NewPlcBOOL(unitFlags.GetUnitGeneratingClock()),
					}),
					"timeFromLastRecoverOfMainsInSeconds": spiValues.NewPlcUSINT(identifyReplyCommand.GetTimeFromLastRecoverOfMainsInSeconds()),
				}
				if gavStoreEnabledByte1 := identifyReplyCommand.GetGavStoreEnabledByte1(); gavStoreEnabledByte1 != nil {
					structContent["gavStoreEnabledByte1"] = spiValues.NewPlcUSINT(*gavStoreEnabledByte1)
				}
				if gavStoreEnabledByte2 := identifyReplyCommand.GetGavStoreEnabledByte2(); gavStoreEnabledByte2 != nil {
					structContent["gavStoreEnabledByte2"] = spiValues.NewPlcUSINT(*gavStoreEnabledByte2)
				}
				addPlcValue(tagName, spiValues.NewPlcStruct(structContent))
			case readWriteModel.IdentifyReplyCommandTerminalLevelsExactly:
				addPlcValue(tagName, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetTerminalLevels()))
			case readWriteModel.IdentifyReplyCommandTypeExactly:
				addPlcValue(tagName, spiValues.NewPlcSTRING(identifyReplyCommand.GetUnitType()))
			default:
				addResponseCode(tagName, apiModel.PlcResponseCode_INVALID_DATA)
				return transaction.FailRequest(errors.Errorf("Unmapped type %T", identifyReplyCommand))
			}
		default:
			wbpcb := spiValues.NewWriteBufferPlcValueBased()
			if err := calData.SerializeWithWriteBuffer(context.Background(), wbpcb); err != nil {
				log.Warn().Err(err).Msgf("Unmapped cal data type %T. Returning raw to string", calData)
				addPlcValue(tagName, spiValues.NewPlcSTRING(fmt.Sprintf("%s", calData)))
			} else {
				addPlcValue(tagName, wbpcb.GetPlcValue())
			}
		}
	default:
		return transaction.FailRequest(errors.Errorf("All types should be mapped here. Not mapped: %T", reply))
	}
	return nil
}
