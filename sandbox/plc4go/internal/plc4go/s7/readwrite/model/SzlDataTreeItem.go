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
type SzlDataTreeItem struct {
	itemIndex    uint16
	mlfb         []int8
	moduleTypeId uint16
	ausbg        uint16
	ausbe        uint16
}

// The corresponding interface
type ISzlDataTreeItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewSzlDataTreeItem(itemIndex uint16, mlfb []int8, moduleTypeId uint16, ausbg uint16, ausbe uint16) spi.Message {
	return &SzlDataTreeItem{itemIndex: itemIndex, mlfb: mlfb, moduleTypeId: moduleTypeId, ausbg: ausbg, ausbe: ausbe}
}

func CastISzlDataTreeItem(structType interface{}) ISzlDataTreeItem {
	castFunc := func(typ interface{}) ISzlDataTreeItem {
		if iSzlDataTreeItem, ok := typ.(ISzlDataTreeItem); ok {
			return iSzlDataTreeItem
		}
		return nil
	}
	return castFunc(structType)
}

func CastSzlDataTreeItem(structType interface{}) SzlDataTreeItem {
	castFunc := func(typ interface{}) SzlDataTreeItem {
		if sSzlDataTreeItem, ok := typ.(SzlDataTreeItem); ok {
			return sSzlDataTreeItem
		}
		return SzlDataTreeItem{}
	}
	return castFunc(structType)
}

func (m SzlDataTreeItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (itemIndex)
	lengthInBits += 16

	// Array field
	if len(m.mlfb) > 0 {
		lengthInBits += 8 * uint16(len(m.mlfb))
	}

	// Simple field (moduleTypeId)
	lengthInBits += 16

	// Simple field (ausbg)
	lengthInBits += 16

	// Simple field (ausbe)
	lengthInBits += 16

	return lengthInBits
}

func (m SzlDataTreeItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func SzlDataTreeItemParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (itemIndex)
	var itemIndex uint16 = io.ReadUint16(16)

	// Array field (mlfb)
	var mlfb []int8
	// Count array
	{
		mlfb := make([]int8, uint16(20))
		for curItem := uint16(0); curItem < uint16(uint16(20)); curItem++ {

			mlfb = append(mlfb, io.ReadInt8(8))
		}
	}

	// Simple Field (moduleTypeId)
	var moduleTypeId uint16 = io.ReadUint16(16)

	// Simple Field (ausbg)
	var ausbg uint16 = io.ReadUint16(16)

	// Simple Field (ausbe)
	var ausbe uint16 = io.ReadUint16(16)

	// Create the instance
	return NewSzlDataTreeItem(itemIndex, mlfb, moduleTypeId, ausbg, ausbe), nil
}

func (m SzlDataTreeItem) Serialize(io spi.WriteBuffer) {

	// Simple Field (itemIndex)
	itemIndex := uint16(m.itemIndex)
	io.WriteUint16(16, (itemIndex))

	// Array Field (mlfb)
	if m.mlfb != nil {
		for _, _element := range m.mlfb {
			io.WriteInt8(8, _element)
		}
	}

	// Simple Field (moduleTypeId)
	moduleTypeId := uint16(m.moduleTypeId)
	io.WriteUint16(16, (moduleTypeId))

	// Simple Field (ausbg)
	ausbg := uint16(m.ausbg)
	io.WriteUint16(16, (ausbg))

	// Simple Field (ausbe)
	ausbe := uint16(m.ausbe)
	io.WriteUint16(16, (ausbe))
}
