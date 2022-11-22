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
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func (m *Reader) Read(ctx context.Context, readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	// TODO: handle ctx
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {
		if len(readRequest.GetTagNames()) != 1 {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("modbus only supports single-item requests"),
			}
			log.Debug().Msgf("modbus only supports single-item requests. Got %d tags", len(readRequest.GetTagNames()))
			return
		}
		// If we are requesting only one tag, use a
		tagName := readRequest.GetTagNames()[0]
		tag := readRequest.GetTag(tagName)
		modbusTagVar, err := CastToModbusTagFromPlcTag(tag)
		if err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid tag item type"),
			}
			log.Debug().Msgf("Invalid tag item type %T", tag)
			return
		}
		numWords := uint16(math.Ceil(float64(modbusTagVar.Quantity*uint16(modbusTagVar.Datatype.DataTypeSize())) / float64(2)))
		log.Debug().Msgf("Working with %d words", numWords)
		var pdu readWriteModel.ModbusPDU = nil
		switch modbusTagVar.TagType {
		case Coil:
			pdu = readWriteModel.NewModbusPDUReadCoilsRequest(modbusTagVar.Address, modbusTagVar.Quantity)
		case DiscreteInput:
			pdu = readWriteModel.NewModbusPDUReadDiscreteInputsRequest(modbusTagVar.Address, modbusTagVar.Quantity)
		case InputRegister:
			pdu = readWriteModel.NewModbusPDUReadInputRegistersRequest(modbusTagVar.Address, numWords)
		case HoldingRegister:
			pdu = readWriteModel.NewModbusPDUReadHoldingRegistersRequest(modbusTagVar.Address, numWords)
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
				Err:      errors.Errorf("unsupported tag type %x", modbusTagVar.TagType),
			}
			log.Debug().Msgf("Unsupported tag type %x", modbusTagVar.TagType)
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
		requestAdu := readWriteModel.NewModbusTcpADU(uint16(transactionIdentifier), m.unitIdentifier, pdu, false)

		// Send the ADU over the wire
		log.Trace().Msg("Send ADU")
		if err = m.messageCodec.SendRequest(ctx, requestAdu, func(message spi.Message) bool {
			responseAdu := message.(readWriteModel.ModbusTcpADU)
			return responseAdu.GetTransactionIdentifier() == uint16(transactionIdentifier) &&
				responseAdu.GetUnitIdentifier() == requestAdu.UnitIdentifier
		}, func(message spi.Message) error {
			// Convert the response into an ADU
			log.Trace().Msg("convert response to ADU")
			responseAdu := message.(readWriteModel.ModbusTcpADU)
			// Convert the modbus response into a PLC4X response
			log.Trace().Msg("convert response to PLC4X response")
			readResponse, err := m.ToPlc4xReadResponse(responseAdu, readRequest)

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
		}, func(err error) error {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request: readRequest,
				Err:     errors.Wrap(err, "got timeout while waiting for response"),
			}
			return nil
		}, time.Second*1); err != nil {
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
	switch pdu := responseAdu.GetPdu().(type) {
	case readWriteModel.ModbusPDUReadDiscreteInputsResponse:
		data = pdu.GetValue()
		// Pure Boolean ...
	case readWriteModel.ModbusPDUReadCoilsResponse:
		data = pdu.GetValue()
		// Pure Boolean ...
	case readWriteModel.ModbusPDUReadInputRegistersResponse:
		data = pdu.GetValue()
		// DataIo ...
	case readWriteModel.ModbusPDUReadHoldingRegistersResponse:
		data = pdu.GetValue()
	case readWriteModel.ModbusPDUError:
		return nil, errors.Errorf("got an error from remote. Errorcode %x", pdu.GetExceptionCode())
	default:
		return nil, errors.Errorf("unsupported response type %T", pdu)
	}

	// Get the tag from the request
	log.Trace().Msg("get a tag from request")
	tagName := readRequest.GetTagNames()[0]
	tag, err := CastToModbusTagFromPlcTag(readRequest.GetTag(tagName))
	if err != nil {
		return nil, errors.Wrap(err, "error casting to modbus-tag")
	}

	// Decode the data according to the information from the request
	log.Trace().Msg("decode data")
	value, err := readWriteModel.DataItemParse(data, tag.Datatype, tag.Quantity)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing data item")
	}
	responseCodes := map[string]model.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}
	plcValues[tagName] = value
	responseCodes[tagName] = model.PlcResponseCode_OK

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}
