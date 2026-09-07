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

package org.apache.plc4x.java.ads.readwrite.utils;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class StaticHelperTest {

    private static final WithOption[] OPTIONS = new WithOption[]{
        WithOption.WithEncoding("ASCII"),
        WithOption.WithSignedIntegerEncoding("twos-complement")
    };

    @Test
    void roundTripZeroTerminatedString() throws Exception {
        WriteBufferByteBased wb = new WriteBufferByteBased(new byte[16], OPTIONS);
        StaticHelper.serializeZeroTerminatedString(wb, "hello");

        ReadBufferByteBased rb = new ReadBufferByteBased(wb.getBytes(), OPTIONS);
        assertEquals("hello", StaticHelper.parseZeroTerminatedString(rb, 5));
    }

    @Test
    void parseZeroTerminatedStringFailsOnMissingTerminator() {
        byte[] data = {'h', 'e', 'l', 'l', 'o', 0x01};
        ReadBufferByteBased rb = new ReadBufferByteBased(data, OPTIONS);
        assertThrows(BufferException.class, () -> StaticHelper.parseZeroTerminatedString(rb, 5));
    }

    @Test
    void lengthZeroTerminatedString() {
        assertEquals((4 + 1) * 8, StaticHelper.lengthZeroTerminatedString("test"));
        assertEquals((0 + 1) * 8, StaticHelper.lengthZeroTerminatedString(""));
    }

    @Test
    void instantiationForCoverage() {
        assertNotNull(new StaticHelper());
    }
}
