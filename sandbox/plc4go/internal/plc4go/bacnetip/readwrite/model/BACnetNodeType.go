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

type BACnetNodeType uint8

const (
	BACnetNodeType_UNKNOWN        BACnetNodeType = 0x00
	BACnetNodeType_SYSTEM         BACnetNodeType = 0x01
	BACnetNodeType_NETWORK        BACnetNodeType = 0x02
	BACnetNodeType_DEVICE         BACnetNodeType = 0x03
	BACnetNodeType_ORGANIZATIONAL BACnetNodeType = 0x04
	BACnetNodeType_AREA           BACnetNodeType = 0x05
	BACnetNodeType_EQUIPMENT      BACnetNodeType = 0x06
	BACnetNodeType_POINT          BACnetNodeType = 0x07
	BACnetNodeType_COLLECTION     BACnetNodeType = 0x08
	BACnetNodeType_PROPERTY       BACnetNodeType = 0x09
	BACnetNodeType_FUNCTIONAL     BACnetNodeType = 0x0A
	BACnetNodeType_OTHER          BACnetNodeType = 0x0B
	BACnetNodeType_SUBSYSTEM      BACnetNodeType = 0x0C
	BACnetNodeType_BUILDING       BACnetNodeType = 0x0D
	BACnetNodeType_FLOOR          BACnetNodeType = 0x0E
	BACnetNodeType_SECTION        BACnetNodeType = 0x0F
	BACnetNodeType_MODULE         BACnetNodeType = 0x10
	BACnetNodeType_TREE           BACnetNodeType = 0x11
	BACnetNodeType_MEMBER         BACnetNodeType = 0x12
	BACnetNodeType_PROTOCOL       BACnetNodeType = 0x13
	BACnetNodeType_ROOM           BACnetNodeType = 0x14
	BACnetNodeType_ZONE           BACnetNodeType = 0x15
)

func BACnetNodeTypeParse(io spi.ReadBuffer) (BACnetNodeType, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e BACnetNodeType) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
