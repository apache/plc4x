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

type DataTransportSize uint8

type IDataTransportSize interface {
	spi.Message
	SizeInBits() bool
	Serialize(io spi.WriteBuffer)
}

const (
	DataTransportSize_NULL            DataTransportSize = 0x00
	DataTransportSize_BIT             DataTransportSize = 0x03
	DataTransportSize_BYTE_WORD_DWORD DataTransportSize = 0x04
	DataTransportSize_INTEGER         DataTransportSize = 0x05
	DataTransportSize_DINTEGER        DataTransportSize = 0x06
	DataTransportSize_REAL            DataTransportSize = 0x07
	DataTransportSize_OCTET_STRING    DataTransportSize = 0x09
)

func (e DataTransportSize) SizeInBits() bool {
	switch e {
	case 0x00:
		{ /* '0x00' */
			return false
		}
	case 0x03:
		{ /* '0x03' */
			return true
		}
	case 0x04:
		{ /* '0x04' */
			return true
		}
	case 0x05:
		{ /* '0x05' */
			return true
		}
	case 0x06:
		{ /* '0x06' */
			return false
		}
	case 0x07:
		{ /* '0x07' */
			return false
		}
	case 0x09:
		{ /* '0x09' */
			return false
		}
	default:
		{
			return false
		}
	}
}

func CastDataTransportSize(structType interface{}) DataTransportSize {
	castFunc := func(typ interface{}) DataTransportSize {
		if sDataTransportSize, ok := typ.(DataTransportSize); ok {
			return sDataTransportSize
		}
		return 0
	}
	return castFunc(structType)
}

func (m DataTransportSize) LengthInBits() uint16 {
	return 8
}

func (m DataTransportSize) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DataTransportSizeParse(io spi.ReadBuffer) (DataTransportSize, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e DataTransportSize) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
