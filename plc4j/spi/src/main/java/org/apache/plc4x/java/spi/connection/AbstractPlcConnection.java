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
package org.apache.plc4x.java.spi.connection;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

/**
 * Base class for implementing connections.
 * Per default, all operations (read, write, subscribe) are unsupported.
 * Concrete implementations should override the methods indicating connection capabilities
 * and for obtaining respective request builders.
 */
public abstract class AbstractPlcConnection implements PlcConnection, PlcConnectionMetadata, PlcPinger, PlcReader, PlcWriter, PlcSubscriber, PlcBrowser {

    protected static final Logger logger = LoggerFactory.getLogger(AbstractPlcConnection.class);

    private boolean canPing = false;
    private boolean canRead = false;
    private boolean canWrite = false;
    private boolean canSubscribe = false;
    private boolean canBrowse = false;
    private PlcTagHandler tagHandler;
    private PlcValueHandler valueHandler;
    private Plc4xProtocolBase<? extends Message> protocol;
    private BaseOptimizer optimizer;
    private PlcAuthentication authentication;

    /**
     * @deprecated only for compatibility reasons.
     */
    @Deprecated
    protected AbstractPlcConnection() {
    }

    protected AbstractPlcConnection(boolean canPing, boolean canRead, boolean canWrite,
                                    boolean canSubscribe, boolean canBrowse,
                                    PlcTagHandler tagHandler, PlcValueHandler valueHandler,
                                    BaseOptimizer optimizer, PlcAuthentication authentication) {
        this.canPing = canPing;
        this.canRead = canRead;
        this.canWrite = canWrite;
        this.canSubscribe = canSubscribe;
        this.canBrowse = canBrowse;
        this.tagHandler = tagHandler;
        this.valueHandler = valueHandler;
        this.optimizer = optimizer;
        this.authentication = authentication;
    }

    public void setProtocol(Plc4xProtocolBase<? extends Message> protocol) {
        this.protocol = protocol;
    }

    public Plc4xProtocolBase<? extends Message> getProtocol() {
        return protocol;
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return this;
    }

    @Override
    public CompletableFuture<? extends PlcPingResponse> ping() {
        CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();
        future.completeExceptionally(new PlcUnsupportedOperationException("The connection does not support pinging"));
        return future;
    }

    @Override
    public boolean isReadSupported() {
        return canRead;
    }

    @Override
    public boolean isWriteSupported() {
        return canWrite;
    }

    @Override
    public boolean isSubscribeSupported() {
        return canSubscribe;
    }

    @Override
    public boolean isBrowseSupported() {
        return canBrowse;
    }

    public PlcTagHandler getPlcTagHandler() {
        return this.tagHandler;
    }

    public PlcValueHandler getPlcValueHandler() {
        return this.valueHandler;
    }

    protected PlcAuthentication getAuthentication() {
        return authentication;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        if (!isReadSupported()) {
            throw new PlcUnsupportedOperationException("The connection does not support reading");
        }
        return new DefaultPlcReadRequest.Builder(this, getPlcTagHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        if (!isWriteSupported()) {
            throw new PlcUnsupportedOperationException("The connection does not support writing");
        }
        return new DefaultPlcWriteRequest.Builder(this, getPlcTagHandler(), getPlcValueHandler());
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        if (!isSubscribeSupported()) {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }
        return new DefaultPlcSubscriptionRequest.Builder(this, getPlcTagHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        if (!canSubscribe) {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        if (!canBrowse) {
            throw new PlcUnsupportedOperationException("The connection does not support browsing");
        }
        return new DefaultPlcBrowseRequest.Builder(this, getPlcTagHandler());
    }

    @Override
    public CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        if (!canPing) {
            throw new PlcUnsupportedOperationException("The connection does not support pinging");
        }
        return protocol.ping(pingRequest);
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
            return optimizer.optimizedUnsubscribe(unsubscriptionRequest, protocol);
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

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        return protocol.browse(browseRequest);
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        return protocol.browseWithInterceptor(browseRequest, interceptor);
    }

    @Override
    public Optional<PlcTag> parseTagAddress(String tagAddress) {
        PlcTag plcTag;
        try {
            plcTag = tagHandler.parseTag(tagAddress);
        } catch (Exception e) {
            logger.error("Error parsing tag address {}", tagAddress);
            return Optional.empty();
        }
        return Optional.ofNullable(plcTag);
    }
}
