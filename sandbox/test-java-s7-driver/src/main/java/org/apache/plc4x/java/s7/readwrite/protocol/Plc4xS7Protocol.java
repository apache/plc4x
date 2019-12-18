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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.*;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.events.IsoTPConnectedEvent;
import org.apache.plc4x.java.s7.readwrite.types.*;
import org.apache.plc4x.java.s7.readwrite.utils.S7Field;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Array;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

public class Plc4xS7Protocol extends Plc4xProtocolBase<TPKTPacket> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xS7Protocol.class);

    private final int callingTsapId;
    private int calledTsapId;
    private COTPTpduSize cotpTpduSize;
    private int pduSize;
    private int maxAmqCaller;
    private int maxAmqCallee;
    private S7ControllerType controllerType;

    private static final AtomicInteger tpduGenerator = new AtomicInteger(10);
    private final Map<Integer, PlcRequestContainer> requests;

    public Plc4xS7Protocol(int callingTsapId, int calledTsapId, COTPTpduSize tpduSize,
                           int maxAmqCaller, int maxAmqCallee, S7ControllerType controllerType) {
        // FIXME REMOVE
        super(TPKTPacket.class);
        this.callingTsapId = callingTsapId;
        this.calledTsapId = calledTsapId;
        this.cotpTpduSize = tpduSize;
        this.pduSize = tpduSize.getSizeInBytes() - 16;
        this.maxAmqCaller = maxAmqCaller;
        this.maxAmqCallee = maxAmqCallee;
        this.controllerType = controllerType;

        requests = new HashMap<>();
    }

    @Override
    public void onConnect() {
        logger.debug("ISO Transport Protocol Sending Connection Request");
        // Open the session on ISO Transport Protocol first.
        COTPPacketConnectionRequest connectionRequest = new COTPPacketConnectionRequest(
            new COTPParameter[] {
                new COTPParameterCalledTsap(calledTsapId),
                new COTPParameterCallingTsap(callingTsapId),
                new COTPParameterTpduSize(cotpTpduSize)
            }, null, (short) 0x0000, (short) 0x000F, COTPProtocolClass.CLASS_0);
        TPKTPacket packet = new TPKTPacket(connectionRequest);
        context.send(packet);
    }

    @Override
    protected void encode(PlcRequestContainer msg, Consumer<TPKTPacket> sendHandler) throws Exception {
        if(msg.getRequest() instanceof DefaultPlcReadRequest) {
            DefaultPlcReadRequest request = (DefaultPlcReadRequest) msg.getRequest();
            List<S7VarRequestParameterItem> requestItems = new ArrayList<>(request.getNumberOfFields());
            for (PlcField field : request.getFields()) {
                requestItems.add(new S7VarRequestParameterItemAddress(toS7Address(field)));
            }
            final int tpduId = tpduGenerator.getAndIncrement();
            sendHandler.accept(new TPKTPacket(new COTPPacketData(null,
                new S7MessageRequest(tpduId,
                    new S7ParameterReadVarRequest(requestItems.toArray(new S7VarRequestParameterItem[0])),
                    new S7PayloadReadVarRequest()),
                true, (short) tpduId)));
            requests.put(tpduId, msg);
        }
    }

    @Override
    protected void decode(TPKTPacket msg) throws Exception {
        if((msg == null) || (msg.getPayload() == null)) {
            return;
        }

        // When getting a response to the connection request on COTP layer, extract some
        // data and continue logging in on the S7 protocol.
        if(msg.getPayload() instanceof COTPPacketConnectionResponse) {
            COTPPacketConnectionResponse cotpPacketConnectionResponse = (COTPPacketConnectionResponse) msg.getPayload();
            for (COTPParameter parameter : cotpPacketConnectionResponse.getParameters()) {
                if(parameter instanceof COTPParameterCalledTsap) {
                    COTPParameterCalledTsap cotpParameterCalledTsap = (COTPParameterCalledTsap) parameter;
                    calledTsapId = cotpParameterCalledTsap.getTsapId();
                } else if(parameter instanceof COTPParameterTpduSize) {
                    COTPParameterTpduSize cotpParameterTpduSize = (COTPParameterTpduSize) parameter;
                    cotpTpduSize = cotpParameterTpduSize.getTpduSize();
                } else if(parameter instanceof COTPParameterCallingTsap) {
                    // Ignore this ...
                } else {
                    logger.warn("Got unknown parameter type '" + parameter.getClass().getName() + "'");
                }
            }

            // Send an S7 login message.
            S7ParameterSetupCommunication s7ParameterSetupCommunication =
                new S7ParameterSetupCommunication(maxAmqCaller, maxAmqCallee, pduSize);
            S7Message s7Message = new S7MessageRequest(0, s7ParameterSetupCommunication,
                new S7PayloadSetupCommunication());
            COTPPacketData cotpPacketData = new COTPPacketData(null, s7Message, true, (short) 1);
            TPKTPacket tpktPacket = new TPKTPacket(cotpPacketData);
            context.send(tpktPacket);
        }

        else if(msg.getPayload() instanceof COTPPacketData) {
            COTPPacketData packetData = (COTPPacketData) msg.getPayload();
            if(packetData.getPayload() instanceof S7MessageResponse) {
                S7MessageResponse s7MessageResponse = (S7MessageResponse) packetData.getPayload();
                final S7Parameter parameter = s7MessageResponse.getParameter();
                if(parameter instanceof S7ParameterSetupCommunication) {
                    S7ParameterSetupCommunication setupCommunication = (S7ParameterSetupCommunication) parameter;

                    // Save some data from the response.
                    maxAmqCaller = setupCommunication.getMaxAmqCaller();
                    maxAmqCallee = setupCommunication.getMaxAmqCallee();
                    pduSize = setupCommunication.getPduLength();

                    // Only if the controller type is set to "ANY", then try to identify the PLC type.
                    if(controllerType == S7ControllerType.ANY) {
                        // Prepare a message to request the remote to identify itself.
                        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(1, new S7ParameterUserData(new S7ParameterUserDataItem[] {
                            new S7ParameterUserDataItemCPUFunctions((short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00, null, null, null)
                        }), new S7PayloadUserData( new S7PayloadUserDataItem[] {
                            new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, new SzlId(SzlModuleTypeClass.CPU, (byte) 0x00, SzlSublist.MODULE_IDENTIFICATION), 0x0000)
                        }));
                        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (short) 2);
                        TPKTPacket tpktPacket = new TPKTPacket(cotpPacketData);
                        context.send(tpktPacket);
                    } else {
                        // Send an event that connection setup is complete.
                        context.fireConnected();

                    }
                } else if (parameter instanceof S7ParameterReadVarResponse) {
                    final PlcRequestContainer requestContainer = requests.remove(s7MessageResponse.getTpduReference());
                    final PlcResponse response = decodeReadResponse(s7MessageResponse, requestContainer);
                    requestContainer.getResponseFuture().complete(response);
                } else if (parameter instanceof S7ParameterWriteVarResponse) {
                    S7ParameterWriteVarResponse writeResponseParameter = (S7ParameterWriteVarResponse) parameter;

                    System.out.println(writeResponseParameter);
                }
            } else if(packetData.getPayload() instanceof S7MessageUserData) {
                S7MessageUserData messageUserData = (S7MessageUserData) packetData.getPayload();
                if(messageUserData.getPayload() instanceof S7PayloadUserData) {
                    S7PayloadUserData payloadUserData = (S7PayloadUserData) messageUserData.getPayload();
                    for (S7PayloadUserDataItem item : payloadUserData.getItems()) {
                        if(item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse) {
                            S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponseItem =
                                (S7PayloadUserDataItemCpuFunctionReadSzlResponse) item;
                            for (SzlDataTreeItem readSzlResponseItemItem : readSzlResponseItem.getItems()) {
                                if(readSzlResponseItemItem.getItemIndex() == 0x0001) {
                                    final String articleNumber = new String(readSzlResponseItemItem.getMlfb());
                                    controllerType = lookupControllerType(articleNumber);

                                    // Send an event that connection setup is complete.
                                    context.fireConnected();
                                }
                            }
                        }
                    }
                }
            } else {
                System.out.println(packetData);
            }
        }

        else {
            System.out.println(msg);
        }
    }

    private PlcResponse decodeReadResponse(S7MessageResponse responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        InternalPlcReadRequest plcReadRequest = (InternalPlcReadRequest) requestContainer.getRequest();

        S7PayloadReadVarResponse payload = (S7PayloadReadVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> values = new HashMap<>();
        S7VarPayloadDataItem[] payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
            S7VarPayloadDataItem payloadItem = payloadItems[index];

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            BaseDefaultFieldItem fieldItem = null;
            ByteBuf data = Unpooled.wrappedBuffer(payloadItem.getData());
            if (responseCode == PlcResponseCode.OK) {
                try {
                    switch (field.getDataType()) {
                        // -----------------------------------------
                        // Bit
                        // -----------------------------------------
                        case BOOL:
                            fieldItem = decodeReadResponseBitField(field, data);
                            break;
                        // -----------------------------------------
                        // Bit-strings
                        // -----------------------------------------
                        case BYTE:  // 1 byte
                            fieldItem = decodeReadResponseByteBitStringField(field, data);
                            break;
                        case WORD:  // 2 byte (16 bit)
                            fieldItem = decodeReadResponseShortBitStringField(field, data);
                            break;
                        case DWORD:  // 4 byte (32 bit)
                            fieldItem = decodeReadResponseIntegerBitStringField(field, data);
                            break;
                        case LWORD:  // 8 byte (64 bit)
                            fieldItem = decodeReadResponseLongBitStringField(field, data);
                            break;
                        // -----------------------------------------
                        // Integers
                        // -----------------------------------------
                        // 8 bit:
                        case SINT:
                            fieldItem = decodeReadResponseSignedByteField(field, data);
                            break;
                        case USINT:
                            fieldItem = decodeReadResponseUnsignedByteField(field, data);
                            break;
                        // 16 bit:
                        case INT:
                            fieldItem = decodeReadResponseSignedShortField(field, data);
                            break;
                        case UINT:
                            fieldItem = decodeReadResponseUnsignedShortField(field, data);
                            break;
                        // 32 bit:
                        case DINT:
                            fieldItem = decodeReadResponseSignedIntegerField(field, data);
                            break;
                        case UDINT:
                            fieldItem = decodeReadResponseUnsignedIntegerField(field, data);
                            break;
                        // 64 bit:
                        case LINT:
                            fieldItem = decodeReadResponseSignedLongField(field, data);
                            break;
                        case ULINT:
                            fieldItem = decodeReadResponseUnsignedLongField(field, data);
                            break;
                        // -----------------------------------------
                        // Floating point values
                        // -----------------------------------------
                        case REAL:
                            fieldItem = decodeReadResponseFloatField(field, data);
                            break;
                        case LREAL:
                            fieldItem = decodeReadResponseDoubleField(field, data);
                            break;
                        // -----------------------------------------
                        // Characters & Strings
                        // -----------------------------------------
                        case CHAR: // 1 byte (8 bit)
                            fieldItem = decodeReadResponseFixedLengthStringField(1, false, data);
                            break;
                        case WCHAR: // 2 byte
                            fieldItem = decodeReadResponseFixedLengthStringField(1, true, data);
                            break;
                        case STRING:
                            fieldItem = decodeReadResponseVarLengthStringField(false, data);
                            break;
                        case WSTRING:
                            fieldItem = decodeReadResponseVarLengthStringField(true, data);
                            break;
                        // -----------------------------------------
                        // TIA Date-Formats
                        // -----------------------------------------
                        case DATE_AND_TIME:
                            fieldItem = decodeReadResponseDateAndTime(field, data);
                            break;
                        case TIME_OF_DAY:
                            fieldItem = decodeReadResponseTimeOfDay(field, data);
                            break;
                        case DATE:
                            fieldItem = decodeReadResponseDate(field, data);
                            break;
                        default:
                            throw new PlcProtocolException("Unsupported type " + field.getDataType());
                    }
                }
                catch (Exception e){
                    logger.warn("Some other error occurred casting field {}, FieldInformation: {}",fieldName, field,e);
                }
            }
            Pair<PlcResponseCode, BaseDefaultFieldItem> result = new ImmutablePair<>(responseCode, fieldItem);
            values.put(fieldName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    private PlcResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return PlcResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return PlcResponseCode.OK;
            case NOT_FOUND:
                return PlcResponseCode.NOT_FOUND;
            case INVALID_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED:
                return PlcResponseCode.INVALID_DATATYPE;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    BaseDefaultFieldItem decodeReadResponseBitField(S7Field field, ByteBuf data) {
        Boolean[] booleans = readAllValues(Boolean.class, field, i -> data.readByte() != 0x00);
        return new DefaultBooleanFieldItem(booleans);
    }

    BaseDefaultFieldItem decodeReadResponseByteBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements()];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseShortBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 2];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseIntegerBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 4];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseLongBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 8];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeBitStringField(byte[] bytes) {
        BitSet bitSet = BitSet.valueOf(bytes);
        Boolean[] booleanValues = new Boolean[8 * bytes.length];
        int k = 0;
        for(int i = bytes.length - 1; i >= 0; i--) {
            for(int j = 0; j < 8; j++) {
                booleanValues[k++] = bitSet.get(8 * i + j);
            }
        }
        return new DefaultBooleanFieldItem(booleanValues);
    }

    BaseDefaultFieldItem decodeReadResponseSignedByteField(S7Field field, ByteBuf data) {
        Byte[] bytes = readAllValues(Byte.class, field, i -> data.readByte());
        return new DefaultByteFieldItem(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedByteField(S7Field field, ByteBuf data) {
        Short[] shorts = readAllValues(Short.class, field, i -> data.readUnsignedByte());
        return new DefaultShortFieldItem(shorts);
    }

    BaseDefaultFieldItem decodeReadResponseSignedShortField(S7Field field, ByteBuf data) {
        Short[] shorts = readAllValues(Short.class, field, i -> data.readShort());
        return new DefaultShortFieldItem(shorts);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedShortField(S7Field field, ByteBuf data) {
        Integer[] ints = readAllValues(Integer.class, field, i -> data.readUnsignedShort());
        return new DefaultIntegerFieldItem(ints);
    }

    BaseDefaultFieldItem decodeReadResponseSignedIntegerField(S7Field field, ByteBuf data) {
        Integer[] ints = readAllValues(Integer.class, field, i -> data.readInt());
        return new DefaultIntegerFieldItem(ints);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedIntegerField(S7Field field, ByteBuf data) {
        Long[] longs = readAllValues(Long.class, field, i -> data.readUnsignedInt());
        return new DefaultLongFieldItem(longs);
    }

    BaseDefaultFieldItem decodeReadResponseSignedLongField(S7Field field, ByteBuf data) {
        Long[] longs = readAllValues(Long.class, field, i -> data.readLong());
        return new DefaultLongFieldItem(longs);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedLongField(S7Field field, ByteBuf data) {
        BigInteger[] bigIntegers = readAllValues(BigInteger.class, field, i -> readUnsigned64BitInteger(data));
        return new DefaultBigIntegerFieldItem(bigIntegers);
    }

    BaseDefaultFieldItem decodeReadResponseFloatField(S7Field field, ByteBuf data) {
        Float[] floats = readAllValues(Float.class, field, i -> data.readFloat());
        return new DefaultFloatFieldItem(floats);
    }

    BaseDefaultFieldItem decodeReadResponseDoubleField(S7Field field, ByteBuf data) {
        Double[] doubles = readAllValues(Double.class, field, i -> data.readDouble());
        return new DefaultDoubleFieldItem(doubles);
    }

    BaseDefaultFieldItem decodeReadResponseFixedLengthStringField(int numChars, boolean isUtf16, ByteBuf data) {
        int numBytes = isUtf16 ? numChars * 2 : numChars;
        String stringValue = data.readCharSequence(numBytes, StandardCharsets.UTF_8).toString();
        return new DefaultStringFieldItem(stringValue);
    }

    BaseDefaultFieldItem decodeReadResponseVarLengthStringField(boolean isUtf16, ByteBuf data) {
        // Max length ... ignored.
        data.skipBytes(1);

        //reading out byte and transforming that to an unsigned byte within an integer, otherwise longer strings are failing
        byte currentLengthByte = data.readByte();
        int currentLength = currentLengthByte & 0xFF;
        return decodeReadResponseFixedLengthStringField(currentLength, isUtf16, data);
    }

    BaseDefaultFieldItem decodeReadResponseDateAndTime(S7Field field,ByteBuf data) {
        LocalDateTime[] localDateTimes = readAllValues(LocalDateTime.class,field, i -> readDateAndTime(data));
        return new DefaultLocalDateTimeFieldItem(localDateTimes);
    }

    BaseDefaultFieldItem decodeReadResponseTimeOfDay(S7Field field,ByteBuf data) {
        LocalTime[] localTimes = readAllValues(LocalTime.class,field, i -> readTimeOfDay(data));
        return new DefaultLocalTimeFieldItem(localTimes);
    }

    BaseDefaultFieldItem decodeReadResponseDate(S7Field field,ByteBuf data) {
        LocalDate[] localTimes = readAllValues(LocalDate.class,field, i -> readDate(data));
        return new DefaultLocalDateFieldItem(localTimes);
    }

    private static <T> T[] readAllValues(Class<T> clazz, S7Field field, Function<Integer, T> extract) {
        try {
            return IntStream.rangeClosed(1, field.getNumElements())
                .mapToObj(extract::apply)
                .collect(Collectors.toList())
                .toArray((T[]) Array.newInstance(clazz, 0));
        } catch (IndexOutOfBoundsException e) {
            throw new PlcRuntimeException("To few bytes in the buffer to read requested type", e);
        }
    }

    private static BigInteger readUnsigned64BitInteger(ByteBuf data) {
        byte[] bytes = new byte[9];
        // Set the first byte to 0
        bytes[0] = 0;
        // Read the next 8 bytes into the rest.
        data.readBytes(bytes, 1, 8);
        return new BigInteger(bytes);
    }

    LocalDateTime readDateAndTime(ByteBuf data) {
        //per definition for Date_And_Time only the first 6 bytes are used

        int year=convertByteToBcd(data.readByte());
        int month=convertByteToBcd(data.readByte());
        int day=convertByteToBcd(data.readByte());
        int hour=convertByteToBcd(data.readByte());
        int minute=convertByteToBcd(data.readByte());
        int second=convertByteToBcd(data.readByte());
        //skip the last 2 bytes no information present
        data.readByte();
        data.readByte();

        //data-type ranges from 1990 up to 2089
        if(year>=90){
            year+=1900;
        }
        else{
            year+=2000;
        }

        return LocalDateTime.of(year,month,day,hour,minute,second);
    }

    LocalTime readTimeOfDay(ByteBuf data) {
        //per definition for Date_And_Time only the first 6 bytes are used
        int millisSinsMidnight = data.readInt();
        return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(millisSinsMidnight, ChronoUnit.MILLIS);
    }

    LocalDate readDate(ByteBuf data) {
        //per definition for Date_And_Time only the first 6 bytes are used
        int daysSince1990 = data.readUnsignedShort();
        return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);
    }

    /**
     * converts incoming byte to an integer regarding used BCD format
     * @param incomingByte
     * @return converted BCD number
     */
    private static int convertByteToBcd(byte incomingByte) {
        int dec = (incomingByte >> 4) * 10;
        return dec + (incomingByte & 0x0f);
    }

    /**
     * Little helper method to parse Siemens article numbers and extract the type of controller.
     * @param articleNumber article number string.
     * @return type of controller.
     */
    private S7ControllerType lookupControllerType(String articleNumber) {
        if (!articleNumber.startsWith("6ES7 ")) {
            return S7ControllerType.ANY;
        }
        String model = articleNumber.substring(articleNumber.indexOf(' ') + 1, articleNumber.indexOf(' ') + 2);
        switch (model) {
            case "2":
                return S7ControllerType.S7_1200;
            case "5":
                return S7ControllerType.S7_1500;
            case "3":
                return S7ControllerType.S7_300;
            case "4":
                return S7ControllerType.S7_400;
            default:
                if (logger.isInfoEnabled()) {
                    logger.info(String.format("Looking up unknown article number %s", articleNumber));
                }
                return S7ControllerType.ANY;
        }
    }

    protected S7Address toS7Address(PlcField field) {
        if(!(field instanceof S7Field)) {
            throw new RuntimeException("Unsupported address type " + field.getClass().getName());
        }
        S7Field s7Field = (S7Field) field;
        return new S7AddressAny(s7Field.getDataType(), s7Field.getNumElements(), s7Field.getBlockNumber(),
            s7Field.getMemoryArea(), s7Field.getByteOffset(), s7Field.getBitOffset());
    }

}
