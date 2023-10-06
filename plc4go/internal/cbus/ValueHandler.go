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
	"github.com/apache/plc4x/plc4go/spi/options"
	"reflect"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
)

type ValueHandler struct {
	spiValues.DefaultValueHandler
}

func NewValueHandler(_options ...options.WithOption) ValueHandler {
	return ValueHandler{
		spiValues.NewDefaultValueHandler(_options...),
	}
}

func (m ValueHandler) NewPlcValue(tag apiModel.PlcTag, value any) (apiValues.PlcValue, error) {
	if cbusTag, ok := tag.(Tag); ok {
		switch cbusTag.GetTagType() {
		case
			CAL_WRITE,
			CAL_IDENTIFY_REPLY,
			CAL_STATUS,
			CAL_STATUS_EXTENDED:
			return nil, errors.New("Implement me") //TODO: implement me
		case SAL:
			var curValues []any
			if len(tag.GetArrayInfo()) > 0 && tag.GetArrayInfo()[0].GetSize() > 1 {
				s := reflect.ValueOf(value)
				if s.Kind() != reflect.Slice {
					return nil, errors.New("couldn't cast value to []any")
				}
				curValues = make([]any, s.Len())
				for i := 0; i < s.Len(); i++ {
					curValues[i] = s.Index(i).Interface()
				}
			} else {
				curValues = append(curValues, value)
			}

			tmpSalTag := tag.(*salTag)
			salCommand := tmpSalTag.salCommand
			switch tmpSalTag.application.ApplicationId() {
			case readWriteModel.ApplicationId_FREE_USAGE:
				switch salCommand {
				// TODO:implement
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_TEMPERATURE_BROADCAST:
				switch salCommand {
				case readWriteModel.TemperatureBroadcastCommandType_BROADCAST_EVENT.PLC4XEnumName():
					if len(curValues) != 2 {
						return nil, errors.Errorf("%s requires exactly 2 arguments [temperatureGroup,temperatureByte]", salCommand)
					}
					temperatureGroup, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for temperatureGroup")
					}
					temperatureByte, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[1])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for temperatureByte")
					}
					return spiValues.NewPlcList([]apiValues.PlcValue{temperatureGroup, temperatureByte}), nil
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_ROOM_CONTROL_SYSTEM:
				switch salCommand {
				// TODO:implement
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
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
					group, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for group")
					}
					return group, nil
				case readWriteModel.LightingCommandType_ON.PLC4XEnumName():
					if len(curValues) != 1 {
						return nil, errors.Errorf("%s requires exactly 1 arguments [groupe]", salCommand)
					}
					group, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for group")
					}
					return group, nil
				case readWriteModel.LightingCommandType_RAMP_TO_LEVEL.PLC4XEnumName():
					if len(curValues) != 2 {
						return nil, errors.Errorf("%s requires exactly 2 arguments [group,level]", salCommand)
					}
					group, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for group")
					}
					level, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[1])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for level")
					}
					return spiValues.NewPlcList([]apiValues.PlcValue{group, level}), nil
				case readWriteModel.LightingCommandType_TERMINATE_RAMP.PLC4XEnumName():
					if len(curValues) != 1 {
						return nil, errors.Errorf("%s requires exactly 1 arguments [groupe]", salCommand)
					}
					group, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for group")
					}
					return group, nil
				case readWriteModel.LightingCommandType_LABEL.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_AIR_CONDITIONING:
				switch salCommand {
				case readWriteModel.AirConditioningCommandType_SET_ZONE_GROUP_OFF.PLC4XEnumName():
					zoneGroup, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for zoneGroup")
					}
					return zoneGroup, nil
				case readWriteModel.AirConditioningCommandType_ZONE_HVAC_PLANT_STATUS.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_ZONE_HUMIDITY_PLANT_STATUS.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_ZONE_TEMPERATURE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_ZONE_HUMIDITY.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_REFRESH.PLC4XEnumName():
					zoneGroup, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for zoneGroup")
					}
					return zoneGroup, nil
				case readWriteModel.AirConditioningCommandType_SET_ZONE_HVAC_MODE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_PLANT_HVAC_LEVEL.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_ZONE_HUMIDITY_MODE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_PLANT_HUMIDITY_LEVEL.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_HVAC_UPPER_GUARD_LIMIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_HVAC_LOWER_GUARD_LIMIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_HVAC_SETBACK_LIMIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_HUMIDITY_UPPER_GUARD_LIMIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_HUMIDITY_LOWER_GUARD_LIMIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_SET_ZONE_GROUP_ON.PLC4XEnumName():
					zoneGroup, err := m.DefaultValueHandler.NewPlcValueFromType(apiValues.BYTE, curValues[0])
					if err != nil {
						return nil, errors.Wrap(err, "error creating value for zoneGroup")
					}
					return zoneGroup, nil
				case readWriteModel.AirConditioningCommandType_SET_HUMIDITY_SETBACK_LIMIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_HVAC_SCHEDULE_ENTRY.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AirConditioningCommandType_HUMIDITY_SCHEDULE_ENTRY.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_TRIGGER_CONTROL:
				switch salCommand {
				case readWriteModel.TriggerControlCommandType_TRIGGER_EVENT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.TriggerControlCommandType_TRIGGER_MIN.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.TriggerControlCommandType_TRIGGER_MAX.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.TriggerControlCommandType_INDICATOR_KILL.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.TriggerControlCommandType_LABEL.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_ENABLE_CONTROL:
				switch salCommand {
				case readWriteModel.EnableControlCommandType_SET_NETWORK_VARIABLE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_SECURITY:
				switch salCommand {
				case readWriteModel.SecurityCommandType_OFF.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.SecurityCommandType_ON.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.SecurityCommandType_EVENT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_METERING:
				switch salCommand {
				case readWriteModel.MeteringCommandType_EVENT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_ACCESS_CONTROL:
				switch salCommand {
				case readWriteModel.AccessControlCommandType_CLOSE_ACCESS_POINT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_LOCK_ACCESS_POINT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_ACCESS_POINT_LEFT_OPEN.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_ACCESS_POINT_FORCED_OPEN.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_ACCESS_POINT_CLOSED.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_REQUEST_TO_EXIT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_VALID_ACCESS.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.AccessControlCommandType_INVALID_ACCESS.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_CLOCK_AND_TIMEKEEPING:
				switch salCommand {
				case readWriteModel.ClockAndTimekeepingCommandType_UPDATE_NETWORK_VARIABLE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				case readWriteModel.ClockAndTimekeepingCommandType_REQUEST_REFRESH.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_TELEPHONY_STATUS_AND_CONTROL:
				switch salCommand {
				case readWriteModel.TelephonyCommandType_EVENT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_MEASUREMENT:
				switch salCommand {
				case readWriteModel.MeasurementCommandType_MEASUREMENT_EVENT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_TESTING:
				switch salCommand {
				// TODO:implement
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_MEDIA_TRANSPORT_CONTROL:
				switch salCommand {
				case readWriteModel.MediaTransportControlCommandType_STOP.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_PLAY.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_PAUSE_RESUME.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_SELECT_CATEGORY.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_SELECT_SELECTION.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_SELECT_TRACK.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_SHUFFLE_ON_OFF.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_REPEAT_ON_OFF.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_NEXT_PREVIOUS_CATEGORY.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_NEXT_PREVIOUS_SELECTION.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_NEXT_PREVIOUS_TRACK.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_FAST_FORWARD.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_REWIND.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_SOURCE_POWER_CONTROL.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_TOTAL_TRACKS.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_STATUS_REQUEST.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_ENUMERATE_CATEGORIES_SELECTIONS_TRACKS.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_ENUMERATION_SIZE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_TRACK_NAME.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_SELECTION_NAME.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.MediaTransportControlCommandType_CATEGORY_NAME.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			case readWriteModel.ApplicationId_ERROR_REPORTING:
				switch salCommand {
				case readWriteModel.ErrorReportingCommandType_DEPRECATED.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.ErrorReportingCommandType_ERROR_REPORT.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.ErrorReportingCommandType_ACKNOWLEDGE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				case readWriteModel.ErrorReportingCommandType_CLEAR_MOST_SEVERE.PLC4XEnumName():
					return nil, errors.New("Implement me") //TODO: implement me me
				default:
					return nil, errors.Errorf("Unsupported command %s for %s", salCommand, tmpSalTag.application.ApplicationId())
				}
			default:
				return nil, errors.Errorf("No support for %s", tmpSalTag.application)
			}
		}
	}
	return m.DefaultValueHandler.NewPlcValue(tag, value)
}
