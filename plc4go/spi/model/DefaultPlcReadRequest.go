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

package model

import (
	"context"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"time"
)

type DefaultPlcReadRequestBuilder struct {
	reader                 spi.PlcReader
	fieldHandler           spi.PlcFieldHandler
	queries                map[string]string
	queryNames             []string
	fields                 map[string]model.PlcField
	fieldNames             []string
	readRequestInterceptor interceptors.ReadRequestInterceptor
}

func NewDefaultPlcReadRequestBuilder(fieldHandler spi.PlcFieldHandler, reader spi.PlcReader) *DefaultPlcReadRequestBuilder {
	return NewDefaultPlcReadRequestBuilderWithInterceptor(fieldHandler, reader, nil)
}

func NewDefaultPlcReadRequestBuilderWithInterceptor(fieldHandler spi.PlcFieldHandler, reader spi.PlcReader, readRequestInterceptor interceptors.ReadRequestInterceptor) *DefaultPlcReadRequestBuilder {
	return &DefaultPlcReadRequestBuilder{
		reader:                 reader,
		fieldHandler:           fieldHandler,
		queries:                map[string]string{},
		queryNames:             make([]string, 0),
		fields:                 map[string]model.PlcField{},
		fieldNames:             make([]string, 0),
		readRequestInterceptor: readRequestInterceptor,
	}
}

func (m *DefaultPlcReadRequestBuilder) AddQuery(name string, query string) model.PlcReadRequestBuilder {
	m.queryNames = append(m.queryNames, name)
	m.queries[name] = query
	return m
}

func (m *DefaultPlcReadRequestBuilder) AddField(name string, field model.PlcField) model.PlcReadRequestBuilder {
	m.fieldNames = append(m.fieldNames, name)
	m.fields[name] = field
	return m
}

func (m *DefaultPlcReadRequestBuilder) Build() (model.PlcReadRequest, error) {
	for _, name := range m.queryNames {
		query := m.queries[name]
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing query: %s", query)
		}
		m.AddField(name, field)
	}
	return NewDefaultPlcReadRequest(m.fields, m.fieldNames, m.reader, m.readRequestInterceptor), nil
}

type DefaultPlcReadRequest struct {
	DefaultRequest
	reader                 spi.PlcReader
	readRequestInterceptor interceptors.ReadRequestInterceptor
}

func NewDefaultPlcReadRequest(fields map[string]model.PlcField, fieldNames []string, reader spi.PlcReader, readRequestInterceptor interceptors.ReadRequestInterceptor) model.PlcReadRequest {
	return DefaultPlcReadRequest{
		DefaultRequest:         NewDefaultRequest(fields, fieldNames),
		reader:                 reader,
		readRequestInterceptor: readRequestInterceptor,
	}
}

func (m DefaultPlcReadRequest) GetReader() spi.PlcReader {
	return m.reader
}

func (m DefaultPlcReadRequest) GetReadRequestInterceptor() interceptors.ReadRequestInterceptor {
	return m.readRequestInterceptor
}

func (m DefaultPlcReadRequest) Execute() <-chan model.PlcReadRequestResult {
	// Shortcut, if no interceptor is defined
	if m.readRequestInterceptor == nil {
		return m.reader.Read(m)
	}

	// Split the requests up into multiple ones.
	readRequests := m.readRequestInterceptor.InterceptReadRequest(m)
	// Shortcut for single-request-requests
	if len(readRequests) == 1 {
		return m.reader.Read(readRequests[0])
	}
	// Create a sub-result-channel slice
	var subResultChannels []<-chan model.PlcReadRequestResult

	// Iterate over all requests and add the result-channels to the list
	for _, subRequest := range readRequests {
		subResultChannels = append(subResultChannels, m.reader.Read(subRequest))
		// TODO: Replace this with a real queueing of requests. Later on we need throttling. At the moment this avoids race condition as the read above writes to fast on the line which is a problem for the test
		time.Sleep(time.Millisecond * 4)
	}

	// Create a new result-channel, which completes as soon as all sub-result-channels have returned
	resultChannel := make(chan model.PlcReadRequestResult)
	go func() {
		var subResults []model.PlcReadRequestResult
		// Iterate over all sub-results
		for _, subResultChannel := range subResultChannels {
			subResult := <-subResultChannel
			subResults = append(subResults, subResult)
		}
		// As soon as all are done, process the results
		result := m.readRequestInterceptor.ProcessReadResponses(m, subResults)
		// Return the final result
		resultChannel <- result
	}()

	return resultChannel
}

func (m DefaultPlcReadRequest) ExecuteWithContext(_ context.Context) <-chan model.PlcReadRequestResult {
	return m.Execute()
}

func (m DefaultPlcReadRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcReadRequest"); err != nil {
		return err
	}

	if err := writeBuffer.PushContext("fields"); err != nil {
		return err
	}
	for _, fieldName := range m.GetFieldNames() {
		if err := writeBuffer.PushContext(fieldName); err != nil {
			return err
		}
		field := m.GetField(fieldName)
		if serializableField, ok := field.(utils.Serializable); ok {
			if err := serializableField.Serialize(writeBuffer); err != nil {
				return err
			}
		} else {
			return errors.Errorf("Error serializing. Field %T doesn't implement Serializable", field)
		}
		if err := writeBuffer.PopContext(fieldName); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("fields"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcReadRequest"); err != nil {
		return err
	}
	return nil
}

func (m DefaultPlcReadRequest) String() string {
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
