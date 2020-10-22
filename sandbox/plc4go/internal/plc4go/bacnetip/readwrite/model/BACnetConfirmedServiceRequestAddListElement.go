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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
)

// The data-structure of this message
type BACnetConfirmedServiceRequestAddListElement struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestAddListElement interface {
    IBACnetConfirmedServiceRequest
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestAddListElement) ServiceChoice() uint8 {
    return 0x08
}

func (m BACnetConfirmedServiceRequestAddListElement) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestAddListElement() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestAddListElement{}
}

func CastIBACnetConfirmedServiceRequestAddListElement(structType interface{}) IBACnetConfirmedServiceRequestAddListElement {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestAddListElement {
        if iBACnetConfirmedServiceRequestAddListElement, ok := typ.(IBACnetConfirmedServiceRequestAddListElement); ok {
            return iBACnetConfirmedServiceRequestAddListElement
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestAddListElement(structType interface{}) BACnetConfirmedServiceRequestAddListElement {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestAddListElement {
        if sBACnetConfirmedServiceRequestAddListElement, ok := typ.(BACnetConfirmedServiceRequestAddListElement); ok {
            return sBACnetConfirmedServiceRequestAddListElement
        }
        if sBACnetConfirmedServiceRequestAddListElement, ok := typ.(*BACnetConfirmedServiceRequestAddListElement); ok {
            return *sBACnetConfirmedServiceRequestAddListElement
        }
        return BACnetConfirmedServiceRequestAddListElement{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestAddListElement) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestAddListElement) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestAddListElementParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestAddListElement(), nil
}

func (m BACnetConfirmedServiceRequestAddListElement) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
