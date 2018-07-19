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
package org.apache.plc4x.java.ads.protocol;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.plc4x.java.ads.api.commands.AdsReadResponse;
import org.apache.plc4x.java.ads.api.commands.AdsWriteRequest;
import org.apache.plc4x.java.ads.api.commands.AdsWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.AdsAddress;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.messages.items.ResponseItem;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.ArrayList;
import java.util.Collection;
import java.util.GregorianCalendar;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes.defaultAssert;
import static org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes.streamOfLittleEndianDataTypePairs;
import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertThat;

@RunWith(Parameterized.class)
public class Plc4x2AdsProtocolTest {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4x2AdsProtocolTest.class);

    private Plc4x2AdsProtocol SUT;

    @Parameterized.Parameter
    public String payloadClazzName;

    @Parameterized.Parameter(1)
    public PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer;

    @Parameterized.Parameter(2)
    public CompletableFuture completableFuture;

    @Parameterized.Parameter(3)
    public String plcRequestContainerClassName;

    @Parameterized.Parameter(4)
    public AmsPacket amsPacket;

    @Parameterized.Parameter(5)
    public String amsPacketClassName;

    @Parameterized.Parameters(name = "{index} Type:{0} {3} {5}")
    public static Collection<Object[]> data() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        Invoke invokeId = Invoke.of(2);
        return streamOfLittleEndianDataTypePairs()
            .map(pair -> Stream.of(
                ImmutablePair.of(
                    new PlcRequestContainer<>(
                        PlcWriteRequest
                            .builder()
                            .addItem(AdsAddress.of(1, 2), pair.getLeft())
                            .build(), new CompletableFuture<>()),
                    AdsWriteResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0))
                ),
                ImmutablePair.of(
                    new PlcRequestContainer<>(
                        PlcReadRequest
                            .builder()
                            .addItem(pair.getLeft().getClass(), AdsAddress.of(1, 2))
                            .build(), new CompletableFuture<>()),
                    AdsReadResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0), Data.of(pair.getRight()))
                )
            ))
            .flatMap(stream -> stream)
            .map(pair -> new Object[]{pair.left.getRequest().getRequestItem().orElseThrow(IllegalStateException::new).getDatatype().getSimpleName(), pair.left, pair.left.getResponseFuture(), pair.left.getRequest().getClass().getSimpleName(), pair.right, pair.right.getClass().getSimpleName()}).collect(Collectors.toList());
    }

    @Before
    public void setUp() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        SUT = new Plc4x2AdsProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, new ConcurrentHashMap<>());
    }

    @Test
    public void encode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, out);
        assertThat(out, hasSize(1));
        assertThat(out.get(0), instanceOf(AmsPacket.class));
        AmsPacket amsPacket = (AmsPacket) out.get(0);
        LOGGER.info("{}\nHexDump:\n{}", amsPacket, amsPacket.dump());
        if (amsPacket instanceof AdsWriteRequest) {
            AdsWriteRequest adsWriteRequest = (AdsWriteRequest) amsPacket;
            byte[] value = adsWriteRequest.getData().getBytes();
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1}));
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1}));
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1, 0x0}));
            } else if (payloadClazzName.equals(GregorianCalendar.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, (byte) 0x80, 0x3E, 0x15, (byte) 0xAB, 0X47, (byte) 0xFC, 0x28}));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, 0x0, (byte) 0x80, 0x3F}));
            } else if (payloadClazzName.equals(Double.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, (byte) 0xF0, 0x3F}));
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1, 0x0, 0x0, 0x0}));
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            }
        }
    }

    @Test
    public void decode() throws Exception {
        ArrayList<Object> in = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, in);
        assertThat(in, hasSize(1));
        syncInvoiceId();
        ArrayList<Object> out = new ArrayList<>();
        LOGGER.info("{}\nHexDump:\n{}", amsPacket, amsPacket.dump());
        SUT.decode(null, amsPacket, out);
        assertThat(out, hasSize(0));
        LOGGER.info("PlcRequestContainer {}", plcRequestContainer);
        PlcResponse plcResponse = plcRequestContainer.getResponseFuture().get();
        ResponseItem responseItem = (ResponseItem) plcResponse.getResponseItem().get();
        LOGGER.info("ResponseItem {}", responseItem);
        if (amsPacket instanceof AdsReadResponse) {
            ReadResponseItem readResponseItem = (ReadResponseItem) responseItem;
            Object value = readResponseItem.getValues().get(0);
            defaultAssert(value);
        }
    }

    private void syncInvoiceId() throws Exception {
        Field correlationBuilderField = SUT.getClass().getDeclaredField("correlationBuilder");
        correlationBuilderField.setAccessible(true);
        AtomicLong correlationBuilder = (AtomicLong) correlationBuilderField.get(SUT);

        AmsHeader amsHeader = amsPacket.getAmsHeader();
        Field invokeIdField = amsHeader.getClass().getDeclaredField("invokeId");
        Field modifiersField = Field.class.getDeclaredField("modifiers");
        modifiersField.setAccessible(true);
        modifiersField.setInt(invokeIdField, invokeIdField.getModifiers() & ~Modifier.FINAL);
        invokeIdField.setAccessible(true);
        invokeIdField.set(amsHeader, Invoke.of(correlationBuilder.get()));
    }
}