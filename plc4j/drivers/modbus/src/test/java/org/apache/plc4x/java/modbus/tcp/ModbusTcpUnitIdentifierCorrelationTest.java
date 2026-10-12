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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.Function;

import static org.junit.jupiter.api.Assertions.*;

/**
 * GH-2798: a Modbus TCP response is only the answer to a request when BOTH the transaction
 * identifier and the unit identifier match. Behind a TCP-to-serial gateway the unit identifier
 * names the slave, so a frame for unit 2 carrying the transaction id of a request to unit 1 must
 * not complete that request. Such a frame is discarded (as plc4go and the RTU flavor do) and the
 * request keeps waiting for a matching response until its timeout.
 */
class ModbusTcpUnitIdentifierCorrelationTest {

    private static final byte[] READ_RESPONSE_PDU = {0x03, 0x02, 0x00, (byte) 250};

    private FakeModbusTcpServer server;

    @AfterEach
    void tearDown() throws Exception {
        if (server != null) {
            server.close();
        }
    }

    @Test
    void readWithMatchingUnitIdentifierSucceeds() throws Exception {
        server = new FakeModbusTcpServer(request -> List.of(request.answer(request.unitId(), READ_RESPONSE_PDU)));
        try (PlcConnection connection = connect(1)) {
            PlcReadResponse response = read(connection, "holding-register:1:INT");
            assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
            assertEquals(250, response.getInteger("value"));
        }
    }

    @Test
    void readWithMismatchedUnitIdentifierIsRejected() throws Exception {
        server = new FakeModbusTcpServer(request -> List.of(request.answer(2, READ_RESPONSE_PDU)));
        try (PlcConnection connection = connect(1)) {
            assertRejected(readAsync(connection, "holding-register:1:INT"));
            assertEquals(1, server.requests.size());
            assertEquals(1, server.requests.getFirst().unitId());
        }
    }

    @Test
    void writeWithMismatchedUnitIdentifierIsRejectedAndNotReplayed() throws Exception {
        server = new FakeModbusTcpServer(request -> List.of(request.answer(2, request.pdu())));
        try (PlcConnection connection = connect(1)) {
            assertRejected(connection.writeRequestBuilder()
                .addTagAddress("value", "holding-register:1:INT", 456)
                .build().execute());
            // An invalid acknowledgement doesn't prove the device didn't execute the write.
            assertEquals(1, server.requests.size());
        }
    }

    @Test
    void writeWithMatchingUnitIdentifierSucceeds() throws Exception {
        server = new FakeModbusTcpServer(request -> List.of(request.answer(request.unitId(), request.pdu())));
        try (PlcConnection connection = connect(1)) {
            PlcWriteResponse response = connection.writeRequestBuilder()
                .addTagAddress("value", "holding-register:1:INT", 456)
                .build().execute().get(5, TimeUnit.SECONDS);
            assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
        }
    }

    @Test
    void perTagUnitIdentifierOverrideIsWhatTheResponseIsCheckedAgainst() throws Exception {
        // The tag addresses unit 5 on a connection defaulting to unit 1: echoing unit 5 must be
        // accepted and answering as unit 1 rejected - a check against the default would flip both.
        server = new FakeModbusTcpServer(request -> List.of(request.answer(request.unitId(), READ_RESPONSE_PDU)));
        try (PlcConnection connection = connect(1)) {
            PlcReadResponse response = read(connection, "holding-register:1:INT{unit-id: 5}");
            assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
            assertEquals(5, server.requests.getFirst().unitId());
        }
        server.close();

        server = new FakeModbusTcpServer(request -> List.of(request.answer(1, READ_RESPONSE_PDU)));
        try (PlcConnection connection = connect(1)) {
            assertRejected(readAsync(connection, "holding-register:1:INT{unit-id: 5}"));
        }
    }

    @Test
    void mismatchedFrameIsDiscardedAndTheMatchingOneStillCompletesTheRequest() throws Exception {
        server = new FakeModbusTcpServer(request -> List.of(
            request.answer(2, new byte[]{0x03, 0x02, 0x00, 0x01}),
            request.answer(request.unitId(), READ_RESPONSE_PDU)));
        try (PlcConnection connection = connect(1)) {
            PlcReadResponse response = read(connection, "holding-register:1:INT");
            assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
            assertEquals(250, response.getInteger("value"));
        }
    }

    @Test
    void connectionKeepsWorkingAfterARejectedResponse() throws Exception {
        server = new FakeModbusTcpServer(request -> List.of(
            request.answer(request.index() == 0 ? 2 : request.unitId(), READ_RESPONSE_PDU)));
        try (PlcConnection connection = connect(1)) {
            assertRejected(readAsync(connection, "holding-register:1:INT"));
            PlcReadResponse response = read(connection, "holding-register:1:INT");
            assertEquals(PlcResponseCode.OK, response.getResponseCode("value"));
            assertEquals(250, response.getInteger("value"));
        }
    }

    private PlcConnection connect(int defaultUnitIdentifier) throws Exception {
        PlcConnection connection = new DefaultPlcDriverManager().getConnection(
            "modbus-tcp:tcp://127.0.0.1:" + server.port()
                + "?default-unit-identifier=" + defaultUnitIdentifier + "&request-timeout-ms=500");
        if (!connection.isConnected()) {
            connection.connect();
        }
        return connection;
    }

    private static PlcReadResponse read(PlcConnection connection, String address) throws Exception {
        return readAsync(connection, address).get(5, TimeUnit.SECONDS);
    }

    private static CompletableFuture<? extends PlcReadResponse> readAsync(PlcConnection connection, String address) {
        return connection.readRequestBuilder()
            .addTagAddress("value", address)
            .build().execute();
    }

    /**
     * A request whose only answer was discarded must not succeed: it runs into the request
     * timeout, which - depending on which layer's timer fires first - fails the request or
     * reports a non-OK code for the tag.
     */
    private static void assertRejected(CompletableFuture<? extends PlcResponse> execution) throws Exception {
        PlcResponse response;
        try {
            response = execution.get(5, TimeUnit.SECONDS);
        } catch (ExecutionException e) {
            assertInstanceOf(TimeoutException.class, e.getCause());
            return;
        }
        PlcResponseCode code = response instanceof PlcReadResponse readResponse
            ? readResponse.getResponseCode("value")
            : ((PlcWriteResponse) response).getResponseCode("value");
        assertNotEquals(PlcResponseCode.OK, code);
    }

    /**
     * A received request: its position on the connection, MBAP header and PDU.
     */
    private record Request(int index, byte[] header, byte[] pdu) {
        int unitId() {
            return header[6] & 0xFF;
        }

        /**
         * Builds a response frame carrying this request's transaction identifier.
         */
        byte[] answer(int unitId, byte[] responsePdu) {
            byte[] frame = new byte[7 + responsePdu.length];
            System.arraycopy(header, 0, frame, 0, 4);
            int length = responsePdu.length + 1;
            frame[4] = (byte) (length >>> 8);
            frame[5] = (byte) length;
            frame[6] = (byte) unitId;
            System.arraycopy(responsePdu, 0, frame, 7, responsePdu.length);
            return frame;
        }
    }

    /**
     * Single-connection Modbus TCP peer answering each request with the frames its responder returns.
     */
    private static final class FakeModbusTcpServer implements AutoCloseable {
        private final ServerSocket serverSocket;
        private final Thread thread;
        final List<Request> requests = new CopyOnWriteArrayList<>();

        FakeModbusTcpServer(Function<Request, List<byte[]>> responder) throws IOException {
            serverSocket = new ServerSocket(0, 1, InetAddress.getLoopbackAddress());
            thread = new Thread(() -> serve(responder), "fake-modbus-tcp-server");
            thread.setDaemon(true);
            thread.start();
        }

        int port() {
            return serverSocket.getLocalPort();
        }

        private void serve(Function<Request, List<byte[]>> responder) {
            try (Socket socket = serverSocket.accept()) {
                InputStream in = socket.getInputStream();
                OutputStream out = socket.getOutputStream();
                while (true) {
                    byte[] header = in.readNBytes(7);
                    if (header.length != 7) {
                        return;
                    }
                    int length = ((header[4] & 0xFF) << 8) | (header[5] & 0xFF);
                    byte[] pdu = in.readNBytes(length - 1);
                    if (pdu.length != length - 1) {
                        return;
                    }
                    Request request = new Request(requests.size(), header, pdu);
                    requests.add(request);
                    for (byte[] frame : responder.apply(request)) {
                        out.write(frame);
                    }
                    out.flush();
                }
            } catch (IOException e) {
                // Connection closed by the client or the test - nothing left to serve.
            }
        }

        @Override
        public void close() throws Exception {
            serverSocket.close();
            thread.join(TimeUnit.SECONDS.toMillis(5));
        }
    }
}
