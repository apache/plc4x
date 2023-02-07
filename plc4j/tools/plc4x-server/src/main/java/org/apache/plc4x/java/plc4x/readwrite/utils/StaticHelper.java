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
package org.apache.plc4x.java.plc4x.readwrite.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.logging.Level;
import java.util.logging.Logger;

public class StaticHelper {

    public static String parseString(ReadBuffer io,  String encoding) {
        try {
            if ("UTF-8".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                short stringLength = io.readUnsignedShort(8);
                // This is the total length of the string on the PLC (Not necessarily the number of characters read)
                final byte[] byteArray = new byte[stringLength];
                for (int i = 0; (i < stringLength) && io.hasMore(8); i++) {
                    final byte curByte = io.readByte();
                    byteArray[i] = curByte;
                }
                return new String(byteArray, StandardCharsets.UTF_8);
            } else if ("UTF-16".equalsIgnoreCase(encoding)) {
                // This is the maximum number of bytes a string can be long.
                int stringLength = io.readUnsignedInt(16);
                final byte[] byteArray = new byte[stringLength * 2];
                for (int i = 0; (i < stringLength) && io.hasMore(16); i++) {
                    final short curShort = io.readShort(16);
                    byteArray[i * 2] = (byte) (curShort >>> 8);
                    byteArray[(i * 2) + 1] = (byte) (curShort & 0xFF);
                }
                return new String(byteArray, StandardCharsets.UTF_16);
            } else {
                throw new PlcRuntimeException("Unsupported string encoding " + encoding);
            }
        } catch (ParseException e) {
            throw new PlcRuntimeException("Error parsing string", e);
        }
    }
    public static void serializeString(WriteBuffer io, PlcValue value, String encoding) {
        String valueString = value.getString();
        valueString = valueString == null ? "" : valueString;

        if ("UTF-8".equalsIgnoreCase(encoding)) {
            final byte[] raw = valueString.getBytes(StandardCharsets.UTF_8);
            try {
                io.writeByte((byte) raw.length);
                for (byte b : raw) {
                    io.writeByte(b);
                }
            }
            catch (SerializationException ex) {
                Logger.getLogger(StaticHelper.class.getName()).log(Level.SEVERE, null, ex);
            }
        } else if ("UTF-16".equalsIgnoreCase(encoding)) {
            final byte[] raw = valueString.getBytes(StandardCharsets.UTF_16);
            try {
                io.writeUnsignedInt(16, raw.length);
                for (int i = 0; i < raw.length; i++) {
                    io.writeByte( raw[i]);
                }
            }
            catch (SerializationException ex) {
                Logger.getLogger(StaticHelper.class.getName()).log(Level.SEVERE, null, ex);
            }
        } else {
            throw new PlcRuntimeException("Unsupported string encoding " + encoding);
        }

    }
}