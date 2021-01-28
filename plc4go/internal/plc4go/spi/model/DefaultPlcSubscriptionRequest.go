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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"time"
)

type SubscriptionType uint8

const (
	SUBSCRIPTION_CYCLIC          SubscriptionType = 0x01
	SUBSCRIPTION_CHANGE_OF_STATE SubscriptionType = 0x02
	SUBSCRIPTION_EVENT           SubscriptionType = 0x03
)

type DefaultPlcSubscriptionRequestBuilder struct {
	subscriber   spi.PlcSubscriber
	fieldHandler spi.PlcFieldHandler
	valueHandler spi.PlcValueHandler
	eventHandler model.PlcSubscriptionEventHandler
	queries      map[string]string
	types        map[string]SubscriptionType
	intervals    map[string]time.Duration
}

func NewDefaultPlcSubscriptionRequestBuilder(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, subscriber spi.PlcSubscriber) *DefaultPlcSubscriptionRequestBuilder {
	return &DefaultPlcSubscriptionRequestBuilder{
		subscriber:   subscriber,
		fieldHandler: fieldHandler,
		valueHandler: valueHandler,
		queries:      map[string]string{},
		types:        map[string]SubscriptionType{},
		intervals:    map[string]time.Duration{},
	}
}

func (m *DefaultPlcSubscriptionRequestBuilder) AddCyclicItem(name string, query string, interval time.Duration) {
	m.queries[name] = query
	m.types[name] = SUBSCRIPTION_CYCLIC
	m.intervals[name] = interval
}

func (m *DefaultPlcSubscriptionRequestBuilder) AddChangeOfStateItem(name string, query string) {
	m.queries[name] = query
	m.types[name] = SUBSCRIPTION_CHANGE_OF_STATE
}

func (m *DefaultPlcSubscriptionRequestBuilder) AddEventItem(name string, query string) {
	m.queries[name] = query
	m.types[name] = SUBSCRIPTION_EVENT
}

func (m *DefaultPlcSubscriptionRequestBuilder) AddItemHandler(eventHandler model.PlcSubscriptionEventHandler) {
	m.eventHandler = eventHandler
}

func (m *DefaultPlcSubscriptionRequestBuilder) Build() (model.PlcSubscriptionRequest, error) {
	fields := make(map[string]model.PlcField)
	for name, query := range m.queries {
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.New("Error parsing query: " + query + ". Got error: " + err.Error())
		}
		fields[name] = field
	}
	return DefaultPlcSubscriptionRequest{
		fields:       fields,
		types:        m.types,
		intervals:    m.intervals,
		subscriber:   m.subscriber,
		eventHandler: m.eventHandler,
	}, nil
}

type DefaultPlcSubscriptionRequest struct {
	fields       map[string]model.PlcField
	types        map[string]SubscriptionType
	intervals    map[string]time.Duration
	eventHandler model.PlcSubscriptionEventHandler
	subscriber   spi.PlcSubscriber
	model.PlcSubscriptionRequest
}

func (m DefaultPlcSubscriptionRequest) Execute() <-chan model.PlcSubscriptionRequestResult {
	return m.subscriber.Subscribe(m)
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

func (m DefaultPlcSubscriptionRequest) GetEventHandler() model.PlcSubscriptionEventHandler {
	return m.eventHandler
}

func (m DefaultPlcSubscriptionRequest) GetType(name string) SubscriptionType {
	return m.types[name]
}

func (m DefaultPlcSubscriptionRequest) GetInterval(name string) time.Duration {
	return m.intervals[name]
}

func (m DefaultPlcSubscriptionRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcSubscriptionRequest"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}
	for fieldName, field := range m.fields {
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

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcSubscriptionRequest"}}); err != nil {
		return err
	}
	return nil
}
