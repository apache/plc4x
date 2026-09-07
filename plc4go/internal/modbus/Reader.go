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
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Reader struct {
	transactionIdentifier int32
	configuration         Configuration
	messageCodec          spi.MessageCodec

	wg sync.WaitGroup // use to track spawned go routines

	passLogToModel bool
	log            zerolog.Logger
}

func NewReader(configuration Configuration, messageCodec spi.MessageCodec, _options ...options.WithOption) *Reader {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		transactionIdentifier: 0,
		configuration:         configuration,
		messageCodec:          messageCodec,
		passLogToModel:        passLoggerToModel,
		log:                   customLogger,
	}
}

// readRequestPdu is the request that reads what the tag addresses (plc4j
// ModbusTcpConnection.getReadRequestPdu).
func readRequestPdu(tag modbusTag) (readWriteModel.ModbusPDU, error) {
	switch tag.TagType {
	case Coil:
		return readWriteModel.NewModbusPDUReadCoilsRequest(tag.Address, tag.Quantity), nil
	case DiscreteInput:
		return readWriteModel.NewModbusPDUReadDiscreteInputsRequest(tag.Address, tag.Quantity), nil
	case InputRegister, HoldingRegister:
		numWords, err := tag.lengthWords()
		if err != nil {
			return nil, err
		}
		if tag.TagType == InputRegister {
			return readWriteModel.NewModbusPDUReadInputRegistersRequest(tag.Address, numWords), nil
		}
		return readWriteModel.NewModbusPDUReadHoldingRegistersRequest(tag.Address, numWords), nil
	case ExtendedRegister:
		// The extended register area is read with FC 0x14, which addresses it as a set of files
		// rather than flat, so the address turns into one item per file it touches (plc4j
		// ModbusTcpConnection.getReadRequestPdu).
		numWords, err := tag.lengthWords()
		if err != nil {
			return nil, err
		}
		groups := splitExtendedRegister(tag.Address, numWords)
		items := make([]readWriteModel.ModbusPDUReadFileRecordRequestItem, len(groups))
		for i, group := range groups {
			items[i] = readWriteModel.NewModbusPDUReadFileRecordRequestItem(
				extendedRegisterReferenceType, group.fileNumber, group.recordNumber, group.lengthWords)
		}
		return readWriteModel.NewModbusPDUReadFileRecordRequest(items), nil
	default:
		return nil, errors.Errorf("unsupported tag type %x", tag.TagType)
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	// TODO: handle ctx
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		// deliver hands one result to the caller. It is only set up once the request is under way;
		// until then the few rejections below deliver directly.
		var deliver func(response apiModel.PlcReadResponse, err error)
		if len(readRequest.GetTagNames()) != 1 {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.New("modbus only supports single-item requests")))
			m.log.Debug().Int("nTags", len(readRequest.GetTagNames())).Msg("modbus only supports single-item requests. Got nTags tags")
			return
		}
		// If we are requesting only one tag, use a
		tagName := readRequest.GetTagNames()[0]
		tag := readRequest.GetTag(tagName)
		modbusTagVar, err := castToModbusTagFromPlcTag(tag)
		if err != nil {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(
				readRequest,
				nil,
				errors.Wrap(err, "invalid tag item type"),
			))
			m.log.Debug().Type("tagType", tag).Msg("Invalid tag item type")
			return
		}

		// A request that nobody answers must not wait forever; the codec derives the lifetime of the
		// expectation it registers from the deadline of the context it is handed.
		// The cancel function must not be deferred here: SendRequest only registers the
		// expectation and returns, so this goroutine finishes long before the response arrives, and
		// cancelling here would take the expectation down with it. It is called once a result has
		// been delivered instead, which every path below does; a request that is never answered is
		// released by its own deadline.
		requestCtx, cancelRequest := withRequestTimeout(ctx, m.configuration.requestTimeout)
		deliver = func(response apiModel.PlcReadResponse, err error) {
			cancelRequest()
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, response, err))
		}
		pdu, err := readRequestPdu(modbusTagVar)
		if err != nil {
			deliver(nil, err)
			m.log.Debug().Err(err).Stringer("tagType", modbusTagVar.TagType).Msg("Couldn't build a read request")
			return
		}

		// Calculate a new transaction identifier
		transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 1
			atomic.StoreInt32(&m.transactionIdentifier, 1)
		}
		m.log.Debug().Int32("transactionIdentifier", transactionIdentifier).Msg("Calculated transaction identifier")

		// Assemble the finished ADU. Which one that is depends on the flavor the connection
		// speaks - the TCP one carries the transaction identifier, the RTU one a CRC.
		m.log.Trace().Msg("Assemble ADU")
		adus := m.configuration.adus()
		requestAdu := adus.buildRequest(uint16(transactionIdentifier), modbusTagVar.resolveUnitId(m.configuration.unitIdentifier), pdu)

		// Send the ADU over the wire
		m.log.Trace().Msg("Send ADU")
		if err = m.messageCodec.SendRequest(requestCtx, "read", requestAdu, func(message spi.Message) bool {
			return adus.acceptsResponse(requestAdu, message)
		}, func(message spi.Message) error {
			// Convert the response into an ADU
			m.log.Trace().Msg("convert response to ADU")
			responsePdu, err := adus.extractPdu(message)
			if err != nil {
				deliver(nil, err)
				return nil
			}
			// Convert the modbus response into a PLC4X response
			m.log.Trace().Msg("convert response to PLC4X response")
			readResponse, err := m.toPlc4xReadResponse(responsePdu, readRequest)

			if err != nil {
				deliver(nil, errors.Wrap(err, "Error decoding response"))
				// TODO: should we return the error here?
				return nil
			}
			deliver(readResponse, nil)
			return nil
		}, func(err error) error {
			deliver(nil, errors.Wrap(err, "got timeout while waiting for response"))
			return nil
		}); err != nil {
			deliver(nil, errors.Wrap(err, "error sending message"))
		}
	})
	return result
}

// ToPlc4xReadResponse turns the ADU a device answered with into a PLC4X response. It takes the
// ADU through the little bit of it this needs, so that it serves every flavor - the generated
// ModbusADU parent type carries no PDU accessor of its own.
func (m *Reader) ToPlc4xReadResponse(responseAdu aduWithPdu, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	return m.toPlc4xReadResponse(responseAdu.GetPdu(), readRequest)
}

func (m *Reader) toPlc4xReadResponse(responsePdu readWriteModel.ModbusPDU, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	ctx := context.TODO()
	var data []uint8
	switch pdu := responsePdu.(type) {
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
	case readWriteModel.ModbusPDUReadFileRecordResponse:
		// A request that crossed a file boundary was sent as one item per file, and the response
		// carries one item per request item, so the value is the registers of all of them in
		// order. plc4j's extractResponseData only looks at the first item, which silently drops
		// everything behind the boundary.
		for _, item := range pdu.GetItems() {
			data = append(data, item.GetData()...)
		}
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
	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	value, err := ParseRegisters(ctxForModel, data, tag.Datatype, tag.Quantity, tag.resolveByteOrder(m.configuration.defaultPayloadByteOrder), tag.StringLength)
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
