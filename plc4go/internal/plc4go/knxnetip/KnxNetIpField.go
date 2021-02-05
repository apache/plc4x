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
    driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
    apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "strconv"
    "strings"
)

type KnxNetIpField interface {
    apiModel.PlcField
}

type KnxNetIpGroupAddressField interface {
    GetFieldType() *driverModel.KnxDatapointType
    IsPatternField() bool
    matches(knxGroupAddress *driverModel.KnxGroupAddress) bool
    toGroupAddress() *driverModel.KnxGroupAddress
    KnxNetIpField
}

type KnxNetIpDeviceField interface {
    toKnxAddress() *driverModel.KnxAddress
    KnxNetIpField
}

type KnxNetIpGroupAddress3LevelPlcField struct {
    MainGroup   string // 5 Bits: Values 0-31
    MiddleGroup string // 3 Bits: values 0-7
    SubGroup    string // 8 Bits
    FieldType   *driverModel.KnxDatapointType
}

func NewKnxNetIpGroupAddress3LevelPlcField(mainGroup string, middleGroup string, subGroup string, fieldType *driverModel.KnxDatapointType) KnxNetIpGroupAddress3LevelPlcField {
    return KnxNetIpGroupAddress3LevelPlcField{
        MainGroup:   mainGroup,
        MiddleGroup: middleGroup,
        SubGroup:    subGroup,
        FieldType:   fieldType,
    }
}

func (k KnxNetIpGroupAddress3LevelPlcField) GetTypeName() string {
    return k.FieldType.Name()
}

func (k KnxNetIpGroupAddress3LevelPlcField) GetFieldType() *driverModel.KnxDatapointType {
    return k.FieldType
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
    middleGroup, err := strconv.Atoi(k.MiddleGroup)
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
        MiddleGroup: uint8(middleGroup),
        SubGroup:    uint8(subGroup),
        Parent:      ga,
    }
    ga.Child = l3
    return ga
}

type KnxNetIpGroupAddress2LevelPlcField struct {
    MainGroup string // 5 Bits: Values 0-31
    SubGroup  string // 11 Bits
    FieldType *driverModel.KnxDatapointType
    KnxNetIpField
}

func NewKnxNetIpGroupAddress2LevelPlcField(mainGroup string, subGroup string, fieldType *driverModel.KnxDatapointType) KnxNetIpGroupAddress2LevelPlcField {
    return KnxNetIpGroupAddress2LevelPlcField{
        MainGroup: mainGroup,
        SubGroup:  subGroup,
        FieldType: fieldType,
    }
}

func (k KnxNetIpGroupAddress2LevelPlcField) GetTypeName() string {
    return k.FieldType.Name()
}

func (k KnxNetIpGroupAddress2LevelPlcField) GetFieldType() *driverModel.KnxDatapointType {
    return k.FieldType
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
    MainGroup string // 16 Bits
    FieldType *driverModel.KnxDatapointType
    KnxNetIpField
}

func NewKnxNetIpGroupAddress1LevelPlcField(mainGroup string, fieldType *driverModel.KnxDatapointType) KnxNetIpGroupAddress1LevelPlcField {
    return KnxNetIpGroupAddress1LevelPlcField{
        MainGroup: mainGroup,
        FieldType: fieldType,
    }
}

func (k KnxNetIpGroupAddress1LevelPlcField) GetTypeName() string {
    return k.FieldType.Name()
}

func (k KnxNetIpGroupAddress1LevelPlcField) GetFieldType() *driverModel.KnxDatapointType {
    return k.FieldType
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

type KnxNetIpDeviceQueryField struct {
    MainGroup   string // 5 Bits: Values 0-31
    MiddleGroup string // 3 Bits: values 0-7
    SubGroup    string // 8 Bits
    KnxNetIpDeviceField
}

func NewKnxNetIpDeviceQueryField(mainGroup string, middleGroup string, subGroup string) KnxNetIpDeviceQueryField {
    return KnxNetIpDeviceQueryField{
        MainGroup:   mainGroup,
        MiddleGroup: middleGroup,
        SubGroup:    subGroup,
    }
}

func (k KnxNetIpDeviceQueryField) GetTypeName() string {
    return ""
}

func (k KnxNetIpDeviceQueryField) GetQuantity() uint16 {
    return 1
}

func (k KnxNetIpDeviceQueryField) toKnxAddress() *driverModel.KnxAddress {
    return nil
}

type KnxNetIpDevicePropertyAddressPlcField struct {
    MainGroup   uint8 // 5 Bits: Values 0-31
    MiddleGroup uint8 // 3 Bits: values 0-7
    SubGroup    uint8 // 8 Bits
    ObjectId    uint8
    PropertyId  uint8
    KnxNetIpDeviceField
}

func NewKnxNetIpDevicePropertyAddressPlcField(mainGroup uint8, middleGroup uint8, subGroup uint8, objectId uint8, propertyId uint8) KnxNetIpDevicePropertyAddressPlcField {
    return KnxNetIpDevicePropertyAddressPlcField{
        MainGroup:   mainGroup,
        MiddleGroup: middleGroup,
        SubGroup:    subGroup,
        ObjectId:    objectId,
        PropertyId:  propertyId,
    }
}

func (k KnxNetIpDevicePropertyAddressPlcField) GetTypeName() string {
    return ""
}

func (k KnxNetIpDevicePropertyAddressPlcField) GetQuantity() uint16 {
    return 1
}

func (k KnxNetIpDevicePropertyAddressPlcField) toKnxAddress() *driverModel.KnxAddress {
    ga := &driverModel.KnxAddress{
        MainGroup:   k.MainGroup,
        MiddleGroup: k.MiddleGroup,
        SubGroup:    k.SubGroup,
    }
    return ga
}

type KnxNetIpDeviceMemoryAddressPlcField struct {
    MainGroup        uint8 // 5 Bits: Values 0-31
    MiddleGroup      uint8 // 3 Bits: values 0-7
    SubGroup         uint8 // 8 Bits
    Address          uint16
    NumberOfElements uint8
    FieldType        *driverModel.KnxDatapointType
    KnxNetIpDeviceField
}

func NewKnxNetIpDeviceMemoryAddressPlcField(mainGroup uint8, middleGroup uint8, subGroup uint8, address uint16, numberOfElements uint8, fieldType *driverModel.KnxDatapointType) KnxNetIpDeviceMemoryAddressPlcField {
    return KnxNetIpDeviceMemoryAddressPlcField{
        MainGroup:        mainGroup,
        MiddleGroup:      middleGroup,
        SubGroup:         subGroup,
        Address:          address,
        NumberOfElements: numberOfElements,
        FieldType:        fieldType,
    }
}

func (k KnxNetIpDeviceMemoryAddressPlcField) GetTypeName() string {
    if k.FieldType != nil {
        return k.FieldType.Name()
    }
    return ""
}

func (k KnxNetIpDeviceMemoryAddressPlcField) GetFieldType() *driverModel.KnxDatapointType {
    return k.FieldType
}

func (k KnxNetIpDeviceMemoryAddressPlcField) GetQuantity() uint16 {
    return uint16(k.NumberOfElements)
}

func (k KnxNetIpDeviceMemoryAddressPlcField) toKnxAddress() *driverModel.KnxAddress {
    individualAddress := &driverModel.KnxAddress{
        MainGroup:   k.MainGroup,
        MiddleGroup: k.MiddleGroup,
        SubGroup:    k.SubGroup,
    }
    return individualAddress
}

type KnxNetIpCommunicationObjectQueryField struct {
    MainGroup   uint8 // 5 Bits: Values 0-31
    MiddleGroup uint8 // 3 Bits: values 0-7
    SubGroup    uint8 // 8 Bits
    KnxNetIpDeviceField
}

func NewKnxNetIpCommunicationObjectQueryField(mainGroup uint8, middleGroup uint8, subGroup uint8) KnxNetIpCommunicationObjectQueryField {
    return KnxNetIpCommunicationObjectQueryField{
        MainGroup:   mainGroup,
        MiddleGroup: middleGroup,
        SubGroup:    subGroup,
    }
}

func (k KnxNetIpCommunicationObjectQueryField) GetTypeName() string {
    return ""
}

func (k KnxNetIpCommunicationObjectQueryField) GetQuantity() uint16 {
    return 1
}

func (k KnxNetIpCommunicationObjectQueryField) toKnxAddress() *driverModel.KnxAddress {
    individualAddress := &driverModel.KnxAddress{
        MainGroup:   k.MainGroup,
        MiddleGroup: k.MiddleGroup,
        SubGroup:    k.SubGroup,
    }
    return individualAddress
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
}
