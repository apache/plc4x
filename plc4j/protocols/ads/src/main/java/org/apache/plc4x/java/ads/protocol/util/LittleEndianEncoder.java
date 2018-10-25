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

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.model.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.time.*;
import java.time.temporal.ChronoUnit;
import java.util.Arrays;
import java.util.stream.Stream;

// TODO: we might user ByteBuffer.wrap(buffer).order(ByteOrder.LITTLE_ENDIAN).putInt(port).asArray() etc
public class LittleEndianEncoder {

    private LittleEndianEncoder() {
        // Utility class
    }

    public static byte[] encodeData(AdsDataType adsDataType, Object... values) throws PlcProtocolException {
        if (values.length == 0) {
            return new byte[]{};
        }
        Class<?> valueType = values[0].getClass();
        Stream<byte[]> result;
        if (valueType == Boolean.class) {
            result = encodeBoolean(adsDataType, Arrays.stream(values).map(Boolean.class::cast));
        } else if (valueType == Byte.class) {
            result = encodeByte(adsDataType, Arrays.stream(values).map(Byte.class::cast));
        } else if (valueType == Short.class) {
            result = encodeShort(adsDataType, Arrays.stream(values).map(Short.class::cast));
        } else if (valueType == Integer.class) {
            result = encodeInteger(adsDataType, Arrays.stream(values).map(Integer.class::cast));
        } else if (valueType == Long.class) {
            result = encodeLong(adsDataType, Arrays.stream(values).map(Long.class::cast));
        } else if (valueType == BigInteger.class) {
            result = encodeBigInteger(adsDataType, Arrays.stream(values).map(BigInteger.class::cast));
        } else if (valueType == LocalTime.class) {
            result = encodeLocalTime(adsDataType, Arrays.stream(values).map(LocalTime.class::cast));
        } else if (valueType == LocalDate.class) {
            result = encodeLocalDate(adsDataType, Arrays.stream(values).map(LocalDate.class::cast));
        } else if (valueType == LocalDateTime.class) {
            result = encodeLocalDateTime(adsDataType, Arrays.stream(values).map(LocalDateTime.class::cast));
        } else if (valueType == Float.class) {
            result = encodeFloat(adsDataType, Arrays.stream(values).map(Float.class::cast));
        } else if (valueType == Double.class) {
            result = encodeDouble(adsDataType, Arrays.stream(values).map(Double.class::cast));
        } else if (valueType == String.class) {
            result = encodeString(adsDataType, Arrays.stream(values).map(String.class::cast));
        } else if (valueType == byte[].class) {
            result = encodeByteArray(adsDataType, Arrays.stream(values).map(byte[].class::cast));
        } else if (valueType == Byte[].class) {
            result = encodeBigByteArray(adsDataType, Arrays.stream(values).map(Byte[].class::cast));
        } else {
            throw new PlcUnsupportedDataTypeException(valueType);
        }

        // TODO: maybe we can replace this by a smarter flatmap
        try {
            return result.collect(
                ByteArrayOutputStream::new,
                (bos, byteValue) -> {
                    try {
                        bos.write(byteValue);
                    } catch (IOException e) {
                        throw new PlcRuntimeException(e);
                    }
                },
                (a, b) -> {
                }).toByteArray();
        } catch (PlcRuntimeException e) {
            throw new PlcProtocolException("Error encoding data", e);
        }
    }

    private static Stream<byte[]> encodeString(AdsDataType adsDataType, Stream<String> stringStream) {
        // TODO: add boundchecks and add optional extension
        // TODO: what do we do with utf-8 values with 2 bytes? what is the charset here?
        return stringStream
            .map(s -> s.getBytes(Charset.defaultCharset()))
            // TODO: this 0 termination is from s7 but might be completly wrong in ads. Guess its a terminator
            .map(bytes -> ArrayUtils.add(bytes, (byte) 0x0));
    }

    private static Stream<byte[]> encodeByteArray(AdsDataType adsDataType, Stream<byte[]> byteArrayStream) {
        // TODO: add boundchecks and add optional extension
        return byteArrayStream;
    }

    private static Stream<byte[]> encodeBigByteArray(AdsDataType adsDataType, Stream<Byte[]> byteArrayStream) {
        // TODO: add boundchecks and add optional extension
        return byteArrayStream.map(ArrayUtils::toPrimitive);
    }

    private static Stream<byte[]> encodeFloat(AdsDataType adsDataType, Stream<Float> floatStream) {
        return floatStream
            .peek(value -> checkBound(adsDataType, value))
            // TODO: check how ads expects this data
            .map(Float::floatToIntBits)
            .map(intValue -> new byte[]{
                (byte) (intValue & 0x000000ff),
                (byte) ((intValue & 0x0000ff00) >> 8),
                (byte) ((intValue & 0x00ff0000) >> 16),
                (byte) ((intValue & 0xff000000) >> 24),
            });
    }

    private static Stream<byte[]> encodeDouble(AdsDataType adsDataType, Stream<Double> doubleStream) {
        return doubleStream
            .peek(value -> checkBound(adsDataType, value))
            // TODO: check how ads expects this data
            .map(Double::doubleToLongBits)
            .map(longValue -> new byte[]{
                (byte) (longValue & 0x00000000_000000ffL),
                (byte) ((longValue & 0x00000000_0000ff00L) >> 8),
                (byte) ((longValue & 0x00000000_00ff0000L) >> 16),
                (byte) ((longValue & 0x00000000_ff000000L) >> 24),
                (byte) ((longValue & 0x000000ff_00000000L) >> 32),
                (byte) ((longValue & 0x0000ff00_00000000L) >> 40),
                (byte) ((longValue & 0x00ff0000_00000000L) >> 48),
                (byte) ((longValue & 0xff000000_00000000L) >> 56),
            });
    }

    private static Stream<byte[]> encodeInteger(AdsDataType adsDataType, Stream<Integer> integerStream) {
        return integerStream
            .peek(value -> checkBound(adsDataType, value))
            .map(intValue -> new byte[]{
                (byte) (intValue & 0x000000ff),
                (byte) ((intValue & 0x0000ff00) >> 8),
                (byte) ((intValue & 0x00ff0000) >> 16),
                (byte) ((intValue & 0xff000000) >> 24),
            });
    }

    private static Stream<byte[]> encodeLong(AdsDataType adsDataType, Stream<Long> integerStream) {
        return integerStream
            .peek(value -> checkBound(adsDataType, value))
            .map(longValue -> new byte[]{
                (byte) (longValue & 0x00000000_000000ffL),
                (byte) ((longValue & 0x00000000_0000ff00L) >> 8),
                (byte) ((longValue & 0x00000000_00ff0000L) >> 16),
                (byte) ((longValue & 0x00000000_ff000000L) >> 24),
                (byte) ((longValue & 0x000000ff_00000000L) >> 32),
                (byte) ((longValue & 0x0000ff00_00000000L) >> 40),
                (byte) ((longValue & 0x00ff0000_00000000L) >> 48),
                (byte) ((longValue & 0xff000000_00000000L) >> 56),
            });
    }

    private static Stream<byte[]> encodeBigInteger(AdsDataType adsDataType, Stream<BigInteger> bigIntegerStream) {
        // TODO: add boundchecks and add optional extension
        return bigIntegerStream
            .map(bigIntValue -> {
                byte[] bytes = bigIntValue.toByteArray();
                if (bytes[0] == 0x0) {
                    byte[] subArray = ArrayUtils.subarray(bytes, 1, bytes.length);
                    ArrayUtils.reverse(subArray);
                    return subArray;
                } else {
                    ArrayUtils.reverse(bytes);
                    return bytes;
                }
            });
    }

    private static Stream<byte[]> encodeLocalTime(AdsDataType adsDataType, Stream<LocalTime> localTimeStream) {
        return localTimeStream
            .map(localTime -> ChronoUnit.MILLIS.between(LocalTime.of(0, 0), localTime))
            .peek(value -> checkBound(adsDataType, value))
            .map(Long::intValue)
            .map(time -> new byte[]{
                (byte) (time & 0x00000000_000000ffL),
                (byte) ((time & 0x00000000_0000ff00L) >> 8),
                (byte) ((time & 0x00000000_00ff0000L) >> 16),
                (byte) ((time & 0x00000000_ff000000L) >> 24),
            });
    }

    private static Stream<byte[]> encodeLocalDate(AdsDataType adsDataType, Stream<LocalDate> localDateStream) {
        return localDateStream
            // TODO: fixme: which offset should we use?
            .map(localDate -> localDate.atTime(0, 0).toInstant(ZoneOffset.UTC))
            .map(Instant::getEpochSecond)
            .peek(value -> checkBound(adsDataType, value))
            .map(Long::intValue)
            .map(time -> new byte[]{
                (byte) (time & 0x00000000_000000ffL),
                (byte) ((time & 0x00000000_0000ff00L) >> 8),
                (byte) ((time & 0x00000000_00ff0000L) >> 16),
                (byte) ((time & 0x00000000_ff000000L) >> 24),
            });
    }

    private static Stream<byte[]> encodeLocalDateTime(AdsDataType adsDataType, Stream<LocalDateTime> localDateTimeStream) {
        return localDateTimeStream
            // TODO: fixme: which offset should we use?
            .map(localDateTime -> localDateTime.toInstant(ZoneOffset.UTC))
            .map(Instant::getEpochSecond)
            .peek(value -> checkBound(adsDataType, value))
            .map(time -> new byte[]{
                (byte) (time & 0x00000000_000000ffL),
                (byte) ((time & 0x00000000_0000ff00L) >> 8),
                (byte) ((time & 0x00000000_00ff0000L) >> 16),
                (byte) ((time & 0x00000000_ff000000L) >> 24),

                (byte) ((time & 0x000000ff_00000000L) >> 32),
                (byte) ((time & 0x0000ff00_00000000L) >> 40),
                (byte) ((time & 0x00ff0000_00000000L) >> 48),
                (byte) ((time & 0xff000000_00000000L) >> 56),
            });
    }


    private static Stream<byte[]> encodeShort(AdsDataType adsDataType, Stream<Short> shortStream) {
        return shortStream
            .peek(value -> checkBound(adsDataType, value))
            .map(shortValue -> new byte[]{
                (byte) (shortValue & 0x00ff),
                (byte) ((shortValue & 0xff00) >> 8),
            });
    }

    private static Stream<byte[]> encodeByte(AdsDataType adsDataType, Stream<Byte> byteStream) {
        return byteStream
            .peek(value -> checkBound(adsDataType, value))
            .map(aByte -> new byte[]{aByte});
    }

    private static Stream<byte[]> encodeBoolean(AdsDataType adsDataType, Stream<Boolean> booleanStream) {
        // TODO: add boundchecks and add optional extension
        return booleanStream
            .map(booleanValue -> new byte[]{booleanValue ? (byte) 0x01 : (byte) 0x00});
    }

    private static void checkBound(AdsDataType adsDataType, double other) {
        if (!adsDataType.withinBounds(other)) {
            throw new PlcRuntimeException(other + " not within bounds of " + adsDataType);
        }
    }
}
