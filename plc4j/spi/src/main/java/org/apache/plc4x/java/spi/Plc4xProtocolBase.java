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
package org.apache.plc4x.java.spi;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.util.concurrent.CompletableFuture;

public abstract class Plc4xProtocolBase<T> {

    protected ConversationContext<T> context;

    protected DriverContext driverContext;

    public void setDriverContext(DriverContext driverContext) {
        this.driverContext = driverContext;
    }

    public DriverContext getDriverContext() {
        return driverContext;
    }

    public void setContext(ConversationContext<T> context) {
        this.context = context;
    }

    public void onConnect(ConversationContext<T> context) {
        // Intentionally do nothing here
    }

    public void onDisconnect(ConversationContext<T> context) {
        // Intentionally do nothing here
    }

    public void onDiscover(ConversationContext<T> context) {
        // Intentionally do nothing here
    }

    /**
     * TODO document me
     * <p>
     * Can be used for non requested incoming messages
     *
     * @param context
     * @param msg
     * @throws Exception
     */
    protected void decode(ConversationContext<T> context, T msg) throws Exception {
    }

    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        throw new NotImplementedException("");
    }

    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        throw new NotImplementedException("");
    }

    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        throw new NotImplementedException("");
    }

    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        throw new NotImplementedException("");
    }

    public abstract void close(ConversationContext<T> context);

}
