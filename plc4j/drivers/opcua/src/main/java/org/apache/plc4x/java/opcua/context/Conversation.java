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
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.ConversationContext.SendRequestContext;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Conversation {
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


    private final Logger logger = LoggerFactory.getLogger(Conversation.class);
    private final AtomicReference<SecurityHeader> securityHeader = new AtomicReference<>(new SecurityHeader(1, 1));
    private final AtomicLong senderSequenceNumber = new AtomicLong(-1);

    private final AtomicReference<NodeIdTypeDefinition> authenticationToken = new AtomicReference<>(new NodeIdTwoByte((short) 0));

    private final ConversationContext<OpcuaAPU> context;
    private final SecureChannelTransactionManager tm;

    private final SecurityPolicy securityPolicy;
    private final MessageSecurity messageSecurity;
    private final EncryptionHandler encryptionHandler;
    private final OpcuaDriverContext driverContext;
    private final OpcuaConfiguration configuration;

    private OpcuaProtocolLimits limits;

    private X509Certificate localCertificate = null;
    private X509Certificate remoteCertificate = null;
    private byte[] remoteNonce;
    private byte[] localNonce;

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

    public Conversation(ConversationContext<OpcuaAPU> context, OpcuaDriverContext driverContext, OpcuaConfiguration configuration) {
        this.context = context;
        this.tm = new SecureChannelTransactionManager();
        this.driverContext = driverContext;
        this.configuration = configuration;

        this.securityPolicy = determineSecurityPolicy(configuration);
        CertificateKeyPair senderKeyPair = driverContext.getCertificateKeyPair();

        if (this.securityPolicy != SecurityPolicy.NONE) {
            //Sender Certificate gets populated during the 'discover' phase when encryption is enabled.
            this.messageSecurity = configuration.getMessageSecurity();
            this.remoteCertificate = configuration.getServerCertificate();
            this.encryptionHandler = new EncryptionHandler(this, senderKeyPair.getPrivateKey());
            this.localCertificate = senderKeyPair.getCertificate();
        } else {
            this.messageSecurity = MessageSecurity.NONE;
            this.encryptionHandler = new EncryptionHandler(this, null);
        }

        Limits encodingLimits = configuration.getEncodingLimits();
        limits = new OpcuaProtocolLimits(
            encodingLimits.getReceiveBufferSize(),
            encodingLimits.getSendBufferSize(),
            encodingLimits.getMaxMessageSize(),
            encodingLimits.getMaxChunkCount()
        );
    }

    public CompletableFuture<OpcuaAcknowledgeResponse> requestHello() {
        logger.debug("Sending hello message to {}", this.driverContext.getEndpoint());
        OpcuaHelloRequest request = new OpcuaHelloRequest(FINAL,
            OpcuaConstants.PROTOCOLVERSION,
            new OpcuaProtocolLimits(
                limits.getReceiveBufferSize(),
                limits.getSendBufferSize(),
                limits.getMaxMessageSize(),
                limits.getMaxChunkCount()
            ),
            new PascalString(driverContext.getEndpoint())
        );

        // open messages are guaranteed to fit into 8192 bytes limit
        //CompletableFuture<OpcuaAcknowledgeResponse> future = new CompletableFuture<>();

        CompletableFuture<OpcuaAcknowledgeResponse> future = new CompletableFuture<>();
        sendRequest(request, future, configuration.getNegotiationTimeout())
            .unwrap(OpcuaAPU::getMessage)
            .check(OpcuaAcknowledgeResponse.class::isInstance)
            .unwrap(OpcuaAcknowledgeResponse.class::cast)
            .handle(opcuaAcknowledgeResponse -> {
                OpcuaProtocolLimits limits = opcuaAcknowledgeResponse.getLimits();
                // merge encoding limits to match common minimum:
                // our receipt buffer should not exceed server send buffer size,
                // our send buffer size should not exceed server receive buffer size
                // chunks and message sizes should match too
                this.limits = new OpcuaProtocolLimits(
                    Math.min(this.limits.getReceiveBufferSize(), limits.getSendBufferSize()),
                    Math.min(this.limits.getSendBufferSize(), limits.getReceiveBufferSize()),
                    Math.min(this.limits.getMaxMessageSize(), limits.getMaxMessageSize()),
                    Math.min(this.limits.getMaxChunkCount(), limits.getMaxChunkCount())
                );
                future.complete(opcuaAcknowledgeResponse);
            });
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
            context.fireDisconnected();
        }).thenApply(r -> null);
    }

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
        for (int count = chunks.size(), index = 0; index < count; index++) {
            boolean last = index + 1 == count;
            if (last) {
                sendRequest(chunks.get(index), future, configuration.getNegotiationTimeout())
                    .unwrap(OpcuaAPU::getMessage)
                    .check(replyType::isInstance)
                    .unwrap(replyType::cast)
                    .unwrap(msg -> encryptionHandler.decodeMessage(msg))
                    .check(replyType::isInstance)
                    .unwrap(replyType::cast)
                    .check(reply -> requestId == sequenceHeaderExtractor.apply(reply).getRequestId())
                    .check(reply -> sequenceValidator.test(sequenceHeaderExtractor.apply(reply), future))
                    .check(msg -> accumulateChunkUntilFinal(chunkStorage, msg.getChunk(), chunkExtractor.apply(msg)))
                    .unwrap(msg -> mergeChunks(chunkStorage, msg, sequenceHeaderExtractor.apply(msg), chunkAssembler))
                    .handle(response -> {
                        future.complete(response);
                    });
            } else {
                context.sendToWire(new OpcuaAPU(chunks.get(index)));
            }
        }
        return future;
    }

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
            false,           //Namespace Uri Specified
            false,            //Server Index Specified
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
        for (int count = chunks.size(), index = 0; index < count; index++) {
            boolean last = index + 1 == count;
            if (last) {
                BiFunction<OpcuaMessageResponse, BinaryPayload, OpcuaMessageResponse> chunkAssembler = (src, chunkPayload) ->
                    new OpcuaMessageResponse(src.getChunk(), src.getSecurityHeader(), chunkPayload);

                sendRequest(chunks.get(index), future, configuration.getRequestTimeout())
                    .unwrap(OpcuaAPU::getMessage)
                    .check(OpcuaMessageResponse.class::isInstance)
                    .unwrap(OpcuaMessageResponse.class::cast)
                    .unwrap(msg -> encryptionHandler.decodeMessage(msg))
                    .check(OpcuaMessageResponse.class::isInstance)
                    .unwrap(OpcuaMessageResponse.class::cast)
                    .check(OpcuaMessageResponse.class::isInstance)
                    .unwrap(OpcuaMessageResponse.class::cast)
                    .check(msg -> msg.getMessage().getSequenceHeader().getRequestId() == requestId)
                    .check(reply -> sequenceValidator.test(reply.getMessage().getSequenceHeader(), future))
                    .check(msg -> accumulateChunkUntilFinal(chunkStorage, msg.getChunk(), msg.getMessage()))
                    .unwrap(msg -> mergeChunks(chunkStorage, msg, msg.getMessage().getSequenceHeader(), chunkAssembler))
                    .handle(response -> {
                        if (response.getChunk().equals(FINAL)) {
                            logger.debug("Received response made of {} bytes for message id: {}, channel id:{}, token:{}",
                                response.getLengthInBytes(), requestId, response.getSecurityHeader().getSecureChannelId(),
                                response.getSecurityHeader().getSecureTokenId()
                            );
                            securityHeader.set(response.getSecurityHeader());

                            Payload message = response.getMessage();
                            ExtensionObjectDefinition extensionObjectBody;
                            if (message instanceof ExtensiblePayload) {
                                extensionObjectBody = (((ExtensiblePayload) message).getPayload()).getBody();
                            } else {
                                try {
                                    BinaryPayload binary = (BinaryPayload) message;
                                    ReadBufferByteBased buffer = new ReadBufferByteBased(binary.getPayload(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                                    extensionObjectBody = ExtensionObject.staticParse(buffer, false).getBody();
                                } catch (ParseException e) {
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
                    });

            } else {
                context.sendToWire(new OpcuaAPU(chunks.get(index)));
            }
        }
        return future;
    }

    private SendRequestContext<OpcuaAPU> sendRequest(MessagePDU messagePDU, CompletableFuture<?> future, long timeout) {
        return context.sendRequest(new OpcuaAPU(messagePDU))
            .onError((req, err) -> future.completeExceptionally(err))
            .expectResponse(OpcuaAPU.class, Duration.ofMillis(timeout))
            .onTimeout((e) -> future.completeExceptionally(e));
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

        return FINAL.equals(chunkType);
    }

    public void setLocalNonce(byte[] localNonce) {
        this.localNonce = localNonce;
    }

    // generate nonce used for setting up signing/encryption keys
    byte[] createNonce() {
        return createNonce(securityPolicy.getNonceLength());
    }

    byte[] createNonce(int nonceLength) {
        return RandomUtils.nextBytes(nonceLength);
    }

    public boolean isSymmetricEncryptionEnabled() {
        return messageSecurity == MessageSecurity.SIGN_ENCRYPT;
    }

    public boolean isSymmetricSigningEnabled() {
        return (messageSecurity == MessageSecurity.SIGN_ENCRYPT || messageSecurity == MessageSecurity.SIGN);
    }

    static SecurityPolicy determineSecurityPolicy(OpcuaConfiguration configuration) {
        if (configuration.isDiscovery() && configuration.getServerCertificate() == null) {
            // discovery is enabled and sender certificate is not known yet
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

    public OpcuaProtocolLimits getLimits() {
        return limits;
    }

    public byte[] getLocalNonce() {
        return localNonce;
    }

    public X509Certificate getLocalCertificate() {
        return localCertificate;
    }

    public void setRemoteNonce(byte[] remoteNonce) {
        this.remoteNonce = remoteNonce;
    }

    public byte[] getRemoteNonce() {
        return remoteNonce;
    }

    public X509Certificate getRemoteCertificate() {
        return remoteCertificate;
    }

    public SecurityPolicy getSecurityPolicy() {
        return securityPolicy;
    }

    public MessageSecurity getMessageSecurity() {
        return messageSecurity;
    }

    public byte[] encryptPassword(byte[] encodeablePassword) {
        return encryptionHandler.encryptPassword(encodeablePassword);
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
            requestHandle,                                         //RequestHandle
            0L,
            SecureChannel.NULL_STRING,
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
