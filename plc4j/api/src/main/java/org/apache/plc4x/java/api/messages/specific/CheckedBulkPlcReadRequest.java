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

import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.model.Address;

import java.util.LinkedList;
import java.util.List;

public class CheckedBulkPlcReadRequest<T> extends BulkPlcReadRequest {

    private final List<ReadRequestItem<T>> readRequestItems;

    private Class<T> datatype;

    public CheckedBulkPlcReadRequest(Class<T> type) {
        this.datatype = type;
        this.readRequestItems = new LinkedList<>();
    }

    public CheckedBulkPlcReadRequest(Class<T> dataType, Address address) {
        this(dataType);
        addItem(new ReadRequestItem<>(dataType, address));
    }

    public CheckedBulkPlcReadRequest(Class<T> dataType, Address address, int size) {
        this(dataType);
        addItem(new ReadRequestItem<>(dataType, address, size));
    }

    public void addCheckedItem(ReadRequestItem<T> readRequestItem) {
        readRequestItems.add(readRequestItem);
    }

    @SuppressWarnings("unchecked")
    public void addItem(ReadRequestItem<?> readRequestItem) {
        if (readRequestItem == null) {
            return;
        }
        if (readRequestItem.getDatatype() != datatype) {
            throw new IllegalArgumentException("Incompatible datatype " + readRequestItem.getDatatype());
        }
        readRequestItems.add((ReadRequestItem<T>) readRequestItem);
    }

    public List<? extends ReadRequestItem<T>> getReadRequestItems() {
        return readRequestItems;
    }

}
