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
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.model.Address;

import java.util.List;

public class PlcReadRequest extends PlcRequest<ReadRequestItem<?>> {

    public PlcReadRequest() {
    }

    public PlcReadRequest(ReadRequestItem<?> requestItem) {
        requestItems.add(requestItem);
    }

    public PlcReadRequest(Class<?> dataType, Address address) {
        addItem(new ReadRequestItem<>(dataType, address));
    }

    public PlcReadRequest(Class<?> dataType, Address address, int size) {
        addItem(new ReadRequestItem<>(dataType, address, size));
    }

    public PlcReadRequest(List<ReadRequestItem<?>> requestItems) {
        super(requestItems);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static class Builder extends PlcRequest.Builder<ReadRequestItem<?>> {

        public Builder addItem(Class<?> dataType, Address address) {
            checkType(dataType);
            requests.add(new ReadRequestItem<>(dataType, address));
            return this;
        }

        public Builder addItem(Class<?> dataType, Address address, int size) {
            checkType(dataType);
            requests.add(new ReadRequestItem<>(dataType, address, size));
            return this;
        }

        public Builder addItem(ReadRequestItem readRequestItem) {
            checkType(readRequestItem.getDatatype());
            requests.add(readRequestItem);
            return this;
        }

        @SuppressWarnings("unchecked")
        public PlcReadRequest build() {
            if (requests.isEmpty()) {
                throw new IllegalStateException("No requests added");
            }
            PlcReadRequest plcReadRequest;
            if (mixed) {
                plcReadRequest = new PlcReadRequest();
            } else {
                plcReadRequest = new TypeSafePlcReadRequest<>(firstType);
            }
            for (ReadRequestItem request : requests) {
                plcReadRequest.addItem(request);
            }
            return plcReadRequest;
        }

        @SuppressWarnings("unchecked")
        public <T> TypeSafePlcReadRequest<T> build(Class<T> type) {
            if (firstType != type) {
                throw new ClassCastException("Incompatible type " + type + ". Required " + firstType);
            }
            if (mixed) {
                throw new IllegalStateException("Mixed types contained");
            }
            return (TypeSafePlcReadRequest<T>) build();
        }

    }

}

