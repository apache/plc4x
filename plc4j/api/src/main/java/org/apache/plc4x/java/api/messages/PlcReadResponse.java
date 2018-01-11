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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

public class PlcReadResponse implements PlcResponse {

    private final PlcReadRequest request;

    private final List<? extends ReadResponseItem<?>> responseItems;

    public PlcReadResponse(PlcReadRequest request, ReadResponseItem<?> responseItems) {
        this.request = request;
        this.responseItems = Collections.singletonList(responseItems);
    }

    public PlcReadResponse(PlcReadRequest request, List<? extends ReadResponseItem<?>> responseItems) {
        this.request = request;
        this.responseItems = responseItems;
    }

    public PlcReadRequest getRequest() {
        return request;
    }

    public List<? extends ReadResponseItem<?>> getResponseItems() {
        return responseItems;
    }

    public Optional<? extends ReadResponseItem<?>> getResponseItem() {
        if (isMultiValue()) {
            throw new IllegalStateException("too many items " + getNumberOfItems());
        }
        if (isEmpty()) {
            return Optional.empty();
        }
        return Optional.<ReadResponseItem<?>>of(getResponseItems().get(0));
    }

    public int getNumberOfItems() {
        return getResponseItems().size();
    }

    public boolean isMultiValue() {
        return getNumberOfItems() > 1;
    }

    public boolean isEmpty() {
        return getNumberOfItems() < 1;
    }

    @SuppressWarnings("unchecked")
    public <T> Optional<ReadResponseItem<T>> getValue(ReadRequestItem<T> item) {
        return getResponseItems().stream()
            .filter(x -> x.getRequestItem().equals(item))
            .map(e -> (ReadResponseItem<T>) e)
            .findAny();
    }
}
