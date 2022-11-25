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

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.Map;

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
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=" + directory);
        Map<String, ProfinetISO15745Profile> gsdFiles = configuration.readGsdFiles();
        assertEquals(directory, configuration.getGsdDirectory());
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
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();


        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }

    @Test
    public void readProfinetDevicesMultiple() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"CA:FE:00:00:00:01","CA:FE:00:00:00:02","CA:FE:00:00:00:03"};
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "")));
        }
    }

    @Test
    public void readProfinetLowerCase() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]");

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();

        for (String mac : macAddresses) {
            assert(devices.containsKey(mac.replace(":", "").toUpperCase()));
        }
    }

    @Test
    public void readProfinetSubModules() throws DecoderException, PlcConnectionException {

        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        String subModules = "[[PLC4X_01, PLC4X_02, PLC4X_01, PLC4X_02]]";
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "devices=[" + String.join(",", macAddresses) + "]&submodules=" + subModules);

        ProfinetDriverContext context = new ProfinetDriverContext();
        context.setConfiguration(configuration);

        Map<String, ProfinetDevice> devices = configuration.getConfiguredDevices();
        configuration.setSubModules();

        for (String mac : macAddresses) {
            String[] test = devices.get(mac.replace(":", "").toUpperCase()).getSubModules();
            assertEquals("PLC4X_01", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[0]);
            assertEquals("PLC4X_02", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[1]);
            assertEquals("PLC4X_01", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[2]);
            assertEquals("PLC4X_02", devices.get(mac.replace(":", "").toUpperCase()).getSubModules()[3]);
        }
    }


}
