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
package org.apache.plc4x.java.opcua.protocol;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.opcua.connection.OpcuaTcpPlcConnection;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.eclipse.milo.opcua.sdk.client.api.subscriptions.UaMonitoredItem;
import org.eclipse.milo.opcua.stack.core.types.builtin.DataValue;
import org.eclipse.milo.opcua.stack.core.types.builtin.StatusCode;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.UInteger;

import java.time.Instant;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.function.Consumer;

/**
 */
public class OpcuaSubsriptionHandle implements PlcSubscriptionHandle {

    private Set<Consumer<PlcSubscriptionEvent>> consumers = new HashSet<>();
    private String fieldName;
    private UInteger clientHandle;

    /**
     * @param fieldName    corresponding map key in the PLC4X request/reply map
     * @param clientHandle
     */
    public OpcuaSubsriptionHandle(String fieldName, UInteger clientHandle) {
        this.fieldName = fieldName;
        this.clientHandle = clientHandle;
    }

    public UInteger getClientHandle() {
        return clientHandle;
    }

    /**
     * @param item
     * @param value
     */
    public void onSubscriptionValue(UaMonitoredItem item, DataValue value) {
        consumers.forEach(plcSubscriptionEventConsumer -> {
            PlcResponseCode resultCode = PlcResponseCode.OK;
            PlcValue stringItem = null;
            if (value.getStatusCode() != StatusCode.GOOD) {
                resultCode = PlcResponseCode.NOT_FOUND;
            } else {
                stringItem = OpcuaTcpPlcConnection.encodePlcValue(value);

            }
            Map<String, ResponseItem<PlcValue>> fields = new HashMap<>();
            ResponseItem<PlcValue> newPair = new ResponseItem<>(resultCode, stringItem);
            fields.put(fieldName, newPair);
            PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), fields);
            plcSubscriptionEventConsumer.accept(event);
        });
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        consumers.add(consumer);
        return null;
//        return () -> consumers.remove(consumer);
    }

}
