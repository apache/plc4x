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
import org.apache.plc4x.java.api.messages.specific.BulkPlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.CheckedBulkPlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.SinglePlcReadRequest;
import org.apache.plc4x.java.api.model.Address;

import java.util.LinkedList;
import java.util.List;
import java.util.Optional;

public interface PlcReadRequest extends PlcRequest {
    void addItem(ReadRequestItem<?> readRequestItem);

    List<? extends ReadRequestItem<?>> getReadRequestItems();

    default Optional<? extends ReadRequestItem<?>> getReadRequestItem() {
        if (getNumberOfItems() > 1) {
            throw new IllegalStateException("too many items " + getNumberOfItems());
        }
        if (getNumberOfItems() < 1) {
            return Optional.empty();
        }
        return Optional.<ReadRequestItem<?>>of(getReadRequestItems().get(0));
    }

    default int getNumberOfItems() {
        return getReadRequestItems().size();
    }

    static Builder builder() {
        return new Builder();
    }

    class Builder {

        private Class firstType;

        private boolean mixed = false;

        private List<ReadRequestItem> requests = new LinkedList<>();

        public <T> Builder addItem(Class<T> dataType, Address address) {
            checkType(dataType);
            requests.add(new ReadRequestItem<>(dataType, address));
            return this;
        }

        public <T> Builder addItem(Class<T> dataType, Address address, int size) {
            checkType(dataType);
            requests.add(new ReadRequestItem<>(dataType, address, size));
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
        public PlcReadRequest build() {
            if (requests.size() < 1) {
                throw new IllegalStateException("No requests added");
            }
            if (requests.size() < 2) {
                SinglePlcReadRequest<?> singlePlcReadRequest = new SinglePlcReadRequest<>();
                singlePlcReadRequest.addItem(requests.get(0));
                return singlePlcReadRequest;
            }
            PlcReadRequest plcReadRequest;
            if (mixed) {
                plcReadRequest = new BulkPlcReadRequest();
            } else {
                plcReadRequest = new CheckedBulkPlcReadRequest<>(firstType);
            }
            for (ReadRequestItem request : requests) {
                plcReadRequest.addItem(request);
            }
            return plcReadRequest;
        }

        @SuppressWarnings("unchecked")
        public BulkPlcReadRequest buildBulk() {
            if (requests.size() < 2) {
                throw new IllegalStateException("Bulk request needs more than one request");
            }
            return (BulkPlcReadRequest) build();
        }

        @SuppressWarnings("unchecked")
        public <T> SinglePlcReadRequest<T> build(Class<T> type) {
            if (requests.size() != 1) {
                throw new IllegalStateException("Checked request needs exactly one request");
            }
            if (firstType != type) {
                throw new ClassCastException("Incompatible type " + type + ". Required " + firstType);
            }
            return (SinglePlcReadRequest<T>) build();
        }

        @SuppressWarnings("unchecked")
        public <T> CheckedBulkPlcReadRequest<T> buildBulk(Class<T> type) {
            if (requests.size() < 2) {
                throw new IllegalStateException("Checked bulk request needs more than one request");
            }
            if (firstType != type) {
                throw new ClassCastException("Incompatible type " + type + ". Required " + firstType);
            }
            if (mixed) {
                throw new IllegalStateException("Mixed types contained");
            }
            return (CheckedBulkPlcReadRequest<T>) build();
        }

    }

}

