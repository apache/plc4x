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
	api "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcByteArray struct {
	Values []byte
	PlcValueAdapter
}

func NewPlcByteArray(values []byte) PlcByteArray {
	return PlcByteArray{
		Values: values,
	}
}

func (m PlcByteArray) GetRaw() []byte {
	return m.Values
}

func (m PlcByteArray) IsList() bool {
	return true
}

func (m PlcByteArray) GetLength() uint32 {
	return uint32(len(m.Values))
}

func (m PlcByteArray) GetIndex(i uint32) api.PlcValue {
	return NewPlcUSINT(m.Values[i])
}

func (m PlcByteArray) GetList() []api.PlcValue {
	var plcValues []api.PlcValue
	for _, value := range m.Values {
		plcValues = append(plcValues, NewPlcUSINT(value))
	}
	return plcValues
}

func (m PlcByteArray) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcByteArray"); err != nil {
		return err
	}
	for _, value := range m.Values {
		if err := writeBuffer.WriteByte("value", value); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("PlcByteArray"); err != nil {
		return err
	}
	return nil
}
