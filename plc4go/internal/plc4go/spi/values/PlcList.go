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
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/pkg/errors"
)

type PlcList struct {
	Values []values.PlcValue
	PlcValueAdapter
}

func NewPlcList(values []values.PlcValue) values.PlcValue {
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

func (m PlcList) GetIndex(i uint32) values.PlcValue {
	return m.Values[i]
}

func (m PlcList) GetList() []values.PlcValue {
	return m.Values
}

func (m PlcList) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcList"); err != nil {
		return err
	}
	for _, listItem := range m.GetList() {
		if listItemSerializable, ok := listItem.(utils.Serializable); ok {
			if err := listItemSerializable.Serialize(writeBuffer); err != nil {
				return err
			}
		} else {
			return errors.New("Error serializing. List item doesn't implement Serializable")
		}
	}
	if err := writeBuffer.PopContext("PlcList"); err != nil {
		return err
	}
	return nil
}
