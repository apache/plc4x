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
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItemAddress;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A BOOL array is requested as the bytes its bits are packed into - eight elements to a byte,
 * least significant bit first - so it has to be unpacked bit by bit rather than parsed one byte
 * per element. Parsing it per element ran out of data after the first byte, left the remaining
 * elements null and failed the tag with INTERNAL_ERROR (GH-2742).
 * <p>
 * The requests are answered from a simulated data block, so what is checked is what a PLC would
 * hand back for the bytes the driver actually asked for.
 */
class S7BoolArrayReadTest {

    private static final int DATA_BLOCK_SIZE = 1024;

    @Test
    void unpacksABoolArrayFromPackedBytes() throws Exception {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        // Elements 1, 2 and 9 set: bits 1 and 2 of the first byte, bit 1 of the second.
        db[16] = 0b0000_0110;
        db[17] = 0b0000_0010;

        assertEquals(List.of(false, true, true, false, false, false, false, false, false, true),
            readBools("%DB1:16.0[0..9]:BOOL", db));
    }

    @Test
    void unpacksAFullByteOfBools() throws Exception {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        db[16] = (byte) 0b1000_0001;

        assertEquals(List.of(true, false, false, false, false, false, false, true),
            readBools("%DB1:16.0[0..7]:BOOL", db));
    }

    /**
     * A byte read has no bit address, so an array starting at a bit offset other than 0 is read
     * from bit 0 of its start byte, over every byte its bits touch, and the leading bits skipped.
     */
    @Test
    void readsABoolArrayStartingInsideAByte() throws Exception {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        // 16.6, 16.7, 17.0, 17.1: set, set, clear, set - with the bits after them set as well, so
        // reading from the wrong bit shows.
        db[16] = (byte) 0b1100_0000;
        db[17] = (byte) 0b1111_1010;

        assertEquals(List.of(true, true, false, true), readBools("%DB1:16.6[0..3]:BOOL", db));
    }

    @Test
    void readsTheBitsASelectionStartsAt() throws Exception {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        db[16] = (byte) 0b0101_0000;

        assertEquals(List.of(true, false, true, false), readBools("%DB1:16.0[4..7]:BOOL", db));
    }

    /**
     * An array too big for one PDU is fetched in byte fragments and unpacked once all of them are
     * in. It used to be split as if each BOOL were a byte, and came back as raw bytes.
     */
    @Test
    void readsABoolArrayLargerThanOnePdu() throws Exception {
        int elements = 2000;
        byte[] db = new byte[DATA_BLOCK_SIZE];
        List<Boolean> expected = new ArrayList<>(elements);
        for (int i = 0; i < elements; i++) {
            boolean value = (i % 3 == 0) || (i % 7 == 0);
            expected.add(value);
            int bit = 16 * 8 + 3 + i;
            if (value) {
                db[bit / 8] |= (byte) (1 << (bit % 8));
            }
        }

        List<byte[]> requests = new ArrayList<>();
        assertEquals(expected, readBools("%DB1:16.3[0.." + (elements - 1) + "]:BOOL", db, requests));
        assertTrue(requests.size() > 1, "the array should not fit into one 240 byte PDU");
    }

    private List<Boolean> readBools(String address, byte[] db) throws Exception {
        return readBools(address, db, new ArrayList<>());
    }

    private List<Boolean> readBools(String address, byte[] db, List<byte[]> requests) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<? extends PlcReadResponse> future = connection.readRequestBuilder()
            .addTagAddress("bools", address)
            .build()
            .execute();

        // Chunks go out one after another, so each request is answered as it arrives.
        long deadline = System.currentTimeMillis() + 5_000;
        while (!future.isDone() && System.currentTimeMillis() < deadline) {
            if (transport.writeCount() > requests.size()) {
                byte[] request = transport.writtenFrames().get(requests.size());
                requests.add(request);
                transport.deliver(answerFromDataBlock(request, db));
                transport.runDataListener();
            } else {
                Thread.sleep(5);
            }
        }

        PlcReadResponse response = future.get(5, TimeUnit.SECONDS);
        connection.close();
        assertEquals(PlcResponseCode.OK, response.getResponseCode("bools"));
        return response.getPlcValue("bools").getList().stream().map(PlcValue::getBoolean).toList();
    }

    /** What a PLC sends back for a single-item byte read of data block 1. */
    private static byte[] answerFromDataBlock(byte[] request, byte[] db) throws Exception {
        S7ParameterReadVarRequest parameter = (S7ParameterReadVarRequest)
            S7ScriptedConnectionHarness.parseFrame(request).getParameter();
        assertEquals(1, parameter.getItems().size());
        S7AddressAny address = (S7AddressAny)
            ((S7VarRequestParameterItemAddress) parameter.getItems().get(0)).getAddress();
        assertEquals(TransportSize.BYTE, address.getTransportSize());
        assertEquals(0, address.getBitAddress(), "a byte read has no bit address");
        int start = address.getByteAddress();
        return S7ScriptedConnectionHarness.readVarResponse(
            S7ScriptedConnectionHarness.tpduReferenceOfFrame(request), DataTransportSize.BYTE_WORD_DWORD,
            Arrays.copyOfRange(db, start, start + address.getNumberOfElements()));
    }
}
