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
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	"github.com/pkg/errors"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcWriteRequestBuilder
type DefaultPlcWriteRequestBuilder struct {
	writer                  spi.PlcWriter       `ignore:"true"`
	tagHandler              spi.PlcTagHandler   `ignore:"true"`
	valueHandler            spi.PlcValueHandler `ignore:"true"`
	tagNames                []string
	tagAddresses            map[string]string
	tags                    map[string]model.PlcTag
	values                  map[string]any
	writeRequestInterceptor interceptors.WriteRequestInterceptor `ignore:"true"`
}

func NewDefaultPlcWriteRequestBuilder(tagHandler spi.PlcTagHandler, valueHandler spi.PlcValueHandler, writer spi.PlcWriter) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		writer:       writer,
		tagHandler:   tagHandler,
		valueHandler: valueHandler,
		tagNames:     make([]string, 0),
		tagAddresses: map[string]string{},
		tags:         map[string]model.PlcTag{},
		values:       map[string]any{},
	}
}

func NewDefaultPlcWriteRequestBuilderWithInterceptor(tagHandler spi.PlcTagHandler, valueHandler spi.PlcValueHandler, writer spi.PlcWriter, writeRequestInterceptor interceptors.WriteRequestInterceptor) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		writer:                  writer,
		tagHandler:              tagHandler,
		valueHandler:            valueHandler,
		tagNames:                make([]string, 0),
		tagAddresses:            map[string]string{},
		tags:                    map[string]model.PlcTag{},
		values:                  map[string]any{},
		writeRequestInterceptor: writeRequestInterceptor,
	}
}

func (m *DefaultPlcWriteRequestBuilder) GetWriter() spi.PlcWriter {
	return m.writer
}

func (m *DefaultPlcWriteRequestBuilder) GetWriteRequestInterceptor() interceptors.WriteRequestInterceptor {
	return m.writeRequestInterceptor
}

func (m *DefaultPlcWriteRequestBuilder) AddTagAddress(name string, tagAddress string, value any) model.PlcWriteRequestBuilder {
	m.tagNames = append(m.tagNames, name)
	m.tagAddresses[name] = tagAddress
	m.values[name] = value
	return m
}

func (m *DefaultPlcWriteRequestBuilder) AddTag(name string, tag model.PlcTag, value any) model.PlcWriteRequestBuilder {
	m.tagNames = append(m.tagNames, name)
	m.tags[name] = tag
	m.values[name] = value
	return m
}

func (m *DefaultPlcWriteRequestBuilder) Build() (model.PlcWriteRequest, error) {
	// Parse any unparsed tagAddresses
	for _, name := range m.tagNames {
		if tagAddress, ok := m.tagAddresses[name]; ok {
			tag, err := m.tagHandler.ParseTag(tagAddress)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing tag query: %s", tagAddress)
			}
			m.tags[name] = tag
		}
	}
	// Reset the queries
	m.tagAddresses = map[string]string{}

	// Process the values for tags.
	plcValues := make(map[string]values.PlcValue)
	for name, tag := range m.tags {
		value, err := m.valueHandler.NewPlcValue(tag, m.values[name])
		if err != nil {
			//			return nil, errors.Wrapf(err, "Error parsing value of type: %s", tag.GetTypeName())
		}
		plcValues[name] = value
	}
	return NewDefaultPlcWriteRequest(m.tags, m.tagNames, plcValues, m.writer, m.writeRequestInterceptor), nil
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcWriteRequest
type DefaultPlcWriteRequest struct {
	DefaultPlcTagRequest
	values                  map[string]values.PlcValue
	writer                  spi.PlcWriter                        `ignore:"true"`
	writeRequestInterceptor interceptors.WriteRequestInterceptor `ignore:"true"`
}

func NewDefaultPlcWriteRequest(tags map[string]model.PlcTag, tagNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor interceptors.WriteRequestInterceptor) model.PlcWriteRequest {
	return &DefaultPlcWriteRequest{NewDefaultPlcTagRequest(tags, tagNames), values, writer, writeRequestInterceptor}
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
