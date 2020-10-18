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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
)

// The data-structure of this message
type CEMIFrameDataExt struct {
	GroupAddress        bool
	HopCount            uint8
	ExtendedFrameFormat uint8
	SourceAddress       IKNXAddress
	DestinationAddress  []int8
	DataLength          uint8
	Tcpi                ITPCI
	Counter             uint8
	Apci                IAPCI
	DataFirstByte       int8
	Data                []int8
	Crc                 uint8
	CEMIFrame
}

// The corresponding interface
type ICEMIFrameDataExt interface {
	ICEMIFrame
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIFrameDataExt) NotAckFrame() bool {
	return true
}

func (m CEMIFrameDataExt) StandardFrame() bool {
	return false
}

func (m CEMIFrameDataExt) Polling() bool {
	return false
}

func (m CEMIFrameDataExt) initialize(repeated bool, priority ICEMIPriority, acknowledgeRequested bool, errorFlag bool) spi.Message {
	m.Repeated = repeated
	m.Priority = priority
	m.AcknowledgeRequested = acknowledgeRequested
	m.ErrorFlag = errorFlag
	return m
}

func NewCEMIFrameDataExt(groupAddress bool, hopCount uint8, extendedFrameFormat uint8, sourceAddress IKNXAddress, destinationAddress []int8, dataLength uint8, tcpi ITPCI, counter uint8, apci IAPCI, dataFirstByte int8, data []int8, crc uint8) CEMIFrameInitializer {
	return &CEMIFrameDataExt{GroupAddress: groupAddress, HopCount: hopCount, ExtendedFrameFormat: extendedFrameFormat, SourceAddress: sourceAddress, DestinationAddress: destinationAddress, DataLength: dataLength, Tcpi: tcpi, Counter: counter, Apci: apci, DataFirstByte: dataFirstByte, Data: data, Crc: crc}
}

func CastICEMIFrameDataExt(structType interface{}) ICEMIFrameDataExt {
	castFunc := func(typ interface{}) ICEMIFrameDataExt {
		if iCEMIFrameDataExt, ok := typ.(ICEMIFrameDataExt); ok {
			return iCEMIFrameDataExt
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIFrameDataExt(structType interface{}) CEMIFrameDataExt {
	castFunc := func(typ interface{}) CEMIFrameDataExt {
		if sCEMIFrameDataExt, ok := typ.(CEMIFrameDataExt); ok {
			return sCEMIFrameDataExt
		}
		return CEMIFrameDataExt{}
	}
	return castFunc(structType)
}

func (m CEMIFrameDataExt) LengthInBits() uint16 {
	var lengthInBits = m.CEMIFrame.LengthInBits()

	// Simple field (groupAddress)
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

	// Simple field (crc)
	lengthInBits += 8

	return lengthInBits
}

func (m CEMIFrameDataExt) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIFrameDataExtParse(io *spi.ReadBuffer) (CEMIFrameInitializer, error) {

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
	data := make([]int8, uint16(dataLength)-uint16(uint16(1)))
	for curItem := uint16(0); curItem < uint16(uint16(dataLength)-uint16(uint16(1))); curItem++ {

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
	return NewCEMIFrameDataExt(groupAddress, hopCount, extendedFrameFormat, sourceAddress, destinationAddress, dataLength, tcpi, counter, apci, dataFirstByte, data, crc), nil
}

func (m CEMIFrameDataExt) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (groupAddress)
		groupAddress := bool(m.GroupAddress)
		_groupAddressErr := io.WriteBit((bool)(groupAddress))
		if _groupAddressErr != nil {
			return errors.New("Error serializing 'groupAddress' field " + _groupAddressErr.Error())
		}

		// Simple Field (hopCount)
		hopCount := uint8(m.HopCount)
		_hopCountErr := io.WriteUint8(3, hopCount)
		if _hopCountErr != nil {
			return errors.New("Error serializing 'hopCount' field " + _hopCountErr.Error())
		}

		// Simple Field (extendedFrameFormat)
		extendedFrameFormat := uint8(m.ExtendedFrameFormat)
		_extendedFrameFormatErr := io.WriteUint8(4, extendedFrameFormat)
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
		_dataLengthErr := io.WriteUint8(8, dataLength)
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
		_counterErr := io.WriteUint8(4, counter)
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
		_dataFirstByteErr := io.WriteInt8(6, dataFirstByte)
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
		_crcErr := io.WriteUint8(8, crc)
		if _crcErr != nil {
			return errors.New("Error serializing 'crc' field " + _crcErr.Error())
		}

		return nil
	}
	return CEMIFrameSerialize(io, m.CEMIFrame, CastICEMIFrame(m), ser)
}
