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
)

// The data-structure of this message
type AdsReadWriteRequest struct {
	IndexGroup  uint32
	IndexOffset uint32
	ReadLength  uint32
	Items       []*AdsMultiRequestItem
	Data        []int8
	Parent      *AdsData
	IAdsReadWriteRequest
}

// The corresponding interface
type IAdsReadWriteRequest interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *AdsReadWriteRequest) CommandId() CommandId {
	return CommandId_ADS_READ_WRITE
}

func (m *AdsReadWriteRequest) Response() bool {
	return false
}

func (m *AdsReadWriteRequest) InitializeParent(parent *AdsData) {
}

func NewAdsReadWriteRequest(indexGroup uint32, indexOffset uint32, readLength uint32, items []*AdsMultiRequestItem, data []int8) *AdsData {
	child := &AdsReadWriteRequest{
		IndexGroup:  indexGroup,
		IndexOffset: indexOffset,
		ReadLength:  readLength,
		Items:       items,
		Data:        data,
		Parent:      NewAdsData(),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastAdsReadWriteRequest(structType interface{}) *AdsReadWriteRequest {
	castFunc := func(typ interface{}) *AdsReadWriteRequest {
		if casted, ok := typ.(AdsReadWriteRequest); ok {
			return &casted
		}
		if casted, ok := typ.(*AdsReadWriteRequest); ok {
			return casted
		}
		if casted, ok := typ.(AdsData); ok {
			return CastAdsReadWriteRequest(casted.Child)
		}
		if casted, ok := typ.(*AdsData); ok {
			return CastAdsReadWriteRequest(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *AdsReadWriteRequest) GetTypeName() string {
	return "AdsReadWriteRequest"
}

func (m *AdsReadWriteRequest) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Simple field (indexGroup)
	lengthInBits += 32

	// Simple field (indexOffset)
	lengthInBits += 32

	// Simple field (readLength)
	lengthInBits += 32

	// Implicit Field (writeLength)
	lengthInBits += 32

	// Array field
	if len(m.Items) > 0 {
		for _, element := range m.Items {
			lengthInBits += element.LengthInBits()
		}
	}

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *AdsReadWriteRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func AdsReadWriteRequestParse(io *utils.ReadBuffer) (*AdsData, error) {

	// Simple Field (indexGroup)
	indexGroup, _indexGroupErr := io.ReadUint32(32)
	if _indexGroupErr != nil {
		return nil, errors.New("Error parsing 'indexGroup' field " + _indexGroupErr.Error())
	}

	// Simple Field (indexOffset)
	indexOffset, _indexOffsetErr := io.ReadUint32(32)
	if _indexOffsetErr != nil {
		return nil, errors.New("Error parsing 'indexOffset' field " + _indexOffsetErr.Error())
	}

	// Simple Field (readLength)
	readLength, _readLengthErr := io.ReadUint32(32)
	if _readLengthErr != nil {
		return nil, errors.New("Error parsing 'readLength' field " + _readLengthErr.Error())
	}

	// Implicit Field (writeLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	writeLength, _writeLengthErr := io.ReadUint32(32)
	if _writeLengthErr != nil {
		return nil, errors.New("Error parsing 'writeLength' field " + _writeLengthErr.Error())
	}

	// Array field (items)
	// Count array
	items := make([]*AdsMultiRequestItem, utils.InlineIf(bool(bool(bool(bool(bool((indexGroup) == (61568)))) || bool(bool(bool((indexGroup) == (61569))))) || bool(bool(bool((indexGroup) == (61570))))), uint16(indexOffset), uint16(uint16(0))))
	for curItem := uint16(0); curItem < uint16(utils.InlineIf(bool(bool(bool(bool(bool((indexGroup) == (61568)))) || bool(bool(bool((indexGroup) == (61569))))) || bool(bool(bool((indexGroup) == (61570))))), uint16(indexOffset), uint16(uint16(0)))); curItem++ {
		_item, _err := AdsMultiRequestItemParse(io, indexGroup)
		if _err != nil {
			return nil, errors.New("Error parsing 'items' field " + _err.Error())
		}
		items[curItem] = _item
	}

	// Array field (data)
	// Count array
	data := make([]int8, uint16(writeLength)-uint16(uint16(uint16(uint16(len(items)))*uint16(uint16(12)))))
	for curItem := uint16(0); curItem < uint16(uint16(writeLength)-uint16(uint16(uint16(uint16(len(items)))*uint16(uint16(12))))); curItem++ {
		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'data' field " + _err.Error())
		}
		data[curItem] = _item
	}

	// Create a partially initialized instance
	_child := &AdsReadWriteRequest{
		IndexGroup:  indexGroup,
		IndexOffset: indexOffset,
		ReadLength:  readLength,
		Items:       items,
		Data:        data,
		Parent:      &AdsData{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *AdsReadWriteRequest) Serialize(io utils.WriteBuffer) error {
	ser := func() error {

		// Simple Field (indexGroup)
		indexGroup := uint32(m.IndexGroup)
		_indexGroupErr := io.WriteUint32(32, (indexGroup))
		if _indexGroupErr != nil {
			return errors.New("Error serializing 'indexGroup' field " + _indexGroupErr.Error())
		}

		// Simple Field (indexOffset)
		indexOffset := uint32(m.IndexOffset)
		_indexOffsetErr := io.WriteUint32(32, (indexOffset))
		if _indexOffsetErr != nil {
			return errors.New("Error serializing 'indexOffset' field " + _indexOffsetErr.Error())
		}

		// Simple Field (readLength)
		readLength := uint32(m.ReadLength)
		_readLengthErr := io.WriteUint32(32, (readLength))
		if _readLengthErr != nil {
			return errors.New("Error serializing 'readLength' field " + _readLengthErr.Error())
		}

		// Implicit Field (writeLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		writeLength := uint32(uint32(uint32(uint32(uint32(len(m.Items)))*uint32(uint32(utils.InlineIf(bool(bool((m.IndexGroup) == (61570))), uint16(uint32(16)), uint16(uint32(12))))))) + uint32(uint32(len(m.Data))))
		_writeLengthErr := io.WriteUint32(32, (writeLength))
		if _writeLengthErr != nil {
			return errors.New("Error serializing 'writeLength' field " + _writeLengthErr.Error())
		}

		// Array Field (items)
		if m.Items != nil {
			for _, _element := range m.Items {
				_elementErr := _element.Serialize(io)
				if _elementErr != nil {
					return errors.New("Error serializing 'items' field " + _elementErr.Error())
				}
			}
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

		return nil
	}
	return m.Parent.SerializeParent(io, m, ser)
}

func (m *AdsReadWriteRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var token xml.Token
	var err error
	token = start
	for {
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			case "indexGroup":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.IndexGroup = data
			case "indexOffset":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.IndexOffset = data
			case "readLength":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ReadLength = data
			case "items":
				var _values []*AdsMultiRequestItem
				var dt *AdsMultiRequestItem
				if err := d.DecodeElement(&dt, &tok); err != nil {
					return err
				}
				_values = append(_values, dt)
				m.Items = _values
			case "data":
				var _encoded string
				if err := d.DecodeElement(&_encoded, &tok); err != nil {
					return err
				}
				_decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
				_len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
				if err != nil {
					return err
				}
				m.Data = utils.ByteArrayToInt8Array(_decoded[0:_len])
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

func (m *AdsReadWriteRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.IndexGroup, xml.StartElement{Name: xml.Name{Local: "indexGroup"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.IndexOffset, xml.StartElement{Name: xml.Name{Local: "indexOffset"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ReadLength, xml.StartElement{Name: xml.Name{Local: "readLength"}}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "items"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Items, xml.StartElement{Name: xml.Name{Local: "items"}}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "items"}}); err != nil {
		return err
	}
	_encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(m.Data)))
	base64.StdEncoding.Encode(_encodedData, utils.Int8ArrayToByteArray(m.Data))
	if err := e.EncodeElement(_encodedData, xml.StartElement{Name: xml.Name{Local: "data"}}); err != nil {
		return err
	}
	return nil
}
