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

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWrite "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type EIPPlcTag interface {
	model.PlcTag
	utils.Serializable

	GetTag() string
	GetType() readWrite.CIPDataTypeCode
	GetElementNb() uint16
}

type plcTag struct {
	Tag       string
	Type      readWrite.CIPDataTypeCode
	ElementNb uint16
}

func NewTag(tag string, _type readWrite.CIPDataTypeCode, elementNb uint16) plcTag {
	return plcTag{
		Tag:       tag,
		Type:      _type,
		ElementNb: elementNb,
	}
}

func (m plcTag) GetAddressString() string {
	return m.GetTag()
}

func (m plcTag) GetValueType() values.PlcValueType {
	if plcValueType, ok := values.PlcValueByName(m.GetType().String()); !ok {
		return values.NULL
	} else {
		return plcValueType
	}
}

func (m plcTag) GetArrayInfo() []model.ArrayInfo {
	return []model.ArrayInfo{}
}

func (m plcTag) GetTag() string {
	return m.Tag
}

func (m plcTag) GetType() readWrite.CIPDataTypeCode {
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

func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("EipTag"); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("node", uint32(len([]rune(m.Tag))*8), "UTF-8", m.Tag); err != nil {
		return err
	}

	if m.Type != 0 {
		if err := writeBuffer.WriteString("type", uint32(len([]rune(m.Type.String()))*8), "UTF-8", m.Type.String()); err != nil {
			return err
		}
	}

	if err := writeBuffer.WriteUint16("elementNb", 16, m.ElementNb); err != nil {
		return err
	}

	if err := writeBuffer.PopContext("EipTag"); err != nil {
		return err
	}
	return nil
}
