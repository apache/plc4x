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

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class StaticHelperTest {

    private static final WithOption[] OPTS = {
        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    @Test
    void writeSizeIndexToByteCountReturnsExpectedSizes() {
        // sizeIndex 17 is the STRING marker, 1 byte per char.
        assertThat(StaticHelper.writeSizeIndexToByteCount(17)).isEqualTo(1);
        // Indices 1..4 are powers of two: 1, 2, 4, 8.
        assertThat(StaticHelper.writeSizeIndexToByteCount(1)).isEqualTo(1);
        assertThat(StaticHelper.writeSizeIndexToByteCount(2)).isEqualTo(2);
        assertThat(StaticHelper.writeSizeIndexToByteCount(3)).isEqualTo(4);
        assertThat(StaticHelper.writeSizeIndexToByteCount(4)).isEqualTo(8);
        // Outside the known range: fallback to the index itself.
        assertThat(StaticHelper.writeSizeIndexToByteCount(99)).isEqualTo(99);
    }

    @Test
    void parseTerminatedStringVariableLengthScansUntilNull() throws Exception {
        // stringLength == -1 means "scan byte-by-byte until null terminator".
        byte[] data = "hello\0xyz".getBytes();
        ReadBufferByteBased buf = new ReadBufferByteBased(data, OPTS);
        assertThat(StaticHelper.parseTerminatedString(buf, -1)).isEqualTo("hello");
    }

    @Test
    void parseTerminatedStringBytesStopsAtNullByte() throws Exception {
        byte[] data = "hi\0tail".getBytes();
        ReadBufferByteBased buf = new ReadBufferByteBased(data, OPTS);
        // Length parameter passes through to readBits as a bit count.
        assertThat(StaticHelper.parseTerminatedStringBytes(buf, data.length * 8)).isEqualTo("hi");
    }

    @Test
    void parseTerminatedStringBytesUsesEntireBufferWhenUnterminated() throws Exception {
        byte[] data = "abc".getBytes();
        ReadBufferByteBased buf = new ReadBufferByteBased(data, OPTS);
        assertThat(StaticHelper.parseTerminatedStringBytes(buf, data.length * 8)).isEqualTo("abc");
    }

    @Test
    void serializeTerminatedStringPaddedWithNulls() throws Exception {
        WriteBufferByteBased buf = new WriteBufferByteBased(new byte[8], OPTS);
        StaticHelper.serializeTerminatedString(buf, "hi", 8);
        byte[] out = buf.getBytes();
        assertThat(out[0]).isEqualTo((byte) 'h');
        assertThat(out[1]).isEqualTo((byte) 'i');
        assertThat(out[2]).isZero();
    }

    @Test
    void serializePlcValueTerminatedStringDelegatesToStringOverload() throws Exception {
        WriteBufferByteBased buf = new WriteBufferByteBased(new byte[8], OPTS);
        StaticHelper.serializeTerminatedString(buf, new PlcSTRING("hi"), 8);
        byte[] out = buf.getBytes();
        assertThat(out[0]).isEqualTo((byte) 'h');
        assertThat(out[1]).isEqualTo((byte) 'i');
    }
}
