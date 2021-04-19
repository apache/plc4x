//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

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

type ReadRequestInterceptorExposer interface {
	GetReadRequestInterceptor() ReadRequestInterceptor
}

type requestFactory func(fields map[string]model.PlcField, fieldNames []string, reader spi.PlcReader, readRequestInterceptor ReadRequestInterceptor) model.PlcReadRequest

type responseFactory func(request model.PlcReadRequest, responseCodes map[string]model.PlcResponseCode, values map[string]values.PlcValue) model.PlcReadResponse

type SingleItemRequestInterceptor struct {
	requestFactory  requestFactory
	responseFactory responseFactory
}

func NewSingleItemRequestInterceptor(requestFactory requestFactory, responseFactory responseFactory) SingleItemRequestInterceptor {
	return SingleItemRequestInterceptor{requestFactory, responseFactory}
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
		subReadRequest := m.requestFactory(
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
		Response: m.responseFactory(readRequest, responseCodes, val),
		Err:      err,
	}
}

func (m SingleItemRequestInterceptor) InterceptWriteRequest(writeRequest model.PlcWriteRequest) []model.PlcWriteRequest {
	return []model.PlcWriteRequest{writeRequest}
}

func (m SingleItemRequestInterceptor) ProcessWriteResponses(writeRequest model.PlcWriteRequest, writeResponses []model.PlcWriteRequestResult) model.PlcWriteRequestResult {
	// TODO: unfinished implementation
	return model.PlcWriteRequestResult{}
}
