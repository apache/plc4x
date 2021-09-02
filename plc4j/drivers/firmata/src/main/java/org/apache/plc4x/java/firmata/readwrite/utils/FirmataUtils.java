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
package org.apache.plc4x.java.firmata.readwrite.utils;

import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class FirmataUtils {

    public static boolean isSysexEnd(ReadBuffer io) {
        return ((ReadBufferByteBased) io).getBytes(io.getPos(), io.getPos() + 1)[0] == (byte) 0xF7;
    }

    public static byte parseSysexString(ReadBuffer io) {
        try {
            byte b = io.readByte();
            // Skip the empty byte.
            io.readByte();
            return b;
        } catch (ParseException e) {
            return 0;
        }
    }

    public static void serializeSysexString(WriteBuffer io, byte data) {
        try {
            io.writeByte(data);
            io.writeByte((byte) 0x00);
        } catch (ParseException e) {
        }
    }

    public static int lengthSysexString(byte[] data) {
        return data.length * 2;
    }

}
