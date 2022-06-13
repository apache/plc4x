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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class VintIo {

    public static int getLengthInBytes(long value) {
        int numBytes = 1;
        // Keep on shifting the value left 7 bits, till the value is 0.
        while((value & 0xFFFFFF80) != 0) {
            numBytes++;
            value = value >> 7;
        }
        return numBytes;
    }

    public static void serialize(long value, WriteBuffer writeBuffer) throws SerializationException {
        serializeInternally(value, writeBuffer, true);
    }

    private static void serializeInternally(long value, WriteBuffer writeBuffer, boolean lastByte) throws SerializationException {
        byte curValue = (byte) ((value & 0x7F) | (lastByte ? 0x00 : 0x80));
        long restValue = value >> 7;
        // Recursively serialize the others first.
        if(restValue != 0) {
            serializeInternally(restValue, writeBuffer, false);
        }
        // Now serialize the last byte.
        writeBuffer.writeUnsignedShort(8, curValue);
    }

    public static long parse(ReadBuffer readBuffer) throws ParseException {
        long curValue = 0;
        short curByte;
        do {
            // Get the next byte
            curByte = readBuffer.readUnsignedShort(8);
            // Add the 7 least significant bits to the value.
            curValue |= (curByte & 0x7F);
            // If the highest bit is set, shift the value right 7 bits and read the next byte.
            if((curByte & 0x80) != 0) {
                curValue = curValue << 7;
            }
        } while ((curByte & 0x80) != 0);
        return curValue;
    }

}
