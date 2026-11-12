/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.ctrlx.readwrite.rest;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.net.http.HttpClient;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Runs the client against the JDK's own HTTP server standing in for a ctrlX device.
 */
class CtrlXRestClientTest {

    private static final String TOKEN_JSON = "{\"access_token\":\"t0k3n\",\"token_type\":\"Bearer\"}";

    private HttpServer server;
    private String baseUrl;

    /** Answers by request path (without query): status and body. */
    private final Map<String, Answer> answers = new ConcurrentHashMap<>();
    /** What the server saw, by request path. */
    private final Map<String, Seen> seen = new ConcurrentHashMap<>();

    private record Answer(int status, String body) {
    }

    private record Seen(String method, String rawPath, String rawQuery, String authorization,
                        String contentType, String body) {
    }

    @BeforeEach
    void startServer() throws IOException {
        server = HttpServer.create(new InetSocketAddress(InetAddress.getLoopbackAddress(), 0), 0);
        server.createContext("/", this::handle);
        server.start();
        baseUrl = "http://localhost:" + server.getAddress().getPort();
    }

    @AfterEach
    void stopServer() {
        server.stop(0);
    }

    private void handle(HttpExchange exchange) throws IOException {
        String path = exchange.getRequestURI().getRawPath();
        seen.put(path, new Seen(exchange.getRequestMethod(), path, exchange.getRequestURI().getRawQuery(),
            exchange.getRequestHeaders().getFirst("Authorization"),
            exchange.getRequestHeaders().getFirst("Content-Type"),
            new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8)));
        Answer answer = answers.getOrDefault(path, new Answer(404, "{\"title\":\"not found\"}"));
        byte[] body = answer.body().getBytes(StandardCharsets.UTF_8);
        exchange.sendResponseHeaders(answer.status(), body.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(body);
        }
    }

    private CtrlXRestClient login() throws CtrlXRestException {
        answers.put(CtrlXRestClient.TOKEN_PATH, new Answer(201, TOKEN_JSON));
        return CtrlXRestClient.login(HttpClient.newHttpClient(), baseUrl, "user name", "p&ss=word");
    }

    @Test
    void logsInWithAPasswordGrant() throws Exception {
        try (CtrlXRestClient ignored = login()) {
            Seen login = seen.get(CtrlXRestClient.TOKEN_PATH);
            assertEquals("POST", login.method());
            assertTrue(login.contentType().startsWith("application/x-www-form-urlencoded"), login.contentType());
            assertEquals("username=user+name&password=p%26ss%3Dword&grant_type=password", login.body());
        }
    }

    @Test
    void listsTheRootNodesWithTheToken() throws Exception {
        answers.put("/automation/api/v2/nodes/datalayer/nodes",
            new Answer(200, "{\"type\":\"arstring\",\"value\":[\"framework\",\"motion\"]}"));
        answers.put("/automation/api/v2/nodes/datalayer/nodesrt",
            new Answer(200, "{\"type\":\"arstring\",\"value\":[\"rt\"]}"));
        try (CtrlXRestClient client = login()) {
            assertEquals(List.of("framework", "motion"), client.getNodeNames());
            assertEquals(List.of("rt"), client.getRealtimeNodeNames());
        }
        Seen nodes = seen.get("/automation/api/v2/nodes/datalayer/nodes");
        assertEquals("GET", nodes.method());
        assertEquals("Bearer t0k3n", nodes.authorization());
    }

    @Test
    void browsesANodeWithItsPathInOneSegment() throws Exception {
        String rawPath = "/automation/api/v2/nodes/framework%2Fmetrics%20x";
        answers.put(rawPath, new Answer(200, "{\"type\":\"arstring\",\"value\":[\"cpu\"]}"));
        try (CtrlXRestClient client = login()) {
            assertEquals(List.of("cpu"), client.browse("framework/metrics x"));
        }
        assertEquals("type=browse", seen.get(rawPath).rawQuery());
    }

    @Test
    void aMissingValueIsAnEmptyList() throws Exception {
        answers.put("/automation/api/v2/nodes/datalayer/nodes", new Answer(200, "{\"type\":\"arstring\"}"));
        answers.put("/automation/api/v2/nodes/datalayer/nodesrt", new Answer(200, "{\"value\":null}"));
        try (CtrlXRestClient client = login()) {
            assertEquals(List.of(), client.getNodeNames());
            assertEquals(List.of(), client.getRealtimeNodeNames());
        }
    }

    @Test
    void anErrorStatusCarriesStatusAndBody() throws Exception {
        try (CtrlXRestClient client = login()) {
            CtrlXRestException e = assertThrows(CtrlXRestException.class, () -> client.browse("nope"));
            assertEquals(404, e.getStatusCode());
            assertTrue(e.getMessage().contains("not found"), e.getMessage());
        }
    }

    @Test
    void answersOfTheWrongShapeAreRejected() throws Exception {
        answers.put("/automation/api/v2/nodes/a", new Answer(200, "not json"));
        answers.put("/automation/api/v2/nodes/b", new Answer(200, "[\"an array\"]"));
        answers.put("/automation/api/v2/nodes/c", new Answer(200, "{\"value\":\"no list\"}"));
        answers.put("/automation/api/v2/nodes/d", new Answer(200, "{\"value\":[1]}"));
        try (CtrlXRestClient client = login()) {
            CtrlXRestException invalid = assertThrows(CtrlXRestException.class, () -> client.browse("a"));
            assertEquals(200, invalid.getStatusCode());
            for (String path : List.of("b", "c", "d")) {
                CtrlXRestException e = assertThrows(CtrlXRestException.class, () -> client.browse(path));
                assertEquals(0, e.getStatusCode());
                assertNull(e.getCause());
            }
        }
    }

    @Test
    void aLoginWithoutTokenFails() {
        answers.put(CtrlXRestClient.TOKEN_PATH, new Answer(200, "{\"token_type\":\"Bearer\"}"));
        assertThrows(CtrlXRestException.class,
            () -> CtrlXRestClient.login(HttpClient.newHttpClient(), baseUrl, "u", "p"));
        answers.put(CtrlXRestClient.TOKEN_PATH, new Answer(200, "{\"access_token\":\"\"}"));
        assertThrows(CtrlXRestException.class,
            () -> CtrlXRestClient.login(HttpClient.newHttpClient(), baseUrl, "u", "p"));
    }

    @Test
    void aRejectedLoginFails() {
        answers.put(CtrlXRestClient.TOKEN_PATH, new Answer(401, "{\"title\":\"wrong password\"}"));
        CtrlXRestException e = assertThrows(CtrlXRestException.class,
            () -> CtrlXRestClient.login(HttpClient.newHttpClient(), baseUrl, "u", "p"));
        assertEquals(401, e.getStatusCode());
    }

    @Test
    void anUnreachableDeviceFails() throws IOException {
        int closedPort;
        try (ServerSocket socket = new ServerSocket(0, 1, InetAddress.getLoopbackAddress())) {
            closedPort = socket.getLocalPort();
        }
        CtrlXRestException e = assertThrows(CtrlXRestException.class,
            () -> CtrlXRestClient.login(HttpClient.newHttpClient(), "http://localhost:" + closedPort, "u", "p"));
        assertTrue(e.getCause() instanceof IOException, String.valueOf(e.getCause()));
    }

    @Test
    void anInterruptedCallFailsAndKeepsTheInterrupt() {
        Thread.currentThread().interrupt();
        try {
            assertThrows(CtrlXRestException.class, this::login);
            assertTrue(Thread.currentThread().isInterrupted());
        } finally {
            Thread.interrupted();
        }
    }

    @Test
    void exceptionConstructors() {
        CtrlXRestException plain = new CtrlXRestException("m");
        assertEquals(0, plain.getStatusCode());
        assertNull(plain.getCause());
        IOException cause = new IOException("c");
        assertEquals(cause, new CtrlXRestException("m", cause).getCause());
    }

}
