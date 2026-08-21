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

// Reader turns plc4x read requests into SLMP Batch Read (0x0401) commands in word units.
//
// A 3E device answers one request at a time and the frame carries no correlation id, so the tags of
// a request are read one after the other - one frame per tag - exactly as plc4j's SlmpConnection
// chains them through executeThrottled. There is no request optimizer, so two adjacent tags are two
// frames even when one Batch Read could have covered both.
type Reader struct {
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	configuration Configuration

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewReader(messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, configuration Configuration, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,

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
			tag, err := castToSlmpTagFromPlcTag(readRequest.GetTag(tagName))
			if err != nil {
				m.log.Debug().Err(err).Str("tagName", tagName).Msg("not an slmp tag")
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			code, value := m.readSingleTag(ctx, tagName, tag)
			responseCodes[tagName] = code
			if value != nil {
				plcValues[tagName] = value
			}
		}
		// Partial-failure isolation, the way plc4j does it: a tag that failed carries its own code
		// and the rest of the request still comes back.
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest,
			spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil))
	})
	return result
}

// readSingleTag sends one Batch Read for one tag and decodes the answer.
func (m *Reader) readSingleTag(ctx context.Context, tagName string, tag PlcTag) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	requestCtx, cancelRequest := context.WithTimeout(ctx, m.configuration.requestTimeout)
	defer cancelRequest()

	frame := newRequestFrame(m.configuration, commandBatchRead,
		readWriteModel.NewSlmpReadRequest(tag.GetDeviceNumber(), tag.GetDeviceCode(), tag.GetNumberOfPoints()))

	response, err := sendTransactedRequestAndWait(requestCtx, m.log, m.messageCodec, m.tm, "read", frame)
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("error reading tag")
		if utils.IsTimeoutError(err) && ctx.Err() == nil {
			// plc4j reports a timeout as REMOTE_ERROR because its PlcResponseCode enum is what it
			// is; plc4go has a code that says what actually happened.
			return apiModel.PlcResponseCode_REQUEST_TIMEOUT, nil
		}
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	return m.decodeReadResponse(tagName, tag, response)
}

// decodeReadResponse maps one Batch Read response onto a response code and a value. Ported from
// plc4j's SlmpResponseMapper.mapTag.
func (m *Reader) decodeReadResponse(tagName string, tag PlcTag, response readWriteModel.SlmpResponseFrame3E) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	if endCode := response.GetEndCode(); endCode != endCodeNormalCompletion {
		// On abnormal completion the payload carries error information rather than device data; the
		// driver has no use for it beyond the log line, the way plc4j's mapper logs the end code.
		m.log.Warn().
			Str("tagName", tagName).
			Str("endCode", formatEndCode(endCode)).
			Msg("SLMP device answered the read with an error end code")
		return apiModel.PlcResponseCode_REMOTE_ERROR, nil
	}
	value, err := tag.GetDataType().Decode(response.GetResponseData(), tag.GetQuantity())
	if err != nil {
		m.log.Warn().Err(err).Str("tagName", tagName).Msg("SLMP response can't be decoded for this tag")
		return apiModel.PlcResponseCode_INVALID_DATA, nil
	}
	return apiModel.PlcResponseCode_OK, value
}
