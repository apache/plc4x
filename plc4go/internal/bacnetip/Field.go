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

package bacnetip

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"strconv"
)

type BacNetPlcField interface {
	GetDeviceIdentifier() uint32
	GetObjectType() uint16
	GetObjectInstance() uint32
	GetPropertyIdentifier() uint32
}

type PlcField struct {
	DeviceIdentifier   uint32
	ObjectType         uint16
	ObjectInstance     uint32
	PropertyIdentifier uint32
}

func (m PlcField) GetAddressString() string {
	return strconv.Itoa(int(m.DeviceIdentifier))
}

func (m PlcField) GetTypeName() string {
	return strconv.Itoa(int(m.ObjectType))
}

func (m PlcField) GetQuantity() uint16 {
	return 1
}

func NewField(deviceIdentifier uint32, objectType uint16, objectInstance uint32, propertyIdentifier uint32) PlcField {
	return PlcField{
		DeviceIdentifier:   deviceIdentifier,
		ObjectType:         objectType,
		ObjectInstance:     objectInstance,
		PropertyIdentifier: propertyIdentifier,
	}
}

func (m PlcField) GetDeviceIdentifier() uint32 {
	return m.DeviceIdentifier
}

func (m PlcField) GetObjectType() uint16 {
	return m.ObjectType
}

func (m PlcField) GetObjectInstance() uint32 {
	return m.ObjectInstance
}

func (m PlcField) GetPropertyIdentifier() uint32 {
	return m.PropertyIdentifier
}

func (m PlcField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("BacNetPlcField"); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint32("deviceIdentifier", 32, m.DeviceIdentifier); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint16("objectType", 16, m.ObjectType); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint32("objectInstance", 32, m.ObjectInstance); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint32("propertyIdentifier", 32, m.PropertyIdentifier); err != nil {
		return err
	}

	if err := writeBuffer.PopContext("BacNetPlcField"); err != nil {
		return err
	}
	return nil
}
