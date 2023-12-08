/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.readwrite.utils;

import org.apache.plc4x.java.s7.readwrite.EventType;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

public class S7PlcSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final EventType eventType;
    private final String id;

    public S7PlcSubscriptionHandle(EventType eventType, PlcSubscriber plcSubscriber) {
        super(plcSubscriber);
        this.eventType = eventType;
        this.id = null;
    }

    public S7PlcSubscriptionHandle(String id, EventType eventType, PlcSubscriber plcSubscriber) {
        super(plcSubscriber);
        this.eventType = eventType;
        this.id = id;
    }

    public EventType getEventType() {
        return eventType;
    }

    public String getEventId() {
        return id;
    }

}
