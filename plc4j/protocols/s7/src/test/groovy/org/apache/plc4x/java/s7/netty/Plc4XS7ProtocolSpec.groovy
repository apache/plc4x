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

package org.apache.plc4x.java.s7.netty

import io.netty.channel.ChannelFuture
import io.netty.channel.ChannelHandlerContext
import io.netty.channel.ChannelOutboundHandlerAdapter
import io.netty.channel.ChannelPromise
import io.netty.channel.embedded.EmbeddedChannel
import org.apache.plc4x.java.api.exceptions.PlcProtocolException
import org.apache.plc4x.java.api.messages.PlcReadRequest
import org.apache.plc4x.java.api.messages.PlcReadResponse
import org.apache.plc4x.java.api.messages.PlcResponse
import org.apache.plc4x.java.api.messages.PlcWriteRequest
import org.apache.plc4x.java.api.types.PlcResponseCode
import org.apache.plc4x.java.base.messages.*
import org.apache.plc4x.java.s7.netty.model.messages.S7Message
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter
import org.apache.plc4x.java.s7.netty.model.params.VarParameter
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem
import org.apache.plc4x.java.s7.netty.model.types.*
import org.apache.plc4x.java.s7.netty.util.S7PlcFieldHandler
import spock.lang.Specification
import spock.lang.Unroll

import java.lang.reflect.Field
import java.util.concurrent.CompletableFuture
import java.util.concurrent.TimeUnit

import static org.mockito.Mockito.mock

class Plc4XS7ProtocolSpec extends Specification {

    private EmbeddedChannel SUT
    private PlcReadRequest.Builder readRequestBuilder
    private PlcWriteRequest.Builder writeRequestBuilder
    private CompletableFuture<S7Message> writeFuture
    private CompletableFuture<PlcResponse> readFuture

    @Unroll
    def "Test encoding the different PLC read requests: '#address'"(String address, boolean expectedSuccess, MemoryArea expectedMemoryArea, int dataBlockNumber, int byteOffset, int bitOffset, int numItems, TransportSize datatype) {
        setup:
        initTest(null)
        CompletableFuture<InternalPlcReadRequest> future = new CompletableFuture<>()
        PlcRequestContainer container = new PlcRequestContainer(
            (DefaultPlcReadRequest) readRequestBuilder.addItem("foo", address).build(), future)
        ChannelFuture channelFuture = SUT.writeOneOutbound(container)
        S7Message writtenMessage = writeFuture.get(100, TimeUnit.MILLISECONDS)
        Throwable exception = channelFuture.cause()

        expect:
        assert channelFuture.isSuccess() == expectedSuccess
        assert expectedSuccess ? exception == null : exception != null
        assert expectedSuccess ? writtenMessage != null : writtenMessage == null
        if(expectedSuccess)
            assertReadRequestSuccess(writtenMessage, expectedMemoryArea, dataBlockNumber, byteOffset, bitOffset, numItems, datatype)
        else
            assertFailure(exception)

        where:
        address    | expectedSuccess | expectedMemoryArea | dataBlockNumber | byteOffset | bitOffset | numItems | datatype
        "%Q0:BYTE" | true            | MemoryArea.OUTPUTS | 0               | 0          | 0         | 1        | TransportSize.BYTE
    }

    void assertReadRequestSuccess(S7Message writtenMessage, MemoryArea memoryArea, int dataBlockNumber, int byteOffset, int bitOffset, int numItems, TransportSize datatype) {
        assert writtenMessage instanceof S7RequestMessage
        assert writtenMessage.getParameters().size() == 1
        assert writtenMessage.getPayloads().isEmpty()
        assert writtenMessage.getParent() != null

        assert writtenMessage.getParameters().get(0) instanceof VarParameter
        VarParameter varParameter = (VarParameter) writtenMessage.getParameters().get(0)
        assert varParameter.getItems().size() == 1

        assert varParameter.getItems().get(0) instanceof S7AnyVarParameterItem
        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameter.getItems().get(0)
        assert s7AnyVarParameterItem.getSpecificationType() == SpecificationType.VARIABLE_SPECIFICATION
        assert s7AnyVarParameterItem.getMemoryArea() == memoryArea
        assert s7AnyVarParameterItem.getDataBlockNumber() == (short) dataBlockNumber
        assert s7AnyVarParameterItem.getByteOffset() == (short) byteOffset
        assert s7AnyVarParameterItem.getBitOffset() == (byte) bitOffset
        assert s7AnyVarParameterItem.getNumElements() == numItems
        assert s7AnyVarParameterItem.getDataType() == datatype
    }

    @Unroll
    def "Test encoding the different PLC write requests: '#address'"(String address, Object value, boolean expectedSuccess, MemoryArea expectedMemoryArea, int dataBlockNumber, int byteOffset, int bitOffset, int numItems, TransportSize datatype) {
        setup:
        initTest(null)
        CompletableFuture<InternalPlcWriteRequest> future = new CompletableFuture<>()
        PlcRequestContainer container = new PlcRequestContainer(
            (DefaultPlcWriteRequest) writeRequestBuilder.addItem("foo", address, value).build(), future)
        ChannelFuture channelFuture = SUT.writeOneOutbound(container)
        S7Message writtenMessage = writeFuture.get(100, TimeUnit.MILLISECONDS)
        Throwable exception = channelFuture.cause()

        expect:
        assert channelFuture.isSuccess() == expectedSuccess
        assert expectedSuccess ? exception == null : exception != null
        assert expectedSuccess ? writtenMessage != null : writtenMessage == null
        if(expectedSuccess)
            assertWriteRequestSuccess(writtenMessage, expectedMemoryArea, dataBlockNumber, byteOffset, bitOffset, numItems, datatype)
        else
            assertFailure(exception)

        where:
        address             | value               || expectedSuccess | expectedMemoryArea | dataBlockNumber | byteOffset | bitOffset | numItems | datatype
        "%Q3.4:BOOL"        | true                || true            | MemoryArea.OUTPUTS | 0               | 3          | 4         | 1        | TransportSize.BOOL
        "%Q3:SINT"          | (byte) 0x42         || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.SINT
        "%Q3:INT"           | (short) 0x4223      || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.INT
        "%Q3:DINT"          | (int) 4223          || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.DINT
        "%Q3:LINT"          | (long) 4223123      || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.LINT
        "%Q3:USINT"         | (byte) 0x42         || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.USINT
        "%Q3:UINT"          | (short) 0x4223      || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.UINT
        "%Q3:UDINT"         | (int) 4223          || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.UDINT
        "%Q3:ULINT"         | (long) 4223123      || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.ULINT
        "%Q3:REAL"          | (float) 42.312      || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.REAL
        "%Q3:LREAL"         | (double) 42.32      || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.LREAL
        "%Q3:STRING"        | "foo"               || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 256      | TransportSize.STRING
        "%Q3:WSTRING"       | "bar"               || true            | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.WSTRING
        //"%Q3:DATE_AND_TIME" | LocalDateTime.now() || false           | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.DATE_AND_TIME
    }

    void assertWriteRequestSuccess(S7Message writtenMessage, MemoryArea memoryArea, int dataBlockNumber, int byteOffset, int bitOffset, int numItems, TransportSize datatype) {
        assert writtenMessage instanceof S7RequestMessage
        assert writtenMessage.getParameters().size() == 1
        assert writtenMessage.getPayloads().size() == 1
        assert writtenMessage.getParent() != null

        assert writtenMessage.getParameters().get(0) instanceof VarParameter
        VarParameter varParameter = (VarParameter) writtenMessage.getParameters().get(0)
        assert varParameter.getItems().size() == 1

        assert varParameter.getItems().get(0) instanceof S7AnyVarParameterItem
        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameter.getItems().get(0)
        assert s7AnyVarParameterItem.getSpecificationType() == SpecificationType.VARIABLE_SPECIFICATION
        assert s7AnyVarParameterItem.getMemoryArea() == memoryArea
        assert s7AnyVarParameterItem.getDataBlockNumber() == (short) dataBlockNumber
        assert s7AnyVarParameterItem.getByteOffset() == (short) byteOffset
        assert s7AnyVarParameterItem.getBitOffset() == (byte) bitOffset
        assert s7AnyVarParameterItem.getNumElements() == numItems
        assert s7AnyVarParameterItem.getDataType() == datatype

        assert writtenMessage.getPayloads().get(0) instanceof VarPayload
        VarPayload varPayload = (VarPayload) writtenMessage.getPayloads().get(0)
        assert varPayload.getItems().size() == 1
    }



    @Unroll
    def "Test decoding the different PLC read responses: '#address'"(String address, MemoryArea memoryArea,
                  int dataBlockNumber, int byteOffset, int bitOffset, int numItems, TransportSize datatype,
                  DataTransportSize dataTransportSize, byte[] data, boolean expectedSuccess,
                  PlcResponseCode expectedResponseCode, int expectedNumItems, Object[] expectedValues) {
        setup:
        initTest(address)
        S7Parameter parameter = new VarParameter(ParameterType.READ_VAR, Collections.singletonList(
            new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, memoryArea, datatype, numItems,
                (short) dataBlockNumber, (short) byteOffset, (byte) bitOffset)))
        S7Payload payload = new VarPayload(ParameterType.READ_VAR, Collections.singletonList(
            new VarPayloadItem(DataTransportErrorCode.OK, dataTransportSize, data)))
        // The tpdu has to match the id manually injected in the setup method.
        S7Message response = new S7ResponseMessage(MessageType.JOB, (short) 1,
            Collections.singletonList(parameter), Collections.singletonList(payload), (byte) 0, (byte) 0)
        SUT.writeInbound(response)
        PlcResponse receivedMessage = readFuture.get(100, TimeUnit.MILLISECONDS)

        expect:
        //assert expectedSuccess ? exception == null : exception != null
        assert expectedSuccess ? receivedMessage != null : receivedMessage == null
        if(expectedSuccess)
            assertReadResponseSuccess(receivedMessage, expectedResponseCode, expectedNumItems, expectedValues)
        else
            assertFailure(exception)

        where:

        // Bit-String values
        address      | memoryArea         | dataBlockNumber | byteOffset | bitOffset | numItems | datatype            | dataTransportSize                 | data                                             || expectedSuccess | expectedResponseCode | expectedNumItems | expectedValues
        "%Q3.4:BOOL" | MemoryArea.OUTPUTS | 0               | 3          | 4         | 1        | TransportSize.BOOL  | DataTransportSize.BIT             | [0x01]                                           || true            | PlcResponseCode.OK   | 1                | [true]
        "%Q3:BYTE"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.BYTE  | DataTransportSize.BYTE_WORD_DWORD | [0x23]                                           || true            | PlcResponseCode.OK   | 8                | [true, true, false, false, false, true, false, false]
        "%Q3:WORD"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.WORD  | DataTransportSize.BYTE_WORD_DWORD | [0x43, 0xA1]                                     || true            | PlcResponseCode.OK   | 16               | [true, false, false, false, false, true, false, true, true, true, false, false, false, false, true, false]
        "%Q3:DWORD"  | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.DWORD | DataTransportSize.BYTE_WORD_DWORD | [0x43, 0xA1, 0x11, 0xFA]                         || true            | PlcResponseCode.OK   | 32               | [false, true, false, true, true, true, true, true, true, false, false, false, true, false, false, false, true, false, false, false, false, true, false, true, true, true, false, false, false, false, true, false]
        "%Q3:LWORD"  | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.LWORD | DataTransportSize.BYTE_WORD_DWORD | [0x43, 0xA1, 0x11, 0xFA, 0x53, 0x82, 0x13, 0x85] || true            | PlcResponseCode.OK   | 64               | [true, false, true, false, false, false, false, true, true, true, false, false, true, false, false, false, false, true, false, false, false, false, false, true, true, true, false, false, true, false, true, false, false, true, false, true, true, true, true, true, true, false, false, false, true, false, false, false, true, false, false, false, false, true, false, true, true, true, false, false, false, false, true, false]
        // Integer values
        //address    | memoryArea         | dataBlockNumber | byteOffset | bitOffset | numItems | datatype            | dataTransportSize                 | data                                             || expectedSuccess | expectedResponseCode | expectedNumItems | expectedValues
        "%Q3:SINT"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.SINT  | DataTransportSize.INTEGER         | [0x43]                                           || true            | PlcResponseCode.OK   | 1                | [(byte) 0x43]
        "%Q3:USINT"  | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.USINT | DataTransportSize.INTEGER         | [0xFF]                                           || true            | PlcResponseCode.OK   | 1                | [(short) 255]
        "%Q3:INT"    | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.INT   | DataTransportSize.INTEGER         | [0x42, 0x23]                                     || true            | PlcResponseCode.OK   | 1                | [(short) 16931]
        "%Q3:UINT"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.UINT  | DataTransportSize.INTEGER         | [0xFF, 0xFF]                                     || true            | PlcResponseCode.OK   | 1                | [(int) 65535]
        "%Q3:DINT"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.DINT  | DataTransportSize.INTEGER         | [0x42, 0x23, 0x12, 0x74]                         || true            | PlcResponseCode.OK   | 1                | [(int) 1109594740]
        "%Q3:UDINT"  | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.UDINT | DataTransportSize.INTEGER         | [0xFF, 0xFF, 0xFF, 0xFF]                         || true            | PlcResponseCode.OK   | 1                | [(long) 4294967295]
        "%Q3:LINT"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.LINT  | DataTransportSize.INTEGER         | [0x43, 0xA1, 0x11, 0xFA, 0x53, 0x82, 0x13, 0x85] || true            | PlcResponseCode.OK   | 1                | [(long) 4873196038632117125]
        //ULINT
        // Floating point values
        "%Q3:REAL"   | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.REAL  | DataTransportSize.REAL            | [0x42, 0x23, 0x12, 0x74]                         || true            | PlcResponseCode.OK   | 1                | [40.76802F]
        "%Q3:LREAL"  | MemoryArea.OUTPUTS | 0               | 3          | 0         | 1        | TransportSize.LREAL | DataTransportSize.REAL            | [0x43, 0xA1, 0x11, 0xFA, 0x53, 0x82, 0x13, 0x85] || true            | PlcResponseCode.OK   | 1                | [6.150197049102015e+17D]
    }

    void assertReadResponseSuccess(PlcResponse receivedMessage, PlcResponseCode expectedResponseCode, int expectedNumItems, Object[] expectedValues) {
        assert receivedMessage instanceof PlcReadResponse
        PlcReadResponse readResponse = (PlcReadResponse) receivedMessage

        assert readResponse.getNumberOfValues("bar") == expectedNumItems
        assert readResponse.getResponseCode("bar") == expectedResponseCode
        Object[] actualValues = readResponse.getAllObjects("bar")
        if(expectedValues == null) {
            assert actualValues == null
        } else {
            assert actualValues.length == expectedValues.length
            for (int i = 0; i < actualValues.length; i++) {
                assert actualValues[i].class == expectedValues[i].class
                assert actualValues[i] == expectedValues[i]
            }
        }
    }

    void assertFailure(Throwable exception) {
        assert exception != null
        assert exception instanceof PlcProtocolException
    }

    void initTest(String address) {
        writeFuture = new CompletableFuture<>()
        readRequestBuilder = new DefaultPlcReadRequest.Builder(mock(PlcReader.class), new S7PlcFieldHandler())
        writeRequestBuilder = new DefaultPlcWriteRequest.Builder(mock(PlcWriter.class), new S7PlcFieldHandler())

        Plc4XS7Protocol protocol = new Plc4XS7Protocol()

        if(address) {
            readFuture = new CompletableFuture<>()
            // Populate the 'requests' field with some manually crafted requests.
            Field requestsFiled = Plc4XS7Protocol.class.getDeclaredField("requests")
            requestsFiled.setAccessible(true)
            Map<Short, PlcRequestContainer> requestMap = (Map<Short, PlcRequestContainer>) requestsFiled.get(protocol)
            requestMap.put((short) 1, new PlcRequestContainer<>(readRequestBuilder.addItem("bar", address).build(), readFuture))
            //requestMap.put((short) 2, new PlcRequestContainer<>(writeRequestBuilder.addItem("bar", address, 42).build(), readFuture))
        }

        SUT = new EmbeddedChannel(new ChannelOutboundHandlerAdapter() {
            @Override
            void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) {
                if(msg instanceof S7Message) {
                    writeFuture.complete((S7Message) msg)
                    promise.setSuccess()
                } else {
                    promise.setFailure(new PlcProtocolException(
                        "Got message of type " + msg.getClass().getSimpleName()))
                }
            }
        }, protocol)
    }


}
