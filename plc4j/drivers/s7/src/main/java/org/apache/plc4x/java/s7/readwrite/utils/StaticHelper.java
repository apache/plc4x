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
package org.apache.plc4x.java.s7.readwrite.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.values.PlcDATE;
import org.apache.plc4x.java.spi.values.PlcTIME;

import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDate;
import java.time.temporal.ChronoUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

public class StaticHelper {

    /**
     * Overrides the active per-call encoding options whenever we drop into a manual
     * parser/serializer from a string-typed mspec context. The generated
     * {@code readManualField}/{@code writeManualField} stack pushes
     * {@code WithUnsignedIntegerEncoding("UTF8")} (and friends) onto the buffer because
     * the mspec field is a string; if we then call {@code readUnsignedShort} for the
     * length prefix, that integer-read picks up the string encoding from context and
     * blows up. These constants force the right encoding for each numeric primitive.
     */
    private static final WithOption[] UINT_OPT = {
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary")
    };
    private static final WithOption[] SINT_OPT = {
        WithOption.WithSignedIntegerEncoding("twos-complement")
    };

    public static void leftShift3(final WriteBuffer buffer, int _value) throws BufferException {
        int valor = _value << 3;
        buffer.writeUnsignedInt(16, valor);
    }

    public static int rightShift3(final ReadBuffer buffer, DataTransportSize tsize) throws BufferException {
        int value = 0;
        if ((tsize == DataTransportSize.OCTET_STRING) ||
            (tsize == DataTransportSize.REAL) ||
            (tsize == DataTransportSize.BIT)) {
            value = buffer.readUnsignedInt(16);
        } else {
            value = buffer.readUnsignedInt(16) >> 3;
        }
        return value;
    }

    public static String parseS7String(ReadBuffer io, int stringLength, String encoding) {
        try {
            if ("UTF8".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                short maxLength = io.readUnsignedShort(8, UINT_OPT);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                short totalStringLength = io.readUnsignedShort(8, UINT_OPT);

                final byte[] byteArray = new byte[totalStringLength];
                for (int i = 0; (i < stringLength) && (io.getRemainingBits() >= 8); i++) {
                    final byte curByte = io.readSignedByte(8, SINT_OPT);
                    if (i < totalStringLength) {
                        byteArray[i] = curByte;
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for (; (i < stringLength) && (io.getRemainingBits() >= 8); i++) {
                            io.readSignedByte(8, SINT_OPT);
                        }
                        break;
                    }
                }
                return new String(byteArray, StandardCharsets.UTF_8);
            } else if ("UTF16BE".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                int maxLength = io.readUnsignedInt(16, UINT_OPT);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                int totalStringLength = io.readUnsignedInt(16, UINT_OPT);

                final byte[] byteArray = new byte[totalStringLength * 2];
                for (int i = 0; (i < stringLength) && (io.getRemainingBits() >= 16); i++) {
                    final short curShort = io.readSignedShort(16, SINT_OPT);
                    if (i < totalStringLength) {
                        byteArray[i * 2] = (byte) (curShort >>> 8);
                        byteArray[(i * 2) + 1] = (byte) (curShort & 0xFF);
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for (; (i < stringLength) && (io.getRemainingBits() >= 16); i++) {
                            io.readSignedShort(16, SINT_OPT);
                        }
                        break;
                    }
                }
                return new String(byteArray, StandardCharsets.UTF_16BE);
            } else {
                throw new PlcRuntimeException("Unsupported string encoding " + encoding);
            }
        } catch (BufferException e) {
            throw new PlcRuntimeException("Error parsing string", e);
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
        switch (encoding) {
            case "UTF8": {
                // In the case of STRING, the default and the max stringLength is 254.
                int maxStringLength = 0xFF & Math.min(stringLength, 254);
                int actStringLength = 0xFF & value.getString().length();
                actStringLength = Math.min(maxStringLength, actStringLength);

                byte[] chars = new byte[maxStringLength];
                byte[] actChars = value.getString().substring(0, actStringLength).getBytes(StandardCharsets.UTF_8);
                System.arraycopy(actChars, 0, chars, 0, actChars.length);
                try {
                    io.writeUnsignedInt(8, maxStringLength, UINT_OPT);
                    io.writeUnsignedInt(8, actStringLength, UINT_OPT);
                    io.writeBits(chars.length * 8, chars);
                } catch (BufferException ex) {
                    Logger.getLogger(StaticHelper.class.getName()).log(Level.SEVERE, null, ex);
                }
                break;
            }
            case "UTF16BE": {
                // In the case of WSTRING the default is also 254. However, the max stringLength is 16382.
                // As we've settled the default handling in S7ProtocolLogic, we'll only handle the max here.
                int maxStringLength = 0xFFFF & Math.min(stringLength, 16382);
                int actStringLength = 0xFFFF & value.getString().length();
                actStringLength = Math.min(maxStringLength, actStringLength);

                byte[] chars = new byte[maxStringLength * 2];
                byte[] actChars = value.getString().substring(0, actStringLength).getBytes(StandardCharsets.UTF_16BE);
                System.arraycopy(actChars, 0, chars, 0, actChars.length);
                try {
                    io.writeUnsignedInt(16, maxStringLength, UINT_OPT);
                    io.writeUnsignedInt(16, actStringLength, UINT_OPT);
                    io.writeBits(chars.length * 8, chars);
                } catch (BufferException ex) {
                    Logger.getLogger(StaticHelper.class.getName()).log(Level.SEVERE, null, ex);
                }
                break;
            }
            default:
                throw new PlcRuntimeException("Unsupported encoding: " + encoding);
        }
    }

    public static Long parseS5Time(ReadBuffer io) {
        try {
            short s5time = (short) io.readSignedInt(16);
            return s5TimeToDuration(s5time);
        } catch (BufferException e) {
            throw new RuntimeException(e);
        }
    }

    public static Long s5TimeToDuration(Short data) {
        short t = data;
        long tv = (short) (((t & 0x000F)) + ((t & 0x00F0) >> 4) * 10 + ((t & 0x0F00) >> 8) * 100);
        long tb = (short) (10 * Math.pow(10, ((t & 0xF000) >> 12)));
        long totalms = tv * tb;
        return (totalms <= 9990000) ? totalms : 9990000;
    }

    public static void serializeS5Time(final WriteBuffer io, PlcValue value) {
        final PlcTIME time = (PlcTIME) value;
        Short shortValue = durationToS5Time(time.getDuration());
        try {
            io.writeUnsignedInt(16, shortValue);
        } catch (BufferException e) {
            throw new RuntimeException(e);
        }
    }

    public static Short durationToS5Time(Duration duration) {
        short tv;
        short tb;
        short s5time = 0x0000;
        long totalms = duration.toMillis();

        if ((totalms >= 0) && (totalms <= 9990000)) {
            if (totalms <= 9990) {
                tb = 0x0000_0000; //10 ms
                tv = (short) (totalms / 10);
            } else if (totalms <= 99900) {
                tb = 0x0000_0001;// 100 ms
                tv = (short) (totalms / 100);
            } else if (totalms <= 999000) {
                tb = 0x0000_0002;//1000 ms
                tv = (short) (totalms / 1000);
            } else {
                tb = 0x0000_0003;//10000 ms
                tv = (short) (totalms / 10000);
            }

            short uni = (short) (tv % 10);
            short dec = (short) ((tv / 10) % 10);
            short cen = (short) ((tv / 100) % 10);

            return (short) (((tb) << 12) | (cen << 8) | (dec << 4) | (uni));
        }
        return s5time;
    }

    //TODO: apply only if not the last item
    public static int eventItemLength(final ReadBuffer buffer, int valueLength) {
        return ((valueLength % 2 == 0) || (buffer.getRemainingBits() < (valueLength + 1) * 8)) ? valueLength : valueLength + 1;
    }

    private static final LocalDate siemensEpoch = LocalDate.of(1990, 1, 1);
    private static final int daysBetweenUnixAndSiemensEpoch = (int) ChronoUnit.DAYS.between(LocalDate.EPOCH, siemensEpoch);

    public static Integer parseTiaDate(ReadBuffer io) {
        try {
            // Dates in Siemens PLCs are stored relative to "Siemens Epoch", which is 1990-01-01
            int daysSinceSiemensEpoch = io.readUnsignedInt(16);
            return daysSinceSiemensEpoch + daysBetweenUnixAndSiemensEpoch;
        } catch (BufferException e) {
            throw new RuntimeException(e);
        }
    }

    public static void serializeTiaDate(WriteBuffer io, PlcValue value) {
        final PlcDATE userDate = (PlcDATE) value;

        int daysSince1990 = userDate.getDaysSinceEpoch() - daysBetweenUnixAndSiemensEpoch;
        try {
            io.writeUnsignedInt(16, daysSince1990);
        } catch (BufferException e) {
            throw new RuntimeException(e);
        }
    }

    /**
     * Siemens numbers the DATE_AND_TIME day-of-week nibble 1 == Sunday .. 7 == Saturday; the mspec
     * spells that out for the DTL variant of the same field. {@code java.time.DayOfWeek} numbers the
     * same days 1 == Monday .. 7 == Sunday, which is what {@code PlcDATE_AND_TIME#getDayOfWeek}
     * returns and what plc4j shares with KNX DPT 19.001. Only the S7 wire format counts from Sunday,
     * so the rotation belongs here rather than in the shared value.
     */
    public static short parseSiemensDayOfWeek(ReadBuffer readBuffer) {
        try {
            short dayOfWeek = readBuffer.readUnsignedShort(4, WithOption.WithName("dayOfWeek"),
                WithOption.WithUnsignedIntegerEncoding("BCD"));
            if (dayOfWeek < 1 || dayOfWeek > 7) {
                throw new RuntimeException("day of week " + dayOfWeek
                    + " is outside the range [1, 7] the Siemens DATE_AND_TIME nibble can represent");
            }
            // Siemens Sunday is the ISO week's last day.
            return dayOfWeek == 1 ? (short) 7 : (short) (dayOfWeek - 1);
        } catch (BufferException e) {
            throw new RuntimeException("Error parsing dayOfWeek", e);
        }
    }

    /** Writes the day of week implied by the timestamp, numbered the way an S7 expects it. */
    public static void serializeSiemensDayOfWeek(WriteBuffer writeBuffer, PlcValue dateTime) {
        try {
            // DayOfWeek.getValue() is 1 == Monday .. 7 == Sunday; Siemens wants Sunday first.
            int iso = dateTime.getDateTime().getDayOfWeek().getValue();
            short siemens = (short) (iso == 7 ? 1 : iso + 1);
            writeBuffer.writeUnsignedShort(4, siemens, WithOption.WithName("dayOfWeek"),
                WithOption.WithUnsignedIntegerEncoding("BCD"));
        } catch (BufferException e) {
            throw new RuntimeException("Error serializing dayOfWeek", e);
        }
    }

    /**
     * The single BCD byte of a DATE_AND_TIME encodes 00-89 as 2000-2089 and 90-99 as 1990-1999, so
     * only years in [1990, 2089] are representable at all.
     */
    private static final int MIN_SIEMENS_YEAR = 1990;

    private static final int MAX_SIEMENS_YEAR = 2089;

    public static short parseSiemensYear(ReadBuffer readBuffer) {
        try {
            short year = readBuffer.readUnsignedShort(8, WithOption.WithName("year"), WithOption.WithUnsignedIntegerEncoding("BCD"));
            if (year < 90) {
                return (short) (2000 + year);
            } else {
                return (short) (1900 + year);
            }
        } catch (BufferException e) {
            throw new RuntimeException("Error parsing year", e);
        }
    }

    /**
     * The exact inverse of {@link #parseSiemensYear(ReadBuffer)}: 1990-1999 go out as 90-99 and
     * 2000-2089 as 00-89. Anything outside that window is rejected instead of silently wrapping -
     * 2090 used to be written as BCD 90 and read back as 1990, and 2000 took the 1900 branch and
     * asked {@code EncodingBCD} for a two digit encoding of 100, which threw an
     * {@link IllegalArgumentException} straight past the {@code BufferException} handler below.
     */
    public static void serializeSiemensYear(WriteBuffer writeBuffer, PlcValue dateTime) {
        try {
            int year = dateTime.getDateTime().getYear();
            if ((year < MIN_SIEMENS_YEAR) || (year > MAX_SIEMENS_YEAR)) {
                throw new RuntimeException("year " + year + " is outside the range ["
                    + MIN_SIEMENS_YEAR + ", " + MAX_SIEMENS_YEAR
                    + "] the Siemens DATE_AND_TIME year byte can represent");
            }
            if (year >= 2000) {
                writeBuffer.writeUnsignedShort(8, (short) (year - 2000), WithOption.WithName("year"), WithOption.WithUnsignedIntegerEncoding("BCD"));
            } else {
                writeBuffer.writeUnsignedShort(8, (short) (year - 1900), WithOption.WithName("year"), WithOption.WithUnsignedIntegerEncoding("BCD"));
            }
        } catch (BufferException e) {
            throw new RuntimeException("Error serializing year", e);
        }
    }

    public static boolean nextByteDoesNotMatch(ReadBuffer readBuffer, int referenceValue) {
        try {
            int startPos = readBuffer.getPositionInBits();
            int nextByte = readBuffer.readUnsignedShort(8);
            readBuffer.setPositionInBits(startPos);
            return nextByte != referenceValue;
        } catch (BufferException e) {
            return true;
        }
    }

    public static boolean nextByteMatches(ReadBuffer readBuffer, int referenceValue) {
        try {
            int startPos = readBuffer.getPositionInBits();
            int nextByte = readBuffer.readUnsignedShort(8);
            readBuffer.setPositionInBits(startPos);
            return nextByte == referenceValue;
        } catch (BufferException e) {
            return true;
        }
    }

    public static boolean nextWordMatches(ReadBuffer readBuffer, int referenceValue) {
        try {
            int startPos = readBuffer.getPositionInBits();
            int nextWord = readBuffer.readUnsignedInt(16);
            readBuffer.setPositionInBits(startPos);
            return nextWord == referenceValue;
        } catch (BufferException e) {
            return true;
        }
    }

}
