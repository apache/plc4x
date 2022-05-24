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

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithAdditionalStringRepresentation;

public class StaticHelper {

    public static final Logger LOGGER = LoggerFactory.getLogger(StaticHelper.class);

    public static Object readEnumGenericFailing(ReadBuffer readBuffer, Long actualLength, Enum<?> template) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        long rawValue = readBuffer.readUnsignedLong("value", bitsToRead);
        Class<?> declaringClass = template.getDeclaringClass();
        if (declaringClass == BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class) {
            if (!BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.isDefined((short) rawValue))
                throw new ParseException("Invalid value " + rawValue + " for " + BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class.getSimpleName());
            return BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class) {
            if (!BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.isDefined((short) rawValue))
                throw new ParseException("Invalid value " + rawValue + " for " + BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class.getSimpleName());
            return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetSegmentation.class) {
            if (!BACnetSegmentation.isDefined((short) rawValue))
                throw new ParseException("Invalid value " + rawValue + " for " + BACnetSegmentation.class.getSimpleName());
            return BACnetSegmentation.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetAction.class) {
            if (!BACnetAction.isDefined((short) rawValue))
                throw new ParseException("Invalid value " + rawValue + " for " + BACnetAction.class.getSimpleName());
            return BACnetAction.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetNotifyType.class) {
            if (!BACnetNotifyType.isDefined((short) rawValue))
                throw new ParseException("Invalid value " + rawValue + " for " + BACnetBinaryPV.class.getSimpleName());
            return BACnetNotifyType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetBinaryPV.class) {
            if (!BACnetBinaryPV.isDefined((short) rawValue))
                throw new ParseException("Invalid value " + rawValue + " for " + BACnetBinaryPV.class.getSimpleName());
            return BACnetBinaryPV.enumForValue((short) rawValue);
        }
        throw new ParseException("Unmapped type " + declaringClass);
    }

    public static Object readEnumGeneric(ReadBuffer readBuffer, Long actualLength, Enum<?> template) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        long rawValue = readBuffer.readUnsignedLong("value", bitsToRead);
        // TODO: map types here for better performance which doesn't use reflection
        Class<?> declaringClass = template.getDeclaringClass();
        if (declaringClass == BACnetAbortReason.class) {
            if (!BACnetAbortReason.isDefined((short) rawValue)) return BACnetAbortReason.VENDOR_PROPRIETARY_VALUE;
            return BACnetAbortReason.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetEventState.class) {
            if (!BACnetEventState.isDefined((short) rawValue)) return BACnetEventState.VENDOR_PROPRIETARY_VALUE;
            return BACnetEventState.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetEventType.class) {
            if (!BACnetEventType.isDefined((short) rawValue)) return BACnetEventType.VENDOR_PROPRIETARY_VALUE;
            return BACnetEventType.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetLifeSafetyMode.class) {
            if (!BACnetLifeSafetyMode.isDefined((short) rawValue)) return BACnetLifeSafetyMode.VENDOR_PROPRIETARY_VALUE;
            return BACnetLifeSafetyMode.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetLifeSafetyState.class) {
            if (!BACnetLifeSafetyState.isDefined((short) rawValue))
                return BACnetLifeSafetyState.VENDOR_PROPRIETARY_VALUE;
            return BACnetLifeSafetyState.enumForValue((int) rawValue);
        } else if (declaringClass == BACnetNetworkType.class) {
            if (!BACnetNetworkType.isDefined((short) rawValue)) return BACnetNetworkType.VENDOR_PROPRIETARY_VALUE;
            return BACnetNetworkType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetObjectType.class) {
            if (!BACnetObjectType.isDefined((short) rawValue)) return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            return BACnetObjectType.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetPropertyIdentifier.class) {
            if (!BACnetPropertyIdentifier.isDefined((short) rawValue))
                return BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
            return BACnetPropertyIdentifier.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetReliability.class) {
            if (!BACnetReliability.isDefined((short) rawValue)) return BACnetReliability.VENDOR_PROPRIETARY_VALUE;
            return BACnetReliability.enumForValue((short) rawValue);
        } else if (declaringClass == ErrorClass.class) {
            if (!ErrorClass.isDefined((short) rawValue)) return ErrorClass.VENDOR_PROPRIETARY_VALUE;
            return ErrorClass.enumForValue((short) rawValue);
        } else if (declaringClass == ErrorCode.class) {
            if (!ErrorCode.isDefined((short) rawValue)) return ErrorCode.VENDOR_PROPRIETARY_VALUE;
            return ErrorCode.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetRejectReason.class) {
            if (!BACnetRejectReason.isDefined((short) rawValue)) return BACnetRejectReason.VENDOR_PROPRIETARY_VALUE;
            return BACnetRejectReason.enumForValue((short) rawValue);
        } else if (declaringClass == BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class) {
            if (!BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.isDefined((short) rawValue))
                return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.VENDOR_PROPRIETARY_VALUE;
            return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.enumForValue((short) rawValue);
        } else {
            LOGGER.warn("using reflection for {}", declaringClass);
            Optional<Method> enumForValue = Arrays.stream(declaringClass.getDeclaredMethods()).filter(method -> method.getName().equals("enumForValue")).findAny();
            if (!enumForValue.isPresent()) {
                throw new ParseException("No enumForValue available");
            }
            Method method = enumForValue.get();
            try {
                Class<?> parameterType = method.getParameterTypes()[0];
                Object paramValue = null;
                if (parameterType == byte.class || parameterType == Byte.class) {
                    paramValue = (byte) rawValue;
                } else if (parameterType == short.class || parameterType == Short.class) {
                    paramValue = (short) rawValue;
                } else if (parameterType == int.class || parameterType == Integer.class) {
                    paramValue = (int) rawValue;
                }
                Object result = method.invoke(null, paramValue);
                if (result == null) {
                    return Enum.valueOf(template.getDeclaringClass(), "VENDOR_PROPRIETARY_VALUE");
                }
                return result;
            } catch (IllegalAccessException | InvocationTargetException e) {
                throw new ParseException("error invoking method", e);
            }
        }
    }

    public static long readProprietaryEnumGeneric(ReadBuffer readBuffer, Long actualLength, boolean shouldRead) throws ParseException {
        if (!shouldRead) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryValue", bitsToRead);
    }

    public static void writeEnumGeneric(WriteBuffer writeBuffer, Enum<?> value) throws SerializationException {
        if (value == null) {
            return;
        }
        int bitsToWrite;
        long valueValue;
        // TODO: map types here for better performance which doesn't use reflection
        if (value.getDeclaringClass() == BACnetAbortReason.class) {
            valueValue = ((BACnetAbortReason) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEventState.class) {
            valueValue = ((BACnetEventState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetEventType.class) {
            valueValue = ((BACnetEventType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLifeSafetyMode.class) {
            valueValue = ((BACnetLifeSafetyMode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetLifeSafetyState.class) {
            valueValue = ((BACnetLifeSafetyState) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNetworkType.class) {
            valueValue = ((BACnetNetworkType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNodeType.class) {
            valueValue = ((BACnetNodeType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetNotifyType.class) {
            valueValue = ((BACnetNotifyType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetObjectType.class) {
            valueValue = ((BACnetObjectType) value).getValue();
        } else if (value.getDeclaringClass() == BACnetPropertyIdentifier.class) {
            valueValue = ((BACnetPropertyIdentifier) value).getValue();
        } else if (value.getDeclaringClass() == BACnetReliability.class) {
            valueValue = ((BACnetReliability) value).getValue();
        } else if (value.getDeclaringClass() == ErrorClass.class) {
            valueValue = ((ErrorClass) value).getValue();
        } else if (value.getDeclaringClass() == ErrorCode.class) {
            valueValue = ((ErrorCode) value).getValue();
        } else if (value.getDeclaringClass() == BACnetRejectReason.class) {
            valueValue = ((BACnetRejectReason) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice.class) {
            valueValue = ((BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice) value).getValue();
        } else if (value.getDeclaringClass() == BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable.class) {
            valueValue = ((BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable) value).getValue();
        } else {
            LOGGER.warn("using reflection for {}", value.getDeclaringClass());
            try {
                valueValue = ((Number) FieldUtils.getDeclaredField(value.getDeclaringClass(), "value", true).get(value)).longValue();
            } catch (IllegalAccessException e) {
                throw new SerializationException("error accessing value", e);
            }
        }

        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("value", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static void writeProprietaryEnumGeneric(WriteBuffer writeBuffer, long value, boolean shouldWrite) throws SerializationException {
        if (!shouldWrite) {
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
        writeBuffer.writeUnsignedLong("proprietaryValue", bitsToWrite, value, WithAdditionalStringRepresentation("VENDOR_PROPRIETARY_VALUE"));
    }

    @Deprecated
    public static BACnetObjectType readObjectType(ReadBuffer readBuffer) throws ParseException {
        int readUnsignedLong = readBuffer.readUnsignedInt("objectType", 10);
        if (!BACnetObjectType.isDefined(readUnsignedLong)) {
            return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetObjectType.enumForValue(readUnsignedLong);
    }

    @Deprecated
    public static void writeObjectType(WriteBuffer writeBuffer, BACnetObjectType value) throws SerializationException {
        if (value == null || value == BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedLong("objectType", 10, value.getValue(), WithAdditionalStringRepresentation(value.name()));
    }

    @Deprecated
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

    @Deprecated
    public static void writeProprietaryObjectType(WriteBuffer writeBuffer, BACnetObjectType objectType, int value) throws SerializationException {
        if (objectType != null && objectType != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedInt("proprietaryObjectType", 10, value, WithAdditionalStringRepresentation(BACnetObjectType.VENDOR_PROPRIETARY_VALUE.name()));
    }

    @Deprecated
    public static BACnetObjectType mapBACnetObjectType(BACnetContextTagEnumerated rawObjectType) {
        if (rawObjectType == null) return null;
        BACnetObjectType baCnetObjectType = BACnetObjectType.enumForValue((int) rawObjectType.getActualValue());
        if (baCnetObjectType == null) return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        return baCnetObjectType;
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

    public static long parseVarUint(byte[] data) {
        if (data.length == 0) {
            return 0;
        }
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

    public static BACnetApplicationTagNull createBACnetApplicationTagNull() {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.NULL.getValue(), 0);
        return new BACnetApplicationTagNull(header);
    }

    public static BACnetContextTagNull createBACnetContextTagNull(byte tagNumber) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 0);
        return new BACnetContextTagNull(header, (short) tagNumber);
    }

    public static BACnetOpeningTag createBACnetOpeningTag(short tagNum) {
        byte tagNumber;
        Short extTagNumber = null;
        if (tagNum <= 14) {
            tagNumber = (byte) tagNum;
        } else {
            tagNumber = 0xF;
            extTagNumber = tagNum;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNumber, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 0x6, extTagNumber, null, null, null);
        return new BACnetOpeningTag(header, tagNum);
    }

    public static BACnetClosingTag createBACnetClosingTag(short tagNum) {
        byte tagNumber;
        Short extTagNumber = null;
        if (tagNum <= 14) {
            tagNumber = (byte) tagNum;
        } else {
            tagNumber = 0xF;
            extTagNumber = tagNum;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNumber, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 0x7, extTagNumber, null, null, null);
        return new BACnetClosingTag(header, tagNum);
    }

    public static BACnetApplicationTagObjectIdentifier createBACnetApplicationTagObjectIdentifier(int objectType, long instance) {
        BACnetTagHeader header = new BACnetTagHeader((byte) BACnetDataType.SIGNED_INTEGER.getValue(), TagClass.APPLICATION_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        int proprietaryValue = 0;
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = objectType;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance);
        return new BACnetApplicationTagObjectIdentifier(header, payload);
    }

    public static BACnetContextTagObjectIdentifier createBACnetContextTagObjectIdentifier(byte tagNum, int objectType, long instance) {
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        int proprietaryValue = 0;
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = objectType;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance);
        return new BACnetContextTagObjectIdentifier(header, payload, (short) tagNum);
    }

    public static BACnetPropertyIdentifierTagged createBACnetPropertyIdentifierTagged(byte tagNum, int propertyType) {
        BACnetPropertyIdentifier propertyIdentifier = BACnetPropertyIdentifier.enumForValue(propertyType);
        int proprietaryValue = 0;
        if (!BACnetPropertyIdentifier.isDefined(propertyType)) {
            propertyIdentifier = BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = propertyType;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) requiredLength(propertyType), null, null, null, null);
        return new BACnetPropertyIdentifierTagged(header, propertyIdentifier, proprietaryValue, (short) tagNum, TagClass.CONTEXT_SPECIFIC_TAGS);
    }

    public static BACnetVendorIdTagged createBACnetVendorIdApplicationTagged(int vendorId) {
        BACnetVendorId baCnetVendorId = BACnetVendorId.enumForValue(vendorId);
        int unknownVendorId = 0;
        if (!BACnetVendorId.isDefined(vendorId)) {
            baCnetVendorId = BACnetVendorId.UNKNOWN_VENDOR;
            unknownVendorId = vendorId;
        }
        BACnetTagHeader header = new BACnetTagHeader((byte) 0x2, TagClass.APPLICATION_TAGS, (byte) requiredLength(vendorId), null, null, null, null);
        return new BACnetVendorIdTagged(header, baCnetVendorId, unknownVendorId, (short) 0x2, TagClass.APPLICATION_TAGS);
    }

    public static BACnetVendorIdTagged createBACnetVendorIdContextTagged(byte tagNum, int vendorId) {
        BACnetVendorId baCnetVendorId = BACnetVendorId.enumForValue(vendorId);
        int unknownVendorId = 0;
        if (!BACnetVendorId.isDefined(vendorId)) {
            baCnetVendorId = BACnetVendorId.UNKNOWN_VENDOR;
            unknownVendorId = vendorId;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) requiredLength(vendorId), null, null, null, null);
        return new BACnetVendorIdTagged(header, baCnetVendorId, unknownVendorId, (short) tagNum, TagClass.CONTEXT_SPECIFIC_TAGS);
    }

    public static BACnetSegmentationTagged creatBACnetSegmentationTagged(BACnetSegmentation value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) 0, 1);
        return new BACnetSegmentationTagged(header, value, (short) 0, TagClass.APPLICATION_TAGS);
    }

    public static BACnetApplicationTagBoolean createBACnetApplicationTagBoolean(boolean value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.BOOLEAN.getValue(), value ? 1L : 0L);
        return new BACnetApplicationTagBoolean(header, new BACnetTagPayloadBoolean(value ? 1L : 0L));
    }

    public static BACnetContextTagBoolean createBACnetContextTagBoolean(byte tagNumber, boolean value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 1);
        return new BACnetContextTagBoolean(header, (short) (value ? 1 : 0), new BACnetTagPayloadBoolean(value ? 1L : 0L), (short) tagNumber);
    }

    public static BACnetApplicationTagUnsignedInteger createBACnetApplicationTagUnsignedInteger(long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.UNSIGNED_INTEGER.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagUnsignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagUnsignedInteger createBACnetContextTagUnsignedInteger(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagUnsignedInteger(header, lengthPayload.getRight(), (short) tagNumber);
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
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.SIGNED_INTEGER.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagSignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagSignedInteger createBACnetContextTagSignedInteger(short tagNumber, long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagSignedInteger(header, lengthPayload.getRight(), tagNumber);
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

    public static BACnetApplicationTagReal createBACnetApplicationTagReal(float value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.REAL.getValue(), 4);
        return new BACnetApplicationTagReal(header, new BACnetTagPayloadReal(value));
    }

    public static BACnetContextTagReal createBACnetContextTagReal(byte tagNumber, float value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        return new BACnetContextTagReal(header, new BACnetTagPayloadReal(value), (short) tagNumber);
    }

    public static BACnetApplicationTagDouble createBACnetApplicationTagDouble(double value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.DOUBLE.getValue(), 8);
        return new BACnetApplicationTagDouble(header, new BACnetTagPayloadDouble(value));
    }

    public static BACnetContextTagDouble createBACnetContextTagDouble(byte tagNumber, double value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 8);
        return new BACnetContextTagDouble(header, new BACnetTagPayloadDouble(value), (short) tagNumber);
    }

    public static BACnetApplicationTagOctetString createBACnetApplicationTagOctetString(byte[] octets) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.OCTET_STRING.getValue(), octets.length + 1);
        return new BACnetApplicationTagOctetString(header, new BACnetTagPayloadOctetString(octets, (long) octets.length + 1));
    }

    public static BACnetContextTagOctetString createBACnetContextTagOctetString(byte tagNumber, byte[] octets) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, octets.length + 1);
        return new BACnetContextTagOctetString(header, new BACnetTagPayloadOctetString(octets, (long) octets.length + 1), (short) tagNumber);
    }

    public static BACnetApplicationTagCharacterString createBACnetApplicationTagCharacterString(BACnetCharacterEncoding baCnetCharacterEncoding, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.CHARACTER_STRING.getValue(), value.length() + 1);
        return new BACnetApplicationTagCharacterString(header, new BACnetTagPayloadCharacterString(baCnetCharacterEncoding, value, (long) value.length() + 1));
    }

    public static BACnetContextTagCharacterString createBACnetContextTagCharacterString(byte tagNumber, BACnetCharacterEncoding baCnetCharacterEncoding, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, value.length() + 1);
        return new BACnetContextTagCharacterString(header, new BACnetTagPayloadCharacterString(baCnetCharacterEncoding, value, (long) value.length() + 1), (short) tagNumber);
    }

    public static BACnetApplicationTagBitString createBACnetApplicationTagBitString(List<Boolean> value) {
        long numberOfBytesNeeded = (value.size() + 7) / 8;
        short unusedBits = (short) (8 - (value.size() % 8));
        if (unusedBits == 8) {
            unusedBits = 0;
        }
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.BIT_STRING.getValue(), numberOfBytesNeeded + 1);
        return new BACnetApplicationTagBitString(header, new BACnetTagPayloadBitString(unusedBits, value, new ArrayList<>(unusedBits), numberOfBytesNeeded + 1));
    }

    public static BACnetContextTagBitString createBACnetContextTagBitString(byte tagNumber, List<Boolean> value) {
        long numberOfBytesNeeded = (value.size() + 7) / 8;
        short unusedBits = (short) (8 - (value.size() % 8));
        if (unusedBits == 8) {
            unusedBits = 0;
        }
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, numberOfBytesNeeded + 1);
        return new BACnetContextTagBitString(header, new BACnetTagPayloadBitString(unusedBits, value, new ArrayList<>(unusedBits), numberOfBytesNeeded + 1), (short) tagNumber);
    }

    public static BACnetApplicationTagEnumerated createBACnetApplicationTagEnumerated(long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.ENUMERATED.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagEnumerated(header, lengthPayload.getRight());
    }

    public static BACnetContextTagEnumerated createBACnetContextTagEnumerated(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagEnumerated(header, lengthPayload.getRight(), (short) tagNumber);
    }

    public static Pair<Long, BACnetTagPayloadEnumerated> CreateEnumeratedPayload(long value) {
        long length = requiredLength(value);
        byte[] data = writeVarUint(value);
        BACnetTagPayloadEnumerated payload = new BACnetTagPayloadEnumerated(data, length);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagDate createBACnetApplicationTagDate(int year, short month, short dayOfMonth, short dayOfWeek) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.DATE.getValue(), 4);
        short yearMinus1900 = (short) (year - 1900);
        if (year == 0xFF) {
            yearMinus1900 = 0xFF;
        }
        return new BACnetApplicationTagDate(header, new BACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek));
    }

    public static BACnetContextTagDate createBACnetContextTagDate(byte tagNumber, int year, short month, short dayOfMonth, short dayOfWeek) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        short yearMinus1900 = (short) (year - 1900);
        if (year == 0xFF) {
            yearMinus1900 = 0xFF;
        }
        return new BACnetContextTagDate(header, new BACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek), (short) tagNumber);
    }

    public static BACnetApplicationTagTime createBACnetApplicationTagTime(short hour, short minute, short second, short fractional) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.TIME.getValue(), 4);
        return new BACnetApplicationTagTime(header, new BACnetTagPayloadTime(hour, minute, second, fractional));
    }

    public static BACnetContextTagTime createBACnetContextTagTime(byte tagNumber, short hour, short minute, short second, short fractional) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        return new BACnetContextTagTime(header, new BACnetTagPayloadTime(hour, minute, second, fractional), (short) tagNumber);
    }

    private static long requiredLength(long value) {
        long length;
        if (value < 0x100) length = 1;
        else if (value < 0x10000) length = 2;
        else if (value < 0x1000000) length = 3;
        else length = 4;
        return length;
    }

}
