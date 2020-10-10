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
type BVLCSecureBVLL struct {
	BVLC
}

// The corresponding interface
type IBVLCSecureBVLL interface {
	IBVLC
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BVLCSecureBVLL) BvlcFunction() uint8 {
	return 0x0C
}

func (m BVLCSecureBVLL) initialize() spi.Message {
	return m
}

func NewBVLCSecureBVLL() BVLCInitializer {
	return &BVLCSecureBVLL{}
}

func (m BVLCSecureBVLL) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BVLC.LengthInBits()

	return lengthInBits
}

func (m BVLCSecureBVLL) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BVLCSecureBVLLParse(io spi.ReadBuffer) (BVLCInitializer, error) {

	// Create the instance
	return NewBVLCSecureBVLL(), nil
}

func (m BVLCSecureBVLL) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IBVLCSecureBVLL); ok {
		}
	}
	serializeFunc(m)
}
