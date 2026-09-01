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

package abeth

import (
	"context"
	"fmt"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// PlcTag is a parsed ab-eth tag address: one element of one data file, optionally narrowed down to
// a single bit. Ported from plc4j's AbEthTag.
//
// It is also an apiModel.PlcSubscriptionTag, because ab-eth subscriptions are emulated by polling
// the read path and the subscription-request builder only accepts tags which are one. plc4j gets
// this for free: its AbEthConnection extends PollingSubscriptionConnectionBase, which wraps the
// plain tag itself.
type PlcTag interface {
	apiModel.PlcSubscriptionTag
	utils.Serializable

	// GetByteSize is how many bytes the read asks the PLC for.
	GetByteSize() uint8
	// GetFileNumber is the number of the data file (the N<x> part of the address).
	GetFileNumber() uint8
	// GetFileType is what kind of data file it is and therefore how the answer is decoded.
	GetFileType() FileType
	// GetElementNumber is the element inside the file the read starts at.
	GetElementNumber() uint8
	// GetBitNumber is the bit a SINGLEBIT tag picks out of the element. Zero for every other tag.
	GetBitNumber() uint8
}

type plcTag struct {
	ByteSize      uint8
	FileNumber    uint8
	FileType      FileType
	ElementNumber uint8
	BitNumber     uint8
}

var _ PlcTag = plcTag{}

func NewTag(byteSize uint8, fileNumber uint8, fileType FileType, elementNumber uint8, bitNumber uint8) PlcTag {
	return plcTag{
		ByteSize:      byteSize,
		FileNumber:    fileNumber,
		FileType:      fileType,
		ElementNumber: elementNumber,
		BitNumber:     bitNumber,
	}
}

func (m plcTag) GetByteSize() uint8 {
	return m.ByteSize
}

func (m plcTag) GetFileNumber() uint8 {
	return m.FileNumber
}

func (m plcTag) GetFileType() FileType {
	return m.FileType
}

func (m plcTag) GetElementNumber() uint8 {
	return m.ElementNumber
}

func (m plcTag) GetBitNumber() uint8 {
	return m.BitNumber
}

// GetPlcSubscriptionType is what a tag which wasn't added through one of the typed builder methods
// defaults to. Polling can only emulate cyclic and change-of-state subscriptions, and reporting only
// what moved is the cheaper of the two.
func (m plcTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionChangeOfState
}

// GetDuration is not applicable: ab-eth has no per-tag subscription duration, the poll interval comes
// from the subscription request.
func (m plcTag) GetDuration() time.Duration {
	return 0
}

// GetAddressString spells the tag the way the tag handler parses it back.
//
// plc4j's AbEthTag.getAddressString drops the size suffix whenever the size is one, which makes the
// address unparseable for every real file type - those take their width from exactly that suffix.
// Since plc4go re-parses address strings (the polling subscriber rebuilds its read requests from
// them), the suffix is kept here and only omitted for the synthetic file types, which derive their
// width from the type name and only accept a suffix which agrees with it.
func (m plcTag) GetAddressString() string {
	address := fmt.Sprintf("N%d:%d", m.FileNumber, m.ElementNumber)
	if m.BitNumber != 0 {
		address += fmt.Sprintf("/%d", m.BitNumber)
	}
	address += ":" + m.FileType.String()
	if !m.FileType.hasFixedWidth() {
		address += fmt.Sprintf("[%d]", m.ByteSize)
	}
	return address
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	return m.FileType.GetPlcValueType()
}

// GetArrayInfo reports every ab-eth tag as a scalar, the way plc4j's AbEthTag does: the byte size
// is the size of the single value that comes back, not a number of elements.
func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

// SerializeWithWriteBuffer renders the tag wrapped as PlcTagItem/tag/AbEthTag, matching how plc4j
// wraps a tag inside a request (PlcTagItem) and what the driver testsuite golden files expect.
func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, wb utils.WriteBuffer) error {
	if err := wb.PushContext("PlcTagItem"); err != nil {
		return err
	}
	if err := wb.PushContext("tag"); err != nil {
		return err
	}
	if err := wb.PushContext("AbEthTag"); err != nil {
		return err
	}
	if err := wb.WriteUint8("byteSize", 8, m.ByteSize); err != nil {
		return err
	}
	if err := wb.WriteUint8("fileNumber", 8, m.FileNumber); err != nil {
		return err
	}
	fileTypeName := m.FileType.String()
	if err := wb.WriteString("fileType", uint32(len(fileTypeName)*8), fileTypeName, utils.WithEncoding("UTF-8")); err != nil {
		return err
	}
	if err := wb.WriteUint8("elementNumber", 8, m.ElementNumber); err != nil {
		return err
	}
	if err := wb.WriteUint8("bitNumber", 8, m.BitNumber); err != nil {
		return err
	}
	if err := wb.PopContext("AbEthTag"); err != nil {
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
