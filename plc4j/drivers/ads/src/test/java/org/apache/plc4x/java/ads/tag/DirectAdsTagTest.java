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
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DirectAdsTagTest {

    @Test
    void ofParsesDecimalAddress() {
        DirectAdsTag tag = DirectAdsTag.of("16448/1234:DINT");
        assertEquals(16448L, tag.getIndexGroup());
        assertEquals(1234L, tag.getIndexOffset());
        assertEquals("DINT", tag.getPlcDataType());
        assertEquals(1, tag.getNumberOfElements());
    }

    @Test
    void ofParsesHexAddress() {
        DirectAdsTag tag = DirectAdsTag.of("0x4040/0xFF:LREAL[5]");
        assertEquals(0x4040L, tag.getIndexGroup());
        assertEquals(0xFFL, tag.getIndexOffset());
        assertEquals(5, tag.getNumberOfElements());
    }

    @Test
    void ofRejectsInvalidAddress() {
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("not-an-address"));
    }

    @Test
    void factoryWithFields() {
        DirectAdsTag tag = DirectAdsTag.of(1L, 2L, "BOOL", null);
        assertEquals(1, tag.getNumberOfElements());
        assertEquals("BOOL", tag.getPlcDataType());
    }

    @Test
    void constructorRejectsZeroElements() {
        assertThrows(IllegalArgumentException.class, () -> new DirectAdsTag(1, 1, "DINT", 0));
    }

    @Test
    void constructorRequiresDataType() {
        assertThrows(NullPointerException.class, () -> new DirectAdsTag(1, 1, null, 1));
    }

    @Test
    void matchesValidAndInvalid() {
        assertTrue(DirectAdsTag.matches("100/200:DINT"));
        assertFalse(DirectAdsTag.matches("foo"));
    }

    @Test
    void getPlcValueTypeKnownAndUnknown() {
        assertEquals(PlcValueType.DINT, DirectAdsTag.of(1, 1, "DINT", 1).getPlcValueType());
        assertEquals(PlcValueType.Struct, DirectAdsTag.of(1, 1, "MyType", 1).getPlcValueType());
    }

    @Test
    void getAddressStringFormatsBase() {
        DirectAdsTag single = DirectAdsTag.of(1, 2, "DINT", 1);
        assertFalse(single.getAddressString().contains("["));
        DirectAdsTag array = DirectAdsTag.of(1, 2, "DINT", 4);
        assertTrue(array.getAddressString().contains("[4]"));
    }

    @Test
    void arrayInfoOnlyForArrays() {
        assertTrue(DirectAdsTag.of(1, 1, "DINT", 1).getArrayInfo().isEmpty());
        assertEquals(1, DirectAdsTag.of(1, 1, "DINT", 5).getArrayInfo().size());
    }

    @Test
    void equalsAndHashCode() {
        DirectAdsTag a = DirectAdsTag.of(1, 2, "DINT", 1);
        DirectAdsTag b = DirectAdsTag.of(1, 2, "INT", 1);
        DirectAdsTag c = DirectAdsTag.of(3, 4, "DINT", 1);
        assertEquals(a, a);
        assertEquals(a, b);
        assertNotEquals(a, c);
        assertNotEquals(a, "x");
        assertEquals(a.hashCode(), b.hashCode());
    }

    @Test
    void toStringContainsFields() {
        assertTrue(DirectAdsTag.of(1, 2, "DINT", 1).toString().contains("indexGroup=1"));
    }

    @Test
    void serialize() throws Exception {
        DirectAdsTag tag = DirectAdsTag.of(1, 2, "DINT", 1);
        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[64],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"));
        tag.serialize(buffer);
        assertNotNull(buffer.getBytes());
    }
}
