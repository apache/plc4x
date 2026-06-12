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

package org.apache.plc4x.java.ads.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DirectAdsStringTagTest {

    @Test
    void ofParsesStringAddress() {
        DirectAdsStringTag tag = DirectAdsStringTag.of("100/200:STRING(80)");
        assertEquals(100L, tag.getIndexGroup());
        assertEquals(200L, tag.getIndexOffset());
        assertEquals("STRING", tag.getPlcDataType());
        assertEquals(80, tag.getStringLength());
        assertEquals(1, tag.getNumberOfElements());
    }

    @Test
    void ofParsesHexStringAddressWithArray() {
        DirectAdsStringTag tag = DirectAdsStringTag.of("0x10/0x20:WSTRING(40)[3]");
        assertEquals(0x10L, tag.getIndexGroup());
        assertEquals(0x20L, tag.getIndexOffset());
        assertEquals("WSTRING", tag.getPlcDataType());
        assertEquals(40, tag.getStringLength());
        assertEquals(3, tag.getNumberOfElements());
    }

    @Test
    void ofRejectsInvalidAddress() {
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsStringTag.of("100/200:DINT"));
    }

    @Test
    void factoryWithFields() {
        DirectAdsStringTag tag = DirectAdsStringTag.of(1, 2, "STRING", 10, null);
        assertEquals(10, tag.getStringLength());
        assertEquals(1, tag.getNumberOfElements());
    }

    @Test
    void matchesValidAndInvalid() {
        assertTrue(DirectAdsStringTag.matches("100/200:STRING(80)"));
        assertFalse(DirectAdsStringTag.matches("100/200:DINT"));
    }

    @Test
    void getAddressStringFormatsBase() {
        DirectAdsStringTag single = DirectAdsStringTag.of(1, 2, "STRING", 10, 1);
        assertFalse(single.getAddressString().contains("["));
        DirectAdsStringTag array = DirectAdsStringTag.of(1, 2, "STRING", 10, 4);
        assertTrue(array.getAddressString().contains("[4]"));
    }

    @Test
    void toStringContainsFields() {
        assertTrue(DirectAdsStringTag.of(1, 2, "STRING", 10, 1).toString().contains("stringLength=10"));
    }

    @Test
    void serialize() throws Exception {
        DirectAdsStringTag tag = DirectAdsStringTag.of(1, 2, "STRING", 10, 1);
        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[64],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"));
        tag.serialize(buffer);
        assertNotNull(buffer.getBytes());
    }
}
