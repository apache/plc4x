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
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * A BOOL array is requested as packed bytes - eight elements to a byte, least significant bit
 * first - so it has to be unpacked bit by bit rather than parsed one byte per element. Parsing it
 * per element ran out of data after the first byte, left the remaining elements null and failed
 * the tag with INTERNAL_ERROR (GH-2742).
 */
class S7BoolArrayReadTest {

    @Test
    void unpacksABoolArrayFromPackedBytes() throws Exception {
        // Elements 1, 2 and 9 set: bits 1 and 2 of the first byte, bit 1 of the second.
        List<? extends PlcValue> values = readBoolArray("%DB1:16.0[0..9]:BOOL",
            new byte[]{0b0000_0110, 0b0000_0010});

        assertEquals(List.of(false, true, true, false, false, false, false, false, false, true),
            values.stream().map(PlcValue::getBoolean).toList());
    }

    @Test
    void unpacksAFullByteOfBools() throws Exception {
        List<? extends PlcValue> values = readBoolArray("%DB1:16.0[0..7]:BOOL",
            new byte[]{(byte) 0b1000_0001});

        assertEquals(List.of(true, false, false, false, false, false, false, true),
            values.stream().map(PlcValue::getBoolean).toList());
    }

    private List<? extends PlcValue> readBoolArray(String address, byte[] packed) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<? extends PlcReadResponse> future = connection.readRequestBuilder()
            .addTagAddress("bools", address)
            .build()
            .execute();

        byte[] request = awaitFrame(transport);
        transport.deliver(S7ScriptedConnectionHarness.readVarResponse(
            S7ScriptedConnectionHarness.tpduReferenceOfFrame(request),
            DataTransportSize.BYTE_WORD_DWORD, packed));
        transport.runDataListener();

        PlcReadResponse response = future.get(5, TimeUnit.SECONDS);
        connection.close();
        assertEquals(PlcResponseCode.OK, response.getResponseCode("bools"));
        return response.getPlcValue("bools").getList();
    }

    private static byte[] awaitFrame(ScriptedS7Transport transport) throws InterruptedException {
        long deadline = System.currentTimeMillis() + 5_000;
        while (transport.writeCount() == 0 && System.currentTimeMillis() < deadline) {
            Thread.sleep(10);
        }
        if (transport.writeCount() == 0) {
            throw new IllegalStateException("the driver never sent the read request");
        }
        return transport.writtenFrames().get(0);
    }
}
