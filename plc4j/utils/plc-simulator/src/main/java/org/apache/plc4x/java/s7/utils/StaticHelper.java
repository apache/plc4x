/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;

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

            return LocalDateTime.of(year, month, day, hour, minute, second, nanosecond);
        } catch (Exception e) {
            return null;
        }
    }

    public static void serializeTiaDateTime(WriteBuffer io, PlcValue value) {
        throw new NotImplementedException("Serializing DATE_AND_TIME not implemented");
    }

    public static String parseS7String(ReadBuffer io, int stringLength, Object encoding) {
        try {
            // This is the maximum number of bytes a string can be long.
            short maxLength = io.readUnsignedShort(8);
            // This is the total length of the string on the PLC (Not necessarily the number of characters read)
            short totalStringLength = io.readShort(8);
            // Read the full size of the string.
            String str = io.readString(stringLength * 8, WithOption.WithEncoding((String) encoding));
            // Cut off the parts that don't belong to it.
            return str.substring(0, totalStringLength);
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializeS7String(WriteBuffer io, PlcValue value, int stringLength, Object encoding) {
        // TODO: Need to implement the serialization or we can't write strings
        throw new NotImplementedException("Serializing STRING not implemented");
    }

    public static String parseS7Char(ReadBuffer io, Object encoding) throws ParseException {
        // Read the full size of the string.
        return io.readString(8, WithOption.WithEncoding((String) encoding));
    }

    public static void serializeS7Char(WriteBuffer io, PlcValue value, Object encoding) {
        // TODO: Need to implement the serialization or we can't write strings
        throw new NotImplementedException("Serializing STRING not implemented");
    }

}
