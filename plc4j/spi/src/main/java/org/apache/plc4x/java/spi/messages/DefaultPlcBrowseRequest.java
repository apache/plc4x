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

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcBrowseRequest implements PlcBrowseRequest, Serializable {

    private final PlcBrowser browser;

    private final LinkedHashMap<String, PlcQuery> queries;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcBrowseRequest(@JsonProperty("browser") PlcBrowser browser,
                                   @JsonProperty("queries") LinkedHashMap<String, PlcQuery> queries) {
        this.browser = browser;
        this.queries = queries;
    }

    @Override
    @JsonIgnore
    public CompletableFuture<PlcBrowseResponse> execute() {
        return browser.browse(this);
    }

    @Override
    @JsonIgnore
    public CompletableFuture<? extends PlcBrowseResponse> executeWithInterceptor(PlcBrowseRequestInterceptor interceptor) {
        return browser.browseWithInterceptor(this, interceptor);
    }

    @JsonIgnore
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
        private final PlcFieldHandler fieldHandler;
        private final LinkedHashMap<String, Supplier<PlcQuery>> queries;

        public Builder(PlcBrowser browser, PlcFieldHandler fieldHandler) {
            this.browser = browser;
            this.fieldHandler = fieldHandler;
            queries = new LinkedHashMap<>();
        }

        @Override
        public Builder addQuery(String name, String query) {
            if (queries.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate query definition '" + name + "'");
            }
            queries.put(name, () -> fieldHandler.parseQuery(query));
            return this;
        }

        @Override
        public PlcBrowseRequest build() {
            LinkedHashMap<String, PlcQuery> parsedQueries = new LinkedHashMap<>();
            queries.forEach((name, fieldQuery) -> {
                parsedQueries.put(name, fieldQuery.get());
            });
            return new DefaultPlcBrowseRequest(browser, parsedQueries);
        }

    }

}
