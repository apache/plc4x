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
	"strconv"
)

// Constant values.
const BACnetUnconfirmedServiceRequestIAm_OBJECTIDENTIFIERHEADER uint8 = 0xC4
const BACnetUnconfirmedServiceRequestIAm_MAXIMUMAPDULENGTHACCEPTEDHEADER uint8 = 0x04
const BACnetUnconfirmedServiceRequestIAm_SEGMENTATIONSUPPORTEDHEADER uint8 = 0x91
const BACnetUnconfirmedServiceRequestIAm_VENDORIDHEADER uint8 = 0x21

// The data-structure of this message
type BACnetUnconfirmedServiceRequestIAm struct {
	objectType                      uint16
	objectInstanceNumber            uint32
	maximumApduLengthAcceptedLength uint8
	maximumApduLengthAccepted       []int8
	segmentationSupported           uint8
	vendorId                        uint8
	BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestIAm interface {
	IBACnetUnconfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestIAm) ServiceChoice() uint8 {
	return 0x00
}

func (m BACnetUnconfirmedServiceRequestIAm) initialize() spi.Message {
	return m
}

func NewBACnetUnconfirmedServiceRequestIAm(objectType uint16, objectInstanceNumber uint32, maximumApduLengthAcceptedLength uint8, maximumApduLengthAccepted []int8, segmentationSupported uint8, vendorId uint8) BACnetUnconfirmedServiceRequestInitializer {
	return &BACnetUnconfirmedServiceRequestIAm{objectType: objectType, objectInstanceNumber: objectInstanceNumber, maximumApduLengthAcceptedLength: maximumApduLengthAcceptedLength, maximumApduLengthAccepted: maximumApduLengthAccepted, segmentationSupported: segmentationSupported, vendorId: vendorId}
}

func CastIBACnetUnconfirmedServiceRequestIAm(structType interface{}) IBACnetUnconfirmedServiceRequestIAm {
	castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestIAm {
		if iBACnetUnconfirmedServiceRequestIAm, ok := typ.(IBACnetUnconfirmedServiceRequestIAm); ok {
			return iBACnetUnconfirmedServiceRequestIAm
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestIAm(structType interface{}) BACnetUnconfirmedServiceRequestIAm {
	castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestIAm {
		if sBACnetUnconfirmedServiceRequestIAm, ok := typ.(BACnetUnconfirmedServiceRequestIAm); ok {
			return sBACnetUnconfirmedServiceRequestIAm
		}
		return BACnetUnconfirmedServiceRequestIAm{}
	}
	return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestIAm) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

	// Const Field (objectIdentifierHeader)
	lengthInBits += 8

	// Simple field (objectType)
	lengthInBits += 10

	// Simple field (objectInstanceNumber)
	lengthInBits += 22

	// Const Field (maximumApduLengthAcceptedHeader)
	lengthInBits += 5

	// Simple field (maximumApduLengthAcceptedLength)
	lengthInBits += 3

	// Array field
	if len(m.maximumApduLengthAccepted) > 0 {
		lengthInBits += 8 * uint16(len(m.maximumApduLengthAccepted))
	}

	// Const Field (segmentationSupportedHeader)
	lengthInBits += 8

	// Simple field (segmentationSupported)
	lengthInBits += 8

	// Const Field (vendorIdHeader)
	lengthInBits += 8

	// Simple field (vendorId)
	lengthInBits += 8

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestIAm) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestIAmParse(io *spi.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

	// Const Field (objectIdentifierHeader)
	objectIdentifierHeader, _objectIdentifierHeaderErr := io.ReadUint8(8)
	if _objectIdentifierHeaderErr != nil {
		return nil, errors.New("Error parsing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
	}
	if objectIdentifierHeader != BACnetUnconfirmedServiceRequestIAm_OBJECTIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
	}

	// Simple Field (objectType)
	objectType, _objectTypeErr := io.ReadUint16(10)
	if _objectTypeErr != nil {
		return nil, errors.New("Error parsing 'objectType' field " + _objectTypeErr.Error())
	}

	// Simple Field (objectInstanceNumber)
	objectInstanceNumber, _objectInstanceNumberErr := io.ReadUint32(22)
	if _objectInstanceNumberErr != nil {
		return nil, errors.New("Error parsing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
	}

	// Const Field (maximumApduLengthAcceptedHeader)
	maximumApduLengthAcceptedHeader, _maximumApduLengthAcceptedHeaderErr := io.ReadUint8(5)
	if _maximumApduLengthAcceptedHeaderErr != nil {
		return nil, errors.New("Error parsing 'maximumApduLengthAcceptedHeader' field " + _maximumApduLengthAcceptedHeaderErr.Error())
	}
	if maximumApduLengthAcceptedHeader != BACnetUnconfirmedServiceRequestIAm_MAXIMUMAPDULENGTHACCEPTEDHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_MAXIMUMAPDULENGTHACCEPTEDHEADER)) + " but got " + strconv.Itoa(int(maximumApduLengthAcceptedHeader)))
	}

	// Simple Field (maximumApduLengthAcceptedLength)
	maximumApduLengthAcceptedLength, _maximumApduLengthAcceptedLengthErr := io.ReadUint8(3)
	if _maximumApduLengthAcceptedLengthErr != nil {
		return nil, errors.New("Error parsing 'maximumApduLengthAcceptedLength' field " + _maximumApduLengthAcceptedLengthErr.Error())
	}

	// Array field (maximumApduLengthAccepted)
	// Count array
	maximumApduLengthAccepted := make([]int8, maximumApduLengthAcceptedLength)
	for curItem := uint16(0); curItem < uint16(maximumApduLengthAcceptedLength); curItem++ {

		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'maximumApduLengthAccepted' field " + _err.Error())
		}
		maximumApduLengthAccepted[curItem] = _item
	}

	// Const Field (segmentationSupportedHeader)
	segmentationSupportedHeader, _segmentationSupportedHeaderErr := io.ReadUint8(8)
	if _segmentationSupportedHeaderErr != nil {
		return nil, errors.New("Error parsing 'segmentationSupportedHeader' field " + _segmentationSupportedHeaderErr.Error())
	}
	if segmentationSupportedHeader != BACnetUnconfirmedServiceRequestIAm_SEGMENTATIONSUPPORTEDHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_SEGMENTATIONSUPPORTEDHEADER)) + " but got " + strconv.Itoa(int(segmentationSupportedHeader)))
	}

	// Simple Field (segmentationSupported)
	segmentationSupported, _segmentationSupportedErr := io.ReadUint8(8)
	if _segmentationSupportedErr != nil {
		return nil, errors.New("Error parsing 'segmentationSupported' field " + _segmentationSupportedErr.Error())
	}

	// Const Field (vendorIdHeader)
	vendorIdHeader, _vendorIdHeaderErr := io.ReadUint8(8)
	if _vendorIdHeaderErr != nil {
		return nil, errors.New("Error parsing 'vendorIdHeader' field " + _vendorIdHeaderErr.Error())
	}
	if vendorIdHeader != BACnetUnconfirmedServiceRequestIAm_VENDORIDHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_VENDORIDHEADER)) + " but got " + strconv.Itoa(int(vendorIdHeader)))
	}

	// Simple Field (vendorId)
	vendorId, _vendorIdErr := io.ReadUint8(8)
	if _vendorIdErr != nil {
		return nil, errors.New("Error parsing 'vendorId' field " + _vendorIdErr.Error())
	}

	// Create the instance
	return NewBACnetUnconfirmedServiceRequestIAm(objectType, objectInstanceNumber, maximumApduLengthAcceptedLength, maximumApduLengthAccepted, segmentationSupported, vendorId), nil
}

func (m BACnetUnconfirmedServiceRequestIAm) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Const Field (objectIdentifierHeader)
		io.WriteUint8(8, 0xC4)

		// Simple Field (objectType)
		objectType := uint16(m.objectType)
		io.WriteUint16(10, (objectType))

		// Simple Field (objectInstanceNumber)
		objectInstanceNumber := uint32(m.objectInstanceNumber)
		io.WriteUint32(22, (objectInstanceNumber))

		// Const Field (maximumApduLengthAcceptedHeader)
		io.WriteUint8(5, 0x04)

		// Simple Field (maximumApduLengthAcceptedLength)
		maximumApduLengthAcceptedLength := uint8(m.maximumApduLengthAcceptedLength)
		io.WriteUint8(3, (maximumApduLengthAcceptedLength))

		// Array Field (maximumApduLengthAccepted)
		if m.maximumApduLengthAccepted != nil {
			for _, _element := range m.maximumApduLengthAccepted {
				io.WriteInt8(8, _element)
			}
		}

		// Const Field (segmentationSupportedHeader)
		io.WriteUint8(8, 0x91)

		// Simple Field (segmentationSupported)
		segmentationSupported := uint8(m.segmentationSupported)
		io.WriteUint8(8, (segmentationSupported))

		// Const Field (vendorIdHeader)
		io.WriteUint8(8, 0x21)

		// Simple Field (vendorId)
		vendorId := uint8(m.vendorId)
		io.WriteUint8(8, (vendorId))

	}
	BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}
