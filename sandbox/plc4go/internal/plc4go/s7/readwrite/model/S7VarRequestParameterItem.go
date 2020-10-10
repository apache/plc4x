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
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type S7VarRequestParameterItem struct {
}

// The corresponding interface
type IS7VarRequestParameterItem interface {
	spi.Message
	ItemType() uint8
	Serialize(io spi.WriteBuffer)
}

type S7VarRequestParameterItemInitializer interface {
	initialize() spi.Message
}

func S7VarRequestParameterItemItemType(m IS7VarRequestParameterItem) uint8 {
	return m.ItemType()
}

func (m S7VarRequestParameterItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (itemType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m S7VarRequestParameterItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7VarRequestParameterItemParse(io spi.ReadBuffer) (spi.Message, error) {

	// Discriminator Field (itemType) (Used as input to a switch field)
	var itemType uint8 = io.ReadUint8(8)

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer S7VarRequestParameterItemInitializer
	var typeSwitchError error
	switch {
	case itemType == 0x12:
		initializer, typeSwitchError = S7VarRequestParameterItemAddressParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m S7VarRequestParameterItem) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if iS7VarRequestParameterItem, ok := typ.(IS7VarRequestParameterItem); ok {

			// Discriminator Field (itemType) (Used as input to a switch field)
			itemType := S7VarRequestParameterItemItemType(iS7VarRequestParameterItem)
			io.WriteUint8(8, (itemType))

			// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
			iS7VarRequestParameterItem.Serialize(io)
		}
	}
	serializeFunc(m)
}
