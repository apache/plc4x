/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.S7ScriptedConnectionHarness.ScriptedS7Transport;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.util.HexFormat;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Exercises BOOL reads through the connected driver, including the packed byte responses
 * requested by the optimizer for arrays (GH-2742).
 */
class S7ConnectionBoolArrayTest {

    /**
     * Expected bits are listed in array order, not the usual most-significant-bit-first
     * notation. Set padding bits must not appear as extra elements in the returned list.
     */
    @ParameterizedTest(name = "{0}: payload {1} yields {2}")
    @CsvSource({
        "[1..8;1], 45, 10100010",
        "[0..7], 01, 10000000",
        "[0..7], 80, 00000001",
        "[0..7], 00, 00000000",
        "[0..7], ff, 11111111",
        "[0..1], fd, 10",
        "[0..6], 80, 0000000",
        "[0..8], 8101, 100000011",
        "[0..8], 81fe, 100000010",
        "[0..15], 8142, 1000000101000010",
        "[0..16], 8142ff, 10000001010000101"
    })
    void readsPackedBoolArray(String array, String payload, String expected) throws Exception {
        PlcReadResponse response = read("%DB42:214.0" + array + ":BOOL",
            DataTransportSize.BYTE_WORD_DWORD, payload);

        assertBoolArray(response, expected);
    }

    @ParameterizedTest(name = "{0}: payload {1} yields {2}")
    @CsvSource({
        "%DB1.DBX13.3[0..7]:BOOL, 1f7c, 11000001",
        "%DB1.DBX13.5[0..2]:BOOL, 1f, 000",
        "%DB1.DBX13.1[0..6]:BOOL, 81, 0000001",
        "%DB1.DBX13.2[0..5]:BOOL, 83, 000001",
        "%DB1.DBX13.3[0..4]:BOOL, 87, 00001",
        "%DB1.DBX13.4[0..3]:BOOL, 8f, 0001",
        "%DB1.DBX13.5[0..2]:BOOL, 9f, 001",
        "%DB1.DBX13.6[0..1]:BOOL, bf, 01",
        "%DB1.DBX13.7[0..1]:BOOL, 7ffe, 00",
        "%DB1.DBX13.7[0..1]:BOOL, 8001, 11",
        "%DB1:13.7[0..8]:BOOL, 8042, 101000010",
        "%DB1:13.7[0..9]:BOOL, 8042fe, 1010000100",
        "%DB1:13.7[0..9]:BOOL, 804201, 1010000101"
    })
    void readsPackedBoolArrayFromNonzeroBit(String address, String payload, String expected) throws Exception {
        PlcReadResponse response = read(address, DataTransportSize.BYTE_WORD_DWORD, payload);

        assertBoolArray(response, expected);
    }

    private static void assertBoolArray(PlcReadResponse response, String expected) {
        assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
        PlcValue value = response.getPlcValue("value");
        assertTrue(value.isList());
        assertEquals(expected.length(), value.getLength());
        for (int i = 0; i < expected.length(); i++) {
            assertEquals(expected.charAt(i) == '1', value.getIndex(i).getBoolean(), "element " + i);
        }
    }

    @ParameterizedTest
    @CsvSource({"00, false", "01, true"})
    void stillReadsScalarBool(String payload, boolean expected) throws Exception {
        // A scalar BIT response carries the selected bit in bit zero, even for address .3.
        PlcReadResponse response = read("%DB42:214.3:BOOL", DataTransportSize.BIT, payload);

        assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
        assertEquals(expected, response.getPlcValue("value").getBoolean());
        assertFalse(response.getPlcValue("value").isList());
    }

    @Test
    void stillReadsByteArray() throws Exception {
        PlcReadResponse response = read("%DB42:216.0[0..3]:BYTE",
            DataTransportSize.BYTE_WORD_DWORD, "01807fff");

        assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
        assertArrayEquals(new byte[]{0x01, (byte) 0x80, 0x7f, (byte) 0xff},
            response.getPlcValue("value").getRaw());
    }

    @Test
    void stillReadsIntArray() throws Exception {
        PlcReadResponse response = read("%DB42:216[0..1]:INT",
            DataTransportSize.INTEGER, "1234fffe");

        assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
        PlcValue value = response.getPlcValue("value");
        assertTrue(value.isList());
        assertEquals(2, value.getLength());
        assertEquals(0x1234, value.getIndex(0).getInteger());
        assertEquals(-2, value.getIndex(1).getInteger());
    }

    private static PlcReadResponse read(String address, DataTransportSize transportSize,
                                       String payload) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        try (S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport)) {
            CompletableFuture<? extends PlcReadResponse> response = connection.readRequestBuilder()
                .addTagAddress("value", address)
                .build()
                .execute();

            long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(5);
            while (transport.writeCount() == 0 && System.nanoTime() < deadline) {
                Thread.sleep(10);
            }
            assertEquals(1, transport.writeCount(), "the read request should reach the transport");
            S7Message request = S7Message.staticParse(new ReadBufferByteBased(transport.writtenFrames().get(0),
                WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                WithOption.WithSignedIntegerEncoding("twos-complement"),
                WithOption.WithFloatEncoding("IEEE754"),
                WithByteBasedOption.WithByteOrder("BIG_ENDIAN")));
            S7Tag tag = S7Tag.of(address);
            if (tag.getDataType() == TransportSize.BOOL) {
                S7ParameterReadVarRequest parameter = assertInstanceOf(S7ParameterReadVarRequest.class,
                    request.getParameter());
                assertEquals(1, parameter.getItems().size());
                S7VarRequestParameterItemAddress item = assertInstanceOf(S7VarRequestParameterItemAddress.class,
                    parameter.getItems().get(0));
                S7AddressAny wireAddress = assertInstanceOf(S7AddressAny.class, item.getAddress());
                boolean array = tag.getNumberOfElements() > 1;
                // Packed arrays read whole bytes; scalar BIT reads retain their selected bit.
                // The fixture contains exactly the bytes needed, including leading/trailing bits.
                assertAll("BOOL wire address",
                    () -> assertEquals(array ? TransportSize.BYTE : TransportSize.BOOL, wireAddress.getTransportSize()),
                    () -> assertEquals(payload.length() / 2, wireAddress.getNumberOfElements()),
                    () -> assertEquals(tag.getBlockNumber(), wireAddress.getDbNumber()),
                    () -> assertEquals(tag.getByteOffset(), wireAddress.getByteAddress()),
                    () -> assertEquals(array ? 0 : tag.getBitOffset(), wireAddress.getBitAddress()));
            }
            int reference = request.getTpduReference();
            // Echo the request reference so the real decoder and response correlation are exercised.
            S7Message reply = new S7MessageResponseData(reference,
                new S7ParameterReadVarResponse((short) 1),
                new S7PayloadReadVarResponse(List.of(new S7VarPayloadDataItem(
                    DataTransportErrorCode.OK, transportSize, HexFormat.of().parseHex(payload)))),
                (short) 0, (short) 0);
            WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[reply.getLengthInBytes()],
                WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                WithOption.WithSignedIntegerEncoding("twos-complement"),
                WithOption.WithFloatEncoding("IEEE754"),
                WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
            reply.serialize(buffer);
            transport.deliver(buffer.getBytes());
            transport.runDataListener();

            return response.get(5, TimeUnit.SECONDS);
        }
    }
}
