/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package modbus

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Writer struct {
	transactionIdentifier int32
	unitIdentifier        uint8
	messageCodec          spi.MessageCodec
}

func NewWriter(unitIdentifier uint8, messageCodec spi.MessageCodec) Writer {
	return Writer{
		transactionIdentifier: 0,
		unitIdentifier:        unitIdentifier,
		messageCodec:          messageCodec,
	}
}

func (m Writer) Write(ctx context.Context, writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	// TODO: handle context
	result := make(chan model.PlcWriteRequestResult, 1)
	go func() {
		// If we are requesting only one tag, use a
		if len(writeRequest.GetTagNames()) != 1 {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.New("modbus only supports single-item requests"))
			return
		}
		tagName := writeRequest.GetTagNames()[0]

		// Get the modbus tag instance from the request
		tag := writeRequest.GetTag(tagName)
		modbusTag, err := CastToModbusTagFromPlcTag(tag)
		if err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "invalid tag item type"))
			return
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(tagName)
		data, err := readWriteModel.DataItemSerialize(value, modbusTag.Datatype, modbusTag.Quantity)
		if err != nil {
			result <- &spiModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error serializing value"),
			}
			return
		}

		// Calculate the number of words needed to send the data
		numWords := uint16(math.Ceil(float64(len(data)) / 2))

		var pdu readWriteModel.ModbusPDU
		switch modbusTag.TagType {
		case Coil:
			pdu = readWriteModel.NewModbusPDUWriteMultipleCoilsRequest(
				modbusTag.Address,
				modbusTag.Quantity,
				data)
		case HoldingRegister:
			pdu = readWriteModel.NewModbusPDUWriteMultipleHoldingRegistersRequest(
				modbusTag.Address,
				numWords,
				data)
		case ExtendedRegister:
			result <- &spiModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("modbus currently doesn't support extended register requests"),
			}
			return
		default:
			result <- &spiModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("unsupported tag type"),
			}
			return
		}

		// Calculate a new unit identifier
		transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 0
			atomic.StoreInt32(&m.transactionIdentifier, 0)
		}

		// Assemble the finished ADU
		requestAdu := readWriteModel.NewModbusTcpADU(uint16(transactionIdentifier), m.unitIdentifier, pdu, false)

		// Send the ADU over the wire
		err = m.messageCodec.SendRequest(ctx, requestAdu, func(message spi.Message) bool {
			responseAdu := message.(readWriteModel.ModbusTcpADU)
			return responseAdu.GetTransactionIdentifier() == uint16(transactionIdentifier) &&
				responseAdu.GetUnitIdentifier() == requestAdu.UnitIdentifier
		}, func(message spi.Message) error {
			// Convert the response into an ADU
			responseAdu := message.(readWriteModel.ModbusTcpADU)
			// Convert the modbus response into a PLC4X response
			readResponse, err := m.ToPlc4xWriteResponse(requestAdu, responseAdu, writeRequest)

			if err != nil {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request: writeRequest,
					Err:     errors.Wrap(err, "Error decoding response"),
				}
			} else {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: readResponse,
				}
			}
			return nil
		}, func(err error) error {
			result <- &spiModel.DefaultPlcWriteRequestResult{
				Request: writeRequest,
				Err:     errors.New("got timeout while waiting for response"),
			}
			return nil
		}, time.Second*1)
	}()
	return result
}

func (m Writer) ToPlc4xWriteResponse(requestAdu readWriteModel.ModbusTcpADU, responseAdu readWriteModel.ModbusTcpADU, writeRequest model.PlcWriteRequest) (model.PlcWriteResponse, error) {
	responseCodes := map[string]model.PlcResponseCode{}
	tagName := writeRequest.GetTagNames()[0]

	// we default to an error until its proven wrong
	responseCodes[tagName] = model.PlcResponseCode_INTERNAL_ERROR
	switch resp := responseAdu.GetPdu().(type) {
	case readWriteModel.ModbusPDUWriteMultipleCoilsResponse:
		req := requestAdu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleCoilsRequest)
		if req.GetQuantity() == resp.GetQuantity() {
			responseCodes[tagName] = model.PlcResponseCode_OK
		}
	case readWriteModel.ModbusPDUWriteMultipleHoldingRegistersResponse:
		req := requestAdu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest)
		if req.GetQuantity() == resp.GetQuantity() {
			responseCodes[tagName] = model.PlcResponseCode_OK
		}
	case readWriteModel.ModbusPDUError:
		switch resp.GetExceptionCode() {
		case readWriteModel.ModbusErrorCode_ILLEGAL_FUNCTION:
			responseCodes[tagName] = model.PlcResponseCode_UNSUPPORTED
		case readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS:
			responseCodes[tagName] = model.PlcResponseCode_INVALID_ADDRESS
		case readWriteModel.ModbusErrorCode_ILLEGAL_DATA_VALUE:
			responseCodes[tagName] = model.PlcResponseCode_INVALID_DATA
		case readWriteModel.ModbusErrorCode_SLAVE_DEVICE_FAILURE:
			responseCodes[tagName] = model.PlcResponseCode_REMOTE_ERROR
		case readWriteModel.ModbusErrorCode_ACKNOWLEDGE:
			responseCodes[tagName] = model.PlcResponseCode_OK
		case readWriteModel.ModbusErrorCode_SLAVE_DEVICE_BUSY:
			responseCodes[tagName] = model.PlcResponseCode_REMOTE_BUSY
		case readWriteModel.ModbusErrorCode_NEGATIVE_ACKNOWLEDGE:
			responseCodes[tagName] = model.PlcResponseCode_REMOTE_ERROR
		case readWriteModel.ModbusErrorCode_MEMORY_PARITY_ERROR:
			responseCodes[tagName] = model.PlcResponseCode_INTERNAL_ERROR
		case readWriteModel.ModbusErrorCode_GATEWAY_PATH_UNAVAILABLE:
			responseCodes[tagName] = model.PlcResponseCode_INTERNAL_ERROR
		case readWriteModel.ModbusErrorCode_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND:
			responseCodes[tagName] = model.PlcResponseCode_REMOTE_ERROR
		default:
			log.Debug().Msgf("Unmapped exception code %x", resp.GetExceptionCode())
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", resp)
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
