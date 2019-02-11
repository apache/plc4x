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

package org.apache.plc4x.sandbox.java.dynamic.utils;

import org.apache.commons.scxml2.env.AbstractSCXMLListener;
import org.apache.commons.scxml2.model.EnterableState;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;

public class RequestRegistry extends AbstractSCXMLListener {

    private final long maxValue;
    private AtomicLong idGenerator = new AtomicLong(0);

    // TODO: Add some timeout stuff here ...
    private ConcurrentHashMap<String, PlcRequestContainer> requestContainers;

    public RequestRegistry(long maxValue) {
        this.maxValue = maxValue;
        requestContainers = new ConcurrentHashMap<>();
    }

    public String generateRequestId() {
        // If we reached the max value, reset to 0.
        if(idGenerator.get() == maxValue) {
            idGenerator.set(0);
        }
        return Long.toString(idGenerator.getAndIncrement());
    }

    public void addContainer(String requestId, PlcRequestContainer container) {
        requestContainers.put(requestId, container);
    }

    public PlcRequestContainer removeContainer(String requestId) {
        if(requestId != null) {
            return requestContainers.remove(requestId);
        }
        return null;
    }

    @Override
    public void onExit(EnterableState state) {
        super.onExit(state);
    }

}
