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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go/model"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go/values"
)

type DefaultPlcReadRequestBuilder struct {
	names        []string
	queries      []string
	fieldHandler plc4go.PlcFieldHandler
	reader       PlcReader
}

func NewDefaultPlcReadRequestBuilder(fieldHandler plc4go.PlcFieldHandler, reader PlcReader) *DefaultPlcReadRequestBuilder {
	return &DefaultPlcReadRequestBuilder{
		names:        []string{},
		queries:      []string{},
		fieldHandler: fieldHandler,
		reader:       reader,
	}
}

func (m *DefaultPlcReadRequestBuilder) AddItem(name string, query string) *DefaultPlcReadRequestBuilder {
	m.names = append(m.names, name)
	m.queries = append(m.queries, query)
	return m
}

func (m *DefaultPlcReadRequestBuilder) Build() model.PlcReadRequest {
	fields := make(map[string]values.PlcValue)
	for i, name := range m.names {
		query := m.queries[i]
		field := m.fieldHandler.ParseQuery(query)
		fields[name] = field
	}
	return defaultPlcReadRequest{
		fields: fields,
		reader: m.reader,
	}
}

type defaultPlcReadRequest struct {
	fields map[string]values.PlcValue
	reader PlcReader
}

func (m defaultPlcReadRequest) Execute() (<-chan model.PlcReadResponse, error) {
	return m.reader.Read()
}
