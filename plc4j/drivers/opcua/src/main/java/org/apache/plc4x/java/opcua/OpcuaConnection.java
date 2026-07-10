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
import java.util.concurrent.TimeoutException;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.metadata.time.TimeSource;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.context.Conversation;
import org.apache.plc4x.java.opcua.context.OpcuaDriverContext;
import org.apache.plc4x.java.opcua.context.OpcuaWire;
import org.apache.plc4x.java.opcua.context.SecureChannel;
import org.apache.plc4x.java.opcua.protocol.OpcuaSubscriptionHandle;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.tag.OpcuaPlcTagHandler;
import org.apache.plc4x.java.opcua.tag.OpcuaQualityStatus;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.transport.tcp.TcpTransportInstance;

import java.net.InetSocketAddress;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Predicate;

public class OpcuaConnection extends ConnectionBase<OpcuaConfiguration> implements OpcuaWire {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaConnection.class);
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
            "",
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
        this.channel = new SecureChannel(conversation, driverContext, configuration, getAuthentication());

        try {
            // Discovery only carries information needed for the encrypted modes
            // (it fetches the server certificate). For SecurityPolicy.NONE the
            // GetEndpoints call adds nothing and — running on the same TCP
            // socket as the real handshake — a second HELLO after the discovery
            // CloseSecureChannel makes most servers either reset the connection
            // or hang waiting for nothing.
            if (configuration.isDiscovery()
                && configuration.getSecurityPolicy() != null
                && configuration.getSecurityPolicy() != org.apache.plc4x.java.opcua.security.SecurityPolicy.NONE) {
                LOGGER.debug("Discovering endpoints before connecting");
                EndpointDescription endpoint = channel.onDiscover()
                    .get(configuration.getNegotiationTimeout(), TimeUnit.MILLISECONDS);
                configuration.setServerCertificate(
                    getX509Certificate(endpoint.getServerCertificate().getStringValue()));
            }
            channel.onConnect().get(configuration.getNegotiationTimeout(), TimeUnit.MILLISECONDS);
            connected = true;
            LOGGER.info("Established connection to server");
        } catch (Exception e) {
            throw new PlcConnectionException("Failed to establish OPC UA connection", e);
        }
    }

    private PlcAuthentication getAuthentication() {
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
                channel.onDisconnect();
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

        List<ReadValueId> readValueArray = new ArrayList<>(request.getTagNames().size());
        Iterator<String> iterator = request.getTagNames().iterator();
        Map<String, PlcTag> tagMap = new LinkedHashMap<>();
        for (int i = 0; i < request.getTagNames().size(); i++) {
            String tagName = iterator.next();
            // TODO: We need to check that the tag-return-code is OK as it could also be INVALID_TAG
            OpcuaTag tag = (OpcuaTag) request.getTag(tagName);
            tagMap.put(tagName, tag);

            NodeId nodeId = generateNodeId(tag);

            readValueArray.add(new ReadValueId(nodeId,
                tag.getAttributeId().getValue(),
                NULL_STRING,
                new QualifiedName(0, NULL_STRING)));
        }

        ReadRequest opcuaReadRequest = new ReadRequest(
            requestHeader,
            0.0d,
            TimestampsToReturn.timestampsToReturnBoth,
            readValueArray
        );

        return conversation.submit(opcuaReadRequest, ReadResponse.class).thenApply(response -> {
            Map<String, PlcResponseItem<PlcValue>> mappedResponse = readResponse(tagMap, response.getResults());
            return new DefaultPlcReadResponse(request, mappedResponse);
        });
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

    public Map<String, PlcResponseItem<PlcValue>> readResponse(Map<String, PlcTag> tagMap, List<DataValue> results) {
        Map<String, PlcResponseItem<PlcValue>> response = new HashMap<>();
        int index = 0;
        for (String tagName : tagMap.keySet()) {
            PlcTag tag = tagMap.get(tagName);
            PlcValue value = null;
            DataValue dataValue = results.get(index++);
            PlcResponseCode responseCode = PlcResponseCode.OK;
            if (dataValue.getValueSpecified()) {
                value = variantToPlcValue(tag, dataValue.getValue());
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

    public static PlcValue variantToPlcValue(PlcTag tag, Variant variant) {
        PlcValue value = null;
        if (variant instanceof VariantBoolean) {
            byte[] array = ((VariantBoolean) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.length);
            for (byte b : array) {
                values.add(new PlcBOOL(b != 0));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantSByte) {
            byte[] array = ((VariantSByte) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.length);
            for (byte b : array) {
                values.add(new PlcSINT(b));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantByte) {
            List<Short> array = ((VariantByte) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Short s : array) {
                values.add(new PlcUSINT(s));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantInt16) {
            List<Short> array = ((VariantInt16) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Short s : array) {
                values.add(new PlcINT(s));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantUInt16) {
            List<Integer> array = ((VariantUInt16) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Integer i : array) {
                values.add(new PlcUINT(i));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantInt32) {
            List<Integer> array = ((VariantInt32) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Integer i : array) {
                values.add(new PlcDINT(i));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantUInt32) {
            List<Long> array = ((VariantUInt32) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Long l : array) {
                values.add(new PlcUDINT(l));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantInt64) {
            List<Long> array = ((VariantInt64) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Long l : array) {
                values.add(new PlcLINT(l));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantUInt64) {
            List<BigInteger> array = ((VariantUInt64) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (BigInteger bi : array) {
                values.add(new PlcULINT(bi));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantFloat) {
            List<Float> array = ((VariantFloat) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Float f : array) {
                values.add(new PlcREAL(f));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantDouble) {
            List<Double> array = ((VariantDouble) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Double d : array) {
                values.add(new PlcLREAL(d));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantString) {
            List<PascalString> stringArray = ((VariantString) variant).getValue();
            List<PlcValue> values = new ArrayList<>(stringArray.size());
            for (PascalString ps : stringArray) {
                values.add(new PlcSTRING(ps.getStringValue()));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantDateTime) {
            List<Long> array = ((VariantDateTime) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (Long l : array) {
                values.add(DefaultPlcValueHandler.of(tag, LocalDateTime.ofInstant(Instant.ofEpochMilli(getDateTime(l)), ZoneOffset.UTC)));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantGuid) {
            List<GuidValue> array = ((VariantGuid) variant).getValue();
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
        } else if (variant instanceof VariantXmlElement) {
            List<PascalString> strings = ((VariantXmlElement) variant).getValue();
            List<PlcValue> values = new ArrayList<>(strings.size());
            for (PascalString ps : strings) {
                values.add(new PlcSTRING(ps.getStringValue()));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantLocalizedText) {
            List<LocalizedText> strings = ((VariantLocalizedText) variant).getValue();
            List<PlcValue> values = new ArrayList<>(strings.size());
            for (LocalizedText lt : strings) {
                String s = "";
                s += lt.getLocaleSpecified() ? lt.getLocale().getStringValue() + "|" : "";
                s += lt.getTextSpecified() ? lt.getText().getStringValue() : "";
                values.add(new PlcSTRING(s));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantQualifiedName) {
            List<QualifiedName> strings = ((VariantQualifiedName) variant).getValue();
            List<PlcValue> values = new ArrayList<>(strings.size());
            for (QualifiedName qn : strings) {
                values.add(new PlcSTRING("ns=" + qn.getNamespaceIndex() + ";s=" + qn.getName().getStringValue()));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantExtensionObject) {
            List<ExtensionObject> objects = ((VariantExtensionObject) variant).getValue();
            List<PlcValue> values = new ArrayList<>(objects.size());
            for (ExtensionObject eo : objects) {
                values.add(new PlcSTRING(eo.toString()));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantNodeId) {
            List<NodeId> nodeIds = ((VariantNodeId) variant).getValue();
            List<PlcValue> values = new ArrayList<>(nodeIds.size());
            for (NodeId nid : nodeIds) {
                values.add(new PlcSTRING(nid.toString()));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantStatusCode) {
            List<StatusCode> statusCodes = ((VariantStatusCode) variant).getValue();
            List<PlcValue> values = new ArrayList<>(statusCodes.size());
            for (StatusCode sc : statusCodes) {
                values.add(new PlcSTRING(sc.toString()));
            }
            value = structurePlcValues(values, variant);
        } else if (variant instanceof VariantByteString) {
            List<ByteStringArray> array = ((VariantByteString) variant).getValue();
            List<PlcValue> values = new ArrayList<>(array.size());
            for (ByteStringArray byteStringArray : array) {
                Short[] tmpValue = byteStringArray.getValue().toArray(new Short[0]);
                values.add(DefaultPlcValueHandler.of(tag, tmpValue));
            }
            value = structurePlcValues(values, variant);
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
        if (currentType == targetType || isTemporalType(currentType)) {
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

    private Variant fromPlcValue(String tagName, OpcuaTag tag, PlcWriteRequest request) {
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
            // When the tag address didn't carry a type suffix, infer from the
            // PlcValue's own type. Falling back on getObject().getClass()
            // collapses BYTE/USINT/WORD/UINT (all back-by Short) into INT,
            // so the server rejects the write with INVALID_DATATYPE.
            PlcValueType inferred = plcValueList.get(0).getPlcValueType();
            if (inferred != null && inferred != PlcValueType.NULL) {
                dataType = inferred;
            } else if (plcValueList.get(0).getObject() instanceof Boolean) {
                dataType = PlcValueType.BOOL;
            } else if (plcValueList.get(0).getObject() instanceof Byte) {
                dataType = PlcValueType.SINT;
            } else if (plcValueList.get(0).getObject() instanceof Short) {
                dataType = PlcValueType.INT;
            } else if (plcValueList.get(0).getObject() instanceof Integer) {
                dataType = PlcValueType.DINT;
            } else if (plcValueList.get(0).getObject() instanceof Long) {
                dataType = PlcValueType.LINT;
            } else if (plcValueList.get(0).getObject() instanceof Float) {
                dataType = PlcValueType.REAL;
            } else if (plcValueList.get(0).getObject() instanceof Double) {
                dataType = PlcValueType.LREAL;
            } else if (plcValueList.get(0).getObject() instanceof String) {
                dataType = PlcValueType.STRING;
            }
        }
        int length = valueObject.getLength();
        switch (dataType) {
            // Simple boolean values
            case BOOL:
                byte[] tmpBOOL = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpBOOL[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantBoolean(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpBOOL);

            // 8-Bit Bit-Strings (Groups of Boolean Values)
            case BYTE:
                List<Short> tmpBYTE = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpBYTE.add(valueObject.getIndex(i).getShort());
                }
                return new VariantByte(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpBYTE);

            // 16-Bit Bit-Strings (Groups of Boolean Values)
            case WORD:
                List<Integer> tmpWORD = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpWORD.add(valueObject.getIndex(i).getInteger());
                }
                return new VariantUInt16(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpWORD);

            // 32-Bit Bit-Strings (Groups of Boolean Values)
            case DWORD:
                List<Long> tmpDWORD = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDWORD.add(valueObject.getIndex(i).getLong());
                }
                return new VariantUInt32(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpDWORD);

            // 64-Bit Bit-Strings (Groups of Boolean Values)
            case LWORD:
                List<BigInteger> tmpLWORD = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLWORD.add(valueObject.getIndex(i).getBigInteger());
                }
                return new VariantUInt64(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpLWORD);

            // 8-Bit Unsigned Integers
            case USINT:
                List<Short> tmpUSINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUSINT.add(valueObject.getIndex(i).getShort());
                }
                return new VariantByte(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpUSINT);

            // 8-Bit Signed Integers
            case SINT:
                byte[] tmpSINT = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpSINT[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantSByte(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpSINT);

            // 16-Bit Unsigned Integers
            case UINT:
                List<Integer> tmpUINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUINT.add(valueObject.getIndex(i).getInt());
                }
                return new VariantUInt16(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpUINT);

            // 16-Bit Signed Integers
            case INT:
                List<Short> tmpINT16 = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpINT16.add(valueObject.getIndex(i).getShort());
                }
                return new VariantInt16(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpINT16);

            // 32-Bit Unsigned Integers
            case UDINT:
                List<Long> tmpUDINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUDINT.add(valueObject.getIndex(i).getLong());
                }
                return new VariantUInt32(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpUDINT);

            // 32-Bit Signed Integers
            case DINT:
                List<Integer> tmpDINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDINT.add(valueObject.getIndex(i).getInt());
                }
                return new VariantInt32(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpDINT);

            // 64-Bit Unsigned Integers
            case ULINT:
                List<BigInteger> tmpULINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpULINT.add(valueObject.getIndex(i).getBigInteger());
                }
                return new VariantUInt64(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpULINT);

            // 64-Bit Signed Integers
            case LINT:
                List<Long> tmpLINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLINT.add(valueObject.getIndex(i).getLong());
                }
                return new VariantInt64(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpLINT);

            // 32-Bit Floating Point Values
            case REAL:
                List<Float> tmpREAL = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpREAL.add(valueObject.getIndex(i).getFloat());
                }
                return new VariantFloat(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpREAL);

            // 64-Bit Floating Point Values
            case LREAL:
                List<Double> tmpLREAL = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLREAL.add(valueObject.getIndex(i).getDouble());
                }
                return new VariantDouble(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
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
                return new VariantString(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpString);

            case DATE_AND_TIME:
                List<Long> tmpDateTime = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDateTime.add(valueObject.getIndex(i).getDateTime().toEpochSecond(ZoneOffset.UTC));
                }
                return new VariantDateTime(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpDateTime);

            // IEC 61131-3 TIME is modelled by S7-1500 OPC UA as a signed
            // 32-bit integer holding the duration in milliseconds.
            case TIME:
                List<Integer> tmpTime = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpTime.add((int) valueObject.getIndex(i).getDuration().toMillis());
                }
                return new VariantInt32(length != 1,
                    dimsSpec,
                    noOfDims,
                    arrayDims,
                    length == 1 ? null : length,
                    tmpTime);
            default:
                throw new PlcRuntimeException("Unsupported write tag type " + dataType);
        }
    }

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        LOGGER.trace("Writing Value");
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        RequestHeader requestHeader = conversation.createRequestHeader();
        List<WriteValue> writeValueList = new ArrayList<>(request.getTagNames().size());
        for (String tagName : request.getTagNames()) {
            OpcuaTag tag = (OpcuaTag) request.getTag(tagName);

            NodeId nodeId = generateNodeId(tag);

            writeValueList.add(new WriteValue(nodeId,
                tag.getAttributeId().getValue(),
                NULL_STRING,
                new DataValue(
                    false,
                    false,
                    false,
                    false,
                    false,
                    true,
                    fromPlcValue(tagName, tag, writeRequest),
                    null,
                    null,
                    null,
                    null,
                    null)));
        }

        WriteRequest opcuaWriteRequest = new WriteRequest(requestHeader, writeValueList);

        return conversation.submit(opcuaWriteRequest, WriteResponse.class)
            .thenApply(response -> writeResponse(request, response));
    }

    private PlcWriteResponse writeResponse(DefaultPlcWriteRequest request, WriteResponse writeResponse) {
        Map<String, PlcResponseCode> responseMap = new HashMap<>();
        List<StatusCode> results = writeResponse.getResults();
        Iterator<String> responseIterator = request.getTagNames().iterator();
        for (int i = 0; i < request.getTagNames().size(); i++) {
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
        long cycleTime = (subscriptionRequest.getTag(tagNames.get(0))).getDuration().orElse(Duration.ofMillis(1000)).toMillis();

        return onSubscribeCreateSubscription(cycleTime).thenApply(response -> {
                long subscriptionId = response.getSubscriptionId();
                OpcuaSubscriptionHandle handle = new OpcuaSubscriptionHandle(this,
                    conversation, subscriptionRequest, subscriptionId, cycleTime);
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
            .thenCompose(handle -> handle.onSubscribeCreateMonitoredItemsRequest())
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

    private GuidValue toGuidValue(String identifier) {
        LOGGER.error("Querying Guid nodes is not supported");
        byte[] data4 = new byte[]{0, 0};
        byte[] data5 = new byte[]{0, 0, 0, 0, 0, 0};
        return new GuidValue(0L, 0, 0, data4, data5);
    }


}
