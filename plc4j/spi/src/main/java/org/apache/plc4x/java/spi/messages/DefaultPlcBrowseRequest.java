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
package org.apache.plc4x.java.spi.messages;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;

public class DefaultPlcBrowseRequest implements PlcBrowseRequest, Serializable {

    private final PlcBrowser browser;

    private final LinkedHashMap<String, PlcQuery> queries;

    public DefaultPlcBrowseRequest(PlcBrowser browser,
                                   LinkedHashMap<String, PlcQuery> queries) {
        this.browser = browser;
        this.queries = queries;
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> execute() {
        return browser.browse(this);
    }

    @Override
    public CompletableFuture<? extends PlcBrowseResponse> executeWithInterceptor(PlcBrowseRequestInterceptor interceptor) {
        return browser.browseWithInterceptor(this, interceptor);
    }

    public PlcBrowser getBrowser() {
        return browser;
    }

    @Override
    public LinkedHashSet<String> getQueryNames() {
        return new LinkedHashSet<>(queries.keySet());
    }

    @Override
    public PlcQuery getQuery(String name) {
        return queries.get(name);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcBrowseRequest");
        writeBuffer.popContext("PlcBrowseRequest");
    }

    public static class Builder implements PlcBrowseRequest.Builder {

        private final PlcBrowser browser;
        private final PlcTagHandler tagHandler;
        private final LinkedHashMap<String, Supplier<PlcQuery>> queries;

        public Builder(PlcBrowser browser, PlcTagHandler tagHandler) {
            this.browser = browser;
            this.tagHandler = tagHandler;
            queries = new LinkedHashMap<>();
        }

        @Override
        public Builder addQuery(String name, String query) {
            if (queries.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate query definition '" + name + "'");
            }
            queries.put(name, () -> tagHandler.parseQuery(query));
            return this;
        }

        @Override
        public PlcBrowseRequest build() {
            LinkedHashMap<String, PlcQuery> parsedQueries = new LinkedHashMap<>();
            queries.forEach((name, tagQuery) -> {
                parsedQueries.put(name, tagQuery.get());
            });
            return new DefaultPlcBrowseRequest(browser, parsedQueries);
        }

    }

}
