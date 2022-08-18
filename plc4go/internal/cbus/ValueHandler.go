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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
	"reflect"
)

type ValueHandler struct {
	spiValues.IEC61131ValueHandler
}

func NewValueHandler() ValueHandler {
	return ValueHandler{}
}

func (m ValueHandler) NewPlcValue(field apiModel.PlcField, value interface{}) (apiValues.PlcValue, error) {
	switch field.GetTypeName() {
	case
		CAL_WRITE.GetName(),
		CAL_IDENTIFY_REPLY.GetName(),
		CAL_STATUS.GetName(),
		CAL_STATUS_EXTENDED.GetName():
		panic("implement me")
	case SAL.GetName():
		var curValues []any
		if field.GetQuantity() > 1 {
			s := reflect.ValueOf(value)
			if s.Kind() != reflect.Slice {
				return nil, errors.New("couldn't cast value to []interface{}")
			}
			curValues = make([]interface{}, s.Len())
			for i := 0; i < s.Len(); i++ {
				curValues[i] = s.Index(i).Interface()
			}
		} else {
			curValues = append(curValues, value)
		}

		field := field.(*salField)
		salCommand := field.salCommand
		switch field.application.ApplicationId() {
		case readWriteModel.ApplicationId_FREE_USAGE:
			panic("Not yet implemented") // TODO: implement
		case readWriteModel.ApplicationId_TEMPERATURE_BROADCAST:
			switch salCommand {
			case readWriteModel.TemperatureBroadcastCommandType_BROADCAST_EVENT.PLC4XEnumName():
				if len(curValues) != 2 {
					return nil, errors.Errorf("%s requires exactly 2 arguments [temperatureGroup,temperatureByte]", salCommand)
				}
				temperatureGroup, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[0])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for temperatureGroup")
				}
				temperatureByte, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[1])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for temperatureByte")
				}
				return spiValues.NewPlcList([]apiValues.PlcValue{temperatureGroup, temperatureByte}), nil
			default:
				return nil, errors.Errorf("Unsupported command %s for %s", salCommand, field.application.ApplicationId())
			}
		case readWriteModel.ApplicationId_ROOM_CONTROL_SYSTEM:
			panic("Implement me")
		case readWriteModel.ApplicationId_LIGHTING,
			readWriteModel.ApplicationId_VENTILATION,
			readWriteModel.ApplicationId_IRRIGATION_CONTROL,
			readWriteModel.ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL,
			readWriteModel.ApplicationId_HEATING,
			readWriteModel.ApplicationId_AUDIO_AND_VIDEO,
			readWriteModel.ApplicationId_HVAC_ACTUATOR:
			switch salCommand {
			case readWriteModel.LightingCommandType_OFF.PLC4XEnumName():
				if len(curValues) != 1 {
					return nil, errors.Errorf("%s requires exactly 1 arguments [groupe]", salCommand)
				}
				group, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[0])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for group")
				}
				return group, nil
			case readWriteModel.LightingCommandType_ON.PLC4XEnumName():
				if len(curValues) != 1 {
					return nil, errors.Errorf("%s requires exactly 1 arguments [groupe]", salCommand)
				}
				group, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[0])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for group")
				}
				return group, nil
			case readWriteModel.LightingCommandType_RAMP_TO_LEVEL.PLC4XEnumName():
				if len(curValues) != 2 {
					return nil, errors.Errorf("%s requires exactly 2 arguments [group,level]", salCommand)
				}
				group, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[0])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for group")
				}
				level, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[0])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for level")
				}
				return spiValues.NewPlcList([]apiValues.PlcValue{group, level}), nil
			case readWriteModel.LightingCommandType_TERMINATE_RAMP.PLC4XEnumName():
				if len(curValues) != 1 {
					return nil, errors.Errorf("%s requires exactly 1 arguments [groupe]", salCommand)
				}
				group, err := m.IEC61131ValueHandler.NewPlcValueFromType(spiValues.IEC61131_BYTE, curValues[0])
				if err != nil {
					return nil, errors.Wrap(err, "error creating value for group")
				}
				return group, nil
			case readWriteModel.LightingCommandType_LABEL.PLC4XEnumName():
				panic("Implement me")
			default:
				return nil, errors.Errorf("Unsupported command %s for %s", salCommand, field.application.ApplicationId())
			}
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
			return nil, errors.Errorf("No support for %s", field.application)
		}
	}
	return m.IEC61131ValueHandler.NewPlcValue(field, value)
}
