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

type PlcLTIME struct {
	value uint64
	PlcSimpleValueAdapter
}

func NewPlcLTIME(value uint64) PlcLTIME {
	return PlcLTIME{
		value: value,
	}
}

func (m PlcLTIME) IsDuration() bool {
	return true
}

func (m PlcLTIME) GetDuration() time.Duration {
	return time.Duration(m.value)
}

func (m PlcLTIME) IsString() bool {
	return true
}

func (m PlcLTIME) GetString() string {
	return fmt.Sprintf("PT%0.fS", m.GetDuration().Seconds())
}

func (m PlcLTIME) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.value, xml.StartElement{Name: xml.Name{Local: "PlcLTIME"}}); err != nil {
		return err
	}
	return nil
}
