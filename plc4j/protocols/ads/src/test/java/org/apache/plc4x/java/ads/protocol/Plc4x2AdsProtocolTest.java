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

import org.apache.commons.lang3.NotImplementedException;
import org.apache.commons.lang3.RandomStringUtils;
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
import org.apache.plc4x.java.ads.model.AdsDataType;
import org.apache.plc4x.java.ads.model.AdsPlcFieldHandler;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes.streamOfLittleEndianDataTypePairs;
import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertThat;
import static org.junit.Assume.assumeThat;

@RunWith(Parameterized.class)
public class Plc4x2AdsProtocolTest {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4x2AdsProtocolTest.class);

    private Plc4x2AdsProtocol SUT;

    @Rule
    public ExpectedException expectedException = ExpectedException.none();
    private List<String> notYetSupportedDataType = Stream.of(
        // We can add types which are not implemented in this protocol. The exception is currently a placeholder (empty list).
        NotImplementedException.class
    ).map(Class::getSimpleName).collect(Collectors.toList());

    @Parameterized.Parameter
    public String payloadClazzName;

    @Parameterized.Parameter(1)
    public PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> plcRequestContainer;

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
            // TODO: calender doesnt work anymore so we might need to adjust the generator above.
            .filter(o -> o.getDataTypeClass() != LocalDateTime.class)
            .filter(o -> o.getDataTypeClass() != Byte[].class)
            .filter(o -> o.getDataTypeClass() != byte[].class)
            .map(Plc4x2AdsProtocolTest::mapToAdsDataType)
            .map(pair -> Stream.of(
                ImmutablePair.of(
                    new PlcRequestContainer<>(
                        (InternalPlcRequest) new DefaultPlcWriteRequest.Builder(null, new AdsPlcFieldHandler()) // TODO: remove null
                            .addItem(RandomStringUtils.randomAscii(10), "1/1:" + pair.adsDataType.name(), pair.getValue())
                            .build(), new CompletableFuture<>()),
                    AdsWriteResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0))
                ),
                ImmutablePair.of(
                    new PlcRequestContainer<>(
                        (InternalPlcRequest) new DefaultPlcReadRequest.Builder(null, new AdsPlcFieldHandler()) // TODO: remove null
                            .addItem(RandomStringUtils.randomAscii(10), "1/1:" + pair.adsDataType.name())
                            .build(), new CompletableFuture<>()),
                    AdsReadResponse.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, Result.of(0), Data.of(pair.getByteRepresentation()))
                )
            ))
            .flatMap(stream -> stream)
            // TODO: request doesn't know its type anymore... fixme
            .map(pair -> new Object[]{"???", pair.left, pair.left.getResponseFuture(), pair.left.getRequest().getClass().getSimpleName(), pair.right, pair.right.getClass().getSimpleName()}).collect(Collectors.toList());
    }

    private static AdsDataTypePair mapToAdsDataType(Plc4XSupportedDataTypes.DataTypePair dataTypePair) {
        Map<Class<?>, AdsDataType> dataTypeMap = new HashMap<>();
        dataTypeMap.put(Boolean.class, AdsDataType.BOOL);
        dataTypeMap.put(Byte.class, AdsDataType.BYTE);
        dataTypeMap.put(Short.class, AdsDataType.INT);
        dataTypeMap.put(Integer.class, AdsDataType.INT32);
        dataTypeMap.put(Long.class, AdsDataType.INT64);
        dataTypeMap.put(BigInteger.class, AdsDataType.INT64);
        dataTypeMap.put(Float.class, AdsDataType.REAL);
        dataTypeMap.put(Double.class, AdsDataType.LREAL);
        dataTypeMap.put(BigDecimal.class, AdsDataType.LREAL);
        dataTypeMap.put(String.class, AdsDataType.STRING);
        dataTypeMap.put(LocalTime.class, AdsDataType.TIME);
        dataTypeMap.put(LocalDate.class, AdsDataType.DATE);
        dataTypeMap.put(LocalDateTime.class, AdsDataType.DATE_AND_TIME);
        dataTypeMap.put(byte[].class, AdsDataType.BYTE);
        dataTypeMap.put(Byte[].class, AdsDataType.BYTE);
        return new AdsDataTypePair(dataTypePair, dataTypeMap.get(dataTypePair.getDataTypeClass()));
    }

    private static class AdsDataTypePair extends Plc4XSupportedDataTypes.DataTypePair {

        private final AdsDataType adsDataType;

        private AdsDataTypePair(Plc4XSupportedDataTypes.DataTypePair dataTypePair, AdsDataType adsDataType) {
            super(dataTypePair.getDataTypePair());
            this.adsDataType = Objects.requireNonNull(adsDataType);
        }
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
        assumeThat(payloadClazzName + " not yet implemented", notYetSupportedDataType, not(hasItem(payloadClazzName)));
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
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1, 0x0, 0x0, 0x0}));
            } else if (payloadClazzName.equals(Long.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1, 0x0, 0x0, 0x0}));
            } else if (payloadClazzName.equals(BigInteger.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x1, 0x0, 0x0, 0x0}));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, 0x0, (byte) 0x80, 0x3F}));
            } else if (payloadClazzName.equals(Double.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, (byte) 0xF0, 0x3F}));
            } else if (payloadClazzName.equals(BigDecimal.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, (byte) 0xF0, 0x3F}));
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            } else if (payloadClazzName.equals(LocalTime.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            } else if (payloadClazzName.equals(LocalDate.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            } else if (payloadClazzName.equals(LocalDateTime.class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            } else if (payloadClazzName.equals(byte[].class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            } else if (payloadClazzName.equals(Byte[].class.getSimpleName())) {
                assertThat(value, equalTo(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
            }
        }
    }

    @Test
    public void decode() throws Exception {
        assumeThat(payloadClazzName + " not yet implemented", notYetSupportedDataType, not(hasItem(payloadClazzName)));
        ArrayList<Object> in = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, in);
        assertThat(in, hasSize(1));
        syncInvoiceId();
        ArrayList<Object> out = new ArrayList<>();
        LOGGER.info("{}\nHexDump:\n{}", amsPacket, amsPacket.dump());
        SUT.decode(null, amsPacket, out);
        assertThat(out, hasSize(0));
        LOGGER.info("PlcRequestContainer {}", plcRequestContainer);
        InternalPlcResponse plcResponse = plcRequestContainer.getResponseFuture().get();
        // TODO: FIXME: this is different now after refactoring
        //ResponseItem responseItem = (ResponseItem) plcResponse.getResponseItem().get();
        //LOGGER.info("ResponseItem {}", responseItem);
        if (amsPacket instanceof AdsReadResponse) {
            // TODO: FIXME: this is different now after refactoring
            //PlcReadResponseItem readResponseItem = (PlcReadResponseItem) responseItem;
            //Object value = readResponseItem.getValues().get(0);
            //defaultAssert(value);
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