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

import static java.util.Map.entry;
import static java.util.concurrent.Executors.newSingleThreadScheduledExecutor;
import static org.apache.plc4x.java.opcua.readwrite.ChunkType.*;

import java.io.ByteArrayInputStream;
import java.security.GeneralSecurityException;
import java.security.Signature;
import java.security.cert.CertificateEncodingException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.Comparator;
import java.util.Map.Entry;
import java.util.Set;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Collectors;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcCertificateAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.SignatureAlgorithm;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Locale;
import java.util.Optional;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.regex.Pattern;

public class SecureChannel {

    private static final Logger LOGGER = LoggerFactory.getLogger(SecureChannel.class);
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
    public static final ExtensionObjectEncodingMask BINARY_ENCODING_MASK = new ExtensionObjectEncodingMask(
        false, false, true);
    private final String sessionName = "UaSession:" + APPLICATION_TEXT.getStringValue() + ":" + RandomStringUtils.random(20, true, true);
    private final PascalByteString localCertificateString;
    private final PascalByteString remoteCertificateThumbprint;
    private final PascalString endpoint;
    private final String username;
    private final String password;
    private final OpcuaConfiguration configuration;
    private final OpcuaDriverContext driverContext;
    /** The user certificate and key of an X509IdentityToken, or null when not authenticating with one. */
    private final CertificateKeyPair userIdentity;
    private final Conversation conversation;
    private ScheduledFuture<?> keepAlive;
    private double sessionTimeout;
    private long revisedLifetime;

    /**
     * Set once the server's lifetime has been raised to the configured minimum, so the warning is
     * emitted once per channel rather than on every renewal - which, by the nature of the
     * condition, would be every few seconds.
     */
    private boolean shortLifetimeWarned;

    public SecureChannel(Conversation conversation, OpcuaDriverContext driverContext, OpcuaConfiguration configuration, PlcAuthentication authentication) {
        this.conversation = conversation;
        this.configuration = configuration;
        this.driverContext = driverContext;
        this.endpoint = new PascalString(driverContext.getEndpoint());
        this.sessionTimeout = configuration.getSessionTimeout();
        if (authentication instanceof PlcUsernamePasswordAuthentication) {
            this.username = ((PlcUsernamePasswordAuthentication) authentication).getUsername();
            this.password = ((PlcUsernamePasswordAuthentication) authentication).getPassword();
            this.userIdentity = null;
        } else if (authentication instanceof PlcCertificateAuthentication certificateAuthentication) {
            // The user identity of an X509IdentityToken. It is deliberately separate from the
            // application instance certificate in driverContext: that one secures the channel and
            // says which installation is talking, this one says who is talking (OPC UA Part 4).
            this.username = null;
            this.password = null;
            this.userIdentity = loadUserIdentity(certificateAuthentication);
        } else if (authentication == null) {
            this.username = configuration.getUsername();
            this.password = configuration.getPassword();
            this.userIdentity = null;
        } else {
            throw new PlcRuntimeException("This type of connection only supports anonymous,"
                + " username-password and certificate authentication, got "
                + authentication.getClass().getSimpleName());
        }

        if (conversation.getSecurityPolicy() == SecurityPolicy.NONE) {
            this.localCertificateString = NULL_BYTE_STRING;
            this.remoteCertificateThumbprint = NULL_BYTE_STRING;
        } else {
            CertificateKeyPair keyPair = driverContext.getCertificateKeyPair();
            this.remoteCertificateThumbprint = driverContext.getThumbprint();
            // Send the whole chain: a CA-signed certificate is not verifiable on its own by a
            // server that only trusts the issuing CA (OPC UA Part 6, SenderCertificate). For a
            // self-signed certificate this is just that one certificate.
            byte[] encoded = keyPair.getEncodedCertificateChain();
            this.localCertificateString = new PascalByteString(encoded.length, encoded);
        }
    }

    private static CertificateKeyPair loadUserIdentity(PlcCertificateAuthentication authentication) {
        try {
            return KeyStoreCredentials.load(authentication.getKeyStore(), authentication.getKeyStorePassword(),
                authentication.getKeyAlias(), "the supplied PlcCertificateAuthentication", "user certificates");
        } catch (GeneralSecurityException e) {
            throw new PlcRuntimeException("Could not read the user certificate to authenticate with", e);
        }
    }

    public CompletableFuture<ActivateSessionResponse> onConnect() {
        // Only the TCP transport supports login.
        LOGGER.debug("Opcua Driver running in ACTIVE mode.");
        return conversation.requestHello()
            .thenCompose(r -> onConnectOpenSecureChannel(SecurityTokenRequestType.securityTokenRequestTypeIssue, 0, 0))
            .thenCompose(r -> onConnectSession());
    }

    /**
     * Establishes the session (CreateSession + ActivateSession) on an already-open secure
     * channel. Used after {@link #onDiscover()}, which has already performed the
     * Hello/OpenSecureChannel exchange: sending a second Hello on the same TCP connection
     * would stall, since Hello is a once-per-connection message.
     */
    public CompletableFuture<ActivateSessionResponse> onConnectSession() {
        return onConnectCreateSessionRequest()
            .thenCompose(this::onConnectActivateSessionRequest)
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
                (long) OpcuaConstants.PROTOCOLVERSION,
                securityTokenRequestType,
                configuration.getMessageSecurity().getMode(),
                new PascalByteString(localNonce.length, localNonce),
                configuration.getChannelLifetime() // lifetime
            );
        } else {
            openSecureChannelRequest = new OpenSecureChannelRequest(
                requestHeader,
                (long) OpcuaConstants.PROTOCOLVERSION,
                securityTokenRequestType,
                MessageSecurityMode.messageSecurityModeNone,
                NULL_BYTE_STRING,
                configuration.getChannelLifetime() // lifetime
            );
        }

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false, false,
            new NodeIdFourByte((short) 0, openSecureChannelRequest.getExtensionId()),
            null, null
        );

        Function<CallContext, OpcuaOpenRequest> openRequest = context -> {
            LOGGER.debug("Submitting OpenSecureChannel with id of {}", context.getRequestId());
            return new OpcuaOpenRequest(FINAL, new OpenChannelMessageRequest(secureChannelId,
                new PascalString(conversation.getSecurityPolicy().getSecurityPolicyUri()),
                this.localCertificateString,
                this.remoteCertificateThumbprint
            ),
            new ExtensiblePayload(
                new SequenceHeader(context.getNextSequenceNumber(), context.getRequestId()),
                new RootExtensionObject(expandedNodeId, openSecureChannelRequest)
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
                LOGGER.debug("Opened secure response id: {}, channel id:{}, token:{} lifetime:{}", openSecureChannelResponse.getExtensionId(),
                    securityToken.getChannelId(), securityToken.getTokenId(), securityToken.getRevisedLifetime());

                // store server and client nonce
                conversation.setRemoteNonce(openSecureChannelResponse.getServerNonce().getStringValue());
                conversation.setLocalNonce(localNonce);
                conversation.setSecurityHeader(new SecurityHeader(securityToken.getChannelId(), securityToken.getTokenId()));
                revisedLifetime = adoptChannelLifetime(securityToken.getRevisedLifetime());
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

        List<PascalString> discoveryUrls = new ArrayList<>(0);

        ApplicationDescription clientDescription = new ApplicationDescription(
            driverContext.getApplicationUri().map(PascalString::new).orElse(APPLICATION_URI),
            PRODUCT_URI,
            applicationName,
            ApplicationType.applicationTypeClient,
            NULL_STRING,
            NULL_STRING,
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

        Entry<EndpointDescription, UserTokenPolicy> selectedEndpoint = selectEndpoint(sessionResponse.getServerEndpoints(),
            configuration.getSecurityPolicy(), configuration.getMessageSecurity());
        if (selectedEndpoint == null) {
            // Endpoints are matched on the host as written in the connection string; the driver
            // performs no name resolution. Servers commonly advertise a different name than the
            // one that was dialed, so say what was offered and how to point the driver at it.
            String advertised = sessionResponse.getServerEndpoints().stream()
                .map(endpoint -> endpoint.getEndpointUrl().getStringValue()
                    + " (" + endpoint.getSecurityPolicyUri().getStringValue()
                    + ", " + endpoint.getSecurityMode()
                    + ", user tokens: " + endpoint.getUserIdentityTokens().stream()
                        .map(policy -> String.valueOf(policy.getTokenType()))
                        .collect(Collectors.joining("/")) + ")")
                .collect(Collectors.joining("\n  "));
            throw new PlcRuntimeException("Unable to find an endpoint matching " + driverContext.getEndpoint()
                + " with security policy " + configuration.getSecurityPolicy()
                + ", message security " + configuration.getMessageSecurity()
                + " and " + requestedTokenType() + " authentication"
                + ". The server offered:\n  " + advertised
                + "\nIf the server advertises a different host name than the one you connect to, set"
                + " 'endpoint-host' (and 'endpoint-port' if it differs) to the name shown above.");
        }

        ActivateSessionRequest activateSessionRequest = createActivateSessionRequest(selectedEndpoint);

        return conversation.submit(activateSessionRequest, ActivateSessionResponse.class).thenApply(responseMessage -> {
            conversation.setRemoteNonce(responseMessage.getServerNonce().getStringValue());
            return responseMessage;
        });
    }

    /**
     * Closes the session and the secure channel on the server. The returned future
     * completes once the {@code CloseSession} has been acknowledged and the
     * {@code CloseSecureChannel} has been handed to the wire, so callers must await it
     * before tearing down the socket — otherwise the server never sees the close, leaks
     * the session/channel, and eventually refuses new channels once its concurrent-channel
     * limit is reached.
     *
     * <p>Note that {@code CloseSecureChannel} is not awaited for a reply: per the OPC UA
     * spec the server simply closes the channel without responding, so we only ensure its
     * bytes are flushed (which {@code requestChannelClose} does synchronously).</p>
     */
    public CompletableFuture<Void> onDisconnect() {
        LOGGER.info("Disconnecting");

        if (keepAlive != null) {
            keepAlive.cancel(true);
            keepAlive = null;
        }

        RequestHeader requestHeader = conversation.createRequestHeader(50000L);
        CloseSessionRequest closeSessionRequest = new CloseSessionRequest(requestHeader, true);
        return conversation.submit(closeSessionRequest, CloseSessionResponse.class)
            // Proceed to close the channel even if the session close failed/timed out;
            // the important thing is that we still tell the server to drop the channel.
            .handle((responseMessage, error) -> {
                LOGGER.trace("Got Close Session Response {}", responseMessage);
                return null;
            })
            .thenRun(this::sendCloseSecureChannel);
    }

    private void sendCloseSecureChannel() {
        RequestHeader requestHeader = conversation.createRequestHeader();
        CloseSecureChannelRequest closeSecureChannelRequest = new CloseSecureChannelRequest(requestHeader);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false, false,
            new NodeIdFourByte((short) 0, closeSecureChannelRequest.getExtensionId()),
            null, null
        );

        Function<CallContext, OpcuaCloseRequest> closeRequest = ctx ->
            new OpcuaCloseRequest(FINAL, ctx.getSecurityHeader(),
            new ExtensiblePayload(
                new SequenceHeader(ctx.getNextSequenceNumber(), ctx.getRequestId()),
                new RootExtensionObject(expandedNodeId, closeSecureChannelRequest)
            )
        );

        // Fire-and-forget: the bytes are flushed synchronously; no response is expected.
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
            null,
            null
        );

        return conversation.submit(endpointsRequest, GetEndpointsResponse.class).thenApply(response -> {
            Entry<EndpointDescription, UserTokenPolicy> entry = selectEndpoint(response.getEndpoints(),
                this.configuration.getSecurityPolicy(), this.configuration.getMessageSecurity());

            if (entry == null) {
                Set<String> endpointUris = response.getEndpoints().stream()
                    .map(EndpointDescription::getEndpointUrl)
                    .map(PascalString::getStringValue)
                    .collect(Collectors.toSet());
                throw new IllegalArgumentException("Could not find endpoint matching client configuration. Tested " + endpointUris + ". "
                    + "Was looking for " + this.endpoint.getStringValue() + " " + this.configuration.getSecurityPolicy().getSecurityPolicyUri() + " " + this.configuration.getMessageSecurity().getMode());
            }
            return entry.getKey();
        });
    }

    private OpenSecureChannelResponse onOpenResponse(OpcuaOpenResponse opcuaOpenResponse) {
        try {
            ReadBufferByteBased readBuffer = toBuffer(opcuaOpenResponse::getMessage);
            ExtensionObject message = ExtensionObject.staticParse(readBuffer, false);

            if (message.getBody() instanceof ServiceFault fault) {
                throw new PlcRuntimeException(Conversation.toProtocolException(fault));
            }

            LOGGER.debug("Received valid answer for open secure channel request, forwarding it to call initiator");
            return (OpenSecureChannelResponse) message.getBody();
        } catch (BufferException e) {
            throw new IllegalArgumentException("Could not handle response", e);
        }
    }

    /**
     * Apply {@link #effectiveChannelLifetime(long, long, long)} to a server-supplied lifetime and
     * tell the operator when their configured minimum overruled the server - including what to
     * change if they would rather honour it.
     */
    private long adoptChannelLifetime(long revisedLifetime) {
        long requested = configuration.getChannelLifetime();
        long minimum = configuration.getMinChannelLifetime();
        long effective = effectiveChannelLifetime(revisedLifetime, requested, minimum);

        if (revisedLifetime > 0 && revisedLifetime < effective && !shortLifetimeWarned) {
            shortLifetimeWarned = true;
            LOGGER.warn("Server asked for a secure channel lifetime of {} ms; using {} ms instead, "
                    + "because min-channel-lifetime is {} ms. Renewals share one executor with every "
                    + "OPC UA connection in this JVM, which is what that minimum protects. The server "
                    + "may treat the channel as expired before the first renewal - if this server "
                    + "genuinely needs renewal that often, lower min-channel-lifetime to {} or less.",
                revisedLifetime, effective, minimum, revisedLifetime);
        }
        return effective;
    }

    /**
     * Reconcile the lifetime the server came back with against the one we asked for.
     *
     * <p>Two different situations, deliberately handled differently:</p>
     * <ul>
     *   <li><strong>Not a lifetime at all</strong> - zero, negative, or longer than we offered.
     *       OPC UA lets a server revise the requested lifetime <em>downwards</em>; none of these
     *       is a revision, so what we requested stands. Nothing is lost by ignoring them.</li>
     *   <li><strong>A lifetime shorter than {@code minimumLifetime}</strong> - a real answer, but
     *       one that would put the renewal schedule on a treadmill. It is raised to the minimum
     *       and the caller warns. Note the consequence: the server considers the channel expired
     *       before our first renewal is due, so the connection may fail at that point. That is the
     *       trade being made - the renewals run on an executor shared by every OPC UA connection
     *       in the JVM, so one peer does not get to set the pace for all of them. An operator who
     *       needs such a server lowers {@code min-channel-lifetime} and accepts the cost
     *       knowingly.</li>
     * </ul>
     *
     * <p>The minimum is bounded by the requested lifetime, so a deliberately short
     * {@code channel-lifetime} is still honoured: this only ever declines to go <em>below</em>
     * what the operator asked for, never above it.</p>
     */
    static long effectiveChannelLifetime(long revisedLifetime, long requestedLifetime, long minimumLifetime) {
        if (revisedLifetime <= 0 || revisedLifetime > requestedLifetime) {
            return requestedLifetime;
        }
        return Math.max(revisedLifetime, Math.min(minimumLifetime, requestedLifetime));
    }

    /**
     * Renewal interval for a channel lifetime: three quarters of it, leaving a quarter of the
     * lifetime as margin for the renewal exchange itself. Never returns a non-positive period -
     * {@code scheduleAtFixedRate} rejects those, and it would do so from inside a completion stage
     * where the failure is easy to lose.
     */
    static long keepAliveInterval(long channelLifetime) {
        return Math.max(1L, (long) Math.ceil(channelLifetime * 0.75f));
    }

    private void renewToken() {
        if (keepAlive != null) {
            // cancel earlier renew feature
            keepAlive.cancel(true);
        }
        long keepAliveTime = keepAliveInterval(revisedLifetime);
        LOGGER.debug("Scheduling session keep alive to happen within {}s", TimeUnit.MILLISECONDS.toSeconds(keepAliveTime));
        keepAlive = KEEP_ALIVE_EXECUTOR.scheduleAtFixedRate(() -> {
            int securityChannelId = this.conversation.getSecurityChannelId();
            int requestId = this.conversation.getRequestId();
            onConnectOpenSecureChannel(SecurityTokenRequestType.securityTokenRequestTypeRenew, securityChannelId, requestId)
                .whenComplete((response, error) -> {
                    if (error != null) {
                        LOGGER.warn("Token renewal failed", error);
                        return;
                    }
                    // Honor any new lifetime the server gave us — if it differs
                    // from what's currently scheduled, reschedule the next renew.
                    long newKeepAliveTime = keepAliveInterval(revisedLifetime);
                    if (newKeepAliveTime != keepAliveTime) {
                        renewToken();
                    }
                });
        }, keepAliveTime, keepAliveTime, TimeUnit.MILLISECONDS);
    }

    private static ReadBufferByteBased toBuffer(Supplier<Payload> supplier) {
        Payload payload = supplier.get();
        if (!(payload instanceof BinaryPayload)) {
            throw new IllegalArgumentException("Unexpected payload kind");
        }
        return new ReadBufferByteBased(((BinaryPayload) payload).getPayload(), PayloadConverter.LITTLE_ENDIAN);
    }

    /**
     * Selects the endpoint and authentication policy based on client settings.
     *
     * @param extensionObjects Endpoint descriptions returned by the server.
     * @param securityPolicy Security policy searched in endpoints.
     * @param messageSecurity Message security needed by client.
     * @return Endpoint matching given.
     */
    private Entry<EndpointDescription, UserTokenPolicy> selectEndpoint(List<EndpointDescription> extensionObjects,
        SecurityPolicy securityPolicy, MessageSecurity messageSecurity) throws PlcRuntimeException {
        // Get a list of the endpoints which match ours.
        MessageSecurityMode effectiveMessageSecurity = SecurityPolicy.NONE == securityPolicy ? MessageSecurityMode.messageSecurityModeNone : messageSecurity.getMode();
        List<Entry<EndpointDescription, UserTokenPolicy>> serverEndpoints = new ArrayList<>();

        for (EndpointDescription endpointDescription : extensionObjects) {
            if (isMatchingEndpointDescription(endpointDescription)) {
                boolean policyMatch = endpointDescription.getSecurityPolicyUri().getStringValue().equals(securityPolicy.getSecurityPolicyUri());
                boolean msgSecurityMatch = endpointDescription.getSecurityMode().equals(effectiveMessageSecurity);

                if (!policyMatch && !msgSecurityMatch) {
                    continue;
                }

                for (UserTokenPolicy userTokenPolicy : endpointDescription.getUserIdentityTokens()) {
                    if (isUserTokenPolicyCompatible(userTokenPolicy, this.username, this.userIdentity != null)) {
                        serverEndpoints.add(entry(endpointDescription, userTokenPolicy));
                    }
                }
            }
        }

        if (serverEndpoints.isEmpty()) {
            return null;
        }

        serverEndpoints.sort(Comparator.comparing(e -> e.getKey().getSecurityLevel()));
        return serverEndpoints.getFirst();
    }

    private boolean isMatchingEndpointDescription(EndpointDescription endpointDescription) {
        if (isMatchingEndpoint(endpointDescription, driverContext.getHost(), driverContext.getPort(), driverContext.getTransportEndpoint())) {
            return true;
        }
        if (configuration.getEndpointHost() != null) {
            return isMatchingEndpoint(endpointDescription, configuration.getEndpointHost(), configuration.getEndpointPort() == null ? driverContext.getPort() : String.valueOf(configuration.getEndpointPort()), driverContext.getTransportEndpoint());
        } else if (configuration.getEndpointPort() != null) {
            return isMatchingEndpoint(endpointDescription, driverContext.getHost(), configuration.getEndpointPort().toString(), driverContext.getTransportEndpoint());
        }
        return false;
    }

    /**
     * Checks each component of the return endpoint description against the connection string.
     * If all are correct then return true.
     *
     * @param endpoint - EndpointDescription returned from server
     * @param host Permitted host
     * @param port Permitted port
     * @param transportEndpoint Transport endpoint
     * @return true if this endpoint matches our configuration
     * @throws PlcRuntimeException - If the returned endpoint string doesn't match the format expected
     */
    private static boolean isMatchingEndpoint(EndpointDescription endpoint, String host, String port, String transportEndpoint) throws PlcRuntimeException {
        String portAddition = port == null ? "" : ":" + port;
        String expected = "opc.tcp://" + host + portAddition + transportEndpoint;
        // Host names are case insensitive, so a server advertising "opc.tcp://MyServer:4840" has
        // to match a connection string that says "myserver".
        return endpoint.getEndpointUrl().getStringValue().toLowerCase(Locale.ROOT)
            .startsWith(expected.toLowerCase(Locale.ROOT));
    }

    /**
     * The token type the supplied credentials ask for, named in the error when no endpoint offers it.
     */
    private String requestedTokenType() {
        if (userIdentity != null) {
            return "certificate";
        }
        return username != null ? "username" : "anonymous";
    }

    /**
     * Confirms that the given policy matches the credentials the client was given. Which of the
     * three token types is wanted follows from those credentials alone: a user certificate asks for
     * a certificate policy, a user name for a username policy, and neither for anonymous access.
     *
     * @param policy                 UserTokenPolicy configured for server endpoint.
     * @param username               the user name to authenticate with, or null for none.
     * @param userCertificatePresent whether a user certificate was supplied to authenticate with.
     * @return True if given token policy matches client configuration.
     */
    static boolean isUserTokenPolicyCompatible(UserTokenPolicy policy, String username,
        boolean userCertificatePresent) {
        if (userCertificatePresent) {
            return policy.getTokenType() == UserTokenType.userTokenTypeCertificate;
        }
        if (username != null) {
            return policy.getTokenType() == UserTokenType.userTokenTypeUserName;
        }
        return policy.getTokenType() == UserTokenType.userTokenTypeAnonymous;
    }

    /**
     * Builds the {@code ActivateSession} request for the user token policy that was selected.
     */
    ActivateSessionRequest createActivateSessionRequest(Entry<EndpointDescription, UserTokenPolicy> selectedEndpoint) {
        PascalString policyId = selectedEndpoint.getValue().getPolicyId();
        UserTokenType tokenType = selectedEndpoint.getValue().getTokenType();
        SecurityPolicy tokenSecurityPolicy = userTokenSecurityPolicy(selectedEndpoint);
        ExtensionObject userIdentityToken = getIdentityToken(tokenType, policyId.getStringValue(), tokenSecurityPolicy);
        RequestHeader requestHeader = conversation.createRequestHeader();

        SignatureData clientSignature = new SignatureData(NULL_STRING, NULL_BYTE_STRING);
        if (conversation.getSecurityPolicy() != SecurityPolicy.NONE) {
            try {
                clientSignature = conversation.createClientSignature();
            } catch (GeneralSecurityException e) {
                throw new PlcRuntimeException("Could not create client signature", e);
            }
        }

        return new ActivateSessionRequest(
            requestHeader,
            clientSignature,
            null,
            null,
            userIdentityToken,
            userTokenSignature(tokenType, tokenSecurityPolicy, clientSignature)
        );
    }

    /**
     * An X509IdentityToken is only accepted together with proof that the client holds the private
     * key of the certificate it just sent (OPC UA Part 4, 5.6.3). The other token types carry no
     * such proof, and keep sending the application instance signature the driver has always sent.
     */
    private SignatureData userTokenSignature(UserTokenType tokenType, SecurityPolicy tokenSecurityPolicy,
                                             SignatureData clientSignature) {
        if (tokenType != UserTokenType.userTokenTypeCertificate) {
            return clientSignature;
        }
        if (tokenSecurityPolicy == SecurityPolicy.NONE) {
            throw new PlcRuntimeException("The server offers certificate authentication with a security policy"
                + " of None, which leaves no algorithm to prove possession of the user certificate's private key."
                + " Connect with a security policy other than NONE, or authenticate with a username instead.");
        }
        try {
            return conversation.createUserTokenSignature(userIdentity.getPrivateKey(), tokenSecurityPolicy);
        } catch (GeneralSecurityException e) {
            throw new PlcRuntimeException("Could not sign the user certificate identity token", e);
        }
    }

    /**
     * The security policy governing the user identity token. A UserTokenPolicy may name its own
     * securityPolicyUri; when it doesn't, the endpoint's own policy applies.
     */
    private SecurityPolicy userTokenSecurityPolicy(Entry<EndpointDescription, UserTokenPolicy> selectedEndpoint) {
        PascalString tokenPolicyUri = selectedEndpoint.getValue().getSecurityPolicyUri();
        if (tokenPolicyUri != null && tokenPolicyUri.getStringValue() != null
            && !tokenPolicyUri.getStringValue().isEmpty()) {
            Optional<SecurityPolicy> policy = SecurityPolicy.findByUri(tokenPolicyUri.getStringValue());
            if (policy.isPresent()) {
                LOGGER.debug("User token policy declares security policy {}", policy.get());
                return policy.get();
            }
            LOGGER.warn("Unknown user token security policy '{}', falling back to the endpoint policy",
                tokenPolicyUri.getStringValue());
        }
        return SecurityPolicy.findByUri(selectedEndpoint.getKey().getSecurityPolicyUri().getStringValue())
            .orElseGet(conversation::getSecurityPolicy);
    }

    /**
     * Creates an IdentityToken to authenticate with a server.
     *
     * @param tokenType      the token type
     * @param securityPolicy the security policy
     * @return returns an ExtensionObject with an IdentityToken.
     * <p>
     * Builds the user identity token for the selected user token policy.
     * <p>
     * The password encryption is dictated by the security policy of the <em>user token policy</em>,
     * which is not necessarily the one securing the channel: a token policy may name its own, and
     * only when it doesn't, does the endpoint's policy apply (OPC UA Part 4, UserNameIdentityToken).
     * A policy of None means the password travels in plain text - see GH-2154.
     */
    private ExtensionObject getIdentityToken(UserTokenType tokenType, String securityPolicy,
                                             SecurityPolicy tokenSecurityPolicy) {
        ExpandedNodeId extExpandedNodeId;
        switch (tokenType) {
            case userTokenTypeAnonymous:
                //If we aren't using authentication tell the server we would like to log in anonymously
                AnonymousIdentityToken anonymousIdentityToken = new AnonymousIdentityToken(new PascalString(securityPolicy));

                extExpandedNodeId = new ExpandedNodeId(
                    false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, anonymousIdentityToken.getExtensionId()),
                    null,
                    null
                );

                return new BinaryExtensionObjectWithMask(extExpandedNodeId, BINARY_ENCODING_MASK, anonymousIdentityToken);
            case userTokenTypeUserName:
                //Encrypt the password using the server nonce and server public key
                byte[] remoteNonce = conversation.getRemoteNonce();
                byte[] passwordBytes = this.password == null ? new byte[0] : this.password.getBytes(StandardCharsets.UTF_8);
                ByteBuffer encodeableBuffer = ByteBuffer.allocate(4 + passwordBytes.length + remoteNonce.length);
                encodeableBuffer.order(java.nio.ByteOrder.LITTLE_ENDIAN);
                encodeableBuffer.putInt(passwordBytes.length + remoteNonce.length);
                encodeableBuffer.put(passwordBytes);
                encodeableBuffer.put(remoteNonce);
                byte[] encodeablePassword = new byte[4 + passwordBytes.length + remoteNonce.length];
                encodeableBuffer.position(0);
                encodeableBuffer.get(encodeablePassword);

                byte[] tokenPassword;
                String encryptionAlgorithm;
                if (tokenSecurityPolicy == SecurityPolicy.NONE) {
                    // No encryption: the password is sent as-is and no algorithm is declared.
                    tokenPassword = encodeablePassword;
                    encryptionAlgorithm = "";
                } else {
                    tokenPassword = conversation.encryptPassword(encodeablePassword, tokenSecurityPolicy);
                    encryptionAlgorithm = tokenSecurityPolicy.getAsymmetricEncryptionAlgorithm().getUri();
                }
                UserNameIdentityToken userNameIdentityToken = new UserNameIdentityToken(
                    new PascalString(securityPolicy),
                    new PascalString(this.username),
                    new PascalByteString(tokenPassword.length, tokenPassword),
                    new PascalString(encryptionAlgorithm)
                );

                extExpandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, userNameIdentityToken.getExtensionId()),
                    null,
                    null);

                return new BinaryExtensionObjectWithMask(extExpandedNodeId, BINARY_ENCODING_MASK, userNameIdentityToken);
            case userTokenTypeCertificate:
                // Only the user certificate itself travels here - OPC UA Part 4 defines
                // certificateData as the certificate, not the chain that signed it. Possession of
                // the matching private key is proven by the userTokenSignature of ActivateSession.
                byte[] userCertificate;
                try {
                    userCertificate = userIdentity.getCertificate().getEncoded();
                } catch (CertificateEncodingException e) {
                    throw new PlcRuntimeException("Could not encode the user certificate to authenticate with", e);
                }
                X509IdentityToken x509IdentityToken = new X509IdentityToken(
                    new PascalString(securityPolicy),
                    new PascalByteString(userCertificate.length, userCertificate)
                );

                extExpandedNodeId = new ExpandedNodeId(
                    false,           //Namespace Uri Specified
                    false,            //Server Index Specified
                    new NodeIdFourByte((short) 0, x509IdentityToken.getExtensionId()),
                    null,
                    null
                );

                return new BinaryExtensionObjectWithMask(extExpandedNodeId, BINARY_ENCODING_MASK, x509IdentityToken);
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
