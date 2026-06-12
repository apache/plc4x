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
package org.apache.plc4x.java.modbus.tcp.config;

import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ModbusTcpConfigurationTest {

    @Test
    void testGettersAndSetters() {
        ModbusTcpConfiguration config = new ModbusTcpConfiguration();

        config.setRequestTimeout(10000);
        assertEquals(10000, config.getRequestTimeout());

        config.setDefaultUnitIdentifier(5);
        assertEquals(5, config.getDefaultUnitIdentifier());

        config.setPingAddress("4x00010:BOOL");
        assertEquals("4x00010:BOOL", config.getPingAddress());

        config.setDefaultPayloadByteOrder(ModbusByteOrder.LITTLE_ENDIAN);
        assertEquals(ModbusByteOrder.LITTLE_ENDIAN, config.getDefaultPayloadByteOrder());

        config.setMaxCoilsPerRequest(1000);
        assertEquals(1000, config.getMaxCoilsPerRequest());

        config.setMaxRegistersPerRequest(50);
        assertEquals(50, config.getMaxRegistersPerRequest());
    }

    @Test
    void testToString() {
        ModbusTcpConfiguration config = new ModbusTcpConfiguration();
        config.setRequestTimeout(5000);
        config.setDefaultUnitIdentifier(1);
        config.setPingAddress("4x00001:BOOL");
        config.setDefaultPayloadByteOrder(ModbusByteOrder.BIG_ENDIAN);
        config.setMaxCoilsPerRequest(2000);
        config.setMaxRegistersPerRequest(125);

        String result = config.toString();
        assertTrue(result.contains("requestTimeout=5000"));
        assertTrue(result.contains("unitIdentifier=1"));
        assertTrue(result.contains("pingAddress=4x00001:BOOL"));
        assertTrue(result.contains("BIG_ENDIAN"));
        assertTrue(result.contains("maxCoilsPerRequest=2000"));
        assertTrue(result.contains("maxRegistersPerRequest=125"));
    }
}
