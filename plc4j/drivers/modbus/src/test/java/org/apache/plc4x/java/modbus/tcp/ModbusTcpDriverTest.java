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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.readwrite.Constants;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class ModbusTcpDriverTest {

    private ModbusTcpDriver driver;

    @BeforeEach
    void setUp() {
        driver = new ModbusTcpDriver();
    }

    @Test
    void testGetProtocolCode() {
        assertEquals("modbus-tcp", driver.getProtocolCode());
    }

    @Test
    void testGetProtocolName() {
        assertEquals("Modbus TCP", driver.getProtocolName());
    }

    @Test
    void testGetMetadata() {
        var metadata = driver.getMetadata();
        assertNotNull(metadata);
        assertTrue(metadata.getDefaultTransportCode().isPresent());
        assertEquals("tcp", metadata.getDefaultTransportCode().get());
        assertTrue(metadata.getSupportedTransportCodes().contains("tcp"));
        assertTrue(metadata.getSupportedTransportCodes().contains("tls"));
        assertTrue(metadata.getSupportedTransportCodes().contains("tls-psk"));
        assertTrue(metadata.getSupportedTransportCodes().contains("test"));
        assertTrue(metadata.isDiscoverySupported());
    }

    @Test
    void testDefaultPorts() {
        Set<Integer> tcpPorts = driver.defaultPorts("tcp");
        assertTrue(tcpPorts.contains(Constants.MODBUSTCPDEFAULTPORT));

        Set<Integer> tlsPorts = driver.defaultPorts("tls");
        assertTrue(tlsPorts.contains(Constants.MODBUSTCPTLSDEFAULTPORT));

        Set<Integer> tlsPskPorts = driver.defaultPorts("tls-psk");
        assertTrue(tlsPskPorts.contains(Constants.MODBUSTCPTLSDEFAULTPORT));

        Set<Integer> unknownPorts = driver.defaultPorts("serial");
        assertTrue(unknownPorts.isEmpty());
    }

    @Test
    void testPrepareTag() {
        ModbusTag tag = driver.prepareTag("4x00001:INT");
        assertNotNull(tag);
        assertInstanceOf(ModbusTagHoldingRegister.class, tag);
    }

    @Test
    void testDiscoveryRequestBuilder() {
        assertNotNull(driver.discoveryRequestBuilder());
    }
}
