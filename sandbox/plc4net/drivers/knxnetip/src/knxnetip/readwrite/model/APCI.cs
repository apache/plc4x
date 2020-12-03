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

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public enum APCI : byte
    {

        GROUP_VALUE_READ_PDU = 0x0,
        GROUP_VALUE_RESPONSE_PDU = 0x1,
        GROUP_VALUE_WRITE_PDU = 0x2,
        INDIVIDUAL_ADDRESS_WRITE_PDU = 0x3,
        INDIVIDUAL_ADDRESS_READ_PDU = 0x4,
        INDIVIDUAL_ADDRESS_RESPONSE_PDU = 0x5,
        ADC_READ_PDU = 0x6,
        ADC_RESPONSE_PDU = 0x7,
        MEMORY_READ_PDU = 0x8,
        MEMORY_RESPONSE_PDU = 0x9,
        MEMORY_WRITE_PDU = 0xA,
        USER_MESSAGE_PDU = 0xB,
        DEVICE_DESCRIPTOR_READ_PDU = 0xC,
        DEVICE_DESCRIPTOR_RESPONSE_PDU = 0xD,
        RESTART_PDU = 0xE,
        OTHER_PDU = 0xF,

    }

}

