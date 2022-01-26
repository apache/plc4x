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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.context.SecureChannel;
import org.apache.plc4x.java.opcua.field.OpcuaField;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.math.BigInteger;
import java.time.Duration;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class OpcuaProtocolLogic extends Plc4xProtocolBase<OpcuaAPU> implements HasConfiguration<OpcuaConfiguration>, PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocolLogic.class);
    protected static final PascalString NULL_STRING = new PascalString( "");
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
    private OpcuaConfiguration configuration;
    private Map<Long, OpcuaSubscriptionHandle> subscriptions = new HashMap<>();
    private SecureChannel channel;
    private AtomicBoolean securedConnection = new AtomicBoolean(false);

    @Override
    public void setConfiguration(OpcuaConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void close(ConversationContext<OpcuaAPU> context) {
        //Nothing
    }

    @Override
    public void onDisconnect(ConversationContext<OpcuaAPU> context) {
        for (Map.Entry<Long, OpcuaSubscriptionHandle> subscriber : subscriptions.entrySet()) {
            subscriber.getValue().stopSubscriber();
        }
        channel.onDisconnect(context);
    }

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.channel = new SecureChannel(driverContext, this.configuration);
    }

    @Override
    public void onConnect(ConversationContext<OpcuaAPU> context) {
        LOGGER.debug("Opcua Driver running in ACTIVE mode.");

        if (this.channel == null) {
            this.channel = new SecureChannel(driverContext, this.configuration);
        }
        this.channel.onConnect(context);
    }

    @Override
    public void onDiscover(ConversationContext<OpcuaAPU> context) {
        // Only the TCP transport supports login.
        LOGGER.debug("Opcua Driver running in ACTIVE mode, discovering endpoints");
        channel.onDiscover(context);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        LOGGER.trace("Reading Value");

        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        RequestHeader requestHeader = new RequestHeader(channel.getAuthenticationToken(),
            SecureChannel.getCurrentDateTime(),
            channel.getRequestHandle(),
            0L,
            NULL_STRING,
            SecureChannel.REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        List<ExtensionObjectDefinition> readValueArray = new ArrayList<>(request.getFieldNames().size());
        Iterator<String> iterator = request.getFieldNames().iterator();
        for (int i = 0; i < request.getFieldNames().size(); i++ ) {
            String fieldName = iterator.next();
            OpcuaField field = (OpcuaField) request.getField(fieldName);

            NodeId nodeId = generateNodeId(field);

            readValueArray.add(new ReadValueId(nodeId,
                0xD,
                NULL_STRING,
                new QualifiedName(0, NULL_STRING)));
        }

        ReadRequest opcuaReadRequest = new ReadRequest(
            requestHeader,
            0.0d,
            TimestampsToReturn.timestampsToReturnNeither,
            readValueArray.size(),
            readValueArray);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.parseInt(opcuaReadRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            opcuaReadRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            extObject.serialize(buffer);

            /* Functional Consumer example using inner class */
            Consumer<byte []> consumer = opcuaResponse -> {
                PlcReadResponse response = null;
                try {
                    ExtensionObjectDefinition reply = ExtensionObject.staticParse(new ReadBufferByteBased(opcuaResponse, ByteOrder.LITTLE_ENDIAN), false).getBody();
                    if (reply instanceof ReadResponse) {
                        future.complete(new DefaultPlcReadResponse(request, readResponse(request.getFieldNames(), ((ReadResponse) reply).getResults())));
                    } else {
                        if (reply instanceof ServiceFault) {
                            ExtensionObjectDefinition header = ((ServiceFault) reply).getResponseHeader();
                            LOGGER.error("Read request ended up with ServiceFault: {}", header);
                        } else {
                            LOGGER.error("Remote party returned an error '{}'", reply);
                        }

                        Map<String, ResponseItem<PlcValue>> status = new LinkedHashMap<>();
                        for (String key : request.getFieldNames()) {
                            status.put(key, new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                        }
                        future.complete(new DefaultPlcReadResponse(request, status));
                        return;
                    }
                } catch (ParseException e) {
                    future.completeExceptionally(new PlcRuntimeException(e));
                };
            };

            /* Functional Consumer example using inner class */
            // Pass the response back to the application.
            Consumer<TimeoutException> timeout = future::completeExceptionally;

            /* Functional Consumer example using inner class */
            BiConsumer<OpcuaAPU, Throwable> error = (message, t) -> {

                // Pass the response back to the application.
                future.completeExceptionally(t);
            };

            channel.submit(context, timeout, error, consumer, buffer);

        } catch (SerializationException e) {
            LOGGER.error("Unable to serialise the ReadRequest");
        }

        return future;
    }

    private NodeId generateNodeId(OpcuaField field) {
        NodeId nodeId = null;
        if (field.getIdentifierType() == OpcuaIdentifierType.BINARY_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdTwoByte(Short.parseShort(field.getIdentifier())));
        } else if (field.getIdentifierType() == OpcuaIdentifierType.NUMBER_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdNumeric((short) field.getNamespace(), Long.parseLong(field.getIdentifier())));
        } else if (field.getIdentifierType() == OpcuaIdentifierType.GUID_IDENTIFIER) {
            UUID guid = UUID.fromString(field.getIdentifier());
            byte[] guidBytes = new byte[16];
            System.arraycopy(guid.getMostSignificantBits(), 0, guidBytes, 0, 8);
            System.arraycopy(guid.getLeastSignificantBits(), 0, guidBytes, 8, 8);
            nodeId = new NodeId(new NodeIdGuid((short) field.getNamespace(), guidBytes));
        } else if (field.getIdentifierType() == OpcuaIdentifierType.STRING_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdString((short) field.getNamespace(), new PascalString(field.getIdentifier())));
        }
        return nodeId;
    }

    public Map<String, ResponseItem<PlcValue>> readResponse(LinkedHashSet<String> fieldNames, List<DataValue> results) {
        PlcResponseCode responseCode = PlcResponseCode.OK;
        Map<String, ResponseItem<PlcValue>> response = new HashMap<>();
        int count = 0;
        for ( String field : fieldNames ) {
            PlcValue value = null;
            if (results.get(count).getValueSpecified()) {
                Variant variant = results.get(count).getValue();
                LOGGER.trace("Response of type {}", variant.getClass().toString());
                if (variant instanceof VariantBoolean) {
                    byte[] array = ((VariantBoolean) variant).getValue();
                    int length = array.length;
                    Boolean[] tmpValue = new Boolean[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = array[i] != 0;
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantSByte) {
                    byte[] array = ((VariantSByte) variant).getValue();
                    value = IEC61131ValueHandler.of(array);
                } else if (variant instanceof VariantByte) {
                    List<Short> array = ((VariantByte) variant).getValue();
                    Short[] tmpValue = array.toArray(new Short[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantInt16) {
                    List<Short> array = ((VariantInt16) variant).getValue();
                    Short[] tmpValue = array.toArray(new Short[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantUInt16) {
                    List<Integer> array = ((VariantUInt16) variant).getValue();
                    Integer[] tmpValue = array.toArray(new Integer[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantInt32) {
                    List<Integer> array = ((VariantInt32) variant).getValue();
                    Integer[] tmpValue = array.toArray(new Integer[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantUInt32) {
                    List<Long> array = ((VariantUInt32) variant).getValue();
                    Long[] tmpValue = array.toArray(new Long[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantInt64) {
                    List<Long> array = ((VariantInt64) variant).getValue();
                    Long[] tmpValue = array.toArray(new Long[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantUInt64) {
                    value = IEC61131ValueHandler.of(((VariantUInt64) variant).getValue());
                } else if (variant instanceof VariantFloat) {
                    List<Float> array = ((VariantFloat) variant).getValue();
                    Float[] tmpValue = array.toArray(new Float[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantDouble) {
                    List<Double> array = ((VariantDouble) variant).getValue();
                    Double[] tmpValue = array.toArray(new Double[0]);
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantString) {
                    int length = ((VariantString) variant).getValue().size();
                    List<PascalString> stringArray = ((VariantString) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = stringArray.get(i).getStringValue();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantDateTime) {
                    List<Long> array = ((VariantDateTime) variant).getValue();
                    int length = array.size();
                    LocalDateTime[] tmpValue = new LocalDateTime[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = LocalDateTime.ofInstant(Instant.ofEpochMilli(getDateTime(array.get(i))), ZoneOffset.UTC);
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantGuid) {
                    List<GuidValue> array = ((VariantGuid) variant).getValue();
                    int length = array.size();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        //These two data section aren't little endian like the rest.
                        byte[] data4Bytes = array.get(i).getData4();
                        int data4 = 0;
                        for (byte data4Byte : data4Bytes) {
                            data4 = (data4 << 8) + (data4Byte & 0xff);
                        }
                        byte[] data5Bytes = array.get(i).getData5();
                        long data5 = 0;
                        for (byte data5Byte : data5Bytes) {
                            data5 = (data5 << 8) + (data5Byte & 0xff);
                        }
                        tmpValue[i] = Long.toHexString(array.get(i).getData1()) + "-" + Integer.toHexString(array.get(i).getData2()) + "-" + Integer.toHexString(array.get(i).getData3()) + "-" + Integer.toHexString(data4) + "-" + Long.toHexString(data5);
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantXmlElement) {
                    int length = ((VariantXmlElement) variant).getValue().size();
                    List<PascalString> strings = ((VariantXmlElement) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = strings.get(i).getStringValue();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantLocalizedText) {
                    int length = ((VariantLocalizedText) variant).getValue().size();
                    List<LocalizedText> strings = ((VariantLocalizedText) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = "";
                        tmpValue[i] += strings.get(i).getLocaleSpecified() ? strings.get(i).getLocale().getStringValue() + "|" : "";
                        tmpValue[i] += strings.get(i).getTextSpecified() ? strings.get(i).getText().getStringValue() : "";
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantQualifiedName) {
                    int length = ((VariantQualifiedName) variant).getValue().size();
                    List<QualifiedName> strings = ((VariantQualifiedName) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = "ns=" + strings.get(i).getNamespaceIndex() + ";s=" + strings.get(i).getName().getStringValue();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantExtensionObject) {
                    int length = ((VariantExtensionObject) variant).getValue().size();
                    List<ExtensionObject> strings = ((VariantExtensionObject) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = strings.get(i).toString();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantNodeId) {
                    int length = ((VariantNodeId) variant).getValue().size();
                    List<NodeId> strings = ((VariantNodeId) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = strings.get(i).toString();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                }else if (variant instanceof VariantStatusCode) {
                    int length = ((VariantStatusCode) variant).getValue().size();
                    List<StatusCode> strings = ((VariantStatusCode) variant).getValue();
                    String[] tmpValue = new String[length];
                    for (int i = 0; i < length; i++) {
                        tmpValue[i] = strings.get(i).toString();
                    }
                    value = IEC61131ValueHandler.of(tmpValue);
                } else if (variant instanceof VariantByteString) {
                    PlcList plcList = new PlcList();
                    List<ByteStringArray> array = ((VariantByteString) variant).getValue();
                    for (int k = 0; k < array.size(); k++) {
                        int length = array.get(k).getValue().size();
                        Short[] tmpValue = new Short[length];
                        for (int i = 0; i < length; i++) {
                            tmpValue[i] = array.get(k).getValue().get(i);
                        }
                        plcList.add(IEC61131ValueHandler.of(tmpValue));
                    }
                    value = plcList;
                } else {
                    responseCode = PlcResponseCode.UNSUPPORTED;
                    LOGGER.error("Data type - " +  variant.getClass() + " is not supported ");
                }
            } else {
                if (results.get(count).getStatusCode().getStatusCode() == OpcuaStatusCode.BadNodeIdUnknown.getValue()) {
                    responseCode = PlcResponseCode.NOT_FOUND;
                } else {
                    responseCode = PlcResponseCode.UNSUPPORTED;
                }
                LOGGER.error("Error while reading value from OPC UA server error code:- " + results.get(count).getStatusCode().toString());
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
        int length = valueObject.getLength();
        switch (dataType) {
            case "BOOL":
            case "BIT":
                byte[] tmpBOOL = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpBOOL[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantBoolean(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpBOOL);
            case "BYTE":
            case "BITARR8":
            case "USINT":
            case "UINT8":
            case "BIT8":
                List<Short> tmpBYTE = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpBYTE.add((short) valueObject.getIndex(i).getByte());
                }
                return new VariantByte(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpBYTE);
            case "SINT":
            case "INT8":
                byte[] tmpSINT = new byte[length];
                for (int i = 0; i < length; i++) {
                    tmpSINT[i] = valueObject.getIndex(i).getByte();
                }
                return new VariantSByte(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpSINT);
            case "INT":
            case "INT16":
                List<Short> tmpINT16 = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpINT16.add(valueObject.getIndex(i).getShort());
                }
                return new VariantInt16(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpINT16);
            case "UINT":
            case "UINT16":
            case "WORD":
            case "BITARR16":
                List<Integer> tmpUINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUINT.add(valueObject.getIndex(i).getInt());
                }
                return new VariantUInt16(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpUINT);
            case "DINT":
            case "INT32":
                List<Integer> tmpDINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDINT.add(valueObject.getIndex(i).getInt());
                }
                return new VariantInt32(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpDINT);
            case "UDINT":
            case "UINT32":
            case "DWORD":
            case "BITARR32":
                List<Long> tmpUDINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpUDINT.add(valueObject.getIndex(i).getLong());
                }
                return new VariantUInt32(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpUDINT);
            case "LINT":
            case "INT64":
                List<Long> tmpLINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLINT.add(valueObject.getIndex(i).getLong());
                }
                return new VariantInt64(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpLINT);
            case "ULINT":
            case "UINT64":
            case "LWORD":
            case "BITARR64":
                List<BigInteger> tmpULINT = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpULINT.add(valueObject.getIndex(i).getBigInteger());
                }
                return new VariantUInt64(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpULINT);
            case "REAL":
            case "FLOAT":
                List<Float> tmpREAL = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpREAL.add(valueObject.getIndex(i).getFloat());
                }
                return new VariantFloat(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpREAL);
            case "LREAL":
            case "DOUBLE":
                List<Double> tmpLREAL = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpLREAL.add(valueObject.getIndex(i).getDouble());
                }
                return new VariantDouble(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpLREAL);
            case "CHAR":
            case "WCHAR":
            case "STRING":
            case "WSTRING":
            case "STRING16":
                List<PascalString> tmpString = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    String s = valueObject.getIndex(i).getString();
                    tmpString.add(new PascalString(s));
                }
                return new VariantString(length != 1,
                    false,
                    null,
                    null,
                    length == 1 ? null : length,
                    tmpString);
            case "DATE_AND_TIME":
                List<Long> tmpDateTime = new ArrayList<>(length);
                for (int i = 0; i < length; i++) {
                    tmpDateTime.add(valueObject.getIndex(i).getDateTime().toEpochSecond(ZoneOffset.UTC));
                }
                return new VariantDateTime(length != 1,
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
        LOGGER.trace("Writing Value");
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        RequestHeader requestHeader = new RequestHeader(channel.getAuthenticationToken(),
            SecureChannel.getCurrentDateTime(),
            channel.getRequestHandle(),
            0L,
            NULL_STRING,
            SecureChannel.REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        List<ExtensionObjectDefinition> writeValueList = new ArrayList<>(request.getFieldNames().size());
        for (String fieldName : request.getFieldNames()) {
            OpcuaField field = (OpcuaField) request.getField(fieldName);

            NodeId nodeId = generateNodeId(field);

            writeValueList.add(new WriteValue(nodeId,
                0xD,
                NULL_STRING,
                new DataValue(
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
                    null)));
        }

        WriteRequest opcuaWriteRequest = new WriteRequest(
            requestHeader,
            writeValueList.size(),
            writeValueList);

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.parseInt(opcuaWriteRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            opcuaWriteRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            extObject.serialize(buffer);

            /* Functional Consumer example using inner class */
            Consumer<byte[]> consumer = opcuaResponse -> {
                WriteResponse responseMessage = null;
                try {
                    responseMessage = (WriteResponse) ExtensionObject.staticParse(new ReadBufferByteBased(opcuaResponse, ByteOrder.LITTLE_ENDIAN), false).getBody();
                } catch (ParseException e) {
                    e.printStackTrace();
                }
                PlcWriteResponse response = writeResponse(request, responseMessage);

                // Pass the response back to the application.
                future.complete(response);
            };

            /* Functional Consumer example using inner class */
            // Pass the response back to the application.
            Consumer<TimeoutException> timeout = future::completeExceptionally;

            /* Functional Consumer example using inner class */
            BiConsumer<OpcuaAPU, Throwable> error = (message, t) -> {
                // Pass the response back to the application.
                future.completeExceptionally(t);
            };

            channel.submit(context, timeout, error, consumer, buffer);

        } catch (SerializationException e) {
            LOGGER.error("Unable to serialise the ReadRequest");
        }

        return future;
    }

    private PlcWriteResponse writeResponse(DefaultPlcWriteRequest request, WriteResponse writeResponse) {
        Map<String, PlcResponseCode> responseMap = new HashMap<>();
        List<StatusCode> results = writeResponse.getResults();
        Iterator<String> responseIterator = request.getFieldNames().iterator();
        for (int i = 0; i < request.getFieldNames().size(); i++ ) {
            String fieldName = responseIterator.next();
            OpcuaStatusCode statusCode = OpcuaStatusCode.enumForValue(results.get(i).getStatusCode());
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


    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
            long subscriptionId;
            ArrayList<String> fields = new ArrayList<>( subscriptionRequest.getFieldNames() );
            long cycleTime = (subscriptionRequest.getField(fields.get(0))).getDuration().orElse(Duration.ofMillis(1000)).toMillis();

            try {
                CompletableFuture<CreateSubscriptionResponse> subscription = onSubscribeCreateSubscription(cycleTime);
                CreateSubscriptionResponse response = subscription.get(SecureChannel.REQUEST_TIMEOUT_LONG, TimeUnit.MILLISECONDS);
                subscriptionId = response.getSubscriptionId();
                subscriptions.put(subscriptionId, new OpcuaSubscriptionHandle(context, this, channel, subscriptionRequest, subscriptionId, cycleTime));
            } catch (Exception e) {
                throw new PlcRuntimeException("Unable to subscribe because of: " + e.getMessage());
            }

            for (String fieldName : subscriptionRequest.getFieldNames()) {
                final DefaultPlcSubscriptionField fieldDefaultPlcSubscription = (DefaultPlcSubscriptionField) subscriptionRequest.getField(fieldName);
                if (!(fieldDefaultPlcSubscription.getPlcField() instanceof OpcuaField)) {
                    values.put(fieldName, new ResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                } else {
                    values.put(fieldName, new ResponseItem<>(PlcResponseCode.OK, subscriptions.get(subscriptionId)));
                }
            }
            return new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
        });
    }

    private CompletableFuture<CreateSubscriptionResponse> onSubscribeCreateSubscription(long cycleTime) {
        CompletableFuture<CreateSubscriptionResponse> future = new CompletableFuture<>();
        LOGGER.trace("Entering creating subscription request");

        RequestHeader requestHeader = new RequestHeader(channel.getAuthenticationToken(),
            SecureChannel.getCurrentDateTime(),
            channel.getRequestHandle(),
            0L,
            NULL_STRING,
            SecureChannel.REQUEST_TIMEOUT_LONG,
            NULL_EXTENSION_OBJECT);

        CreateSubscriptionRequest createSubscriptionRequest = new CreateSubscriptionRequest(
            requestHeader,
            cycleTime,
            12000,
            5,
            65536,
            true,
            (short) 0
        );

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.parseInt(createSubscriptionRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            createSubscriptionRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            extObject.serialize(buffer);

            /* Functional Consumer example using inner class */
            Consumer<byte[]> consumer = opcuaResponse -> {
                CreateSubscriptionResponse responseMessage = null;
                try {
                    responseMessage = (CreateSubscriptionResponse) ExtensionObject.staticParse(new ReadBufferByteBased(opcuaResponse, ByteOrder.LITTLE_ENDIAN), false).getBody();
                } catch (ParseException e) {
                    e.printStackTrace();
                }

                // Pass the response back to the application.
                future.complete(responseMessage);

            };

            /* Functional Consumer example using inner class */
            Consumer<TimeoutException> timeout = e -> {
                LOGGER.error("Timeout while waiting on the crate subscription response");
                e.printStackTrace();
                // Pass the response back to the application.
                future.completeExceptionally(e);
            };

            /* Functional Consumer example using inner class */
            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                LOGGER.error("Error while creating the subscription");
                e.printStackTrace();
                // Pass the response back to the application.
                future.completeExceptionally(e);
            };

            channel.submit(context, timeout, error, consumer, buffer);
        } catch (SerializationException e) {
            LOGGER.error("Error while creating the subscription");
            e.printStackTrace();
            future.completeExceptionally(e);
        }
        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        unsubscriptionRequest.getSubscriptionHandles().forEach(o -> {
            OpcuaSubscriptionHandle opcuaSubHandle = (OpcuaSubscriptionHandle) o;
            opcuaSubHandle.stopSubscriber();
        });
        return null;
    }

    public void removeSubscription(Long subscriptionId) {
        subscriptions.remove(subscriptionId);
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        List<PlcConsumerRegistration> registrations = new LinkedList<>();
        // Register the current consumer for each of the given subscription handles
        for (PlcSubscriptionHandle subscriptionHandle : handles) {
            LOGGER.debug("Registering Consumer");
            final PlcConsumerRegistration consumerRegistration = subscriptionHandle.register(consumer);
            registrations.add(consumerRegistration);
        }
        return new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        registration.unregister();
    }

    public static long getDateTime(long dateTime) {
        return (dateTime - EPOCH_OFFSET) / 10000;
    }

    private GuidValue toGuidValue(String identifier) {
        LOGGER.error("Querying Guid nodes is not supported");
        byte[] data4 = new byte[] {0,0};
        byte[] data5 = new byte[] {0,0,0,0,0,0};
        return new GuidValue(0L,0,0,data4, data5);
    }
}
