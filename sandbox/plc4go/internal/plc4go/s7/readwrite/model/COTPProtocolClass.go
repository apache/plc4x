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

type COTPProtocolClass int8

type ICOTPProtocolClass interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

const (
	COTPProtocolClass_CLASS_0 COTPProtocolClass = 0x00
	COTPProtocolClass_CLASS_1 COTPProtocolClass = 0x10
	COTPProtocolClass_CLASS_2 COTPProtocolClass = 0x20
	COTPProtocolClass_CLASS_3 COTPProtocolClass = 0x30
	COTPProtocolClass_CLASS_4 COTPProtocolClass = 0x40
)

func CastCOTPProtocolClass(structType interface{}) COTPProtocolClass {
	castFunc := func(typ interface{}) COTPProtocolClass {
		if sCOTPProtocolClass, ok := typ.(COTPProtocolClass); ok {
			return sCOTPProtocolClass
		}
		return 0
	}
	return castFunc(structType)
}

func (m COTPProtocolClass) LengthInBits() uint16 {
	return 8
}

func (m COTPProtocolClass) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPProtocolClassParse(io *spi.ReadBuffer) (COTPProtocolClass, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e COTPProtocolClass) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
