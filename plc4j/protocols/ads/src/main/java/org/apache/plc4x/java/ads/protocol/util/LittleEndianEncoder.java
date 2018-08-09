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
import org.apache.plc4x.java.ads.api.commands.types.TimeStamp;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Date;
import java.util.stream.Stream;

// TODO: we might user ByteBuffer.wrap(buffer).order(ByteOrder.LITTLE_ENDIAN).putInt(port).asArray() etc
public class LittleEndianEncoder {

    private LittleEndianEncoder() {
        // Utility class
    }

    // TODO: add bound checking
    public static byte[] encodeData(Class<?> valueType, Object... values) throws PlcProtocolException {
        if (values.length == 0) {
            return new byte[]{};
        }
        Stream<byte[]> result;
        if (valueType == Boolean.class) {
            result = encodeBoolean(Arrays.stream(values).map(Boolean.class::cast));
        } else if (valueType == Byte.class) {
            result = encodeByte(Arrays.stream(values).map(Byte.class::cast));
        } else if (valueType == Short.class) {
            result = encodeShort(Arrays.stream(values).map(Short.class::cast));
        } else if (valueType == Integer.class) {
            result = encodeInteger(Arrays.stream(values).map(Integer.class::cast));
        } else if (valueType == BigInteger.class) {
            result = encodeBigInteger(Arrays.stream(values).map(BigInteger.class::cast));
        } else if (valueType == Calendar.class || Calendar.class.isAssignableFrom(valueType)) {
            result = encodeCalendar(Arrays.stream(values).map(Calendar.class::cast));
        } else if (valueType == Float.class) {
            result = encodeFloat(Arrays.stream(values).map(Float.class::cast));
        } else if (valueType == Double.class) {
            result = encodeDouble(Arrays.stream(values).map(Double.class::cast));
        } else if (valueType == String.class) {
            result = encodeString(Arrays.stream(values).map(String.class::cast));
        } else if (valueType == byte[].class) {
            result = encodeByteArray(Arrays.stream(values).map(byte[].class::cast));
        } else if (valueType == Byte[].class) {
            result = encodeBigByteArray(Arrays.stream(values).map(Byte[].class::cast));
        } else {
            throw new PlcProtocolException("Unsupported datatype " + valueType);
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

    private static Stream<byte[]> encodeString(Stream<String> stringStream) {
        // TODO: what do we do with utf-8 values with 2 bytes? what is the charset here?
        return stringStream
            .map(s -> s.getBytes(Charset.defaultCharset()))
            // TODO: this 0 termination is from s7 but might be completly wrong in ads. Guess its a terminator
            .map(bytes -> ArrayUtils.add(bytes, (byte) 0x0));
    }

    private static Stream<byte[]> encodeByteArray(Stream<byte[]> byteArrayStream) {
        return byteArrayStream;
    }

    private static Stream<byte[]> encodeBigByteArray(Stream<Byte[]> byteArrayStream) {
        return byteArrayStream.map(ArrayUtils::toPrimitive);
    }

    private static Stream<byte[]> encodeFloat(Stream<Float> floatStream) {
        return floatStream
            // TODO: check how ads expects this data
            .map(Float::floatToIntBits)
            .map(intValue -> new byte[]{
                (byte) (intValue & 0x000000ff),
                (byte) ((intValue & 0x0000ff00) >> 8),
                (byte) ((intValue & 0x00ff0000) >> 16),
                (byte) ((intValue & 0xff000000) >> 24),
            });
    }

    private static Stream<byte[]> encodeDouble(Stream<Double> doubleStream) {
        return doubleStream
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

    private static Stream<byte[]> encodeInteger(Stream<Integer> integerStream) {
        return integerStream
            .map(intValue -> new byte[]{
                (byte) (intValue & 0x000000ff),
                (byte) ((intValue & 0x0000ff00) >> 8),
                (byte) ((intValue & 0x00ff0000) >> 16),
                (byte) ((intValue & 0xff000000) >> 24),
            });
    }

    private static Stream<byte[]> encodeBigInteger(Stream<BigInteger> bigIntegerStream) {
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

    private static Stream<byte[]> encodeCalendar(Stream<Calendar> calendarStream) {
        return calendarStream
            .map(Calendar.class::cast)
            .map(Calendar::getTime)
            .map(Date::getTime)
            .map(BigInteger::valueOf)
            .map(TimeStamp::javaToWinTime)
            .map(BigInteger::longValue)
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


    private static Stream<byte[]> encodeShort(Stream<Short> shortStream) {
        return shortStream
            .map(shortValue -> new byte[]{
                (byte) (shortValue & 0x00ff),
                (byte) ((shortValue & 0xff00) >> 8),
            });
    }

    private static Stream<byte[]> encodeByte(Stream<Byte> byteStream) {
        return byteStream
            .map(aByte -> new byte[]{aByte});
    }

    private static Stream<byte[]> encodeBoolean(Stream<Boolean> booleanStream) {
        return booleanStream
            .map(booleanValue -> new byte[]{booleanValue ? (byte) 0x01 : (byte) 0x00});
    }
}
