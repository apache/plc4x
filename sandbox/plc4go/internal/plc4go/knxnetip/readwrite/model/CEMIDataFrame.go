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
    "reflect"
)

// The data-structure of this message
type CEMIDataFrame struct {
    StandardFrame bool
    Polling bool
    NotRepeated bool
    NotAckFrame bool
    Priority ICEMIPriority
    AcknowledgeRequested bool
    ErrorFlag bool
    GroupDestinationAddress bool
    HopCount uint8
    ExtendedFrameFormat uint8
    SourceAddress IKNXAddress
    DestinationAddress []int8
    DataLength uint8
    Tcpi ITPCI
    Counter uint8
    Apci IAPCI
    DataFirstByte int8
    Data []int8

}

// The corresponding interface
type ICEMIDataFrame interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewCEMIDataFrame(standardFrame bool, polling bool, notRepeated bool, notAckFrame bool, priority ICEMIPriority, acknowledgeRequested bool, errorFlag bool, groupDestinationAddress bool, hopCount uint8, extendedFrameFormat uint8, sourceAddress IKNXAddress, destinationAddress []int8, dataLength uint8, tcpi ITPCI, counter uint8, apci IAPCI, dataFirstByte int8, data []int8) spi.Message {
    return &CEMIDataFrame{StandardFrame: standardFrame, Polling: polling, NotRepeated: notRepeated, NotAckFrame: notAckFrame, Priority: priority, AcknowledgeRequested: acknowledgeRequested, ErrorFlag: errorFlag, GroupDestinationAddress: groupDestinationAddress, HopCount: hopCount, ExtendedFrameFormat: extendedFrameFormat, SourceAddress: sourceAddress, DestinationAddress: destinationAddress, DataLength: dataLength, Tcpi: tcpi, Counter: counter, Apci: apci, DataFirstByte: dataFirstByte, Data: data}
}

func CastICEMIDataFrame(structType interface{}) ICEMIDataFrame {
    castFunc := func(typ interface{}) ICEMIDataFrame {
        if iCEMIDataFrame, ok := typ.(ICEMIDataFrame); ok {
            return iCEMIDataFrame
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIDataFrame(structType interface{}) CEMIDataFrame {
    castFunc := func(typ interface{}) CEMIDataFrame {
        if sCEMIDataFrame, ok := typ.(CEMIDataFrame); ok {
            return sCEMIDataFrame
        }
        if sCEMIDataFrame, ok := typ.(*CEMIDataFrame); ok {
            return *sCEMIDataFrame
        }
        return CEMIDataFrame{}
    }
    return castFunc(structType)
}

func (m CEMIDataFrame) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Simple field (standardFrame)
    lengthInBits += 1

    // Simple field (polling)
    lengthInBits += 1

    // Simple field (notRepeated)
    lengthInBits += 1

    // Simple field (notAckFrame)
    lengthInBits += 1

    // Enum Field (priority)
    lengthInBits += 2

    // Simple field (acknowledgeRequested)
    lengthInBits += 1

    // Simple field (errorFlag)
    lengthInBits += 1

    // Simple field (groupDestinationAddress)
    lengthInBits += 1

    // Simple field (hopCount)
    lengthInBits += 3

    // Simple field (extendedFrameFormat)
    lengthInBits += 4

    // Simple field (sourceAddress)
    lengthInBits += m.SourceAddress.LengthInBits()

    // Array field
    if len(m.DestinationAddress) > 0 {
        lengthInBits += 8 * uint16(len(m.DestinationAddress))
    }

    // Simple field (dataLength)
    lengthInBits += 8

    // Enum Field (tcpi)
    lengthInBits += 2

    // Simple field (counter)
    lengthInBits += 4

    // Enum Field (apci)
    lengthInBits += 4

    // Simple field (dataFirstByte)
    lengthInBits += 6

    // Array field
    if len(m.Data) > 0 {
        lengthInBits += 8 * uint16(len(m.Data))
    }

    return lengthInBits
}

func (m CEMIDataFrame) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIDataFrameParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Simple Field (standardFrame)
    standardFrame, _standardFrameErr := io.ReadBit()
    if _standardFrameErr != nil {
        return nil, errors.New("Error parsing 'standardFrame' field " + _standardFrameErr.Error())
    }

    // Simple Field (polling)
    polling, _pollingErr := io.ReadBit()
    if _pollingErr != nil {
        return nil, errors.New("Error parsing 'polling' field " + _pollingErr.Error())
    }

    // Simple Field (notRepeated)
    notRepeated, _notRepeatedErr := io.ReadBit()
    if _notRepeatedErr != nil {
        return nil, errors.New("Error parsing 'notRepeated' field " + _notRepeatedErr.Error())
    }

    // Simple Field (notAckFrame)
    notAckFrame, _notAckFrameErr := io.ReadBit()
    if _notAckFrameErr != nil {
        return nil, errors.New("Error parsing 'notAckFrame' field " + _notAckFrameErr.Error())
    }

    // Enum field (priority)
    priority, _priorityErr := CEMIPriorityParse(io)
    if _priorityErr != nil {
        return nil, errors.New("Error parsing 'priority' field " + _priorityErr.Error())
    }

    // Simple Field (acknowledgeRequested)
    acknowledgeRequested, _acknowledgeRequestedErr := io.ReadBit()
    if _acknowledgeRequestedErr != nil {
        return nil, errors.New("Error parsing 'acknowledgeRequested' field " + _acknowledgeRequestedErr.Error())
    }

    // Simple Field (errorFlag)
    errorFlag, _errorFlagErr := io.ReadBit()
    if _errorFlagErr != nil {
        return nil, errors.New("Error parsing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Simple Field (groupDestinationAddress)
    groupDestinationAddress, _groupDestinationAddressErr := io.ReadBit()
    if _groupDestinationAddressErr != nil {
        return nil, errors.New("Error parsing 'groupDestinationAddress' field " + _groupDestinationAddressErr.Error())
    }

    // Simple Field (hopCount)
    hopCount, _hopCountErr := io.ReadUint8(3)
    if _hopCountErr != nil {
        return nil, errors.New("Error parsing 'hopCount' field " + _hopCountErr.Error())
    }

    // Simple Field (extendedFrameFormat)
    extendedFrameFormat, _extendedFrameFormatErr := io.ReadUint8(4)
    if _extendedFrameFormatErr != nil {
        return nil, errors.New("Error parsing 'extendedFrameFormat' field " + _extendedFrameFormatErr.Error())
    }

    // Simple Field (sourceAddress)
    _sourceAddressMessage, _err := KNXAddressParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'sourceAddress'. " + _err.Error())
    }
    var sourceAddress IKNXAddress
    sourceAddress, _sourceAddressOk := _sourceAddressMessage.(IKNXAddress)
    if !_sourceAddressOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_sourceAddressMessage).Name() + " to IKNXAddress")
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

    // Simple Field (dataLength)
    dataLength, _dataLengthErr := io.ReadUint8(8)
    if _dataLengthErr != nil {
        return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Enum field (tcpi)
    tcpi, _tcpiErr := TPCIParse(io)
    if _tcpiErr != nil {
        return nil, errors.New("Error parsing 'tcpi' field " + _tcpiErr.Error())
    }

    // Simple Field (counter)
    counter, _counterErr := io.ReadUint8(4)
    if _counterErr != nil {
        return nil, errors.New("Error parsing 'counter' field " + _counterErr.Error())
    }

    // Enum field (apci)
    apci, _apciErr := APCIParse(io)
    if _apciErr != nil {
        return nil, errors.New("Error parsing 'apci' field " + _apciErr.Error())
    }

    // Simple Field (dataFirstByte)
    dataFirstByte, _dataFirstByteErr := io.ReadInt8(6)
    if _dataFirstByteErr != nil {
        return nil, errors.New("Error parsing 'dataFirstByte' field " + _dataFirstByteErr.Error())
    }

    // Array field (data)
    // Count array
    data := make([]int8, uint16(dataLength) - uint16(uint16(1)))
    for curItem := uint16(0); curItem < uint16(uint16(dataLength) - uint16(uint16(1))); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'data' field " + _err.Error())
        }
        data[curItem] = _item
    }

    // Create the instance
    return NewCEMIDataFrame(standardFrame, polling, notRepeated, notAckFrame, priority, acknowledgeRequested, errorFlag, groupDestinationAddress, hopCount, extendedFrameFormat, sourceAddress, destinationAddress, dataLength, tcpi, counter, apci, dataFirstByte, data), nil
}

func (m CEMIDataFrame) Serialize(io utils.WriteBuffer) error {

    // Simple Field (standardFrame)
    standardFrame := bool(m.StandardFrame)
    _standardFrameErr := io.WriteBit((bool) (standardFrame))
    if _standardFrameErr != nil {
        return errors.New("Error serializing 'standardFrame' field " + _standardFrameErr.Error())
    }

    // Simple Field (polling)
    polling := bool(m.Polling)
    _pollingErr := io.WriteBit((bool) (polling))
    if _pollingErr != nil {
        return errors.New("Error serializing 'polling' field " + _pollingErr.Error())
    }

    // Simple Field (notRepeated)
    notRepeated := bool(m.NotRepeated)
    _notRepeatedErr := io.WriteBit((bool) (notRepeated))
    if _notRepeatedErr != nil {
        return errors.New("Error serializing 'notRepeated' field " + _notRepeatedErr.Error())
    }

    // Simple Field (notAckFrame)
    notAckFrame := bool(m.NotAckFrame)
    _notAckFrameErr := io.WriteBit((bool) (notAckFrame))
    if _notAckFrameErr != nil {
        return errors.New("Error serializing 'notAckFrame' field " + _notAckFrameErr.Error())
    }

    // Enum field (priority)
    priority := CastCEMIPriority(m.Priority)
    _priorityErr := priority.Serialize(io)
    if _priorityErr != nil {
        return errors.New("Error serializing 'priority' field " + _priorityErr.Error())
    }

    // Simple Field (acknowledgeRequested)
    acknowledgeRequested := bool(m.AcknowledgeRequested)
    _acknowledgeRequestedErr := io.WriteBit((bool) (acknowledgeRequested))
    if _acknowledgeRequestedErr != nil {
        return errors.New("Error serializing 'acknowledgeRequested' field " + _acknowledgeRequestedErr.Error())
    }

    // Simple Field (errorFlag)
    errorFlag := bool(m.ErrorFlag)
    _errorFlagErr := io.WriteBit((bool) (errorFlag))
    if _errorFlagErr != nil {
        return errors.New("Error serializing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Simple Field (groupDestinationAddress)
    groupDestinationAddress := bool(m.GroupDestinationAddress)
    _groupDestinationAddressErr := io.WriteBit((bool) (groupDestinationAddress))
    if _groupDestinationAddressErr != nil {
        return errors.New("Error serializing 'groupDestinationAddress' field " + _groupDestinationAddressErr.Error())
    }

    // Simple Field (hopCount)
    hopCount := uint8(m.HopCount)
    _hopCountErr := io.WriteUint8(3, (hopCount))
    if _hopCountErr != nil {
        return errors.New("Error serializing 'hopCount' field " + _hopCountErr.Error())
    }

    // Simple Field (extendedFrameFormat)
    extendedFrameFormat := uint8(m.ExtendedFrameFormat)
    _extendedFrameFormatErr := io.WriteUint8(4, (extendedFrameFormat))
    if _extendedFrameFormatErr != nil {
        return errors.New("Error serializing 'extendedFrameFormat' field " + _extendedFrameFormatErr.Error())
    }

    // Simple Field (sourceAddress)
    sourceAddress := CastIKNXAddress(m.SourceAddress)
    _sourceAddressErr := sourceAddress.Serialize(io)
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

    // Simple Field (dataLength)
    dataLength := uint8(m.DataLength)
    _dataLengthErr := io.WriteUint8(8, (dataLength))
    if _dataLengthErr != nil {
        return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Enum field (tcpi)
    tcpi := CastTPCI(m.Tcpi)
    _tcpiErr := tcpi.Serialize(io)
    if _tcpiErr != nil {
        return errors.New("Error serializing 'tcpi' field " + _tcpiErr.Error())
    }

    // Simple Field (counter)
    counter := uint8(m.Counter)
    _counterErr := io.WriteUint8(4, (counter))
    if _counterErr != nil {
        return errors.New("Error serializing 'counter' field " + _counterErr.Error())
    }

    // Enum field (apci)
    apci := CastAPCI(m.Apci)
    _apciErr := apci.Serialize(io)
    if _apciErr != nil {
        return errors.New("Error serializing 'apci' field " + _apciErr.Error())
    }

    // Simple Field (dataFirstByte)
    dataFirstByte := int8(m.DataFirstByte)
    _dataFirstByteErr := io.WriteInt8(6, (dataFirstByte))
    if _dataFirstByteErr != nil {
        return errors.New("Error serializing 'dataFirstByte' field " + _dataFirstByteErr.Error())
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

func (m *CEMIDataFrame) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "standardFrame":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.StandardFrame = data
            case "polling":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Polling = data
            case "notRepeated":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.NotRepeated = data
            case "notAckFrame":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.NotAckFrame = data
            case "priority":
                var data *CEMIPriority
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Priority = data
            case "acknowledgeRequested":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.AcknowledgeRequested = data
            case "errorFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ErrorFlag = data
            case "groupDestinationAddress":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.GroupDestinationAddress = data
            case "hopCount":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.HopCount = data
            case "extendedFrameFormat":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ExtendedFrameFormat = data
            case "sourceAddress":
                var data *KNXAddress
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SourceAddress = CastIKNXAddress(data)
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
                m.DestinationAddress = utils.ByteToInt8(_decoded[0:_len])
            case "dataLength":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DataLength = data
            case "tcpi":
                var data *TPCI
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Tcpi = data
            case "counter":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Counter = data
            case "apci":
                var data *APCI
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Apci = data
            case "dataFirstByte":
                var data int8
                if err := d.DecodeElement(&data, &tok); err != nil {
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
                m.Data = utils.ByteToInt8(_decoded[0:_len])
            }
        }
    }
}

func (m CEMIDataFrame) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataFrame"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.StandardFrame, xml.StartElement{Name: xml.Name{Local: "standardFrame"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Polling, xml.StartElement{Name: xml.Name{Local: "polling"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.NotRepeated, xml.StartElement{Name: xml.Name{Local: "notRepeated"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.NotAckFrame, xml.StartElement{Name: xml.Name{Local: "notAckFrame"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Priority, xml.StartElement{Name: xml.Name{Local: "priority"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.AcknowledgeRequested, xml.StartElement{Name: xml.Name{Local: "acknowledgeRequested"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ErrorFlag, xml.StartElement{Name: xml.Name{Local: "errorFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.GroupDestinationAddress, xml.StartElement{Name: xml.Name{Local: "groupDestinationAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HopCount, xml.StartElement{Name: xml.Name{Local: "hopCount"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ExtendedFrameFormat, xml.StartElement{Name: xml.Name{Local: "extendedFrameFormat"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SourceAddress, xml.StartElement{Name: xml.Name{Local: "sourceAddress"}}); err != nil {
        return err
    }
    _encodedDestinationAddress := make([]byte, base64.StdEncoding.EncodedLen(len(m.DestinationAddress)))
    base64.StdEncoding.Encode(_encodedDestinationAddress, utils.Int8ToByte(m.DestinationAddress))
    if err := e.EncodeElement(_encodedDestinationAddress, xml.StartElement{Name: xml.Name{Local: "destinationAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DataLength, xml.StartElement{Name: xml.Name{Local: "dataLength"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Tcpi, xml.StartElement{Name: xml.Name{Local: "tcpi"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Counter, xml.StartElement{Name: xml.Name{Local: "counter"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Apci, xml.StartElement{Name: xml.Name{Local: "apci"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DataFirstByte, xml.StartElement{Name: xml.Name{Local: "dataFirstByte"}}); err != nil {
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

