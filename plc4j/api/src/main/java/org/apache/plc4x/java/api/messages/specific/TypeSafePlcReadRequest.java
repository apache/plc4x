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

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.model.Address;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

public class TypeSafePlcReadRequest<T> extends PlcReadRequest {

    private Class<T> dataType;

    public TypeSafePlcReadRequest(Class<T> dataType) {
        this.dataType = dataType;
    }

    public TypeSafePlcReadRequest(Class<T> dataType, PlcReadRequest plcReadRequest) {
        this(dataType);
        for (ReadRequestItem<?> readRequestItem : plcReadRequest.getRequestItems()) {
            addItem(readRequestItem);
        }
    }

    public TypeSafePlcReadRequest(Class<T> dataType, Address address) {
        this(dataType);
        addItem(new ReadRequestItem<>(dataType, address));
    }

    public TypeSafePlcReadRequest(Class<T> dataType, Address address, int size) {
        this(dataType);
        addItem(new ReadRequestItem<>(dataType, address, size));
    }

    public TypeSafePlcReadRequest(Class<T> dataType, ReadRequestItem<T>... requestItems) {
        this(dataType);
        Objects.requireNonNull(requestItems);
        for (ReadRequestItem<T> readRequestItem : requestItems) {
            addItem(readRequestItem);
        }
    }

    @Override
    public void addItem(ReadRequestItem<?> readRequestItem) {
        Objects.requireNonNull(readRequestItem);
        if (readRequestItem.getDatatype() != dataType) {
            throw new IllegalArgumentException("Unexpected data type " + readRequestItem.getDatatype() + " on readRequestItem. Expected " + dataType);
        }
        super.addItem(readRequestItem);
    }

    @SuppressWarnings("unchecked")
    public List<ReadRequestItem<T>> getCheckedReadRequestItems() {
        return (List) getRequestItems();
    }

    @SuppressWarnings("unchecked")
    @Override
    public Optional<ReadRequestItem<T>> getRequestItem() {
        return (Optional<ReadRequestItem<T>>) super.getRequestItem();
    }

    public Class<T> getDataType() {
        return dataType;
    }
}
