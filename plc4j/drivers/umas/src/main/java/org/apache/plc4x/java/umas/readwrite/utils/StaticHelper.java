/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.umas.readwrite.utils;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.*;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

/**
 * Helper methods for UMAS protocol parsing and serialization, invoked
 * via STATIC_CALL expressions in the UMAS mspec.
 */
public class StaticHelper {

    /**
     * Parses a null-terminated string from the buffer.
     *
     * <p>If {@code stringLength >= 0}, reads exactly {@code stringLength} bytes and
     * returns the content up to the first null terminator (or the full length if none found).
     *
     * <p>If {@code stringLength == -1}, scans forward byte-by-byte until a null terminator
     * is found (variable-length mode, used by UmasUDTDefinition).
     *
     * @param readBuffer   the buffer to read from
     * @param stringLength the fixed number of bytes to read, or -1 for variable-length
     * @return the decoded string (without null terminator)
     * @throws ParseException if reading fails
     */
    public static String parseTerminatedString(ReadBuffer readBuffer, int stringLength) throws ParseException  {
        if (stringLength == -1) {
            // Variable-length: scan for null terminator
            List<Byte> stringBytes = new ArrayList<>();
            for (byte curByte = readBuffer.readSignedByte(8); curByte != (byte) 0x00;
                 curByte = readBuffer.readSignedByte(8)) {
                stringBytes.add(curByte);
            }
            byte[] bytes = new byte[stringBytes.size()];
            for (int i = 0; i < stringBytes.size(); i++) {
                bytes[i] = stringBytes.get(i);
            }
            return new String(bytes, StandardCharsets.UTF_8);
        }

        // Fixed-length: read exactly stringLength bytes, trim at null
        return parseTerminatedStringBytes(readBuffer, stringLength);
    }

    /**
     * Serializes a null-terminated string to the buffer.
     *
     * <p>If {@code stringLength >= 0}, writes the string content padded with null bytes
     * to fill exactly {@code stringLength} bytes.
     *
     * <p>If {@code stringLength == -1}, writes the string followed by a single null terminator.
     *
     * @param writeBuffer  the buffer to write to
     * @param value        the string value to serialize
     * @param stringLength the fixed field width, or -1 for variable-length
     * @throws SerializationException if writing fails
     */
    public static void serializeTerminatedString(WriteBuffer writeBuffer, String value, int stringLength) throws SerializationException {
        byte[] stringBytes = value.getBytes(StandardCharsets.UTF_8);

        if (stringLength == -1) {
            // Variable-length: write string + null terminator
            writeBuffer.writeString("value", stringBytes.length * 8, value, WithReaderWriterArgs.WithEncoding("UTF-8"));
            writeBuffer.writeSignedByte("terminator", 8, (byte) 0x00);
        } else {
            // Fixed-length: write string + null padding to fill stringLength bytes
            int bytesToWrite = Math.min(stringBytes.length, stringLength);
            if (bytesToWrite > 0) {
                writeBuffer.writeString("value", bytesToWrite * 8, value.substring(0, bytesToWrite), WithReaderWriterArgs.WithEncoding("UTF-8"));
            }
            for (int i = bytesToWrite; i < stringLength; i++) {
                writeBuffer.writeSignedByte("padding", 8, (byte) 0x00);
            }
        }
    }

    /**
     * Overload accepting {@link PlcValue} for use in DataItem serialization context,
     * where the code generator passes a PlcValue rather than a raw String.
     */
    public static void serializeTerminatedString(WriteBuffer writeBuffer, PlcValue value, int stringLength) throws SerializationException {
        serializeTerminatedString(writeBuffer, value.getString(), stringLength);
    }

    /**
     * Parses a null-terminated string from the buffer as raw bytes, used by the STRING
     * data type in the DataItem dataIo. Reads up to {@code numberOfValues} bytes and
     * returns the content up to the first null terminator as a String.
     *
     * @param readBuffer     the buffer to read from
     * @param numberOfValues the maximum number of bytes to read
     * @return the decoded string
     * @throws ParseException if reading fails
     */
    public static String parseTerminatedStringBytes(ReadBuffer readBuffer, int numberOfValues) throws ParseException {
        byte[] bytes = readBuffer.readByteArray("value", numberOfValues);
        for (int i = 0; i < bytes.length; i++) {
            if (bytes[i] == 0x00) {
                return new String(bytes, 0, i, StandardCharsets.UTF_8);
            }
        }
        return new String(bytes, StandardCharsets.UTF_8);
    }

}
