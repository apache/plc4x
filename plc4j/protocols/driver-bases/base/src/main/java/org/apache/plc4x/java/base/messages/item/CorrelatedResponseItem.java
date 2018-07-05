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
package org.apache.plc4x.java.base.messages.item;

import org.apache.plc4x.java.api.messages.items.ResponseItem;

import java.util.Objects;

public class CorrelatedResponseItem<RESPONSE_ITEM extends ResponseItem<?>> {

    private final int correlationId;

    private final RESPONSE_ITEM responseItem;

    public CorrelatedResponseItem(int correlationId, RESPONSE_ITEM responseItem) {
        this.correlationId = correlationId;
        this.responseItem = responseItem;
    }

    public int getCorrelationId() {
        return correlationId;
    }

    public RESPONSE_ITEM getResponseItem() {
        return responseItem;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof CorrelatedResponseItem)) {
            return false;
        }
        CorrelatedResponseItem<?> that = (CorrelatedResponseItem<?>) o;
        return correlationId == that.correlationId &&
            Objects.equals(responseItem, that.responseItem);
    }

    @Override
    public int hashCode() {

        return Objects.hash(correlationId, responseItem);
    }

    @Override
    public String toString() {
        return "CorrelatedResponseItem{" +
            "correlationId=" + correlationId +
            ", responseItem=" + responseItem +
            '}';
    }
}
