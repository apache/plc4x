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
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"time"
)

type DefaultPlcSubscriptionEvent struct {
	fields        map[string]model.PlcField
	types         map[string]SubscriptionType
	intervals     map[string]time.Duration
	responseCodes map[string]model.PlcResponseCode
	values        map[string]values.PlcValue
}

func NewDefaultPlcSubscriptionEvent(fields map[string]model.PlcField, types map[string]SubscriptionType,
	intervals map[string]time.Duration, responseCodes map[string]model.PlcResponseCode,
	values map[string]values.PlcValue) DefaultPlcSubscriptionEvent {
	return DefaultPlcSubscriptionEvent{
		fields:        fields,
		types:         types,
		intervals:     intervals,
		responseCodes: responseCodes,
		values:        values,
	}
}

func (m DefaultPlcSubscriptionEvent) GetFieldNames() []string {
	var fieldNames []string
	for fieldName := range m.fields {
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames
}

func (m DefaultPlcSubscriptionEvent) GetField(name string) model.PlcField {
	return m.fields[name]
}

func (m DefaultPlcSubscriptionEvent) GetType(name string) SubscriptionType {
	return m.types[name]
}

func (m DefaultPlcSubscriptionEvent) GetInterval(name string) time.Duration {
	return m.intervals[name]
}

func (m DefaultPlcSubscriptionEvent) GetResponseCode(name string) model.PlcResponseCode {
	return m.responseCodes[name]
}

func (m DefaultPlcSubscriptionEvent) GetAddress(name string) string {
	panic("GetAddress not implemented")
}

func (m DefaultPlcSubscriptionEvent) GetValue(name string) values.PlcValue {
	return m.values[name]
}

func (m DefaultPlcSubscriptionEvent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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
