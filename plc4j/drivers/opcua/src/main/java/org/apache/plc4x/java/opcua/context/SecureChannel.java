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
package org.apache.plc4x.java.opcua.context;

import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.readwrite.io.ExtensionObjectIO;
import org.apache.plc4x.java.opcua.readwrite.io.OpcuaAPUIO;
import org.apache.plc4x.java.opcua.readwrite.types.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateEncodingException;
import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class SecureChannel {

    private static final Logger LOGGER = LoggerFactory.getLogger(SecureChannel.class);
    private static final String FINAL_CHUNK = "F";
    private static final String CONTINUATION_CHUNK = "C";
    private static final String ABORT_CHUNK = "A";
    private static final int VERSION = 0;
    private static final int DEFAULT_MAX_CHUNK_COUNT = 64;
    private static final int DEFAULT_MAX_MESSAGE_SIZE = 2097152;
    private static final int DEFAULT_RECEIVE_BUFFER_SIZE = 65535;
    private static final int DEFAULT_SEND_BUFFER_SIZE = 65535;
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(1000000);
    public static final long REQUEST_TIMEOUT_LONG = 10000L;
    private static final String PASSWORD_ENCRYPTION_ALGORITHM = "http://www.w3.org/2001/04/xmlenc#rsa-oaep";
    private static final PascalString SECURITY_POLICY_NONE = new PascalString("http://opcfoundation.org/UA/SecurityPolicy#None");
    protected static final PascalString NULL_STRING = new PascalString( "");
    private static final PascalByteString NULL_BYTE_STRING = new PascalByteString( -1, null);
    private static ExpandedNodeId NULL_EXPANDED_NODEID = new ExpandedNodeId(false,
        false,
        new NodeIdTwoByte((short) 0),
        null,
        null
    );

    protected static final ExtensionObject NULL_EXTENSION_OBJECT = new ExtensionObject(
        NULL_EXPANDED_NODEID,
        new ExtensionObjectEncodingMask(false, false, false),
        new NullExtension());               // Body

    private static final long EPOCH_OFFSET = 116444736000000000L;         //Offset between OPC UA epoch time and linux epoch time.
    private static final PascalString APPLICATION_URI = new PascalString("urn:apache:plc4x:client");
    private static final PascalString PRODUCT_URI = new PascalString("urn:apache:plc4x:client");
    private static final PascalString APPLICATION_TEXT = new PascalString("OPCUA client for the Apache PLC4X:PLC4J project");
    private static final long DEFAULT_CONNECTION_LIFETIME = 36000000;
    private final String sessionName = "UaSession:" + APPLICATION_TEXT.getStringValue() + ":" + RandomStringUtils.random(20, true, true);
    private final byte[] clientNonce = RandomUtils.nextBytes(40);
    private AtomicInteger requestHandleGenerator = new AtomicInteger(1);
    private PascalString policyId;
    private PascalString endpoint;
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
    private byte[] senderCertificate = null;
    private byte[] senderNonce = null;
    private PascalByteString certificateThumbprint = null;
    private boolean checkedEndpoints = false;
    private EncryptionHandler encryptionHandler = null;
    private OpcuaConfiguration configuration;
    private AtomicInteger channelId = new AtomicInteger(1);
    private AtomicInteger tokenId = new AtomicInteger(1);
    private NodeIdTypeDefinition authenticationToken = new NodeIdTwoByte((short) 0);
    private DriverContext driverContext;
    ConversationContext<OpcuaAPU> context;
    private SecureChannelTransactionManager channelTransactionManager = new SecureChannelTransactionManager();
    private long lifetime = DEFAULT_CONNECTION_LIFETIME;
    private CompletableFuture<Void> keepAlive;
    private int sendBufferSize;
    private int maxMessageSize;
    private AtomicLong senderSequenceNumber = new AtomicLong();

    public SecureChannel(DriverContext driverContext, OpcuaConfiguration configuration) {
        this.driverContext = driverContext;
        this.configuration = configuration;

        this.endpoint = new PascalString(configuration.getEndpoint());
        this.discovery = configuration.isDiscovery();
        this.username = configuration.getUsername();
        this.password = configuration.getPassword();
        this.certFile = configuration.getCertDirectory();
        this.securityPolicy = "http://opcfoundation.org/UA/SecurityPolicy#" + configuration.getSecurityPolicy();
        this.ckp = configuration.getCertificateKeyPair();

        if (configuration.getSecurityPolicy() != null && configuration.getSecurityPolicy().equals("Basic256Sha256")) {
            //Sender Certificate gets populated during the discover phase when encryption is enabled.
            this.senderCertificate = configuration.getSenderCertificate();
            this.encryptionHandler = new EncryptionHandler(this.ckp, this.senderCertificate, configuration.getSecurityPolicy());
            try {
                this.publicCertificate = new PascalByteString(this.ckp.getCertificate().getEncoded().length, this.ckp.getCertificate().getEncoded());
                this.isEncrypted = true;
            } catch (CertificateEncodingException e) {
                throw new PlcRuntimeException("Failed to encode the certificate");
            }
            this.thumbprint = configuration.getThumbprint();
        } else {
            this.encryptionHandler = new EncryptionHandler(this.ckp, this.senderCertificate, configuration.getSecurityPolicy());
            this.publicCertificate = NULL_BYTE_STRING;
            this.thumbprint = NULL_BYTE_STRING;
            this.isEncrypted = false;
        }
        this.keyStoreFile = configuration.getKeyStoreFile();
    }

    public void submit(ConversationContext<OpcuaAPU> context, Consumer<TimeoutException> onTimeout, BiConsumer<OpcuaAPU, Throwable> error, Consumer<byte[]> consumer, WriteBufferByteBased buffer) {
        int transactionId = channelTransactionManager.getTransactionIdentifier();

        //TODO: We need to split large messages up into chunks if it is larger than the sendBufferSize
        //      This value is negotiated when opening a channel

        OpcuaMessageRequest messageRequest = new OpcuaMessageRequest(FINAL_CHUNK,
            channelId.get(),
            tokenId.get(),
            transactionId,
            transactionId,
            buffer.getData());

        final OpcuaAPU apu;
        try {
            if (this.isEncrypted) {
                apu = OpcuaAPUIO.staticParse(encryptionHandler.encodeMessage(messageRequest, buffer.getData()), false);
            } else {
                apu = new OpcuaAPU(messageRequest);
            }
        } catch (ParseException e) {
           throw new PlcRuntimeException("Unable to encrypt message before sending");
        }

        Consumer<Integer> requestConsumer = t -> {
            try {
                ByteArrayOutputStream messageBuffer = new ByteArrayOutputStream();
                context.sendRequest(apu)
                    .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                    .onTimeout(onTimeout)
                    .onError(error)
                    .unwrap(apuMessage -> encryptionHandler.decodeMessage(apuMessage))
                    .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                    .check(p -> {
                        if (p.getRequestId() == transactionId) {
                            try {
                                messageBuffer.write(p.getMessage());
                                if (!(senderSequenceNumber.incrementAndGet() == (p.getSequenceNumber()))) {
                                    LOGGER.error("Sequence number isn't as expected, we might have missed a packet. - {} != {}", senderSequenceNumber.incrementAndGet(), p.getSequenceNumber());
                                    context.fireDisconnected();
                                }
                            } catch (IOException e) {
                                LOGGER.debug("Failed to store incoming message in buffer {}");
                                throw new PlcRuntimeException("Error while sending message");
                            }
                            if (p.getChunk().equals(FINAL_CHUNK)) {
                                return true;
                            } else {
                                return false;
                            }
                        } else {
                            return false;
                        }
                    })
                    .handle(opcuaResponse -> {
                        if (opcuaResponse.getChunk().equals(FINAL_CHUNK)) {
                            tokenId.set(opcuaResponse.getSecureTokenId());
                            channelId.set(opcuaResponse.getSecureChannelId());


                            consumer.accept(messageBuffer.toByteArray());
                        }
                    });
            } catch (Exception e) {
                throw new PlcRuntimeException("Error while sending message");
            }
        };
        LOGGER.debug("Submitting Transaction to TransactionManager {}", transactionId);
        channelTransactionManager.submit(requestConsumer, transactionId);
    }

    public void onConnect(ConversationContext<OpcuaAPU> context) {
        // Only the TCP transport supports login.
        LOGGER.debug("Opcua Driver running in ACTIVE mode.");
        this.context = context;

        OpcuaHelloRequest hello = new OpcuaHelloRequest(FINAL_CHUNK,
            VERSION,
            DEFAULT_RECEIVE_BUFFER_SIZE,
            DEFAULT_SEND_BUFFER_SIZE,
            DEFAULT_MAX_MESSAGE_SIZE,
            DEFAULT_MAX_CHUNK_COUNT,
            this.endpoint);

        Consumer<Integer> requestConsumer = t -> {
            context.sendRequest(new OpcuaAPU(hello))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaAcknowledgeResponse)
                .unwrap(p -> (OpcuaAcknowledgeResponse) p.getMessage())
                .handle(opcuaAcknowledgeResponse -> {
                    sendBufferSize = Math.min(opcuaAcknowledgeResponse.getReceiveBufferSize(), DEFAULT_SEND_BUFFER_SIZE);
                    maxMessageSize = Math.min(opcuaAcknowledgeResponse.getMaxMessageSize(), DEFAULT_MAX_MESSAGE_SIZE);
                    onConnectOpenSecureChannel(context, opcuaAcknowledgeResponse);
                });
        };
        channelTransactionManager.submit(requestConsumer, channelTransactionManager.getTransactionIdentifier());
    }

    public void onConnectOpenSecureChannel(ConversationContext<OpcuaAPU> context, OpcuaAcknowledgeResponse opcuaAcknowledgeResponse) {

        int transactionId = channelTransactionManager.getTransactionIdentifier();

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        OpenSecureChannelRequest openSecureChannelRequest = null;
        if (this.isEncrypted) {
            openSecureChannelRequest = new OpenSecureChannelRequest(
                requestHeader,
                VERSION,
                SecurityTokenRequestType.securityTokenRequestTypeIssue,
                MessageSecurityMode.messageSecurityModeSignAndEncrypt,
                new PascalByteString(clientNonce.length, clientNonce),
                lifetime);
        } else {
            openSecureChannelRequest = new OpenSecureChannelRequest(
                requestHeader,
                VERSION,
                SecurityTokenRequestType.securityTokenRequestTypeIssue,
                MessageSecurityMode.messageSecurityModeNone,
                NULL_BYTE_STRING,
                lifetime);
        }

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(openSecureChannelRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            openSecureChannelRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, extObject);

            OpcuaOpenRequest openRequest = new OpcuaOpenRequest(FINAL_CHUNK,
                0,
                new PascalString(this.securityPolicy),
                this.publicCertificate,
                this.thumbprint,
                transactionId,
                transactionId,
                buffer.getData());

            final OpcuaAPU apu;

            if (this.isEncrypted) {
                apu = OpcuaAPUIO.staticParse(encryptionHandler.encodeMessage(openRequest, buffer.getData()), false);
            } else {
                apu = new OpcuaAPU(openRequest);
            }

            Consumer<Integer> requestConsumer = t -> {
                context.sendRequest(apu)
                    .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                    .unwrap(apuMessage -> encryptionHandler.decodeMessage(apuMessage))
                    .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
                    .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
                    .check(p -> {
                        if (p.getRequestId() == transactionId) {
                            return true;
                        } else {
                            return false;
                        }
                    })
                    .handle(opcuaOpenResponse -> {
                        try {
                            ReadBuffer readBuffer = new ReadBufferByteBased(opcuaOpenResponse.getMessage(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                            ExtensionObject message = ExtensionObjectIO.staticParse(readBuffer, false);
                            //Store the initial sequence number from the server. there's no requirement for the server and client to use the same starting number.
                            senderSequenceNumber.set(opcuaOpenResponse.getSequenceNumber());
                            certificateThumbprint = opcuaOpenResponse.getReceiverCertificateThumbprint();

                            if (message.getBody() instanceof ServiceFault) {
                                ServiceFault fault = (ServiceFault) message.getBody();
                                LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                            } else {
                                LOGGER.debug("Got Secure Response Connection Response");
                                try {
                                    OpenSecureChannelResponse openSecureChannelResponse = (OpenSecureChannelResponse) message.getBody();
                                    tokenId.set((int) ((ChannelSecurityToken) openSecureChannelResponse.getSecurityToken()).getTokenId());
                                    channelId.set((int) ((ChannelSecurityToken) openSecureChannelResponse.getSecurityToken()).getChannelId());
                                    onConnectCreateSessionRequest(context);
                                } catch (PlcConnectionException e) {
                                    LOGGER.error("Error occurred while connecting to OPC UA server");
                                    e.printStackTrace();
                                }
                            }
                        } catch (ParseException e) {
                            e.printStackTrace();
                        }
                    });
            };
            LOGGER.debug("Submitting OpenSecureChannel with id of {}", transactionId);
            channelTransactionManager.submit(requestConsumer, transactionId);
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Open Secure Request");
        }
    }

    public void onConnectCreateSessionRequest(ConversationContext<OpcuaAPU> context) throws PlcConnectionException {

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            0L,
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        LocalizedText applicationName = new LocalizedText(
            true,
            true,
            new PascalString("en"),
            APPLICATION_TEXT);

        PascalString gatewayServerUri = NULL_STRING;
        PascalString discoveryProfileUri = NULL_STRING;
        int noOfDiscoveryUrls = -1;
        PascalString[] discoveryUrls = new PascalString[0];

        ApplicationDescription clientDescription = new ApplicationDescription(APPLICATION_URI,
            PRODUCT_URI,
            applicationName,
            ApplicationType.applicationTypeClient,
            gatewayServerUri,
            discoveryProfileUri,
            noOfDiscoveryUrls,
            discoveryUrls);

        CreateSessionRequest createSessionRequest = new CreateSessionRequest(
            requestHeader,
            clientDescription,
            NULL_STRING,
            this.endpoint,
            new PascalString(sessionName),
            new PascalByteString(clientNonce.length, clientNonce),
            NULL_BYTE_STRING,
            120000L,
            0L);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(createSessionRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            createSessionRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, extObject);

            Consumer<byte[]> consumer = opcuaResponse -> {
                        try {
                            ExtensionObject message = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false);
                            if (message.getBody() instanceof ServiceFault) {
                                ServiceFault fault = (ServiceFault) message.getBody();
                                LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                            } else {
                                LOGGER.debug("Got Create Session Response Connection Response");
                                try {
                                    CreateSessionResponse responseMessage;

                                    ExtensionObjectDefinition unknownExtensionObject = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false).getBody();
                                    if (unknownExtensionObject instanceof CreateSessionResponse) {
                                        responseMessage = (CreateSessionResponse) unknownExtensionObject;

                                        authenticationToken = responseMessage.getAuthenticationToken().getNodeId();

                                        onConnectActivateSessionRequest(context, responseMessage, (CreateSessionResponse) message.getBody());
                                    } else {
                                        ServiceFault serviceFault = (ServiceFault) unknownExtensionObject;
                                        ResponseHeader header = (ResponseHeader) serviceFault.getResponseHeader();
                                        LOGGER.error("Subscription ServiceFault returned from server with error code,  '{}'", header.getServiceResult().toString());

                                    }

                                } catch (PlcConnectionException e) {
                                    LOGGER.error("Error occurred while connecting to OPC UA server");
                                } catch (ParseException e) {
                                    LOGGER.error("Unable to parse the returned Subscription response");
                                    e.printStackTrace();
                                }
                            }
                        } catch (ParseException e) {
                            e.printStackTrace();
                        }

                    };

            Consumer<TimeoutException> timeout = e -> {
                LOGGER.error("Timeout while waiting for subscription response");
                e.printStackTrace();
            };

            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                LOGGER.error("Error while waiting for subscription response");
                e.printStackTrace();
            };

            submit(context, timeout, error, consumer, buffer);
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Create Session Request");
        }
    }

    private void onConnectActivateSessionRequest(ConversationContext<OpcuaAPU> context, CreateSessionResponse opcuaMessageResponse, CreateSessionResponse sessionResponse) throws PlcConnectionException, ParseException {

        senderCertificate = sessionResponse.getServerCertificate().getStringValue();
        encryptionHandler.setServerCertificate(EncryptionHandler.getCertificateX509(senderCertificate));
        this.senderNonce = sessionResponse.getServerNonce().getStringValue();
        UserTokenType tokenType = UserTokenType.userTokenTypeAnonymous;
        String[] endpoints = new String[3];
        try {
            InetAddress address = InetAddress.getByName(this.configuration.getHost());
            endpoints[0] = "opc.tcp://" + address.getHostAddress() + ":" + configuration.getPort() +  configuration.getTransportEndpoint();
            endpoints[1] = "opc.tcp://" + address.getHostName() + ":" + configuration.getPort() +  configuration.getTransportEndpoint();
            endpoints[2] = "opc.tcp://" + address.getCanonicalHostName() + ":" + configuration.getPort() +  configuration.getTransportEndpoint();
        } catch (UnknownHostException e) {
            e.printStackTrace();
        }

        for (String hostEndpoints : endpoints) {
            for (ExtensionObjectDefinition extensionObject : sessionResponse.getServerEndpoints()) {
                EndpointDescription endpointDescription = (EndpointDescription) extensionObject;
                if (endpointDescription.getEndpointUrl().getStringValue().equals(hostEndpoints)) {
                    for (ExtensionObjectDefinition userTokenCast : endpointDescription.getUserIdentityTokens()) {
                        UserTokenPolicy identityToken = (UserTokenPolicy) userTokenCast;
                        if ((identityToken.getTokenType() == UserTokenType.userTokenTypeAnonymous) && (this.username == null)) {
                            LOGGER.info("Using Endpoint {} with security {}", endpointDescription.getEndpointUrl().getStringValue(), identityToken.getPolicyId().getStringValue());
                            policyId = identityToken.getPolicyId();
                            tokenType = identityToken.getTokenType();
                        } else if ((identityToken.getTokenType() == UserTokenType.userTokenTypeUserName) && (this.username != null)) {
                            LOGGER.info("Using Endpoint {} with security {}", endpointDescription.getEndpointUrl().getStringValue(), identityToken.getPolicyId().getStringValue());
                            policyId = identityToken.getPolicyId();
                            tokenType = identityToken.getTokenType();
                        }
                    }
                }
            }
        }

        if (this.policyId == null) {
            throw new PlcRuntimeException("Unable to find endpoint - " + endpoints[1]);
        }

        ExtensionObject userIdentityToken = getIdentityToken(tokenType, policyId.getStringValue());

        int requestHandle = getRequestHandle();

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            requestHandle,
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        SignatureData clientSignature = new SignatureData(NULL_STRING, NULL_BYTE_STRING);

        SignedSoftwareCertificate[] signedSoftwareCertificate = new SignedSoftwareCertificate[1];

        signedSoftwareCertificate[0] = new SignedSoftwareCertificate(NULL_BYTE_STRING, NULL_BYTE_STRING);

        ActivateSessionRequest activateSessionRequest = new ActivateSessionRequest(
            requestHeader,
            clientSignature,
            0,
            null,
            0,
            null,
            userIdentityToken,
            clientSignature);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(activateSessionRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            activateSessionRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, extObject);

            Consumer<byte[]> consumer = opcuaResponse -> {
                try {
                    ExtensionObject message = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false);
                    if (message.getBody() instanceof ServiceFault) {
                        ServiceFault fault = (ServiceFault) message.getBody();
                        LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                    } else {
                        LOGGER.debug("Got Activate Session Response Connection Response");
                        try {
                            ActivateSessionResponse responseMessage;

                            ExtensionObjectDefinition unknownExtensionObject = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false).getBody();
                            if (unknownExtensionObject instanceof ActivateSessionResponse) {
                                responseMessage = (ActivateSessionResponse) unknownExtensionObject;

                                long returnedRequestHandle = ((ResponseHeader) responseMessage.getResponseHeader()).getRequestHandle();
                                if (!(requestHandle == returnedRequestHandle)) {
                                    LOGGER.error("Request handle isn't as expected, we might have missed a packet. {} != {}", requestHandle, returnedRequestHandle);
                                }

                                // Send an event that connection setup is complete.
                                keepAlive();
                                context.fireConnected();
                            } else {
                                ServiceFault serviceFault = (ServiceFault) unknownExtensionObject;
                                ResponseHeader header = (ResponseHeader) serviceFault.getResponseHeader();
                                LOGGER.error("Subscription ServiceFault returned from server with error code,  '{}'", header.getServiceResult().toString());
                            }
                        } catch (ParseException e) {
                            LOGGER.error("Unable to parse the returned Subscription response");
                            e.printStackTrace();
                        }
                    }
                } catch (ParseException e) {
                    e.printStackTrace();
                }

            };

            Consumer<TimeoutException> timeout = e -> {
                LOGGER.error("Timeout while waiting for activate session response");
                e.printStackTrace();
            };

            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                LOGGER.error("Error while waiting for activate session response");
                e.printStackTrace();
            };

            submit(context, timeout, error, consumer, buffer);
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Activate Session Request");
        }
    }

    public void onDisconnect(ConversationContext<OpcuaAPU> context) {
        LOGGER.info("Disconnecting");
        int requestHandle = getRequestHandle();

        if (keepAlive != null) {
            keepAlive.complete(null);
        }

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, 473),
            null,
            null);    //Identifier for OpenSecureChannel

        RequestHeader requestHeader = new RequestHeader(
            new NodeId(authenticationToken),
            getCurrentDateTime(),
            requestHandle,                                         //RequestHandle
            0L,
            NULL_STRING,
            5000L,
            NULL_EXTENSION_OBJECT);

        CloseSessionRequest closeSessionRequest = new CloseSessionRequest(
            requestHeader,
            true);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            closeSessionRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, extObject);

            Consumer<byte[]> consumer = opcuaResponse -> {
                try {
                    ExtensionObject message = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false);
                    if (message.getBody() instanceof ServiceFault) {
                        ServiceFault fault = (ServiceFault) message.getBody();
                        LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                    } else {
                        LOGGER.debug("Got Close Session Response Connection Response");
                        try {
                            CloseSessionResponse responseMessage;

                            ExtensionObjectDefinition unknownExtensionObject = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false).getBody();
                            if (unknownExtensionObject instanceof CloseSessionResponse) {
                                responseMessage = (CloseSessionResponse) unknownExtensionObject;

                                LOGGER.trace("Got Close Session Response Connection Response" + responseMessage.toString());
                                onDisconnectCloseSecureChannel(context);
                            } else {
                                ServiceFault serviceFault = (ServiceFault) unknownExtensionObject;
                                ResponseHeader header = (ResponseHeader) serviceFault.getResponseHeader();
                                LOGGER.error("Subscription ServiceFault returned from server with error code,  '{}'", header.getServiceResult().toString());
                            }
                        } catch (ParseException e) {
                            LOGGER.error("Unable to parse the returned Close Session response");
                            e.printStackTrace();
                        }
                    }
                } catch (ParseException e) {
                    e.printStackTrace();
                }

            };

            Consumer<TimeoutException> timeout = e -> {
                LOGGER.error("Timeout while waiting for close session response");
                e.printStackTrace();
            };

            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                LOGGER.error("Error while waiting for close session response");
                e.printStackTrace();
            };

            submit(context, timeout, error, consumer, buffer);
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Close Session Request");
        }
    }

    private void onDisconnectCloseSecureChannel(ConversationContext<OpcuaAPU> context) {

        int transactionId = channelTransactionManager.getTransactionIdentifier();

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        CloseSecureChannelRequest closeSecureChannelRequest = new CloseSecureChannelRequest(requestHeader);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(closeSecureChannelRequest.getIdentifier())),
            null,
            null);

        OpcuaCloseRequest closeRequest = new OpcuaCloseRequest(FINAL_CHUNK,
            channelId.get(),
            tokenId.get(),
            transactionId,
            transactionId,
            new ExtensionObject(
                expandedNodeId,
                null,
                closeSecureChannelRequest));

        Consumer<Integer> requestConsumer = t -> {
            context.sendRequest(new OpcuaAPU(closeRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .check(p -> {
                    if (p.getRequestId() == transactionId) {
                        return true;
                    } else {
                        return false;
                    }
                })
                .handle(opcuaMessageResponse -> {
                    LOGGER.trace("Got Close Secure Channel Response" + opcuaMessageResponse.toString());
                });

            context.fireDisconnected();
        };

        channelTransactionManager.submit(requestConsumer, transactionId);

    }

    public void onDiscover(ConversationContext<OpcuaAPU> context) {
        // Only the TCP transport supports login.
        LOGGER.debug("Opcua Driver running in ACTIVE mode, discovering endpoints");

        OpcuaHelloRequest hello = new OpcuaHelloRequest(FINAL_CHUNK,
            VERSION,
            DEFAULT_RECEIVE_BUFFER_SIZE,
            DEFAULT_SEND_BUFFER_SIZE,
            DEFAULT_MAX_MESSAGE_SIZE,
            DEFAULT_MAX_CHUNK_COUNT,
            this.endpoint);

        Consumer<Integer> requestConsumer = t -> {
            context.sendRequest(new OpcuaAPU(hello))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaAcknowledgeResponse)
                .unwrap(p -> (OpcuaAcknowledgeResponse) p.getMessage())
                .handle(opcuaAcknowledgeResponse -> {
                    LOGGER.debug("Got Hello Response Connection Response");
                    onDiscoverOpenSecureChannel(context, opcuaAcknowledgeResponse);
                });
        };

        channelTransactionManager.submit(requestConsumer, 1);
    }

    public void onDiscoverOpenSecureChannel(ConversationContext<OpcuaAPU> context, OpcuaAcknowledgeResponse opcuaAcknowledgeResponse) {
        int transactionId = channelTransactionManager.getTransactionIdentifier();

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        OpenSecureChannelRequest openSecureChannelRequest = new OpenSecureChannelRequest(
            requestHeader,
            VERSION,
            SecurityTokenRequestType.securityTokenRequestTypeIssue,
            MessageSecurityMode.messageSecurityModeNone,
            NULL_BYTE_STRING,
            lifetime);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(openSecureChannelRequest.getIdentifier())),
            null,
            null);


        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(openSecureChannelRequest.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, new ExtensionObject(
                expandedNodeId,
                null,
                openSecureChannelRequest));

            OpcuaOpenRequest openRequest = new OpcuaOpenRequest(FINAL_CHUNK,
                0,
                SECURITY_POLICY_NONE,
                NULL_BYTE_STRING,
                NULL_BYTE_STRING,
                transactionId,
                transactionId,
                buffer.getData());

            Consumer<Integer> requestConsumer = t -> {
                context.sendRequest(new OpcuaAPU(openRequest))
                    .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                    .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
                    .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
                    .check(p -> {
                        if (p.getRequestId() == transactionId) {
                            return true;
                        } else {
                            return false;
                        }
                    })
                    .handle(opcuaOpenResponse -> {
                        try {
                            ExtensionObject message = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaOpenResponse.getMessage(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false);
                            if (message.getBody() instanceof ServiceFault) {
                                ServiceFault fault = (ServiceFault) message.getBody();
                                LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                            } else {
                                LOGGER.debug("Got Secure Response Connection Response");
                                try {
                                    onDiscoverGetEndpointsRequest(context, opcuaOpenResponse, (OpenSecureChannelResponse) message.getBody());
                                } catch (PlcConnectionException e) {
                                    LOGGER.error("Error occurred while connecting to OPC UA server");
                                }
                            }
                        } catch (ParseException e) {
                            e.printStackTrace();
                        }
                    });
            };

            channelTransactionManager.submit(requestConsumer, transactionId);
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Create Session Request");
        }
    }

    public void onDiscoverGetEndpointsRequest(ConversationContext<OpcuaAPU> context, OpcuaOpenResponse opcuaOpenResponse, OpenSecureChannelResponse openSecureChannelResponse) throws PlcConnectionException {
        certificateThumbprint = opcuaOpenResponse.getReceiverCertificateThumbprint();
        tokenId.set((int) ((ChannelSecurityToken) openSecureChannelResponse.getSecurityToken()).getTokenId());
        channelId.set((int) ((ChannelSecurityToken) openSecureChannelResponse.getSecurityToken()).getChannelId());

        int transactionId = channelTransactionManager.getTransactionIdentifier();

        Integer nextSequenceNumber = opcuaOpenResponse.getSequenceNumber() + 1;
        Integer nextRequestId = opcuaOpenResponse.getRequestId() + 1;

        if (!(transactionId == nextSequenceNumber)) {
            LOGGER.error("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
            throw new PlcConnectionException("Sequence number isn't as expected, we might have missed a packet. - " +  transactionId + " != " + nextSequenceNumber);
        }

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            0L,
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        GetEndpointsRequest endpointsRequest = new GetEndpointsRequest(
            requestHeader,
            this.endpoint,
            0,
            null,
            0,
            null);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(endpointsRequest.getIdentifier())),
            null,
            null);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(endpointsRequest.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, new ExtensionObject(
                expandedNodeId,
                null,
                endpointsRequest));

            OpcuaMessageRequest messageRequest = new OpcuaMessageRequest(FINAL_CHUNK,
                channelId.get(),
                tokenId.get(),
                nextSequenceNumber,
                nextRequestId,
                buffer.getData());

            Consumer<Integer> requestConsumer = t -> {
                context.sendRequest(new OpcuaAPU(messageRequest))
                    .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                    .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                    .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                    .check(p -> {
                        if (p.getRequestId() == transactionId) {
                            return true;
                        } else {
                            return false;
                        }
                    })
                    .handle(opcuaMessageResponse -> {
                        try {
                            ExtensionObject message = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaMessageResponse.getMessage(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN), false);
                            if (message.getBody() instanceof ServiceFault) {
                                ServiceFault fault = (ServiceFault) message.getBody();
                                LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                            } else {
                                LOGGER.debug("Got Create Session Response Connection Response");
                                GetEndpointsResponse response = (GetEndpointsResponse) message.getBody();

                                EndpointDescription[] endpoints = (EndpointDescription[]) response.getEndpoints();
                                for (EndpointDescription endpoint : endpoints) {
                                    if (endpoint.getEndpointUrl().getStringValue().equals(this.endpoint.getStringValue()) && endpoint.getSecurityPolicyUri().getStringValue().equals(this.securityPolicy)) {
                                        LOGGER.info("Found OPC UA endpoint {}", this.endpoint.getStringValue());
                                        this.configuration.setSenderCertificate(endpoint.getServerCertificate().getStringValue());
                                    }
                                }

                                try {
                                    MessageDigest messageDigest = MessageDigest.getInstance("SHA-1");
                                    byte[] digest = messageDigest.digest(this.configuration.getSenderCertificate());
                                    this.configuration.setThumbprint(new PascalByteString(digest.length, digest));
                                } catch (NoSuchAlgorithmException e) {
                                    LOGGER.error("Failed to find hashing algorithm");
                                }
                                onDiscoverCloseSecureChannel(context, response);
                            }
                        } catch (ParseException e) {
                            e.printStackTrace();
                        }
                    });
            };

            channelTransactionManager.submit(requestConsumer, transactionId);
        } catch (ParseException e) {
            LOGGER.error("Unable to to Parse Create Session Request");
        }
    }

    private void onDiscoverCloseSecureChannel(ConversationContext<OpcuaAPU> context, GetEndpointsResponse message) {

        int transactionId = channelTransactionManager.getTransactionIdentifier();

        RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
            getCurrentDateTime(),
            0L,                                         //RequestHandle
            0L,
            NULL_STRING,
            REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        CloseSecureChannelRequest closeSecureChannelRequest = new CloseSecureChannelRequest(requestHeader);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(closeSecureChannelRequest.getIdentifier())),
            null,
            null);

        OpcuaCloseRequest closeRequest = new OpcuaCloseRequest(FINAL_CHUNK,
            channelId.get(),
            tokenId.get(),
            transactionId,
            transactionId,
            new ExtensionObject(
                expandedNodeId,
                null,
                closeSecureChannelRequest));

        Consumer<Integer> requestConsumer = t -> {
            context.sendRequest(new OpcuaAPU(closeRequest))
                .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                .check(p -> p.getMessage() instanceof OpcuaMessageResponse)
                .unwrap(p -> (OpcuaMessageResponse) p.getMessage())
                .check(p -> {
                    if (p.getRequestId() == transactionId) {
                        return true;
                    } else {
                        return false;
                    }
                })
                .handle(opcuaMessageResponse -> {
                    LOGGER.trace("Got Close Secure Channel Response" + opcuaMessageResponse.toString());
                    // Send an event that connection setup is complete.
                    context.fireDiscovered(this.configuration);
                });
        };

        channelTransactionManager.submit(requestConsumer, transactionId);
    }

    private void keepAlive() {
        keepAlive = CompletableFuture.supplyAsync(() -> {
            while(true) {

                try {
                    Thread.sleep((long) Math.ceil(this.lifetime * 0.75f));
                } catch (InterruptedException e) {
                    LOGGER.trace("Interrupted Exception");
                }

                int transactionId = channelTransactionManager.getTransactionIdentifier();

                RequestHeader requestHeader = new RequestHeader(new NodeId(authenticationToken),
                    getCurrentDateTime(),
                    0L,                                         //RequestHandle
                    0L,
                    NULL_STRING,
                    REQUEST_TIMEOUT_LONG,
                    NULL_EXTENSION_OBJECT);

                OpenSecureChannelRequest openSecureChannelRequest = null;
                if (this.isEncrypted) {
                    openSecureChannelRequest = new OpenSecureChannelRequest(
                        requestHeader,
                        VERSION,
                        SecurityTokenRequestType.securityTokenRequestTypeIssue,
                        MessageSecurityMode.messageSecurityModeSignAndEncrypt,
                        new PascalByteString(clientNonce.length, clientNonce),
                        lifetime);
                } else {
                    openSecureChannelRequest = new OpenSecureChannelRequest(
                        requestHeader,
                        VERSION,
                        SecurityTokenRequestType.securityTokenRequestTypeIssue,
                        MessageSecurityMode.messageSecurityModeNone,
                        NULL_BYTE_STRING,
                        lifetime);
                }

                ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, Integer.valueOf(openSecureChannelRequest.getIdentifier())),
                    null,
                    null);

                ExtensionObject extObject = new ExtensionObject(
                    expandedNodeId,
                    null,
                    openSecureChannelRequest);

                try {
                    WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                    ExtensionObjectIO.staticSerialize(buffer, extObject);

                    OpcuaOpenRequest openRequest = new OpcuaOpenRequest(FINAL_CHUNK,
                        0,
                        new PascalString(this.securityPolicy),
                        this.publicCertificate,
                        this.thumbprint,
                        transactionId,
                        transactionId,
                        buffer.getData());

                    final OpcuaAPU apu;

                    if (this.isEncrypted) {
                        apu = OpcuaAPUIO.staticParse(encryptionHandler.encodeMessage(openRequest, buffer.getData()), false);
                    } else {
                        apu = new OpcuaAPU(openRequest);
                    }

                    Consumer<Integer> requestConsumer = t -> {
                        context.sendRequest(apu)
                            .expectResponse(OpcuaAPU.class, REQUEST_TIMEOUT)
                            .unwrap(apuMessage -> encryptionHandler.decodeMessage(apuMessage))
                            .check(p -> p.getMessage() instanceof OpcuaOpenResponse)
                            .unwrap(p -> (OpcuaOpenResponse) p.getMessage())
                            .check(p -> {
                                if (p.getRequestId() == transactionId) {
                                    return true;
                                } else {
                                    return false;
                                }
                            })
                            .handle(opcuaOpenResponse -> {
                                try {
                                    ReadBufferByteBased readBuffer = new ReadBufferByteBased(opcuaOpenResponse.getMessage(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                                    ExtensionObject message = ExtensionObjectIO.staticParse(readBuffer, false);

                                    if (message.getBody() instanceof ServiceFault) {
                                        ServiceFault fault = (ServiceFault) message.getBody();
                                        LOGGER.error("Failed to connect to opc ua server for the following reason:- {}, {}", ((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode(), OpcuaStatusCode.enumForValue(((ResponseHeader) fault.getResponseHeader()).getServiceResult().getStatusCode()));
                                    } else {
                                        LOGGER.debug("Got Secure Response Connection Response");
                                        OpenSecureChannelResponse openSecureChannelResponse = (OpenSecureChannelResponse) message.getBody();
                                        ChannelSecurityToken token = (ChannelSecurityToken) openSecureChannelResponse.getSecurityToken();
                                        certificateThumbprint = opcuaOpenResponse.getReceiverCertificateThumbprint();
                                        tokenId.set((int) token.getTokenId());
                                        channelId.set((int) token.getChannelId());
                                        lifetime = token.getRevisedLifetime();
                                    }
                                } catch (ParseException e) {
                                    e.printStackTrace();
                                }
                            });
                    };
                    channelTransactionManager.submit(requestConsumer, transactionId);
                } catch (ParseException e) {
                    LOGGER.error("Unable to to Parse Open Secure Request");
                }
            }
            }
        );
    }

    /**
     * Returns the next request handle
     *
     * @return the next sequential request handle
     */
    public int getRequestHandle() {
        int transactionId = requestHandleGenerator.getAndIncrement();
        if(requestHandleGenerator.get() == SecureChannelTransactionManager.DEFAULT_MAX_REQUEST_ID) {
            requestHandleGenerator.set(1);
        }
        return transactionId;
    }

    /**
     * Returns the authentication token for the current connection
     *
     * @return a NodeId Authentication token
     */
    public NodeId getAuthenticationToken() {
        return new NodeId(this.authenticationToken);
    }

    /**
     * Gets the Channel identifier for the current channel
     *
     * @return int representing the channel identifier
     */
    public int getChannelId() {
        return this.channelId.get();
    }

    /**
     * Gets the Token Identifier
     *
     * @return int representing the token identifier
     */
    public int getTokenId() {
        return this.tokenId.get();
    }

    /**
     * Creates an IdentityToken to authenticate with a server.
     * @param securityPolicy
     * @return returns an ExtensionObject with an IdentityToken.
     */
    private ExtensionObject getIdentityToken(UserTokenType tokenType, String securityPolicy) {
        ExpandedNodeId extExpandedNodeId = null;
        ExtensionObject userIdentityToken = null;
        switch (tokenType) {
            case userTokenTypeAnonymous:
                //If we aren't using authentication tell the server we would like to login anonymously
                AnonymousIdentityToken anonymousIdentityToken = new AnonymousIdentityToken();

                extExpandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, OpcuaNodeIdServices.AnonymousIdentityToken_Encoding_DefaultBinary.getValue()),
                    null,
                    null);

                return new ExtensionObject(
                    extExpandedNodeId,
                    new ExtensionObjectEncodingMask(false, false, true),
                    new UserIdentityToken(new PascalString(securityPolicy), anonymousIdentityToken));
            case userTokenTypeUserName:
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

                byte[] encryptedPassword = encryptionHandler.encryptPassword(encodeablePassword);
                UserNameIdentityToken userNameIdentityToken =  new UserNameIdentityToken(
                    new PascalString(this.username),
                    new PascalByteString(encryptedPassword.length, encryptedPassword),
                    new PascalString(PASSWORD_ENCRYPTION_ALGORITHM)
                );

                extExpandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, OpcuaNodeIdServices.UserNameIdentityToken_Encoding_DefaultBinary.getValue()),
                    null,
                    null);

                return new ExtensionObject(
                    extExpandedNodeId,
                    new ExtensionObjectEncodingMask(false, false, true),
                    new UserIdentityToken(new PascalString(securityPolicy), userNameIdentityToken));
        }
        return null;
    }

    public static long getCurrentDateTime() {
        return (System.currentTimeMillis() * 10000) + EPOCH_OFFSET;
    }
}
