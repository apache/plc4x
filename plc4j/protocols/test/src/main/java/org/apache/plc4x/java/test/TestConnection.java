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

import org.apache.plc4x.java.api.connection.*;
import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.messages.items.WriteResponseItem;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.types.ResponseCode;

import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
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
    public Address parseAddress(String addressString) throws PlcInvalidAddressException {
        if (!TestAddress.isValid(addressString)) {
            throw new PlcInvalidAddressException("Address must contain a single word.");
        }

        return TestAddress.of(addressString);
    }

    @Override
    public Optional<PlcLister> getLister() {
        return Optional.empty();
    }

    @Override
    public Optional<PlcReader> getReader() {
        return Optional.of(this);
    }

    @Override
    public Optional<PlcWriter> getWriter() {
        return Optional.of(this);
    }

    @Override
    public Optional<PlcSubscriber> getSubscriber() {
        return Optional.empty(); // TODO: implement this
    }

    @Override
    @SuppressWarnings("unchecked")
    public CompletableFuture<? extends PlcReadResponse> read(PlcReadRequest readRequest) {
        List<ReadResponseItem<?>> responseItems = new LinkedList<>();
        for (ReadRequestItem<?> requestItem : readRequest.getRequestItems()) {
            TestAddress address = (TestAddress) requestItem.getAddress();
            Optional<?> value = device.get(requestItem.getDatatype(), address);
            ReadResponseItem<?> responseItem;
            if (value.isPresent()) {
                responseItem = new ReadResponseItem(requestItem, ResponseCode.OK, Collections.singletonList(value.get()));
            } else {
                responseItem = new ReadResponseItem(requestItem, ResponseCode.NOT_FOUND, Collections.emptyList());
            }
            responseItems.add(responseItem);
        }
        PlcReadResponse response = new PlcReadResponse(readRequest, responseItems);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public <T> CompletableFuture<TypeSafePlcReadResponse<T>> read(TypeSafePlcReadRequest<T> readRequest) {
        return null; // TODO: implement this
    }

    @Override
    @SuppressWarnings("unchecked")
    public CompletableFuture<? extends PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        List<WriteResponseItem<?>> responseItems = new LinkedList<>();
        for (WriteRequestItem<?> requestItem : writeRequest.getRequestItems()) {
            TestAddress address = (TestAddress) requestItem.getAddress();
            Object value = requestItem.getValues().get(0);
            device.set(address, value);
            WriteResponseItem<?> responseItem = new WriteResponseItem(requestItem, ResponseCode.OK);
            responseItems.add(responseItem);
        }
        PlcWriteResponse response = new PlcWriteResponse(writeRequest, responseItems);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public <T> CompletableFuture<TypeSafePlcWriteResponse<T>> write(TypeSafePlcWriteRequest<T> writeRequest) {
        return null; // TODO: implement this
    }

    @Override
    public String toString() {
        return String.format("test:%s", device);
    }

}
