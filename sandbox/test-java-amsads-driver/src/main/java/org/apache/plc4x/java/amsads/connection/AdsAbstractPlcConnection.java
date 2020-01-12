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
package org.apache.plc4x.java.amsads.connection;

import io.netty.channel.ChannelFuture;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.plc4x.java.amsads.field.AdsFieldHandler;
import org.apache.plc4x.java.amsads.field.DirectAdsField;
import org.apache.plc4x.java.amsads.field.SymbolicAdsField;
import org.apache.plc4x.java.amsads.readwrite.*;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.DefaultNettyPlcConnection;
import org.apache.plc4x.java.spi.messages.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.concurrent.*;

@Deprecated
public abstract class AdsAbstractPlcConnection extends DefaultNettyPlcConnection implements PlcReader, PlcWriter, PlcProprietarySender {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsAbstractPlcConnection.class);

    protected static final Configuration CONF = new SystemConfiguration();
    protected static final long SYMBOL_RESOLVE_TIMEOUT = CONF.getLong("plc4x.adsconnection.symbol.resolve,timeout", 3000);

    protected final AmsNetId targetAmsNetId;

    protected final int targetAmsPort;

    protected final AmsNetId sourceAmsNetId;

    protected final int sourceAmsPort;

    protected final ConcurrentMap<SymbolicAdsField, DirectAdsField> fieldMapping;

    protected AdsAbstractPlcConnection(ChannelFactory channelFactory, AmsNetId targetAmsNetId, int targetAmsPort) {
        this(channelFactory, targetAmsNetId, targetAmsPort, generateAmsNetId(), generateAMSPort());
    }

    protected AdsAbstractPlcConnection(ChannelFactory channelFactory, AmsNetId targetAmsNetId, int targetAmsPort, AmsNetId sourceAmsNetId, int sourceAmsPort) {
        super(null, channelFactory, false, null);
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
        this.fieldMapping = new ConcurrentHashMap<>();
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public boolean canWrite() {
        return true;
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
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new AdsFieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new AdsFieldHandler());
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
            AdsReadWriteRequest adsReadWriteRequest = new AdsReadWriteRequest(
                0xF003L,
                0L,
                4L,
                symbolicAdsFieldInternal.getSymbolicField().getBytes().length,
                symbolicAdsFieldInternal.getSymbolicField().getBytes()
            );

            // TODO: This is blocking, should be changed to be async.
            CompletableFuture<InternalPlcProprietaryResponse<AdsReadWriteResponse>> getHandelFuture = new CompletableFuture<>();
            channel.writeAndFlush(new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsReadWriteRequest), getHandelFuture));
            InternalPlcProprietaryResponse<AdsReadWriteResponse> getHandleResponse = getFromFuture(getHandelFuture, SYMBOL_RESOLVE_TIMEOUT);
            AdsReadWriteResponse response = getHandleResponse.getResponse();

            if (response.getResult() != 0L) {
                throw new PlcRuntimeException("Non error code received " + response.getResult());
            }

            ByteBuffer buffer = ByteBuffer.allocate(Long.BYTES);
            buffer.order(ByteOrder.LITTLE_ENDIAN);
            buffer.put(response.getData());
            Long symbolHandle = buffer.getLong();
            return DirectAdsField.of(0xF005, symbolHandle, symbolicAdsFieldInternal.getAdsDataType(), symbolicAdsFieldInternal.getNumberOfElements());
        });
    }

    protected static AmsNetId generateAmsNetId() {
        return new AmsNetId((short) 0, (short) 0, (short) 0, (short) 0, (short) 0, (short) 0);
    }

    protected static int generateAMSPort() {
        return 0;
    }

    @Override
    public void close() throws PlcConnectionException {
        fieldMapping.values().stream()
            .parallel()
            .map(adsField -> {
                ByteBuffer buffer = ByteBuffer.allocate(Long.BYTES);
                buffer.order(ByteOrder.LITTLE_ENDIAN);
                buffer.putLong(adsField.getIndexGroup());
                byte[] bytes = buffer.array();
                return new AdsWriteRequest(
                    0xF006L,
                    0L,
                    bytes.length,
                    bytes
                );
            })
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
