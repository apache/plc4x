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

package model

import (
	"encoding/xml"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
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
	readRequestInterceptor ReadRequestInterceptor
}

func NewDefaultPlcReadRequestBuilder(fieldHandler spi.PlcFieldHandler, reader spi.PlcReader) *DefaultPlcReadRequestBuilder {
	return NewDefaultPlcReadRequestBuilderWithInterceptor(fieldHandler, reader, nil)
}

func NewDefaultPlcReadRequestBuilderWithInterceptor(fieldHandler spi.PlcFieldHandler, reader spi.PlcReader, readRequestInterceptor ReadRequestInterceptor) *DefaultPlcReadRequestBuilder {
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

func (m *DefaultPlcReadRequestBuilder) AddQuery(name string, query string) {
	m.queryNames = append(m.queryNames, name)
	m.queries[name] = query
}

func (m *DefaultPlcReadRequestBuilder) AddField(name string, field model.PlcField) {
	m.fieldNames = append(m.fieldNames, name)
	m.fields[name] = field
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
	return DefaultPlcReadRequest{
		Fields:                 m.fields,
		FieldNames:             m.fieldNames,
		Reader:                 m.reader,
		ReadRequestInterceptor: m.readRequestInterceptor,
	}, nil
}

type DefaultPlcReadRequest struct {
	Fields                 map[string]model.PlcField
	FieldNames             []string
	Reader                 spi.PlcReader
	ReadRequestInterceptor ReadRequestInterceptor
}

func NewDefaultPlcReadRequest(fields map[string]model.PlcField, fieldNames []string, reader spi.PlcReader, readRequestInterceptor ReadRequestInterceptor) DefaultPlcReadRequest {
	return DefaultPlcReadRequest{
		Fields:                 fields,
		FieldNames:             fieldNames,
		Reader:                 reader,
		ReadRequestInterceptor: readRequestInterceptor,
	}
}

func (m DefaultPlcReadRequest) Execute() <-chan model.PlcReadRequestResult {
	// Shortcut, if no interceptor is defined
	if m.ReadRequestInterceptor == nil {
		return m.Reader.Read(m)
	}

	// Split the requests up into multiple ones.
	readRequests := m.ReadRequestInterceptor.InterceptReadRequest(m)
	// Shortcut for single-request-requests
	if len(readRequests) == 1 {
		return m.Reader.Read(readRequests[0])
	}
	// Create a sub-result-channel slice
	var subResultChannels []<-chan model.PlcReadRequestResult

	// Iterate over all requests and add the result-channels to the list
	for _, subRequest := range readRequests {
		subResultChannels = append(subResultChannels, m.Reader.Read(subRequest))
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
		result := m.ReadRequestInterceptor.ProcessReadResponses(m, subResults)
		// Return the final result
		resultChannel <- result
	}()

	return resultChannel
}

func (m DefaultPlcReadRequest) GetFieldNames() []string {
	return m.FieldNames
}

func (m DefaultPlcReadRequest) GetField(name string) model.PlcField {
	if field, ok := m.Fields[name]; ok {
		return field
	}
	return nil
}

func (m DefaultPlcReadRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcReadRequest"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}
	for _, fieldName := range m.FieldNames {
		field := m.Fields[fieldName]
		if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
		if err := e.EncodeElement(field, xml.StartElement{Name: xml.Name{Local: "field"}}); err != nil {
			return err
		}
		if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcReadRequest"}}); err != nil {
		return err
	}
	return nil
}
