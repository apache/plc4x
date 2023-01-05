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
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
)

type Tag interface {
	apiModel.PlcTag
}

type GroupAddressTag interface {
	Tag

	GetTagType() *driverModel.KnxDatapointType
	IsPatternTag() bool
	matches(knxGroupAddress driverModel.KnxGroupAddress) bool
	toGroupAddress() driverModel.KnxGroupAddress
}

type DeviceTag interface {
	Tag

	toKnxAddress() driverModel.KnxAddress
}

type GroupAddress3LevelPlcTag struct {
	GroupAddressTag

	MainGroup   string // 5 Bits: Values 0-31
	MiddleGroup string // 3 Bits: values 0-7
	SubGroup    string // 8 Bits
	TagType     *driverModel.KnxDatapointType
}

func NewGroupAddress3LevelPlcTag(mainGroup string, middleGroup string, subGroup string, tagType *driverModel.KnxDatapointType) GroupAddress3LevelPlcTag {
	return GroupAddress3LevelPlcTag{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
		TagType:     tagType,
	}
}

func (k GroupAddress3LevelPlcTag) GetAddressString() string {
	return fmt.Sprintf("%s/%s/%s:%s", k.MainGroup, k.MiddleGroup, k.SubGroup, k.TagType.String())
}

func (k GroupAddress3LevelPlcTag) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k GroupAddress3LevelPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (k GroupAddress3LevelPlcTag) GetTagType() *driverModel.KnxDatapointType {
	return k.TagType
}

func (k GroupAddress3LevelPlcTag) IsPatternTag() bool {
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

func (k GroupAddress3LevelPlcTag) matches(knxGroupAddress driverModel.KnxGroupAddress) bool {
	level3KnxGroupAddress, ok := knxGroupAddress.(driverModel.KnxGroupAddress3LevelExactly)
	if !ok {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level3KnxGroupAddress.GetMainGroup()))) &&
		matches(k.MiddleGroup, strconv.Itoa(int(level3KnxGroupAddress.GetMiddleGroup()))) &&
		matches(k.SubGroup, strconv.Itoa(int(level3KnxGroupAddress.GetSubGroup())))
}

func (k GroupAddress3LevelPlcTag) toGroupAddress() driverModel.KnxGroupAddress {
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

type GroupAddress2LevelPlcTag struct {
	GroupAddressTag

	MainGroup string // 5 Bits: Values 0-31
	SubGroup  string // 11 Bits
	TagType   *driverModel.KnxDatapointType
}

func NewGroupAddress2LevelPlcTag(mainGroup string, subGroup string, tagType *driverModel.KnxDatapointType) GroupAddress2LevelPlcTag {
	return GroupAddress2LevelPlcTag{
		MainGroup: mainGroup,
		SubGroup:  subGroup,
		TagType:   tagType,
	}
}

func (k GroupAddress2LevelPlcTag) GetAddressString() string {
	return fmt.Sprintf("%s/%s:%s", k.MainGroup, k.SubGroup, k.TagType.String())
}

func (k GroupAddress2LevelPlcTag) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k GroupAddress2LevelPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (k GroupAddress2LevelPlcTag) GetTagType() *driverModel.KnxDatapointType {
	return k.TagType
}

func (k GroupAddress2LevelPlcTag) IsPatternTag() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		_, err = strconv.Atoi(k.SubGroup)
		if err == nil {
			return false
		}
	}
	return true
}

func (k GroupAddress2LevelPlcTag) matches(knxGroupAddress driverModel.KnxGroupAddress) bool {
	level2KnxGroupAddress, ok := knxGroupAddress.(driverModel.KnxGroupAddress2LevelExactly)
	if !ok {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level2KnxGroupAddress.GetMainGroup()))) &&
		matches(k.SubGroup, strconv.Itoa(int(level2KnxGroupAddress.GetSubGroup())))
}

func (k GroupAddress2LevelPlcTag) toGroupAddress() driverModel.KnxGroupAddress {
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

type GroupAddress1LevelPlcTag struct {
	GroupAddressTag

	MainGroup string // 16 Bits
	TagType   *driverModel.KnxDatapointType
}

func NewGroupAddress1LevelPlcTag(mainGroup string, tagType *driverModel.KnxDatapointType) GroupAddress1LevelPlcTag {
	return GroupAddress1LevelPlcTag{
		MainGroup: mainGroup,
		TagType:   tagType,
	}
}

func (k GroupAddress1LevelPlcTag) GetAddressString() string {
	return fmt.Sprintf("%s:%s", k.MainGroup, k.TagType.String())
}

func (k GroupAddress1LevelPlcTag) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k GroupAddress1LevelPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (k GroupAddress1LevelPlcTag) GetTagType() *driverModel.KnxDatapointType {
	return k.TagType
}

func (k GroupAddress1LevelPlcTag) IsPatternTag() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		return false
	}
	return true
}

func (k GroupAddress1LevelPlcTag) matches(knxGroupAddress driverModel.KnxGroupAddress) bool {
	level1KnxGroupAddress, ok := knxGroupAddress.(driverModel.KnxGroupAddressFreeLevel)
	if !ok {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level1KnxGroupAddress.GetSubGroup())))
}

func (k GroupAddress1LevelPlcTag) toGroupAddress() driverModel.KnxGroupAddress {
	mainGroup, err := strconv.ParseUint(k.MainGroup, 10, 16)
	if err != nil {
		return nil
	}
	return driverModel.NewKnxGroupAddressFreeLevel(uint16(mainGroup))
}

type DevicePropertyAddressPlcTag struct {
	MainGroup     uint8 // 5 Bits: Values 0-31
	MiddleGroup   uint8 // 3 Bits: values 0-7
	SubGroup      uint8 // 8 Bits
	ObjectId      uint8
	PropertyId    uint8
	PropertyIndex uint16
	NumElements   uint8
	DeviceTag
}

func NewDevicePropertyAddressPlcTag(mainGroup uint8, middleGroup uint8, subGroup uint8, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) DevicePropertyAddressPlcTag {
	return DevicePropertyAddressPlcTag{
		MainGroup:     mainGroup,
		MiddleGroup:   middleGroup,
		SubGroup:      subGroup,
		ObjectId:      objectId,
		PropertyId:    propertyId,
		PropertyIndex: propertyIndex,
		NumElements:   numElements,
	}
}

func (k DevicePropertyAddressPlcTag) GetAddressString() string {
	return fmt.Sprintf("%d/%d/%d#%d/%d/%d[%d]",
		k.MainGroup, k.MiddleGroup, k.SubGroup, k.ObjectId, k.PropertyId, k.PropertyIndex, k.NumElements)
}

func (k DevicePropertyAddressPlcTag) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k DevicePropertyAddressPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (k DevicePropertyAddressPlcTag) toKnxAddress() driverModel.KnxAddress {
	return driverModel.NewKnxAddress(
		k.MainGroup,
		k.MiddleGroup,
		k.SubGroup,
	)
}

type DeviceMemoryAddressPlcTag struct {
	MainGroup   uint8 // 5 Bits: Values 0-31
	MiddleGroup uint8 // 3 Bits: values 0-7
	SubGroup    uint8 // 8 Bits
	Address     uint16
	NumElements uint8
	TagType     *driverModel.KnxDatapointType
	DeviceTag
}

func NewDeviceMemoryAddressPlcTag(mainGroup uint8, middleGroup uint8, subGroup uint8, address uint16, numElements uint8, tagType *driverModel.KnxDatapointType) DeviceMemoryAddressPlcTag {
	return DeviceMemoryAddressPlcTag{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
		Address:     address,
		NumElements: numElements,
		TagType:     tagType,
	}
}

func (k DeviceMemoryAddressPlcTag) GetAddressString() string {
	return fmt.Sprintf("%d/%d/%d#%d:%s[%d]",
		k.MainGroup, k.MiddleGroup, k.SubGroup, k.Address, k.TagType.String(), k.NumElements)
}

func (k DeviceMemoryAddressPlcTag) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k DeviceMemoryAddressPlcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (k DeviceMemoryAddressPlcTag) toKnxAddress() driverModel.KnxAddress {
	individualAddress := driverModel.NewKnxAddress(
		k.MainGroup,
		k.MiddleGroup,
		k.SubGroup,
	)
	return individualAddress
}

func CastToKnxTagFromPlcTag(plcTag apiModel.PlcTag) (Tag, error) {
	if knxTag, ok := plcTag.(Tag); ok {
		return knxTag, nil
	}
	return nil, errors.New("couldn't cast to KnxNetIpTag")
}

func CastToGroupAddressTagFromPlcTag(plcTag apiModel.PlcTag) (GroupAddressTag, error) {
	if groupAddressTag, ok := plcTag.(GroupAddressTag); ok {
		return groupAddressTag, nil
	}
	return nil, errors.New("couldn't cast to KnxNetIpTag")
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
