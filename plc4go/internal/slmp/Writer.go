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

package slmp

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Writer turns plc4x write requests into SLMP Batch Write (0x1401) commands in word units. It is a
// structural mirror of Reader - one frame per tag, per-tag failure isolation - because plc4j's
// onWrite is a structural mirror of its onRead.
type Writer struct {
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	configuration Configuration

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewWriter(messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, configuration Configuration, _options ...options.WithOption) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,

		log: customLogger,
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
			tag, err := castToSlmpTagFromPlcTag(writeRequest.GetTag(tagName))
			if err != nil {
				m.log.Debug().Err(err).Str("tagName", tagName).Msg("not an slmp tag")
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

// writeSingleTag encodes one value and sends one Batch Write for it.
func (m *Writer) writeSingleTag(ctx context.Context, tagName string, tag PlcTag, value apiValues.PlcValue) apiModel.PlcResponseCode {
	payload, err := tag.GetDataType().Encode(value, tag.GetQuantity())
	if err != nil {
		// Nothing is sent for a value that can't be encoded, so the device is left alone.
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("value can't be encoded for this tag")
		return apiModel.PlcResponseCode_INVALID_DATA
	}
	// Invariant DataType.Encode guarantees: len(payload) == numberOfPoints * 2, so the frame's
	// announced point count and its payload length agree. The generated serializer derives the
	// payload length from numberOfPoints, so a disagreement would go out as a frame no device
	// accepts.
	requestCtx, cancelRequest := context.WithTimeout(ctx, m.configuration.requestTimeout)
	defer cancelRequest()

	frame := newRequestFrame(m.configuration, commandBatchWrite,
		readWriteModel.NewSlmpWriteRequest(tag.GetDeviceNumber(), tag.GetDeviceCode(), tag.GetNumberOfPoints(), payload))

	response, err := sendTransactedRequestAndWait(requestCtx, m.log, m.messageCodec, m.tm, "write", frame)
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("error writing tag")
		if utils.IsTimeoutError(err) && !errors.Is(ctx.Err(), context.Canceled) {
			// A write that timed out may still have been applied by the device - the driver cannot
			// tell, because a 3E frame has no correlation id - so a caller that retries it can
			// apply it twice. plc4j carries the same caveat.
			return apiModel.PlcResponseCode_REQUEST_TIMEOUT
		}
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
	return m.decodeWriteResponse(tagName, response)
}

// decodeWriteResponse maps one Batch Write response onto a response code. Ported from plc4j's
// SlmpResponseMapper.mapWriteTag.
func (m *Writer) decodeWriteResponse(tagName string, response readWriteModel.SlmpResponseFrame3E) apiModel.PlcResponseCode {
	if endCode := response.GetEndCode(); endCode != endCodeNormalCompletion {
		m.log.Warn().
			Str("tagName", tagName).
			Str("endCode", formatEndCode(endCode)).
			Msg("SLMP device answered the write with an error end code")
		return apiModel.PlcResponseCode_REMOTE_ERROR
	}
	if payloadLength := len(response.GetResponseData()); payloadLength > 0 {
		// SH-080008: a Batch Write success response carries no data. A payload here is the signature
		// of a mis-attributed response - a late read answer arriving after its own request timed out,
		// say - so reporting OK would claim a write happened on no evidence.
		m.log.Warn().
			Str("tagName", tagName).
			Int("payloadLength", payloadLength).
			Msg("SLMP write success frame unexpectedly carried a payload; treating it as a possibly mis-attributed response")
		return apiModel.PlcResponseCode_REMOTE_ERROR
	}
	return apiModel.PlcResponseCode_OK
}
