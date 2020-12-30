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
package modbus

import (
	"errors"
	"fmt"
	modbusModel "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"math"
	"sync/atomic"
)

type ModbusWriter struct {
	transactionIdentifier int32
	unitIdentifier        uint8
	messageCodec          spi.MessageCodec
	spi.PlcWriter
}

func NewModbusWriter(unitIdentifier uint8, messageCodec spi.MessageCodec) ModbusWriter {
	return ModbusWriter{
		transactionIdentifier: 0,
		unitIdentifier:        unitIdentifier,
		messageCodec:          messageCodec,
	}
}

func (m ModbusWriter) Write(writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	result := make(chan model.PlcWriteRequestResult)
	// If we are requesting only one field, use a
	if len(writeRequest.GetFieldNames()) == 1 {
		fieldName := writeRequest.GetFieldNames()[0]

		// Get the modbus field instance from the request
		field := writeRequest.GetField(fieldName)
		modbusField, err := CastToModbusFieldFromPlcField(field)
		if err != nil {
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("invalid field item type"),
			}
			return result
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(fieldName)
		io := utils.NewWriteBuffer()
		if err := modbusModel.DataItemSerialize(io, value, modbusField.Datatype, modbusField.Quantity); err != nil {
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("error serializing value: " + err.Error()),
			}
			return result
		}
		data := utils.Uint8ArrayToInt8Array(io.GetBytes())

		// Calculate the number of words needed to send the data
		numWords := uint16(math.Ceil(float64(len(data)) / 2))

		var pdu *modbusModel.ModbusPDU
		switch modbusField.FieldType {
		case MODBUS_FIELD_COIL:
			pdu = modbusModel.NewModbusPDUWriteMultipleCoilsRequest(
				modbusField.Address,
				modbusField.Quantity,
				data)
		case MODBUS_FIELD_HOLDING_REGISTER:
			pdu = modbusModel.NewModbusPDUWriteMultipleHoldingRegistersRequest(
			    modbusField.Address,
			    numWords,
			    data)
		case MODBUS_FIELD_EXTENDED_REGISTER:
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("modbus currently doesn't support extended register requests"),
			}
			return result
		default:
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("unsupported field type"),
			}
			return result
		}

		// Calculate a new unit identifier
		transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 0
			atomic.StoreInt32(&m.transactionIdentifier, 0)
		}

		// Assemble the finished ADU
		requestAdu := modbusModel.ModbusTcpADU{
			TransactionIdentifier: uint16(transactionIdentifier),
			UnitIdentifier:        m.unitIdentifier,
			Pdu:                   pdu,
		}

		// Send the ADU over the wire
		err = m.messageCodec.Send(requestAdu)
		if err != nil {
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("error sending message: " + err.Error()),
			}
		}

		// Register an expected response
		check := func(response interface{}) (bool, bool) {
			responseAdu := modbusModel.CastModbusTcpADU(response)
			return responseAdu.TransactionIdentifier == uint16(transactionIdentifier) &&
				responseAdu.UnitIdentifier == requestAdu.UnitIdentifier, false
		}
		// Register a callback to handle the response
		responseChan := m.messageCodec.Expect(check)
		go func() {
			response := <-responseChan
			// Convert the response into an ADU
			responseAdu := modbusModel.CastModbusTcpADU(response)
			// Convert the modbus response into a PLC4X response
			readResponse, err := m.ToPlc4xWriteResponse(requestAdu, *responseAdu, writeRequest)

			if err != nil {
				result <- model.PlcWriteRequestResult{
					Request: writeRequest,
					Err:     errors.New("Error decoding response: " + err.Error()),
				}
			} else {
				result <- model.PlcWriteRequestResult{
					Request:  writeRequest,
					Response: readResponse,
				}
			}
		}()
	} else {
		result <- model.PlcWriteRequestResult{
			Request:  writeRequest,
			Response: nil,
			Err:      errors.New("modbus only supports single-item requests"),
		}
	}
	fmt.Printf("Write Request %s", writeRequest)
	return result
}

func (m ModbusWriter) ToPlc4xWriteResponse(requestAdu modbusModel.ModbusTcpADU, responseAdu modbusModel.ModbusTcpADU, writeRequest model.PlcWriteRequest) (model.PlcWriteResponse, error) {
	responseCodes := map[string]model.PlcResponseCode{}
	fieldName := writeRequest.GetFieldNames()[0]

	responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
	switch responseAdu.Pdu.Child.(type) {
	case *modbusModel.ModbusPDUWriteMultipleCoilsResponse:
		req := modbusModel.CastModbusPDUWriteMultipleCoilsRequest(requestAdu.Pdu)
		resp := modbusModel.CastModbusPDUWriteMultipleCoilsResponse(responseAdu.Pdu)
		if req.Quantity == resp.Quantity {
			responseCodes[fieldName] = model.PlcResponseCode_OK
		}
	case *modbusModel.ModbusPDUWriteMultipleHoldingRegistersResponse:
		req := modbusModel.CastModbusPDUWriteMultipleHoldingRegistersRequest(requestAdu.Pdu)
		resp := modbusModel.CastModbusPDUWriteMultipleHoldingRegistersResponse(responseAdu.Pdu)
		if req.Quantity == resp.Quantity {
			responseCodes[fieldName] = model.PlcResponseCode_OK
		}
	case *modbusModel.ModbusPDUError:
		resp := modbusModel.CastModbusPDUError(&responseAdu.Pdu)
		switch resp.ExceptionCode {
		case modbusModel.ModbusErrorCode_ILLEGAL_FUNCTION:
			responseCodes[fieldName] = model.PlcResponseCode_UNSUPPORTED
		case modbusModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS:
			responseCodes[fieldName] = model.PlcResponseCode_INVALID_ADDRESS
		case modbusModel.ModbusErrorCode_ILLEGAL_DATA_VALUE:
			responseCodes[fieldName] = model.PlcResponseCode_INVALID_DATA
		case modbusModel.ModbusErrorCode_SLAVE_DEVICE_FAILURE:
			responseCodes[fieldName] = model.PlcResponseCode_REMOTE_ERROR
		case modbusModel.ModbusErrorCode_ACKNOWLEDGE:
			responseCodes[fieldName] = model.PlcResponseCode_OK
		case modbusModel.ModbusErrorCode_SLAVE_DEVICE_BUSY:
			responseCodes[fieldName] = model.PlcResponseCode_REMOTE_BUSY
		case modbusModel.ModbusErrorCode_NEGATIVE_ACKNOWLEDGE:
			responseCodes[fieldName] = model.PlcResponseCode_REMOTE_ERROR
		case modbusModel.ModbusErrorCode_MEMORY_PARITY_ERROR:
			responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
		case modbusModel.ModbusErrorCode_GATEWAY_PATH_UNAVAILABLE:
			responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
		case modbusModel.ModbusErrorCode_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND:
			responseCodes[fieldName] = model.PlcResponseCode_REMOTE_ERROR
		}
	}

	// Return the response
	return plc4goModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
