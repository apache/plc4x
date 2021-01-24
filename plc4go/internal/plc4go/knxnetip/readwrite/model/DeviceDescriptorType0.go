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

type DeviceDescriptorType0 uint16

type IDeviceDescriptorType0 interface {
    FirmwareType() FirmwareType
    MediumType() DeviceDescriptorMediumType
    Serialize(io utils.WriteBuffer) error
}

const(
    DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_0 DeviceDescriptorType0 = 0x0010
    DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_1 DeviceDescriptorType0 = 0x0011
    DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_2 DeviceDescriptorType0 = 0x0012
    DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_3 DeviceDescriptorType0 = 0x0013
    DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_0 DeviceDescriptorType0 = 0x0020
    DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_1 DeviceDescriptorType0 = 0x0021
    DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_5 DeviceDescriptorType0 = 0x0025
    DeviceDescriptorType0_TP1_SYSTEM_300 DeviceDescriptorType0 = 0x0300
    DeviceDescriptorType0_TP1_BIM_M112_0 DeviceDescriptorType0 = 0x0700
    DeviceDescriptorType0_TP1_BIM_M112_1 DeviceDescriptorType0 = 0x0701
    DeviceDescriptorType0_TP1_BIM_M112_5 DeviceDescriptorType0 = 0x0705
    DeviceDescriptorType0_TP1_SYSTEM_B DeviceDescriptorType0 = 0x07B0
    DeviceDescriptorType0_TP1_IR_DECODER_0 DeviceDescriptorType0 = 0x0810
    DeviceDescriptorType0_TP1_IR_DECODER_1 DeviceDescriptorType0 = 0x0811
    DeviceDescriptorType0_TP1_COUPLER_0 DeviceDescriptorType0 = 0x0910
    DeviceDescriptorType0_TP1_COUPLER_1 DeviceDescriptorType0 = 0x0911
    DeviceDescriptorType0_TP1_COUPLER_2 DeviceDescriptorType0 = 0x0912
    DeviceDescriptorType0_TP1_KNXNETIP_ROUTER DeviceDescriptorType0 = 0x091A
    DeviceDescriptorType0_TP1_NONE_D DeviceDescriptorType0 = 0x0AFD
    DeviceDescriptorType0_TP1_NONE_E DeviceDescriptorType0 = 0x0AFE
    DeviceDescriptorType0_PL110_BCU_1_2 DeviceDescriptorType0 = 0x1012
    DeviceDescriptorType0_PL110_BCU_1_3 DeviceDescriptorType0 = 0x1013
    DeviceDescriptorType0_PL110_SYSTEM_B DeviceDescriptorType0 = 0x17B0
    DeviceDescriptorType0_PL110_MEDIA_COUPLER_PL_TP DeviceDescriptorType0 = 0x1900
    DeviceDescriptorType0_RF_BI_DIRECTIONAL_DEVICES DeviceDescriptorType0 = 0x2010
    DeviceDescriptorType0_RF_UNI_DIRECTIONAL_DEVICES DeviceDescriptorType0 = 0x2110
    DeviceDescriptorType0_TP0_BCU_1 DeviceDescriptorType0 = 0x3012
    DeviceDescriptorType0_PL132_BCU_1 DeviceDescriptorType0 = 0x4012
    DeviceDescriptorType0_KNX_IP_SYSTEM7 DeviceDescriptorType0 = 0x5705
)


func (e DeviceDescriptorType0) FirmwareType() FirmwareType {
    switch e  {
        case 0x0010: { /* '0x0010' */
            return FirmwareType_BCU_1_SYSTEM_1
        }
        case 0x0011: { /* '0x0011' */
            return FirmwareType_BCU_1_SYSTEM_1
        }
        case 0x0012: { /* '0x0012' */
            return FirmwareType_BCU_1_SYSTEM_1
        }
        case 0x0013: { /* '0x0013' */
            return FirmwareType_BCU_1_SYSTEM_1
        }
        case 0x0020: { /* '0x0020' */
            return FirmwareType_BCU_2_SYSTEM_2
        }
        case 0x0021: { /* '0x0021' */
            return FirmwareType_BCU_2_SYSTEM_2
        }
        case 0x0025: { /* '0x0025' */
            return FirmwareType_BCU_2_SYSTEM_2
        }
        case 0x0300: { /* '0x0300' */
            return FirmwareType_SYSTEM_300
        }
        case 0x0700: { /* '0x0700' */
            return FirmwareType_BIM_M112
        }
        case 0x0701: { /* '0x0701' */
            return FirmwareType_BIM_M112
        }
        case 0x0705: { /* '0x0705' */
            return FirmwareType_BIM_M112
        }
        case 0x07B0: { /* '0x07B0' */
            return FirmwareType_SYSTEM_B
        }
        case 0x0810: { /* '0x0810' */
            return FirmwareType_IR_DECODER
        }
        case 0x0811: { /* '0x0811' */
            return FirmwareType_IR_DECODER
        }
        case 0x0910: { /* '0x0910' */
            return FirmwareType_COUPLER
        }
        case 0x0911: { /* '0x0911' */
            return FirmwareType_COUPLER
        }
        case 0x0912: { /* '0x0912' */
            return FirmwareType_COUPLER
        }
        case 0x091A: { /* '0x091A' */
            return FirmwareType_COUPLER
        }
        case 0x0AFD: { /* '0x0AFD' */
            return FirmwareType_NONE
        }
        case 0x0AFE: { /* '0x0AFE' */
            return FirmwareType_NONE
        }
        case 0x1012: { /* '0x1012' */
            return FirmwareType_BCU_1_SYSTEM_1
        }
        case 0x1013: { /* '0x1013' */
            return FirmwareType_BCU_1_SYSTEM_1
        }
        case 0x17B0: { /* '0x17B0' */
            return FirmwareType_SYSTEM_B
        }
        case 0x1900: { /* '0x1900' */
            return FirmwareType_MEDIA_COUPLER_PL_TP
        }
        case 0x2010: { /* '0x2010' */
            return FirmwareType_RF_BI_DIRECTIONAL_DEVICES
        }
        case 0x2110: { /* '0x2110' */
            return FirmwareType_RF_UNI_DIRECTIONAL_DEVICES
        }
        case 0x3012: { /* '0x3012' */
            return FirmwareType_BCU_1
        }
        case 0x4012: { /* '0x4012' */
            return FirmwareType_BCU_1
        }
        case 0x5705: { /* '0x5705' */
            return FirmwareType_SYSTEM_7
        }
        default: {
            return 0
        }
    }
}

func (e DeviceDescriptorType0) MediumType() DeviceDescriptorMediumType {
    switch e  {
        case 0x0010: { /* '0x0010' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0011: { /* '0x0011' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0012: { /* '0x0012' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0013: { /* '0x0013' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0020: { /* '0x0020' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0021: { /* '0x0021' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0025: { /* '0x0025' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0300: { /* '0x0300' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0700: { /* '0x0700' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0701: { /* '0x0701' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0705: { /* '0x0705' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x07B0: { /* '0x07B0' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0810: { /* '0x0810' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0811: { /* '0x0811' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0910: { /* '0x0910' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0911: { /* '0x0911' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0912: { /* '0x0912' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x091A: { /* '0x091A' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0AFD: { /* '0x0AFD' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x0AFE: { /* '0x0AFE' */
            return DeviceDescriptorMediumType_TP1
        }
        case 0x1012: { /* '0x1012' */
            return DeviceDescriptorMediumType_PL110
        }
        case 0x1013: { /* '0x1013' */
            return DeviceDescriptorMediumType_PL110
        }
        case 0x17B0: { /* '0x17B0' */
            return DeviceDescriptorMediumType_PL110
        }
        case 0x1900: { /* '0x1900' */
            return DeviceDescriptorMediumType_PL110
        }
        case 0x2010: { /* '0x2010' */
            return DeviceDescriptorMediumType_RF
        }
        case 0x2110: { /* '0x2110' */
            return DeviceDescriptorMediumType_RF
        }
        case 0x3012: { /* '0x3012' */
            return DeviceDescriptorMediumType_TP0
        }
        case 0x4012: { /* '0x4012' */
            return DeviceDescriptorMediumType_PL132
        }
        case 0x5705: { /* '0x5705' */
            return DeviceDescriptorMediumType_KNX_IP
        }
        default: {
            return 0
        }
    }
}
func DeviceDescriptorType0ByValue(value uint16) DeviceDescriptorType0 {
    switch value {
        case 0x0010:
            return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_0
        case 0x0011:
            return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_1
        case 0x0012:
            return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_2
        case 0x0013:
            return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_3
        case 0x0020:
            return DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_0
        case 0x0021:
            return DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_1
        case 0x0025:
            return DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_5
        case 0x0300:
            return DeviceDescriptorType0_TP1_SYSTEM_300
        case 0x0700:
            return DeviceDescriptorType0_TP1_BIM_M112_0
        case 0x0701:
            return DeviceDescriptorType0_TP1_BIM_M112_1
        case 0x0705:
            return DeviceDescriptorType0_TP1_BIM_M112_5
        case 0x07B0:
            return DeviceDescriptorType0_TP1_SYSTEM_B
        case 0x0810:
            return DeviceDescriptorType0_TP1_IR_DECODER_0
        case 0x0811:
            return DeviceDescriptorType0_TP1_IR_DECODER_1
        case 0x0910:
            return DeviceDescriptorType0_TP1_COUPLER_0
        case 0x0911:
            return DeviceDescriptorType0_TP1_COUPLER_1
        case 0x0912:
            return DeviceDescriptorType0_TP1_COUPLER_2
        case 0x091A:
            return DeviceDescriptorType0_TP1_KNXNETIP_ROUTER
        case 0x0AFD:
            return DeviceDescriptorType0_TP1_NONE_D
        case 0x0AFE:
            return DeviceDescriptorType0_TP1_NONE_E
        case 0x1012:
            return DeviceDescriptorType0_PL110_BCU_1_2
        case 0x1013:
            return DeviceDescriptorType0_PL110_BCU_1_3
        case 0x17B0:
            return DeviceDescriptorType0_PL110_SYSTEM_B
        case 0x1900:
            return DeviceDescriptorType0_PL110_MEDIA_COUPLER_PL_TP
        case 0x2010:
            return DeviceDescriptorType0_RF_BI_DIRECTIONAL_DEVICES
        case 0x2110:
            return DeviceDescriptorType0_RF_UNI_DIRECTIONAL_DEVICES
        case 0x3012:
            return DeviceDescriptorType0_TP0_BCU_1
        case 0x4012:
            return DeviceDescriptorType0_PL132_BCU_1
        case 0x5705:
            return DeviceDescriptorType0_KNX_IP_SYSTEM7
    }
    return 0
}

func DeviceDescriptorType0ByName(value string) DeviceDescriptorType0 {
    switch value {
    case "TP1_BCU_1_SYSTEM_1_0":
        return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_0
    case "TP1_BCU_1_SYSTEM_1_1":
        return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_1
    case "TP1_BCU_1_SYSTEM_1_2":
        return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_2
    case "TP1_BCU_1_SYSTEM_1_3":
        return DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_3
    case "TP1_BCU_2_SYSTEM_2_0":
        return DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_0
    case "TP1_BCU_2_SYSTEM_2_1":
        return DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_1
    case "TP1_BCU_2_SYSTEM_2_5":
        return DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_5
    case "TP1_SYSTEM_300":
        return DeviceDescriptorType0_TP1_SYSTEM_300
    case "TP1_BIM_M112_0":
        return DeviceDescriptorType0_TP1_BIM_M112_0
    case "TP1_BIM_M112_1":
        return DeviceDescriptorType0_TP1_BIM_M112_1
    case "TP1_BIM_M112_5":
        return DeviceDescriptorType0_TP1_BIM_M112_5
    case "TP1_SYSTEM_B":
        return DeviceDescriptorType0_TP1_SYSTEM_B
    case "TP1_IR_DECODER_0":
        return DeviceDescriptorType0_TP1_IR_DECODER_0
    case "TP1_IR_DECODER_1":
        return DeviceDescriptorType0_TP1_IR_DECODER_1
    case "TP1_COUPLER_0":
        return DeviceDescriptorType0_TP1_COUPLER_0
    case "TP1_COUPLER_1":
        return DeviceDescriptorType0_TP1_COUPLER_1
    case "TP1_COUPLER_2":
        return DeviceDescriptorType0_TP1_COUPLER_2
    case "TP1_KNXNETIP_ROUTER":
        return DeviceDescriptorType0_TP1_KNXNETIP_ROUTER
    case "TP1_NONE_D":
        return DeviceDescriptorType0_TP1_NONE_D
    case "TP1_NONE_E":
        return DeviceDescriptorType0_TP1_NONE_E
    case "PL110_BCU_1_2":
        return DeviceDescriptorType0_PL110_BCU_1_2
    case "PL110_BCU_1_3":
        return DeviceDescriptorType0_PL110_BCU_1_3
    case "PL110_SYSTEM_B":
        return DeviceDescriptorType0_PL110_SYSTEM_B
    case "PL110_MEDIA_COUPLER_PL_TP":
        return DeviceDescriptorType0_PL110_MEDIA_COUPLER_PL_TP
    case "RF_BI_DIRECTIONAL_DEVICES":
        return DeviceDescriptorType0_RF_BI_DIRECTIONAL_DEVICES
    case "RF_UNI_DIRECTIONAL_DEVICES":
        return DeviceDescriptorType0_RF_UNI_DIRECTIONAL_DEVICES
    case "TP0_BCU_1":
        return DeviceDescriptorType0_TP0_BCU_1
    case "PL132_BCU_1":
        return DeviceDescriptorType0_PL132_BCU_1
    case "KNX_IP_SYSTEM7":
        return DeviceDescriptorType0_KNX_IP_SYSTEM7
    }
    return 0
}

func CastDeviceDescriptorType0(structType interface{}) DeviceDescriptorType0 {
    castFunc := func(typ interface{}) DeviceDescriptorType0 {
        if sDeviceDescriptorType0, ok := typ.(DeviceDescriptorType0); ok {
            return sDeviceDescriptorType0
        }
        return 0
    }
    return castFunc(structType)
}

func (m DeviceDescriptorType0) LengthInBits() uint16 {
    return 16
}

func (m DeviceDescriptorType0) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DeviceDescriptorType0Parse(io *utils.ReadBuffer) (DeviceDescriptorType0, error) {
    val, err := io.ReadUint16(16)
    if err != nil {
        return 0, nil
    }
    return DeviceDescriptorType0ByValue(val), nil
}

func (e DeviceDescriptorType0) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint16(16, uint16(e))
    return err
}

func (e DeviceDescriptorType0) String() string {
    switch e {
    case DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_0:
        return "TP1_BCU_1_SYSTEM_1_0"
    case DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_1:
        return "TP1_BCU_1_SYSTEM_1_1"
    case DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_2:
        return "TP1_BCU_1_SYSTEM_1_2"
    case DeviceDescriptorType0_TP1_BCU_1_SYSTEM_1_3:
        return "TP1_BCU_1_SYSTEM_1_3"
    case DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_0:
        return "TP1_BCU_2_SYSTEM_2_0"
    case DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_1:
        return "TP1_BCU_2_SYSTEM_2_1"
    case DeviceDescriptorType0_TP1_BCU_2_SYSTEM_2_5:
        return "TP1_BCU_2_SYSTEM_2_5"
    case DeviceDescriptorType0_TP1_SYSTEM_300:
        return "TP1_SYSTEM_300"
    case DeviceDescriptorType0_TP1_BIM_M112_0:
        return "TP1_BIM_M112_0"
    case DeviceDescriptorType0_TP1_BIM_M112_1:
        return "TP1_BIM_M112_1"
    case DeviceDescriptorType0_TP1_BIM_M112_5:
        return "TP1_BIM_M112_5"
    case DeviceDescriptorType0_TP1_SYSTEM_B:
        return "TP1_SYSTEM_B"
    case DeviceDescriptorType0_TP1_IR_DECODER_0:
        return "TP1_IR_DECODER_0"
    case DeviceDescriptorType0_TP1_IR_DECODER_1:
        return "TP1_IR_DECODER_1"
    case DeviceDescriptorType0_TP1_COUPLER_0:
        return "TP1_COUPLER_0"
    case DeviceDescriptorType0_TP1_COUPLER_1:
        return "TP1_COUPLER_1"
    case DeviceDescriptorType0_TP1_COUPLER_2:
        return "TP1_COUPLER_2"
    case DeviceDescriptorType0_TP1_KNXNETIP_ROUTER:
        return "TP1_KNXNETIP_ROUTER"
    case DeviceDescriptorType0_TP1_NONE_D:
        return "TP1_NONE_D"
    case DeviceDescriptorType0_TP1_NONE_E:
        return "TP1_NONE_E"
    case DeviceDescriptorType0_PL110_BCU_1_2:
        return "PL110_BCU_1_2"
    case DeviceDescriptorType0_PL110_BCU_1_3:
        return "PL110_BCU_1_3"
    case DeviceDescriptorType0_PL110_SYSTEM_B:
        return "PL110_SYSTEM_B"
    case DeviceDescriptorType0_PL110_MEDIA_COUPLER_PL_TP:
        return "PL110_MEDIA_COUPLER_PL_TP"
    case DeviceDescriptorType0_RF_BI_DIRECTIONAL_DEVICES:
        return "RF_BI_DIRECTIONAL_DEVICES"
    case DeviceDescriptorType0_RF_UNI_DIRECTIONAL_DEVICES:
        return "RF_UNI_DIRECTIONAL_DEVICES"
    case DeviceDescriptorType0_TP0_BCU_1:
        return "TP0_BCU_1"
    case DeviceDescriptorType0_PL132_BCU_1:
        return "PL132_BCU_1"
    case DeviceDescriptorType0_KNX_IP_SYSTEM7:
        return "KNX_IP_SYSTEM7"
    }
    return ""
}
