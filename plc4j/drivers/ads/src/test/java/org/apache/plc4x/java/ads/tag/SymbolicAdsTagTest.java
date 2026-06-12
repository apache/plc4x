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
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.junit.jupiter.api.Assertions.*;

class SymbolicAdsTagTest {

    @Test
    void ofParsesValidAddress() {
        SymbolicAdsTag tag = SymbolicAdsTag.of("Main.value");
        assertEquals("Main.value", tag.getSymbolicAddress());
        assertEquals("Main.value", tag.getAddressString());
        assertNull(tag.getPlcValueType());
        assertTrue(tag.getArrayInfo().isEmpty());
    }

    @Test
    void ofRejectsInvalidAddress() {
        assertThrows(PlcInvalidTagException.class, () -> SymbolicAdsTag.of("!!!"));
    }

    @Test
    void matchesValidAndInvalid() {
        assertTrue(SymbolicAdsTag.matches("Main.value"));
        assertFalse(SymbolicAdsTag.matches("!!!"));
    }

    @Test
    void constructorRequiresAddress() {
        assertThrows(NullPointerException.class,
            () -> new SymbolicAdsTag(null, null, Collections.emptyList()));
    }

    @Test
    void equalsAndHashCode() {
        SymbolicAdsTag a = SymbolicAdsTag.of("Main.value");
        SymbolicAdsTag b = SymbolicAdsTag.of("Main.value");
        SymbolicAdsTag c = SymbolicAdsTag.of("Main.other");
        assertEquals(a, a);
        assertEquals(a, b);
        assertNotEquals(a, c);
        assertNotEquals(a, "not a tag");
        assertEquals(a.hashCode(), b.hashCode());
    }

    @Test
    void toStringContainsAddress() {
        assertTrue(SymbolicAdsTag.of("Main.value").toString().contains("Main.value"));
    }

    @Test
    void serialize() throws Exception {
        SymbolicAdsTag tag = SymbolicAdsTag.of("Main.value");
        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[64]);
        tag.serialize(buffer);
        assertNotNull(buffer.getBytes());
    }
}
