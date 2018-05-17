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
import org.apache.plc4x.java.api.model.Address;

/**
 * This request can be used to pass proprietary protocol specific messages to the underlying layers.
 *
 * @param <CUSTOM_REQUEST> request
 */
public class PlcProprietaryRequest<CUSTOM_REQUEST> extends PlcRequest<PlcProprietaryRequest.DummyRequestItem> {

    private final CUSTOM_REQUEST request;

    public PlcProprietaryRequest(CUSTOM_REQUEST request) {
        this.request = request;
    }

    public CUSTOM_REQUEST getRequest() {
        return request;
    }

    protected static class DummyRequestItem extends RequestItem<Void> {

        public DummyRequestItem() {
            super(Void.class, new Address() {
            });
        }
    }

    @Override
    public String toString() {
        return "PlcProprietaryRequest{" +
            "request=" + request +
            "} " + super.toString();
    }
}
