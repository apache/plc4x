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
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.model.Address;

import java.util.List;
import java.util.Objects;

public class PlcWriteRequest extends PlcRequest<WriteRequestItem<?>> {

    public PlcWriteRequest() {
    }

    public PlcWriteRequest(WriteRequestItem<?> requestItem) {
        addItem(requestItem);
    }

    @SafeVarargs
    public <T> PlcWriteRequest(Class<T> dataType, Address address, T... values) {
        addItem(new WriteRequestItem<>(dataType, address, values));
    }

    public PlcWriteRequest(List<WriteRequestItem<?>> requestItems) {
        super(requestItems);
    }

    public static PlcWriteRequest.Builder builder() {
        return new Builder();
    }

    public static class Builder extends PlcRequest.Builder<WriteRequestItem<?>> {

        @SuppressWarnings("unchecked")
        public <T> PlcWriteRequest.Builder addItem(Address address, T value) {
            Objects.requireNonNull(value);
            checkType(value.getClass());
            requests.add(new WriteRequestItem<>((Class<T>) value.getClass(), address, (T) value));
            return this;
        }

        @SafeVarargs
        public final <T> PlcWriteRequest.Builder addItem(Class<T> dataType, Address address, T... values) {
            checkType(dataType);
            requests.add(new WriteRequestItem<>(dataType, address, values));
            return this;
        }

        public final PlcWriteRequest.Builder addItem(WriteRequestItem<?> writeRequestItem) {
            checkType(writeRequestItem.getDatatype());
            requests.add(writeRequestItem);
            return this;
        }

        public final PlcWriteRequest build() {
            if (requests.isEmpty()) {
                throw new IllegalStateException("No requests added");
            }
            PlcWriteRequest plcWriteRequest;
            if (mixed) {
                plcWriteRequest = new PlcWriteRequest();
            } else {
                plcWriteRequest = new TypeSafePlcWriteRequest<>(firstType);
            }
            for (WriteRequestItem request : requests) {
                plcWriteRequest.addItem(request);
            }
            return plcWriteRequest;
        }

        @SuppressWarnings("unchecked")
        public final <T> TypeSafePlcWriteRequest<T> build(Class<T> type) {
            if (firstType != type) {
                throw new ClassCastException("Incompatible type " + type + ". Required " + firstType);
            }
            if (mixed) {
                throw new IllegalStateException("Mixed types contained");
            }
            return (TypeSafePlcWriteRequest<T>) build();
        }

    }
}
