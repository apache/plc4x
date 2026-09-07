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

// Writer turns plc4x write requests into UMAS variable writes (FC 0x23).
//
// Like the reader it writes one tag per exchange, because UMAS answers one request at a time.
type Writer struct {
	requester *requester
	session   *session

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewWriter(requester *requester, session *session, _options ...options.WithOption) *Writer {
	return &Writer{
		requester: requester,
		session:   session,
		log:       options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

func (m *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	m.log.Trace().Msg("Writing")
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil,
					errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		responseCodes := map[string]apiModel.PlcResponseCode{}
		for _, tagName := range writeRequest.GetTagNames() {
			tag, ok := writeRequest.GetTag(tagName).(PlcTag)
			if !ok {
				m.log.Warn().Str("tagName", tagName).Type("tag", writeRequest.GetTag(tagName)).Msg("Not a UMAS tag")
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			responseCodes[tagName] = m.writeSingleTag(ctx, tagName, tag, writeRequest.GetValue(tagName))
		}
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest,
			spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil))
	})
	return result
}

// writeSingleTag writes one symbol. A failure is reported as a response code for that one tag rather
// than failing the whole request, the way plc4j's writeSingleTag does it.
func (m *Writer) writeSingleTag(ctx context.Context, tagName string, tag PlcTag, value apiValues.PlcValue) apiModel.PlcResponseCode {
	symbol, ok := m.session.lookupSymbol(tag.GetSymbolicAddress())
	if !ok {
		m.log.Warn().
			Str("tagName", tagName).
			Str("symbolicAddress", tag.GetSymbolicAddress()).
			Msg("The symbol is not in the symbol table")
		return apiModel.PlcResponseCode_NOT_FOUND
	}

	data, err := encodeWriteValue(symbol.GetDataType(), value)
	if err != nil {
		m.log.Warn().Err(err).
			Str("tagName", tagName).
			Uint16("dataType", symbol.GetDataType()).
			Msg("Can't serialize the value for the symbol")
		return apiModel.PlcResponseCode_INVALID_DATATYPE
	}
	reference, err := buildWriteReference(ctx, symbol, data)
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("Can't build a write reference for the symbol")
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	}

	item, err := m.requester.exchange(ctx, "WriteVariable("+tagName+")",
		readWriteModel.NewUmasPDUWriteVariableRequest(
			m.session.getPairingKey(), m.session.getProjectCrc(), singleVariableCount,
			[]readWriteModel.VariableWriteRequestReference{reference}))
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("error writing tag")
		if isTimeout(err) && ctx.Err() == nil {
			// The exchange ran out its own request timeout, the caller is still waiting.
			return apiModel.PlcResponseCode_REQUEST_TIMEOUT
		}
		return apiModel.PlcResponseCode_REMOTE_ERROR
	}
	if _, ok := item.(readWriteModel.UmasPDUWriteVariableResponse); ok {
		return apiModel.PlcResponseCode_OK
	}
	if isUmasErrorResponse(item) {
		m.log.Warn().Str("tagName", tagName).Msg("The PLC refused the write")
		return apiModel.PlcResponseCode_REMOTE_ERROR
	}
	m.log.Warn().Str("tagName", tagName).Type("responseType", item).Msg("Unexpected write response type")
	return apiModel.PlcResponseCode_INTERNAL_ERROR
}

// buildWriteReference spells out which symbol to write and with what. Ported from plc4j's
// buildWriteReference.
//
// A write splits the symbol's 32 bit offset differently than a read does: the low 16 bits go into
// baseOffset and the high 16 into offset, where a read puts the low 8 bits into offset and the rest
// into baseOffset. plc4j documents the asymmetry as observed, with no explanation of it, so it is
// reproduced as it is.
func buildWriteReference(ctx context.Context, symbol readWriteModel.UmasUnlocatedVariableReference, data []byte) (readWriteModel.VariableWriteRequestReference, error) {
	symbolOffset := symbol.GetOffset()
	baseOffset := uint16(symbolOffset & 0xFFFF)
	offset := uint16((symbolOffset >> 16) & 0xFFFF)

	if len(data) > math.MaxUint16 {
		return nil, errors.Errorf("a write payload of %d bytes doesn't fit the length field of a write reference", len(data))
	}

	if isStringType(symbol.GetDataType()) {
		// Same reason as on the read side: STRING's request size of 17 doesn't fit the 4 bit size
		// index field, so the value goes over as a byte array of its own length.
		arrayLength := uint16(len(data))
		return readWriteModel.NewVariableWriteRequestReference(
			arrayReadFlag, stringReadSizeIndex, symbol.GetBlock(),
			baseOffset, offset, &arrayLength, data), nil
	}
	// For a fixed width type the payload length is implied by the size index; the model derives it
	// with the same writeSizeIndexToByteCount mapping the PLC uses, so a mismatch between the two
	// would serialize the wrong number of bytes.
	sizeIndex := dataSizeIndexFor(symbol.GetDataType())
	if expected := readWriteModel.WriteSizeIndexToByteCount(ctx, sizeIndex); int(expected) != len(data) {
		return nil, errors.Errorf("a %d byte payload doesn't match the %d bytes the size index %d declares",
			len(data), expected, sizeIndex)
	}
	return readWriteModel.NewVariableWriteRequestReference(
		scalarReadFlag, sizeIndex, symbol.GetBlock(), baseOffset, offset, nil, data), nil
}
