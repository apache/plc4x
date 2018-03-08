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

import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.model.ADSAddress;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.ChannelFactory;

import java.util.concurrent.CompletableFuture;

public abstract class ADSAbstractPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    protected final AMSNetId targetAmsNetId;

    protected final AMSPort targetAmsPort;

    protected final AMSNetId sourceAmsNetId;

    protected final AMSPort sourceAmsPort;

    protected ADSAbstractPlcConnection(ChannelFactory channelFactory, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(channelFactory, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    protected ADSAbstractPlcConnection(ChannelFactory channelFactory, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        super(channelFactory);
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
    }

    public AMSNetId getTargetAmsNetId() {
        return targetAmsNetId;
    }

    public AMSPort getTargetAmsPort() {
        return targetAmsPort;
    }

    public AMSNetId getSourceAmsNetId() {
        return sourceAmsNetId;
    }

    public AMSPort getSourceAmsPort() {
        return sourceAmsPort;
    }


    @Override
    public Address parseAddress(String addressString) throws PlcException {
        return ADSAddress.of(addressString);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> readFuture = new CompletableFuture<>();
        channel.writeAndFlush(new PlcRequestContainer<>(readRequest, readFuture));
        return readFuture;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> writeFuture = new CompletableFuture<>();
        channel.writeAndFlush(new PlcRequestContainer<>(writeRequest, writeFuture));
        return writeFuture;
    }

    protected static AMSNetId generateAMSNetId() {
        return AMSNetId.of("0.0.0.0.0.0");
    }

    protected static AMSPort generateAMSPort() {
        return AMSPort.of(0);
    }

}
