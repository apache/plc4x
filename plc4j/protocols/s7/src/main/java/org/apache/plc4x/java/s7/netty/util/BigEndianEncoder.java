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

import java.util.Calendar;

public class BigEndianEncoder {

    private BigEndianEncoder() {
        // Utility class
    }

    public static byte[] encodeData(Object[] values) {
        final int length = values.length;
        if (length == 0) {
            return new byte[]{};
        }
        byte[] result = null;
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
            result = null;
        } else if (valueType == Float.class) {
            result = encodeFloat(values, length);
        } else if (valueType == String.class) {
            result = encodeString(values, length);
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
