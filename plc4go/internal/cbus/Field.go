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
	GetApplicationId() readWriteModel.ApplicationId
}

func NewStatusField(statusRequestType StatusRequestType, level *byte, applicationId readWriteModel.ApplicationId, numElements uint16) StatusField {
	return &statusField{
		fieldType:         STATUS,
		statusRequestType: statusRequestType,
		applicationId:     applicationId,
		numElements:       numElements,
	}
}

// CALField can be used to get device/network management fields
type CALField interface {
	model.PlcField
}

func NewCALField(numElements uint16) CALField {
	return &calField{
		fieldType:   CAL,
		numElements: numElements,
	}
}

type statusField struct {
	fieldType         FieldType
	statusRequestType StatusRequestType
	level             *byte
	applicationId     readWriteModel.ApplicationId
	numElements       uint16
}

func (m statusField) GetAddressString() string {
	return fmt.Sprintf("%d[%d]", m.fieldType, m.numElements)
}

func (m statusField) GetStatusRequestType() StatusRequestType {
	return m.statusRequestType
}

func (m statusField) GetApplicationId() readWriteModel.ApplicationId {
	return m.applicationId
}

func (m statusField) GetTypeName() string {
	return STATUS.GetName()
}

func (m statusField) GetQuantity() uint16 {
	return m.numElements
}

type calField struct {
	fieldType   FieldType
	numElements uint16
}

func (m calField) GetAddressString() string {
	return fmt.Sprintf("%d[%d]", m.fieldType, m.numElements)
}

func (m calField) GetTypeName() string {
	return CAL.GetName()
}

func (m calField) GetQuantity() uint16 {
	return m.numElements
}

func (m calField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.fieldType.GetName()); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.fieldType.GetName()); err != nil {
		return err
	}
	return nil
}
