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
	"encoding/binary"
	"math"
	"runtime/debug"
	"slices"
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

type Writer struct {
	transactionIdentifier int32
	configuration         Configuration
	messageCodec          spi.MessageCodec

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewWriter(configuration Configuration, messageCodec spi.MessageCodec, _options ...options.WithOption) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		transactionIdentifier: 0,
		configuration:         configuration,
		messageCodec:          messageCodec,
		log:                   customLogger,
	}
}

// serializeCoils packs the coil states into the bit field FC 0x0F carries: one bit per coil, the
// first coil in the least significant bit of the first byte. Coils don't go through the register
// codec - a coil is a bit, not a register - which is why this is spelled out here the way plc4j
// spells it out in ModbusTcpConnection.fromPlcValueCoil.
func serializeCoils(value apiValues.PlcValue, byteOrder ByteOrder) ([]byte, error) {
	if value == nil {
		return nil, errors.New("expecting a BOOL or a list of them when writing coils")
	}
	if !value.IsList() {
		if !value.IsBool() {
			return nil, errors.New("expecting only BOOL or list values when writing coils")
		}
		if value.GetBool() {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	}
	elements := value.GetList()
	data := make([]byte, (len(elements)+7)/8)
	// plc4j fills the buffer back to front - leading padding, then the values in reverse order -
	// and reverses the bytes afterwards. The bit of element i therefore ends up this many bits
	// behind the most significant bit of the first byte, counted before that reversal.
	padding := (8 - len(elements)%8) % 8
	for i, element := range elements {
		if element == nil || !element.IsBool() {
			return nil, errors.New("expecting only BOOL values when writing coils")
		}
		if !element.GetBool() {
			continue
		}
		bit := padding + (len(elements) - 1 - i)
		data[bit/8] |= 1 << (7 - bit%8)
	}
	if byteOrder.swapsBytes() {
		data = byteSwap(data)
	}
	slices.Reverse(data)
	return data, nil
}

func (m *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		// deliver hands one result to the caller. It is only set up once the request is under way;
		// until then the few rejections below deliver directly.
		var deliver func(response apiModel.PlcWriteResponse, err error)
		// If we are requesting only one tag, use a
		if len(writeRequest.GetTagNames()) != 1 {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.New("modbus only supports single-item requests")))
			return
		}
		tagName := writeRequest.GetTagNames()[0]

		// Get the modbus tag instance from the request
		tag := writeRequest.GetTag(tagName)
		modbusTag, err := castToModbusTagFromPlcTag(tag)
		if err != nil {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "invalid tag item type")))
			return
		}

		// A request that nobody answers must not wait forever; see the same spot in Reader.go for
		// why the cancel function is called from deliver rather than deferred.
		requestCtx, cancelRequest := withRequestTimeout(ctx, m.configuration.requestTimeout)
		deliver = func(response apiModel.PlcWriteResponse, err error) {
			cancelRequest()
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, response, err))
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(tagName)
		byteOrder := modbusTag.resolveByteOrder(m.configuration.defaultPayloadByteOrder)

		var pdu readWriteModel.ModbusPDU
		switch modbusTag.TagType {
		case Coil:
			// Every coil write, one coil or many, goes out as FC 0x0F. plc4j's
			// getWriteRequestPdu builds a WriteMultipleCoilsRequest for all of them; FC 0x05
			// appears only in the response validation, for a device that answers with it.
			data, err := serializeCoils(value, byteOrder)
			if err != nil {
				deliver(nil, errors.Wrap(err, "error serializing value"))
				return
			}
			// The byte count the PDU announces is the length of the payload, so a payload that
			// doesn't cover the addressed coils would go out as a frame no device can make sense
			// of. plc4j rejects the same mismatch in getWriteRequestPdu.
			if expected := (int(modbusTag.Quantity) + 7) / 8; len(data) != expected {
				deliver(nil, errors.Errorf("number of values (%d bytes) doesn't match the number of addressed coils (%d)", len(data), modbusTag.Quantity))
				return
			}
			pdu = readWriteModel.NewModbusPDUWriteMultipleCoilsRequest(
				modbusTag.Address,
				modbusTag.Quantity,
				data)
		case HoldingRegister, ExtendedRegister:
			// The quantity a register write announces is the number of registers the tag covers,
			// while the byte count is the length of the payload. A payload that isn't exactly that
			// many whole registers - an odd-length string, say - would go out as a frame no
			// conforming device accepts, so it is refused here the way plc4j's getWriteRequestPdu
			// refuses it.
			numWords, err := modbusTag.lengthWords()
			if err != nil {
				deliver(nil, err)
				return
			}
			data, err := SerializeRegisters(ctx, value, modbusTag.Datatype, modbusTag.Quantity,
				byteOrder, modbusTag.StringLength)
			if err != nil {
				deliver(nil, errors.Wrap(err, "error serializing value"))
				return
			}
			if len(data) != int(numWords)*2 {
				deliver(nil, errors.Errorf("number of values (%d bytes) doesn't match the number of addressed registers (%d)", len(data), numWords))
				return
			}
			if modbusTag.TagType == ExtendedRegister {
				// The extended register area is written with FC 0x15, which addresses it as a set
				// of files rather than flat, so the payload is cut up along the file boundaries it
				// crosses (plc4j ModbusTcpConnection.getWriteRequestPdu).
				groups := splitExtendedRegister(modbusTag.Address, numWords)
				items := make([]readWriteModel.ModbusPDUWriteFileRecordRequestItem, len(groups))
				offset := 0
				for i, group := range groups {
					end := offset + int(group.lengthWords)*2
					items[i] = readWriteModel.NewModbusPDUWriteFileRecordRequestItem(
						extendedRegisterReferenceType, group.fileNumber, group.recordNumber, data[offset:end])
					offset = end
				}
				pdu = readWriteModel.NewModbusPDUWriteFileRecordRequest(items)
			} else if len(data) == 2 {
				// Exactly one register goes out as FC 0x06, mirroring
				// ModbusTcpConnection.getWriteRequestPdu in plc4j.
				pdu = readWriteModel.NewModbusPDUWriteSingleRegisterRequest(
					modbusTag.Address,
					binary.BigEndian.Uint16(data))
			} else {
				pdu = readWriteModel.NewModbusPDUWriteMultipleHoldingRegistersRequest(
					modbusTag.Address,
					numWords,
					data)
			}
		default:
			deliver(nil, errors.New("unsupported tag type"))
			return
		}

		// Calculate a new unit identifier
		transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 0
			atomic.StoreInt32(&m.transactionIdentifier, 0)
		}

		// Assemble the finished ADU. Which one that is depends on the flavor the connection
		// speaks - the TCP one carries the transaction identifier, the RTU one a CRC.
		adus := m.configuration.adus()
		requestAdu := adus.buildRequest(uint16(transactionIdentifier), modbusTag.resolveUnitId(m.configuration.unitIdentifier), pdu)

		// Send the ADU over the wire
		if err = m.messageCodec.SendRequest(requestCtx, "write", requestAdu, func(message spi.Message) bool {
			return adus.acceptsResponse(requestAdu, message)
		}, func(message spi.Message) error {
			// Convert the response into an ADU
			responsePdu, err := adus.extractPdu(message)
			if err != nil {
				deliver(nil, err)
				return nil
			}
			// Convert the modbus response into a PLC4X response
			readResponse, err := m.toPlc4xWriteResponse(pdu, responsePdu, writeRequest)

			if err != nil {
				deliver(nil, errors.Wrap(err, "Error decoding response"))
			} else {
				deliver(readResponse, nil)
			}
			return nil
		}, func(err error) error {
			deliver(nil, errors.Wrap(err, "got timeout while waiting for response"))
			return nil
		}); err != nil {
			m.log.Debug().Err(err).Msg("error sending message")
			deliver(nil, errors.Wrap(err, "error sending message"))
		}
	})
	return result
}

// ToPlc4xWriteResponse turns the ADU a device answered with into a PLC4X response. It takes the
// ADUs through the little bit of them this needs, so that it serves every flavor - the generated
// ModbusADU parent type carries no PDU accessor of its own.
func (m *Writer) ToPlc4xWriteResponse(requestAdu aduWithPdu, responseAdu aduWithPdu, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	return m.toPlc4xWriteResponse(requestAdu.GetPdu(), responseAdu.GetPdu(), writeRequest)
}

func (m *Writer) toPlc4xWriteResponse(requestPdu readWriteModel.ModbusPDU, responsePdu readWriteModel.ModbusPDU, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	responseCodes := map[string]apiModel.PlcResponseCode{}
	tagName := writeRequest.GetTagNames()[0]

	// we default to an error until its proven wrong
	responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
	switch resp := responsePdu.(type) {
	case readWriteModel.ModbusPDUWriteMultipleCoilsResponse:
		req, ok := requestPdu.(readWriteModel.ModbusPDUWriteMultipleCoilsRequest)
		if !ok {
			return nil, errors.Errorf("got a write-multiple-coils response for a %T request", requestPdu)
		}
		if req.GetQuantity() == resp.GetQuantity() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		}
	case readWriteModel.ModbusPDUWriteMultipleHoldingRegistersResponse:
		req, ok := requestPdu.(readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest)
		if !ok {
			return nil, errors.Errorf("got a write-multiple-holding-registers response for a %T request", requestPdu)
		}
		if req.GetQuantity() == resp.GetQuantity() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		}
	case readWriteModel.ModbusPDUWriteSingleCoilResponse:
		// FC 0x05 echoes address and value back; anything else means the device did something
		// other than what we asked for (plc4j ModbusTcpConnection.writeSingleTag).
		req, ok := requestPdu.(readWriteModel.ModbusPDUWriteSingleCoilRequest)
		if !ok {
			return nil, errors.Errorf("got a write-single-coil response for a %T request", requestPdu)
		}
		if req.GetAddress() == resp.GetAddress() && req.GetValue() == resp.GetValue() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		} else {
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		}
	case readWriteModel.ModbusPDUWriteSingleRegisterResponse:
		// FC 0x06 echoes address and value back, same as FC 0x05.
		req, ok := requestPdu.(readWriteModel.ModbusPDUWriteSingleRegisterRequest)
		if !ok {
			return nil, errors.Errorf("got a write-single-register response for a %T request", requestPdu)
		}
		if req.GetAddress() == resp.GetAddress() && req.GetValue() == resp.GetValue() {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		} else {
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		}
	case readWriteModel.ModbusPDUWriteFileRecordResponse:
		// FC 0x15 echoes the groups it wrote back. plc4j treats every response that isn't an
		// exception as OK here, and there is nothing in the echo a mismatch could be told from
		// that the exception codes don't already cover.
		if _, ok := requestPdu.(readWriteModel.ModbusPDUWriteFileRecordRequest); !ok {
			return nil, errors.Errorf("got a write-file-record response for a %T request", requestPdu)
		}
		responseCodes[tagName] = apiModel.PlcResponseCode_OK
	case readWriteModel.ModbusPDUError:
		switch resp.GetExceptionCode() {
		case readWriteModel.ModbusErrorCode_ILLEGAL_FUNCTION:
			responseCodes[tagName] = apiModel.PlcResponseCode_UNSUPPORTED
		case readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS:
			responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
		case readWriteModel.ModbusErrorCode_ILLEGAL_DATA_VALUE:
			responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_DATA
		case readWriteModel.ModbusErrorCode_SLAVE_DEVICE_FAILURE:
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		case readWriteModel.ModbusErrorCode_ACKNOWLEDGE:
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		case readWriteModel.ModbusErrorCode_SLAVE_DEVICE_BUSY:
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_BUSY
		case readWriteModel.ModbusErrorCode_NEGATIVE_ACKNOWLEDGE:
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		case readWriteModel.ModbusErrorCode_MEMORY_PARITY_ERROR:
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
		case readWriteModel.ModbusErrorCode_GATEWAY_PATH_UNAVAILABLE:
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
		case readWriteModel.ModbusErrorCode_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND:
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		default:
			m.log.Debug().Interface("exceptionCode", resp.GetExceptionCode()).Msg("Unmapped exception code")
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", resp)
	}

	// Return the response
	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
