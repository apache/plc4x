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
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	"github.com/pkg/errors"
	"time"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcWriteRequestBuilder
type DefaultPlcWriteRequestBuilder struct {
	writer                  spi.PlcWriter
	fieldHandler            spi.PlcFieldHandler
	valueHandler            spi.PlcValueHandler
	queries                 map[string]string
	queryNames              []string
	fields                  map[string]model.PlcField
	fieldNames              []string
	values                  map[string]interface{}
	writeRequestInterceptor interceptors.WriteRequestInterceptor
}

func NewDefaultPlcWriteRequestBuilder(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, writer spi.PlcWriter) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		writer:       writer,
		fieldHandler: fieldHandler,
		valueHandler: valueHandler,
		queries:      map[string]string{},
		queryNames:   make([]string, 0),
		fields:       map[string]model.PlcField{},
		fieldNames:   make([]string, 0),
		values:       map[string]interface{}{},
	}
}

func NewDefaultPlcWriteRequestBuilderWithInterceptor(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, writer spi.PlcWriter, writeRequestInterceptor interceptors.WriteRequestInterceptor) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		writer:                  writer,
		fieldHandler:            fieldHandler,
		valueHandler:            valueHandler,
		queries:                 map[string]string{},
		queryNames:              make([]string, 0),
		fields:                  map[string]model.PlcField{},
		fieldNames:              make([]string, 0),
		values:                  map[string]interface{}{},
		writeRequestInterceptor: writeRequestInterceptor,
	}
}

func (m *DefaultPlcWriteRequestBuilder) GetWriter() spi.PlcWriter {
	return m.writer
}

func (m *DefaultPlcWriteRequestBuilder) GetWriteRequestInterceptor() interceptors.WriteRequestInterceptor {
	return m.writeRequestInterceptor
}

func (m *DefaultPlcWriteRequestBuilder) AddQuery(name string, query string, value interface{}) model.PlcWriteRequestBuilder {
	m.queryNames = append(m.queryNames, name)
	m.queries[name] = query
	m.values[name] = value
	return m
}

func (m *DefaultPlcWriteRequestBuilder) AddField(name string, field model.PlcField, value interface{}) model.PlcWriteRequestBuilder {
	m.fieldNames = append(m.fieldNames, name)
	m.fields[name] = field
	m.values[name] = value
	return m
}

func (m *DefaultPlcWriteRequestBuilder) Build() (model.PlcWriteRequest, error) {
	// Parse the queries as well as pro
	for _, name := range m.queryNames {
		query := m.queries[name]
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing query: %s", query)
		}
		m.AddField(name, field, m.values[name])
	}

	// Process the values for fields.
	plcValues := make(map[string]values.PlcValue)
	for name, field := range m.fields {
		value, err := m.valueHandler.NewPlcValue(field, m.values[name])
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing value of type: %s", field.GetTypeName())
		}
		plcValues[name] = value
	}
	return NewDefaultPlcWriteRequest(m.fields, m.fieldNames, plcValues, m.writer, m.writeRequestInterceptor), nil
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcWriteRequest
type DefaultPlcWriteRequest struct {
	DefaultRequest
	values                  map[string]values.PlcValue
	writer                  spi.PlcWriter
	writeRequestInterceptor interceptors.WriteRequestInterceptor
}

func NewDefaultPlcWriteRequest(fields map[string]model.PlcField, fieldNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor interceptors.WriteRequestInterceptor) model.PlcWriteRequest {
	return &DefaultPlcWriteRequest{NewDefaultRequest(fields, fieldNames), values, writer, writeRequestInterceptor}
}

func (d *DefaultPlcWriteRequest) Execute() <-chan model.PlcWriteRequestResult {
	return d.ExecuteWithContext(context.TODO())
}

func (d *DefaultPlcWriteRequest) ExecuteWithContext(ctx context.Context) <-chan model.PlcWriteRequestResult {
	// Shortcut, if no interceptor is defined
	if d.writeRequestInterceptor == nil {
		return d.writer.Write(ctx, d)
	}

	// Split the requests up into multiple ones.
	writeRequests := d.writeRequestInterceptor.InterceptWriteRequest(ctx, d)
	// Shortcut for single-request-requests
	if len(writeRequests) == 1 {
		return d.writer.Write(ctx, writeRequests[0])
	}
	// Create a sub-result-channel slice
	var subResultChannels []<-chan model.PlcWriteRequestResult

	// Iterate over all requests and add the result-channels to the list
	for _, subRequest := range writeRequests {
		subResultChannels = append(subResultChannels, d.writer.Write(ctx, subRequest))
		// TODO: Replace this with a real queueing of requests. Later on we need throttling. At the moment this avoids race condition as the read above writes to fast on the line which is a problem for the test
		time.Sleep(time.Millisecond * 4)
	}

	// Create a new result-channel, which completes as soon as all sub-result-channels have returned
	resultChannel := make(chan model.PlcWriteRequestResult)
	go func() {
		var subResults []model.PlcWriteRequestResult
		// Iterate over all sub-results
		for _, subResultChannel := range subResultChannels {
			select {
			case <-ctx.Done():
				resultChannel <- &DefaultPlcWriteRequestResult{Request: d, Err: ctx.Err()}
				return
			case subResult := <-subResultChannel:
				subResults = append(subResults, subResult)
			}
		}
		// As soon as all are done, process the results
		result := d.writeRequestInterceptor.ProcessWriteResponses(ctx, d, subResults)
		// Return the final result
		resultChannel <- result
	}()

	return resultChannel
}

func (d *DefaultPlcWriteRequest) GetWriter() spi.PlcWriter {
	return d.writer
}

func (d *DefaultPlcWriteRequest) GetWriteRequestInterceptor() interceptors.WriteRequestInterceptor {
	return d.writeRequestInterceptor
}

func (d *DefaultPlcWriteRequest) GetValue(name string) values.PlcValue {
	return d.values[name]
}
