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

package values

import (
	"encoding/xml"
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

type PlcList struct {
	Values []api.PlcValue
	PlcValueAdapter
}

func NewPlcList(values []api.PlcValue) api.PlcValue {
	return PlcList{
		Values: values,
	}
}

func (m PlcList) IsList() bool {
	return true
}

func (m PlcList) GetLength() uint32 {
	return uint32(len(m.Values))
}

func (m PlcList) GetIndex(i uint32) api.PlcValue {
	return m.Values[i]
}

func (m PlcList) GetList() []api.PlcValue {
	return m.Values
}

func (m PlcList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcList"}}); err != nil {
		return err
	}

	for _, value := range m.Values {
		if err := e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: "-set-by-element-"}}); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcList"}}); err != nil {
		return err
	}
	return nil
}
