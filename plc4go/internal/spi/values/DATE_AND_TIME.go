/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"time"
)

type PlcDATE_AND_TIME struct {
	value time.Time
	PlcValueAdapter
}

func NewPlcDATE_AND_TIME(value time.Time) PlcDATE_AND_TIME {
	return PlcDATE_AND_TIME{
		value: value,
	}
}

func (m PlcDATE_AND_TIME) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	m.Serialize(buf)
	return buf.GetBytes()
}

func (m PlcDATE_AND_TIME) IsDateTime() bool {
	return true
}
func (m PlcDATE_AND_TIME) GetDateTime() time.Time {
	return time.Time{}.Add(m.GetDuration())
}

func (m PlcDATE_AND_TIME) GetString() string {
	return fmt.Sprintf("%v", m.GetDateTime())
}

func (m PlcDATE_AND_TIME) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcDATE_AND_TIME", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}
