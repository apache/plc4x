//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package knxnetip

import (
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"strconv"
	"strings"
)

type KnxNetIpField interface {
	IsPatternField() bool
	matches(knxGroupAddress *driverModel.KnxGroupAddress) bool
	toGroupAddress() *driverModel.KnxGroupAddress
	apiModel.PlcField
}

type KnxNetIpGroupAddress3LevelPlcField struct {
	FieldType *model.KnxDatapointType
	// 5 Bits: Values 0-31
	MainGroup string
	// 3 Bits: values 0-7
	MiddleGroup string
	// 8 Bits
	SubGroup string
	KnxNetIpField
}

func NewKnxNetIpGroupAddress3LevelPlcField(fieldType *model.KnxDatapointType, mainGroup string, middleGroup string, subGroup string) KnxNetIpGroupAddress3LevelPlcField {
	return KnxNetIpGroupAddress3LevelPlcField{
		FieldType:   fieldType,
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
	}
}

func (k KnxNetIpGroupAddress3LevelPlcField) GetTypeName() string {
	return k.FieldType.FormatName()
}

func (k KnxNetIpGroupAddress3LevelPlcField) GetQuantity() uint16 {
	return 1
}

func (k KnxNetIpGroupAddress3LevelPlcField) IsPatternField() bool {
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

func (k KnxNetIpGroupAddress3LevelPlcField) matches(knxGroupAddress *driverModel.KnxGroupAddress) bool {
	level3KnxGroupAddress := driverModel.CastKnxGroupAddress3Level(knxGroupAddress)
	if level3KnxGroupAddress == nil {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level3KnxGroupAddress.MainGroup))) &&
		matches(k.MiddleGroup, strconv.Itoa(int(level3KnxGroupAddress.MiddleGroup))) &&
		matches(k.SubGroup, strconv.Itoa(int(level3KnxGroupAddress.SubGroup)))
}

func (k KnxNetIpGroupAddress3LevelPlcField) toGroupAddress() *driverModel.KnxGroupAddress {
	mainGroup, err := strconv.Atoi(k.MainGroup)
	if err != nil {
		return nil
	}
	midleGroup, err := strconv.Atoi(k.MiddleGroup)
	if err != nil {
		return nil
	}
	subGroup, err := strconv.Atoi(k.SubGroup)
	if err != nil {
		return nil
	}
	ga := &driverModel.KnxGroupAddress{}
	l3 := &driverModel.KnxGroupAddress3Level{
		MainGroup:   uint8(mainGroup),
		MiddleGroup: uint8(midleGroup),
		SubGroup:    uint8(subGroup),
		Parent:      ga,
	}
	ga.Child = l3
	return ga
}

type KnxNetIpGroupAddress2LevelPlcField struct {
	FieldType *model.KnxDatapointType
	// 5 Bits: Values 0-31
	MainGroup string
	// 11 Bits
	SubGroup string
	KnxNetIpField
}

func NewKnxNetIpGroupAddress2LevelPlcField(fieldType *model.KnxDatapointType, mainGroup string, subGroup string) KnxNetIpGroupAddress2LevelPlcField {
	return KnxNetIpGroupAddress2LevelPlcField{
		FieldType: fieldType,
		MainGroup: mainGroup,
		SubGroup:  subGroup,
	}
}

func (k KnxNetIpGroupAddress2LevelPlcField) GetTypeName() string {
	return k.FieldType.FormatName()
}

func (k KnxNetIpGroupAddress2LevelPlcField) GetQuantity() uint16 {
	return 1
}

func (k KnxNetIpGroupAddress2LevelPlcField) IsPatternField() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		_, err = strconv.Atoi(k.SubGroup)
		if err == nil {
			return false
		}
	}
	return true
}

func (k KnxNetIpGroupAddress2LevelPlcField) matches(knxGroupAddress *driverModel.KnxGroupAddress) bool {
	level2KnxGroupAddress := driverModel.CastKnxGroupAddress2Level(knxGroupAddress)
	if level2KnxGroupAddress == nil {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level2KnxGroupAddress.MainGroup))) &&
		matches(k.SubGroup, strconv.Itoa(int(level2KnxGroupAddress.SubGroup)))
}

func (k KnxNetIpGroupAddress2LevelPlcField) toGroupAddress() *driverModel.KnxGroupAddress {
	mainGroup, err := strconv.Atoi(k.MainGroup)
	if err != nil {
		return nil
	}
	subGroup, err := strconv.Atoi(k.SubGroup)
	if err != nil {
		return nil
	}
	ga := &driverModel.KnxGroupAddress{}
	l3 := &driverModel.KnxGroupAddress2Level{
		MainGroup: uint8(mainGroup),
		SubGroup:  uint16(subGroup),
		Parent:    ga,
	}
	ga.Child = l3
	return ga
}

type KnxNetIpGroupAddress1LevelPlcField struct {
	FieldType *model.KnxDatapointType
	// 16 Bits
	MainGroup string
	KnxNetIpField
}

func NewKnxNetIpGroupAddress1LevelPlcField(fieldType *model.KnxDatapointType, mainGroup string) KnxNetIpGroupAddress1LevelPlcField {
	return KnxNetIpGroupAddress1LevelPlcField{
		FieldType: fieldType,
		MainGroup: mainGroup,
	}
}

func (k KnxNetIpGroupAddress1LevelPlcField) GetTypeName() string {
	return k.FieldType.FormatName()
}

func (k KnxNetIpGroupAddress1LevelPlcField) GetQuantity() uint16 {
	return 1
}

func (k KnxNetIpGroupAddress1LevelPlcField) IsPatternField() bool {
	_, err := strconv.Atoi(k.MainGroup)
	if err == nil {
		return false
	}
	return true
}

func (k KnxNetIpGroupAddress1LevelPlcField) matches(knxGroupAddress *driverModel.KnxGroupAddress) bool {
	level1KnxGroupAddress := driverModel.CastKnxGroupAddressFreeLevel(knxGroupAddress)
	if level1KnxGroupAddress == nil {
		return false
	}
	return matches(k.MainGroup, strconv.Itoa(int(level1KnxGroupAddress.SubGroup)))
}

func (k KnxNetIpGroupAddress1LevelPlcField) toGroupAddress() *driverModel.KnxGroupAddress {
	mainGroup, err := strconv.Atoi(k.MainGroup)
	if err != nil {
		return nil
	}
	ga := &driverModel.KnxGroupAddress{}
	l3 := &driverModel.KnxGroupAddressFreeLevel{
		SubGroup: uint16(mainGroup),
		Parent:   ga,
	}
	ga.Child = l3
	return ga
}

func CastToKnxNetIpFieldFromPlcField(plcField apiModel.PlcField) (KnxNetIpField, error) {
	if knxNetIpField, ok := plcField.(KnxNetIpField); ok {
		return knxNetIpField, nil
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
	return false
}

func GroupAddressToString(groupAddress *driverModel.KnxGroupAddress) string {
	if groupAddress != nil {
		switch groupAddress.Child.(type) {
		case *driverModel.KnxGroupAddress3Level:
			level3 := driverModel.CastKnxGroupAddress3Level(groupAddress)
			return strconv.Itoa(int(level3.MainGroup)) + "/" + strconv.Itoa(int(level3.MiddleGroup)) + "/" + strconv.Itoa(int(level3.SubGroup))
		case *driverModel.KnxGroupAddress2Level:
			level2 := driverModel.CastKnxGroupAddress2Level(groupAddress)
			return strconv.Itoa(int(level2.MainGroup)) + "/" + strconv.Itoa(int(level2.SubGroup))
		case *driverModel.KnxGroupAddressFreeLevel:
			level1 := driverModel.CastKnxGroupAddressFreeLevel(groupAddress)
			return strconv.Itoa(int(level1.SubGroup))
		}
	}
	return ""
}
