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

import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.concurrent.CompletableFuture;

public class CachedWriteRequest implements PlcWriteRequest {

    private final CachedPlcConnection parent;
    private final PlcWriteRequest innerRequest;

    public CachedWriteRequest(CachedPlcConnection parent, PlcWriteRequest innerRequest) {
        this.parent = parent;
        this.innerRequest = innerRequest;
    }

    @Override
    public CompletableFuture<? extends PlcWriteResponse> execute() {
        // Only allowed if connection is still active
        return parent.execute(innerRequest);
    }

    @Override
    public int getNumberOfTags() {
        return innerRequest.getNumberOfTags();
    }

    @Override
    public LinkedHashSet<String> getTagNames() {
        return innerRequest.getTagNames();
    }

    @Override
    public PlcTag getTag(String s) {
        return innerRequest.getTag(s);
    }

    @Override
    public List<PlcTag> getTags() {
        return innerRequest.getTags();
    }

    @Override
    public int getNumberOfValues(String name) {
        return innerRequest.getNumberOfValues(name);
    }

    @Override
    public PlcValue getPlcValue(String name) {
        return innerRequest.getPlcValue(name);
    }
}
