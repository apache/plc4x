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
type BVLCWideBroadcastDistributionTable struct {
	BVLC
}

// The corresponding interface
type IBVLCWideBroadcastDistributionTable interface {
	IBVLC
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BVLCWideBroadcastDistributionTable) BvlcFunction() uint8 {
	return 0x01
}

func (m BVLCWideBroadcastDistributionTable) initialize() spi.Message {
	return m
}

func NewBVLCWideBroadcastDistributionTable() BVLCInitializer {
	return &BVLCWideBroadcastDistributionTable{}
}

func CastIBVLCWideBroadcastDistributionTable(structType interface{}) IBVLCWideBroadcastDistributionTable {
	castFunc := func(typ interface{}) IBVLCWideBroadcastDistributionTable {
		if iBVLCWideBroadcastDistributionTable, ok := typ.(IBVLCWideBroadcastDistributionTable); ok {
			return iBVLCWideBroadcastDistributionTable
		}
		return nil
	}
	return castFunc(structType)
}

func CastBVLCWideBroadcastDistributionTable(structType interface{}) BVLCWideBroadcastDistributionTable {
	castFunc := func(typ interface{}) BVLCWideBroadcastDistributionTable {
		if sBVLCWideBroadcastDistributionTable, ok := typ.(BVLCWideBroadcastDistributionTable); ok {
			return sBVLCWideBroadcastDistributionTable
		}
		return BVLCWideBroadcastDistributionTable{}
	}
	return castFunc(structType)
}

func (m BVLCWideBroadcastDistributionTable) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BVLC.LengthInBits()

	return lengthInBits
}

func (m BVLCWideBroadcastDistributionTable) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BVLCWideBroadcastDistributionTableParse(io *spi.ReadBuffer) (BVLCInitializer, error) {

	// Create the instance
	return NewBVLCWideBroadcastDistributionTable(), nil
}

func (m BVLCWideBroadcastDistributionTable) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BVLCSerialize(io, m.BVLC, CastIBVLC(m), ser)
}
