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
package org.apache.plc4x.java.canopen.readwrite.utils;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateExpeditedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateUploadResponsePayload;
import org.apache.plc4x.java.canopen.readwrite.SDOSegmentUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import static org.apache.plc4x.java.spi.generation.StaticHelper.COUNT;

public class StaticHelper {

    public static CANOpenService serviceId(short identifier) {
        return CANOpenService.enumForValue((byte) (identifier >> 7));
    }

    public static int uploadPadding(SDOSegmentUploadResponse payload) {
        return 7 - payload.getData().length;
    }

    public static int count(boolean expedited, boolean indicated, SDOInitiateUploadResponsePayload payload) {
        return expedited && indicated && payload instanceof SDOInitiateExpeditedUploadResponse ? 4 - COUNT(((SDOInitiateExpeditedUploadResponse) payload).getData()) : 0;
    }

    public static void writeFunction(WriteBuffer io, short identifier) {
        // NOOP - a placeholder to let mspec compile
    }

    public static Object parseString(ReadBuffer io, int length, String charset) {
        return io.readString(8 * length, charset);
    }

    public static void serializeString(WriteBuffer io, PlcValue value, String charset) throws ParseException {
        io.writeString(8, charset, value.getString());
    }

    public static byte[] parseByteArray(ReadBuffer io, Integer length) {
        return new byte[0];
    }
}
