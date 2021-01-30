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

import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.context.CertificateKeyPair;
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
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.crypto.Cipher;
import java.io.ByteArrayInputStream;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.security.cert.CertificateEncodingException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.time.Duration;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ThreadLocalRandom;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;

public class OpcuaProtocolLogic extends Plc4xProtocolBase<OpcuaAPU> implements HasConfiguration<OpcuaConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(1000000);
    public static final long REQUEST_TIMEOUT_LONG = 10000L;

    private static final int DEFAULT_CONNECTION_LIFETIME = 36000000;
    private static final int DEFAULT_MAX_CHUNK_COUNT = 64;
    private static final int DEFAULT_MAX_REQUEST_ID = 0xFFFFFFFF;
    private static final int DEFAULT_MAX_MESSAGE_SIZE = 2097152;
    private static final int DEFAULT_RECEIVE_BUFFER_SIZE = 65535;
    private static final int DEFAULT_SEND_BUFFER_SIZE = 65535;
    private static final int VERSION = 0;
    private static final String PASSWORD_ENCRYPTION_ALGORITHM = "http://www.w3.org/2001/04/xmlenc#rsa-oaep";

    private NodeId authenticationToken = new NodeIdTwoByte(NodeIdType.nodeIdTypeTwoByte, new TwoByteNodeId((short) 0));
    private static final PascalString NULL_STRING = new PascalString(-1,null);
    private static final PascalByteString NULL_BYTE_STRING = new PascalByteString(-1, new byte[0]);
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

    private static final String CHUNK = "F";

    private static final String applicationUri = "urn:apache:plc4x:client";
    private static final String productUri = "urn:apache:plc4x:client";
    private static final String applicationText = "OPCUA client for the Apache PLC4X:PLC4J project";

    private final String sessionName = "UaSession:" + applicationText + ":" + RandomStringUtils.random(20, true, true);
    private final byte[] clientNonce = RandomUtils.nextBytes(40);
    private RequestTransactionManager tm;

    private PascalString policyId = null;
    private String endpoint;
    private boolean discovery;
    private String username;
    private String password;
    private String certFile;
    private String securityPolicy;
    private String keyStoreFile;
    private CertificateKeyPair ckp;
    private PascalByteString publicCertificate;
    private PascalByteString thumbprint;
    private boolean isEncrypted;
    private boolean checkedEndpoints = false;
    private AtomicInteger transactionIdentifierGenerator = new AtomicInteger(1);
    private AtomicInteger requestHandleGenerator = new AtomicInteger(1);
    private AtomicInteger tokenId = new AtomicInteger(1);
    private AtomicInteger channelId = new AtomicInteger(1);
    private byte[] senderCertificate = null;
    private byte[] senderNonce = null;
    private String certificateThumbprint = null;

    private AtomicBoolean securedConnection = new AtomicBoolean(false);

    @Override
    public void setConfiguration(OpcuaConfiguration configuration) {
        this.endpoint = configuration.getEndpoint();
        this.discovery = configuration.isDiscovery();
        this.username = configuration.getUsername();
        this.password = configuration.getPassword();
        this.certFile = configuration.getCertDirectory();
        this.securityPolicy = "http://opcfoundation.org/UA/SecurityPolicy#" + configuration.getSecurityPolicy();
        LOGGER.info("---------------------------------------");
        LOGGER.info(configuration.getSecurityPolicy());
        this.ckp = configuration.getCertificateKeyPair();
        if (configuration.getSecurityPolicy().equals("Basic256Sha256")) {
            try {
                this.publicCertificate = new PascalByteString(this.ckp.getCertificate().getEncoded().length, this.ckp.getCertificate().getEncoded());
                this.isEncrypted = true;
            } catch (CertificateEncodingException e) {
                LOGGER.error("Failed to encode the certificate");
            }
            //this.thumbprint = new PascalByteString(this.ckp.getThumbPrint().length, this.ckp.getThumbPrint());
            this.thumbprint =NULL_BYTE_STRING;
        } else {
            this.publicCertificate = NULL_BYTE_STRING;
            this.thumbprint = NULL_BYTE_STRING;
            this.isEncrypted = false;
        }
        this.keyStoreFile = configuration.getKeyStoreFile();

        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void close(ConversationContext<OpcuaAPU> context) {
        //Nothing
    }

    @Override
    public void onDisconnect(ConversationContext<OpcuaAPU> context) {
        int transactionId = getTransactionIdentifier(securedConnection.get());

        int requestHandle = getRequestHandle(securedConnection.get());

        ExpandedNodeId expandedNodeId = new ExpandedNodeIdFourByte(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            NULL_STRING,                      //Namespace Uri
            1L,                     //Server Index
            new FourByteNodeId((short) 0, 473));    //Identifier for OpenSecureChannel

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            requestHandle,                                         //RequestHandle
            0L,
            NULL_STRING,
            5000L,
            NULL_EXTENSION_OBJECT);

        CloseSessionRequest closeSessionRequest = new CloseSessionRequest((byte) 1,
            (byte) 0,
            requestHeader,
            true);

        try {
            WriteBuffer buffer = new WriteBuffer(closeSessionRequest.getLengthInBytes(), true);
            OpcuaMessageIO.staticSerialize(buffer, closeSessionRequest);

            OpcuaMessageRequest messageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                transactionId,
                transactionId,
                buffer.getData());

            context.sendRequest(new OpcuaAPU(messageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaMessageResponse -> {
                        LOGGER.info("Got Close Session Response Connection Response" + opcuaMessageResponse.toString());
                        onDisconnectCloseSecureChannel(context);
                    });
        } catch (ParseException e) {
            LOGGER.error("Failed to parse the Message Request");
        }
    }

    private void onDisconnectCloseSecureChannel(ConversationContext<OpcuaAPU> context) {

        int transactionId = getTransactionIdentifier(securedConnection.get());

        ExpandedNodeId expandedNodeId = new ExpandedNodeIdFourByte(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            NULL_STRING,                      //Namespace Uri
            1L,                     //Server Index
            new FourByteNodeId((short) 0, 452));    //Identifier for OpenSecureChannel

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        CloseSecureChannelRequest closeSecureChannelRequest = new CloseSecureChannelRequest((byte) 1,
            (byte) 0,
            requestHeader);

        OpcuaCloseRequest closeRequest = new OpcuaCloseRequest(CHUNK,
            channelId.get(),
            tokenId.get(),
            transactionId,
            transactionId,
            closeSecureChannelRequest);

        context.sendRequest(new OpcuaAPU(closeRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
            .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
            .handle(opcuaMessageResponse -> {
                LOGGER.info("Got Close Secure Channel Response" + opcuaMessageResponse.toString());
            });
        context.fireDisconnected();
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
                if (this.isEncrypted & !this.checkedEndpoints) {
                    onConnectGetEndpointsOpenSecureChannel(context, opcuaAcknowledgeResponse);
                } else {
                    onConnectOpenSecureChannel(context, opcuaAcknowledgeResponse);
                }
            });
    }

    public void onConnectGetEndpointsOpenSecureChannel(ConversationContext<OpcuaAPU> context, OpcuaAcknowledgeResponse opcuaAcknowledgeResponse) {
        int transactionId = getTransactionIdentifier(securedConnection.get());

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
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        OpenSecureChannelRequest openSecureChannelRequest = new OpenSecureChannelRequest((byte) 1,
            (byte) 0,
            requestHeader,
            VERSION,
            SecurityTokenRequestType.securityTokenRequestTypeIssue,
            MessageSecurityMode.messageSecurityModeNone,
            NULL_BYTE_STRING,
            DEFAULT_CONNECTION_LIFETIME);

        OpcuaOpenRequest openRequest = new OpcuaOpenRequest(CHUNK,
            0,
            new PascalString("http://opcfoundation.org/UA/SecurityPolicy#None".length(), "http://opcfoundation.org/UA/SecurityPolicy#None"),
            NULL_BYTE_STRING,
            NULL_BYTE_STRING,
            transactionId,
            transactionId,
            openSecureChannelRequest);

        context.sendRequest(new OpcuaAPU(openRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
            .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
            .handle(opcuaOpenResponse -> {
                if (opcuaOpenResponse.getMessage() instanceof ServiceFault) {
                    ServiceFault fault = (ServiceFault) opcuaOpenResponse.getMessage();
                    LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", fault.getResponseHeader().getServiceResult().getStatusCode(), OpcuaStatusCodes.enumForValue(fault.getResponseHeader().getServiceResult().getStatusCode()));
                } else {
                    LOGGER.debug("Got Secure Response Connection Response");
                    try {
                        onConnectGetEndpointsRequest(context, opcuaOpenResponse);
                    } catch (PlcConnectionException e) {
                        LOGGER.error("Error occurred while connecting to OPC UA server");
                    }
                }
            });
    }

    public void onConnectGetEndpointsRequest(ConversationContext<OpcuaAPU> context, OpcuaOpenResponse opcuaOpenResponse) throws PlcConnectionException {
        certificateThumbprint = opcuaOpenResponse.getReceiverCertificateThumbprint();
        OpenSecureChannelResponse openSecureChannelResponse = (OpenSecureChannelResponse) opcuaOpenResponse.getMessage();
        tokenId.set((int) openSecureChannelResponse.getSecurityToken().getTokenId());
        channelId.set((int) openSecureChannelResponse.getSecurityToken().getChannelId());

        int transactionId = getTransactionIdentifier(securedConnection.get());

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
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        GetEndpointsRequest endpointsRequest = new GetEndpointsRequest((byte) 1,
            (byte) 0,
            requestHeader,
            new PascalString(this.endpoint.length(), this.endpoint),
            0,
            null,
            0,
            null);

        try {
            WriteBuffer buffer = new WriteBuffer(endpointsRequest.getLengthInBytes(), true);
            OpcuaMessageIO.staticSerialize(buffer, endpointsRequest);

            OpcuaMessageRequest messageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                nextSequenceNumber,
                nextRequestId,
                buffer.getData());

            context.sendRequest(new OpcuaAPU(messageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaMessageResponse -> {
                    try {
                        OpcuaMessage message = OpcuaMessageIO.staticParse(new ReadBuffer(opcuaMessageResponse.getMessage(), true));
                        if (message instanceof ServiceFault) {
                            ServiceFault fault = (ServiceFault) message;
                            LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", fault.getResponseHeader().getServiceResult().getStatusCode(), OpcuaStatusCodes.enumForValue(fault.getResponseHeader().getServiceResult().getStatusCode()));
                        } else {
                            LOGGER.debug("Got Create Session Response Connection Response");
                            onConnectGetEndpointsCloseSecureChannel(context, (GetEndpointsResponse) message);
                        }
                    } catch (ParseException e) {
                        e.printStackTrace();
                    }

                });
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Create Session Request");
        }
    }

    private void onConnectGetEndpointsCloseSecureChannel(ConversationContext<OpcuaAPU> context, GetEndpointsResponse message) {

        int transactionId = getTransactionIdentifier(securedConnection.get());
        EndpointDescription[] endpoints = message.getEndpoints();
        for (EndpointDescription endpoint : endpoints) {
            if (endpoint.getEndpointUrl().getStringValue().equals(this.endpoint) && endpoint.getSecurityPolicyUri().getStringValue().equals(this.securityPolicy)) {
                LOGGER.info("Found OPC UA endpoint {}", this.endpoint);
                this.senderCertificate = endpoint.getServerCertificate().getStringValue();
            }
        }

        ExpandedNodeId expandedNodeId = new ExpandedNodeIdFourByte(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            NULL_STRING,                      //Namespace Uri
            1L,                     //Server Index
            new FourByteNodeId((short) 0, 452));    //Identifier for CloseSecureChannel

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        CloseSecureChannelRequest closeSecureChannelRequest = new CloseSecureChannelRequest((byte) 1,
            (byte) 0,
            requestHeader);

        OpcuaCloseRequest closeRequest = new OpcuaCloseRequest(CHUNK,
            channelId.get(),
            tokenId.get(),
            transactionId,
            transactionId,
            closeSecureChannelRequest);

        context.sendRequest(new OpcuaAPU(closeRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
            .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
            .handle(opcuaMessageResponse -> {
                LOGGER.info("Got Close Secure Channel Response" + opcuaMessageResponse.toString());
                this.checkedEndpoints = true;
                channelId.set(0);
                tokenId.set(0);
                onConnect(context);
            });
    }

    public void onConnectOpenSecureChannel(ConversationContext<OpcuaAPU> context, OpcuaAcknowledgeResponse opcuaAcknowledgeResponse) {

        int transactionId = getTransactionIdentifier(securedConnection.get());

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
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        OpenSecureChannelRequest openSecureChannelRequest = new OpenSecureChannelRequest((byte) 1,
            (byte) 0,
            requestHeader,
            VERSION,
            SecurityTokenRequestType.securityTokenRequestTypeIssue,
            MessageSecurityMode.messageSecurityModeNone,
            NULL_BYTE_STRING,
            DEFAULT_CONNECTION_LIFETIME);

        OpcuaOpenRequest openRequest = new OpcuaOpenRequest(CHUNK,
            0,
            new PascalString(this.securityPolicy.length(), this.securityPolicy),
            this.publicCertificate,
            this.thumbprint,
            transactionId,
            transactionId,
            openSecureChannelRequest);

        context.sendRequest(new OpcuaAPU(openRequest))
            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
            .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
            .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
            .handle(opcuaOpenResponse -> {
                if (opcuaOpenResponse.getMessage() instanceof ServiceFault) {
                    ServiceFault fault = (ServiceFault) opcuaOpenResponse.getMessage();
                    LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", fault.getResponseHeader().getServiceResult().getStatusCode(), OpcuaStatusCodes.enumForValue(fault.getResponseHeader().getServiceResult().getStatusCode()));
                } else {
                    LOGGER.debug("Got Secure Response Connection Response");
                    try {
                        onConnectCreateSessionRequest(context, opcuaOpenResponse);
                    } catch (PlcConnectionException e) {
                        LOGGER.error("Error occurred while connecting to OPC UA server");
                    }
                }
            });

    }

    public void onConnectCreateSessionRequest(ConversationContext<OpcuaAPU> context, OpcuaOpenResponse opcuaOpenResponse) throws PlcConnectionException {

        certificateThumbprint = opcuaOpenResponse.getReceiverCertificateThumbprint();
        OpenSecureChannelResponse openSecureChannelResponse = (OpenSecureChannelResponse) opcuaOpenResponse.getMessage();
        tokenId.set((int) openSecureChannelResponse.getSecurityToken().getTokenId());
        channelId.set((int) openSecureChannelResponse.getSecurityToken().getChannelId());


        int transactionId = getTransactionIdentifier(securedConnection.get());

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
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        LocalizedText applicationName = new LocalizedText((short) 0,
            true,
            true,
            new PascalString("en".length(), "en"),
            new PascalString(applicationText.length(), applicationText));

        PascalString gatewayServerUri = NULL_STRING;
        PascalString discoveryProfileUri = NULL_STRING;
        int noOfDiscoveryUrls = -1;
        PascalString[] discoveryUrls = new PascalString[0];

        ApplicationDescription clientDescription = new ApplicationDescription(new PascalString(applicationUri.length(), applicationUri),
            new PascalString(productUri.length(), productUri),
            applicationName,
            ApplicationType.applicationTypeClient,
            gatewayServerUri,
            discoveryProfileUri,
            noOfDiscoveryUrls,
            discoveryUrls);

        CreateSessionRequest createSessionRequest = new CreateSessionRequest((byte) 1,
            (byte) 0,
            requestHeader,
            clientDescription,
            NULL_STRING,
            new PascalString(endpoint.length(), endpoint),
            new PascalString(sessionName.length(), sessionName),
            new PascalByteString(clientNonce.length, clientNonce),
            NULL_BYTE_STRING,
            120000L,
            0L);

        try {
            WriteBuffer buffer = new WriteBuffer(createSessionRequest.getLengthInBytes(), true);
            OpcuaMessageIO.staticSerialize(buffer, createSessionRequest);

            OpcuaMessageRequest messageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                nextSequenceNumber,
                nextRequestId,
                buffer.getData());

            context.sendRequest(new OpcuaAPU(messageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaMessageResponse -> {
                    try {
                        OpcuaMessage message = OpcuaMessageIO.staticParse(new ReadBuffer(opcuaMessageResponse.getMessage(), true));
                        if (message instanceof ServiceFault) {
                            ServiceFault fault = (ServiceFault) message;
                            LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", fault.getResponseHeader().getServiceResult().getStatusCode(), OpcuaStatusCodes.enumForValue(fault.getResponseHeader().getServiceResult().getStatusCode()));
                        } else {
                            LOGGER.debug("Got Create Session Response Connection Response");
                            try {
                                onConnectActivateSessionRequest(context, opcuaMessageResponse, (CreateSessionResponse) message);
                            } catch (PlcConnectionException e) {
                                LOGGER.error("Error occurred while connecting to OPC UA server");
                            }
                        }
                    } catch (ParseException e) {
                        e.printStackTrace();
                    }

                });
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Create Session Request");
        }
    }

    private void onConnectActivateSessionRequest(ConversationContext<OpcuaAPU> context, OpcuaMessageResponse opcuaMessageResponse, CreateSessionResponse sessionResponse) throws PlcConnectionException {

        CreateSessionResponse createSessionResponse = sessionResponse;
        senderCertificate = createSessionResponse.getServerCertificate().getStringValue();
        senderNonce = createSessionResponse.getServerNonce().getStringValue();

        for (EndpointDescription endpointDescription: createSessionResponse.getServerEndpoints()) {
            LOGGER.info("{} - {}", endpointDescription.getEndpointUrl().getStringValue(), this.endpoint);
            if (endpointDescription.getEndpointUrl().getStringValue().equals(this.endpoint)) {
                for (UserTokenPolicy identityToken : endpointDescription.getUserIdentityTokens()) {
                    if (identityToken.getTokenType() == UserTokenType.userTokenTypeAnonymous) {
                        if (this.username == null) {
                            policyId = identityToken.getPolicyId();
                        }
                    } else if (identityToken.getTokenType() == UserTokenType.userTokenTypeUserName) {
                        if (this.username != null) {
                            policyId = identityToken.getPolicyId();
                        }
                    }
                }
            }
        }
        LOGGER.info(policyId.getStringValue());

        authenticationToken = createSessionResponse.getAuthenticationToken();
        tokenId.set((int) opcuaMessageResponse.getSecureTokenId());
        channelId.set((int) opcuaMessageResponse.getSecureChannelId());

        int transactionId = getTransactionIdentifier(securedConnection.get());

        Integer nextSequenceNumber = opcuaMessageResponse.getSequenceNumber() + 1;
        Integer nextRequestId = opcuaMessageResponse.getRequestId() + 1;

        if (!(transactionId == nextSequenceNumber)) {
            LOGGER.error("Sequence number isn't as expected, we might have missed a packet. - {} != {}" , transactionId, nextSequenceNumber);
            throw new PlcConnectionException("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
        }

        int requestHandle = getRequestHandle(securedConnection.get());

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            requestHandle,
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        SignatureData clientSignature = new SignatureData(NULL_STRING, NULL_BYTE_STRING);

        SignedSoftwareCertificate[] signedSoftwareCertificate = new SignedSoftwareCertificate[1];

        signedSoftwareCertificate[0] = new SignedSoftwareCertificate(NULL_BYTE_STRING, NULL_BYTE_STRING);


        ExtensionObject userIdentityToken = null;
        if (this.username == null) {
            userIdentityToken = getIdentityToken("none");
        } else {
            userIdentityToken = getIdentityToken("username");
        }

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

        try {
            WriteBuffer buffer = new WriteBuffer(activateSessionRequest.getLengthInBytes(), true);
            OpcuaMessageIO.staticSerialize(buffer, activateSessionRequest);

            OpcuaMessageRequest activateMessageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                nextSequenceNumber,
                nextRequestId,
                buffer.getData());

            context.sendRequest(new OpcuaAPU(activateMessageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaActivateResponse -> {
                    LOGGER.debug("Got Activate Session Response Connection Response");
                    try {
                        OpcuaMessage message = OpcuaMessageIO.staticParse(new ReadBuffer(opcuaActivateResponse.getMessage(), true));
                        if (message instanceof ServiceFault) {
                            ServiceFault fault = (ServiceFault) message;
                            LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", fault.getResponseHeader().getServiceResult().getStatusCode(), OpcuaStatusCodes.enumForValue(fault.getResponseHeader().getServiceResult().getStatusCode()));
                        } else {
                            ActivateSessionResponse activateMessageResponse = (ActivateSessionResponse) message;

                            long returnedRequestHandle = activateMessageResponse.getResponseHeader().getRequestHandle();
                            if (!(requestHandle == returnedRequestHandle)) {
                                LOGGER.error("Request handle isn't as expected, we might have missed a packet. {} != {}", requestHandle, returnedRequestHandle);
                            }

                            // Send an event that connection setup is complete.
                            context.fireConnected();
                        }
                    } catch (ParseException e) {
                        e.printStackTrace();
                    }

                });
        } catch (ParseException e) {
            LOGGER.info("Unable to serialise the ActivateSessionRequest");
        }
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        LOGGER.info("Reading Value");
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;


        int requestHandle = getRequestHandle(securedConnection.get());

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            requestHandle,
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        ReadValueId[] readValueArray = new ReadValueId[request.getFieldNames().size()];
        Iterator<String> iterator = request.getFieldNames().iterator();
        for (int i = 0; i < request.getFieldNames().size(); i++ ) {
            String fieldName = iterator.next();
            OpcuaField field = (OpcuaField) request.getField(fieldName);

            NodeId nodeId = null;
            if (field.getIdentifierType() == OpcuaIdentifierType.BINARY_IDENTIFIER) {
                nodeId = new NodeIdTwoByte(NodeIdType.nodeIdTypeTwoByte, new TwoByteNodeId(Short.valueOf(field.getIdentifier())));
            } else if (field.getIdentifierType() == OpcuaIdentifierType.NUMBER_IDENTIFIER) {
                nodeId = new NodeIdNumeric(NodeIdType.nodeIdTypeNumeric, new NumericNodeId(field.getNamespace(),Long.valueOf(field.getIdentifier())));
            } else if (field.getIdentifierType() == OpcuaIdentifierType.GUID_IDENTIFIER) {
                nodeId = new NodeIdGuid(NodeIdType.nodeIdTypeGuid, new GuidNodeId(field.getNamespace(), toGuidValue(field.getIdentifier())));
            } else if (field.getIdentifierType() == OpcuaIdentifierType.STRING_IDENTIFIER) {
                nodeId = new NodeIdString(NodeIdType.nodeIdTypeString, new StringNodeId(field.getNamespace(), new PascalString(field.getIdentifier().length(), field.getIdentifier())));
            }
            readValueArray[i] = new ReadValueId(nodeId,
                0xD,
                NULL_STRING,
                new QualifiedName(0, NULL_STRING));
        }

        ReadRequest opcuaReadRequest = new ReadRequest((byte) 1,
            (byte) 0,
            requestHeader,
            0.0d,
            TimestampsToReturn.timestampsToReturnNeither,
            readValueArray.length,
            readValueArray);

        int transactionId = getTransactionIdentifier(securedConnection.get());

        try {
            WriteBuffer buffer = new WriteBuffer(opcuaReadRequest.getLengthInBytes(), true);
            OpcuaMessageIO.staticSerialize(buffer, opcuaReadRequest);

            OpcuaMessageRequest readMessageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                transactionId,
                transactionId,
                buffer.getData());

            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> context.sendRequest(new OpcuaAPU(readMessageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaResponse -> {
                    // Prepare the response.
                    PlcReadResponse response = null;
                    try {
                        response = new DefaultPlcReadResponse(request, readResponse(request.getFieldNames(), (ReadResponse) OpcuaMessageIO.staticParse(new ReadBuffer(opcuaResponse.getMessage(), true))));
                    } catch (ParseException e) {
                        e.printStackTrace();
                    };

                    // Pass the response back to the application.
                    future.complete(response);

                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        } catch (ParseException e) {
            LOGGER.error("Unable to serialise the ReadRequest");
        }

        return future;
    }

    private Map<String, ResponseItem<PlcValue>> readResponse(LinkedHashSet<String> fieldNames, ReadResponse readResponse) {
        DataValue[] results = readResponse.getResults();

        PlcResponseCode responseCode = PlcResponseCode.OK;
        Map<String, ResponseItem<PlcValue>> response = new HashMap<>();
        int count = 0;
        for ( String field : fieldNames ) {
            PlcValue value = null;
            if (results[count].getValueSpecified()) {
                Variant variant = results[count].getValue();
                LOGGER.info("Response of type {}", variant.getClass().toString());
                if (variant instanceof VariantBoolean) {
                    byte[] array = ((VariantBoolean) variant).getValue();
                    int length = array.length;
                    Byte[] tmpValue = new Byte[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantSByte) {
                    byte[] array = ((VariantSByte) variant).getValue();
                    int length = array.length;
                    Byte[] tmpValue = new Byte[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantByte) {
                    short[] array = ((VariantByte) variant).getValue();
                    int length = array.length;
                    Short[] tmpValue = new Short[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantInt16) {
                    short[] array = ((VariantInt16) variant).getValue();
                    int length = array.length;
                    Short[] tmpValue = new Short[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantUInt16) {
                    int[] array = ((VariantUInt16) variant).getValue();
                    int length = array.length;
                    Integer[] tmpValue = new Integer[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantInt32) {
                    int[] array = ((VariantInt32) variant).getValue();
                    int length = array.length;
                    Integer[] tmpValue = new Integer[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantUInt32) {
                    long[] array = ((VariantUInt32) variant).getValue();
                    int length = array.length;
                    Long[] tmpValue = new Long[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantInt64) {
                    long[] array = ((VariantInt64) variant).getValue();
                    int length = array.length;
                    Long[] tmpValue = new Long[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantUInt64) {
                    value = IEC61131ValueHandler.of(((VariantUInt64) variant).getValue());
                } else if (variant instanceof VariantFloat) {
                    float[] array = ((VariantFloat) variant).getValue();
                    int length = array.length;
                    Float[] tmpValue = new Float[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantDouble) {
                    double[] array = ((VariantDouble) variant).getValue();
                    int length = array.length;
                    Double[] tmpValue = new Double[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i];
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantString) {
                    int length = ((VariantString) variant).getValue().length;
                    PascalString[] stringArray = ((VariantString) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = stringArray[i].getStringValue();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantDateTime) {
                    long[] array = ((VariantDateTime) variant).getValue();
                    int length = array.length;
                    LocalDateTime[] tmpValue = new LocalDateTime[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = LocalDateTime.ofInstant(Instant.ofEpochMilli(getDateTime(array[i])), ZoneOffset.UTC);
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantGuid) {
                    GuidValue[] array = ((VariantGuid) variant).getValue();
                    int length = array.length;
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        //These two data section aren't little endian like the rest.
                        byte[] data4Bytes = array[i].getData4();
                        int data4 = 0;
                        for (int k = 0; k < data4Bytes.length; k++)
                        {
                            data4 = (data4 << 8) + (data4Bytes[k] & 0xff);
                        }
                        byte[] data5Bytes = array[i].getData5();
                        long data5 = 0;
                        for (int k = 0; k < data5Bytes.length; k++)
                        {
                            data5 = (data5 << 8) + (data5Bytes[k] & 0xff);
                        }
                        tmpValue[i] = Long.toHexString(array[i].getData1()) + "-" + Integer.toHexString(array[i].getData2()) + "-" + Integer.toHexString(array[i].getData3()) + "-" + Integer.toHexString(data4) + "-" + Long.toHexString(data5);
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantXmlElement) {
                    int length = ((VariantXmlElement) variant).getValue().length;
                    PascalString[] stringArray = ((VariantXmlElement) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = stringArray[i].getStringValue();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantLocalizedText) {
                    int length = ((VariantLocalizedText) variant).getValue().length;
                    LocalizedText[] stringArray = ((VariantLocalizedText) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = "";
                        tmpValue[i] += stringArray[i].getLocaleSpecified() ? stringArray[i].getLocale().getStringValue() + "|" : "";
                        tmpValue[i] += stringArray[i].getTextSpecified() ? stringArray[i].getText().getStringValue() : "";
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantQualifiedName) {
                    int length = ((VariantQualifiedName) variant).getValue().length;
                    QualifiedName[] stringArray = ((VariantQualifiedName) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = "ns=" + stringArray[i].getNamespaceIndex() + ";s=" + stringArray[i].getName().getStringValue();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantExtensionObject) {
                    int length = ((VariantExtensionObject) variant).getValue().length;
                    ExtensionObject[] stringArray = ((VariantExtensionObject) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = stringArray[i].toString();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantNodeId) {
                    int length = ((VariantNodeId) variant).getValue().length;
                    NodeId[] stringArray = ((VariantNodeId) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = stringArray[i].toString();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                }else if (variant instanceof VariantStatusCode) {
                    int length = ((VariantStatusCode) variant).getValue().length;
                    StatusCode[] stringArray = ((VariantStatusCode) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = stringArray[i].toString();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantByteString) {
                    PlcList plcList = new PlcList();
                    ByteStringArray[] array = ((VariantByteString) variant).getValue();
                    for (int k = 0; k < array.length; k++) {
                        int length = array[k].getValue().length;
                        Short[] tmpValue = new Short[length];
                        for (int i = 0; i < length; i++) {
                            tmpValue[i] = array[k].getValue()[i];
                        }
                        plcList.add(IEC61131ValueHandler.of(tmpValue));
                    }
                    value = plcList;
                } else {
                    responseCode = PlcResponseCode.UNSUPPORTED;
                    LOGGER.error("Data type - " +  variant.getClass() + " is not supported ");
                }
            } else {
                if (results[count].getStatusCode().getStatusCode() == OpcuaStatusCodes.BadNodeIdUnknown.getValue()) {
                    responseCode = PlcResponseCode.NOT_FOUND;
                } else {
                    responseCode = PlcResponseCode.UNSUPPORTED;
                }
                LOGGER.error("Error while reading value from OPC UA server error code:- " + results[count].getStatusCode().toString());
            }
            count++;
            response.put(field, new ResponseItem<>(responseCode, value));
        }
        return response;
    }

    private Variant fromPlcValue(String fieldName, OpcuaField field, PlcWriteRequest request) {

        PlcList valueObject;
        if (request.getPlcValue(fieldName).getObject() instanceof ArrayList) {
            valueObject = (PlcList) request.getPlcValue(fieldName);
        } else {
            ArrayList<PlcValue> list = new ArrayList<>();
            list.add(request.getPlcValue(fieldName));
            valueObject = new PlcList(list);
        }

        List<PlcValue> plcValueList = valueObject.getList();
        String dataType = field.getPlcDataType();
        if (dataType.equals("IEC61131_NULL")) {
            if (plcValueList.get(0).getObject() instanceof Boolean) {
                dataType = "IEC61131_BOOL";
            } else if (plcValueList.get(0).getObject() instanceof Byte) {
                dataType = "IEC61131_SINT";
            } else if (plcValueList.get(0).getObject() instanceof Short) {
                dataType = "IEC61131_INT";
            } else if (plcValueList.get(0).getObject() instanceof Integer) {
                dataType = "IEC61131_DINT";
            } else if (plcValueList.get(0).getObject() instanceof Long) {
                dataType = "IEC61131_LINT";
            } else if (plcValueList.get(0).getObject() instanceof Float) {
                dataType = "IEC61131_REAL";
            } else if (plcValueList.get(0).getObject() instanceof Double) {
                dataType = "IEC61131_LREAL";
            } else if (plcValueList.get(0).getObject() instanceof String) {
                dataType = "IEC61131_STRING";
            }
        }
        int length = valueObject.getLength();
        switch (dataType) {
            case "IEC61131_BOOL":
            case "IEC61131_BIT":
                byte[] tmpBOOL = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpBOOL[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantBoolean(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpBOOL);
            case "IEC61131_BYTE":
            case "IEC61131_BITARR8":
            case "IEC61131_USINT":
            case "IEC61131_UINT8":
            case "IEC61131_BIT8":
                short[] tmpBYTE = new short[length];
                for (int i = 0; i < length; i++) {
                    tmpBYTE[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantByte(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpBYTE);
            case "IEC61131_SINT":
            case "IEC61131_INT8":
                byte[] tmpSINT = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpSINT[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantSByte(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpSINT);
            case "IEC61131_INT":
            case "IEC61131_INT16":
                short[] tmpINT16 = new short[length];
                for (int i = 0; i < length; i++) {
                    tmpINT16[i] = valueObject.getIndex(i).getShort();
                }
                return new VariantInt16(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpINT16);
            case "IEC61131_UINT":
            case "IEC61131_UINT16":
            case "IEC61131_WORD":
            case "IEC61131_BITARR16":
                int[] tmpUINT = new int[length];
                for (int i = 0; i < length; i++) {
                    tmpUINT[i] = valueObject.getIndex(i).getInt();
                }
                return new VariantUInt16(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpUINT);
            case "IEC61131_DINT":
            case "IEC61131_INT32":
                int[] tmpDINT = new int[length];
                for (int i = 0; i < length; i++) {
                    tmpDINT[i] = valueObject.getIndex(i).getInt();
                }
                return new VariantInt32(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpDINT);
            case "IEC61131_UDINT":
            case "IEC61131_UINT32":
            case "IEC61131_DWORD":
            case "IEC61131_BITARR32":
                long[] tmpUDINT = new long[length];
                for (int i = 0; i < length; i++) {
                    tmpUDINT[i] = valueObject.getIndex(i).getLong();
                }
                return new VariantUInt32(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpUDINT);
            case "IEC61131_LINT":
            case "IEC61131_INT64":
                long[] tmpLINT = new long[length];
                for (int i = 0; i < length; i++) {
                    tmpLINT[i] = valueObject.getIndex(i).getLong();
                }
                return new VariantInt64(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpLINT);
            case "IEC61131_ULINT":
            case "IEC61131_UINT64":
            case "IEC61131_LWORD":
            case "IEC61131_BITARR64":
                BigInteger[] tmpULINT = new BigInteger[length];
                for (int i = 0; i < length; i++) {
                    tmpULINT[i] = valueObject.getIndex(i).getBigInteger();
                }
                return new VariantUInt64(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpULINT);
            case "IEC61131_REAL":
            case "IEC61131_FLOAT":
                float[] tmpREAL = new float[length];
                for (int i = 0; i < length; i++) {
                    tmpREAL[i] = valueObject.getIndex(i).getFloat();
                }
                return new VariantFloat(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpREAL);
            case "IEC61131_LREAL":
            case "IEC61131_DOUBLE":
                double[] tmpLREAL = new double[length];
                for (int i = 0; i < length; i++) {
                    tmpLREAL[i] = valueObject.getIndex(i).getDouble();
                }
                return new VariantDouble(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpLREAL);
            case "IEC61131_CHAR":
            case "IEC61131_WCHAR":
            case "IEC61131_STRING":
            case "IEC61131_WSTRING":
            case "IEC61131_STRING16":
                PascalString[] tmpString = new PascalString[length];
                for (int i = 0; i < length; i++) {
                    String s = valueObject.getIndex(i).getString();
                    tmpString[i] = new PascalString(s.length(), s);
                }
                return new VariantString(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpString);
            case "IEC61131_DATE_AND_TIME":
                long[] tmpDateTime = new long[length];
                for (int i = 0; i < length; i++) {
                    tmpDateTime[i] = valueObject.getIndex(i).getDateTime().toEpochSecond(ZoneOffset.UTC);
                }
                return new VariantDateTime(length == 1 ? false : true,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpDateTime);
            default:
                throw new PlcRuntimeException("Unsupported write field type " + dataType);
        }

    }


    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        LOGGER.info("Writing Value");
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        int requestHandle = getRequestHandle(securedConnection.get());

        RequestHeader requestHeader = new RequestHeader(authenticationToken,
            getCurrentDateTime(),
            requestHandle,
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        WriteValue[] writeValueArray = new WriteValue[request.getFieldNames().size()];
        Iterator<String> iterator = request.getFieldNames().iterator();
        for (int i = 0; i < request.getFieldNames().size(); i++ ) {
            String fieldName = iterator.next();
            OpcuaField field = (OpcuaField) request.getField(fieldName);

            NodeId nodeId = null;
            if (field.getIdentifierType() == OpcuaIdentifierType.BINARY_IDENTIFIER) {
                nodeId = new NodeIdTwoByte(NodeIdType.nodeIdTypeTwoByte, new TwoByteNodeId(Short.valueOf(field.getIdentifier())));
            } else if (field.getIdentifierType() == OpcuaIdentifierType.NUMBER_IDENTIFIER) {
                nodeId = new NodeIdNumeric(NodeIdType.nodeIdTypeNumeric, new NumericNodeId(field.getNamespace(),Long.valueOf(field.getIdentifier())));
            } else if (field.getIdentifierType() == OpcuaIdentifierType.GUID_IDENTIFIER) {
                nodeId = new NodeIdGuid(NodeIdType.nodeIdTypeGuid, new GuidNodeId(field.getNamespace(), toGuidValue(field.getIdentifier())));
            } else if (field.getIdentifierType() == OpcuaIdentifierType.STRING_IDENTIFIER) {
                nodeId = new NodeIdString(NodeIdType.nodeIdTypeString, new StringNodeId(field.getNamespace(), new PascalString(field.getIdentifier().length(), field.getIdentifier())));
            }
            writeValueArray[i] = new WriteValue(nodeId,
                0xD,
                NULL_STRING,
                new DataValue(false,
                    false,
                    false,
                    false,
                    false,
                    false,
                    false,
                    true,
                    fromPlcValue(fieldName, field, writeRequest),
                    null,
                    null,
                    null,
                    null,
                    null));
        }

        WriteRequest opcuaWriteRequest = new WriteRequest((byte) 1,
            (byte) 0,
            requestHeader,
            writeValueArray.length,
            writeValueArray);

        int transactionId = getTransactionIdentifier(securedConnection.get());

        try {
            WriteBuffer buffer = new WriteBuffer(opcuaWriteRequest.getLengthInBytes(), true);
            OpcuaMessageIO.staticSerialize(buffer, opcuaWriteRequest);

            OpcuaMessageRequest writeMessageRequest = new OpcuaMessageRequest(CHUNK,
                channelId.get(),
                tokenId.get(),
                transactionId,
                transactionId,
                buffer.getData());

            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> context.sendRequest(new OpcuaAPU(writeMessageRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .handle(opcuaResponse -> {
                    WriteResponse responseMessage = null;
                    try {
                        responseMessage = (WriteResponse) OpcuaMessageIO.staticParse(new ReadBuffer(opcuaResponse.getMessage(), true));
                    } catch (ParseException e) {
                        e.printStackTrace();
                    }
                    PlcWriteResponse response = writeResponse(request, responseMessage);

                    // Pass the response back to the application.
                    future.complete(response);

                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        } catch (ParseException e) {
            LOGGER.info("Unable to serialize write request");
        }

        return future;
    }

    private PlcWriteResponse writeResponse(DefaultPlcWriteRequest request, WriteResponse writeResponse) {
        Map<String, PlcResponseCode> responseMap = new HashMap<>();

        StatusCode[] results = writeResponse.getResults();
        Iterator<String> responseIterator = request.getFieldNames().iterator();
        for (int i = 0; i < request.getFieldNames().size(); i++ ) {
            String fieldName = responseIterator.next();
            OpcuaStatusCodes statusCode = OpcuaStatusCodes.enumForValue(results[i].getStatusCode());
            switch (statusCode) {
                case Good:
                    responseMap.put(fieldName, PlcResponseCode.OK);
                    break;
                case BadNodeIdUnknown:
                    responseMap.put(fieldName, PlcResponseCode.NOT_FOUND);
                    break;
                default:
                    responseMap.put(fieldName, PlcResponseCode.REMOTE_ERROR);
            }
        }

        return new DefaultPlcWriteResponse(request, responseMap);
    }

    /**
     * Returns the next transaction identifier, for a secured connection we should return a random number.
     * @param random - Select if we return a random number or the next value
     * @return
     */
    private int getTransactionIdentifier(boolean random) {
        if (random) {
            return ThreadLocalRandom.current().nextInt();
        } else {
            int transactionId = transactionIdentifierGenerator.getAndIncrement();
            if(transactionIdentifierGenerator.get() == DEFAULT_MAX_REQUEST_ID) {
                transactionIdentifierGenerator.set(1);
            }
            return transactionId;
        }
    }

    /**
     * Returns the next request handle, for a secured connection we should return a random number.
     * @param random - Select if we return a random number or the next value
     * @return
     */
    private int getRequestHandle(boolean random) {
        if (random) {
            return ThreadLocalRandom.current().nextInt();
        } else {
            int transactionId = requestHandleGenerator.getAndIncrement();
            if(requestHandleGenerator.get() == DEFAULT_MAX_REQUEST_ID) {
                requestHandleGenerator.set(1);
            }
            return transactionId;
        }
    }

    private long getCurrentDateTime() {
        return (System.currentTimeMillis() * 10000) + epochOffset;
    }

    private long getDateTime(long dateTime) {
        return (dateTime - epochOffset) / 10000;
    }

    /**
     * Creates an IdentityToken to authenticate with a server.
     * @param securityPolicy
     * @return returns an ExtensionObject with an IdentityToken.
     */
    private ExtensionObject getIdentityToken(String securityPolicy) {
        ExpandedNodeId extExpandedNodeId = null;
        ExtensionObject userIdentityToken = null;
        switch (securityPolicy) {
            case "none":
                //If we aren't using authentication tell the server we would like to login anonymously
                PascalString anonymousIdentityToken = this.policyId;

                WriteBuffer buffer = new WriteBuffer(anonymousIdentityToken.getLengthInBytes(), true);
                try{
                    PascalStringIO.staticSerialize(buffer, anonymousIdentityToken);
                } catch (ParseException e) {
                    LOGGER.error("Failed to serialize the user identity token - {}", anonymousIdentityToken.getStringValue());
                }
                extExpandedNodeId = new ExpandedNodeIdFourByte(false,
                    false,
                    null,
                    null,
                    new FourByteNodeId((short) 0,  OpcuaNodeIdServices.AnonymousIdentityToken_Encoding_DefaultBinary.getValue()));
                return new ExtensionObject(extExpandedNodeId, (short) 1, buffer.getData().length, buffer.getData());
            case "username":
                //Encrypt the password using the server nonce and server public key
                byte[] passwordBytes = this.password.getBytes();
                ByteBuffer encodeableBuffer = ByteBuffer.allocate(4 + passwordBytes.length + this.senderNonce.length);
                encodeableBuffer.order(ByteOrder.LITTLE_ENDIAN);
                encodeableBuffer.putInt(passwordBytes.length + this.senderNonce.length);
                encodeableBuffer.put(passwordBytes);
                encodeableBuffer.put(this.senderNonce);
                byte[] encodeablePassword = new byte[4 + passwordBytes.length + this.senderNonce.length];
                encodeableBuffer.position(0);
                encodeableBuffer.get(encodeablePassword);

                byte[] encryptedPassword = encrypt(encodeablePassword, senderCertificate);
                UserNameIdentityToken userNameIdentityToken =  new UserNameIdentityToken(
                    new PascalString("username".length(), "username"),
                    new PascalString(this.username.length(), this.username),
                    new PascalByteString(encryptedPassword.length, encryptedPassword),
                    new PascalString(PASSWORD_ENCRYPTION_ALGORITHM.length(), PASSWORD_ENCRYPTION_ALGORITHM)
                );
                WriteBuffer bufferUserName = new WriteBuffer(userNameIdentityToken.getLengthInBytes(), true);
                try{
                    UserNameIdentityTokenIO.staticSerialize(bufferUserName, userNameIdentityToken);
                } catch (ParseException e) {
                    LOGGER.error("Failed to serialize the user identity token - {}", userNameIdentityToken);
                }
                extExpandedNodeId = new ExpandedNodeIdFourByte(false,
                    false,
                    null,
                    null,
                    new FourByteNodeId((short) 0,  OpcuaNodeIdServices.UserNameIdentityToken_Encoding_DefaultBinary.getValue()));
                return new ExtensionObject(extExpandedNodeId, (short) 1, bufferUserName.getData().length, bufferUserName.getData());
        }
        return null;
    }

    public byte[] encrypt(byte[] data, byte[] publicKey) {
        try {
            Cipher cipher = Cipher.getInstance("RSA/ECB/OAEPWithSHA-1AndMGF1Padding");
            cipher.init(Cipher.ENCRYPT_MODE, getCertificateX509().getPublicKey());
            return cipher.doFinal(data);
        } catch (Exception e) {
            LOGGER.error("Unable to encrypt Data");
            return null;
        }
    }

    public X509Certificate getCertificateX509() {
        try {
            CertificateFactory factory =  CertificateFactory.getInstance("X.509");
            return (X509Certificate) factory.generateCertificate(new ByteArrayInputStream(this.senderCertificate));
        } catch (Exception e) {
            LOGGER.error("Unable to get certificate from String {}", this.senderCertificate);
            return null;
        }
    }

    private GuidValue toGuidValue(String identifier) {
        LOGGER.error("Querying Guid nodes is not supported");
        byte[] data4 = new byte[] {0,0};
        byte[] data5 = new byte[] {0,0,0,0,0,0};
        return new GuidValue(0L,0,0,data4, data5);

    }
}
