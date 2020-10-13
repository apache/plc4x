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
type BACnetErrorAtomicReadFile struct {
	BACnetError
}

// The corresponding interface
type IBACnetErrorAtomicReadFile interface {
	IBACnetError
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetErrorAtomicReadFile) ServiceChoice() uint8 {
	return 0x06
}

func (m BACnetErrorAtomicReadFile) initialize() spi.Message {
	return m
}

func NewBACnetErrorAtomicReadFile() BACnetErrorInitializer {
	return &BACnetErrorAtomicReadFile{}
}

func CastIBACnetErrorAtomicReadFile(structType interface{}) IBACnetErrorAtomicReadFile {
	castFunc := func(typ interface{}) IBACnetErrorAtomicReadFile {
		if iBACnetErrorAtomicReadFile, ok := typ.(IBACnetErrorAtomicReadFile); ok {
			return iBACnetErrorAtomicReadFile
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetErrorAtomicReadFile(structType interface{}) BACnetErrorAtomicReadFile {
	castFunc := func(typ interface{}) BACnetErrorAtomicReadFile {
		if sBACnetErrorAtomicReadFile, ok := typ.(BACnetErrorAtomicReadFile); ok {
			return sBACnetErrorAtomicReadFile
		}
		return BACnetErrorAtomicReadFile{}
	}
	return castFunc(structType)
}

func (m BACnetErrorAtomicReadFile) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetError.LengthInBits()

	return lengthInBits
}

func (m BACnetErrorAtomicReadFile) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetErrorAtomicReadFileParse(io *spi.ReadBuffer) (BACnetErrorInitializer, error) {

	// Create the instance
	return NewBACnetErrorAtomicReadFile(), nil
}

func (m BACnetErrorAtomicReadFile) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetErrorSerialize(io, m.BACnetError, CastIBACnetError(m), ser)
}
