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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "io"
)

// The data-structure of this message
type LDataFrameData struct {
    SourceAddress *KnxAddress
    DestinationAddress []int8
    GroupAddress bool
    HopCount uint8
    DataLength uint8
    Control bool
    Numbered bool
    Counter uint8
    ControlType *ControlType
    Apci *APCI
    ExtendedApci *ExtendedAPCI
    DataFirstByte *int8
    Data []int8
    Parent *LDataFrame
    ILDataFrameData
}

// The corresponding interface
type ILDataFrameData interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *LDataFrameData) NotAckFrame() bool {
    return true
}

func (m *LDataFrameData) ExtendedFrame() bool {
    return false
}

func (m *LDataFrameData) Polling() bool {
    return false
}


func (m *LDataFrameData) InitializeParent(parent *LDataFrame, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) {
    m.Parent.Repeated = repeated
    m.Parent.Priority = priority
    m.Parent.AcknowledgeRequested = acknowledgeRequested
    m.Parent.ErrorFlag = errorFlag
}

func NewLDataFrameData(sourceAddress *KnxAddress, destinationAddress []int8, groupAddress bool, hopCount uint8, dataLength uint8, control bool, numbered bool, counter uint8, controlType *ControlType, apci *APCI, extendedApci *ExtendedAPCI, dataFirstByte *int8, data []int8, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) *LDataFrame {
    child := &LDataFrameData{
        SourceAddress: sourceAddress,
        DestinationAddress: destinationAddress,
        GroupAddress: groupAddress,
        HopCount: hopCount,
        DataLength: dataLength,
        Control: control,
        Numbered: numbered,
        Counter: counter,
        ControlType: controlType,
        Apci: apci,
        ExtendedApci: extendedApci,
        DataFirstByte: dataFirstByte,
        Data: data,
        Parent: NewLDataFrame(repeated, priority, acknowledgeRequested, errorFlag),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastLDataFrameData(structType interface{}) *LDataFrameData {
    castFunc := func(typ interface{}) *LDataFrameData {
        if casted, ok := typ.(LDataFrameData); ok {
            return &casted
        }
        if casted, ok := typ.(*LDataFrameData); ok {
            return casted
        }
        if casted, ok := typ.(LDataFrame); ok {
            return CastLDataFrameData(casted.Child)
        }
        if casted, ok := typ.(*LDataFrame); ok {
            return CastLDataFrameData(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *LDataFrameData) GetTypeName() string {
    return "LDataFrameData"
}

func (m *LDataFrameData) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (sourceAddress)
    lengthInBits += m.SourceAddress.LengthInBits()

    // Array field
    if len(m.DestinationAddress) > 0 {
        lengthInBits += 8 * uint16(len(m.DestinationAddress))
    }

    // Simple field (groupAddress)
    lengthInBits += 1

    // Simple field (hopCount)
    lengthInBits += 3

    // Simple field (dataLength)
    lengthInBits += 4

    // Simple field (control)
    lengthInBits += 1

    // Simple field (numbered)
    lengthInBits += 1

    // Simple field (counter)
    lengthInBits += 4

    // Optional Field (controlType)
    if m.ControlType != nil {
        lengthInBits += 2
    }

    // Optional Field (apci)
    if m.Apci != nil {
        lengthInBits += 4
    }

    // Optional Field (extendedApci)
    if m.ExtendedApci != nil {
        lengthInBits += 6
    }

    // Optional Field (dataFirstByte)
    if m.DataFirstByte != nil {
        lengthInBits += 6
    }

    // Array field
    if len(m.Data) > 0 {
        lengthInBits += 8 * uint16(len(m.Data))
    }

    return lengthInBits
}

func (m *LDataFrameData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func LDataFrameDataParse(io *utils.ReadBuffer) (*LDataFrame, error) {

    // Simple Field (sourceAddress)
    sourceAddress, _sourceAddressErr := KnxAddressParse(io)
    if _sourceAddressErr != nil {
        return nil, errors.New("Error parsing 'sourceAddress' field " + _sourceAddressErr.Error())
    }

    // Array field (destinationAddress)
    // Count array
    destinationAddress := make([]int8, uint16(2))
    for curItem := uint16(0); curItem < uint16(uint16(2)); curItem++ {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'destinationAddress' field " + _err.Error())
        }
        destinationAddress[curItem] = _item
    }

    // Simple Field (groupAddress)
    groupAddress, _groupAddressErr := io.ReadBit()
    if _groupAddressErr != nil {
        return nil, errors.New("Error parsing 'groupAddress' field " + _groupAddressErr.Error())
    }

    // Simple Field (hopCount)
    hopCount, _hopCountErr := io.ReadUint8(3)
    if _hopCountErr != nil {
        return nil, errors.New("Error parsing 'hopCount' field " + _hopCountErr.Error())
    }

    // Simple Field (dataLength)
    dataLength, _dataLengthErr := io.ReadUint8(4)
    if _dataLengthErr != nil {
        return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Simple Field (control)
    control, _controlErr := io.ReadBit()
    if _controlErr != nil {
        return nil, errors.New("Error parsing 'control' field " + _controlErr.Error())
    }

    // Simple Field (numbered)
    numbered, _numberedErr := io.ReadBit()
    if _numberedErr != nil {
        return nil, errors.New("Error parsing 'numbered' field " + _numberedErr.Error())
    }

    // Simple Field (counter)
    counter, _counterErr := io.ReadUint8(4)
    if _counterErr != nil {
        return nil, errors.New("Error parsing 'counter' field " + _counterErr.Error())
    }

    // Optional Field (controlType) (Can be skipped, if a given expression evaluates to false)
    var controlType *ControlType = nil
    if control {
        _val, _err := ControlTypeParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'controlType' field " + _err.Error())
        }
        controlType = &_val
    }

    // Optional Field (apci) (Can be skipped, if a given expression evaluates to false)
    var apci *APCI = nil
    if !(control) {
        _val, _err := APCIParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'apci' field " + _err.Error())
        }
        apci = &_val
    }

    // Optional Field (extendedApci) (Can be skipped, if a given expression evaluates to false)
    var extendedApci *ExtendedAPCI = nil
    if bool(!(control)) && bool(bool(((*apci)) == (APCI_OTHER_PDU))) {
        _val, _err := ExtendedAPCIParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'extendedApci' field " + _err.Error())
        }
        extendedApci = &_val
    }

    // Optional Field (dataFirstByte) (Can be skipped, if a given expression evaluates to false)
    var dataFirstByte *int8 = nil
    if bool(!(control)) && bool(bool(((*apci)) != (APCI_OTHER_PDU))) {
        _val, _err := io.ReadInt8(6)
        if _err != nil {
            return nil, errors.New("Error parsing 'dataFirstByte' field " + _err.Error())
        }
        dataFirstByte = &_val
    }

    // Array field (data)
    // Count array
    data := make([]int8, utils.InlineIf(bool(bool((dataLength) < ((1)))), uint16(uint16(0)), uint16(uint16(dataLength) - uint16(uint16(1)))))
    for curItem := uint16(0); curItem < uint16(utils.InlineIf(bool(bool((dataLength) < ((1)))), uint16(uint16(0)), uint16(uint16(dataLength) - uint16(uint16(1))))); curItem++ {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'data' field " + _err.Error())
        }
        data[curItem] = _item
    }

    // Create a partially initialized instance
    _child := &LDataFrameData{
        SourceAddress: sourceAddress,
        DestinationAddress: destinationAddress,
        GroupAddress: groupAddress,
        HopCount: hopCount,
        DataLength: dataLength,
        Control: control,
        Numbered: numbered,
        Counter: counter,
        ControlType: controlType,
        Apci: apci,
        ExtendedApci: extendedApci,
        DataFirstByte: dataFirstByte,
        Data: data,
        Parent: &LDataFrame{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *LDataFrameData) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (sourceAddress)
    _sourceAddressErr := m.SourceAddress.Serialize(io)
    if _sourceAddressErr != nil {
        return errors.New("Error serializing 'sourceAddress' field " + _sourceAddressErr.Error())
    }

    // Array Field (destinationAddress)
    if m.DestinationAddress != nil {
        for _, _element := range m.DestinationAddress {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'destinationAddress' field " + _elementErr.Error())
            }
        }
    }

    // Simple Field (groupAddress)
    groupAddress := bool(m.GroupAddress)
    _groupAddressErr := io.WriteBit((groupAddress))
    if _groupAddressErr != nil {
        return errors.New("Error serializing 'groupAddress' field " + _groupAddressErr.Error())
    }

    // Simple Field (hopCount)
    hopCount := uint8(m.HopCount)
    _hopCountErr := io.WriteUint8(3, (hopCount))
    if _hopCountErr != nil {
        return errors.New("Error serializing 'hopCount' field " + _hopCountErr.Error())
    }

    // Simple Field (dataLength)
    dataLength := uint8(m.DataLength)
    _dataLengthErr := io.WriteUint8(4, (dataLength))
    if _dataLengthErr != nil {
        return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Simple Field (control)
    control := bool(m.Control)
    _controlErr := io.WriteBit((control))
    if _controlErr != nil {
        return errors.New("Error serializing 'control' field " + _controlErr.Error())
    }

    // Simple Field (numbered)
    numbered := bool(m.Numbered)
    _numberedErr := io.WriteBit((numbered))
    if _numberedErr != nil {
        return errors.New("Error serializing 'numbered' field " + _numberedErr.Error())
    }

    // Simple Field (counter)
    counter := uint8(m.Counter)
    _counterErr := io.WriteUint8(4, (counter))
    if _counterErr != nil {
        return errors.New("Error serializing 'counter' field " + _counterErr.Error())
    }

    // Optional Field (controlType) (Can be skipped, if the value is null)
    var controlType *ControlType = nil
    if m.ControlType != nil {
        controlType = m.ControlType
        _controlTypeErr := controlType.Serialize(io)
        if _controlTypeErr != nil {
            return errors.New("Error serializing 'controlType' field " + _controlTypeErr.Error())
        }
    }

    // Optional Field (apci) (Can be skipped, if the value is null)
    var apci *APCI = nil
    if m.Apci != nil {
        apci = m.Apci
        _apciErr := apci.Serialize(io)
        if _apciErr != nil {
            return errors.New("Error serializing 'apci' field " + _apciErr.Error())
        }
    }

    // Optional Field (extendedApci) (Can be skipped, if the value is null)
    var extendedApci *ExtendedAPCI = nil
    if m.ExtendedApci != nil {
        extendedApci = m.ExtendedApci
        _extendedApciErr := extendedApci.Serialize(io)
        if _extendedApciErr != nil {
            return errors.New("Error serializing 'extendedApci' field " + _extendedApciErr.Error())
        }
    }

    // Optional Field (dataFirstByte) (Can be skipped, if the value is null)
    var dataFirstByte *int8 = nil
    if m.DataFirstByte != nil {
        dataFirstByte = m.DataFirstByte
        _dataFirstByteErr := io.WriteInt8(6, *(dataFirstByte))
        if _dataFirstByteErr != nil {
            return errors.New("Error serializing 'dataFirstByte' field " + _dataFirstByteErr.Error())
        }
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
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *LDataFrameData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "sourceAddress":
                var data *KnxAddress
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.SourceAddress = data
            case "destinationAddress":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.DestinationAddress = utils.ByteArrayToInt8Array(_decoded[0:_len])
            case "groupAddress":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.GroupAddress = data
            case "hopCount":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.HopCount = data
            case "dataLength":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DataLength = data
            case "control":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Control = data
            case "numbered":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Numbered = data
            case "counter":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Counter = data
            case "controlType":
                var data *ControlType
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ControlType = data
            case "apci":
                var data *APCI
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.Apci = data
            case "extendedApci":
                var data *ExtendedAPCI
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ExtendedApci = data
            case "dataFirstByte":
                var data *int8
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.DataFirstByte = data
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
                m.Data = utils.ByteArrayToInt8Array(_decoded[0:_len])
            }
        }
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *LDataFrameData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.SourceAddress, xml.StartElement{Name: xml.Name{Local: "sourceAddress"}}); err != nil {
        return err
    }
    _encodedDestinationAddress := make([]byte, base64.StdEncoding.EncodedLen(len(m.DestinationAddress)))
    base64.StdEncoding.Encode(_encodedDestinationAddress, utils.Int8ArrayToByteArray(m.DestinationAddress))
    if err := e.EncodeElement(_encodedDestinationAddress, xml.StartElement{Name: xml.Name{Local: "destinationAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.GroupAddress, xml.StartElement{Name: xml.Name{Local: "groupAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HopCount, xml.StartElement{Name: xml.Name{Local: "hopCount"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DataLength, xml.StartElement{Name: xml.Name{Local: "dataLength"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Control, xml.StartElement{Name: xml.Name{Local: "control"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Numbered, xml.StartElement{Name: xml.Name{Local: "numbered"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Counter, xml.StartElement{Name: xml.Name{Local: "counter"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ControlType, xml.StartElement{Name: xml.Name{Local: "controlType"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Apci, xml.StartElement{Name: xml.Name{Local: "apci"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ExtendedApci, xml.StartElement{Name: xml.Name{Local: "extendedApci"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DataFirstByte, xml.StartElement{Name: xml.Name{Local: "dataFirstByte"}}); err != nil {
        return err
    }
    _encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(m.Data)))
    base64.StdEncoding.Encode(_encodedData, utils.Int8ArrayToByteArray(m.Data))
    if err := e.EncodeElement(_encodedData, xml.StartElement{Name: xml.Name{Local: "data"}}); err != nil {
        return err
    }
    return nil
}

