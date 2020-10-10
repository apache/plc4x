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
type TunnelingRequest struct {
	tunnelingRequestDataBlock TunnelingRequestDataBlock
	cemi                      CEMI
	KNXNetIPMessage
}

// The corresponding interface
type ITunnelingRequest interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m TunnelingRequest) MsgType() uint16 {
	return 0x0420
}

func (m TunnelingRequest) initialize() spi.Message {
	return m
}

func NewTunnelingRequest(tunnelingRequestDataBlock TunnelingRequestDataBlock, cemi CEMI) KNXNetIPMessageInitializer {
	return &TunnelingRequest{tunnelingRequestDataBlock: tunnelingRequestDataBlock, cemi: cemi}
}

func (m TunnelingRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (tunnelingRequestDataBlock)
	lengthInBits += m.tunnelingRequestDataBlock.LengthInBits()

	// Simple field (cemi)
	lengthInBits += m.cemi.LengthInBits()

	return lengthInBits
}

func (m TunnelingRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func TunnelingRequestParse(io spi.ReadBuffer, totalLength uint16) (KNXNetIPMessageInitializer, error) {

	// Simple Field (tunnelingRequestDataBlock)
	_tunnelingRequestDataBlockMessage, _err := TunnelingRequestDataBlockParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'tunnelingRequestDataBlock'. " + _err.Error())
	}
	var tunnelingRequestDataBlock TunnelingRequestDataBlock
	tunnelingRequestDataBlock, _tunnelingRequestDataBlockOk := _tunnelingRequestDataBlockMessage.(TunnelingRequestDataBlock)
	if !_tunnelingRequestDataBlockOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_tunnelingRequestDataBlockMessage).Name() + " to TunnelingRequestDataBlock")
	}

	// Simple Field (cemi)
	_cemiMessage, _err := CEMIParse(io, uint8((totalLength)-((6)+(tunnelingRequestDataBlock.LengthInBytes()))))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'cemi'. " + _err.Error())
	}
	var cemi CEMI
	cemi, _cemiOk := _cemiMessage.(CEMI)
	if !_cemiOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_cemiMessage).Name() + " to CEMI")
	}

	// Create the instance
	return NewTunnelingRequest(tunnelingRequestDataBlock, cemi), nil
}

func (m TunnelingRequest) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ITunnelingRequest); ok {

			// Simple Field (tunnelingRequestDataBlock)
			var tunnelingRequestDataBlock TunnelingRequestDataBlock = m.tunnelingRequestDataBlock
			tunnelingRequestDataBlock.Serialize(io)

			// Simple Field (cemi)
			var cemi CEMI = m.cemi
			cemi.Serialize(io)
		}
	}
	serializeFunc(m)
}
