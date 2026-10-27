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
package org.apache.plc4x.java.s7.optimizer;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.s7.context.S7DriverContext;
import org.apache.plc4x.java.s7.tag.S7StringFixedLengthTag;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagValueItem;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.junit.jupiter.api.Test;

import java.util.LinkedHashMap;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/** Direct tests for the base optimizer (no block-merging). */
class S7OptimizerTest {

    private final S7Optimizer optimizer = new S7Optimizer();
    private final S7DriverContext context = new S7DriverContext();

    @Test
    void splitReadOneTag() {
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put("a", new DefaultPlcTagItem<>(S7Tag.of("%DB1.DBW0:INT")));
        PlcReadRequest req = new DefaultPlcReadRequest(null, tags);
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(req, context);
        assertEquals(1, chunks.size());
        S7ReadChunk.Slot slot = chunks.get(0).slots().get(0);
        assertEquals("a", slot.bindings().get(0).tagName());
        assertFalse(slot.bindings().get(0).isSplitFragment());
    }

    @Test
    void splitReadFixedLengthString() {
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        S7StringFixedLengthTag stringTag = new S7StringFixedLengthTag(
            TransportSize.STRING, MemoryArea.DATA_BLOCKS, 1, 0, (byte) 0, 1, 80);
        tags.put("s", new DefaultPlcTagItem<>(stringTag));
        PlcReadRequest req = new DefaultPlcReadRequest(null, tags);
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(req, context);
        assertEquals(1, chunks.size());
    }

    @Test
    void splitReadRejectsNonS7Tag() {
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        // A bare PlcTag implementation that the optimizer should reject.
        tags.put("x", new DefaultPlcTagItem<>(new PlcTag() {
            @Override public String getAddressString() { return ""; }
        }));
        PlcReadRequest req = new DefaultPlcReadRequest(null, tags);
        assertThrows(PlcRuntimeException.class, () -> optimizer.splitReadRequest(req, context));
    }

    @Test
    void splitWriteBoolAndString() {
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put("flag", new DefaultPlcTagValueItem<>(S7Tag.of("%M0.0:BOOL"), PlcBOOL.of(true)));
        tags.put("name", new DefaultPlcTagValueItem<>(
            new S7StringFixedLengthTag(TransportSize.STRING, MemoryArea.DATA_BLOCKS, 1, 0, (byte) 0, 1, 16),
            new PlcSTRING("hello")));
        tags.put("count", new DefaultPlcTagValueItem<>(S7Tag.of("%MW0:INT"), new PlcINT((short) 42)));
        PlcWriteRequest req = new DefaultPlcWriteRequest(null, tags);
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(req, context);
        assertEquals(1, chunks.size());
        assertEquals(3, chunks.get(0).tagNames().size());
    }

    @Test
    void splitWriteByteArrayCheckedSize() {
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put("buf", new DefaultPlcTagValueItem<>(
            S7Tag.of("%DB1.DBB0[0..7]:BYTE"), new PlcRawByteArray(new byte[]{1, 2, 3, 4, 5, 6, 7, 8})));
        PlcWriteRequest req = new DefaultPlcWriteRequest(null, tags);
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(req, context);
        assertEquals(1, chunks.size());
    }

    /**
     * A write of whole bytes overwrites every bit in them, so only the bytes a BOOL array covers
     * completely are written packed; each bit it shares a byte with goes out on its own.
     */
    @Test
    void boolArrayWriteItems() {
        assertEquals(List.of("16.0 x16 @0"), describe("%DB1:16.0[0..15]:BOOL"));
        assertEquals(List.of("16.0 x8 @0", "17.0 x1 @8", "17.1 x1 @9"), describe("%DB1:16.0[0..9]:BOOL"));
        assertEquals(List.of("16.6 x1 @0", "16.7 x1 @1", "17.0 x8 @2", "18.0 x1 @10", "18.1 x1 @11"),
            describe("%DB1:16.6[0..11]:BOOL"));
        assertEquals(List.of("16.3 x1 @0", "16.4 x1 @1", "16.5 x1 @2"), describe("%DB1:16.3[0..2]:BOOL"));
    }

    @Test
    void otherTagsAreWrittenAsOneItem() {
        S7Tag tag = S7Tag.of("%DB1.DBW0[0..3]:INT");
        assertEquals(List.of(new S7WriteChunk.Item("t", tag, 0)), S7Optimizer.writeItemsFor("t", tag));
        S7Tag bool = S7Tag.of("%DB1:16.3:BOOL");
        assertEquals(List.of(new S7WriteChunk.Item("t", bool, 0)), S7Optimizer.writeItemsFor("t", bool));
    }

    @Test
    void boolArrayWriteItemsSplitAcrossPdus() {
        // 7 leading bits, one whole byte, 7 trailing bits: 15 items, more than one 240 byte PDU holds.
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put("bools", new DefaultPlcTagValueItem<>(S7Tag.of("%DB1:16.1[0..21]:BOOL"), null));
        S7DriverContext pdu240 = new S7DriverContext();
        pdu240.setPduSize(240);
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(new DefaultPlcWriteRequest(null, tags), pdu240);
        assertTrue(chunks.size() > 1);
        assertEquals(15, chunks.stream().mapToInt(c -> c.items().size()).sum());
        chunks.forEach(c -> assertEquals(List.of("bools"), c.tagNames()));
    }

    private static List<String> describe(String address) {
        return S7Optimizer.writeItemsFor("t", S7Tag.of(address)).stream()
            .map(i -> i.tag().getByteOffset() + "." + i.tag().getBitOffset() + " x"
                + i.tag().getNumberOfElements() + " @" + i.firstElement())
            .toList();
    }

    @Test
    void emptyWriteRequestProducesNoChunks() {
        PlcWriteRequest req = new DefaultPlcWriteRequest(null, new LinkedHashMap<>());
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(req, context);
        assertTrue(chunks.isEmpty());
    }

    @Test
    void chunkRecordAccessors() {
        S7ReadChunk.Binding b = new S7ReadChunk.Binding("name", S7Tag.of("%MW0:INT"), 4, 1, true);
        assertEquals("name", b.tagName());
        assertEquals(4, b.payloadByteOffset());
        assertEquals(1, b.fragmentByteOffset());
        assertTrue(b.isSplitFragment());

        S7WriteChunk wc = new S7WriteChunk(java.util.List.of(
            new S7WriteChunk.Item("x", S7Tag.of("%MW0:INT"), 0), new S7WriteChunk.Item("x", S7Tag.of("%MW2:INT"), 1)));
        assertEquals(1, wc.tagNames().size());
    }
}
