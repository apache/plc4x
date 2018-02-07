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
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.ADSAddress;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.messages.items.WriteResponseItem;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicLong;

public class Plc4XADSProtocol extends MessageToMessageCodec<AMSTCPPaket, PlcRequestContainer> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XADSProtocol.class);

    private static final AtomicLong correlationBuilder = new AtomicLong(1);

    private ConcurrentMap<Long, PlcRequestContainer<PlcRequest, PlcResponse>> requests;

    private final AMSNetId targetAmsNetId;
    private final AMSPort targetAmsPort;
    private final AMSNetId sourceAmsNetId;
    private final AMSPort sourceAmsPort;

    public Plc4XADSProtocol(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
        this.requests = new ConcurrentHashMap<>();
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
        // TODO: implement serialization
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

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, AMSTCPPaket amstcpPaket, List<Object> out) throws Exception {
        PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer = requests.remove(amstcpPaket.getAmsHeader().getInvokeId().getAsLong());
        if (plcRequestContainer == null) {
            LOGGER.info("Unmapped packet received {}", amstcpPaket);
            return;
        }
        PlcRequest request = plcRequestContainer.getRequest();
        PlcResponse response = null;

        // Handle the response to a read request.
        if (request instanceof PlcReadRequest) {
            if (amstcpPaket instanceof ADSReadResponse) {
                response = decodeReadResponse((ADSReadResponse) amstcpPaket, plcRequestContainer);
            } else {
                throw new PlcProtocolException("Wrong type correlated " + amstcpPaket);
            }
        } else if (request instanceof PlcWriteRequest) {
            if (amstcpPaket instanceof ADSWriteResponse) {
                response = decodeWriteResponse((ADSWriteResponse) amstcpPaket, plcRequestContainer);
            } else {
                throw new PlcProtocolException("Wrong type correlated " + amstcpPaket);
            }
        }

        // Confirm the response being handled.
        if (response != null) {
            plcRequestContainer.getResponseFuture().complete(response);
        }
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeWriteResponse(ADSWriteResponse responseMessage, PlcRequestContainer<PlcRequest, PlcResponse> requestContainer) throws PlcProtocolException {
        PlcWriteRequest plcWriteRequest = (PlcWriteRequest) requestContainer.getRequest();
        WriteRequestItem requestItem = plcWriteRequest.getRequestItems().get(0);

        ResponseCode responseCode = decodeResponseCode(responseMessage.getResult());

        if (plcWriteRequest instanceof TypeSafePlcWriteRequest) {
            return new TypeSafePlcWriteResponse((TypeSafePlcWriteRequest) plcWriteRequest, Collections.singletonList(new WriteResponseItem<>(requestItem, responseCode)));
        } else {
            return new PlcWriteResponse(plcWriteRequest, Collections.singletonList(new WriteResponseItem<>(requestItem, responseCode)));
        }
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeReadResponse(ADSReadResponse responseMessage, PlcRequestContainer<PlcRequest, PlcResponse> requestContainer) throws PlcProtocolException {
        PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();
        ReadRequestItem requestItem = plcReadRequest.getRequestItems().get(0);

        ResponseCode responseCode = decodeResponseCode(responseMessage.getResult());
        byte[] bytes = responseMessage.getData().getBytes();

        // TODO: implement deserialization
        if (plcReadRequest instanceof TypeSafePlcReadRequest) {
            return new TypeSafePlcReadResponse((TypeSafePlcReadRequest) plcReadRequest, Collections.singletonList(new ReadResponseItem<>(requestItem, responseCode, Collections.singletonList(bytes))));
        } else {
            return new PlcReadResponse(plcReadRequest, Collections.singletonList(new ReadResponseItem<>(requestItem, responseCode, Collections.singletonList(bytes))));
        }
    }

    private ResponseCode decodeResponseCode(Result result) {
        // TODO: complete mapping
        return result.getAsLong() == 0 ? ResponseCode.OK : ResponseCode.INTERNAL_ERROR;
    }

}
