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
)

// The data-structure of this message
type S7VarPayloadStatusItem struct {
	returnCode DataTransportErrorCode
}

// The corresponding interface
type IS7VarPayloadStatusItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewS7VarPayloadStatusItem(returnCode DataTransportErrorCode) spi.Message {
	return &S7VarPayloadStatusItem{returnCode: returnCode}
}

func (m S7VarPayloadStatusItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Enum Field (returnCode)
	lengthInBits += 8

	return lengthInBits
}

func (m S7VarPayloadStatusItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7VarPayloadStatusItemParse(io spi.ReadBuffer) (spi.Message, error) {

	// Enum field (returnCode)
	returnCode, _returnCodeErr := DataTransportErrorCodeParse(io)
	if _returnCodeErr != nil {
		return nil, errors.New("Error parsing 'returnCode' field " + _returnCodeErr.Error())
	}

	// Create the instance
	return NewS7VarPayloadStatusItem(returnCode), nil
}

func (m S7VarPayloadStatusItem) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IS7VarPayloadStatusItem); ok {

			// Enum field (returnCode)
			returnCode := m.returnCode
			returnCode.Serialize(io)
		}
	}
	serializeFunc(m)
}
