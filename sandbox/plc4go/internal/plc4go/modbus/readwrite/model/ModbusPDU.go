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
	"io"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDU struct {
}

// The corresponding interface
type IModbusPDU interface {
	spi.Message
	ErrorFlag() bool
	FunctionFlag() uint8
	Response() bool
	Serialize(io utils.WriteBuffer) error
}

type ModbusPDUInitializer interface {
	initialize() spi.Message
}

func ModbusPDUErrorFlag(m IModbusPDU) bool {
	return m.ErrorFlag()
}

func ModbusPDUFunctionFlag(m IModbusPDU) uint8 {
	return m.FunctionFlag()
}

func ModbusPDUResponse(m IModbusPDU) bool {
	return m.Response()
}

func CastIModbusPDU(structType interface{}) IModbusPDU {
	castFunc := func(typ interface{}) IModbusPDU {
		if iModbusPDU, ok := typ.(IModbusPDU); ok {
			return iModbusPDU
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDU(structType interface{}) ModbusPDU {
	castFunc := func(typ interface{}) ModbusPDU {
		if sModbusPDU, ok := typ.(ModbusPDU); ok {
			return sModbusPDU
		}
		if sModbusPDU, ok := typ.(*ModbusPDU); ok {
			return *sModbusPDU
		}
		return ModbusPDU{}
	}
	return castFunc(structType)
}

func (m ModbusPDU) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (errorFlag)
	lengthInBits += 1

	// Discriminator Field (functionFlag)
	lengthInBits += 7

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m ModbusPDU) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUParse(io *utils.ReadBuffer, response bool) (spi.Message, error) {

	// Discriminator Field (errorFlag) (Used as input to a switch field)
	errorFlag, _errorFlagErr := io.ReadBit()
	if _errorFlagErr != nil {
		return nil, errors.New("Error parsing 'errorFlag' field " + _errorFlagErr.Error())
	}

	// Discriminator Field (functionFlag) (Used as input to a switch field)
	functionFlag, _functionFlagErr := io.ReadUint8(7)
	if _functionFlagErr != nil {
		return nil, errors.New("Error parsing 'functionFlag' field " + _functionFlagErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer ModbusPDUInitializer
	var typeSwitchError error
	switch {
	case errorFlag == true:
		initializer, typeSwitchError = ModbusPDUErrorParse(io)
	case errorFlag == false && functionFlag == 0x02 && response == false:
		initializer, typeSwitchError = ModbusPDUReadDiscreteInputsRequestParse(io)
	case errorFlag == false && functionFlag == 0x02 && response == true:
		initializer, typeSwitchError = ModbusPDUReadDiscreteInputsResponseParse(io)
	case errorFlag == false && functionFlag == 0x01 && response == false:
		initializer, typeSwitchError = ModbusPDUReadCoilsRequestParse(io)
	case errorFlag == false && functionFlag == 0x01 && response == true:
		initializer, typeSwitchError = ModbusPDUReadCoilsResponseParse(io)
	case errorFlag == false && functionFlag == 0x05 && response == false:
		initializer, typeSwitchError = ModbusPDUWriteSingleCoilRequestParse(io)
	case errorFlag == false && functionFlag == 0x05 && response == true:
		initializer, typeSwitchError = ModbusPDUWriteSingleCoilResponseParse(io)
	case errorFlag == false && functionFlag == 0x0F && response == false:
		initializer, typeSwitchError = ModbusPDUWriteMultipleCoilsRequestParse(io)
	case errorFlag == false && functionFlag == 0x0F && response == true:
		initializer, typeSwitchError = ModbusPDUWriteMultipleCoilsResponseParse(io)
	case errorFlag == false && functionFlag == 0x04 && response == false:
		initializer, typeSwitchError = ModbusPDUReadInputRegistersRequestParse(io)
	case errorFlag == false && functionFlag == 0x04 && response == true:
		initializer, typeSwitchError = ModbusPDUReadInputRegistersResponseParse(io)
	case errorFlag == false && functionFlag == 0x03 && response == false:
		initializer, typeSwitchError = ModbusPDUReadHoldingRegistersRequestParse(io)
	case errorFlag == false && functionFlag == 0x03 && response == true:
		initializer, typeSwitchError = ModbusPDUReadHoldingRegistersResponseParse(io)
	case errorFlag == false && functionFlag == 0x06 && response == false:
		initializer, typeSwitchError = ModbusPDUWriteSingleRegisterRequestParse(io)
	case errorFlag == false && functionFlag == 0x06 && response == true:
		initializer, typeSwitchError = ModbusPDUWriteSingleRegisterResponseParse(io)
	case errorFlag == false && functionFlag == 0x10 && response == false:
		initializer, typeSwitchError = ModbusPDUWriteMultipleHoldingRegistersRequestParse(io)
	case errorFlag == false && functionFlag == 0x10 && response == true:
		initializer, typeSwitchError = ModbusPDUWriteMultipleHoldingRegistersResponseParse(io)
	case errorFlag == false && functionFlag == 0x17 && response == false:
		initializer, typeSwitchError = ModbusPDUReadWriteMultipleHoldingRegistersRequestParse(io)
	case errorFlag == false && functionFlag == 0x17 && response == true:
		initializer, typeSwitchError = ModbusPDUReadWriteMultipleHoldingRegistersResponseParse(io)
	case errorFlag == false && functionFlag == 0x16 && response == false:
		initializer, typeSwitchError = ModbusPDUMaskWriteHoldingRegisterRequestParse(io)
	case errorFlag == false && functionFlag == 0x16 && response == true:
		initializer, typeSwitchError = ModbusPDUMaskWriteHoldingRegisterResponseParse(io)
	case errorFlag == false && functionFlag == 0x18 && response == false:
		initializer, typeSwitchError = ModbusPDUReadFifoQueueRequestParse(io)
	case errorFlag == false && functionFlag == 0x18 && response == true:
		initializer, typeSwitchError = ModbusPDUReadFifoQueueResponseParse(io)
	case errorFlag == false && functionFlag == 0x14 && response == false:
		initializer, typeSwitchError = ModbusPDUReadFileRecordRequestParse(io)
	case errorFlag == false && functionFlag == 0x14 && response == true:
		initializer, typeSwitchError = ModbusPDUReadFileRecordResponseParse(io)
	case errorFlag == false && functionFlag == 0x15 && response == false:
		initializer, typeSwitchError = ModbusPDUWriteFileRecordRequestParse(io)
	case errorFlag == false && functionFlag == 0x15 && response == true:
		initializer, typeSwitchError = ModbusPDUWriteFileRecordResponseParse(io)
	case errorFlag == false && functionFlag == 0x07 && response == false:
		initializer, typeSwitchError = ModbusPDUReadExceptionStatusRequestParse(io)
	case errorFlag == false && functionFlag == 0x07 && response == true:
		initializer, typeSwitchError = ModbusPDUReadExceptionStatusResponseParse(io)
	case errorFlag == false && functionFlag == 0x08 && response == false:
		initializer, typeSwitchError = ModbusPDUDiagnosticRequestParse(io)
	case errorFlag == false && functionFlag == 0x08 && response == true:
		initializer, typeSwitchError = ModbusPDUDiagnosticResponseParse(io)
	case errorFlag == false && functionFlag == 0x0B && response == false:
		initializer, typeSwitchError = ModbusPDUGetComEventCounterRequestParse(io)
	case errorFlag == false && functionFlag == 0x0B && response == true:
		initializer, typeSwitchError = ModbusPDUGetComEventCounterResponseParse(io)
	case errorFlag == false && functionFlag == 0x0C && response == false:
		initializer, typeSwitchError = ModbusPDUGetComEventLogRequestParse(io)
	case errorFlag == false && functionFlag == 0x0C && response == true:
		initializer, typeSwitchError = ModbusPDUGetComEventLogResponseParse(io)
	case errorFlag == false && functionFlag == 0x11 && response == false:
		initializer, typeSwitchError = ModbusPDUReportServerIdRequestParse(io)
	case errorFlag == false && functionFlag == 0x11 && response == true:
		initializer, typeSwitchError = ModbusPDUReportServerIdResponseParse(io)
	case errorFlag == false && functionFlag == 0x2B && response == false:
		initializer, typeSwitchError = ModbusPDUReadDeviceIdentificationRequestParse(io)
	case errorFlag == false && functionFlag == 0x2B && response == true:
		initializer, typeSwitchError = ModbusPDUReadDeviceIdentificationResponseParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func ModbusPDUSerialize(io utils.WriteBuffer, m ModbusPDU, i IModbusPDU, childSerialize func() error) error {

	// Discriminator Field (errorFlag) (Used as input to a switch field)
	errorFlag := bool(i.ErrorFlag())
	_errorFlagErr := io.WriteBit((errorFlag))
	if _errorFlagErr != nil {
		return errors.New("Error serializing 'errorFlag' field " + _errorFlagErr.Error())
	}

	// Discriminator Field (functionFlag) (Used as input to a switch field)
	functionFlag := uint8(i.FunctionFlag())
	_functionFlagErr := io.WriteUint8(7, (functionFlag))
	if _functionFlagErr != nil {
		return errors.New("Error serializing 'functionFlag' field " + _functionFlagErr.Error())
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	_typeSwitchErr := childSerialize()
	if _typeSwitchErr != nil {
		return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
	}

	return nil
}

func (m *ModbusPDU) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			}
		}
	}
}

func (m ModbusPDU) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
		{Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDU"},
	}}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}
