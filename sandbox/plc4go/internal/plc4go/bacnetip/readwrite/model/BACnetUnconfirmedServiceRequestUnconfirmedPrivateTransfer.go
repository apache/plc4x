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
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_VENDORIDHEADER uint8 = 0x09
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_SERVICENUMBERHEADER uint8 = 0x1A
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESOPENINGTAG uint8 = 0x2E
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESCLOSINGTAG uint8 = 0x2F

// The data-structure of this message
type BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer struct {
	vendorId      uint8
	serviceNumber uint16
	values        []int8
	BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer interface {
	IBACnetUnconfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) ServiceChoice() uint8 {
	return 0x04
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) initialize() spi.Message {
	return m
}

func NewBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(vendorId uint8, serviceNumber uint16, values []int8) BACnetUnconfirmedServiceRequestInitializer {
	return &BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer{vendorId: vendorId, serviceNumber: serviceNumber, values: values}
}

func CastIBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(structType interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer {
	castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer {
		if iBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer, ok := typ.(IBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer); ok {
			return iBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(structType interface{}) BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer {
	castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer {
		if sBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer, ok := typ.(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer); ok {
			return sBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
		}
		return BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer{}
	}
	return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

	// Const Field (vendorIdHeader)
	lengthInBits += 8

	// Simple field (vendorId)
	lengthInBits += 8

	// Const Field (serviceNumberHeader)
	lengthInBits += 8

	// Simple field (serviceNumber)
	lengthInBits += 16

	// Const Field (listOfValuesOpeningTag)
	lengthInBits += 8

	// Array field
	if len(m.values) > 0 {
		lengthInBits += 8 * uint16(len(m.values))
	}

	// Const Field (listOfValuesClosingTag)
	lengthInBits += 8

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransferParse(io spi.ReadBuffer, len uint16) (BACnetUnconfirmedServiceRequestInitializer, error) {

	// Const Field (vendorIdHeader)
	vendorIdHeader, _vendorIdHeaderErr := io.ReadUint8(8)
	if _vendorIdHeaderErr != nil {
		return nil, errors.New("Error parsing 'vendorIdHeader' field " + _vendorIdHeaderErr.Error())
	}
	if vendorIdHeader != BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_VENDORIDHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_VENDORIDHEADER)) + " but got " + strconv.Itoa(int(vendorIdHeader)))
	}

	// Simple Field (vendorId)
	vendorId, _vendorIdErr := io.ReadUint8(8)
	if _vendorIdErr != nil {
		return nil, errors.New("Error parsing 'vendorId' field " + _vendorIdErr.Error())
	}

	// Const Field (serviceNumberHeader)
	serviceNumberHeader, _serviceNumberHeaderErr := io.ReadUint8(8)
	if _serviceNumberHeaderErr != nil {
		return nil, errors.New("Error parsing 'serviceNumberHeader' field " + _serviceNumberHeaderErr.Error())
	}
	if serviceNumberHeader != BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_SERVICENUMBERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_SERVICENUMBERHEADER)) + " but got " + strconv.Itoa(int(serviceNumberHeader)))
	}

	// Simple Field (serviceNumber)
	serviceNumber, _serviceNumberErr := io.ReadUint16(16)
	if _serviceNumberErr != nil {
		return nil, errors.New("Error parsing 'serviceNumber' field " + _serviceNumberErr.Error())
	}

	// Const Field (listOfValuesOpeningTag)
	listOfValuesOpeningTag, _listOfValuesOpeningTagErr := io.ReadUint8(8)
	if _listOfValuesOpeningTagErr != nil {
		return nil, errors.New("Error parsing 'listOfValuesOpeningTag' field " + _listOfValuesOpeningTagErr.Error())
	}
	if listOfValuesOpeningTag != BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESOPENINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESOPENINGTAG)) + " but got " + strconv.Itoa(int(listOfValuesOpeningTag)))
	}

	// Array field (values)
	var values []int8
	// Length array
	_valuesLength := uint16(len) - uint16(uint16(8))
	_valuesEndPos := io.GetPos() + uint16(_valuesLength)
	for io.GetPos() < _valuesEndPos {
		_valuesVal, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'values' field " + _err.Error())
		}
		values = append(values, _valuesVal)
	}

	// Const Field (listOfValuesClosingTag)
	listOfValuesClosingTag, _listOfValuesClosingTagErr := io.ReadUint8(8)
	if _listOfValuesClosingTagErr != nil {
		return nil, errors.New("Error parsing 'listOfValuesClosingTag' field " + _listOfValuesClosingTagErr.Error())
	}
	if listOfValuesClosingTag != BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESCLOSINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESCLOSINGTAG)) + " but got " + strconv.Itoa(int(listOfValuesClosingTag)))
	}

	// Create the instance
	return NewBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(vendorId, serviceNumber, values), nil
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Const Field (vendorIdHeader)
		io.WriteUint8(8, 0x09)

		// Simple Field (vendorId)
		vendorId := uint8(m.vendorId)
		io.WriteUint8(8, (vendorId))

		// Const Field (serviceNumberHeader)
		io.WriteUint8(8, 0x1A)

		// Simple Field (serviceNumber)
		serviceNumber := uint16(m.serviceNumber)
		io.WriteUint16(16, (serviceNumber))

		// Const Field (listOfValuesOpeningTag)
		io.WriteUint8(8, 0x2E)

		// Array Field (values)
		if m.values != nil {
			for _, _element := range m.values {
				io.WriteInt8(8, _element)
			}
		}

		// Const Field (listOfValuesClosingTag)
		io.WriteUint8(8, 0x2F)

	}
	BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}
