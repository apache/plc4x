//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

type KnxDatapointType uint16

type IKnxDatapointType interface {
	Text() string
	SizeInBits() uint8
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxDatapointType_DPT_1_BIT                                         KnxDatapointType = 1
	KnxDatapointType_DPT_1_BIT_CONTROLLED                              KnxDatapointType = 2
	KnxDatapointType_DPT_3_BIT_CONTROLLED                              KnxDatapointType = 3
	KnxDatapointType_DPT_CHARACTER                                     KnxDatapointType = 4
	KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE                          KnxDatapointType = 5
	KnxDatapointType_DPT_8_BIT_SIGNED_VALUE                            KnxDatapointType = 6
	KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE                         KnxDatapointType = 7
	KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE                           KnxDatapointType = 8
	KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE                            KnxDatapointType = 9
	KnxDatapointType_DPT_TIME                                          KnxDatapointType = 10
	KnxDatapointType_DPT_DATE                                          KnxDatapointType = 11
	KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE                         KnxDatapointType = 12
	KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE                           KnxDatapointType = 13
	KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE                            KnxDatapointType = 14
	KnxDatapointType_DPT_ENTRANCE_ACCESS                               KnxDatapointType = 15
	KnxDatapointType_DPT_CHARACTER_STRING                              KnxDatapointType = 16
	KnxDatapointType_DPT_SCENE_NUMBER                                  KnxDatapointType = 17
	KnxDatapointType_DPT_SCENE_CONTROL                                 KnxDatapointType = 18
	KnxDatapointType_DPT_DATE_TIME                                     KnxDatapointType = 19
	KnxDatapointType_DPT_1_BYTE                                        KnxDatapointType = 20
	KnxDatapointType_DPT_8_BIT_SET                                     KnxDatapointType = 21
	KnxDatapointType_DPT_16_BIT_SET                                    KnxDatapointType = 22
	KnxDatapointType_DPT_2_BIT_SET                                     KnxDatapointType = 23
	KnxDatapointType_DPT_2_NIBBLE_SET                                  KnxDatapointType = 25
	KnxDatapointType_DPT_8_BIT_SET_2                                   KnxDatapointType = 26
	KnxDatapointType_DPT_32_BIT_SET                                    KnxDatapointType = 27
	KnxDatapointType_DPT_ELECTRICAL_ENERGY                             KnxDatapointType = 29
	KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION                   KnxDatapointType = 30
	KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM          KnxDatapointType = 206
	KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION                        KnxDatapointType = 217
	KnxDatapointType_DPT_ALARM_INFO                                    KnxDatapointType = 219
	KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE                         KnxDatapointType = 222
	KnxDatapointType_DPT_SCALING_SPEED                                 KnxDatapointType = 225
	KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION               KnxDatapointType = 229
	KnxDatapointType_DPT_MBUS_ADDRESS                                  KnxDatapointType = 230
	KnxDatapointType_DPT_3_BYTE_COLOUR_RGB                             KnxDatapointType = 232
	KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1                       KnxDatapointType = 234
	KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY KnxDatapointType = 235
	KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL                      KnxDatapointType = 236
	KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT              KnxDatapointType = 237
	KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT               KnxDatapointType = 238
	KnxDatapointType_DPT_POSITIONS                                     KnxDatapointType = 240
	KnxDatapointType_DPT_STATUS_32_BIT                                 KnxDatapointType = 241
	KnxDatapointType_DPT_STATUS_48_BIT                                 KnxDatapointType = 242
	KnxDatapointType_DPT_CONVERTER_STATUS                              KnxDatapointType = 244
	KnxDatapointType_DPT_CONVERTER_TEST_RESULT                         KnxDatapointType = 245
	KnxDatapointType_DPT_BATTERY_INFORMATION                           KnxDatapointType = 246
	KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION      KnxDatapointType = 249
	KnxDatapointType_DPT_STATUS_24_BIT                                 KnxDatapointType = 250
	KnxDatapointType_DPT_COLOUR_RGBW                                   KnxDatapointType = 251
	KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW                         KnxDatapointType = 252
	KnxDatapointType_DPT_RELATIVE_CONTROL_RGB                          KnxDatapointType = 254
	KnxDatapointType_DPT_F32F32                                        KnxDatapointType = 255
	KnxDatapointType_DPT_F16F16F16F16                                  KnxDatapointType = 275
)

func (e KnxDatapointType) Text() string {
	switch e {
	case 1:
		{ /* '1' */
			return "1-bit"
		}
	case 10:
		{ /* '10' */
			return "time"
		}
	case 11:
		{ /* '11' */
			return "date"
		}
	case 12:
		{ /* '12' */
			return "4-byte unsigned value"
		}
	case 13:
		{ /* '13' */
			return "4-byte signed value"
		}
	case 14:
		{ /* '14' */
			return "4-byte float value"
		}
	case 15:
		{ /* '15' */
			return "entrance access"
		}
	case 16:
		{ /* '16' */
			return "character string"
		}
	case 17:
		{ /* '17' */
			return "scene number"
		}
	case 18:
		{ /* '18' */
			return "scene control"
		}
	case 19:
		{ /* '19' */
			return "Date Time"
		}
	case 2:
		{ /* '2' */
			return "1-bit controlled"
		}
	case 20:
		{ /* '20' */
			return "1-byte"
		}
	case 206:
		{ /* '206' */
			return "16-bit unsigned value & 8-bit enum"
		}
	case 21:
		{ /* '21' */
			return "8-bit set"
		}
	case 217:
		{ /* '217' */
			return "datapoint type version"
		}
	case 219:
		{ /* '219' */
			return "alarm info"
		}
	case 22:
		{ /* '22' */
			return "16-bit set"
		}
	case 222:
		{ /* '222' */
			return "3x 2-byte float value"
		}
	case 225:
		{ /* '225' */
			return "scaling speed"
		}
	case 229:
		{ /* '229' */
			return "4-1-1 byte combined information"
		}
	case 23:
		{ /* '23' */
			return "2-bit set"
		}
	case 230:
		{ /* '230' */
			return "MBus address"
		}
	case 232:
		{ /* '232' */
			return "3-byte colour RGB"
		}
	case 234:
		{ /* '234' */
			return "language code ISO 639-1"
		}
	case 235:
		{ /* '235' */
			return "Signed value with classification and validity"
		}
	case 236:
		{ /* '236' */
			return "Prioritised Mode Control"
		}
	case 237:
		{ /* '237' */
			return "configuration/ diagnostics"
		}
	case 238:
		{ /* '238' */
			return "configuration/ diagnostics"
		}
	case 240:
		{ /* '240' */
			return "positions"
		}
	case 241:
		{ /* '241' */
			return "status"
		}
	case 242:
		{ /* '242' */
			return "status"
		}
	case 244:
		{ /* '244' */
			return "Converter Status"
		}
	case 245:
		{ /* '245' */
			return "Converter test result"
		}
	case 246:
		{ /* '246' */
			return "Battery Information"
		}
	case 249:
		{ /* '249' */
			return "brightness colour temperature transition"
		}
	case 25:
		{ /* '25' */
			return "2-nibble set"
		}
	case 250:
		{ /* '250' */
			return "status"
		}
	case 251:
		{ /* '251' */
			return "Colour RGBW"
		}
	case 252:
		{ /* '252' */
			return "Relative Control RGBW"
		}
	case 254:
		{ /* '254' */
			return "Relative Control RGB"
		}
	case 255:
		{ /* '255' */
			return "F32F32"
		}
	case 26:
		{ /* '26' */
			return "8-bit set"
		}
	case 27:
		{ /* '27' */
			return "32-bit set"
		}
	case 275:
		{ /* '275' */
			return "F16F16F16F16"
		}
	case 29:
		{ /* '29' */
			return "electrical energy"
		}
	case 3:
		{ /* '3' */
			return "3-bit controlled"
		}
	case 30:
		{ /* '30' */
			return "24 times channel activation"
		}
	case 4:
		{ /* '4' */
			return "character"
		}
	case 5:
		{ /* '5' */
			return "8-bit unsigned value"
		}
	case 6:
		{ /* '6' */
			return "8-bit signed value"
		}
	case 7:
		{ /* '7' */
			return "2-byte unsigned value"
		}
	case 8:
		{ /* '8' */
			return "2-byte signed value"
		}
	case 9:
		{ /* '9' */
			return "2-byte float value"
		}
	default:
		{
			return ""
		}
	}
}

func (e KnxDatapointType) SizeInBits() uint8 {
	switch e {
	case 1:
		{ /* '1' */
			return 1
		}
	case 10:
		{ /* '10' */
			return 24
		}
	case 11:
		{ /* '11' */
			return 24
		}
	case 12:
		{ /* '12' */
			return 32
		}
	case 13:
		{ /* '13' */
			return 32
		}
	case 14:
		{ /* '14' */
			return 32
		}
	case 15:
		{ /* '15' */
			return 32
		}
	case 16:
		{ /* '16' */
			return 112
		}
	case 17:
		{ /* '17' */
			return 8
		}
	case 18:
		{ /* '18' */
			return 8
		}
	case 19:
		{ /* '19' */
			return 64
		}
	case 2:
		{ /* '2' */
			return 2
		}
	case 20:
		{ /* '20' */
			return 8
		}
	case 206:
		{ /* '206' */
			return 24
		}
	case 21:
		{ /* '21' */
			return 8
		}
	case 217:
		{ /* '217' */
			return 16
		}
	case 219:
		{ /* '219' */
			return 48
		}
	case 22:
		{ /* '22' */
			return 16
		}
	case 222:
		{ /* '222' */
			return 48
		}
	case 225:
		{ /* '225' */
			return 24
		}
	case 229:
		{ /* '229' */
			return 48
		}
	case 23:
		{ /* '23' */
			return 2
		}
	case 230:
		{ /* '230' */
			return 64
		}
	case 232:
		{ /* '232' */
			return 24
		}
	case 234:
		{ /* '234' */
			return 16
		}
	case 235:
		{ /* '235' */
			return 48
		}
	case 236:
		{ /* '236' */
			return 8
		}
	case 237:
		{ /* '237' */
			return 16
		}
	case 238:
		{ /* '238' */
			return 8
		}
	case 240:
		{ /* '240' */
			return 24
		}
	case 241:
		{ /* '241' */
			return 32
		}
	case 242:
		{ /* '242' */
			return 48
		}
	case 244:
		{ /* '244' */
			return 16
		}
	case 245:
		{ /* '245' */
			return 48
		}
	case 246:
		{ /* '246' */
			return 16
		}
	case 249:
		{ /* '249' */
			return 48
		}
	case 25:
		{ /* '25' */
			return 8
		}
	case 250:
		{ /* '250' */
			return 24
		}
	case 251:
		{ /* '251' */
			return 48
		}
	case 252:
		{ /* '252' */
			return 40
		}
	case 254:
		{ /* '254' */
			return 24
		}
	case 255:
		{ /* '255' */
			return 64
		}
	case 26:
		{ /* '26' */
			return 8
		}
	case 27:
		{ /* '27' */
			return 32
		}
	case 275:
		{ /* '275' */
			return 64
		}
	case 29:
		{ /* '29' */
			return 64
		}
	case 3:
		{ /* '3' */
			return 4
		}
	case 30:
		{ /* '30' */
			return 24
		}
	case 4:
		{ /* '4' */
			return 8
		}
	case 5:
		{ /* '5' */
			return 8
		}
	case 6:
		{ /* '6' */
			return 8
		}
	case 7:
		{ /* '7' */
			return 16
		}
	case 8:
		{ /* '8' */
			return 16
		}
	case 9:
		{ /* '9' */
			return 16
		}
	default:
		{
			return 0
		}
	}
}
func KnxDatapointTypeValueOf(value uint16) KnxDatapointType {
	switch value {
	case 1:
		return KnxDatapointType_DPT_1_BIT
	case 10:
		return KnxDatapointType_DPT_TIME
	case 11:
		return KnxDatapointType_DPT_DATE
	case 12:
		return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
	case 13:
		return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
	case 14:
		return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
	case 15:
		return KnxDatapointType_DPT_ENTRANCE_ACCESS
	case 16:
		return KnxDatapointType_DPT_CHARACTER_STRING
	case 17:
		return KnxDatapointType_DPT_SCENE_NUMBER
	case 18:
		return KnxDatapointType_DPT_SCENE_CONTROL
	case 19:
		return KnxDatapointType_DPT_DATE_TIME
	case 2:
		return KnxDatapointType_DPT_1_BIT_CONTROLLED
	case 20:
		return KnxDatapointType_DPT_1_BYTE
	case 206:
		return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
	case 21:
		return KnxDatapointType_DPT_8_BIT_SET
	case 217:
		return KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION
	case 219:
		return KnxDatapointType_DPT_ALARM_INFO
	case 22:
		return KnxDatapointType_DPT_16_BIT_SET
	case 222:
		return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
	case 225:
		return KnxDatapointType_DPT_SCALING_SPEED
	case 229:
		return KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION
	case 23:
		return KnxDatapointType_DPT_2_BIT_SET
	case 230:
		return KnxDatapointType_DPT_MBUS_ADDRESS
	case 232:
		return KnxDatapointType_DPT_3_BYTE_COLOUR_RGB
	case 234:
		return KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1
	case 235:
		return KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY
	case 236:
		return KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL
	case 237:
		return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT
	case 238:
		return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT
	case 240:
		return KnxDatapointType_DPT_POSITIONS
	case 241:
		return KnxDatapointType_DPT_STATUS_32_BIT
	case 242:
		return KnxDatapointType_DPT_STATUS_48_BIT
	case 244:
		return KnxDatapointType_DPT_CONVERTER_STATUS
	case 245:
		return KnxDatapointType_DPT_CONVERTER_TEST_RESULT
	case 246:
		return KnxDatapointType_DPT_BATTERY_INFORMATION
	case 249:
		return KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION
	case 25:
		return KnxDatapointType_DPT_2_NIBBLE_SET
	case 250:
		return KnxDatapointType_DPT_STATUS_24_BIT
	case 251:
		return KnxDatapointType_DPT_COLOUR_RGBW
	case 252:
		return KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW
	case 254:
		return KnxDatapointType_DPT_RELATIVE_CONTROL_RGB
	case 255:
		return KnxDatapointType_DPT_F32F32
	case 26:
		return KnxDatapointType_DPT_8_BIT_SET_2
	case 27:
		return KnxDatapointType_DPT_32_BIT_SET
	case 275:
		return KnxDatapointType_DPT_F16F16F16F16
	case 29:
		return KnxDatapointType_DPT_ELECTRICAL_ENERGY
	case 3:
		return KnxDatapointType_DPT_3_BIT_CONTROLLED
	case 30:
		return KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION
	case 4:
		return KnxDatapointType_DPT_CHARACTER
	case 5:
		return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
	case 6:
		return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
	case 7:
		return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
	case 8:
		return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
	case 9:
		return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
	}
	return 0
}

func CastKnxDatapointType(structType interface{}) KnxDatapointType {
	castFunc := func(typ interface{}) KnxDatapointType {
		if sKnxDatapointType, ok := typ.(KnxDatapointType); ok {
			return sKnxDatapointType
		}
		return 0
	}
	return castFunc(structType)
}

func (m KnxDatapointType) LengthInBits() uint16 {
	return 16
}

func (m KnxDatapointType) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxDatapointTypeParse(io *utils.ReadBuffer) (KnxDatapointType, error) {
	val, err := io.ReadUint16(16)
	if err != nil {
		return 0, nil
	}
	return KnxDatapointTypeValueOf(val), nil
}

func (e KnxDatapointType) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint16(16, uint16(e))
	return err
}

func (e KnxDatapointType) String() string {
	switch e {
	case KnxDatapointType_DPT_1_BIT:
		return "DPT_1_BIT"
	case KnxDatapointType_DPT_TIME:
		return "DPT_TIME"
	case KnxDatapointType_DPT_DATE:
		return "DPT_DATE"
	case KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE:
		return "DPT_4_BYTE_UNSIGNED_VALUE"
	case KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE:
		return "DPT_4_BYTE_SIGNED_VALUE"
	case KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE:
		return "DPT_4_BYTE_FLOAT_VALUE"
	case KnxDatapointType_DPT_ENTRANCE_ACCESS:
		return "DPT_ENTRANCE_ACCESS"
	case KnxDatapointType_DPT_CHARACTER_STRING:
		return "DPT_CHARACTER_STRING"
	case KnxDatapointType_DPT_SCENE_NUMBER:
		return "DPT_SCENE_NUMBER"
	case KnxDatapointType_DPT_SCENE_CONTROL:
		return "DPT_SCENE_CONTROL"
	case KnxDatapointType_DPT_DATE_TIME:
		return "DPT_DATE_TIME"
	case KnxDatapointType_DPT_1_BIT_CONTROLLED:
		return "DPT_1_BIT_CONTROLLED"
	case KnxDatapointType_DPT_1_BYTE:
		return "DPT_1_BYTE"
	case KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM:
		return "DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM"
	case KnxDatapointType_DPT_8_BIT_SET:
		return "DPT_8_BIT_SET"
	case KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION:
		return "DPT_DATAPOINT_TYPE_VERSION"
	case KnxDatapointType_DPT_ALARM_INFO:
		return "DPT_ALARM_INFO"
	case KnxDatapointType_DPT_16_BIT_SET:
		return "DPT_16_BIT_SET"
	case KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE:
		return "DPT_3X_2_BYTE_FLOAT_VALUE"
	case KnxDatapointType_DPT_SCALING_SPEED:
		return "DPT_SCALING_SPEED"
	case KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION:
		return "DPT_4_1_1_BYTE_COMBINED_INFORMATION"
	case KnxDatapointType_DPT_2_BIT_SET:
		return "DPT_2_BIT_SET"
	case KnxDatapointType_DPT_MBUS_ADDRESS:
		return "DPT_MBUS_ADDRESS"
	case KnxDatapointType_DPT_3_BYTE_COLOUR_RGB:
		return "DPT_3_BYTE_COLOUR_RGB"
	case KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1:
		return "DPT_LANGUAGE_CODE_ISO_639_1"
	case KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY:
		return "DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY"
	case KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL:
		return "DPT_PRIORITISED_MODE_CONTROL"
	case KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT:
		return "DPT_CONFIGURATION_DIAGNOSTICS_16_BIT"
	case KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT:
		return "DPT_CONFIGURATION_DIAGNOSTICS_8_BIT"
	case KnxDatapointType_DPT_POSITIONS:
		return "DPT_POSITIONS"
	case KnxDatapointType_DPT_STATUS_32_BIT:
		return "DPT_STATUS_32_BIT"
	case KnxDatapointType_DPT_STATUS_48_BIT:
		return "DPT_STATUS_48_BIT"
	case KnxDatapointType_DPT_CONVERTER_STATUS:
		return "DPT_CONVERTER_STATUS"
	case KnxDatapointType_DPT_CONVERTER_TEST_RESULT:
		return "DPT_CONVERTER_TEST_RESULT"
	case KnxDatapointType_DPT_BATTERY_INFORMATION:
		return "DPT_BATTERY_INFORMATION"
	case KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION:
		return "DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION"
	case KnxDatapointType_DPT_2_NIBBLE_SET:
		return "DPT_2_NIBBLE_SET"
	case KnxDatapointType_DPT_STATUS_24_BIT:
		return "DPT_STATUS_24_BIT"
	case KnxDatapointType_DPT_COLOUR_RGBW:
		return "DPT_COLOUR_RGBW"
	case KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW:
		return "DPT_RELATIVE_CONTROL_RGBW"
	case KnxDatapointType_DPT_RELATIVE_CONTROL_RGB:
		return "DPT_RELATIVE_CONTROL_RGB"
	case KnxDatapointType_DPT_F32F32:
		return "DPT_F32F32"
	case KnxDatapointType_DPT_8_BIT_SET_2:
		return "DPT_8_BIT_SET_2"
	case KnxDatapointType_DPT_32_BIT_SET:
		return "DPT_32_BIT_SET"
	case KnxDatapointType_DPT_F16F16F16F16:
		return "DPT_F16F16F16F16"
	case KnxDatapointType_DPT_ELECTRICAL_ENERGY:
		return "DPT_ELECTRICAL_ENERGY"
	case KnxDatapointType_DPT_3_BIT_CONTROLLED:
		return "DPT_3_BIT_CONTROLLED"
	case KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION:
		return "DPT_24_TIMES_CHANNEL_ACTIVATION"
	case KnxDatapointType_DPT_CHARACTER:
		return "DPT_CHARACTER"
	case KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE:
		return "DPT_8_BIT_UNSIGNED_VALUE"
	case KnxDatapointType_DPT_8_BIT_SIGNED_VALUE:
		return "DPT_8_BIT_SIGNED_VALUE"
	case KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE:
		return "DPT_2_BYTE_UNSIGNED_VALUE"
	case KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE:
		return "DPT_2_BYTE_SIGNED_VALUE"
	case KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE:
		return "DPT_2_BYTE_FLOAT_VALUE"
	}
	return ""
}
