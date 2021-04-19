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

package utils

import (
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/big"
)

func NewXmlReadBuffer(reader io.Reader) ReadBuffer {
	return &xmlReadBuffer{
		xml.NewDecoder(reader),
		1,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type xmlReadBuffer struct {
	*xml.Decoder
	pos uint
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (x *xmlReadBuffer) Reset() {
	panic("implement me")
}

func (x *xmlReadBuffer) GetPos() uint16 {
	return uint16(x.pos / 8)
}

func (x *xmlReadBuffer) GetBytes() []uint8 {
	panic("implement me")
}

func (x *xmlReadBuffer) GetTotalBytes() uint64 {
	panic("implement me")
}

func (x *xmlReadBuffer) HasMore(bitLength uint8) bool {
	return true
}

func (x *xmlReadBuffer) PeekByte(offset uint8) uint8 {
	panic("implement me")
}

func (x *xmlReadBuffer) PullContext(logicalName string, readerArgs ...WithReaderArgs) error {
	startElement, err := x.travelToNextStartElement()
	if err != nil {
		return err
	}
	if startElement.Name.Local != logicalName {
		return errors.Errorf("Unexpected Start element '%s'. Expected '%s'", startElement.Name.Local, logicalName)
	}
	return nil
}

func (x *xmlReadBuffer) ReadBit(logicalName string, readerArgs ...WithReaderArgs) (bool, error) {
	var value bool
	err := x.decode(logicalName, "bit", 1, readerArgs, &value)
	if err != nil {
		return false, err
	}
	x.move(1)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint8, error) {
	var value uint8
	err := x.decode(logicalName, "uint8", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint16, error) {
	var value uint16
	err := x.decode(logicalName, "uint16", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint32, error) {
	var value uint32
	err := x.decode(logicalName, "uint32", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint64, error) {
	var value uint64
	err := x.decode(logicalName, "uint64", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int8, error) {
	var value int8
	err := x.decode(logicalName, "int8", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int16, error) {
	var value int16
	err := x.decode(logicalName, "int16", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int32, error) {
	var value int32
	err := x.decode(logicalName, "int32", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int64, error) {
	var value int64
	err := x.decode(logicalName, "int64", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadBigInt(logicalName string, bitLength uint64, readerArgs ...WithReaderArgs) (*big.Int, error) {
	var value big.Int
	// TODO: bitLength is too short for a big int
	err := x.decode(logicalName, "bigInt", uint8(bitLength), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	x.move(uint8(bitLength))
	return &value, nil
}

func (x *xmlReadBuffer) ReadFloat32(logicalName string, signed bool, exponentBitLength uint8, mantissaBitLength uint8, readerArgs ...WithReaderArgs) (float32, error) {
	bitLength := exponentBitLength + mantissaBitLength
	if signed {
		bitLength += 1
	}
	var value float32
	err := x.decode(logicalName, "float32", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadFloat64(logicalName string, signed bool, exponentBitLength uint8, mantissaBitLength uint8, readerArgs ...WithReaderArgs) (float64, error) {
	bitLength := exponentBitLength + mantissaBitLength
	if signed {
		bitLength += 1
	}
	var value float64
	err := x.decode(logicalName, "float64", bitLength, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadBigFloat(logicalName string, signed bool, exponentBitLength uint8, mantissaBitLength uint8, readerArgs ...WithReaderArgs) (*big.Float, error) {
	bitLength := exponentBitLength + mantissaBitLength
	if signed {
		bitLength += 1
	}
	var value big.Float
	err := x.decode(logicalName, "bigFloat", bitLength, readerArgs, &value)
	if err != nil {
		return nil, err
	}
	x.move(bitLength)
	return &value, nil
}

func (x *xmlReadBuffer) ReadString(logicalName string, bitLength uint32, readerArgs ...WithReaderArgs) (string, error) {
	var value string
	// TODO: bitlength too short
	err := x.decode(logicalName, "string", uint8(bitLength), readerArgs, &value)
	if err != nil {
		return "", err
	}
	x.move(uint8(bitLength))
	return value, nil
}

func (x *xmlReadBuffer) CloseContext(logicalName string, _ ...WithReaderArgs) error {
	endElement, err := x.travelToNextEndElement()
	if err != nil {
		return err
	}
	if endElement.Name.Local != logicalName {
		return errors.Errorf("Unexpected End element '%s'. Expected '%s'", endElement.Name.Local, logicalName)
	}
	return nil
}

func (x *xmlReadBuffer) move(bits uint8) {
	x.pos += uint(bits)
}

func (x *xmlReadBuffer) travelToNextStartElement() (xml.StartElement, error) {
	var startElement xml.StartElement
findTheStartToken:
	for {
		token, err := x.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch token.(type) {
		case xml.StartElement:
			startElement = token.(xml.StartElement)
			break findTheStartToken
		case xml.EndElement:
			return xml.StartElement{}, errors.Errorf("unexpected end element %s", token.(xml.EndElement).Name)
		}
	}
	return startElement, nil
}

func (x *xmlReadBuffer) travelToNextEndElement() (xml.EndElement, error) {
	var endElement xml.EndElement
findTheEndToken:
	for {
		token, err := x.Token()
		if err != nil {
			return xml.EndElement{}, err
		}
		switch token.(type) {
		case xml.EndElement:
			endElement = token.(xml.EndElement)
			break findTheEndToken
		case xml.StartElement:
			return xml.EndElement{}, errors.Errorf("unexpected start element %s", token.(xml.StartElement).Name)
		}
	}
	return endElement, nil
}

func (x *xmlReadBuffer) decode(logicalName string, dataType string, bitLength uint8, readerArgs []WithReaderArgs, targetValue interface{}) error {
	startElement, err := x.travelToNextStartElement()
	if err != nil {
		return err
	}
	err = validateStartElement(startElement, logicalName, dataType, bitLength, readerArgs...)
	if err != nil {
		return err
	}
	err = x.DecodeElement(targetValue, &startElement)
	if err != nil {
		return err
	}
	return nil
}

func validateStartElement(startElement xml.StartElement, logicalName string, dataType string, bitLength uint8, readerArgs ...WithReaderArgs) error {
	logicalName = sanitizeLogicalName(logicalName)
	if startElement.Name.Local != logicalName {
		return errors.Errorf("unexpected element '%s'. Expected '%s'", startElement.Name.Local, logicalName)
	} else if err := validateAttr(startElement.Attr, dataType, bitLength, readerArgs...); err != nil {
		return errors.Wrap(err, "Error validating Attributes")
	}
	return nil
}

func validateAttr(attr []xml.Attr, dataType string, bitLength uint8, readerArgs ...WithReaderArgs) error {
	dataTypeValidated := false
	bitLengthValidate := false
	for _, attribute := range attr {
		switch attribute.Name.Local {
		case "dataType":
			if attribute.Value != dataType {
				return errors.Errorf("Unexpected dataType :%s. Want %s", attribute.Value, dataType)
			}
			dataTypeValidated = true
		case "bitLength":
			if attribute.Value != fmt.Sprintf("%d", bitLength) {
				return errors.Errorf("Unexpected bitLength '%s'. Want '%d'", attribute.Value, bitLength)
			}
			bitLengthValidate = true
		}
	}
	if !dataTypeValidated {
		return errors.New("required attribute dataType missing")
	}
	if !bitLengthValidate {
		return errors.New("required attribute bitLength missing")
	}
	return nil
}
