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
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.model.AdsAddress;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.ChannelFactory;

import java.util.concurrent.CompletableFuture;

public abstract class AdsAbstractPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    protected final AmsNetId targetAmsNetId;

    protected final AmsPort targetAmsPort;

    protected final AmsNetId sourceAmsNetId;

    protected final AmsPort sourceAmsPort;

    protected AdsAbstractPlcConnection(ChannelFactory channelFactory, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        this(channelFactory, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    protected AdsAbstractPlcConnection(ChannelFactory channelFactory, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        super(channelFactory);
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
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
    public Address parseAddress(String addressString) throws PlcException {
        return AdsAddress.of(addressString);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> readFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>(readRequest, readFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                readFuture.completeExceptionally(future.cause());
            }
        });
        return readFuture;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> writeFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>(writeRequest, writeFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                writeFuture.completeExceptionally(future.cause());
            }
        });
        return writeFuture;
    }

    protected static AmsNetId generateAMSNetId() {
        return AmsNetId.of("0.0.0.0.0.0");
    }

    protected static AmsPort generateAMSPort() {
        return AmsPort.of(0);
    }

}
