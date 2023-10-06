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
package org.apache.plc4x.java.canopen;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class ManualParserTest {

    public static final int EXTENDED_FRAME_FORMAT_FLAG = 0b10000000_00000000_00000000_00000000;
    public static final int REMOTE_TRANSMISSION_FLAG = 0b01000000_00000000_00000000_00000000;
    public static final int ERROR_FRAME_FLAG = 0b00100000_00000000_00000000_00000000;

    public static final int STANDARD_FORMAT_IDENTIFIER_MASK = 0b00000000_00000000_00000111_11111111;
    public static final int EXTENDED_FORMAT_IDENTIFIER_MASK = 0b00011111_11111111_11111111_11111111;

    // cansend 5A1#11.2233.44556677.88
    final String STANDARD = "a1050000080000001122334455667788";

    // cansend 5A1#R
    String STANDARD_REPLY = "a1050040000000000000000000000000";

    // cansend 1E6EC676#05.05.1F.26.C3
    final String EXTENDED = "76c66e9e0500000005051f26c3000000";

    @Test
    public void readBufferTest() throws Exception {
        ReadBuffer buffer = new ReadBufferByteBased(new byte[]{(byte) 0xA1, 0x05, 0x00, 0x00}, ByteOrder.LITTLE_ENDIAN);
        int value = buffer.readInt(32);

        assertEquals(value, 0x5A1);
    }

    @Test
    public void standardFrameParser() throws Exception {
        SocketCanFrameStub frame = parse(STANDARD);
        //System.out.println(frame);

        assertEquals(frame.id, 0x5A1);
        assertEquals(frame.extended, false);
        assertEquals(frame.remote, false);
        assertEquals(frame.error, false);
        assertEquals(frame.data.length, 8);
    }

    @Test
    public void extendedFrameParser() throws Exception {
        SocketCanFrameStub frame = parse(EXTENDED);
        //System.out.println(frame);

        assertEquals(frame.id, 0x1e6ec676);
        assertEquals(frame.extended, true);
        assertEquals(frame.remote, false);
        assertEquals(frame.error, false);
        assertEquals(frame.data.length, 5);
    }

    public static SocketCanFrameStub parse(String hex) throws Exception {
        byte[] input = Hex.decodeHex(hex.toCharArray());

        ReadBufferByteBased readBuffer = new ReadBufferByteBased(input, ByteOrder.LITTLE_ENDIAN);
        int rawId = readBuffer.readInt(32);
        boolean extended = (rawId & EXTENDED_FRAME_FORMAT_FLAG) != 0;
        boolean remote = (rawId & REMOTE_TRANSMISSION_FLAG) != 0;
        boolean error = (rawId & ERROR_FRAME_FLAG) != 0;
        int id = extended ? (rawId & EXTENDED_FORMAT_IDENTIFIER_MASK) : (rawId & STANDARD_FORMAT_IDENTIFIER_MASK);
        int length = readBuffer.readByte();
        byte[] data = readBuffer.getBytes(8, 8 + length);

        return new SocketCanFrameStub(
            id, extended, remote, error, data
        );
    }

    static class SocketCanFrameStub {
        public final int id;
        public final boolean extended;
        public final boolean remote;
        public final boolean error;
        public final byte[] data;

        public SocketCanFrameStub(int id, boolean extended, boolean remote, boolean error, byte[] data) {
            this.id = id;
            this.extended = extended;
            this.remote = remote;
            this.error = error;
            this.data = data;
        }

        public String toString() {
            return "CAN Frame ID=" + Integer.toHexString(id) + ", extended=" + extended + ", remote=" + remote + ", error=" + error + ", data=" + Hex.encodeHexString(data);
        }
    }
}
