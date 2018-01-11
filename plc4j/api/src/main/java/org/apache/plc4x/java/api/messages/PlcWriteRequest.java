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
import org.apache.plc4x.java.api.messages.specific.BulkPlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.CheckedBulkPlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.SinglePlcWriteRequest;
import org.apache.plc4x.java.api.model.Address;

import java.util.LinkedList;
import java.util.List;
import java.util.Optional;

public class PlcWriteRequest implements PlcRequest {

    private final List<WriteRequestItem<?>> requestItems;

    public PlcWriteRequest() {
        this.requestItems = new LinkedList<>();
    }

    public <T> PlcWriteRequest(Class<T> dataType, Address address, T... values) {
        this();
        addItem(new WriteRequestItem<>(dataType, address, values));
    }

    public PlcWriteRequest(List<WriteRequestItem<?>> requestItems) {
        this.requestItems = requestItems;
    }

    public void addItem(WriteRequestItem<?> requestItem) {
        getRequestItems().add(requestItem);
    }

    public List<WriteRequestItem<?>> getRequestItems() {
        return requestItems;
    }

    public Optional<? extends WriteRequestItem<?>> getRequestItem() {
        if (getNumberOfItems() > 1) {
            throw new IllegalStateException("too many items " + getNumberOfItems());
        }
        if (getNumberOfItems() < 1) {
            return Optional.empty();
        }
        return Optional.<WriteRequestItem<?>>of(getRequestItems().get(0));
    }

    public int getNumberOfItems() {
        return getRequestItems().size();
    }

    public static PlcWriteRequest.Builder builder() {
        return new Builder();
    }

    public static class Builder {

        private Class firstType;

        private boolean mixed = false;

        private List<WriteRequestItem> requests = new LinkedList<>();

        public <T> PlcWriteRequest.Builder addItem(Class<T> dataType, Address address, T... values) {
            checkType(dataType);
            requests.add(new WriteRequestItem<>(dataType, address, values));
            return this;
        }

        private void checkType(Class dataType) {
            if (firstType == null) {
                firstType = dataType;
            }
            if (firstType != dataType) {
                mixed = true;
            }
        }

        @SuppressWarnings("unchecked")
        public PlcWriteRequest build() {
            if (requests.size() < 1) {
                throw new IllegalStateException("No requests added");
            }
            if (requests.size() < 2) {
                return new SinglePlcWriteRequest<>(requests.get(0));
            }
            PlcWriteRequest plcWriteRequest;
            if (mixed) {
                plcWriteRequest = new BulkPlcWriteRequest();
            } else {
                plcWriteRequest = new CheckedBulkPlcWriteRequest<>(firstType);
            }
            for (WriteRequestItem request : requests) {
                plcWriteRequest.addItem(request);
            }
            return plcWriteRequest;
        }

        @SuppressWarnings("unchecked")
        public BulkPlcWriteRequest buildBulk() {
            if (requests.size() < 2) {
                throw new IllegalStateException("Bulk request needs more than one request");
            }
            return (BulkPlcWriteRequest) build();
        }

        @SuppressWarnings("unchecked")
        public <T> SinglePlcWriteRequest<T> build(Class<T> type) {
            if (requests.size() != 1) {
                throw new IllegalStateException("Checked request needs exactly one request");
            }
            if (firstType != type) {
                throw new ClassCastException("Incompatible type " + type + ". Required " + firstType);
            }
            return (SinglePlcWriteRequest<T>) build();
        }

        @SuppressWarnings("unchecked")
        public <T> CheckedBulkPlcWriteRequest<T> buildBulk(Class<T> type) {
            if (requests.size() < 2) {
                throw new IllegalStateException("Checked bulk request needs more than one request");
            }
            if (firstType != type) {
                throw new ClassCastException("Incompatible type " + type + ". Required " + firstType);
            }
            if (mixed) {
                throw new IllegalStateException("Mixed types contained");
            }
            return (CheckedBulkPlcWriteRequest<T>) build();
        }

    }
}
