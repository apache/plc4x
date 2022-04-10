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
package org.apache.plc4x.java.knxnetip.utils;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.generation.*;

public class KnxHelper {

    public static float bytesToF16(ReadBuffer io) {
        try {
            boolean negative = io.readBit();
            long exponent = io.readUnsignedLong(4);
            long mantissa = io.readUnsignedLong(11);
            return (float) ((negative ? -1 : 1) * (0.01 * mantissa) * Math.pow(2, exponent));
        } catch(ParseException e) {
            throw new PlcRuntimeException("Error parsing F16 value", e);
        }
    }

    public static void f16toBytes(WriteBuffer io, Object param) {
        if(!(param instanceof Float)) {
            throw new PlcRuntimeException("Invalid datatype");
        }
        try {
            float value = (float) param;
            boolean negative = value < 0;
            final int exponent = Math.getExponent(value);
            final double mantissa = value / Math.pow(2, exponent);
            //io.writeBit(negative);
            //io.writeInt(4, exponent);
            //io.writeDouble(11, Math.getExponent(mantissa), mantissa); //Don't think this works BH

            String mantissaString = Double.toString(mantissa);
            int mantissaInteger = Integer.parseInt(mantissaString.substring(mantissaString.indexOf('.')));
            io.writeBit(negative);
            io.writeUnsignedInt(4,exponent + 15);
            io.writeUnsignedInt(11, mantissaInteger);

        } catch(SerializationException e) {
            throw new PlcRuntimeException("Error serializing F16 value", e);
        }
    }


    public static void main(String[] args) throws Exception {
        final byte[] bytes = Hex.decodeHex("0C65");
        ReadBuffer io = new ReadBufferByteBased(bytes);
        final double v = bytesToF16(io);
        System.out.println(v);
    }

}
