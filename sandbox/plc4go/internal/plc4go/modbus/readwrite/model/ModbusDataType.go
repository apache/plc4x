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

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"

type ModbusDataType uint8

const (
	NULL           ModbusDataType = 00
	BOOL           ModbusDataType = 01
	BYTE           ModbusDataType = 10
	WORD           ModbusDataType = 11
	DWORD          ModbusDataType = 12
	LWORD          ModbusDataType = 13
	SINT           ModbusDataType = 20
	INT            ModbusDataType = 21
	DINT           ModbusDataType = 22
	LINT           ModbusDataType = 23
	USINT          ModbusDataType = 24
	UINT           ModbusDataType = 25
	UDINT          ModbusDataType = 26
	ULINT          ModbusDataType = 27
	REAL           ModbusDataType = 30
	LREAL          ModbusDataType = 31
	TIME           ModbusDataType = 40
	LTIME          ModbusDataType = 41
	DATE           ModbusDataType = 50
	LDATE          ModbusDataType = 51
	TIME_OF_DAY    ModbusDataType = 60
	LTIME_OF_DAY   ModbusDataType = 61
	DATE_AND_TIME  ModbusDataType = 70
	LDATE_AND_TIME ModbusDataType = 71
	CHAR           ModbusDataType = 80
	WCHAR          ModbusDataType = 81
	STRING         ModbusDataType = 82
	WSTRING        ModbusDataType = 83
)

func (e ModbusDataType) GetDataTypeSize() uint8 {
	switch e {
	case 00:
		{ /* '00' */
			return 0
		}
	case 01:
		{ /* '01' */
			return 1
		}
	case 10:
		{ /* '10' */
			return 1
		}
	case 11:
		{ /* '11' */
			return 2
		}
	case 12:
		{ /* '12' */
			return 4
		}
	case 13:
		{ /* '13' */
			return 8
		}
	case 20:
		{ /* '20' */
			return 1
		}
	case 21:
		{ /* '21' */
			return 2
		}
	case 22:
		{ /* '22' */
			return 4
		}
	case 23:
		{ /* '23' */
			return 8
		}
	case 24:
		{ /* '24' */
			return 1
		}
	case 25:
		{ /* '25' */
			return 2
		}
	case 26:
		{ /* '26' */
			return 4
		}
	case 27:
		{ /* '27' */
			return 8
		}
	case 30:
		{ /* '30' */
			return 4
		}
	case 31:
		{ /* '31' */
			return 8
		}
	case 40:
		{ /* '40' */
			return 8
		}
	case 41:
		{ /* '41' */
			return 8
		}
	case 50:
		{ /* '50' */
			return 8
		}
	case 51:
		{ /* '51' */
			return 8
		}
	case 60:
		{ /* '60' */
			return 8
		}
	case 61:
		{ /* '61' */
			return 8
		}
	case 70:
		{ /* '70' */
			return 8
		}
	case 71:
		{ /* '71' */
			return 8
		}
	case 80:
		{ /* '80' */
			return 1
		}
	case 81:
		{ /* '81' */
			return 2
		}
	case 82:
		{ /* '82' */
			return 1
		}
	case 83:
		{ /* '83' */
			return 2
		}
	default:
		{
			return 0
		}
	}
}

func (e *ModbusDataType) Parse(io spi.ReadBuffer) {
	// TODO: Implement ...
}

func (e ModbusDataType) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
