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

import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.model.Address;

import java.util.Collections;
import java.util.List;

public class SinglePlcWriteRequest<T> implements PlcWriteRequest {

    private WriteRequestItem<T> writeRequestItem;

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
    public void addItem(WriteRequestItem writeRequestItem) {
        if (this.writeRequestItem != null && writeRequestItem != null) {
            throw new IllegalStateException(SinglePlcReadRequest.class.getName() + " can only contain on readRequestItem");
        }
        this.writeRequestItem = (WriteRequestItem<T>) writeRequestItem;
    }

    @Override
    public List<WriteRequestItem<T>> getRequestItems() {
        return (writeRequestItem != null) ? Collections.singletonList(writeRequestItem) : Collections.emptyList();
    }

    public WriteRequestItem<T> getWriteRequestItem() {
        return writeRequestItem;
    }

    public void setWriteRequestItem(WriteRequestItem<T> writeRequestItem) {
        this.writeRequestItem = writeRequestItem;
    }

    public int getNumberOfItems() {
        return writeRequestItem != null ? 1 : 0;
    }
}
