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
package org.apache.plc4x.java.simulated.utils;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.Charset;


public class StaticHelper {

    public static String parsePascalString(ReadBuffer io, String encoding) {
        try {
            // This is the maximum number of bytes a string can be long.
            short stringLength = io.readUnsignedShort(8);
            // Read the full size of the string.
            String str = io.readString(stringLength * 8, (String) encoding);
            // Cut off the parts that don't belong to it.
            return str;
        } catch (ParseException e) {
            return null;
        }
    }

    public static void serializePascalString(WriteBuffer io, PlcValue value, String encoding) throws ParseException {
        final byte[] bytes = value.getString().getBytes(Charset.forName(encoding));
        try {
            if (bytes.length < 256) {
                io.writeByte((byte) bytes.length);
                for (byte aByte : bytes) {
                    io.writeByte(aByte);
                }
            } else {
                throw new ParseException("Error writing string, string > 255 bytes");
            }
        } catch (ParseException e) {
            throw new ParseException("Error writing string", e);
        }
    }

}
