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
package org.apache.plc4x.java.ads.connection;

import io.netty.channel.ChannelFuture;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.plc4x.java.ads.api.commands.AdsReadWriteRequest;
import org.apache.plc4x.java.ads.api.commands.AdsReadWriteResponse;
import org.apache.plc4x.java.ads.api.commands.AdsWriteRequest;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.AdsPlcFieldHandler;
import org.apache.plc4x.java.ads.model.DirectAdsField;
import org.apache.plc4x.java.ads.model.SymbolicAdsField;
import org.apache.plc4x.java.base.messages.PlcProprietarySender;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.messages.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;
import java.util.concurrent.*;

public abstract class AdsAbstractPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter, PlcProprietarySender {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsAbstractPlcConnection.class);

    protected static final Configuration CONF = new SystemConfiguration();
    protected static final long SYMBOL_RESOLVE_TIMEOUT = CONF.getLong("plc4x.adsconnection.symbol.resolve,timeout", 3000);

    protected final AmsNetId targetAmsNetId;

    protected final AmsPort targetAmsPort;

    protected final AmsNetId sourceAmsNetId;

    protected final AmsPort sourceAmsPort;

    protected final ConcurrentMap<SymbolicAdsField, DirectAdsField> fieldMapping;

    protected AdsAbstractPlcConnection(ChannelFactory channelFactory, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        this(channelFactory, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    protected AdsAbstractPlcConnection(ChannelFactory channelFactory, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        super(channelFactory);
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
        this.fieldMapping = new ConcurrentHashMap<>();
    }

    public AmsNetId getTargetAmsNetId() {
        return targetAmsNetId;
    }

    public AmsPort getTargetAmsPort() {
        return targetAmsPort;
    }

    public AmsNetId getSourceAmsNetId() {
        return sourceAmsNetId;
    }

    public AmsPort getSourceAmsPort() {
        return sourceAmsPort;
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        mapFields(readRequest);
        CompletableFuture<InternalPlcReadResponse> readFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>((InternalPlcReadRequest) readRequest, readFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                readFuture.completeExceptionally(future.cause());
            }
        });
        return readFuture
            .thenApply(PlcReadResponse.class::cast);
    }

    @Override
    public Optional<PlcReadRequest.Builder> readRequestBuilder() {
        return Optional.of(new DefaultPlcReadRequest.Builder(this, new AdsPlcFieldHandler()));
    }

    @Override
    public Optional<PlcWriteRequest.Builder> writeRequestBuilder() {
        return Optional.of(new DefaultPlcWriteRequest.Builder(this, new AdsPlcFieldHandler()));
    }

    @Override
    public Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        mapFields(writeRequest);
        CompletableFuture<InternalPlcWriteResponse> writeFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>((InternalPlcWriteRequest) writeRequest, writeFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                writeFuture.completeExceptionally(future.cause());
            }
        });
        return writeFuture
            .thenApply(PlcWriteResponse.class::cast);
    }

    @Override
    public <T> CompletableFuture<PlcProprietaryResponse<T>> send(PlcProprietaryRequest proprietaryRequest) {
        CompletableFuture<InternalPlcProprietaryResponse<T>> sendFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>((InternalPlcProprietaryRequest) proprietaryRequest, sendFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                sendFuture.completeExceptionally(future.cause());
            }
        });
        return sendFuture
            .thenApply(PlcProprietaryResponse.class::cast);
    }

    protected void mapFields(PlcFieldRequest request) {
        request.getFields().stream()
            .parallel()
            .filter(SymbolicAdsField.class::isInstance)
            .map(SymbolicAdsField.class::cast)
            .forEach(this::mapFields);
    }

    protected void mapFields(SymbolicAdsField symbolicAdsField) {
        // If the map doesn't contain an entry for the given symbolicAdsField,
        // resolve it and add it to the map.
        fieldMapping.computeIfAbsent(symbolicAdsField, symbolicAdsFieldInternal -> {
            LOGGER.debug("Resolving {}", symbolicAdsFieldInternal);
            AdsReadWriteRequest adsReadWriteRequest = AdsReadWriteRequest.of(
                targetAmsNetId,
                targetAmsPort,
                sourceAmsNetId,
                sourceAmsPort,
                Invoke.NONE,
                IndexGroup.ReservedGroups.ADSIGRP_SYM_HNDBYNAME,
                IndexOffset.NONE,
                ReadLength.of(IndexOffset.NUM_BYTES),
                Data.of(symbolicAdsFieldInternal.getSymbolicField())
            );

            // TODO: This is blocking, should be changed to be async.
            CompletableFuture<InternalPlcProprietaryResponse<AdsReadWriteResponse>> getHandelFuture = new CompletableFuture<>();
            channel.writeAndFlush(new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsReadWriteRequest), getHandelFuture));
            InternalPlcProprietaryResponse<AdsReadWriteResponse> getHandleResponse = getFromFuture(getHandelFuture, SYMBOL_RESOLVE_TIMEOUT);
            AdsReadWriteResponse response = getHandleResponse.getResponse();

            if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                throw new PlcRuntimeException("Non error code received " + response.getResult());
            }

            IndexOffset symbolHandle = IndexOffset.of(response.getData().getBytes());
            return DirectAdsField.of(IndexGroup.ReservedGroups.ADSIGRP_SYM_VALBYHND.getAsLong(), symbolHandle.getAsLong(), symbolicAdsFieldInternal.getAdsDataType(), symbolicAdsFieldInternal.getNumberOfElements());
        });
    }

    protected static AmsNetId generateAMSNetId() {
        return AmsNetId.of("0.0.0.0.0.0");
    }

    protected static AmsPort generateAMSPort() {
        return AmsPort.of(0);
    }

    @Override
    public void close() throws PlcConnectionException {
        fieldMapping.values().stream()
            .parallel()
            .map(adsField -> AdsWriteRequest.of(
                targetAmsNetId,
                targetAmsPort,
                sourceAmsNetId,
                sourceAmsPort,
                Invoke.NONE,
                IndexGroup.ReservedGroups.ADSIGRP_SYM_RELEASEHND,
                IndexOffset.NONE,
                Data.of(IndexGroup.of(adsField.getIndexGroup()).getBytes())
            ))
            .map(adsWriteRequest -> new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsWriteRequest), new CompletableFuture<>()))
            // We don't need a response so we just supply a throw away future.
            .forEach(channel::write);
        channel.flush();
        super.close();
    }

    /**
     * Clears the fieldMapping.
     */
    public void clearMapping() {
        fieldMapping.clear();
    }

    protected <T> T getFromFuture(CompletableFuture<T> future, long timeout) {
        try {
            return future.get(timeout, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            LOGGER.warn("Interrupted!", e);
            Thread.currentThread().interrupt();
            throw new PlcRuntimeException(e);
        } catch (ExecutionException | TimeoutException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public String toString() {
        return "AdsAbstractPlcConnection{" +
            "targetAmsNetId=" + targetAmsNetId +
            ", targetAmsPort=" + targetAmsPort +
            ", sourceAmsNetId=" + sourceAmsNetId +
            ", sourceAmsPort=" + sourceAmsPort +
            "} " + super.toString();
    }
}
