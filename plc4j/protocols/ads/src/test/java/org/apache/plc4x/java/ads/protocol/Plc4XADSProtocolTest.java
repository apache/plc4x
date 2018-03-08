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
import org.apache.plc4x.java.ads.api.commands.ADSReadResponse;
import org.apache.plc4x.java.ads.api.commands.ADSWriteRequest;
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
import java.util.Calendar;
import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertThat;

@RunWith(Parameterized.class)
public class Plc4XADSProtocolTest {
    private static final Logger LOGGER = LoggerFactory.getLogger(ADSProtocolTest.class);

    public static final Calendar calenderInstance = Calendar.getInstance();

    private Plc4XADSProtocol SUT;

    @Parameterized.Parameter
    public String payloadClazzName;

    @Parameterized.Parameter(1)
    public PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer;

    @Parameterized.Parameter(2)
    public CompletableFuture completableFuture;

    @Parameterized.Parameter(3)
    public String plcRequestContainerClassName;

    @Parameterized.Parameter(4)
    public AMSTCPPacket amstcpPacket;

    @Parameterized.Parameter(5)
    public String aMSTCPPacketClassName;

    @Parameterized.Parameters(name = "{index} Type:{0} {3} {5}")
    public static Collection<Object[]> data() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(2);
        return Stream.of(Boolean.class,
            Byte.class,
            Short.class,
            Calendar.class,
            Float.class,
            Integer.class,
            String.class)
            .map(clazz -> {
                if (clazz == Boolean.class) {
                    return ImmutablePair.of(Boolean.TRUE, new byte[]{0x01});
                } else if (clazz == Byte.class) {
                    return ImmutablePair.of(Byte.valueOf("1"), new byte[]{0x1});
                } else if (clazz == Short.class) {
                    return ImmutablePair.of(Short.valueOf("1"), new byte[]{0x1, 0x0});
                } else if (clazz == Calendar.class) {
                    return ImmutablePair.of(calenderInstance, new byte[]{0x0, 0x0, 0x0, 0x0, 0x4, 0x3, 0x2, 0x1});
                } else if (clazz == Float.class) {
                    return ImmutablePair.of(Float.valueOf("1"), new byte[]{0x0, 0x0, (byte) 0x80, 0x3F});
                } else if (clazz == Integer.class) {
                    return ImmutablePair.of(Integer.valueOf("1"), new byte[]{0x1, 0x0, 0x0, 0x0});
                } else if (clazz == String.class) {
                    return ImmutablePair.of(String.valueOf("Hello World!"), new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00});
                } else {
                    throw new IllegalArgumentException("Unmapped type " + clazz);
                }
            })
            .map(pair -> Stream.of(
                ImmutablePair.of(
                    new PlcRequestContainer<>(
                        PlcWriteRequest
                            .builder()
                            .addItem(ADSAddress.of(1, 2), pair.left)
                            .build(), new CompletableFuture<>()),
                    ADSWriteResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0))
                ),
                ImmutablePair.of(
                    new PlcRequestContainer<>(
                        PlcReadRequest
                            .builder()
                            .addItem(pair.left.getClass(), ADSAddress.of(1, 2))
                            .build(), new CompletableFuture<>()),
                    ADSReadResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0), Data.of(pair.right))
                )
            ))
            .flatMap(stream -> stream)
            .map(pair -> new Object[]{pair.left.getRequest().getRequestItem().orElseThrow(IllegalStateException::new).getDatatype().getSimpleName(), pair.left, pair.left.getResponseFuture(), pair.left.getRequest().getClass().getSimpleName(), pair.right, pair.right.getClass().getSimpleName()}).collect(Collectors.toList());
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
        assertThat(out.get(0), instanceOf(AMSTCPPacket.class));
        AMSTCPPacket amstcpPacket = (AMSTCPPacket) out.get(0);
        LOGGER.info("{}\nHexDump:\n{}", amstcpPacket, amstcpPacket.dump());
        if (amstcpPacket instanceof ADSWriteRequest) {
            ADSWriteRequest adsWriteRequest = (ADSWriteRequest) amstcpPacket;
            byte[] value = adsWriteRequest.getData().getBytes();
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1}));
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1}));
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1, 0x0}));
            } else if (payloadClazzName.equals(Calendar.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0}));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, 0x0, (byte) 0x80, 0x3F}));
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
        LOGGER.info("{}\nHexDump:\n{}", amstcpPacket, amstcpPacket.dump());
        SUT.decode(null, amstcpPacket, out);
        assertThat(out, hasSize(1));
        assertThat(out.get(0), instanceOf(PlcRequestContainer.class));
        PlcRequestContainer<?, ?> plcRequestContainer = (PlcRequestContainer) out.get(0);
        LOGGER.info("PlcRequestContainer {}", plcRequestContainer);
        PlcResponse plcResponse = plcRequestContainer.getResponseFuture().get();
        ResponseItem responseItem = (ResponseItem) plcResponse.getResponseItem().get();
        LOGGER.info("ResponseItem {}", responseItem);
        if (amstcpPacket instanceof ADSReadResponse) {
            ReadResponseItem readResponseItem = (ReadResponseItem) responseItem;
            Object value = readResponseItem.getValues().get(0);
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertThat(value, equalTo(Boolean.TRUE));
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertThat(value, equalTo(Byte.valueOf("1")));
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertThat(value, equalTo(Short.valueOf("1")));
            } else if (payloadClazzName.equals(Calendar.class.getSimpleName())) {
                assertThat(value, equalTo(calenderInstance));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(value, equalTo(Float.valueOf("1")));
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertThat(value, equalTo(Integer.valueOf("1")));
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertThat(value, equalTo(String.valueOf("Hello World!")));
            }
        }
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