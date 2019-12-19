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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.connection.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.types.COTPProtocolClass;
import org.apache.plc4x.java.s7.readwrite.types.COTPTpduSize;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.types.DeviceGroup;
import org.apache.plc4x.java.s7.readwrite.types.S7ControllerType;
import org.apache.plc4x.java.s7.readwrite.types.SzlModuleTypeClass;
import org.apache.plc4x.java.s7.readwrite.types.SzlSublist;
import org.apache.plc4x.java.s7.readwrite.utils.S7Field;
import org.apache.plc4x.java.s7.readwrite.utils.S7TsapIdEncoder;
import org.apache.plc4x.java.s7data.readwrite.*;
import org.apache.plc4x.java.s7data.readwrite.io.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.items.*;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4xS7Protocol extends Plc4xProtocolBase<TPKTPacket> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xS7Protocol.class);

    private final int callingTsapId;
    private int calledTsapId;
    private COTPTpduSize cotpTpduSize;
    private int pduSize;
    private int maxAmqCaller;
    private int maxAmqCallee;
    private S7ControllerType controllerType;

    private static final AtomicInteger tpduGenerator = new AtomicInteger(10);

    public Plc4xS7Protocol(S7Configuration configuration) {
        this.callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, configuration.rack, configuration.slot);
        this.calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OS, 0, 0);
        this.controllerType = configuration.controllerType == null ? S7ControllerType.ANY : S7ControllerType.valueOf(configuration.controllerType);
        if (this.controllerType == S7ControllerType.LOGO && configuration.pduSize == 1024) {
            configuration.pduSize = 480;
        }
        this.cotpTpduSize = getNearestMatchingTpduSize(configuration.pduSize);
        this.pduSize = cotpTpduSize.getSizeInBytes() - 16;
        this.maxAmqCaller = configuration.maxAmqCaller;
        this.maxAmqCallee = configuration.maxAmqCallee;
    }

    /**
     * Iterate over all values until one is found that the given tpdu size will fit.
     *
     * @param tpduSizeParameter requested tpdu size.
     * @return smallest {@link COTPTpduSize} which will fit a given size of tpdu.
     */
    protected COTPTpduSize getNearestMatchingTpduSize(short tpduSizeParameter) {
        for (COTPTpduSize value : COTPTpduSize.values()) {
            if (value.getSizeInBytes() >= tpduSizeParameter) {
                return value;
            }
        }
        return null;
    }

    @Override
    public void onConnect(ConversationContext<TPKTPacket> context) {
        logger.debug("ISO Transport Protocol Sending Connection Request");
        // Open the session on ISO Transport Protocol first.
        TPKTPacket packet = new TPKTPacket(createCOTPConnectionRequest(calledTsapId, callingTsapId, cotpTpduSize));

        context.sendRequest(packet)
            .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
            .check(p -> p.getPayload() instanceof COTPPacketConnectionResponse)
            .unwrap(p -> (COTPPacketConnectionResponse) p.getPayload())
            .handle(cotpPacketConnectionResponse -> {
                context.sendRequest(createS7ConnectionRequest(cotpPacketConnectionResponse))
                    .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
                    .unwrap(TPKTPacket::getPayload)
                    .only(COTPPacketData.class)
                    .unwrap(COTPPacket::getPayload)
                    .only(S7MessageResponse.class)
                    .unwrap(S7Message::getParameter)
                    .only(S7ParameterSetupCommunication.class)
                    .handle(setupCommunication -> {
                        // Save some data from the response.
                        maxAmqCaller = setupCommunication.getMaxAmqCaller();
                        maxAmqCallee = setupCommunication.getMaxAmqCallee();
                        pduSize = setupCommunication.getPduLength();

                        // Only if the controller type is set to "ANY", then try to identify the PLC type.
                        if (controllerType != S7ControllerType.ANY) {
                            // Send an event that connection setup is complete.
                            context.fireConnected();
                            return;
                        }
                        // Prepare a message to request the remote to identify itself.
                        TPKTPacket tpktPacket = createIdentifyRemoteMessage();
                        context.sendRequest(tpktPacket)
                            .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
                            .check(p -> p.getPayload() instanceof COTPPacketData)
                            .unwrap(p -> ((COTPPacketData) p.getPayload()))
                            .check(p -> p.getPayload() instanceof S7MessageUserData)
                            .unwrap(p -> ((S7MessageUserData) p.getPayload()))
                            .check(p -> p.getPayload() instanceof S7PayloadUserData)
                            .handle(messageUserData -> {
                                S7PayloadUserData payloadUserData = (S7PayloadUserData) messageUserData.getPayload();
                                extractControllerTypeAndFireConnected(context, payloadUserData);
                            });
                    });
            });
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<S7VarRequestParameterItem> requestItems = new ArrayList<>(request.getNumberOfFields());
        for (PlcField field : request.getFields()) {
            requestItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(field)));
        }
        final int tpduId = tpduGenerator.getAndIncrement();
        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
            new S7MessageRequest(tpduId,
                new S7ParameterReadVarRequest(requestItems.toArray(new S7VarRequestParameterItem[0])),
                new S7PayloadReadVarRequest()),
            true, (short) tpduId));

        context.sendRequest(tpktPacket)
            .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> ((COTPPacketData) p.getPayload()))
            .check(p -> p.getPayload() instanceof S7MessageResponse)
            .unwrap(p -> ((S7MessageResponse) p.getPayload()))
            .check(p -> p.getTpduReference() == tpduId)
            .check(p -> p.getParameter() instanceof S7ParameterReadVarResponse)
            .handle(p -> {
                try {
                    future.complete(((PlcReadResponse) decodeReadResponse(p, ((InternalPlcReadRequest) readRequest))));
                } catch (PlcProtocolException e) {
                    e.printStackTrace();
                }
            });
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<S7VarRequestParameterItem> parameterItems = new ArrayList<>(request.getNumberOfFields());
        List<S7VarPayloadDataItem> payloadItems = new ArrayList<>(request.getNumberOfFields());
        for (String fieldName : request.getFieldNames()) {
            final S7Field field = (S7Field) request.getField(fieldName);
            final BaseDefaultFieldItem fieldItem = request.getFieldItem(fieldName);
            parameterItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(field)));
            payloadItems.add(encodeFieldItem(field, fieldItem));
        }
        final int tpduId = tpduGenerator.getAndIncrement();
        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
            new S7MessageRequest(tpduId,
                new S7ParameterWriteVarRequest(parameterItems.toArray(new S7VarRequestParameterItem[0])),
                new S7PayloadWriteVarRequest(payloadItems.toArray(new S7VarPayloadDataItem[0]))),
            true, (short) tpduId));

        context.sendRequest(tpktPacket)
            .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> ((COTPPacketData) p.getPayload()))
            .check(p -> p.getPayload() instanceof S7MessageResponse)
            .unwrap(p -> ((S7MessageResponse) p.getPayload()))
            .check(p -> p.getTpduReference() == tpduId)
            .check(p -> p.getParameter() instanceof S7ParameterWriteVarResponse)
            .handle(p -> {
                try {
                    future.complete(((PlcWriteResponse) decodeWriteResponse(p, ((InternalPlcWriteRequest) writeRequest))));
                } catch (PlcProtocolException e) {
                    e.printStackTrace();
                }
            });
        return future;
    }

    private void extractControllerTypeAndFireConnected(ConversationContext<TPKTPacket> context, S7PayloadUserData payloadUserData) {
        for (S7PayloadUserDataItem item : payloadUserData.getItems()) {
            if (!(item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse)) {
                continue;
            }
            S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponseItem =
                (S7PayloadUserDataItemCpuFunctionReadSzlResponse) item;
            for (SzlDataTreeItem readSzlResponseItemItem : readSzlResponseItem.getItems()) {
                if (readSzlResponseItemItem.getItemIndex() != 0x0001) {
                    continue;
                }
                final String articleNumber = new String(readSzlResponseItemItem.getMlfb());
                controllerType = decodeControllerType(articleNumber);

                // Send an event that connection setup is complete.
                context.fireConnected();
            }
        }
    }

    private TPKTPacket createIdentifyRemoteMessage() {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(1, new S7ParameterUserData(new S7ParameterUserDataItem[]{
            new S7ParameterUserDataItemCPUFunctions((short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00, null, null, null)
        }), new S7PayloadUserData(new S7PayloadUserDataItem[]{
            new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, new SzlId(SzlModuleTypeClass.CPU, (byte) 0x00, SzlSublist.MODULE_IDENTIFICATION), 0x0000)
        }));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (short) 2);
        return new TPKTPacket(cotpPacketData);
    }

    private TPKTPacket createS7ConnectionRequest(COTPPacketConnectionResponse cotpPacketConnectionResponse) {
        for (COTPParameter parameter : cotpPacketConnectionResponse.getParameters()) {
            if (parameter instanceof COTPParameterCalledTsap) {
                COTPParameterCalledTsap cotpParameterCalledTsap = (COTPParameterCalledTsap) parameter;
                calledTsapId = cotpParameterCalledTsap.getTsapId();
            } else if (parameter instanceof COTPParameterTpduSize) {
                COTPParameterTpduSize cotpParameterTpduSize = (COTPParameterTpduSize) parameter;
                cotpTpduSize = cotpParameterTpduSize.getTpduSize();
            } else if (parameter instanceof COTPParameterCallingTsap) {
                // Ignore this ...
            } else {
                logger.warn("Got unknown parameter type '" + parameter.getClass().getName() + "'");
            }
        }

        // Send an S7 login message.
        S7ParameterSetupCommunication s7ParameterSetupCommunication =
            new S7ParameterSetupCommunication(maxAmqCaller, maxAmqCallee, pduSize);
        S7Message s7Message = new S7MessageRequest(0, s7ParameterSetupCommunication,
            new S7PayloadSetupCommunication());
        COTPPacketData cotpPacketData = new COTPPacketData(null, s7Message, true, (short) 1);
        return new TPKTPacket(cotpPacketData);
    }

    private COTPPacketConnectionRequest createCOTPConnectionRequest(int calledTsapId, int callingTsapId, COTPTpduSize cotpTpduSize) {
        return new COTPPacketConnectionRequest(
            new COTPParameter[]{
                new COTPParameterCalledTsap(calledTsapId),
                new COTPParameterCallingTsap(callingTsapId),
                new COTPParameterTpduSize(cotpTpduSize)
            }, null, (short) 0x0000, (short) 0x000F, COTPProtocolClass.CLASS_0);
    }

    private PlcResponse decodeReadResponse(S7MessageResponse responseMessage, InternalPlcReadRequest plcReadRequest) throws PlcProtocolException {
        S7PayloadReadVarResponse payload = (S7PayloadReadVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> values = new HashMap<>();
        S7VarPayloadDataItem[] payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
            S7VarPayloadDataItem payloadItem = payloadItems[index];

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            BaseDefaultFieldItem fieldItem = null;
            ByteBuf data = Unpooled.wrappedBuffer(payloadItem.getData());
            if (responseCode == PlcResponseCode.OK) {
                fieldItem = decodeFieldItem(fieldName, field, data);
            }
            Pair<PlcResponseCode, BaseDefaultFieldItem> result = new ImmutablePair<>(responseCode, fieldItem);
            values.put(fieldName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    private PlcResponse decodeWriteResponse(S7MessageResponse responseMessage, InternalPlcWriteRequest plcWriteRequest) throws PlcProtocolException {
        S7PayloadWriteVarResponse payload = (S7PayloadWriteVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, PlcResponseCode> responses = new HashMap<>();
        S7VarPayloadStatusItem[] payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcWriteRequest.getFieldNames()) {
            S7VarPayloadStatusItem payloadItem = payloadItems[index];

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            responses.put(fieldName, responseCode);
            index++;
        }

        return new DefaultPlcWriteResponse(plcWriteRequest, responses);
    }

    private S7VarPayloadDataItem encodeFieldItem(S7Field field, BaseDefaultFieldItem fieldItem) {
        if(fieldItem.getNumberOfValues() > 0) {
            throw new NotImplementedException("Writing more than one element currently not supported");
        }

        try {
            WriteBuffer writeBuffer;
            DataTransportSize transportSize;
            switch (field.getDataType().getDataProtocolId()) {
                case 1: // 1 Bit
                    writeBuffer = new WriteBuffer(1);
                    DataItemBOOLIO.serialize(writeBuffer, new DataItemBOOL(fieldItem.getBoolean(0)));
                    transportSize = DataTransportSize.BIT;
                    break;

                case 11: // BYTE
                    writeBuffer = new WriteBuffer(1);
                    DataItemBYTEIO.serialize(writeBuffer, new DataItemBYTE(ArrayUtils.toPrimitive((Boolean[]) fieldItem.getValues())));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 12: // WORD
                    writeBuffer = new WriteBuffer(2);
                    DataItemWORDIO.serialize(writeBuffer, new DataItemWORD(ArrayUtils.toPrimitive((Boolean[]) fieldItem.getValues())));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 13: // DWORD
                    writeBuffer = new WriteBuffer(4);
                    DataItemDWORDIO.serialize(writeBuffer, new DataItemDWORD(ArrayUtils.toPrimitive((Boolean[]) fieldItem.getValues())));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 14: // LWORD
                    writeBuffer = new WriteBuffer(8);
                    DataItemLWORDIO.serialize(writeBuffer, new DataItemLWORD(ArrayUtils.toPrimitive((Boolean[]) fieldItem.getValues())));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;

                case 21: // SINT
                    writeBuffer = new WriteBuffer(1);
                    DataItemSINTIO.serialize(writeBuffer, new DataItemSINT(fieldItem.getByte(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 22: // USINT
                    writeBuffer = new WriteBuffer(1);
                    DataItemUSINTIO.serialize(writeBuffer, new DataItemUSINT(fieldItem.getShort(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 23: // INT
                    writeBuffer = new WriteBuffer(2);
                    DataItemINTIO.serialize(writeBuffer, new DataItemINT(fieldItem.getShort(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 24: // UINT
                    writeBuffer = new WriteBuffer(2);
                    DataItemUINTIO.serialize(writeBuffer, new DataItemUINT(fieldItem.getInteger(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 25: // DINT
                    writeBuffer = new WriteBuffer(4);
                    DataItemDINTIO.serialize(writeBuffer, new DataItemDINT(fieldItem.getInteger(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 26: // UDINT
                    writeBuffer = new WriteBuffer(4);
                    DataItemUDINTIO.serialize(writeBuffer, new DataItemUDINT(fieldItem.getLong(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 27: // LINT
                    writeBuffer = new WriteBuffer(8);
                    DataItemLINTIO.serialize(writeBuffer, new DataItemLINT(fieldItem.getLong(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 28: // ULINT
                    writeBuffer = new WriteBuffer(8);
                    DataItemULINTIO.serialize(writeBuffer, new DataItemULINT(fieldItem.getBigInteger(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;

                case 31: // REAL
                    writeBuffer = new WriteBuffer(2);
                    DataItemREALIO.serialize(writeBuffer, new DataItemREAL(fieldItem.getFloat(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 32: // LREAL
                    writeBuffer = new WriteBuffer(4);
                    DataItemLREALIO.serialize(writeBuffer, new DataItemLREAL(fieldItem.getDouble(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;

/*                case 41: // CHAR
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());
                case 42: // WCHAR
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());
                case 43: // STRING
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());
                case 44: // WSTRING
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());*/

                case 51: // Time
                    writeBuffer = new WriteBuffer(4);
                    DataItemTimeIO.serialize(writeBuffer, new DataItemTime(fieldItem.getTime(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 52: // LTime
                    writeBuffer = new WriteBuffer(4);
                    DataItemLTimeIO.serialize(writeBuffer, new DataItemLTime(fieldItem.getTime(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 53: // Date
                    writeBuffer = new WriteBuffer(2);
                    DataItemDateIO.serialize(writeBuffer, new DataItemDate(fieldItem.getDate(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 54: // TimeOfDay
                    writeBuffer = new WriteBuffer(4);
                    DataItemTimeOfDayIO.serialize(writeBuffer, new DataItemTimeOfDay(fieldItem.getTime(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;
                case 55: // DateAndTime
                    writeBuffer = new WriteBuffer(8);
                    DataItemDateAndTimeIO.serialize(writeBuffer, new DataItemDateAndTime(fieldItem.getDateTime(0)));
                    transportSize = DataTransportSize.BYTE_WORD_DWORD;
                    break;

                default:
                    throw new NotImplementedException(
                        "Processing of datatype " + field.getDataType() + " currently not supported");
            }

            byte[] data = writeBuffer.getData();
            return new S7VarPayloadDataItem(DataTransportErrorCode.OK, transportSize, data.length, data);
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return null;
    }

    private BaseDefaultFieldItem decodeFieldItem(String fieldName, S7Field field, ByteBuf data) {
        ReadBuffer readBuffer = new ReadBuffer(data.array());
        try {
            DataItem item = DataItemIO.parse(readBuffer, field.getDataType().getDataProtocolId());
            switch (field.getDataType().getDataProtocolId()) {
                case 1: // 1 Bit
                    return new DefaultBooleanFieldItem(((DataItemBOOL) item).getValue());
                case 11: // BYTE
                    return new DefaultBooleanFieldItem(ArrayUtils.toObject(((DataItemBYTE) item).getValue()));
                case 12: // WORD
                    return new DefaultBooleanFieldItem(ArrayUtils.toObject(((DataItemWORD) item).getValue()));
                case 13: // DWORD
                    return new DefaultBooleanFieldItem(ArrayUtils.toObject(((DataItemDWORD) item).getValue()));
                case 14: // LWORD
                    return new DefaultBooleanFieldItem(ArrayUtils.toObject(((DataItemLWORD) item).getValue()));

                case 21: // SINT
                    return new DefaultByteFieldItem(((DataItemSINT) item).getValue());
                case 22: // USINT
                    return new DefaultShortFieldItem(((DataItemUSINT) item).getValue());
                case 23: // INT
                    return new DefaultShortFieldItem(((DataItemINT) item).getValue());
                case 24: // UINT
                    return new DefaultIntegerFieldItem(((DataItemUINT) item).getValue());
                case 25: // DINT
                    return new DefaultIntegerFieldItem(((DataItemDINT) item).getValue());
                case 26: // UDINT
                    return new DefaultLongFieldItem(((DataItemUDINT) item).getValue());
                case 27: // LINT
                    return new DefaultLongFieldItem(((DataItemLINT) item).getValue());
                case 28: // ULINT
                    return new DefaultBigIntegerFieldItem(((DataItemULINT) item).getValue());

                case 31: // REAL
                    return new DefaultFloatFieldItem(((DataItemREAL) item).getValue());
                case 32: // LREAL
                    return new DefaultDoubleFieldItem(((DataItemLREAL) item).getValue());

/*                case 41: // CHAR
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());
                case 42: // WCHAR
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());
                case 43: // STRING
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());
                case 44: // WSTRING
                    return new DefaultStringFieldItem(((DataItemCHAR) item).getValue());*/

                case 51: // Time
                    return new DefaultLocalTimeFieldItem(((DataItemTime) item).getValue());
                case 52: // LTime
                    return new DefaultLocalTimeFieldItem(((DataItemLTime) item).getValue());
                case 53: // Date
                    return new DefaultLocalDateFieldItem(((DataItemDate) item).getValue());
                case 54: // TimeOfDay
                    return new DefaultLocalTimeFieldItem(((DataItemTimeOfDay) item).getValue());
                case 55: // DateAndTime
                    return new DefaultLocalDateTimeFieldItem(((DataItemDateAndTime) item).getValue());

                default:
                    throw new NotImplementedException(
                        "Processing of datatype " + field.getDataType() + " currently not supported");
            }
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return null;
    }

    /**
     * Helper to convert the return codes returned from the S7 into one of our standard
     * PLC4X return codes
     * @param dataTransportErrorCode S7 return code
     * @return PLC4X return code.
     */
    private PlcResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return PlcResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return PlcResponseCode.OK;
            case NOT_FOUND:
                return PlcResponseCode.NOT_FOUND;
            case INVALID_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED:
                return PlcResponseCode.INVALID_DATATYPE;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    /**
     * Little helper method to parse Siemens article numbers and extract the type of controller.
     *
     * @param articleNumber article number string.
     * @return type of controller.
     */
    private S7ControllerType decodeControllerType(String articleNumber) {
        if (!articleNumber.startsWith("6ES7 ")) {
            return S7ControllerType.ANY;
        }
        String model = articleNumber.substring(articleNumber.indexOf(' ') + 1, articleNumber.indexOf(' ') + 2);
        switch (model) {
            case "2":
                return S7ControllerType.S7_1200;
            case "5":
                return S7ControllerType.S7_1500;
            case "3":
                return S7ControllerType.S7_300;
            case "4":
                return S7ControllerType.S7_400;
            default:
                if (logger.isInfoEnabled()) {
                    logger.info(String.format("Looking up unknown article number %s", articleNumber));
                }
                return S7ControllerType.ANY;
        }
    }

    /**
     * Currently we only support the S7 Any type of addresses. This helper simply converts the S7Field
     * from PLC4X into S7Address objects.
     * @param field S7Field instance we need to convert into an S7Address
     * @return the S7Address
     */
    protected S7Address encodeS7Address(PlcField field) {
        if (!(field instanceof S7Field)) {
            throw new RuntimeException("Unsupported address type " + field.getClass().getName());
        }
        S7Field s7Field = (S7Field) field;
        return new S7AddressAny(s7Field.getDataType(), s7Field.getNumElements(), s7Field.getBlockNumber(),
            s7Field.getMemoryArea(), s7Field.getByteOffset(), s7Field.getBitOffset());
    }

}
