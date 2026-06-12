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
package org.apache.plc4x.java.can.generic.transport;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class GenericFrameTest {

    private static final WithOption[] OPTS = {
        WithByteBasedOption.WithByteOrder("BIG_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    @Test
    void accessorsAndLengthFields() {
        GenericFrame frame = new GenericFrame(0x123, new byte[]{1, 2, 3});
        assertThat(frame.getNodeId()).isEqualTo(0x123);
        assertThat(frame.getData()).containsExactly(1, 2, 3);
        assertThat(frame.getLengthInBytes()).isEqualTo(16);
        assertThat(frame.getLengthInBits()).isEqualTo(16 * 8);
    }

    @Test
    void serializeProducesCanIdDlcPaddingAndPaddedData() throws Exception {
        WriteBufferByteBased buf = new WriteBufferByteBased(new byte[16], OPTS);
        new GenericFrame(0x123, new byte[]{1, 2, 3}).serialize(buf);
        byte[] out = buf.getBytes();
        // CAN ID big-endian uint32
        assertThat(out[0]).isZero();
        assertThat(out[1]).isZero();
        assertThat(out[2]).isEqualTo((byte) 0x01);
        assertThat(out[3]).isEqualTo((byte) 0x23);
        // DLC = payload length
        assertThat(out[4]).isEqualTo((byte) 3);
        // 3 bytes of padding follow
        assertThat(out[5]).isZero();
        assertThat(out[6]).isZero();
        assertThat(out[7]).isZero();
        // First 3 payload bytes
        assertThat(out[8]).isEqualTo((byte) 1);
        assertThat(out[9]).isEqualTo((byte) 2);
        assertThat(out[10]).isEqualTo((byte) 3);
        // Remaining 5 bytes are zero-padded to fill the 16-byte frame.
        assertThat(out[11]).isZero();
        assertThat(out[15]).isZero();
    }

    @Test
    void serializeClampsDlcToEightBytes() throws Exception {
        // CAN payload is at most 8 bytes — the extra trailing bytes are
        // silently truncated.
        WriteBufferByteBased buf = new WriteBufferByteBased(new byte[16], OPTS);
        new GenericFrame(0, new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).serialize(buf);
        byte[] out = buf.getBytes();
        assertThat(out[4]).isEqualTo((byte) 8);
    }
}
