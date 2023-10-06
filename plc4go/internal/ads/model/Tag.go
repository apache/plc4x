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

package model

import (
	"context"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

const NONE = int32(-1)

type PlcTag struct {
	apiModel.PlcSubscriptionTag

	ArrayInfo []apiModel.ArrayInfo
}

func NeedsResolving(plcTag apiModel.PlcTag) bool {
	switch plcTag.(type) {
	case SymbolicPlcTag:
		return true
	case DirectPlcTag:
		return false
	default:
		return false
	}
}

type DirectPlcTag struct {
	PlcTag

	IndexGroup   uint32
	IndexOffset  uint32
	ValueType    apiValues.PlcValueType
	StringLength int32
	DataType     readWriteModel.AdsDataTypeTableEntry
}

func NewDirectAdsPlcTag(indexGroup uint32, indexOffset uint32, valueType apiValues.PlcValueType, stringLength int32, arrayInfo []apiModel.ArrayInfo) (apiModel.PlcTag, error) {
	return DirectPlcTag{
		IndexGroup:   indexGroup,
		IndexOffset:  indexOffset,
		ValueType:    valueType,
		StringLength: stringLength,
		PlcTag: PlcTag{
			ArrayInfo: arrayInfo,
		},
	}, nil
}

func CastToDirectAdsTagFromPlcTag(plcTag apiModel.PlcTag) (DirectPlcTag, error) {
	if adsTag, ok := plcTag.(DirectPlcTag); ok {
		return adsTag, nil
	}
	return DirectPlcTag{}, errors.Errorf("couldn't %T cast to DirectPlcTag", plcTag)
}

func (m DirectPlcTag) GetAddressString() string {
	address := fmt.Sprintf("0x%d/%d:%s", m.IndexGroup, m.IndexOffset, m.ValueType.String())
	if m.ValueType == apiValues.STRING || m.ValueType == apiValues.WSTRING {
		address = address + "(" + strconv.Itoa(int(m.StringLength)) + ")"
	}
	if len(m.ArrayInfo) > 0 {
		for _, ai := range m.ArrayInfo {
			address = address + "[" + strconv.Itoa(int(ai.GetLowerBound())) + ".." + strconv.Itoa(int(ai.GetUpperBound())) + "]"
		}
	}
	return address
}

func (m DirectPlcTag) GetValueType() apiValues.PlcValueType {
	return m.ValueType
}

func (m DirectPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (m DirectPlcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m DirectPlcTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("DirectPlcTag"); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint32("indexGroup", 32, m.IndexGroup); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint32("indexOffset", 32, m.IndexOffset); err != nil {
		return err
	}
	if err := writeBuffer.WriteString("adsDatatypeName", uint32(len([]rune(m.ValueType.String()))*8), "UTF-8", m.ValueType.String()); err != nil {
		return err
	}
	if (m.ValueType == apiValues.STRING || m.ValueType == apiValues.WSTRING) && (m.StringLength != NONE) {
		if err := writeBuffer.WriteInt32("stringLength", 32, m.StringLength); err != nil {
			return err
		}
	}
	if len(m.ArrayInfo) > 0 {
		if err := writeBuffer.PushContext("ArrayInfo"); err != nil {
			return err
		}
		for _, ai := range m.ArrayInfo {
			if err := writeBuffer.PushContext("ArrayInfo"); err != nil {
				return err
			}
			if err := writeBuffer.WriteInt32("lowerBound", 32, int32(ai.GetLowerBound())); err != nil {
				return err
			}
			if err := writeBuffer.WriteInt32("upperBound", 32, int32(ai.GetUpperBound())); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("ArrayInfo"); err != nil {
				return err
			}
		}
		if err := writeBuffer.PopContext("ArrayInfo"); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext("DirectPlcTag"); err != nil {
		return err
	}
	return nil
}

func (m DirectPlcTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (m DirectPlcTag) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{}, errors.Errorf("%s", name) // TODO: why did this panic before
}

type SymbolicPlcTag struct {
	PlcTag

	SymbolicAddress string
}

func NewAdsSymbolicPlcTag(symbolicAddress string, arrayInfo []apiModel.ArrayInfo) (apiModel.PlcTag, error) {
	return SymbolicPlcTag{
		SymbolicAddress: symbolicAddress,
		PlcTag: PlcTag{
			ArrayInfo: arrayInfo,
		},
	}, nil
}

func CastToSymbolicPlcTagFromPlcTag(plcTag apiModel.PlcTag) (SymbolicPlcTag, error) {
	if adsTag, ok := plcTag.(SymbolicPlcTag); ok {
		return adsTag, nil
	}
	return SymbolicPlcTag{}, errors.Errorf("couldn't cast %T to SymbolicPlcTag", plcTag)
}

func (m SymbolicPlcTag) GetAddressString() string {
	return m.SymbolicAddress
}

func (m SymbolicPlcTag) GetValueType() apiValues.PlcValueType {
	return apiValues.NULL
}

func (m SymbolicPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (m SymbolicPlcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m SymbolicPlcTag) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("SymbolicPlcTag"); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("symbolicAddress", uint32(len([]rune(m.SymbolicAddress))*8), "UTF-8", m.SymbolicAddress); err != nil {
		return err
	}
	if len(m.ArrayInfo) > 0 {
		if err := writeBuffer.PushContext("ArrayInfo"); err != nil {
			return err
		}
		for _, ai := range m.ArrayInfo {
			if err := writeBuffer.PushContext("ArrayInfo"); err != nil {
				return err
			}
			if err := writeBuffer.WriteInt32("lowerBound", 32, int32(ai.GetLowerBound())); err != nil {
				return err
			}
			if err := writeBuffer.WriteInt32("upperBound", 32, int32(ai.GetUpperBound())); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("ArrayInfo"); err != nil {
				return err
			}
		}
		if err := writeBuffer.PopContext("ArrayInfo"); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext("SymbolicPlcTag"); err != nil {
		return err
	}
	return nil
}
