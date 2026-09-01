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

package umas

import (
	"context"
	"math"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// singleVariableCount is the variableCount field of a read or write request. Both this driver
	// and plc4j send one variable per request, even for a multi-tag plc4x request: nothing has ever
	// been captured that says how a PLC answers several at once, and the response is a single
	// undelimited block of bytes which would be impossible to split apart afterwards.
	singleVariableCount = uint8(1)
	// stringReadSizeIndex is the dataSizeIndex a STRING read uses. STRING's own request size is 17,
	// which doesn't fit the 4 bit field, so the read asks for a byte array of single bytes instead
	// and says how many with the arrayLength field.
	stringReadSizeIndex = uint8(1)
	// arrayReadFlag is the isArray field of a reference which reads a byte array.
	arrayReadFlag = uint8(1)
	// scalarReadFlag is the isArray field of a reference which reads one fixed width value.
	scalarReadFlag = uint8(0)
)

// Reader turns plc4x read requests into UMAS variable reads (FC 0x22).
//
// UMAS answers one request at a time - plc4j's UmasConnection.getMaxConcurrentRequests returns 1 -
// so the tags of a request are read one after the other, each in its own exchange.
type Reader struct {
	requester *requester
	session   *session

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewReader(requester *requester, session *session, _options ...options.WithOption) *Reader {
	return &Reader{
		requester: requester,
		session:   session,
		log:       options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil,
					errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		responseCodes := map[string]apiModel.PlcResponseCode{}
		plcValues := map[string]apiValues.PlcValue{}
		for _, tagName := range readRequest.GetTagNames() {
			tag, ok := readRequest.GetTag(tagName).(PlcTag)
			if !ok {
				m.log.Warn().Str("tagName", tagName).Type("tag", readRequest.GetTag(tagName)).Msg("Not a UMAS tag")
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			code, value := m.readSingleTag(ctx, tagName, tag)
			responseCodes[tagName] = code
			if value != nil {
				plcValues[tagName] = value
			}
		}
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest,
			spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil))
	})
	return result
}

// readSingleTag reads one symbol. A failure is reported as a response code for that one tag rather
// than failing the whole request, the way plc4j's readSingleTag does it.
func (m *Reader) readSingleTag(ctx context.Context, tagName string, tag PlcTag) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	symbol, ok := m.session.lookupSymbol(tag.GetSymbolicAddress())
	if !ok {
		m.log.Warn().
			Str("tagName", tagName).
			Str("symbolicAddress", tag.GetSymbolicAddress()).
			Msg("The symbol is not in the symbol table")
		return apiModel.PlcResponseCode_NOT_FOUND, nil
	}

	reference, err := m.buildReadReference(symbol)
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("Can't build a read reference for the symbol")
		return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
	}

	item, err := m.requester.exchange(ctx, "ReadVariable("+tagName+")",
		readWriteModel.NewUmasPDUReadVariableRequest(
			m.session.getPairingKey(), m.session.getProjectCrc(), singleVariableCount,
			[]readWriteModel.VariableReadRequestReference{reference}))
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("error reading tag")
		if isTimeout(err) && ctx.Err() == nil {
			// The exchange ran out its own request timeout, the caller is still waiting.
			return apiModel.PlcResponseCode_REQUEST_TIMEOUT, nil
		}
		return apiModel.PlcResponseCode_REMOTE_ERROR, nil
	}
	response, ok := item.(readWriteModel.UmasPDUReadVariableResponse)
	if !ok {
		if isUmasErrorResponse(item) {
			m.log.Warn().Str("tagName", tagName).Msg("The PLC refused the read")
			return apiModel.PlcResponseCode_REMOTE_ERROR, nil
		}
		m.log.Warn().Str("tagName", tagName).Type("responseType", item).Msg("Unexpected read response type")
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}

	value, err := decodeReadResponse(ctx, symbol.GetDataType(), response.GetBlock())
	if err != nil {
		m.log.Warn().Err(err).
			Str("tagName", tagName).
			Uint16("dataType", symbol.GetDataType()).
			Msg("Can't decode the read response")
		return apiModel.PlcResponseCode_INVALID_DATA, nil
	}
	return apiModel.PlcResponseCode_OK, value
}

// buildReadReference spells out which symbol to read. Ported from plc4j's buildReadReference.
//
// The 32 bit offset of a symbol is split over the reference's two offset fields, and a read splits
// it differently than a write does: the low 8 bits go into the 8 bit offset field and everything
// above them into the 16 bit baseOffset. plc4j documents both splits without saying where they come
// from other than "this is what works", so they are reproduced as they are.
func (m *Reader) buildReadReference(symbol readWriteModel.UmasUnlocatedVariableReference) (readWriteModel.VariableReadRequestReference, error) {
	symbolOffset := symbol.GetOffset()
	baseOffset := symbolOffset >> 8
	if baseOffset > math.MaxUint16 {
		// The reference has no room for it, and truncating would read some other symbol's memory.
		return nil, errors.Errorf("the symbol offset 0x%08X doesn't fit the offset fields of a read reference", symbolOffset)
	}
	offset := uint8(symbolOffset & 0xFF)

	if isStringType(symbol.GetDataType()) {
		// A STRING is read as a byte array of its buffer size, because the size index field is only
		// 4 bits wide and STRING's request size is 17.
		stringSize := m.session.stringBufferSize(symbol)
		return readWriteModel.NewVariableReadRequestReference(
			arrayReadFlag, stringReadSizeIndex, symbol.GetBlock(),
			uint16(baseOffset), offset, &stringSize), nil
	}
	return readWriteModel.NewVariableReadRequestReference(
		scalarReadFlag, dataSizeIndexFor(symbol.GetDataType()), symbol.GetBlock(),
		uint16(baseOffset), offset, nil), nil
}
