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

type KnxPropertyDataType uint8

type IKnxPropertyDataType interface {
	SizeInBytes() uint8
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxPropertyDataType_PDT_CONTROL             KnxPropertyDataType = 0
	KnxPropertyDataType_PDT_CHAR                KnxPropertyDataType = 1
	KnxPropertyDataType_PDT_UNSIGNED_CHAR       KnxPropertyDataType = 2
	KnxPropertyDataType_PDT_INT                 KnxPropertyDataType = 3
	KnxPropertyDataType_PDT_UNSIGNED_INT        KnxPropertyDataType = 4
	KnxPropertyDataType_PDT_KNX_FLOAT           KnxPropertyDataType = 5
	KnxPropertyDataType_PDT_DATE                KnxPropertyDataType = 6
	KnxPropertyDataType_PDT_TIME                KnxPropertyDataType = 7
	KnxPropertyDataType_PDT_LONG                KnxPropertyDataType = 8
	KnxPropertyDataType_PDT_UNSIGNED_LONG       KnxPropertyDataType = 9
	KnxPropertyDataType_PDT_FLOAT               KnxPropertyDataType = 10
	KnxPropertyDataType_PDT_DOUBLE              KnxPropertyDataType = 11
	KnxPropertyDataType_PDT_CHAR_BLOCK          KnxPropertyDataType = 12
	KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS KnxPropertyDataType = 13
	KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK    KnxPropertyDataType = 14
	KnxPropertyDataType_PDT_DATE_TIME           KnxPropertyDataType = 15
	KnxPropertyDataType_PDT_VARIABLE_LENGTH     KnxPropertyDataType = 16
	KnxPropertyDataType_PDT_GENERIC_01          KnxPropertyDataType = 17
	KnxPropertyDataType_PDT_GENERIC_02          KnxPropertyDataType = 18
	KnxPropertyDataType_PDT_GENERIC_03          KnxPropertyDataType = 19
	KnxPropertyDataType_PDT_GENERIC_04          KnxPropertyDataType = 20
	KnxPropertyDataType_PDT_GENERIC_05          KnxPropertyDataType = 21
	KnxPropertyDataType_PDT_GENERIC_06          KnxPropertyDataType = 22
	KnxPropertyDataType_PDT_GENERIC_07          KnxPropertyDataType = 23
	KnxPropertyDataType_PDT_GENERIC_08          KnxPropertyDataType = 24
	KnxPropertyDataType_PDT_GENERIC_09          KnxPropertyDataType = 25
	KnxPropertyDataType_PDT_GENERIC_10          KnxPropertyDataType = 26
	KnxPropertyDataType_PDT_GENERIC_11          KnxPropertyDataType = 27
	KnxPropertyDataType_PDT_GENERIC_12          KnxPropertyDataType = 28
	KnxPropertyDataType_PDT_GENERIC_13          KnxPropertyDataType = 29
	KnxPropertyDataType_PDT_GENERIC_14          KnxPropertyDataType = 30
	KnxPropertyDataType_PDT_GENERIC_15          KnxPropertyDataType = 31
	KnxPropertyDataType_PDT_GENERIC_16          KnxPropertyDataType = 32
	KnxPropertyDataType_PDT_GENERIC_17          KnxPropertyDataType = 33
	KnxPropertyDataType_PDT_GENERIC_18          KnxPropertyDataType = 34
	KnxPropertyDataType_PDT_GENERIC_19          KnxPropertyDataType = 35
	KnxPropertyDataType_PDT_GENERIC_20          KnxPropertyDataType = 36
	KnxPropertyDataType_PDT_UTF_8               KnxPropertyDataType = 47
	KnxPropertyDataType_PDT_VERSION             KnxPropertyDataType = 48
	KnxPropertyDataType_PDT_ALARM_INFO          KnxPropertyDataType = 49
	KnxPropertyDataType_PDT_BINARY_INFORMATION  KnxPropertyDataType = 50
	KnxPropertyDataType_PDT_BITSET8             KnxPropertyDataType = 51
	KnxPropertyDataType_PDT_BITSET16            KnxPropertyDataType = 52
	KnxPropertyDataType_PDT_ENUM8               KnxPropertyDataType = 53
	KnxPropertyDataType_PDT_SCALING             KnxPropertyDataType = 54
	KnxPropertyDataType_PDT_NE_VL               KnxPropertyDataType = 60
	KnxPropertyDataType_PDT_NE_FL               KnxPropertyDataType = 61
	KnxPropertyDataType_PDT_FUNCTION            KnxPropertyDataType = 62
	KnxPropertyDataType_PDT_ESCAPE              KnxPropertyDataType = 63
)

func (e KnxPropertyDataType) SizeInBytes() uint8 {
	switch e {
	case 0:
		{ /* '0' */
			return 10
		}
	case 1:
		{ /* '1' */
			return 1
		}
	case 10:
		{ /* '10' */
			return 4
		}
	case 11:
		{ /* '11' */
			return 8
		}
	case 12:
		{ /* '12' */
			return 10
		}
	case 13:
		{ /* '13' */
			return 3
		}
	case 14:
		{ /* '14' */
			return 5
		}
	case 15:
		{ /* '15' */
			return 8
		}
	case 16:
		{ /* '16' */
			return 0
		}
	case 17:
		{ /* '17' */
			return 1
		}
	case 18:
		{ /* '18' */
			return 2
		}
	case 19:
		{ /* '19' */
			return 3
		}
	case 2:
		{ /* '2' */
			return 1
		}
	case 20:
		{ /* '20' */
			return 4
		}
	case 21:
		{ /* '21' */
			return 5
		}
	case 22:
		{ /* '22' */
			return 6
		}
	case 23:
		{ /* '23' */
			return 7
		}
	case 24:
		{ /* '24' */
			return 8
		}
	case 25:
		{ /* '25' */
			return 9
		}
	case 26:
		{ /* '26' */
			return 10
		}
	case 27:
		{ /* '27' */
			return 11
		}
	case 28:
		{ /* '28' */
			return 12
		}
	case 29:
		{ /* '29' */
			return 13
		}
	case 3:
		{ /* '3' */
			return 2
		}
	case 30:
		{ /* '30' */
			return 14
		}
	case 31:
		{ /* '31' */
			return 15
		}
	case 32:
		{ /* '32' */
			return 16
		}
	case 33:
		{ /* '33' */
			return 17
		}
	case 34:
		{ /* '34' */
			return 18
		}
	case 35:
		{ /* '35' */
			return 19
		}
	case 36:
		{ /* '36' */
			return 20
		}
	case 4:
		{ /* '4' */
			return 2
		}
	case 47:
		{ /* '47' */
			return 0
		}
	case 48:
		{ /* '48' */
			return 2
		}
	case 49:
		{ /* '49' */
			return 6
		}
	case 5:
		{ /* '5' */
			return 2
		}
	case 50:
		{ /* '50' */
			return 1
		}
	case 51:
		{ /* '51' */
			return 1
		}
	case 52:
		{ /* '52' */
			return 2
		}
	case 53:
		{ /* '53' */
			return 1
		}
	case 54:
		{ /* '54' */
			return 1
		}
	case 6:
		{ /* '6' */
			return 3
		}
	case 60:
		{ /* '60' */
			return 0
		}
	case 61:
		{ /* '61' */
			return 0
		}
	case 62:
		{ /* '62' */
			return 0
		}
	case 63:
		{ /* '63' */
			return 0
		}
	case 7:
		{ /* '7' */
			return 3
		}
	case 8:
		{ /* '8' */
			return 4
		}
	case 9:
		{ /* '9' */
			return 4
		}
	default:
		{
			return 0
		}
	}
}
func KnxPropertyDataTypeByValue(value uint8) KnxPropertyDataType {
	switch value {
	case 0:
		return KnxPropertyDataType_PDT_CONTROL
	case 1:
		return KnxPropertyDataType_PDT_CHAR
	case 10:
		return KnxPropertyDataType_PDT_FLOAT
	case 11:
		return KnxPropertyDataType_PDT_DOUBLE
	case 12:
		return KnxPropertyDataType_PDT_CHAR_BLOCK
	case 13:
		return KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS
	case 14:
		return KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK
	case 15:
		return KnxPropertyDataType_PDT_DATE_TIME
	case 16:
		return KnxPropertyDataType_PDT_VARIABLE_LENGTH
	case 17:
		return KnxPropertyDataType_PDT_GENERIC_01
	case 18:
		return KnxPropertyDataType_PDT_GENERIC_02
	case 19:
		return KnxPropertyDataType_PDT_GENERIC_03
	case 2:
		return KnxPropertyDataType_PDT_UNSIGNED_CHAR
	case 20:
		return KnxPropertyDataType_PDT_GENERIC_04
	case 21:
		return KnxPropertyDataType_PDT_GENERIC_05
	case 22:
		return KnxPropertyDataType_PDT_GENERIC_06
	case 23:
		return KnxPropertyDataType_PDT_GENERIC_07
	case 24:
		return KnxPropertyDataType_PDT_GENERIC_08
	case 25:
		return KnxPropertyDataType_PDT_GENERIC_09
	case 26:
		return KnxPropertyDataType_PDT_GENERIC_10
	case 27:
		return KnxPropertyDataType_PDT_GENERIC_11
	case 28:
		return KnxPropertyDataType_PDT_GENERIC_12
	case 29:
		return KnxPropertyDataType_PDT_GENERIC_13
	case 3:
		return KnxPropertyDataType_PDT_INT
	case 30:
		return KnxPropertyDataType_PDT_GENERIC_14
	case 31:
		return KnxPropertyDataType_PDT_GENERIC_15
	case 32:
		return KnxPropertyDataType_PDT_GENERIC_16
	case 33:
		return KnxPropertyDataType_PDT_GENERIC_17
	case 34:
		return KnxPropertyDataType_PDT_GENERIC_18
	case 35:
		return KnxPropertyDataType_PDT_GENERIC_19
	case 36:
		return KnxPropertyDataType_PDT_GENERIC_20
	case 4:
		return KnxPropertyDataType_PDT_UNSIGNED_INT
	case 47:
		return KnxPropertyDataType_PDT_UTF_8
	case 48:
		return KnxPropertyDataType_PDT_VERSION
	case 49:
		return KnxPropertyDataType_PDT_ALARM_INFO
	case 5:
		return KnxPropertyDataType_PDT_KNX_FLOAT
	case 50:
		return KnxPropertyDataType_PDT_BINARY_INFORMATION
	case 51:
		return KnxPropertyDataType_PDT_BITSET8
	case 52:
		return KnxPropertyDataType_PDT_BITSET16
	case 53:
		return KnxPropertyDataType_PDT_ENUM8
	case 54:
		return KnxPropertyDataType_PDT_SCALING
	case 6:
		return KnxPropertyDataType_PDT_DATE
	case 60:
		return KnxPropertyDataType_PDT_NE_VL
	case 61:
		return KnxPropertyDataType_PDT_NE_FL
	case 62:
		return KnxPropertyDataType_PDT_FUNCTION
	case 63:
		return KnxPropertyDataType_PDT_ESCAPE
	case 7:
		return KnxPropertyDataType_PDT_TIME
	case 8:
		return KnxPropertyDataType_PDT_LONG
	case 9:
		return KnxPropertyDataType_PDT_UNSIGNED_LONG
	}
	return 0
}

func KnxPropertyDataTypeByName(value string) KnxPropertyDataType {
	switch value {
	case "PDT_CONTROL":
		return KnxPropertyDataType_PDT_CONTROL
	case "PDT_CHAR":
		return KnxPropertyDataType_PDT_CHAR
	case "PDT_FLOAT":
		return KnxPropertyDataType_PDT_FLOAT
	case "PDT_DOUBLE":
		return KnxPropertyDataType_PDT_DOUBLE
	case "PDT_CHAR_BLOCK":
		return KnxPropertyDataType_PDT_CHAR_BLOCK
	case "PDT_POLL_GROUP_SETTINGS":
		return KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS
	case "PDT_SHORT_CHAR_BLOCK":
		return KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK
	case "PDT_DATE_TIME":
		return KnxPropertyDataType_PDT_DATE_TIME
	case "PDT_VARIABLE_LENGTH":
		return KnxPropertyDataType_PDT_VARIABLE_LENGTH
	case "PDT_GENERIC_01":
		return KnxPropertyDataType_PDT_GENERIC_01
	case "PDT_GENERIC_02":
		return KnxPropertyDataType_PDT_GENERIC_02
	case "PDT_GENERIC_03":
		return KnxPropertyDataType_PDT_GENERIC_03
	case "PDT_UNSIGNED_CHAR":
		return KnxPropertyDataType_PDT_UNSIGNED_CHAR
	case "PDT_GENERIC_04":
		return KnxPropertyDataType_PDT_GENERIC_04
	case "PDT_GENERIC_05":
		return KnxPropertyDataType_PDT_GENERIC_05
	case "PDT_GENERIC_06":
		return KnxPropertyDataType_PDT_GENERIC_06
	case "PDT_GENERIC_07":
		return KnxPropertyDataType_PDT_GENERIC_07
	case "PDT_GENERIC_08":
		return KnxPropertyDataType_PDT_GENERIC_08
	case "PDT_GENERIC_09":
		return KnxPropertyDataType_PDT_GENERIC_09
	case "PDT_GENERIC_10":
		return KnxPropertyDataType_PDT_GENERIC_10
	case "PDT_GENERIC_11":
		return KnxPropertyDataType_PDT_GENERIC_11
	case "PDT_GENERIC_12":
		return KnxPropertyDataType_PDT_GENERIC_12
	case "PDT_GENERIC_13":
		return KnxPropertyDataType_PDT_GENERIC_13
	case "PDT_INT":
		return KnxPropertyDataType_PDT_INT
	case "PDT_GENERIC_14":
		return KnxPropertyDataType_PDT_GENERIC_14
	case "PDT_GENERIC_15":
		return KnxPropertyDataType_PDT_GENERIC_15
	case "PDT_GENERIC_16":
		return KnxPropertyDataType_PDT_GENERIC_16
	case "PDT_GENERIC_17":
		return KnxPropertyDataType_PDT_GENERIC_17
	case "PDT_GENERIC_18":
		return KnxPropertyDataType_PDT_GENERIC_18
	case "PDT_GENERIC_19":
		return KnxPropertyDataType_PDT_GENERIC_19
	case "PDT_GENERIC_20":
		return KnxPropertyDataType_PDT_GENERIC_20
	case "PDT_UNSIGNED_INT":
		return KnxPropertyDataType_PDT_UNSIGNED_INT
	case "PDT_UTF_8":
		return KnxPropertyDataType_PDT_UTF_8
	case "PDT_VERSION":
		return KnxPropertyDataType_PDT_VERSION
	case "PDT_ALARM_INFO":
		return KnxPropertyDataType_PDT_ALARM_INFO
	case "PDT_KNX_FLOAT":
		return KnxPropertyDataType_PDT_KNX_FLOAT
	case "PDT_BINARY_INFORMATION":
		return KnxPropertyDataType_PDT_BINARY_INFORMATION
	case "PDT_BITSET8":
		return KnxPropertyDataType_PDT_BITSET8
	case "PDT_BITSET16":
		return KnxPropertyDataType_PDT_BITSET16
	case "PDT_ENUM8":
		return KnxPropertyDataType_PDT_ENUM8
	case "PDT_SCALING":
		return KnxPropertyDataType_PDT_SCALING
	case "PDT_DATE":
		return KnxPropertyDataType_PDT_DATE
	case "PDT_NE_VL":
		return KnxPropertyDataType_PDT_NE_VL
	case "PDT_NE_FL":
		return KnxPropertyDataType_PDT_NE_FL
	case "PDT_FUNCTION":
		return KnxPropertyDataType_PDT_FUNCTION
	case "PDT_ESCAPE":
		return KnxPropertyDataType_PDT_ESCAPE
	case "PDT_TIME":
		return KnxPropertyDataType_PDT_TIME
	case "PDT_LONG":
		return KnxPropertyDataType_PDT_LONG
	case "PDT_UNSIGNED_LONG":
		return KnxPropertyDataType_PDT_UNSIGNED_LONG
	}
	return 0
}

func CastKnxPropertyDataType(structType interface{}) KnxPropertyDataType {
	castFunc := func(typ interface{}) KnxPropertyDataType {
		if sKnxPropertyDataType, ok := typ.(KnxPropertyDataType); ok {
			return sKnxPropertyDataType
		}
		return 0
	}
	return castFunc(structType)
}

func (m KnxPropertyDataType) LengthInBits() uint16 {
	return 8
}

func (m KnxPropertyDataType) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxPropertyDataTypeParse(io *utils.ReadBuffer) (KnxPropertyDataType, error) {
	val, err := io.ReadUint8(8)
	if err != nil {
		return 0, nil
	}
	return KnxPropertyDataTypeByValue(val), nil
}

func (e KnxPropertyDataType) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint8(8, uint8(e))
	return err
}

func (e KnxPropertyDataType) String() string {
	switch e {
	case KnxPropertyDataType_PDT_CONTROL:
		return "PDT_CONTROL"
	case KnxPropertyDataType_PDT_CHAR:
		return "PDT_CHAR"
	case KnxPropertyDataType_PDT_FLOAT:
		return "PDT_FLOAT"
	case KnxPropertyDataType_PDT_DOUBLE:
		return "PDT_DOUBLE"
	case KnxPropertyDataType_PDT_CHAR_BLOCK:
		return "PDT_CHAR_BLOCK"
	case KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS:
		return "PDT_POLL_GROUP_SETTINGS"
	case KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK:
		return "PDT_SHORT_CHAR_BLOCK"
	case KnxPropertyDataType_PDT_DATE_TIME:
		return "PDT_DATE_TIME"
	case KnxPropertyDataType_PDT_VARIABLE_LENGTH:
		return "PDT_VARIABLE_LENGTH"
	case KnxPropertyDataType_PDT_GENERIC_01:
		return "PDT_GENERIC_01"
	case KnxPropertyDataType_PDT_GENERIC_02:
		return "PDT_GENERIC_02"
	case KnxPropertyDataType_PDT_GENERIC_03:
		return "PDT_GENERIC_03"
	case KnxPropertyDataType_PDT_UNSIGNED_CHAR:
		return "PDT_UNSIGNED_CHAR"
	case KnxPropertyDataType_PDT_GENERIC_04:
		return "PDT_GENERIC_04"
	case KnxPropertyDataType_PDT_GENERIC_05:
		return "PDT_GENERIC_05"
	case KnxPropertyDataType_PDT_GENERIC_06:
		return "PDT_GENERIC_06"
	case KnxPropertyDataType_PDT_GENERIC_07:
		return "PDT_GENERIC_07"
	case KnxPropertyDataType_PDT_GENERIC_08:
		return "PDT_GENERIC_08"
	case KnxPropertyDataType_PDT_GENERIC_09:
		return "PDT_GENERIC_09"
	case KnxPropertyDataType_PDT_GENERIC_10:
		return "PDT_GENERIC_10"
	case KnxPropertyDataType_PDT_GENERIC_11:
		return "PDT_GENERIC_11"
	case KnxPropertyDataType_PDT_GENERIC_12:
		return "PDT_GENERIC_12"
	case KnxPropertyDataType_PDT_GENERIC_13:
		return "PDT_GENERIC_13"
	case KnxPropertyDataType_PDT_INT:
		return "PDT_INT"
	case KnxPropertyDataType_PDT_GENERIC_14:
		return "PDT_GENERIC_14"
	case KnxPropertyDataType_PDT_GENERIC_15:
		return "PDT_GENERIC_15"
	case KnxPropertyDataType_PDT_GENERIC_16:
		return "PDT_GENERIC_16"
	case KnxPropertyDataType_PDT_GENERIC_17:
		return "PDT_GENERIC_17"
	case KnxPropertyDataType_PDT_GENERIC_18:
		return "PDT_GENERIC_18"
	case KnxPropertyDataType_PDT_GENERIC_19:
		return "PDT_GENERIC_19"
	case KnxPropertyDataType_PDT_GENERIC_20:
		return "PDT_GENERIC_20"
	case KnxPropertyDataType_PDT_UNSIGNED_INT:
		return "PDT_UNSIGNED_INT"
	case KnxPropertyDataType_PDT_UTF_8:
		return "PDT_UTF_8"
	case KnxPropertyDataType_PDT_VERSION:
		return "PDT_VERSION"
	case KnxPropertyDataType_PDT_ALARM_INFO:
		return "PDT_ALARM_INFO"
	case KnxPropertyDataType_PDT_KNX_FLOAT:
		return "PDT_KNX_FLOAT"
	case KnxPropertyDataType_PDT_BINARY_INFORMATION:
		return "PDT_BINARY_INFORMATION"
	case KnxPropertyDataType_PDT_BITSET8:
		return "PDT_BITSET8"
	case KnxPropertyDataType_PDT_BITSET16:
		return "PDT_BITSET16"
	case KnxPropertyDataType_PDT_ENUM8:
		return "PDT_ENUM8"
	case KnxPropertyDataType_PDT_SCALING:
		return "PDT_SCALING"
	case KnxPropertyDataType_PDT_DATE:
		return "PDT_DATE"
	case KnxPropertyDataType_PDT_NE_VL:
		return "PDT_NE_VL"
	case KnxPropertyDataType_PDT_NE_FL:
		return "PDT_NE_FL"
	case KnxPropertyDataType_PDT_FUNCTION:
		return "PDT_FUNCTION"
	case KnxPropertyDataType_PDT_ESCAPE:
		return "PDT_ESCAPE"
	case KnxPropertyDataType_PDT_TIME:
		return "PDT_TIME"
	case KnxPropertyDataType_PDT_LONG:
		return "PDT_LONG"
	case KnxPropertyDataType_PDT_UNSIGNED_LONG:
		return "PDT_UNSIGNED_LONG"
	}
	return ""
}
