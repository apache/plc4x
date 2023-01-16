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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

enum Module {
    NIC(0x00000002),
    Time(0x00000003),
    UserManagement(0x00000004),
    RAS(0x00000005),
    FTP(0x00000006),
    SMB(0x00000007),
    TwinCat(0x00000008),
    Software(0x0000000A),
    COU(0x0000000B),
    Memory(0x0000000C),
    FirewallWinCE(0x0000000E),
    FileSystemObject(0x00000010),
    DisplayDevice(0x00000013),
    EWF(0x00000014),
    FBWF(0x00000015),
    OS(0x00000018),
    RAID(0x00000019),
    Fan(0x0000001B),
    Mainboard(0x0000001C),
    DiskManagement(0x0000001D),
    UPS(0x0000001E),
    PhysicalDrive(0x0000001F),
    MassStorageDrive(0x00000020),
    UnifiedWriteFilter(0x00000021),
    IO(0x00000022),
    Misc(0x00000100);

    final int typeNumber;
    Module(int typeNumber) {
        this.typeNumber = typeNumber;
    }

}
public class HelloAdsTelemetry {

    private static final Logger logger = LoggerFactory.getLogger(HelloAdsTelemetry.class);
    public static void main(String[] args) {
        if(args.length != 2) {
            logger.error("Usage: HelloAdsTelemetry {ip-address of PLC} {local ip-address}");
            System.exit(1);
        }

        String remoteIp = args[0];
        String localIp = args[1];
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection(String.format("ads:tcp://%s?targetAmsNetId=%s.1.1&targetAmsPort=10000&sourceAmsNetId=%s.1.1&sourceAmsPort=65534&load-symbol-and-data-type-tables=false", remoteIp, remoteIp, localIp))) {
            // Load the number of modules:
            int numModules = connection.readRequestBuilder().addTagAddress("numberOfModules", "0x0000F302/0xF0200000:UINT").build().execute().get().getInteger("numberOfModules");

            // Load the mdp type and index for each module
            // NOTE: We need to do this without using multi-item-requests as it seems that this part of the system doesn't support this.
            Map<Integer, Integer> moduleTypeIdMap = new HashMap<>();
            for(int i = 1; i < numModules; i++) {
                String name = "module-" + i;
                String address = String.format("0x0000F302/0xF020%04X:UDINT", i);
                PlcReadResponse plcReadResponse = connection.readRequestBuilder().addTagAddress(name, address).build().execute().get();
                if (plcReadResponse.getResponseCode(name) == PlcResponseCode.OK) {
                    int value = plcReadResponse.getInteger(name);
                    int mdpType = ((value & 0xFFFF0000) >> 16);
                    int mdpId = value & 0x0000FFFF;
                    logger.info("Module {} has mdp type {} and mdp id {}", i, mdpType, mdpId);
                    moduleTypeIdMap.put(mdpType, mdpId);
                }
            }

            // Read the ADS Version information.
            if(moduleTypeIdMap.containsKey(Module.TwinCat.typeNumber)) {
                Integer mdpId = moduleTypeIdMap.get(Module.TwinCat.typeNumber);
                int addrAdsTypeMain = (mdpId << 20) | 0x80010001;
                int addrAdsTypeMinor = (mdpId << 20) | 0x80010002;
                int addrAdsTypeBuild = (mdpId << 20) | 0x80010003;
                int twinCatMainVersion = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UINT", addrAdsTypeMain)).build().execute().get().getInteger("value");
                int twinCatMinorVersion = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UINT", addrAdsTypeMinor)).build().execute().get().getInteger("value");
                int twinCatBuildVersion = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UINT", addrAdsTypeBuild)).build().execute().get().getInteger("value");
                logger.info("TwinCat Version: {}.{}.{}", twinCatMainVersion, twinCatMinorVersion, twinCatBuildVersion);
            }
            // Read the CPU Frequency and Utilization.
            if(moduleTypeIdMap.containsKey(Module.COU.typeNumber)) {
                Integer mdpId = moduleTypeIdMap.get(Module.COU.typeNumber);
                int addrCpuFrequency = (mdpId << 20) | 0x80010001;
                int addrCpuUsage = (mdpId << 20) | 0x80010002;
                int cpuFrequency = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UDINT", addrCpuFrequency)).build().execute().get().getInteger("value");
                int cpuUsage = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UINT", addrCpuUsage)).build().execute().get().getInteger("value");
                logger.info("CPU: Frequency: {}MHz Usage: {}%", cpuFrequency, cpuUsage);
            }
            // Read the Memory usage.
            if(moduleTypeIdMap.containsKey(Module.Memory.typeNumber)) {
                Integer mdpId = moduleTypeIdMap.get(Module.Memory.typeNumber);
                int addrMemoryAllocated = (mdpId << 20) | 0x80010001;
                int addrMemoryAvailable = (mdpId << 20) | 0x80010002;
                int memoryAllocated = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UDINT", addrMemoryAllocated)).build().execute().get().getInteger("value");
                int memoryAvailable = connection.readRequestBuilder().addTagAddress("value", String.format("0x0000F302/0x%8X:UDINT", addrMemoryAvailable)).build().execute().get().getInteger("value");
                logger.info("Memory: Allocated: {}MB, Available: {}MB, Usage {}%", memoryAllocated / (1024 * 1024), memoryAvailable / (1024 * 1024), ((float) 100 / memoryAllocated) * memoryAvailable);
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

}
