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
)

type PlcCHAR struct {
	value []byte
	PlcSimpleValueAdapter
}

func NewPlcCHAR(value uint8) PlcCHAR {
	return PlcCHAR{
		value: []byte{value},
	}
}

func (m PlcCHAR) IsString() bool {
	return true
}

func (m PlcCHAR) GetString() string {
	return string(m.value)
}

func (m PlcCHAR) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcBYTE", 16, "UTF-8", string(m.value))
}
