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
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

public class HelloAdsEtherCatTelemetry {

    private static final Logger logger = LoggerFactory.getLogger(HelloAdsTelemetry.class);

    private static final int AoEGroupIndex = 0x0000F302;

    protected void outputEtherCatData(String remoteIp, String localAmsNetId, String remoteAmdNetId) {
        Map<Integer, EtherCatDevice> devices = new HashMap<>();
        // The AmsNetId of the PLC usually is {ip}.1.1 and that of the EtherCAT master is {ip}.3.1
        // The port number equals the EtherCAT address. For the EtherCAT master, this port is 0xFFFF = 65535
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection(String.format("ads:tcp://%s?targetAmsNetId=%s&targetAmsPort=65535&sourceAmsNetId=%s&sourceAmsPort=65534&load-symbol-and-data-type-tables=false", remoteIp, remoteAmdNetId, localAmsNetId))) {
            String manufacturerDeviceName = connection.readRequestBuilder().addTagAddress("manufacturerDeviceName", getAddress(EtherCatMasterConstants.ManufacturerDeviceName)).build().execute().get().getString("manufacturerDeviceName");
            int hardwareVersion = connection.readRequestBuilder().addTagAddress("hardwareVersion", getAddress(EtherCatMasterConstants.HardwareVersion)).build().execute().get().getInteger("hardwareVersion");
            String softwareVersion = connection.readRequestBuilder().addTagAddress("softwareVersion", getAddress(EtherCatMasterConstants.SoftwareVersion)).build().execute().get().getString("softwareVersion");
            logger.info("Found Device: {} Hardware Version {}, Software Version {}", manufacturerDeviceName, hardwareVersion, softwareVersion);

            logger.info("Identity Object:");
            outputEtherCatSection(connection, EtherCatMasterConstants.IdentityObjectNum.offset, EtherCatMasterConstants.IdentityObjectNum.offset);

            // Load the number of EtherCAT slaves:
            int numSlaves = connection.readRequestBuilder().addTagAddress("numberOfSlaves", "0x0000F302/0xF0200000:USINT").build().execute().get().getInteger("numberOfSlaves");

            // Load the number of slaves and their etherCatAddresses
            // NOTE: We need to do this without using multi-item-requests as it seems that this part of the system doesn't support this.
            for(int i = 0; i < numSlaves; i++) {
                logger.info("Slave {}", i);
                int etherCatAddressOffset = EtherCatMasterConstants.ConfiguredSlavesNum.offset | (i + 1);
                int configDataOffset = EtherCatMasterConstants.ConfigurationDataNum.offset | (i << 16);
                int stateMachineOffset = EtherCatMasterConstants.StateMachineNum.offset | (i << 16);

                String etherCatAddressAddress = String.format("0x%08X/0x%08X:%s", AoEGroupIndex, etherCatAddressOffset, PlcValueType.UINT.name());
                int etherCatAddress = connection.readRequestBuilder().addTagAddress("etherCatAddress", etherCatAddressAddress).build().execute().get().getInteger("etherCatAddress");
                logger.info(" - EtherCat Address: {}", etherCatAddress);

                logger.info(" - Configuration Data:");
                EtherCatDevice device = outputEtherCatSection(connection, configDataOffset, EtherCatMasterConstants.ConfigurationDataNum.offset);
                //logger.info(" - State Machine Data");
                //outputEtherCatSection(connection, stateMachineOffset, EtherCatConstant.StateMachineNum.offset);
                devices.put(i, device);

                //String manufacturerDeviceName = connection.readRequestBuilder().addTagAddress("lalala", getAddress(EtherCatConstant.ManufacturerDeviceName)).build().execute().get().getString("manufacturerDeviceName");

            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        logger.info("Loading EtherCAT device information:");
        for (Map.Entry<Integer, EtherCatDevice> device : devices.entrySet()) {
            int deviceIndex = device.getKey();
            EtherCatDevice etherCatDevice = device.getValue();
            logger.info(" - Connecting with device {} on EtherCAT address {}", deviceIndex, 1001);
            try (PlcConnection etherCatConnection = PlcDriverManager.getDefault().getConnectionManager().getConnection(String.format("ads:tcp://%s?targetAmsNetId=%s&targetAmsPort=%d&sourceAmsNetId=%s&sourceAmsPort=65534&load-symbol-and-data-type-tables=false", remoteIp, remoteAmdNetId, etherCatDevice.getEtherCatAddress(), localAmsNetId))) {
                String etherCatAddressAddress = String.format("0x%08X/0x%08X:%s", AoEGroupIndex, 0x60000001, PlcValueType.BOOL.name());
                PlcReadRequest build = etherCatConnection.readRequestBuilder()
                    .addTagAddress("Channel 1", etherCatAddressAddress)
                    .build();
                PlcReadResponse plcReadResponse = build.execute().get();
                System.out.println(plcReadResponse);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    }
    
    protected EtherCatDevice outputEtherCatSection(PlcConnection connection, int baseOffset, int baseTypeOffset) throws Exception {
        String sectionNumAddress = String.format("0x%08X/0x%08X:%s", AoEGroupIndex, baseOffset, PlcValueType.USINT.name());

        int etherCatAddress = 0;
        String deviceName = null;
        int vendorId = 0;
        int productCode = 0;
        int revisionNumber = 0;

        int identityObjectNum = connection.readRequestBuilder().addTagAddress("num", sectionNumAddress).build().execute().get().getInteger("num");
        for (int i = 1; i < identityObjectNum; i++) {
            int offset = baseOffset | i;
            int typeOffset = baseTypeOffset | i;
            EtherCatMasterConstants etherCatConstantAddress = EtherCatMasterConstants.enumForValue(typeOffset);
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
                        switch (etherCatConstantAddress) {
                            case ConfigurationDataAddress:
                                etherCatAddress = value.getInteger();
                                break;
                            case ConfigurationDataName:
                                deviceName = value.getString();
                                break;
                            case ConfigurationDataVendorId:
                                vendorId = value.getInteger();
                                break;
                            case ConfigurationDataProductCode:
                                productCode = value.getInteger();
                                break;
                            case ConfigurationDataRevisionNumber:
                                revisionNumber = value.getInteger();
                                break;
                        }
                        logger.info("    - {}: {}", etherCatConstantAddress.name(), value.toString());
                    } else {
                        logger.info("    - Unknown: {}", value.toString());
                    }
                }
            } catch (Exception e) {
                // Ignore this ...
            }
        }
        return new EtherCatDevice(etherCatAddress, deviceName, vendorId, productCode, revisionNumber);
    }
    
    
    protected String getAddress(EtherCatMasterConstants variable) {
        String dataTypeName = variable.plcValueType.name();
        if (variable.plcValueType == PlcValueType.STRING) {
            dataTypeName += "(255)";
        }
        return String.format("0x%08X/0x%08X:%s", AoEGroupIndex, variable.offset, dataTypeName);
    }

    public static void main(String[] args) {
        if(args.length != 3) {
            logger.error("Usage: HelloAdsTelemetry {remote ip-address} {local-ams-net-id} {remote-ams-net-id}");
            System.exit(1);
        }

        String remoteIp = args[0];
        String localAmsNetId = args[1];
        String remoteAmsNetId = args[2];
        new HelloAdsEtherCatTelemetry().outputEtherCatData(remoteIp, localAmsNetId, remoteAmsNetId);
    }
    
}