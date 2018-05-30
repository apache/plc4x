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

import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.commands.types.TimeStamp;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;

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

    public static Length getLengthFor(Class<?> clazz, long defaultValue) {
        if (Calendar.class.isAssignableFrom(clazz)) {
            return Length.of(8);
        }
        return Length.of(LENGTH_MAP.getOrDefault(clazz, defaultValue));
    }

    @SuppressWarnings("unchecked")
    public static <T> List<T> decodeData(Class<T> datatype, byte[] adsData) throws PlcProtocolException {
        List<Object> result = new LinkedList<>();
        int i = 0;
        final int length = adsData.length;

        // Expand arrays to avoid IndexOutOfBoundsException
        final byte[] safeLengthAdsData;
        if (datatype == Short.class && length < 2) {
            safeLengthAdsData = new byte[2];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if (datatype == Integer.class && length < 4) {
            safeLengthAdsData = new byte[4];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if (datatype == Float.class && length < 4) {
            safeLengthAdsData = new byte[4];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if (datatype == Double.class && length < 8) {
            safeLengthAdsData = new byte[8];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else if ((datatype == Calendar.class || Calendar.class.isAssignableFrom(datatype)) && length < 8) {
            safeLengthAdsData = new byte[8];
            System.arraycopy(adsData, 0, safeLengthAdsData, 0, length);
        } else {
            safeLengthAdsData = adsData;
        }

        while (i < length) {
            byte byteOne = safeLengthAdsData[i];
            if (datatype == String.class) {
                i = decodeString(safeLengthAdsData, i, length, result);
            } else if (datatype == Boolean.class) {
                result.add((byteOne & 0x01) == 0x01);
                i += 1;
            } else if (datatype == Byte.class) {
                result.add(byteOne);
                i += 1;
            } else if (datatype == Short.class) {
                decodeShort(safeLengthAdsData, i, result);
                i += 2;
            } else if (datatype == Integer.class) {
                decodeInteger(safeLengthAdsData, i, result);
                i += 4;
            } else if (datatype == Float.class) {
                decodeFloat(safeLengthAdsData, i, result);
                i += 4;
            } else if (datatype == Double.class) {
                decodeDouble(safeLengthAdsData, i, result);
                i += 8;
            } else if (datatype == Calendar.class || Calendar.class.isAssignableFrom(datatype)) {
                extractCalendar(safeLengthAdsData, i, result);
                i += 8;
            } else {
                throw new PlcProtocolException("Unsupported datatype " + datatype.getSimpleName());
            }
        }
        return (List<T>) result;
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
