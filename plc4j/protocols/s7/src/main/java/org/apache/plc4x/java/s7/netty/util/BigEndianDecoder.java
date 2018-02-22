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

import java.util.LinkedList;
import java.util.List;

public class BigEndianDecoder {

    private BigEndianDecoder() {
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
                int intValue = ((s7Data[i] & 0xff) << 24) | ((s7Data[i + 1] & 0xff) << 16) |
                    ((s7Data[i + 2] & 0xff) << 8) | (s7Data[i + 3] & 0xff);
                result.add(Float.intBitsToFloat(intValue));
                i += 4;
            } else if (datatype == String.class) {
                StringBuilder builder = new StringBuilder();
                while (s7Data[i] != (byte) 0x0 && i < length) {
                    builder.append((char) s7Data[i]);
                    i++;
                }
                i++; // skip terminating character
                result.add(builder.toString());
            } else {
                throw new PlcProtocolException("Unsupported datatype " + datatype.getSimpleName());
            }
        }
        return (List<T>) result;
    }
}
