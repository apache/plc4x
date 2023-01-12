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
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetBrowseTests {

    @Test
    public void readProfinetBrowseTags()  {
        ProfinetConfiguration configuration = new ConfigurationFactory().createConfiguration(
            ProfinetConfiguration.class, "gsddirectory=src/test/resources/&devices=[[device_name_1, PLC4X_1, {PLC4X_DUMMY_MODULE, , PLC4X_DUMMY_MODULE, }]]");

        ProfinetDevice device = configuration.getDevices().getConfiguredDevices().get("DEVICE_NAME_1");
        device.setVendorDeviceId("CAFE", "0001");

        PlcBrowseItem response = device.browseTags();
        response.getTag();
    }

}
