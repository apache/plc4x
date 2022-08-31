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
package org.apache.plc4x.java.ads.readwrite.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

public class StaticHelper {

    public static String parseAmsString(ReadBuffer readBuffer, int stringLength, String encoding) {
        try {
            if ("UTF-8".equalsIgnoreCase(encoding)) {
                List<Byte> bytes = new ArrayList<>();
                for(int i = 0; (i < stringLength) && readBuffer.hasMore(8); i++) {
                    final byte curByte = readBuffer.readByte();
                    if (curByte != 0) {
                        bytes.add(curByte);
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for(; (i < stringLength) && readBuffer.hasMore(8); i++) {
                            readBuffer.readByte();
                        }
                        break;
                    }
                }
                // Read the terminating byte.
                readBuffer.readByte();
                final byte[] byteArray = new byte[bytes.size()];
                for (int i = 0; i < bytes.size(); i++) {
                    byteArray[i] = bytes.get(i);
                }
                return new String(byteArray, StandardCharsets.UTF_8);
            } else if ("UTF-16".equalsIgnoreCase(encoding)) {
                List<Byte> bytes = new ArrayList<>();
                for(int i = 0; (i < stringLength) && readBuffer.hasMore(16); i++) {
                    final short curShort = readBuffer.readShort(16);
                    if (curShort != 0) {
                        bytes.add((byte) (curShort >>> 8));
                        bytes.add((byte) (curShort & 0xFF));
                    } else {
                        // Gobble up the remaining data, which is not added to the string.
                        i++;
                        for(; (i < stringLength) && readBuffer.hasMore(16); i++) {
                            readBuffer.readShort(16);
                        }
                        break;
                    }
                }
                // Read the terminating byte.
                readBuffer.readByte();
                final byte[] byteArray = new byte[bytes.size()];
                for (int i = 0; i < bytes.size(); i++) {
                    byteArray[i] = bytes.get(i);
                }
                return new String(byteArray, StandardCharsets.UTF_16);
            } else {
                throw new PlcRuntimeException("Unsupported string encoding " + encoding);
            }
        } catch (ParseException e) {
            throw new PlcRuntimeException("Error parsing string", e);
        }
    }

    public static void serializeAmsString(WriteBuffer writeBuffer, PlcValue value, int stringLength, Object encoding) {
        // TODO: Need to implement the serialization or we can't write strings
        throw new PlcRuntimeException("Not implemented yet");
    }

}
