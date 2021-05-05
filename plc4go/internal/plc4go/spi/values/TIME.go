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
	"fmt"
	"time"
)

type PlcTIME struct {
	value uint32
	PlcSimpleValueAdapter
}

func NewPlcTIME(value uint32) PlcTIME {
	return PlcTIME{
		value: value,
	}
}

func (m PlcTIME) IsDuration() bool {
	return true
}

func (m PlcTIME) GetDuration() time.Duration {
	return time.Duration(m.value)
}

func (m PlcTIME) IsString() bool {
	return true
}

func (m PlcTIME) GetString() string {
	return fmt.Sprintf("PT%0.fS", m.GetDuration().Seconds())
}

func (m PlcTIME) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.value, xml.StartElement{Name: xml.Name{Local: "PlcTIME"}}); err != nil {
		return err
	}
	return nil
}
