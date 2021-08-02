/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/rs/zerolog/log"
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

type readRequestFactory func(fields map[string]model.PlcField, fieldNames []string, reader spi.PlcReader, readRequestInterceptor ReadRequestInterceptor) model.PlcReadRequest
type writeRequestFactory func(fields map[string]model.PlcField, fieldNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor WriteRequestInterceptor) model.PlcWriteRequest

type readResponseFactory func(request model.PlcReadRequest, responseCodes map[string]model.PlcResponseCode, values map[string]values.PlcValue) model.PlcReadResponse
type writeResponseFactory func(request model.PlcWriteRequest, responseCodes map[string]model.PlcResponseCode) model.PlcWriteResponse

type SingleItemRequestInterceptor struct {
	readRequestFactory   readRequestFactory
	writeRequestFactory  writeRequestFactory
	readResponseFactory  readResponseFactory
	writeResponseFactory writeResponseFactory
}

func NewSingleItemRequestInterceptor(readRequestFactory readRequestFactory, writeRequestFactory writeRequestFactory, readResponseFactory readResponseFactory, writeResponseFactory writeResponseFactory) SingleItemRequestInterceptor {
	return SingleItemRequestInterceptor{readRequestFactory, writeRequestFactory, readResponseFactory, writeResponseFactory}
}

func (m SingleItemRequestInterceptor) InterceptReadRequest(readRequest model.PlcReadRequest) []model.PlcReadRequest {
	// If this request just has one field, go the shortcut
	if len(readRequest.GetFieldNames()) == 1 {
		log.Debug().Msg("We got only one request, no splitting required")
		return []model.PlcReadRequest{readRequest}
	}
	log.Trace().Msg("Splitting requests")
	// In all other cases, create a new read request containing only one item
	var readRequests []model.PlcReadRequest
	for _, fieldName := range readRequest.GetFieldNames() {
		log.Debug().Str("fieldName", fieldName).Msg("Splitting into own request")
		field := readRequest.GetField(fieldName)
		subReadRequest := m.readRequestFactory(
			map[string]model.PlcField{fieldName: field},
			[]string{fieldName},
			readRequest.(ReaderExposer).GetReader(),
			readRequest.(ReadRequestInterceptorExposer).GetReadRequestInterceptor(),
		)
		readRequests = append(readRequests, subReadRequest)
	}
	return readRequests
}

func (m SingleItemRequestInterceptor) ProcessReadResponses(readRequest model.PlcReadRequest, readResults []model.PlcReadRequestResult) model.PlcReadRequestResult {
	if len(readResults) == 1 {
		log.Debug().Msg("We got only one response, no merging required")
		return readResults[0]
	}
	log.Trace().Msg("Merging requests")
	responseCodes := map[string]model.PlcResponseCode{}
	val := map[string]values.PlcValue{}
	var err error = nil
	for _, readResult := range readResults {
		if readResult.Err != nil {
			log.Debug().Err(readResult.Err).Msgf("Error during read")
			if err == nil {
				// Lazy initialization of multi error
				err = utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{readResult.Err}}
			} else {
				multiError := err.(utils.MultiError)
				multiError.Errors = append(multiError.Errors, readResult.Err)
			}
		} else if readResult.Response != nil {
			if len(readResult.Response.GetRequest().GetFieldNames()) > 1 {
				log.Error().Int("numberOfFields", len(readResult.Response.GetRequest().GetFieldNames())).Msg("We should only get 1")
			}
			for _, fieldName := range readResult.Response.GetRequest().GetFieldNames() {
				responseCodes[fieldName] = readResult.Response.GetResponseCode(fieldName)
				val[fieldName] = readResult.Response.GetValue(fieldName)
			}
		}
	}
	return model.PlcReadRequestResult{
		Request:  readRequest,
		Response: m.readResponseFactory(readRequest, responseCodes, val),
		Err:      err,
	}
}

func (m SingleItemRequestInterceptor) InterceptWriteRequest(writeRequest model.PlcWriteRequest) []model.PlcWriteRequest {
	// If this request just has one field, go the shortcut
	if len(writeRequest.GetFieldNames()) == 1 {
		log.Debug().Msg("We got only one request, no splitting required")
		return []model.PlcWriteRequest{writeRequest}
	}
	log.Trace().Msg("Splitting requests")
	// In all other cases, create a new write request containing only one item
	var writeRequests []model.PlcWriteRequest
	for _, fieldName := range writeRequest.GetFieldNames() {
		log.Debug().Str("fieldName", fieldName).Msg("Splitting into own request")
		field := writeRequest.GetField(fieldName)
		subWriteRequest := m.writeRequestFactory(
			map[string]model.PlcField{fieldName: field},
			[]string{fieldName},
			map[string]values.PlcValue{fieldName: writeRequest.GetValue(fieldName)},
			writeRequest.(WriterExposer).GetWriter(),
			writeRequest.(WriteRequestInterceptorExposer).GetWriteRequestInterceptor(),
		)
		writeRequests = append(writeRequests, subWriteRequest)
	}
	return writeRequests
}

func (m SingleItemRequestInterceptor) ProcessWriteResponses(writeRequest model.PlcWriteRequest, writeResults []model.PlcWriteRequestResult) model.PlcWriteRequestResult {
	if len(writeResults) == 1 {
		log.Debug().Msg("We got only one response, no merging required")
		return writeResults[0]
	}
	log.Trace().Msg("Merging requests")
	responseCodes := map[string]model.PlcResponseCode{}
	var err error = nil
	for _, writeResult := range writeResults {
		if writeResult.Err != nil {
			log.Debug().Err(writeResult.Err).Msgf("Error during write")
			if err == nil {
				// Lazy initialization of multi error
				err = utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{writeResult.Err}}
			} else {
				multiError := err.(utils.MultiError)
				multiError.Errors = append(multiError.Errors, writeResult.Err)
			}
		} else if writeResult.Response != nil {
			if len(writeResult.Response.GetRequest().GetFieldNames()) > 1 {
				log.Error().Int("numberOfFields", len(writeResult.Response.GetRequest().GetFieldNames())).Msg("We should only get 1")
			}
			for _, fieldName := range writeResult.Response.GetRequest().GetFieldNames() {
				responseCodes[fieldName] = writeResult.Response.GetResponseCode(fieldName)
			}
		}
	}
	return model.PlcWriteRequestResult{
		Request:  writeRequest,
		Response: m.writeResponseFactory(writeRequest, responseCodes),
		Err:      err,
	}
}
