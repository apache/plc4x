/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.connection;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcUnsubscriptionRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.PlcWriter;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.api.value.PlcValueHandler;


import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

/**
 * Base class for implementing connections.
 * Per default, all operations (read, write, subscribe) are unsupported.
 * Concrete implementations should override the methods indicating connection capabilities
 * and for obtaining respective request builders.
 */
public abstract class AbstractPlcConnection implements PlcConnection, PlcConnectionMetadata, PlcReader, PlcWriter , PlcSubscriber {

    private boolean canRead = false;
    private boolean canWrite = false;
    private boolean canSubscribe = false;
    private PlcFieldHandler fieldHandler;
    private PlcValueHandler valueHandler;
    private Plc4xProtocolBase<?> protocol;
    private BaseOptimizer optimizer;

    /**
     * @deprecated only for compatibility reasons.
     */
    @Deprecated
    protected AbstractPlcConnection() {
    }

    protected AbstractPlcConnection(boolean canRead, boolean canWrite, boolean canSubscribe, PlcFieldHandler fieldHandler, PlcValueHandler valueHandler,
                                 BaseOptimizer optimizer) {
        this.canRead = canRead;
        this.canWrite = canWrite;
        this.canSubscribe = canSubscribe;
        this.fieldHandler = fieldHandler;
        this.valueHandler = valueHandler;
        this.optimizer = optimizer;
    }

    public void setProtocol(Plc4xProtocolBase<?> protocol) {
        this.protocol = protocol;
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return this;
    }

    @Override
    public CompletableFuture<Void> ping() {
        CompletableFuture<Void> future = new CompletableFuture<>();
        future.completeExceptionally(new PlcUnsupportedOperationException("The connection does not support pinging"));
        return future;
    }

    @Override
    public boolean canRead() {
        return canRead;
    }

    @Override
    public boolean canWrite() {
        return canWrite;
    }

    @Override
    public boolean canSubscribe() {
        return canSubscribe;
    }

    public PlcFieldHandler getPlcFieldHandler() {
        return this.fieldHandler;
    }

    public PlcValueHandler getPlcValueHandler() {
        return this.valueHandler;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        if (!canRead()) {
            throw new PlcUnsupportedOperationException("The connection does not support reading");
        }
        return new DefaultPlcReadRequest.Builder(this, getPlcFieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        if (!canWrite()) {
            throw new PlcUnsupportedOperationException("The connection does not support writing");
        }
        return new DefaultPlcWriteRequest.Builder(this, getPlcFieldHandler(), getPlcValueHandler());
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        if (!canSubscribe()) {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }
        return new DefaultPlcSubscriptionRequest.Builder(this, getPlcFieldHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        if (!canSubscribe) {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        if(optimizer != null) {
            return optimizer.optimizedRead(readRequest, protocol);
        }
        return protocol.read(readRequest);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        if(optimizer != null) {
            return optimizer.optimizedWrite(writeRequest, protocol);
        }
        return protocol.write(writeRequest);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        if(optimizer != null) {
            return optimizer.optimizedSubscribe(subscriptionRequest, protocol);
        }
        return protocol.subscribe(subscriptionRequest);
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        if(optimizer != null) {
            return optimizer.optmizedUnsubscribe(unsubscriptionRequest, protocol);
        }
        return protocol.unsubscribe(unsubscriptionRequest);
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        throw new NotImplementedException("");
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        throw new NotImplementedException("");
    }

}
