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
	log "github.com/sirupsen/logrus"
	"math"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
)

type ModbusPDU struct {
}

type ModbusPDUInitializer interface {
	initialize() spi.Message
}

func (m ModbusPDU) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (error)
	lengthInBits += 1

	// Discriminator Field (function)
	lengthInBits += 7

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m ModbusPDU) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUParse(io spi.ReadBuffer, response bool) spi.Message {
	var startPos = io.GetPos()
	var curPos uint16

	// Discriminator Field (error) (Used as input to a switch field)
	var error bool = io.ReadBit()

	// Discriminator Field (function) (Used as input to a switch field)
	var function uint8 = io.ReadUint8(7)

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer ModbusPDUInitializer
	switch {
	case error == true:
		initializer = ModbusPDUErrorParse(io)
	case error == false && function == 0x02 && response == false:
		initializer = ModbusPDUReadDiscreteInputsRequestParse(io)
	case error == false && function == 0x02 && response == true:
		initializer = ModbusPDUReadDiscreteInputsResponseParse(io)
	case error == false && function == 0x01 && response == false:
		initializer = ModbusPDUReadCoilsRequestParse(io)
	case error == false && function == 0x01 && response == true:
		initializer = ModbusPDUReadCoilsResponseParse(io)
	case error == false && function == 0x05 && response == false:
		initializer = ModbusPDUWriteSingleCoilRequestParse(io)
	case error == false && function == 0x05 && response == true:
		initializer = ModbusPDUWriteSingleCoilResponseParse(io)
	case error == false && function == 0x0F && response == false:
		initializer = ModbusPDUWriteMultipleCoilsRequestParse(io)
	case error == false && function == 0x0F && response == true:
		initializer = ModbusPDUWriteMultipleCoilsResponseParse(io)
	case error == false && function == 0x04 && response == false:
		initializer = ModbusPDUReadInputRegistersRequestParse(io)
	case error == false && function == 0x04 && response == true:
		initializer = ModbusPDUReadInputRegistersResponseParse(io)
	case error == false && function == 0x03 && response == false:
		initializer = ModbusPDUReadHoldingRegistersRequestParse(io)
	case error == false && function == 0x03 && response == true:
		initializer = ModbusPDUReadHoldingRegistersResponseParse(io)
	case error == false && function == 0x06 && response == false:
		initializer = ModbusPDUWriteSingleRegisterRequestParse(io)
	case error == false && function == 0x06 && response == true:
		initializer = ModbusPDUWriteSingleRegisterResponseParse(io)
	case error == false && function == 0x10 && response == false:
		initializer = ModbusPDUWriteMultipleHoldingRegistersRequestParse(io)
	case error == false && function == 0x10 && response == true:
		initializer = ModbusPDUWriteMultipleHoldingRegistersResponseParse(io)
	case error == false && function == 0x17 && response == false:
		initializer = ModbusPDUReadWriteMultipleHoldingRegistersRequestParse(io)
	case error == false && function == 0x17 && response == true:
		initializer = ModbusPDUReadWriteMultipleHoldingRegistersResponseParse(io)
	case error == false && function == 0x16 && response == false:
		initializer = ModbusPDUMaskWriteHoldingRegisterRequestParse(io)
	case error == false && function == 0x16 && response == true:
		initializer = ModbusPDUMaskWriteHoldingRegisterResponseParse(io)
	case error == false && function == 0x18 && response == false:
		initializer = ModbusPDUReadFifoQueueRequestParse(io)
	case error == false && function == 0x18 && response == true:
		initializer = ModbusPDUReadFifoQueueResponseParse(io)
	case error == false && function == 0x14 && response == false:
		initializer = ModbusPDUReadFileRecordRequestParse(io)
	case error == false && function == 0x14 && response == true:
		initializer = ModbusPDUReadFileRecordResponseParse(io)
	case error == false && function == 0x15 && response == false:
		initializer = ModbusPDUWriteFileRecordRequestParse(io)
	case error == false && function == 0x15 && response == true:
		initializer = ModbusPDUWriteFileRecordResponseParse(io)
	case error == false && function == 0x07 && response == false:
		initializer = ModbusPDUReadExceptionStatusRequestParse(io)
	case error == false && function == 0x07 && response == true:
		initializer = ModbusPDUReadExceptionStatusResponseParse(io)
	case error == false && function == 0x08 && response == false:
		initializer = ModbusPDUDiagnosticRequestParse(io)
	case error == false && function == 0x0C && response == false:
		initializer = ModbusPDUGetComEventLogRequestParse(io)
	case error == false && function == 0x0C && response == true:
		initializer = ModbusPDUGetComEventLogResponseParse(io)
	case error == false && function == 0x11 && response == false:
		initializer = ModbusPDUReportServerIdRequestParse(io)
	case error == false && function == 0x11 && response == true:
		initializer = ModbusPDUReportServerIdResponseParse(io)
	case error == false && function == 0x2B && response == false:
		initializer = ModbusPDUReadDeviceIdentificationRequestParse(io)
	case error == false && function == 0x2B && response == true:
		initializer = ModbusPDUReadDeviceIdentificationResponseParse(io)
	}

	// Create the instance
	return initializer.initialize()
}
