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

import org.apache.plc4x.java.api.messages.items.UnsubscriptionRequestItem;
import org.apache.plc4x.java.api.model.SubscriptionHandle;

import java.util.LinkedList;
import java.util.List;
import java.util.Objects;

public class PlcUnsubscriptionRequest implements PlcMessage {

    protected final List<UnsubscriptionRequestItem> requestItems;

    public PlcUnsubscriptionRequest() {
        this.requestItems = new LinkedList<>();
    }

    public PlcUnsubscriptionRequest(List<UnsubscriptionRequestItem> requestItems) {
        Objects.requireNonNull(requestItems, "Request items must not be null");
        this.requestItems = requestItems;
    }

    public void addItem(UnsubscriptionRequestItem unsubscriptionRequestItem) {
        Objects.requireNonNull(unsubscriptionRequestItem, "Request item must not be null");
        getRequestItems().add(unsubscriptionRequestItem);
    }

    public List<UnsubscriptionRequestItem> getRequestItems() {
        return requestItems;
    }

    public int getNumberOfItems() {
        return getRequestItems().size();
    }

    @Override
    public String toString() {
        return "PlcUnsubscriptionRequest{" +
            "requestItems=" + requestItems +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcUnsubscriptionRequest)) {
            return false;
        }
        PlcUnsubscriptionRequest that = (PlcUnsubscriptionRequest) o;
        return Objects.equals(requestItems, that.requestItems);
    }

    @Override
    public int hashCode() {
        return Objects.hash(requestItems);
    }
}
