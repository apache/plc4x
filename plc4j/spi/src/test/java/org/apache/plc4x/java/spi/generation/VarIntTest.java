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
package org.apache.plc4x.java.spi.generation;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class VarIntTest {

    @ParameterizedTest
    @CsvSource({
        "'00', 0, 'Single‑byte encoding for 0'",
        "'01', 1, 'Single‑byte encoding for 1'",
        "'7F', 127, 'Single‑byte maximum (all 7 data bits set)'",
        "'8100', 128, 'Two‑byte encoding for 128; groups: [1, 0] (first byte 0x81, second byte 0x00)'",
        "'8148', 200, 'Two‑byte encoding for 200; groups: [1, 72] (0x81, 0x48)'",
        "'8600', 768, 'Two‑byte encoding for 768; groups: [6, 0] (0x86, 0x00) – as you indicated'",
        "'CE10', 10000, 'Two‑byte encoding for 10000; groups: [78, 16] (first byte: 78→0xCE, second: 16→0x10)'",
        "'FF7F', 16383, 'Two‑byte maximum (groups: [127, 127]; first byte 0xFF, second byte 0x7F)'",
        "'818000', 16384, 'Three‑byte encoding; smallest number needing three bytes; groups: [1, 0, 0]'",
        "'FFFF7F', 2097151, 'Three‑byte maximum; groups: [127, 127, 127] (encoded as 0xFF, 0xFF, 0x7F)'",
        "'82B19640', 5000000, 'Four‑byte encoding for 5000000; groups: [2, 49, 22, 64] (0x82, 0xB1, 0x96, 0x40)'",
        "'FFFFFF7F', 268435455, 'Four‑byte maximum; groups: [127, 127, 127, 127] (encoded as 0xFF, 0xFF, 0xFF, 0x7F)'",
    })
    void writeVarUintRoundtrip(String hexString, long expectedValue, String description) throws Exception {
        byte[] serialized = Hex.decodeHex(hexString);

        // Parse the given array into a value
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(serialized);
        long value = readBuffer.readUnsignedLong("", 32, WithOption.WithEncoding("VARUDINT"));
        assertEquals(expectedValue, value);

        // Serialize the given value into a byte array
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(serialized.length);
        writeBuffer.writeUnsignedLong("", 32, expectedValue, WithOption.WithEncoding("VARUDINT"));
        byte[] result = writeBuffer.getBytes();
        assertArrayEquals(serialized, result, description);
    }

    @ParameterizedTest
    @CsvSource({
        "'00', 0, '0 encoded in one byte.'",
        "'01', 1, 'Minimal positive value.'",
        "'3F', 63, 'Maximum positive (0x3F = 63).'",
        "'40', -64, '0x40 = 64 becomes 64–128 = –64.'",
        "'7F', -1, '0x7F = 127 becomes 127–128 = –1.'",
        "'8040', 64, 'Groups: 1 and 0 → (1<<7) + 0 = 128.'",
        "'8148', 200, 'Groups: 1 and 0x48 (72) → (1<<7) + 72 = 128 + 72 = 200.'",
        "'FF1C', -100, 'Groups: 0xFF (127) and 0x1C (28) → (127<<7) + 28 = 16284; since 16284 ≥ 8192, subtract 16384: 16284–16384 = –100.'",
        "'819C20', 20000, 'Groups: 1, 28, 32 → (1<<14) + (28<<7) + 32 = 16384 + 3584 + 32 = 20000.'",
        "'FEE360', -20000, 'Groups: 0xFE (126), 0xE3 (99), 0x60 (96) → (126<<14) + (99<<7) + 96 = 2077152; 2077152 – 2097152 = –20000.'",
        "'82B19640', 5000000, 'Groups: 2, 49, 22, 64 → (2<<21) + (49<<14) + (22<<7) + 64 = 5000000.'",
        "'FDCEE940', -5000000, 'Groups: 0xFD (125), 0xCF (79), 0xA8 (40), 0x00 (0) → (125<<21) + (79<<14) + (40<<7) + 0 = 263435456; 263435456 – 268435456 = –5000000.'",
        "'BFFFFF7F', 134217727, 'Maximum positive value in 28 bits. (Groups: 63, 127, 127, 127.)'",
        "'C0808000', -134217728, 'Minimum negative value in 28 bits. (Groups: 64, 0, 0, 0 → 64<<21 = 134217728; then 134217728 – 268435456 = –134217728.)'",
        "'FF8CB2BE48', -242442424, 'Example'",
        "'F880808000', -2147483648, 'Min 32 bit signed integer'",
        "'87FFFFFF7F', 2147483647, 'Max 32 bit signed integer'"
    })
    void testVarIntRoundtrip(String hexString, int expectedValue, String description) throws Exception {
        byte[] serialized = Hex.decodeHex(hexString);

        // Parse the given array into a value
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(serialized);
        int value = readBuffer.readInt("", 32, WithOption.WithEncoding("VARDINT"));
        assertEquals(expectedValue, value);

        // Serialize the given value into a byte array
        WriteBufferByteBased buffer = new WriteBufferByteBased(serialized.length);
        buffer.writeInt("", 32, expectedValue, WithOption.WithEncoding("VARDINT"));
        byte[] result = buffer.getBytes();
        assertArrayEquals(serialized, result, description);
    }

}
