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

package cbus

import (
	"context"
	"encoding/binary"
	"fmt"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// UnitInfoQuery can be used to get information about unit(s)
type UnitInfoQuery interface {
	apiModel.PlcQuery

	GetUnitAddress() *readWriteModel.UnitAddress
	GetAttribute() *readWriteModel.Attribute
}

func NewUnitInfoQuery(unitAddress *readWriteModel.UnitAddress, attribute *readWriteModel.Attribute, numElements uint16) UnitInfoQuery {
	return &unitInfoQuery{
		unitAddress: unitAddress,
		tagType:     UNIT_INFO,
		attribute:   attribute,
		numElements: numElements,
	}
}

type unitInfoQuery struct {
	tagType     TagType
	unitAddress *readWriteModel.UnitAddress
	attribute   *readWriteModel.Attribute
	numElements uint16
}

func (u unitInfoQuery) GetQueryString() string {
	unitAddressString := "*"
	if u.unitAddress != nil {
		unitAddressString = fmt.Sprintf("%d", (*u.unitAddress).GetAddress())
	}
	attributeString := "*"
	if u.attribute != nil {
		unitAddressString = u.attribute.String()
	}
	return fmt.Sprintf("cal/%s/identify=%s", unitAddressString, attributeString)
}

func (u unitInfoQuery) GetValueType() values.PlcValueType {
	return values.Struct
}

func (u unitInfoQuery) GetArrayInfo() []apiModel.ArrayInfo {
	if u.numElements != 1 {
		return []apiModel.ArrayInfo{
			spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(u.numElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (u unitInfoQuery) GetTagType() TagType {
	return u.tagType
}

func (u unitInfoQuery) GetUnitAddress() *readWriteModel.UnitAddress {
	return u.unitAddress
}

func (u unitInfoQuery) GetAttribute() *readWriteModel.Attribute {
	return u.attribute
}

func (u unitInfoQuery) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := u.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (u unitInfoQuery) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(u.tagType.GetName()); err != nil {
		return err
	}

	if unitAddress := u.unitAddress; unitAddress != nil {
		if err := (*unitAddress).SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}

	if attribute := u.attribute; attribute != nil {
		if err := (*attribute).SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(u.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (u unitInfoQuery) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), u); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
