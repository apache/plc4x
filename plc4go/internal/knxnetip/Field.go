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

package knxnetip

import (
	"errors"
	"fmt"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"strconv"
	"strings"
)

type Field interface {
	apiModel.PlcField
}

type GroupAddressField interface {
	GetFieldType() *driverModel.KnxDatapointType
	IsPatternField() bool
	matches(knxGroupAddress driverModel.KnxGroupAddress) bool
	toGroupAddress() driverModel.KnxGroupAddress
	Field
}

type DeviceField interface {
	toKnxAddress() driverModel.KnxAddress
	Field
}

type GroupAddress3LevelPlcField struct {
	MainGroup   string // 5 Bits: Values 0-31
	MiddleGroup string // 3 Bits: values 0-7
	SubGroup    string // 8 Bits
	FieldType   *driverModel.KnxDatapointType
}

func NewGroupAddress3LevelPlcField(mainGroup string, middleGroup string, subGroup string, fieldType *driverModel.KnxDatapointType) GroupAddress3LevelPlcField {
	return GroupAddress3LevelPlcField{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
		FieldType:   fieldType,
	}
}

func (k GroupAddress3LevelPlcField) GetAddressString() string {
	return fmt.Sprintf("%s/%s/%s:%s", k.MainGroup, k.MiddleGroup, k.SubGroup, k.FieldType.String())
}

func (k GroupAddress3LevelPlcField) GetTypeName() string {
	return k.FieldType.Name()
}

func (k GroupAddress3LevelPlcField) GetFieldType() *driverModel.KnxDatapointType {
	return k.FieldType
}

func (k GroupAddress3LevelPlcField) GetQuantity() uint16 {
	return 1
}

func (k GroupAddress3LevelPlcField) IsPatternField() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		_, err = strconv.Atoi(k.MiddleGroup)
		if err == nil {
			_, err = strconv.Atoi(k.SubGroup)
			if err == nil {
				return false
			}
		}
	}
	return true
}

func (k GroupAddress3LevelPlcField) matches(knxGroupAddress driverModel.KnxGroupAddress) bool {
	level3KnxGroupAddress := driverModel.CastKnxGroupAddress3Level(knxGroupAddress)
	if level3KnxGroupAddress == nil {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level3KnxGroupAddress.GetMainGroup()))) &&
		matches(k.MiddleGroup, strconv.Itoa(int(level3KnxGroupAddress.GetMiddleGroup()))) &&
		matches(k.SubGroup, strconv.Itoa(int(level3KnxGroupAddress.GetSubGroup())))
}

func (k GroupAddress3LevelPlcField) toGroupAddress() driverModel.KnxGroupAddress {
	mainGroup, err := strconv.ParseUint(k.MainGroup, 10, 8)
	if err != nil {
		return nil
	}
	middleGroup, err := strconv.ParseUint(k.MiddleGroup, 10, 8)
	if err != nil {
		return nil
	}
	subGroup, err := strconv.ParseUint(k.SubGroup, 10, 8)
	if err != nil {
		return nil
	}
	return driverModel.NewKnxGroupAddress3Level(uint8(mainGroup), uint8(middleGroup), uint8(subGroup))
}

type GroupAddress2LevelPlcField struct {
	MainGroup string // 5 Bits: Values 0-31
	SubGroup  string // 11 Bits
	FieldType *driverModel.KnxDatapointType
	Field
}

func NewGroupAddress2LevelPlcField(mainGroup string, subGroup string, fieldType *driverModel.KnxDatapointType) GroupAddress2LevelPlcField {
	return GroupAddress2LevelPlcField{
		MainGroup: mainGroup,
		SubGroup:  subGroup,
		FieldType: fieldType,
	}
}

func (k GroupAddress2LevelPlcField) GetAddressString() string {
	return fmt.Sprintf("%s/%s:%s", k.MainGroup, k.SubGroup, k.FieldType.String())
}

func (k GroupAddress2LevelPlcField) GetTypeName() string {
	return k.FieldType.Name()
}

func (k GroupAddress2LevelPlcField) GetFieldType() *driverModel.KnxDatapointType {
	return k.FieldType
}

func (k GroupAddress2LevelPlcField) GetQuantity() uint16 {
	return 1
}

func (k GroupAddress2LevelPlcField) IsPatternField() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		_, err = strconv.Atoi(k.SubGroup)
		if err == nil {
			return false
		}
	}
	return true
}

func (k GroupAddress2LevelPlcField) matches(knxGroupAddress driverModel.KnxGroupAddress) bool {
	level2KnxGroupAddress := driverModel.CastKnxGroupAddress2Level(knxGroupAddress)
	if level2KnxGroupAddress == nil {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level2KnxGroupAddress.GetMainGroup()))) &&
		matches(k.SubGroup, strconv.Itoa(int(level2KnxGroupAddress.GetSubGroup())))
}

func (k GroupAddress2LevelPlcField) toGroupAddress() driverModel.KnxGroupAddress {
	mainGroup, err := strconv.ParseUint(k.MainGroup, 10, 8)
	if err != nil {
		return nil
	}
	subGroup, err := strconv.ParseUint(k.SubGroup, 10, 16)
	if err != nil {
		return nil
	}
	return driverModel.NewKnxGroupAddress2Level(uint8(mainGroup), uint16(subGroup))
}

type GroupAddress1LevelPlcField struct {
	MainGroup string // 16 Bits
	FieldType *driverModel.KnxDatapointType
	Field
}

func NewGroupAddress1LevelPlcField(mainGroup string, fieldType *driverModel.KnxDatapointType) GroupAddress1LevelPlcField {
	return GroupAddress1LevelPlcField{
		MainGroup: mainGroup,
		FieldType: fieldType,
	}
}

func (k GroupAddress1LevelPlcField) GetAddressString() string {
	return fmt.Sprintf("%s:%s", k.MainGroup, k.FieldType.String())
}

func (k GroupAddress1LevelPlcField) GetTypeName() string {
	return k.FieldType.Name()
}

func (k GroupAddress1LevelPlcField) GetFieldType() *driverModel.KnxDatapointType {
	return k.FieldType
}

func (k GroupAddress1LevelPlcField) GetQuantity() uint16 {
	return 1
}

func (k GroupAddress1LevelPlcField) IsPatternField() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		return false
	}
	return true
}

func (k GroupAddress1LevelPlcField) matches(knxGroupAddress driverModel.KnxGroupAddress) bool {
	level1KnxGroupAddress := driverModel.CastKnxGroupAddressFreeLevel(knxGroupAddress)
	if level1KnxGroupAddress == nil {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level1KnxGroupAddress.GetSubGroup())))
}

func (k GroupAddress1LevelPlcField) toGroupAddress() driverModel.KnxGroupAddress {
	mainGroup, err := strconv.ParseUint(k.MainGroup, 10, 16)
	if err != nil {
		return nil
	}
	return driverModel.NewKnxGroupAddressFreeLevel(uint16(mainGroup))
}

type DeviceQueryField struct {
	MainGroup   string // 5 Bits: Values 0-31
	MiddleGroup string // 3 Bits: values 0-7
	SubGroup    string // 8 Bits
	DeviceField
}

func NewDeviceQueryField(mainGroup string, middleGroup string, subGroup string) DeviceQueryField {
	return DeviceQueryField{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
	}
}

func (k DeviceQueryField) GetAddressString() string {
	return fmt.Sprintf("%s.%s.%s", k.MainGroup, k.MiddleGroup, k.SubGroup)
}

func (k DeviceQueryField) GetTypeName() string {
	return ""
}

func (k DeviceQueryField) GetQuantity() uint16 {
	return 1
}

func (k DeviceQueryField) toKnxAddress() driverModel.KnxAddress {
	return nil
}

type DevicePropertyAddressPlcField struct {
	MainGroup     uint8 // 5 Bits: Values 0-31
	MiddleGroup   uint8 // 3 Bits: values 0-7
	SubGroup      uint8 // 8 Bits
	ObjectId      uint8
	PropertyId    uint8
	PropertyIndex uint16
	NumElements   uint8
	DeviceField
}

func NewDevicePropertyAddressPlcField(mainGroup uint8, middleGroup uint8, subGroup uint8, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) DevicePropertyAddressPlcField {
	return DevicePropertyAddressPlcField{
		MainGroup:     mainGroup,
		MiddleGroup:   middleGroup,
		SubGroup:      subGroup,
		ObjectId:      objectId,
		PropertyId:    propertyId,
		PropertyIndex: propertyIndex,
		NumElements:   numElements,
	}
}

func (k DevicePropertyAddressPlcField) GetAddressString() string {
	return fmt.Sprintf("%d/%d/%d#%d/%d/%d[%d]",
		k.MainGroup, k.MiddleGroup, k.SubGroup, k.ObjectId, k.PropertyId, k.PropertyIndex, k.NumElements)
}

func (k DevicePropertyAddressPlcField) GetTypeName() string {
	return ""
}

func (k DevicePropertyAddressPlcField) GetQuantity() uint16 {
	return 1
}

func (k DevicePropertyAddressPlcField) toKnxAddress() driverModel.KnxAddress {
	return driverModel.NewKnxAddress(
		k.MainGroup,
		k.MiddleGroup,
		k.SubGroup,
	)
}

type DeviceMemoryAddressPlcField struct {
	MainGroup   uint8 // 5 Bits: Values 0-31
	MiddleGroup uint8 // 3 Bits: values 0-7
	SubGroup    uint8 // 8 Bits
	Address     uint16
	NumElements uint8
	FieldType   *driverModel.KnxDatapointType
	DeviceField
}

func NewDeviceMemoryAddressPlcField(mainGroup uint8, middleGroup uint8, subGroup uint8, address uint16, numElements uint8, fieldType *driverModel.KnxDatapointType) DeviceMemoryAddressPlcField {
	return DeviceMemoryAddressPlcField{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
		Address:     address,
		NumElements: numElements,
		FieldType:   fieldType,
	}
}

func (k DeviceMemoryAddressPlcField) GetAddressString() string {
	return fmt.Sprintf("%d/%d/%d#%d:%s[%d]",
		k.MainGroup, k.MiddleGroup, k.SubGroup, k.Address, k.FieldType.String(), k.NumElements)
}

func (k DeviceMemoryAddressPlcField) GetTypeName() string {
	if k.FieldType != nil {
		return k.FieldType.Name()
	}
	return ""
}

func (k DeviceMemoryAddressPlcField) GetFieldType() *driverModel.KnxDatapointType {
	return k.FieldType
}

func (k DeviceMemoryAddressPlcField) GetQuantity() uint16 {
	return uint16(k.NumElements)
}

func (k DeviceMemoryAddressPlcField) toKnxAddress() driverModel.KnxAddress {
	individualAddress := driverModel.NewKnxAddress(
		k.MainGroup,
		k.MiddleGroup,
		k.SubGroup,
	)
	return individualAddress
}

type CommunicationObjectQueryField struct {
	MainGroup   uint8 // 5 Bits: Values 0-31
	MiddleGroup uint8 // 3 Bits: values 0-7
	SubGroup    uint8 // 8 Bits
	DeviceField
}

func NewCommunicationObjectQueryField(mainGroup uint8, middleGroup uint8, subGroup uint8) CommunicationObjectQueryField {
	return CommunicationObjectQueryField{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
	}
}

func (k CommunicationObjectQueryField) GetAddressString() string {
	return fmt.Sprintf("%d.%d.%d#com-obj",
		k.MainGroup, k.MiddleGroup, k.SubGroup)
}

func (k CommunicationObjectQueryField) GetTypeName() string {
	return ""
}

func (k CommunicationObjectQueryField) GetQuantity() uint16 {
	return 1
}

func (k CommunicationObjectQueryField) toKnxAddress() driverModel.KnxAddress {
	individualAddress := driverModel.NewKnxAddress(
		k.MainGroup,
		k.MiddleGroup,
		k.SubGroup,
	)
	return individualAddress
}

func CastToFieldFromPlcField(plcField apiModel.PlcField) (Field, error) {
	if field, ok := plcField.(Field); ok {
		return field, nil
	}
	return nil, errors.New("couldn't cast to KnxNetIpField")
}

func matches(pattern string, groupAddressPart string) bool {
	// A "*" simply matches everything
	if pattern == "*" {
		return true
	}
	// If the pattern starts and ends with square brackets, it's a list of values or range queries
	if strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]") {
		matches := false
		for _, segment := range strings.Split(pattern, ",") {
			if strings.Contains(segment, "-") {
				// If the segment contains a "-", then it's a range query
				split := strings.Split(segment, "-")
				if len(split) == 2 {
					if val, err := strconv.Atoi(groupAddressPart); err != nil {
						var err error
						var from int
						if from, err = strconv.Atoi(split[0]); err != nil {
							continue
						}
						if val < from {
							continue
						}
						var to int
						if to, err = strconv.Atoi(split[1]); err == nil {
							continue
						}
						if val > to {
							continue
						}
						matches = true
					}
				}
			} else if segment == groupAddressPart {
				// In all other cases it's an explicit value
				matches = true
			}
		}
		return matches
	} else {
		return pattern == groupAddressPart
	}
}
