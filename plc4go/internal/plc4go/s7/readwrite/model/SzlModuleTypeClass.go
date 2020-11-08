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
package model

import (
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

type SzlModuleTypeClass uint8

type ISzlModuleTypeClass interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    SzlModuleTypeClass_CPU SzlModuleTypeClass = 0x0
    SzlModuleTypeClass_IM SzlModuleTypeClass = 0x4
    SzlModuleTypeClass_FM SzlModuleTypeClass = 0x8
    SzlModuleTypeClass_CP SzlModuleTypeClass = 0xC
)

func SzlModuleTypeClassValueOf(value uint8) SzlModuleTypeClass {
    switch value {
        case 0x0:
            return SzlModuleTypeClass_CPU
        case 0x4:
            return SzlModuleTypeClass_IM
        case 0x8:
            return SzlModuleTypeClass_FM
        case 0xC:
            return SzlModuleTypeClass_CP
    }
    return 0
}

func CastSzlModuleTypeClass(structType interface{}) SzlModuleTypeClass {
    castFunc := func(typ interface{}) SzlModuleTypeClass {
        if sSzlModuleTypeClass, ok := typ.(SzlModuleTypeClass); ok {
            return sSzlModuleTypeClass
        }
        return 0
    }
    return castFunc(structType)
}

func (m SzlModuleTypeClass) LengthInBits() uint16 {
    return 4
}

func (m SzlModuleTypeClass) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func SzlModuleTypeClassParse(io *utils.ReadBuffer) (SzlModuleTypeClass, error) {
    // TODO: Implement ...
    return 0, nil
}

func (e SzlModuleTypeClass) Serialize(io utils.WriteBuffer) error {
    // TODO: Implement ...
    return nil
}
