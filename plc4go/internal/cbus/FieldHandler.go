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
	"encoding/hex"
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type FieldType uint8

//go:generate stringer -type FieldType
const (
	STATUS FieldType = iota
	// TODO: implement
	CAL_RESET
	CAL_RECALL
	CAL_IDENTIFY
	CAL_GETSTATUS
	// TODO: implement
	CAL_WRITE
	// TODO: implement
	CAL_IDENTIFY_REPLY
	// TODO: implement
	CAL_STATUS
	// TODO: implement
	CAL_STATUS_EXTENDED
	SAL
	SAL_MONITOR
	MMI_STATUS_MONITOR
	UNIT_INFO
)

func (i FieldType) GetName() string {
	return i.String()
}

type FieldHandler struct {
	statusRequestPattern *regexp.Regexp
	calPattern           *regexp.Regexp
	salPattern           *regexp.Regexp
	salMonitorPattern    *regexp.Regexp
	mmiMonitorPattern    *regexp.Regexp
	unityQuery           *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		statusRequestPattern: regexp.MustCompile(`^status/(?P<statusRequestType>(?P<binary>binary)|level=0x(?P<startingGroupAddressLabel>00|20|40|60|80|A0|C0|E0))/(?P<application>.*)`),
		calPattern:           regexp.MustCompile(`^cal/(?P<unitAddress>.+)/(?P<calType>reset|recall=\[(?P<recallParamNo>\w+), ?(?P<recallCount>\d+)]|identify=(?P<identifyAttribute>\w+), ?(?P<getstatusCount>\d+)|getstatus=(?P<getstatusParamNo>\w+)|write=\[(?P<writeParamNo>\w+), ?(?P<writeCode>0[xX][0-9a-fA-F][0-9a-fA-F])]|identifyReply=(?P<replyAttribute>\w+)|reply=(?P<replyParamNo>\w+)|status=(?P<statusApplication>.*)|statusExtended=(?P<statusExtendedApplication>.*))`),
		salPattern:           regexp.MustCompile(`^sal/(?P<application>.*)/(?P<salCommand>.*)`),
		salMonitorPattern:    regexp.MustCompile(`^salmonitor/(?P<unitAddress>.+)/(?P<application>.+)`),
		mmiMonitorPattern:    regexp.MustCompile(`^mmimonitor/(?P<unitAddress>.+)/(?P<application>.+)`),
		unityQuery:           regexp.MustCompile(`^info/(?P<unitAddress>.+)/(?P<identifyAttribute>.+)`),
	}
}

type CommandAndArgumentsCount interface {
	fmt.Stringer
	PLC4XEnumName() string
	NumberOfArguments() uint8
}

var PossibleSalCommands = map[readWriteModel.ApplicationId][]CommandAndArgumentsCount{
	readWriteModel.ApplicationId_RESERVED:                           nil, // TODO: Not yet implemented
	readWriteModel.ApplicationId_FREE_USAGE:                         nil, // TODO: Not yet implemented
	readWriteModel.ApplicationId_TEMPERATURE_BROADCAST:              c2nl(readWriteModel.TemperatureBroadcastCommandTypeValues),
	readWriteModel.ApplicationId_ROOM_CONTROL_SYSTEM:                nil, // TODO: Not yet implemented
	readWriteModel.ApplicationId_LIGHTING:                           c2nl(readWriteModel.LightingCommandTypeValues),
	readWriteModel.ApplicationId_VENTILATION:                        c2nl(readWriteModel.LightingCommandTypeValues),
	readWriteModel.ApplicationId_IRRIGATION_CONTROL:                 c2nl(readWriteModel.LightingCommandTypeValues),
	readWriteModel.ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL: c2nl(readWriteModel.LightingCommandTypeValues),
	readWriteModel.ApplicationId_HEATING:                            c2nl(readWriteModel.LightingCommandTypeValues),
	readWriteModel.ApplicationId_AIR_CONDITIONING:                   c2nl(readWriteModel.AirConditioningCommandTypeValues),
	readWriteModel.ApplicationId_TRIGGER_CONTROL:                    c2nl(readWriteModel.TriggerControlCommandTypeValues),
	readWriteModel.ApplicationId_ENABLE_CONTROL:                     c2nl(readWriteModel.EnableControlCommandTypeValues),
	readWriteModel.ApplicationId_AUDIO_AND_VIDEO:                    c2nl(readWriteModel.LightingCommandTypeValues),
	readWriteModel.ApplicationId_SECURITY:                           c2nl(readWriteModel.SecurityCommandTypeValues),
	readWriteModel.ApplicationId_METERING:                           c2nl(readWriteModel.MeteringCommandTypeValues),
	readWriteModel.ApplicationId_ACCESS_CONTROL:                     c2nl(readWriteModel.AccessControlCommandTypeValues),
	readWriteModel.ApplicationId_CLOCK_AND_TIMEKEEPING:              c2nl(readWriteModel.ClockAndTimekeepingCommandTypeValues),
	readWriteModel.ApplicationId_TELEPHONY_STATUS_AND_CONTROL:       c2nl(readWriteModel.TelephonyCommandTypeValues),
	readWriteModel.ApplicationId_MEASUREMENT:                        c2nl(readWriteModel.MeasurementCommandTypeValues),
	readWriteModel.ApplicationId_TESTING:                            nil, // TODO: Not yet implemented
	readWriteModel.ApplicationId_MEDIA_TRANSPORT_CONTROL:            c2nl(readWriteModel.MediaTransportControlCommandTypeValues),
	readWriteModel.ApplicationId_ERROR_REPORTING:                    c2nl(readWriteModel.ErrorReportingCommandTypeValues),
	readWriteModel.ApplicationId_HVAC_ACTUATOR:                      c2nl(readWriteModel.LightingCommandTypeValues),
}

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.statusRequestPattern, query); match != nil {
		return m.handleStatusRequestPattern(match)
	} else if match := utils.GetSubgroupMatches(m.calPattern, query); match != nil {
		return m.handleCalPattern(match)
	} else if match := utils.GetSubgroupMatches(m.salPattern, query); match != nil {
		return m.handleSALPattern(match)
	} else if match := utils.GetSubgroupMatches(m.salMonitorPattern, query); match != nil {
		return m.handleSALMonitorPattern(match)
	} else if match := utils.GetSubgroupMatches(m.mmiMonitorPattern, query); match != nil {
		return m.handleMMIMonitorPattern(match)
	} else if match := utils.GetSubgroupMatches(m.unityQuery, query); match != nil {
		return m.handleUnitQuery(match)
	} else {
		return nil, errors.Errorf("Unable to parse %s", query)
	}
}

func (m FieldHandler) handleStatusRequestPattern(match map[string]string) (model.PlcField, error) {
	var startingGroupAddressLabel *byte
	var statusRequestType StatusRequestType
	statusRequestArgument := match["statusRequestType"]
	if statusRequestArgument != "" {
		if match["binary"] != "" {
			statusRequestType = StatusRequestTypeBinaryState
		} else if levelArgument := match["startingGroupAddressLabel"]; levelArgument != "" {
			statusRequestType = StatusRequestTypeLevel
			startingGroupAddressLabelArgument := match["startingGroupAddressLabel"]
			decodedHex, _ := hex.DecodeString(startingGroupAddressLabelArgument)
			if len(decodedHex) != 1 {
				panic("invalid state. Should have exactly 1")
			}
			startingGroupAddressLabel = &decodedHex[0]
		} else {
			return nil, errors.Errorf("Unknown statusRequestType%s", statusRequestArgument)
		}
	}
	application, err := applicationIdFromArgument(match["application"])
	if err != nil {
		return nil, errors.Wrap(err, "Error getting application id from argument")
	}
	return NewStatusField(statusRequestType, startingGroupAddressLabel, application, 1), nil
}

func (m FieldHandler) handleCalPattern(match map[string]string) (model.PlcField, error) {
	var unitAddress readWriteModel.UnitAddress
	unitAddressArgument := match["unitAddress"]
	if strings.HasPrefix(unitAddressArgument, "0x") {
		decodedHex, err := hex.DecodeString(unitAddressArgument[2:])
		if err != nil {
			return nil, errors.Wrap(err, "Not a valid hex")
		}
		if len(decodedHex) != 1 {
			return nil, errors.Errorf("Hex must be exatly one byte")
		}
		unitAddress = readWriteModel.NewUnitAddress(decodedHex[0])
	} else {
		atoi, err := strconv.ParseUint(unitAddressArgument, 10, 8)
		if err != nil {
			return nil, errors.Errorf("Unknown unit address %s", unitAddressArgument)
		}
		unitAddress = readWriteModel.NewUnitAddress(byte(atoi))
	}

	calTypeArgument := match["calType"]
	switch {
	case strings.HasPrefix(calTypeArgument, "reset"):
		panic("Not implemented") // TODO: implement me
	case strings.HasPrefix(calTypeArgument, "recall="):
		var recalParamNo readWriteModel.Parameter
		recallParamNoArgument := match["recallParamNo"]
		if strings.HasPrefix(recallParamNoArgument, "0x") {
			decodedHex, err := hex.DecodeString(recallParamNoArgument[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Not a valid hex")
			}
			if len(decodedHex) != 1 {
				return nil, errors.Errorf("Hex must be exatly one byte")
			}
			recalParamNo = readWriteModel.Parameter(decodedHex[0])
		} else {

			if atoi, err := strconv.ParseUint(recallParamNoArgument, 10, 8); err == nil {
				recalParamNo = readWriteModel.Parameter(atoi)
			} else {
				parameterByName, ok := readWriteModel.ParameterByName(recallParamNoArgument)
				if !ok {
					return nil, errors.Errorf("Unknown recallParamNo %s", recallParamNoArgument)
				}
				recalParamNo = parameterByName
			}
		}
		var count uint8
		atoi, err := strconv.ParseUint(match["recallCount"], 10, 8)
		if err != nil {
			return nil, errors.Wrap(err, "recallCount not a valid number")
		}
		count = uint8(atoi)
		return NewCALRecallField(unitAddress, recalParamNo, count, 1), nil
	case strings.HasPrefix(calTypeArgument, "identify="):
		var attribute readWriteModel.Attribute
		attributeArgument := match["identifyAttribute"]
		if strings.HasPrefix(attributeArgument, "0x") {
			decodedHex, err := hex.DecodeString(attributeArgument[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Not a valid hex")
			}
			if len(decodedHex) != 1 {
				return nil, errors.Errorf("Hex must be exatly one byte")
			}
			attribute = readWriteModel.Attribute(decodedHex[0])
		} else {
			if atoi, err := strconv.ParseUint(attributeArgument, 10, 8); err == nil {
				attribute = readWriteModel.Attribute(atoi)
			} else {
				parameterByName, ok := readWriteModel.AttributeByName(attributeArgument)
				if !ok {
					return nil, errors.Errorf("Unknown attributeArgument %s", attributeArgument)
				}
				attribute = parameterByName
			}
		}
		return NewCALIdentifyField(unitAddress, attribute, 1), nil
	case strings.HasPrefix(calTypeArgument, "getstatus="):
		var recalParamNo readWriteModel.Parameter
		recallParamNoArgument := match["getstatusParamNo"]
		if strings.HasPrefix(recallParamNoArgument, "0x") {
			decodedHex, err := hex.DecodeString(recallParamNoArgument[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Not a valid hex")
			}
			if len(decodedHex) != 1 {
				return nil, errors.Errorf("Hex must be exatly one byte")
			}
			recalParamNo = readWriteModel.Parameter(decodedHex[0])
		} else {
			if atoi, err := strconv.ParseUint(recallParamNoArgument, 10, 8); err == nil {
				recalParamNo = readWriteModel.Parameter(atoi)
			} else {
				parameterByName, ok := readWriteModel.ParameterByName(recallParamNoArgument)
				if !ok {
					return nil, errors.Errorf("Unknown getstatusParamNo %s", recallParamNoArgument)
				}
				recalParamNo = parameterByName
			}
		}
		var count uint8
		atoi, err := strconv.ParseUint(match["getstatusCount"], 10, 8)
		if err != nil {
			return nil, errors.Wrap(err, "getstatusCount not a valid number")
		}
		count = uint8(atoi)
		return NewCALGetstatusField(unitAddress, recalParamNo, count, 1), nil
	case strings.HasPrefix(calTypeArgument, "write="):
		panic("Not implemented") // TODO: implement me
	case strings.HasPrefix(calTypeArgument, "identifyReply="):
		panic("Not implemented") // TODO: implement me
	case strings.HasPrefix(calTypeArgument, "reply="):
		panic("Not implemented") // TODO: implement me
	case strings.HasPrefix(calTypeArgument, "status="):
		panic("Not implemented") // TODO: implement me
	case strings.HasPrefix(calTypeArgument, "statusExtended="):
		panic("Not implemented") // TODO: implement me
	default:
		return nil, errors.Errorf("Invalid cal type %s", calTypeArgument)
	}
}

func (m FieldHandler) handleSALPattern(match map[string]string) (model.PlcField, error) {
	application, err := applicationIdFromArgument(match["application"])
	if err != nil {
		return nil, errors.Wrap(err, "Error getting application id from argument")
	}
	salCommand := match["salCommand"]
	if salCommand == "" {
		return nil, errors.Wrap(err, "Error getting salCommand from argument")
	}
	isValid := false
	numElements := uint16(0)
	for _, request := range PossibleSalCommands[application.ApplicationId()] {
		if salCommand == request.PLC4XEnumName() {
			isValid = true
			numElements = uint16(request.NumberOfArguments())
			break
		}
	}
	if !isValid {
		return nil, errors.Errorf("Invalid sal command %s for %s. Allowed requests: %s", salCommand, application, PossibleSalCommands[application.ApplicationId()])
	}
	return NewSALField(application, salCommand, numElements), nil
}

func (m FieldHandler) handleSALMonitorPattern(match map[string]string) (model.PlcField, error) {
	var unitAddress readWriteModel.UnitAddress
	{
		unitAddressArgument := match["unitAddress"]
		if unitAddressArgument == "*" {
			unitAddress = nil
		} else if strings.HasPrefix(unitAddressArgument, "0x") {
			decodedHex, err := hex.DecodeString(unitAddressArgument[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Not a valid hex")
			}
			if len(decodedHex) != 1 {
				return nil, errors.Errorf("Hex must be exatly one byte")
			}
			unitAddress = readWriteModel.NewUnitAddress(decodedHex[0])
		} else {
			atoi, err := strconv.ParseUint(unitAddressArgument, 10, 8)
			if err != nil {
				return nil, errors.Errorf("Unknown unit address %s", unitAddressArgument)
			}
			unitAddress = readWriteModel.NewUnitAddress(byte(atoi))
		}
	}

	var application readWriteModel.ApplicationIdContainer
	{
		applicationIdArgument := match["application"]
		if applicationIdArgument == "*" {
			application = readWriteModel.ApplicationIdContainer_RESERVED_FF
		} else {
			var err error
			application, err = applicationIdFromArgument(applicationIdArgument)
			if err != nil {
				return nil, errors.Wrap(err, "Error getting application id from argument")
			}
		}
	}

	return NewSALMonitorField(unitAddress, application, 1), nil
}

func (m FieldHandler) handleMMIMonitorPattern(match map[string]string) (model.PlcField, error) {
	var unitAddress readWriteModel.UnitAddress
	{
		unitAddressArgument := match["unitAddress"]
		if unitAddressArgument == "*" {
			unitAddress = nil
		} else if strings.HasPrefix(unitAddressArgument, "0x") {
			decodedHex, err := hex.DecodeString(unitAddressArgument[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Not a valid hex")
			}
			if len(decodedHex) != 1 {
				return nil, errors.Errorf("Hex must be exatly one byte")
			}
			unitAddress = readWriteModel.NewUnitAddress(decodedHex[0])
		} else {
			atoi, err := strconv.ParseUint(unitAddressArgument, 10, 8)
			if err != nil {
				return nil, errors.Errorf("Unknown unit address %s", unitAddressArgument)
			}
			unitAddress = readWriteModel.NewUnitAddress(byte(atoi))
		}
	}

	var application readWriteModel.ApplicationIdContainer
	{
		applicationIdArgument := match["application"]
		if applicationIdArgument == "*" {
			application = readWriteModel.ApplicationIdContainer_RESERVED_FF
		} else {
			var err error
			application, err = applicationIdFromArgument(applicationIdArgument)
			if err != nil {
				return nil, errors.Wrap(err, "Error getting application id from argument")
			}
		}
	}

	return NewMMIMonitorField(unitAddress, application, 1), nil
}

func (m FieldHandler) handleUnitQuery(match map[string]string) (model.PlcField, error) {
	var unitAddress *readWriteModel.UnitAddress
	unitAddressArgument := match["unitAddress"]
	if unitAddressArgument == "*" {
		unitAddress = nil
	} else if strings.HasPrefix(unitAddressArgument, "0x") {
		decodedHex, err := hex.DecodeString(unitAddressArgument[2:])
		if err != nil {
			return nil, errors.Wrap(err, "Not a valid hex")
		}
		if len(decodedHex) != 1 {
			return nil, errors.Errorf("Hex must be exatly one byte")
		}
		var unitAddressVar readWriteModel.UnitAddress
		unitAddressVar = readWriteModel.NewUnitAddress(decodedHex[0])
		unitAddress = &unitAddressVar
	} else {
		atoi, err := strconv.ParseUint(unitAddressArgument, 10, 8)
		if err != nil {
			return nil, errors.Errorf("Unknown unit address %s", unitAddressArgument)
		}
		var unitAddressVar readWriteModel.UnitAddress
		unitAddressVar = readWriteModel.NewUnitAddress(byte(atoi))
		unitAddress = &unitAddressVar
	}

	var attribute *readWriteModel.Attribute
	attributeArgument := match["identifyAttribute"]
	if attributeArgument == "*" {
		attribute = nil
	} else if strings.HasPrefix(attributeArgument, "0x") {
		decodedHex, err := hex.DecodeString(attributeArgument[2:])
		if err != nil {
			return nil, errors.Wrap(err, "Not a valid hex")
		}
		if len(decodedHex) != 1 {
			return nil, errors.Errorf("Hex must be exatly one byte")
		}
		var attributeVar readWriteModel.Attribute
		attributeVar = readWriteModel.Attribute(decodedHex[0])
		attribute = &attributeVar
	} else {
		if atoi, err := strconv.ParseUint(attributeArgument, 10, 8); err == nil {
			var attributeVar readWriteModel.Attribute
			attributeVar = readWriteModel.Attribute(atoi)
			attribute = &attributeVar
		} else {
			parameterByName, ok := readWriteModel.AttributeByName(attributeArgument)
			if !ok {
				return nil, errors.Errorf("Unknown attributeArgument %s", attributeArgument)
			}
			var attributeVar readWriteModel.Attribute
			attributeVar = parameterByName
			attribute = &attributeVar
		}
	}
	return NewUnitInfoField(unitAddress, attribute, 1), nil
}

func applicationIdFromArgument(applicationIdArgument string) (readWriteModel.ApplicationIdContainer, error) {
	if strings.HasPrefix(applicationIdArgument, "0x") {
		decodedHex, err := hex.DecodeString(applicationIdArgument[2:])
		if err != nil {
			return 0, errors.Wrap(err, "Not a valid hex")
		}
		if len(decodedHex) != 1 {
			return 0, errors.Errorf("Hex must be exatly one byte")
		}
		return readWriteModel.ApplicationIdContainer(decodedHex[0]), nil
	}
	if atoi, err := strconv.ParseUint(applicationIdArgument, 10, 8); err == nil {
		return readWriteModel.ApplicationIdContainer(atoi), nil
	}
	// We try first the application id
	applicationId, ok := readWriteModel.ApplicationIdByName(applicationIdArgument)
	if ok {
		switch applicationId {
		case readWriteModel.ApplicationId_TEMPERATURE_BROADCAST:
			return readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19, nil
		case readWriteModel.ApplicationId_ROOM_CONTROL_SYSTEM:
			return readWriteModel.ApplicationIdContainer_ROOM_CONTROL_SYSTEM_26, nil
		case readWriteModel.ApplicationId_LIGHTING:
			return readWriteModel.ApplicationIdContainer_LIGHTING_38, nil
		case readWriteModel.ApplicationId_VENTILATION:
			return readWriteModel.ApplicationIdContainer_VENTILATION_70, nil
		case readWriteModel.ApplicationId_IRRIGATION_CONTROL:
			return readWriteModel.ApplicationIdContainer_IRRIGATION_CONTROL_71, nil
		case readWriteModel.ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL:
			return readWriteModel.ApplicationIdContainer_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL_72, nil
		case readWriteModel.ApplicationId_HEATING:
			return readWriteModel.ApplicationIdContainer_HEATING_88, nil
		case readWriteModel.ApplicationId_AIR_CONDITIONING:
			return readWriteModel.ApplicationIdContainer_AIR_CONDITIONING_AC, nil
		case readWriteModel.ApplicationId_TRIGGER_CONTROL:
			return readWriteModel.ApplicationIdContainer_TRIGGER_CONTROL_CA, nil
		case readWriteModel.ApplicationId_ENABLE_CONTROL:
			return readWriteModel.ApplicationIdContainer_ENABLE_CONTROL_CB, nil
		case readWriteModel.ApplicationId_AUDIO_AND_VIDEO:
			return readWriteModel.ApplicationIdContainer_AUDIO_AND_VIDEO_CD, nil
		case readWriteModel.ApplicationId_SECURITY:
			return readWriteModel.ApplicationIdContainer_SECURITY_D0, nil
		case readWriteModel.ApplicationId_METERING:
			return readWriteModel.ApplicationIdContainer_METERING_D1, nil
		case readWriteModel.ApplicationId_ACCESS_CONTROL:
			return readWriteModel.ApplicationIdContainer_ACCESS_CONTROL_D5, nil
		case readWriteModel.ApplicationId_CLOCK_AND_TIMEKEEPING:
			return readWriteModel.ApplicationIdContainer_CLOCK_AND_TIMEKEEPING_DF, nil
		case readWriteModel.ApplicationId_TELEPHONY_STATUS_AND_CONTROL:
			return readWriteModel.ApplicationIdContainer_TELEPHONY_STATUS_AND_CONTROL_E0, nil
		case readWriteModel.ApplicationId_MEASUREMENT:
			return readWriteModel.ApplicationIdContainer_MEASUREMENT_E4, nil
		case readWriteModel.ApplicationId_TESTING:
			return readWriteModel.ApplicationIdContainer_TESTING_FA, nil
		case readWriteModel.ApplicationId_MEDIA_TRANSPORT_CONTROL:
			return readWriteModel.ApplicationIdContainer_MEDIA_TRANSPORT_CONTROL_C0, nil
		case readWriteModel.ApplicationId_ERROR_REPORTING:
			return readWriteModel.ApplicationIdContainer_ERROR_REPORTING_CE, nil
		case readWriteModel.ApplicationId_HVAC_ACTUATOR:
			return readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_73, nil
		default:
			return 0, errors.Errorf("%s can't be used directly... select proper application id container", applicationId)
		}
	} else {
		applicationIdByName, ok := readWriteModel.ApplicationIdContainerByName(applicationIdArgument)
		if !ok {
			return 0, errors.Errorf("Unknown applicationId%s", applicationIdArgument)
		}
		return applicationIdByName, nil
	}
}

func c2nl[T CommandAndArgumentsCount](t []T) []CommandAndArgumentsCount {
	result := make([]CommandAndArgumentsCount, len(t))
	for i, e := range t {
		result[i] = e
	}
	return result
}
