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
import org.junit.jupiter.api.Test;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * A controller that refuses a write because PUT/GET communication is not enabled answers with an
 * error class in the S7 header, and the caller has to be able to tell that apart from a generic
 * failure (GH-599). An S7-300 sends 0x83/0x04 where other controllers send 0x81/0x04, so both
 * have to come out as ACCESS_DENIED.
 */
class S7AccessDeniedTest {

    @Test
    void reportsAccessDeniedForTheS7300ErrorClass() throws Exception {
        assertEquals(PlcResponseCode.ACCESS_DENIED, writeAnsweredWith(0x83, 0x04));
    }

    @Test
    void reportsAccessDeniedForTheOtherErrorClass() throws Exception {
        assertEquals(PlcResponseCode.ACCESS_DENIED, writeAnsweredWith(0x81, 0x04));
    }

    /**
     * Class 0x83 on its own is "no resources available", so only the 0x04 pairing is a refusal.
     */
    @Test
    void keepsOtherCodesInThatClassGeneric() throws Exception {
        assertEquals(PlcResponseCode.INTERNAL_ERROR, writeAnsweredWith(0x83, 0x01));
    }

    /**
     * Writes a tag, answers the request with the given header error, and returns the code the
     * caller ends up seeing for that tag.
     */
    private PlcResponseCode writeAnsweredWith(int errorClass, int errorCode) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<? extends PlcWriteResponse> future = connection.writeRequestBuilder()
            .addTagAddress("value", "%DB1.DBW0:INT", 42)
            .build()
            .execute();

        byte[] request = awaitFrame(transport);
        transport.deliver(S7ScriptedConnectionHarness.headerErrorResponse(
            S7ScriptedConnectionHarness.tpduReferenceOfFrame(request), errorClass, errorCode));
        transport.runDataListener();

        PlcResponseCode code = future.get(5, TimeUnit.SECONDS).getResponseCode("value");
        connection.close();
        return code;
    }

    private static byte[] awaitFrame(ScriptedS7Transport transport) throws InterruptedException {
        long deadline = System.currentTimeMillis() + 5_000;
        while (transport.writeCount() == 0 && System.currentTimeMillis() < deadline) {
            Thread.sleep(10);
        }
        if (transport.writeCount() == 0) {
            throw new IllegalStateException("the driver never sent the write request");
        }
        return transport.writtenFrames().get(0);
    }
}
