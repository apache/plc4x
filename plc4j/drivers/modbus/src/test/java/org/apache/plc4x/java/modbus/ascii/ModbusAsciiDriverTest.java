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
package org.apache.plc4x.java.modbus.ascii;

import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagInputRegister;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ModbusAsciiDriverTest {

    private ModbusAsciiDriver driver;

    @BeforeEach
    void setUp() {
        driver = new ModbusAsciiDriver();
    }

    @Test
    void testGetProtocolCode() {
        assertEquals("modbus-ascii", driver.getProtocolCode());
    }

    @Test
    void testGetProtocolName() {
        assertEquals("Modbus ASCII", driver.getProtocolName());
    }

    @Test
    void testGetMetadata() {
        var metadata = driver.getMetadata();
        assertNotNull(metadata);
        assertTrue(metadata.getDefaultTransportCode().isPresent());
        assertEquals("serial", metadata.getDefaultTransportCode().get());
        assertTrue(metadata.getSupportedTransportCodes().contains("tcp"));
        assertTrue(metadata.getSupportedTransportCodes().contains("serial"));
        assertTrue(metadata.getSupportedTransportCodes().contains("test"));
    }

    @Test
    void testPrepareTag() {
        ModbusTag tag = driver.prepareTag("3x00001:INT");
        assertNotNull(tag);
        assertInstanceOf(ModbusTagInputRegister.class, tag);
    }
}
