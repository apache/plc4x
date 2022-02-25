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
package org.apache.plc4x.java.test.readwrite.utils;

import org.apache.plc4x.java.spi.generation.*;

public class StaticHelper {

    public static boolean parseBit(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeBit(WriteBuffer io, boolean data) {
    }

    public static byte parseByte(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeByte(WriteBuffer io, byte data) {
    }

    public static byte parseInt8(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeInt8(WriteBuffer io, byte data) {
    }

    public static short parseUint8(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeUint8(WriteBuffer io, short data) {
    }

    public static float parseFloat(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeFloat(WriteBuffer io, float data) {
    }

    public static double parseDouble(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeDouble(WriteBuffer io, double data) {
    }

    public static String parseString(ReadBuffer io) {
        throw new IllegalArgumentException("Hurz!");
    }

    public static void serializeString(WriteBuffer io, String data) {
    }

}
