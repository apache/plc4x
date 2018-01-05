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
package org.apache.plc4x.java.s7.netty;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcRequestContainer;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.netty.NettyTestBase;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.lang.reflect.Field;
import java.util.*;
import java.util.function.Consumer;
import java.util.stream.Stream;

import static java.util.Collections.singletonList;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.*;

@SuppressWarnings("unchecked")
public class Plc4XS7ProtocolTest extends NettyTestBase {

    private Plc4XS7Protocol SUT;

    @BeforeEach
    void setUp() {
        SUT = new Plc4XS7Protocol();
    }

    @ParameterizedTest
    @MethodSource("typeAndAddressProvider")
    @Tag("fast")
    public void encode(Class<?> type, S7Address address) throws Exception {
        // TODO: finish me
        // Read Request Tests
        {
            LinkedList<Object> out = new LinkedList<>();
            SUT.encode(null, createMockedContainer(new PlcReadRequest(type, address)), out);
            // TODO: finish the asserts
            assertThat(out).hasSize(1);
        }
        // Write Request Tests
        {
            LinkedList<Object> out = new LinkedList<>();
            SUT.encode(null, createMockedContainer(new PlcWriteRequest(type, address, fakeValueFor(type))), out);
            // TODO: finish the asserts
            assertThat(out).hasSize(1);
        }
    }

    @ParameterizedTest
    @MethodSource("typeAndAddressProvider")
    @Tag("fast")
    public void decode(Class<?> type, S7Address address) throws Exception {
        // TODO: finish me
        if (type == String.class) {
            // String seems not yet decodable
            return;
        }
        // Read Test
        {
            short fakeTpduReference = (short) 1;
            {
                // We need to put in a fake tpdu reference
                Field requests = Plc4XS7Protocol.class.getDeclaredField("requests");
                requests.setAccessible(true);
                Map<Short, PlcRequestContainer> requestContainerMap = (Map<Short, PlcRequestContainer>) requests.get(SUT);
                requestContainerMap.put(fakeTpduReference, createMockedContainer(new PlcReadRequest(type, address)));
            }
            S7ResponseMessage msg = new S7ResponseMessage(
                MessageType.ACK,
                fakeTpduReference,
                singletonList(mock(VarParameter.class)),
                singletonList(new VarPayload(ParameterType.READ_VAR, singletonList(varPayloadItemFor(type)))),
                (byte) 0x00,
                (byte) 0x00);
            LinkedList<Object> out = new LinkedList<>();
            SUT.decode(null, msg, out);
            // TODO: finish the asserts
            assertThat(out).hasSize(0);
        }
        // Write Test
        {
            short fakeTpduReference = (short) 2;
            {
                // We need to put in a fake tpdu reference
                Field requests = Plc4XS7Protocol.class.getDeclaredField("requests");
                requests.setAccessible(true);
                Map<Short, PlcRequestContainer> requestContainerMap = (Map<Short, PlcRequestContainer>) requests.get(SUT);
                requestContainerMap.put(fakeTpduReference, createMockedContainer(new PlcWriteRequest(type, address, fakeValueFor(type))));
            }
            S7ResponseMessage msg = new S7ResponseMessage(
                MessageType.ACK,
                fakeTpduReference,
                singletonList(mock(VarParameter.class)),
                singletonList(new VarPayload(ParameterType.WRITE_VAR, singletonList(varPayloadItemFor(type)))),
                (byte) 0x00,
                (byte) 0x00);
            LinkedList<Object> out = new LinkedList<>();
            SUT.decode(null, msg, out);
            // TODO: finish the asserts
            assertThat(out).hasSize(0);
        }
    }

    private static Stream<Arguments> typeAndAddressProvider() {
        List<Arguments> arguments = new LinkedList<>();
        Arrays.asList(
            Boolean.class,
            Byte.class,
            Short.class,
            // TODO: enable once Calender in implemented
            //Calendar.class,
            Float.class,
            Integer.class,
            String.class)
            .forEach(
                aClass -> Arrays.asList(
                    mock(S7Address.class),
                    mock(S7BitAddress.class),
                    mock(S7DataBlockAddress.class))
                    .forEach(s7Address -> arguments.add(Arguments.of(aClass, s7Address)))
            );
        return arguments.stream();
    }

    private <T> T fakeValueFor(Class<T> type) {
        if (type == Boolean.class) {
            return (T) Boolean.TRUE;
        } else if (type == Byte.class) {
            return (T) Byte.valueOf((byte) 0x0000_0000);
        } else if (type == Short.class) {
            return (T) Short.valueOf((short) 123);
        } else if (type == Calendar.class) {
            return (T) Calendar.getInstance();
        } else if (type == Float.class) {
            return (T) Float.valueOf(123f);
        } else if (type == Integer.class) {
            return (T) Integer.valueOf(123);
        } else if (type == String.class) {
            return (T) "string";
        } else {
            throw new IllegalArgumentException("Type t not supported " + type);
        }
    }

    private VarPayloadItem varPayloadItemFor(Class type) {
        // TODO: fix example
        final DataTransportSize size;
        final byte[] data;
        if (type == Boolean.class) {
            size = DataTransportSize.BIT;
            data = new byte[]{(byte) 0b0};
        } else if (type == Byte.class) {
            size = DataTransportSize.BYTE_WORD_DWORD;
            data = new byte[]{(byte) 0b0000_0000};
        } else if (type == Short.class) {
            size = DataTransportSize.BYTE_WORD_DWORD;
            data = new byte[]{(byte) 0b0000_0000, (byte) 0b0000_0000};
        } else if (type == Calendar.class) {
            size = DataTransportSize.BYTE_WORD_DWORD;
            // TODO: what size is calender?
            data = new byte[]{(byte) 0b0000_0000};
        } else if (type == Float.class) {
            size = DataTransportSize.BYTE_WORD_DWORD;
            data = new byte[]{(byte) 0b0000_0000, (byte) 0b0000_0000, (byte) 0b0000_0000, (byte) 0b0000_0000};
        } else if (type == Integer.class) {
            size = DataTransportSize.INTEGER;
            data = new byte[]{(byte) 0b0000_0000, (byte) 0b0000_0000, (byte) 0b0000_0000, (byte) 0b0000_0000};
        } else if (type == String.class) {
            size = DataTransportSize.BYTE_WORD_DWORD;
            // TODO: what size is string?
            data = new byte[]{(byte) 0b0000_0000};
        } else {
            throw new IllegalArgumentException("Type t not supported " + type);
        }
        return new VarPayloadItem(DataTransportErrorCode.OK, size, data);
    }

    private <T extends PlcRequest> PlcRequestContainer createMockedContainer(T initialRequest) {
        return createMockedContainer(initialRequest, null);
    }

    private <T extends PlcRequest> PlcRequestContainer createMockedContainer(T initialRequest, Consumer<T> requestEnricher) {
        Objects.requireNonNull(initialRequest);
        PlcRequestContainer mock = mock(PlcRequestContainer.class, RETURNS_DEEP_STUBS);
        if (requestEnricher != null) {
            requestEnricher.accept(initialRequest);
        }
        when(mock.getRequest()).thenReturn(initialRequest);
        if (initialRequest.getClass() == PlcReadRequest.class) {
            return mock;
        } else if (initialRequest.getClass() == PlcWriteRequest.class) {
            return mock;
        } else {
            throw new IllegalArgumentException("Unsupported Type: " + initialRequest.getClass());
        }
    }

}
