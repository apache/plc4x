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
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.util.concurrent.CompletableFuture;

public abstract class Plc4xProtocolBase<T> {

    protected ConversationContext<T> conversationContext;

    protected DriverContext driverContext;

    public void setDriverContext(DriverContext driverContext) {
        this.driverContext = driverContext;
    }

    public DriverContext getDriverContext() {
        return driverContext;
    }

    public void setConversationContext(ConversationContext<T> conversationContext) {
        this.conversationContext = conversationContext;
    }

    public ConversationContext<T> getConversationContext() {
        return conversationContext;
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

    public abstract PlcTagHandler getTagHandler();

    /**
     * Default callback, called if an incoming message can't be correlated with an expected response.
     *
     * @param context
     * @param msg
     * @throws Exception
     */
    protected void decode(ConversationContext<T> context, T msg) throws Exception {
    }

    public CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        throw new NotImplementedException("");
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

    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        throw new NotImplementedException("");
    }

    public CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        throw new NotImplementedException("");
    }

    public abstract void close(ConversationContext<T> context);

}
