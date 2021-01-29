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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

type FirmwareType uint8

type IFirmwareType interface {
    Code() uint8
    Serialize(io utils.WriteBuffer) error
}

const(
    FirmwareType_NONE FirmwareType = 0x1
    FirmwareType_BCU_1 FirmwareType = 0x2
    FirmwareType_BCU_1_SYSTEM_1 FirmwareType = 0x3
    FirmwareType_BCU_2_SYSTEM_2 FirmwareType = 0x4
    FirmwareType_BIM_M112 FirmwareType = 0x5
    FirmwareType_SYSTEM_B FirmwareType = 0x6
    FirmwareType_IR_DECODER FirmwareType = 0x7
    FirmwareType_MEDIA_COUPLER_PL_TP FirmwareType = 0x8
    FirmwareType_COUPLER FirmwareType = 0x9
    FirmwareType_RF_BI_DIRECTIONAL_DEVICES FirmwareType = 0xA
    FirmwareType_RF_UNI_DIRECTIONAL_DEVICES FirmwareType = 0xB
    FirmwareType_SYSTEM_300 FirmwareType = 0xC
    FirmwareType_SYSTEM_7 FirmwareType = 0xD
)


func (e FirmwareType) Code() uint8 {
    switch e  {
        case 0x1: { /* '0x1' */
            return 0xAF
        }
        case 0x2: { /* '0x2' */
            return 0x00
        }
        case 0x3: { /* '0x3' */
            return 0x01
        }
        case 0x4: { /* '0x4' */
            return 0x02
        }
        case 0x5: { /* '0x5' */
            return 0x70
        }
        case 0x6: { /* '0x6' */
            return 0x7B
        }
        case 0x7: { /* '0x7' */
            return 0x81
        }
        case 0x8: { /* '0x8' */
            return 0x90
        }
        case 0x9: { /* '0x9' */
            return 0x91
        }
        case 0xA: { /* '0xA' */
            return 0x01
        }
        case 0xB: { /* '0xB' */
            return 0x11
        }
        case 0xC: { /* '0xC' */
            return 0x30
        }
        case 0xD: { /* '0xD' */
            return 0x70
        }
        default: {
            return 0
        }
    }
}
func FirmwareTypeByValue(value uint8) FirmwareType {
    switch value {
        case 0x1:
            return FirmwareType_NONE
        case 0x2:
            return FirmwareType_BCU_1
        case 0x3:
            return FirmwareType_BCU_1_SYSTEM_1
        case 0x4:
            return FirmwareType_BCU_2_SYSTEM_2
        case 0x5:
            return FirmwareType_BIM_M112
        case 0x6:
            return FirmwareType_SYSTEM_B
        case 0x7:
            return FirmwareType_IR_DECODER
        case 0x8:
            return FirmwareType_MEDIA_COUPLER_PL_TP
        case 0x9:
            return FirmwareType_COUPLER
        case 0xA:
            return FirmwareType_RF_BI_DIRECTIONAL_DEVICES
        case 0xB:
            return FirmwareType_RF_UNI_DIRECTIONAL_DEVICES
        case 0xC:
            return FirmwareType_SYSTEM_300
        case 0xD:
            return FirmwareType_SYSTEM_7
    }
    return 0
}

func FirmwareTypeByName(value string) FirmwareType {
    switch value {
    case "NONE":
        return FirmwareType_NONE
    case "BCU_1":
        return FirmwareType_BCU_1
    case "BCU_1_SYSTEM_1":
        return FirmwareType_BCU_1_SYSTEM_1
    case "BCU_2_SYSTEM_2":
        return FirmwareType_BCU_2_SYSTEM_2
    case "BIM_M112":
        return FirmwareType_BIM_M112
    case "SYSTEM_B":
        return FirmwareType_SYSTEM_B
    case "IR_DECODER":
        return FirmwareType_IR_DECODER
    case "MEDIA_COUPLER_PL_TP":
        return FirmwareType_MEDIA_COUPLER_PL_TP
    case "COUPLER":
        return FirmwareType_COUPLER
    case "RF_BI_DIRECTIONAL_DEVICES":
        return FirmwareType_RF_BI_DIRECTIONAL_DEVICES
    case "RF_UNI_DIRECTIONAL_DEVICES":
        return FirmwareType_RF_UNI_DIRECTIONAL_DEVICES
    case "SYSTEM_300":
        return FirmwareType_SYSTEM_300
    case "SYSTEM_7":
        return FirmwareType_SYSTEM_7
    }
    return 0
}

func CastFirmwareType(structType interface{}) FirmwareType {
    castFunc := func(typ interface{}) FirmwareType {
        if sFirmwareType, ok := typ.(FirmwareType); ok {
            return sFirmwareType
        }
        return 0
    }
    return castFunc(structType)
}

func (m FirmwareType) LengthInBits() uint16 {
    return 4
}

func (m FirmwareType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func FirmwareTypeParse(io *utils.ReadBuffer) (FirmwareType, error) {
    val, err := io.ReadUint8(4)
    if err != nil {
        return 0, nil
    }
    return FirmwareTypeByValue(val), nil
}

func (e FirmwareType) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(4, uint8(e))
    return err
}

func (e FirmwareType) String() string {
    switch e {
    case FirmwareType_NONE:
        return "NONE"
    case FirmwareType_BCU_1:
        return "BCU_1"
    case FirmwareType_BCU_1_SYSTEM_1:
        return "BCU_1_SYSTEM_1"
    case FirmwareType_BCU_2_SYSTEM_2:
        return "BCU_2_SYSTEM_2"
    case FirmwareType_BIM_M112:
        return "BIM_M112"
    case FirmwareType_SYSTEM_B:
        return "SYSTEM_B"
    case FirmwareType_IR_DECODER:
        return "IR_DECODER"
    case FirmwareType_MEDIA_COUPLER_PL_TP:
        return "MEDIA_COUPLER_PL_TP"
    case FirmwareType_COUPLER:
        return "COUPLER"
    case FirmwareType_RF_BI_DIRECTIONAL_DEVICES:
        return "RF_BI_DIRECTIONAL_DEVICES"
    case FirmwareType_RF_UNI_DIRECTIONAL_DEVICES:
        return "RF_UNI_DIRECTIONAL_DEVICES"
    case FirmwareType_SYSTEM_300:
        return "SYSTEM_300"
    case FirmwareType_SYSTEM_7:
        return "SYSTEM_7"
    }
    return ""
}
