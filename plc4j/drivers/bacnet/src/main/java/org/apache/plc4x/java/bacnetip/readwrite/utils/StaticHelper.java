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
package org.apache.plc4x.java.bacnetip.readwrite.utils;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithAdditionalStringRepresentation;

public class StaticHelper {

    public static final Logger LOGGER = LoggerFactory.getLogger(StaticHelper.class);

    public static BACnetPropertyIdentifier readPropertyIdentifier(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        long readUnsignedLong = readBuffer.readUnsignedLong("propertyIdentifier", bitsToRead);
        if (!BACnetPropertyIdentifier.isDefined(readUnsignedLong)) {
            return BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetPropertyIdentifier.enumForValue(readUnsignedLong);
    }

    public static void writePropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier value) throws SerializationException {
        if (value == null || value == BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("propertyIdentifier", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static void writeProprietaryPropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier baCnetPropertyIdentifier, long value) throws SerializationException {
        if (baCnetPropertyIdentifier != null && baCnetPropertyIdentifier != BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryPropertyIdentifier", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static Long readProprietaryPropertyIdentifier(ReadBuffer readBuffer, BACnetPropertyIdentifier value, Long actualLength) throws ParseException {
        if (value != null && value != BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryPropertyIdentifier", bitsToRead);
    }

    public static BACnetEventType readEventType(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        int readUnsignedLong = readBuffer.readUnsignedInt("eventType", bitsToRead);
        if (!BACnetEventType.isDefined(readUnsignedLong)) {
            return BACnetEventType.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetEventType.enumForValue(readUnsignedLong);
    }

    public static void writeEventType(WriteBuffer writeBuffer, BACnetEventType value) throws SerializationException {
        if (value == null || value == BACnetEventType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("eventType", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static Long readProprietaryEventType(ReadBuffer readBuffer, BACnetEventType value, Long actualLength) throws ParseException {
        if (value != null && value != BACnetEventType.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryEventType", bitsToRead);
    }

    public static void writeProprietaryEventType(WriteBuffer writeBuffer, BACnetEventType eventType, long value) throws SerializationException {
        if (eventType != null && eventType != BACnetEventType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryEventType", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetEventType.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static BACnetEventState readEventState(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        int readUnsignedLong = readBuffer.readUnsignedInt("eventState", bitsToRead);
        if (!BACnetEventState.isDefined(readUnsignedLong)) {
            return BACnetEventState.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetEventState.enumForValue(readUnsignedLong);
    }

    public static void writeEventState(WriteBuffer writeBuffer, BACnetEventState value) throws SerializationException {
        if (value == null || value == BACnetEventState.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("eventState", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static Long readProprietaryEventState(ReadBuffer readBuffer, BACnetEventState value, Long actualLength) throws ParseException {
        if (value != null && value != BACnetEventState.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryEventState", bitsToRead);
    }

    public static void writeProprietaryEventState(WriteBuffer writeBuffer, BACnetEventState eventState, long value) throws SerializationException {
        if (eventState != null && eventState != BACnetEventState.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryEventState", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetEventState.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static BACnetObjectType readObjectType(ReadBuffer readBuffer) throws ParseException {
        int readUnsignedLong = readBuffer.readUnsignedInt("objectType", 10);
        if (!BACnetObjectType.isDefined(readUnsignedLong)) {
            return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetObjectType.enumForValue(readUnsignedLong);
    }

    public static void writeObjectType(WriteBuffer writeBuffer, BACnetObjectType value) throws SerializationException {
        if (value == null || value == BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedLong("objectType", 10, value.getValue(), WithAdditionalStringRepresentation(value.name()));
    }

    public static Integer readProprietaryObjectType(ReadBuffer readBuffer, BACnetObjectType value) throws ParseException {
        if (value != null && value != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return 0;
        }
        // We need to reset our reader to the position we read before
        // TODO: maybe we reset to much here because pos is byte based
        // we consume the leftover bits before we reset to avoid trouble // TODO: we really need bit precision on resetting
        readBuffer.readUnsignedInt(6);
        readBuffer.reset(readBuffer.getPos() - 2);
        return readBuffer.readUnsignedInt("proprietaryObjectType", 10);
    }

    public static void writeProprietaryObjectType(WriteBuffer writeBuffer, BACnetObjectType objectType, int value) throws SerializationException {
        if (objectType != null && objectType != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedInt("proprietaryObjectType", 10, value, WithAdditionalStringRepresentation(BACnetObjectType.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static boolean isBACnetConstructedDataClosingTag(ReadBuffer readBuffer, boolean instantTerminate, int expectedTagNumber) {
        if (instantTerminate) {
            return true;
        }
        int oldPos = readBuffer.getPos();
        try {
            // TODO: add graceful exit if we know already that we are at the end (we might need to add available bytes to reader)
            int tagNumber = readBuffer.readUnsignedInt(4);
            boolean isContextTag = readBuffer.readBit();
            int tagValue = readBuffer.readUnsignedInt(3);

            boolean foundOurClosingTag = isContextTag && tagNumber == expectedTagNumber && tagValue == 0x7;
            LOGGER.debug("Checking at pos pos:{}: tagNumber:{}, isContextTag:{}, tagValue:{}, expectedTagNumber:{}. foundOurClosingTag:{}", oldPos, tagNumber, isContextTag, tagValue, expectedTagNumber, foundOurClosingTag);
            return foundOurClosingTag;
        } catch (ParseException e) {
            LOGGER.warn("Error reading termination bit", e);
            return true;
        } catch (ArrayIndexOutOfBoundsException e) {
            LOGGER.debug("Reached EOF at {}", oldPos, e);
            return true;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static BACnetDataType guessDataType(BACnetObjectType objectType) {
        switch (objectType) {
            case ACCESS_CREDENTIAL:
            case ACCESS_DOOR:
            case ACCESS_POINT:
            case ACCESS_RIGHTS:
            case ACCESS_USER:
            case ACCESS_ZONE:
            case ACCUMULATOR:
                // TODO: temporary
                return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
            case ALERT_ENROLLMENT:
                // TODO: temporary
                return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
            case ANALOG_INPUT:
            case ANALOG_OUTPUT:
                return BACnetDataType.REAL;
            case ANALOG_VALUE:
            case AVERAGING:
            case BINARY_INPUT:
            case BINARY_LIGHTING_OUTPUT:
            case BINARY_OUTPUT:
            case BINARY_VALUE:
                // TODO: temporary
                return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
            case BITSTRING_VALUE:
                return BACnetDataType.BIT_STRING;
            case CALENDAR:
            case CHANNEL:
            case CHARACTERSTRING_VALUE:
            case COMMAND:
            case CREDENTIAL_DATA_INPUT:
            case DATEPATTERN_VALUE:
            case DATE_VALUE:
            case DATETIMEPATTERN_VALUE:
            case DATETIME_VALUE:
            case DEVICE:
            case ELEVATOR_GROUP:
            case ESCALATOR:
            case EVENT_ENROLLMENT:
            case EVENT_LOG:
            case FILE:
            case GLOBAL_GROUP:
            case GROUP:
                // TODO: temporary
                return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
            case INTEGER_VALUE:
                return BACnetDataType.SIGNED_INTEGER;
            case LARGE_ANALOG_VALUE:
            case LIFE_SAFETY_POINT:
                // TODO: temporary
                return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
            case LIFE_SAFETY_ZONE:
                return BACnetDataType.BACNET_OBJECT_IDENTIFIER;
            case LIFT:
            case LIGHTING_OUTPUT:
            case LOAD_CONTROL:
                // TODO: temporary // FIXME: this is just so tags get consumed
                return BACnetDataType.ENUMERATED;
            case LOOP:
            case MULTI_STATE_INPUT:
            case MULTI_STATE_OUTPUT:
            case MULTI_STATE_VALUE:
            case NETWORK_PORT:
            case NETWORK_SECURITY:
            case NOTIFICATION_CLASS:
            case NOTIFICATION_FORWARDER:
            case OCTETSTRING_VALUE:
            case POSITIVE_INTEGER_VALUE:
            case PROGRAM:
            case PULSE_CONVERTER:
            case SCHEDULE:
            case STRUCTURED_VIEW:
                // TODO: temporary
                return BACnetDataType.BACNET_OBJECT_IDENTIFIER;
            case TIMEPATTERN_VALUE:
            case TIME_VALUE:
            case TIMER:
            case TREND_LOG:
            case TREND_LOG_MULTIPLE:
                return BACnetDataType.ENUMERATED;
        }
        // TODO: temporary
        return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
    }

    public static long parseVarUint(byte[] data) {
        return new BigInteger(data).longValue();
    }

    public static byte[] writeVarUint(long value) {
        return BigInteger.valueOf(value).toByteArray();
    }

    public static BACnetTagHeader createBACnetTagHeaderBalanced(boolean isContext, short id, long value) {
        TagClass tagClass = TagClass.APPLICATION_TAGS;
        if (isContext) {
            tagClass = TagClass.CONTEXT_SPECIFIC_TAGS;
        }

        byte tagNumber;
        Short extTagNumber = null;
        if (id <= 14) {
            tagNumber = (byte) id;
        } else {
            tagNumber = 0xF;
            extTagNumber = id;
        }

        byte lengthValueType;
        Short extLength = null;
        Integer extExtLength = null;
        Long extExtExtLength = null;
        if (value <= 4) {
            lengthValueType = (byte) value;
        } else {
            lengthValueType = 5;
            // Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
            if (value <= 253) {
                extLength = (short) value;
            } else if (value <= 65535) {
                extLength = 254;
                extExtLength = (int) value;
            } else {
                extLength = 255;
                extExtExtLength = value;
            }
        }

        return new BACnetTagHeader(tagNumber, tagClass, lengthValueType, extTagNumber, extLength, extExtLength, extExtExtLength);
    }

    public static BACnetApplicationTagObjectIdentifier createBACnetApplicationTagObjectIdentifier(int objectType, long instance) {
        BACnetTagHeader header = new BACnetTagHeader((byte) 0xC, TagClass.APPLICATION_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, objectType, instance);
        return new BACnetApplicationTagObjectIdentifier(header, payload);
    }

    public static BACnetContextTagObjectIdentifier createBACnetContextTagObjectIdentifier(byte tagNum, int objectType, long instance) {
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, objectType, instance);
        return new BACnetContextTagObjectIdentifier(header, payload, (short) tagNum, true);
    }

    public static BACnetContextTagPropertyIdentifier createBACnetContextTagPropertyIdentifier(byte tagNum, int propertyType) {
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 4, null, null, null, null);
        BACnetPropertyIdentifier propertyIdentifier = BACnetPropertyIdentifier.enumForValue(propertyType);
        if (!BACnetPropertyIdentifier.isDefined(propertyType)) {
            propertyIdentifier = BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
        }
        return new BACnetContextTagPropertyIdentifier(header, propertyIdentifier, propertyType, (short) tagNum, true, 0L);
    }

    public static BACnetApplicationTagEnumerated createBACnetApplicationTagEnumerated(long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) 0x9, lengthPayload.getLeft());
        return new BACnetApplicationTagEnumerated(header, lengthPayload.getRight());
    }

    public static BACnetContextTagEnumerated CreateBACnetContextTagEnumerated(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagEnumerated(header, lengthPayload.getRight(), (short) tagNumber, true);
    }

    public static Pair<Long, BACnetTagPayloadEnumerated> CreateEnumeratedPayload(long value) {
        long length;
        if (value < 0x100)
            length = 1;
        else if (value < 0x10000)
            length = 2;
        else if (value < 0x1000000)
            length = 3;
        else
            length = 4;
        byte[] data = writeVarUint(value);
        BACnetTagPayloadEnumerated payload = new BACnetTagPayloadEnumerated(data, length);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagUnsignedInteger createBACnetApplicationTagUnsignedInteger(long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) 0x2, lengthPayload.getLeft());
        return new BACnetApplicationTagUnsignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagUnsignedInteger createBACnetContextTagUnsignedInteger(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagUnsignedInteger(header, lengthPayload.getRight(), (short) tagNumber, true);
    }

    public static Pair<Long, BACnetTagPayloadUnsignedInteger> createUnsignedPayload(long value) {
        long length;
        Short valueUint8 = null;
        Integer valueUint16 = null;
        Long valueUint24 = null;
        Long valueUint32 = null;
        BigInteger valueUint40 = null;
        BigInteger valueUint48 = null;
        BigInteger valueUint56 = null;
        BigInteger valueUint64 = null;
        if (value < 0x100) {
            length = 1;
            valueUint8 = (short) value;
        } else if (value < 0x10000) {
            length = 2;
            valueUint16 = (int) value;
        } else if (value < 0x1000000) {
            length = 3;
            valueUint24 = value;
        } else {
            length = 4;
            valueUint32 = value;
        }
        BACnetTagPayloadUnsignedInteger payload = new BACnetTagPayloadUnsignedInteger(valueUint8, valueUint16, valueUint24, valueUint32, valueUint40, valueUint48, valueUint56, valueUint64, length);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagSignedInteger createBACnetApplicationTagSignedInteger(long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) 0x3, lengthPayload.getLeft());
        return new BACnetApplicationTagSignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagSignedInteger createBACnetContextTagSignedInteger(short tagNumber, long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagSignedInteger(header, lengthPayload.getRight(), tagNumber, true);
    }

    public static Pair<Long, BACnetTagPayloadSignedInteger> createSignedPayload(long value) {
        long length;
        Byte valueInt8 = null;
        Short valueInt16 = null;
        Integer valueInt24 = null;
        Integer valueInt32 = null;
        if (value < 0x100) {
            length = 1;
            valueInt8 = (byte) value;
        } else if (value < 0x10000) {
            length = 2;
            valueInt16 = (short) value;
        } else if (value < 0x1000000) {
            length = 3;
            valueInt24 = (int) value;
        } else {
            length = 4;
            valueInt32 = (int) value;
        }
        BACnetTagPayloadSignedInteger payload = new BACnetTagPayloadSignedInteger(valueInt8, valueInt16, valueInt24, valueInt32, null, null, null, null, length);
        return Pair.of(length, payload);
    }

}
