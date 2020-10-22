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
    "reflect"
)

// The data-structure of this message
type CEMIFrameData struct {
    SourceAddress IKNXAddress
    DestinationAddress []int8
    GroupAddress bool
    HopCount uint8
    DataLength uint8
    Tcpi ITPCI
    Counter uint8
    Apci IAPCI
    DataFirstByte int8
    Data []int8
    Crc uint8
    CEMIFrame
}

// The corresponding interface
type ICEMIFrameData interface {
    ICEMIFrame
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIFrameData) NotAckFrame() bool {
    return true
}

func (m CEMIFrameData) StandardFrame() bool {
    return true
}

func (m CEMIFrameData) Polling() bool {
    return false
}

func (m CEMIFrameData) initialize(repeated bool, priority ICEMIPriority, acknowledgeRequested bool, errorFlag bool) spi.Message {
    m.Repeated = repeated
    m.Priority = priority
    m.AcknowledgeRequested = acknowledgeRequested
    m.ErrorFlag = errorFlag
    return m
}

func NewCEMIFrameData(sourceAddress IKNXAddress, destinationAddress []int8, groupAddress bool, hopCount uint8, dataLength uint8, tcpi ITPCI, counter uint8, apci IAPCI, dataFirstByte int8, data []int8, crc uint8) CEMIFrameInitializer {
    return &CEMIFrameData{SourceAddress: sourceAddress, DestinationAddress: destinationAddress, GroupAddress: groupAddress, HopCount: hopCount, DataLength: dataLength, Tcpi: tcpi, Counter: counter, Apci: apci, DataFirstByte: dataFirstByte, Data: data, Crc: crc}
}

func CastICEMIFrameData(structType interface{}) ICEMIFrameData {
    castFunc := func(typ interface{}) ICEMIFrameData {
        if iCEMIFrameData, ok := typ.(ICEMIFrameData); ok {
            return iCEMIFrameData
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIFrameData(structType interface{}) CEMIFrameData {
    castFunc := func(typ interface{}) CEMIFrameData {
        if sCEMIFrameData, ok := typ.(CEMIFrameData); ok {
            return sCEMIFrameData
        }
        if sCEMIFrameData, ok := typ.(*CEMIFrameData); ok {
            return *sCEMIFrameData
        }
        return CEMIFrameData{}
    }
    return castFunc(structType)
}

func (m CEMIFrameData) LengthInBits() uint16 {
    var lengthInBits uint16 = m.CEMIFrame.LengthInBits()

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

    // Simple field (crc)
    lengthInBits += 8

    return lengthInBits
}

func (m CEMIFrameData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIFrameDataParse(io *spi.ReadBuffer) (CEMIFrameInitializer, error) {

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

    // Simple Field (crc)
    crc, _crcErr := io.ReadUint8(8)
    if _crcErr != nil {
        return nil, errors.New("Error parsing 'crc' field " + _crcErr.Error())
    }

    // Create the instance
    return NewCEMIFrameData(sourceAddress, destinationAddress, groupAddress, hopCount, dataLength, tcpi, counter, apci, dataFirstByte, data, crc), nil
}

func (m CEMIFrameData) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

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

    // Simple Field (groupAddress)
    groupAddress := bool(m.GroupAddress)
    _groupAddressErr := io.WriteBit((bool) (groupAddress))
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

    // Simple Field (crc)
    crc := uint8(m.Crc)
    _crcErr := io.WriteUint8(8, (crc))
    if _crcErr != nil {
        return errors.New("Error serializing 'crc' field " + _crcErr.Error())
    }

        return nil
    }
    return CEMIFrameSerialize(io, m.CEMIFrame, CastICEMIFrame(m), ser)
}
