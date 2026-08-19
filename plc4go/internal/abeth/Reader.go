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

package abeth

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const (
	// df1SourceAddress is the DF1 source address the driver announces itself with. The legacy
	// driver hard-coded a 5 here and plc4j's AbEthConnection kept it.
	df1SourceAddress = uint8(5)
	// df1RequestStatus is the status byte of an outgoing request, which is always zero: the field
	// carries a status only on the way back.
	df1RequestStatus = uint8(0)
	// df1SubElementNumber is the sub element a read addresses. plc4j never addresses one, even
	// though the address syntax once had a slot for it.
	df1SubElementNumber = uint8(0)
)

// emptySenderContext is the sender context a read request carries. plc4j's AbEthConnection sends
// all zeros for everything but the connection request.
var emptySenderContext = []uint8{0, 0, 0, 0, 0, 0, 0, 0}

// Reader turns plc4x read requests into DF1 "protected typed logical read" commands.
//
// ab-eth answers one request at a time - the legacy driver capped concurrent requests at one and
// plc4j's AbEthConnection still returns 1 from getMaxConcurrentRequests - so the tags of a request
// are read one after the other and each read is matched to its response by transaction counter.
type Reader struct {
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	configuration Configuration
	session       *session

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewReader(messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, configuration Configuration, session *session, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,
		session:       session,

		log: customLogger,
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

// readSingleTag sends one read for one tag and decodes the answer. A failure is reported as a
// response code for that one tag rather than failing the whole request, the way plc4j does it.
func (m *Reader) readSingleTag(ctx context.Context, tagName string, tag PlcTag) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	requestCtx, cancelRequest := context.WithTimeout(ctx, m.configuration.requestTimeout)
	defer cancelRequest()

	transactionCounter := m.session.nextTransactionCounter()
	request := readWriteModel.NewCIPEncapsulationReadRequest(
		m.session.getSessionHandle(), 0, emptySenderContext, 0,
		readWriteModel.NewDF1CommandRequestMessage(
			m.configuration.station, df1SourceAddress, df1RequestStatus, transactionCounter,
			readWriteModel.NewDF1RequestProtectedTypedLogicalRead(
				tag.GetByteSize(), tag.GetFileNumber(), tag.GetFileType().GetTypeCode(),
				tag.GetElementNumber(), df1SubElementNumber,
			),
		),
	)

	message, err := sendTransactedRequestAndWait(requestCtx, m.log, m.messageCodec, m.tm, "read", request,
		func(message spi.Message) bool {
			readResponse, ok := message.(readWriteModel.CIPEncapsulationReadResponse)
			if !ok {
				return false
			}
			// Responses are matched to requests by transaction counter, which is the only thing
			// tying them together: the session handle is the same on every packet of a connection.
			return readResponse.GetResponse().GetTransactionCounter() == transactionCounter
		})
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("error reading tag")
		if requestCtx.Err() != nil && ctx.Err() == nil {
			return apiModel.PlcResponseCode_REQUEST_TIMEOUT, nil
		}
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	readResponse, ok := message.(readWriteModel.CIPEncapsulationReadResponse)
	if !ok {
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	return m.decodeReadResponse(tagName, tag, readResponse)
}

// decodeReadResponse turns the payload bytes of a DF1 response into a plc4x value. Lifted from
// plc4j's AbEthConnection.decodeReadResponse, which in turn lifted it from the legacy
// AbEthProtocolLogic: what the bytes mean depends entirely on the tag's file type.
func (m *Reader) decodeReadResponse(tagName string, tag PlcTag, readResponse readWriteModel.CIPEncapsulationReadResponse) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	response := readResponse.GetResponse()
	// plc4j maps every non-zero DF1 status to NOT_FOUND. The DF1 status byte does distinguish more
	// than that, but nothing in the driver ever looked at the detail.
	if response.GetStatus() != 0 {
		m.log.Debug().
			Str("tagName", tagName).
			Uint8("status", response.GetStatus()).
			Msg("PLC answered the read with a non-zero DF1 status")
		return apiModel.PlcResponseCode_NOT_FOUND, nil
	}
	logicalRead, ok := response.(readWriteModel.DF1CommandResponseMessageProtectedTypedLogicalRead)
	if !ok {
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	data := logicalRead.GetData()

	switch tag.GetFileType() {
	case FileTypeInteger:
		// A single byte is the value itself; anything wider comes back as the list of its bytes.
		// That is what plc4j does (DefaultPlcValueHandler.of(tag, List<Short>) on an INT tag), even
		// though an integer file really holds two byte elements.
		if len(data) == 0 {
			return apiModel.PlcResponseCode_INVALID_DATA, nil
		}
		if len(data) == 1 {
			return apiModel.PlcResponseCode_OK, spiValues.NewPlcINT(int16(data[0]))
		}
		elements := make([]apiValues.PlcValue, 0, len(data))
		for _, dataByte := range data {
			elements = append(elements, spiValues.NewPlcINT(int16(dataByte)))
		}
		return apiModel.PlcResponseCode_OK, spiValues.NewPlcList(elements)
	case FileTypeWord:
		if len(data) < 2 {
			return apiModel.PlcResponseCode_INVALID_DATA, nil
		}
		// Deliberate deviation from plc4j, which sign-extends the two bytes and then hands the
		// result to PlcWORD - a value type whose lower bound is zero, so every negative word ends
		// up as a PlcInvalidTagException swallowed into a null value. A WORD is a bit pattern, so
		// the unsigned reassembly is what the declared value type actually means.
		return apiModel.PlcResponseCode_OK, spiValues.NewPlcWORD(uint16(data[0]) | uint16(data[1])<<8)
	case FileTypeDword:
		if len(data) < 4 {
			return apiModel.PlcResponseCode_INVALID_DATA, nil
		}
		return apiModel.PlcResponseCode_OK, spiValues.NewPlcDWORD(
			uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24)
	case FileTypeSinglebit:
		if len(data) < 2 {
			return apiModel.PlcResponseCode_INVALID_DATA, nil
		}
		// The two bytes are one integer element in little-endian order, so bits 0-7 live in the
		// first byte and bits 8-15 in the second.
		bitNumber := tag.GetBitNumber()
		dataByte := data[bitNumber/8]
		return apiModel.PlcResponseCode_OK, spiValues.NewPlcBOOL(dataByte&(1<<(bitNumber%8)) != 0)
	default:
		// plc4j logs this case and then returns OK with a null value, which tells the caller the
		// read worked while handing it nothing. Reporting UNSUPPORTED says what actually happened.
		m.log.Warn().
			Str("tagName", tagName).
			Stringer("fileType", tag.GetFileType()).
			Msg("Decoding of this file type is not implemented")
		return apiModel.PlcResponseCode_UNSUPPORTED, nil
	}
}
