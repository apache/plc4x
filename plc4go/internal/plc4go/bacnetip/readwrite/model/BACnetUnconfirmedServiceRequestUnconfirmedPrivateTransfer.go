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
	"strconv"
)

// Constant values.
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_VENDORIDHEADER uint8 = 0x09
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_SERVICENUMBERHEADER uint8 = 0x1A
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESOPENINGTAG uint8 = 0x2E
const BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESCLOSINGTAG uint8 = 0x2F

// The data-structure of this message
type BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer struct {
	VendorId      uint8
	ServiceNumber uint16
	Values        []int8
	Parent        *BACnetUnconfirmedServiceRequest
	IBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) ServiceChoice() uint8 {
	return 0x04
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) InitializeParent(parent *BACnetUnconfirmedServiceRequest) {
}

func NewBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(vendorId uint8, serviceNumber uint16, values []int8) *BACnetUnconfirmedServiceRequest {
	child := &BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer{
		VendorId:      vendorId,
		ServiceNumber: serviceNumber,
		Values:        values,
		Parent:        NewBACnetUnconfirmedServiceRequest(),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(structType interface{}) *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer {
	castFunc := func(typ interface{}) *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer {
		if casted, ok := typ.(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer); ok {
			return &casted
		}
		if casted, ok := typ.(*BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer); ok {
			return casted
		}
		if casted, ok := typ.(BACnetUnconfirmedServiceRequest); ok {
			return CastBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(casted.Child)
		}
		if casted, ok := typ.(*BACnetUnconfirmedServiceRequest); ok {
			return CastBACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) GetTypeName() string {
	return "BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer"
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) LengthInBits() uint16 {
	lengthInBits := uint16(0)

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
	if len(m.Values) > 0 {
		lengthInBits += 8 * uint16(len(m.Values))
	}

	// Const Field (listOfValuesClosingTag)
	lengthInBits += 8

	return lengthInBits
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransferParse(io *utils.ReadBuffer, len uint16) (*BACnetUnconfirmedServiceRequest, error) {

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
	// Length array
	values := make([]int8, 0)
	_valuesLength := uint16(len) - uint16(uint16(8))
	_valuesEndPos := io.GetPos() + uint16(_valuesLength)
	for io.GetPos() < _valuesEndPos {
		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'values' field " + _err.Error())
		}
		values = append(values, _item)
	}

	// Const Field (listOfValuesClosingTag)
	listOfValuesClosingTag, _listOfValuesClosingTagErr := io.ReadUint8(8)
	if _listOfValuesClosingTagErr != nil {
		return nil, errors.New("Error parsing 'listOfValuesClosingTag' field " + _listOfValuesClosingTagErr.Error())
	}
	if listOfValuesClosingTag != BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESCLOSINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer_LISTOFVALUESCLOSINGTAG)) + " but got " + strconv.Itoa(int(listOfValuesClosingTag)))
	}

	// Create a partially initialized instance
	_child := &BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer{
		VendorId:      vendorId,
		ServiceNumber: serviceNumber,
		Values:        values,
		Parent:        &BACnetUnconfirmedServiceRequest{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) Serialize(io utils.WriteBuffer) error {
	ser := func() error {

		// Const Field (vendorIdHeader)
		_vendorIdHeaderErr := io.WriteUint8(8, 0x09)
		if _vendorIdHeaderErr != nil {
			return errors.New("Error serializing 'vendorIdHeader' field " + _vendorIdHeaderErr.Error())
		}

		// Simple Field (vendorId)
		vendorId := uint8(m.VendorId)
		_vendorIdErr := io.WriteUint8(8, (vendorId))
		if _vendorIdErr != nil {
			return errors.New("Error serializing 'vendorId' field " + _vendorIdErr.Error())
		}

		// Const Field (serviceNumberHeader)
		_serviceNumberHeaderErr := io.WriteUint8(8, 0x1A)
		if _serviceNumberHeaderErr != nil {
			return errors.New("Error serializing 'serviceNumberHeader' field " + _serviceNumberHeaderErr.Error())
		}

		// Simple Field (serviceNumber)
		serviceNumber := uint16(m.ServiceNumber)
		_serviceNumberErr := io.WriteUint16(16, (serviceNumber))
		if _serviceNumberErr != nil {
			return errors.New("Error serializing 'serviceNumber' field " + _serviceNumberErr.Error())
		}

		// Const Field (listOfValuesOpeningTag)
		_listOfValuesOpeningTagErr := io.WriteUint8(8, 0x2E)
		if _listOfValuesOpeningTagErr != nil {
			return errors.New("Error serializing 'listOfValuesOpeningTag' field " + _listOfValuesOpeningTagErr.Error())
		}

		// Array Field (values)
		if m.Values != nil {
			for _, _element := range m.Values {
				_elementErr := io.WriteInt8(8, _element)
				if _elementErr != nil {
					return errors.New("Error serializing 'values' field " + _elementErr.Error())
				}
			}
		}

		// Const Field (listOfValuesClosingTag)
		_listOfValuesClosingTagErr := io.WriteUint8(8, 0x2F)
		if _listOfValuesClosingTagErr != nil {
			return errors.New("Error serializing 'listOfValuesClosingTag' field " + _listOfValuesClosingTagErr.Error())
		}

		return nil
	}
	return m.Parent.SerializeParent(io, m, ser)
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var token xml.Token
	var err error
	token = start
	for {
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			case "vendorId":
				var data uint8
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.VendorId = data
			case "serviceNumber":
				var data uint16
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ServiceNumber = data
			case "values":
				var _encoded string
				if err := d.DecodeElement(&_encoded, &tok); err != nil {
					return err
				}
				_decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
				_len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
				if err != nil {
					return err
				}
				m.Values = utils.ByteArrayToInt8Array(_decoded[0:_len])
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

func (m *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.VendorId, xml.StartElement{Name: xml.Name{Local: "vendorId"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ServiceNumber, xml.StartElement{Name: xml.Name{Local: "serviceNumber"}}); err != nil {
		return err
	}
	_encodedValues := make([]byte, base64.StdEncoding.EncodedLen(len(m.Values)))
	base64.StdEncoding.Encode(_encodedValues, utils.Int8ArrayToByteArray(m.Values))
	if err := e.EncodeElement(_encodedValues, xml.StartElement{Name: xml.Name{Local: "values"}}); err != nil {
		return err
	}
	return nil
}
