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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
)

type StatusRequestType uint8

const (
	StatusRequestTypeBinaryState StatusRequestType = iota
	StatusRequestTypeLevel
)

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

// CALRecallField can be used to get device/network management fields
type CALRecallField interface {
	model.PlcField
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

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (m statusField) GetAddressString() string {
	return fmt.Sprintf("%d[%d]", m.fieldType, m.numElements)
}

func (m statusField) GetStatusRequestType() StatusRequestType {
	return m.statusRequestType
}

func (m statusField) GetStartingGroupAddressLabel() *byte {
	return m.startingGroupAddressLabel
}

func (m statusField) GetApplication() readWriteModel.ApplicationIdContainer {
	return m.application
}

func (m statusField) GetTypeName() string {
	return STATUS.GetName()
}

func (m statusField) GetQuantity() uint16 {
	return m.numElements
}

func (m statusField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.fieldType.GetName()); err != nil {
		return err
	}

	// TODO: add string representation
	if err := writeBuffer.WriteUint8("statusRequestType", 8, uint8(m.statusRequestType)); err != nil {
		return err
	}
	if m.startingGroupAddressLabel != nil {
		if err := writeBuffer.WriteUint8("startingGroupAddressLabel", 8, *m.startingGroupAddressLabel); err != nil {
			return err
		}
	}
	if err := writeBuffer.WriteUint8("application", 8, uint8(m.application), utils.WithAdditionalStringRepresentation(m.application.String())); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m calField) Serialize(writeBuffer utils.WriteBuffer) error {
	return m.unitAddress.Serialize(writeBuffer)
}

func (m calRecallField) GetParameter() readWriteModel.Parameter {
	return m.parameter
}

func (m calRecallField) GetCount() uint8 {
	return m.count
}

func (m calRecallField) GetAddressString() string {
	return fmt.Sprintf("%d[%d]", m.fieldType, m.numElements)
}

func (m calRecallField) GetTypeName() string {
	return m.fieldType.GetName()
}

func (m calRecallField) GetQuantity() uint16 {
	return m.numElements
}

func (m calRecallField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.fieldType.GetName()); err != nil {
		return err
	}

	if err := m.calField.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("parameter", 8, uint8(m.parameter), utils.WithAdditionalStringRepresentation(m.parameter.String())); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint8("count", 8, m.count); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (c calIdentifyField) GetAttribute() readWriteModel.Attribute {
	return c.attribute
}

func (c calIdentifyField) GetAddressString() string {
	return fmt.Sprintf("%d[%d]", c.fieldType, c.numElements)
}

func (c calIdentifyField) GetTypeName() string {
	return c.fieldType.GetName()
}

func (c calIdentifyField) GetQuantity() uint16 {
	return c.numElements
}

func (m calIdentifyField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.fieldType.GetName()); err != nil {
		return err
	}

	if err := m.calField.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("attribute", 8, uint8(m.attribute), utils.WithAdditionalStringRepresentation(m.attribute.String())); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m calGetstatusField) GetParameter() readWriteModel.Parameter {
	return m.parameter
}

func (m calGetstatusField) GetCount() uint8 {
	return m.count
}

func (c calGetstatusField) GetAddressString() string {
	return fmt.Sprintf("%d[%d]", c.fieldType, c.numElements)
}

func (c calGetstatusField) GetTypeName() string {
	return c.fieldType.GetName()
}

func (c calGetstatusField) GetQuantity() uint16 {
	return c.numElements
}

func (m calGetstatusField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.fieldType.GetName()); err != nil {
		return err
	}

	if err := m.calField.Serialize(writeBuffer); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint8("parameter", 8, uint8(m.parameter), utils.WithAdditionalStringRepresentation(m.parameter.String())); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint8("count", 8, m.count); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}
