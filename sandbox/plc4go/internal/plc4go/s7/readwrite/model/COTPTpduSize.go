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

type COTPTpduSize int8

type ICOTPTpduSize interface {
	spi.Message
	SizeInBytes() uint16
	Serialize(io spi.WriteBuffer) error
}

const (
	COTPTpduSize_SIZE_128  COTPTpduSize = 0x07
	COTPTpduSize_SIZE_256  COTPTpduSize = 0x08
	COTPTpduSize_SIZE_512  COTPTpduSize = 0x09
	COTPTpduSize_SIZE_1024 COTPTpduSize = 0x0a
	COTPTpduSize_SIZE_2048 COTPTpduSize = 0x0b
	COTPTpduSize_SIZE_4096 COTPTpduSize = 0x0c
	COTPTpduSize_SIZE_8192 COTPTpduSize = 0x0d
)

func (e COTPTpduSize) SizeInBytes() uint16 {
	switch e {
	case 0x07:
		{ /* '0x07' */
			return 128
		}
	case 0x08:
		{ /* '0x08' */
			return 256
		}
	case 0x09:
		{ /* '0x09' */
			return 512
		}
	case 0x0a:
		{ /* '0x0a' */
			return 1024
		}
	case 0x0b:
		{ /* '0x0b' */
			return 2048
		}
	case 0x0c:
		{ /* '0x0c' */
			return 4096
		}
	case 0x0d:
		{ /* '0x0d' */
			return 8192
		}
	default:
		{
			return 0
		}
	}
}

func CastCOTPTpduSize(structType interface{}) COTPTpduSize {
	castFunc := func(typ interface{}) COTPTpduSize {
		if sCOTPTpduSize, ok := typ.(COTPTpduSize); ok {
			return sCOTPTpduSize
		}
		return 0
	}
	return castFunc(structType)
}

func (m COTPTpduSize) LengthInBits() uint16 {
	return 8
}

func (m COTPTpduSize) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPTpduSizeParse(io *spi.ReadBuffer) (COTPTpduSize, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e COTPTpduSize) Serialize(io spi.WriteBuffer) error {
	// TODO: Implement ...
	return nil
}
