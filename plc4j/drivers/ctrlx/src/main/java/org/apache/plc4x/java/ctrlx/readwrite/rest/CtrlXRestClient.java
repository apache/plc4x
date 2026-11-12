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

import org.apache.plc4x.java.spi.utils.json.Json;
import org.apache.plc4x.java.spi.utils.json.JsonException;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * The few calls of the ctrlX Data Layer REST API (v2) the driver uses, on the JDK's HTTP client.
 * <p>
 * Every Data Layer answer used here has the shape {@code {"type": "arstring", "value": [...]}};
 * the client hands back the {@code value} list.
 */
public class CtrlXRestClient implements AutoCloseable {

    /** Where the identity manager hands out tokens, relative to the device's base URL. */
    static final String TOKEN_PATH = "/identity-manager/api/v1.0/auth/token";

    /** Root of the Data Layer API, relative to the device's base URL. */
    static final String DATA_LAYER_PATH = "/automation/api/v2";

    private final HttpClient httpClient;
    private final String dataLayerUrl;
    private final String bearerToken;

    private CtrlXRestClient(HttpClient httpClient, String baseUrl, String bearerToken) {
        this.httpClient = httpClient;
        this.dataLayerUrl = baseUrl + DATA_LAYER_PATH;
        this.bearerToken = bearerToken;
    }

    /**
     * Logs in with the given credentials and returns a client using the token the device handed
     * out. The HTTP client is owned by the returned client from then on, and closed with it.
     *
     * @param httpClient HTTP client, already set up with the TLS trust for this device
     * @param baseUrl    the device, e.g. {@code https://192.168.1.1}
     * @param username   ctrlX user
     * @param password   its password
     * @return a logged-in client
     * @throws CtrlXRestException if the login fails
     */
    public static CtrlXRestClient login(HttpClient httpClient, String baseUrl, String username, String password)
        throws CtrlXRestException {
        String form = "username=" + formEncode(username)
            + "&password=" + formEncode(password)
            + "&grant_type=password";
        HttpRequest request = HttpRequest.newBuilder(URI.create(baseUrl + TOKEN_PATH))
            .header("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
            .header("Accept", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(form, StandardCharsets.UTF_8))
            .build();
        Map<?, ?> answer = asObject(send(httpClient, request, "Login"), "Login");
        if (!(answer.get("access_token") instanceof String token) || token.isEmpty()) {
            throw new CtrlXRestException("Login answer carries no access_token");
        }
        return new CtrlXRestClient(httpClient, baseUrl, token);
    }

    /**
     * Names of the Data Layer's top-level nodes ({@code GET /nodes/datalayer/nodes}).
     */
    public List<String> getNodeNames() throws CtrlXRestException {
        return stringList(get("/nodes/datalayer/nodes"), "Listing the nodes");
    }

    /**
     * Names of the Data Layer's real-time nodes ({@code GET /nodes/datalayer/nodesrt}).
     */
    public List<String> getRealtimeNodeNames() throws CtrlXRestException {
        return stringList(get("/nodes/datalayer/nodesrt"), "Listing the real-time nodes");
    }

    /**
     * Names of a node's children ({@code GET /nodes/{path}?type=browse}).
     * <p>
     * The whole path goes into one URL segment, its slashes encoded, the way the client generated
     * from the device's OpenAPI description sent it.
     *
     * @param path the node, e.g. {@code framework/metrics}
     */
    public List<String> browse(String path) throws CtrlXRestException {
        return stringList(get("/nodes/" + pathEncode(path) + "?type=browse"), "Browsing " + path);
    }

    @Override
    public void close() {
        httpClient.close();
    }

    private Object get(String relativeUrl) throws CtrlXRestException {
        HttpRequest request = HttpRequest.newBuilder(URI.create(dataLayerUrl + relativeUrl))
            .header("Authorization", "Bearer " + bearerToken)
            .header("Accept", "application/json")
            .GET()
            .build();
        return send(httpClient, request, "GET " + relativeUrl);
    }

    private static Object send(HttpClient httpClient, HttpRequest request, String what) throws CtrlXRestException {
        HttpResponse<String> response;
        try {
            response = httpClient.send(request, HttpResponse.BodyHandlers.ofString(StandardCharsets.UTF_8));
        } catch (IOException e) {
            throw new CtrlXRestException(what + " failed: " + e.getMessage(), e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new CtrlXRestException(what + " was interrupted", e);
        }
        int status = response.statusCode();
        if (status < 200 || status >= 300) {
            throw new CtrlXRestException(what + " answered with HTTP " + status + ": " + response.body(), status, null);
        }
        try {
            return Json.parse(response.body());
        } catch (JsonException e) {
            throw new CtrlXRestException(what + " answered with invalid JSON: " + e.getMessage(), status, e);
        }
    }

    private static Map<?, ?> asObject(Object json, String what) throws CtrlXRestException {
        if (json instanceof Map<?, ?> map) {
            return map;
        }
        throw new CtrlXRestException(what + " answered with something other than a JSON object");
    }

    /**
     * The {@code value} of a Data Layer answer, which must be a list of strings (or absent/null,
     * which reads as an empty list).
     */
    private static List<String> stringList(Object json, String what) throws CtrlXRestException {
        Object value = asObject(json, what).get("value");
        if (value == null) {
            return List.of();
        }
        if (!(value instanceof List<?> list)) {
            throw new CtrlXRestException(what + " answered with a value that isn't a list");
        }
        List<String> strings = new ArrayList<>(list.size());
        for (Object item : list) {
            if (!(item instanceof String s)) {
                throw new CtrlXRestException(what + " answered with a list holding something other than strings");
            }
            strings.add(s);
        }
        return strings;
    }

    private static String formEncode(String s) {
        return URLEncoder.encode(s, StandardCharsets.UTF_8);
    }

    static String pathEncode(String s) {
        return URLEncoder.encode(s, StandardCharsets.UTF_8).replace("+", "%20");
    }

}
