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
package org.apache.plc4x.java.ads.netty;

import org.apache.commons.io.HexDump;
import org.apache.commons.lang3.SerializationUtils;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.plc4x.java.ads.api.commands.ADSReadResponse;
import org.apache.plc4x.java.ads.api.commands.ADSWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPacket;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.ADSAddress;
import org.apache.plc4x.java.api.messages.*;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.ArrayList;
import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.hasSize;
import static org.hamcrest.Matchers.instanceOf;
import static org.junit.Assert.assertThat;

@RunWith(Parameterized.class)
public class Plc4XADSProtocolTest {
    private static final Logger LOGGER = LoggerFactory.getLogger(ADSProtocolTest.class);

    private Plc4XADSProtocol SUT;

    @Parameterized.Parameter()
    public PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer;

    @Parameterized.Parameter(1)
    public CompletableFuture completableFuture;

    @Parameterized.Parameter(2)
    public AMSTCPPacket amstcpPacket;

    @Parameterized.Parameters(name = "{2} {index}")
    public static Collection<Object[]> data() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(2);
        Data data = Data.of(SerializationUtils.serialize("Hello World!"));
        return Stream.of(
            ImmutablePair.of(
                new PlcRequestContainer<>(
                    PlcWriteRequest
                        .builder()
                        .addItem(ADSAddress.of(1, 2), "HelloWorld!")
                        .build(), new CompletableFuture<>()),
                ADSWriteResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0))
            ),
            ImmutablePair.of(
                new PlcRequestContainer<>(
                    PlcReadRequest
                        .builder()
                        .addItem(String.class, ADSAddress.of(1, 2))
                        .build(), new CompletableFuture<>()),
                ADSReadResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0), data)
            )
        ).map(pair -> new Object[]{pair.left, pair.left.getResponseFuture(), pair.right}).collect(Collectors.toList());
    }

    @Before
    public void setUp() throws Exception {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        SUT = new Plc4XADSProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    @Test
    public void encode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, out);
        assertThat(out, hasSize(1));
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            byte[] bytes = ((AMSTCPPacket) out.get(0)).getByteBuf().array();
            HexDump.dump(bytes, 0, byteArrayOutputStream, 0);
            byteArrayOutputStream.flush();
            LOGGER.info("{}\nHexDump:\n{}", amstcpPacket, byteArrayOutputStream);
        }
    }

    @Test
    public void decode() throws Exception {
        ArrayList<Object> in = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, in);
        assertThat(in, hasSize(1));
        syncInvoiceId();
        ArrayList<Object> out = new ArrayList<>();
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            byte[] bytes = amstcpPacket.getByteBuf().array();
            HexDump.dump(bytes, 0, byteArrayOutputStream, 0);
            LOGGER.info("{}\nHexDump:\n{}", amstcpPacket, byteArrayOutputStream);
        }
        SUT.decode(null, amstcpPacket, out);
        assertThat(out, hasSize(1));
        assertThat(out.get(0), instanceOf(PlcRequestContainer.class));
    }

    private void syncInvoiceId() throws Exception {
        Field correlationBuilderField = SUT.getClass().getDeclaredField("correlationBuilder");
        correlationBuilderField.setAccessible(true);
        AtomicLong correlationBuilder = (AtomicLong) correlationBuilderField.get(SUT);

        AMSHeader amsHeader = amstcpPacket.getAmsHeader();
        Field invokeIdField = amsHeader.getClass().getDeclaredField("invokeId");
        Field modifiersField = Field.class.getDeclaredField("modifiers");
        modifiersField.setAccessible(true);
        modifiersField.setInt(invokeIdField, invokeIdField.getModifiers() & ~Modifier.FINAL);
        invokeIdField.setAccessible(true);
        invokeIdField.set(amsHeader, Invoke.of(correlationBuilder.get()));
    }
}