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
type CEMIFrameData struct {
	sourceAddress      KNXAddress
	destinationAddress []int8
	groupAddress       bool
	hopCount           uint8
	dataLength         uint8
	tcpi               TPCI
	counter            uint8
	apci               APCI
	dataFirstByte      int8
	data               []int8
	crc                uint8
	CEMIFrame
}

// The corresponding interface
type ICEMIFrameData interface {
	ICEMIFrame
	Serialize(io spi.WriteBuffer)
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

func (m CEMIFrameData) initialize(repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) spi.Message {
	m.repeated = repeated
	m.priority = priority
	m.acknowledgeRequested = acknowledgeRequested
	m.errorFlag = errorFlag
	return m
}

func NewCEMIFrameData(sourceAddress KNXAddress, destinationAddress []int8, groupAddress bool, hopCount uint8, dataLength uint8, tcpi TPCI, counter uint8, apci APCI, dataFirstByte int8, data []int8, crc uint8) CEMIFrameInitializer {
	return &CEMIFrameData{sourceAddress: sourceAddress, destinationAddress: destinationAddress, groupAddress: groupAddress, hopCount: hopCount, dataLength: dataLength, tcpi: tcpi, counter: counter, apci: apci, dataFirstByte: dataFirstByte, data: data, crc: crc}
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
		return CEMIFrameData{}
	}
	return castFunc(structType)
}

func (m CEMIFrameData) LengthInBits() uint16 {
	var lengthInBits uint16 = m.CEMIFrame.LengthInBits()

	// Simple field (sourceAddress)
	lengthInBits += m.sourceAddress.LengthInBits()

	// Array field
	if len(m.destinationAddress) > 0 {
		lengthInBits += 8 * uint16(len(m.destinationAddress))
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
	if len(m.data) > 0 {
		lengthInBits += 8 * uint16(len(m.data))
	}

	// Simple field (crc)
	lengthInBits += 8

	return lengthInBits
}

func (m CEMIFrameData) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIFrameDataParse(io spi.ReadBuffer) (CEMIFrameInitializer, error) {

	// Simple Field (sourceAddress)
	_sourceAddressMessage, _err := KNXAddressParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'sourceAddress'. " + _err.Error())
	}
	var sourceAddress KNXAddress
	sourceAddress, _sourceAddressOk := _sourceAddressMessage.(KNXAddress)
	if !_sourceAddressOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_sourceAddressMessage).Name() + " to KNXAddress")
	}

	// Array field (destinationAddress)
	var destinationAddress []int8
	// Count array
	{
		destinationAddress := make([]int8, uint16(2))
		for curItem := uint16(0); curItem < uint16(uint16(2)); curItem++ {

			destinationAddress = append(destinationAddress, io.ReadInt8(8))
		}
	}

	// Simple Field (groupAddress)
	var groupAddress bool = io.ReadBit()

	// Simple Field (hopCount)
	var hopCount uint8 = io.ReadUint8(3)

	// Simple Field (dataLength)
	var dataLength uint8 = io.ReadUint8(4)

	// Enum field (tcpi)
	tcpi, _tcpiErr := TPCIParse(io)
	if _tcpiErr != nil {
		return nil, errors.New("Error parsing 'tcpi' field " + _tcpiErr.Error())
	}

	// Simple Field (counter)
	var counter uint8 = io.ReadUint8(4)

	// Enum field (apci)
	apci, _apciErr := APCIParse(io)
	if _apciErr != nil {
		return nil, errors.New("Error parsing 'apci' field " + _apciErr.Error())
	}

	// Simple Field (dataFirstByte)
	var dataFirstByte int8 = io.ReadInt8(6)

	// Array field (data)
	var data []int8
	// Count array
	{
		data := make([]int8, uint16(dataLength)-uint16(uint16(1)))
		for curItem := uint16(0); curItem < uint16(uint16(dataLength)-uint16(uint16(1))); curItem++ {

			data = append(data, io.ReadInt8(8))
		}
	}

	// Simple Field (crc)
	var crc uint8 = io.ReadUint8(8)

	// Create the instance
	return NewCEMIFrameData(sourceAddress, destinationAddress, groupAddress, hopCount, dataLength, tcpi, counter, apci, dataFirstByte, data, crc), nil
}

func (m CEMIFrameData) Serialize(io spi.WriteBuffer) {

	// Simple Field (sourceAddress)
	sourceAddress := KNXAddress(m.sourceAddress)
	sourceAddress.Serialize(io)

	// Array Field (destinationAddress)
	if m.destinationAddress != nil {
		for _, _element := range m.destinationAddress {
			io.WriteInt8(8, _element)
		}
	}

	// Simple Field (groupAddress)
	groupAddress := bool(m.groupAddress)
	io.WriteBit((bool)(groupAddress))

	// Simple Field (hopCount)
	hopCount := uint8(m.hopCount)
	io.WriteUint8(3, (hopCount))

	// Simple Field (dataLength)
	dataLength := uint8(m.dataLength)
	io.WriteUint8(4, (dataLength))

	// Enum field (tcpi)
	tcpi := TPCI(m.tcpi)
	tcpi.Serialize(io)

	// Simple Field (counter)
	counter := uint8(m.counter)
	io.WriteUint8(4, (counter))

	// Enum field (apci)
	apci := APCI(m.apci)
	apci.Serialize(io)

	// Simple Field (dataFirstByte)
	dataFirstByte := int8(m.dataFirstByte)
	io.WriteInt8(6, (dataFirstByte))

	// Array Field (data)
	if m.data != nil {
		for _, _element := range m.data {
			io.WriteInt8(8, _element)
		}
	}

	// Simple Field (crc)
	crc := uint8(m.crc)
	io.WriteUint8(8, (crc))
}
