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
package org.apache.plc4x.java.s7.netty.util;

import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;

import java.util.Calendar;

public class S7TypeEncoder {

    private S7TypeEncoder() {
        // Utility class
    }

    public static byte[] encodeData(Object[] values) throws PlcProtocolException {
        final int length = values.length;
        if (length == 0) {
            return new byte[]{};
        }
        final byte[] result;
        Class valueType = values[0].getClass();
        if (valueType == Boolean.class) {
            result = encodeBoolean(values, length);
        } else if (valueType == Byte.class) {
            result = encodeByte(values, length);
        } else if (valueType == Short.class) {
            result = encodeShort(values, length);
        } else if (valueType == Integer.class) {
            result = encodeInteger(values, length);
        } else if (valueType == Calendar.class) {
            // TODO: Decide what to do here ...
            throw new PlcNotImplementedException("calender not yet implemented in s7");
        } else if (valueType == Float.class) {
            result = encodeFloat(values, length);
        } else if (valueType == Double.class) {
            result = encodeDouble(values, length);
        } else if (valueType == String.class) {
            result = encodeString(values, length);
        } else {
            throw new PlcProtocolException("Unsupported data type " + valueType);
        }
        return result;
    }

    public static byte[] encodeString(Object[] values, int length) {
        byte[] result;
        int size = 0;
        for (Object value : values) {
            size = size + ((String) value).length();
        }
        result = new byte[size + length];
        int j = 0;
        for (Object value : values) {
            String str = (String) value;
            for (int i = 0; i < str.length(); i++) {
                result[j++] = (byte) str.charAt(i);
            }
            result[j++] = (byte) 0x0;
        }
        return result;
    }

    public static byte[] encodeFloat(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 4];
        for (int i = 0; i < length; i++) {
            float floatValue = (float) values[i];
            int intValue = Float.floatToIntBits(floatValue);
            result[i * 4] = (byte) ((intValue & 0xff000000) >> 24);
            result[(i * 4) + 1] = (byte) ((intValue & 0x00ff0000) >> 16);
            result[(i * 4) + 2] = (byte) ((intValue & 0x0000ff00) >> 8);
            result[(i * 4) + 3] = (byte) (intValue & 0xff);
        }
        return result;
    }

    public static byte[] encodeDouble(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 8];
        for (int i = 0; i < length; i++) {
            double doubleValue = (double) values[i];
            long longValue = Double.doubleToLongBits(doubleValue);
            result[(i * 8)] = (byte) ((longValue & 0xFF000000_00000000L) >> 56);
            result[(i * 8) + 1] = (byte) ((longValue & 0x00FF0000_00000000L) >> 48);
            result[(i * 8) + 2] = (byte) ((longValue & 0x0000FF00_00000000L) >> 40);
            result[(i * 8) + 3] = (byte) ((longValue & 0x000000FF_00000000L) >> 32);
            result[(i * 8) + 4] = (byte) ((longValue & 0x00000000_FF000000L) >> 24);
            result[(i * 8) + 5] = (byte) ((longValue & 0x00000000_00FF0000L) >> 16);
            result[(i * 8) + 6] = (byte) ((longValue & 0x00000000_0000FF00L) >> 8);
            result[(i * 8) + 7] = (byte) (longValue & 0x00000000_000000FFL);
        }
        return result;
    }

    public static byte[] encodeInteger(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 4];
        for (int i = 0; i < length; i++) {
            int intValue = (int) values[i];
            result[i * 4] = (byte) ((intValue & 0xff000000) >> 24);
            result[(i * 4) + 1] = (byte) ((intValue & 0x00ff0000) >> 16);
            result[(i * 4) + 2] = (byte) ((intValue & 0x0000ff00) >> 8);
            result[(i * 4) + 3] = (byte) (intValue & 0xff);
        }
        return result;
    }

    public static byte[] encodeShort(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 2];
        for (int i = 0; i < length; i++) {
            short intValue = (short) values[i];
            result[i * 2] = (byte) ((intValue & 0xff00) >> 8);
            result[(i * 2) + 1] = (byte) (intValue & 0xff);
        }
        return result;
    }

    public static byte[] encodeByte(Object[] values, int length) {
        byte[] result;
        result = new byte[length];
        for (int i = 0; i < length; i++) {
            result[i] = (byte) values[i];
        }
        return result;
    }

    public static byte[] encodeBoolean(Object[] values, int length) {
        byte[] result;// TODO: Check if this is true and the result is not Math.ceil(values.lenght / 8)
        result = new byte[length];
        for (int i = 0; i < length; i++) {
            result[i] = (byte) (((Boolean) values[i]) ? 0x01 : 0x00);
        }
        return result;
    }
}
