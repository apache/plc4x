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
package org.apache.plc4x.java.firmata.readwrite.utils;

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

public class StaticHelper {

    public static boolean isSysexEnd(ReadBuffer io) {
        int positionInBits = io.getPositionInBits();
        try {
            byte[] bytes = io.readBits(8);
            return bytes.length == 1 && bytes[0] == (byte) 0xF7;
        } catch (BufferException e) {
            return false;
        } finally {
            io.setPositionInBits(positionInBits);
        }
    }

    public static byte parseSysexString(ReadBuffer io) throws BufferException {
        // Each "sysex string" byte is followed by a zero padding byte (Firmata
        // 7-bit MIDI framing). Read both via readBits so the WriteBufferXmlBased
        // counterpart renders them with the matching dataType="byte" form.
        //
        // A failed read has to be reported: the manual array reading this runs until isSysexEnd
        // finds the end marker, so returning a value for the end of the data would never end.
        byte[] valueBytes = io.readBits(8);
        byte[] padBytes = io.readBits(8);
        // padding byte intentionally discarded
        if (padBytes.length == 1) {
            // no-op
        }
        return valueBytes.length == 1 ? valueBytes[0] : 0;
    }

    public static void serializeSysexString(WriteBuffer io, byte data) {
        try {
            io.writeBits(8, new byte[]{data});
            io.writeBits(8, new byte[]{(byte) 0x00});
        } catch (BufferException e) {
            // intentionally ignored — legacy contract
        }
    }

    public static int lengthSysexString(byte[] data) {
        // Each byte of the logical sysex "string" is sent on the wire as two
        // bytes (value + 0x00 padding). The generated SysexCommand length
        // computation expects this contribution in *bits*, so multiply by 8.
        return data.length * 2 * 8;
    }

}
