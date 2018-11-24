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

package org.apache.plc4x.java.scraper;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

/**
 * Plc Scraper that scrapes one source.
 */
public class Scraper implements Runnable {

    private final PlcDriverManager driverManager;
    private final String connectionString;
    private final long requestTimeoutMs;
    private final ResultHandler handler;

    public Scraper(PlcDriverManager driverManager, String connectionString, long requestTimeoutMs, ResultHandler handler) {
        this.driverManager = driverManager;
        this.connectionString = connectionString;
        this.requestTimeoutMs = requestTimeoutMs;
        this.handler = handler;
    }

    @Override
    public void run() {
        // Does a single fetch
        try (PlcConnection connection = driverManager.getConnection(connectionString)) {
            PlcReadResponse response;
            try {
                response = connection.readRequestBuilder()
                    .addItem("item1", "add1")
                    .build()
                    .execute()
                    .get(requestTimeoutMs, TimeUnit.MILLISECONDS);
            } catch (ExecutionException e) {
                // Handle execution exception
                handler.handleException(e);
                return;
            }
            CompletableFuture.runAsync(() -> handler.handle(transformResponseToMap(response)));
        } catch (PlcConnectionException e) {
            throw new PlcRuntimeException("Unable to fetch", e);
        } catch (Exception e) {
            throw new PlcRuntimeException("Unexpected exception during fetch", e);
        }
    }

    private Map<String, Object> transformResponseToMap(PlcReadResponse response) {
        return response.getFieldNames().stream()
            .collect(Collectors.toMap(
                name -> name,
                response::getObject
            ));
    }

    public interface ResultHandler {

        void handle(Map<String, Object> result);

        void handleException(Exception e);

    }
}
