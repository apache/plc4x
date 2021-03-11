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
	log "github.com/sirupsen/logrus"
	"io"
)

// The data-structure of this message
type APDUAbort struct {
	Server           bool
	OriginalInvokeId uint8
	AbortReason      uint8
	Parent           *APDU
	IAPDUAbort
}

// The corresponding interface
type IAPDUAbort interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *APDUAbort) ApduType() uint8 {
	return 0x7
}

func (m *APDUAbort) InitializeParent(parent *APDU) {
}

func NewAPDUAbort(server bool, originalInvokeId uint8, abortReason uint8) *APDU {
	child := &APDUAbort{
		Server:           server,
		OriginalInvokeId: originalInvokeId,
		AbortReason:      abortReason,
		Parent:           NewAPDU(),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastAPDUAbort(structType interface{}) *APDUAbort {
	castFunc := func(typ interface{}) *APDUAbort {
		if casted, ok := typ.(APDUAbort); ok {
			return &casted
		}
		if casted, ok := typ.(*APDUAbort); ok {
			return casted
		}
		if casted, ok := typ.(APDU); ok {
			return CastAPDUAbort(casted.Child)
		}
		if casted, ok := typ.(*APDU); ok {
			return CastAPDUAbort(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *APDUAbort) GetTypeName() string {
	return "APDUAbort"
}

func (m *APDUAbort) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Reserved Field (reserved)
	lengthInBits += 3

	// Simple field (server)
	lengthInBits += 1

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (abortReason)
	lengthInBits += 8

	return lengthInBits
}

func (m *APDUAbort) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUAbortParse(io *utils.ReadBuffer) (*APDU, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(3)
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

	// Simple Field (server)
	server, _serverErr := io.ReadBit()
	if _serverErr != nil {
		return nil, errors.New("Error parsing 'server' field " + _serverErr.Error())
	}

	// Simple Field (originalInvokeId)
	originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
	if _originalInvokeIdErr != nil {
		return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
	}

	// Simple Field (abortReason)
	abortReason, _abortReasonErr := io.ReadUint8(8)
	if _abortReasonErr != nil {
		return nil, errors.New("Error parsing 'abortReason' field " + _abortReasonErr.Error())
	}

	// Create a partially initialized instance
	_child := &APDUAbort{
		Server:           server,
		OriginalInvokeId: originalInvokeId,
		AbortReason:      abortReason,
		Parent:           &APDU{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *APDUAbort) Serialize(io utils.WriteBuffer) error {
	ser := func() error {

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(3, uint8(0x00))
			if _err != nil {
				return errors.New("Error serializing 'reserved' field " + _err.Error())
			}
		}

		// Simple Field (server)
		server := bool(m.Server)
		_serverErr := io.WriteBit((server))
		if _serverErr != nil {
			return errors.New("Error serializing 'server' field " + _serverErr.Error())
		}

		// Simple Field (originalInvokeId)
		originalInvokeId := uint8(m.OriginalInvokeId)
		_originalInvokeIdErr := io.WriteUint8(8, (originalInvokeId))
		if _originalInvokeIdErr != nil {
			return errors.New("Error serializing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
		}

		// Simple Field (abortReason)
		abortReason := uint8(m.AbortReason)
		_abortReasonErr := io.WriteUint8(8, (abortReason))
		if _abortReasonErr != nil {
			return errors.New("Error serializing 'abortReason' field " + _abortReasonErr.Error())
		}

		return nil
	}
	return m.Parent.SerializeParent(io, m, ser)
}

func (m *APDUAbort) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var token xml.Token
	var err error
	token = start
	for {
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			case "server":
				var data bool
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.Server = data
			case "originalInvokeId":
				var data uint8
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.OriginalInvokeId = data
			case "abortReason":
				var data uint8
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.AbortReason = data
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

func (m *APDUAbort) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.Server, xml.StartElement{Name: xml.Name{Local: "server"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.OriginalInvokeId, xml.StartElement{Name: xml.Name{Local: "originalInvokeId"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.AbortReason, xml.StartElement{Name: xml.Name{Local: "abortReason"}}); err != nil {
		return err
	}
	return nil
}
