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

import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import org.apache.plc4x.java.ads.api.commands.AdsReadWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.model.SymbolicAdsAddress;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.junit.Before;
import org.junit.Test;

import java.util.concurrent.CompletableFuture;

import static org.junit.Assert.assertNotNull;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

public class AdsAbstractPlcConnectionTest {

    private AdsAbstractPlcConnection SUT;

    @Before
    public void setUp() throws Exception {
        ChannelFactory channelFactory = mock(ChannelFactory.class, RETURNS_DEEP_STUBS);
        SUT = new AdsAbstractPlcConnection(channelFactory, mock(AmsNetId.class), mock(AmsPort.class), mock(AmsNetId.class), mock(AmsPort.class)) {
            @Override
            protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
                return null;
            }
        };

        // Specific to mapAddress
        Channel channel = mock(Channel.class, RETURNS_DEEP_STUBS);
        when(channelFactory.createChannel(any())).thenReturn(channel);
        when(channel.writeAndFlush(any(PlcRequestContainer.class))).then(invocation -> {
            PlcRequestContainer plcRequestContainer = invocation.getArgument(0);
            PlcProprietaryResponse plcProprietaryResponse = mock(PlcProprietaryResponse.class, RETURNS_DEEP_STUBS);
            AdsReadWriteResponse adsReadWriteResponse = mock(AdsReadWriteResponse.class, RETURNS_DEEP_STUBS);
            when(adsReadWriteResponse.getResult()).thenReturn(Result.of(0));
            when(adsReadWriteResponse.getData()).thenReturn(Data.of(new byte[]{1, 2, 3, 4}));
            when(plcProprietaryResponse.getResponse()).thenReturn(adsReadWriteResponse);
            plcRequestContainer.getResponseFuture().complete(plcProprietaryResponse);
            return mock(ChannelFuture.class);
        });

        SUT.connect();
    }

    @Test
    public void getTargetAmsNetId() {
        AmsNetId targetAmsNetId = SUT.getTargetAmsNetId();
        assertNotNull(targetAmsNetId);
    }

    @Test
    public void getTargetAmsPort() {
        AmsPort targetAmsPort = SUT.getTargetAmsPort();
        assertNotNull(targetAmsPort);
    }

    @Test
    public void getSourceAmsNetId() {
        AmsNetId sourceAmsNetId = SUT.getSourceAmsNetId();
        assertNotNull(sourceAmsNetId);
    }

    @Test
    public void getSourceAmsPort() {
        AmsPort sourceAmsPort = SUT.getSourceAmsPort();
        assertNotNull(sourceAmsPort);
    }

    @Test
    public void parseAddress() {
        Address address = SUT.parseAddress("0/0");
        assertNotNull(address);
        Address SymbolicAddress = SUT.parseAddress("Main.byByte[0]");
        assertNotNull(SymbolicAddress);
    }

    @Test
    public void read() {
        CompletableFuture<PlcReadResponse> read = SUT.read(mock(PlcReadRequest.class));
        assertNotNull(read);
        CompletableFuture<TypeSafePlcReadResponse<Object>> typeSafeRead = SUT.read(mock(TypeSafePlcReadRequest.class));
        assertNotNull(typeSafeRead);
    }

    @Test
    public void write() {
        CompletableFuture<PlcWriteResponse> write = SUT.write(mock(PlcWriteRequest.class));
        assertNotNull(write);
        CompletableFuture<TypeSafePlcWriteResponse<Object>> typeSafeWrite = SUT.write(mock(TypeSafePlcWriteRequest.class));
        assertNotNull(typeSafeWrite);
    }

    @Test
    public void send() {
        CompletableFuture send = SUT.send(mock(PlcProprietaryRequest.class));
        assertNotNull(send);
    }

    @Test
    public void mapAddresses() {
        SUT.mapAddresses(mock(PlcRequest.class));
    }

    @Test
    public void mapAddress() {
        SUT.mapAddress(SymbolicAdsAddress.of("Main.byByte[0]"));
    }

    @Test
    public void generateAMSNetId() {
        AmsNetId targetAmsNetId = SUT.getTargetAmsNetId();
        assertNotNull(targetAmsNetId);
    }

    @Test
    public void generateAMSPort() {
        AmsPort amsPort = AdsAbstractPlcConnection.generateAMSPort();
        assertNotNull(amsPort);
    }

    @Test
    public void close() {
        SUT.close();
    }

    @Test
    public void getFromFuture() {
        Object fromFuture = SUT.getFromFuture(mock(CompletableFuture.class, RETURNS_DEEP_STUBS), 1);
        assertNotNull(fromFuture);
    }

    @Test
    public void testToString() {
        String s = SUT.toString();
        assertNotNull(s);
    }
}