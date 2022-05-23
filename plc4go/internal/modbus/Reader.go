/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"sync/atomic"
	"time"
)

type Reader struct {
	transactionIdentifier int32
	unitIdentifier        uint8
	messageCodec          spi.MessageCodec
}

func NewReader(unitIdentifier uint8, messageCodec spi.MessageCodec) *Reader {
	return &Reader{
		transactionIdentifier: 0,
		unitIdentifier:        unitIdentifier,
		messageCodec:          messageCodec,
	}
}

func (m *Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {
		if len(readRequest.GetFieldNames()) != 1 {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("modbus only supports single-item requests"),
			}
			log.Debug().Msgf("modbus only supports single-item requests. Got %d fields", len(readRequest.GetFieldNames()))
			return
		}
		// If we are requesting only one field, use a
		fieldName := readRequest.GetFieldNames()[0]
		field := readRequest.GetField(fieldName)
		modbusField, err := CastToModbusFieldFromPlcField(field)
		if err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			log.Debug().Msgf("Invalid field item type %T", field)
			return
		}
		numWords := uint16(math.Ceil(float64(modbusField.Quantity*uint16(modbusField.Datatype.DataTypeSize())) / float64(2)))
		log.Debug().Msgf("Working with %d words", numWords)
		var pdu *readWriteModel.ModbusPDU = nil
		switch modbusField.FieldType {
		case Coil:
			pdu = readWriteModel.NewModbusPDUReadCoilsRequest(modbusField.Address, modbusField.Quantity).GetParent()
		case DiscreteInput:
			pdu = readWriteModel.NewModbusPDUReadDiscreteInputsRequest(modbusField.Address, modbusField.Quantity).GetParent()
		case InputRegister:
			pdu = readWriteModel.NewModbusPDUReadInputRegistersRequest(modbusField.Address, numWords).GetParent()
		case HoldingRegister:
			pdu = readWriteModel.NewModbusPDUReadHoldingRegistersRequest(modbusField.Address, numWords).GetParent()
		case ExtendedRegister:
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("modbus currently doesn't support extended register requests"),
			}
			return
		default:
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Errorf("unsupported field type %x", modbusField.FieldType),
			}
			log.Debug().Msgf("Unsupported field type %x", modbusField.FieldType)
			return
		}

		// Calculate a new transaction identifier
		transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 1
			atomic.StoreInt32(&m.transactionIdentifier, 1)
		}
		log.Debug().Msgf("Calculated transaction identifier %x", transactionIdentifier)

		// Assemble the finished ADU
		log.Trace().Msg("Assemble ADU")
		requestAdu := readWriteModel.ModbusTcpADU{
			TransactionIdentifier: uint16(transactionIdentifier),
			UnitIdentifier:        m.unitIdentifier,
			Pdu:                   pdu,
		}

		// Send the ADU over the wire
		log.Trace().Msg("Send ADU")
		if err = m.messageCodec.SendRequest(
			requestAdu,
			func(message interface{}) bool {
				responseAdu := readWriteModel.CastModbusTcpADU(message)
				return responseAdu.TransactionIdentifier == uint16(transactionIdentifier) &&
					responseAdu.UnitIdentifier == requestAdu.UnitIdentifier
			},
			func(message interface{}) error {
				// Convert the response into an ADU
				log.Trace().Msg("convert response to ADU")
				responseAdu := readWriteModel.CastModbusTcpADU(message)
				// Convert the modbus response into a PLC4X response
				log.Trace().Msg("convert response to PLC4X response")
				readResponse, err := m.ToPlc4xReadResponse(*responseAdu, readRequest)

				if err != nil {
					result <- &plc4goModel.DefaultPlcReadRequestResult{
						Request: readRequest,
						Err:     errors.Wrap(err, "Error decoding response"),
					}
					// TODO: should we return the error here?
					return nil
				}
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: readResponse,
				}
				return nil
			},
			func(err error) error {
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request: readRequest,
					Err:     errors.Wrap(err, "got timeout while waiting for response"),
				}
				return nil
			},
			time.Second*1); err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error sending message"),
			}
		}
	}()
	return result
}

func (m *Reader) ToPlc4xReadResponse(responseAdu readWriteModel.ModbusTcpADU, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	var data []uint8
	switch responseAdu.Pdu.Child.(type) {
	case *readWriteModel.ModbusPDUReadDiscreteInputsResponse:
		pdu := readWriteModel.CastModbusPDUReadDiscreteInputsResponse(responseAdu.Pdu)
		data = pdu.Value
		// Pure Boolean ...
	case *readWriteModel.ModbusPDUReadCoilsResponse:
		pdu := readWriteModel.CastModbusPDUReadCoilsResponse(responseAdu.Pdu)
		data = pdu.Value
		// Pure Boolean ...
	case *readWriteModel.ModbusPDUReadInputRegistersResponse:
		pdu := readWriteModel.CastModbusPDUReadInputRegistersResponse(responseAdu.Pdu)
		data = pdu.Value
		// DataIo ...
	case *readWriteModel.ModbusPDUReadHoldingRegistersResponse:
		pdu := readWriteModel.CastModbusPDUReadHoldingRegistersResponse(responseAdu.Pdu)
		data = pdu.Value
	case *readWriteModel.ModbusPDUError:
		return nil, errors.Errorf("got an error from remote. Errorcode %x", responseAdu.Pdu.Child.(*readWriteModel.ModbusPDUError).ExceptionCode)
	default:
		return nil, errors.Errorf("unsupported response type %T", responseAdu.Pdu.Child)
	}

	// Get the field from the request
	log.Trace().Msg("get a field from request")
	fieldName := readRequest.GetFieldNames()[0]
	field, err := CastToModbusFieldFromPlcField(readRequest.GetField(fieldName))
	if err != nil {
		return nil, errors.Wrap(err, "error casting to modbus-field")
	}

	// Decode the data according to the information from the request
	log.Trace().Msg("decode data")
	rb := utils.NewReadBufferByteBased(data)
	value, err := readWriteModel.DataItemParse(rb, field.Datatype, field.Quantity)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing data item")
	}
	responseCodes := map[string]model.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}
	plcValues[fieldName] = value
	responseCodes[fieldName] = model.PlcResponseCode_OK

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}
