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

import java.util.LinkedList;
import java.util.List;

public class CheckedPlcWriteRequest<T> extends PlcWriteRequest {

    private final List<WriteRequestItem<T>> requestItems;

    private final Class<T> datatype;

    public CheckedPlcWriteRequest(Class<T> type) {
        this.datatype = type;
        this.requestItems = new LinkedList<>();
    }

    public CheckedPlcWriteRequest(Class<T> dataType, Address address, T... values) {
        this(dataType);
        addItem(new WriteRequestItem<>(dataType, address, values));
    }

    public CheckedPlcWriteRequest(Class<T> dataType, List<WriteRequestItem<T>> requestItems) {
        this(dataType);
        this.requestItems.addAll(requestItems);
    }

    public void addCheckedItem(WriteRequestItem<T> writeRequestItem) {
        addItem(writeRequestItem);
    }

    @SuppressWarnings("unchecked")
    public void addItem(WriteRequestItem<?> writeRequestItem) {
        if (writeRequestItem == null) {
            return;
        }
        if (writeRequestItem.getDatatype() != datatype) {
            throw new IllegalArgumentException("Incompatible datatype " + writeRequestItem.getDatatype());
        }
        super.addItem(writeRequestItem);
    }

    public List<WriteRequestItem<T>> getCheckedRequestItems() {
        return requestItems;
    }

    @SuppressWarnings("unchecked")
    @Override
    public List<WriteRequestItem<?>> getRequestItems() {
        return (List) getCheckedRequestItems();
    }

}
