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
type SzlDataTreeItem struct {
    ItemIndex uint16
    Mlfb []int8
    ModuleTypeId uint16
    Ausbg uint16
    Ausbe uint16

}

// The corresponding interface
type ISzlDataTreeItem interface {
    spi.Message
    Serialize(io spi.WriteBuffer) error
}


func NewSzlDataTreeItem(itemIndex uint16, mlfb []int8, moduleTypeId uint16, ausbg uint16, ausbe uint16) spi.Message {
    return &SzlDataTreeItem{ItemIndex: itemIndex, Mlfb: mlfb, ModuleTypeId: moduleTypeId, Ausbg: ausbg, Ausbe: ausbe}
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
        if sSzlDataTreeItem, ok := typ.(*SzlDataTreeItem); ok {
            return *sSzlDataTreeItem
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
    if len(m.Mlfb) > 0 {
        lengthInBits += 8 * uint16(len(m.Mlfb))
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

func SzlDataTreeItemParse(io *spi.ReadBuffer) (spi.Message, error) {

    // Simple Field (itemIndex)
    itemIndex, _itemIndexErr := io.ReadUint16(16)
    if _itemIndexErr != nil {
        return nil, errors.New("Error parsing 'itemIndex' field " + _itemIndexErr.Error())
    }

    // Array field (mlfb)
    // Count array
    mlfb := make([]int8, uint16(20))
    for curItem := uint16(0); curItem < uint16(uint16(20)); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'mlfb' field " + _err.Error())
        }
        mlfb[curItem] = _item
    }

    // Simple Field (moduleTypeId)
    moduleTypeId, _moduleTypeIdErr := io.ReadUint16(16)
    if _moduleTypeIdErr != nil {
        return nil, errors.New("Error parsing 'moduleTypeId' field " + _moduleTypeIdErr.Error())
    }

    // Simple Field (ausbg)
    ausbg, _ausbgErr := io.ReadUint16(16)
    if _ausbgErr != nil {
        return nil, errors.New("Error parsing 'ausbg' field " + _ausbgErr.Error())
    }

    // Simple Field (ausbe)
    ausbe, _ausbeErr := io.ReadUint16(16)
    if _ausbeErr != nil {
        return nil, errors.New("Error parsing 'ausbe' field " + _ausbeErr.Error())
    }

    // Create the instance
    return NewSzlDataTreeItem(itemIndex, mlfb, moduleTypeId, ausbg, ausbe), nil
}

func (m SzlDataTreeItem) Serialize(io spi.WriteBuffer) error {

    // Simple Field (itemIndex)
    itemIndex := uint16(m.ItemIndex)
    _itemIndexErr := io.WriteUint16(16, (itemIndex))
    if _itemIndexErr != nil {
        return errors.New("Error serializing 'itemIndex' field " + _itemIndexErr.Error())
    }

    // Array Field (mlfb)
    if m.Mlfb != nil {
        for _, _element := range m.Mlfb {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'mlfb' field " + _elementErr.Error())
            }
        }
    }

    // Simple Field (moduleTypeId)
    moduleTypeId := uint16(m.ModuleTypeId)
    _moduleTypeIdErr := io.WriteUint16(16, (moduleTypeId))
    if _moduleTypeIdErr != nil {
        return errors.New("Error serializing 'moduleTypeId' field " + _moduleTypeIdErr.Error())
    }

    // Simple Field (ausbg)
    ausbg := uint16(m.Ausbg)
    _ausbgErr := io.WriteUint16(16, (ausbg))
    if _ausbgErr != nil {
        return errors.New("Error serializing 'ausbg' field " + _ausbgErr.Error())
    }

    // Simple Field (ausbe)
    ausbe := uint16(m.Ausbe)
    _ausbeErr := io.WriteUint16(16, (ausbe))
    if _ausbeErr != nil {
        return errors.New("Error serializing 'ausbe' field " + _ausbeErr.Error())
    }

    return nil
}
