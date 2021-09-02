/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package values

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"unicode/utf16"
)

type PlcWCHAR struct {
	value []rune
	PlcSimpleValueAdapter
}

func NewPlcWCHAR(value uint16) PlcWCHAR {
	return PlcWCHAR{
		value: utf16.Decode([]uint16{value}),
	}
}

func (m PlcWCHAR) IsString() bool {
	return true
}

func (m PlcWCHAR) GetString() string {
	return string(m.value)
}

func (m PlcWCHAR) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcSTRING", uint8(len(m.value)*8), "UTF-8", string(m.value))
}
