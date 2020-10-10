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
type CEMIDataFrame struct {
	standardFrame           bool
	polling                 bool
	notRepeated             bool
	notAckFrame             bool
	priority                CEMIPriority
	acknowledgeRequested    bool
	errorFlag               bool
	groupDestinationAddress bool
	hopCount                uint8
	extendedFrameFormat     uint8
	sourceAddress           KNXAddress
	destinationAddress      []int8
	dataLength              uint8
	tcpi                    TPCI
	counter                 uint8
	apci                    APCI
	dataFirstByte           int8
	data                    []int8
}

// The corresponding interface
type ICEMIDataFrame interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewCEMIDataFrame(standardFrame bool, polling bool, notRepeated bool, notAckFrame bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool, groupDestinationAddress bool, hopCount uint8, extendedFrameFormat uint8, sourceAddress KNXAddress, destinationAddress []int8, dataLength uint8, tcpi TPCI, counter uint8, apci APCI, dataFirstByte int8, data []int8) spi.Message {
	return &CEMIDataFrame{standardFrame: standardFrame, polling: polling, notRepeated: notRepeated, notAckFrame: notAckFrame, priority: priority, acknowledgeRequested: acknowledgeRequested, errorFlag: errorFlag, groupDestinationAddress: groupDestinationAddress, hopCount: hopCount, extendedFrameFormat: extendedFrameFormat, sourceAddress: sourceAddress, destinationAddress: destinationAddress, dataLength: dataLength, tcpi: tcpi, counter: counter, apci: apci, dataFirstByte: dataFirstByte, data: data}
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
	lengthInBits += m.sourceAddress.LengthInBits()

	// Array field
	if len(m.destinationAddress) > 0 {
		lengthInBits += 8 * uint16(len(m.destinationAddress))
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
	if len(m.data) > 0 {
		lengthInBits += 8 * uint16(len(m.data))
	}

	return lengthInBits
}

func (m CEMIDataFrame) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIDataFrameParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (standardFrame)
	var standardFrame bool = io.ReadBit()

	// Simple Field (polling)
	var polling bool = io.ReadBit()

	// Simple Field (notRepeated)
	var notRepeated bool = io.ReadBit()

	// Simple Field (notAckFrame)
	var notAckFrame bool = io.ReadBit()

	// Enum field (priority)
	priority, _priorityErr := CEMIPriorityParse(io)
	if _priorityErr != nil {
		return nil, errors.New("Error parsing 'priority' field " + _priorityErr.Error())
	}

	// Simple Field (acknowledgeRequested)
	var acknowledgeRequested bool = io.ReadBit()

	// Simple Field (errorFlag)
	var errorFlag bool = io.ReadBit()

	// Simple Field (groupDestinationAddress)
	var groupDestinationAddress bool = io.ReadBit()

	// Simple Field (hopCount)
	var hopCount uint8 = io.ReadUint8(3)

	// Simple Field (extendedFrameFormat)
	var extendedFrameFormat uint8 = io.ReadUint8(4)

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
		destinationAddress := make([]int8, 2)
		for curItem := uint16(0); curItem < uint16(2); curItem++ {

			destinationAddress = append(destinationAddress, io.ReadInt8(8))
		}
	}

	// Simple Field (dataLength)
	var dataLength uint8 = io.ReadUint8(8)

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
		data := make([]int8, (dataLength)-(1))
		for curItem := uint16(0); curItem < uint16((dataLength)-(1)); curItem++ {

			data = append(data, io.ReadInt8(8))
		}
	}

	// Create the instance
	return NewCEMIDataFrame(standardFrame, polling, notRepeated, notAckFrame, priority, acknowledgeRequested, errorFlag, groupDestinationAddress, hopCount, extendedFrameFormat, sourceAddress, destinationAddress, dataLength, tcpi, counter, apci, dataFirstByte, data), nil
}

func (m CEMIDataFrame) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ICEMIDataFrame); ok {

			// Simple Field (standardFrame)
			var standardFrame bool = m.standardFrame
			io.WriteBit((bool)(standardFrame))

			// Simple Field (polling)
			var polling bool = m.polling
			io.WriteBit((bool)(polling))

			// Simple Field (notRepeated)
			var notRepeated bool = m.notRepeated
			io.WriteBit((bool)(notRepeated))

			// Simple Field (notAckFrame)
			var notAckFrame bool = m.notAckFrame
			io.WriteBit((bool)(notAckFrame))

			// Enum field (priority)
			priority := m.priority
			priority.Serialize(io)

			// Simple Field (acknowledgeRequested)
			var acknowledgeRequested bool = m.acknowledgeRequested
			io.WriteBit((bool)(acknowledgeRequested))

			// Simple Field (errorFlag)
			var errorFlag bool = m.errorFlag
			io.WriteBit((bool)(errorFlag))

			// Simple Field (groupDestinationAddress)
			var groupDestinationAddress bool = m.groupDestinationAddress
			io.WriteBit((bool)(groupDestinationAddress))

			// Simple Field (hopCount)
			var hopCount uint8 = m.hopCount
			io.WriteUint8(3, (hopCount))

			// Simple Field (extendedFrameFormat)
			var extendedFrameFormat uint8 = m.extendedFrameFormat
			io.WriteUint8(4, (extendedFrameFormat))

			// Simple Field (sourceAddress)
			var sourceAddress KNXAddress = m.sourceAddress
			sourceAddress.Serialize(io)

			// Array Field (destinationAddress)
			if m.destinationAddress != nil {
				for _, _element := range m.destinationAddress {
					io.WriteInt8(8, _element)
				}
			}

			// Simple Field (dataLength)
			var dataLength uint8 = m.dataLength
			io.WriteUint8(8, (dataLength))

			// Enum field (tcpi)
			tcpi := m.tcpi
			tcpi.Serialize(io)

			// Simple Field (counter)
			var counter uint8 = m.counter
			io.WriteUint8(4, (counter))

			// Enum field (apci)
			apci := m.apci
			apci.Serialize(io)

			// Simple Field (dataFirstByte)
			var dataFirstByte int8 = m.dataFirstByte
			io.WriteInt8(6, (dataFirstByte))

			// Array Field (data)
			if m.data != nil {
				for _, _element := range m.data {
					io.WriteInt8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
