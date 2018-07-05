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

import org.apache.plc4x.java.api.messages.PlcRequestContainer;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.items.RequestItem;

import java.util.Objects;

public class CorrelatedRequestItem<REQUEST_ITEM extends RequestItem<?>> {

    private final int correlationId;

    private final REQUEST_ITEM requestItem;

    private final PlcRequestContainer<?, PlcResponse<?, ?, ?>> plcRequestContainer;

    public CorrelatedRequestItem(int correlationId, REQUEST_ITEM requestItem, PlcRequestContainer<?, PlcResponse<?, ?, ?>> plcRequestContainer) {
        this.correlationId = correlationId;
        this.requestItem = requestItem;
        this.plcRequestContainer = plcRequestContainer;
    }

    public int getCorrelationId() {
        return correlationId;
    }

    public REQUEST_ITEM getRequestItem() {
        return requestItem;
    }

    public PlcRequestContainer<?, PlcResponse<?, ?, ?>> getPlcRequestContainer() {
        return plcRequestContainer;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof CorrelatedRequestItem)) {
            return false;
        }
        CorrelatedRequestItem<?> that = (CorrelatedRequestItem<?>) o;
        return correlationId == that.correlationId &&
            Objects.equals(requestItem, that.requestItem) &&
            Objects.equals(plcRequestContainer, that.plcRequestContainer);
    }

    @Override
    public int hashCode() {

        return Objects.hash(correlationId, requestItem, plcRequestContainer);
    }

    @Override
    public String toString() {
        return "CorrelatedRequestItem{" +
            "correlationId=" + correlationId +
            ", requestItem=" + requestItem +
            ", plcRequestContainer=" + plcRequestContainer +
            '}';
    }
}
