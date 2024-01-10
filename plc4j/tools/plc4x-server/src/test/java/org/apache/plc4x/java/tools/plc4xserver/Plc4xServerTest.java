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

package org.apache.plc4x.java.tools.plc4xserver;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertNotNull;

import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

public class Plc4xServerTest {

    private static final Plc4xServer SERVER = new Plc4xServer();
    private static final String CONNECTION_STRING_TEMPLATE = "plc4x://localhost:%d?remote-connection-string=%s";
    private static final String CONNECTION_STRING_SIMULATED_ENCODED = "simulated%3A%2F%2Flocalhost";
    private static final long TIMEOUT_VALUE = 10;
    private static final TimeUnit TIMEOUT_UNIT = TimeUnit.SECONDS;

    private final PlcConnectionManager connectionManager = new DefaultPlcDriverManager();

    @BeforeAll
    public static void setUp() throws ExecutionException, InterruptedException, TimeoutException {
        SERVER.start().get(TIMEOUT_VALUE, TIMEOUT_UNIT);
    }

    @AfterAll
    public static void tearDown() {
        SERVER.stop();
    }

    @Test
    public void testWrite() throws Exception {
        final PlcWriteResponse response;

        try (PlcConnection connection = connectionManager.getConnection(
                String.format(CONNECTION_STRING_TEMPLATE, SERVER.getPort(), CONNECTION_STRING_SIMULATED_ENCODED))) {
            final PlcWriteRequest request = connection.writeRequestBuilder()
                    .addTagAddress(
                            "foo",
                            "STATE/foo:DINT",
                            42
                    )
                    .build();
            response = request.execute().get(TIMEOUT_VALUE, TIMEOUT_UNIT);
        }

        assertEquals(PlcResponseCode.OK, response.getResponseCode("foo"));
    }

    @Test
    public void testRead() throws Exception {
        final PlcReadResponse response;

        try (PlcConnection connection = connectionManager.getConnection(
                String.format(CONNECTION_STRING_TEMPLATE, SERVER.getPort(), CONNECTION_STRING_SIMULATED_ENCODED))) {
            final PlcReadRequest request = connection.readRequestBuilder()
                    .addTagAddress(
                            "foo",
                            "RANDOM/foo:DINT"
                    )
                    .build();
            response = request.execute().get(TIMEOUT_VALUE, TIMEOUT_UNIT);
        }

        assertEquals(PlcResponseCode.OK, response.getResponseCode("foo"));
        assertNotNull(response.getPlcValue("foo"));
        assertInstanceOf(Integer.class, response.getPlcValue("foo").getObject());
    }

    @Test
    public void testReadWriteSameConnection() throws Exception {
        final PlcWriteResponse writeResponse;
        final PlcReadResponse readResponse;

        try (PlcConnection connection = connectionManager.getConnection(
                String.format(CONNECTION_STRING_TEMPLATE, SERVER.getPort(), CONNECTION_STRING_SIMULATED_ENCODED))) {
            final PlcWriteRequest writeRequest = connection.writeRequestBuilder()
                    .addTagAddress(
                            "foo",
                            "STATE/foo:DINT",
                            21
                    )
                    .build();
            writeResponse = writeRequest.execute().get(TIMEOUT_VALUE, TIMEOUT_UNIT);

            final PlcReadRequest readRequest = connection.readRequestBuilder()
                    .addTagAddress(
                            "foo",
                            "STATE/foo:DINT"
                    )
                    .build();
            readResponse = readRequest.execute().get(TIMEOUT_VALUE, TIMEOUT_UNIT);
        }

        assertEquals(PlcResponseCode.OK, writeResponse.getResponseCode("foo"));
        assertEquals(PlcResponseCode.OK, readResponse.getResponseCode("foo"));

        assertInstanceOf(Integer.class, readResponse.getPlcValue("foo").getObject());
        assertEquals(21, readResponse.getInteger("foo"));
    }
}
