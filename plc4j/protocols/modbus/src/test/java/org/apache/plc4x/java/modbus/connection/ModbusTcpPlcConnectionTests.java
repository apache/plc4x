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
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.modbus.model.*;
import org.hamcrest.Matchers;
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
    public void prepareEmptyField() {
        try {
            SUT.prepareField("");
        } catch (PlcInvalidFieldException exception) {
            assertThat(exception.getMessage(), Matchers.startsWith(" invalid"));
        }
    }

    @Test
    public void prepareCoilModbusField() throws Exception {
        try {
            CoilModbusField field = (CoilModbusField) SUT.prepareField("coil:0");
            assertEquals(field.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void prepareMaskWriteRegisterModbusField() throws Exception {
        try {
            MaskWriteRegisterModbusField field = (MaskWriteRegisterModbusField) SUT.prepareField("maskwrite:1/2/3");
            assertEquals(field.getAddress(), 1);
            assertEquals(field.getAndMask(), 2);
            assertEquals(field.getOrMask(), 3);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void prepareReadDiscreteInputsModbusField() throws Exception {
        try {
            ReadDiscreteInputsModbusField field = (ReadDiscreteInputsModbusField) SUT.prepareField("readdiscreteinputs:0");
            assertEquals(field.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void prepareReadHoldingRegistersModbusField() throws Exception {
        try {
            ReadHoldingRegistersModbusField field = (ReadHoldingRegistersModbusField) SUT.prepareField("readholdingregisters:0");
            assertEquals(field.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void prepareReadInputRegistersModbusField() throws Exception {
        try {
            ReadInputRegistersModbusField field = (ReadInputRegistersModbusField) SUT.prepareField("readinputregisters:0");
            assertEquals(field.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void prepareRegisterField() throws Exception {
        try {
            RegisterModbusField field = (RegisterModbusField) SUT.prepareField("register:0");
            assertEquals(field.getAddress(), 0);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }

}