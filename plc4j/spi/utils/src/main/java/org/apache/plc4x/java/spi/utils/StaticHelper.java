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

package org.apache.plc4x.java.spi.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.Message;

import java.math.BigInteger;
import java.util.Collection;
import java.util.List;

public class StaticHelper {

    public static int ARRAY_SIZE_IN_BYTES(Object obj) {
        if (obj instanceof List<?> list) {
            int numBytes = 0;
            for (Object element : list) {
                if (!(element instanceof Message)) {
                    throw new RuntimeException(
                        "Array elements for array size in bytes must implement Message interface");
                }
                numBytes += ((Message) element).getLengthInBytes();
            }
            return numBytes;
        }
        if (obj.getClass().isArray() && !obj.getClass().getComponentType().isPrimitive()) {
            Object[] arr = (Object[]) obj;
            int numBytes = 0;
            for (Object element : arr) {
                if (!(element instanceof Message)) {
                    throw new RuntimeException(
                        "Array elements for array size in bytes must implement Message interface");
                }
                numBytes += ((Message) element).getLengthInBytes();
            }
            return numBytes;
        }
        throw new RuntimeException("Unable to calculate array size in bytes for type " + obj.getClass().getName());
    }

    public static int COUNT(Object obj) {
        if (obj == null) {
            return 0;
        }
        if (obj.getClass().isArray()) {
            if (obj.getClass().getComponentType() != null && obj.getClass().getComponentType().isPrimitive()) {
                if (obj.getClass().getComponentType() == boolean.class) {
                    boolean[] arr = (boolean[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == byte.class) {
                    byte[] arr = (byte[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == short.class) {
                    short[] arr = (short[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == int.class) {
                    int[] arr = (int[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == long.class) {
                    long[] arr = (long[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == float.class) {
                    float[] arr = (float[]) obj;
                    return arr.length;
                }
                if (obj.getClass().getComponentType() == double.class) {
                    double[] arr = (double[]) obj;
                    return arr.length;
                }
            } else {
                Object[] arr = (Object[]) obj;
                return arr.length;
            }
        } else if (obj instanceof Collection<?> col) {
            return col.size();
        }
        throw new PlcRuntimeException("Unable to count object of type " + obj.getClass().getName());
    }

    public static int STR_LEN(Object str) {
        if (str == null) {
            return 0;
        }
        if (str instanceof PlcValue plcValue) {
            return plcValue.getString().length();
        }
        return str.toString().length();
    }

    public static <T> T CAST(Object obj, Class<T> clazz) {
        try {
            return clazz.cast(obj);
        } catch (ClassCastException e) {
            throw new PlcRuntimeException("Unable to cast object of type " + obj.getClass().getName() + " to " + clazz.getName());
        }
    }

    public static int CEIL(double value) {
        return (int) Math.ceil(value);
    }

    public static int PADCOUNT(Object obj, boolean hasNext) {
        return hasNext ? COUNT(obj) : 0;
    }

    public static int GET_VARDUINT_LENGTH_IN_BITS(long value) {
        int curFieldLengthInBits = 0;
        long temp = value;
        do {
            curFieldLengthInBits += 8;
            temp >>>= 7;
        } while (temp != 0);
        return curFieldLengthInBits;
    }

    public static int GET_VARDINT_LENGTH_IN_BITS(long value) {
        int curFieldLengthInBits = 8;
        boolean positive = value >= 0;
        long tmpValue = value;
        while (tmpValue >> 6 != (positive ? 0 : -1)) {
            curFieldLengthInBits += 8;
            tmpValue >>= 7;
        }
        return curFieldLengthInBits;
    }

    public static int GET_VARDUINT_LENGTH_IN_BITS(BigInteger value) {
        int curFieldLengthInBits = 0;
        long temp = value.longValue();
        do {
            curFieldLengthInBits += 8;
            temp >>>= 7;
        } while (temp != 0);
        return curFieldLengthInBits;
    }

    public static String ENCODE_HEX(byte[] data) {
        if (data == null || data.length == 0) {
            return "";
        }
        final char[] HEX = "0123456789ABCDEF".toCharArray();
        char[] out = new char[data.length * 2];
        int i = 0;
        for (byte b : data) {
            int v = b & 0xFF;
            out[i++] = HEX[v >>> 4];
            out[i++] = HEX[v & 0x0F];
        }
        return new String(out);
    }

    public static byte[] DECODE_HEX(String hex) {
        if (hex == null || hex.isEmpty()) {
            return new byte[0];
        }
        int len = hex.length();
        if ((len & 1) != 0) {
            throw new IllegalArgumentException("Hex string must have even length");
        }
        byte[] out = new byte[len / 2];
        for (int i = 0, j = 0; i < len; i += 2, j++) {
            int hi = Character.digit(hex.charAt(i), 16);
            int lo = Character.digit(hex.charAt(i + 1), 16);
            if (hi == -1 || lo == -1) {
                throw new IllegalArgumentException("Invalid hex character at position " + i);
            }
            out[j] = (byte) ((hi << 4) + lo);
        }
        return out;
    }

    public static String CAPITALIZE(String str) {
        return str.substring(0, 1).toUpperCase() + str.substring(1);
    }

    /**
     * This helper is only required for var-length encoded integers.
     * @param value the value to calculate the length for.
     * @return number of bits required to encode the value.
     */
    public static int GET_VAR_LENGTH_SINT_IN_BITS(long value) {
        int numBytes = 0;
        long temp = value;

        while (true) {
            // Extract the low 7 bits
            int low7 = (int) (temp & 0x7F);
            boolean signBitSet = (low7 & 0x40) != 0; // bit 6

            // Arithmetic shift right by 7
            temp >>= 7;
            numBytes++;

            boolean done =
                (temp == 0 && !signBitSet) ||
                    (temp == -1 && signBitSet);

            if (done) break;
        }

        return numBytes * 8;
    }

    /**
     * This helper is only required for var-length encoded integers.
     * @param value the value to calculate the length for.
     * @return number of bits required to encode the value.
     */
    public static int GET_VAR_LENGTH_UINT_IN_BITS(long value) {
        if (value < 0) {
            throw new IllegalArgumentException("Unsigned value must be >= 0");
        }
        int numBytes = 0;
        long temp = value;

        do {
            numBytes++;
            temp >>>= 7; // unsigned shift
        } while (temp != 0);

        return numBytes * 8;
    }

    /**
     * This helper is only required for var-length encoded integers.
     * @param value the value to calculate the length for.
     * @return number of bits required to encode the value.
     */
    public static int GET_VAR_LENGTH_SINT_IN_BITS(BigInteger value) {
        int numBytes = 0;
        BigInteger temp = value;

        while (true) {
            // Take the low 7 bits of the current chunk
            int low7 = temp.and(BigInteger.valueOf(0x7F)).intValue();
            boolean signBitSet = (low7 & 0x40) != 0; // bit 6

            // Arithmetic shift by 7 (BigInteger.shiftRight is arithmetic)
            temp = temp.shiftRight(7);

            numBytes++;

            boolean done =
                (temp.equals(BigInteger.ZERO) && !signBitSet) ||
                    (temp.equals(BigInteger.valueOf(-1)) && signBitSet);

            if (done) break;
        }
        return numBytes * 8;
    }

    /**
     * This helper is only required for var-length encoded integers.
     * @param value the value to calculate the length for.
     * @return number of bits required to encode the value.
     */
    public static int GET_VAR_LENGTH_UINT_IN_BITS(BigInteger value) {
        if (value.signum() < 0) {
            throw new IllegalArgumentException("Unsigned value must be >= 0");
        }
        int numBytes = 0;
        BigInteger temp = value;

        do {
            numBytes++;
            temp = temp.shiftRight(7);
        } while (temp.signum() != 0);

        return numBytes * 8;
    }

}
