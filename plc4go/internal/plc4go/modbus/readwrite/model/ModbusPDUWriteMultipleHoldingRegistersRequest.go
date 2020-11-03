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
    "encoding/base64"
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDUWriteMultipleHoldingRegistersRequest struct {
    StartingAddress uint16
    Quantity uint16
    Value []int8
    ModbusPDU
}

// The corresponding interface
type IModbusPDUWriteMultipleHoldingRegistersRequest interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUWriteMultipleHoldingRegistersRequest) ErrorFlag() bool {
    return false
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) FunctionFlag() uint8 {
    return 0x10
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) Response() bool {
    return false
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) initialize() spi.Message {
    return m
}

func NewModbusPDUWriteMultipleHoldingRegistersRequest(startingAddress uint16, quantity uint16, value []int8) ModbusPDUInitializer {
    return &ModbusPDUWriteMultipleHoldingRegistersRequest{StartingAddress: startingAddress, Quantity: quantity, Value: value}
}

func CastIModbusPDUWriteMultipleHoldingRegistersRequest(structType interface{}) IModbusPDUWriteMultipleHoldingRegistersRequest {
    castFunc := func(typ interface{}) IModbusPDUWriteMultipleHoldingRegistersRequest {
        if iModbusPDUWriteMultipleHoldingRegistersRequest, ok := typ.(IModbusPDUWriteMultipleHoldingRegistersRequest); ok {
            return iModbusPDUWriteMultipleHoldingRegistersRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUWriteMultipleHoldingRegistersRequest(structType interface{}) ModbusPDUWriteMultipleHoldingRegistersRequest {
    castFunc := func(typ interface{}) ModbusPDUWriteMultipleHoldingRegistersRequest {
        if sModbusPDUWriteMultipleHoldingRegistersRequest, ok := typ.(ModbusPDUWriteMultipleHoldingRegistersRequest); ok {
            return sModbusPDUWriteMultipleHoldingRegistersRequest
        }
        if sModbusPDUWriteMultipleHoldingRegistersRequest, ok := typ.(*ModbusPDUWriteMultipleHoldingRegistersRequest); ok {
            return *sModbusPDUWriteMultipleHoldingRegistersRequest
        }
        return ModbusPDUWriteMultipleHoldingRegistersRequest{}
    }
    return castFunc(structType)
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (startingAddress)
    lengthInBits += 16

    // Simple field (quantity)
    lengthInBits += 16

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.Value) > 0 {
        lengthInBits += 8 * uint16(len(m.Value))
    }

    return lengthInBits
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUWriteMultipleHoldingRegistersRequestParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (startingAddress)
    startingAddress, _startingAddressErr := io.ReadUint16(16)
    if _startingAddressErr != nil {
        return nil, errors.New("Error parsing 'startingAddress' field " + _startingAddressErr.Error())
    }

    // Simple Field (quantity)
    quantity, _quantityErr := io.ReadUint16(16)
    if _quantityErr != nil {
        return nil, errors.New("Error parsing 'quantity' field " + _quantityErr.Error())
    }

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount, _byteCountErr := io.ReadUint8(8)
    if _byteCountErr != nil {
        return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array field (value)
    // Count array
    value := make([]int8, byteCount)
    for curItem := uint16(0); curItem < uint16(byteCount); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'value' field " + _err.Error())
        }
        value[curItem] = _item
    }

    // Create the instance
    return NewModbusPDUWriteMultipleHoldingRegistersRequest(startingAddress, quantity, value), nil
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (startingAddress)
    startingAddress := uint16(m.StartingAddress)
    _startingAddressErr := io.WriteUint16(16, (startingAddress))
    if _startingAddressErr != nil {
        return errors.New("Error serializing 'startingAddress' field " + _startingAddressErr.Error())
    }

    // Simple Field (quantity)
    quantity := uint16(m.Quantity)
    _quantityErr := io.WriteUint16(16, (quantity))
    if _quantityErr != nil {
        return errors.New("Error serializing 'quantity' field " + _quantityErr.Error())
    }

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount := uint8(uint8(len(m.Value)))
    _byteCountErr := io.WriteUint8(8, (byteCount))
    if _byteCountErr != nil {
        return errors.New("Error serializing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array Field (value)
    if m.Value != nil {
        for _, _element := range m.Value {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'value' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}

func (m *ModbusPDUWriteMultipleHoldingRegistersRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "startingAddress":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.StartingAddress = data
            case "quantity":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Quantity = data
            case "value":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.Value = utils.ByteToInt8(_decoded[0:_len])
            }
        }
    }
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.StartingAddress, xml.StartElement{Name: xml.Name{Local: "startingAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Quantity, xml.StartElement{Name: xml.Name{Local: "quantity"}}); err != nil {
        return err
    }
    _encodedValue := make([]byte, base64.StdEncoding.EncodedLen(len(m.Value)))
    base64.StdEncoding.Encode(_encodedValue, utils.Int8ToByte(m.Value))
    if err := e.EncodeElement(_encodedValue, xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

