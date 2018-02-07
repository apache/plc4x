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
package org.apache.plc4x.java.ads.netty;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.ads.api.commands.ADSReadRequest;
import org.apache.plc4x.java.ads.api.commands.ADSReadResponse;
import org.apache.plc4x.java.ads.api.commands.ADSWriteRequest;
import org.apache.plc4x.java.ads.api.commands.ADSWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.IndexGroup;
import org.apache.plc4x.java.ads.api.commands.types.IndexOffset;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.ADSAddress;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcRequestContainer;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.model.Address;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;

public class Plc4XADSProtocol extends MessageToMessageCodec<AMSTCPPaket, PlcRequestContainer> {

    private static final AtomicLong correlationBuilder = new AtomicLong(1);

    private Map<Short, PlcRequestContainer> requests;

    private final AMSNetId targetAmsNetId;
    private final AMSPort targetAmsPort;
    private final AMSNetId sourceAmsNetId;
    private final AMSPort sourceAmsPort;

    public Plc4XADSProtocol(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
        this.requests = new HashMap<>();
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        }
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, AMSTCPPaket amstcpPaket, List<Object> out) throws Exception {
        if (amstcpPaket instanceof ADSReadResponse) {
            // TODO: implement me
        } else if (amstcpPaket instanceof ADSWriteResponse) {
            // TODO: implement me
        }
    }

    private void encodeWriteRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        PlcWriteRequest writeRequest = (PlcWriteRequest) msg.getRequest();
        if (writeRequest.getRequestItems().size() != 1) {
            throw new PlcProtocolException("Only one item supported");
        }
        WriteRequestItem<?> writeRequestItem = writeRequest.getRequestItems().get(0);
        Address address = writeRequestItem.getAddress();
        if (!(address instanceof ADSAddress)) {
            throw new PlcProtocolException("Address not of type ADSAddress: " + address.getClass());
        }
        ADSAddress adsAddress = (ADSAddress) address;
        Invoke invokeId = Invoke.of(correlationBuilder.incrementAndGet());
        IndexGroup indexGroup = IndexGroup.of(adsAddress.getIndexGroup());
        IndexOffset indexOffset = IndexOffset.of(adsAddress.getIndexOffset());
        // TODO: how to get length and data. Serialization of plc is the problem here
        Length length = Length.of(1);
        Data data = Data.of(new byte[]{0x42});
        AMSTCPPaket amstcpPaket = new ADSWriteRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, length, data);
        out.add(amstcpPaket);
    }

    private void encodeReadRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();

        if (readRequest.getRequestItems().size() != 1) {
            throw new PlcProtocolException("Only one item supported");
        }
        ReadRequestItem<?> readRequestItem = readRequest.getRequestItems().get(0);
        Address address = readRequestItem.getAddress();
        if (!(address instanceof ADSAddress)) {
            throw new PlcProtocolException("Address not of type ADSAddress: " + address.getClass());
        }
        ADSAddress adsAddress = (ADSAddress) address;
        Invoke invokeId = Invoke.of(correlationBuilder.incrementAndGet());
        IndexGroup indexGroup = IndexGroup.of(adsAddress.getIndexGroup());
        IndexOffset indexOffset = IndexOffset.of(adsAddress.getIndexOffset());
        Length length = Length.of(readRequestItem.getSize());
        AMSTCPPaket amstcpPaket = new ADSReadRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, length);
        out.add(amstcpPaket);
    }

}
