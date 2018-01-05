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
package org.apache.plc4x.java.api.messages.specific;

import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.model.Address;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

public class SinglePlcWriteRequest<T> implements PlcWriteRequest {

    private WriteRequestItem<T> requestItem;

    public SinglePlcWriteRequest() {
    }

    public SinglePlcWriteRequest(Class<T> dataType, Address address, T value) {
        addItem(new WriteRequestItem<>(dataType, address, value));
    }

    public SinglePlcWriteRequest(Class<T> dataType, Address address, T[] values) {
        addItem(new WriteRequestItem<>(dataType, address, values));
    }

    @Override
    @SuppressWarnings("unchecked")
    public void addItem(WriteRequestItem requestItem) {
        if (this.requestItem != null && requestItem != null) {
            throw new IllegalStateException(SinglePlcReadRequest.class.getName() + " can only contain on readRequestItem");
        }
        this.requestItem = (WriteRequestItem<T>) requestItem;
    }

    @Override
    public List<WriteRequestItem<T>> getRequestItems() {
        return (requestItem != null) ? Collections.singletonList(requestItem) : Collections.emptyList();
    }

    public Optional<WriteRequestItem<T>> getRequestItem() {
        return Optional.ofNullable(requestItem);
    }

    public void setWriteRequestItem(WriteRequestItem<T> requestItem) {
        this.requestItem = requestItem;
    }

    public int getNumberOfItems() {
        return requestItem != null ? 1 : 0;
    }
}
