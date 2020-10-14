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
	log "github.com/sirupsen/logrus"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type S7AddressAny struct {
	TransportSize    ITransportSize
	NumberOfElements uint16
	DbNumber         uint16
	Area             IMemoryArea
	ByteAddress      uint16
	BitAddress       uint8
	S7Address
}

// The corresponding interface
type IS7AddressAny interface {
	IS7Address
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7AddressAny) AddressType() uint8 {
	return 0x10
}

func (m S7AddressAny) initialize() spi.Message {
	return m
}

func NewS7AddressAny(transportSize ITransportSize, numberOfElements uint16, dbNumber uint16, area IMemoryArea, byteAddress uint16, bitAddress uint8) S7AddressInitializer {
	return &S7AddressAny{TransportSize: transportSize, NumberOfElements: numberOfElements, DbNumber: dbNumber, Area: area, ByteAddress: byteAddress, BitAddress: bitAddress}
}

func CastIS7AddressAny(structType interface{}) IS7AddressAny {
	castFunc := func(typ interface{}) IS7AddressAny {
		if iS7AddressAny, ok := typ.(IS7AddressAny); ok {
			return iS7AddressAny
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7AddressAny(structType interface{}) S7AddressAny {
	castFunc := func(typ interface{}) S7AddressAny {
		if sS7AddressAny, ok := typ.(S7AddressAny); ok {
			return sS7AddressAny
		}
		return S7AddressAny{}
	}
	return castFunc(structType)
}

func (m S7AddressAny) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Address.LengthInBits()

	// Enum Field (transportSize)
	lengthInBits += 8

	// Simple field (numberOfElements)
	lengthInBits += 16

	// Simple field (dbNumber)
	lengthInBits += 16

	// Enum Field (area)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 5

	// Simple field (byteAddress)
	lengthInBits += 16

	// Simple field (bitAddress)
	lengthInBits += 3

	return lengthInBits
}

func (m S7AddressAny) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7AddressAnyParse(io *spi.ReadBuffer) (S7AddressInitializer, error) {

	// Enum field (transportSize)
	transportSize, _transportSizeErr := TransportSizeParse(io)
	if _transportSizeErr != nil {
		return nil, errors.New("Error parsing 'transportSize' field " + _transportSizeErr.Error())
	}

	// Simple Field (numberOfElements)
	numberOfElements, _numberOfElementsErr := io.ReadUint16(16)
	if _numberOfElementsErr != nil {
		return nil, errors.New("Error parsing 'numberOfElements' field " + _numberOfElementsErr.Error())
	}

	// Simple Field (dbNumber)
	dbNumber, _dbNumberErr := io.ReadUint16(16)
	if _dbNumberErr != nil {
		return nil, errors.New("Error parsing 'dbNumber' field " + _dbNumberErr.Error())
	}

	// Enum field (area)
	area, _areaErr := MemoryAreaParse(io)
	if _areaErr != nil {
		return nil, errors.New("Error parsing 'area' field " + _areaErr.Error())
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(5)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (byteAddress)
	byteAddress, _byteAddressErr := io.ReadUint16(16)
	if _byteAddressErr != nil {
		return nil, errors.New("Error parsing 'byteAddress' field " + _byteAddressErr.Error())
	}

	// Simple Field (bitAddress)
	bitAddress, _bitAddressErr := io.ReadUint8(3)
	if _bitAddressErr != nil {
		return nil, errors.New("Error parsing 'bitAddress' field " + _bitAddressErr.Error())
	}

	// Create the instance
	return NewS7AddressAny(transportSize, numberOfElements, dbNumber, area, byteAddress, bitAddress), nil
}

func (m S7AddressAny) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Enum field (transportSize)
		transportSize := CastTransportSize(m.TransportSize)
		_transportSizeErr := transportSize.Serialize(io)
		if _transportSizeErr != nil {
			return errors.New("Error serializing 'transportSize' field " + _transportSizeErr.Error())
		}

		// Simple Field (numberOfElements)
		numberOfElements := uint16(m.NumberOfElements)
		_numberOfElementsErr := io.WriteUint16(16, (numberOfElements))
		if _numberOfElementsErr != nil {
			return errors.New("Error serializing 'numberOfElements' field " + _numberOfElementsErr.Error())
		}

		// Simple Field (dbNumber)
		dbNumber := uint16(m.DbNumber)
		_dbNumberErr := io.WriteUint16(16, (dbNumber))
		if _dbNumberErr != nil {
			return errors.New("Error serializing 'dbNumber' field " + _dbNumberErr.Error())
		}

		// Enum field (area)
		area := CastMemoryArea(m.Area)
		_areaErr := area.Serialize(io)
		if _areaErr != nil {
			return errors.New("Error serializing 'area' field " + _areaErr.Error())
		}

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(5, uint8(0x00))
			if _err != nil {
				return errors.New("Error serializing 'reserved' field " + _err.Error())
			}
		}

		// Simple Field (byteAddress)
		byteAddress := uint16(m.ByteAddress)
		_byteAddressErr := io.WriteUint16(16, (byteAddress))
		if _byteAddressErr != nil {
			return errors.New("Error serializing 'byteAddress' field " + _byteAddressErr.Error())
		}

		// Simple Field (bitAddress)
		bitAddress := uint8(m.BitAddress)
		_bitAddressErr := io.WriteUint8(3, (bitAddress))
		if _bitAddressErr != nil {
			return errors.New("Error serializing 'bitAddress' field " + _bitAddressErr.Error())
		}

		return nil
	}
	return S7AddressSerialize(io, m.S7Address, CastIS7Address(m), ser)
}
