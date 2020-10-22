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
    "math"
    modbusModel "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/modbus/readwrite/model"
    plc4goModel "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/model"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/model"
    "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
    "sync/atomic"
)

type ModbusReader struct {
    transactionIdentifier uint16
    unitIdentifier int32
    messageCodec spi.MessageCodec
	spi.PlcReader
}

func NewModbusReader(transactionIdentifier uint16, messageCodec spi.MessageCodec) ModbusReader {
	return ModbusReader{
        transactionIdentifier: transactionIdentifier,
        unitIdentifier: 0,
        messageCodec: messageCodec,
    }
}

func (m ModbusReader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
    result := make(chan model.PlcReadRequestResult)
    // If we are requesting only one field, use a
    if len(readRequest.GetFieldNames()) == 1 {
        fieldName := readRequest.GetFieldNames()[0]
        field := readRequest.GetField(fieldName)
        modbusField, err := CastFromPlcField(field)
        if err != nil {
            result <- model.PlcReadRequestResult{
                Request: readRequest,
                Response: nil,
                Err: errors.New("invalid field item type"),
            }
            return result
        }
        var pdu modbusModel.IModbusPDU = nil
        switch modbusField.FieldType {
        case MODBUS_FIELD_COIL:
            pdu = modbusModel.ModbusPDUReadInputRegistersRequest{
                StartingAddress: modbusField.Address,
                Quantity:        modbusField.Quantity,
            }
        case MODBUS_FIELD_DISCRETE_INPUT:
            pdu = modbusModel.ModbusPDUReadDiscreteInputsRequest{
                StartingAddress: modbusField.Address,
                Quantity:        modbusField.Quantity,
            }
        case MODBUS_FIELD_INPUT_REGISTER:
            pdu = modbusModel.ModbusPDUReadInputRegistersRequest{
                StartingAddress: modbusField.Address,
                Quantity:        modbusField.Quantity,
            }
        case MODBUS_FIELD_HOLDING_REGISTER:
            pdu = modbusModel.ModbusPDUReadHoldingRegistersRequest{
                StartingAddress: modbusField.Address,
                Quantity:        modbusField.Quantity,
            }
        case MODBUS_FIELD_EXTENDED_REGISTER:
            result <- model.PlcReadRequestResult{
                Request: readRequest,
                Response: nil,
                Err: errors.New("modbus currently doesn't support extended register requests"),
            }
            return result
        default:
            result <- model.PlcReadRequestResult{
                Request: readRequest,
                Response: nil,
                Err: errors.New("unsupported field type"),
            }
            return result
        }

        // Calculate a new unit identifier
        unitIdentifier := atomic.AddInt32(&m.unitIdentifier, 1)
        if unitIdentifier > math.MaxUint8 {
            unitIdentifier = 0
            atomic.StoreInt32(&m.unitIdentifier, 0)
        }

        // Assemble the finished ADU
        requestAdu := modbusModel.ModbusTcpADU{
            TransactionIdentifier: m.transactionIdentifier,
            UnitIdentifier:        uint8(unitIdentifier),
            Pdu:                   pdu,
        }

        // Send the ADU over the wire
        err = m.messageCodec.Send(requestAdu)
        if err != nil {
            result <- model.PlcReadRequestResult{
                Request: readRequest,
                Response: nil,
                Err: errors.New("error sending message: " + err.Error()),
            }
        }

        // Register an expected response
        check := func(response interface{})bool {
            responseAdu := modbusModel.CastModbusTcpADU(response)
            return responseAdu.TransactionIdentifier == m.transactionIdentifier &&
                responseAdu.UnitIdentifier == requestAdu.UnitIdentifier
        }
        // Register a callback to handle the response
        responseChan := m.messageCodec.Expect(check)
        go func() {
            response := <-responseChan
            // Convert the response into an ADU
            responseAdu := modbusModel.CastModbusTcpADU(response)
            // Convert the modbus response into a PLC4X response
            readResponse, err := toPlc4xResponse(requestAdu, responseAdu, readRequest)

            if err != nil {
                result <- model.PlcReadRequestResult{
                    Request: readRequest,
                    Err: errors.New("Error decoding response: " + err.Error()),
                }
            } else {
                result <- model.PlcReadRequestResult{
                    Request:  readRequest,
                    Response: readResponse,
                }
            }
        }()
    } else {
        result <- model.PlcReadRequestResult{
            Request: readRequest,
            Response: nil,
            Err: errors.New("modbus only supports single-item requests"),
        }
    }
    fmt.Printf("Read Request %s", readRequest)
	return result
}

func toPlc4xResponse(requestAdu modbusModel.ModbusTcpADU, responseAdu modbusModel.ModbusTcpADU, readRequest model.PlcReadRequest) (model.PlcReadResponse,error) {
    switch responseAdu.Pdu.(type) {
    case modbusModel.ModbusPDUReadDiscreteInputsResponse:
        pdu := modbusModel.CastModbusPDUReadDiscreteInputsResponse(responseAdu.Pdu)
        fmt.Printf("ModbusPDUReadDiscreteInputsResponse %d", pdu)
        // Pure Boolean ...
    case modbusModel.ModbusPDUReadCoilsResponse:
        pdu := modbusModel.CastModbusPDUReadCoilsResponse(&responseAdu.Pdu)
        fmt.Printf("ModbusPDUReadCoilsResponse %d", pdu)
        // Pure Boolean ...
    case modbusModel.ModbusPDUReadInputRegistersResponse:
        pdu := modbusModel.CastModbusPDUReadInputRegistersResponse(responseAdu.Pdu)
        fmt.Printf("ModbusPDUReadInputRegistersResponse %d", pdu)
        // DataIo ...
    case modbusModel.ModbusPDUReadHoldingRegistersResponse:
        pdu := modbusModel.CastModbusPDUReadHoldingRegistersResponse(responseAdu.Pdu)
        _casted := make([]uint8, len(pdu.Value))
        for i,_val := range pdu.Value {
            _casted[i] = uint8(_val)
        }
        rb := spi.ReadBufferNew(_casted)
        fieldName := readRequest.GetFieldNames()[0]
        field := readRequest.GetField(fieldName)

        value, err := modbusModel.DataItemParse(rb, field.GetTypeName(), 1)
        if err != nil {
            return nil, err
        }
        values := map[string]values.PlcValue{}
        values[fieldName] = value
        return plc4goModel.NewDefaultPlcReadResponse(values), nil
    }
    return nil, errors.New("unsupported response type")
}

