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

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public enum SzlSublist : byte
    {
        NONE = 0x00,
        MODULE_IDENTIFICATION = 0x11,
        CPU_FEATURES = 0x12,
        USER_MEMORY_AREA = 0x13,
        SYSTEM_AREAS = 0x14,
        BLOCK_TYPES = 0x15,
        PRIORITY_CLASSES = 0x16,
        EXTENDED_PRIORITY_CLASSES = 0x17,
        OPERATING_SYSTEM = 0x18,
        STATUS_MODULE_LEDS = 0x19,
        COMPONENT_IDENTIFICATION = 0x1C,
        INTERRUPT_INFO = 0x21,
        INTERRUPT_STATUS = 0x22,
        OPERATION_HISTORY = 0x23,
        OPERATING_MODES = 0x24,
        ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS = 0x25,
        COMMUNICATION_MODE_DATA = 0x31,
        COMMUNICATION_STATUS_DATA = 0x32,
        DIAGNOSTIC_DATA_ON_COMM_OBJECTS = 0x33,
        STATUS_H_SYSTEM = 0x70,
        H_CPU_GROUP_INFORMATION = 0x71,
        STATUS_SINGLE_MODULE_LED = 0x74,
        SWITCHED_DP_SLAVES_H_SYSTEM = 0x75,
        DP_MASTER_SYSTEM_INFORMATION = 0x90,
        MODULE_STATUS_INFORMATION = 0x91,
        RACK_OR_STATION_STATUS_INFORMATION = 0x92,
        RACK_OR_STATION_STATUS_INFORMATION_2 = 0x94,
        ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION = 0x95,
        MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP = 0x96,
        TOOL_CHANGER_INFORMATION_PROFINET = 0x9C,
        DIAGNOSTIC_BUFFER = 0xA0,
        MODULE_DIAGNOSTIC_INFORMATION_DR0 = 0xB1,
        MODULE_DIAGNOSTIC_INFORMATION_DR1_GI = 0xB2,
        MODULE_DIAGNOSTIC_INFORMATION_DR1_LA = 0xB3,
        DIAGNOSTIC_DATA_DP_SLAVE = 0xB4,
    }
}
