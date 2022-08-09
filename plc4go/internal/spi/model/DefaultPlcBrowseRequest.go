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
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/pkg/errors"
)

type DefaultPlcBrowseRequestBuilder struct {
	fieldHandler spi.PlcFieldHandler
	browser      spi.PlcBrowser
	queries      map[string]string
	queryNames   []string
	fields       map[string]model.PlcField
	fieldNames   []string
}

func NewDefaultPlcBrowseRequestBuilder(fieldHandler spi.PlcFieldHandler, browser spi.PlcBrowser) *DefaultPlcBrowseRequestBuilder {
	return &DefaultPlcBrowseRequestBuilder{
		fieldHandler: fieldHandler,
		browser:      browser,
		queries:      map[string]string{},
		fields:       map[string]model.PlcField{},
	}
}

func (d *DefaultPlcBrowseRequestBuilder) AddQuery(name string, query string) model.PlcBrowseRequestBuilder {
	d.queryNames = append(d.queryNames, name)
	d.queries[name] = query
	return d
}

func (d *DefaultPlcBrowseRequestBuilder) AddField(name string, field model.PlcField) model.PlcBrowseRequestBuilder {
	d.fieldNames = append(d.fieldNames, name)
	d.fields[name] = field
	return d
}

func (d *DefaultPlcBrowseRequestBuilder) Build() (model.PlcBrowseRequest, error) {
	for _, name := range d.queryNames {
		query := d.queries[name]
		field, err := d.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing query: %s", query)
		}
		d.AddField(name, field)
	}
	return NewDefaultPlcBrowseRequest(d.fields, d.fieldNames, d.browser), nil
}

type DefaultPlcBrowseRequest struct {
	DefaultRequest
	browser spi.PlcBrowser
}

func NewDefaultPlcBrowseRequest(fields map[string]model.PlcField, fieldNames []string, browser spi.PlcBrowser) model.PlcBrowseRequest {
	return DefaultPlcBrowseRequest{
		DefaultRequest: NewDefaultRequest(fields, fieldNames),
		browser:        browser,
	}
}

func (d DefaultPlcBrowseRequest) Execute() <-chan model.PlcBrowseRequestResult {
	return d.browser.Browse(d)
}

func (d DefaultPlcBrowseRequest) ExecuteWithInterceptor(interceptor func(result model.PlcBrowseEvent) bool) <-chan model.PlcBrowseRequestResult {
	return d.browser.BrowseWithInterceptor(d, interceptor)
}

func (d DefaultPlcBrowseRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseRequest"); err != nil {
		return err
	}

	if err := writeBuffer.PushContext("fields"); err != nil {
		return err
	}
	for _, fieldName := range d.GetFieldNames() {
		if err := writeBuffer.PushContext(fieldName); err != nil {
			return err
		}
		field := d.GetField(fieldName)
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
	if err := writeBuffer.PopContext("PlcBrowseRequest"); err != nil {
		return err
	}
	return nil
}

func (d DefaultPlcBrowseRequest) String() string {
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
