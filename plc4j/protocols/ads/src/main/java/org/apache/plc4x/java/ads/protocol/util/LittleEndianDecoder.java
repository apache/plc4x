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

import org.apache.plc4x.java.ads.api.commands.types.TimeStamp;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;

import java.math.BigInteger;
import java.util.Calendar;
import java.util.Date;
import java.util.LinkedList;
import java.util.List;

public class LittleEndianDecoder {

    private LittleEndianDecoder() {
        // Utility class
    }

    @SuppressWarnings("unchecked")
    public static <T> List<T> decodeData(Class<T> datatype, byte[] adsData) throws PlcProtocolException {
        List<Object> result = new LinkedList<>();
        int i = 0;
        final int length = adsData.length;

        while (i < length) {
            byte byteOne = adsData[i];
            if (datatype == String.class) {
                i = decodeString(adsData, i, length, result);
            } else if (datatype == Boolean.class) {
                result.add((byteOne & 0x01) == 0x01);
                i += 1;
            } else if (datatype == Byte.class) {
                result.add(byteOne);
                i += 1;
            } else  if (datatype == Short.class) {
                decodeShort(adsData, i, result);
                i += 2;
            } else if (datatype == Integer.class) {
                decodeInteger(adsData, i, result);
                i += 4;
            } else if (datatype == Float.class) {
                decodeFloat(adsData, i, result);
                i += 4;
            } else if (datatype == Calendar.class || Calendar.class.isAssignableFrom(datatype)) {
                extractCalendar(adsData, i, result);
                i += 8;
            } else {
                throw new PlcProtocolException("Unsupported datatype " + datatype.getSimpleName());
            }
        }
        return (List<T>) result;
    }

    private static int decodeString(byte[] adsData, int i, int length, List<Object> result) {
        int pos = i;
        StringBuilder builder = new StringBuilder();
        while (adsData[pos] != (byte) 0x0 && pos < length) {
            builder.append((char) adsData[pos]);
            pos++;
        }
        pos++; // skip terminating character
        result.add(builder.toString());
        return pos;
    }

    private static void decodeShort(byte[] adsData, int i, List<Object> result) {
        byte byteOne = adsData[i];
        byte byteTwo = adsData[i+1];
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
