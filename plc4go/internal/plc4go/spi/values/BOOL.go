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
	"strconv"
)

type PlcBOOL struct {
	value bool
	PlcSimpleValueAdapter
}

func NewPlcBOOL(value bool) PlcBOOL {
	return PlcBOOL{
		value: value,
	}
}

func (m PlcBOOL) IsBool() bool {
	return true
}

func (m PlcBOOL) GetBoolLength() uint32 {
	return 1
}

func (m PlcBOOL) GetBool() bool {
	return m.value
}

func (m PlcBOOL) GetBoolAt(index uint32) bool {
	if index == 0 {
		return m.value
	}
	return false
}

func (m PlcBOOL) GetBoolArray() []bool {
	return []bool{m.value}
}

func (m PlcBOOL) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	size := 40
	if m.value {
		size = 32
	}
	if err := e.EncodeElement(m.value, xml.StartElement{Name: xml.Name{Local: "PlcBOOL"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "string"},
		{Name: xml.Name{Local: "bitLength"}, Value: strconv.Itoa(size)},
	}}); err != nil {
		return err
	}
	return nil
}

func (m PlcBOOL) IsString() bool {
	return true
}

func (m PlcBOOL) GetString() string {
	if m.GetBool() {
		return "true"
	} else {
		return "false"
	}
}
