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
package org.apache.plc4x.java.ads.readwrite;

import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Flag layout tests for the ADS symbol table (GH-2773).
 * <p>
 * The 16 bit symbol flags are defined by {@code ADSSYMBOLFLAG_*} in Beckhoff's {@code TcAdsDef.h}.
 * Bits 11..8 are no flags at all but a 4 bit task/context id, so a symbol living in a non-zero
 * context must neither be mistaken for one carrying attributes (bit 12) nor for one followed by a
 * 32 bit context mask. Either misreading makes the parser run past the end of the entry into the
 * next one.
 */
class AdsSymbolTableEntryFlagsTest {

    private static final int ADSSYMBOLFLAG_ATTRIBUTES = 0x1000;
    private static final int ADSSYMBOLFLAG_STATIC = 0x2000;
    private static final int ADSSYMBOLFLAG_INITONRESET = 0x4000;
    private static final int ADSSYMBOLFLAG_EXTENDEDFLAGS = 0x8000;

    @Test
    @DisplayName("the context id is read from bits 11..8 and does not trigger the attribute parser")
    void contextIdIsNotMistakenForAttributes() throws Exception {
        // 0x0200 is context 2: the exact value that used to set flagAttributes.
        byte[] first = entry("MAIN.a", 0x0200, null);
        byte[] second = entry("MAIN.b", 0x0000, null);
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(concat(first, second));

        AdsSymbolTableEntry a = AdsSymbolTableEntry.staticParse(readBuffer);
        assertEquals("MAIN.a", a.getName());
        assertEquals(2, a.getContextMask());
        assertNull(a.getAttributes());
        assertFalse(a.getFlagStatic());
        assertFalse(a.getFlagInitOnReset());
        assertFalse(a.getFlagExtendedFlags());
        assertEquals(first.length * 8, readBuffer.getPositionInBits(), "the first entry must end where its entryLength says");

        AdsSymbolTableEntry b = AdsSymbolTableEntry.staticParse(readBuffer);
        assertEquals("MAIN.b", b.getName());
        assertEquals(0, b.getContextMask());
    }

    @Test
    @DisplayName("every context id 0..15 round-trips")
    void allContextIdsRoundTrip() throws Exception {
        for (int context = 0; context < 16; context++) {
            byte[] encoded = entry("MAIN.x", context << 8, null);
            AdsSymbolTableEntry parsed = AdsSymbolTableEntry.staticParse(new ReadBufferByteBased(encoded));
            assertEquals(context, parsed.getContextMask(), "context " + context);
            assertNull(parsed.getAttributes(), "context " + context);
            assertArrayEquals(encoded, serialize(parsed), "context " + context);
        }
    }

    @Test
    @DisplayName("bits 15..12 map to the TcAdsDef.h flags")
    void upperFlagBitsMatchTcAdsDef() throws Exception {
        AdsSymbolTableEntry parsed = AdsSymbolTableEntry.staticParse(new ReadBufferByteBased(
            entry("MAIN.x", ADSSYMBOLFLAG_STATIC, null)));
        assertTrue(parsed.getFlagStatic());
        assertFalse(parsed.getFlagInitOnReset());
        assertFalse(parsed.getFlagExtendedFlags());

        parsed = AdsSymbolTableEntry.staticParse(new ReadBufferByteBased(
            entry("MAIN.x", ADSSYMBOLFLAG_INITONRESET, null)));
        assertFalse(parsed.getFlagStatic());
        assertTrue(parsed.getFlagInitOnReset());
        assertFalse(parsed.getFlagExtendedFlags());

        parsed = AdsSymbolTableEntry.staticParse(new ReadBufferByteBased(
            entry("MAIN.x", ADSSYMBOLFLAG_EXTENDEDFLAGS, null)));
        assertFalse(parsed.getFlagStatic());
        assertFalse(parsed.getFlagInitOnReset());
        assertTrue(parsed.getFlagExtendedFlags());
    }

    @Test
    @DisplayName("the attribute flag (bit 12) still selects the attribute block")
    void attributeFlagSelectsAttributes() throws Exception {
        byte[] encoded = entry("MAIN.x", ADSSYMBOLFLAG_ATTRIBUTES | 0x0300, attribute("TcDisplayName", "x"));
        AdsSymbolTableEntry parsed = AdsSymbolTableEntry.staticParse(new ReadBufferByteBased(encoded));

        assertEquals(3, parsed.getContextMask());
        assertNotNull(parsed.getAttributes());
        assertEquals(1, parsed.getAttributes().getAttributes().size());
        assertEquals("TcDisplayName", parsed.getAttributes().getAttributes().get(0).getName());
        assertArrayEquals(encoded, serialize(parsed));
    }

    private static byte[] serialize(AdsSymbolTableEntry entry) throws Exception {
        byte[] bytes = new byte[entry.getLengthInBytes()];
        entry.serialize(new WriteBufferByteBased(bytes));
        return bytes;
    }

    /** One attribute block holding a single name/value pair, as sent when bit 12 is set. */
    private static byte[] attribute(String name, String value) {
        byte[] n = name.getBytes(StandardCharsets.UTF_8);
        byte[] v = value.getBytes(StandardCharsets.UTF_8);
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        writeUInt16(out, 1); // numAttributes
        out.write(n.length);
        out.write(v.length);
        out.write(n, 0, n.length);
        out.write(0);
        out.write(v, 0, v.length);
        out.write(0);
        return out.toByteArray();
    }

    private static byte[] entry(String name, int flags, byte[] attributes) {
        byte[] n = name.getBytes(StandardCharsets.UTF_8);
        byte[] type = "BOOL".getBytes(StandardCharsets.UTF_8);
        ByteArrayOutputStream body = new ByteArrayOutputStream();
        writeUInt32(body, 0x4040); // group
        writeUInt32(body, 0);      // offset
        writeUInt32(body, 1);      // size
        writeUInt32(body, 33);     // dataType: ADST_BIT
        writeUInt16(body, flags);
        writeUInt16(body, 0);      // second flag word
        writeUInt16(body, n.length);
        writeUInt16(body, type.length);
        writeUInt16(body, 0);      // commentLength
        body.write(n, 0, n.length);
        body.write(0);
        body.write(type, 0, type.length);
        body.write(0);
        body.write(0);             // empty comment
        if (attributes != null) {
            body.write(attributes, 0, attributes.length);
        }
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        writeUInt32(out, 4 + body.size()); // entryLength
        out.write(body.toByteArray(), 0, body.size());
        return out.toByteArray();
    }

    private static byte[] concat(byte[] a, byte[] b) {
        byte[] result = new byte[a.length + b.length];
        System.arraycopy(a, 0, result, 0, a.length);
        System.arraycopy(b, 0, result, a.length, b.length);
        return result;
    }

    private static void writeUInt32(ByteArrayOutputStream out, int value) {
        out.write(value & 0xFF);
        out.write((value >>> 8) & 0xFF);
        out.write((value >>> 16) & 0xFF);
        out.write((value >>> 24) & 0xFF);
    }

    private static void writeUInt16(ByteArrayOutputStream out, int value) {
        out.write(value & 0xFF);
        out.write((value >>> 8) & 0xFF);
    }
}
