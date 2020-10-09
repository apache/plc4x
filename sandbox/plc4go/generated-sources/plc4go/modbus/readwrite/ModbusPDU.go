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
package readwrite

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
)

type ModbusPDU struct {
}

type IModbusPDU interface {
	spi.Message
	ErrorFlag() bool
	FunctionFlag() uint8
	Response() bool
	Serialize(io spi.WriteBuffer)
}

type ModbusPDUInitializer interface {
	initialize() spi.Message
}

// Dummy implementation ...
func (m ModbusPDU) ErrorFlag() bool {
	return false
}

func ModbusPDUErrorFlag(m IModbusPDU) bool {
	return m.ErrorFlag()
}

// Dummy implementation ...
func (m ModbusPDU) FunctionFlag() uint8 {
	return 0
}

func ModbusPDUFunctionFlag(m IModbusPDU) uint8 {
	return m.FunctionFlag()
}

// Dummy implementation ...
func (m ModbusPDU) Response() bool {
	return false
}

func ModbusPDUResponse(m IModbusPDU) bool {
	return m.Response()
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

func ModbusPDUParse(io spi.ReadBuffer, response bool) (spi.Message, error) {

	// Discriminator Field (errorFlag) (Used as input to a switch field)
	var errorFlag bool = io.ReadBit()

	// Discriminator Field (functionFlag) (Used as input to a switch field)
	var functionFlag uint8 = io.ReadUint8(7)

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

func (m ModbusPDU) Serialize(io spi.WriteBuffer) {

	// Discriminator Field (errorFlag) (Used as input to a switch field)
	errorFlag := ModbusPDUErrorFlag(m)
	io.WriteBit((bool)(errorFlag))

	// Discriminator Field (functionFlag) (Used as input to a switch field)
	functionFlag := ModbusPDUFunctionFlag(m)
	io.WriteUint8(7, (functionFlag))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	switch m.(type) {
	case IModbusPDUError:
		IModbusPDUError(m).Serialize(io)
	case IModbusPDUReadDiscreteInputsRequest:
		IModbusPDUReadDiscreteInputsRequest(m).Serialize(io)
	case IModbusPDUReadDiscreteInputsResponse:
		IModbusPDUReadDiscreteInputsResponse(m).Serialize(io)
	case IModbusPDUReadCoilsRequest:
		IModbusPDUReadCoilsRequest(m).Serialize(io)
	case IModbusPDUReadCoilsResponse:
		IModbusPDUReadCoilsResponse(m).Serialize(io)
	case IModbusPDUWriteSingleCoilRequest:
		IModbusPDUWriteSingleCoilRequest(m).Serialize(io)
	case IModbusPDUWriteSingleCoilResponse:
		IModbusPDUWriteSingleCoilResponse(m).Serialize(io)
	case IModbusPDUWriteMultipleCoilsRequest:
		IModbusPDUWriteMultipleCoilsRequest(m).Serialize(io)
	case IModbusPDUWriteMultipleCoilsResponse:
		IModbusPDUWriteMultipleCoilsResponse(m).Serialize(io)
	case IModbusPDUReadInputRegistersRequest:
		IModbusPDUReadInputRegistersRequest(m).Serialize(io)
	case IModbusPDUReadInputRegistersResponse:
		IModbusPDUReadInputRegistersResponse(m).Serialize(io)
	case IModbusPDUReadHoldingRegistersRequest:
		IModbusPDUReadHoldingRegistersRequest(m).Serialize(io)
	case IModbusPDUReadHoldingRegistersResponse:
		IModbusPDUReadHoldingRegistersResponse(m).Serialize(io)
	case IModbusPDUWriteSingleRegisterRequest:
		IModbusPDUWriteSingleRegisterRequest(m).Serialize(io)
	case IModbusPDUWriteSingleRegisterResponse:
		IModbusPDUWriteSingleRegisterResponse(m).Serialize(io)
	case IModbusPDUWriteMultipleHoldingRegistersRequest:
		IModbusPDUWriteMultipleHoldingRegistersRequest(m).Serialize(io)
	case IModbusPDUWriteMultipleHoldingRegistersResponse:
		IModbusPDUWriteMultipleHoldingRegistersResponse(m).Serialize(io)
	case IModbusPDUReadWriteMultipleHoldingRegistersRequest:
		IModbusPDUReadWriteMultipleHoldingRegistersRequest(m).Serialize(io)
	case IModbusPDUReadWriteMultipleHoldingRegistersResponse:
		IModbusPDUReadWriteMultipleHoldingRegistersResponse(m).Serialize(io)
	case IModbusPDUMaskWriteHoldingRegisterRequest:
		IModbusPDUMaskWriteHoldingRegisterRequest(m).Serialize(io)
	case IModbusPDUMaskWriteHoldingRegisterResponse:
		IModbusPDUMaskWriteHoldingRegisterResponse(m).Serialize(io)
	case IModbusPDUReadFifoQueueRequest:
		IModbusPDUReadFifoQueueRequest(m).Serialize(io)
	case IModbusPDUReadFifoQueueResponse:
		IModbusPDUReadFifoQueueResponse(m).Serialize(io)
	case IModbusPDUReadFileRecordRequest:
		IModbusPDUReadFileRecordRequest(m).Serialize(io)
	case IModbusPDUReadFileRecordResponse:
		IModbusPDUReadFileRecordResponse(m).Serialize(io)
	case IModbusPDUWriteFileRecordRequest:
		IModbusPDUWriteFileRecordRequest(m).Serialize(io)
	case IModbusPDUWriteFileRecordResponse:
		IModbusPDUWriteFileRecordResponse(m).Serialize(io)
	case IModbusPDUReadExceptionStatusRequest:
		IModbusPDUReadExceptionStatusRequest(m).Serialize(io)
	case IModbusPDUReadExceptionStatusResponse:
		IModbusPDUReadExceptionStatusResponse(m).Serialize(io)
	case IModbusPDUDiagnosticRequest:
		IModbusPDUDiagnosticRequest(m).Serialize(io)
	case IModbusPDUGetComEventLogRequest:
		IModbusPDUGetComEventLogRequest(m).Serialize(io)
	case IModbusPDUGetComEventLogResponse:
		IModbusPDUGetComEventLogResponse(m).Serialize(io)
	case IModbusPDUReportServerIdRequest:
		IModbusPDUReportServerIdRequest(m).Serialize(io)
	case IModbusPDUReportServerIdResponse:
		IModbusPDUReportServerIdResponse(m).Serialize(io)
	case IModbusPDUReadDeviceIdentificationRequest:
		IModbusPDUReadDeviceIdentificationRequest(m).Serialize(io)
	case IModbusPDUReadDeviceIdentificationResponse:
		IModbusPDUReadDeviceIdentificationResponse(m).Serialize(io)
	}
}
