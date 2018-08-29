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

import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;

public class MockConnection extends org.apache.plc4x.java.base.connection.MockConnection implements PlcReader, PlcWriter {

    private final String url;
    private final PlcAuthentication authentication;
    private final Map<PlcField, Object> dataValueMap = new HashMap<>();
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
    public PlcField prepareField(String fieldString) {
        return new MockField(fieldString);
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
        List<PlcReadResponseItem<Object>> responseItems = new LinkedList<>();
        for (PlcReadRequestItem<?> reqItem : readRequest.getRequestItems()) {
            @SuppressWarnings("unchecked")
            PlcReadRequestItem<Object> requestItem = reqItem;
            PlcReadResponseItem<Object> responseItem = new PlcReadResponseItem<>(requestItem, PlcResponseCode.OK,
                Collections.singletonList(getDataValue(requestItem.getField())));
            responseItems.add(responseItem);
        }
        PlcReadResponse response;
        if (readRequest instanceof TypeSafePlcReadRequest) {
            @SuppressWarnings("unchecked")
            TypeSafePlcReadRequest<Object> readReq = (TypeSafePlcReadRequest<Object>) readRequest;
            response = new TypeSafePlcReadResponse<>(readReq, responseItems);
        } else {
            response = new PlcReadResponse(readRequest, responseItems);
        }
        return CompletableFuture.completedFuture(response);
    }

    @SuppressWarnings("unchecked")
    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        curWriteCnt++;
        if (writeExceptionTriggerCount > 0 && curWriteCnt == writeExceptionTriggerCount) {
            curWriteCnt = 0;
            CompletableFuture<PlcWriteResponse> cf = new CompletableFuture<>();
            cf.completeExceptionally(new PlcIoException(writeExceptionMsg));
            return cf;
        }
        List<PlcWriteResponseItem<Object>> responseItems = new LinkedList<>();
        for (PlcWriteRequestItem<?> reqItem : writeRequest.getRequestItems()) {
            PlcWriteRequestItem<Object> requestItem = reqItem;
            setDataValue(requestItem.getField(), requestItem.getValues());
            PlcWriteResponseItem<Object> responseItem = new PlcWriteResponseItem<>(requestItem, PlcResponseCode.OK);
            responseItems.add(responseItem);
        }
        PlcWriteResponse response;
        if (writeRequest instanceof TypeSafePlcWriteRequest) {
            TypeSafePlcWriteRequest<Object> writeReq = (TypeSafePlcWriteRequest<Object>) writeRequest;
            response = new TypeSafePlcWriteResponse<>(writeReq, responseItems);
        } else {
            response = new PlcWriteResponse(writeRequest, responseItems);
        }

        return CompletableFuture.completedFuture(response);
    }

    public void setDataValue(PlcField field, Object o) {
        dataValueMap.put(field, o);
    }

    public Object getDataValue(PlcField field) {
        return dataValueMap.get(field);
    }

    public Map<PlcField, Object> getAllDataValues() {
        return dataValueMap;
    }

    public void clearAllDataValues() {
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
