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

type BACnetNotifyType uint8

type IBACnetNotifyType interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

const (
	BACnetNotifyType_ALARM            BACnetNotifyType = 0x0
	BACnetNotifyType_EVENT            BACnetNotifyType = 0x1
	BACnetNotifyType_ACK_NOTIFICATION BACnetNotifyType = 0x2
)

func CastBACnetNotifyType(structType interface{}) BACnetNotifyType {
	castFunc := func(typ interface{}) BACnetNotifyType {
		if sBACnetNotifyType, ok := typ.(BACnetNotifyType); ok {
			return sBACnetNotifyType
		}
		return 0
	}
	return castFunc(structType)
}

func (m BACnetNotifyType) LengthInBits() uint16 {
	return 4
}

func (m BACnetNotifyType) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetNotifyTypeParse(io spi.ReadBuffer) (BACnetNotifyType, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e BACnetNotifyType) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
