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
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/model"
	"plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
)

type DefaultPlcWriteRequestBuilder struct {
	writer       spi.PlcWriter
	fieldHandler spi.PlcFieldHandler
	valueHandler spi.PlcValueHandler
	queries      map[string]string
	values       map[string]interface{}
}

func NewDefaultPlcWriteRequestBuilder(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, writer spi.PlcWriter) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		writer:       writer,
		fieldHandler: fieldHandler,
		valueHandler: valueHandler,
		queries:      map[string]string{},
		values:       map[string]interface{}{},
	}
}

func (m *DefaultPlcWriteRequestBuilder) AddItem(name string, query string, value interface{}) {
	m.queries[name] = query
	m.values[name] = value
}

func (m *DefaultPlcWriteRequestBuilder) Build() (model.PlcWriteRequest, error) {
	fields := make(map[string]model.PlcField)
	values := make(map[string]values.PlcValue)
	for _, name := range m.queries {
		query := m.queries[name]
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.New("Error parsing query: " + query + ". Got error: " + err.Error())
		}
		fields[name] = field
		value, err := m.valueHandler.NewPlcValue(field.GetTypeName(), m.values[name])
		if err != nil {
			return nil, errors.New("Error parsing value of type: " + field.GetTypeName() + ". Got error: " + err.Error())
		}
		values[name] = value
	}
	return DefaultPlcWriteRequest{
		fields: fields,
		values: values,
		writer: m.writer,
	}, nil
}

type DefaultPlcWriteRequest struct {
	fields map[string]model.PlcField
	values map[string]values.PlcValue
	writer spi.PlcWriter
	model.PlcWriteRequest
}

func (m DefaultPlcWriteRequest) Execute() <-chan model.PlcWriteRequestResult {
	return m.writer.Write(m)
}
