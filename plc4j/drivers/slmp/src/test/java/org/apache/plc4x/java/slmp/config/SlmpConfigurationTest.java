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
package org.apache.plc4x.java.slmp.config;

import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class SlmpConfigurationTest {

    @Test
    void gettersAndSetters() {
        SlmpConfiguration config = new SlmpConfiguration();
        config.setMonitoringTimer(4);
        assertEquals(4, config.getMonitoringTimer());
        config.setRequestTimeout(2000);
        assertEquals(2000, config.getRequestTimeout());
    }

    @Test
    void defaultsAreApplied() {
        SlmpConfiguration config = new ConfigurationFactory()
            .createConfiguration(SlmpConfiguration.class, "");
        assertEquals(0x0000, config.getMonitoringTimer());
        assertEquals(5_000, config.getRequestTimeout());
    }

    @Test
    void overridesAreParsed() {
        SlmpConfiguration config = new ConfigurationFactory()
            .createConfiguration(SlmpConfiguration.class, "monitoring-timer=4&request-timeout-ms=2000");
        assertEquals(4, config.getMonitoringTimer());
        assertEquals(2_000, config.getRequestTimeout());
    }
}
