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
package org.apache.plc4x.java.s7.utils;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

public class StaticHelper {

    public static LocalTime parseTiaTime(ReadBuffer io) {
        try {
            int millisSinceMidnight = io.readInt(32);
            return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
                millisSinceMidnight, ChronoUnit.MILLIS);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeTiaTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing TIME not implemented");
    }

    public static LocalTime parseS5Time(ReadBuffer io) {
        try {
            int stuff = io.readInt(16);
            // TODO: Implement this correctly.
            throw new NotImplementedException("S5TIME not implemented");
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeS5Time(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing S5TIME not implemented");

    }

    public static LocalTime parseTiaLTime(ReadBuffer io) {
        throw new NotImplementedException("LTIME not implemented");
    }

    public static void serializeTiaLTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing LTIME not implemented");
    }

    public static LocalTime parseTiaTimeOfDay(ReadBuffer io) {
        try {
            long millisSinceMidnight = io.readUnsignedLong(32);
            return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
                millisSinceMidnight, ChronoUnit.MILLIS);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeTiaTimeOfDay(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing TIME_OF_DAY not implemented");
    }

    public static LocalDate parseTiaDate(ReadBuffer io) {
        try {
            int daysSince1990 = io.readUnsignedInt(16);
            return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeTiaDate(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing DATE not implemented");
    }

    //TODO: Call BCD converter
    public static LocalDateTime parseTiaDateTime(ReadBuffer io) {
        try {
            int year = io.readUnsignedInt(16);
            int month = io.readUnsignedInt(8);
            int day = io.readUnsignedInt(8);
            // Skip day-of-week
            io.readByte();
            int hour = io.readByte();
            int minute = io.readByte();
            int second = io.readByte();
            int nanosecond = io.readUnsignedInt(24);
            // Skip day-of-week
            io.readByte();

            return LocalDateTime.of(year, month, day, hour, minute, second, nanosecond);
        } catch (Exception e) {
            return null;
        }
    }

    public static void serializeTiaDateTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing DATE_AND_TIME not implemented");
    }

    public static String parseS7Char(ReadBuffer io, String encoding) {
        if ("UTF-8".equalsIgnoreCase(encoding)) {
            return io.readString(8, encoding);
        } else if ("UTF-16".equalsIgnoreCase(encoding)) {
            return io.readString(16, encoding);
        } else {
            throw new PlcRuntimeException("Unsupported encoding");
        }
    }

    public static String parseS7String(ReadBuffer io, int stringLength, String encoding) {
        try {
            if ("UTF-8".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                short maxLength = io.readUnsignedShort(8);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                short totalStringLength = io.readUnsignedShort(8);

                final byte[] byteArray = new byte[totalStringLength];
                for (int i = 0; (i < stringLength) && io.hasMore(8); i++) {
                    final byte curByte = io.readByte();
                    if (i < totalStringLength) {
                        byteArray[i] = curByte;
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for (; (i < stringLength) && io.hasMore(8); i++) {
                            io.readByte();
                        }
                        break;
                    }
                }
                return new String(byteArray, StandardCharsets.UTF_8);
            } else if ("UTF-16".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                int maxLength = io.readUnsignedInt(16);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                int totalStringLength = io.readUnsignedInt(16);

                final byte[] byteArray = new byte[totalStringLength * 2];
                for (int i = 0; (i < stringLength) && io.hasMore(16); i++) {
                    final short curShort = io.readShort(16);
                    if (i < totalStringLength) {
                        byteArray[i * 2] = (byte) (curShort >>> 8);
                        byteArray[(i * 2) + 1] = (byte) (curShort & 0xFF);
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for (; (i < stringLength) && io.hasMore(16); i++) {
                            io.readShort(16);
                        }
                        break;
                    }
                }
                return new String(byteArray, StandardCharsets.UTF_16);
            } else {
                throw new PlcRuntimeException("Unsupported string encoding " + encoding);
            }
        } catch (ParseException e) {
            throw new PlcRuntimeException("Error parsing string", e);
        }
    }

    /*
     * A variable of data type CHAR (character) occupies one byte.
     */

    public static void serializeS7Char(WriteBuffer io, PlcValue value, String encoding) {
        // TODO: Need to implement the serialization or we can't write strings
        if ("UTF-8".equalsIgnoreCase(encoding)) {
            //return io.readString(8, encoding);
        } else if ("UTF-16".equalsIgnoreCase(encoding)) {
            //return io.readString(16, encoding);
        } else {
            throw new PlcRuntimeException("Unsupported encoding");
        }
    }

    /*           +-------------------+
     * Byte n     | Maximum length    | (k)
     *            +-------------------+
     * Byte n+1   | Current Length    | (m)
     *            +-------------------+
     * Byte n+2   | 1st character     | \         \
     *            +-------------------+  |         |
     * Byte n+3   | 2st character     |  | Current |
     *            +-------------------+   >        |
     * Byte ...   | ...               |  | length  |  Maximum
     *            +-------------------+  |          >
     * Byte n+m+1 | mth character     | /          |  length
     *            +-------------------+            |
     * Byte ...   | ...               |            |
     *            +-------------------+            |
     * Byte ...   | ...               |           /
     *            +-------------------+
     * For this version, the user must read the maximum acceptable length in
     * the string in a first instance.
     * Then the user application should avoid the envelope of the adjacent
     * fields passing the maximum length in "stringLength".
     * If your application does not handle S7string, you can handle
     * the String as char arrays from your application.
     */
    public static void serializeS7String(WriteBuffer io, PlcValue value, int stringLength, String encoding) {
        byte k = (byte) ((stringLength > 250) ? 250 : stringLength);
        byte m = (byte) value.getString().length();
        m = (m > k) ? k : m;
        byte[] chars = new byte[m];
        for (int i = 0; i < m; ++i) {
            char c = value.getString().charAt(i);
            chars[i] = (byte) c;
        }

        try {
            io.writeByte(k);
            io.writeByte(m);
            io.writeByteArray(chars);
        } catch (ParseException ex) {
            Logger.getLogger(StaticHelper.class.getName()).log(Level.SEVERE, null, ex);
        }
    }
}
