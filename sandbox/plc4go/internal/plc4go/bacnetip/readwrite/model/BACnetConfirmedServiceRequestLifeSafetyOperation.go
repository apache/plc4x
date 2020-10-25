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
    "encoding/xml"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetConfirmedServiceRequestLifeSafetyOperation struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestLifeSafetyOperation interface {
    IBACnetConfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestLifeSafetyOperation) ServiceChoice() uint8 {
    return 0x1B
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestLifeSafetyOperation() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestLifeSafetyOperation{}
}

func CastIBACnetConfirmedServiceRequestLifeSafetyOperation(structType interface{}) IBACnetConfirmedServiceRequestLifeSafetyOperation {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestLifeSafetyOperation {
        if iBACnetConfirmedServiceRequestLifeSafetyOperation, ok := typ.(IBACnetConfirmedServiceRequestLifeSafetyOperation); ok {
            return iBACnetConfirmedServiceRequestLifeSafetyOperation
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestLifeSafetyOperation(structType interface{}) BACnetConfirmedServiceRequestLifeSafetyOperation {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestLifeSafetyOperation {
        if sBACnetConfirmedServiceRequestLifeSafetyOperation, ok := typ.(BACnetConfirmedServiceRequestLifeSafetyOperation); ok {
            return sBACnetConfirmedServiceRequestLifeSafetyOperation
        }
        if sBACnetConfirmedServiceRequestLifeSafetyOperation, ok := typ.(*BACnetConfirmedServiceRequestLifeSafetyOperation); ok {
            return *sBACnetConfirmedServiceRequestLifeSafetyOperation
        }
        return BACnetConfirmedServiceRequestLifeSafetyOperation{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestLifeSafetyOperationParse(io *utils.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestLifeSafetyOperation(), nil
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}

func (m *BACnetConfirmedServiceRequestLifeSafetyOperation) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            }
        }
    }
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestLifeSafetyOperation"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

