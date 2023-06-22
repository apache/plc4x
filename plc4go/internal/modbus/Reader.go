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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
	"math"
	"runtime/debug"
	"sync/atomic"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"

	"github.com/pkg/errors"
)

type Reader struct {
	transactionIdentifier int32
	unitIdentifier        uint8
	messageCodec          spi.MessageCodec

	passLogToModel bool
	log            zerolog.Logger
}

func NewReader(unitIdentifier uint8, messageCodec spi.MessageCodec, _options ...options.WithOption) *Reader {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		transactionIdentifier: 0,
		unitIdentifier:        unitIdentifier,
		messageCodec:          messageCodec,
		passLogToModel:        passLoggerToModel,
		log:                   customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	// TODO: handle ctx
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		if len(readRequest.GetTagNames()) != 1 {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.New("modbus only supports single-item requests"))
			m.log.Debug().Msgf("modbus only supports single-item requests. Got %d tags", len(readRequest.GetTagNames()))
			return
		}
		// If we are requesting only one tag, use a
		tagName := readRequest.GetTagNames()[0]
		tag := readRequest.GetTag(tagName)
		modbusTagVar, err := castToModbusTagFromPlcTag(tag)
		if err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				nil,
				errors.Wrap(err, "invalid tag item type"),
			)
			m.log.Debug().Msgf("Invalid tag item type %T", tag)
			return
		}
		numWords := uint16(math.Ceil(float64(modbusTagVar.Quantity*uint16(modbusTagVar.Datatype.DataTypeSize())) / float64(2)))
		m.log.Debug().Msgf("Working with %d words", numWords)
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
			result <- spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				nil,
				errors.New("modbus currently doesn't support extended register requests"),
			)
			return
		default:
			result <- spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				nil,
				errors.Errorf("unsupported tag type %x", modbusTagVar.TagType),
			)
			m.log.Debug().Msgf("Unsupported tag type %x", modbusTagVar.TagType)
			return
		}

		// Calculate a new transaction identifier
		transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 1
			atomic.StoreInt32(&m.transactionIdentifier, 1)
		}
		m.log.Debug().Msgf("Calculated transaction identifier %x", transactionIdentifier)

		// Assemble the finished ADU
		m.log.Trace().Msg("Assemble ADU")
		requestAdu := readWriteModel.NewModbusTcpADU(uint16(transactionIdentifier), m.unitIdentifier, pdu, false)

		// Send the ADU over the wire
		m.log.Trace().Msg("Send ADU")
		if err = m.messageCodec.SendRequest(ctx, requestAdu, func(message spi.Message) bool {
			responseAdu := message.(readWriteModel.ModbusTcpADU)
			return responseAdu.GetTransactionIdentifier() == uint16(transactionIdentifier) &&
				responseAdu.GetUnitIdentifier() == requestAdu.UnitIdentifier
		}, func(message spi.Message) error {
			// Convert the response into an ADU
			m.log.Trace().Msg("convert response to ADU")
			responseAdu := message.(readWriteModel.ModbusTcpADU)
			// Convert the modbus response into a PLC4X response
			m.log.Trace().Msg("convert response to PLC4X response")
			readResponse, err := m.ToPlc4xReadResponse(responseAdu, readRequest)

			if err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					nil,
					errors.Wrap(err, "Error decoding response"),
				)
				// TODO: should we return the error here?
				return nil
			}
			result <- spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				readResponse,
				nil,
			)
			return nil
		}, func(err error) error {
			result <- spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				nil,
				errors.Wrap(err, "got timeout while waiting for response"),
			)
			return nil
		}, time.Second*1); err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				nil,
				errors.Wrap(err, "error sending message"),
			)
		}
	}()
	return result
}

func (m *Reader) ToPlc4xReadResponse(responseAdu readWriteModel.ModbusTcpADU, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
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
	m.log.Trace().Msg("get a tag from request")
	tagName := readRequest.GetTagNames()[0]
	tag, err := castToModbusTagFromPlcTag(readRequest.GetTag(tagName))
	if err != nil {
		return nil, errors.Wrap(err, "error casting to modbus-tag")
	}

	// Decode the data according to the information from the request
	m.log.Trace().Msg("decode data")
	ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
	value, err := readWriteModel.DataItemParse(ctxForModel, data, tag.Datatype, tag.Quantity)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing data item")
	}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]apiValues.PlcValue{}
	plcValues[tagName] = value
	responseCodes[tagName] = apiModel.PlcResponseCode_OK

	// Return the response
	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}
