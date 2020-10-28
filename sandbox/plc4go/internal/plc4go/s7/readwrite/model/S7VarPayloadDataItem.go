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
    "math"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type S7VarPayloadDataItem struct {
    ReturnCode IDataTransportErrorCode
    TransportSize IDataTransportSize
    Data []int8

}

// The corresponding interface
type IS7VarPayloadDataItem interface {
    spi.Message
    Serialize(io utils.WriteBuffer, lastItem bool) error
}


func NewS7VarPayloadDataItem(returnCode IDataTransportErrorCode, transportSize IDataTransportSize, data []int8) spi.Message {
    return &S7VarPayloadDataItem{ReturnCode: returnCode, TransportSize: transportSize, Data: data}
}

func CastIS7VarPayloadDataItem(structType interface{}) IS7VarPayloadDataItem {
    castFunc := func(typ interface{}) IS7VarPayloadDataItem {
        if iS7VarPayloadDataItem, ok := typ.(IS7VarPayloadDataItem); ok {
            return iS7VarPayloadDataItem
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7VarPayloadDataItem(structType interface{}) S7VarPayloadDataItem {
    castFunc := func(typ interface{}) S7VarPayloadDataItem {
        if sS7VarPayloadDataItem, ok := typ.(S7VarPayloadDataItem); ok {
            return sS7VarPayloadDataItem
        }
        if sS7VarPayloadDataItem, ok := typ.(*S7VarPayloadDataItem); ok {
            return *sS7VarPayloadDataItem
        }
        return S7VarPayloadDataItem{}
    }
    return castFunc(structType)
}

func (m S7VarPayloadDataItem) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Enum Field (returnCode)
    lengthInBits += 8

    // Enum Field (transportSize)
    lengthInBits += 8

    // Implicit Field (dataLength)
    lengthInBits += 16

    // Array field
    if len(m.Data) > 0 {
        lengthInBits += 8 * uint16(len(m.Data))
    }

    // Padding Field (padding)
    _timesPadding := uint8(utils.InlineIf(false, uint16(uint8(0)), uint16(uint8(uint8(len(m.Data))) % uint8(uint8(2)))))
    for ;_timesPadding > 0; _timesPadding-- {
        lengthInBits += 8
    }

    return lengthInBits
}

func (m S7VarPayloadDataItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7VarPayloadDataItemParse(io *utils.ReadBuffer, lastItem bool) (spi.Message, error) {

    // Enum field (returnCode)
    returnCode, _returnCodeErr := DataTransportErrorCodeParse(io)
    if _returnCodeErr != nil {
        return nil, errors.New("Error parsing 'returnCode' field " + _returnCodeErr.Error())
    }

    // Enum field (transportSize)
    transportSize, _transportSizeErr := DataTransportSizeParse(io)
    if _transportSizeErr != nil {
        return nil, errors.New("Error parsing 'transportSize' field " + _transportSizeErr.Error())
    }

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    dataLength, _dataLengthErr := io.ReadUint16(16)
    if _dataLengthErr != nil {
        return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Array field (data)
    // Count array
    data := make([]int8, utils.InlineIf(transportSize.SizeInBits(), uint16(math.Ceil(float64(dataLength) / float64(float64(8.0)))), uint16(dataLength)))
    for curItem := uint16(0); curItem < uint16(utils.InlineIf(transportSize.SizeInBits(), uint16(math.Ceil(float64(dataLength) / float64(float64(8.0)))), uint16(dataLength))); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'data' field " + _err.Error())
        }
        data[curItem] = _item
    }

    // Padding Field (padding)
    {
        _timesPadding := (utils.InlineIf(lastItem, uint16(uint8(0)), uint16(uint8(uint8(len(data))) % uint8(uint8(2)))))
        for ;(io.HasMore(8)) && (_timesPadding > 0);_timesPadding-- {
            // Just read the padding data and ignore it
            _, _err := io.ReadUint8(8)
            if _err != nil {
                return nil, errors.New("Error parsing 'padding' field " + _err.Error())
            }
        }
    }

    // Create the instance
    return NewS7VarPayloadDataItem(returnCode, transportSize, data), nil
}

func (m S7VarPayloadDataItem) Serialize(io utils.WriteBuffer, lastItem bool) error {

    // Enum field (returnCode)
    returnCode := CastDataTransportErrorCode(m.ReturnCode)
    _returnCodeErr := returnCode.Serialize(io)
    if _returnCodeErr != nil {
        return errors.New("Error serializing 'returnCode' field " + _returnCodeErr.Error())
    }

    // Enum field (transportSize)
    transportSize := CastDataTransportSize(m.TransportSize)
    _transportSizeErr := transportSize.Serialize(io)
    if _transportSizeErr != nil {
        return errors.New("Error serializing 'transportSize' field " + _transportSizeErr.Error())
    }

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    dataLength := uint16(uint16(uint16(len(m.Data))) * uint16(uint16(utils.InlineIf(bool(bool((m.TransportSize) == (DataTransportSize_BIT))), uint16(uint16(1)), uint16(uint16(utils.InlineIf(transportSize.SizeInBits(), uint16(uint16(8)), uint16(uint16(1)))))))))
    _dataLengthErr := io.WriteUint16(16, (dataLength))
    if _dataLengthErr != nil {
        return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
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

    // Padding Field (padding)
    {
        _timesPadding := uint8(utils.InlineIf(lastItem, uint16(uint8(0)), uint16(uint8(uint8(len(m.Data))) % uint8(uint8(2)))))
        for ;_timesPadding > 0; _timesPadding-- {
            _paddingValue := uint8(uint8(0))
            _paddingErr := io.WriteUint8(8, (_paddingValue))
            if _paddingErr != nil {
                return errors.New("Error serializing 'padding' field " + _paddingErr.Error())
            }
        }
    }

    return nil
}

func (m *S7VarPayloadDataItem) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "returnCode":
                var data *DataTransportErrorCode
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ReturnCode = data
            case "transportSize":
                var data *DataTransportSize
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TransportSize = data
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

func (m S7VarPayloadDataItem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7VarPayloadDataItem"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ReturnCode, xml.StartElement{Name: xml.Name{Local: "returnCode"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.TransportSize, xml.StartElement{Name: xml.Name{Local: "transportSize"}}); err != nil {
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

