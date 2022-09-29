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
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
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

// StatusField can be used to query status using a P-to-MP-StatusRequest command
type StatusField interface {
	model.PlcField
	GetStatusRequestType() StatusRequestType
	GetStartingGroupAddressLabel() *byte
	GetApplication() readWriteModel.ApplicationIdContainer
}

func NewStatusField(statusRequestType StatusRequestType, startingGroupAddressLabel *byte, application readWriteModel.ApplicationIdContainer, numElements uint16) StatusField {
	return &statusField{
		fieldType:                 STATUS,
		startingGroupAddressLabel: startingGroupAddressLabel,
		statusRequestType:         statusRequestType,
		application:               application,
		numElements:               numElements,
	}
}

type CalField interface {
	GetUnitAddress() readWriteModel.UnitAddress
}

// CALRecallField can be used to get device/network management fields
type CALRecallField interface {
	model.PlcField
	CalField
	GetParameter() readWriteModel.Parameter
	GetCount() uint8
}

func NewCALRecallField(unitAddress readWriteModel.UnitAddress, parameter readWriteModel.Parameter, count uint8, numElements uint16) CALRecallField {
	return &calRecallField{
		calField:    calField{unitAddress: unitAddress},
		fieldType:   CAL_RECALL,
		parameter:   parameter,
		count:       count,
		numElements: numElements,
	}
}

// CALIdentifyField can be used to get device/network management fields
type CALIdentifyField interface {
	model.PlcField
	CalField
	GetAttribute() readWriteModel.Attribute
}

func NewCALIdentifyField(unitAddress readWriteModel.UnitAddress, attribute readWriteModel.Attribute, numElements uint16) CALIdentifyField {
	return &calIdentifyField{
		calField:    calField{unitAddress: unitAddress},
		fieldType:   CAL_IDENTIFY,
		attribute:   attribute,
		numElements: numElements,
	}
}

// CALGetstatusField can be used to get device/network management fields
type CALGetstatusField interface {
	model.PlcField
	CalField
	GetParameter() readWriteModel.Parameter
	GetCount() uint8
}

func NewCALGetstatusField(unitAddress readWriteModel.UnitAddress, parameter readWriteModel.Parameter, count uint8, numElements uint16) CALGetstatusField {
	return &calGetstatusField{
		calField:    calField{unitAddress: unitAddress},
		fieldType:   CAL_RECALL,
		parameter:   parameter,
		count:       count,
		numElements: numElements,
	}
}

// SALField can be used to send SAL commands
type SALField interface {
	model.PlcField
	GetApplication() readWriteModel.ApplicationIdContainer
	GetSALCommand() string
}

func NewSALField(application readWriteModel.ApplicationIdContainer, salCommand string, numElements uint16) SALField {
	return &salField{
		fieldType:   SAL,
		application: application,
		salCommand:  salCommand,
		numElements: numElements,
	}
}

// SALMonitorField can be used to monitor sal fields
type SALMonitorField interface {
	model.PlcField
	GetUnitAddress() *readWriteModel.UnitAddress
	GetApplication() *readWriteModel.ApplicationIdContainer
}

func NewSALMonitorField(unitAddress *readWriteModel.UnitAddress, application *readWriteModel.ApplicationIdContainer, numElements uint16) SALMonitorField {
	return &salMonitorField{
		fieldType:   SAL_MONITOR,
		unitAddress: unitAddress,
		application: application,
		numElements: numElements,
	}
}

// MMIMonitorField can be used to monitor mmi fields
type MMIMonitorField interface {
	model.PlcField
	GetUnitAddress() *readWriteModel.UnitAddress
	GetApplication() *readWriteModel.ApplicationIdContainer
}

func NewMMIMonitorField(unitAddress *readWriteModel.UnitAddress, application *readWriteModel.ApplicationIdContainer, numElements uint16) SALMonitorField {
	return &mmiMonitorField{
		fieldType:   MMI_STATUS_MONITOR,
		unitAddress: unitAddress,
		application: application,
		numElements: numElements,
	}
}

// UnitInfoField can be used to get information about unit(s)
type UnitInfoField interface {
	model.PlcField
	GetUnitAddress() *readWriteModel.UnitAddress
	GetAttribute() *readWriteModel.Attribute
}

func NewUnitInfoField(unitAddress *readWriteModel.UnitAddress, attribute *readWriteModel.Attribute, numElements uint16) UnitInfoField {
	return &unitInfoField{
		unitAddress: unitAddress,
		fieldType:   UNIT_INFO,
		attribute:   attribute,
		numElements: numElements,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type statusField struct {
	fieldType                 FieldType
	statusRequestType         StatusRequestType
	startingGroupAddressLabel *byte
	application               readWriteModel.ApplicationIdContainer
	numElements               uint16
}

type calField struct {
	unitAddress readWriteModel.UnitAddress
}

type calRecallField struct {
	calField
	fieldType   FieldType
	parameter   readWriteModel.Parameter
	count       uint8
	numElements uint16
}

type calIdentifyField struct {
	calField
	fieldType   FieldType
	attribute   readWriteModel.Attribute
	numElements uint16
}

type calGetstatusField struct {
	calField
	fieldType   FieldType
	parameter   readWriteModel.Parameter
	count       uint8
	numElements uint16
}

type salField struct {
	fieldType   FieldType
	application readWriteModel.ApplicationIdContainer
	salCommand  string
	numElements uint16
}
type salMonitorField struct {
	fieldType   FieldType
	unitAddress *readWriteModel.UnitAddress
	application *readWriteModel.ApplicationIdContainer
	numElements uint16
}

type mmiMonitorField struct {
	fieldType   FieldType
	unitAddress *readWriteModel.UnitAddress
	application *readWriteModel.ApplicationIdContainer
	numElements uint16
}

type unitInfoField struct {
	fieldType   FieldType
	unitAddress *readWriteModel.UnitAddress
	attribute   *readWriteModel.Attribute
	numElements uint16
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (s statusField) GetAddressString() string {
	statusRequestType := ""
	switch s.statusRequestType {
	case StatusRequestTypeBinaryState:
		statusRequestType = "binary"
	case StatusRequestTypeLevel:
		statusRequestType = "level"
		statusRequestType += fmt.Sprintf("=0x%x", *s.startingGroupAddressLabel)
	}
	return fmt.Sprintf("status/%s/%s", statusRequestType, s.application)
}

func (s statusField) GetStatusRequestType() StatusRequestType {
	return s.statusRequestType
}

func (s statusField) GetStartingGroupAddressLabel() *byte {
	return s.startingGroupAddressLabel
}

func (s statusField) GetApplication() readWriteModel.ApplicationIdContainer {
	return s.application
}

func (s statusField) GetTypeName() string {
	return STATUS.GetName()
}

func (s statusField) GetQuantity() uint16 {
	return s.numElements
}

func (s statusField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(s.fieldType.GetName()); err != nil {
		return err
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

	if err := writeBuffer.PopContext(s.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (s statusField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(s); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calField) GetUnitAddress() readWriteModel.UnitAddress {
	return c.unitAddress
}

func (c calField) Serialize(writeBuffer utils.WriteBuffer) error {
	if unitAddress := c.unitAddress; unitAddress != nil {
		return c.unitAddress.Serialize(writeBuffer)
	}
	return nil
}

func (c calField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calRecallField) GetParameter() readWriteModel.Parameter {
	return c.parameter
}

func (c calRecallField) GetCount() uint8 {
	return c.count
}

func (c calRecallField) GetAddressString() string {
	return fmt.Sprintf("cal/%d/recall=%s", c.unitAddress.GetAddress(), c.parameter)
}

func (c calRecallField) GetTypeName() string {
	return c.fieldType.GetName()
}

func (c calRecallField) GetQuantity() uint16 {
	return c.numElements
}

func (c calRecallField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(c.fieldType.GetName()); err != nil {
		return err
	}

	if err := c.calField.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := c.parameter.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("count", 8, c.count); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(c.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calRecallField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calIdentifyField) GetAttribute() readWriteModel.Attribute {
	return c.attribute
}

func (c calIdentifyField) GetAddressString() string {
	return fmt.Sprintf("cal/%d/identify=%s", c.unitAddress.GetAddress(), c.GetAttribute())
}

func (c calIdentifyField) GetTypeName() string {
	return c.fieldType.GetName()
}

func (c calIdentifyField) GetQuantity() uint16 {
	return c.numElements
}

func (c calIdentifyField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(c.fieldType.GetName()); err != nil {
		return err
	}

	if err := c.calField.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := c.attribute.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(c.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calIdentifyField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (c calGetstatusField) GetParameter() readWriteModel.Parameter {
	return c.parameter
}

func (c calGetstatusField) GetCount() uint8 {
	return c.count
}

func (c calGetstatusField) GetAddressString() string {
	return fmt.Sprintf("cal/getstatus=%s, %d", c.parameter, c.GetCount())
}

func (c calGetstatusField) GetTypeName() string {
	return c.fieldType.GetName()
}

func (c calGetstatusField) GetQuantity() uint16 {
	return c.numElements
}

func (c calGetstatusField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(c.fieldType.GetName()); err != nil {
		return err
	}

	if err := c.calField.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := c.parameter.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("count", 8, c.count); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(c.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calGetstatusField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(c); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (s salField) GetApplication() readWriteModel.ApplicationIdContainer {
	return s.application
}

func (s salField) GetSALCommand() string {
	return s.salCommand
}

func (s salField) GetAddressString() string {
	return fmt.Sprintf("sal/%s/%s", s.application, s.salCommand)
}

func (s salField) GetTypeName() string {
	return s.fieldType.GetName()
}

func (s salField) GetQuantity() uint16 {
	return s.numElements
}

func (s salField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(s.fieldType.GetName()); err != nil {
		return err
	}

	if err := s.application.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("salCommand", uint32(len(s.salCommand)*8), "UTF-8", s.salCommand); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(s.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (s salField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(s); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (s salMonitorField) GetAddressString() string {
	unitAddress := "*"
	if s.unitAddress != nil {
		unitAddress = fmt.Sprintf("%d", (*s.unitAddress).GetAddress())
	}
	application := "*"
	if s.application != nil {
		application = fmt.Sprintf("%d", *s.application)
	}
	return fmt.Sprintf("salmonitor/%s/%s", unitAddress, application)
}

func (s salMonitorField) GetTypeName() string {
	return s.fieldType.GetName()
}

func (s salMonitorField) GetQuantity() uint16 {
	return s.numElements
}

func (s salMonitorField) GetUnitAddress() *readWriteModel.UnitAddress {
	return s.unitAddress
}

func (s salMonitorField) GetApplication() *readWriteModel.ApplicationIdContainer {
	return s.application
}

func (s salMonitorField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(s.fieldType.GetName()); err != nil {
		return err
	}

	if unitAddress := s.unitAddress; unitAddress != nil {
		if err := (*unitAddress).Serialize(writeBuffer); err != nil {
			return err
		}
	}
	if application := s.application; application != nil {
		if err := application.Serialize(writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(s.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (s salMonitorField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(s); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (m mmiMonitorField) GetAddressString() string {
	unitAddress := "*"
	if m.unitAddress != nil {
		unitAddress = fmt.Sprintf("%d", (*m.unitAddress).GetAddress())
	}
	application := "*"
	if m.application != nil {
		application = fmt.Sprintf("%d", *m.application)
	}
	return fmt.Sprintf("mmimonitor/%s/%s", unitAddress, application)
}

func (m mmiMonitorField) GetTypeName() string {
	return m.fieldType.GetName()
}

func (m mmiMonitorField) GetQuantity() uint16 {
	return m.numElements
}

func (m mmiMonitorField) GetUnitAddress() *readWriteModel.UnitAddress {
	return m.unitAddress
}

func (m mmiMonitorField) GetApplication() *readWriteModel.ApplicationIdContainer {
	return m.application
}

func (m mmiMonitorField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.fieldType.GetName()); err != nil {
		return err
	}

	if unitAddress := m.unitAddress; unitAddress != nil {
		if err := (*unitAddress).Serialize(writeBuffer); err != nil {
			return err
		}
	}
	if application := m.application; application != nil {
		if err := application.Serialize(writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(m.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m mmiMonitorField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (u unitInfoField) GetUnitAddress() *readWriteModel.UnitAddress {
	return u.unitAddress
}

func (u unitInfoField) GetAttribute() *readWriteModel.Attribute {
	return u.attribute
}

func (u unitInfoField) GetAddressString() string {
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

func (u unitInfoField) GetTypeName() string {
	return u.fieldType.GetName()
}

func (u unitInfoField) GetQuantity() uint16 {
	return u.numElements
}

func (u unitInfoField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(u.fieldType.GetName()); err != nil {
		return err
	}

	if unitAddress := u.unitAddress; unitAddress != nil {
		if err := (*unitAddress).Serialize(writeBuffer); err != nil {
			return err
		}
	}

	if attribute := u.attribute; attribute != nil {
		if err := (*attribute).Serialize(writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(u.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (u unitInfoField) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(u); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
