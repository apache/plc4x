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
package org.apache.plc4x.java.examples.helloads.telemetry;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.HashMap;
import java.util.Map;

public enum EtherCatMasterConstants {
    DeviceType(0x10000000, PlcValueType.UDINT),
    ManufacturerDeviceName(0x10080000, PlcValueType.STRING),
    HardwareVersion(0x10090000, PlcValueType.UINT),
    SoftwareVersion(0x100A0000, PlcValueType.STRING),
    IdentityObjectNum(0x10180000, PlcValueType.USINT),
    IdentityObjectVendorId(0x10180001, PlcValueType.UDINT),
    IdentityObjectProductCode(0x10180002, PlcValueType.UDINT),
    IdentityObjectRevisionNumber(0x10180003, PlcValueType.UDINT),
    IdentityObjectSerialNumber(0x10180004, PlcValueType.UDINT),
    ConfigurationDataNum(0x80000000, PlcValueType.USINT), // 0x8xxx0000 xxx = device index (starting with 0)
    ConfigurationDataAddress(0x80000001, PlcValueType.UINT),
    ConfigurationDataType(0x80000002, PlcValueType.STRING), // Not sure ...
    ConfigurationDataName(0x80000003, PlcValueType.STRING),
    ConfigurationDataDeviceType(0x80000004, PlcValueType.UDINT),
    ConfigurationDataVendorId(0x80000005, PlcValueType.UDINT),
    ConfigurationDataProductCode(0x80000006, PlcValueType.UDINT),
    ConfigurationDataRevisionNumber(0x80000007, PlcValueType.UDINT),
    ConfigurationDataSerialNumber(0x80000008, PlcValueType.UDINT),
    ConfigurationDataMailboxOutSize(0x80000021, PlcValueType.UINT),
    ConfigurationDataMailboxInSize(0x80000022, PlcValueType.UINT),
    ConfigurationDataLinkStatus(0x80000023, PlcValueType.USINT),
    ConfigurationDataLinkPreset(0x80000024, PlcValueType.USINT),
    ConfigurationDataFlags(0x80000025, PlcValueType.UINT),
    StateMachineNum(0xA0000000, PlcValueType.USINT),      // 0xAxxx0000 xxx = device index (starting with 0)
    StateMachineAlStatus(0xA0000001, PlcValueType.UINT),
    StateMachineAlControl(0xA0000002, PlcValueType.UINT),
    StateMachineLastAlStatusCode(0xA0000003, PlcValueType.UINT),
    StateMachineLinkConnectionStatus(0xA0000004, PlcValueType.USINT),
    StateMachineLinkControl(0xA0000005, PlcValueType.USINT),
    StateMachineFixedAddressPort0(0xA0000006, PlcValueType.UINT),
    StateMachineFixedAddressPort1(0xA0000007, PlcValueType.UINT),
    StateMachineFixedAddressPort2(0xA0000008, PlcValueType.UINT),
    StateMachineFixedAddressPort3(0xA0000009, PlcValueType.UINT),
    StateMachineCrcErrorCountPort0(0xA000000A, PlcValueType.UDINT),
    StateMachineCrcErrorCountPort1(0xA000000B, PlcValueType.UDINT),
    StateMachineCrcErrorCountPort2(0xA000000C, PlcValueType.UDINT),
    StateMachineCrcErrorCountPort3(0xA000000D, PlcValueType.UDINT),
    StateMachineCyclicWcErrorCount(0xA000000E, PlcValueType.UDINT),
    StateMachineSlaveNotPresentCount(0xA000000F, PlcValueType.UDINT),
    StateMachineAbnormalStateChangeCount(0xA0000010, PlcValueType.UDINT),
    //StateMachineDisableAutomaticLinkControl(0xA0000011, PlcValueType.BOOL),
    //ScanSlavesNum(0xF002000, PlcValueType.),
    //ScanSlavesScanCommand(0xF002001, PlcValueType.),
    ScanSlavesScanStatus(0xF002002, PlcValueType.USINT),
    //ScanSlavesScanReply(0xF002003, PlcValueType.),
    ConfiguredSlavesNum(0xF0200000, PlcValueType.USINT), // Get the EtherCAT address by looping through the modules
    //    FrameStatisticsNum(0xF1200000, PlcValueType.);
    FrameStatisticsCyclicLostFramesCount(0xF1200001, PlcValueType.UDINT),
    FrameStatisticsAcyclicLostFramesCount(0xF1200002, PlcValueType.UDINT);
    //DiagnosticNum(0xF2000000, PlcValueType.),
    //DiagnosticResetDiagnosticCounters(0xF2000001, PlcValueType.);

    private static final Map<Integer, EtherCatMasterConstants> map;
    static {
        map = new HashMap<>();
        for (EtherCatMasterConstants value : EtherCatMasterConstants.values()) {
            map.put(value.offset, value);
        }
    }

    final int offset;
    final PlcValueType plcValueType;
    EtherCatMasterConstants(int offset, PlcValueType plcValueType) {
        this.offset = offset;
        this.plcValueType = plcValueType;
    }

    public static EtherCatMasterConstants enumForValue(int offset) {
        return map.get(offset);
    }

}
