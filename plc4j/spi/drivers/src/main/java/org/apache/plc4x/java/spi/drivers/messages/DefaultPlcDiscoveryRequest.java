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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

public class DefaultPlcDiscoveryRequest implements PlcDiscoveryRequest {

    private final PlcDiscoverer discoverer;
    private final Map<String, String> queries;

    public DefaultPlcDiscoveryRequest(PlcDiscoverer discoverer, Map<String, String> queries) {
        this.discoverer = discoverer;
        this.queries = queries;
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> execute() {
        return discoverer.discover(this);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> executeWithHandler(PlcDiscoveryItemHandler handler) {
        return discoverer.discoverWithHandler(this, handler);
    }

    public Map<String, String> getQueries() {
        return queries;
    }

    public static class Builder implements PlcDiscoveryRequest.Builder {

        private final PlcDiscoverer discoverer;
        private final LinkedHashMap<String, String> queries = new LinkedHashMap<>();

        public Builder(PlcDiscoverer discoverer) {
            this.discoverer = discoverer;
        }

        @Override
        public PlcDiscoveryRequest.Builder addQuery(String name, String query) {
            queries.put(name, query);
            return this;
        }

        @Override
        public PlcDiscoveryRequest build() {
            return new DefaultPlcDiscoveryRequest(discoverer, queries);
        }
    }

}
