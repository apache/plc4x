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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.model.PlcQuery;

import java.util.LinkedHashSet;
import java.util.concurrent.CompletableFuture;

public interface PlcBrowseRequest extends PlcRequest {

    CompletableFuture<? extends PlcBrowseResponse> execute();

    /**
     * In contrast to the default execute method, the executeWithInterceptor allows passing in a so-called
     * interceptor. This can be used for two different situations:
     * 1. Filter which items go into the final PlcBrowseResponse (Items are added, if interceptor returns true).
     * 2. Allow accessing found items on the fly while the system is still processing, hereby allowing a more
     *    asynchronous consumption of found tags.
     *
     * @param interceptor interceptor for intercepting found items
     * @return future for the final PlcBrowseResponse
     */
    CompletableFuture<? extends PlcBrowseResponse> executeWithInterceptor(PlcBrowseRequestInterceptor interceptor);

    LinkedHashSet<String> getQueryNames();

    PlcQuery getQuery(String name);

    interface Builder extends PlcRequestBuilder {

        @Override
        PlcBrowseRequest build();

        Builder addQuery(String name, String query);

    }

}
