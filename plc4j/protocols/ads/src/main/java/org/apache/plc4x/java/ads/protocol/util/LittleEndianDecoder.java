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
package org.apache.plc4x.java.ads.protocol.util;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.ads.api.commands.types.TimeStamp;
import org.apache.plc4x.java.ads.model.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.base.messages.items.BooleanFieldItem;
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.apache.plc4x.java.base.messages.items.IntegerFieldItem;

import java.math.BigInteger;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

// TODO: we might user ByteBuffer.wrap(buffer).order(ByteOrder.LITTLE_ENDIAN).putInt(port).asArray() etc
public class LittleEndianDecoder {

    private static final Map<Class<?>, Long> LENGTH_MAP;

    static {
        Map<Class<?>, Long> lengthMap = new ConcurrentHashMap<>();
        lengthMap.put(Boolean.class, 1L);
        lengthMap.put(Byte.class, 1L);
        lengthMap.put(Short.class, 2L);
        lengthMap.put(Integer.class, 4L);
        lengthMap.put(Float.class, 4L);
        lengthMap.put(Double.class, 8L);
        lengthMap.put(Calendar.class, 8L);
        LENGTH_MAP = Collections.unmodifiableMap(lengthMap);
    }

    private LittleEndianDecoder() {
        // Utility class
    }

    public static long getLengthFor(Class<?> clazz, long defaultValue) {
        if (Calendar.class.isAssignableFrom(clazz)) {
            return 8;
        }
        return LENGTH_MAP.getOrDefault(clazz, defaultValue);
    }

    @SuppressWarnings("unchecked")
    public static FieldItem<?> decodeData(AdsDataType adsDataType, byte[] adsData) {
        ByteBuf wrappedBuffer = Unpooled.wrappedBuffer(adsData);
        switch (adsDataType) {
            case BIT: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                return new BooleanFieldItem(values.toArray(new Boolean[0]));
            }
            case BIT8: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                return new BooleanFieldItem(values.toArray(new Boolean[0]));
            }
            case BITARR8: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer((long) aByte);
                }
                return new IntegerFieldItem(values.toArray(new Long[0]));
            }
            case BITARR16: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aLong);
                }
                return new IntegerFieldItem(values.toArray(new Long[0]));
            }
            case BITARR32: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case INT8: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case INT16: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case INT32: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case INT64: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UINT8: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UINT16: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UINT32: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UINT64: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case FLOAT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case DOUBLE: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case BOOL: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case BYTE: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case WORD: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case DWORD: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case SINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case USINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case INT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case DINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UDINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case LINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case ULINT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case REAL: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case LREAL: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case STRING: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case TIME: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case TIME_OF_DAY: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case DATE: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case DATE_AND_TIME: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            default:
                throw new IllegalArgumentException("Unsupported adsDataType " + adsDataType);
        }
        /*
        if (dataType == byte[].class) {
            return (List<T>) Collections.singletonList(adsData);
        }
        if (dataType == Byte[].class) {
            return (List<T>) Collections.singletonList(ArrayUtils.toObject(adsData));
        }
        List<Object> result = new LinkedList<>();
        int i = 0;
        final int length = adsData.length;

        // Expand arrays to avoid IndexOutOfBoundsException
        final byte[] safeLengthAdsData;
        if (dataType == Short.class && length < 2) {
            safeLengthAdsData = new byte[2];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if (dataType == Integer.class && length < 4) {
            safeLengthAdsData = new byte[4];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if (dataType == Float.class && length < 4) {
            safeLengthAdsData = new byte[4];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if (dataType == Double.class && length < 8) {
            safeLengthAdsData = new byte[8];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if ((dataType == Calendar.class || Calendar.class.isAssignableFrom(dataType)) && length < 8) {
            safeLengthAdsData = new byte[8];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else {
            safeLengthAdsData = adsData;
        }

        while (i < length) {
            byte byteOne = safeLengthAdsData[i];
            if (dataType == String.class) {
                i = decodeString(safeLengthAdsData, i, length, result);
            } else if (dataType == Boolean.class) {
                result.add((byteOne & 0x01) == 0x01);
                i += 1;
            } else if (dataType == Byte.class) {
                result.add(byteOne);
                i += 1;
            } else if (dataType == Short.class) {
                decodeShort(safeLengthAdsData, i, result);
                i += 2;
            } else if (dataType == Integer.class) {
                decodeInteger(safeLengthAdsData, i, result);
                i += 4;
            } else if (dataType == BigInteger.class) {
                decodeBigInteger(safeLengthAdsData, result);
                // A big integer can consume the whole stream
                i = length;
            } else if (dataType == Float.class) {
                decodeFloat(safeLengthAdsData, i, result);
                i += 4;
            } else if (dataType == Double.class) {
                decodeDouble(safeLengthAdsData, i, result);
                i += 8;
            } else if (dataType == Calendar.class || Calendar.class.isAssignableFrom(dataType)) {
                extractCalendar(safeLengthAdsData, i, result);
                i += 8;
            } else {
                throw new PlcUnsupportedDataTypeException(dataType);
            }
        }
        return (List<T>) result;
        */
    }

    private static int decodeString(byte[] adsData, int i, int length, List<Object> result) throws PlcProtocolException {
        if (length < 2) {
            throw new PlcProtocolException("String must be a null terminated byte array");
        }
        int pos = i;
        StringBuilder builder = new StringBuilder();
        // TODO: check if we have at least a 0x0 space
        while (pos < length && adsData[pos] != (byte) 0x0) {
            builder.append((char) adsData[pos]);
            pos++;
        }
        pos++; // skip terminating character
        result.add(builder.toString());
        return pos;
    }

    private static void decodeShort(byte[] adsData, int i, List<Object> result) {
        byte byteOne = adsData[i];
        byte byteTwo = adsData[i + 1];
        result.add((short) ((byteOne & 0xff) | ((byteTwo & 0xff) << 8)));
    }

    private static void decodeInteger(byte[] adsData, int i, List<Object> result) {
        byte byteOne = adsData[i];
        byte byteTwo = adsData[i + 1];
        byte byteThree = adsData[i + 2];
        byte byteFour = adsData[i + 3];
        result.add((byteOne & 0xff) | ((byteTwo & 0xff) << 8) | ((byteThree & 0xff) << 16) | (byteFour & 0xff) << 24);
    }

    private static void decodeBigInteger(byte[] adsData, List<Object> result) {
        byte[] clone = ArrayUtils.clone(adsData);
        // In ADS data is transferred Little Endian
        ArrayUtils.reverse(clone);
        // Adding a 0 ensures that this is interpreted as a positive
        // number as the most significant bit is guaranteed to be 0.
        byte[] bigIntegerByteArray = ArrayUtils.insert(0, clone, (byte) 0x0);
        result.add(new BigInteger(bigIntegerByteArray));
    }

    private static void decodeFloat(byte[] adsData, int i, List<Object> result) {
        byte byteOne = adsData[i];
        byte byteTwo = adsData[i + 1];
        byte byteThree = adsData[i + 2];
        byte byteFour = adsData[i + 3];
        // TODO: check how ads expects this data
        // Description of the Real number format:
        // https://www.sps-lehrgang.de/zahlenformate-step7/#c144
        // https://de.wikipedia.org/wiki/IEEE_754
        int intValue = (byteOne & 0xff) | ((byteTwo & 0xff) << 8) | ((byteThree & 0xff) << 16) | ((byteFour & 0xff) << 24);
        result.add(Float.intBitsToFloat(intValue));
    }

    private static void decodeDouble(byte[] adsData, int i, List<Object> result) {
        byte byteOne = adsData[i];
        byte byteTwo = adsData[i + 1];
        byte byteThree = adsData[i + 2];
        byte byteFour = adsData[i + 3];
        byte byteFive = adsData[i + 4];
        byte byteSix = adsData[i + 5];
        byte byteSeven = adsData[i + 6];
        byte byteEigth = adsData[i + 7];
        // TODO: check how ads expects this data
        // Description of the Real number format:
        // https://www.sps-lehrgang.de/zahlenformate-step7/#c144
        // https://de.wikipedia.org/wiki/IEEE_754
        long longValue = (long) (byteOne & 0xff) | ((long) (byteTwo & 0xff) << 8) | ((long) (byteThree & 0xff) << 16) | ((long) (byteFour & 0xff) << 24)
            | (long) (byteFive & 0xff) << 32 | ((long) (byteSix & 0xff) << 40) | ((long) (byteSeven & 0xff) << 48) | ((long) (byteEigth & 0xff) << 56);
        result.add(Double.longBitsToDouble(longValue));
    }

    private static void extractCalendar(byte[] adsData, int i, List<Object> result) {
        byte byteOne = adsData[i];
        byte byteTwo = adsData[i + 1];
        byte byteThree = adsData[i + 2];
        byte byteFour = adsData[i + 3];
        byte byteFive = adsData[i + 4];
        byte byteSix = adsData[i + 5];
        byte byteSeven = adsData[i + 6];
        byte byteEight = adsData[i + 7];
        Calendar calendar = Calendar.getInstance();
        calendar.setTime(new Date(TimeStamp.winTimeToJava(
            new BigInteger(new byte[]{
                // LE
                byteEight,
                byteSeven,
                byteSix,
                byteFive,
                byteFour,
                byteThree,
                byteTwo,
                byteOne,
            })).longValue()));
        result.add(calendar);
    }
}
