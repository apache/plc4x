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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
)

// The data-structure of this message
type S7ParameterUserDataItem struct {

}

// The corresponding interface
type IS7ParameterUserDataItem interface {
    spi.Message
    ItemType() uint8
    Serialize(io spi.WriteBuffer) error
}

type S7ParameterUserDataItemInitializer interface {
    initialize() spi.Message
}

func S7ParameterUserDataItemItemType(m IS7ParameterUserDataItem) uint8 {
    return m.ItemType()
}


func CastIS7ParameterUserDataItem(structType interface{}) IS7ParameterUserDataItem {
    castFunc := func(typ interface{}) IS7ParameterUserDataItem {
        if iS7ParameterUserDataItem, ok := typ.(IS7ParameterUserDataItem); ok {
            return iS7ParameterUserDataItem
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7ParameterUserDataItem(structType interface{}) S7ParameterUserDataItem {
    castFunc := func(typ interface{}) S7ParameterUserDataItem {
        if sS7ParameterUserDataItem, ok := typ.(S7ParameterUserDataItem); ok {
            return sS7ParameterUserDataItem
        }
        if sS7ParameterUserDataItem, ok := typ.(*S7ParameterUserDataItem); ok {
            return *sS7ParameterUserDataItem
        }
        return S7ParameterUserDataItem{}
    }
    return castFunc(structType)
}

func (m S7ParameterUserDataItem) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Discriminator Field (itemType)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits
}

func (m S7ParameterUserDataItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7ParameterUserDataItemParse(io *spi.ReadBuffer) (spi.Message, error) {

    // Discriminator Field (itemType) (Used as input to a switch field)
    itemType, _itemTypeErr := io.ReadUint8(8)
    if _itemTypeErr != nil {
        return nil, errors.New("Error parsing 'itemType' field " + _itemTypeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer S7ParameterUserDataItemInitializer
    var typeSwitchError error
    switch {
    case itemType == 0x12:
        initializer, typeSwitchError = S7ParameterUserDataItemCPUFunctionsParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Create the instance
    return initializer.initialize(), nil
}

func S7ParameterUserDataItemSerialize(io spi.WriteBuffer, m S7ParameterUserDataItem, i IS7ParameterUserDataItem, childSerialize func() error) error {

    // Discriminator Field (itemType) (Used as input to a switch field)
    itemType := uint8(i.ItemType())
    _itemTypeErr := io.WriteUint8(8, (itemType))
    if _itemTypeErr != nil {
        return errors.New("Error serializing 'itemType' field " + _itemTypeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}
