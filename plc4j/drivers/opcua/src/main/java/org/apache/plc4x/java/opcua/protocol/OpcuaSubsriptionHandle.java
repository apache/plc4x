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

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.opcua.connection.OpcuaTcpPlcConnection;
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
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
public class OpcuaSubsriptionHandle implements PlcSubscriptionHandle {
    Set< Consumer<PlcSubscriptionEvent>> consumers = new HashSet<>();
    String fieldName;
    public UInteger getClientHandle() {
        return clientHandle;
    }

    UInteger clientHandle;

    public  OpcuaSubsriptionHandle(String fieldName, UInteger clientHandle){
        this.clientHandle = clientHandle;
    }
    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        consumers.add(consumer);
        return () -> {consumers.remove(consumer);};
    }

    public void onSubscriptionValue(UaMonitoredItem item, DataValue value) {
        consumers.forEach(plcSubscriptionEventConsumer -> {
            PlcResponseCode resultCode = PlcResponseCode.OK;
            BaseDefaultFieldItem stringItem = null;
            if(value.getStatusCode() != StatusCode.GOOD){
                resultCode = PlcResponseCode.NOT_FOUND;
            }else{
                stringItem = OpcuaTcpPlcConnection.encodeFieldItem(value);

            }
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();
            Pair<PlcResponseCode, BaseDefaultFieldItem> newPair = new ImmutablePair<>(resultCode, stringItem);
            fields.put(fieldName, newPair);
            PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), fields);
            plcSubscriptionEventConsumer.accept(event);
        });
    }

}
