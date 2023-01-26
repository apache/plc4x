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

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

enum EtherCatConstant {
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

    private static final Map<Integer, EtherCatConstant> map;
    static {
        map = new HashMap<>();
        for (EtherCatConstant value : EtherCatConstant.values()) {
            map.put(value.offset, value);
        }
    }

    final int offset;
    final PlcValueType plcValueType;
    EtherCatConstant(int offset, PlcValueType plcValueType) {
        this.offset = offset;
        this.plcValueType = plcValueType;
    }

    public static EtherCatConstant enumForValue(int offset) {
        return map.get(offset);
    }

}
public class HelloAdsEtherCatTelemetry {

    private static final Logger logger = LoggerFactory.getLogger(HelloAdsTelemetry.class);

    private static final int AoEGroupIndex = 0x0000F302;

    protected void outputEtherCatData(String localIp, String remoteIp, String remoteAmdNetId) {
        // The AmsNetId of the PLC usually is {ip}.1.1 and that of the EtherCAT master is {ip}.3.1
        // The port number equals the EtherCAT address. For the EtherCAT master, this port is 0xFFFF = 65535
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection(String.format("ads:tcp://%s?targetAmsNetId=%s.3.1&targetAmsPort=65535&sourceAmsNetId=%s.1.1&sourceAmsPort=65534&load-symbol-and-data-type-tables=false", remoteIp, remoteIp, localIp))) {
            String manufacturerDeviceName = connection.readRequestBuilder().addTagAddress("manufacturerDeviceName", getAddress(EtherCatConstant.ManufacturerDeviceName)).build().execute().get().getString("manufacturerDeviceName");
            int hardwareVersion = connection.readRequestBuilder().addTagAddress("hardwareVersion", getAddress(EtherCatConstant.HardwareVersion)).build().execute().get().getInteger("hardwareVersion");
            String softwareVersion = connection.readRequestBuilder().addTagAddress("softwareVersion", getAddress(EtherCatConstant.SoftwareVersion)).build().execute().get().getString("softwareVersion");
            logger.info("Found Device: {} Hardware Version {}, Software Version {}", manufacturerDeviceName, hardwareVersion, softwareVersion);

            logger.info("Identity Object:");
            outputEtherCatSection(connection, EtherCatConstant.IdentityObjectNum.offset, EtherCatConstant.IdentityObjectNum.offset);

            // Load the number of EtherCAT slaves:
            int numSlaves = connection.readRequestBuilder().addTagAddress("numberOfSlaves", "0x0000F302/0xF0200000:USINT").build().execute().get().getInteger("numberOfSlaves");

            // Load the number of slaves and their etherCatAddresses
            // NOTE: We need to do this without using multi-item-requests as it seems that this part of the system doesn't support this.
            Map<Integer, Integer> etherCatAddresses = new HashMap<>();
            for(int i = 0; i < numSlaves; i++) {
                logger.info("Slave {}", i);
                int etherCatAddressOffset = EtherCatConstant.ConfiguredSlavesNum.offset | (i + 1);
                int configDataOffset = EtherCatConstant.ConfigurationDataNum.offset | (i << 16);
                int stateMachineOffset = EtherCatConstant.StateMachineNum.offset | (i << 16);

                String etherCatAddressAddress = String.format("0x%08X/0x%08X:%s", AoEGroupIndex, etherCatAddressOffset, PlcValueType.UINT.name());
                int etherCatAddress = connection.readRequestBuilder().addTagAddress("etherCatAddress", etherCatAddressAddress).build().execute().get().getInteger("etherCatAddress");
                logger.info(" - EtherCat Address: {}", etherCatAddress);

                logger.info(" - Configuration Data:");
                outputEtherCatSection(connection, configDataOffset, EtherCatConstant.ConfigurationDataNum.offset);
                logger.info(" - State Machine Data");
                outputEtherCatSection(connection, stateMachineOffset, EtherCatConstant.StateMachineNum.offset);
            }

        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    } 
    
    protected void outputEtherCatSection(PlcConnection connection, int baseOffset, int baseTypeOffset) throws Exception {
        String sectionNumAddress = String.format("0x%08X/0x%08X:%s", AoEGroupIndex, baseOffset, PlcValueType.USINT.name());
        int identityObjectNum = connection.readRequestBuilder().addTagAddress("num", sectionNumAddress).build().execute().get().getInteger("num");
        for (int i = 1; i < identityObjectNum; i++) {
            int offset = baseOffset | i;
            int typeOffset = baseTypeOffset | i;
            EtherCatConstant etherCatConstantAddress = EtherCatConstant.enumForValue(typeOffset);
            PlcValueType etherCatConstantType = (etherCatConstantAddress != null) ? etherCatConstantAddress.plcValueType : PlcValueType.USINT;
            String address = String.format("0x%08X/0x%08X:%s", AoEGroupIndex, offset, etherCatConstantType.name());
            if(etherCatConstantType == PlcValueType.STRING) {
                address += "(255)";
            }
            try {
                PlcReadResponse readResponse = connection.readRequestBuilder().addTagAddress("value", address).build().execute().get();
                if (readResponse.getResponseCode("value") == PlcResponseCode.OK) {
                    PlcValue value = readResponse.getPlcValue("value");
                    if (etherCatConstantAddress != null) {
                        logger.info("    - {}: {}", etherCatConstantAddress.name(), value.toString());
                    } else {
                        logger.info("    - Unknown: {}", value.toString());
                    }
                }
            } catch (Exception e) {
                // Ignore this ...
            }
        }
    }
    
    
    protected String getAddress(EtherCatConstant variable) {
        String dataTypeName = variable.plcValueType.name();
        if (variable.plcValueType == PlcValueType.STRING) {
            dataTypeName += "(255)";
        }
        return String.format("0x%08X/0x%08X:%s", AoEGroupIndex, variable.offset, dataTypeName);
    }

    public static void main(String[] args) {
        if(args.length != 2) {
            logger.error("Usage: HelloAdsTelemetry {ip-address of PLC} {local ip-address}");
            System.exit(1);
        }

        String remoteIp = args[0];
        String localIp = args[1];
        new HelloAdsEtherCatTelemetry().outputEtherCatData(localIp, remoteIp, remoteIp + ".3.1");
    }
    
}