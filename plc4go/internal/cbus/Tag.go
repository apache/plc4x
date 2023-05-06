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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type StatusRequestType uint8

const (
	StatusRequestTypeBinaryState StatusRequestType = iota
	StatusRequestTypeLevel
)

func (s StatusRequestType) String() string {
	switch s {
	case StatusRequestTypeBinaryState:
		return "StatusRequestTypeBinaryState"
	case StatusRequestTypeLevel:
		return "StatusRequestTypeLevel"
	}
	return ""
}

type Tag interface {
	apiModel.PlcTag

	GetTagType() TagType
}

// StatusTag can be used to query status using a P-to-MP-StatusRequest command
type StatusTag interface {
	Tag

	GetBridgeAddresses() []readWriteModel.BridgeAddress
	GetStatusRequestType() StatusRequestType
	GetStartingGroupAddressLabel() *byte
	GetApplication() readWriteModel.ApplicationIdContainer
}

func NewStatusTag(bridgeAddresses []readWriteModel.BridgeAddress, statusRequestType StatusRequestType, startingGroupAddressLabel *byte, application readWriteModel.ApplicationIdContainer, numElements uint16) StatusTag {
	return &statusTag{
		bridgeAddresses:           bridgeAddresses,
		tagType:                   STATUS,
		startingGroupAddressLabel: startingGroupAddressLabel,
		statusRequestType:         statusRequestType,
		application:               application,
		numElements:               numElements,
	}
}

type CalTag interface {
	GetBridgeAddresses() []readWriteModel.BridgeAddress
	GetUnitAddress() readWriteModel.UnitAddress
}

// CALRecallTag can be used to get device/network management tags
type CALRecallTag interface {
	Tag
	CalTag

	GetParameter() readWriteModel.Parameter
	GetCount() uint8
}

func NewCALRecallTag(unitAddress readWriteModel.UnitAddress, bridgeAddresses []readWriteModel.BridgeAddress, parameter readWriteModel.Parameter, count uint8, numElements uint16) CALRecallTag {
	return &calRecallTag{
		calTag:      calTag{bridgeAddresses, unitAddress},
		tagType:     CAL_RECALL,
		parameter:   parameter,
		count:       count,
		numElements: numElements,
	}
}

// CALIdentifyTag can be used to get device/network management tags
type CALIdentifyTag interface {
	Tag
	CalTag

	GetAttribute() readWriteModel.Attribute
}

func NewCALIdentifyTag(unitAddress readWriteModel.UnitAddress, bridgeAddresses []readWriteModel.BridgeAddress, attribute readWriteModel.Attribute, numElements uint16) CALIdentifyTag {
	return &calIdentifyTag{
		calTag:      calTag{bridgeAddresses, unitAddress},
		tagType:     CAL_IDENTIFY,
		attribute:   attribute,
		numElements: numElements,
	}
}

// CALGetStatusTag can be used to get device/network management tags
type CALGetStatusTag interface {
	Tag
	CalTag

	GetParameter() readWriteModel.Parameter
	GetCount() uint8
}

func NewCALGetStatusTag(unitAddress readWriteModel.UnitAddress, bridgeAddresses []readWriteModel.BridgeAddress, parameter readWriteModel.Parameter, count uint8, numElements uint16) CALGetStatusTag {
	return &calGetStatusTag{
		calTag:      calTag{bridgeAddresses, unitAddress},
		tagType:     CAL_GETSTATUS,
		parameter:   parameter,
		count:       count,
		numElements: numElements,
	}
}

// SALTag can be used to send SAL commands
type SALTag interface {
	Tag

	GetBridgeAddresses() []readWriteModel.BridgeAddress
	GetApplication() readWriteModel.ApplicationIdContainer
	GetSALCommand() string
}

func NewSALTag(bridgeAddresses []readWriteModel.BridgeAddress, application readWriteModel.ApplicationIdContainer, salCommand string, numElements uint16) SALTag {
	return &salTag{
		bridgeAddresses: bridgeAddresses,
		tagType:         SAL,
		application:     application,
		salCommand:      salCommand,
		numElements:     numElements,
	}
}

// SALMonitorTag can be used to monitor sal tags
type SALMonitorTag interface {
	Tag

	GetUnitAddress() readWriteModel.UnitAddress
	GetApplication() *readWriteModel.ApplicationIdContainer
}

func NewSALMonitorTag(unitAddress readWriteModel.UnitAddress, application *readWriteModel.ApplicationIdContainer, numElements uint16) SALMonitorTag {
	return &salMonitorTag{
		tagType:     SAL_MONITOR,
		unitAddress: unitAddress,
		application: application,
		numElements: numElements,
	}
}

// MMIMonitorTag can be used to monitor mmi tags
type MMIMonitorTag interface {
	Tag

	GetUnitAddress() readWriteModel.UnitAddress
	GetApplication() *readWriteModel.ApplicationIdContainer
}

func NewMMIMonitorTag(unitAddress readWriteModel.UnitAddress, application *readWriteModel.ApplicationIdContainer, numElements uint16) MMIMonitorTag {
	return &mmiMonitorTag{
		tagType:     MMI_STATUS_MONITOR,
		unitAddress: unitAddress,
		application: application,
		numElements: numElements,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type statusTag struct {
	bridgeAddresses           []readWriteModel.BridgeAddress
	tagType                   TagType
	statusRequestType         StatusRequestType
	startingGroupAddressLabel *byte
	application               readWriteModel.ApplicationIdContainer
	numElements               uint16
}

type calTag struct {
	bridgeAddresses []readWriteModel.BridgeAddress
	unitAddress     readWriteModel.UnitAddress
}

type calRecallTag struct {
	calTag
	tagType     TagType
	parameter   readWriteModel.Parameter
	count       uint8
	numElements uint16
}

type calIdentifyTag struct {
	calTag
	tagType     TagType
	attribute   readWriteModel.Attribute
	numElements uint16
}

type calGetStatusTag struct {
	calTag
	tagType     TagType
	parameter   readWriteModel.Parameter
	count       uint8
	numElements uint16
}

type salTag struct {
	bridgeAddresses []readWriteModel.BridgeAddress
	tagType         TagType
	application     readWriteModel.ApplicationIdContainer
	salCommand      string
	numElements     uint16
}

type salMonitorTag struct {
	tagType     TagType
	unitAddress readWriteModel.UnitAddress
	application *readWriteModel.ApplicationIdContainer
	numElements uint16
}

type mmiMonitorTag struct {
	tagType     TagType
	unitAddress readWriteModel.UnitAddress
	application *readWriteModel.ApplicationIdContainer
	numElements uint16
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (s statusTag) GetBridgeAddresses() []readWriteModel.BridgeAddress {
	return s.bridgeAddresses
}

func (s statusTag) GetAddressString() string {
	statusRequestType := ""
	switch s.statusRequestType {
	case StatusRequestTypeBinaryState:
		statusRequestType = "binary"
	case StatusRequestTypeLevel:
		statusRequestType = "level"
		if s.startingGroupAddressLabel != nil {
			statusRequestType += fmt.Sprintf("=%#02x", *s.startingGroupAddressLabel)
		}
	default:
		statusRequestType = "invalid"
	}
	return fmt.Sprintf("status/%s/%s", statusRequestType, s.application)
}

func (s statusTag) GetValueType() apiValues.PlcValueType {
	return apiValues.NULL
}

func (s statusTag) GetArrayInfo() []apiModel.ArrayInfo {
	if s.numElements != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(s.numElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (s statusTag) GetTagType() TagType {
	return s.tagType
}

func (s statusTag) GetStatusRequestType() StatusRequestType {
	return s.statusRequestType
}

func (s statusTag) GetStartingGroupAddressLabel() *byte {
	return s.startingGroupAddressLabel
}

func (s statusTag) GetApplication() readWriteModel.ApplicationIdContainer {
	return s.application
}

func (s statusTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := s.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (s statusTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(s.tagType.GetName()); err != nil {
		return err
	}

	if len(s.bridgeAddresses) > 0 {
		if err := writeBuffer.PushContext("bridgeAddresses"); err != nil {
			return err
		}
		for _, address := range s.bridgeAddresses {
			if err := address.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
		}
		if err := writeBuffer.PopContext("bridgeAddresses"); err != nil {
			return err
		}
	}

	if err := writeBuffer.WriteUint8("statusRequestType", 8, uint8(s.statusRequestType), utils.WithAdditionalStringRepresentation(s.statusRequestType.String())); err != nil {
		return err
	}
	if s.startingGroupAddressLabel != nil {
		if err := writeBuffer.WriteUint8("startingGroupAddressLabel", 8, *s.startingGroupAddressLabel); err != nil {
			return err
		}
	}
	if err := writeBuffer.WriteUint8("application", 8, uint8(s.application), utils.WithAdditionalStringRepresentation(s.application.String())); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(s.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (s statusTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), s); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calTag) GetBridgeAddresses() []readWriteModel.BridgeAddress {
	return c.bridgeAddresses
}

func (c calTag) GetUnitAddress() readWriteModel.UnitAddress {
	return c.unitAddress
}

func (c calTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := c.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (c calTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if len(c.bridgeAddresses) > 0 {
		if err := writeBuffer.PushContext("bridgeAddresses"); err != nil {
			return err
		}
		for _, address := range c.bridgeAddresses {
			if err := address.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
		}
		if err := writeBuffer.PopContext("bridgeAddresses"); err != nil {
			return err
		}
	}
	if unitAddress := c.unitAddress; unitAddress != nil {
		return c.unitAddress.SerializeWithWriteBuffer(ctx, writeBuffer)
	}
	return nil
}

func (c calTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calRecallTag) GetParameter() readWriteModel.Parameter {
	return c.parameter
}

func (c calRecallTag) GetCount() uint8 {
	return c.count
}

func (c calRecallTag) GetAddressString() string {
	return fmt.Sprintf("cal/%d/recall=%s", c.unitAddress.GetAddress(), c.parameter)
}

func (c calRecallTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (c calRecallTag) GetArrayInfo() []apiModel.ArrayInfo {
	if c.count != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(c.count),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (c calRecallTag) GetTagType() TagType {
	return c.tagType
}

func (c calRecallTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := c.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (c calRecallTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(c.tagType.GetName()); err != nil {
		return err
	}

	if err := c.calTag.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := c.parameter.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("count", 8, c.count); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(c.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calRecallTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calIdentifyTag) GetAttribute() readWriteModel.Attribute {
	return c.attribute
}

func (c calIdentifyTag) GetAddressString() string {
	return fmt.Sprintf("cal/%d/identify=%s", c.unitAddress.GetAddress(), c.GetAttribute())
}

func (c calIdentifyTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (c calIdentifyTag) GetArrayInfo() []apiModel.ArrayInfo {
	if c.numElements != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(c.numElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (c calIdentifyTag) GetTagType() TagType {
	return c.tagType
}

func (c calIdentifyTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := c.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (c calIdentifyTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(c.tagType.GetName()); err != nil {
		return err
	}

	if err := c.calTag.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := c.attribute.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(c.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calIdentifyTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calGetStatusTag) GetParameter() readWriteModel.Parameter {
	return c.parameter
}

func (c calGetStatusTag) GetCount() uint8 {
	return c.count
}

func (c calGetStatusTag) GetAddressString() string {
	return fmt.Sprintf("cal/getstatus=%s, %d", c.parameter, c.GetCount())
}

func (c calGetStatusTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (c calGetStatusTag) GetArrayInfo() []apiModel.ArrayInfo {
	if c.count != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(c.count),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (c calGetStatusTag) GetTagType() TagType {
	return c.tagType
}

func (c calGetStatusTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := c.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (c calGetStatusTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(c.tagType.GetName()); err != nil {
		return err
	}

	if err := c.calTag.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := c.parameter.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("count", 8, c.count); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(c.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calGetStatusTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (s salTag) GetBridgeAddresses() []readWriteModel.BridgeAddress {
	return s.bridgeAddresses
}

func (s salTag) GetApplication() readWriteModel.ApplicationIdContainer {
	return s.application
}

func (s salTag) GetSALCommand() string {
	return s.salCommand
}

func (s salTag) GetAddressString() string {
	return fmt.Sprintf("sal/%s/%s", s.application, s.salCommand)
}

func (s salTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (s salTag) GetArrayInfo() []apiModel.ArrayInfo {
	if s.numElements != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(s.numElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (s salTag) GetTagType() TagType {
	return s.tagType
}

func (s salTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := s.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (s salTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(s.tagType.GetName()); err != nil {
		return err
	}

	if len(s.bridgeAddresses) > 0 {
		if err := writeBuffer.PushContext("bridgeAddresses"); err != nil {
			return err
		}
		for _, address := range s.bridgeAddresses {
			if err := address.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
		}
		if err := writeBuffer.PopContext("bridgeAddresses"); err != nil {
			return err
		}
	}

	if err := s.application.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("salCommand", uint32(len(s.salCommand)*8), "UTF-8", s.salCommand); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(s.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (s salTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), s); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (s salMonitorTag) GetAddressString() string {
	unitAddress := "*"
	if s.unitAddress != nil {
		unitAddress = fmt.Sprintf("%#02x", s.unitAddress.GetAddress())
	}
	application := "*"
	if s.application != nil {
		application = fmt.Sprintf("%s", *s.application)
	}
	return fmt.Sprintf("salmonitor/%s/%s", unitAddress, application)
}

func (s salMonitorTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (s salMonitorTag) GetArrayInfo() []apiModel.ArrayInfo {
	if s.numElements != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(s.numElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (s salMonitorTag) GetTagType() TagType {
	return s.tagType
}

func (s salMonitorTag) GetUnitAddress() readWriteModel.UnitAddress {
	return s.unitAddress
}

func (s salMonitorTag) GetApplication() *readWriteModel.ApplicationIdContainer {
	return s.application
}

func (s salMonitorTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := s.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (s salMonitorTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(s.tagType.GetName()); err != nil {
		return err
	}

	if unitAddress := s.unitAddress; unitAddress != nil {
		if err := unitAddress.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}
	if application := s.application; application != nil {
		if err := application.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(s.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (s salMonitorTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), s); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (m mmiMonitorTag) GetAddressString() string {
	unitAddress := "*"
	if m.unitAddress != nil {
		unitAddress = fmt.Sprintf("%#02x", m.unitAddress.GetAddress())
	}
	application := "*"
	if m.application != nil {
		application = fmt.Sprintf("%s", *m.application)
	}
	return fmt.Sprintf("mmimonitor/%s/%s", unitAddress, application)
}

func (m mmiMonitorTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (m mmiMonitorTag) GetArrayInfo() []apiModel.ArrayInfo {
	if m.numElements != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(m.numElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (m mmiMonitorTag) GetTagType() TagType {
	return m.tagType
}

func (m mmiMonitorTag) GetUnitAddress() readWriteModel.UnitAddress {
	return m.unitAddress
}

func (m mmiMonitorTag) GetApplication() *readWriteModel.ApplicationIdContainer {
	return m.application
}

func (m mmiMonitorTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m mmiMonitorTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.tagType.GetName()); err != nil {
		return err
	}

	if unitAddress := m.unitAddress; unitAddress != nil {
		if err := unitAddress.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}
	if application := m.application; application != nil {
		if err := application.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(m.tagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m mmiMonitorTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
