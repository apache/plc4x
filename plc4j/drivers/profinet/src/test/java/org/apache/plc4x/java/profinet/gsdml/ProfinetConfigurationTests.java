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

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.io.File;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetConfigurationTests {

    /*
        Profinet GSD File Directory Configuration Test
     */
    @Test
    public void readGsdDirectory()  {
        String directory = "src/test/resources";
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, PLC4X, {PLC4X_1}]]&gsddirectory=" + directory);

        assertEquals(1, configuration.getGsdFiles().getGsdFiles().size());
    }

    @Test
    public void readGsdFilesInDirectory()  {
        String directory = "src/test/resources";

        new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, PLC4X, {PLC4X_1}]]&gsddirectory=" + directory);

        Map<String, ProfinetISO15745Profile> gsdFiles = ProfinetConfiguration.getGsdFiles().getGsdFiles();
        assertEquals(gsdFiles.size(), 1);
    }

    @Test
    public void parseJoinedDeviceConfiguration() {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources&devices=[[device_name_1,device_access_1,{submodule_1,submodule_2}]]");
        assertEquals(1, configuration.getDevices().getConfiguredDevices().size());
    }

    /*
        Profinet GSD File Directory Configuration Test
     */
    @Test
    public void parseJoinedDeviceConfigurationExtraSpaces() {
        String[] deviceNames = new String[] {"DEVICE_NAME_1"};

        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, device_access_1, {submodule_1, submodule_2}]]&gsddirectory=src/test/resources");

        Map<String, ProfinetDevice> devices = configuration.getDevices().getConfiguredDevices();

        for (String deviceName : deviceNames) {
            assert(devices.containsKey(deviceName));
        }
    }

    @Test
    public void readProfinetDevicesMultiple() {
        String[] deviceNames = new String[] {"DEVICE_NAME_1","DEVICE_NAME_2","DEVICE_NAME_3"};

        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, PLC4X_1, {PLC4X_01,PLC4X_02,PLC4X_01,PLC4X_02}],[device_name_2, PLC4X_1, {PLC4X_01,PLC4X_02,PLC4X_01,PLC4X_02}],[device_name_3, PLC4X_1, {PLC4X_01,PLC4X_02,PLC4X_01,PLC4X_02}]]&gsddirectory=src/test/resources");

        Map<String, ProfinetDevice> devices = configuration.getDevices().getConfiguredDevices();

        for (String deviceName : deviceNames) {
            assert(devices.containsKey(deviceName));
        }
    }

    @Test
    public void readProfinetLowerCase() {
        String[] deviceName = new String[] {"device_Name_1"};
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, device_access_1, {submodule_1, submodule_2}]]&gsddirectory=src/test/resources");

        Map<String, ProfinetDevice> devices = configuration.getDevices().getConfiguredDevices();

        for (String mac : deviceName) {
            assert(devices.containsKey(mac.replace(":", "").toUpperCase()));
        }
    }

    @Test
    public void setIncorrectSubModule() {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, PLC4X_1, {PLC4X_01, PLC4X_02, PLC4X_01, PLC4X_02}]]&gsddirectory=src/test/resources");

        Map<String, ProfinetDevice> devices = configuration.getDevices().getConfiguredDevices();

        XmlMapper xmlMapper = new XmlMapper();
        assertThrows(PlcException.class, () -> devices.get("DEVICE_NAME_1").getDeviceContext().setGsdFile(xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class)));
    }


    @Test
    public void setCorrectSubModule() {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, PLC4X_1, {PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE}]]&gsddirectory=src/test/resources");

        Map<String, ProfinetDevice> devices = configuration.getDevices().getConfiguredDevices();

        XmlMapper xmlMapper = new XmlMapper();
        assertDoesNotThrow(() -> devices.get("DEVICE_NAME_1").getDeviceContext().setGsdFile(xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class)));
    }

    @Test
    public void setCorrectSubModuleCaseInsensitive() {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[[device_name_1, PLC4X_1, {PLC4X_DUMMY_MODULE, PLC4X_dummy_MODULE, PLC4X_DUMMY_MODULE, PLC4X_DUMMY_MODULE}]]&gsddirectory=src/test/resources");

        Map<String, ProfinetDevice> devices = configuration.getDevices().getConfiguredDevices();

        XmlMapper xmlMapper = new XmlMapper();
        assertDoesNotThrow(() -> devices.get("DEVICE_NAME_1").getDeviceContext().setGsdFile(xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class)));
    }
}
