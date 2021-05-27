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

import "encoding/xml"

type PlcLWORD struct {
	value uint64
	PlcSimpleValueAdapter
}

func NewPlcLWORD(value uint64) PlcLWORD {
	return PlcLWORD{
		value: value,
	}
}

func (m PlcLWORD) IsBool() bool {
	return true
}

func (m PlcLWORD) GetBoolLength() uint32 {
	return 64
}

func (m PlcLWORD) GetBool() bool {
	return m.value&1 == 1
}

func (m PlcLWORD) GetBoolAt(index uint32) bool {
	if index > 63 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcLWORD) GetBoolArray() []bool {
	return []bool{m.value&1 == 1, m.value>>1&1 == 1,
		m.value>>2&1 == 1, m.value>>3&1 == 1,
		m.value>>4&1 == 1, m.value>>5&1 == 1,
		m.value>>6&1 == 1, m.value>>7&1 == 1,
		m.value>>8&1 == 1, m.value>>9&1 == 1,
		m.value>>10&1 == 1, m.value>>11&1 == 1,
		m.value>>12&1 == 1, m.value>>13&1 == 1,
		m.value>>14&1 == 1, m.value>>15&1 == 1,
		m.value>>16&1 == 1, m.value>>17&1 == 1,
		m.value>>18&1 == 1, m.value>>19&1 == 1,
		m.value>>20&1 == 1, m.value>>21&1 == 1,
		m.value>>22&1 == 1, m.value>>23&1 == 1,
		m.value>>24&1 == 1, m.value>>25&1 == 1,
		m.value>>26&1 == 1, m.value>>27&1 == 1,
		m.value>>28&1 == 1, m.value>>29&1 == 1,
		m.value>>30&1 == 1, m.value>>31&1 == 1,
		m.value>>32&1 == 1, m.value>>33&1 == 1,
		m.value>>34&1 == 1, m.value>>35&1 == 1,
		m.value>>36&1 == 1, m.value>>37&1 == 1,
		m.value>>38&1 == 1, m.value>>39&1 == 1,
		m.value>>40&1 == 1, m.value>>41&1 == 1,
		m.value>>42&1 == 1, m.value>>43&1 == 1,
		m.value>>44&1 == 1, m.value>>45&1 == 1,
		m.value>>46&1 == 1, m.value>>47&1 == 1,
		m.value>>48&1 == 1, m.value>>49&1 == 1,
		m.value>>50&1 == 1, m.value>>51&1 == 1,
		m.value>>52&1 == 1, m.value>>53&1 == 1,
		m.value>>54&1 == 1, m.value>>55&1 == 1,
		m.value>>56&1 == 1, m.value>>57&1 == 1,
		m.value>>58&1 == 1, m.value>>59&1 == 1,
		m.value>>60&1 == 1, m.value>>61&1 == 1,
		m.value>>62&1 == 1, m.value>>63&1 == 1}
}

func (m PlcLWORD) IsString() bool {
	return true
}

func (m PlcLWORD) GetString() string {
	var strVal string
	for i, val := range m.GetBoolArray() {
		if i > 0 {
			strVal = strVal + ", "
		}
		if val {
			strVal = strVal + "true"
		} else {
			strVal = strVal + "false"
		}
	}
	return strVal
}

func (m PlcLWORD) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.value, xml.StartElement{Name: xml.Name{Local: "PlcLWORD"}}); err != nil {
		return err
	}
	return nil
}
