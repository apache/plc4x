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

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.concurrent.CompletableFuture;

public class CachedReadRequest implements PlcReadRequest {

    private final CachedPlcConnection parent;
    private final PlcReadRequest innerRequest;

    public CachedReadRequest(CachedPlcConnection parent, PlcReadRequest innerRequest) {
        this.parent = parent;
        this.innerRequest = innerRequest;
    }

    @Override
    public CompletableFuture<? extends PlcReadResponse> execute() {
        // Only allowed if connection is still active
        return parent.execute(innerRequest);
    }

    @Override
    public int getNumberOfFields() {
        return innerRequest.getNumberOfFields();
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        return innerRequest.getFieldNames();
    }

    @Override
    public PlcField getField(String s) {
        return innerRequest.getField(s);
    }

    @Override
    public List<PlcField> getFields() {
        return innerRequest.getFields();
    }
}
