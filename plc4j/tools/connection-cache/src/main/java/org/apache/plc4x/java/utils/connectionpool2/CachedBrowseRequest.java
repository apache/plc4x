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
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;

import java.util.concurrent.CompletableFuture;

public class CachedBrowseRequest implements PlcBrowseRequest {

    private final CachedPlcConnection parent;
    private final PlcBrowseRequest innerRequest;

    public CachedBrowseRequest(CachedPlcConnection parent, PlcBrowseRequest innerRequest) {
        this.parent = parent;
        this.innerRequest = innerRequest;
    }

    @Override
    public CompletableFuture<? extends PlcBrowseResponse> execute() {
        // Only allowed if connection is still active
        return parent.execute(innerRequest);
    }

}
