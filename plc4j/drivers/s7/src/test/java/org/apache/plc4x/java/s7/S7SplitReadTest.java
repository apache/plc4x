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
import org.junit.jupiter.api.Test;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A tag too big for one PDU is fetched in byte fragments and decoded once they are all in. It
 * used to be split by elements, and the reassembled bytes were handed back raw for every type -
 * an INT array came back as a byte array.
 * <p>
 * The requests are answered from a simulated data block with the 240 byte PDU the harness
 * negotiates.
 */
class S7SplitReadTest {

    private static final int DATA_BLOCK_SIZE = 1024;

    @Test
    void decodesAnIntArraySplitAcrossPdus() throws Exception {
        ByteBuffer db = ByteBuffer.allocate(DATA_BLOCK_SIZE);
        List<Short> expected = new ArrayList<>();
        for (int i = 0; i < 200; i++) {
            short value = (short) (i * 3 - 100);
            expected.add(value);
            db.putShort(i * 2, value);
        }
        List<byte[]> requests = new ArrayList<>();

        PlcValue value = read("%DB1:0[0..199]:INT", db.array(), requests, -1, PlcResponseCode.OK);

        assertEquals(expected, value.getList().stream().map(PlcValue::getShort).toList());
        assertTrue(requests.size() > 1, "400 bytes should not fit into one 240 byte PDU");
    }

    @Test
    void decodesARealArraySplitAcrossPdus() throws Exception {
        ByteBuffer db = ByteBuffer.allocate(DATA_BLOCK_SIZE);
        List<Float> expected = new ArrayList<>();
        for (int i = 0; i < 100; i++) {
            float value = i * 1.5f;
            expected.add(value);
            db.putFloat(400 + i * 4, value);
        }

        PlcValue value = read("%DB1:400[0..99]:REAL", db.array(), new ArrayList<>(), -1, PlcResponseCode.OK);

        assertEquals(expected, value.getList().stream().map(PlcValue::getFloat).toList());
    }

    @Test
    void decodesAFixedLengthStringArraySplitAcrossPdus() throws Exception {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        putString(db, 0, 200, "first");
        putString(db, 202, 200, "second");

        PlcValue value = read("%DB1:0[0..1]:STRING(200)", db, new ArrayList<>(), -1, PlcResponseCode.OK);

        assertEquals(List.of("first", "second"), value.getList().stream().map(PlcValue::getString).toList());
    }

    /**
     * A STRING without a declared length is read as 254 characters plus its two length bytes,
     * which alone is more than a 240 byte PDU holds. It used to be sized as two bytes.
     */
    @Test
    void decodesAVarLengthStringBiggerThanOnePdu() throws Exception {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        putString(db, 0, 254, "hello");
        List<byte[]> requests = new ArrayList<>();

        PlcValue value = read("%DB1:0:STRING", db, requests, -1, PlcResponseCode.OK);

        assertEquals("hello", value.getString());
        assertTrue(requests.size() > 1);
    }

    @Test
    void aFailedFirstFragmentFailsTheTag() throws Exception {
        read("%DB1:0[0..199]:INT", new byte[DATA_BLOCK_SIZE], new ArrayList<>(), 0, PlcResponseCode.ACCESS_DENIED);
    }

    /** A later fragment that succeeds must not cover up an earlier one that failed. */
    @Test
    void aFailedLaterFragmentFailsTheTag() throws Exception {
        read("%DB1:0[0..199]:INT", new byte[DATA_BLOCK_SIZE], new ArrayList<>(), 1, PlcResponseCode.ACCESS_DENIED);
    }

    private static void putString(byte[] db, int offset, int maxLength, String text) {
        byte[] chars = text.getBytes(StandardCharsets.US_ASCII);
        db[offset] = (byte) maxLength;
        db[offset + 1] = (byte) chars.length;
        System.arraycopy(chars, 0, db, offset + 2, chars.length);
    }

    private static PlcValue read(String address, byte[] db, List<byte[]> requests, int failing,
                                 PlcResponseCode expectedCode) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<? extends PlcReadResponse> future = connection.readRequestBuilder()
            .addTagAddress("value", address)
            .build()
            .execute();
        S7ScriptedConnectionHarness.answerReadsFromDataBlock(transport, future, db, requests, failing);

        PlcReadResponse response = future.get(5, TimeUnit.SECONDS);
        connection.close();
        assertEquals(expectedCode, response.getResponseCode("value"));
        return response.getPlcValue("value");
    }
}
