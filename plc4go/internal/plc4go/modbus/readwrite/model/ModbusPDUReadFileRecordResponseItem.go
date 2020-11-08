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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDUReadFileRecordResponseItem struct {
    ReferenceType uint8
    Data []int8
    IModbusPDUReadFileRecordResponseItem
}

// The corresponding interface
type IModbusPDUReadFileRecordResponseItem interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewModbusPDUReadFileRecordResponseItem(referenceType uint8, data []int8) *ModbusPDUReadFileRecordResponseItem {
    return &ModbusPDUReadFileRecordResponseItem{ReferenceType: referenceType, Data: data}
}

func CastModbusPDUReadFileRecordResponseItem(structType interface{}) ModbusPDUReadFileRecordResponseItem {
    castFunc := func(typ interface{}) ModbusPDUReadFileRecordResponseItem {
        if casted, ok := typ.(ModbusPDUReadFileRecordResponseItem); ok {
            return casted
        }
        if casted, ok := typ.(*ModbusPDUReadFileRecordResponseItem); ok {
            return *casted
        }
        return ModbusPDUReadFileRecordResponseItem{}
    }
    return castFunc(structType)
}

func (m *ModbusPDUReadFileRecordResponseItem) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Implicit Field (dataLength)
    lengthInBits += 8

    // Simple field (referenceType)
    lengthInBits += 8

    // Array field
    if len(m.Data) > 0 {
        lengthInBits += 8 * uint16(len(m.Data))
    }

    return lengthInBits
}

func (m *ModbusPDUReadFileRecordResponseItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordResponseItemParse(io *utils.ReadBuffer) (*ModbusPDUReadFileRecordResponseItem, error) {

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    dataLength, _dataLengthErr := io.ReadUint8(8)
    if _dataLengthErr != nil {
        return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Simple Field (referenceType)
    referenceType, _referenceTypeErr := io.ReadUint8(8)
    if _referenceTypeErr != nil {
        return nil, errors.New("Error parsing 'referenceType' field " + _referenceTypeErr.Error())
    }

    // Array field (data)
    // Length array
    data := make([]int8, 0)
    _dataLength := uint16(dataLength) - uint16(uint16(1))
    _dataEndPos := io.GetPos() + uint16(_dataLength)
    for ;io.GetPos() < _dataEndPos; {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'data' field " + _err.Error())
        }
        data = append(data, _item)
    }

    // Create the instance
    return NewModbusPDUReadFileRecordResponseItem(referenceType, data), nil
}

func (m *ModbusPDUReadFileRecordResponseItem) Serialize(io utils.WriteBuffer) error {

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    dataLength := uint8(uint8(uint8(len(m.Data))) + uint8(uint8(1)))
    _dataLengthErr := io.WriteUint8(8, (dataLength))
    if _dataLengthErr != nil {
        return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Simple Field (referenceType)
    referenceType := uint8(m.ReferenceType)
    _referenceTypeErr := io.WriteUint8(8, (referenceType))
    if _referenceTypeErr != nil {
        return errors.New("Error serializing 'referenceType' field " + _referenceTypeErr.Error())
    }

    // Array Field (data)
    if m.Data != nil {
        for _, _element := range m.Data {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'data' field " + _elementErr.Error())
            }
        }
    }

    return nil
}

func (m *ModbusPDUReadFileRecordResponseItem) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
            case "referenceType":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ReferenceType = data
            case "data":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.Data = utils.ByteToInt8(_decoded[0:_len])
            }
        }
    }
}

func (m *ModbusPDUReadFileRecordResponseItem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFileRecordResponseItem"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ReferenceType, xml.StartElement{Name: xml.Name{Local: "referenceType"}}); err != nil {
        return err
    }
    _encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(m.Data)))
    base64.StdEncoding.Encode(_encodedData, utils.Int8ToByte(m.Data))
    if err := e.EncodeElement(_encodedData, xml.StartElement{Name: xml.Name{Local: "data"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

