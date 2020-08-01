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
package org.apache.plc4x.java.amsads.attic.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.amsads.types.AdsDataType;
import org.apache.plc4x.java.amsads.field.AdsField;
import org.apache.plc4x.java.amsads.field.DirectAdsField;
import org.apache.plc4x.java.amsads.field.SymbolicAdsField;
import org.apache.plc4x.java.amsads.protocol.exception.AdsException;
import org.apache.plc4x.java.amsads.readwrite.*;
import org.apache.plc4x.java.amsads.readwrite.types.CommandId;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;
import java.util.stream.Collectors;

import static org.apache.plc4x.java.amsads.protocol.util.LittleEndianDecoder.decodeData;
import static org.apache.plc4x.java.amsads.protocol.util.LittleEndianEncoder.encodeData;

@Deprecated
public class Plc4x2AdsProtocol extends MessageToMessageCodec<AmsPacket, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4x2AdsProtocol.class);

    private static final AtomicLong correlationBuilder = new AtomicLong(1);

    private final ConcurrentMap<Long, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> requests;

    private final ConcurrentMap<SymbolicAdsField, DirectAdsField> fieldMapping;

    private List<Consumer<AdsDeviceNotificationRequest>> deviceNotificationListeners;

    private final AmsNetId targetAmsNetId;
    private final int targetAmsPort;
    private final AmsNetId sourceAmsNetId;
    private final int sourceAmsPort;

    public Plc4x2AdsProtocol(AmsNetId targetAmsNetId, int targetAmsPort, AmsNetId sourceAmsNetId, int sourceAmsPort, ConcurrentMap<SymbolicAdsField, DirectAdsField> fieldMapping) {
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
        this.requests = new ConcurrentHashMap<>();
        this.fieldMapping = fieldMapping;
        this.deviceNotificationListeners = new LinkedList<>();
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws Exception {
        LOGGER.trace("(<--OUT): {}, {}, {}", ctx, msg, out);
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        } else if (request instanceof PlcProprietaryRequest) {
            encodeProprietaryRequest(msg, out);
        } else {
            throw new PlcProtocolException("Unknown type " + request.getClass());
        }
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        LOGGER.trace("(-->ERR): {}", ctx, cause);
        if (cause instanceof AdsException) {
            Long invokeId = ((AdsException) cause).getInvokeId();
            if (invokeId != null) {
                PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> remove = requests.remove(invokeId);
                if (remove != null) {
                    remove.getResponseFuture().completeExceptionally(new PlcIoException(cause));
                } else {
                    LOGGER.warn("Unrelated exception received {}", invokeId, cause);
                }
            } else {
                super.exceptionCaught(ctx, cause);
            }
        } else if ((cause instanceof IOException) && (cause.getMessage().contains("Connection reset by peer") ||
            cause.getMessage().contains("Operation timed out"))) {
            String reason = cause.getMessage().contains("Connection reset by peer") ?
                "Connection terminated unexpectedly" : "Remote host not responding";
            if (!requests.isEmpty()) {
                // If the connection is hung up, all still pending requests can be closed.
                for (PlcRequestContainer requestContainer : requests.values()) {
                    requestContainer.getResponseFuture().completeExceptionally(new PlcIoException(reason));
                }
                // Clear the list
                requests.clear();
            }
        } else {
            super.exceptionCaught(ctx, cause);
        }
    }

    private void encodeWriteRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws PlcException {
        InternalPlcWriteRequest writeRequest = (InternalPlcWriteRequest) msg.getRequest();
        if (writeRequest.getFields().size() != 1) {
            throw new PlcProtocolException("Only one item supported");
        }
        PlcField field = writeRequest.getFields().get(0);
        if (field instanceof SymbolicAdsField) {
            DirectAdsField mappedField = fieldMapping.get(field);
            LOGGER.debug("Replacing {} with {}", field, mappedField);
            field = mappedField;
        }
        if (!(field instanceof DirectAdsField)) {
            throw new PlcProtocolException("PlcField not of type DirectAdsField: " + field.getClass());
        }
        DirectAdsField directAdsField = (DirectAdsField) field;
        long invokeId = correlationBuilder.incrementAndGet();
        long indexGroup = directAdsField.getIndexGroup();
        long indexOffset = directAdsField.getIndexOffset();

        PlcValue plcValue = writeRequest.getPlcValues().get(0);
        Object[] plcValues;
        if(plcValue instanceof PlcList) {
            PlcList plcList = (PlcList) plcValue;
            plcValues = plcList.getList().toArray(new Object[0]);
        } else {
            plcValues = new Object[1];
            plcValues[0] = plcValue.getObject();
        }

        byte[] bytes = encodeData(directAdsField.getAdsDataType(), plcValues);
        int bytesToBeWritten = bytes.length;
        int maxTheoreticalSize = directAdsField.getAdsDataType().getTargetByteSize() * directAdsField.getNumberOfElements();
        if (bytesToBeWritten > maxTheoreticalSize) {
            LOGGER.debug("Requested AdsDatatype {} is exceeded by number of bytes {}. Limit {}.", directAdsField.getAdsDataType(), bytesToBeWritten, maxTheoreticalSize);
            throw new PlcProtocolPayloadTooBigException("ADS", maxTheoreticalSize, bytesToBeWritten, plcValues);
        }
        AdsWriteRequest data = new AdsWriteRequest(indexGroup, indexOffset, bytes);
        AmsPacket amsPacket = new AmsPacket(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, CommandId.ADS_WRITE, new State(false, false, false, false, false, false, true, false, false), 0, invokeId, data);
        LOGGER.debug("encoded write request {}", amsPacket);
        out.add(amsPacket);
        requests.put(invokeId, msg);
    }

    private void encodeReadRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws PlcException {
        PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();

        if (readRequest.getFields().size() != 1) {
            throw new PlcProtocolException("Only one item supported");
        }
        PlcField field = readRequest.getFields().get(0);
        if (field instanceof SymbolicAdsField) {
            DirectAdsField mappedField = fieldMapping.get(field);
            if (mappedField == null) {
                throw new PlcProtocolException("No field mapping for " + field);
            }
            LOGGER.debug("Replacing {} with {}", field, mappedField);
            field = mappedField;
        }
        if (!(field instanceof DirectAdsField)) {
            throw new PlcProtocolException("PlcField not of type DirectAdsField: " + field.getClass());
        }
        DirectAdsField directAdsField = (DirectAdsField) field;
        long invokeId = correlationBuilder.incrementAndGet();
        long indexGroup = directAdsField.getIndexGroup();
        long indexOffset = directAdsField.getIndexOffset();
        AdsDataType adsDataType = directAdsField.getAdsDataType();
        int numberOfElements = directAdsField.getNumberOfElements();
        int readLength = adsDataType.getTargetByteSize() * numberOfElements;
        AdsReadWriteRequest data = new AdsReadWriteRequest(indexGroup, indexOffset, readLength, new AdsReadRequest[0], new byte[0]);
        AmsPacket amsPacket = new AmsPacket(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, CommandId.ADS_READ, new State(false, false, false, false, false, false, true, false, false), 0, invokeId, data);
        LOGGER.debug("encoded read request {}", amsPacket);
        out.add(amsPacket);
        requests.put(invokeId, msg);
    }

    private void encodeProprietaryRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws PlcProtocolException {
        PlcProprietaryRequest plcProprietaryRequest = (PlcProprietaryRequest) msg.getRequest();
        if (!(plcProprietaryRequest.getProprietaryRequest() instanceof AmsPacket)) {
            throw new PlcProtocolException("Unsupported proprietary type for this driver " + plcProprietaryRequest.getProprietaryRequest().getClass());
        }
        AmsPacket amsPacket = (AmsPacket) plcProprietaryRequest.getProprietaryRequest();
        LOGGER.debug("encoded proprietary request {}", amsPacket);
        out.add(amsPacket);
        requests.put(amsPacket.getInvokeId(), msg);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, AmsPacket amsPacket, List<Object> out) throws Exception {
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, amsPacket, out);
        AdsData data = amsPacket.getData();
        if (data instanceof AdsDeviceNotificationRequest) {
            LOGGER.debug("Received notification {}", amsPacket);
            handleAdsDeviceNotificationRequest((AdsDeviceNotificationRequest) data);
            return;
        }
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> plcRequestContainer = requests.remove(amsPacket.getInvokeId());
        if (plcRequestContainer == null) {
            LOGGER.info("Unmapped packet received {}", amsPacket);
            return;
        }
        PlcRequest request = plcRequestContainer.getRequest();
        final InternalPlcResponse response;

        // Handle the response to a read request.
        if (request instanceof PlcReadRequest) {
            if (data instanceof AdsReadResponse) {
                response = decodeReadResponse((AdsReadResponse) data, plcRequestContainer);
            } else {
                throw new PlcProtocolException("Wrong type correlated " + amsPacket);
            }
        } else if (request instanceof PlcWriteRequest) {
            if (data instanceof AdsWriteResponse) {
                response = decodeWriteResponse((AdsWriteResponse) data, plcRequestContainer);
            } else {
                throw new PlcProtocolException("Wrong type correlated " + amsPacket);
            }
        } else if (request instanceof PlcProprietaryRequest) {
            response = decodeProprietaryResponse(amsPacket, plcRequestContainer);
        } else {
            response = null;
        }
        LOGGER.debug("Plc4x response {}", response);

        // Confirm the response being handled.
        if (response != null) {
            plcRequestContainer.getResponseFuture().complete(response);
        }
    }

    private void handleAdsDeviceNotificationRequest(AdsDeviceNotificationRequest adsDeviceNotificationRequest) {
        for (Consumer<AdsDeviceNotificationRequest> deviceNotificationListener : deviceNotificationListeners) {
            try {
                deviceNotificationListener.accept(adsDeviceNotificationRequest);
            } catch (RuntimeException e) {
                LOGGER.error("Exception received from {} while handling {}", deviceNotificationListener, adsDeviceNotificationRequest, e);
            }
        }
    }

    public boolean addConsumer(Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer) {
        return deviceNotificationListeners.add(adsDeviceNotificationRequestConsumer);
    }

    public boolean removeConsumer(Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer) {
        return deviceNotificationListeners.remove(adsDeviceNotificationRequestConsumer);
    }


    @SuppressWarnings("unchecked")
    private InternalPlcResponse decodeWriteResponse(AdsWriteResponse responseMessage, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> requestContainer) {
        InternalPlcWriteRequest plcWriteRequest = (InternalPlcWriteRequest) requestContainer.getRequest();
        PlcResponseCode responseCode = decodeResponseCode(responseMessage.getResult());

        // TODO: does every item has the same ads response or is this whole aggregation broken?
        Map<String, PlcResponseCode> responseItems = plcWriteRequest.getFieldNames()
            .stream()
            .collect(Collectors.toMap(
                fieldName -> fieldName,
                ignore -> responseCode
            ));
        return new DefaultPlcWriteResponse(plcWriteRequest, responseItems);
    }

    @SuppressWarnings("unchecked")
    private InternalPlcResponse decodeReadResponse(AdsReadResponse responseMessage, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> requestContainer) {
        InternalPlcReadRequest plcReadRequest = (InternalPlcReadRequest) requestContainer.getRequest();

        // TODO: only single requests supported for now
        AdsField field = (AdsField) plcReadRequest.getFields().get(0);

        PlcResponseCode responseCode = decodeResponseCode(responseMessage.getResult());
        byte[] bytes = responseMessage.getData();
        PlcValue plcValue = decodeData(field.getAdsDataType(), bytes);

        // TODO: does every item has the same ads response or is this whole aggregation broken?
        Map<String, ResponseItem<PlcValue>> responseItems = plcReadRequest.getFieldNames()
            .stream()
            .collect(Collectors.toMap(
                fieldName -> fieldName,
                ignore -> new ResponseItem(responseCode, plcValue)
            ));

        return new DefaultPlcReadResponse(plcReadRequest, responseItems);
    }

    @SuppressWarnings("unchecked")
    private InternalPlcResponse decodeProprietaryResponse(AmsPacket amsPacket, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> plcRequestContainer) {
        return new DefaultPlcProprietaryResponse<>((InternalPlcProprietaryRequest) plcRequestContainer.getRequest(), amsPacket);
    }

    private PlcResponseCode decodeResponseCode(long result) {
        if (result == 0L) {
            return PlcResponseCode.OK;
        } else if (result == 0x7L) {
            return PlcResponseCode.INVALID_ADDRESS;
        } else {
            return PlcResponseCode.INTERNAL_ERROR;
        }
    }

}
