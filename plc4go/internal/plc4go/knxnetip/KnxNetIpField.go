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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/knxnetip/readwrite/model"
    model2 "plc4x.apache.org/plc4go/v0/pkg/plc4go/model"
)

type KnxNetIpField interface {
     model2.PlcField
}

type KnxNetIpGroupAddress3LevelPlcField struct {
    FieldType   *model.KnxDatapointType
    // 5 Bits: Values 0-31
    MainGroup   string
    // 3 Bits: values 0-7
    MiddleGroup string
    // 8 Bits
    SubGroup    string
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

type KnxNetIpGroupAddress2LevelPlcField struct {
    FieldType   *model.KnxDatapointType
    // 5 Bits: Values 0-31
    MainGroup   string
    // 11 Bits
    SubGroup    string
    KnxNetIpField
}

func NewKnxNetIpGroupAddress2LevelPlcField(fieldType *model.KnxDatapointType, mainGroup string, subGroup string) KnxNetIpGroupAddress2LevelPlcField {
    return KnxNetIpGroupAddress2LevelPlcField{
        FieldType:   fieldType,
        MainGroup:   mainGroup,
        SubGroup:    subGroup,
    }
}

func (k KnxNetIpGroupAddress2LevelPlcField) GetTypeName() string {
    return k.FieldType.FormatName()
}

func (k KnxNetIpGroupAddress2LevelPlcField) GetQuantity() uint16 {
    return 1
}

type KnxNetIpGroupAddress1LevelPlcField struct {
    FieldType   *model.KnxDatapointType
    // 16 Bits
    MainGroup   string
    KnxNetIpField
}

func NewKnxNetIpGroupAddress1LevelPlcField(fieldType *model.KnxDatapointType, mainGroup string) KnxNetIpGroupAddress1LevelPlcField {
    return KnxNetIpGroupAddress1LevelPlcField{
        FieldType:   fieldType,
        MainGroup:   mainGroup,
    }
}

func (k KnxNetIpGroupAddress1LevelPlcField) GetTypeName() string {
    return k.FieldType.FormatName()
}

func (k KnxNetIpGroupAddress1LevelPlcField) GetQuantity() uint16 {
    return 1
}

func CastToKnxNetIpFieldFromPlcField(plcField model2.PlcField) (KnxNetIpField, error) {
    if knxNetIpField, ok := plcField.(KnxNetIpField); ok {
        return knxNetIpField, nil
    }
    return nil, errors.New("couldn't cast to KnxNetIpField")
}
