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
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;

public class CachedReadRequestBuilder implements PlcReadRequest.Builder {

    private final CachedPlcConnection parent;
    private final PlcReadRequest.Builder builder;

    public CachedReadRequestBuilder(CachedPlcConnection parent, PlcReadRequest.Builder builder) {
        this.parent = parent;
        this.builder = builder;
    }

    @Override
    public PlcReadRequest.Builder addItem(String s, String s1) {
        builder.addItem(s, s1);
        return this;
    }

    @Override
    public PlcReadRequest.Builder addItem(String name, PlcField fieldQuery) {
        builder.addItem(name, fieldQuery);
        return this;
    }

    @Override
    public PlcReadRequest build() {
        return new CachedReadRequest(parent, builder.build());
    }

}
