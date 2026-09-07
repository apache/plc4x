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

package org.apache.plc4x.java.ads.readwrite.utils;

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

public class StaticHelper {

    public static String parseZeroTerminatedString(ReadBuffer io, int stringValueLength) throws BufferException {
        byte[] bytes = io.readBits(stringValueLength * 8, WithOption.WithName("stringValueLength"));
        String stringValue = new String(bytes, 0, bytes.length);
        // Read the string
        //String stringValue = io.readString(stringValueLength * 8, WithOption.WithName("stringValue"), WithOption.WithEncoding("ASCII"));
        // Consume the terminator
        byte terminatorByte = io.readSignedByte(8);
        if(terminatorByte != (byte) 0x00) {
            throw new BufferException("Expected 0x00, but found " + terminatorByte);
        }

        return stringValue;
    }

    public static void serializeZeroTerminatedString(WriteBuffer io, String data) throws BufferException {
        io.writeString(data.length()  * 8, data, WithOption.WithName("stringValue"), WithOption.WithEncoding("ASCII"));
        io.writeSignedByte(8, (byte) 0x00, WithOption.WithName("terminator"));
    }

    public static int lengthZeroTerminatedString(String data) {
        return (data.length() + 1) * 8;
    }

}
