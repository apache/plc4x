/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.java.modbus.connection;

import io.netty.channel.Channel;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.modbus.model.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import static org.junit.Assert.*;
import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;

public class ModbusTcpPlcConnectionTests {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusTcpPlcConnectionTests.class);

    private ModbusTcpPlcConnection SUT;

    private Channel channelMock;

    private ExecutorService executorService;

    @Before
    public void setUp() throws Exception {
        SUT = ModbusTcpPlcConnection.of(InetAddress.getByName("localhost"), null);
        channelMock = mock(Channel.class, RETURNS_DEEP_STUBS);
        FieldUtils.writeField(SUT, "channel", channelMock, true);
        executorService = Executors.newFixedThreadPool(10);
    }

    @After
    public void tearDown() {
        executorService.shutdownNow();
        SUT = null;
    }

    @Test
    public void emptyParseAddress() {
        try {
            SUT.parseAddress("");
        } catch (IllegalArgumentException exception) {
            assertTrue("Unexpected exception", exception.getMessage().startsWith("address  doesn't match "));
        }
    }

    @Test
    public void parseCoilModbusAddress() {
        try {
            CoilModbusAddress address = (CoilModbusAddress) SUT.parseAddress("coil:0");
            assertEquals(address.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseMaskWriteRegisterModbusAddress() {
        try {
            MaskWriteRegisterModbusAddress address = (MaskWriteRegisterModbusAddress) SUT.parseAddress("maskwrite:1/2/3");
            assertEquals(address.getAddress(), 1);
            assertEquals(address.getAndMask(), 2);
            assertEquals(address.getOrMask(), 3);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseReadDiscreteInputsModbusAddress() {
        try {
            ReadDiscreteInputsModbusAddress address = (ReadDiscreteInputsModbusAddress) SUT.parseAddress("readdiscreteinputs:0");
            assertEquals(address.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseReadHoldingRegistersModbusAddress() {
        try {
            ReadHoldingRegistersModbusAddress address = (ReadHoldingRegistersModbusAddress) SUT.parseAddress("readholdingregisters:0");
            assertEquals(address.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseReadInputRegistersModbusAddress() {
        try {
            ReadInputRegistersModbusAddress address = (ReadInputRegistersModbusAddress) SUT.parseAddress("readinputregisters:0");
            assertEquals(address.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseRegisterAddress() {
        try {
            RegisterModbusAddress address = (RegisterModbusAddress) SUT.parseAddress("register:0");
            assertEquals(address.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

}