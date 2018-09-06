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
package org.apache.plc4x.java.base.messages;

import org.apache.plc4x.java.api.messages.PlcMessageBuilder;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Collection;
import java.util.LinkedHashSet;

public class DefaultPlcSubscriptionResponse implements PlcSubscriptionResponse {

    @Override
    public PlcSubscriptionHandle getSubscriptionHandle(String name) {
        return null;
    }

    @Override
    public Collection<String> getFieldNames() {
        return null;
    }

    @Override
    public PlcField getField(String name) {
        return null;
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        return null;
    }

    @Override
    public PlcSubscriptionRequest getRequest() {
        return null;
    }

}
