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
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.S7ScriptedConnectionHarness.ScriptedS7Transport;
import org.junit.jupiter.api.Test;

import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A tag whose address the builder couldn't parse stays in the request with an error code and a
 * null tag. It has to be reported per tag as INVALID_ADDRESS - the driver used to fail the
 * whole request with a PlcProtocolException, which took the request's valid tags down with it.
 */
class S7ConnectionInvalidTagTest {

    @Test
    void harnessBringsUpAConnection() throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);

        assertTrue(connection.isConnected(), "the scripted handshake should leave the connection up");
        connection.close();
    }

    @Test
    void readWithInvalidTagAddressIsReportedAndNeverSent() throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);

        PlcReadResponse response = connection.readRequestBuilder()
            .addTagAddress("bad", "%DB1:NOSUCHTHING")
            .build()
            .execute()
            .get(5, TimeUnit.SECONDS);

        assertEquals(PlcResponseCode.INVALID_ADDRESS, response.getResponseCode("bad"));
        assertEquals(0, transport.writeCount(), "a rejected tag must not produce a request");

        connection.close();
    }

    @Test
    void writeWithInvalidTagAddressIsReportedAndNeverSent() throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);

        PlcWriteResponse response = connection.writeRequestBuilder()
            .addTagAddress("bad", "%DB1:NOSUCHTHING", 42)
            .build()
            .execute()
            .get(5, TimeUnit.SECONDS);

        assertEquals(PlcResponseCode.INVALID_ADDRESS, response.getResponseCode("bad"));
        assertEquals(0, transport.writeCount(), "a rejected tag must not produce a request");

        connection.close();
    }

    /**
     * The point of per-tag codes: one typo must not stop the rest of the request from being
     * read. The valid tag has to reach the wire even though a sibling was rejected.
     */
    @Test
    void readMixesValidAndInvalidTagAddresses() throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);

        connection.readRequestBuilder()
            .addTagAddress("good", "%DB1.DBW0:INT")
            .addTagAddress("bad", "%DB1:NOSUCHTHING")
            .build()
            .execute();

        // The valid tag still goes out on its own.
        long deadline = System.currentTimeMillis() + 5_000;
        while (transport.writeCount() == 0 && System.currentTimeMillis() < deadline) {
            Thread.sleep(10);
        }
        assertEquals(1, transport.writeCount(), "the valid tag must still be requested");

        connection.close();
    }
}
