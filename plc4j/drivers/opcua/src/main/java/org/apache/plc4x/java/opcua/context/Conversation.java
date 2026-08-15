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

import static org.apache.plc4x.java.opcua.readwrite.ChunkType.ABORT;
import static org.apache.plc4x.java.opcua.readwrite.ChunkType.FINAL;

import java.security.GeneralSecurityException;
import java.security.cert.X509Certificate;
import java.time.Duration;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.BiFunction;
import java.util.function.BiPredicate;
import java.util.function.Function;
import org.apache.commons.lang3.RandomUtils;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.opcua.config.Limits;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.protocol.chunk.ChunkStorage;
import org.apache.plc4x.java.opcua.protocol.chunk.MemoryChunkStorage;
import org.apache.plc4x.java.opcua.readwrite.BinaryPayload;
import org.apache.plc4x.java.opcua.readwrite.ChunkType;
import org.apache.plc4x.java.opcua.readwrite.ExpandedNodeId;
import org.apache.plc4x.java.opcua.readwrite.ExtensiblePayload;
import org.apache.plc4x.java.opcua.readwrite.ExtensionObject;
import org.apache.plc4x.java.opcua.readwrite.ExtensionObjectDefinition;
import org.apache.plc4x.java.opcua.readwrite.ExtensionObjectEncodingMask;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.NodeId;
import org.apache.plc4x.java.opcua.readwrite.NodeIdFourByte;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTwoByte;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTypeDefinition;
import org.apache.plc4x.java.opcua.readwrite.NullExtensionObjectWithMask;
import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaAcknowledgeResponse;
import org.apache.plc4x.java.opcua.readwrite.OpcuaCloseRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaConstants;
import org.apache.plc4x.java.opcua.readwrite.OpcuaHelloRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaMessageRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaMessageResponse;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenResponse;
import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.readwrite.OpcuaStatusCode;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.apache.plc4x.java.opcua.readwrite.Payload;
import org.apache.plc4x.java.opcua.readwrite.RequestHeader;
import org.apache.plc4x.java.opcua.readwrite.ResponseHeader;
import org.apache.plc4x.java.opcua.readwrite.RootExtensionObject;
import org.apache.plc4x.java.opcua.readwrite.SecurityHeader;
import org.apache.plc4x.java.opcua.readwrite.SequenceHeader;
import org.apache.plc4x.java.opcua.readwrite.ServiceFault;
import org.apache.plc4x.java.opcua.readwrite.SignatureData;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Owns the per-connection OPC UA conversation state: the security header /
 * token id, the sequence number cursor, and the encryption handler. Provides
 * three primitives on top of an {@link OpcuaWire}:
 *
 * <ul>
 *   <li>{@link #requestHello()} — the initial Hello/Acknowledge exchange</li>
 *   <li>{@link #requestChannelOpen(Function)} / {@link #requestChannelClose(Function)} —
 *       the {@code OPN}/{@code CLO} exchange that sets up / tears down the
 *       secure channel</li>
 *   <li>{@link #submit(ExtensionObjectDefinition, Class)} — every other service
 *       call (CreateSession, ActivateSession, Read, Write, Browse, Publish, ...)</li>
 * </ul>
 *
 * <p>Ported from the old SPI's fluent
 * {@code context.sendRequest(...).expectResponse(...).check(...).unwrap(...).handle(...)}
 * pattern to the new SPI's listener+future model: each request registers a
 * multi-fire listener (so multi-chunk responses can be accumulated) and
 * sends out its chunks, then unsubscribes once the FINAL chunk lands.</p>
 */
public class Conversation implements SecureChannelState {
    private static final long EPOCH_OFFSET = 116444736000000000L;         //Offset between OPC UA epoch time and linux epoch time.

    private static final ExpandedNodeId NULL_EXPANDED_NODE_ID = new ExpandedNodeId(false,
        false,
        new NodeIdTwoByte((short) 0),
        null,
        null
    );

    protected static final ExtensionObject NULL_EXTENSION_OBJECT = new NullExtensionObjectWithMask(
        NULL_EXPANDED_NODE_ID,
        new ExtensionObjectEncodingMask(false, false, false)
    );

    /** Empty PascalString — handy for request-header placeholders. */
    public static final PascalString NULL_STRING = new PascalString("");


    private final Logger logger = LoggerFactory.getLogger(Conversation.class);
    private final AtomicReference<SecurityHeader> securityHeader = new AtomicReference<>(new SecurityHeader(1L, 1L));
    private final AtomicLong senderSequenceNumber = new AtomicLong(-1);

    private final AtomicReference<NodeIdTypeDefinition> authenticationToken = new AtomicReference<>(new NodeIdTwoByte((short) 0));

    private final OpcuaWire wire;
    private final SecureChannelTransactionManager tm;

    private final SecurityPolicy securityPolicy;
    private final MessageSecurity messageSecurity;
    private final EncryptionHandler encryptionHandler;
    private final OpcuaDriverContext driverContext;
    private final OpcuaConfiguration configuration;

    private OpcuaProtocolLimits limits;

    private X509Certificate localCertificate = null;
    private int localCertificateChainSize = 0;
    private X509Certificate remoteCertificate = null;
    private byte[] remoteNonce;
    private byte[] localNonce;

    /**
     * Validates the incoming sequence number for monotonic increase. The first
     * received sequence number anchors the counter; every subsequent number
     * must equal the previous one plus one.
     */
    private final BiPredicate<SequenceHeader, CompletableFuture<?>> sequenceValidator = (sequenceHeader, callback) -> {
        if (senderSequenceNumber.get() == -1L) {
            senderSequenceNumber.set(sequenceHeader.getSequenceNumber());
            return true;
        }
        int expectedSequence = sequenceHeader.getSequenceNumber() - 1;
        if (!senderSequenceNumber.compareAndSet(expectedSequence, sequenceHeader.getSequenceNumber())) {
            callback.completeExceptionally(
                new PlcProtocolException("Lost sequence, expected " + expectedSequence + " but received " + sequenceHeader.getSequenceNumber())
            );
            return false;
        }
        return true;
    };

    public Conversation(OpcuaWire wire, OpcuaDriverContext driverContext, OpcuaConfiguration configuration) {
        this.wire = wire;
        this.tm = new SecureChannelTransactionManager();
        this.driverContext = driverContext;
        this.configuration = configuration;

        this.securityPolicy = determineSecurityPolicy(configuration);
        CertificateKeyPair senderKeyPair = driverContext.getCertificateKeyPair();

        if (this.securityPolicy != SecurityPolicy.NONE) {
            // The remote certificate gets populated during the 'discover' phase
            // when encryption is enabled.
            this.messageSecurity = configuration.getMessageSecurity();
            this.remoteCertificate = configuration.getServerCertificate();
            this.encryptionHandler = new EncryptionHandler(this, senderKeyPair.getPrivateKey());
            this.localCertificate = senderKeyPair.getCertificate();
            // The header is sized against what actually goes on the wire, which for a CA-signed
            // certificate is the whole chain rather than the certificate alone.
            this.localCertificateChainSize = senderKeyPair.getEncodedCertificateChain().length;
        } else {
            this.messageSecurity = MessageSecurity.NONE;
            this.encryptionHandler = new EncryptionHandler(this, null);
        }

        Limits encodingLimits = configuration.getEncodingLimits();
        limits = new OpcuaProtocolLimits(
            (long) encodingLimits.getReceiveBufferSize(),
            (long) encodingLimits.getSendBufferSize(),
            (long) encodingLimits.getMaxMessageSize(),
            (long) encodingLimits.getMaxChunkCount()
        );
    }

    public CompletableFuture<OpcuaAcknowledgeResponse> requestHello() {
        logger.debug("Sending hello message to {}", this.driverContext.getEndpoint());
        OpcuaHelloRequest request = new OpcuaHelloRequest(FINAL,
            (long) OpcuaConstants.PROTOCOLVERSION,
            new OpcuaProtocolLimits(
                limits.getReceiveBufferSize(),
                limits.getSendBufferSize(),
                limits.getMaxMessageSize(),
                limits.getMaxChunkCount()
            ),
            new PascalString(driverContext.getEndpoint())
        );

        CompletableFuture<OpcuaAcknowledgeResponse> future = new CompletableFuture<>();

        // Hello/Acknowledge is a simple single-frame exchange — no chunking,
        // no encryption — so the one-shot {@link OpcuaWire#expect} primitive
        // is enough.
        wire.expect(apu -> apu.getMessage() instanceof OpcuaAcknowledgeResponse,
                Duration.ofMillis(configuration.getNegotiationTimeout()))
            .whenComplete((apu, error) -> {
                if (error != null) {
                    future.completeExceptionally(error);
                    return;
                }
                OpcuaAcknowledgeResponse opcuaAcknowledgeResponse = (OpcuaAcknowledgeResponse) apu.getMessage();
                OpcuaProtocolLimits ackLimits = opcuaAcknowledgeResponse.getLimits();
                // Merge encoding limits to match common minimum:
                //  - our receive buffer must not exceed the server's send buffer
                //  - our send buffer must not exceed the server's receive buffer
                //  - chunks + message sizes negotiate down to the smaller of the two
                this.limits = new OpcuaProtocolLimits(
                    Math.min(this.limits.getReceiveBufferSize(), ackLimits.getSendBufferSize()),
                    Math.min(this.limits.getSendBufferSize(), ackLimits.getReceiveBufferSize()),
                    Math.min(this.limits.getMaxMessageSize(), ackLimits.getMaxMessageSize()),
                    Math.min(this.limits.getMaxChunkCount(), ackLimits.getMaxChunkCount())
                );
                future.complete(opcuaAcknowledgeResponse);
            });

        wire.sendToWire(new OpcuaAPU(request));
        return future;
    }

    public CompletableFuture<OpcuaOpenResponse> requestChannelOpen(Function<CallContext, OpcuaOpenRequest> request) {
        return request(
            OpcuaOpenResponse.class, request,
            (rsp, chunk) -> new OpcuaOpenResponse(rsp.getChunk(), rsp.getOpenResponse(), chunk),
            (rsp) -> rsp.getMessage().getSequenceHeader(),
            OpcuaOpenResponse::getMessage
        );
    }

    public CompletableFuture<Void> requestChannelClose(Function<CallContext, OpcuaCloseRequest> request) {
        logger.trace("Got close secure channel request");
        return request(
            OpcuaMessageResponse.class, request,
            (rsp, chunk) -> new OpcuaMessageResponse(rsp.getChunk(), rsp.getSecurityHeader(), chunk),
            (rsp) -> rsp.getMessage().getSequenceHeader(),
            OpcuaMessageResponse::getMessage
        ).whenComplete((r, e) -> {
            wire.fireDisconnected();
        }).thenApply(r -> null);
    }

    /**
     * Generic request/response that handles encryption + multi-chunk response
     * accumulation. The caller supplies:
     * <ul>
     *   <li>{@code replyType} — the expected response PDU class</li>
     *   <li>{@code request} — given a {@link CallContext}, builds the outgoing PDU</li>
     *   <li>{@code chunkAssembler} — combines a partial reply + accumulated body
     *       bytes into a fully reassembled reply</li>
     *   <li>{@code sequenceHeaderExtractor} / {@code chunkExtractor} — pull the
     *       sequence header and payload out of an incoming chunk</li>
     * </ul>
     *
     * <p>The wire-level flow is: serialise the request, split into chunks via
     * the encryption handler, fire the first N-1 chunks with
     * {@link OpcuaWire#sendToWire}, then for the last chunk both send it and
     * wait for the response.</p>
     */
    private <T extends MessagePDU, R extends MessagePDU> CompletableFuture<R> request(
        Class<R> replyType, Function<CallContext, T> request,
        BiFunction<R, BinaryPayload, R> chunkAssembler,
        Function<R, SequenceHeader> sequenceHeaderExtractor,
        Function<R, Payload> chunkExtractor
    ) {
        int requestId = tm.getTransactionIdentifier();
        logger.debug("Firing request {}", requestId);
        T messagePDU = request.apply(
            new CallContext(securityHeader.get(), tm.getSequenceSupplier(), requestId)
        );

        MemoryChunkStorage chunkStorage = new MemoryChunkStorage();
        List<MessagePDU> chunks = encryptionHandler.encodeMessage(messagePDU, tm.getSequenceSupplier());
        CompletableFuture<R> future = new CompletableFuture<>();

        // Subscribe before sending the final chunk so we don't lose a fast reply.
        // The listener stays registered until either FINAL chunk arrives, an
        // error fires, or the timeout elapses.
        AtomicReference<OpcuaWire.Subscription> subRef = new AtomicReference<>();
        OpcuaWire.Subscription sub = wire.subscribe(
            apu -> {
                MessagePDU msg = apu.getMessage();
                if (!replyType.isInstance(msg)) {
                    return false;
                }
                MessagePDU decoded;
                try {
                    decoded = encryptionHandler.decodeMessage(replyType.cast(msg));
                } catch (Exception e) {
                    // Decryption failure isn't a "skip this packet" — surface it.
                    future.completeExceptionally(e);
                    return true;
                }
                if (!replyType.isInstance(decoded)) {
                    return false;
                }
                R reply = replyType.cast(decoded);
                return requestId == sequenceHeaderExtractor.apply(reply).getRequestId();
            },
            apu -> {
                try {
                    R reply = replyType.cast(encryptionHandler.decodeMessage(replyType.cast(apu.getMessage())));
                    if (!sequenceValidator.test(sequenceHeaderExtractor.apply(reply), future)) {
                        unsubscribeQuietly(subRef);
                        return;
                    }
                    boolean done = accumulateChunkUntilFinal(chunkStorage, reply.getChunk(), chunkExtractor.apply(reply));
                    if (done) {
                        unsubscribeQuietly(subRef);
                        R merged = mergeChunks(chunkStorage, reply, sequenceHeaderExtractor.apply(reply), chunkAssembler);
                        future.complete(merged);
                    }
                } catch (Exception e) {
                    unsubscribeQuietly(subRef);
                    future.completeExceptionally(e);
                }
            });
        subRef.set(sub);

        // If the caller never completes the future, fail it via timeout — and
        // make sure to drop the listener so we don't leak.
        future.orTimeout(configuration.getNegotiationTimeout(), java.util.concurrent.TimeUnit.MILLISECONDS)
            .exceptionally(error -> {
                if (error instanceof TimeoutException) {
                    unsubscribeQuietly(subRef);
                }
                return null;
            });

        for (MessagePDU chunk : chunks) {
            wire.sendToWire(new OpcuaAPU(chunk));
        }
        return future;
    }

    private static void unsubscribeQuietly(AtomicReference<OpcuaWire.Subscription> ref) {
        OpcuaWire.Subscription sub = ref.get();
        if (sub != null) {
            sub.unsubscribe();
        }
    }

    /**
     * Sends an OPC UA service request (Read, Write, Browse, Publish, ...) and
     * returns its response. Wraps the request as a {@link RootExtensionObject}
     * inside an {@link ExtensiblePayload}, applies symmetric encryption per the
     * current security context, and reassembles multi-chunk replies.
     */
    public <T extends ExtensionObjectDefinition, R extends ExtensionObjectDefinition> CompletableFuture<R> submit(T object, Class<R> replyType) {
        return submit(object).thenApply(response -> {
            if (replyType.isInstance(response)) {
                return replyType.cast(response);
            }
            throw new IllegalStateException("Received reply of unexpected type " + response.getClass().getName() + " while " + replyType.getName() + " has been expected");
        });
    }

    private CompletableFuture<Object> submit(ExtensionObjectDefinition requestDefinition) {
        Integer requestId = tm.getTransactionIdentifier();

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(
            false,
            false,
            new NodeIdFourByte((short) 0, requestDefinition.getExtensionId()),
            null,
            null
        );
        ExtensiblePayload payload = new ExtensiblePayload(
            new SequenceHeader(tm.getSequenceSupplier().get(), requestId),
            new RootExtensionObject(expandedNodeId, requestDefinition)
        );

        MemoryChunkStorage chunkStorage = new MemoryChunkStorage();
        SecurityHeader securityHeaderValue = securityHeader.get();
        OpcuaMessageRequest request = new OpcuaMessageRequest(FINAL, securityHeaderValue, payload);

        logger.debug("Submitting Transaction to TransactionManager {}, security channel {}, token {}", requestId,
            securityHeaderValue.getSecureChannelId(), securityHeaderValue.getSecureTokenId());

        List<MessagePDU> chunks = encryptionHandler.encodeMessage(request, tm.getSequenceSupplier());
        CompletableFuture<Object> future = new CompletableFuture<>();

        BiFunction<OpcuaMessageResponse, BinaryPayload, OpcuaMessageResponse> chunkAssembler = (src, chunkPayload) ->
            new OpcuaMessageResponse(src.getChunk(), src.getSecurityHeader(), chunkPayload);

        AtomicReference<OpcuaWire.Subscription> subRef = new AtomicReference<>();
        OpcuaWire.Subscription sub = wire.subscribe(
            apu -> {
                MessagePDU msg = apu.getMessage();
                if (!(msg instanceof OpcuaMessageResponse)) {
                    return false;
                }
                try {
                    MessagePDU decoded = encryptionHandler.decodeMessage((OpcuaMessageResponse) msg);
                    if (!(decoded instanceof OpcuaMessageResponse)) {
                        return false;
                    }
                    return ((OpcuaMessageResponse) decoded).getMessage().getSequenceHeader().getRequestId() == requestId;
                } catch (Exception e) {
                    future.completeExceptionally(e);
                    return true;
                }
            },
            apu -> {
                try {
                    OpcuaMessageResponse response = (OpcuaMessageResponse) encryptionHandler.decodeMessage((OpcuaMessageResponse) apu.getMessage());
                    if (!sequenceValidator.test(response.getMessage().getSequenceHeader(), future)) {
                        unsubscribeQuietly(subRef);
                        return;
                    }
                    boolean done = accumulateChunkUntilFinal(chunkStorage, response.getChunk(), response.getMessage());
                    if (!done) {
                        return;
                    }
                    OpcuaMessageResponse merged = mergeChunks(chunkStorage, response, response.getMessage().getSequenceHeader(), chunkAssembler);
                    unsubscribeQuietly(subRef);

                    if (merged.getChunk().equals(FINAL)) {
                        logger.debug("Received response made of {} bytes for message id: {}, channel id:{}, token:{}",
                            merged.getLengthInBytes(), requestId, merged.getSecurityHeader().getSecureChannelId(),
                            merged.getSecurityHeader().getSecureTokenId()
                        );
                        securityHeader.set(merged.getSecurityHeader());

                        Payload message = merged.getMessage();
                        ExtensionObjectDefinition extensionObjectBody;
                        if (message instanceof ExtensiblePayload) {
                            extensionObjectBody = (((ExtensiblePayload) message).getPayload()).getBody();
                        } else {
                            try {
                                BinaryPayload binary = (BinaryPayload) message;
                                ReadBufferByteBased buffer = new ReadBufferByteBased(binary.getPayload(),
                                    PayloadConverter.LITTLE_ENDIAN);
                                extensionObjectBody = ExtensionObject.staticParse(buffer, false).getBody();
                            } catch (BufferException e) {
                                future.completeExceptionally(e);
                                return;
                            }
                        }

                        if (extensionObjectBody instanceof ServiceFault) {
                            ServiceFault fault = (ServiceFault) extensionObjectBody;
                            future.completeExceptionally(toProtocolException(fault));
                        } else {
                            future.complete(extensionObjectBody);
                        }
                    }
                } catch (Exception e) {
                    unsubscribeQuietly(subRef);
                    future.completeExceptionally(e);
                }
            }
        );
        subRef.set(sub);

        future.orTimeout(configuration.getRequestTimeout(), java.util.concurrent.TimeUnit.MILLISECONDS)
            .exceptionally(error -> {
                if (error instanceof TimeoutException) {
                    unsubscribeQuietly(subRef);
                }
                return null;
            });

        for (MessagePDU chunk : chunks) {
            wire.sendToWire(new OpcuaAPU(chunk));
        }
        return future;
    }

    private <T> T mergeChunks(ChunkStorage chunkStorage, T source, SequenceHeader sequenceHeader, BiFunction<T, BinaryPayload, T> producer) {
        byte[] message = chunkStorage.get();
        return producer.apply(source,
            new BinaryPayload(
                sequenceHeader,
                message
            )
        );
    }

    private boolean accumulateChunkUntilFinal(ChunkStorage storage, ChunkType chunkType, Payload data) {
        if (ABORT.equals(chunkType)) {
            storage.reset();
            return true;
        }

        if (!(data instanceof BinaryPayload)) {
            throw new IllegalArgumentException("Unexpected payload type " + data.getClass());
        }
        storage.append(((BinaryPayload) data).getPayload());

        // Enforce the negotiated receive limits while accumulating chunks. Without this a malicious
        // or faulty server could stream an unbounded number of CONTINUE chunks (never sending a
        // FINAL), growing the in-memory buffer without bound until the client runs out of memory.
        // A limit of 0 means "no limit" per the OPC UA spec, so only enforce positive limits.
        if (limits != null) {
            long maxChunkCount = limits.getMaxChunkCount();
            if (maxChunkCount > 0 && storage.count() > maxChunkCount) {
                storage.reset();
                throw new IllegalStateException("Received message exceeds the negotiated maximum chunk count of " + maxChunkCount);
            }
            long maxMessageSize = limits.getMaxMessageSize();
            if (maxMessageSize > 0 && storage.size() > maxMessageSize) {
                storage.reset();
                throw new IllegalStateException("Received message size " + storage.size()
                    + " exceeds the negotiated maximum message size of " + maxMessageSize + " bytes");
            }
        }

        return FINAL.equals(chunkType);
    }

    public void setLocalNonce(byte[] localNonce) {
        this.localNonce = localNonce;
    }

    // Generate a nonce used for setting up signing/encryption keys.
    byte[] createNonce() {
        return createNonce(securityPolicy.getNonceLength());
    }

    byte[] createNonce(int nonceLength) {
        return RandomUtils.nextBytes(nonceLength);
    }

    @Override
    public boolean isSymmetricEncryptionEnabled() {
        return messageSecurity == MessageSecurity.SIGN_ENCRYPT;
    }

    @Override
    public boolean isSymmetricSigningEnabled() {
        return (messageSecurity == MessageSecurity.SIGN_ENCRYPT || messageSecurity == MessageSecurity.SIGN);
    }

    static SecurityPolicy determineSecurityPolicy(OpcuaConfiguration configuration) {
        if (configuration.isDiscovery() && configuration.getServerCertificate() == null) {
            // Discovery is enabled and the sender certificate isn't known yet,
            // so the discovery phase always runs with security disabled.
            return SecurityPolicy.NONE;
        }

        return configuration.getSecurityPolicy();
    }

    static PlcProtocolException toProtocolException(ServiceFault fault) {
        if (fault.getResponseHeader() instanceof ResponseHeader) {
            ResponseHeader responseHeader = (ResponseHeader) fault.getResponseHeader();
            long statusCode = responseHeader.getServiceResult().getStatusCode();
            String statusName = OpcuaStatusCode.isDefined(statusCode) ? OpcuaStatusCode.enumForValue(statusCode).name() : "<unknown>";
            return new PlcProtocolException("Server returned error " + statusName + " (0x" + Long.toHexString(statusCode) + ")");
        }
        return new PlcProtocolException("Unexpected service fault");
    }

    @Override
    public OpcuaProtocolLimits getLimits() {
        return limits;
    }

    @Override
    public byte[] getLocalNonce() {
        return localNonce;
    }

    @Override
    public X509Certificate getLocalCertificate() {
        return localCertificate;
    }

    /**
     * Number of bytes the local certificate occupies in the asymmetric security header, i.e. the
     * whole certificate chain when there is one.
     */
    public int getLocalCertificateChainSize() {
        return localCertificateChainSize;
    }

    public void setRemoteNonce(byte[] remoteNonce) {
        this.remoteNonce = remoteNonce;
    }

    @Override
    public byte[] getRemoteNonce() {
        return remoteNonce;
    }

    @Override
    public X509Certificate getRemoteCertificate() {
        return remoteCertificate;
    }

    @Override
    public SecurityPolicy getSecurityPolicy() {
        return securityPolicy;
    }

    @Override
    public MessageSecurity getMessageSecurity() {
        return messageSecurity;
    }

    public byte[] encryptPassword(byte[] encodeablePassword, SecurityPolicy policy) {
        return encryptionHandler.encryptPassword(encodeablePassword, policy);
    }

    public void setSecurityHeader(SecurityHeader securityHeader) {
        this.securityHeader.set(securityHeader);
    }

    public SignatureData createClientSignature() throws GeneralSecurityException {
        return encryptionHandler.createClientSignature();
    }

    public void setRemoteCertificate(X509Certificate certificate) {
        this.remoteCertificate = certificate;
    }

    public RequestHeader createRequestHeader(long requestTimeout) {
        return createRequestHeader(requestTimeout, tm.getRequestHandle());
    }

    protected RequestHeader createRequestHeader(long requestTimeout, int requestHandle) {
        return new RequestHeader(
            new NodeId(authenticationToken.get()),
            getCurrentDateTime(),
            (long) requestHandle,
            0L,
            NULL_STRING,
            requestTimeout,
            NULL_EXTENSION_OBJECT
        );
    }

    public RequestHeader createRequestHeader() {
        return createRequestHeader(configuration.getRequestTimeout());
    }

    public static long getCurrentDateTime() {
        return (System.currentTimeMillis() * 10000) + EPOCH_OFFSET;
    }

    public void setAuthenticationToken(NodeIdTypeDefinition authenticationToken) {
        this.authenticationToken.set(authenticationToken);
    }

    public int getSecurityChannelId() {
        return Long.valueOf(securityHeader.get().getSecureChannelId()).intValue();
    }

    public int getRequestId() {
        return tm.getRequestHandle();
    }

}
