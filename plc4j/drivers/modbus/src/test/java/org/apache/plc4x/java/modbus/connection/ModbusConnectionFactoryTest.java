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

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.serial.connection.connection.SerialChannelFactory;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import java.net.InetAddress;

import static org.junit.Assert.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

@RunWith(MockitoJUnitRunner.class)
public class ModbusConnectionFactoryTest {

    @InjectMocks
    private ModbusConnectionFactory SUT;

    @Mock
    private InetAddress inetAddress;


    @Test
    public void modbusTcpPlcConnectionOf() throws Exception {
        {
            assertThrows(NullPointerException.class, () -> SUT.modbusTcpPlcConnectionOf(null, null, null));
        }
        {
            ModbusTcpPlcConnection modbusTcpPlcConnection = SUT.modbusTcpPlcConnectionOf(inetAddress, null, null);
            assertGeneratedPort(modbusTcpPlcConnection);
        }
        {
            ModbusTcpPlcConnection modbusTcpPlcConnection = SUT.modbusTcpPlcConnectionOf(inetAddress, 13, null);
            assertEquals(inetAddress, modbusTcpPlcConnection.getRemoteAddress());
            assertPort(modbusTcpPlcConnection, 13);
        }
        {
            ModbusTcpPlcConnection modbusTcpPlcConnection = SUT.modbusTcpPlcConnectionOf(inetAddress, null, "xyz");
            assertEquals(inetAddress, modbusTcpPlcConnection.getRemoteAddress());
            assertGeneratedPort(modbusTcpPlcConnection);
        }
    }

    public void assertGeneratedPort(ModbusTcpPlcConnection modbusTcpPlcConnection) throws Exception {
        assertPort(modbusTcpPlcConnection, 502);
    }

    public void assertPort(ModbusTcpPlcConnection modbusTcpPlcConnection, int port) throws Exception {
        TcpSocketChannelFactory channelFactory = (TcpSocketChannelFactory) FieldUtils
            .getDeclaredField(NettyPlcConnection.class, "channelFactory", true)
            .get(modbusTcpPlcConnection);
        assertEquals(port, channelFactory.getPort());
    }

    @Test
    public void modbusSerialPlcConnectionOf() throws Exception {
        {
            assertThrows(NullPointerException.class, () -> SUT.modbusSerialPlcConnectionOf(null, null));
        }
        {
            ModbusSerialPlcConnection modbusSerialPlcConnection = SUT.modbusSerialPlcConnectionOf("/dev/ttyS01", null);
            assertPort(modbusSerialPlcConnection, "/dev/ttyS01");
        }
    }

    public void assertPort(ModbusSerialPlcConnection modbusSerialPlcConnection, String serialPort) throws Exception {
        SerialChannelFactory channelFactory = (SerialChannelFactory) FieldUtils
            .getDeclaredField(NettyPlcConnection.class, "channelFactory", true)
            .get(modbusSerialPlcConnection);
        assertEquals(serialPort, channelFactory.getSerialPort());
    }
}