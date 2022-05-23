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

package simulated

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
)

type Field interface {
	GetFieldType() *FieldType
	GetName() string
	GetDataTypeSize() *model.SimulatedDataTypeSizes
}

type SimulatedField struct {
	FieldType    FieldType
	Name         string
	DataTypeSize model.SimulatedDataTypeSizes
	Quantity     uint16
}

func NewSimulatedField(fieldType FieldType, name string, dataTypeSize model.SimulatedDataTypeSizes, quantity uint16) SimulatedField {
	return SimulatedField{
		FieldType:    fieldType,
		Name:         name,
		DataTypeSize: dataTypeSize,
		Quantity:     quantity,
	}
}

func (t SimulatedField) GetFieldType() FieldType {
	return t.FieldType
}

func (t SimulatedField) GetName() string {
	return t.Name
}

func (t SimulatedField) GetDataTypeSize() model.SimulatedDataTypeSizes {
	return t.DataTypeSize
}

func (t SimulatedField) GetAddressString() string {
	return fmt.Sprintf("%s/%s:%s[%d]", t.FieldType.Name(), t.Name, t.DataTypeSize.String(), t.Quantity)
}

func (t SimulatedField) GetTypeName() string {
	return t.DataTypeSize.String()
}

func (t SimulatedField) GetQuantity() uint16 {
	return t.Quantity
}
