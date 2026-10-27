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

import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.S7ScriptedConnectionHarness.ScriptedS7Transport;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7VarPayloadDataItem;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItemAddress;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A BOOL array used to be written as the whole bytes its bits are packed into, which overwrote
 * every other bit in those bytes, and wrote from bit 0 whatever the tag's bit offset. Only the
 * bytes the array covers completely may be written packed; every bit it shares a byte with has to
 * go out on its own, so the rest of that byte is left alone.
 * <p>
 * The requests are applied to a simulated data block, so what is checked is the memory a PLC ends
 * up with.
 */
class S7BoolArrayWriteTest {

    private static final int DATA_BLOCK_SIZE = 64;

    @Test
    void writesAWholeByteArrayPacked() throws Exception {
        byte[] db = filled((byte) 0x5A);
        List<byte[]> requests = new ArrayList<>();

        assertEquals(PlcResponseCode.OK, writeBools("%DB1:16.0[0..15]:BOOL", bools(16, i -> i % 3 == 0), db,
            requests, null));

        byte[] expected = filled((byte) 0x5A);
        expected[16] = (byte) 0b0100_1001;
        expected[17] = (byte) 0b1001_0010;
        assertArrayEquals(expected, db);
        assertEquals(1, requests.size());
        assertEquals(1, items(requests.get(0)), "whole bytes go out as one packed item");
    }

    @Test
    void leavesTheRestOfAPartlyCoveredByteAlone() throws Exception {
        byte[] db = filled((byte) 0xFF);

        // Ten elements: byte 16, then 17.0 and 17.1. 17.2 to 17.7 belong to something else.
        assertEquals(PlcResponseCode.OK, writeBools("%DB1:16.0[0..9]:BOOL", bools(10, i -> false), db,
            new ArrayList<>(), null));

        byte[] expected = filled((byte) 0xFF);
        expected[16] = 0;
        expected[17] = (byte) 0b1111_1100;
        assertArrayEquals(expected, db);
    }

    @Test
    void writesFromTheTagsBitOffset() throws Exception {
        byte[] db = filled((byte) 0x00);

        // 16.6 to 18.1: two leading bits, byte 17, two trailing bits.
        assertEquals(PlcResponseCode.OK, writeBools("%DB1:16.6[0..11]:BOOL", bools(12, i -> true), db,
            new ArrayList<>(), null));

        byte[] expected = filled((byte) 0x00);
        expected[16] = (byte) 0b1100_0000;
        expected[17] = (byte) 0xFF;
        expected[18] = (byte) 0b0000_0011;
        assertArrayEquals(expected, db);
    }

    @Test
    void writesABoolArraySplitAcrossPdus() throws Exception {
        byte[] db = filled((byte) 0x00);
        List<byte[]> requests = new ArrayList<>();

        // 16.1 to 18.6: seven leading bits, byte 17, seven trailing bits - 15 items.
        assertEquals(PlcResponseCode.OK, writeBools("%DB1:16.1[0..21]:BOOL", bools(22, i -> true), db,
            requests, null));

        byte[] expected = filled((byte) 0x00);
        expected[16] = (byte) 0b1111_1110;
        expected[17] = (byte) 0xFF;
        expected[18] = (byte) 0b0111_1111;
        assertArrayEquals(expected, db);
        assertTrue(requests.size() > 1, "15 items should not fit into one 240 byte PDU");
    }

    @Test
    void reportsTheFirstFailedItemForTheTag() throws Exception {
        PlcResponseCode code = writeBools("%DB1:16.0[0..9]:BOOL", bools(10, i -> true), filled((byte) 0),
            new ArrayList<>(), DataTransportErrorCode.ACCESS_DENIED);

        assertEquals(PlcResponseCode.ACCESS_DENIED, code);
    }

    @Test
    void writesNothingForAValueWithTooFewElements() throws Exception {
        List<byte[]> requests = new ArrayList<>();

        PlcResponseCode code = writeBools("%DB1:16.0[0..9]:BOOL", bools(9, i -> true), filled((byte) 0),
            requests, null);

        assertEquals(PlcResponseCode.INVALID_DATA, code);
        assertTrue(requests.isEmpty());
    }

    private interface BitPattern {
        boolean at(int i);
    }

    private static PlcList bools(int count, BitPattern pattern) {
        List<PlcBOOL> values = new ArrayList<>(count);
        for (int i = 0; i < count; i++) {
            values.add(new PlcBOOL(pattern.at(i)));
        }
        return new PlcList(values);
    }

    private static byte[] filled(byte value) {
        byte[] db = new byte[DATA_BLOCK_SIZE];
        Arrays.fill(db, value);
        return db;
    }

    /**
     * Writes the value, applies every request to the data block and returns the tag's code.
     *
     * @param failLastItem when set, the status the last item of the first request is answered
     *                     with instead of applying it
     */
    private static PlcResponseCode writeBools(String address, PlcList value, byte[] db, List<byte[]> requests,
                                              DataTransportErrorCode failLastItem) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<? extends PlcWriteResponse> future = connection.writeRequestBuilder()
            .addTagAddress("bools", address, value)
            .build()
            .execute();

        // Chunks go out one after another, so each request is answered as it arrives.
        long deadline = System.currentTimeMillis() + 5_000;
        while (!future.isDone() && System.currentTimeMillis() < deadline) {
            if (transport.writeCount() > requests.size()) {
                byte[] request = transport.writtenFrames().get(requests.size());
                DataTransportErrorCode failure = requests.isEmpty() ? failLastItem : null;
                requests.add(request);
                transport.deliver(applyToDataBlock(request, db, failure));
                transport.runDataListener();
            } else {
                Thread.sleep(5);
            }
        }

        PlcResponseCode code = future.get(5, TimeUnit.SECONDS).getResponseCode("bools");
        connection.close();
        return code;
    }

    private static int items(byte[] request) throws Exception {
        return ((S7ParameterWriteVarRequest) S7ScriptedConnectionHarness.parseFrame(request).getParameter())
            .getItems().size();
    }

    /** What a PLC does with a write request to data block 1: applies it and reports per item. */
    private static byte[] applyToDataBlock(byte[] request, byte[] db, DataTransportErrorCode failLastItem)
        throws Exception {
        S7Message message = S7ScriptedConnectionHarness.parseFrame(request);
        S7ParameterWriteVarRequest parameter = (S7ParameterWriteVarRequest) message.getParameter();
        List<S7VarPayloadDataItem> payload = ((S7PayloadWriteVarRequest) message.getPayload()).getItems();
        List<DataTransportErrorCode> statuses = new ArrayList<>();
        for (int i = 0; i < parameter.getItems().size(); i++) {
            if (failLastItem != null && i == parameter.getItems().size() - 1) {
                statuses.add(failLastItem);
                continue;
            }
            S7AddressAny address = (S7AddressAny)
                ((S7VarRequestParameterItemAddress) parameter.getItems().get(i)).getAddress();
            byte[] data = payload.get(i).getData();
            int start = address.getByteAddress();
            if (address.getTransportSize() == TransportSize.BOOL) {
                assertEquals(1, address.getNumberOfElements());
                int mask = 1 << address.getBitAddress();
                db[start] = (byte) ((data[0] & 0x01) != 0 ? db[start] | mask : db[start] & ~mask);
            } else {
                assertEquals(TransportSize.BYTE, address.getTransportSize());
                assertEquals(0, address.getBitAddress(), "a byte write has no bit address");
                assertEquals(address.getNumberOfElements(), data.length);
                System.arraycopy(data, 0, db, start, data.length);
            }
            statuses.add(DataTransportErrorCode.OK);
        }
        return S7ScriptedConnectionHarness.writeVarResponse(
            S7ScriptedConnectionHarness.tpduReferenceOfFrame(request), statuses);
    }
}
