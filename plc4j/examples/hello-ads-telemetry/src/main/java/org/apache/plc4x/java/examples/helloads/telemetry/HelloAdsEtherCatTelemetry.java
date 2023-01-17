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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

enum EtherCatConstants {
    DeviceType(0x10000000, PlcValueType.UDINT),
    ManufacturerDeviceName(0x10080000, PlcValueType.STRING),
    //HardwareVersion(0x10090000, PlcValueType.),
    //SoftwareVersion(0x100A0000, PlcValueType.),
    //IdentityObjectNum(0x10180000, PlcValueType.),
    IdentityObjectVendorId(0x10180000, PlcValueType.UDINT),
    IdentityObjectProductCode(0x10180000, PlcValueType.UDINT),
    IdentityObjectRevisionNumber(0x10180000, PlcValueType.UDINT),
    IdentityObjectSoftwareVersion(0x10180000, PlcValueType.UDINT),
    //ConfigurationDataNum(0x80000000, PlcValueType.), // 0x8xxx0000 xxx = device index (starting with 0)
    ConfigurationDataAddress(0x80000001, PlcValueType.UINT),
    ConfigurationDataType(0x8000000, PlcValueType.STRING), // Not sure ...
    ConfigurationDataName(0x8000000, PlcValueType.STRING),
    ConfigurationDataDeviceType(0x8000000, PlcValueType.UDINT),
    ConfigurationDataVendorId(0x8000000, PlcValueType.UDINT),
    ConfigurationDataProductCode(0x8000000, PlcValueType.UDINT),
    ConfigurationDataRevisionNumber(0x8000000, PlcValueType.UDINT),
    ConfigurationDataSerialNumber(0x8000000, PlcValueType.UDINT),
    ConfigurationDataMailboxOutSize(0x8000000, PlcValueType.UINT),
    ConfigurationDataMailboxInSize(0x8000000, PlcValueType.UINT),
    ConfigurationDataLinkStatus(0x8000000, PlcValueType.USINT),
    ConfigurationDataLinkPreset(0x8000000, PlcValueType.USINT),
    ConfigurationDataFlags(0x8000000, PlcValueType.UINT),
//    StateMachine(0xA0000000, PlcValueType.),      // 0xAxxx0000 xxx = device index (starting with 0)
//    ScanSlaves(0xF002000, PlcValueType.),
    ConfiguredSlavesNum(0xF0200000, PlcValueType.USINT);
//    FrameStatistics(0xF1200000, PlcValueType.);

    final int typeNumber;
    final PlcValueType plcValueType;
    EtherCatConstants(int typeNumber, PlcValueType plcValueType) {
        this.typeNumber = typeNumber;
        this.plcValueType = plcValueType;
    }
}
public class HelloAdsEtherCatTelemetry {

    private static final Logger logger = LoggerFactory.getLogger(HelloAdsTelemetry.class);
    public static void main(String[] args) {
        if(args.length != 2) {
            logger.error("Usage: HelloAdsTelemetry {ip-address of PLC} {local ip-address}");
            System.exit(1);
        }

        String remoteIp = args[0];
        String localIp = args[1];
        // The AmsNetId of the PLC usually is {ip}.1.1 and that of the EtherCAT master is {ip}.3.1
        // The port number equals the EtherCAT address. For the EtherCAT master, this port is 0xFFFF = 65535
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection(String.format("ads:tcp://%s?targetAmsNetId=%s.3.1&targetAmsPort=65535&sourceAmsNetId=%s.1.1&sourceAmsPort=65534&load-symbol-and-data-type-tables=false", remoteIp, remoteIp, localIp))) {
            int vendorId = connection.readRequestBuilder().addTagAddress("vendorId", "0x0000F302/0x10180001:UDINT").build().execute().get().getInteger("vendorId");
            int productCode = connection.readRequestBuilder().addTagAddress("productCode", "0x0000F302/0x10180002:UDINT").build().execute().get().getInteger("productCode");
            int revisionNumber = connection.readRequestBuilder().addTagAddress("revisionNumber", "0x0000F302/0x10180003:UDINT").build().execute().get().getInteger("revisionNumber");
            int serialNumber = connection.readRequestBuilder().addTagAddress("serialNumber", "0x0000F302/0x10180004:UDINT").build().execute().get().getInteger("serialNumber");
            logger.info("EtherCAT Master: Vendor Id: {}, Product Code: {}, Revision Number: {}, Serial Number {}", vendorId, productCode, revisionNumber, serialNumber);

            // Load the number of EtherCAT slaves:
            int numSlaves = connection.readRequestBuilder().addTagAddress("numberOfSlaves", "0x0000F302/0xF0200000:USINT").build().execute().get().getInteger("numberOfSlaves");

            // Load the number of slaves and their etherCatAddresses
            // NOTE: We need to do this without using multi-item-requests as it seems that this part of the system doesn't support this.
            Map<Integer, Integer> etherCatAddresses = new HashMap<>();
            for(int i = 1; i < numSlaves; i++) {
                String name = "slave-" + i;
                String address = String.format("0x0000F302/0xF020%04X:UINT", i);
                PlcReadResponse plcReadResponse = connection.readRequestBuilder().addTagAddress(name, address).build().execute().get();
                if (plcReadResponse.getResponseCode(name) == PlcResponseCode.OK) {
                    int etherCatAddress = plcReadResponse.getInteger(name);
                    logger.info("Slave {} has EtherCAT address {}", i, etherCatAddress);
                    etherCatAddresses.put(i, etherCatAddress);

                }
            }

        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

}
