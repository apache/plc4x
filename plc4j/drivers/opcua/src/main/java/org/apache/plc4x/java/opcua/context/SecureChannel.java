/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, PROTOCOL_VERSION_0 2.0 (the
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
package org.apache.plc4x.java.opcua.context;

import static java.util.concurrent.Executors.newSingleThreadScheduledExecutor;
import static org.apache.plc4x.java.opcua.readwrite.ChunkType.*;

import java.io.ByteArrayInputStream;
import java.security.GeneralSecurityException;
import java.security.Signature;
import java.security.cert.CertificateEncodingException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.function.Function;
import java.util.function.Supplier;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.SignatureAlgorithm;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager.RequestTransaction;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Stream;

import static java.util.concurrent.Executors.newSingleThreadExecutor;

public class SecureChannel {

    private static final Logger LOGGER = LoggerFactory.getLogger(SecureChannel.class);
    private static final String PASSWORD_ENCRYPTION_ALGORITHM = "http://www.w3.org/2001/04/xmlenc#rsa-oaep";
    public static final PascalString NULL_STRING = new PascalString("");
    public static final PascalByteString NULL_BYTE_STRING = new PascalByteString(-1, null);
    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("(.(?<transportCode>tcp|https?))?://" +
        "(?<transportHost>[\\w.-]+)(:" +
        "(?<transportPort>\\d*))?");

    public static final Pattern URI_PATTERN = Pattern.compile("^(?<protocolCode>opc)" +
        INET_ADDRESS_PATTERN +
        "(?<transportEndpoint>[\\w/=]*)[?]?"
    );

    private static final PascalString APPLICATION_URI = new PascalString("urn:apache:plc4x:client");
    private static final PascalString PRODUCT_URI = new PascalString("urn:apache:plc4x:client");
    private static final PascalString APPLICATION_TEXT = new PascalString("OPCUA client for the Apache PLC4X:PLC4J project");
    public static final ScheduledExecutorService KEEP_ALIVE_EXECUTOR = newSingleThreadScheduledExecutor(runnable -> new Thread(runnable, "plc4x-opcua-keep-alive"));
    private final String sessionName = "UaSession:" + APPLICATION_TEXT.getStringValue() + ":" + RandomStringUtils.random(20, true, true);
    private final PascalByteString localCertificateString;
    private final PascalByteString remoteCertificateThumbprint;
    private PascalString policyId;
    private UserTokenType tokenType;
    private final PascalString endpoint;
    private final String username;
    private final String password;
    private final RequestTransactionManager tm;
    private final OpcuaConfiguration configuration;
    private final OpcuaDriverContext driverContext;
    private final Conversation conversation;
    private ScheduledFuture<?> keepAlive;
    private final List<String> endpoints = new ArrayList<>();
    private double sessionTimeout;
    private long revisedLifetime;

    public SecureChannel(Conversation conversation, RequestTransactionManager tm, OpcuaDriverContext driverContext, OpcuaConfiguration configuration, PlcAuthentication authentication) {
        this.conversation = conversation;
        this.tm = tm;
        this.configuration = configuration;
        this.driverContext = driverContext;
        this.endpoint = new PascalString(driverContext.getEndpoint());
        this.sessionTimeout = configuration.getSessionTimeout();
        if (authentication != null) {
            if (authentication instanceof PlcUsernamePasswordAuthentication) {
                this.username = ((PlcUsernamePasswordAuthentication) authentication).getUsername();
                this.password = ((PlcUsernamePasswordAuthentication) authentication).getPassword();
            } else {
                throw new PlcRuntimeException("This type of connection only supports username-password authentication");
            }
        } else {
            this.username = configuration.getUsername();
            this.password = configuration.getPassword();
        }

        // Generate a list of endpoints we can use.
        try {
            InetAddress address = InetAddress.getByName(driverContext.getHost());
            this.endpoints.add(address.getHostAddress());
            this.endpoints.add(address.getHostName());
            this.endpoints.add(address.getCanonicalHostName());
        } catch (UnknownHostException e) {
            LOGGER.warn("Unable to resolve host name. Using original host from connection string which may cause issues connecting to server");
            this.endpoints.add(driverContext.getHost());
        }

        if (conversation.getSecurityPolicy() == SecurityPolicy.NONE) {
            this.localCertificateString = NULL_BYTE_STRING;
            this.remoteCertificateThumbprint = NULL_BYTE_STRING;
        } else {
            CertificateKeyPair keyPair = driverContext.getCertificateKeyPair();
            this.remoteCertificateThumbprint = driverContext.getThumbprint();
            try {
                byte[] encoded = keyPair.getCertificate().getEncoded();
                this.localCertificateString = new PascalByteString(encoded.length, encoded);
            } catch (CertificateEncodingException e) {
                throw new PlcRuntimeException("Could not decode certificate", e);
            }
        }
    }

    public CompletableFuture<ActivateSessionResponse> onConnect() {
        // Only the TCP transport supports login.
        LOGGER.debug("Opcua Driver running in ACTIVE mode.");
        return conversation.requestHello()
            .thenCompose(r -> onConnectOpenSecureChannel(SecurityTokenRequestType.securityTokenRequestTypeIssue, 0, 0))
            .thenCompose(r -> onConnectCreateSessionRequest())
            .thenCompose(r -> onConnectActivateSessionRequest(r))
            .thenApply(response -> {
                renewToken();
                return response;
            });
    }

    public CompletableFuture<OpenSecureChannelResponse> onConnectOpenSecureChannel(SecurityTokenRequestType securityTokenRequestType, int secureChannelId, int requestId) {
        LOGGER.debug("Sending open secure channel message to {}", this.driverContext.getEndpoint());

        RequestHeader requestHeader = conversation.createRequestHeader(configuration.getNegotiationTimeout(), requestId);

        OpenSecureChannelRequest openSecureChannelRequest;
        byte[] localNonce = conversation.createNonce();
        if (conversation.getSecurityPolicy() != SecurityPolicy.NONE) {
            openSecureChannelRequest = new OpenSecureChannelRequest(
                requestHeader,
                OpcuaConstants.PROTOCOLVERSION,
                securityTokenRequestType,
                configuration.getMessageSecurity().getMode(),
                new PascalByteString(localNonce.length, localNonce),
                configuration.getChannelLifetime() // lifetime
            );
        } else {
            openSecureChannelRequest = new OpenSecureChannelRequest(
                requestHeader,
                OpcuaConstants.PROTOCOLVERSION,
                securityTokenRequestType,
                MessageSecurityMode.messageSecurityModeNone,
                NULL_BYTE_STRING,
                configuration.getChannelLifetime() // lifetime
            );
        }

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false, false,
            new NodeIdFourByte((short) 0, Integer.parseInt(openSecureChannelRequest.getIdentifier())),
            null, null
        );
        ExtensionObject extObject = new ExtensionObject(expandedNodeId, null, openSecureChannelRequest);

        Function<CallContext, OpcuaOpenRequest> openRequest = context -> {
            LOGGER.debug("Submitting OpenSecureChannel with id of {}", context.getRequestId());
            return new OpcuaOpenRequest(FINAL, new OpenChannelMessageRequest(secureChannelId,
                new PascalString(conversation.getSecurityPolicy().getSecurityPolicyUri()),
                this.localCertificateString,
                this.remoteCertificateThumbprint
            ),
            new ExtensiblePayload(
                new SequenceHeader(context.getNextSequenceNumber(), context.getRequestId()),
                extObject
            ));
        };

        return conversation.requestChannelOpen(openRequest)
            .thenApply(response -> {
                LOGGER.info("Received open channel response {}, parsing it", response.getMessage().getSequenceHeader().getRequestId());
                return response;
            })
            .thenApply(this::onOpenResponse)
            .thenApply(openSecureChannelResponse -> {
                ChannelSecurityToken securityToken = (ChannelSecurityToken) openSecureChannelResponse.getSecurityToken();
                LOGGER.debug("Opened secure response id: {}, channel id:{}, token:{} lifetime:{}", openSecureChannelResponse.getIdentifier(),
                    securityToken.getChannelId(), securityToken.getTokenId(), securityToken.getRevisedLifetime());

                // store server and client nonce
                conversation.setRemoteNonce(openSecureChannelResponse.getServerNonce().getStringValue());
                conversation.setLocalNonce(localNonce);
                conversation.setSecurityHeader(new SecurityHeader(securityToken.getChannelId(), securityToken.getTokenId()));
                revisedLifetime = securityToken.getRevisedLifetime();
                return openSecureChannelResponse;
            });
    }

    public CompletableFuture<CreateSessionResponse> onConnectCreateSessionRequest() {
        LOGGER.debug("Sending create session request to {}", this.driverContext.getEndpoint());
        RequestHeader requestHeader = conversation.createRequestHeader();

        LocalizedText applicationName = new LocalizedText(
            true,
            true,
            new PascalString("en"),
            APPLICATION_TEXT
        );

        int noOfDiscoveryUrls = -1;
        List<PascalString> discoveryUrls = new ArrayList<>(0);

        ApplicationDescription clientDescription = new ApplicationDescription(
            driverContext.getApplicationUri().map(PascalString::new).orElse(APPLICATION_URI),
            PRODUCT_URI,
            applicationName,
            ApplicationType.applicationTypeClient,
            NULL_STRING,
            NULL_STRING,
            noOfDiscoveryUrls,
            discoveryUrls
        );

        byte[] temporaryNonce = conversation.createNonce(32);
        CreateSessionRequest createSessionRequest = new CreateSessionRequest(
            requestHeader,
            clientDescription,
            NULL_STRING,
            this.endpoint,
            new PascalString(sessionName),
            conversation.getSecurityPolicy() == SecurityPolicy.NONE ? NULL_BYTE_STRING : createPascalString(temporaryNonce),
            conversation.getSecurityPolicy() == SecurityPolicy.NONE ? NULL_BYTE_STRING : localCertificateString,
            sessionTimeout,
            0L
        );

        return conversation.submit(createSessionRequest, CreateSessionResponse.class)
            .thenApply(sessionResponse -> {
                if (conversation.getSecurityPolicy() != SecurityPolicy.NONE) {
                    // verify temporaryNonce against server returned data
                    SignatureData signatureData = extractSignatureData(sessionResponse.getServerSignature());
                    if (signatureData == null) {
                        throw new IllegalArgumentException("Returned signature data is not valid");
                    }

                    String algorithm = signatureData.getAlgorithm().getStringValue();

                    SignatureAlgorithm signatureAlgorithm = conversation.getSecurityPolicy().getAsymmetricSignatureAlgorithm();
                    if (!signatureAlgorithm.getUri().equals(algorithm)) {
                        throw new IllegalArgumentException("Invalid signature algorithm. Expected " + signatureAlgorithm.getUri());
                    }
                    try {
                        int certificateLength = localCertificateString.getStringLength();
                        byte[] rawData = new byte[certificateLength + 32];
                        System.arraycopy(localCertificateString.getStringValue(), 0, rawData, 0, certificateLength);
                        System.arraycopy(temporaryNonce, 0, rawData, certificateLength, 32);
                        X509Certificate remoteCertificate = conversation.getRemoteCertificate();
                        // make sure returned certificate is trusted
                        driverContext.getCertificateVerifier().checkCertificateTrusted(remoteCertificate);

                        Signature signature = signatureAlgorithm.getSignature();
                        signature.initVerify(remoteCertificate.getPublicKey());
                        signature.update(rawData);
                        if (!signature.verify(signatureData.getSignature().getStringValue())) {
                            throw new IllegalArgumentException("Could not verify server signature");
                        }
                    } catch (GeneralSecurityException e) {
                        throw new RuntimeException(e);
                    }

                }
                return sessionResponse;
            })
            .thenApply(responseMessage -> {
                conversation.setAuthenticationToken(responseMessage.getAuthenticationToken().getNodeId());
                sessionTimeout = responseMessage.getRevisedSessionTimeout();
                return responseMessage;
            });
    }

    private SignatureData extractSignatureData(ExtensionObjectDefinition object) {
        if (object instanceof SignatureData) {
            return (SignatureData) object;
        }
        return null;
    }

    private CompletableFuture<ActivateSessionResponse> onConnectActivateSessionRequest(CreateSessionResponse sessionResponse) {
        LOGGER.debug("Sending activate session request to {}", this.driverContext.getEndpoint());
        conversation.setRemoteCertificate(getX509Certificate(sessionResponse.getServerCertificate().getStringValue()));
        conversation.setRemoteNonce(sessionResponse.getServerNonce().getStringValue());

        String[] endpoints = new String[3];
        try {
            InetAddress address = InetAddress.getByName(driverContext.getHost());
            endpoints[0] = "opc.tcp://" + address.getHostAddress() + ":" + driverContext.getPort() + driverContext.getTransportEndpoint();
            endpoints[1] = "opc.tcp://" + address.getHostName() + ":" + driverContext.getPort() + driverContext.getTransportEndpoint();
            endpoints[2] = "opc.tcp://" + address.getCanonicalHostName() + ":" + driverContext.getPort() + driverContext.getTransportEndpoint();
        } catch (UnknownHostException e) {
            LOGGER.debug("error getting host", e);
        }

        selectEndpoint(sessionResponse);

        if (this.policyId == null) {
            throw new PlcRuntimeException("Unable to find endpoint - " + endpoints[1]);
        }

        ExtensionObject userIdentityToken = getIdentityToken(this.tokenType, policyId.getStringValue());
        RequestHeader requestHeader = conversation.createRequestHeader();
        SignatureData clientSignature = new SignatureData(NULL_STRING, NULL_BYTE_STRING);
        if (conversation.getSecurityPolicy() != SecurityPolicy.NONE) {
            try {
                clientSignature = conversation.createClientSignature();
            } catch (GeneralSecurityException e) {
                throw new PlcRuntimeException("Could not create client signature", e);
            }
        }

        ActivateSessionRequest activateSessionRequest = new ActivateSessionRequest(
            requestHeader,
            clientSignature,
            0,
            null,
            0,
            null,
            userIdentityToken,
            clientSignature
        );

        return conversation.submit(activateSessionRequest, ActivateSessionResponse.class).thenApply(responseMessage -> {
            conversation.setRemoteNonce(responseMessage.getServerNonce().getStringValue());
            return responseMessage;
        });
    }

    public void onDisconnect() {
        LOGGER.info("Disconnecting");

        if (keepAlive != null) {
            keepAlive.cancel(true);
            keepAlive = null;
        }

        RequestHeader requestHeader = conversation.createRequestHeader(50000L);
        CloseSessionRequest closeSessionRequest = new CloseSessionRequest(requestHeader, true);
        conversation.submit(closeSessionRequest, CloseSessionResponse.class).thenAccept(responseMessage -> {
            LOGGER.trace("Got Close Session Response Connection Response" + responseMessage);
            onDisconnectCloseSecureChannel();
        });
    }

    private void onDisconnectCloseSecureChannel() {
        RequestHeader requestHeader = conversation.createRequestHeader();
        CloseSecureChannelRequest closeSecureChannelRequest = new CloseSecureChannelRequest(requestHeader);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false, false,
            new NodeIdFourByte((short) 0, Integer.parseInt(closeSecureChannelRequest.getIdentifier())),
            null, null
        );

        Function<CallContext, OpcuaCloseRequest> closeRequest = ctx ->
            new OpcuaCloseRequest(FINAL, ctx.getSecurityHeader(),
            new ExtensiblePayload(
                new SequenceHeader(ctx.getNextSequenceNumber(), ctx.getRequestId()),
                new ExtensionObject(expandedNodeId, null, closeSecureChannelRequest)
            )
        );

        conversation.requestChannelClose(closeRequest);
    }

    public CompletableFuture<EndpointDescription> onDiscover() {
        // Only the TCP transport supports login.
        LOGGER.debug("Opcua Driver running in ACTIVE mode, discovering endpoints");

        return conversation.requestHello()
            .thenCompose(ack -> onConnectOpenSecureChannel(SecurityTokenRequestType.securityTokenRequestTypeIssue, 0, 0))
            .thenCompose(scr -> onDiscoverGetEndpointsRequest())
            .thenApply(endpoint -> {
                LOGGER.info("Finished discovery of communication endpoint");
                return endpoint;
            });
    }

    public CompletableFuture<EndpointDescription> onDiscoverGetEndpointsRequest() {
        RequestHeader requestHeader = conversation.createRequestHeader();

        GetEndpointsRequest endpointsRequest = new GetEndpointsRequest(
            requestHeader,
            this.endpoint,
            0,
            null,
            0,
            null
        );

        return conversation.submit(endpointsRequest, GetEndpointsResponse.class).thenApply(response -> {
            List<ExtensionObjectDefinition> endpoints = response.getEndpoints();
            MessageSecurityMode effectiveMode = this.configuration.getSecurityPolicy() == SecurityPolicy.NONE ? MessageSecurityMode.messageSecurityModeNone : this.configuration.getMessageSecurity().getMode();
            for (ExtensionObjectDefinition endpoint : endpoints) {
                EndpointDescription endpointDescription = (EndpointDescription) endpoint;

                boolean urlMatch = endpointDescription.getEndpointUrl().getStringValue().equals(this.endpoint.getStringValue());
                boolean policyMatch = endpointDescription.getSecurityPolicyUri().getStringValue().equals(this.configuration.getSecurityPolicy().getSecurityPolicyUri());
                boolean msgSecurityMatch = endpointDescription.getSecurityMode().equals(effectiveMode);

                LOGGER.debug("Validate OPC UA endpoint {} during discovery phase."
                    + "Expected {}. Endpoint policy {} looking for {}. Message security {}, looking for {}", endpointDescription.getEndpointUrl().getStringValue(), this.endpoint.getStringValue(),
                    endpointDescription.getSecurityPolicyUri().getStringValue(), configuration.getSecurityPolicy().getSecurityPolicyUri(),
                    endpointDescription.getSecurityMode(), configuration.getMessageSecurity().getMode());

                if (urlMatch && policyMatch && msgSecurityMatch) {
                   LOGGER.info("Found OPC UA endpoint {}", this.endpoint.getStringValue());
                   return endpointDescription;
                }
            }

            throw new IllegalArgumentException("Could not find endpoint matching client configuration. Tested " + endpoints.size() + " endpoints. "
                + "None matched " + this.endpoint.getStringValue() + " " + this.configuration.getSecurityPolicy().getSecurityPolicyUri() + " " + this.configuration.getMessageSecurity().getMode());
        });
    }

    private OpenSecureChannelResponse onOpenResponse(OpcuaOpenResponse opcuaOpenResponse) {
        try {
            ReadBuffer readBuffer = toBuffer(opcuaOpenResponse::getMessage);
            ExtensionObject message = ExtensionObject.staticParse(readBuffer, false);

            if (message.getBody() instanceof ServiceFault) {
                ServiceFault fault = (ServiceFault) message.getBody();
                throw new PlcRuntimeException(Conversation.toProtocolException(fault));
            }

            LOGGER.debug("Received valid answer for open secure channel request, forwarding it to call initiator");
            return (OpenSecureChannelResponse) message.getBody();
        } catch (ParseException e) {
            throw new IllegalArgumentException("Could not handle response", e);
        }
    }

    private void renewToken() {
        if (keepAlive != null) {
            // cancel earlier renew feature
            keepAlive.cancel(true);
        }
        long keepAliveTime = (long) Math.ceil(revisedLifetime * 0.75f);
        LOGGER.debug("Scheduling session keep alive to happen within {}s", TimeUnit.MILLISECONDS.toSeconds(keepAliveTime));
        keepAlive = KEEP_ALIVE_EXECUTOR.scheduleAtFixedRate(() -> {
            RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> {
                int securityChannelId = this.conversation.getSecurityChannelId();
                int requestId = this.conversation.getRequestId();
                onConnectOpenSecureChannel(SecurityTokenRequestType.securityTokenRequestTypeRenew, securityChannelId, requestId)
                    .whenComplete((response, error) -> {
                        if (error != null) {
                            transaction.failRequest(error);
                            return;
                        }
                        // make sure we still honor channel lifetime boundary
                        long newKeepAliveTime = (long) Math.ceil(revisedLifetime * 0.75f);
                        if (newKeepAliveTime != keepAliveTime) {
                            renewToken();
                        }
                        transaction.endRequest();

                    });
            });
        }, keepAliveTime, keepAliveTime, TimeUnit.MILLISECONDS);
    }

    private static ReadBufferByteBased toBuffer(Supplier<Payload> supplier) {
        Payload payload = supplier.get();
        if (!(payload instanceof BinaryPayload)) {
            throw new IllegalArgumentException("Unexpected payload kind");
        }
        return new ReadBufferByteBased(((BinaryPayload) payload).getPayload(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
    }

    /**
     * Selects the endpoint to use based on the connection string provided.
     * If Discovery is disabled it will use the host address return from the server
     *
     * @param sessionResponse - The CreateSessionResponse message returned by the server
     * @throws PlcRuntimeException - If no endpoint with a compatible policy is found raise and error.
     */
    private void selectEndpoint(CreateSessionResponse sessionResponse) throws PlcRuntimeException {
        // Get a list of the endpoints which match ours.
        Stream<EndpointDescription> filteredEndpoints = sessionResponse.getServerEndpoints().stream()
            .map(e -> (EndpointDescription) e)
            .filter(this::isEndpoint);

        //Determine if the requested security policy is included in the endpoint
        filteredEndpoints.forEach(endpoint -> hasIdentity(
            endpoint.getUserIdentityTokens().stream()
                .map(p -> (UserTokenPolicy) p)
                .toArray(UserTokenPolicy[]::new)
        ));

        if (this.policyId == null) {
            throw new PlcRuntimeException("Unable to find endpoint - " + this.endpoints.get(0));
        }
        if (this.tokenType == null) {
            throw new PlcRuntimeException("Unable to find Security Policy for endpoint - " + this.endpoints.get(0));
        }
    }

    /**
     * Checks each component of the return endpoint description against the connection string.
     * If all are correct then return true.
     *
     * @param endpoint - EndpointDescription returned from server
     * @return true if this endpoint matches our configuration
     * @throws PlcRuntimeException - If the returned endpoint string doesn't match the format expected
     */
    private boolean isEndpoint(EndpointDescription endpoint) throws PlcRuntimeException {
        // Split up the connection string into it's individual segments.
        String endpointUri = endpoint.getEndpointUrl().getStringValue();
        Matcher matcher = URI_PATTERN.matcher(endpointUri);
        if (!matcher.matches()) {
            throw new PlcRuntimeException(
                "Endpoint " + endpointUri + "  returned from the server doesn't match the format '{protocol-code}:({transport-code})?//{transport-host}(:{transport-port})(/{transport-endpoint})'");
        }
        LOGGER.trace("Using Endpoint {} {} {}", matcher.group("transportHost"), matcher.group("transportPort"), matcher.group("transportEndpoint"));

        //When the parameter discovery=false is configured, prefer using the custom address. If the transportEndpoint is empty,
        // directly replace it with the TransportEndpoint returned by the server.
        if (!configuration.isDiscovery() && StringUtils.isBlank(driverContext.getTransportEndpoint())) {
            driverContext.setTransportEndpoint(matcher.group("transportEndpoint"));
            return true;
        }

        if (configuration.isDiscovery() && !this.endpoints.contains(matcher.group("transportHost"))) {
            return false;
        }

        if (!driverContext.getPort().equals(matcher.group("transportPort"))) {
            return false;
        }

        if (!driverContext.getTransportEndpoint().equals(matcher.group("transportEndpoint"))) {
            return false;
        }

        return true;
    }

    /**
     * Confirms that a policy that matches the connection string is available from
     * the returned endpoints. It sets the policyId and tokenType for the policy to use.
     *
     * @param policies - A list of policies returned with the endpoint description.
     */
    private void hasIdentity(UserTokenPolicy[] policies) {
        for (UserTokenPolicy identityToken : policies) {
            if ((identityToken.getTokenType() == UserTokenType.userTokenTypeAnonymous) && (this.username == null)) {
                policyId = identityToken.getPolicyId();
                tokenType = identityToken.getTokenType();
            } else if ((identityToken.getTokenType() == UserTokenType.userTokenTypeUserName) && (this.username != null)) {
                policyId = identityToken.getPolicyId();
                tokenType = identityToken.getTokenType();
            }
        }
    }

    /**
     * Creates an IdentityToken to authenticate with a server.
     *
     * @param tokenType      the token type
     * @param securityPolicy the security policy
     * @return returns an ExtensionObject with an IdentityToken.
     */
    private ExtensionObject getIdentityToken(UserTokenType tokenType, String securityPolicy) {
        ExpandedNodeId extExpandedNodeId;
        switch (tokenType) {
            case userTokenTypeAnonymous:
                //If we aren't using authentication tell the server we would like to log in anonymously
                AnonymousIdentityToken anonymousIdentityToken = new AnonymousIdentityToken();

                extExpandedNodeId = new ExpandedNodeId(
                    false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, OpcuaNodeIdServicesObject.AnonymousIdentityToken_Encoding_DefaultBinary.getValue()),
                    null,
                    null
                );

                return new ExtensionObject(
                    extExpandedNodeId,
                    new ExtensionObjectEncodingMask(false, false, true),
                    new UserIdentityToken(new PascalString(securityPolicy), anonymousIdentityToken));
            case userTokenTypeUserName:
                //Encrypt the password using the server nonce and server public key
                byte[] remoteNonce = conversation.getRemoteNonce();
                byte[] passwordBytes = this.password == null ? new byte[0] : this.password.getBytes();
                ByteBuffer encodeableBuffer = ByteBuffer.allocate(4 + passwordBytes.length + remoteNonce.length);
                encodeableBuffer.order(ByteOrder.LITTLE_ENDIAN);
                encodeableBuffer.putInt(passwordBytes.length + remoteNonce.length);
                encodeableBuffer.put(passwordBytes);
                encodeableBuffer.put(remoteNonce);
                byte[] encodeablePassword = new byte[4 + passwordBytes.length + remoteNonce.length];
                encodeableBuffer.position(0);
                encodeableBuffer.get(encodeablePassword);

                byte[] encryptedPassword = conversation.encryptPassword(encodeablePassword);
                UserNameIdentityToken userNameIdentityToken = new UserNameIdentityToken(
                    new PascalString(this.username),
                    new PascalByteString(encryptedPassword.length, encryptedPassword),
                    new PascalString(PASSWORD_ENCRYPTION_ALGORITHM)
                );

                extExpandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, OpcuaNodeIdServicesObject.UserNameIdentityToken_Encoding_DefaultBinary.getValue()),
                    null,
                    null);

                return new ExtensionObject(
                    extExpandedNodeId,
                    new ExtensionObjectEncodingMask(false, false, true),
                    new UserIdentityToken(new PascalString(securityPolicy), userNameIdentityToken));
        }
        return null;
    }

    public static X509Certificate getX509Certificate(byte[] certificate) {
        try {
            CertificateFactory factory = CertificateFactory.getInstance("X.509");
            return (X509Certificate) factory.generateCertificate(new ByteArrayInputStream(certificate));
        } catch (Exception e) {
            LOGGER.error("Unable to get certificate from String {}", certificate);
            return null;
        }
    }

    private static PascalByteString createPascalString(byte[] bytes) {
        if (null == bytes) {
            return NULL_BYTE_STRING;
        }
        return new PascalByteString(bytes.length, bytes);
    }

}
