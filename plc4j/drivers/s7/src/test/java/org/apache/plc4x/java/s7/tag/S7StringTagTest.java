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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7StringTagTest {

    @Test
    void parseFixedLengthString() {
        S7StringFixedLengthTag tag = S7StringFixedLengthTag.of("%DB1.DB0:STRING(80)");
        assertNotNull(tag);
        assertEquals(80, tag.getStringLength());
        assertEquals(TransportSize.STRING, tag.getDataType());
        assertEquals(1, tag.getBlockNumber());
        assertEquals(0, tag.getByteOffset());
    }

    @Test
    void parseFixedLengthStringWithCount() {
        S7StringFixedLengthTag tag = S7StringFixedLengthTag.of("%DB1.DB0:STRING(40)[3]");
        assertNotNull(tag);
        assertEquals(40, tag.getStringLength());
        assertEquals(3, tag.getNumberOfElements());
    }

    @Test
    void parseFixedLengthShortForm() {
        S7StringFixedLengthTag tag = S7StringFixedLengthTag.of("%DB1:0:STRING(40)");
        assertNotNull(tag);
        assertEquals(40, tag.getStringLength());
    }

    @Test
    void fixedLengthMatches() {
        assertTrue(S7StringFixedLengthTag.matches("%DB1.DB0:STRING(80)"));
        assertTrue(S7StringFixedLengthTag.matches("%DB1:0:STRING(80)"));
        assertFalse(S7StringFixedLengthTag.matches("%DB1.DBW0:INT"));
    }

    @Test
    void fixedLengthEqualityAndHashCode() {
        S7StringFixedLengthTag a = new S7StringFixedLengthTag(TransportSize.STRING, MemoryArea.DATA_BLOCKS, 1, 0, (byte) 0, 1, 80);
        S7StringFixedLengthTag b = new S7StringFixedLengthTag(TransportSize.STRING, MemoryArea.DATA_BLOCKS, 1, 0, (byte) 0, 1, 80);
        S7StringFixedLengthTag c = new S7StringFixedLengthTag(TransportSize.STRING, MemoryArea.DATA_BLOCKS, 1, 0, (byte) 0, 1, 40);
        assertEquals(a, b);
        assertEquals(a.hashCode(), b.hashCode());
        assertNotEquals(a, c);
    }

    @Test
    void fixedLengthToStringMentionsLength() {
        S7StringFixedLengthTag tag = new S7StringFixedLengthTag(TransportSize.STRING, MemoryArea.DATA_BLOCKS, 1, 0, (byte) 0, 1, 80);
        assertTrue(tag.toString().contains("80"));
    }

    @Test
    void parseVarLengthString() {
        S7StringVarLengthTag tag = S7StringVarLengthTag.of("%DB1.DB0:STRING");
        assertNotNull(tag);
        assertEquals(TransportSize.STRING, tag.getDataType());
    }

    @Test
    void varLengthMatches() {
        assertTrue(S7StringVarLengthTag.matches("%DB1.DB0:STRING"));
        assertTrue(S7StringVarLengthTag.matches("%DB1.DB0:WSTRING"));
        assertFalse(S7StringVarLengthTag.matches("%DB1.DBW0:INT"));
    }

    @Test
    void varLengthEqualityAndHashCode() {
        S7StringVarLengthTag a = S7StringVarLengthTag.of("%DB1.DB0:STRING");
        S7StringVarLengthTag b = S7StringVarLengthTag.of("%DB1.DB0:STRING");
        assertEquals(a, b);
        assertEquals(a.hashCode(), b.hashCode());
    }

    @Test
    void varLengthToStringContainsType() {
        S7StringVarLengthTag tag = S7StringVarLengthTag.of("%DB1.DB0:WSTRING");
        assertNotNull(tag);
        assertTrue(tag.toString().contains("WSTRING"));
    }

    @Test
    void varLengthShortForm() {
        S7StringVarLengthTag tag = S7StringVarLengthTag.of("%DB1:0:STRING");
        assertNotNull(tag);
        assertEquals(TransportSize.STRING, tag.getDataType());
    }

    /**
     * S7Tag.of() has to accept everything S7PlcTagHandler.parseTag() accepts, STRING lengths
     * included - see GH-2388. Otherwise an address that works via addTagAddress(...) blows up
     * when the very same string is handed to S7Tag.of(...).
     */
    @Test
    void s7TagOfParsesFixedLengthStrings() {
        S7Tag longForm = S7Tag.of("%DB1.DB0:STRING(80)");
        assertInstanceOf(S7StringFixedLengthTag.class, longForm);
        assertEquals(80, ((S7StringFixedLengthTag) longForm).getStringLength());

        // the reporter's address
        S7Tag shortForm = S7Tag.of("%DB69:68:STRING(20)");
        assertInstanceOf(S7StringFixedLengthTag.class, shortForm);
        assertEquals(20, ((S7StringFixedLengthTag) shortForm).getStringLength());
        assertEquals(69, shortForm.getBlockNumber());
        assertEquals(68, shortForm.getByteOffset());

        S7Tag withCount = S7Tag.of("%DB1.DB0:WSTRING(40)[3]");
        assertInstanceOf(S7StringFixedLengthTag.class, withCount);
        assertEquals(40, ((S7StringFixedLengthTag) withCount).getStringLength());
        assertEquals(3, withCount.getNumberOfElements());
    }

    @Test
    void s7TagMatchesAcceptsStringLengths() {
        assertTrue(S7Tag.matches("%DB69:68:STRING(20)"));
        assertTrue(S7Tag.matches("%DB1.DB0:WSTRING(40)[3]"));
        assertFalse(S7Tag.matches("not-a-tag"));
    }

    /**
     * Both entry points must return equal tags for the same address.
     */
    @Test
    void s7TagOfAgreesWithTagHandler() {
        S7PlcTagHandler handler = new S7PlcTagHandler();
        for (String address : new String[]{
            "%DB69:68:STRING(20)",
            "%DB1.DB0:STRING(80)",
            "%DB1.DB0:STRING(40)[3]",
            "%DB1.DB0:STRING",
            "%DB1:0:STRING",
            "%MW0:INT",
            "%DB1.DBX0.0:BOOL"}) {
            assertEquals(handler.parseTag(address), S7Tag.of(address), address);
        }
    }

    @Test
    void s7TagOfStillRejectsGarbage() {
        assertThrows(org.apache.plc4x.java.api.exceptions.PlcInvalidTagException.class,
            () -> S7Tag.of("not-a-tag"));
    }

    @Test
    void tagHandlerRoutesToCorrectTagClass() {
        S7PlcTagHandler handler = new S7PlcTagHandler();
        assertInstanceOf(S7StringFixedLengthTag.class, handler.parseTag("%DB1.DB0:STRING(80)"));
        assertInstanceOf(S7StringVarLengthTag.class, handler.parseTag("%DB1.DB0:STRING"));
        assertInstanceOf(S7Tag.class, handler.parseTag("%MW0:INT"));
        assertThrows(org.apache.plc4x.java.api.exceptions.PlcInvalidTagException.class,
            () -> handler.parseTag("not-a-tag"));
        // parseQuery returns null — there's no S7 query language; the connection's
        // onBrowse emits everything it knows regardless of the supplied query string.
        assertNull(handler.parseQuery("any"));
    }
}
