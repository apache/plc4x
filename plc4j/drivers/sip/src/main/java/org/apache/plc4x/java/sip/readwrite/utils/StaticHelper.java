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

package org.apache.plc4x.java.sip.readwrite.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class StaticHelper {

    public static String asString(byte[] bytes) {
        return new String(bytes);
    }

    public static int untilToken(ReadBuffer readBuffer, String terminator, int keep) {
        int start = readBuffer.getPos();
        StringBuilder buffer = new StringBuilder();
        int length = terminator.length();
        int retrieved = 0;
        boolean success = false;
        while (readBuffer.hasMore(8) && retrieved < length) {
            try {
                buffer.append((char) readBuffer.readByte());
                if (buffer.length() >= length) {
                    if (buffer.toString().endsWith(terminator)) {
                        success = true;
                        break;
                    }
                }
            } catch (ParseException e) {
                throw new RuntimeException(e);
            }
        }

        if (!success) {
            throw new PlcRuntimeException("Failed to reach termination sequence for array");
        }

        int end = readBuffer.getPos();
        readBuffer.reset(start);
        return end - start - keep;
    }

    public static boolean until(ReadBuffer readBuffer, String terminator) {
        int start = readBuffer.getPos();
        StringBuilder buffer = new StringBuilder();
        int length = terminator.length();
        int retrieved = 0;
        while (readBuffer.hasMore(8) && retrieved < length) {
            try {
                buffer.append((char) readBuffer.readByte());
                if (buffer.length() == length) {
                    if (buffer.toString().equals(terminator)) {
                        readBuffer.reset(start + length - 1);
                        return false;
                    }
                    break;
                }
            } catch (ParseException e) {
                throw new RuntimeException(e);
            }
        }

        readBuffer.reset(start);
        return true;
    }

    public static String readStringTill(ReadBuffer readBuffer, String terminator) throws ParseException {
        int start = readBuffer.getPos();
        StringBuilder buffer = new StringBuilder();
        int length = terminator.length();
        while (readBuffer.hasMore(8)) {
            buffer.append((char) readBuffer.readByte());
            if (buffer.length() >= length) {
                if (buffer.subSequence(buffer.length() - length, buffer.length()).equals(terminator)) {
                    break;
                }
            }
        }

        int end = readBuffer.getPos();
        readBuffer.reset(end - length);
        if (start == end) {
            return "";
        }
        return buffer.subSequence(0, buffer.length() - length).toString();
    }

    public static String readString(ReadBuffer readBuffer, String terminator) throws ParseException {
        int start = readBuffer.getPos();
        StringBuilder buffer = new StringBuilder();
        int length = terminator.length();
        while (readBuffer.hasMore(8)) {
            buffer.append((char) readBuffer.readByte());
            if (buffer.length() >= length) {
                if (buffer.subSequence(buffer.length() - length, buffer.length()).equals(terminator)) {
                    break;
                }
            }
        }

        int end = readBuffer.getPos();
        //readBuffer.reset(end - length);
        if (start == end) {
            return "";
        }
        //return buffer.subSequence(0, buffer.length() - length).toString();
        return buffer.toString();
    }

    public static void writeString(WriteBuffer writeBuffer, String method/*, String terminator*/) throws SerializationException {
        writeBuffer.writeString(8 * method.length()/* + (terminator.length() * 8)*/ , method/* + terminator*/);
    }

    public static void writeStringTill(WriteBuffer writeBuffer, String method) throws SerializationException {
        writeBuffer.writeString(8 * method.length(), method);
    }
}
