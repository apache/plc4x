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
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
package org.apache.plc4x.java.opcua.connection;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.*;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.apache.plc4x.java.opcua.protocol.OpcuaField;
import org.apache.plc4x.java.opcua.protocol.OpcuaSubsriptionHandle;
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
import org.eclipse.milo.opcua.stack.core.security.SecurityPolicy;
import org.eclipse.milo.opcua.stack.core.types.builtin.*;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.UInteger;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.UShort;
import org.eclipse.milo.opcua.stack.core.types.enumerated.MonitoringMode;
import org.eclipse.milo.opcua.stack.core.types.enumerated.TimestampsToReturn;
import org.eclipse.milo.opcua.stack.core.types.structured.EndpointDescription;
import org.eclipse.milo.opcua.stack.core.types.structured.MonitoredItemCreateRequest;
import org.eclipse.milo.opcua.stack.core.types.structured.MonitoringParameters;
import org.eclipse.milo.opcua.stack.core.types.structured.ReadValueId;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.net.InetAddress;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Predicate;
import java.util.stream.Collectors;

import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.uint;

public class OpcuaTcpPlcConnection extends BaseOpcuaPlcConnection {

    private static final int OPCUA_DEFAULT_TCP_PORT = 4840;

    private static final Logger logger = LoggerFactory.getLogger(OpcuaTcpPlcConnection.class);
    private InetAddress address;
    private int requestTimeout = 5000;
    private int port;
    private String params;
    private OpcUaClient client;
    private boolean isConnected = false;
    private  final AtomicLong clientHandles = new AtomicLong(1L);

    private OpcuaTcpPlcConnection(InetAddress address, String params, int requestTimeout) {
        this( address, OPCUA_DEFAULT_TCP_PORT,  params, requestTimeout);
        logger.info("Configured OpcuaTcpPlcConnection with: host-name {}", address.getHostAddress());
    }

    public OpcuaTcpPlcConnection(InetAddress address, int port, String params, int requestTimeout) {
        this(params);
        logger.info("Configured OpcuaTcpPlcConnection with: host-name {}", address.getHostAddress());
        this.address = address;
        this.port = port;
        this.params = params;
        this.requestTimeout = requestTimeout;
    }

    public OpcuaTcpPlcConnection(String params) {
        super(params);
    }

    public static OpcuaTcpPlcConnection of(InetAddress address, String params, int requestTimeout) {
        return new OpcuaTcpPlcConnection(address, params, requestTimeout);
    }

    public static OpcuaTcpPlcConnection of(InetAddress address, int port, String params, int requestTimeout) {
        return new OpcuaTcpPlcConnection(address, port, params, requestTimeout);
    }

    public static BaseDefaultFieldItem encodeFieldItem(DataValue value){
        NodeId typeNode = value.getValue().getDataType().get();
        Object objValue = value.getValue().getValue();

        if(typeNode.equals(Identifiers.Boolean)){
            return new DefaultBooleanFieldItem((Boolean)objValue);
        }else if (typeNode.equals(Identifiers.ByteString)){
            byte[] array = ((ByteString)objValue).bytes();
            Byte[] byteArry = new Byte[array.length];
            int counter = 0;
            for (byte bytie: array
            ) {
                byteArry[counter] = bytie;
                counter++;
            }
            return new DefaultByteArrayFieldItem(byteArry);
        }else if (typeNode.equals(Identifiers.Integer)){
            return new DefaultIntegerFieldItem((Integer)objValue);
        }else if (typeNode.equals(Identifiers.Int16)){
            return new DefaultShortFieldItem((Short)objValue);
        }else if (typeNode.equals(Identifiers.Int32)){
            return new DefaultIntegerFieldItem((Integer)objValue);
        }else if (typeNode.equals(Identifiers.Int64)){
            return new DefaultLongFieldItem((Long)objValue);
        }else if (typeNode.equals(Identifiers.UInteger)){
            return new DefaultLongFieldItem((Long)objValue);
        }else if (typeNode.equals(Identifiers.UInt16)){
            return new DefaultIntegerFieldItem(((UShort)objValue).intValue());
        }else if (typeNode.equals(Identifiers.UInt32)){
            return new DefaultLongFieldItem(((UInteger)objValue).longValue());
        }else if (typeNode.equals(Identifiers.UInt64)){
            return new DefaultBigIntegerFieldItem(new BigInteger(objValue.toString()));
        }else if (typeNode.equals(Identifiers.Byte)){
            return new DefaultShortFieldItem(Short.valueOf(objValue.toString()));
        }else if (typeNode.equals(Identifiers.Float)){
            return new DefaultFloatFieldItem((Float)objValue);
        }else if (typeNode.equals(Identifiers.Double)){
            return new DefaultDoubleFieldItem((Double)objValue);
        }else if (typeNode.equals(Identifiers.SByte)){
            return new DefaultByteFieldItem((Byte)objValue);
        }else {
            return new DefaultStringFieldItem(objValue.toString());
        }

    }

    public InetAddress getRemoteAddress() {
        return address;
    }

    @Override
    public void connect() throws PlcConnectionException {
        List<EndpointDescription> endpoints =  null;

        try {
            endpoints = DiscoveryClient.getEndpoints(getEndpointUrl(address, port, params)).get();
        //TODO Exception should be handeled better when the Discovery-API of Milo is stable
        } catch (Exception ex) {
            // try the explicit discovery endpoint as well
            String discoveryUrl = getEndpointUrl(address, port, params);

            if (!discoveryUrl.endsWith("/")) {
                discoveryUrl += "/";
            }
            discoveryUrl += "discovery";

            logger.info("Trying explicit discovery URL: {}", discoveryUrl);
            try {
                endpoints = DiscoveryClient.getEndpoints(discoveryUrl).get();
            } catch (InterruptedException | ExecutionException e) {
                throw new PlcConnectionException("Unable to discover URL:" + discoveryUrl);
            }
        }

        EndpointDescription endpoint = endpoints.stream()
                .filter(e -> e.getSecurityPolicyUri().equals(getSecurityPolicy().getUri()))
                .filter(endpointFilter())
                .findFirst()
                .orElseThrow(() -> new PlcConnectionException("No desired endpoints from"));

        OpcUaClientConfig config = OpcUaClientConfig.builder()
            .setApplicationName(LocalizedText.english("eclipse milo opc-ua client of the apache PLC4X:PLC4J project"))
            .setApplicationUri("urn:eclipse:milo:plc4x:client")
            .setEndpoint(endpoint)
            .setIdentityProvider(getIdentityProvider())
            .setRequestTimeout(UInteger.valueOf(requestTimeout))
            .build();

        try {
            this.client =  OpcUaClient.create(config);
            this.client.connect().get();
            isConnected = true;
        } catch (UaException e) {
            isConnected = false;
            String message = (config == null) ? "NULL" : config.toString();
            throw  new PlcConnectionException("The given input values are a not valid OPC UA connection configuration [CONFIG]: " + message);
        } catch (InterruptedException | ExecutionException e) {
            isConnected = false;
            throw  new PlcConnectionException("Error while creation of the connection because of : " + e.getMessage());
        }
    }

    @Override
    public boolean isConnected() {
        return client != null && isConnected;
    }

    @Override
    public void close() throws Exception {
        if(client != null){
            client.disconnect().get();
            isConnected = false;
        }
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        InternalPlcSubscriptionRequest internalPlcSubscriptionRequest = checkInternal(subscriptionRequest, InternalPlcSubscriptionRequest.class);
        CompletableFuture<PlcSubscriptionResponse> future = CompletableFuture.supplyAsync(() ->{
            Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> responseItems = internalPlcSubscriptionRequest.getSubscriptionPlcFieldMap().entrySet().stream()
            .map(subscriptionPlcFieldEntry -> {
                final String plcFieldName = subscriptionPlcFieldEntry.getKey();
                final SubscriptionPlcField subscriptionPlcField = subscriptionPlcFieldEntry.getValue();
                final OpcuaField field = (OpcuaField)Objects.requireNonNull(subscriptionPlcField.getPlcField());
                long cycleTime = subscriptionPlcField.getDuration().orElse(Duration.ofSeconds(1)).toMillis();
                NodeId idNode = generateNodeId(field);
                ReadValueId readValueId = new ReadValueId(
                    idNode,
                    AttributeId.Value.uid(), null, QualifiedName.NULL_VALUE);
                UInteger clientHandle = uint(clientHandles.getAndIncrement());

                MonitoringParameters parameters = new MonitoringParameters(
                    clientHandle,
                    (double) cycleTime,     // sampling interval
                    null,       // filter, null means use default
                    uint(1),   // queue size
                    true        // discard oldest
                );
                MonitoringMode monitoringMode;
                switch (subscriptionPlcField.getPlcSubscriptionType()) {
                    case CYCLIC:
                        monitoringMode = MonitoringMode.Sampling;
                        break;
                    case CHANGE_OF_STATE:
                        monitoringMode = MonitoringMode.Reporting;
                        break;
                    case EVENT:
                        monitoringMode = MonitoringMode.Reporting;
                        break;
                    default: monitoringMode = MonitoringMode.Reporting;
                }

                PlcSubscriptionHandle subHandle = null;
                PlcResponseCode responseCode = PlcResponseCode.ACCESS_DENIED;
                try {
                    UaSubscription subscription = client.getSubscriptionManager().createSubscription(1000.0).get();

                    MonitoredItemCreateRequest request = new MonitoredItemCreateRequest(
                        readValueId, monitoringMode, parameters);
                    List<MonitoredItemCreateRequest> requestList = new LinkedList<>();
                    requestList.add(request);
                    OpcuaSubsriptionHandle subsriptionHandle = new OpcuaSubsriptionHandle(plcFieldName, clientHandle);
                    BiConsumer<UaMonitoredItem, Integer> onItemCreated =
                        (item, id) -> item.setValueConsumer(subsriptionHandle::onSubscriptionValue);

                    List<UaMonitoredItem> items = subscription.createMonitoredItems(
                        TimestampsToReturn.Both,
                        requestList,
                        onItemCreated
                    ).get();

                    subHandle = subsriptionHandle;
                    responseCode = PlcResponseCode.OK;
                } catch (InterruptedException | ExecutionException e) {
                    logger.warn("Unable to subscribe Elements because of: {}", e.getMessage());
                }


                return Pair.of(plcFieldName, Pair.of(responseCode, subHandle));
            })
            .collect(Collectors.toMap(Pair::getKey, Pair::getValue));
            return (PlcSubscriptionResponse) new DefaultPlcSubscriptionResponse(internalPlcSubscriptionRequest, responseItems);
        });

        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        InternalPlcUnsubscriptionRequest internalPlcUnsubscriptionRequest = checkInternal(unsubscriptionRequest, InternalPlcUnsubscriptionRequest.class);
        internalPlcUnsubscriptionRequest.getInternalPlcSubscriptionHandles().forEach(o -> {
            OpcuaSubsriptionHandle opcSubHandle = (OpcuaSubsriptionHandle) o;
            try {
                client.getSubscriptionManager().deleteSubscription(opcSubHandle.getClientHandle()).get();
            } catch (InterruptedException | ExecutionException e) {
                logger.warn("Unable to unsubscribe Elements because of: {}", e.getMessage());
            }
        });

        return null;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        List<PlcConsumerRegistration> unregisters = new LinkedList<>();
        handles.forEach(plcSubscriptionHandle -> unregisters.add(plcSubscriptionHandle.register(consumer)));

        return () -> unregisters.forEach(PlcConsumerRegistration::unregister);
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        registration.unregister();
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = CompletableFuture.supplyAsync(() -> {
            readRequest.getFields();
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();
            List<NodeId> readValueIds = new LinkedList<>();
            List<PlcField> readPLCValues = readRequest.getFields();
            for (PlcField field: readPLCValues) {
                NodeId idNode = generateNodeId((OpcuaField) field);
                readValueIds.add(idNode);
            }

            CompletableFuture<List<DataValue>> dataValueCompletableFuture = client.readValues(0.0, TimestampsToReturn.Both, readValueIds);
            List<DataValue> readValues = null;
            try {
                readValues = dataValueCompletableFuture.get();
            } catch (InterruptedException | ExecutionException e) {
                logger.warn("Unable to read Elements because of: {}", e.getMessage());
            }
            for(int counter = 0; counter < readValueIds.size(); counter++){
                PlcResponseCode resultCode = PlcResponseCode.OK;
                BaseDefaultFieldItem stringItem = null;
                if(readValues == null || readValues.size() <= counter || readValues.get(counter).getStatusCode() != StatusCode.GOOD){
                    resultCode = PlcResponseCode.NOT_FOUND;
                }else{
                    stringItem = encodeFieldItem(readValues.get(counter));

                }
                Pair<PlcResponseCode, BaseDefaultFieldItem> newPair = new ImmutablePair<>(resultCode, stringItem);
                fields.put((String) readRequest.getFieldNames().toArray()[counter], newPair);


            }
            InternalPlcReadRequest internalPlcReadRequest = checkInternal(readRequest, InternalPlcReadRequest.class);
            return (PlcReadResponse) new DefaultPlcReadResponse(internalPlcReadRequest, fields );
        });


        return future;
    }


    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future;
        future = CompletableFuture.supplyAsync(() -> {

            InternalPlcWriteRequest internalPlcWriteRequest = (InternalPlcWriteRequest) writeRequest;

            List<PlcField> writePLCValues = writeRequest.getFields();
            LinkedList<DataValue> values = new LinkedList<>();
            LinkedList<NodeId> ids = new LinkedList<>();
            LinkedList<String> names = new LinkedList<>();
            Map<String, PlcResponseCode> fieldResponse = new HashMap<>();
            for (String fieldName: writeRequest.getFieldNames()) {
                OpcuaField uaField = (OpcuaField) writeRequest.getField(fieldName);
                NodeId idNode = generateNodeId(uaField);
                Variant var = new Variant(internalPlcWriteRequest.getFieldItem(fieldName).getObject(0));
                DataValue value = new DataValue(var, null, null);
                ids.add(idNode);
                names.add(fieldName);
                values.add(value);
            }
            CompletableFuture<List<StatusCode>> opcRequest =
                client.writeValues(ids, values);
            List<StatusCode> statusCodes = null;
            try {
                statusCodes = opcRequest.get();
            } catch (InterruptedException | ExecutionException e) {
                statusCodes = new LinkedList<>();
                for(int counter = 0; counter < ids.size(); counter++){
                    ((LinkedList<StatusCode>) statusCodes).push(StatusCode.BAD);
                }
            }

            for(int counter = 0; counter < names.size(); counter++){
                PlcResponseCode resultCode;
                if(statusCodes != null && statusCodes.size() > counter){
                    if(statusCodes.get(counter).isGood()){
                        resultCode = PlcResponseCode.OK;
                    }else if(statusCodes.get(counter).isUncertain()){
                        resultCode = PlcResponseCode.NOT_FOUND;
                    }else {
                        resultCode = PlcResponseCode.ACCESS_DENIED;
                    }
                }else{
                    resultCode = PlcResponseCode.ACCESS_DENIED;
                }
                fieldResponse.put(names.get(counter), resultCode);
            }
            InternalPlcWriteRequest internalPlcReadRequest = checkInternal(writeRequest, InternalPlcWriteRequest.class);
            PlcWriteResponse response = new DefaultPlcWriteResponse(internalPlcReadRequest, fieldResponse);
            return response;
        });


        return future;
    }


    private NodeId generateNodeId(OpcuaField uaField){
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

            default: idNode = new NodeId(uaField.getNamespace(), uaField.getIdentifier());
        }

        return  idNode;
    }

    private String getEndpointUrl(InetAddress address, Integer port, String params) {
        return "opc.tcp://" + address.getHostAddress() +":" + port + "/" + params;
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
}
