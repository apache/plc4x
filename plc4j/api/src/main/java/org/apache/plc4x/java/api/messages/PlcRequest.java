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

import org.apache.plc4x.java.api.messages.items.RequestItem;

import java.util.LinkedList;
import java.util.List;
import java.util.Optional;

/**
 * Base type for all messages sent from the plc4x system to a connected plc.
 */
public abstract class PlcRequest<REQUEST_ITEM extends RequestItem> implements PlcMessage {

    protected final List<REQUEST_ITEM> requestItems;

    public PlcRequest() {
        this.requestItems = new LinkedList<>();
    }

    public PlcRequest(List<REQUEST_ITEM> requestItems) {
        this.requestItems = requestItems;
    }

    public void addItem(REQUEST_ITEM readRequestItem) {
        getRequestItems().add(readRequestItem);
    }

    public List<REQUEST_ITEM> getRequestItems() {
        return requestItems;
    }

    public Optional<? extends REQUEST_ITEM> getRequestItem() {
        if (isMultiValue()) {
            throw new IllegalStateException("too many items " + getNumberOfItems());
        }
        if (isEmpty()) {
            return Optional.empty();
        }
        return Optional.of(getRequestItems().get(0));
    }

    public void setRequestItem(REQUEST_ITEM requestItem) {
        if (isMultiValue()) {
            throw new IllegalStateException("too many items " + getNumberOfItems());
        }
        addItem(requestItem);
    }

    public int getNumberOfItems() {
        return getRequestItems().size();
    }

    public boolean isMultiValue() {
        return getNumberOfItems() > 1;
    }

    public boolean isEmpty() {
        return getNumberOfItems() < 1;
    }

    public abstract static class Builder<REQUEST_ITEM> {

        Class firstType;

        boolean mixed = false;

        List<REQUEST_ITEM> requests = new LinkedList<>();

        void checkType(Class dataType) {
            if (firstType == null) {
                firstType = dataType;
            }
            if (firstType != dataType) {
                mixed = true;
            }
        }
    }
}
