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

package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.profinet.config.ConfigurationProfinetDevice;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.device.ProfinetEmptyModule;
import org.apache.plc4x.java.profinet.device.ProfinetModule;
import org.apache.plc4x.java.profinet.device.ProfinetModuleImpl;
import org.apache.plc4x.java.profinet.readwrite.PnIoCm_Block_ExpectedSubmoduleReq;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetDeviceContextTests {

    @Test
    public void readProfinetAllocatedSubModulesLengthCheck()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name_1, PLC4X_1, (PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE)]]");

        Map<String, ConfigurationProfinetDevice> configuredDevices = configuration.getDevices().getConfiguredDevices();
        Map<String, ProfinetDevice> devices = new HashMap<>();
        for (Map.Entry<String, ConfigurationProfinetDevice> entry : configuredDevices.entrySet()) {
            devices.put(entry.getKey(),
                new ProfinetDevice(
                    new DummyMessageWrapper(),
                    entry.getValue().getDevicename(),
                    entry.getValue().getDeviceaccess(),
                    entry.getValue().getSubmodules(),
                    entry.getValue().getGsdHandler()
                )
            );
            devices.get(entry.getValue().getDevicename()).setIpAddress(entry.getValue().getIpaddress());
        }

        ProfinetDevice device =devices.get("DEVICE_NAME_1");
        device.setVendorDeviceId("CAFE", "0001");

        ProfinetModule[] modules = device.getDeviceContext().getModules();

        assertEquals(modules.length, 32);
    }

    @Test
    public void readProfinetAllocatedSubModulesTypeCheck()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name_1, PLC4X_1, (PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE)]]");

        Map<String, ConfigurationProfinetDevice> configuredDevices = configuration.getDevices().getConfiguredDevices();
        Map<String, ProfinetDevice> devices = new HashMap<>();
        for (Map.Entry<String, ConfigurationProfinetDevice> entry : configuredDevices.entrySet()) {
            devices.put(entry.getKey(),
                new ProfinetDevice(
                    new DummyMessageWrapper(),
                    entry.getValue().getDevicename(),
                    entry.getValue().getDeviceaccess(),
                    entry.getValue().getSubmodules(),
                    entry.getValue().getGsdHandler()
                )
            );
            devices.get(entry.getValue().getDevicename()).setIpAddress(entry.getValue().getIpaddress());
        }

        ProfinetDevice device = devices.get("DEVICE_NAME_1");
        device.setVendorDeviceId("CAFE", "0001");

        ProfinetModule[] modules = device.getDeviceContext().getModules();

        assertInstanceOf(ProfinetModuleImpl.class, modules[0]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[1]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[2]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[3]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[4]);
        assertInstanceOf(ProfinetEmptyModule.class, modules[5]);
    }

    @Test
    public void readProfinetAllocatedSubModulesTypeCheckEmptyModule()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name_1, PLC4X_1, (PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE)]]");

        Map<String, ConfigurationProfinetDevice> configuredDevices = configuration.getDevices().getConfiguredDevices();
        Map<String, ProfinetDevice> devices = new HashMap<>();
        for (Map.Entry<String, ConfigurationProfinetDevice> entry : configuredDevices.entrySet()) {
            devices.put(entry.getKey(),
                new ProfinetDevice(
                    new DummyMessageWrapper(),
                    entry.getValue().getDevicename(),
                    entry.getValue().getDeviceaccess(),
                    entry.getValue().getSubmodules(),
                    entry.getValue().getGsdHandler()
                )
            );
            devices.get(entry.getValue().getDevicename()).setIpAddress(entry.getValue().getIpaddress());
        }

        ProfinetDevice device = devices.get("DEVICE_NAME_1");
        device.setVendorDeviceId("CAFE", "0001");

        ProfinetModule[] modules = device.getDeviceContext().getModules();

        assertInstanceOf(ProfinetModuleImpl.class, modules[0]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[1]);
        assertInstanceOf(ProfinetEmptyModule.class, modules[2]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[3]);
        assertInstanceOf(ProfinetModuleImpl.class, modules[4]);
        assertInstanceOf(ProfinetEmptyModule.class, modules[5]);
    }

    @Test
    public void readExpectedSubModuleApiBlocks()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name_1, PLC4X_1, (PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE)]]");

        Map<String, ConfigurationProfinetDevice> configuredDevices = configuration.getDevices().getConfiguredDevices();
        Map<String, ProfinetDevice> devices = new HashMap<>();
        for (Map.Entry<String, ConfigurationProfinetDevice> entry : configuredDevices.entrySet()) {
            devices.put(entry.getKey(),
                new ProfinetDevice(
                    new DummyMessageWrapper(),
                    entry.getValue().getDevicename(),
                    entry.getValue().getDeviceaccess(),
                    entry.getValue().getSubmodules(),
                    entry.getValue().getGsdHandler()
                )
            );
            devices.get(entry.getValue().getDevicename()).setIpAddress(entry.getValue().getIpaddress());
        }

        ProfinetDevice device = devices.get("DEVICE_NAME_1");
        device.setVendorDeviceId("CAFE", "0001");

        List<PnIoCm_Block_ExpectedSubmoduleReq> moduleReq = device.getDeviceContext().getExpectedSubmoduleReq();

        assertEquals(3, moduleReq.get(0).getApis().get(0).getSubmodules().size());
        assertEquals(1, moduleReq.get(1).getApis().get(0).getSubmodules().size());
        assertEquals(1, moduleReq.get(2).getApis().get(0).getSubmodules().size());
        assertEquals(1, moduleReq.get(3).getApis().get(0).getSubmodules().size());

    }

}
