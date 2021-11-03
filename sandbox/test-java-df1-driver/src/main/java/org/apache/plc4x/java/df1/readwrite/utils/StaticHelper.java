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
package org.apache.plc4x.java.df1.readwrite.utils;

import com.github.snksoft.crc.CRC;
import org.apache.plc4x.java.df1.readwrite.DF1Command;
import org.apache.plc4x.java.df1.readwrite.io.DF1CommandIO;
import org.apache.plc4x.java.spi.generation.*;

import java.util.List;

public class StaticHelper {

    public static int crcCheck(short destinationAddress, short sourceAddress, DF1Command command) {
        // CRC-16/DF-1
        CRC crc = new CRC(new CRC.Parameters(16, 0x8005, 0x0000, true, true, 0x0000));
        long df1crc = crc.init();
        df1crc = crc.update(df1crc, new byte[]{(byte) destinationAddress, (byte) sourceAddress});
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(command.getLengthInBytes(), ByteOrder.BIG_ENDIAN);
        try {
            DF1CommandIO.staticSerialize(writeBuffer, command);
        } catch (ParseException e) {
            throw new RuntimeException(e);
        }
        df1crc = crc.update(df1crc, writeBuffer.getData());
        df1crc = crc.update(df1crc, new byte[]{0x03});
        return crc.finalCRC16(df1crc) & 0xFFFF;
    }

    public static boolean dataTerminate(ReadBuffer io) {
        ReadBufferByteBased rbbb = (ReadBufferByteBased) io;
        try {
            // The byte sequence 0x10 followed by 0x03 indicates the end of the message,
            // so if we would read this, we abort the loop and stop reading data.
            if ((rbbb.peekByte(0) == (byte) 0x10) && (rbbb.peekByte(1) == (byte) 0x03)) {
                return true;
            }
        } catch (ParseException e) {
            // Just ignore and return false.
        }
        return false;
    }

    public static short readData(ReadBuffer io) {
        ReadBufferByteBased rbbb = (ReadBufferByteBased) io;
        try {
            // If we read a 0x10, this has to be followed by another 0x10, which is how
            // this value is escaped in DF1, so if we encounter two 0x10, we simply ignore the first.
            if ((rbbb.peekByte(0) == (byte) 0x10) && (rbbb.peekByte(1) == 0x10)) {
                io.readByte();
            }
            return io.readUnsignedShort(8);
        } catch (ParseException e) {
            throw new RuntimeException("Error parsing data", e);
        }
    }

    public static void writeData(WriteBuffer io, short data) {
        try {
            // If a value is 0x10, this has to be duplicated in order to be escaped.
            if ((data == (short) 0x10)) {
                io.writeUnsignedShort(8, (short) 0x10);
            }
            io.writeUnsignedShort(8, data);
        } catch (ParseException e) {
            throw new RuntimeException("Error parsing data", e);
        }
    }

    public static int dataLength(byte[] data) {
        int i = 0;
        for (short dataByte : data) {
            // If a value is 0x10, this has to be duplicated which increases the message size by one.
            if (dataByte == 0x10) {
                i++;
            }
            i++;
        }
        return i;
    }

}
