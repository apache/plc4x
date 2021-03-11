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
	"encoding/xml"
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"io"
)

// The data-structure of this message
type AdsMultiRequestItemWrite struct {
	ItemIndexGroup  uint32
	ItemIndexOffset uint32
	ItemWriteLength uint32
	Parent          *AdsMultiRequestItem
	IAdsMultiRequestItemWrite
}

// The corresponding interface
type IAdsMultiRequestItemWrite interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *AdsMultiRequestItemWrite) IndexGroup() uint32 {
	return 61569
}

func (m *AdsMultiRequestItemWrite) InitializeParent(parent *AdsMultiRequestItem) {
}

func NewAdsMultiRequestItemWrite(itemIndexGroup uint32, itemIndexOffset uint32, itemWriteLength uint32) *AdsMultiRequestItem {
	child := &AdsMultiRequestItemWrite{
		ItemIndexGroup:  itemIndexGroup,
		ItemIndexOffset: itemIndexOffset,
		ItemWriteLength: itemWriteLength,
		Parent:          NewAdsMultiRequestItem(),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastAdsMultiRequestItemWrite(structType interface{}) *AdsMultiRequestItemWrite {
	castFunc := func(typ interface{}) *AdsMultiRequestItemWrite {
		if casted, ok := typ.(AdsMultiRequestItemWrite); ok {
			return &casted
		}
		if casted, ok := typ.(*AdsMultiRequestItemWrite); ok {
			return casted
		}
		if casted, ok := typ.(AdsMultiRequestItem); ok {
			return CastAdsMultiRequestItemWrite(casted.Child)
		}
		if casted, ok := typ.(*AdsMultiRequestItem); ok {
			return CastAdsMultiRequestItemWrite(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *AdsMultiRequestItemWrite) GetTypeName() string {
	return "AdsMultiRequestItemWrite"
}

func (m *AdsMultiRequestItemWrite) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Simple field (itemIndexGroup)
	lengthInBits += 32

	// Simple field (itemIndexOffset)
	lengthInBits += 32

	// Simple field (itemWriteLength)
	lengthInBits += 32

	return lengthInBits
}

func (m *AdsMultiRequestItemWrite) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func AdsMultiRequestItemWriteParse(io *utils.ReadBuffer) (*AdsMultiRequestItem, error) {

	// Simple Field (itemIndexGroup)
	itemIndexGroup, _itemIndexGroupErr := io.ReadUint32(32)
	if _itemIndexGroupErr != nil {
		return nil, errors.New("Error parsing 'itemIndexGroup' field " + _itemIndexGroupErr.Error())
	}

	// Simple Field (itemIndexOffset)
	itemIndexOffset, _itemIndexOffsetErr := io.ReadUint32(32)
	if _itemIndexOffsetErr != nil {
		return nil, errors.New("Error parsing 'itemIndexOffset' field " + _itemIndexOffsetErr.Error())
	}

	// Simple Field (itemWriteLength)
	itemWriteLength, _itemWriteLengthErr := io.ReadUint32(32)
	if _itemWriteLengthErr != nil {
		return nil, errors.New("Error parsing 'itemWriteLength' field " + _itemWriteLengthErr.Error())
	}

	// Create a partially initialized instance
	_child := &AdsMultiRequestItemWrite{
		ItemIndexGroup:  itemIndexGroup,
		ItemIndexOffset: itemIndexOffset,
		ItemWriteLength: itemWriteLength,
		Parent:          &AdsMultiRequestItem{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *AdsMultiRequestItemWrite) Serialize(io utils.WriteBuffer) error {
	ser := func() error {

		// Simple Field (itemIndexGroup)
		itemIndexGroup := uint32(m.ItemIndexGroup)
		_itemIndexGroupErr := io.WriteUint32(32, (itemIndexGroup))
		if _itemIndexGroupErr != nil {
			return errors.New("Error serializing 'itemIndexGroup' field " + _itemIndexGroupErr.Error())
		}

		// Simple Field (itemIndexOffset)
		itemIndexOffset := uint32(m.ItemIndexOffset)
		_itemIndexOffsetErr := io.WriteUint32(32, (itemIndexOffset))
		if _itemIndexOffsetErr != nil {
			return errors.New("Error serializing 'itemIndexOffset' field " + _itemIndexOffsetErr.Error())
		}

		// Simple Field (itemWriteLength)
		itemWriteLength := uint32(m.ItemWriteLength)
		_itemWriteLengthErr := io.WriteUint32(32, (itemWriteLength))
		if _itemWriteLengthErr != nil {
			return errors.New("Error serializing 'itemWriteLength' field " + _itemWriteLengthErr.Error())
		}

		return nil
	}
	return m.Parent.SerializeParent(io, m, ser)
}

func (m *AdsMultiRequestItemWrite) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var token xml.Token
	var err error
	token = start
	for {
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			case "itemIndexGroup":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ItemIndexGroup = data
			case "itemIndexOffset":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ItemIndexOffset = data
			case "itemWriteLength":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ItemWriteLength = data
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

func (m *AdsMultiRequestItemWrite) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.ItemIndexGroup, xml.StartElement{Name: xml.Name{Local: "itemIndexGroup"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ItemIndexOffset, xml.StartElement{Name: xml.Name{Local: "itemIndexOffset"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ItemWriteLength, xml.StartElement{Name: xml.Name{Local: "itemWriteLength"}}); err != nil {
		return err
	}
	return nil
}
