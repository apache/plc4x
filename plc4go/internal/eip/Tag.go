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

package eip

import (
	"context"
	"encoding/binary"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcTag interface {
	apiModel.PlcTag
	utils.Serializable

	GetTag() string
	GetType() readWriteModel.CIPDataTypeCode
	GetElementNb() uint16
}

type plcTag struct {
	Tag       string
	Type      readWriteModel.CIPDataTypeCode
	ElementNb uint16
}

func NewTag(tag string, _type readWriteModel.CIPDataTypeCode, elementNb uint16) PlcTag {
	return plcTag{
		Tag:       tag,
		Type:      _type,
		ElementNb: elementNb,
	}
}

func (m plcTag) GetAddressString() string {
	return m.GetTag()
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	if plcValueType, ok := apiValues.PlcValueTypeByName(m.GetType().String()); !ok {
		return apiValues.NULL
	} else {
		return plcValueType
	}
}

func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (m plcTag) GetTag() string {
	return m.Tag
}

func (m plcTag) GetType() readWriteModel.CIPDataTypeCode {
	return m.Type
}

func (m plcTag) GetElementNb() uint16 {
	return m.ElementNb
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, wb utils.WriteBuffer) error {
	if err := wb.PushContext("EipTag"); err != nil {
		return err
	}

	if err := wb.WriteString("node", uint32(len([]rune(m.Tag))*8), m.Tag); err != nil {
		return err
	}

	if m.Type != 0 {
		if err := wb.WriteString("type", uint32(len([]rune(m.Type.String()))*8), m.Type.String()); err != nil {
			return err
		}
	}

	if err := wb.WriteUint16("elementNb", 16, m.ElementNb); err != nil {
		return err
	}

	if err := wb.PopContext("EipTag"); err != nil {
		return err
	}
	return nil
}

func (m plcTag) String() string {
	wb := utils.NewWriteBufferBoxBased(utils.WithWriteBufferBoxBasedOmitEmptyBoxes(), utils.WithWriteBufferBoxBasedMergeSingleBoxes())
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
