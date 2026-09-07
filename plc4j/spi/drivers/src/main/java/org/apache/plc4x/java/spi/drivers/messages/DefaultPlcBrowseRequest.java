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

import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseRequestInterceptor;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.spi.drivers.functions.PlcBrowser;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;

import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.concurrent.CompletableFuture;

public class DefaultPlcBrowseRequest implements PlcBrowseRequest {

    private final PlcBrowser browser;
    private final LinkedHashMap<String, PlcQuery> queries;

    public DefaultPlcBrowseRequest(PlcBrowser browser, LinkedHashMap<String, PlcQuery> queries) {
        this.browser = browser;
        this.queries = queries;
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> execute() {
        return browser.browse(this);
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> executeWithInterceptor(PlcBrowseRequestInterceptor interceptor) {
        return browser.browseWithInterceptor(this, interceptor);
    }

    @Override
    public LinkedHashSet<String> getQueryNames() {
        return new LinkedHashSet<>(queries.keySet());
    }

    @Override
    public PlcQuery getQuery(String name) {
        return queries.get(name);
    }

    public static class Builder implements PlcBrowseRequest.Builder {

        private final PlcBrowser browser;
        private final PlcTagHandler tagHandler;
        private final LinkedHashMap<String, PlcQuery> queries = new LinkedHashMap<>();

        public Builder(PlcBrowser browser, PlcTagHandler tagHandler) {
            this.browser = browser;
            this.tagHandler = tagHandler;
        }

        @Override
        public PlcBrowseRequest.Builder addQuery(String name, String query) {
            queries.put(name, tagHandler.parseQuery(query));
            return this;
        }

        @Override
        public PlcBrowseRequest build() {
            return new DefaultPlcBrowseRequest(browser, queries);
        }
    }

}
