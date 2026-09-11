//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public enum DeviceDescriptor : ushort
    {
        TP1_BCU_1_SYSTEM_1_0 = 0x0010,
        TP1_BCU_1_SYSTEM_1_1 = 0x0011,
        TP1_BCU_1_SYSTEM_1_2 = 0x0012,
        TP1_BCU_1_SYSTEM_1_3 = 0x0013,
        TP1_BCU_2_SYSTEM_2_0 = 0x0020,
        TP1_BCU_2_SYSTEM_2_1 = 0x0021,
        TP1_BCU_2_SYSTEM_2_5 = 0x0025,
        TP1_SYSTEM_300 = 0x0300,
        TP1_BIM_M112_0 = 0x0700,
        TP1_BIM_M112_1 = 0x0701,
        TP1_BIM_M112_5 = 0x0705,
        TP1_SYSTEM_B = 0x07B0,
        TP1_IR_DECODER_0 = 0x0810,
        TP1_IR_DECODER_1 = 0x0811,
        TP1_COUPLER_0 = 0x0910,
        TP1_COUPLER_1 = 0x0911,
        TP1_COUPLER_2 = 0x0912,
        TP1_KNXNETIP_ROUTER = 0x091A,
        TP1_NONE_D = 0x0AFD,
        TP1_NONE_E = 0x0AFE,
        PL110_BCU_1_2 = 0x1012,
        PL110_BCU_1_3 = 0x1013,
        PL110_SYSTEM_B = 0x17B0,
        PL110_MEDIA_COUPLER_PL_TP = 0x1900,
        RF_BI_DIRECTIONAL_DEVICES = 0x2010,
        RF_UNI_DIRECTIONAL_DEVICES = 0x2110,
        TP0_BCU_1 = 0x3012,
        PL132_BCU_1 = 0x4012,
        KNX_IP_SYSTEM7 = 0x5705,
    }

    public static class DeviceDescriptorExtensions
    {
        public static DeviceDescriptorMediumType GetMediumType(this DeviceDescriptor value) => value switch
        {
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_0 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_1 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_2 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_3 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BCU_2_SYSTEM_2_0 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BCU_2_SYSTEM_2_1 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BCU_2_SYSTEM_2_5 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_SYSTEM_300 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BIM_M112_0 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BIM_M112_1 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_BIM_M112_5 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_SYSTEM_B => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_IR_DECODER_0 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_IR_DECODER_1 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_COUPLER_0 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_COUPLER_1 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_COUPLER_2 => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_KNXNETIP_ROUTER => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_NONE_D => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.TP1_NONE_E => DeviceDescriptorMediumType.TP1,
            DeviceDescriptor.PL110_BCU_1_2 => DeviceDescriptorMediumType.PL110,
            DeviceDescriptor.PL110_BCU_1_3 => DeviceDescriptorMediumType.PL110,
            DeviceDescriptor.PL110_SYSTEM_B => DeviceDescriptorMediumType.PL110,
            DeviceDescriptor.PL110_MEDIA_COUPLER_PL_TP => DeviceDescriptorMediumType.PL110,
            DeviceDescriptor.RF_BI_DIRECTIONAL_DEVICES => DeviceDescriptorMediumType.RF,
            DeviceDescriptor.RF_UNI_DIRECTIONAL_DEVICES => DeviceDescriptorMediumType.RF,
            DeviceDescriptor.TP0_BCU_1 => DeviceDescriptorMediumType.TP0,
            DeviceDescriptor.PL132_BCU_1 => DeviceDescriptorMediumType.PL132,
            DeviceDescriptor.KNX_IP_SYSTEM7 => DeviceDescriptorMediumType.KNX_IP,
            _ => default,
        };
        public static FirmwareType GetFirmwareType(this DeviceDescriptor value) => value switch
        {
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_0 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_1 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_2 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.TP1_BCU_1_SYSTEM_1_3 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.TP1_BCU_2_SYSTEM_2_0 => FirmwareType.SYSTEM_2,
            DeviceDescriptor.TP1_BCU_2_SYSTEM_2_1 => FirmwareType.SYSTEM_2,
            DeviceDescriptor.TP1_BCU_2_SYSTEM_2_5 => FirmwareType.SYSTEM_2,
            DeviceDescriptor.TP1_SYSTEM_300 => FirmwareType.SYSTEM_300,
            DeviceDescriptor.TP1_BIM_M112_0 => FirmwareType.SYSTEM_7,
            DeviceDescriptor.TP1_BIM_M112_1 => FirmwareType.SYSTEM_7,
            DeviceDescriptor.TP1_BIM_M112_5 => FirmwareType.SYSTEM_7,
            DeviceDescriptor.TP1_SYSTEM_B => FirmwareType.SYSTEM_B,
            DeviceDescriptor.TP1_IR_DECODER_0 => FirmwareType.IR_DECODER,
            DeviceDescriptor.TP1_IR_DECODER_1 => FirmwareType.IR_DECODER,
            DeviceDescriptor.TP1_COUPLER_0 => FirmwareType.COUPLER,
            DeviceDescriptor.TP1_COUPLER_1 => FirmwareType.COUPLER,
            DeviceDescriptor.TP1_COUPLER_2 => FirmwareType.COUPLER,
            DeviceDescriptor.TP1_KNXNETIP_ROUTER => FirmwareType.COUPLER,
            DeviceDescriptor.TP1_NONE_D => FirmwareType.NONE,
            DeviceDescriptor.TP1_NONE_E => FirmwareType.NONE,
            DeviceDescriptor.PL110_BCU_1_2 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.PL110_BCU_1_3 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.PL110_SYSTEM_B => FirmwareType.SYSTEM_B,
            DeviceDescriptor.PL110_MEDIA_COUPLER_PL_TP => FirmwareType.MEDIA_COUPLER_PL_TP,
            DeviceDescriptor.RF_BI_DIRECTIONAL_DEVICES => FirmwareType.RF_BI_DIRECTIONAL_DEVICES,
            DeviceDescriptor.RF_UNI_DIRECTIONAL_DEVICES => FirmwareType.RF_UNI_DIRECTIONAL_DEVICES,
            DeviceDescriptor.TP0_BCU_1 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.PL132_BCU_1 => FirmwareType.SYSTEM_1,
            DeviceDescriptor.KNX_IP_SYSTEM7 => FirmwareType.SYSTEM_7,
            _ => default,
        };
    }
}
