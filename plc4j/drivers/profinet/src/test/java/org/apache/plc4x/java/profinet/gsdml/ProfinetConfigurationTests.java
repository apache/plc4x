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
import org.apache.commons.codec.DecoderException;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDeviceContext;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseRequest;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

import static org.junit.jupiter.api.Assertions.assertEquals;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetConfigurationTests {

    /*
        Profinet GSD File Directory Configuration Test
     */
    @Test
    public void readGsdDirectory()  {

        String directory = "/home/plc4x/gsd_directory";
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=" + directory);

        assertEquals(directory, configuration.getGsdDirectory());
    }

    @Test
    public void readGsdFilesInDirectory()  {
        String directory = "src/test/resources";
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();

        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=" + directory);

        ProfinetDriverContext driverContext = new ProfinetDriverContext();
        driverContext.setConfiguration(configuration);
        protocolLogic.setDriverContext(driverContext);

        Map<String, ProfinetISO15745Profile> gsdFiles = ((ProfinetDriverContext) protocolLogic.getDriverContext()).getGsdFiles();
        assertEquals(gsdFiles.size(), 1);
    }

    @Test
    public void readGsdFilesInDirectoryUsingTilde()  {
        String directory = "~/Documents/Profinet/gsd";
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=" + directory);

        assertEquals(directory, configuration.getGsdDirectory());
    }

    /*
        Profinet GSD File Directory Configuration Test
     */
    @Test
    public void readProfinetDevices() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"CA:FE:00:00:00:01"};
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "deviceaccess=[PLC4X_1]&devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setConfiguration(configuration);
        protocolLogic.setConfiguration(configuration);

        try {
            protocolLogic.setDevices();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
        Map<String, ProfinetDevice> devices = protocolLogic.getDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }

    @Test
    public void readProfinetDevicesMultiple() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"CA:FE:00:00:00:01","CA:FE:00:00:00:02","CA:FE:00:00:00:03"};
        String subModules = "[[PLC4X_01,PLC4X_02,PLC4X_01,PLC4X_02],[PLC4X_01,PLC4X_02,PLC4X_01,PLC4X_02],[PLC4X_01,PLC4X_02,PLC4X_01,PLC4X_02]]";
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "deviceaccess=[PLC4X_1, PLC4X_1, PLC4X_1]&devices=[" + String.join(",", macAddresses) + "]&submodules=" + subModules);

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setConfiguration(configuration);
        protocolLogic.setConfiguration(configuration);

        try {
            protocolLogic.setDevices();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
        Map<String, ProfinetDevice> devices = protocolLogic.getDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }

    @Test
    public void readProfinetLowerCase() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "deviceaccess=[PLC4X_1]&devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setConfiguration(configuration);
        protocolLogic.setConfiguration(configuration);

        try {
            protocolLogic.setDevices();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
        Map<String, ProfinetDevice> devices = protocolLogic.getDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "").toUpperCase()));
        }
    }

    @Test
    public void readProfinetSubModules() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        String subModules = "[[PLC4X_01, PLC4X_02, PLC4X_01, PLC4X_02]]";
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "deviceaccess=[PLC4X_1]&devices=[" + String.join(",", macAddresses) + "]&submodules=" + subModules);

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setConfiguration(configuration);
        protocolLogic.setConfiguration(configuration);

        try {
            protocolLogic.setDevices();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
        Map<String, ProfinetDevice> devices = protocolLogic.getDevices();

        for (String mac : macAddresses) {
            assertEquals("PLC4X_01", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[0]);
            assertEquals("PLC4X_02", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[1]);
            assertEquals("PLC4X_01", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[2]);
            assertEquals("PLC4X_02", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[3]);
        }
    }

    @Test
    public void readProfinetBrowseTags() throws DecoderException, PlcConnectionException, ExecutionException, InterruptedException {

        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        String subModules = "[[PLC4X_01, PLC4X_02, PLC4X_01, PLC4X_02]]";
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "deviceaccess=[PLC4X_1]&devices=[" + String.join(",", macAddresses) + "]&submodules=" + subModules);

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setConfiguration(configuration);
        protocolLogic.setConfiguration(configuration);


        try {
            protocolLogic.setDevices();
            XmlMapper xmlMapper = new XmlMapper();
            protocolLogic.getDevices().get("000C29752567").getDeviceContext().setGsdFile(xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class));
        } catch (PlcException | IOException e) {
            throw new RuntimeException(e);
        }
        CompletableFuture<PlcBrowseResponse> devices = protocolLogic.browse(new DefaultPlcBrowseRequest(null, new LinkedHashMap<>()));

        PlcBrowseResponse result = devices.get();
    }
}
