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

type Tag interface {
	apiModel.PlcTag

	GetTagType() TagType
	GetName() string
	GetDataTypeSize() model.SimulatedDataTypeSizes
}

type simulatedTag struct {
	TagType      TagType
	Name         string
	DataTypeSize model.SimulatedDataTypeSizes
	Quantity     uint16
	// ExplicitRange records whether the address wrote the selection as a range. A one-element
	// range is still a range - [4] is a scalar and [4..4] a list of one - which no count can say.
	ExplicitRange bool
}

func NewSimulatedTag(tagType TagType, name string, dataTypeSize model.SimulatedDataTypeSizes, quantity uint16) Tag {
	return NewSimulatedTagWithShape(tagType, name, dataTypeSize, quantity, quantity > 1)
}

// NewSimulatedTagWithShape is NewSimulatedTag plus what the address said about its shape: a range
// is an array even when it spans one element, which the quantity alone cannot carry.
func NewSimulatedTagWithShape(tagType TagType, name string, dataTypeSize model.SimulatedDataTypeSizes, quantity uint16, explicitRange bool) Tag {
	return simulatedTag{
		ExplicitRange: explicitRange,
		TagType:       tagType,
		Name:          name,
		DataTypeSize:  dataTypeSize,
		Quantity:      quantity,
	}
}

func (t simulatedTag) GetTagType() TagType {
	return t.TagType
}

func (t simulatedTag) GetName() string {
	return t.Name
}

func (t simulatedTag) GetDataTypeSize() model.SimulatedDataTypeSizes {
	return t.DataTypeSize
}

func (t simulatedTag) GetAddressString() string {
	return fmt.Sprintf("%s/%s%s:%s", t.TagType.Name(), t.Name,
		spiModel.RenderArrayExpression(t.GetArrayInfo()), t.DataTypeSize.String())
}

func (t simulatedTag) GetValueType() values.PlcValueType {
	if plcValueType, ok := values.PlcValueTypeByName(t.DataTypeSize.String()); ok {
		return plcValueType
	}
	return values.NULL
}

func (t simulatedTag) GetArrayInfo() []apiModel.ArrayInfo {
	// The flag decides the shape; the count only sizes it.
	if t.ExplicitRange {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(t.Quantity) - 1,
				Range:      true,
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (t simulatedTag) String() string {
	return "simulated"
}
