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
    Number() uint16
    Name() string
    SizeInBits() uint8
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxDatapointType_DPT_UNKNOWN KnxDatapointType = 0
    KnxDatapointType_DPT_1_BIT KnxDatapointType = 1
    KnxDatapointType_DPT_1_BIT_CONTROLLED KnxDatapointType = 2
    KnxDatapointType_DPT_3_BIT_CONTROLLED KnxDatapointType = 3
    KnxDatapointType_DPT_CHARACTER KnxDatapointType = 4
    KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE KnxDatapointType = 5
    KnxDatapointType_DPT_8_BIT_SIGNED_VALUE KnxDatapointType = 6
    KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE KnxDatapointType = 7
    KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE KnxDatapointType = 8
    KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE KnxDatapointType = 9
    KnxDatapointType_DPT_TIME KnxDatapointType = 10
    KnxDatapointType_DPT_DATE KnxDatapointType = 11
    KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE KnxDatapointType = 12
    KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE KnxDatapointType = 13
    KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE KnxDatapointType = 14
    KnxDatapointType_DPT_ENTRANCE_ACCESS KnxDatapointType = 15
    KnxDatapointType_DPT_CHARACTER_STRING KnxDatapointType = 16
    KnxDatapointType_DPT_SCENE_NUMBER KnxDatapointType = 17
    KnxDatapointType_DPT_SCENE_CONTROL KnxDatapointType = 18
    KnxDatapointType_DPT_DATE_TIME KnxDatapointType = 19
    KnxDatapointType_DPT_1_BYTE KnxDatapointType = 20
    KnxDatapointType_DPT_8_BIT_SET KnxDatapointType = 21
    KnxDatapointType_DPT_16_BIT_SET KnxDatapointType = 22
    KnxDatapointType_DPT_2_BIT_SET KnxDatapointType = 23
    KnxDatapointType_DPT_2_NIBBLE_SET KnxDatapointType = 24
    KnxDatapointType_DPT_8_BIT_SET_2 KnxDatapointType = 25
    KnxDatapointType_DPT_32_BIT_SET KnxDatapointType = 26
    KnxDatapointType_DPT_ELECTRICAL_ENERGY KnxDatapointType = 27
    KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION KnxDatapointType = 28
    KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM KnxDatapointType = 29
    KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION KnxDatapointType = 30
    KnxDatapointType_DPT_ALARM_INFO KnxDatapointType = 31
    KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE KnxDatapointType = 32
    KnxDatapointType_DPT_SCALING_SPEED KnxDatapointType = 33
    KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION KnxDatapointType = 34
    KnxDatapointType_DPT_MBUS_ADDRESS KnxDatapointType = 35
    KnxDatapointType_DPT_3_BYTE_COLOUR_RGB KnxDatapointType = 36
    KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1 KnxDatapointType = 37
    KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY KnxDatapointType = 38
    KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL KnxDatapointType = 39
    KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT KnxDatapointType = 40
    KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT KnxDatapointType = 41
    KnxDatapointType_DPT_POSITIONS KnxDatapointType = 42
    KnxDatapointType_DPT_STATUS_32_BIT KnxDatapointType = 43
    KnxDatapointType_DPT_STATUS_48_BIT KnxDatapointType = 44
    KnxDatapointType_DPT_CONVERTER_STATUS KnxDatapointType = 45
    KnxDatapointType_DPT_CONVERTER_TEST_RESULT KnxDatapointType = 46
    KnxDatapointType_DPT_BATTERY_INFORMATION KnxDatapointType = 47
    KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION KnxDatapointType = 48
    KnxDatapointType_DPT_STATUS_24_BIT KnxDatapointType = 49
    KnxDatapointType_DPT_COLOUR_RGBW KnxDatapointType = 50
    KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW KnxDatapointType = 51
    KnxDatapointType_DPT_RELATIVE_CONTROL_RGB KnxDatapointType = 52
    KnxDatapointType_DPT_F32F32 KnxDatapointType = 53
    KnxDatapointType_DPT_F16F16F16F16 KnxDatapointType = 54
)


func (e KnxDatapointType) Number() uint16 {
    switch e  {
        case 0: { /* '0' */
            return 0
        }
        case 1: { /* '1' */
            return 1
        }
        case 10: { /* '10' */
            return 10
        }
        case 11: { /* '11' */
            return 11
        }
        case 12: { /* '12' */
            return 12
        }
        case 13: { /* '13' */
            return 13
        }
        case 14: { /* '14' */
            return 14
        }
        case 15: { /* '15' */
            return 15
        }
        case 16: { /* '16' */
            return 16
        }
        case 17: { /* '17' */
            return 17
        }
        case 18: { /* '18' */
            return 18
        }
        case 19: { /* '19' */
            return 19
        }
        case 2: { /* '2' */
            return 2
        }
        case 20: { /* '20' */
            return 20
        }
        case 21: { /* '21' */
            return 21
        }
        case 22: { /* '22' */
            return 22
        }
        case 23: { /* '23' */
            return 23
        }
        case 24: { /* '24' */
            return 25
        }
        case 25: { /* '25' */
            return 26
        }
        case 26: { /* '26' */
            return 27
        }
        case 27: { /* '27' */
            return 29
        }
        case 28: { /* '28' */
            return 30
        }
        case 29: { /* '29' */
            return 206
        }
        case 3: { /* '3' */
            return 3
        }
        case 30: { /* '30' */
            return 217
        }
        case 31: { /* '31' */
            return 219
        }
        case 32: { /* '32' */
            return 222
        }
        case 33: { /* '33' */
            return 225
        }
        case 34: { /* '34' */
            return 229
        }
        case 35: { /* '35' */
            return 230
        }
        case 36: { /* '36' */
            return 232
        }
        case 37: { /* '37' */
            return 234
        }
        case 38: { /* '38' */
            return 235
        }
        case 39: { /* '39' */
            return 236
        }
        case 4: { /* '4' */
            return 4
        }
        case 40: { /* '40' */
            return 237
        }
        case 41: { /* '41' */
            return 238
        }
        case 42: { /* '42' */
            return 240
        }
        case 43: { /* '43' */
            return 241
        }
        case 44: { /* '44' */
            return 242
        }
        case 45: { /* '45' */
            return 244
        }
        case 46: { /* '46' */
            return 245
        }
        case 47: { /* '47' */
            return 246
        }
        case 48: { /* '48' */
            return 249
        }
        case 49: { /* '49' */
            return 250
        }
        case 5: { /* '5' */
            return 5
        }
        case 50: { /* '50' */
            return 251
        }
        case 51: { /* '51' */
            return 252
        }
        case 52: { /* '52' */
            return 254
        }
        case 53: { /* '53' */
            return 255
        }
        case 54: { /* '54' */
            return 275
        }
        case 6: { /* '6' */
            return 6
        }
        case 7: { /* '7' */
            return 7
        }
        case 8: { /* '8' */
            return 8
        }
        case 9: { /* '9' */
            return 9
        }
        default: {
            return 0
        }
    }
}

func (e KnxDatapointType) Name() string {
    switch e  {
        case 0: { /* '0' */
            return "Unknown Datapoint Type"
        }
        case 1: { /* '1' */
            return "1-bit"
        }
        case 10: { /* '10' */
            return "time"
        }
        case 11: { /* '11' */
            return "date"
        }
        case 12: { /* '12' */
            return "4-byte unsigned value"
        }
        case 13: { /* '13' */
            return "4-byte signed value"
        }
        case 14: { /* '14' */
            return "4-byte float value"
        }
        case 15: { /* '15' */
            return "entrance access"
        }
        case 16: { /* '16' */
            return "character string"
        }
        case 17: { /* '17' */
            return "scene number"
        }
        case 18: { /* '18' */
            return "scene control"
        }
        case 19: { /* '19' */
            return "Date Time"
        }
        case 2: { /* '2' */
            return "1-bit controlled"
        }
        case 20: { /* '20' */
            return "1-byte"
        }
        case 21: { /* '21' */
            return "8-bit set"
        }
        case 22: { /* '22' */
            return "16-bit set"
        }
        case 23: { /* '23' */
            return "2-bit set"
        }
        case 24: { /* '24' */
            return "2-nibble set"
        }
        case 25: { /* '25' */
            return "8-bit set"
        }
        case 26: { /* '26' */
            return "32-bit set"
        }
        case 27: { /* '27' */
            return "electrical energy"
        }
        case 28: { /* '28' */
            return "24 times channel activation"
        }
        case 29: { /* '29' */
            return "16-bit unsigned value & 8-bit enum"
        }
        case 3: { /* '3' */
            return "3-bit controlled"
        }
        case 30: { /* '30' */
            return "datapoint type version"
        }
        case 31: { /* '31' */
            return "alarm info"
        }
        case 32: { /* '32' */
            return "3x 2-byte float value"
        }
        case 33: { /* '33' */
            return "scaling speed"
        }
        case 34: { /* '34' */
            return "4-1-1 byte combined information"
        }
        case 35: { /* '35' */
            return "MBus address"
        }
        case 36: { /* '36' */
            return "3-byte colour RGB"
        }
        case 37: { /* '37' */
            return "language code ISO 639-1"
        }
        case 38: { /* '38' */
            return "Signed value with classification and validity"
        }
        case 39: { /* '39' */
            return "Prioritised Mode Control"
        }
        case 4: { /* '4' */
            return "character"
        }
        case 40: { /* '40' */
            return "configuration/ diagnostics"
        }
        case 41: { /* '41' */
            return "configuration/ diagnostics"
        }
        case 42: { /* '42' */
            return "positions"
        }
        case 43: { /* '43' */
            return "status"
        }
        case 44: { /* '44' */
            return "status"
        }
        case 45: { /* '45' */
            return "Converter Status"
        }
        case 46: { /* '46' */
            return "Converter test result"
        }
        case 47: { /* '47' */
            return "Battery Information"
        }
        case 48: { /* '48' */
            return "brightness colour temperature transition"
        }
        case 49: { /* '49' */
            return "status"
        }
        case 5: { /* '5' */
            return "8-bit unsigned value"
        }
        case 50: { /* '50' */
            return "Colour RGBW"
        }
        case 51: { /* '51' */
            return "Relative Control RGBW"
        }
        case 52: { /* '52' */
            return "Relative Control RGB"
        }
        case 53: { /* '53' */
            return "F32F32"
        }
        case 54: { /* '54' */
            return "F16F16F16F16"
        }
        case 6: { /* '6' */
            return "8-bit signed value"
        }
        case 7: { /* '7' */
            return "2-byte unsigned value"
        }
        case 8: { /* '8' */
            return "2-byte signed value"
        }
        case 9: { /* '9' */
            return "2-byte float value"
        }
        default: {
            return ""
        }
    }
}

func (e KnxDatapointType) SizeInBits() uint8 {
    switch e  {
        case 0: { /* '0' */
            return 0
        }
        case 1: { /* '1' */
            return 1
        }
        case 10: { /* '10' */
            return 24
        }
        case 11: { /* '11' */
            return 24
        }
        case 12: { /* '12' */
            return 32
        }
        case 13: { /* '13' */
            return 32
        }
        case 14: { /* '14' */
            return 32
        }
        case 15: { /* '15' */
            return 32
        }
        case 16: { /* '16' */
            return 112
        }
        case 17: { /* '17' */
            return 8
        }
        case 18: { /* '18' */
            return 8
        }
        case 19: { /* '19' */
            return 64
        }
        case 2: { /* '2' */
            return 2
        }
        case 20: { /* '20' */
            return 8
        }
        case 21: { /* '21' */
            return 8
        }
        case 22: { /* '22' */
            return 16
        }
        case 23: { /* '23' */
            return 2
        }
        case 24: { /* '24' */
            return 8
        }
        case 25: { /* '25' */
            return 8
        }
        case 26: { /* '26' */
            return 32
        }
        case 27: { /* '27' */
            return 64
        }
        case 28: { /* '28' */
            return 24
        }
        case 29: { /* '29' */
            return 24
        }
        case 3: { /* '3' */
            return 4
        }
        case 30: { /* '30' */
            return 16
        }
        case 31: { /* '31' */
            return 48
        }
        case 32: { /* '32' */
            return 48
        }
        case 33: { /* '33' */
            return 24
        }
        case 34: { /* '34' */
            return 48
        }
        case 35: { /* '35' */
            return 64
        }
        case 36: { /* '36' */
            return 24
        }
        case 37: { /* '37' */
            return 16
        }
        case 38: { /* '38' */
            return 48
        }
        case 39: { /* '39' */
            return 8
        }
        case 4: { /* '4' */
            return 8
        }
        case 40: { /* '40' */
            return 16
        }
        case 41: { /* '41' */
            return 8
        }
        case 42: { /* '42' */
            return 24
        }
        case 43: { /* '43' */
            return 32
        }
        case 44: { /* '44' */
            return 48
        }
        case 45: { /* '45' */
            return 16
        }
        case 46: { /* '46' */
            return 48
        }
        case 47: { /* '47' */
            return 16
        }
        case 48: { /* '48' */
            return 48
        }
        case 49: { /* '49' */
            return 24
        }
        case 5: { /* '5' */
            return 8
        }
        case 50: { /* '50' */
            return 48
        }
        case 51: { /* '51' */
            return 40
        }
        case 52: { /* '52' */
            return 24
        }
        case 53: { /* '53' */
            return 64
        }
        case 54: { /* '54' */
            return 64
        }
        case 6: { /* '6' */
            return 8
        }
        case 7: { /* '7' */
            return 16
        }
        case 8: { /* '8' */
            return 16
        }
        case 9: { /* '9' */
            return 16
        }
        default: {
            return 0
        }
    }
}
func KnxDatapointTypeByValue(value uint16) KnxDatapointType {
    switch value {
        case 0:
            return KnxDatapointType_DPT_UNKNOWN
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
        case 21:
            return KnxDatapointType_DPT_8_BIT_SET
        case 22:
            return KnxDatapointType_DPT_16_BIT_SET
        case 23:
            return KnxDatapointType_DPT_2_BIT_SET
        case 24:
            return KnxDatapointType_DPT_2_NIBBLE_SET
        case 25:
            return KnxDatapointType_DPT_8_BIT_SET_2
        case 26:
            return KnxDatapointType_DPT_32_BIT_SET
        case 27:
            return KnxDatapointType_DPT_ELECTRICAL_ENERGY
        case 28:
            return KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION
        case 29:
            return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
        case 3:
            return KnxDatapointType_DPT_3_BIT_CONTROLLED
        case 30:
            return KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION
        case 31:
            return KnxDatapointType_DPT_ALARM_INFO
        case 32:
            return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
        case 33:
            return KnxDatapointType_DPT_SCALING_SPEED
        case 34:
            return KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION
        case 35:
            return KnxDatapointType_DPT_MBUS_ADDRESS
        case 36:
            return KnxDatapointType_DPT_3_BYTE_COLOUR_RGB
        case 37:
            return KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1
        case 38:
            return KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY
        case 39:
            return KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL
        case 4:
            return KnxDatapointType_DPT_CHARACTER
        case 40:
            return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT
        case 41:
            return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT
        case 42:
            return KnxDatapointType_DPT_POSITIONS
        case 43:
            return KnxDatapointType_DPT_STATUS_32_BIT
        case 44:
            return KnxDatapointType_DPT_STATUS_48_BIT
        case 45:
            return KnxDatapointType_DPT_CONVERTER_STATUS
        case 46:
            return KnxDatapointType_DPT_CONVERTER_TEST_RESULT
        case 47:
            return KnxDatapointType_DPT_BATTERY_INFORMATION
        case 48:
            return KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION
        case 49:
            return KnxDatapointType_DPT_STATUS_24_BIT
        case 5:
            return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
        case 50:
            return KnxDatapointType_DPT_COLOUR_RGBW
        case 51:
            return KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW
        case 52:
            return KnxDatapointType_DPT_RELATIVE_CONTROL_RGB
        case 53:
            return KnxDatapointType_DPT_F32F32
        case 54:
            return KnxDatapointType_DPT_F16F16F16F16
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

func KnxDatapointTypeByName(value string) KnxDatapointType {
    switch value {
    case "DPT_UNKNOWN":
        return KnxDatapointType_DPT_UNKNOWN
    case "DPT_1_BIT":
        return KnxDatapointType_DPT_1_BIT
    case "DPT_TIME":
        return KnxDatapointType_DPT_TIME
    case "DPT_DATE":
        return KnxDatapointType_DPT_DATE
    case "DPT_4_BYTE_UNSIGNED_VALUE":
        return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
    case "DPT_4_BYTE_SIGNED_VALUE":
        return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
    case "DPT_4_BYTE_FLOAT_VALUE":
        return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
    case "DPT_ENTRANCE_ACCESS":
        return KnxDatapointType_DPT_ENTRANCE_ACCESS
    case "DPT_CHARACTER_STRING":
        return KnxDatapointType_DPT_CHARACTER_STRING
    case "DPT_SCENE_NUMBER":
        return KnxDatapointType_DPT_SCENE_NUMBER
    case "DPT_SCENE_CONTROL":
        return KnxDatapointType_DPT_SCENE_CONTROL
    case "DPT_DATE_TIME":
        return KnxDatapointType_DPT_DATE_TIME
    case "DPT_1_BIT_CONTROLLED":
        return KnxDatapointType_DPT_1_BIT_CONTROLLED
    case "DPT_1_BYTE":
        return KnxDatapointType_DPT_1_BYTE
    case "DPT_8_BIT_SET":
        return KnxDatapointType_DPT_8_BIT_SET
    case "DPT_16_BIT_SET":
        return KnxDatapointType_DPT_16_BIT_SET
    case "DPT_2_BIT_SET":
        return KnxDatapointType_DPT_2_BIT_SET
    case "DPT_2_NIBBLE_SET":
        return KnxDatapointType_DPT_2_NIBBLE_SET
    case "DPT_8_BIT_SET_2":
        return KnxDatapointType_DPT_8_BIT_SET_2
    case "DPT_32_BIT_SET":
        return KnxDatapointType_DPT_32_BIT_SET
    case "DPT_ELECTRICAL_ENERGY":
        return KnxDatapointType_DPT_ELECTRICAL_ENERGY
    case "DPT_24_TIMES_CHANNEL_ACTIVATION":
        return KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION
    case "DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM":
        return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
    case "DPT_3_BIT_CONTROLLED":
        return KnxDatapointType_DPT_3_BIT_CONTROLLED
    case "DPT_DATAPOINT_TYPE_VERSION":
        return KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION
    case "DPT_ALARM_INFO":
        return KnxDatapointType_DPT_ALARM_INFO
    case "DPT_3X_2_BYTE_FLOAT_VALUE":
        return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
    case "DPT_SCALING_SPEED":
        return KnxDatapointType_DPT_SCALING_SPEED
    case "DPT_4_1_1_BYTE_COMBINED_INFORMATION":
        return KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION
    case "DPT_MBUS_ADDRESS":
        return KnxDatapointType_DPT_MBUS_ADDRESS
    case "DPT_3_BYTE_COLOUR_RGB":
        return KnxDatapointType_DPT_3_BYTE_COLOUR_RGB
    case "DPT_LANGUAGE_CODE_ISO_639_1":
        return KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1
    case "DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY":
        return KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY
    case "DPT_PRIORITISED_MODE_CONTROL":
        return KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL
    case "DPT_CHARACTER":
        return KnxDatapointType_DPT_CHARACTER
    case "DPT_CONFIGURATION_DIAGNOSTICS_16_BIT":
        return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT
    case "DPT_CONFIGURATION_DIAGNOSTICS_8_BIT":
        return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT
    case "DPT_POSITIONS":
        return KnxDatapointType_DPT_POSITIONS
    case "DPT_STATUS_32_BIT":
        return KnxDatapointType_DPT_STATUS_32_BIT
    case "DPT_STATUS_48_BIT":
        return KnxDatapointType_DPT_STATUS_48_BIT
    case "DPT_CONVERTER_STATUS":
        return KnxDatapointType_DPT_CONVERTER_STATUS
    case "DPT_CONVERTER_TEST_RESULT":
        return KnxDatapointType_DPT_CONVERTER_TEST_RESULT
    case "DPT_BATTERY_INFORMATION":
        return KnxDatapointType_DPT_BATTERY_INFORMATION
    case "DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION":
        return KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION
    case "DPT_STATUS_24_BIT":
        return KnxDatapointType_DPT_STATUS_24_BIT
    case "DPT_8_BIT_UNSIGNED_VALUE":
        return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
    case "DPT_COLOUR_RGBW":
        return KnxDatapointType_DPT_COLOUR_RGBW
    case "DPT_RELATIVE_CONTROL_RGBW":
        return KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW
    case "DPT_RELATIVE_CONTROL_RGB":
        return KnxDatapointType_DPT_RELATIVE_CONTROL_RGB
    case "DPT_F32F32":
        return KnxDatapointType_DPT_F32F32
    case "DPT_F16F16F16F16":
        return KnxDatapointType_DPT_F16F16F16F16
    case "DPT_8_BIT_SIGNED_VALUE":
        return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
    case "DPT_2_BYTE_UNSIGNED_VALUE":
        return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
    case "DPT_2_BYTE_SIGNED_VALUE":
        return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
    case "DPT_2_BYTE_FLOAT_VALUE":
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
    return KnxDatapointTypeByValue(val), nil
}

func (e KnxDatapointType) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint16(16, uint16(e))
    return err
}

func (e KnxDatapointType) String() string {
    switch e {
    case KnxDatapointType_DPT_UNKNOWN:
        return "DPT_UNKNOWN"
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
    case KnxDatapointType_DPT_8_BIT_SET:
        return "DPT_8_BIT_SET"
    case KnxDatapointType_DPT_16_BIT_SET:
        return "DPT_16_BIT_SET"
    case KnxDatapointType_DPT_2_BIT_SET:
        return "DPT_2_BIT_SET"
    case KnxDatapointType_DPT_2_NIBBLE_SET:
        return "DPT_2_NIBBLE_SET"
    case KnxDatapointType_DPT_8_BIT_SET_2:
        return "DPT_8_BIT_SET_2"
    case KnxDatapointType_DPT_32_BIT_SET:
        return "DPT_32_BIT_SET"
    case KnxDatapointType_DPT_ELECTRICAL_ENERGY:
        return "DPT_ELECTRICAL_ENERGY"
    case KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION:
        return "DPT_24_TIMES_CHANNEL_ACTIVATION"
    case KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM:
        return "DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM"
    case KnxDatapointType_DPT_3_BIT_CONTROLLED:
        return "DPT_3_BIT_CONTROLLED"
    case KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION:
        return "DPT_DATAPOINT_TYPE_VERSION"
    case KnxDatapointType_DPT_ALARM_INFO:
        return "DPT_ALARM_INFO"
    case KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE:
        return "DPT_3X_2_BYTE_FLOAT_VALUE"
    case KnxDatapointType_DPT_SCALING_SPEED:
        return "DPT_SCALING_SPEED"
    case KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION:
        return "DPT_4_1_1_BYTE_COMBINED_INFORMATION"
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
    case KnxDatapointType_DPT_CHARACTER:
        return "DPT_CHARACTER"
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
    case KnxDatapointType_DPT_STATUS_24_BIT:
        return "DPT_STATUS_24_BIT"
    case KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE:
        return "DPT_8_BIT_UNSIGNED_VALUE"
    case KnxDatapointType_DPT_COLOUR_RGBW:
        return "DPT_COLOUR_RGBW"
    case KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW:
        return "DPT_RELATIVE_CONTROL_RGBW"
    case KnxDatapointType_DPT_RELATIVE_CONTROL_RGB:
        return "DPT_RELATIVE_CONTROL_RGB"
    case KnxDatapointType_DPT_F32F32:
        return "DPT_F32F32"
    case KnxDatapointType_DPT_F16F16F16F16:
        return "DPT_F16F16F16F16"
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
