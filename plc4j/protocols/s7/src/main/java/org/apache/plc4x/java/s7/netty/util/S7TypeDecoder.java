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

import org.apache.plc4x.java.api.exceptions.PlcProtocolException;

import java.io.UnsupportedEncodingException;
import java.util.LinkedList;
import java.util.List;

public class S7TypeDecoder {

    private S7TypeDecoder() {
        // Utility class
    }

    @SuppressWarnings("unchecked")
    public static <T> List<T> decodeData(Class<T> datatype, byte[] s7Data) throws PlcProtocolException {

        List<Object> result = new LinkedList<>();
        int i = 0;
        final int length = s7Data.length;
        while (i < length) {
            if (datatype == Boolean.class) {
                result.add((s7Data[i] & 0x01) == 0x01);
                i += 1;
            } else if (datatype == Byte.class) {
                result.add(s7Data[i]);
                i += 1;
            } else if (datatype == Short.class) {
                result.add((short) (((s7Data[i] & 0xff) << 8) | (s7Data[i + 1] & 0xff)));
                i += 2;
            } else if (datatype == Integer.class) {
                result.add(((s7Data[i] & 0xff) << 24) | ((s7Data[i + 1] & 0xff) << 16) |
                    ((s7Data[i + 2] & 0xff) << 8) | (s7Data[i + 3] & 0xff));
                i += 4;
            } else if (datatype == Float.class) {
                // Description of the Real number format:
                // https://www.sps-lehrgang.de/zahlenformate-step7/#c144
                // https://de.wikipedia.org/wiki/IEEE_754
                int intValue = ((s7Data[i] & 0xff) << 24)
                    | ((s7Data[i + 1] & 0xff) << 16)
                    | ((s7Data[i + 2] & 0xff) << 8)
                    | (s7Data[i + 3] & 0xff);
                result.add(Float.intBitsToFloat(intValue));
                i += 4;
            } else if (datatype == Double.class) {
                // Description of the Real number format:
                // https://www.sps-lehrgang.de/zahlenformate-step7/#c144
                // https://de.wikipedia.org/wiki/IEEE_754
                long longValue = (((long) (s7Data[i] & 0xff)) << 56)
                    | (((long) (s7Data[i] & 0xff)) << 48)
                    | (((long) (s7Data[i + 1] & 0xff)) << 40)
                    | (((long) (s7Data[i + 2] & 0xff)) << 32)

                    | (((long) (s7Data[i + 3] & 0xff)) << 24)
                    | (((long) (s7Data[i + 4] & 0xff)) << 16)
                    | (((long) (s7Data[i + 5] & 0xff)) << 8)
                    | (((long) s7Data[i + 6] & 0xff));
                result.add(Double.longBitsToDouble(longValue));
                i += 8;
            } else if (datatype == String.class) {
                // Every string value had a prefix of two bytes for which I have no idea, what the meaning is.
                // This code assumes the string values doesn't contain UTF-8 values with a code of 0x00 as it
                // uses this as termination char.
                try {
                    int j = 0;
                    for (; j < s7Data.length; j++) {
                        if (s7Data[j] == 0) {
                            break;
                        }
                    }
                    result.add(new String(s7Data, 2, j - 2, "UTF-8"));
                } catch (UnsupportedEncodingException e) {
                    throw new PlcProtocolException("Error decoding String value", e);
                }
                i += s7Data.length;
            } else {
                throw new PlcProtocolException("Unsupported data type " + datatype.getSimpleName());
            }
        }
        return (List<T>) result;
    }
}
