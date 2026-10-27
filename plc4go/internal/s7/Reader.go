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

package s7

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type Reader struct {
	tpduGenerator *TpduGenerator
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	driverContext *DriverContext

	wg sync.WaitGroup // use to track spawned go routines

	passLogToModel bool
	log            zerolog.Logger
}

func NewReader(tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, driverContext *DriverContext, _options ...options.WithOption) *Reader {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		tpduGenerator:  tpduGenerator,
		messageCodec:   messageCodec,
		tm:             tm,
		driverContext:  driverContext,
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()

		chunks, err := splitRead(readRequest, m.driverContext.PduSize)
		if err != nil {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, err))
			return
		}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		plcValues := map[string]apiValues.PlcValue{}
		// The reassembled payload of each tag that was read in fragments.
		fragmented := map[string][]byte{}
		for _, chunk := range chunks {
			requestItems := make([]readWriteModel.S7VarRequestParameterItem, len(chunk))
			for i, item := range chunk {
				requestItems[i] = readWriteModel.NewS7VarRequestParameterItemAddress(item.address)
			}
			m.log.Trace().Int("items", len(chunk)).Msg("Send read request")
			response, err := exchange(ctx, m.log, m.tpduGenerator, m.messageCodec, m.tm, "read",
				readWriteModel.NewS7ParameterReadVarRequest(requestItems), nil)
			if err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, err))
				return
			}
			if err := m.applyReadResponse(response, chunk, readRequest, responseCodes, plcValues, fragmented); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil,
					errors.Wrap(err, "Error decoding response")))
				return
			}
		}
		// A tag read in fragments can only be decoded once all of them are in.
		for tagName, data := range fragmented {
			if responseCodes[tagName] != apiModel.PlcResponseCode_OK {
				continue
			}
			plcValue, err := m.parse(ctx, readRequest.GetTag(tagName).(PlcTag), data)
			if err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil,
					errors.Wrap(err, "Error decoding response")))
				return
			}
			plcValues[tagName] = plcValue
		}
		m.log.Trace().Msg("Returning the response")
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest,
			spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil))
	})
	return result
}

// applyReadResponse records what the response to one request of a read says about its tags.
func (m *Reader) applyReadResponse(response readWriteModel.S7Message, chunk []readItem, readRequest apiModel.PlcReadRequest,
	responseCodes map[string]apiModel.PlcResponseCode, plcValues map[string]apiValues.PlcValue, fragmented map[string][]byte) error {
	errorClass, errorCode, err := headerError(response)
	if err != nil {
		return err
	}
	// If the result contains any form of non-null error code, handle this instead.
	if (errorClass != 0) || (errorCode != 0) {
		responseCode := mapPlcErrorCode(m.log, errorClass, errorCode)
		for _, item := range chunk {
			mergeResponseCode(responseCodes, item.tagName, responseCode)
			plcValues[item.tagName] = spiValues.NewPlcNULL()
		}
		return nil
	}

	payload, ok := response.GetPayload().(readWriteModel.S7PayloadReadVarResponse)
	if !ok {
		return errors.Errorf("unexpected payload %T", response.GetPayload())
	}
	// If the numbers of items don't match, we're in big trouble as the only
	// way to know how to interpret the responses is by aligning them with the
	// items from the request as this information is not returned by the PLC.
	if len(chunk) != len(payload.GetItems()) {
		return errors.New("The number of requested items doesn't match the number of returned items")
	}
	for i, item := range chunk {
		payloadItem := payload.GetItems()[i]
		responseCode := decodeResponseCode(payloadItem.GetReturnCode())
		mergeResponseCode(responseCodes, item.tagName, responseCode)
		if responseCode != apiModel.PlcResponseCode_OK {
			continue
		}
		if item.split {
			data, ok := fragmented[item.tagName]
			if !ok {
				data = make([]byte, payloadBytes(readAddressOf(readRequest.GetTag(item.tagName).(PlcTag))))
				fragmented[item.tagName] = data
			}
			copy(data[item.fragmentByteOffset:], payloadItem.GetData())
			continue
		}
		plcValue, err := m.parse(context.TODO(), readRequest.GetTag(item.tagName).(PlcTag), payloadItem.GetData())
		if err != nil {
			return errors.Wrap(err, "Error parsing data item")
		}
		plcValues[item.tagName] = plcValue
	}
	return nil
}

func (m *Reader) parse(ctx context.Context, tag PlcTag, data []byte) (apiValues.PlcValue, error) {
	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	return parsePlcValue(ctxForModel, tag, data, m.driverContext.ControllerType)
}

// Currently we only support the S7 Any type of addresses. This helper simply converts the S7Tag from PLC4X into
// S7Address objects.
func encodeS7Address(tag apiModel.PlcTag) (readWriteModel.S7Address, error) {
	s7Tag, ok := tag.(PlcTag)
	if !ok {
		return nil, errors.Errorf("Unsupported address type %t", tag)
	}
	return readAddressOf(s7Tag), nil
}

// readAddressOf is the address a read of the tag asks for. Strings, BOOL arrays and types without
// an own transport size are requested as the characters or bytes they occupy; cyclic
// subscriptions address their tags the same way.
func readAddressOf(s7Tag PlcTag) readWriteModel.S7AddressAny {
	transportSize := s7Tag.GetDataType()
	numElements := s7Tag.GetNumElements()
	if transportSize == readWriteModel.TransportSize_STRING {
		transportSize = readWriteModel.TransportSize_CHAR
		stringLength := uint16(254)
		if s7StringTag, ok := s7Tag.(PlcStringTag); ok {
			stringLength = s7StringTag.stringLength
		}
		numElements = numElements * (stringLength + 2)
	} else if transportSize == readWriteModel.TransportSize_WSTRING {
		transportSize = readWriteModel.TransportSize_CHAR
		stringLength := uint16(254)
		if s7StringTag, ok := s7Tag.(PlcStringTag); ok {
			stringLength = s7StringTag.stringLength
		}
		numElements = numElements * (stringLength + 2) * 2
	} else if transportSize == readWriteModel.TransportSize_BOOL && numElements > 1 {
		// BOOL arrays are read as the bytes their bits are packed into, from bit 0 of the start
		// byte: a byte read has no bit address, so parsePlcValue skips the bits before the tag's
		// own bit offset.
		return readWriteModel.NewS7AddressAny(
			readWriteModel.TransportSize_BYTE,
			packedBoolBytes(s7Tag),
			s7Tag.GetBlockNumber(),
			s7Tag.GetMemoryArea(),
			s7Tag.GetByteOffset(),
			0,
		)
	}
	// Types without an own transport-size code (temporal types, 64-bit types, ...) have to be
	// requested as plain byte arrays, the S7 denies them with "Data type not supported" otherwise.
	if transportSize.Code() == 0x00 {
		numElements = numElements * uint16(transportSize.SizeInBytes())
		transportSize = readWriteModel.TransportSize_BYTE
	}
	return readWriteModel.NewS7AddressAny(
		transportSize,
		numElements,
		s7Tag.GetBlockNumber(),
		s7Tag.GetMemoryArea(),
		s7Tag.GetByteOffset(),
		s7Tag.GetBitOffset(),
	)
}

// packedBoolBytes is how many bytes the bits of a BOOL tag span, counting from bit 0 of its start
// byte: a tag that starts at a bit offset other than 0 can spill into one more byte than its
// element count alone needs.
func packedBoolBytes(tag PlcTag) uint16 {
	return uint16((uint32(tag.GetBitOffset()) + uint32(tag.GetNumElements()) + 7) / 8)
}

// parsePlcValue decodes the data a read of the tag returned. Every request addresses a BOOL array
// from bit 0 of its start byte (see readAddressOf), so its elements begin at the tag's bit offset.
func parsePlcValue(ctx context.Context, tag PlcTag, data []byte, controllerType readWriteModel.ControllerType) (apiValues.PlcValue, error) {
	stringLength := int32(254)
	if s7StringTag, ok := tag.(PlcStringTag); ok {
		stringLength = int32(s7StringTag.stringLength)
	}
	numElements := tag.GetNumElements()
	dataProtocolId := tag.GetDataType().DataProtocolId()
	if numElements == 1 {
		return readWriteModel.DataItemParse(ctx, data, dataProtocolId, controllerType, stringLength)
	}
	switch tag.GetDataType() {
	case readWriteModel.TransportSize_BYTE:
		return spiValues.NewPlcRawByteArray(data), nil
	case readWriteModel.TransportSize_BOOL:
		// BOOL arrays arrive as packed bytes (LSB first), matching the write-side packing.
		items := make([]apiValues.PlcValue, 0, numElements)
		firstBit := tag.GetBitOffset()
		for i := range numElements {
			bit := int(firstBit) + int(i)
			if bit/8 >= len(data) {
				return nil, errors.Errorf("BOOL array response too short: got %d bytes for %d elements from bit %d", len(data), numElements, firstBit)
			}
			items = append(items, spiValues.NewPlcBOOL((data[bit/8]>>(bit%8))&0x01 != 0))
		}
		return spiValues.NewPlcList(items), nil
	default:
		readBuffer := utils.NewReadBufferByteBased(data)
		items := make([]apiValues.PlcValue, 0, numElements)
		for i := range numElements {
			item, err := readWriteModel.DataItemParseWithBuffer(ctx, readBuffer, dataProtocolId, controllerType, stringLength)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing element %d of type %s", i, tag.GetDataType())
			}
			items = append(items, item)
		}
		return spiValues.NewPlcList(items), nil
	}
}

// mapPlcErrorCode maps S7 header error class/code pairs to PLC4X response codes (ported from plc4j).
func mapPlcErrorCode(localLog zerolog.Logger, errorClass uint8, errorCode uint8) apiModel.PlcResponseCode {
	// 0x81/4 means PUT/GET is disabled. An S7-300 reports the same refusal as 0x83/4 (GH-599),
	// and error class 0x85 covers further access-denied situations.
	if (errorClass == 0x81 && errorCode == 4) || (errorClass == 0x83 && errorCode == 4) {
		localLog.Warn().Msg("Got an error response from the PLC. This particular response code usually indicates " +
			"that PUT/GET is not enabled on the PLC.")
		return apiModel.PlcResponseCode_ACCESS_DENIED
	}
	if errorClass == 0x85 {
		return apiModel.PlcResponseCode_ACCESS_DENIED
	}
	localLog.Warn().
		Uint8("errorClass", errorClass).
		Uint8("errorCode", errorCode).
		Msg("Got an unknown error response from the PLC. " +
			"We probably need to implement explicit handling for this, so please file a bug-report " +
			"on https://github.com/apache/plc4x/issues and ideally attach a WireShark dump " +
			"containing a capture of the communication.")
	return apiModel.PlcResponseCode_INTERNAL_ERROR
}

// Helper to convert the return codes returned from the S7 into one of our standard
func decodeResponseCode(dataTransportErrorCode readWriteModel.DataTransportErrorCode) apiModel.PlcResponseCode {
	switch dataTransportErrorCode {
	case readWriteModel.DataTransportErrorCode_OK:
		return apiModel.PlcResponseCode_OK
	case readWriteModel.DataTransportErrorCode_NOT_FOUND, readWriteModel.DataTransportErrorCode_INVALID_ADDRESS:
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	case readWriteModel.DataTransportErrorCode_DATA_TYPE_NOT_SUPPORTED:
		return apiModel.PlcResponseCode_INVALID_DATATYPE
	case readWriteModel.DataTransportErrorCode_ACCESS_DENIED:
		return apiModel.PlcResponseCode_ACCESS_DENIED
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
}
