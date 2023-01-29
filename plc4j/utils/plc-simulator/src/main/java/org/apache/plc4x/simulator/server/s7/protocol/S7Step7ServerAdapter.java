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
package org.apache.plc4x.simulator.server.s7.protocol;

import io.netty.channel.*;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.simulator.model.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.BitSet;
import java.util.List;

public class S7Step7ServerAdapter extends ChannelInboundHandlerAdapter {

    private static final Logger LOGGER = LoggerFactory.getLogger(S7Step7ServerAdapter.class);

    private Context context;

    private State state;

    // COTP parameters
    private static final int localReference = 42;
    private int remoteReference = -1;
    private COTPProtocolClass protocolClass;
    private static final int localTsapId = 1;
    private int remoteTsapId = -1;
    private static final COTPTpduSize maxTpduSize = COTPTpduSize.SIZE_256;
    private COTPTpduSize tpduSize;
    // S7 parameters
    // Set this to 1 as we don't want to handle stuff in parallel
    private static final int maxAmqCaller = 1;
    private int amqCaller;
    // Set this to 1 as we don't want to handle stuff in parallel
    private static final int maxAmqCallee = 1;
    private int amqCallee;
    private static final int maxPduLength = 240;
    private int pduLength;

    public S7Step7ServerAdapter(Context context) {
        this.context = context;
        state = State.INITIAL;
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (msg instanceof TPKTPacket) {
            TPKTPacket packet = (TPKTPacket) msg;
            final COTPPacket cotpPacket = packet.getPayload();
            switch (state) {
                case INITIAL: {
                    if (!(cotpPacket instanceof COTPPacketConnectionRequest)) {
                        LOGGER.error("Expecting COTP Connection-Request");
                        return;
                    }

                    COTPTpduSize proposedTpduSize = null;
                    COTPPacketConnectionRequest cotpConnectionRequest = (COTPPacketConnectionRequest) cotpPacket;
                    for (COTPParameter parameter : cotpConnectionRequest.getParameters()) {
                        if (parameter instanceof COTPParameterCalledTsap) {
                            // this is actually ignored as it doesn't contain any information.
                        } else if (parameter instanceof COTPParameterCallingTsap) {
                            COTPParameterCallingTsap callingTsapParameter = (COTPParameterCallingTsap) parameter;
                            remoteTsapId = callingTsapParameter.getTsapId();
                        } else if (parameter instanceof COTPParameterTpduSize) {
                            COTPParameterTpduSize tpduSizeParameter = (COTPParameterTpduSize) parameter;
                            proposedTpduSize = tpduSizeParameter.getTpduSize();
                        } else {
                            LOGGER.error(String.format("Unexpected COTP Connection-Request Parameter %s",
                                parameter.getClass().getName()));
                            return;
                        }
                    }
                    if (proposedTpduSize == null) {
                        LOGGER.error("Missing COTP Connection-Request Parameter Tpdu Size");
                        return;
                    }

                    remoteReference = cotpConnectionRequest.getSourceReference();
                    protocolClass = cotpConnectionRequest.getProtocolClass();
                    tpduSize = (proposedTpduSize.getSizeInBytes() > maxTpduSize.getSizeInBytes()) ? maxTpduSize : proposedTpduSize;

                    // Prepare a response and send it back to the remote.
                    List<COTPParameter> parameters = new ArrayList<>();
                    parameters.add(new COTPParameterCalledTsap(remoteTsapId));
                    parameters.add(new COTPParameterCallingTsap(localTsapId));
                    parameters.add(new COTPParameterTpduSize(tpduSize));
                    COTPPacketConnectionResponse response = new COTPPacketConnectionResponse(
                        parameters, null, remoteReference, localReference, protocolClass
                    );
                    ctx.writeAndFlush(new TPKTPacket(response));

                    state = State.COTP_CONNECTED;
                    break;
                }
                case COTP_CONNECTED: {
                    if (!(cotpPacket instanceof COTPPacketData)) {
                        LOGGER.error("Expecting COTP Data packet");
                        return;
                    }

                    COTPPacketData packetData = (COTPPacketData) cotpPacket;
                    final short cotpTpduRef = packetData.getTpduRef();
                    final S7Message payload = packetData.getPayload();
                    if (!(payload instanceof S7MessageRequest)) {
                        LOGGER.error("Expecting S7 Message Request");
                        return;
                    }
                    S7MessageRequest s7MessageRequest = (S7MessageRequest) payload;
                    final int s7TpduReference = s7MessageRequest.getTpduReference();
                    final S7Parameter s7Parameter = s7MessageRequest.getParameter();
                    if (!(s7Parameter instanceof S7ParameterSetupCommunication)) {
                        LOGGER.error("Expecting S7 Message Request containing a S7 Setup Communication Parameter");
                        return;
                    }
                    S7ParameterSetupCommunication s7ParameterSetupCommunication =
                        (S7ParameterSetupCommunication) s7Parameter;

                    amqCaller = Math.min(s7ParameterSetupCommunication.getMaxAmqCaller(), maxAmqCaller);
                    amqCallee = Math.min(s7ParameterSetupCommunication.getMaxAmqCallee(), maxAmqCallee);
                    pduLength = Math.min(s7ParameterSetupCommunication.getPduLength(), maxPduLength);

                    S7ParameterSetupCommunication s7ParameterSetupCommunicationResponse =
                        new S7ParameterSetupCommunication(amqCaller, amqCallee, pduLength);
                    // TODO should send S7MessageResponseData
                    S7MessageResponseData s7MessageResponse = new S7MessageResponseData(
                        s7TpduReference, s7ParameterSetupCommunicationResponse, null, (short) 0, (short) 0);
                    ctx.writeAndFlush(new TPKTPacket(new COTPPacketData(null, s7MessageResponse, true, cotpTpduRef)));

                    state = State.S7_CONNECTED;
                    break;
                }
                case S7_CONNECTED: {
                    if (!(cotpPacket instanceof COTPPacketData)) {
                        LOGGER.error("Expecting COTP Data packet");
                        return;
                    }

                    COTPPacketData packetData = (COTPPacketData) cotpPacket;
                    final short cotpTpduRef = packetData.getTpduRef();
                    final S7Message payload = packetData.getPayload();
                    if (payload instanceof S7MessageUserData) {
                        S7MessageUserData s7MessageUserData = (S7MessageUserData) payload;
                        final int s7TpduReference = s7MessageUserData.getTpduReference();
                        final S7Parameter s7Parameter = s7MessageUserData.getParameter();
                        if (s7Parameter instanceof S7ParameterUserData) {
                            S7ParameterUserData userDataParameter = (S7ParameterUserData) s7Parameter;
                            for (S7ParameterUserDataItem item : userDataParameter.getItems()) {
                                if (item instanceof S7ParameterUserDataItemCPUFunctions) {
                                    S7ParameterUserDataItemCPUFunctions function =
                                        (S7ParameterUserDataItemCPUFunctions) item;
                                    final S7PayloadUserData userDataPayload =
                                        (S7PayloadUserData) s7MessageUserData.getPayload();

                                    for (S7PayloadUserDataItem userDataPayloadItem : userDataPayload.getItems()) {
                                        if (userDataPayloadItem instanceof S7PayloadUserDataItemCpuFunctionReadSzlRequest) {
                                            S7PayloadUserDataItemCpuFunctionReadSzlRequest readSzlRequestPayload =
                                                (S7PayloadUserDataItemCpuFunctionReadSzlRequest) userDataPayloadItem;

                                            final SzlId szlId = readSzlRequestPayload.getSzlId();
                                            // This is a request to list the type of device
                                            if ((szlId.getTypeClass() == SzlModuleTypeClass.CPU) &&
                                                (szlId.getSublistList() == SzlSublist.MODULE_IDENTIFICATION)) {

                                                S7ParameterUserDataItemCPUFunctions readSzlResponseParameter =
                                                    new S7ParameterUserDataItemCPUFunctions((short) 0x12,
                                                        (byte) 0x08, function.getCpuFunctionGroup(),
                                                        function.getCpuSubfunction(), (short) 1,
                                                        (short) 0, (short) 0, 0);

                                                // This is the product number of a S7-1200
                                                List<SzlDataTreeItem> items = new ArrayList<>();
                                                items.add(new SzlDataTreeItem((short) 0x0001,
                                                    "6ES7 212-1BD30-0XB0 ".getBytes(), 0x2020, 0x0001, 0x2020));

                                                S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponsePayload =
                                                    new S7PayloadUserDataItemCpuFunctionReadSzlResponse(
                                                        DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, szlId,
                                                        readSzlRequestPayload.getSzlIndex(), items);

                                                List<S7ParameterUserDataItem> responseParameterItems = new ArrayList<>();
                                                responseParameterItems.add(readSzlResponseParameter);
                                                S7ParameterUserData responseParameterUserData =
                                                    new S7ParameterUserData(responseParameterItems);

                                                List<S7PayloadUserDataItem> responsePayloadItems = new ArrayList<>();
                                                responsePayloadItems.add(readSzlResponsePayload);
                                                S7PayloadUserData responsePayloadUserData =
                                                    new S7PayloadUserData(responsePayloadItems);

                                                S7Message s7ResponseMessage = new S7MessageUserData(s7TpduReference,
                                                    responseParameterUserData, responsePayloadUserData);
                                                ctx.writeAndFlush(new TPKTPacket(new COTPPacketData(null, s7ResponseMessage, true, cotpTpduRef)));
                                            } else {
                                                LOGGER.error("Not able to respond to the given request Read SZL with SZL type class " +
                                                    szlId.getTypeClass().name() + " and SZL sublist " + szlId.getSublistList().name());
                                            }

                                        }
                                    }
                                }
                            }
                        } else {
                            LOGGER.error("Unsupported type of S7MessageUserData parameter " +
                                s7Parameter.getClass().getName());
                        }
                    } else {
                        if (cotpPacket.getPayload() instanceof S7MessageRequest) {
                            S7MessageRequest request = (S7MessageRequest) cotpPacket.getPayload();
                            if (request.getParameter() instanceof S7ParameterReadVarRequest) {
                                S7ParameterReadVarRequest readVarRequestParameter =
                                    (S7ParameterReadVarRequest) request.getParameter();
                                List<S7VarRequestParameterItem> items = readVarRequestParameter.getItems();
                                List<S7VarPayloadDataItem> payloadItems = new ArrayList<>();
                                for (S7VarRequestParameterItem item : items) {
                                    if (item instanceof S7VarRequestParameterItemAddress) {
                                        S7VarRequestParameterItemAddress address =
                                            (S7VarRequestParameterItemAddress) item;
                                        final S7Address address1 = address.getAddress();
                                        if (address1 instanceof S7AddressAny) {
                                            S7AddressAny addressAny = (S7AddressAny) address1;
                                            switch (addressAny.getArea()) {
                                                case DATA_BLOCKS: {
                                                    final int dataBlockNumber = addressAny.getDbNumber();
                                                    if (dataBlockNumber != 1) {
                                                        // TODO: Return unknown object.
                                                    }
                                                    final int numberOfElements = addressAny.getNumberOfElements();
                                                    if (numberOfElements != 1) {
                                                        // TODO: Return invalid address.
                                                    }
                                                    final int byteAddress = addressAny.getByteAddress();
                                                    if (byteAddress != 0) {
                                                        // TODO: Return invalid address.
                                                    }
                                                    final byte bitAddress = addressAny.getBitAddress();
                                                    switch (addressAny.getTransportSize()) {
                                                        case BOOL:
                                                            payloadItems.add(new S7VarPayloadDataItem(DataTransportErrorCode.OK, DataTransportSize.BIT, new byte[]{1}, true));
                                                            break;
                                                        case INT:
                                                        case UINT: {
                                                            String firstKey = context.getMemory().keySet().iterator().next();
                                                            Object value = context.getMemory().get(firstKey);
                                                            short shortValue = 42; // ((Number) value).shortValue();
                                                            byte[] data = new byte[2];
                                                            data[1] = (byte) (shortValue & 0xff);
                                                            data[0] = (byte) ((shortValue >> 8) & 0xff);
                                                            payloadItems.add(new S7VarPayloadDataItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, data, true));
                                                            break;
                                                        }
                                                        default: {
                                                            // TODO: Return invalid address.
                                                        }
                                                    }
                                                    break;
                                                }
                                                case INPUTS:
                                                case OUTPUTS: {
                                                    final int ioNumber = (addressAny.getByteAddress() * 8) + addressAny.getBitAddress();
                                                    final int numElements = (addressAny.getTransportSize() == TransportSize.BOOL) ?
                                                        addressAny.getNumberOfElements() : addressAny.getTransportSize().getSizeInBytes() * 8;
                                                    final BitSet bitSet = toBitSet(context.getDigitalInputs(), ioNumber, numElements);
                                                    final byte[] data = Arrays.copyOf(bitSet.toByteArray(), (numElements + 7) / 8);
                                                    payloadItems.add(new S7VarPayloadDataItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, data, true));
                                                    break;
                                                }
                                            }
                                        }
                                    }
                                }
                                S7ParameterReadVarResponse readVarResponseParameter = new S7ParameterReadVarResponse((short) items.size());
                                S7PayloadReadVarResponse readVarResponsePayload = new S7PayloadReadVarResponse(payloadItems);
                                S7MessageResponseData response = new S7MessageResponseData(request.getTpduReference(),
                                    readVarResponseParameter, readVarResponsePayload, (short) 0x00, (short) 0x00);
                                ctx.writeAndFlush(new TPKTPacket(new COTPPacketData(null, response, true, cotpTpduRef)));
                            } else if (request.getParameter() instanceof S7ParameterWriteVarRequest) {
                                S7ParameterWriteVarRequest writeVarRequestParameter =
                                    (S7ParameterWriteVarRequest) request.getParameter();

                            } else {
                                LOGGER.error("Unsupported type of S7MessageRequest parameter " +
                                    request.getParameter().getClass().getName());
                            }
                        } else {
                            LOGGER.error("Unsupported type of message " + payload.getClass().getName());
                        }
                    }
                    break;
                }
                default:
                    throw new IllegalStateException("Unexpected value: " + state);
            }
        }
    }

    private enum State {
        INITIAL,
        COTP_CONNECTED,
        S7_CONNECTED
    }

    private BitSet toBitSet(List<Boolean> booleans, int startIndex, int numElements) {
        BitSet bitSet = new BitSet(booleans.size());
        for (int i = 0; i < Math.min(booleans.size() - startIndex, numElements); i++) {
            bitSet.set(i, booleans.get(i + startIndex));
        }
        return bitSet;
    }

}
