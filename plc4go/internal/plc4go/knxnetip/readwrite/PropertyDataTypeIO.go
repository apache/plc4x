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
package readwrite

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

func ParsePropertyDataType(io utils.ReadBuffer, dataType model.KnxPropertyDataType, numBytes uint8) values.PlcValue {
	// Helper for parsing arrays.
	parseArray := func(num uint8) values2.PlcByteArray {
		var data []byte
		for i := uint8(0); i < num; i++ {
			b, _ := io.ReadUint8("", 8)
			data = append(data, b)
		}
		return values2.NewPlcByteArray(data)
	}

	switch dataType {
	case model.KnxPropertyDataType_PDT_UNKNOWN: // Var Number of Bytes
		return parseArray(numBytes)
	case model.KnxPropertyDataType_PDT_CONTROL: // 10 Bytes
		return parseArray(10)
	case model.KnxPropertyDataType_PDT_CHAR: // 1 Byte
		val, _ := io.ReadUint8("", 8)
		return values2.NewPlcCHAR(val)
	case model.KnxPropertyDataType_PDT_UNSIGNED_CHAR: // 1 Byte
		val, _ := io.ReadUint8("", 8)
		return values2.NewPlcCHAR(val)
	case model.KnxPropertyDataType_PDT_INT: // 2 Bytes
		val, _ := io.ReadInt16("", 16)
		return values2.NewPlcINT(val)
	case model.KnxPropertyDataType_PDT_UNSIGNED_INT: // 2 Bytes
		val, _ := io.ReadUint16("", 16)
		return values2.NewPlcUINT(val)
	case model.KnxPropertyDataType_PDT_KNX_FLOAT: // 2 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_DATE: // 3 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_TIME: // 3 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_LONG: // 4 Bytes
		val, _ := io.ReadInt64("", 64)
		return values2.NewPlcLINT(val)
	case model.KnxPropertyDataType_PDT_UNSIGNED_LONG: // 4 Bytes
		val, _ := io.ReadUint64("", 64)
		return values2.NewPlcULINT(val)
	case model.KnxPropertyDataType_PDT_FLOAT: // 4 Bytes
		val, _ := io.ReadFloat32("", true, 8, 23)
		return values2.NewPlcREAL(val)
	case model.KnxPropertyDataType_PDT_DOUBLE: // 8 Bytes
		val, _ := io.ReadFloat64("", true, 11, 52)
		return values2.NewPlcLREAL(val)
	case model.KnxPropertyDataType_PDT_CHAR_BLOCK: // 10 Bytes
		return parseArray(10)
	case model.KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS: // 3 Bytes
		return parseArray(3)
	case model.KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK: // 5 Bytes
		return parseArray(5)
	case model.KnxPropertyDataType_PDT_DATE_TIME: // 8 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_VARIABLE_LENGTH: // Var Number of Bytes
		return parseArray(numBytes)
	case model.KnxPropertyDataType_PDT_GENERIC_01: // 1 Bytes
		return parseArray(1)
	case model.KnxPropertyDataType_PDT_GENERIC_02: // 2 Bytes
		return parseArray(2)
	case model.KnxPropertyDataType_PDT_GENERIC_03: // 3 Bytes
		return parseArray(3)
	case model.KnxPropertyDataType_PDT_GENERIC_04: // 4 Bytes
		return parseArray(4)
	case model.KnxPropertyDataType_PDT_GENERIC_05: // 5 Bytes
		return parseArray(5)
	case model.KnxPropertyDataType_PDT_GENERIC_06: // 6 Bytes
		return parseArray(6)
	case model.KnxPropertyDataType_PDT_GENERIC_07: // 7 Bytes
		return parseArray(7)
	case model.KnxPropertyDataType_PDT_GENERIC_08: // 8 Bytes
		return parseArray(8)
	case model.KnxPropertyDataType_PDT_GENERIC_09: // 9 Bytes
		return parseArray(9)
	case model.KnxPropertyDataType_PDT_GENERIC_10: // 10 Bytes
		return parseArray(10)
	case model.KnxPropertyDataType_PDT_GENERIC_11: // 11 Bytes
		return parseArray(11)
	case model.KnxPropertyDataType_PDT_GENERIC_12: // 12 Bytes
		return parseArray(12)
	case model.KnxPropertyDataType_PDT_GENERIC_13: // 13 Bytes
		return parseArray(13)
	case model.KnxPropertyDataType_PDT_GENERIC_14: // 14 Bytes
		return parseArray(14)
	case model.KnxPropertyDataType_PDT_GENERIC_15: // 15 Bytes
		return parseArray(15)
	case model.KnxPropertyDataType_PDT_GENERIC_16: // 16 Bytes
		return parseArray(16)
	case model.KnxPropertyDataType_PDT_GENERIC_17: // 17 Bytes
		return parseArray(17)
	case model.KnxPropertyDataType_PDT_GENERIC_18: // 18 Bytes
		return parseArray(18)
	case model.KnxPropertyDataType_PDT_GENERIC_19: // 19 Bytes
		return parseArray(19)
	case model.KnxPropertyDataType_PDT_GENERIC_20: // 20 Bytes
		return parseArray(20)
	case model.KnxPropertyDataType_PDT_UTF_8: // Var Number of Bytes
		return parseArray(numBytes)
	case model.KnxPropertyDataType_PDT_VERSION: // 2 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_ALARM_INFO: // 6 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_BINARY_INFORMATION: // 1 Bytes
		// TODO: Implement this ...
	case model.KnxPropertyDataType_PDT_BITSET8: // 1 Bytes
		val, _ := io.ReadUint8("", 8)
		return values2.NewPlcBitString(val)
	case model.KnxPropertyDataType_PDT_BITSET16: // 2 Bytes
		val, _ := io.ReadUint16("", 16)
		return values2.NewPlcBitString(val)
	case model.KnxPropertyDataType_PDT_ENUM8: // 1 Bytes
		val, _ := io.ReadUint8("", 8)
		return values2.NewPlcCHAR(val)
	case model.KnxPropertyDataType_PDT_SCALING: // 1 Bytes
		val, _ := io.ReadUint8("", 8)
		return values2.NewPlcCHAR(val)
	case model.KnxPropertyDataType_PDT_NE_VL: // Var Number of Bytes
		return parseArray(numBytes)
	case model.KnxPropertyDataType_PDT_NE_FL: // Var Number of Bytes
		return parseArray(numBytes)
	case model.KnxPropertyDataType_PDT_FUNCTION: // Var Number of Bytes
		return parseArray(numBytes)
	case model.KnxPropertyDataType_PDT_ESCAPE: // Var Number of Bytes
		return parseArray(numBytes)
	default:
		return parseArray(numBytes)
	}
	return nil
}
