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

package simulated

import (
	"fmt"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type Field interface {
	apiModel.PlcField

	GetFieldType() *FieldType
	GetName() string
	GetDataTypeSize() *model.SimulatedDataTypeSizes
}

type simulatedField struct {
	FieldType    FieldType
	Name         string
	DataTypeSize model.SimulatedDataTypeSizes
	Quantity     uint16
}

func NewSimulatedField(fieldType FieldType, name string, dataTypeSize model.SimulatedDataTypeSizes, quantity uint16) apiModel.PlcField {
	return simulatedField{
		FieldType:    fieldType,
		Name:         name,
		DataTypeSize: dataTypeSize,
		Quantity:     quantity,
	}
}

func (t simulatedField) GetFieldType() FieldType {
	return t.FieldType
}

func (t simulatedField) GetName() string {
	return t.Name
}

func (t simulatedField) GetDataTypeSize() model.SimulatedDataTypeSizes {
	return t.DataTypeSize
}

func (t simulatedField) GetAddressString() string {
	return fmt.Sprintf("%s/%s:%s[%d]", t.FieldType.Name(), t.Name, t.DataTypeSize.String(), t.Quantity)
}

func (t simulatedField) GetValueType() values.PlcValueType {
	if plcValueType, ok := values.PlcValueByName(t.DataTypeSize.String()); ok {
		return plcValueType
	}
	return values.NULL
}

func (t simulatedField) GetArrayInfo() []apiModel.ArrayInfo {
	if t.Quantity != 1 {
		return []apiModel.ArrayInfo{
			spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(t.Quantity),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}
