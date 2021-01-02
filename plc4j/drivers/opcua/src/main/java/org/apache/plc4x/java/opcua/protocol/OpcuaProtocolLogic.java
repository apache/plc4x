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
package org.apache.plc4x.java.opcua.protocol;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.readwrite.io.*;
import org.apache.plc4x.java.opcua.readwrite.types.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.UnsupportedEncodingException;
import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.IntStream;
import java.nio.charset.StandardCharsets;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 */
public class OpcuaProtocolLogic extends Plc4xProtocolBase<OpcuaAPU> implements HasConfiguration<OpcuaConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private final AtomicInteger tpduGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;

    @Override
    public void setConfiguration(OpcuaConfiguration configuration) {
    }

    @Override
    public void close(ConversationContext<OpcuaAPU> context) {
        // Nothing to do here ...
    }

    @Override
    public void onConnect(ConversationContext<OpcuaAPU> context) {
        // Only the TCP transport supports login.
        LOGGER.info("Opcua Driver running in ACTIVE mode.");

        final String endpoint = "opc.tcp://127.0.0.1:12687/plc4x";
        OpcuaHelloRequest hello = new OpcuaHelloRequest("F",
                                                        0,
                                                        65535,
                                                        65535,
                                                        2097152,
                                                        64,
                                                        endpoint.length(),
                                                        endpoint);

        context.sendRequest(new OpcuaAPU(hello))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaAcknowledgeResponse)
            .unwrap(p -> (OpcuaAcknowledgeResponse) p.getMessage())
            .handle(opcuaAcknowledgeResponse -> {
                LOGGER.debug("Got Hello Response Connection Response");

                NodeIdTwoByte authenticationToken = new NodeIdTwoByte(NodeIdType.nodeIdTypeTwoByte,
                                                                    new TwoByteNodeId((short) 0));

                ExpandedNodeId expandedNodeId = new ExpandedNodeIdFourByte(false,
                                                                    false,
                                                                    new PascalString(-1,null),
                                                                    1L,
                                                                    new FourByteNodeId((short) 0, 466));

                ExpandedNodeId extExpandedNodeId = new ExpandedNodeIdTwoByte(false,
                                                                    false,
                                                                    null,
                                                                    null,
                                                                    new TwoByteNodeId((short) 0));

                ExtensionObject extObject = new ExtensionObject(extExpandedNodeId, (short) 0, null, null);

                RequestHeader requestHeader = new RequestHeader(authenticationToken,
                                                                (System.currentTimeMillis() * 10000) + 116444736000000000L,
                                                                0L,
                                                                0L,
                                                                new PascalString(-1, null),
                                                                10000L,
                                                                extObject);



                OpenSecureChannelRequest openrequest = new OpenSecureChannelRequest((byte) 1,
                                                                (byte) 0,
                                                                requestHeader,
                                                                0L,
                                                                SecurityTokenRequestType.securityTokenRequestTypeIssue,
                                                                MessageSecurityMode.messageSecurityModeNone,
                                                                new PascalString(-1, null),
                                                                36000000);


                String nameSpace = "http://opcfoundation.org/UA/SecurityPolicy#None";
                OpcuaOpenRequest openRequest = new OpcuaOpenRequest("F",
                                                                0,
                                                                nameSpace.length(),
                                                                nameSpace,
                                                                -1,
                                                                "",
                                                                -1,
                                                                "",
                                                                1,
                                                                1,
                                                                openrequest);

                context.sendRequest(new OpcuaAPU(openRequest))
                    .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                    .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
                    .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
                    .handle(opcuaOpenResponse -> {
                        LOGGER.debug("Got Secure Response Connection Response");
                        OpenSecureChannelResponse openSecureChannelResponse = (OpenSecureChannelResponse) opcuaOpenResponse.getMessage();
                        Integer tokenId = (int) openSecureChannelResponse.getSecurityToken().getTokenId();
                        Integer channelId = (int) openSecureChannelResponse.getSecurityToken().getChannelId();
                        Integer nextSequenceNumber = opcuaOpenResponse.getSequenceNumber() + 1;
                        Integer nextRequestId = opcuaOpenResponse.getRequestId() + 1;

                        NodeIdTwoByte authenticationToken2 = new NodeIdTwoByte(NodeIdType.nodeIdTypeTwoByte,
                                                                            new TwoByteNodeId((short) 0));

                        ExpandedNodeId extExpandedNodeId2 = new ExpandedNodeIdTwoByte(false,
                                                                            false,
                                                                            NodeIdType.nodeIdTypeTwoByte,
                                                                            null,
                                                                            null,
                                                                            new TwoByteNodeId((short) 0));

                        ExtensionObject extObject2 = new ExtensionObject(extExpandedNodeId2, (short) 0);

                        RequestHeader requestHeader2 = new RequestHeader(authenticationToken2,
                                                                        (System.currentTimeMillis() * 10000) + 116444736000000000L,
                                                                        0L,
                                                                        0L,
                                                                        new PascalString(-1, null),
                                                                        10000L,
                                                                        extObject2);

                        String applicationUri = "urn:eclipse:milo:plc4x:client";
                        String productUri = "urn:eclipse:milo:plc4x:client";
                        String text = "eclipse milo opc-ua client of the apache PLC4X:PLC4J project";
                        LocalizedText applicationName = new LocalizedText((short) 0,
                                                                          true,
                                                                          true,
                                                                          new PascalString(2, "en"),
                                                                          new PascalString(text.length(), text));
                        PascalString gatewayServerUri = new PascalString(-1, null);
                        PascalString discoveryProfileUri = new PascalString(-1, null);
                        int noOfDiscoveryUrls = -1;
                        PascalString discoveryUrls = null;

                        ApplicationDescription clientDescription = new ApplicationDescription(new PascalString(applicationUri.length(), applicationUri),
                                                                                    new PascalString(productUri.length(), productUri),
                                                                                    applicationName,
                                                                                    ApplicationType.applicationTypeClient,
                                                                                    gatewayServerUri,
                                                                                    discoveryProfileUri,
                                                                                    noOfDiscoveryUrls,
                                                                                    discoveryUrls);

                        String endpoint2 = "opc.tcp://127.0.0.1:12687/plc4x";
                        String sessionName = "UaSession:eclipse milo opc-ua client of the apache PLC4X:PLC4J project:" + System.currentTimeMillis();
                        String clientNonce = "764287368237654873259869867";

                        CreateSessionRequest createSessionRequest = new CreateSessionRequest((byte) 1,
                                                                        (byte) 0,
                                                                        requestHeader2,
                                                                        clientDescription,
                                                                        new PascalString(-1, null),
                                                                        new PascalString(endpoint2.length(), endpoint2),
                                                                        new PascalString(sessionName.length(), sessionName),
                                                                        new PascalString(clientNonce.length(), clientNonce),
                                                                        new PascalString(-1, null),
                                                                        120000L,
                                                                        0L);

                        OpcuaMessageRequest messageRequest = new OpcuaMessageRequest("F",
                                                                        channelId,
                                                                        tokenId,
                                                                        nextSequenceNumber,
                                                                        nextRequestId,
                                                                        createSessionRequest);

                        context.sendRequest(new OpcuaAPU(messageRequest))
                            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                            .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                            .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                            .handle(opcuaMessageResponse -> {
                                LOGGER.debug("Got Create Session Response Connection Response");
                                CreateSessionResponse createSessionResponse = (CreateSessionResponse) opcuaMessageResponse.getMessage();

                                NodeIdByteString authenticationToken3 = (NodeIdByteString) createSessionResponse.getAuthenticationToken();
                                Integer tokenId2 = (int) opcuaMessageResponse.getSecureTokenId();
                                Integer channelId2 = (int) opcuaMessageResponse.getSecureChannelId();
                                Integer nextSequenceNumber2 = opcuaMessageResponse.getSequenceNumber() + 1;
                                Integer nextRequestId2 = opcuaMessageResponse.getRequestId() + 1;


                                ExpandedNodeId extExpandedNodeId3 = new ExpandedNodeIdTwoByte(false,
                                                                                    false,
                                                                                    NodeIdType.nodeIdTypeTwoByte,
                                                                                    null,
                                                                                    null,
                                                                                    new TwoByteNodeId((short) 0));
                                System.out.println("(((((((((((((((" + extExpandedNodeId3.getLengthInBytes());

                                ExtensionObject extObject3 = new ExtensionObject(extExpandedNodeId3, (short) 0);

                                System.out.println("(((((((((((((((" + extObject3.getLengthInBytes());
                                System.out.println("(((((((((((((((" + authenticationToken3.getLengthInBytes());
                                System.out.println("@@@@@@@@@@@@@@@@@" + authenticationToken3.getId().getIdentifier().getStringLength());
                                System.out.println("@@@@@@@@@@@@@@@@@" + authenticationToken3.getId().getIdentifier().getStringValue().length());

                                RequestHeader requestHeader3 = new RequestHeader(authenticationToken3,
                                                                                (System.currentTimeMillis() * 10000) + 116444736000000000L,
                                                                                1L,
                                                                                0L,
                                                                                new PascalString(-1, null),
                                                                                10000L,
                                                                                extObject3);

                                System.out.println("(((((((((((((((" + requestHeader3.getLengthInBytes());

                                SignatureData clientSignature = new SignatureData(new PascalString(-1, null), new PascalString(-1, null));

                                System.out.println("(((((((((((((((" + clientSignature.getLengthInBytes());

                                SignedSoftwareCertificate[] signedSoftwareCertificate = new SignedSoftwareCertificate[1];

                                signedSoftwareCertificate[0] = new SignedSoftwareCertificate(new PascalString(-1, null), new PascalString(-1, null));

                                ExpandedNodeId extExpandedNodeId4 = new ExpandedNodeIdFourByte(false,
                                                                                    false,
                                                                                    NodeIdType.nodeIdTypeFourByte,
                                                                                    null,
                                                                                    null,
                                                                                    new FourByteNodeId((short) 1,  321));

                                System.out.println("(((((((((((((((" + extExpandedNodeId4.getLengthInBytes());


                                ExtensionObject useridentityToken = new ExtensionObject(extExpandedNodeId4, (short) 1);

                                System.out.println("(((((((((((((((" + useridentityToken.getLengthInBytes());

                                String endpoint3 = "opc.tcp://127.0.0.1:12687/plc4x";

                                ActivateSessionRequest activateSessionRequest = new ActivateSessionRequest((byte) 1,
                                                                                (byte) 0,
                                                                                requestHeader3,
                                                                                clientSignature,
                                                                                0,
                                                                                null,
                                                                                0,
                                                                                null,
                                                                                useridentityToken,
                                                                                clientSignature);

                                System.out.println("(((((((((((((((" + activateSessionRequest.getLengthInBytes());

                                OpcuaMessageRequest activateMessageRequest = new OpcuaMessageRequest("F",
                                                                                channelId2,
                                                                                tokenId2,
                                                                                nextSequenceNumber2,
                                                                                nextRequestId2,
                                                                                activateSessionRequest);

                                System.out.println("(((((((((((((((" + activateMessageRequest.getLengthInBytes());

                                context.sendRequest(new OpcuaAPU(activateMessageRequest))
                                    .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                                    .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                                    .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                                    .handle(opcuaActivateResponse -> {
                                        LOGGER.debug("Got Activate Session Response Connection Response");

                                    });


                            });

                    });
            });
    }

}
