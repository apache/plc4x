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

type DefaultPlcWriteRequestBuilder struct {
	names        []string
	queries      []string
	plcValues    []values.PlcValue
	fieldHandler plc4go.PlcFieldHandler
	writer       PlcWriter
}

func NewDefaultPlcWriteRequestBuilder(fieldHandler plc4go.PlcFieldHandler, writer PlcWriter) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		names:        []string{},
		queries:      []string{},
		plcValues:    []values.PlcValue{},
		fieldHandler: fieldHandler,
		writer:       writer,
	}
}

func (m *DefaultPlcWriteRequestBuilder) AddItem(name string, query string, plcValues values.PlcValue) *DefaultPlcWriteRequestBuilder {
	m.names = append(m.names, name)
	m.queries = append(m.queries, query)
	m.plcValues = append(m.plcValues, plcValues)
	return m
}

func (m *DefaultPlcWriteRequestBuilder) Build() model.PlcWriteRequest {
	fields := make(map[string]model.PlcField)
	plcValues := make(map[string]values.PlcValue)
	for i, name := range m.names {
		query := m.queries[i]
		field := m.fieldHandler.ParseQuery(query)
		fields[name] = field
		plcValues[name] = m.plcValues[i]
	}
	return defaultPlcWriteRequest{
		fields:    fields,
		plcValues: plcValues,
		writer:    m.writer,
	}
}

type defaultPlcWriteRequest struct {
	fields    map[string]model.PlcField
	plcValues map[string]values.PlcValue
	writer    PlcWriter
}

func (m defaultPlcWriteRequest) Execute() (<-chan model.PlcWriteResponse, error) {
	return m.writer.Write()
}
