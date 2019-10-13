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
package org.apache.plc4x.java.ads.connection;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.serial.connection.connection.SerialChannelFactory;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.junit.Assert;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import java.net.InetAddress;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

@RunWith(MockitoJUnitRunner.class)
public class AdsConnectionFactoryTest {

    @InjectMocks
    private AdsConnectionFactory SUT;

    @Mock
    private InetAddress inetAddress;

    @Mock
    private AmsNetId targetAmsNetId;

    @Mock
    private AmsPort targetAmsPort;

    @Mock
    private AmsNetId sourceAmsNetId;

    @Mock
    private AmsPort sourceAmsPort;

    @Test
    public void adsTcpPlcConnectionOf() throws Exception {
        {
            assertThrows(NullPointerException.class, () -> SUT.adsTcpPlcConnectionOf(null, null, null, null, null, null));
            assertThrows(NullPointerException.class, () -> SUT.adsTcpPlcConnectionOf(inetAddress, null, null, null, null, null));
            assertThrows(NullPointerException.class, () -> SUT.adsTcpPlcConnectionOf(inetAddress, null, targetAmsNetId, null, null, null));
        }
        {
            AdsTcpPlcConnection adsTcpPlcConnection = SUT.adsTcpPlcConnectionOf(inetAddress, null, targetAmsNetId, targetAmsPort, null, null);
            assertEquals(inetAddress, adsTcpPlcConnection.getRemoteAddress());
            assertEquals(targetAmsNetId, adsTcpPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsTcpPlcConnection.getTargetAmsPort());
            assertNotEquals(sourceAmsNetId, adsTcpPlcConnection.getSourceAmsNetId());
            assertNotEquals(sourceAmsPort, adsTcpPlcConnection.getSourceAmsPort());
            assertGeneratedPort(adsTcpPlcConnection);
        }
        {
            AdsTcpPlcConnection adsTcpPlcConnection = SUT.adsTcpPlcConnectionOf(inetAddress, 13, targetAmsNetId, targetAmsPort, null, null);
            assertEquals(inetAddress, adsTcpPlcConnection.getRemoteAddress());
            assertEquals(targetAmsNetId, adsTcpPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsTcpPlcConnection.getTargetAmsPort());
            assertNotEquals(sourceAmsNetId, adsTcpPlcConnection.getSourceAmsNetId());
            assertNotEquals(sourceAmsPort, adsTcpPlcConnection.getSourceAmsPort());
            assertPort(adsTcpPlcConnection, 13);
        }
        {
            AdsTcpPlcConnection adsTcpPlcConnection = SUT.adsTcpPlcConnectionOf(inetAddress, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, null);
            assertEquals(inetAddress, adsTcpPlcConnection.getRemoteAddress());
            assertEquals(targetAmsNetId, adsTcpPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsTcpPlcConnection.getTargetAmsPort());
            assertNotEquals(sourceAmsNetId, adsTcpPlcConnection.getSourceAmsNetId());
            assertNotEquals(sourceAmsPort, adsTcpPlcConnection.getSourceAmsPort());
            assertGeneratedPort(adsTcpPlcConnection);
        }
        {
            AdsTcpPlcConnection adsTcpPlcConnection = SUT.adsTcpPlcConnectionOf(inetAddress, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
            assertEquals(inetAddress, adsTcpPlcConnection.getRemoteAddress());
            assertEquals(targetAmsNetId, adsTcpPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsTcpPlcConnection.getTargetAmsPort());
            assertEquals(sourceAmsNetId, adsTcpPlcConnection.getSourceAmsNetId());
            assertEquals(sourceAmsPort, adsTcpPlcConnection.getSourceAmsPort());
            assertGeneratedPort(adsTcpPlcConnection);
        }
        {
            AdsTcpPlcConnection adsTcpPlcConnection = SUT.adsTcpPlcConnectionOf(inetAddress, 13, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
            assertEquals(inetAddress, adsTcpPlcConnection.getRemoteAddress());
            assertEquals(targetAmsNetId, adsTcpPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsTcpPlcConnection.getTargetAmsPort());
            assertEquals(sourceAmsNetId, adsTcpPlcConnection.getSourceAmsNetId());
            assertEquals(sourceAmsPort, adsTcpPlcConnection.getSourceAmsPort());
            assertPort(adsTcpPlcConnection, 13);
        }
    }

    public void assertGeneratedPort(AdsTcpPlcConnection adsTcpPlcConnection) throws Exception {
        assertPort(adsTcpPlcConnection, 48898);
    }

    public void assertPort(AdsTcpPlcConnection adsTcpPlcConnection, int port) throws Exception {
        TcpSocketChannelFactory channelFactory = (TcpSocketChannelFactory) FieldUtils
            .getDeclaredField(NettyPlcConnection.class, "channelFactory", true)
            .get(adsTcpPlcConnection);
        Assert.assertEquals(port, channelFactory.getPort());
    }

    @Test
    public void adsSerialPlcConnectionOf() throws Exception {
        {
            assertThrows(NullPointerException.class, () -> SUT.adsSerialPlcConnectionOf(null, null, null, null, null));
            assertThrows(NullPointerException.class, () -> SUT.adsSerialPlcConnectionOf("/dev/ttyS01", null, null, null, null));
            assertThrows(NullPointerException.class, () -> SUT.adsSerialPlcConnectionOf("/dev/ttyS01", targetAmsNetId, null, null, null));
        }
        {
            AdsSerialPlcConnection adsSerialPlcConnection = SUT.adsSerialPlcConnectionOf("/dev/ttyS01", targetAmsNetId, targetAmsPort, null, null);
            assertEquals(targetAmsNetId, adsSerialPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsSerialPlcConnection.getTargetAmsPort());
            assertNotEquals(sourceAmsNetId, adsSerialPlcConnection.getSourceAmsNetId());
            assertNotEquals(sourceAmsPort, adsSerialPlcConnection.getSourceAmsPort());
            assertPort(adsSerialPlcConnection, "/dev/ttyS01");
        }
        {
            AdsSerialPlcConnection adsSerialPlcConnection = SUT.adsSerialPlcConnectionOf("/dev/ttyS01", targetAmsNetId, targetAmsPort, sourceAmsNetId, null);
            assertEquals(targetAmsNetId, adsSerialPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsSerialPlcConnection.getTargetAmsPort());
            assertNotEquals(sourceAmsNetId, adsSerialPlcConnection.getSourceAmsNetId());
            assertNotEquals(sourceAmsPort, adsSerialPlcConnection.getSourceAmsPort());
            assertPort(adsSerialPlcConnection, "/dev/ttyS01");
        }
        {
            AdsSerialPlcConnection adsSerialPlcConnection = SUT.adsSerialPlcConnectionOf("/dev/ttyS01", targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
            assertEquals(targetAmsNetId, adsSerialPlcConnection.getTargetAmsNetId());
            assertEquals(targetAmsPort, adsSerialPlcConnection.getTargetAmsPort());
            assertEquals(sourceAmsNetId, adsSerialPlcConnection.getSourceAmsNetId());
            assertEquals(sourceAmsPort, adsSerialPlcConnection.getSourceAmsPort());
            assertPort(adsSerialPlcConnection, "/dev/ttyS01");
        }
    }

    public void assertPort(AdsSerialPlcConnection adsSerialPlcConnection, String serialPort) throws Exception {
        SerialChannelFactory channelFactory = (SerialChannelFactory) FieldUtils
            .getDeclaredField(NettyPlcConnection.class, "channelFactory", true)
            .get(adsSerialPlcConnection);
        Assert.assertEquals(serialPort, channelFactory.getSerialPort());
    }
}