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
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetTagApplicationDouble struct {
    Value float64
    BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationDouble interface {
    IBACnetTag
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetTagApplicationDouble) ContextSpecificTag() uint8 {
    return 0
}

func (m BACnetTagApplicationDouble) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
    m.TypeOrTagNumber = typeOrTagNumber
    m.LengthValueType = lengthValueType
    m.ExtTagNumber = extTagNumber
    m.ExtLength = extLength
    return m
}

func NewBACnetTagApplicationDouble(value float64) BACnetTagInitializer {
    return &BACnetTagApplicationDouble{Value: value}
}

func CastIBACnetTagApplicationDouble(structType interface{}) IBACnetTagApplicationDouble {
    castFunc := func(typ interface{}) IBACnetTagApplicationDouble {
        if iBACnetTagApplicationDouble, ok := typ.(IBACnetTagApplicationDouble); ok {
            return iBACnetTagApplicationDouble
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetTagApplicationDouble(structType interface{}) BACnetTagApplicationDouble {
    castFunc := func(typ interface{}) BACnetTagApplicationDouble {
        if sBACnetTagApplicationDouble, ok := typ.(BACnetTagApplicationDouble); ok {
            return sBACnetTagApplicationDouble
        }
        if sBACnetTagApplicationDouble, ok := typ.(*BACnetTagApplicationDouble); ok {
            return *sBACnetTagApplicationDouble
        }
        return BACnetTagApplicationDouble{}
    }
    return castFunc(structType)
}

func (m BACnetTagApplicationDouble) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetTag.LengthInBits()

    // Simple field (value)
    lengthInBits += 64

    return lengthInBits
}

func (m BACnetTagApplicationDouble) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagApplicationDoubleParse(io *utils.ReadBuffer, lengthValueType uint8, extLength uint8) (BACnetTagInitializer, error) {

    // Simple Field (value)
    value, _valueErr := io.ReadFloat64(64)
    if _valueErr != nil {
        return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
    }

    // Create the instance
    return NewBACnetTagApplicationDouble(value), nil
}

func (m BACnetTagApplicationDouble) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (value)
    value := float64(m.Value)
    _valueErr := io.WriteFloat64(64, (value))
    if _valueErr != nil {
        return errors.New("Error serializing 'value' field " + _valueErr.Error())
    }

        return nil
    }
    return BACnetTagSerialize(io, m.BACnetTag, CastIBACnetTag(m), ser)
}
