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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type BVLCReadBroadcastDistributionTableAck struct {
	BVLC
}

// The corresponding interface
type IBVLCReadBroadcastDistributionTableAck interface {
	IBVLC
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BVLCReadBroadcastDistributionTableAck) BvlcFunction() uint8 {
	return 0x03
}

func (m BVLCReadBroadcastDistributionTableAck) initialize() spi.Message {
	return m
}

func NewBVLCReadBroadcastDistributionTableAck() BVLCInitializer {
	return &BVLCReadBroadcastDistributionTableAck{}
}

func CastIBVLCReadBroadcastDistributionTableAck(structType interface{}) IBVLCReadBroadcastDistributionTableAck {
	castFunc := func(typ interface{}) IBVLCReadBroadcastDistributionTableAck {
		if iBVLCReadBroadcastDistributionTableAck, ok := typ.(IBVLCReadBroadcastDistributionTableAck); ok {
			return iBVLCReadBroadcastDistributionTableAck
		}
		return nil
	}
	return castFunc(structType)
}

func CastBVLCReadBroadcastDistributionTableAck(structType interface{}) BVLCReadBroadcastDistributionTableAck {
	castFunc := func(typ interface{}) BVLCReadBroadcastDistributionTableAck {
		if sBVLCReadBroadcastDistributionTableAck, ok := typ.(BVLCReadBroadcastDistributionTableAck); ok {
			return sBVLCReadBroadcastDistributionTableAck
		}
		return BVLCReadBroadcastDistributionTableAck{}
	}
	return castFunc(structType)
}

func (m BVLCReadBroadcastDistributionTableAck) LengthInBits() uint16 {
	var lengthInBits = m.BVLC.LengthInBits()

	return lengthInBits
}

func (m BVLCReadBroadcastDistributionTableAck) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BVLCReadBroadcastDistributionTableAckParse(io *spi.ReadBuffer) (BVLCInitializer, error) {

	// Create the instance
	return NewBVLCReadBroadcastDistributionTableAck(), nil
}

func (m BVLCReadBroadcastDistributionTableAck) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BVLCSerialize(io, m.BVLC, CastIBVLC(m), ser)
}
