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

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetBrowseTests {

    @Test
    public void readProfinetBrowseTagsCheckStatus()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name, PLC4X_1, (PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, )]]");

        ProfinetDevice device = configuration.getDevices().getConfiguredDevices().get("DEVICE_NAME");
        device.setVendorDeviceId("CAFE", "0001");

        Map<String, List<PlcBrowseItem>> browseItems = new HashMap<>();
        Map<String, List<PlcBrowseItem>> response = device.browseTags(browseItems);

        assert(response.containsKey("DEVICE_NAME.0.1.PLC4X_1_V0.Status"));
        assert(response.containsKey("DEVICE_NAME.0.32768.PLC4X_1_S0.Status"));
        assert(response.containsKey("DEVICE_NAME.0.32769.PLC4X_1_S1.Status"));
        assert(response.containsKey("DEVICE_NAME.1.1.PLC4X_DUMMY_MODULE_V0.Status"));
        assert(response.containsKey("DEVICE_NAME.3.1.PLC4X_DUMMY_MODULE_V0.Status"));
    }

    @Test
    public void readProfinetBrowseTagsCheckFloat()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name, PLC4X_1, (PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, )]]");

        ProfinetDevice device = configuration.getDevices().getConfiguredDevices().get("DEVICE_NAME");
        device.setVendorDeviceId("CAFE", "0001");

        Map<String, List<PlcBrowseItem>> browseItems = new HashMap<>();
        Map<String, List<PlcBrowseItem>> response = device.browseTags(browseItems);

        String key = "DEVICE_NAME.3.1.PLC4X_INPUT_MODULE_INFO_FLOAT";
        assert(response.containsKey(key));
        assertEquals(response.get(key).get(0).getTag().getAddressString(), key);
        assertEquals(response.get(key).get(0).getTag().getPlcValueType(), PlcValueType.REAL);
    }

    @Test
    public void readProfinetBrowseTagsCheckBoolean()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name, PLC4X_1, (PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, )]]");

        ProfinetDevice device = configuration.getDevices().getConfiguredDevices().get("DEVICE_NAME");
        device.setVendorDeviceId("CAFE", "0001");

        Map<String, List<PlcBrowseItem>> browseItems = new HashMap<>();
        Map<String, List<PlcBrowseItem>> response = device.browseTags(browseItems);

        String key = "DEVICE_NAME.1.1.PLC4X_INPUT_MODULE_INFO_32.1";
        assert(response.containsKey(key));
        assertEquals(response.get(key).get(0).getTag().getAddressString(), key);
        assertEquals(response.get(key).get(0).getTag().getPlcValueType(), PlcValueType.BOOL);

        key = "DEVICE_NAME.1.1.PLC4X_INPUT_MODULE_INFO_32.7";
        assertEquals(response.get(key).get(0).getTag().getPlcValueType(), PlcValueType.BOOL);
    }

}
