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
package org.apache.plc4x.simulator.server.s7.protocol;

import io.netty.channel.*;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class S7Step7ServerProtocol extends ChannelInboundHandlerAdapter {

    private static final Logger LOGGER = LoggerFactory.getLogger(S7Step7ServerProtocol.class);

    private State state = State.INITIAL;

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

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if(msg instanceof TPKTPacket) {
            TPKTPacket packet = (TPKTPacket) msg;
            switch (state) {
                case INITIAL: {
                    final COTPPacket cotpPacket = packet.getPayload();
                    if (!(cotpPacket instanceof COTPPacketConnectionRequest)) {
                        LOGGER.error("Expecting COTP Connection-Request");
                        return;
                    }

                    COTPTpduSize proposedTpduSize = null;
                    COTPPacketConnectionRequest cotpConnectionRequest = (COTPPacketConnectionRequest) cotpPacket;
                    for (COTPParameter parameter : cotpConnectionRequest.getParameters()) {
                        if (parameter instanceof COTPParameterCalledTsap) {
                            COTPParameterCalledTsap calledTsapParameter = (COTPParameterCalledTsap) parameter;
                        } else if (parameter instanceof COTPParameterCallingTsap) {
                            COTPParameterCallingTsap callingTsapParameter = (COTPParameterCallingTsap) parameter;
                            remoteTsapId = callingTsapParameter.getTsapId();
                        } else if (parameter instanceof COTPParameterTpduSize) {
                            COTPParameterTpduSize tpduSizeParameter = (COTPParameterTpduSize) parameter;
                            proposedTpduSize = tpduSizeParameter.getTpduSize();
                        } else {
                            LOGGER.error(
                                "Unexpected COTP Connection-Request Parameter " + parameter.getClass().getName());
                            return;
                        }
                    }

                    remoteReference = cotpConnectionRequest.getSourceReference();
                    protocolClass = cotpConnectionRequest.getProtocolClass();
                    tpduSize = (proposedTpduSize.getSizeInBytes() > maxTpduSize.getSizeInBytes()) ? maxTpduSize : proposedTpduSize;

                    // Prepare a response and send it back to the remote.
                    COTPParameter[] parameters = new COTPParameter[3];
                    parameters[0] = new COTPParameterCalledTsap(remoteTsapId);
                    parameters[1] = new COTPParameterCallingTsap(localTsapId);
                    parameters[2] = new COTPParameterTpduSize(tpduSize);
                    COTPPacketConnectionResponse response = new COTPPacketConnectionResponse(
                        parameters, null, remoteReference, localReference, protocolClass);
                    ctx.writeAndFlush(new TPKTPacket(response));

                    state = State.COTP_CONNECTED;
                    break;
                }
                case COTP_CONNECTED: {
                    final COTPPacket cotpPacket = packet.getPayload();
                    if (!(cotpPacket instanceof COTPPacketData)) {
                        LOGGER.error("Expecting COTP Data packet");
                        return;
                    }

                    COTPPacketData packetData = (COTPPacketData) cotpPacket;
                    final short cotpTpduRef = packetData.getTpduRef();
                    final S7Message payload = packetData.getPayload();
                    if(!(payload instanceof S7MessageRequest)) {
                        LOGGER.error("Expecting S7 Message Request");
                        return;
                    }
                    S7MessageRequest s7MessageRequest = (S7MessageRequest) payload;
                    final int s7TpduReference = s7MessageRequest.getTpduReference();
                    final S7Parameter s7Parameter = s7MessageRequest.getParameter();
                    if(!(s7Parameter instanceof S7ParameterSetupCommunication)) {
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
                    S7MessageResponse s7MessageResponse = new S7MessageResponse(
                        s7TpduReference, s7ParameterSetupCommunicationResponse, new S7PayloadSetupCommunication(),
                        (short) 0, (short) 0);
                    ctx.writeAndFlush(new TPKTPacket(new COTPPacketData(null, s7MessageResponse, true, cotpTpduRef)));

                    state = State.S7_CONNECTED;
                    break;
                }
                case S7_CONNECTED: {
                    final COTPPacket cotpPacket = packet.getPayload();
                    if (!(cotpPacket instanceof COTPPacketData)) {
                        LOGGER.error("Expecting COTP Data packet");
                        return;
                    }

                    COTPPacketData packetData = (COTPPacketData) cotpPacket;
                    final short cotpTpduRef = packetData.getTpduRef();
                    final S7Message payload = packetData.getPayload();
                    if(payload instanceof S7MessageUserData) {
                        S7MessageUserData s7MessageUserData = (S7MessageUserData) payload;
                        final int s7TpduReference = s7MessageUserData.getTpduReference();
                        final S7Parameter s7Parameter = s7MessageUserData.getParameter();
                        if(s7Parameter instanceof S7ParameterUserData) {
                            S7ParameterUserData userDataParameter = (S7ParameterUserData) s7Parameter;
                            for (S7ParameterUserDataItem item : userDataParameter.getItems()) {
                                if(item instanceof S7ParameterUserDataItemCPUFunctions) {
                                    S7ParameterUserDataItemCPUFunctions function =
                                        (S7ParameterUserDataItemCPUFunctions) item;
                                    final S7PayloadUserData userDataPayload =
                                        (S7PayloadUserData) s7MessageUserData.getPayload();

                                    for (S7PayloadUserDataItem userDataPayloadItem : userDataPayload.getItems()) {
                                        if(userDataPayloadItem instanceof S7PayloadUserDataItemCpuFunctionReadSzlRequest) {
                                            S7PayloadUserDataItemCpuFunctionReadSzlRequest readSzlRequestPayload =
                                                (S7PayloadUserDataItemCpuFunctionReadSzlRequest) userDataPayloadItem;

                                            final SzlId szlId = readSzlRequestPayload.getSzlId();
                                            // This is a request to list the type of device
                                            if((szlId.getTypeClass() == SzlModuleTypeClass.CPU) &&
                                                (szlId.getSublistList() == SzlSublist.MODULE_IDENTIFICATION)) {

                                                S7ParameterUserDataItemCPUFunctions readSzlResponseParameter =
                                                    new S7ParameterUserDataItemCPUFunctions((short) 0x12,
                                                        (byte) 0x08, function.getCpuFunctionGroup(),
                                                        function.getCpuSubfunction(), (short) 1,
                                                        (short) 0, (short) 0, 0);

                                                SzlDataTreeItem[] items = new SzlDataTreeItem[1];
                                                items[0] = new SzlDataTreeItem((short) 0x0001,
                                                    "6ES7 212-1BD30-0XB0 ".getBytes(), 0x2020, 0x0001, 0x2020);

                                                S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponsePayload =
                                                    new S7PayloadUserDataItemCpuFunctionReadSzlResponse(
                                                        (short) 0xFF, DataTransportSize.OCTET_STRING, szlId,
                                                        readSzlRequestPayload.getSzlIndex(), items);

                                                S7ParameterUserDataItem[] responseParameterItems =
                                                    new S7ParameterUserDataItem[1];
                                                responseParameterItems[0] = readSzlResponseParameter;
                                                S7ParameterUserData responseParameterUserData =
                                                    new S7ParameterUserData(responseParameterItems);

                                                S7PayloadUserDataItem[] responsePayloadItems =
                                                    new S7PayloadUserDataItem[1];
                                                responsePayloadItems[0] = readSzlResponsePayload;
                                                S7PayloadUserData responsePayloadUserData =
                                                    new S7PayloadUserData(responsePayloadItems);

                                                S7Message s7ResponseMessage = new S7MessageUserData(s7TpduReference,
                                                    responseParameterUserData, responsePayloadUserData);
                                                ctx.writeAndFlush(new TPKTPacket(new COTPPacketData(null, s7ResponseMessage, true, cotpTpduRef)));
                                            } else {
                                                LOGGER.error("Not able to respond to the given request Read SZL with SZL type class " +
                                                    szlId.getTypeClass().name() + " and SZL sublise " + szlId.getSublistList().name());
                                            }

                                        }
                                    }
                                }
                            }
                        } else {
                            LOGGER.error("Unsupported type of S7MessageUserData parameter " +
                                s7Parameter.getClass().getName());
                            return;
                        }
                    } else {
                        LOGGER.error("Unsupported type of message " + payload.getClass().getName());
                        return;
                    }
                    break;
                }
            }
        }
    }

    private enum State {
        INITIAL,
        COTP_CONNECTED,
        S7_CONNECTED
    }

}
