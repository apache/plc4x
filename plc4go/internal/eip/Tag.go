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
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcTag interface {
	apiModel.PlcTag
	utils.Serializable

	GetTag() string
	GetType() readWriteModel.CIPDataTypeCode
	GetElementNb() uint16
	// GetSelection reports what the address selects, which drives the CIP path and the element
	// count. It is not what GetArrayInfo reports - see there.
	GetSelection() []apiModel.ArrayInfo
}

// maxElements is what a CIP request can ask for: the element count is carried in 16 bits.
const maxElements = 65535

type plcTag struct {
	Tag       string
	Type      readWriteModel.CIPDataTypeCode
	Selection []apiModel.ArrayInfo
}

// NewTag builds a tag selecting elementNb elements from the start of tag, which is the shape the
// element-count form described. An address selecting a range is built with NewTagWithSelection.
func NewTag(tag string, _type readWriteModel.CIPDataTypeCode, elementNb uint16) PlcTag {
	if elementNb < 1 {
		elementNb = 1
	}
	var selection []apiModel.ArrayInfo
	if elementNb > 1 {
		selection = []apiModel.ArrayInfo{&spiModel.DefaultArrayInfo{UpperBound: uint32(elementNb) - 1, Range: true}}
	}
	return NewTagWithSelection(tag, _type, selection)
}

func NewTagWithSelection(tag string, _type readWriteModel.CIPDataTypeCode, selection []apiModel.ArrayInfo) PlcTag {
	return plcTag{
		Tag:       tag,
		Type:      _type,
		Selection: selection,
	}
}

// GetAddressString mirrors what the address pattern accepts, so re-parsing it yields an equal
// tag: tag[selection][:dataType], with the selection before the type.
func (m plcTag) GetAddressString() string {
	address := "%" + m.Tag + spiModel.RenderArrayExpression(m.Selection)
	if m.Type != 0 {
		address = address + ":" + m.Type.String()
	}
	return address
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	if plcValueType, ok := apiValues.PlcValueTypeByName(m.GetType().String()); !ok {
		return apiValues.NULL
	} else {
		return plcValueType
	}
}

// GetArrayInfo reports the shape of the value the caller receives, so a consumer can tell a
// scalar from a list without knowing the protocol: empty for a scalar, one entry per dimension
// for an array. A bare index selects one element and so reports empty; a range reports its
// dimensions even when it spans a single element.
//
// This is not what the driver fetches - reading one element of an array still walks a member
// path, which the selection describes.
func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	for _, dimension := range m.Selection {
		if dimension.IsRange() {
			return m.Selection
		}
	}
	return []apiModel.ArrayInfo{}
}

func (m plcTag) GetSelection() []apiModel.ArrayInfo {
	return m.Selection
}

func (m plcTag) GetTag() string {
	return m.Tag
}

func (m plcTag) GetType() readWriteModel.CIPDataTypeCode {
	return m.Type
}

// GetElementNb is how many elements the request asks the device for, derived from the selection;
// a tag that selects nothing explicitly reads a single element.
func (m plcTag) GetElementNb() uint16 {
	// Computed as a uint64: the product of several dimensions wraps a uint32 long before it
	// reaches the wire, and the count is carried in 16 bits there. The handler refuses a selection
	// larger than that, so the conversion below cannot narrow a value anyone asked for.
	elements := uint64(1)
	for _, dimension := range m.Selection {
		elements *= uint64(dimension.GetSize())
	}
	if elements < 1 {
		return 1
	}
	if elements > maxElements {
		return maxElements
	}
	return uint16(elements)
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, wb utils.WriteBuffer) error {
	// The driver testsuite golden expects tags wrapped as PlcTagItem/tag/EipTag
	// (mirroring plc4j's PlcTagItem), rather than a bare EipTag.
	if err := wb.PushContext("PlcTagItem"); err != nil {
		return err
	}
	if err := wb.PushContext("tag"); err != nil {
		return err
	}
	if err := wb.PushContext("EipTag"); err != nil {
		return err
	}

	// encoding="UTF8" matches plc4j's EipTag.serialize (WithOption.WithEncoding("UTF8")
	// on both the "node" and "type" fields).
	if err := wb.WriteString("node", uint32(len([]rune(m.Tag))*8), m.Tag, utils.WithEncoding("UTF8")); err != nil {
		return err
	}

	if m.Type != 0 {
		if err := wb.WriteString("type", uint32(len([]rune(m.Type.String()))*8), m.Type.String(), utils.WithEncoding("UTF8")); err != nil {
			return err
		}
	}

	if err := wb.WriteUint16("elementNb", 16, m.GetElementNb()); err != nil {
		return err
	}

	if err := wb.PopContext("EipTag"); err != nil {
		return err
	}
	if err := wb.PopContext("tag"); err != nil {
		return err
	}
	if err := wb.PopContext("PlcTagItem"); err != nil {
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
