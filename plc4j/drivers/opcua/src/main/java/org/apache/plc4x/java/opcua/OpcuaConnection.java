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
package org.apache.plc4x.java.opcua;

import static org.apache.plc4x.java.opcua.context.SecureChannel.getX509Certificate;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.time.*;
import java.util.Map.Entry;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.TimeUnit;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcNullAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.context.Conversation;
import org.apache.plc4x.java.opcua.context.OpcuaDriverContext;
import org.apache.plc4x.java.opcua.context.OpcuaWire;
import org.apache.plc4x.java.opcua.context.SecureChannel;
import org.apache.plc4x.java.opcua.protocol.OpcuaSubscriptionHandle;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.opcua.tag.OpcuaPlcTagHandler;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.transport.tcp.TcpTransportInstance;

import java.io.ByteArrayInputStream;
import java.net.InetSocketAddress;
import javax.xml.parsers.DocumentBuilderFactory;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.NodeList;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Predicate;
import java.util.stream.Collectors;

public class OpcuaConnection extends ConnectionBase<OpcuaConfiguration> implements OpcuaWire {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaConnection.class);

    /** Publishing interval used when a subscription tag doesn't ask for one. */
    private static final Duration DEFAULT_CYCLE_TIME = Duration.ofMillis(1000);
    public static final PascalString NULL_STRING = new PascalString(null);
    private static final ExpandedNodeId NULL_EXPANDED_NODEID = new ExpandedNodeId(false,
        false,
        new NodeIdTwoByte((short) 0),
        null,
        null
    );

    public static final ExtensionObject NULL_EXTENSION_OBJECT = new NullExtensionObjectWithMask(
        NULL_EXPANDED_NODEID,
        new ExtensionObjectEncodingMask(false, false, false));

    private static final long EPOCH_OFFSET = 116444736000000000L;         //Offset between OPC UA epoch time and linux epoch time.
    // IEC 61131-3 date types use 1990-01-01 as epoch, PlcDATE etc. use Unix epoch (1970-01-01).
    private static final long IEC_DATE_EPOCH_OFFSET_DAYS = LocalDate.of(1990, 1, 1).toEpochDay();
    private final Map<Long, OpcuaSubscriptionHandle> subscriptions = new ConcurrentHashMap<>();

    // Listener registry backing the OpcuaWire implementation: every incoming
    // OpcuaAPU off the codec is dispatched to each registered listener whose
    // predicate matches. One-shot listeners (expect) self-remove on first match.
    private final List<WireListener> listeners = new CopyOnWriteArrayList<>();
    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    private final OpcuaDriverContext driverContext = new OpcuaDriverContext();
    private OpcuaMessageCodec messageCodec;
    private SecureChannel channel;
    private Conversation conversation;
    private volatile boolean connected = false;

    // Phase 3: per-session cache of server-resolved node attributes (data type, array shape,
    // access rights), keyed by canonical node address ("ns=<ns>;<idType>=<id>"). OPC UA node
    // types are static for the lifetime of a session, so once resolved they never need
    // re-reading — both browsing and (Phase 4) typed read/write resolve through this cache.
    private final Map<String, NodeAttributes> nodeTypeCache = new ConcurrentHashMap<>();

    // Phase 5: per-session cache of resolved StructureDefinitions (the field layout of custom struct
    // types), keyed by the DataType node address. Populated lazily on the first read of a value of
    // that type and reused for all subsequent reads/writes.
    private final Map<String, StructureDefinition> structureDefinitionCache = new ConcurrentHashMap<>();

    public OpcuaConnection(OpcuaConfiguration configuration,
                           TransportInstance<?> transportInstance,
                           AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new OpcuaPlcTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    public boolean isConnected() {
        return connected && messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        LOGGER.debug("Opcua Driver running in ACTIVE mode.");
        // The new SPI hands us a resolved transport instance — pull the host/port
        // back out so the driver context can build the OPC UA endpoint URL.
        if (!(transportInstance instanceof TcpTransportInstance tcp)) {
            throw new PlcConnectionException(
                "OPC UA driver requires a TCP transport instance, got " + transportInstance.getClass().getName());
        }
        InetSocketAddress remote = tcp.getRemoteAddress();
        driverContext.initialize(
            getTransportCode(),
            remote.getHostString(),
            String.valueOf(remote.getPort()),
            // The driver-config (e.g. the OPC UA "/milo" path) is the part of the connection
            // URL after host:port; strict servers reject a Hello whose endpoint URL omits it,
            // so it must be carried into the endpoint the driver advertises. Read from the
            // TransportInstance interface so this works for any transport (tcp, tls, ...).
            transportInstance.getDriverConfig(),
            configuration);

        messageCodec = new OpcuaMessageCodec(transportInstance, this::handleIncoming);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming OPC UA data", e);
            }
        });

        this.conversation = new Conversation(this, driverContext, configuration);
        this.channel = new SecureChannel(conversation, driverContext, configuration, resolveAuthentication());

        try {
            // Discovery only carries information needed for the encrypted modes
            // (it fetches the server certificate). For SecurityPolicy.NONE the
            // GetEndpoints call adds nothing, so we skip it.
            if (configuration.isDiscovery()
                && configuration.getSecurityPolicy() != null
                && configuration.getSecurityPolicy() != org.apache.plc4x.java.opcua.security.SecurityPolicy.NONE) {
                LOGGER.debug("Discovering endpoints before connecting");
                EndpointDescription endpoint = channel.onDiscover()
                    .get(configuration.getNegotiationTimeout(), TimeUnit.MILLISECONDS);
                configuration.setServerCertificate(
                    getX509Certificate(endpoint.getServerCertificate().getStringValue()));
                // onDiscover() already performed Hello + OpenSecureChannel on this TCP
                // connection, so the secure channel is open. Reuse it and only establish the
                // session — a second Hello on the same connection would stall (Hello is a
                // once-per-connection message).
                channel.onConnectSession().get(configuration.getNegotiationTimeout(), TimeUnit.MILLISECONDS);
            } else {
                channel.onConnect().get(configuration.getNegotiationTimeout(), TimeUnit.MILLISECONDS);
            }
            connected = true;
            LOGGER.info("Established connection to server");
        } catch (Exception e) {
            throw new PlcConnectionException("Failed to establish OPC UA connection", e);
        }
    }

    private PlcAuthentication resolveAuthentication() {
        // Authentication passed to getConnection(url, authentication) takes precedence over
        // credentials embedded in the connection string.
        PlcAuthentication passed = getAuthentication();
        // PlcNullAuthentication is the explicit "no credentials / anonymous" marker — treat it the
        // same as no authentication rather than an unsupported token type.
        if (passed != null && !(passed instanceof PlcNullAuthentication)) {
            return passed;
        }
        if (configuration.getUsername() != null && configuration.getPassword() != null) {
            return new PlcUsernamePasswordAuthentication(configuration.getUsername(), configuration.getPassword());
        }
        return null;
    }

    @Override
    public void close() throws Exception {
        connected = false;
        for (Entry<Long, OpcuaSubscriptionHandle> subscriber : subscriptions.entrySet()) {
            subscriber.getValue().stopSubscriber();
        }
        if (channel != null) {
            try {
                // Wait for CloseSession + CloseSecureChannel to actually reach the server
                // before we tear down the socket below; otherwise the server leaks the
                // session/channel and refuses new channels once its limit is hit.
                channel.onDisconnect().get(configuration.getRequestTimeout(), TimeUnit.MILLISECONDS);
            } catch (Exception e) {
                LOGGER.warn("Error during secure channel disconnect", e);
            }
        }
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        listeners.clear();
        consumers.clear();
        super.close();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // OpcuaWire implementation — the conversation layer drives the transport through this.
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncoming(OpcuaAPU apu) {
        for (WireListener listener : listeners) {
            if (listener.predicate.test(apu)) {
                if (listener.oneShot) {
                    listeners.remove(listener);
                }
                listener.handler.accept(apu);
                if (listener.oneShot) {
                    return;
                }
            }
        }
    }

    @Override
    public void sendToWire(OpcuaAPU apu) {
        try {
            messageCodec.send(apu);
        } catch (MessageCodecException e) {
            throw new PlcRuntimeException("Error sending OPC UA message", e);
        }
    }

    @Override
    public CompletableFuture<OpcuaAPU> expect(Predicate<OpcuaAPU> predicate, Duration timeout) {
        CompletableFuture<OpcuaAPU> future = new CompletableFuture<>();
        WireListener listener = new WireListener(predicate, future::complete, true);
        listeners.add(listener);
        future.orTimeout(timeout.toMillis(), TimeUnit.MILLISECONDS)
            .exceptionally(error -> {
                listeners.remove(listener);
                return null;
            });
        return future;
    }

    @Override
    public Subscription subscribe(Predicate<OpcuaAPU> predicate, Consumer<OpcuaAPU> handler) {
        WireListener listener = new WireListener(predicate, handler, false);
        listeners.add(listener);
        return () -> listeners.remove(listener);
    }

    @Override
    public void fireDisconnected() {
        connected = false;
    }

    private static final class WireListener {
        private final Predicate<OpcuaAPU> predicate;
        private final Consumer<OpcuaAPU> handler;
        private final boolean oneShot;

        WireListener(Predicate<OpcuaAPU> predicate, Consumer<OpcuaAPU> handler, boolean oneShot) {
            this.predicate = predicate;
            this.handler = handler;
            this.oneShot = oneShot;
        }
    }

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        LOGGER.trace("Reading Value");

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        RequestHeader requestHeader = conversation.createRequestHeader();

        // Tags the builder rejected (unparseable address) carry their code straight into the
        // response; they have no node id to ask the server for.
        Map<String, PlcResponseItem<PlcValue>> rejectedTags = rejectedReadTags(request);
        List<String> sendable = sendableTagNames(request);
        if (sendable.isEmpty()) {
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(request, rejectedTags));
        }

        List<ReadValueId> readValueArray = new ArrayList<>(sendable.size());
        Map<String, PlcTag> tagMap = new LinkedHashMap<>();
        for (String tagName : sendable) {
            OpcuaTag tag = (OpcuaTag) request.getTag(tagName);
            tagMap.put(tagName, tag);

            NodeId nodeId = generateNodeId(tag);

            readValueArray.add(new ReadValueId(nodeId,
                tag.getAttributeId().getValue(),
                indexRangeOf(tag),
                new QualifiedName(0, NULL_STRING)));
        }

        ReadRequest opcuaReadRequest = new ReadRequest(
            requestHeader,
            0.0d,
            TimestampsToReturn.timestampsToReturnBoth,
            readValueArray
        );

        return conversation.submit(opcuaReadRequest, ReadResponse.class).thenCompose(response -> {
            List<DataValue> results = response.getResults();
            // Any tag whose value came back as a custom struct (captured as raw bytes) needs its
            // StructureDefinition resolved from the server before the raw body can be decoded into a
            // PlcStruct. Resolve those up-front (cached per type), then map the whole response.
            Map<String, CompletableFuture<StructureDefinition>> structFutures = new LinkedHashMap<>();
            int index = 0;
            for (String tagName : tagMap.keySet()) {
                DataValue dataValue = results.get(index++);
                if (dataValue.getValueSpecified() && containsRawStruct(dataValue.getValue())) {
                    ExpandedNodeId encodingNodeId = rawStructEncodingId(dataValue.getValue());
                    structFutures.put(tagName,
                        resolveStructureDefinition((OpcuaTag) tagMap.get(tagName), encodingNodeId));
                }
            }
            if (structFutures.isEmpty()) {
                return CompletableFuture.completedFuture(new DefaultPlcReadResponse(request,
                    inRequestOrder(request, readResponse(tagMap, results, Collections.emptyMap()), rejectedTags)));
            }
            return CompletableFuture.allOf(structFutures.values().toArray(new CompletableFuture[0]))
                .thenApply(v -> {
                    Map<String, StructureDefinition> structDefs = new HashMap<>();
                    structFutures.forEach((name, future) -> structDefs.put(name, future.getNow(null)));
                    return new DefaultPlcReadResponse(request,
                        inRequestOrder(request, readResponse(tagMap, results, structDefs), rejectedTags));
                });
        });
    }

    /**
     * The tag names that can actually be put on the wire. A tag whose address the builder
     * couldn't parse stays in the request with an error code and a {@code null} tag; it must
     * neither be sent nor take up a slot in the response, which OPC UA maps back to tags by
     * position.
     */
    private static List<String> sendableTagNames(PlcTagRequest request) {
        List<String> names = new ArrayList<>(request.getNumberOfTags());
        for (String tagName : request.getTagNames()) {
            if (request.getTagResponseCode(tagName) == PlcResponseCode.OK) {
                names.add(tagName);
            }
        }
        return names;
    }

    /** Response items for the tags the builder rejected, keyed by name. */
    private static Map<String, PlcResponseItem<PlcValue>> rejectedReadTags(PlcTagRequest request) {
        Map<String, PlcResponseItem<PlcValue>> rejected = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcResponseCode code = request.getTagResponseCode(tagName);
            if (code != PlcResponseCode.OK) {
                rejected.put(tagName, new DefaultPlcResponseItem<>(code, null));
            }
        }
        return rejected;
    }

    /** Response codes for the tags the builder rejected, keyed by name. */
    private static Map<String, PlcResponseCode> rejectedWriteTags(PlcTagRequest request) {
        Map<String, PlcResponseCode> rejected = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcResponseCode code = request.getTagResponseCode(tagName);
            if (code != PlcResponseCode.OK) {
                rejected.put(tagName, code);
            }
        }
        return rejected;
    }

    /** The tag's OPC UA IndexRange as a PascalString, or the null string when the whole node is addressed. */
    private static PascalString indexRangeOf(OpcuaTag tag) {
        String indexRange = tag.getIndexRange();
        return indexRange != null ? new PascalString(indexRange) : NULL_STRING;
    }

    public static NodeId generateNodeId(OpcuaTag tag) {
        NodeId nodeId = null;
        if (tag.getIdentifierType() == OpcuaIdentifierType.BINARY_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdTwoByte(Short.parseShort(tag.getIdentifier())));
        } else if (tag.getIdentifierType() == OpcuaIdentifierType.NUMBER_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdNumeric(tag.getNamespace(), Long.parseLong(tag.getIdentifier())));
        } else if (tag.getIdentifierType() == OpcuaIdentifierType.GUID_IDENTIFIER) {
            UUID guid = UUID.fromString(tag.getIdentifier());
            ByteBuffer bb = ByteBuffer.allocate(16)
                    .order(ByteOrder.LITTLE_ENDIAN)
                    .putInt((int)(guid.getMostSignificantBits() >> (4*8)))
                    .putShort((short)(guid.getMostSignificantBits() >> (2*8)))
                    .putShort((short)guid.getMostSignificantBits())
                    .order(ByteOrder.BIG_ENDIAN)
                    .putLong(guid.getLeastSignificantBits());
            nodeId = new NodeId(new NodeIdGuid(tag.getNamespace(), bb.array()));
        } else if (tag.getIdentifierType() == OpcuaIdentifierType.STRING_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdString(tag.getNamespace(), new PascalString(tag.getIdentifier())));
        }
        return nodeId;
    }

    /** Merges decoded and rejected tags back into the order the caller asked for. */
    private static Map<String, PlcResponseItem<PlcValue>> inRequestOrder(PlcTagRequest request,
                                                                        Map<String, PlcResponseItem<PlcValue>> decoded,
                                                                        Map<String, PlcResponseItem<PlcValue>> rejected) {
        Map<String, PlcResponseItem<PlcValue>> ordered = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcResponseItem<PlcValue> item = decoded.get(tagName);
            ordered.put(tagName, item != null ? item : rejected.get(tagName));
        }
        return ordered;
    }

    public Map<String, PlcResponseItem<PlcValue>> readResponse(Map<String, PlcTag> tagMap, List<DataValue> results) {
        return readResponse(tagMap, results, Collections.emptyMap());
    }

    public Map<String, PlcResponseItem<PlcValue>> readResponse(Map<String, PlcTag> tagMap, List<DataValue> results,
                                                               Map<String, StructureDefinition> structDefs) {
        Map<String, PlcResponseItem<PlcValue>> response = new HashMap<>();
        int index = 0;
        for (String tagName : tagMap.keySet()) {
            PlcTag tag = tagMap.get(tagName);
            PlcValue value = null;
            DataValue dataValue = results.get(index++);
            PlcResponseCode responseCode = PlcResponseCode.OK;
            if (dataValue.getValueSpecified()) {
                StructureDefinition structDef = structDefs.get(tagName);
                if (structDef != null && dataValue.getValue() instanceof VariantExtensionObject) {
                    value = extensionObjectToPlcValue((VariantExtensionObject) dataValue.getValue(), structDef);
                } else {
                    value = variantToPlcValue(tag, dataValue.getValue());
                }
                if (value == null) {
                    LOGGER.error("Variant type {} is not supported.", dataValue.getValue().getClass());
                    responseCode = PlcResponseCode.UNSUPPORTED;
                }
            } else {
                StatusCode statusCode = dataValue.getStatusCode();
                responseCode = mapOpcStatusCode(statusCode.getStatusCode(), PlcResponseCode.UNSUPPORTED);
                LOGGER.error("Error while reading value from OPC UA server error code: {}", statusCode.toString());
            }

            response.put(tagName, new DefaultPlcResponseItem<>(responseCode, value));
        }
        return response;
    }

    // ======================================================================================
    // Struct (PlcStruct) support — decode custom structure values captured as raw ExtensionObject
    // bodies against the server-declared StructureDefinition.
    // ======================================================================================

    /** Whether a Variant carries at least one custom struct captured as raw bytes. */
    private static boolean containsRawStruct(Variant variant) {
        if (!(variant instanceof VariantExtensionObject)) {
            return false;
        }
        List<ExtensionObject> extensionObjects = ((VariantExtensionObject) variant).getValue();
        if (extensionObjects == null) {
            return false;
        }
        return extensionObjects.stream().anyMatch(eo -> eo instanceof RawBinaryExtensionObjectWithMask);
    }

    /**
     * Resolves (and caches, by encoding node) the {@link StructureDefinition} of a struct-typed
     * value. Tries the modern DataTypeDefinition attribute first; if the server doesn't expose it
     * (pre-1.04 servers), falls back to the legacy binary type dictionary reachable from the
     * encoding node. Returns null if neither yields a layout.
     */
    private CompletableFuture<StructureDefinition> resolveStructureDefinition(OpcuaTag tag, ExpandedNodeId encodingNodeId) {
        String key = encodingNodeId != null ? nodeIdKey(new NodeId(encodingNodeId.getNodeId())) : cacheKey(tag);
        StructureDefinition cached = structureDefinitionCache.get(key);
        if (cached != null) {
            return CompletableFuture.completedFuture(cached);
        }
        return resolveStructureDefinitionModern(tag).thenCompose(modern -> {
            if (modern != null) {
                structureDefinitionCache.put(key, modern);
                return CompletableFuture.completedFuture(modern);
            }
            if (encodingNodeId == null) {
                return CompletableFuture.completedFuture(null);
            }
            return resolveStructureDefinitionFromDictionary(encodingNodeId).thenApply(legacy -> {
                if (legacy != null) {
                    structureDefinitionCache.put(key, legacy);
                }
                return legacy;
            });
        }).exceptionally(error -> {
            LOGGER.warn("Failed to resolve StructureDefinition for tag '{}'", tag, error);
            return null;
        });
    }

    /** Modern path: variable DataType (attr 14) -&gt; data-type node -&gt; DataTypeDefinition (attr 23). */
    private CompletableFuture<StructureDefinition> resolveStructureDefinitionModern(OpcuaTag tag) {
        NodeId variableNode = generateNodeId(tag);
        ReadValueId dataTypeRead = new ReadValueId(variableNode, AttributeId.DataType.getValue(),
            NULL_STRING, new QualifiedName(0, NULL_STRING));
        ReadRequest request = new ReadRequest(conversation.createRequestHeader(), 0.0d,
            TimestampsToReturn.timestampsToReturnNeither, Collections.singletonList(dataTypeRead));
        return conversation.submit(request, ReadResponse.class).thenCompose(response -> {
            NodeId dataTypeNodeId = dataTypeNodeIdOf(response.getResults().getFirst());
            if (dataTypeNodeId == null) {
                return CompletableFuture.completedFuture(null);
            }
            return readStructureDefinition(dataTypeNodeId);
        });
    }

    /** Reads the DataTypeDefinition attribute (23) of a data-type node and extracts its StructureDefinition. */
    private CompletableFuture<StructureDefinition> readStructureDefinition(NodeId dataTypeNodeId) {
        ReadValueId definitionRead = new ReadValueId(dataTypeNodeId, AttributeId.DataTypeDefinition.getValue(),
            NULL_STRING, new QualifiedName(0, NULL_STRING));
        ReadRequest request = new ReadRequest(conversation.createRequestHeader(), 0.0d,
            TimestampsToReturn.timestampsToReturnNeither, Collections.singletonList(definitionRead));
        return conversation.submit(request, ReadResponse.class).thenApply(response -> {
            DataValue value = response.getResults().getFirst();
            if (!value.getValueSpecified() || !(value.getValue() instanceof VariantExtensionObject)) {
                // Server doesn't expose the DataTypeDefinition attribute (e.g. pre-1.04 servers
                // return BadAttributeIdInvalid); the layout must come from the legacy type
                // dictionary instead (handled elsewhere).
                return null;
            }
            List<ExtensionObject> extensionObjects = ((VariantExtensionObject) value.getValue()).getValue();
            if (extensionObjects == null || extensionObjects.isEmpty()) {
                return null;
            }
            // The DataTypeDefinition is a well-known standard type (ns=0;i=101), so it is parsed
            // into a typed BinaryExtensionObjectWithMask whose body is the StructureDefinition.
            ExtensionObjectDefinition body = extensionObjects.getFirst().getBody();
            return (body instanceof StructureDefinition) ? (StructureDefinition) body : null;
        });
    }

    // Standard reference-type NodeIds used to walk from a struct's binary-encoding node to its
    // entry in the legacy type dictionary: encoding --HasDescription--> description, and the
    // dictionary --HasComponent--> description (so description --HasComponent(inverse)--> dictionary).
    private static final long HAS_COMPONENT = 47L;
    private static final long HAS_DESCRIPTION = 39L;

    /** The encoding NodeId of the first custom struct (raw ExtensionObject) in a Variant, or null. */
    private static ExpandedNodeId rawStructEncodingId(Variant variant) {
        if (!(variant instanceof VariantExtensionObject)) {
            return null;
        }
        for (ExtensionObject extensionObject : ((VariantExtensionObject) variant).getValue()) {
            if (extensionObject instanceof RawBinaryExtensionObjectWithMask) {
                return extensionObject.getTypeId();
            }
        }
        return null;
    }

    /**
     * Legacy path: resolve a struct's field layout from the binary type dictionary. Walks the
     * encoding node to its DataTypeDescription (which names the type in the dictionary) and to the
     * DataTypeDictionary node (whose value is the OPC binary schema XML), then parses that schema.
     */
    private CompletableFuture<StructureDefinition> resolveStructureDefinitionFromDictionary(ExpandedNodeId encodingNodeId) {
        NodeId encodingNode = new NodeId(encodingNodeId.getNodeId());
        return browseFirstTarget(encodingNode, HAS_DESCRIPTION, BrowseDirection.browseDirectionForward)
            .thenCompose(descriptionNode -> {
                if (descriptionNode == null) {
                    return CompletableFuture.completedFuture(null);
                }
                CompletableFuture<DataValue> nameValue = readNodeValue(descriptionNode, AttributeId.Value.getValue());
                CompletableFuture<NodeId> dictionaryNode =
                    browseFirstTarget(descriptionNode, HAS_COMPONENT, BrowseDirection.browseDirectionInverse);
                return nameValue.thenCombine(dictionaryNode, (nv, dn) -> new Object[]{nv, dn}).thenCompose(pair -> {
                    String typeName = stringValueOf((DataValue) pair[0]);
                    NodeId dictNode = (NodeId) pair[1];
                    if (typeName == null || dictNode == null) {
                        return CompletableFuture.completedFuture(null);
                    }
                    return readNodeValue(dictNode, AttributeId.Value.getValue()).thenApply(dictValue -> {
                        byte[] schema = byteStringValueOf(dictValue);
                        return schema == null ? null : parseBsdStructure(schema, typeName, encodingNode);
                    });
                });
            });
    }

    /** Browses one reference of the given type/direction and returns the first target node, or null. */
    private CompletableFuture<NodeId> browseFirstTarget(NodeId source, long referenceTypeId, BrowseDirection direction) {
        BrowseDescription description = new BrowseDescription(source, direction,
            new NodeId(new NodeIdNumeric(0, referenceTypeId)), true, 0L, BROWSE_RESULT_MASK_ALL);
        BrowseRequest request = new BrowseRequest(conversation.createRequestHeader(),
            new ViewDescription(new NodeId(new NodeIdTwoByte((short) 0)), 0L, 0L), 0L,
            Collections.singletonList(description));
        return conversation.submit(request, BrowseResponse.class).thenApply(response -> {
            BrowseResult result = response.getResults().getFirst();
            if (result.getReferences() == null || result.getReferences().isEmpty()) {
                return null;
            }
            ExpandedNodeId target = result.getReferences().getFirst().getNodeId();
            return target == null ? null : new NodeId(target.getNodeId());
        });
    }

    /** Reads a single attribute of a node and returns the raw DataValue. */
    private CompletableFuture<DataValue> readNodeValue(NodeId node, long attributeId) {
        ReadValueId readValueId = new ReadValueId(node, attributeId, NULL_STRING, new QualifiedName(0, NULL_STRING));
        ReadRequest request = new ReadRequest(conversation.createRequestHeader(), 0.0d,
            TimestampsToReturn.timestampsToReturnNeither, Collections.singletonList(readValueId));
        return conversation.submit(request, ReadResponse.class).thenApply(response -> response.getResults().getFirst());
    }

    private static String stringValueOf(DataValue dataValue) {
        if (dataValue == null || !dataValue.getValueSpecified() || !(dataValue.getValue() instanceof VariantString)) {
            return null;
        }
        List<PascalString> values = ((VariantString) dataValue.getValue()).getValue();
        return (values == null || values.isEmpty()) ? null : values.getFirst().getStringValue();
    }

    private static byte[] byteStringValueOf(DataValue dataValue) {
        if (dataValue == null || !dataValue.getValueSpecified() || !(dataValue.getValue() instanceof VariantByteString)) {
            return null;
        }
        List<ByteStringArray> arrays = ((VariantByteString) dataValue.getValue()).getValue();
        if (arrays == null || arrays.isEmpty()) {
            return null;
        }
        List<Short> bytes = arrays.getFirst().getValue();
        byte[] result = new byte[bytes.size()];
        for (int i = 0; i < bytes.size(); i++) {
            result[i] = bytes.get(i).byteValue();
        }
        return result;
    }

    /**
     * Parses the OPC binary schema (bsd.xml) for the named structured type and turns its fields into
     * a {@link StructureDefinition} of built-in scalars/arrays. Nested/custom field types can't be
     * expressed this way yet, so such structures return null (decoded via a placeholder instead).
     */
    private static StructureDefinition parseBsdStructure(byte[] schema, String typeName, NodeId encodingNodeId) {
        try {
            DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
            factory.setNamespaceAware(true);
            factory.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
            Document document = factory.newDocumentBuilder().parse(new ByteArrayInputStream(schema));

            String localName = typeName.contains(":") ? typeName.substring(typeName.indexOf(':') + 1) : typeName;
            Element match = null;
            NodeList structuredTypes = document.getElementsByTagNameNS("*", "StructuredType");
            for (int i = 0; i < structuredTypes.getLength(); i++) {
                Element structuredType = (Element) structuredTypes.item(i);
                String name = structuredType.getAttribute("Name");
                if (name.equals(typeName) || name.equals(localName)) {
                    match = structuredType;
                    break;
                }
            }
            if (match == null) {
                return null;
            }

            NodeList fieldNodes = match.getElementsByTagNameNS("*", "Field");
            // Arrays are encoded as a preceding int32 length field plus the array field (which
            // references that length via LengthField); the length field itself is not a struct
            // member, so collect and skip those names.
            Set<String> lengthFields = new HashSet<>();
            for (int i = 0; i < fieldNodes.getLength(); i++) {
                String lengthField = ((Element) fieldNodes.item(i)).getAttribute("LengthField");
                if (!lengthField.isEmpty()) {
                    lengthFields.add(lengthField);
                }
            }
            List<StructureField> fields = new ArrayList<>();
            for (int i = 0; i < fieldNodes.getLength(); i++) {
                Element field = (Element) fieldNodes.item(i);
                String name = field.getAttribute("Name");
                if (lengthFields.contains(name)) {
                    continue;
                }
                Integer builtInId = bsdTypeNameToBuiltInId(field.getAttribute("TypeName"));
                if (builtInId == null) {
                    return null;
                }
                boolean isArray = !field.getAttribute("LengthField").isEmpty();
                fields.add(new StructureField(new PascalString(name), null,
                    new NodeId(new NodeIdNumeric(0, builtInId.longValue())), isArray ? 1 : -1, null, 0L, false));
            }
            return new StructureDefinition(encodingNodeId, new NodeId(new NodeIdNumeric(0, 22L)),
                StructureType.structureTypeStructure, fields);
        } catch (Exception e) {
            LOGGER.warn("Failed to parse type dictionary for '{}'", typeName, e);
            return null;
        }
    }

    /** Maps an OPC binary-schema TypeName (e.g. "opc:Int32") to its OPC UA built-in data-type id. */
    private static Integer bsdTypeNameToBuiltInId(String typeName) {
        if (typeName == null) {
            return null;
        }
        String local = typeName.contains(":") ? typeName.substring(typeName.indexOf(':') + 1) : typeName;
        return switch (local) {
            case "Boolean" -> 1;
            case "SByte" -> 2;
            case "Byte" -> 3;
            case "Int16" -> 4;
            case "UInt16" -> 5;
            case "Int32" -> 6;
            case "UInt32" -> 7;
            case "Int64" -> 8;
            case "UInt64" -> 9;
            case "Float" -> 10;
            case "Double" -> 11;
            case "String", "CharArray" -> 12;
            case "DateTime" -> 13;
            case "Guid" -> 14;
            case "ByteString" -> 15;
            default -> null;
        };
    }

    /** Turns a struct Variant (single or array of ExtensionObjects) into a PlcStruct / PlcList of PlcStruct. */
    private PlcValue extensionObjectToPlcValue(VariantExtensionObject variant, StructureDefinition definition) {
        List<ExtensionObject> extensionObjects = variant.getValue();
        List<PlcValue> values = new ArrayList<>(extensionObjects.size());
        for (ExtensionObject extensionObject : extensionObjects) {
            if (extensionObject instanceof RawBinaryExtensionObjectWithMask) {
                values.add(decodeStruct(((RawBinaryExtensionObjectWithMask) extensionObject).getRawBody(), definition));
            } else {
                // Standard (well-known) structure types are already parsed; surface them textually
                // for now (decoding those to PlcStruct is a later step).
                values.add(new PlcSTRING(extensionObject.toString()));
            }
        }
        return structurePlcValues(values, variant);
    }

    /** Decodes a struct's raw binary body into a PlcStruct using its StructureDefinition. */
    private PlcValue decodeStruct(byte[] rawBody, StructureDefinition definition) {
        try {
            ReadBufferByteBased buffer = new ReadBufferByteBased(rawBody, PayloadConverter.LITTLE_ENDIAN);
            Map<String, PlcValue> fields = new LinkedHashMap<>();
            boolean withOptionalFields =
                definition.getStructureType() == StructureType.structureTypeStructureWithOptionalFields;
            // Structures with optional fields are prefixed by a bit mask (one bit per optional field
            // in declaration order) telling which ones are actually present in the encoding.
            long optionalMask = withOptionalFields ? buffer.readUnsignedInt(32) : 0L;
            int optionalIndex = 0;
            for (StructureField field : definition.getFields()) {
                String name = field.getName() != null ? field.getName().getStringValue() : null;
                if (field.getIsOptional()) {
                    boolean present = (optionalMask & (1L << optionalIndex)) != 0;
                    optionalIndex++;
                    if (!present) {
                        continue;
                    }
                }
                fields.put(name, decodeField(buffer, field));
            }
            return new PlcStruct(fields);
        } catch (Exception e) {
            LOGGER.error("Failed to decode struct value", e);
            return null;
        }
    }

    /** Decodes a single struct field (scalar or 1-D array) from the buffer. */
    private static PlcValue decodeField(ReadBuffer buffer, StructureField field) throws Exception {
        Long dataTypeId = numericIdentifierOf(field.getDataType());
        if (field.getValueRank() < 0) {
            return decodeScalar(buffer, dataTypeId);
        }
        int length = buffer.readSignedInt(32);
        if (length < 0) {
            return new PlcList(Collections.emptyList());
        }
        List<PlcValue> elements = new ArrayList<>(length);
        for (int i = 0; i < length; i++) {
            elements.add(decodeScalar(buffer, dataTypeId));
        }
        return new PlcList(elements);
    }

    /** Decodes one scalar value of the given OPC UA built-in data type from the buffer. */
    private static PlcValue decodeScalar(ReadBuffer buffer, Long builtInDataTypeId) throws Exception {
        int id = builtInDataTypeId == null ? -1 : builtInDataTypeId.intValue();
        // Each OPC UA built-in is read with the smallest buffer reader that covers its bit width
        // (readUnsignedByte only handles 1-7 bits, so 8-bit and unsigned values step up a size).
        return switch (id) {
            case 1 -> new PlcBOOL(buffer.readUnsignedShort(8) != 0);
            case 2 -> new PlcSINT(buffer.readSignedByte(8));
            case 3 -> new PlcUSINT(buffer.readUnsignedShort(8));
            case 4 -> new PlcINT(buffer.readSignedShort(16));
            case 5 -> new PlcUINT(buffer.readUnsignedInt(16));
            case 6 -> new PlcDINT(buffer.readSignedInt(32));
            case 7 -> new PlcUDINT(buffer.readUnsignedLong(32));
            case 8 -> new PlcLINT(buffer.readSignedLong(64));
            case 9 -> new PlcULINT(buffer.readUnsignedBigInteger(64));
            case 10 -> new PlcREAL(buffer.readFloat(32));
            case 11 -> new PlcLREAL(buffer.readDouble(64));
            case 12 -> new PlcSTRING(PascalString.staticParse(buffer).getStringValue());
            default -> throw new PlcRuntimeException("Unsupported struct field data type id " + id
                + " (nested structs, enums and non-builtin types are not yet decoded)");
        };
    }

    /** Canonical key for a NodeId (namespace + identifier) used to cache resolved type layouts. */
    private static String nodeIdKey(NodeId nodeId) {
        return namespaceOf(nodeId) + ":" + nodeId.getNodeId().getIdentifier();
    }

    // ======================================================================================
    // Struct (PlcStruct) support — WRITE side: encode a PlcStruct back into a custom-struct
    // ExtensionObject (Phase 5d), mirroring the read decoding.
    // ======================================================================================

    private static final long HAS_ENCODING = 38L;
    private static final ExtensionObjectEncodingMask STRUCT_BINARY_ENCODING_MASK =
        new ExtensionObjectEncodingMask(false, false, true);
    // The field layout + binary-encoding node needed to write a value of a custom struct type.
    private final Map<String, StructWriteInfo> structWriteInfoCache = new ConcurrentHashMap<>();

    private static final class StructWriteInfo {
        private final StructureDefinition definition;
        private final ExpandedNodeId encodingNodeId;

        private StructWriteInfo(StructureDefinition definition, ExpandedNodeId encodingNodeId) {
            this.definition = definition;
            this.encodingNodeId = encodingNodeId;
        }
    }

    /** Whether a written value is a struct (a single PlcStruct — arrays of structs are a later step). */
    private static boolean isStructValue(PlcValue value) {
        return value != null && value.isStruct();
    }

    /** Resolves the struct layout + encoding node for every struct-valued write tag. */
    private CompletableFuture<Map<String, StructWriteInfo>> resolveStructWriteInfos(DefaultPlcWriteRequest request) {
        Map<String, StructWriteInfo> resolved = new ConcurrentHashMap<>();
        List<CompletableFuture<?>> futures = new ArrayList<>();
        for (String tagName : sendableTagNames(request)) {
            if (!isStructValue(request.getPlcValue(tagName))) {
                continue;
            }
            futures.add(resolveStructWriteInfo((OpcuaTag) request.getTag(tagName))
                .thenAccept(info -> {
                    if (info != null) {
                        resolved.put(tagName, info);
                    }
                }));
        }
        return CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).thenApply(v -> resolved);
    }

    /**
     * Resolves (and caches) the layout and binary-encoding node of a struct-typed node for writing:
     * reads the node's DataType, finds its "Default Binary" encoding node and its StructureDefinition
     * (modern DataTypeDefinition attribute, else the legacy type dictionary).
     */
    private CompletableFuture<StructWriteInfo> resolveStructWriteInfo(OpcuaTag tag) {
        String key = cacheKey(tag);
        StructWriteInfo cached = structWriteInfoCache.get(key);
        if (cached != null) {
            return CompletableFuture.completedFuture(cached);
        }
        NodeId variableNode = generateNodeId(tag);
        return readNodeValue(variableNode, AttributeId.DataType.getValue()).thenCompose(dataTypeValue -> {
            NodeId dataTypeNode = dataTypeNodeIdOf(dataTypeValue);
            if (dataTypeNode == null) {
                return CompletableFuture.completedFuture(null);
            }
            return browseBinaryEncodingNode(dataTypeNode).thenCompose(encodingNodeId -> {
                if (encodingNodeId == null) {
                    return CompletableFuture.completedFuture(null);
                }
                return readStructureDefinition(dataTypeNode).thenCompose(modern -> {
                    if (modern != null) {
                        return CompletableFuture.completedFuture(new StructWriteInfo(modern, encodingNodeId));
                    }
                    return resolveStructureDefinitionFromDictionary(encodingNodeId)
                        .thenApply(legacy -> legacy == null ? null : new StructWriteInfo(legacy, encodingNodeId));
                });
            });
        }).thenApply(info -> {
            if (info != null) {
                structWriteInfoCache.put(key, info);
            }
            return info;
        }).exceptionally(error -> {
            LOGGER.warn("Failed to resolve struct write info for tag '{}'", tag, error);
            return null;
        });
    }

    /** Browses a data-type node's encodings and returns its "Default Binary" encoding node. */
    private CompletableFuture<ExpandedNodeId> browseBinaryEncodingNode(NodeId dataTypeNode) {
        BrowseDescription description = new BrowseDescription(dataTypeNode, BrowseDirection.browseDirectionForward,
            new NodeId(new NodeIdNumeric(0, HAS_ENCODING)), true, 0L, BROWSE_RESULT_MASK_ALL);
        BrowseRequest request = new BrowseRequest(conversation.createRequestHeader(),
            new ViewDescription(new NodeId(new NodeIdTwoByte((short) 0)), 0L, 0L), 0L,
            Collections.singletonList(description));
        return conversation.submit(request, BrowseResponse.class).thenApply(response -> {
            List<ReferenceDescription> references = response.getResults().getFirst().getReferences();
            if (references == null || references.isEmpty()) {
                return null;
            }
            for (ReferenceDescription reference : references) {
                if ("Default Binary".equals(qualifiedNameValue(reference.getBrowseName()))) {
                    return reference.getNodeId();
                }
            }
            return references.getFirst().getNodeId();
        });
    }

    /** Encodes a PlcStruct as a single-element ExtensionObject Variant of the given custom type. */
    private Variant encodeStructVariant(PlcValue value, StructWriteInfo info) {
        byte[] body = encodeStructBody(value, info.definition);
        ExtensionObject extensionObject =
            new RawBinaryExtensionObjectWithMask(info.encodingNodeId, STRUCT_BINARY_ENCODING_MASK, body);
        return new VariantExtensionObject(false, false, null, Collections.emptyList(), null,
            Collections.singletonList(extensionObject));
    }

    /** Serializes a PlcStruct into the OPC UA binary body described by its StructureDefinition. */
    private static byte[] encodeStructBody(PlcValue value, StructureDefinition definition) {
        try {
            byte[] body = new byte[structBodySize(value, definition)];
            WriteBufferByteBased buffer = new WriteBufferByteBased(body, PayloadConverter.LITTLE_ENDIAN);
            for (StructureField field : definition.getFields()) {
                encodeField(buffer, field, structField(value, field));
            }
            return body;
        } catch (Exception e) {
            throw new PlcRuntimeException("Failed to encode struct value", e);
        }
    }

    private static PlcValue structField(PlcValue struct, StructureField field) {
        String name = field.getName() != null ? field.getName().getStringValue() : null;
        PlcValue value = struct.getValue(name);
        if (value == null) {
            throw new PlcRuntimeException("Missing struct field '" + name + "' in value to write");
        }
        return value;
    }

    private static void encodeField(WriteBuffer buffer, StructureField field, PlcValue value) throws Exception {
        Long dataTypeId = numericIdentifierOf(field.getDataType());
        if (field.getValueRank() < 0) {
            encodeScalar(buffer, dataTypeId, value);
            return;
        }
        List<? extends PlcValue> elements = value.getList();
        buffer.writeSignedInt(32, elements.size());
        for (PlcValue element : elements) {
            encodeScalar(buffer, dataTypeId, element);
        }
    }

    private static void encodeScalar(WriteBuffer buffer, Long builtInDataTypeId, PlcValue value) throws Exception {
        int id = builtInDataTypeId == null ? -1 : builtInDataTypeId.intValue();
        switch (id) {
            case 1:  buffer.writeUnsignedShort(8, (short) (value.getBoolean() ? 1 : 0)); break;
            case 2:  buffer.writeSignedByte(8, value.getByte()); break;
            case 3:  buffer.writeUnsignedShort(8, value.getShort()); break;
            case 4:  buffer.writeSignedShort(16, value.getShort()); break;
            case 5:  buffer.writeUnsignedInt(16, value.getInt()); break;
            case 6:  buffer.writeSignedInt(32, value.getInt()); break;
            case 7:  buffer.writeUnsignedLong(32, value.getLong()); break;
            case 8:  buffer.writeSignedLong(64, value.getLong()); break;
            case 9:  buffer.writeUnsignedBigInteger(64, value.getBigInteger()); break;
            case 10: buffer.writeFloat(32, value.getFloat()); break;
            case 11: buffer.writeDouble(64, value.getDouble()); break;
            case 12: new PascalString(value.getString()).serialize(buffer); break;
            default:
                throw new PlcRuntimeException("Unsupported struct field data type id " + id + " for writing");
        }
    }

    private static int structBodySize(PlcValue value, StructureDefinition definition) {
        int size = 0;
        for (StructureField field : definition.getFields()) {
            PlcValue fieldValue = structField(value, field);
            Long dataTypeId = numericIdentifierOf(field.getDataType());
            if (field.getValueRank() < 0) {
                size += scalarSize(dataTypeId, fieldValue);
            } else {
                size += 4; // int32 array length
                for (PlcValue element : fieldValue.getList()) {
                    size += scalarSize(dataTypeId, element);
                }
            }
        }
        return size;
    }

    private static int scalarSize(Long builtInDataTypeId, PlcValue value) {
        int id = builtInDataTypeId == null ? -1 : builtInDataTypeId.intValue();
        return switch (id) {
            case 1, 2, 3 -> 1;
            case 4, 5 -> 2;
            case 6, 7, 10 -> 4;
            case 8, 9, 11 -> 8;
            case 12 -> new PascalString(value.getString()).getLengthInBytes();
            default -> throw new PlcRuntimeException("Unsupported struct field data type id " + id + " for writing");
        };
    }

    // The standard "Objects" folder — the usual entry point into a server's address space.
    private static final String OBJECTS_FOLDER_ADDRESS = "ns=0;i=85";
    // HierarchicalReferences (i=33): browsing only these (plus subtypes) yields the clean
    // containment tree (Organizes / HasComponent / HasProperty, ...) instead of every
    // reference type (which would drag in type definitions and other non-containment noise).
    private static final NodeId HIERARCHICAL_REFERENCES = new NodeId(new NodeIdNumeric(0, 33L));
    // resultMask 0x3F -> return all reference fields (referenceType, isForward, nodeClass,
    // browseName, displayName, typeDefinition).
    private static final long BROWSE_RESULT_MASK_ALL = 0x3FL;

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return onBrowseWithInterceptor(browseRequest, (queryName, query, item) -> true);
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest,
                                                                           PlcBrowseRequestInterceptor interceptor) {
        Map<String, PlcResponseCode> responseCodes = new ConcurrentHashMap<>();
        Map<String, List<PlcBrowseItem>> values = new ConcurrentHashMap<>();
        List<CompletableFuture<?>> queryFutures = new ArrayList<>();

        for (String queryName : browseRequest.getQueryNames()) {
            PlcQuery query = browseRequest.getQuery(queryName);
            String startAddress = query.getQueryString();
            // An empty query or the "**" wildcard means "browse everything": start at the
            // standard Objects folder and recurse the whole sub-tree. Otherwise the query is
            // the address of the node to start browsing from.
            if (startAddress == null || startAddress.isBlank() || "**".equals(startAddress.trim())) {
                startAddress = OBJECTS_FOLDER_ADDRESS;
            }
            NodeId startNodeId;
            try {
                startNodeId = generateNodeId(OpcuaTag.of(startAddress));
            } catch (Exception e) {
                LOGGER.warn("Invalid browse start node '{}' for query '{}'", startAddress, queryName, e);
                responseCodes.put(queryName, PlcResponseCode.INVALID_ADDRESS);
                values.put(queryName, Collections.emptyList());
                continue;
            }
            // A shared, global set of already-visited nodes provides cycle detection (a node
            // is never expanded twice) and keeps a well-formed tree from being duplicated.
            Set<String> visited = ConcurrentHashMap.newKeySet();
            queryFutures.add(
                browseNode(startNodeId, visited, queryName, query, interceptor)
                    .handle((items, error) -> {
                        if (error != null) {
                            LOGGER.warn("Browsing failed for query '{}'", queryName, error);
                            responseCodes.put(queryName, PlcResponseCode.INTERNAL_ERROR);
                            values.put(queryName, Collections.emptyList());
                        } else {
                            responseCodes.put(queryName, PlcResponseCode.OK);
                            values.put(queryName, items);
                        }
                        return null;
                    }));
        }

        return CompletableFuture.allOf(queryFutures.toArray(new CompletableFuture[0]))
            .thenApply(v -> new DefaultPlcBrowseResponse(browseRequest, responseCodes, values));
    }

    /**
     * Browses the forward hierarchical references of {@code nodeId} and returns each referenced
     * node as a {@link PlcBrowseItem}, recursing into every not-yet-visited child so the full
     * sub-tree is returned. {@code visited} guards against reference cycles.
     */
    private CompletableFuture<List<PlcBrowseItem>> browseNode(NodeId nodeId, Set<String> visited, String queryName,
                                                              PlcQuery query, PlcBrowseRequestInterceptor interceptor) {
        return browseReferences(nodeId, null, new ArrayList<>()).thenCompose(references ->
            // Resolve the server-side type/access attributes of all variable children in one
            // batched Read before turning the references into browse items.
            resolveVariableAttributes(references).thenCompose(attributes -> {
                List<CompletableFuture<PlcBrowseItem>> itemFutures = new ArrayList<>();
                for (ReferenceDescription reference : references) {
                    String childAddress = addressOf(reference.getNodeId());
                    if (childAddress == null) {
                        // Node identifier form we can't turn back into an address (e.g. opaque) — skip.
                        continue;
                    }
                    CompletableFuture<Map<String, PlcBrowseItem>> childrenFuture;
                    if (visited.add(childAddress)) {
                        childrenFuture = browseNode(generateNodeId(OpcuaTag.of(childAddress)),
                            visited, queryName, query, interceptor)
                            .thenApply(childItems -> childItems.stream()
                                .collect(Collectors.toMap(PlcBrowseItem::getName, item -> item,
                                    (a, b) -> a, LinkedHashMap::new)));
                    } else {
                        // Already visited (cycle or shared node) — list it but don't expand again.
                        childrenFuture = CompletableFuture.completedFuture(Collections.emptyMap());
                    }
                    NodeAttributes childAttributes = attributes.get(childAddress);
                    itemFutures.add(childrenFuture.thenApply(children ->
                        buildBrowseItem(reference, childAddress, childAttributes, children)));
                }
                return CompletableFuture.allOf(itemFutures.toArray(new CompletableFuture[0]))
                    .thenApply(v -> itemFutures.stream()
                        .map(CompletableFuture::join)
                        .filter(item -> interceptor.intercept(queryName, query, item))
                        .collect(Collectors.toList()));
            }));
    }

    /**
     * Issues a single Browse (or BrowseNext when {@code continuationPoint} is set) and follows
     * continuation points until the server has returned all references of the node.
     */
    private CompletableFuture<List<ReferenceDescription>> browseReferences(NodeId nodeId, PascalByteString continuationPoint,
                                                                           List<ReferenceDescription> accumulator) {
        CompletableFuture<BrowseResult> resultFuture;
        if (continuationPoint == null) {
            BrowseDescription description = new BrowseDescription(
                nodeId,
                BrowseDirection.browseDirectionForward,
                HIERARCHICAL_REFERENCES,
                true,               // include reference subtypes
                0L,                 // nodeClassMask 0 -> all node classes
                BROWSE_RESULT_MASK_ALL);
            BrowseRequest request = new BrowseRequest(
                conversation.createRequestHeader(),
                new ViewDescription(new NodeId(new NodeIdTwoByte((short) 0)), 0L, 0L),
                0L,                 // requestedMaxReferencesPerNode 0 -> server decides
                Collections.singletonList(description));
            resultFuture = conversation.submit(request, BrowseResponse.class)
                .thenApply(response -> response.getResults().getFirst());
        } else {
            BrowseNextRequest request = new BrowseNextRequest(
                conversation.createRequestHeader(),
                false,              // don't release the continuation point, we want the next batch
                Collections.singletonList(continuationPoint));
            resultFuture = conversation.submit(request, BrowseNextResponse.class)
                .thenApply(response -> response.getResults().getFirst());
        }
        return resultFuture.thenCompose(result -> {
            if (result.getReferences() != null) {
                accumulator.addAll(result.getReferences());
            }
            PascalByteString nextContinuationPoint = result.getContinuationPoint();
            if (nextContinuationPoint != null && nextContinuationPoint.getStringLength() > 0) {
                return browseReferences(nodeId, nextContinuationPoint, accumulator);
            }
            return CompletableFuture.completedFuture(accumulator);
        });
    }

    private PlcBrowseItem buildBrowseItem(ReferenceDescription reference, String address,
                                          NodeAttributes attributes, Map<String, PlcBrowseItem> children) {
        NodeClass nodeClass = reference.getNodeClass();
        boolean isVariable = nodeClass == NodeClass.nodeClassVariable;
        String name = localizedTextValue(reference.getDisplayName());
        if (name == null || name.isEmpty()) {
            name = qualifiedNameValue(reference.getBrowseName());
        }

        // Fold the server-resolved data type into the tag address (e.g. "ns=2;i=3;DINT") so the
        // resulting tag is correctly typed for reads/writes without a manual type suffix.
        OpcuaDataType dataType = attributes != null ? attributes.dataType : null;
        boolean hasDataType = dataType != null && dataType != OpcuaDataType.NULL;
        OpcuaTag tag = hasDataType ? OpcuaTag.of(address + ";" + dataType.name()) : OpcuaTag.of(address);

        boolean readable;
        boolean writable;
        List<ArrayInfo> arrayInfo;
        if (attributes != null) {
            // Real access rights and array shape read from the server.
            readable = attributes.readable;
            writable = attributes.writable;
            arrayInfo = attributes.arrayInfo;
        } else {
            // Non-variable nodes (or variables whose attributes couldn't be read): fall back to the
            // NodeClass heuristic — only variables are considered readable/writable.
            readable = isVariable;
            writable = isVariable;
            arrayInfo = Collections.emptyList();
        }

        Map<String, PlcValue> options = new HashMap<>();
        if (nodeClass != null) {
            options.put("node-class", new PlcSTRING(nodeClass.name()));
        }
        String browseName = qualifiedNameValue(reference.getBrowseName());
        if (browseName != null) {
            options.put("browse-name", new PlcSTRING(browseName));
        }
        if (hasDataType) {
            options.put("data-type", new PlcSTRING(dataType.name()));
        }

        return new DefaultPlcBrowseItem(tag, name, readable, writable,
            Collections.emptySet(), false, arrayInfo, children, options);
    }

    // The node attributes read for every variable child during a browse so items carry a real
    // data type, array shape and access rights rather than the NodeClass-based guess of Phase 1.
    private static final AttributeId[] BROWSE_ATTRIBUTES = {
        AttributeId.DataType, AttributeId.ValueRank, AttributeId.ArrayDimensions, AttributeId.AccessLevel};

    // Abstract numeric supertype DataType NodeIds (namespace 0). A variable typed with one of these
    // accepts any concrete subtype, so there is no single built-in type to resolve to — the concrete
    // type is chosen from the written value instead. Only UInteger needs special handling: plain
    // (signed) value inference already yields a valid subtype for Integer/Number.
    private static final long UINTEGER_DATATYPE_ID = 28L;

    /** The server-resolved metadata of a single variable node. */
    static final class NodeAttributes {
        private final OpcuaDataType dataType;      // resolved built-in type, or null if unresolved
        private final boolean unsigned;            // node is the abstract UInteger supertype
        private final List<ArrayInfo> arrayInfo;   // empty for scalars
        private final boolean readable;
        private final boolean writable;

        private NodeAttributes(OpcuaDataType dataType, boolean unsigned, List<ArrayInfo> arrayInfo,
                               boolean readable, boolean writable) {
            this.dataType = dataType;
            this.unsigned = unsigned;
            this.arrayInfo = arrayInfo;
            this.readable = readable;
            this.writable = writable;
        }

        OpcuaDataType dataType() {
            return dataType;
        }

        boolean unsigned() {
            return unsigned;
        }

        List<ArrayInfo> arrayInfo() {
            return arrayInfo;
        }

        boolean readable() {
            return readable;
        }

        boolean writable() {
            return writable;
        }
    }

    /**
     * Resolves DataType/ValueRank/ArrayDimensions/AccessLevel for every variable node among the
     * given references and returns a map keyed by node address. Nodes already resolved this session
     * are served from {@link #nodeTypeCache}; only the rest incur a (batched) server read. Object/
     * method/other node classes are skipped (they have no such attributes). If the read fails the
     * browse still succeeds — items just fall back to the NodeClass heuristic (and any cached
     * entries are still applied).
     */
    private CompletableFuture<Map<String, NodeAttributes>> resolveVariableAttributes(List<ReferenceDescription> references) {
        List<String> variableAddresses = new ArrayList<>();
        Set<String> seen = new HashSet<>();
        for (ReferenceDescription reference : references) {
            if (reference.getNodeClass() != NodeClass.nodeClassVariable) {
                continue;
            }
            String address = addressOf(reference.getNodeId());
            if (address != null && seen.add(address)) {
                variableAddresses.add(address);
            }
        }
        if (variableAddresses.isEmpty()) {
            return CompletableFuture.completedFuture(Collections.emptyMap());
        }

        // Serve already-resolved nodes straight from the session cache; only read the remainder.
        Map<String, NodeAttributes> resolved = new HashMap<>();
        List<String> toRead = new ArrayList<>();
        for (String address : variableAddresses) {
            NodeAttributes cached = nodeTypeCache.get(address);
            if (cached != null) {
                resolved.put(address, cached);
            } else {
                toRead.add(address);
            }
        }
        if (toRead.isEmpty()) {
            return CompletableFuture.completedFuture(resolved);
        }

        return readNodeAttributes(toRead).thenApply(fresh -> {
            resolved.putAll(fresh);
            return resolved;
        }).exceptionally(error -> {
            LOGGER.warn("Failed to resolve variable node attributes during browse; "
                + "browse items fall back to the NodeClass heuristic", error);
            return resolved;
        });
    }

    /**
     * Lazily resolves the server-side attributes of a single node, caching the result for the
     * session. This is the entry point for typed reads/writes: on a cache hit no server round-trip
     * happens; on a miss the four attributes are read and cached. Returns {@code null} attributes
     * only if the read yields nothing usable.
     */
    CompletableFuture<NodeAttributes> resolveNodeAttributes(OpcuaTag tag) {
        String key = cacheKey(tag);
        NodeAttributes cached = nodeTypeCache.get(key);
        if (cached != null) {
            return CompletableFuture.completedFuture(cached);
        }
        return readNodeAttributes(Collections.singletonList(key)).thenApply(fresh -> fresh.get(key));
    }

    /**
     * Issues one batched {@link ReadRequest} for the {@link #BROWSE_ATTRIBUTES} of every given node
     * address, stores each resolved {@link NodeAttributes} in {@link #nodeTypeCache} and returns
     * them keyed by address.
     */
    private CompletableFuture<Map<String, NodeAttributes>> readNodeAttributes(List<String> addresses) {
        List<ReadValueId> readValueIds = new ArrayList<>(addresses.size() * BROWSE_ATTRIBUTES.length);
        for (String address : addresses) {
            NodeId nodeId = generateNodeId(OpcuaTag.of(address));
            for (AttributeId attributeId : BROWSE_ATTRIBUTES) {
                readValueIds.add(new ReadValueId(nodeId, attributeId.getValue(),
                    NULL_STRING, new QualifiedName(0, NULL_STRING)));
            }
        }
        ReadRequest readRequest = new ReadRequest(conversation.createRequestHeader(), 0.0d,
            TimestampsToReturn.timestampsToReturnNeither, readValueIds);

        return conversation.submit(readRequest, ReadResponse.class).thenApply(response -> {
            List<DataValue> results = response.getResults();
            Map<String, NodeAttributes> fresh = new LinkedHashMap<>();
            for (int i = 0; i < addresses.size(); i++) {
                String address = addresses.get(i);
                NodeAttributes attributes = attributesFrom(results, i * BROWSE_ATTRIBUTES.length);
                nodeTypeCache.put(address, attributes);
                fresh.put(address, attributes);
            }
            return fresh;
        });
    }

    /**
     * Builds a {@link NodeAttributes} from the four consecutive {@link DataValue}s (DataType,
     * ValueRank, ArrayDimensions, AccessLevel — in {@link #BROWSE_ATTRIBUTES} order) starting at
     * {@code offset} in a Read response.
     */
    private static NodeAttributes attributesFrom(List<DataValue> results, int offset) {
        DataValue dataTypeValue = results.get(offset);
        DataValue valueRankValue = results.get(offset + 1);
        DataValue arrayDimensionsValue = results.get(offset + 2);
        DataValue accessLevelValue = results.get(offset + 3);

        Short accessLevel = byteValue(accessLevelValue);
        // AccessLevel bit0 = CurrentRead, bit1 = CurrentWrite. If the attribute couldn't be read,
        // assume read/write (matching the earlier NodeClass-based default).
        boolean readable = accessLevel == null || (accessLevel & 0x01) != 0;
        boolean writable = accessLevel == null || (accessLevel & 0x02) != 0;

        NodeId dataTypeNodeId = dataTypeNodeIdOf(dataTypeValue);
        return new NodeAttributes(
            concreteDataType(dataTypeNodeId),
            isUnsignedAbstract(dataTypeNodeId),
            arrayInfoOf(valueRankValue, arrayDimensionsValue),
            readable, writable);
    }

    /** Canonical cache key for a tag: its node address without the attribute or data-type suffix. */
    private static String cacheKey(OpcuaTag tag) {
        return String.format("ns=%d;%s=%s", tag.getNamespace(),
            tag.getIdentifierType().getValue(), tag.getIdentifier());
    }

    /** Number of nodes whose attributes are currently cached for this session (test/diagnostic hook). */
    int nodeTypeCacheSize() {
        return nodeTypeCache.size();
    }

    /** Extracts the DataType attribute's NodeId from its {@link DataValue}, or null if not present. */
    private static NodeId dataTypeNodeIdOf(DataValue dataValue) {
        if (dataValue == null || !dataValue.getValueSpecified() || !(dataValue.getValue() instanceof VariantNodeId)) {
            return null;
        }
        List<NodeId> nodeIds = ((VariantNodeId) dataValue.getValue()).getValue();
        return (nodeIds == null || nodeIds.isEmpty()) ? null : nodeIds.getFirst();
    }

    /**
     * Resolves a DataType NodeId to a concrete {@link OpcuaDataType}. Built-in DataType NodeIds live
     * in namespace 0 and share their numeric identifier with the OPC UA Variant built-in type
     * (Boolean=1, Int32=6, ...), so we can map straight through. Abstract supertypes and
     * custom/enum/struct data types (non-namespace-0 or non-numeric) have no concrete built-in type
     * and return null.
     */
    private static OpcuaDataType concreteDataType(NodeId dataTypeNodeId) {
        if (dataTypeNodeId == null || namespaceOf(dataTypeNodeId) != 0) {
            return null;
        }
        Long numericId = numericIdentifierOf(dataTypeNodeId);
        if (numericId == null || numericId <= 0 || numericId > Short.MAX_VALUE) {
            return null;
        }
        return OpcuaDataType.firstEnumForFieldVariantType(numericId.shortValue());
    }

    /** Whether a DataType NodeId is the abstract UInteger supertype (so writes must pick an unsigned type). */
    private static boolean isUnsignedAbstract(NodeId dataTypeNodeId) {
        if (dataTypeNodeId == null || namespaceOf(dataTypeNodeId) != 0) {
            return false;
        }
        Long numericId = numericIdentifierOf(dataTypeNodeId);
        return numericId != null && numericId == UINTEGER_DATATYPE_ID;
    }

    /**
     * Builds the array shape from the ValueRank and ArrayDimensions attributes. ValueRank &lt; 0
     * (-1 scalar, -2 any, -3 scalar-or-one-dim) yields no array info; ValueRank &gt;= 0 marks an
     * array, its per-dimension sizes taken from ArrayDimensions when the server supplies them.
     */
    private static List<ArrayInfo> arrayInfoOf(DataValue valueRankValue, DataValue arrayDimensionsValue) {
        Integer valueRank = int32Value(valueRankValue);
        if (valueRank == null || valueRank < 0) {
            return Collections.emptyList();
        }
        List<Long> dimensions = uint32Values(arrayDimensionsValue);
        List<ArrayInfo> arrayInfo = new ArrayList<>();
        if (dimensions != null && !dimensions.isEmpty()) {
            for (Long dimension : dimensions) {
                // A dimension of 0 means "unknown size" per the spec — represent it as an empty range.
                int upper = (dimension == null || dimension <= 0) ? -1 : (int) (dimension - 1);
                arrayInfo.add(new DefaultArrayInfo(0, upper));
            }
        } else {
            // We know the rank but not the sizes (ValueRank 0 => one-or-more dimensions).
            int rank = Math.max(valueRank, 1);
            for (int i = 0; i < rank; i++) {
                arrayInfo.add(new DefaultArrayInfo(0, -1));
            }
        }
        return arrayInfo;
    }

    /** Extracts the first Byte (unsigned, as Short) from a Variant, or null if not a byte scalar. */
    private static Short byteValue(DataValue dataValue) {
        if (dataValue == null || !dataValue.getValueSpecified() || !(dataValue.getValue() instanceof VariantByte)) {
            return null;
        }
        List<Short> values = ((VariantByte) dataValue.getValue()).getValue();
        return (values == null || values.isEmpty()) ? null : values.getFirst();
    }

    /** Extracts the first Int32 from a Variant, or null if not an int32 scalar. */
    private static Integer int32Value(DataValue dataValue) {
        if (dataValue == null || !dataValue.getValueSpecified() || !(dataValue.getValue() instanceof VariantInt32)) {
            return null;
        }
        List<Integer> values = ((VariantInt32) dataValue.getValue()).getValue();
        return (values == null || values.isEmpty()) ? null : values.getFirst();
    }

    /** Extracts the UInt32 list from a Variant, or null if not a uint32 value. */
    private static List<Long> uint32Values(DataValue dataValue) {
        if (dataValue == null || !dataValue.getValueSpecified() || !(dataValue.getValue() instanceof VariantUInt32)) {
            return null;
        }
        return ((VariantUInt32) dataValue.getValue()).getValue();
    }

    /** Namespace index of a NodeId (0 for the implicit-namespace two-byte form; -1 if unknown). */
    private static int namespaceOf(NodeId nodeId) {
        NodeIdTypeDefinition definition = nodeId.getNodeId();
        if (definition instanceof NodeIdTwoByte) {
            return 0;
        } else if (definition instanceof NodeIdFourByte fourByte) {
            return fourByte.getNamespaceIndex();
        } else if (definition instanceof NodeIdNumeric numeric) {
            return numeric.getNamespaceIndex();
        } else if (definition instanceof NodeIdString string) {
            return string.getNamespaceIndex();
        }
        return -1;
    }

    /** Numeric identifier of a NodeId, or null when the identifier isn't numeric (string/guid/opaque). */
    private static Long numericIdentifierOf(NodeId nodeId) {
        NodeIdTypeDefinition definition = nodeId.getNodeId();
        if (definition instanceof NodeIdTwoByte || definition instanceof NodeIdFourByte
                || definition instanceof NodeIdNumeric) {
            try {
                return Long.parseLong(definition.getIdentifier());
            } catch (NumberFormatException e) {
                return null;
            }
        }
        return null;
    }

    /** Formats the target node of a reference back into an OpcuaTag address (or null if unsupported). */
    private static String addressOf(ExpandedNodeId expandedNodeId) {
        if (expandedNodeId == null) {
            return null;
        }
        NodeIdTypeDefinition nodeId = expandedNodeId.getNodeId();
        String identifier = nodeId.getIdentifier();
        return switch (nodeId) {
            case NodeIdTwoByte ignored -> "ns=0;i=" + identifier;
            case NodeIdFourByte fourByte -> "ns=" + fourByte.getNamespaceIndex() + ";i=" + identifier;
            case NodeIdNumeric numeric -> "ns=" + numeric.getNamespaceIndex() + ";i=" + identifier;
            case NodeIdString string -> "ns=" + string.getNamespaceIndex() + ";s=" + identifier;
            default ->
                // GUID and opaque (ByteString) node identifiers are not round-tripped in Phase 1:
                // NodeIdGuid#getIdentifier() returns the raw byte-array toString(), which the tag
                // parser can't turn back into a usable node id. Such nodes are skipped for now.
                null;
        };
    }

    private static String localizedTextValue(LocalizedText localizedText) {
        if (localizedText == null || !localizedText.getTextSpecified() || localizedText.getText() == null) {
            return null;
        }
        return localizedText.getText().getStringValue();
    }

    private static String qualifiedNameValue(QualifiedName qualifiedName) {
        if (qualifiedName == null || qualifiedName.getName() == null) {
            return null;
        }
        return qualifiedName.getName().getStringValue();
    }

    public static PlcValue variantToPlcValue(PlcTag tag, Variant variant) {
        PlcValue value = null;
        // A node that exists but currently carries no value is encoded as a Null variant
        // (VariantType 0). That's a valid, successful read of an empty value, not an
        // unsupported type, so it maps to PlcNull rather than to null.
        switch (variant) {
            case null -> {
                return new PlcNull();
            }
            case VariantNull ignored -> {
                return new PlcNull();
            }
            case VariantBoolean variantBoolean -> {
                byte[] array = variantBoolean.getValue();
                List<PlcValue> values = new ArrayList<>(array.length);
                for (byte b : array) {
                    values.add(new PlcBOOL(b != 0));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantSByte variantSByte -> {
                byte[] array = variantSByte.getValue();
                List<PlcValue> values = new ArrayList<>(array.length);
                for (byte b : array) {
                    values.add(new PlcSINT(b));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantByte variantByte -> {
                List<Short> array = variantByte.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Short s : array) {
                    values.add(new PlcUSINT(s));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantInt16 variantInt16 -> {
                List<Short> array = variantInt16.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Short s : array) {
                    values.add(new PlcINT(s));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantUInt16 variantUInt16 -> {
                List<Integer> array = variantUInt16.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Integer i : array) {
                    values.add(new PlcUINT(i));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantInt32 variantInt32 -> {
                List<Integer> array = variantInt32.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Integer i : array) {
                    values.add(new PlcDINT(i));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantUInt32 variantUInt32 -> {
                List<Long> array = variantUInt32.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Long l : array) {
                    values.add(new PlcUDINT(l));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantInt64 variantInt64 -> {
                List<Long> array = variantInt64.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Long l : array) {
                    values.add(new PlcLINT(l));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantUInt64 variantUInt64 -> {
                List<BigInteger> array = variantUInt64.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (BigInteger bi : array) {
                    values.add(new PlcULINT(bi));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantFloat variantFloat -> {
                List<Float> array = variantFloat.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Float f : array) {
                    values.add(new PlcREAL(f));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantDouble variantDouble -> {
                List<Double> array = variantDouble.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Double d : array) {
                    values.add(new PlcLREAL(d));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantString variantString -> {
                List<PascalString> stringArray = variantString.getValue();
                List<PlcValue> values = new ArrayList<>(stringArray.size());
                for (PascalString ps : stringArray) {
                    values.add(new PlcSTRING(ps.getStringValue()));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantDateTime variantDateTime -> {
                List<Long> array = variantDateTime.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (Long l : array) {
                    values.add(DefaultPlcValueHandler.of(tag, LocalDateTime.ofInstant(Instant.ofEpochMilli(getDateTime(l)), ZoneOffset.UTC)));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantGuid variantGuid -> {
                List<GuidValue> array = variantGuid.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (GuidValue guidValue : array) {
                    //These two data sections aren't little endian like the rest.
                    byte[] data4Bytes = guidValue.getData4();
                    int data4 = 0;
                    for (byte data4Byte : data4Bytes) {
                        data4 = (data4 << 8) + (data4Byte & 0xff);
                    }
                    byte[] data5Bytes = guidValue.getData5();
                    long data5 = 0;
                    for (byte data5Byte : data5Bytes) {
                        data5 = (data5 << 8) + (data5Byte & 0xff);
                    }
                    values.add(new PlcSTRING(Long.toHexString(guidValue.getData1()) + "-" + Integer.toHexString(guidValue.getData2()) + "-" + Integer.toHexString(guidValue.getData3()) + "-" + Integer.toHexString(data4) + "-" + Long.toHexString(data5)));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantXmlElement variantXmlElement -> {
                List<PascalString> strings = variantXmlElement.getValue();
                List<PlcValue> values = new ArrayList<>(strings.size());
                for (PascalString ps : strings) {
                    values.add(new PlcSTRING(ps.getStringValue()));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantLocalizedText variantLocalizedText -> {
                List<LocalizedText> strings = variantLocalizedText.getValue();
                List<PlcValue> values = new ArrayList<>(strings.size());
                for (LocalizedText lt : strings) {
                    String s = "";
                    s += lt.getLocaleSpecified() ? lt.getLocale().getStringValue() + "|" : "";
                    s += lt.getTextSpecified() ? lt.getText().getStringValue() : "";
                    values.add(new PlcSTRING(s));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantQualifiedName variantQualifiedName -> {
                List<QualifiedName> strings = variantQualifiedName.getValue();
                List<PlcValue> values = new ArrayList<>(strings.size());
                for (QualifiedName qn : strings) {
                    values.add(new PlcSTRING("ns=" + qn.getNamespaceIndex() + ";s=" + qn.getName().getStringValue()));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantExtensionObject variantExtensionObject -> {
                List<ExtensionObject> objects = variantExtensionObject.getValue();
                List<PlcValue> values = new ArrayList<>(objects.size());
                for (ExtensionObject eo : objects) {
                    values.add(new PlcSTRING(eo.toString()));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantNodeId variantNodeId -> {
                List<NodeId> nodeIds = variantNodeId.getValue();
                List<PlcValue> values = new ArrayList<>(nodeIds.size());
                for (NodeId nid : nodeIds) {
                    values.add(new PlcSTRING(nid.toString()));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantStatusCode variantStatusCode -> {
                List<StatusCode> statusCodes = variantStatusCode.getValue();
                List<PlcValue> values = new ArrayList<>(statusCodes.size());
                for (StatusCode sc : statusCodes) {
                    values.add(new PlcSTRING(sc.toString()));
                }
                value = structurePlcValues(values, variant);
            }
            case VariantByteString variantByteString -> {
                List<ByteStringArray> array = variantByteString.getValue();
                List<PlcValue> values = new ArrayList<>(array.size());
                for (ByteStringArray byteStringArray : array) {
                    Short[] tmpValue = byteStringArray.getValue().toArray(new Short[0]);
                    values.add(DefaultPlcValueHandler.of(tag, tmpValue));
                }
                value = structurePlcValues(values, variant);
            }
            default -> {
            }
        }

        // If the tag declares a specific type (via suffix like :TIME, :DATE, etc.),
        // re-interpret the raw numeric value as the correct IEC 61131-3 type.
        if (value != null && tag != null) {
            PlcValueType targetType = tag.getPlcValueType();
            if (targetType != PlcValueType.NULL) {
                value = applyTypeOverride(value, targetType);
            }
        }

        return value;
    }

    /**
     * Recursively applies a type override to a PlcValue.
     * For PlcList values, each element is converted individually.
     * For scalar values, the raw numeric is re-interpreted as the target type.
     */
    private static PlcValue applyTypeOverride(PlcValue value, PlcValueType targetType) {
        // If the value is already a properly typed temporal value (e.g., from DTL
        // conversion), skip the override to avoid corrupting it.
        PlcValueType currentType = value.getPlcValueType();
        // There is nothing to re-interpret for an empty value; converting it would turn a
        // "no value" into a bogus zero-valued TIME/DATE/... .
        if (currentType == targetType || currentType == PlcValueType.NULL || isTemporalType(currentType)) {
            return value;
        }
        if (value instanceof PlcList list) {
            List<PlcValue> converted = new ArrayList<>(list.getLength());
            for (PlcValue item : list.getList()) {
                converted.add(applyTypeOverride(item, targetType));
            }
            return new PlcList(converted);
        }
        long raw = value.getLong();
        return switch (targetType) {
            case TIME -> new PlcTIME(raw);
            case LTIME -> new PlcLTIME(raw);
            case DATE ->
                // S7/IEC value is days since 1990-01-01, PlcDATE expects days since 1970-01-01
                new PlcDATE(raw + IEC_DATE_EPOCH_OFFSET_DAYS);
            case LDATE ->
                // PlcLDATE expects seconds since 1970-01-01
                new PlcLDATE(raw + IEC_DATE_EPOCH_OFFSET_DAYS * 86400L);
            case TIME_OF_DAY ->
                // S7/IEC value is milliseconds since midnight
                new PlcTIME_OF_DAY(LocalTime.ofNanoOfDay(raw * 1_000_000L));
            case LTIME_OF_DAY -> new PlcLTIME_OF_DAY(raw);
            case DATE_AND_TIME ->
                // PlcDATE_AND_TIME expects seconds since 1970-01-01
                new PlcDATE_AND_TIME(raw + IEC_DATE_EPOCH_OFFSET_DAYS * 86400L);
            case DATE_AND_LTIME ->
                // PlcDATE_AND_LTIME expects nanoseconds since 1970-01-01
                new PlcDATE_AND_LTIME(raw + IEC_DATE_EPOCH_OFFSET_DAYS * 86400L * 1_000_000_000L);
            case LDATE_AND_TIME ->
                // PlcLDATE_AND_TIME expects milliseconds since 1970-01-01
                new PlcLDATE_AND_TIME(raw + IEC_DATE_EPOCH_OFFSET_DAYS * 86400L * 1000L);
            default -> value;
        };
    }

    private static boolean isTemporalType(PlcValueType type) {
        return switch (type) {
            case TIME, LTIME, DATE, LDATE, TIME_OF_DAY, LTIME_OF_DAY, DATE_AND_TIME, DATE_AND_LTIME, LDATE_AND_TIME ->
                true;
            default -> false;
        };
    }

    /**
     * Structures a flat list of PlcValues according to the variant's dimensionality:
     * - Single value: returns the scalar PlcValue directly
     * - 1D array (no array dimensions specified): returns a flat PlcList
     * - Multi-dimensional array: returns nested PlcLists matching the declared dimensions
     */
    private static PlcValue structurePlcValues(List<PlcValue> values, Variant variant) {
        if (values.size() == 1) {
            return values.getFirst();
        }
        List<Integer> dimensions = variant.getArrayDimensions();
        if (dimensions == null || dimensions.isEmpty()) {
            return new PlcList(values);
        }
        return buildMultiDimensionalList(values, dimensions);
    }

    /**
     * Recursively partitions a flat list of PlcValues into nested PlcLists
     * according to the given dimensions.
     * E.g., 6 values with dimensions [3, 2] produces a PlcList of 3 PlcLists of 2 values each.
     */
    private static PlcValue buildMultiDimensionalList(List<PlcValue> flatValues, List<Integer> dimensions) {
        if (dimensions.size() <= 1) {
            return new PlcList(flatValues);
        }
        int currentDim = dimensions.getFirst();
        // The array dimensions are raw signed int32 values taken straight off the wire, independent
        // of how many values were actually sent. Only partition when the declared dimension is
        // consistent with the materialized value count: this bounds the work to real received data
        // and avoids both a division by zero (currentDim == 0) and an eager multi-GB allocation /
        // loop for an attacker-supplied dimension (e.g. 0x7fffffff). Otherwise fall back to a flat
        // list, matching the behaviour when no dimensions are declared.
        if (currentDim <= 0 || currentDim > flatValues.size() || flatValues.size() % currentDim != 0) {
            return new PlcList(flatValues);
        }
        List<Integer> remainingDims = dimensions.subList(1, dimensions.size());
        int chunkSize = flatValues.size() / currentDim;
        List<PlcValue> result = new ArrayList<>(currentDim);
        for (int i = 0; i < currentDim; i++) {
            result.add(buildMultiDimensionalList(flatValues.subList(i * chunkSize, (i + 1) * chunkSize), remainingDims));
        }
        return new PlcList(result);
    }

    private static PlcResponseCode mapOpcStatusCode(long opcStatusCode, PlcResponseCode fallback) {
        if (!OpcuaStatusCode.isDefined(opcStatusCode)) {
            return PlcResponseCode.INTERNAL_ERROR;
        }

        OpcuaStatusCode statusCode = OpcuaStatusCode.enumForValue(opcStatusCode);
        if (statusCode == OpcuaStatusCode.Good) {
            return PlcResponseCode.OK;
        } else if (statusCode == OpcuaStatusCode.BadNodeIdUnknown) {
            return PlcResponseCode.NOT_FOUND;
        } else if (statusCode == OpcuaStatusCode.BadTypeMismatch) {
            return PlcResponseCode.INVALID_DATATYPE;
        } else if (statusCode == OpcuaStatusCode.BadNotWritable) {
            return PlcResponseCode.ACCESS_DENIED;
        } else if (statusCode == OpcuaStatusCode.BadUserAccessDenied) {
            return PlcResponseCode.ACCESS_DENIED;
        } else if (statusCode == OpcuaStatusCode.BadAttributeIdInvalid) {
            return PlcResponseCode.INVALID_ADDRESS;
        } else if (statusCode == OpcuaStatusCode.BadIndexRangeNoData) {
            return PlcResponseCode.INVALID_ADDRESS;
        }
        return fallback;
    }

    /**
     * Walks a (possibly nested) {@link PlcList} and produces a flat 1-D
     * {@link PlcList} of leaf values along with the per-level dimensions in
     * outermost-first order. For a 1-D input the returned list is the input
     * itself and {@code dimsOut} stays empty.
     */
    private static PlcList flattenMultidim(PlcList input, List<Integer> dimsOut) {
        if (input.getLength() == 0 || !(input.getIndex(0) instanceof PlcList)) {
            return input;
        }
        List<PlcValue> flat = new ArrayList<>();
        List<Integer> innerDims = new ArrayList<>();
        for (int i = 0; i < input.getLength(); i++) {
            PlcList inner = (PlcList) input.getIndex(i);
            List<Integer> tmp = new ArrayList<>();
            PlcList innerFlat = flattenMultidim(inner, tmp);
            if (i == 0) {
                if (tmp.isEmpty()) {
                    innerDims.add(inner.getLength());
                } else {
                    innerDims.addAll(tmp);
                }
            }
            flat.addAll(innerFlat.getList());
        }
        dimsOut.add(input.getLength());
        dimsOut.addAll(innerDims);
        return new PlcList(flat);
    }

    /**
     * Last-resort write type when neither a ;TYPE suffix nor a server-resolved type is available:
     * infer from the PlcValue itself. Prefer the value's own declared PlcValueType; only when that
     * is absent fall back to the backing Java class — which collapses the Short-backed types
     * (BYTE/USINT/WORD/UINT) into INT, so the server may reject the write with INVALID_DATATYPE.
     */
    private static PlcValueType inferWriteType(PlcValue value) {
        PlcValueType inferred = value.getPlcValueType();
        if (inferred != null && inferred != PlcValueType.NULL) {
            return inferred;
        }
        Object object = value.getObject();
        if (object instanceof Boolean) {
            return PlcValueType.BOOL;
        } else if (object instanceof Byte) {
            return PlcValueType.SINT;
        } else if (object instanceof Short) {
            return PlcValueType.INT;
        } else if (object instanceof Integer) {
            return PlcValueType.DINT;
        } else if (object instanceof Long) {
            return PlcValueType.LINT;
        } else if (object instanceof Float) {
            return PlcValueType.REAL;
        } else if (object instanceof Double) {
            return PlcValueType.LREAL;
        } else if (object instanceof String) {
            return PlcValueType.STRING;
        }
        return PlcValueType.NULL;
    }

    /**
     * Write type for a node whose server DataType is the abstract UInteger supertype: pick the
     * unsigned concrete type matching the value's width (any unsigned subtype is accepted). Falls
     * back to plain inference for non-integer values.
     */
    private static PlcValueType inferUnsignedWriteType(PlcValue value) {
        Object object = value.getObject();
        if (object instanceof Byte) {
            return PlcValueType.USINT;
        } else if (object instanceof Short) {
            return PlcValueType.UINT;
        } else if (object instanceof Integer) {
            return PlcValueType.UDINT;
        } else if (object instanceof Long || object instanceof BigInteger) {
            return PlcValueType.ULINT;
        }
        return inferWriteType(value);
    }

    private Variant fromPlcValue(String tagName, OpcuaTag tag, PlcWriteRequest request, NodeAttributes serverAttributes) {
        PlcList valueObject;
        if (request.getPlcValue(tagName).getObject() instanceof List) {
            valueObject = (PlcList) request.getPlcValue(tagName);
        } else {
            List<PlcValue> list = new ArrayList<>();
            list.add(request.getPlcValue(tagName));
            valueObject = new PlcList(list);
        }

        // OPC UA carries matrices/cubes as a flat value array plus an
        // arrayDimensions field; collapse any nested PlcList structure
        // into one buffer and capture the shape. For 1D input dims stays
        // empty and the rest of the path behaves as before.
        List<Integer> dims = new ArrayList<>();
        valueObject = flattenMultidim(valueObject, dims);
        boolean dimsSpec = !dims.isEmpty();
        Integer noOfDims = dimsSpec ? dims.size() : null;
        List<Integer> arrayDims = dimsSpec ? dims : java.util.Collections.emptyList();

        List<PlcValue> plcValueList = valueObject.getList();
        PlcValueType dataType = tag.getPlcValueType();
        if (dataType.equals(PlcValueType.NULL) || dataType.equals(PlcValueType.List)) {
            // The tag address didn't carry a ;TYPE suffix. Phase 4: prefer the type the server
            // declares for this node (resolved and cached up-front in onWrite) — it is
            // authoritative and, unlike the Java-value guess below, unambiguously distinguishes
            // the Short-backed types (BYTE/USINT/WORD/UINT/INT).
            OpcuaDataType serverType = serverAttributes != null ? serverAttributes.dataType() : null;
            if (serverType != null && serverType != OpcuaDataType.NULL) {
                dataType = PlcValueType.valueOf(serverType.name());
            } else if (serverAttributes != null && serverAttributes.unsigned()) {
                // Node is the abstract UInteger type: no concrete built-in to resolve to, so pick
                // an unsigned type sized to the value (an abstract node accepts any subtype).
                dataType = inferUnsignedWriteType(plcValueList.getFirst());
            } else {
                dataType = inferWriteType(plcValueList.getFirst());
            }
        }
        int length = valueObject.getLength();
        // When an IndexRange selects part of an array, the written value must itself be an array
        // matching the range — even a single selected element is a 1-element array, not a scalar.
        boolean arraySpecified = length > 1 || tag.getIndexRange() != null;
        Integer arrayLength = arraySpecified ? length : null;
        switch (dataType) {
            // Simple boolean values
            case BOOL:
                byte[] tmpBOOL = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpBOOL[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantBoolean(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpBOOL);

            // 8-Bit Bit-Strings (Groups of Boolean Values)
            case BYTE:
                List<Short> tmpBYTE = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpBYTE.add(valueObject.getIndex(i).getShort());
                }
                return new VariantByte(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpBYTE);

            // 16-Bit Bit-Strings (Groups of Boolean Values)
            case WORD:
                List<Integer> tmpWORD = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpWORD.add(valueObject.getIndex(i).getInteger());
                }
                return new VariantUInt16(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpWORD);

            // 32-Bit Bit-Strings (Groups of Boolean Values)
            case DWORD:
                List<Long> tmpDWORD = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDWORD.add(valueObject.getIndex(i).getLong());
                }
                return new VariantUInt32(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpDWORD);

            // 64-Bit Bit-Strings (Groups of Boolean Values)
            case LWORD:
                List<BigInteger> tmpLWORD = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLWORD.add(valueObject.getIndex(i).getBigInteger());
                }
                return new VariantUInt64(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpLWORD);

            // 8-Bit Unsigned Integers
            case USINT:
                List<Short> tmpUSINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUSINT.add(valueObject.getIndex(i).getShort());
                }
                return new VariantByte(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpUSINT);

            // 8-Bit Signed Integers
            case SINT:
                byte[] tmpSINT = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpSINT[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantSByte(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpSINT);

            // 16-Bit Unsigned Integers
            case UINT:
                List<Integer> tmpUINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUINT.add(valueObject.getIndex(i).getInt());
                }
                return new VariantUInt16(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpUINT);

            // 16-Bit Signed Integers
            case INT:
                List<Short> tmpINT16 = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpINT16.add(valueObject.getIndex(i).getShort());
                }
                return new VariantInt16(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpINT16);

            // 32-Bit Unsigned Integers
            case UDINT:
                List<Long> tmpUDINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUDINT.add(valueObject.getIndex(i).getLong());
                }
                return new VariantUInt32(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpUDINT);

            // 32-Bit Signed Integers
            case DINT:
                List<Integer> tmpDINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDINT.add(valueObject.getIndex(i).getInt());
                }
                return new VariantInt32(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpDINT);

            // 64-Bit Unsigned Integers
            case ULINT:
                List<BigInteger> tmpULINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpULINT.add(valueObject.getIndex(i).getBigInteger());
                }
                return new VariantUInt64(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpULINT);

            // 64-Bit Signed Integers
            case LINT:
                List<Long> tmpLINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLINT.add(valueObject.getIndex(i).getLong());
                }
                return new VariantInt64(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpLINT);

            // 32-Bit Floating Point Values
            case REAL:
                List<Float> tmpREAL = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpREAL.add(valueObject.getIndex(i).getFloat());
                }
                return new VariantFloat(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpREAL);

            // 64-Bit Floating Point Values
            case LREAL:
                List<Double> tmpLREAL = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLREAL.add(valueObject.getIndex(i).getDouble());
                }
                return new VariantDouble(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpLREAL);

            // UTF-8 Characters and Strings
            case CHAR:
            case STRING:

                // UTF-16 Characters and Strings
            case WCHAR:
            case WSTRING:
                List<PascalString> tmpString = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    String s = valueObject.getIndex(i).getString();
                    tmpString.add(new PascalString(s));
                }
                return new VariantString(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpString);

            case DATE_AND_TIME:
                List<Long> tmpDateTime = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDateTime.add(valueObject.getIndex(i).getDateTime().toEpochSecond(ZoneOffset.UTC));
                }
                return new VariantDateTime(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpDateTime);

            // IEC 61131-3 TIME is modelled by S7-1500 OPC UA as a signed
            // 32-bit integer holding the duration in milliseconds.
            case TIME:
                List<Integer> tmpTime = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpTime.add((int) valueObject.getIndex(i).getDuration().toMillis());
                }
                return new VariantInt32(arraySpecified,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    arrayLength,
                    tmpTime);
            default:
                throw new PlcRuntimeException("Unsupported write tag type " + dataType);
        }
    }

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        LOGGER.trace("Writing Value");
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        if (sendableTagNames(request).isEmpty()) {
            return CompletableFuture.completedFuture(
                new DefaultPlcWriteResponse(request, rejectedWriteTags(request)));
        }

        // Phase 4: for tags without an explicit ;TYPE suffix, resolve the server-declared data
        // type (via the session type cache) up-front so the write is built with the authoritative
        // OPC UA type instead of a lossy Java-value guess. Then assemble and submit the write.
        return resolveWriteTypes(request).thenCompose(serverAttributes ->
            resolveStructWriteInfos(request).thenCompose(structInfos -> {
                RequestHeader requestHeader = conversation.createRequestHeader();
                List<WriteValue> writeValueList = new ArrayList<>(request.getNumberOfTags());
                for (String tagName : sendableTagNames(request)) {
                    OpcuaTag tag = (OpcuaTag) request.getTag(tagName);

                    NodeId nodeId = generateNodeId(tag);

                    Variant variant;
                    if (isStructValue(request.getPlcValue(tagName))) {
                        // Phase 5d: encode a PlcStruct back into a custom-struct ExtensionObject.
                        StructWriteInfo structInfo = structInfos.get(tagName);
                        if (structInfo == null) {
                            throw new PlcRuntimeException("Cannot resolve the structure layout to write "
                                + "tag '" + tagName + "'");
                        }
                        variant = encodeStructVariant(request.getPlcValue(tagName), structInfo);
                    } else {
                        variant = fromPlcValue(tagName, tag, writeRequest, serverAttributes.get(tagName));
                    }

                    writeValueList.add(new WriteValue(nodeId,
                        tag.getAttributeId().getValue(),
                        indexRangeOf(tag),
                        new DataValue(false, false, false, false, false, true,
                            variant, null, null, null, null, null)));
                }

                WriteRequest opcuaWriteRequest = new WriteRequest(requestHeader, writeValueList);

                return conversation.submit(opcuaWriteRequest, WriteResponse.class)
                    .thenApply(response -> writeResponse(request, response));
            }));
    }

    /**
     * Resolves the server-declared attributes for every write tag that lacks an explicit ;TYPE
     * suffix, returning a map (tagName -&gt; {@link NodeAttributes}). Suffix-typed tags are skipped
     * (the suffix wins) and failed lookups are simply omitted, so {@link #fromPlcValue} falls back
     * to Java-value inference for them.
     */
    private CompletableFuture<Map<String, NodeAttributes>> resolveWriteTypes(DefaultPlcWriteRequest request) {
        Map<String, NodeAttributes> serverAttributes = new ConcurrentHashMap<>();
        List<CompletableFuture<?>> futures = new ArrayList<>();
        for (String tagName : sendableTagNames(request)) {
            OpcuaTag tag = (OpcuaTag) request.getTag(tagName);
            // An explicit ;TYPE suffix is authoritative — no server round-trip needed.
            if (tag.getDataType() != OpcuaDataType.NULL) {
                continue;
            }
            futures.add(resolveNodeAttributes(tag)
                .thenAccept(attributes -> {
                    if (attributes != null) {
                        serverAttributes.put(tagName, attributes);
                    }
                })
                .exceptionally(error -> {
                    LOGGER.debug("Could not resolve server type for write tag '{}'; "
                        + "falling back to value inference", tagName, error);
                    return null;
                }));
        }
        return CompletableFuture.allOf(futures.toArray(new CompletableFuture[0]))
            .thenApply(v -> serverAttributes);
    }

    private PlcWriteResponse writeResponse(DefaultPlcWriteRequest request, WriteResponse writeResponse) {
        Map<String, PlcResponseCode> responseMap = new HashMap<>(rejectedWriteTags(request));
        List<StatusCode> results = writeResponse.getResults();
        // Only the tags that were actually sent have a slot in the response.
        List<String> sendable = sendableTagNames(request);
        Iterator<String> responseIterator = sendable.iterator();
        for (int i = 0; i < sendable.size(); i++) {
            String tagName = responseIterator.next();
            long opcStatusCode = results.get(i).getStatusCode();
            PlcResponseCode statusCode = mapOpcStatusCode(opcStatusCode, PlcResponseCode.REMOTE_ERROR);
            responseMap.put(tagName, statusCode);
        }
        return new DefaultPlcWriteResponse(request, responseMap);
    }


    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        List<String> tagNames = new ArrayList<>(subscriptionRequest.getTagNames());
        // A subscription has a single publishing interval, so take the fastest rate any of its tags
        // asked for - taking the first tag's rate would starve every faster tag behind it. The
        // individual sampling rates are set per monitored item (see onSubscribeCreateMonitoredItemsRequest).
        long cycleTime = tagNames.stream()
            .map(subscriptionRequest::getTag)
            .map(tag -> tag.getDuration().orElse(DEFAULT_CYCLE_TIME))
            .min(Duration::compareTo)
            .orElse(DEFAULT_CYCLE_TIME)
            .toMillis();

        // A configurable server-side queue is useful when the sampling rate can exceed the publishing
        // rate - e.g. a change-of-state subscription that publishes once per second but samples at 0
        // (the server's fastest practical rate). The queue then retains the intermediate values that
        // accumulate between publishes instead of collapsing them to the latest.
        long queueSize = configuration.getSubscriptionQueueSize();

        return onSubscribeCreateSubscription(cycleTime).thenApply(response -> {
                long subscriptionId = response.getSubscriptionId();
                // The server may revise the publishing interval; everything we time against it
                // (the publish request cadence and its timeouts) has to follow the revised value,
                // not what we asked for.
                long revisedCycleTime = revisedPublishingInterval(response, cycleTime);
                OpcuaSubscriptionHandle handle = new OpcuaSubscriptionHandle(this,
                    conversation, subscriptionRequest, subscriptionId, cycleTime, revisedCycleTime, queueSize);
                if (subscriptionRequest.getConsumer() != null) {
                    handle.register(subscriptionRequest.getConsumer());
                }
                subscriptionRequest.getTagNames().forEach(tagName -> {
                    Consumer<PlcSubscriptionEvent> tagConsumer = subscriptionRequest.getTagConsumer(tagName);
                    if (tagConsumer != null) {
                        handle.registerTagConsumer(tagName, tagConsumer);
                    }
                });
                subscriptions.put(handle.getSubscriptionId(), handle);
                return handle;
            })
            .thenCompose(OpcuaSubscriptionHandle::onSubscribeCreateMonitoredItemsRequest)
            .thenApply(handle -> {
                Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
                for (String tagName : subscriptionRequest.getTagNames()) {
                    final DefaultPlcSubscriptionTag tagDefaultPlcSubscription = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);
                    if (!(tagDefaultPlcSubscription.getTag() instanceof OpcuaTag)) {
                        values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                    } else {
                        values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
                    }
                }

                return new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
            });
    }

    /**
     * The publishing interval the server actually granted, in milliseconds. Servers are free to
     * revise the interval we requested (typically upwards, to their smallest supported rate).
     * Falls back to the requested value if the server reports something unusable.
     */
    private static long revisedPublishingInterval(CreateSubscriptionResponse response, long requestedCycleTime) {
        double revised = response.getRevisedPublishingInterval();
        if (revised <= 0) {
            return requestedCycleTime;
        }
        long revisedMillis = Math.round(revised);
        if (revisedMillis != requestedCycleTime) {
            LOGGER.debug("Server revised the publishing interval from {}ms to {}ms", requestedCycleTime, revisedMillis);
        }
        return revisedMillis;
    }

    private CompletableFuture<CreateSubscriptionResponse> onSubscribeCreateSubscription(long cycleTime) {
        LOGGER.trace("Entering creating subscription request");

        RequestHeader requestHeader = conversation.createRequestHeader();
        CreateSubscriptionRequest createSubscriptionRequest = new CreateSubscriptionRequest(
            requestHeader,
            (double) cycleTime,
            12000L,
            5L,
            65536L,
            true,
            (short) 0
        );

        return conversation.submit(createSubscriptionRequest, CreateSubscriptionResponse.class);
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        unsubscriptionRequest.getSubscriptionHandles().forEach(o -> {
            OpcuaSubscriptionHandle opcuaSubHandle = (OpcuaSubscriptionHandle) o;
            opcuaSubHandle.stopSubscriber();
        });
        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
    }

    public void removeSubscription(Long subscriptionId) {
        subscriptions.remove(subscriptionId);
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        // Register the current consumer for each of the given subscription handles
        for (PlcSubscriptionHandle subscriptionHandle : handles) {
            LOGGER.debug("Registering Consumer");
            subscriptionHandle.register(consumer);
        }
        DefaultPlcConsumerRegistration registration =
            new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(registration, consumer);
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        registration.unregister();
        consumers.remove(registration);
    }

    public static long getDateTime(long dateTime) {
        return (dateTime - EPOCH_OFFSET) / 10000;
    }

}
