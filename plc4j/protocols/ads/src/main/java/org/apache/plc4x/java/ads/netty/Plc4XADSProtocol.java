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
import org.apache.plc4x.java.ads.api.generic.AMSTCPPacket;
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

import static org.apache.plc4x.java.ads.netty.util.LittleEndianDecoder.decodeData;
import static org.apache.plc4x.java.ads.netty.util.LittleEndianEncoder.encodeData;

public class Plc4XADSProtocol extends MessageToMessageCodec<AMSTCPPacket, PlcRequestContainer<PlcRequest, PlcResponse>> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XADSProtocol.class);

    private static final AtomicLong correlationBuilder = new AtomicLong(1);

    private final ConcurrentMap<Long, PlcRequestContainer<PlcRequest, PlcResponse>> requests;

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
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws Exception {
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        }
    }


    private void encodeWriteRequest(PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws PlcException {
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
        byte[] bytes = encodeData(writeRequestItem.getDatatype(), writeRequestItem.getValues().toArray());
        Data data = Data.of(bytes);
        AMSTCPPacket amstcpPacket = ADSWriteRequest.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, data);
        out.add(amstcpPacket);
        requests.put(invokeId.getAsLong(), msg);
    }

    private void encodeReadRequest(PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws PlcException {
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
        AMSTCPPacket amstcpPacket = ADSReadRequest.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, length);
        out.add(amstcpPacket);
        requests.put(invokeId.getAsLong(), msg);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, AMSTCPPacket amstcpPacket, List<Object> out) throws Exception {
        PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer = requests.remove(amstcpPacket.getAmsHeader().getInvokeId().getAsLong());
        if (plcRequestContainer == null) {
            LOGGER.info("Unmapped packet received {}", amstcpPacket);
            return;
        }
        PlcRequest request = plcRequestContainer.getRequest();
        PlcResponse response = null;

        // Handle the response to a read request.
        if (request instanceof PlcReadRequest) {
            if (amstcpPacket instanceof ADSReadResponse) {
                response = decodeReadResponse((ADSReadResponse) amstcpPacket, plcRequestContainer);
            } else {
                throw new PlcProtocolException("Wrong type correlated " + amstcpPacket);
            }
        } else if (request instanceof PlcWriteRequest) {
            if (amstcpPacket instanceof ADSWriteResponse) {
                response = decodeWriteResponse((ADSWriteResponse) amstcpPacket, plcRequestContainer);
            } else {
                throw new PlcProtocolException("Wrong type correlated " + amstcpPacket);
            }
        }

        // Confirm the response being handled.
        if (response != null) {
            plcRequestContainer.getResponseFuture().complete(response);
        }
        out.add(plcRequestContainer);
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeWriteResponse(ADSWriteResponse responseMessage, PlcRequestContainer<PlcRequest, PlcResponse> requestContainer) {
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
        List decoded = decodeData(requestItem.getDatatype(), bytes);

        if (plcReadRequest instanceof TypeSafePlcReadRequest) {
            return new TypeSafePlcReadResponse((TypeSafePlcReadRequest) plcReadRequest, Collections.singletonList(new ReadResponseItem<>(requestItem, responseCode, decoded)));
        } else {
            return new PlcReadResponse(plcReadRequest, Collections.singletonList(new ReadResponseItem<>(requestItem, responseCode, decoded)));
        }
    }

    private ResponseCode decodeResponseCode(Result result) {
        switch (result.toAdsReturnCode()) {
            case ADS_CODE_0:
                return ResponseCode.OK;
            case ADS_CODE_1:
                return ResponseCode.INTERNAL_ERROR;
            case ADS_CODE_2:
            case ADS_CODE_3:
            case ADS_CODE_4:
            case ADS_CODE_5:
                return ResponseCode.INTERNAL_ERROR;
            case ADS_CODE_6:
            case ADS_CODE_7:
                return ResponseCode.INVALID_ADDRESS;
            case ADS_CODE_8:
            case ADS_CODE_9:
            case ADS_CODE_10:
            case ADS_CODE_11:
            case ADS_CODE_12:
            case ADS_CODE_13:
            case ADS_CODE_14:
            case ADS_CODE_15:
            case ADS_CODE_16:
            case ADS_CODE_17:
            case ADS_CODE_18:
            case ADS_CODE_19:
            case ADS_CODE_20:
            case ADS_CODE_21:
            case ADS_CODE_22:
            case ADS_CODE_23:
            case ADS_CODE_24:
            case ADS_CODE_25:
            case ADS_CODE_26:
            case ADS_CODE_27:
            case ADS_CODE_28:
            case ADS_CODE_1280:
            case ADS_CODE_1281:
            case ADS_CODE_1282:
            case ADS_CODE_1283:
            case ADS_CODE_1284:
            case ADS_CODE_1285:
            case ADS_CODE_1286:
            case ADS_CODE_1287:
            case ADS_CODE_1288:
            case ADS_CODE_1289:
            case ADS_CODE_1290:
            case ADS_CODE_1291:
            case ADS_CODE_1292:
            case ADS_CODE_1293:
            case ADS_CODE_1792:
            case ADS_CODE_1793:
            case ADS_CODE_1794:
            case ADS_CODE_1795:
            case ADS_CODE_1796:
            case ADS_CODE_1797:
            case ADS_CODE_1798:
            case ADS_CODE_1799:
            case ADS_CODE_1800:
            case ADS_CODE_1801:
            case ADS_CODE_1802:
            case ADS_CODE_1803:
            case ADS_CODE_1804:
            case ADS_CODE_1805:
            case ADS_CODE_1806:
            case ADS_CODE_1807:
            case ADS_CODE_1808:
            case ADS_CODE_1809:
            case ADS_CODE_1810:
            case ADS_CODE_1811:
            case ADS_CODE_1812:
            case ADS_CODE_1813:
            case ADS_CODE_1814:
            case ADS_CODE_1815:
            case ADS_CODE_1816:
            case ADS_CODE_1817:
            case ADS_CODE_1818:
            case ADS_CODE_1819:
            case ADS_CODE_1820:
            case ADS_CODE_1821:
            case ADS_CODE_1822:
            case ADS_CODE_1823:
            case ADS_CODE_1824:
            case ADS_CODE_1825:
            case ADS_CODE_1826:
            case ADS_CODE_1827:
            case ADS_CODE_1828:
            case ADS_CODE_1836:
            case ADS_CODE_1856:
            case ADS_CODE_1857:
            case ADS_CODE_1858:
            case ADS_CODE_1859:
            case ADS_CODE_1860:
            case ADS_CODE_1861:
            case ADS_CODE_1862:
            case ADS_CODE_1863:
            case ADS_CODE_1864:
            case ADS_CODE_1872:
            case ADS_CODE_1873:
            case ADS_CODE_1874:
            case ADS_CODE_1875:
            case ADS_CODE_1876:
            case ADS_CODE_1877:
            case ADS_CODE_4096:
            case ADS_CODE_4097:
            case ADS_CODE_4098:
            case ADS_CODE_4099:
            case ADS_CODE_4100:
            case ADS_CODE_4101:
            case ADS_CODE_4102:
            case ADS_CODE_4103:
            case ADS_CODE_4104:
            case ADS_CODE_4105:
            case ADS_CODE_4106:
            case ADS_CODE_4107:
            case ADS_CODE_4108:
            case ADS_CODE_4109:
            case ADS_CODE_4110:
            case ADS_CODE_4111:
            case ADS_CODE_4112:
            case ADS_CODE_4119:
            case ADS_CODE_4120:
            case ADS_CODE_4121:
            case ADS_CODE_4122:
            case ADS_CODE_10060:
            case ADS_CODE_10061:
            case ADS_CODE_10065:
                return ResponseCode.INTERNAL_ERROR;
            case UNKNOWN:
                return ResponseCode.INTERNAL_ERROR;
        }
        throw new IllegalStateException(result.toAdsReturnCode() + " not mapped");
    }

}
