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

package slmp

import (
	"context"
	"fmt"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// maxPoints is the conservative single-frame word ceiling for 3E binary Batch Read/Write plc4j's
	// SlmpTag pins down. It is not the exact device maximum - there is no request optimizer to split
	// a bigger read, so a tag that would need one is refused instead of being sent as a frame the
	// device rejects.
	maxPoints = 960
	// maxDeviceNumber is the largest device address the 3E frame can carry: the head device number
	// is an unsigned 24-bit field.
	maxDeviceNumber = 0xFFFFFF
)

// PlcTag is a parsed slmp tag address: a run of consecutive word devices starting at one head
// device, read or written in word units. Ported from plc4j's SlmpTag.
//
// It is also an apiModel.PlcSubscriptionTag, because slmp subscriptions are emulated by polling the
// read path and the subscription-request builder only accepts tags which are one. plc4j gets this
// for free: its SlmpConnection extends PollingSubscriptionConnectionBase, which wraps the plain tag
// itself.
type PlcTag interface {
	apiModel.PlcSubscriptionTag
	utils.Serializable

	// GetDeviceCode is the 1-byte device code the frame addresses (D, W or R in this version).
	GetDeviceCode() readWriteModel.SlmpDeviceCode
	// GetDeviceNumber is the head device number the run starts at.
	GetDeviceNumber() uint32
	// GetDataType is how the words are interpreted in both directions.
	GetDataType() DataType
	// GetQuantity is how many elements of that type the tag covers.
	GetQuantity() uint16
	// GetNumberOfPoints is how many 16-bit words the Batch Read/Write transfers, which is the
	// quantity times the type's word footprint.
	GetNumberOfPoints() uint16
}

type plcTag struct {
	DeviceCode   readWriteModel.SlmpDeviceCode
	DeviceNumber uint32
	DataType     DataType
	Quantity     uint16
	// ExplicitRange records whether the address wrote the selection as a range. A one-element
	// range is still a range - [4] is a scalar and [4..4] a list of one - which no count can say.
	ExplicitRange bool
}

var _ PlcTag = plcTag{}

func NewTag(deviceCode readWriteModel.SlmpDeviceCode, deviceNumber uint32, dataType DataType, quantity uint16) PlcTag {
	return NewTagWithShape(deviceCode, deviceNumber, dataType, quantity, quantity > 1)
}

// NewTagWithShape is NewTag plus what the address said about its shape: a range is an array even
// when it spans one element, which the quantity alone cannot carry.
func NewTagWithShape(deviceCode readWriteModel.SlmpDeviceCode, deviceNumber uint32, dataType DataType, quantity uint16, explicitRange bool) PlcTag {
	return plcTag{
		ExplicitRange: explicitRange,
		DeviceCode:    deviceCode,
		DeviceNumber:  deviceNumber,
		DataType:      dataType,
		Quantity:      quantity,
	}
}

func (m plcTag) GetDeviceCode() readWriteModel.SlmpDeviceCode {
	return m.DeviceCode
}

func (m plcTag) GetDeviceNumber() uint32 {
	return m.DeviceNumber
}

func (m plcTag) GetDataType() DataType {
	return m.DataType
}

func (m plcTag) GetQuantity() uint16 {
	return m.Quantity
}

func (m plcTag) GetNumberOfPoints() uint16 {
	return m.Quantity * m.DataType.WordsPerElement()
}

// GetPlcSubscriptionType is what a tag which wasn't added through one of the typed builder methods
// defaults to. Polling can only emulate cyclic and change-of-state subscriptions, and reporting only
// what moved is the cheaper of the two.
func (m plcTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionChangeOfState
}

// GetDuration is not applicable: slmp has no per-tag subscription duration, the poll interval comes
// from the subscription request.
func (m plcTag) GetDuration() time.Duration {
	return 0
}

// GetAddressString spells the tag the way the tag handler parses it back. The round trip is load
// bearing: the polling subscriber rebuilds its read requests from this string.
//
// W addresses are written in hex with the 0x prefix, because that is how the handler reads a W
// address back; every other supported device is addressed in decimal.
func (m plcTag) GetAddressString() string {
	var address string
	if m.DeviceCode == readWriteModel.SlmpDeviceCode_W {
		address = fmt.Sprintf("0x%X", m.DeviceNumber)
	} else {
		address = fmt.Sprintf("%d", m.DeviceNumber)
	}
	// Unlike plc4j's SlmpTag.getAddressString, the data type is always spelled out. plc4j omits it
	// for a WORD tag of quantity one and spells it for every other tag, which round-trips too, but
	// an address that always names its type is the one a reader of a log line can act on.
	return fmt.Sprintf("%s%s%s:%s", m.DeviceCode, address,
		spiModel.RenderArrayExpression(m.GetArrayInfo()), m.DataType)
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	return m.DataType.GetValueType()
}

// GetArrayInfo reports the shape of the value the caller receives, as an inclusive range: a
// quantity of 5 yields [0..4], the same as plc4j.
//
// The indices are relative to what the caller receives, not to what was written: an SLMP address
// is a device number, so the driver folds the start of the selection into it when it resolves
// the address.
func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	// The flag decides the shape; the count only sizes it.
	if m.ExplicitRange {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(m.Quantity) - 1,
				Range:      true,
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

// SerializeWithWriteBuffer renders the tag wrapped as PlcTagItem/tag/SlmpTag, matching how plc4j
// wraps a tag inside a request (PlcTagItem).
func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, wb utils.WriteBuffer) error {
	if err := wb.PushContext("PlcTagItem"); err != nil {
		return err
	}
	if err := wb.PushContext("tag"); err != nil {
		return err
	}
	if err := wb.PushContext("SlmpTag"); err != nil {
		return err
	}
	deviceCodeName := m.DeviceCode.String()
	if err := wb.WriteString("deviceCode", uint32(len(deviceCodeName)*8), deviceCodeName, utils.WithEncoding("UTF-8")); err != nil {
		return err
	}
	if err := wb.WriteUint32("deviceNumber", 24, m.DeviceNumber); err != nil {
		return err
	}
	dataTypeName := m.DataType.String()
	if err := wb.WriteString("dataType", uint32(len(dataTypeName)*8), dataTypeName, utils.WithEncoding("UTF-8")); err != nil {
		return err
	}
	if err := wb.WriteUint16("quantity", 16, m.Quantity); err != nil {
		return err
	}
	if err := wb.PopContext("SlmpTag"); err != nil {
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

// castToSlmpTagFromPlcTag narrows a tag out of a request back to an slmp tag.
func castToSlmpTagFromPlcTag(tag apiModel.PlcTag) (PlcTag, error) {
	if slmpTag, ok := tag.(PlcTag); ok {
		return slmpTag, nil
	}
	return nil, errors.Errorf("couldn't cast %T to an slmp tag", tag)
}
