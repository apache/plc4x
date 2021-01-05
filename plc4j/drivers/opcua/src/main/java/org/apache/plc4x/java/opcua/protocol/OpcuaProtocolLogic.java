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
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.DefaultNettyPlcConnection;
import org.apache.plc4x.java.spi.transport.Transport;
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
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.UnsupportedEncodingException;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.regex.Matcher;
import java.util.stream.IntStream;
import java.nio.charset.StandardCharsets;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 */
public class OpcuaProtocolLogic extends Plc4xProtocolBase<OpcuaAPU> implements HasConfiguration<OpcuaConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);
    private static final String CHUNK = "F";
    private static final int VERSION = 0;
    private static final int DEFAULT_RECEIVE_BUFFER_SIZE = 65535;
    private static final int DEFAULT_SEND_BUFFER_SIZE = 65535;
    private static final int DEFAULT_MAX_MESSAGE_SIZE = 2097152;
    private static final int DEFAULT_MAX_CHUNK_COUNT = 64;
    private NodeId authenticationToken = new NodeIdTwoByte(NodeIdType.nodeIdTypeTwoByte, new TwoByteNodeId((short) 0));
    private static final PascalString NULL_STRING = new PascalString(-1,null);
    private static ExpandedNodeId NULL_EXPANDED_NODEID = new ExpandedNodeIdTwoByte(false,
                                                                                    false,
                                                                                    null,
                                                                                    null,
                                                                                    new TwoByteNodeId((short) 0));
    private static final ExtensionObject NULL_EXTENSION_OBJECT = new ExtensionObject(NULL_EXPANDED_NODEID,
                                                                                        (short) 0,
                                                                                null,               //Body Length
                                                                                    null);               // Body
    private static final long epochOffset = 116444736000000000L;         //Offset between OPC UA epoch time and linux epoch time.
    private static final int DEFAULT_CONNECTION_LIFETIME = 36000000;
    private static final String nameSpaceSecurityPolicyNone = "http://opcfoundation.org/UA/SecurityPolicy#None";
    private static final String applicationUri = "urn:apache:plc4x:client";
    private static final String productUri = "urn:apache:plc4x:client";
    private static final String applicationText = "OPCUA client for the Apache PLC4X:PLC4J project";

    private String sessionName = "UaSession:" + applicationText + ":" + RandomStringUtils.random(20, true, true);
    private String clientNonce = RandomStringUtils.random(40, true, true);
    private RequestTransactionManager tm;

    private String endpoint;
    private AtomicInteger transactionIdentifierGenerator = new AtomicInteger(1);
    private AtomicInteger requestHandleGenerator = new AtomicInteger(1);
    private AtomicInteger tokenId = new AtomicInteger(1);
    private AtomicInteger channelId = new AtomicInteger(1);

    @Override
    public void setConfiguration(OpcuaConfiguration configuration) {
        this.endpoint = configuration.getEndpoint();
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void close(ConversationContext<OpcuaAPU> context) {
        // Nothing to do here ...
    }

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<OpcuaAPU> context) {
        // Only the TCP transport supports login.
        LOGGER.info("Opcua Driver running in ACTIVE mode.");

        OpcuaHelloRequest hello = new OpcuaHelloRequest(CHUNK,
            VERSION,
            DEFAULT_RECEIVE_BUFFER_SIZE,
            DEFAULT_SEND_BUFFER_SIZE,
            DEFAULT_MAX_MESSAGE_SIZE,
            DEFAULT_MAX_CHUNK_COUNT,
            this.endpoint.length(),
            this.endpoint);

        context.sendRequest(new OpcuaAPU(hello))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaAcknowledgeResponse)
            .unwrap(p -> (OpcuaAcknowledgeResponse) p.getMessage())
            .handle(opcuaAcknowledgeResponse -> {
                LOGGER.debug("Got Hello Response Connection Response");
                onConnectOpenSecureChannel(context, opcuaAcknowledgeResponse);
            });
    }

    public void onConnectOpenSecureChannel(ConversationContext<OpcuaAPU> context, OpcuaAcknowledgeResponse opcuaAcknowledgeResponse) {

        int transactionId = transactionIdentifierGenerator.getAndIncrement();
        if(transactionIdentifierGenerator.get() == 0xFFFF) {
            transactionIdentifierGenerator.set(1);
        }

        ExpandedNodeId expandedNodeId = new ExpandedNodeIdFourByte(false,           //Namespace Uri Specified
                                                                    false,            //Server Index Specified
                                                                    NULL_STRING,                      //Namespace Uri
                                                                    1L,                     //Server Index
                                                                    new FourByteNodeId((short) 0, 466));    //Identifier for OpenSecureChannel

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            10000L,
            NULL_EXTENSION_OBJECT);

        OpenSecureChannelRequest openSecureChannelRequest = new OpenSecureChannelRequest((byte) 1,
            (byte) 0,
            requestHeader,
            VERSION,
            SecurityTokenRequestType.securityTokenRequestTypeIssue,
            MessageSecurityMode.messageSecurityModeNone,
            NULL_STRING,
            DEFAULT_CONNECTION_LIFETIME);

        OpcuaOpenRequest openRequest = new OpcuaOpenRequest(CHUNK,
            0,
            new PascalString(nameSpaceSecurityPolicyNone.length(), nameSpaceSecurityPolicyNone),
            NULL_STRING,
            NULL_STRING,
            transactionId,
            transactionId,
            openSecureChannelRequest);

        context.sendRequest(new OpcuaAPU(openRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
            .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
            .handle(opcuaOpenResponse -> {
                LOGGER.debug("Got Secure Response Connection Response");
                try {
                    onConnectCreateSessionRequest(context, opcuaOpenResponse);
                } catch (PlcConnectionException e) {
                    LOGGER.error("Error occurred while connecting to OPC UA server");
                }
            });

    }

    public void onConnectCreateSessionRequest(ConversationContext<OpcuaAPU> context, OpcuaOpenResponse opcuaOpenResponse) throws PlcConnectionException {
        OpenSecureChannelResponse openSecureChannelResponse = (OpenSecureChannelResponse) opcuaOpenResponse.getMessage();
        tokenId.set((int) openSecureChannelResponse.getSecurityToken().getTokenId());
        channelId.set((int) openSecureChannelResponse.getSecurityToken().getChannelId());

        int transactionId = transactionIdentifierGenerator.getAndIncrement();
        if(transactionIdentifierGenerator.get() == 0xFFFF) {
            transactionIdentifierGenerator.set(1);
        }

        Integer nextSequenceNumber = opcuaOpenResponse.getSequenceNumber() + 1;
        Integer nextRequestId = opcuaOpenResponse.getRequestId() + 1;

        if (!(transactionId == nextSequenceNumber)) {
            LOGGER.error("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
            throw new PlcConnectionException("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
        }

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            0L,
            0L,
            NULL_STRING,
            10000L,
            NULL_EXTENSION_OBJECT);

        LocalizedText applicationName = new LocalizedText((short) 0,
            true,
            true,
            new PascalString("en".length(), "en"),
            new PascalString(applicationText.length(), applicationText));

        PascalString gatewayServerUri = NULL_STRING;
        PascalString discoveryProfileUri = NULL_STRING;
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

        clientNonce = RandomStringUtils.random(40, true, true);

        CreateSessionRequest createSessionRequest = new CreateSessionRequest((byte) 1,
            (byte) 0,
            requestHeader,
            clientDescription,
            NULL_STRING,
            new PascalString(endpoint.length(), endpoint),
            new PascalString(sessionName.length(), sessionName),
            new PascalString(clientNonce.length(), clientNonce),
            NULL_STRING,
            120000L,
            0L);

        OpcuaMessageRequest messageRequest = new OpcuaMessageRequest(CHUNK,
            channelId.get(),
            tokenId.get(),
            nextSequenceNumber,
            nextRequestId,
            createSessionRequest);

        context.sendRequest(new OpcuaAPU(messageRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
            .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
            .handle(opcuaMessageResponse -> {
                LOGGER.debug("Got Create Session Response Connection Response");
                try {
                    onConnectActivateSessionRequest(context, opcuaMessageResponse);
                } catch (PlcConnectionException e) {
                    LOGGER.error("Error occurred while connecting to OPC UA server");
                }
            });
    }

    private void onConnectActivateSessionRequest(ConversationContext<OpcuaAPU> context, OpcuaMessageResponse opcuaMessageResponse) throws PlcConnectionException {

        CreateSessionResponse createSessionResponse = (CreateSessionResponse) opcuaMessageResponse.getMessage();

        authenticationToken = (NodeIdByteString) createSessionResponse.getAuthenticationToken();
        tokenId.set((int) opcuaMessageResponse.getSecureTokenId());
        channelId.set((int) opcuaMessageResponse.getSecureChannelId());

        int transactionId = transactionIdentifierGenerator.getAndIncrement();
        if(transactionIdentifierGenerator.get() == 0xFFFF) {
            transactionIdentifierGenerator.set(1);
        }

        Integer nextSequenceNumber = opcuaMessageResponse.getSequenceNumber() + 1;
        Integer nextRequestId = opcuaMessageResponse.getRequestId() + 1;

        if (!(transactionId == nextSequenceNumber)) {
            LOGGER.error("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
            throw new PlcConnectionException("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
        }

        int requestHandle = requestHandleGenerator.getAndIncrement();
        if(requestHandleGenerator.get() == 0xFFFF) {
            requestHandleGenerator.set(1);
        }

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            requestHandle,
            0L,
            NULL_STRING,
            10000L,
            NULL_EXTENSION_OBJECT);

        SignatureData clientSignature = new SignatureData(NULL_STRING, NULL_STRING);

        SignedSoftwareCertificate[] signedSoftwareCertificate = new SignedSoftwareCertificate[1];

        signedSoftwareCertificate[0] = new SignedSoftwareCertificate(NULL_STRING, NULL_STRING);

        //Manually serialize this object
        PascalString anonymousIdentityToken = new PascalString("anonymous".length(), "anonymous");
        WriteBuffer buffer = new WriteBuffer(anonymousIdentityToken.getLengthInBytes(), true);
        try{
            PascalStringIO.staticSerialize(buffer, anonymousIdentityToken);
        } catch (ParseException e) {
            LOGGER.error("Failed to serialize the user identity token - " + anonymousIdentityToken.getStringValue());
            throw new PlcConnectionException("Failed to serialize the user identity token - " + anonymousIdentityToken.getStringValue());
        }

        ExpandedNodeId extExpandedNodeId4 = new ExpandedNodeIdFourByte(false,
            false,
            null,
            null,
            new FourByteNodeId((short) 0,  321));

        ExtensionObject userIdentityToken = new ExtensionObject(extExpandedNodeId4, (short) 1, buffer.getData().length, buffer.getData());

        ActivateSessionRequest activateSessionRequest = new ActivateSessionRequest((byte) 1,
            (byte) 0,
            requestHeader,
            clientSignature,
            0,
            null,
            0,
            null,
            userIdentityToken,
            clientSignature);

        OpcuaMessageRequest activateMessageRequest = new OpcuaMessageRequest(CHUNK,
            channelId.get(),
            tokenId.get(),
            nextSequenceNumber,
            nextRequestId,
            activateSessionRequest);

        context.sendRequest(new OpcuaAPU(activateMessageRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
            .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
            .handle(opcuaActivateResponse -> {
                LOGGER.debug("Got Activate Session Response Connection Response");

                ActivateSessionResponse activateMessageResponse = (ActivateSessionResponse) opcuaActivateResponse.getMessage();

                long returnedRequestHandle = activateMessageResponse.getResponseHeader().getRequestHandle();
                if (!(requestHandle == returnedRequestHandle)) {
                    LOGGER.error("Request handle isn't as expected, we might have missed a packet. - " +  requestHandle + " != " + returnedRequestHandle);
                }

                // Send an event that connection setup is complete.
                context.fireConnected();
            });
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        LOGGER.info("Reading Value");
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        if(request.getFieldNames().size() == 1) {
            String fieldName = request.getFieldNames().iterator().next();
            OpcuaField field = (OpcuaField) request.getField(fieldName);

            int requestHandle = requestHandleGenerator.getAndIncrement();
            // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
            if(requestHandleGenerator.get() == 0xFFFF) {
                requestHandleGenerator.set(1);
            }

            RequestHeader requestHeader = new RequestHeader(authenticationToken,
                getCurrentDateTime(),
                requestHandle,
                0L,
                NULL_STRING,
                10000L,
                NULL_EXTENSION_OBJECT);

            ReadValueId[] readValueArray = new ReadValueId[1];

            NodeIdString nodeId = new NodeIdString(NodeIdType.nodeIdTypeString, new StringNodeId(field.getNamespace(), new PascalString(field.getIdentifier().length(), field.getIdentifier())));

            readValueArray[0] = new ReadValueId(nodeId,
                0xD,
                NULL_STRING,
                new QualifiedName(0, NULL_STRING));

            ReadRequest opcuaReadRequest = new ReadRequest((byte) 1,
                (byte) 0,
                requestHeader,
                0.0d,
                TimestampsToReturn.timestampsToReturnNeither,
                readValueArray.length,
                readValueArray);

            int transactionIdentifier = transactionIdentifierGenerator.getAndIncrement();
            // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
            if(transactionIdentifierGenerator.get() == 0xFFFF) {
                transactionIdentifierGenerator.set(1);
            }

            OpcuaMessageRequest readMessageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                transactionIdentifier,
                transactionIdentifier,
                opcuaReadRequest);

            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> context.sendRequest(new OpcuaAPU(readMessageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaResponse -> {
                    // Try to decode the response data based on the corresponding request.
                    ReadResponse readResponse = (ReadResponse) opcuaResponse.getMessage();

                    //TODO;- Fix this
                    DataValue[] results = readResponse.getResults();
                    Integer value = null;
                    if (results.length > 0) {
                        Variant variant = results[0].getValue();
                        LOGGER.info("Repsponse Include Variant of type " + variant.getClass().toString());
                        if (variant instanceof VariantInt32) {
                            value = ((VariantInt32) variant).getValue()[0];
                        }
                    }

                    PlcValue plcValue = new PlcDINT(value);
                    // Prepare the response.
                    PlcReadResponse response = new DefaultPlcReadResponse(request,
                        Collections.singletonMap(fieldName, new ResponseItem<>(PlcResponseCode.OK, plcValue)));

                    // Pass the response back to the application.
                    future.complete(response);

                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        } else {
            future.completeExceptionally(new PlcRuntimeException("Modbus only supports single filed requests"));
        }
        return future;
    }

    @Override
    protected void decode(ConversationContext<OpcuaAPU> context, OpcuaAPU msg) throws Exception {
        super.decode(context, msg);
    }


    private long getCurrentDateTime() {
        return (System.currentTimeMillis() * 10000) + epochOffset;
    }
}
