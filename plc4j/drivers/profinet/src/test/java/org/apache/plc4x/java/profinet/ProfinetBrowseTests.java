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

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDeviceContext;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseRequest;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.io.File;
import java.io.IOException;
import java.util.LinkedHashMap;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetBrowseTests {

    @Test
    public void readProfinetBrowseTags() throws ExecutionException, InterruptedException {
        String[] macAddresses = new String[] {"00:0c:29:75:25:67"};
        String subModules = "[[PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, ]]";
        ProfinetProtocolLogic protocolLogic = new ProfinetProtocolLogic();
        ProfinetConfiguration configuration = (ProfinetConfiguration) new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name_1, PLC4X_1, {PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, }]]");

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setConfiguration(configuration);
        protocolLogic.setConfiguration(configuration);

        try {
            XmlMapper xmlMapper = new XmlMapper();
            ProfinetDevice device = configuration.getDevices().getConfiguredDevices().get("DEVICE_NAME_1");
            device.setDriverContext((ProfinetDriverContext) protocolLogic.getDriverContext());
            device.getDeviceContext().setDeviceName("PLC4X Test Module");
            device.getDeviceContext().setGsdFile(xmlMapper.readValue(new File("src/test/resources/gsdml.xml"), ProfinetISO15745Profile.class));
            device.getDeviceContext().setVendorId("CAFE");
            device.getDeviceContext().setDeviceId("0001");
            device.setSubModulesObjects();
        } catch (PlcException | IOException e) {
            throw new RuntimeException(e);
        }
        CompletableFuture<PlcBrowseResponse> devices = protocolLogic.browse(new DefaultPlcBrowseRequest(null, new LinkedHashMap<>()));

        PlcBrowseResponse result = devices.get();
    }

}
