/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.

 */
package org.apache.plc4x.java.opcua.connection;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.opcua.protocol.OpcuaField;
import org.apache.plc4x.java.opcua.protocol.OpcuaSubsriptionHandle;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.values.*;
import org.eclipse.milo.opcua.sdk.client.OpcUaClient;
import org.eclipse.milo.opcua.sdk.client.api.config.OpcUaClientConfig;
import org.eclipse.milo.opcua.sdk.client.api.identity.AnonymousProvider;
import org.eclipse.milo.opcua.sdk.client.api.identity.IdentityProvider;
import org.eclipse.milo.opcua.sdk.client.api.subscriptions.UaMonitoredItem;
import org.eclipse.milo.opcua.sdk.client.api.subscriptions.UaSubscription;
import org.eclipse.milo.opcua.stack.client.DiscoveryClient;
import org.eclipse.milo.opcua.stack.core.AttributeId;
import org.eclipse.milo.opcua.stack.core.Identifiers;
import org.eclipse.milo.opcua.stack.core.UaException;
import org.eclipse.milo.opcua.stack.core.StatusCodes;
import org.eclipse.milo.opcua.stack.core.security.SecurityPolicy;
import org.eclipse.milo.opcua.stack.core.types.builtin.*;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.*;
import org.eclipse.milo.opcua.stack.core.types.enumerated.MonitoringMode;
import org.eclipse.milo.opcua.stack.core.types.enumerated.TimestampsToReturn;
import org.eclipse.milo.opcua.stack.core.types.structured.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.math.BigDecimal;
import java.util.stream.Collectors;
import java.net.InetAddress;
import java.time.Duration;
import java.time.LocalDateTime;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Predicate;

import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.uint;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ulong;

/**
 * Corresponding implementaion for a TCP-based connection for an OPC UA server.
 * TODO: At the moment are just connections without any security mechanism possible
 * <p>
 */
public class OpcuaTcpPlcConnection extends BaseOpcuaPlcConnection {

    private static final int OPCUA_DEFAULT_TCP_PORT = 4840;

    private static final Logger logger = LoggerFactory.getLogger(OpcuaTcpPlcConnection.class);
    private final AtomicLong clientHandles = new AtomicLong(1L);
    private InetAddress address;
    private int requestTimeout = 5000;
    private int port;
    private String params;
    private OpcUaClient client;
    private boolean isConnected = false;

    private OpcuaTcpPlcConnection(InetAddress address, String params, int requestTimeout) {
        this(address, OPCUA_DEFAULT_TCP_PORT, params, requestTimeout);
        logger.info("Configured OpcuaTcpPlcConnection with: host-name {}", address.getHostAddress());
    }

    private OpcuaTcpPlcConnection(InetAddress address, int port, String params, int requestTimeout) {
        this(params);
        logger.info("Configured OpcuaTcpPlcConnection with: host-name {}", address.getHostAddress());
        this.address = address;
        this.port = port;
        this.params = params;
        this.requestTimeout = requestTimeout;
    }

    private OpcuaTcpPlcConnection(String params) {
        super(getOptionString(params));
    }

    public static OpcuaTcpPlcConnection of(InetAddress address, String params, int requestTimeout) {
        return new OpcuaTcpPlcConnection(address, params, requestTimeout);
    }

    public static OpcuaTcpPlcConnection of(InetAddress address, int port, String params, int requestTimeout) {
        return new OpcuaTcpPlcConnection(address, port, params, requestTimeout);
    }

    public static PlcValue encodePlcValue(DataValue value) {
        ExpandedNodeId typeNode = value.getValue().getDataType().get();
        Object objValue = value.getValue().getValue();

        if (objValue.getClass().isArray()) {
            Object[] objArray = (Object[]) objValue;
            if (objArray[0] instanceof Boolean) {
                Boolean[] obj = (Boolean[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcBOOL((Boolean) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof Integer) {
                Integer[] obj = (Integer[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcDINT((Integer) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof Short) {
                Short[] obj = (Short[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcINT((Short) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof Byte) {
                Byte[] obj = (Byte[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcSINT((Byte) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof Long) {
                Long[] obj = (Long[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcLINT((Long) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof Float) {
                Float[] obj = (Float[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcREAL((Float) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof Double) {
                Double[] obj = (Double[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcLREAL((Double) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else if (objArray[0] instanceof String) {
                String[] obj = (String[]) objValue;
                List<PlcValue> plcValue;
                {
                    int itemCount = (int) obj.length;
                    plcValue = new LinkedList<>();

                    for(int curItem = 0; curItem < itemCount; curItem++) {
                        plcValue.add(new PlcSTRING((String) obj[curItem]));
                    }
                }
                return new PlcList(plcValue);
            } else {
                logger.warn("Node type for " + objArray[0].getClass() + " is not supported");
                return null;
            }

        } else {
            if (typeNode.equals(Identifiers.Boolean)) {
                return new PlcBOOL((Boolean) objValue);
            } else if (typeNode.equals(Identifiers.Integer)) {
                return new PlcDINT((Integer) objValue);
            } else if (typeNode.equals(Identifiers.Int16)) {
                return new PlcINT((Short) objValue);
            } else if (typeNode.equals(Identifiers.Int32)) {
                return new PlcDINT((Integer) objValue);
            } else if (typeNode.equals(Identifiers.Int64)) {
                return new PlcLINT((Long) objValue);
            } else if (typeNode.equals(Identifiers.UInteger)) {
                return new PlcLINT((Long) objValue);
            } else if (typeNode.equals(Identifiers.UInt16)) {
                return new PlcUINT(((UShort) objValue).intValue());
            } else if (typeNode.equals(Identifiers.UInt32)) {
                return new PlcUDINT(((UInteger) objValue).longValue());
            } else if (typeNode.equals(Identifiers.UInt64)) {
                return new PlcULINT(new BigInteger(objValue.toString()));
            } else if (typeNode.equals(Identifiers.Byte)) {
                return new PlcINT(Short.valueOf(objValue.toString()));
            } else if (typeNode.equals(Identifiers.Float)) {
                return new PlcREAL((Float) objValue);
            } else if (typeNode.equals(Identifiers.Double)) {
                return new PlcLREAL((Double) objValue);
            } else if (typeNode.equals(Identifiers.SByte)) {
                return new PlcSINT((Byte) objValue);
            } else {
                return new PlcSTRING(objValue.toString());
            }
        }

    }

    public InetAddress getRemoteAddress() {
        return address;
    }

    @Override
    public void connect() throws PlcConnectionException {
        List<EndpointDescription> endpoints = null;
        EndpointDescription endpoint = null;

        try {
            endpoints = DiscoveryClient.getEndpoints(getEndpointUrl(address, port, getSubPathOfParams(params))).get();
            //TODO Exception should be handeled better when the Discovery-API of Milo is stable
        } catch (Exception ex) {
            logger.info("Failed to discover Endpoint with enabled discovery. If the endpoint does not allow a correct discovery disable this option with the nDiscovery=true option. Failed Endpoint: {}", getEndpointUrl(address, port, params));

            // try the explicit discovery endpoint as well
            String discoveryUrl = getEndpointUrl(address, port, getSubPathOfParams(params));

            if (!discoveryUrl.endsWith("/")) {
                discoveryUrl += "/";
            }
            discoveryUrl += "discovery";

            logger.info("Trying explicit discovery URL: {}", discoveryUrl);
            try {
                endpoints = DiscoveryClient.getEndpoints(discoveryUrl).get();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new PlcConnectionException("Unable to discover URL:" + discoveryUrl);
            } catch (ExecutionException e) {
                throw new PlcConnectionException("Unable to discover URL:" + discoveryUrl);
            }

        }
        endpoint = endpoints.stream()
            .filter(e -> e.getSecurityPolicyUri().equals(getSecurityPolicy().getUri()))
            .filter(endpointFilter())
            .findFirst()
            .orElseThrow(() -> new PlcConnectionException("No desired endpoints from"));

        if (this.skipDiscovery) {
            //ApplicationDescription applicationDescription = new ApplicationDescription();
            //endpoint = new EndpointDescription(address.getHostAddress(),applicationDescription , null, MessageSecurityMode.None, SecurityPolicy.None.getUri(), null , TransportProfile.TCP_UASC_UABINARY.getUri(), UByte.valueOf(0));// TODO hier machen wenn fertig
            ApplicationDescription currentAD = endpoint.getServer();
            ApplicationDescription withoutDiscoveryAD = new ApplicationDescription(
                currentAD.getApplicationUri(),
                currentAD.getProductUri(),
                currentAD.getApplicationName(),
                currentAD.getApplicationType(),
                currentAD.getGatewayServerUri(),
                currentAD.getDiscoveryProfileUri(),
                new String[0]);
            //try to replace the overhanded address
            //any error will result in the overhanded address of the client
            String newEndpointUrl = endpoint.getEndpointUrl(), prefix = "", suffix = "";
            String splitterPrefix = "://";
            String splitterSuffix = ":";
            String[] prefixSplit = newEndpointUrl.split(splitterPrefix);
            if (prefixSplit.length > 1) {
                String[] suffixSplit = prefixSplit[1].split(splitterSuffix);
                //reconstruct the uri
                newEndpointUrl = "";
                newEndpointUrl += prefixSplit[0] + splitterPrefix + address.getHostAddress();
                for (int suffixCounter = 1; suffixCounter < suffixSplit.length; suffixCounter++) {
                    newEndpointUrl += splitterSuffix + suffixSplit[suffixCounter];
                }
                // attach surounding prefix match
                for (int prefixCounter = 2; prefixCounter < prefixSplit.length; prefixCounter++) {
                    newEndpointUrl += splitterPrefix + prefixSplit[prefixCounter];
                }
            }

            EndpointDescription noDiscoverEndpoint = new EndpointDescription(
                newEndpointUrl,
                withoutDiscoveryAD,
                endpoint.getServerCertificate(),
                endpoint.getSecurityMode(),
                endpoint.getSecurityPolicyUri(),
                endpoint.getUserIdentityTokens(),
                endpoint.getTransportProfileUri(),
                endpoint.getSecurityLevel());
            endpoint = noDiscoverEndpoint;
        }


        OpcUaClientConfig config = OpcUaClientConfig.builder()
            .setApplicationName(LocalizedText.english("eclipse milo opc-ua client of the apache PLC4X:PLC4J project"))
            .setApplicationUri("urn:eclipse:milo:plc4x:client")
            .setEndpoint(endpoint)
            .setIdentityProvider(getIdentityProvider())
            .setRequestTimeout(UInteger.valueOf(requestTimeout))
            .build();

        try {
            this.client = OpcUaClient.create(config);
            this.client.connect().get();
            isConnected = true;
        } catch (UaException e) {
            isConnected = false;
            String message = (config == null) ? "NULL" : config.toString();
            throw new PlcConnectionException("The given input values are a not valid OPC UA connection configuration [CONFIG]: " + message);
        } catch (InterruptedException e) {
            isConnected = false;
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Error while creation of the connection because of : " + e.getMessage());
        } catch (ExecutionException e) {
            isConnected = false;
            throw new PlcConnectionException("Error while creation of the connection because of : " + e.getMessage());
        }
    }

    @Override
    public boolean isConnected() {
        return client != null && isConnected;
    }

    @Override
    public void close() throws Exception {
        if (client != null) {
            client.disconnect().get();
            isConnected = false;
        }
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        CompletableFuture<PlcSubscriptionResponse> future = CompletableFuture.supplyAsync(() -> {
            Map<String, ResponseItem<PlcSubscriptionHandle>> responseItems = new HashMap<>();
            for (String fieldName : subscriptionRequest.getFieldNames()) {
                final DefaultPlcSubscriptionField subscriptionField = (DefaultPlcSubscriptionField) subscriptionRequest.getField(fieldName);
                final OpcuaField field = (OpcuaField) Objects.requireNonNull(subscriptionField.getPlcField());
                long cycleTime = subscriptionField.getDuration().orElse(Duration.ofSeconds(1)).toMillis();
                NodeId idNode = generateNodeId(field);
                ReadValueId readValueId = new ReadValueId(
                    idNode,
                    AttributeId.Value.uid(), null, QualifiedName.NULL_VALUE);
                UInteger clientHandle = uint(clientHandles.getAndIncrement());

                MonitoringMode monitoringMode;
                switch (subscriptionField.getPlcSubscriptionType()) {
                    case CYCLIC:
                        monitoringMode = MonitoringMode.Sampling;
                        break;
                    case CHANGE_OF_STATE:
                        monitoringMode = MonitoringMode.Reporting;
                        cycleTime = subscriptionField.getDuration().orElse(Duration.ofSeconds(0)).toMillis();
                        break;
                    case EVENT:
                        monitoringMode = MonitoringMode.Reporting;
                        break;
                    default:
                        monitoringMode = MonitoringMode.Reporting;
                }

                MonitoringParameters parameters = new MonitoringParameters(
                    clientHandle,
                    (double) cycleTime,     // sampling interval
                    null,       // filter, null means use default
                    uint(1),   // queue size
                    true        // discard oldest
                );

                PlcSubscriptionHandle subHandle = null;
                PlcResponseCode responseCode = PlcResponseCode.ACCESS_DENIED;
                try {
                    UaSubscription subscription = client.getSubscriptionManager().createSubscription(cycleTime).get();

                    MonitoredItemCreateRequest request = new MonitoredItemCreateRequest(
                        readValueId, monitoringMode, parameters);
                    List<MonitoredItemCreateRequest> requestList = new LinkedList<>();
                    requestList.add(request);
                    OpcuaSubsriptionHandle subsriptionHandle = new OpcuaSubsriptionHandle(fieldName, clientHandle);
                    UaSubscription.ItemCreationCallback onItemCreated =
                        (item, id) -> item.setValueConsumer(subsriptionHandle::onSubscriptionValue);

                    List<UaMonitoredItem> items = subscription.createMonitoredItems(
                        TimestampsToReturn.Both,
                        requestList,
                        onItemCreated
                    ).get();

                    subHandle = subsriptionHandle;
                    responseCode = PlcResponseCode.OK;
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    logger.warn("Unable to subscribe Elements because of: {}", e.getMessage());
                } catch (ExecutionException e) {
                    logger.warn("Unable to subscribe Elements because of: {}", e.getMessage());
                }

                responseItems.put(fieldName, new ResponseItem(responseCode, subHandle));
            }
            return new DefaultPlcSubscriptionResponse(subscriptionRequest, responseItems);
        });

        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        unsubscriptionRequest.getSubscriptionHandles().forEach(o -> {
            OpcuaSubsriptionHandle opcSubHandle = (OpcuaSubsriptionHandle) o;
            try {
                client.getSubscriptionManager().deleteSubscription(opcSubHandle.getClientHandle()).get();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                logger.warn("Unable to unsubscribe Elements because of: {}", e.getMessage());
            } catch (ExecutionException e) {
                logger.warn("Unable to unsubscribe Elements because of: {}", e.getMessage());
            }
        });

        return null;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        List<PlcConsumerRegistration> registrations = new LinkedList<>();
        // Register the current consumer for each of the given subscription handles
        for (PlcSubscriptionHandle subscriptionHandle : handles) {
            final PlcConsumerRegistration consumerRegistration = subscriptionHandle.register(consumer);
            registrations.add(consumerRegistration);
        }

        return new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        registration.unregister();
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = CompletableFuture.supplyAsync(() -> {
            readRequest.getFields();
            Map<String, ResponseItem<PlcValue>> fields = new HashMap<>();
            List<NodeId> readValueIds = new LinkedList<>();
            List<PlcField> readPLCValues = readRequest.getFields();
            for (PlcField field : readPLCValues) {
                NodeId idNode = generateNodeId((OpcuaField) field);
                readValueIds.add(idNode);
            }

            CompletableFuture<List<DataValue>> dataValueCompletableFuture = client.readValues(0.0, TimestampsToReturn.Both, readValueIds);
            List<DataValue> readValues = null;
            try {
                readValues = dataValueCompletableFuture.get();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                logger.warn("Unable to read Elements because of: {}", e.getMessage());
            } catch (ExecutionException e) {
                logger.warn("Unable to read Elements because of: {}", e.getMessage());
            }
            for (int counter = 0; counter < readValueIds.size(); counter++) {
                PlcResponseCode resultCode = PlcResponseCode.OK;
                PlcValue stringItem = null;
                if (readValues == null || readValues.size() <= counter ||
                    !readValues.get(counter).getStatusCode().equals(StatusCode.GOOD)) {
                    resultCode = PlcResponseCode.NOT_FOUND;
                } else {
                    stringItem = encodePlcValue(readValues.get(counter));

                }
                ResponseItem<PlcValue> newPair = new ResponseItem<>(resultCode, stringItem);
                fields.put((String) readRequest.getFieldNames().toArray()[counter], newPair);


            }
            return new DefaultPlcReadResponse(readRequest, fields);
        });

        return future;
    }


    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future;
        future = CompletableFuture.supplyAsync(() -> {
            List<PlcField> writePLCValues = writeRequest.getFields();
            LinkedList<DataValue> values = new LinkedList<>();
            LinkedList<NodeId> ids = new LinkedList<>();
            LinkedList<String> names = new LinkedList<>();
            Map<String, PlcResponseCode> fieldResponse = new HashMap<>();
            for (String fieldName : writeRequest.getFieldNames()) {
                OpcuaField uaField = (OpcuaField) writeRequest.getField(fieldName);
                NodeId idNode = generateNodeId(uaField);
                Object valueObject = writeRequest.getPlcValue(fieldName).getObject();
                // Added small work around for handling BigIntegers as input type for UInt64
                if (valueObject instanceof BigInteger) valueObject = ulong((BigInteger) valueObject);
                Variant var = null;
                if (valueObject instanceof ArrayList) {
                    List<PlcValue> plcValueList = (List<PlcValue>) valueObject;
                    String dataType = uaField.getPlcDataType();
                    if (dataType.equals("NULL")) {
                        if (plcValueList.get(0).getObject() instanceof Boolean) {
                            dataType = "BOOL";
                        } else if (plcValueList.get(0).getObject() instanceof Byte) {
                            dataType = "SINT";
                        } else if (plcValueList.get(0).getObject() instanceof Short) {
                            dataType = "INT";
                        } else if (plcValueList.get(0).getObject() instanceof Integer) {
                            dataType = "DINT";
                        } else if (plcValueList.get(0).getObject() instanceof Long) {
                            dataType = "LINT";
                        } else if (plcValueList.get(0).getObject() instanceof Float) {
                            dataType = "REAL";
                        } else if (plcValueList.get(0).getObject() instanceof Double) {
                            dataType = "LREAL";
                        } else if (plcValueList.get(0).getObject() instanceof String) {
                            dataType = "STRING";
                        }
                    }
                    switch (dataType) {
                        case "BOOL":
                        case "BIT":
                            List<Boolean> booleanList = (plcValueList).stream().map(
                                    x -> ((PlcBOOL) x).getBoolean()).collect(Collectors.toList());
                            var = new Variant(booleanList.toArray(new Boolean[booleanList.size()]));
                            break;
                        case "BYTE":
                        case "BITARR8":
                            List<UByte> byteList = (plcValueList).stream().map(
                                    x -> UByte.valueOf(((PlcBYTE) x).getShort())).collect(Collectors.toList());
                            var = new Variant(byteList.toArray(new UByte[byteList.size()]));
                            break;
                        case "SINT":
                        case "INT8":
                            List<Byte> sintList = (plcValueList).stream().map(
                                    x -> ((PlcSINT) x).getByte()).collect(Collectors.toList());
                            var = new Variant(sintList.toArray(new Byte[sintList.size()]));
                            break;
                        case "USINT":
                        case "UINT8":
                        case "BIT8":
                            List<UByte> usintList = (plcValueList).stream().map(
                                    x -> UByte.valueOf(((PlcUSINT) x).getShort())).collect(Collectors.toList());
                            var = new Variant(usintList.toArray(new UByte[usintList.size()]));
                            break;
                        case "INT":
                        case "INT16":
                            List<Short> intList = (plcValueList).stream().map(
                                    x -> ((PlcINT) x).getShort()).collect(Collectors.toList());
                            var = new Variant(intList.toArray(new Short[intList.size()]));
                            break;
                        case "UINT":
                        case "UINT16":
                            List<UShort> uintList = (plcValueList).stream().map(
                                    x -> UShort.valueOf(((PlcUINT) x).getInteger())).collect(Collectors.toList());
                            var = new Variant(uintList.toArray(new UShort[uintList.size()]));
                            break;
                        case "WORD":
                        case "BITARR16":
                            List<UShort> wordList = (plcValueList).stream().map(
                                    x -> UShort.valueOf(((PlcWORD) x).getInteger())).collect(Collectors.toList());
                            var = new Variant(wordList.toArray(new UShort[wordList.size()]));
                            break;
                        case "DINT":
                        case "INT32":
                            List<Integer> dintList = (plcValueList).stream().map(
                                    x -> ((PlcDINT) x).getInteger()).collect(Collectors.toList());
                            var = new Variant(dintList.toArray(new Integer[dintList.size()]));
                            break;
                        case "UDINT":
                        case "UINT32":
                            List<UInteger> udintList = (plcValueList).stream().map(
                                    x -> UInteger.valueOf(((PlcUDINT) x).getLong())).collect(Collectors.toList());
                            var = new Variant(udintList.toArray(new UInteger[udintList.size()]));
                            break;
                        case "DWORD":
                        case "BITARR32":
                            List<UInteger> dwordList = (plcValueList).stream().map(
                                    x -> UInteger.valueOf(((PlcDWORD) x).getLong())).collect(Collectors.toList());
                            var = new Variant(dwordList.toArray(new UInteger[dwordList.size()]));
                            break;
                        case "LINT":
                        case "INT64":
                            List<Long> lintList = (plcValueList).stream().map(
                                    x -> ((PlcLINT) x).getLong()).collect(Collectors.toList());
                            var = new Variant(lintList.toArray(new Long[lintList.size()]));
                            break;
                        case "ULINT":
                        case "UINT64":
                            List<ULong> ulintList = (plcValueList).stream().map(
                                    x -> ULong.valueOf(((PlcULINT) x).getBigInteger())).collect(Collectors.toList());
                            var = new Variant(ulintList.toArray(new ULong[ulintList.size()]));
                            break;
                        case "LWORD":
                        case "BITARR64":
                            List<ULong> lwordList = (plcValueList).stream().map(
                                    x -> ULong.valueOf(((PlcLWORD) x).getBigInteger())).collect(Collectors.toList());
                            var = new Variant(lwordList.toArray(new ULong[lwordList.size()]));
                            break;
                        case "REAL":
                        case "FLOAT":
                            List<Float> realList = (plcValueList).stream().map(
                                    x -> ((PlcREAL) x).getFloat()).collect(Collectors.toList());
                            var = new Variant(realList.toArray(new Float[realList.size()]));
                            break;
                        case "LREAL":
                        case "DOUBLE":
                            List<Double> lrealList = (plcValueList).stream().map(
                                    x -> (Double) ((PlcLREAL) x).getDouble()).collect(Collectors.toList());
                            var = new Variant(lrealList.toArray(new Double[lrealList.size()]));
                            break;
                        case "CHAR":
                            List<String> charList = (plcValueList).stream().map(
                                    x -> ((PlcCHAR) x).getString()).collect(Collectors.toList());
                            var = new Variant(charList.toArray(new String[charList.size()]));
                            break;
                        case "WCHAR":
                            List<String> wcharList = (plcValueList).stream().map(
                                    x -> ((PlcWCHAR) x).getString()).collect(Collectors.toList());
                            var = new Variant(wcharList.toArray(new String[wcharList.size()]));
                            break;
                        case "STRING":
                            List<String> stringList = (plcValueList).stream().map(
                                    x -> ((PlcSTRING) x).getString()).collect(Collectors.toList());
                            var = new Variant(stringList.toArray(new String[stringList.size()]));
                            break;
                        case "WSTRING":
                        case "STRING16":
                            List<String> wstringList = (plcValueList).stream().map(
                                    x -> (String) ((PlcSTRING) x).getString()).collect(Collectors.toList());
                            var = new Variant(wstringList.toArray(new String[wstringList.size()]));
                            break;
                        case "DATE_AND_TIME":
                            List<LocalDateTime> dateTimeList = (plcValueList).stream().map(
                                    x -> ((PlcDATE_AND_TIME) x).getDateTime()).collect(Collectors.toList());
                            var = new Variant(dateTimeList.toArray(new LocalDateTime[dateTimeList.size()]));
                            break;
                        default:
                            logger.warn("Unsupported data type : {}, {}", plcValueList.get(0).getClass(), dataType);
                    }
                } else {
                    String dataType = uaField.getPlcDataType();                    
                    PlcValue plcValue = (PlcValue) writeRequest.getPlcValue(fieldName);

                    if (dataType.equals("NULL")) {
                        if (plcValue.getObject() instanceof Boolean) {
                            dataType = "BOOL";
                        } else if (plcValue.getObject() instanceof Byte) {
                            dataType = "SINT";
                        } else if (plcValue.getObject() instanceof Short) {
                            dataType = "INT";
                        } else if (plcValue.getObject() instanceof Integer) {
                            dataType = "DINT";
                        } else if (plcValue.getObject() instanceof Long) {
                            dataType = "LINT";
                        } else if (plcValue.getObject() instanceof Float) {
                            dataType = "REAL";
                        } else if (plcValue.getObject() instanceof Double) {
                            dataType = "LREAL";
                        } else if (plcValue.getObject() instanceof String) {
                            dataType = "STRING";
                        }
                    }
                    switch (dataType) {
                        case "BOOL":
                        case "BIT":
                            var = new Variant(plcValue.getBoolean());
                            break;
                        case "BYTE":
                        case "BITARR8":
                            var = new Variant(UByte.valueOf(plcValue.getShort()));
                            break;
                        case "SINT":
                        case "INT8":
                            var = new Variant(plcValue.getByte());
                            break;
                        case "USINT":
                        case "UINT8":
                        case "BIT8":
                            var = new Variant(UByte.valueOf(plcValue.getShort()));
                            break;
                        case "INT":
                        case "INT16":
                            var = new Variant(plcValue.getShort());
                            break;
                        case "UINT":
                        case "UINT16":
                            var = new Variant(UShort.valueOf(plcValue.getInteger()));
                            break;
                        case "WORD":
                        case "BITARR16":
                            var = new Variant(UShort.valueOf(plcValue.getInteger()));
                            break;
                        case "DINT":
                        case "INT32":
                            var = new Variant(plcValue.getInteger());
                            break;
                        case "UDINT":
                        case "UINT32":
                            var = new Variant(UInteger.valueOf(plcValue.getLong()));
                            break;
                        case "DWORD":
                        case "BITARR32":
                            var = new Variant(UInteger.valueOf(plcValue.getLong()));
                            break;
                        case "LINT":
                        case "INT64":
                            var = new Variant(plcValue.getLong());
                            break;
                        case "ULINT":
                        case "UINT64":
                            var = new Variant(ULong.valueOf(plcValue.getBigInteger()));
                            break;
                        case "LWORD":
                        case "BITARR64":
                            var = new Variant(ULong.valueOf(plcValue.getBigInteger()));
                            break;
                        case "REAL":
                        case "FLOAT":
                            var = new Variant(plcValue.getFloat());
                            break;
                        case "LREAL":
                        case "DOUBLE":
                            var = new Variant(plcValue.getDouble());
                            break;
                        case "CHAR":
                            var = new Variant(plcValue.getString());
                            break;
                        case "WCHAR":
                            var = new Variant(plcValue.getString());
                            break;
                        case "STRING":
                            var = new Variant(plcValue.getString());
                            break;
                        case "WSTRING":
                        case "STRING16":
                            var = new Variant(plcValue.getString());
                            break;
                        case "DATE_AND_TIME":
                            var = new Variant(plcValue.getDateTime());
                            break;
                        default:
                            logger.warn("Unsupported data type : {}, {}", plcValue.getClass(), dataType);
                    }
                }
                DataValue value = new DataValue(var, StatusCode.GOOD, null, null);
                ids.add(idNode);
                names.add(fieldName);
                values.add(value);
            }
            CompletableFuture<List<StatusCode>> opcRequest =
                client.writeValues(ids, values);
            List<StatusCode> statusCodes = null;
            try {
                statusCodes = opcRequest.get();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                statusCodes = new LinkedList<>();
                for (int counter = 0; counter < ids.size(); counter++) {
                    ((LinkedList<StatusCode>) statusCodes).push(StatusCode.BAD);
                }
            } catch (ExecutionException e) {
                statusCodes = new LinkedList<>();
                for (int counter = 0; counter < ids.size(); counter++) {
                    ((LinkedList<StatusCode>) statusCodes).push(StatusCode.BAD);
                }
            }

            for (int counter = 0; counter < names.size(); counter++) {
                final PlcResponseCode resultCode;
                if (statusCodes != null && statusCodes.size() > counter) {
                    Optional<String[]> status = StatusCodes.lookup(statusCodes.get(counter).getValue());
                    if (status.isPresent()) {
                        if (status.get()[0].equals("Good")) {
                            resultCode = PlcResponseCode.OK;
                        } else if (status.get()[0].equals("Uncertain")) {
                            resultCode = PlcResponseCode.NOT_FOUND;
                        } else if (status.get()[0].equals("Bad")) {
                            resultCode = PlcResponseCode.INVALID_DATATYPE;
                        } else if (status.get()[0].equals("Bad_NodeIdUnknown")) {
                            resultCode = PlcResponseCode.NOT_FOUND;
                        } else {
                            resultCode = PlcResponseCode.ACCESS_DENIED;
                        }
                    } else {
                        resultCode = PlcResponseCode.ACCESS_DENIED;
                    }
                } else {
                    resultCode = PlcResponseCode.ACCESS_DENIED;
                }
                fieldResponse.put(names.get(counter), resultCode);
            }
            PlcWriteResponse response = new DefaultPlcWriteResponse(writeRequest, fieldResponse);
            return response;
        });


        return future;
    }


    private NodeId generateNodeId(OpcuaField uaField) {
        NodeId idNode = null;
        switch (uaField.getIdentifierType()) {
            case STRING_IDENTIFIER:
                idNode = new NodeId(uaField.getNamespace(), uaField.getIdentifier());
                break;
            case NUMBER_IDENTIFIER:
                idNode = new NodeId(uaField.getNamespace(), UInteger.valueOf(uaField.getIdentifier()));
                break;
            case GUID_IDENTIFIER:
                idNode = new NodeId(uaField.getNamespace(), UUID.fromString(uaField.getIdentifier()));
                break;
            case BINARY_IDENTIFIER:
                idNode = new NodeId(uaField.getNamespace(), new ByteString(uaField.getIdentifier().getBytes()));
                break;

            default:
                idNode = new NodeId(uaField.getNamespace(), uaField.getIdentifier());
        }

        return idNode;
    }

    private String getEndpointUrl(InetAddress address, Integer port, String params) {
        return "opc.tcp://" + address.getHostAddress() + ":" + port + "/" + params;
    }

    private Predicate<EndpointDescription> endpointFilter() {
        return e -> true;
    }

    private SecurityPolicy getSecurityPolicy() {
        return SecurityPolicy.None;
    }

    private IdentityProvider getIdentityProvider() {
        return new AnonymousProvider();
    }

    private static String getSubPathOfParams(String params){
        if(params.contains("=")){
            if(params.contains("?")){
                return params.split("\\?")[0];
            }else{
                return "";
            }

        }else {
            return params;
        }
    }

    private static String getOptionString(String params){
        if(params.contains("=")){
            if(params.contains("?")){
                return params.split("\\?")[1];
            }else{
                return params;
            }

        }else {
            return "";
        }
    }
}
