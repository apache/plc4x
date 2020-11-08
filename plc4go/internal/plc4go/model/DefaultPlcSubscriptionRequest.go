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
	"errors"
	values2 "plc4x.apache.org/plc4go/v0/internal/plc4go/model/values"
	"plc4x.apache.org/plc4go/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go/v0/pkg/plc4go/model"
	"plc4x.apache.org/plc4go/v0/pkg/plc4go/values"
)

type DefaultPlcSubscriptionRequestBuilder struct {
	subscriber   spi.PlcSubscriber
	fieldHandler spi.PlcFieldHandler
	valueHandler spi.PlcValueHandler
	queries      map[string]string
	values       map[string]interface{}
}

func NewDefaultPlcSubscriptionRequestBuilder(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, subscriber spi.PlcSubscriber) *DefaultPlcSubscriptionRequestBuilder {
	return &DefaultPlcSubscriptionRequestBuilder{
        subscriber:   subscriber,
		fieldHandler: fieldHandler,
		valueHandler: valueHandler,
		queries:      map[string]string{},
		values:       map[string]interface{}{},
	}
}

func (m *DefaultPlcSubscriptionRequestBuilder) AddItem(name string, query string, value interface{}) {
	m.queries[name] = query
	m.values[name] = value
}

func (m *DefaultPlcSubscriptionRequestBuilder) Build() (model.PlcSubscriptionRequest, error) {
	fields := make(map[string]model.PlcField)
	values := make(map[string]values.PlcValue)
	for name, query := range m.queries {
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.New("Error parsing query: " + query + ". Got error: " + err.Error())
		}
		fields[name] = field
		value, err := m.valueHandler.NewPlcValue(field, m.values[name])
		if err != nil {
			return nil, errors.New("Error parsing value of type: " + field.GetTypeName() + ". Got error: " + err.Error())
		}
		values[name] = value
	}
	return DefaultPlcSubscriptionRequest{
		fields: fields,
		values: values,
        subscriber: m.subscriber,
	}, nil
}

type DefaultPlcSubscriptionRequest struct {
	fields map[string]model.PlcField
	values map[string]values.PlcValue
    subscriber spi.PlcSubscriber
	model.PlcSubscriptionRequest
}

func (m DefaultPlcSubscriptionRequest) Execute() <-chan model.PlcSubscriptionRequestResult {
    panic("implement me")
}

func (m DefaultPlcSubscriptionRequest) GetFieldNames() []string {
	var fieldNames []string
	for fieldName, _ := range m.fields {
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames
}

func (m DefaultPlcSubscriptionRequest) GetField(name string) model.PlcField {
	return m.fields[name]
}

func (m DefaultPlcSubscriptionRequest) GetValue(name string) values.PlcValue {
	return m.values[name]
}

func (m DefaultPlcSubscriptionRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcSubscriptionRequest"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}
	for fieldName, field := range m.fields {
		value := m.values[fieldName]
		if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
		if err := e.EncodeElement(field, xml.StartElement{Name: xml.Name{Local: "field"}}); err != nil {
			return err
		}
		switch value.(type) {
		case values2.PlcList:
			listValue, ok := value.(values2.PlcList)
			if !ok {
				return errors.New("couldn't cast PlcValue to PlcList")
			}
			for _, subValue := range listValue.Values {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
					return err
				}
				if !subValue.IsString() {
					return errors.New("value not serializable to string")
				}
				e.EncodeToken(xml.CharData(subValue.GetString()))
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "value"}}); err != nil {
					return err
				}
			}
		default:
			if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
				return err
			}
			if !value.IsString() {
				return errors.New("value not serializable to string")
			}
			e.EncodeToken(xml.CharData(value.GetString()))
			if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "value"}}); err != nil {
				return err
			}
		}
		if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcSubscriptionRequest"}}); err != nil {
		return err
	}
	return nil
}
