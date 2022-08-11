/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package utils

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/big"
	"strings"
)

// NewXmlReadBuffer return as ReadBuffer which doesn't validate attributes and lists
func NewXmlReadBuffer(reader io.Reader) ReadBuffer {
	return &xmlReadBuffer{
		Decoder:        xml.NewDecoder(reader),
		pos:            1,
		doValidateList: false,
		doValidateAttr: false,
	}
}

// NewStrictXmlReadBuffer return as ReadBuffer which does validate attributes and lists depending on the setting
func NewStrictXmlReadBuffer(reader io.Reader, validateAttr bool, validateList bool) ReadBuffer {
	return &xmlReadBuffer{
		Decoder:        xml.NewDecoder(reader),
		pos:            1,
		doValidateAttr: validateAttr,
		doValidateList: validateList,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type xmlReadBuffer struct {
	bufferCommons
	*xml.Decoder
	pos            uint
	doValidateAttr bool
	doValidateList bool
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (x *xmlReadBuffer) GetPos() uint16 {
	return uint16(x.pos / 8)
}

func (x *xmlReadBuffer) Reset(pos uint16) {
	x.pos = uint(pos * 8)
}

func (x *xmlReadBuffer) HasMore(bitLength uint8) bool {
	// TODO: work with x.InputOffset() and check if we are at EOF
	return true
}

func (x *xmlReadBuffer) PullContext(logicalName string, readerArgs ...WithReaderArgs) error {
	startElement, err := x.travelToNextStartElement()
	if err != nil {
		return err
	}
	if startElement.Name.Local != logicalName {
		return errors.Errorf("Unexpected Start element '%s'. Expected '%s'", startElement.Name.Local, logicalName)
	}
	if err := x.validateIfList(readerArgs, startElement); err != nil {
		return err
	}
	return nil
}

func (x *xmlReadBuffer) ReadBit(logicalName string, readerArgs ...WithReaderArgs) (bool, error) {
	var value bool
	err := x.decode(logicalName, rwBitKey, 1, readerArgs, &value)
	if err != nil {
		return false, err
	}
	x.move(1)
	return value, nil
}

func (x *xmlReadBuffer) ReadByte(logicalName string, readerArgs ...WithReaderArgs) (byte, error) {
	var value string
	err := x.decode(logicalName, rwByteKey, 8, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	hexString := value
	if !strings.HasPrefix(hexString, "0x") {
		return 0, errors.Errorf("Hex string should start with 0x. Actual value %s", hexString)
	}
	hexString = strings.Replace(hexString, "0x", "", 1)
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return 0, err
	}
	x.move(8)
	return decoded[0], nil
}

func (x *xmlReadBuffer) ReadByteArray(logicalName string, numberOfBytes int, readerArgs ...WithReaderArgs) ([]byte, error) {
	var value string
	err := x.decode(logicalName, rwByteKey, uint(numberOfBytes/8), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	hexString := value
	if !strings.HasPrefix(hexString, "0x") {
		return nil, errors.Errorf("Hex string should start with 0x. Actual value %s", hexString)
	}
	hexString = strings.Replace(hexString, "0x", "", 1)
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	x.move(uint8(numberOfBytes / 8))
	return decoded, nil
}

func (x *xmlReadBuffer) ReadUint8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint8, error) {
	var value uint8
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint16, error) {
	var value uint16
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint32, error) {
	var value uint32
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint64, error) {
	var value uint64
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int8, error) {
	var value int8
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int16, error) {
	var value int16
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int32, error) {
	var value int32
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int64, error) {
	var value int64
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadBigInt(logicalName string, bitLength uint64, readerArgs ...WithReaderArgs) (*big.Int, error) {
	var value big.Int
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	x.move(uint8(bitLength))
	return &value, nil
}

func (x *xmlReadBuffer) ReadFloat32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float32, error) {
	var value float32
	err := x.decode(logicalName, rwFloatKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadFloat64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float64, error) {
	var value float64
	err := x.decode(logicalName, rwFloatKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadBigFloat(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (*big.Float, error) {
	var value big.Float
	err := x.decode(logicalName, rwFloatKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	x.move(bitLength)
	return &value, nil
}

func (x *xmlReadBuffer) ReadString(logicalName string, bitLength uint32, readerArgs ...WithReaderArgs) (string, error) {
	var value string
	// TODO: bitlength too short
	err := x.decode(logicalName, rwStringKey, uint(bitLength), readerArgs, &value)
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
	for {
		token, err := x.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch token.(type) {
		case xml.StartElement:
			return token.(xml.StartElement), nil
		case xml.EndElement:
			return xml.StartElement{}, errors.Errorf("unexpected end element %s", token.(xml.EndElement).Name)
		}
	}
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

func (x *xmlReadBuffer) decode(logicalName string, dataType string, bitLength uint, readerArgs []WithReaderArgs, targetValue interface{}) error {
	startElement, err := x.travelToNextStartElement()
	if err != nil {
		return err
	}
	err = x.validateStartElement(startElement, logicalName, dataType, bitLength, readerArgs...)
	if err != nil {
		return err
	}
	err = x.DecodeElement(targetValue, &startElement)
	if err != nil {
		return err
	}
	return nil
}

func (x *xmlReadBuffer) validateIfList(readerArgs []WithReaderArgs, startElement xml.StartElement) error {
	if !x.doValidateList {
		return nil
	}
	if x.isToBeRenderedAsList(upcastReaderArgs(readerArgs...)...) {
		for _, attr := range startElement.Attr {
			switch attr.Name.Local {
			case rwIsListKey:
				if attr.Value != "true" {
					return errors.Errorf("Startelement should be marked as %s=true", rwIsListKey)
				}
			}
		}
	}
	return nil
}

func (x *xmlReadBuffer) validateStartElement(startElement xml.StartElement, logicalName string, dataType string, bitLength uint, readerArgs ...WithReaderArgs) error {
	logicalName = x.sanitizeLogicalName(logicalName)
	if startElement.Name.Local != logicalName {
		return errors.Errorf("unexpected start element '%s'. Expected '%s'", startElement.Name.Local, logicalName)
	} else if err := x.validateAttr(startElement.Attr, dataType, bitLength, readerArgs...); err != nil {
		return errors.Wrap(err, "Error validating Attributes")
	}
	return nil
}

func (x *xmlReadBuffer) validateAttr(attr []xml.Attr, dataType string, bitLength uint, _ ...WithReaderArgs) error {
	if !x.doValidateAttr {
		return nil
	}
	dataTypeValidated := false
	bitLengthValidate := false
	for _, attribute := range attr {
		switch attribute.Name.Local {
		case rwDataTypeKey:
			if attribute.Value != dataType {
				return errors.Errorf("Unexpected %s :%s. Want %s", rwDataTypeKey, attribute.Value, dataType)
			}
			dataTypeValidated = true
		case rwBitLengthKey:
			if attribute.Value != fmt.Sprintf("%d", bitLength) {
				return errors.Errorf("Unexpected %s '%s'. Want '%d'", rwBitLengthKey, attribute.Value, bitLength)
			}
			bitLengthValidate = true
		}
	}
	if !dataTypeValidated {
		return errors.Errorf("required attribute %s missing", rwDataTypeKey)
	}
	if !bitLengthValidate {
		return errors.Errorf("required attribute %s missing", rwBitLengthKey)
	}
	return nil
}
