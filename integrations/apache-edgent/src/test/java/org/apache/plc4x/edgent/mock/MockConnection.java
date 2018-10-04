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
package org.apache.plc4x.edgent.mock;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class MockConnection extends org.apache.plc4x.java.base.connection.MockConnection implements PlcReader, PlcWriter {

    private final String url;
    private final PlcAuthentication authentication;
    private final Map<PlcField, FieldItem<?>> dataValueMap = new HashMap<>();
    private long curReadCnt;
    private int readExceptionTriggerCount;
    private String readExceptionMsg;
    private long curWriteCnt;
    private int writeExceptionTriggerCount;
    private String writeExceptionMsg;
    private boolean connected;

    public MockConnection(String url) {
        this(url, null);
    }

    public MockConnection(String url, PlcAuthentication authentication) {
        super();
        this.url = url;
        this.authentication = authentication;
        this.connected = false;
    }

    public String getUrl() {
        return url;
    }

    @Override
    public Optional<PlcReadRequest.Builder> readRequestBuilder() {
        return Optional.of(new DefaultPlcReadRequest.Builder(this, new MockFieldHandler()));
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        curReadCnt++;
        if (readExceptionTriggerCount > 0 && curReadCnt == readExceptionTriggerCount) {
            curReadCnt = 0;
            CompletableFuture<PlcReadResponse> cf = new CompletableFuture<>();
            cf.completeExceptionally(new PlcIoException(readExceptionMsg));
            return cf;
        }
        Map<String, Pair<PlcResponseCode, FieldItem>> fields = new LinkedHashMap<>();
        for (String fieldName : readRequest.getFieldNames()) {
            PlcField field = readRequest.getField(fieldName);
            fields.put(fieldName, new ImmutablePair<>(PlcResponseCode.OK, getFieldItem(field)));
        }
        PlcReadResponse response = new DefaultPlcReadResponse((InternalPlcReadRequest) readRequest, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public Optional<PlcWriteRequest.Builder> writeRequestBuilder() {
        return Optional.of(new DefaultPlcWriteRequest.Builder(this, new MockFieldHandler()));
    }

    @SuppressWarnings("unchecked")
    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        DefaultPlcWriteRequest defaultPlcWriteRequest = (DefaultPlcWriteRequest) writeRequest;
        curWriteCnt++;
        if (writeExceptionTriggerCount > 0 && curWriteCnt == writeExceptionTriggerCount) {
            curWriteCnt = 0;
            CompletableFuture<PlcWriteResponse> cf = new CompletableFuture<>();
            cf.completeExceptionally(new PlcIoException(writeExceptionMsg));
            return cf;
        }
        Map<String, PlcResponseCode> fields = new LinkedHashMap<>();
        for (String fieldName : defaultPlcWriteRequest.getFieldNames()) {
            PlcField field = defaultPlcWriteRequest.getField(fieldName);
            setFieldItem(field, defaultPlcWriteRequest.getFieldItem(fieldName));
            fields.put(fieldName, PlcResponseCode.OK);
        }
        PlcWriteResponse response = new DefaultPlcWriteResponse(defaultPlcWriteRequest, fields);

        return CompletableFuture.completedFuture(response);
    }

    public void setFieldItem(PlcField field, FieldItem<?> fieldItem) {
        dataValueMap.put(field, fieldItem);
    }

    public FieldItem<?> getFieldItem(PlcField field) {
        return dataValueMap.get(field);
    }

    public Map<PlcField, FieldItem<?>> getAllFieldItems() {
        return dataValueMap;
    }

    public void clearAllFieldItems() {
        dataValueMap.clear();
    }

    public void setReadException(int readTriggerCount, String msg) {
        readExceptionTriggerCount = readTriggerCount;
        readExceptionMsg = msg;
        curReadCnt = 0;
    }

    public void setWriteException(int writeTriggerCount, String msg) {
        writeExceptionTriggerCount = writeTriggerCount;
        writeExceptionMsg = msg;
        curWriteCnt = 0;
    }

}
