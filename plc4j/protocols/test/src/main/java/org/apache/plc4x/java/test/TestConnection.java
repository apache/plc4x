/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.test;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcWriter;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

/**
 * Connection to a test device.
 * This class is not thread-safe.
 */
class TestConnection implements PlcConnection, PlcReader, PlcWriter {
    private final TestDevice device;
    private boolean connected = false;

    TestConnection(TestDevice device) {
        this.device = device;
    }

    @Override
    public void connect() {
        connected = true;
    }

    @Override
    public boolean isConnected() {
        return connected;
    }

    @Override
    public void close() {
        connected = false;
    }

    @Override
    public Optional<PlcReadRequest.Builder> readRequestBuilder() {
        return Optional.of(new DefaultPlcReadRequest.Builder(this, new TestFieldHandler()));
    }

    @Override
    public Optional<PlcWriteRequest.Builder> writeRequestBuilder() {
        return Optional.of(new DefaultPlcWriteRequest.Builder(this, new TestFieldHandler()));
    }

    @Override
    public Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        if(!(readRequest instanceof InternalPlcReadRequest)) {
            throw new IllegalArgumentException("Read request doesn't implement InternalPlcReadRequest");
        }
        InternalPlcReadRequest request = (InternalPlcReadRequest) readRequest;
        Map<String, Pair<PlcResponseCode, FieldItem>> fields = new HashMap<>();
        for (String fieldName : request.getFieldNames()) {
            TestField field = (TestField) request.getField(fieldName);
            Optional<FieldItem> fieldItemOptional = device.get(field);
            ImmutablePair<PlcResponseCode, FieldItem> fieldPair;
            boolean present = fieldItemOptional.isPresent();
            fieldPair = present
                ? new ImmutablePair<>(PlcResponseCode.OK, fieldItemOptional.get())
                : new ImmutablePair<>(PlcResponseCode.NOT_FOUND, null);
            fields.put(fieldName, fieldPair);
        }
        PlcReadResponse response = new DefaultPlcReadResponse(request, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        if(!(writeRequest instanceof InternalPlcWriteRequest)) {
            throw new IllegalArgumentException("Read request doesn't implement InternalPlcWriteRequest");
        }
        InternalPlcWriteRequest request = (InternalPlcWriteRequest) writeRequest;
        Map<String, PlcResponseCode> fields = new HashMap<>();
        for (String fieldName : request.getFieldNames()) {
            TestField field = (TestField) request.getField(fieldName);
            FieldItem fieldItem = request.getFieldItem(fieldName);
            device.set(field, fieldItem);
            fields.put(fieldName, PlcResponseCode.OK);
        }
        PlcWriteResponse response = new DefaultPlcWriteResponse(request, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public String toString() {
        return String.format("test:%s", device);
    }

}
