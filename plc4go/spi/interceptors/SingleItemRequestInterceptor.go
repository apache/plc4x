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

package interceptors

import (
	"context"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ReaderExposer interface {
	GetReader() spi.PlcReader
}

type WriterExposer interface {
	GetWriter() spi.PlcWriter
}

type ReadRequestInterceptorExposer interface {
	GetReadRequestInterceptor() ReadRequestInterceptor
}

type WriteRequestInterceptorExposer interface {
	GetWriteRequestInterceptor() WriteRequestInterceptor
}

type readRequestFactory func(
	tags map[string]apiModel.PlcTag,
	tagNames []string,
	reader spi.PlcReader,
	readRequestInterceptor ReadRequestInterceptor,
) apiModel.PlcReadRequest
type writeRequestFactory func(
	tags map[string]apiModel.PlcTag,
	tagNames []string,
	values map[string]values.PlcValue,
	writer spi.PlcWriter,
	writeRequestInterceptor WriteRequestInterceptor,
) apiModel.PlcWriteRequest

type readResponseFactory func(
	request apiModel.PlcReadRequest,
	responseCodes map[string]apiModel.PlcResponseCode,
	values map[string]values.PlcValue,
) apiModel.PlcReadResponse
type writeResponseFactory func(
	request apiModel.PlcWriteRequest,
	responseCodes map[string]apiModel.PlcResponseCode,
) apiModel.PlcWriteResponse

type SingleItemRequestInterceptor struct {
	readRequestFactory   readRequestFactory
	writeRequestFactory  writeRequestFactory
	readResponseFactory  readResponseFactory
	writeResponseFactory writeResponseFactory

	log zerolog.Logger
}

func NewSingleItemRequestInterceptor(readRequestFactory readRequestFactory, writeRequestFactory writeRequestFactory, readResponseFactory readResponseFactory, writeResponseFactory writeResponseFactory, _options ...options.WithOption) SingleItemRequestInterceptor {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return SingleItemRequestInterceptor{
		readRequestFactory:   readRequestFactory,
		writeRequestFactory:  writeRequestFactory,
		readResponseFactory:  readResponseFactory,
		writeResponseFactory: writeResponseFactory,
		log:                  customLogger,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=interceptedPlcReadRequestResult
type interceptedPlcReadRequestResult struct {
	Request  apiModel.PlcReadRequest
	Response apiModel.PlcReadResponse
	Err      error
}

func (d *interceptedPlcReadRequestResult) GetRequest() apiModel.PlcReadRequest {
	return d.Request
}

func (d *interceptedPlcReadRequestResult) GetResponse() apiModel.PlcReadResponse {
	return d.Response
}

func (d *interceptedPlcReadRequestResult) GetErr() error {
	return d.Err
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=interceptedPlcWriteRequestResult
type interceptedPlcWriteRequestResult struct {
	Request  apiModel.PlcWriteRequest
	Response apiModel.PlcWriteResponse
	Err      error
}

func (d *interceptedPlcWriteRequestResult) GetRequest() apiModel.PlcWriteRequest {
	return d.Request
}

func (d *interceptedPlcWriteRequestResult) GetResponse() apiModel.PlcWriteResponse {
	return d.Response
}

func (d *interceptedPlcWriteRequestResult) GetErr() error {
	return d.Err
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (m SingleItemRequestInterceptor) InterceptReadRequest(ctx context.Context, readRequest apiModel.PlcReadRequest) []apiModel.PlcReadRequest {
	if readRequest == nil || len(readRequest.GetTagNames()) == 0 {
		return nil
	}
	// If this request just has one tag, go the shortcut
	if len(readRequest.GetTagNames()) == 1 {
		m.log.Debug().Msg("We got only one request, no splitting required")
		return []apiModel.PlcReadRequest{readRequest}
	}
	m.log.Trace().Msg("Splitting requests")
	// In all other cases, create a new read request containing only one item
	var readRequests []apiModel.PlcReadRequest
	for _, tagName := range readRequest.GetTagNames() {
		if err := ctx.Err(); err != nil {
			m.log.Warn().Err(err).Msg("aborting early")
			return nil
		}
		m.log.Debug().Str("tagName", tagName).Msg("Splitting into own request")
		tag := readRequest.GetTag(tagName)
		subReadRequest := m.readRequestFactory(
			map[string]apiModel.PlcTag{tagName: tag},
			[]string{tagName},
			readRequest.(ReaderExposer).GetReader(),
			readRequest.(ReadRequestInterceptorExposer).GetReadRequestInterceptor(),
		)
		readRequests = append(readRequests, subReadRequest)
	}
	return readRequests
}

func (m SingleItemRequestInterceptor) ProcessReadResponses(ctx context.Context, readRequest apiModel.PlcReadRequest, readResults []apiModel.PlcReadRequestResult) apiModel.PlcReadRequestResult {
	if len(readResults) == 1 {
		m.log.Debug().Msg("We got only one response, no merging required")
		return readResults[0]
	}
	m.log.Trace().Msg("Merging requests")
	responseCodes := map[string]apiModel.PlcResponseCode{}
	val := map[string]values.PlcValue{}
	var collectedErrors []error
	for _, readResult := range readResults {
		if err := ctx.Err(); err != nil {
			m.log.Warn().Err(err).Msg("aborting early")
			collectedErrors = append(collectedErrors, err)
			break
		}
		if err := readResult.GetErr(); err != nil {
			m.log.Debug().Err(err).Msg("Error during read")
			collectedErrors = append(collectedErrors, err)
		} else if response := readResult.GetResponse(); response != nil {
			request := response.GetRequest()
			if len(request.GetTagNames()) > 1 {
				m.log.Error().Int("numberOfTags", len(request.GetTagNames())).Msg("We should only get 1")
			}
			for _, tagName := range request.GetTagNames() {
				responseCodes[tagName] = response.GetResponseCode(tagName)
				val[tagName] = response.GetValue(tagName)
			}
		}
	}
	var err error
	if len(collectedErrors) > 0 {
		err = utils.MultiError{MainError: errors.New("error aggregating"), Errors: collectedErrors}
	}
	return &interceptedPlcReadRequestResult{
		Request:  readRequest,
		Response: m.readResponseFactory(readRequest, responseCodes, val),
		Err:      err,
	}
}

func (m SingleItemRequestInterceptor) InterceptWriteRequest(ctx context.Context, writeRequest apiModel.PlcWriteRequest) []apiModel.PlcWriteRequest {
	if writeRequest == nil {
		return nil
	}
	// If this request just has one tag, go the shortcut
	if len(writeRequest.GetTagNames()) == 1 {
		m.log.Debug().Msg("We got only one request, no splitting required")
		return []apiModel.PlcWriteRequest{writeRequest}
	}
	m.log.Trace().Msg("Splitting requests")
	// In all other cases, create a new write request containing only one item
	var writeRequests []apiModel.PlcWriteRequest
	for _, tagName := range writeRequest.GetTagNames() {
		if err := ctx.Err(); err != nil {
			m.log.Warn().Err(err).Msg("aborting early")
			return nil
		}
		m.log.Debug().Str("tagName", tagName).Msg("Splitting into own request")
		tag := writeRequest.GetTag(tagName)
		subWriteRequest := m.writeRequestFactory(
			map[string]apiModel.PlcTag{tagName: tag},
			[]string{tagName},
			map[string]values.PlcValue{tagName: writeRequest.GetValue(tagName)},
			writeRequest.(WriterExposer).GetWriter(),
			writeRequest.(WriteRequestInterceptorExposer).GetWriteRequestInterceptor(),
		)
		writeRequests = append(writeRequests, subWriteRequest)
	}
	return writeRequests
}

func (m SingleItemRequestInterceptor) ProcessWriteResponses(ctx context.Context, writeRequest apiModel.PlcWriteRequest, writeResults []apiModel.PlcWriteRequestResult) apiModel.PlcWriteRequestResult {
	if len(writeResults) == 1 {
		m.log.Debug().Msg("We got only one response, no merging required")
		return writeResults[0]
	}
	m.log.Trace().Msg("Merging requests")
	responseCodes := map[string]apiModel.PlcResponseCode{}
	var collectedErrors []error
	for _, writeResult := range writeResults {
		if err := ctx.Err(); err != nil {
			m.log.Warn().Err(err).Msg("aborting early")
			collectedErrors = append(collectedErrors, err)
			break
		}
		if err := writeResult.GetErr(); err != nil {
			m.log.Debug().Err(err).Msg("Error during write")
			collectedErrors = append(collectedErrors, err)
		} else if writeResult.GetResponse() != nil {
			if len(writeResult.GetResponse().GetRequest().GetTagNames()) > 1 {
				m.log.Error().Int("numberOfTags", len(writeResult.GetResponse().GetRequest().GetTagNames())).Msg("We should only get 1")
			}
			for _, tagName := range writeResult.GetResponse().GetRequest().GetTagNames() {
				responseCodes[tagName] = writeResult.GetResponse().GetResponseCode(tagName)
			}
		}
	}
	var err error
	if len(collectedErrors) > 0 {
		err = utils.MultiError{MainError: errors.New("while aggregating results"), Errors: collectedErrors}
	}
	return &interceptedPlcWriteRequestResult{
		Request:  writeRequest,
		Response: m.writeResponseFactory(writeRequest, responseCodes),
		Err:      err,
	}
}
