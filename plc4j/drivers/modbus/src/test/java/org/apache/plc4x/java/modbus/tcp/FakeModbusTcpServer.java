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

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.TimeUnit;
import java.util.function.Function;

/**
 * Single-connection Modbus TCP peer on the loopback interface, answering each request with the
 * frames its responder returns - none, to let a request run into its timeout.
 */
final class FakeModbusTcpServer implements AutoCloseable {

    /**
     * A received request: its position on the connection, MBAP header and PDU.
     */
    record Request(int index, byte[] header, byte[] pdu) {
        int unitId() {
            return header[6] & 0xFF;
        }

        int functionCode() {
            return pdu[0] & 0xFF;
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

    private final ServerSocket serverSocket;
    private final Thread thread;
    final List<Request> requests = new CopyOnWriteArrayList<>();

    FakeModbusTcpServer(Function<Request, List<byte[]>> responder) throws IOException {
        serverSocket = new ServerSocket(0, 1, InetAddress.getLoopbackAddress());
        thread = new Thread(() -> serve(responder), "fake-modbus-tcp-server");
        thread.setDaemon(true);
        thread.start();
    }

    /**
     * Opens a connection to this server with a short request timeout.
     */
    PlcConnection connect(int defaultUnitIdentifier) throws Exception {
        PlcConnection connection = new DefaultPlcDriverManager().getConnection(
            "modbus-tcp:tcp://127.0.0.1:" + serverSocket.getLocalPort()
                + "?default-unit-identifier=" + defaultUnitIdentifier + "&request-timeout-ms=500");
        if (!connection.isConnected()) {
            connection.connect();
        }
        return connection;
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
