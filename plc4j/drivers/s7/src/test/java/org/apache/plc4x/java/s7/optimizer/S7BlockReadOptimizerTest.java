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

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.context.S7DriverContext;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagValueItem;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.junit.jupiter.api.Test;

import java.util.LinkedHashMap;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class S7BlockReadOptimizerTest {

    private final S7BlockReadOptimizer optimizer = new S7BlockReadOptimizer();
    private final S7DriverContext context = new S7DriverContext();

    private static PlcReadRequest req(String[][] tagSpecs) {
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        for (String[] spec : tagSpecs) {
            tags.put(spec[0], new DefaultPlcTagItem<>(S7Tag.of(spec[1])));
        }
        return new DefaultPlcReadRequest(null, tags);
    }

    @Test
    void emptyRequestReturnsEmpty() {
        PlcReadRequest r = new DefaultPlcReadRequest(null, new LinkedHashMap<>());
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, context);
        assertTrue(chunks.isEmpty());
    }

    @Test
    void singleSmallReadReturnsOneChunk() {
        PlcReadRequest r = req(new String[][] {{"a", "%DB1.DBW0:INT"}});
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, context);
        assertEquals(1, chunks.size());
        assertEquals(1, chunks.get(0).slots().size());
        assertEquals("a", chunks.get(0).slots().get(0).bindings().get(0).tagName());
    }

    @Test
    void manySmallReadsWithTinyPduSplitsAcrossChunks() {
        S7DriverContext tiny = new S7DriverContext();
        // Just large enough for the empty request/response headers plus one item.
        tiny.setPduSize(S7Optimizer.EMPTY_READ_REQUEST_SIZE + S7Optimizer.S7_ADDRESS_ANY_SIZE + 1);
        // Six tags, each in a different DB so they cannot block-merge.
        PlcReadRequest r = req(new String[][] {
            {"a", "%DB1.DBW0:INT"},
            {"b", "%DB2.DBW0:INT"},
            {"c", "%DB3.DBW0:INT"},
            {"d", "%DB4.DBW0:INT"},
            {"e", "%DB5.DBW0:INT"},
            {"f", "%DB6.DBW0:INT"},
        });
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, tiny);
        assertEquals(6, chunks.size(), "expected one chunk per tag");
        for (S7ReadChunk c : chunks) {
            assertEquals(1, c.slots().size());
        }
    }

    @Test
    void adjacentTagsInSameDbAreBlockMerged() {
        // Two ints in DB1, byte-offsets 0 and 2: adjacent and same area -> merged into one block.
        PlcReadRequest r = req(new String[][] {
            {"a", "%DB1.DBW0:INT"},
            {"b", "%DB1.DBW2:INT"},
        });
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, context);
        assertEquals(1, chunks.size());
        // The merged chunk has exactly one slot fetching the byte-block, but two bindings.
        S7ReadChunk.Slot slot = chunks.get(0).slots().get(0);
        assertEquals(2, slot.bindings().size());
        assertEquals("a", slot.bindings().get(0).tagName());
        assertEquals("b", slot.bindings().get(1).tagName());
        assertEquals(0, slot.bindings().get(0).payloadByteOffset());
        assertEquals(2, slot.bindings().get(1).payloadByteOffset());
    }

    @Test
    void distantTagsAreNotMerged() {
        // Two ints far apart in DB1 (offset 0 and offset 100) -> separate slots.
        PlcReadRequest r = req(new String[][] {
            {"a", "%DB1.DBW0:INT"},
            {"b", "%DB1.DBW100:INT"},
        });
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, context);
        // Both fit in one PDU so one chunk; but two slots (no merge).
        assertEquals(1, chunks.size());
        assertEquals(2, chunks.get(0).slots().size());
    }

    @Test
    void singleReadLargerThanPduIsSplit() {
        S7DriverContext tiny = new S7DriverContext();
        tiny.setPduSize(64);
        // A 200-byte array is larger than the PDU's available payload so it must be split.
        PlcReadRequest r = req(new String[][] {{"big", "%DB1.DBB0:BYTE[200]"}});
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, tiny);
        assertTrue(chunks.size() > 1, "expected multiple chunks, got " + chunks.size());
        // Each chunk's binding flags as a split fragment of the same original tag.
        int totalElems = 0;
        for (S7ReadChunk c : chunks) {
            for (S7ReadChunk.Slot s : c.slots()) {
                for (S7ReadChunk.Binding b : s.bindings()) {
                    assertEquals("big", b.tagName());
                    assertTrue(b.isSplitFragment());
                    totalElems += s.fragmentTag().getNumberOfElements();
                }
            }
        }
        assertEquals(200, totalElems);
    }

    @Test
    void splitWriteFitsAllInOneChunk() {
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put("a", new DefaultPlcTagValueItem<>(S7Tag.of("%DB1.DBW0:INT"), new PlcINT((short) 1)));
        tags.put("b", new DefaultPlcTagValueItem<>(S7Tag.of("%DB1.DBW2:INT"), new PlcINT((short) 2)));
        PlcWriteRequest req = new DefaultPlcWriteRequest(null, tags);
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(req, context);
        assertEquals(1, chunks.size());
        assertEquals(java.util.List.of("a", "b"), chunks.get(0).tagNames());
    }

    @Test
    void splitWriteSplitsAcrossPdus() {
        S7DriverContext tiny = new S7DriverContext();
        tiny.setPduSize(S7Optimizer.EMPTY_WRITE_REQUEST_SIZE + S7Optimizer.S7_ADDRESS_ANY_SIZE + 6 + 1);
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        for (int i = 0; i < 4; i++) {
            tags.put("t" + i, new DefaultPlcTagValueItem<>(
                S7Tag.of("%DB" + (i + 1) + ".DBW0:INT"), new PlcINT((short) i)));
        }
        PlcWriteRequest req = new DefaultPlcWriteRequest(null, tags);
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(req, tiny);
        assertEquals(4, chunks.size());
    }

    @Test
    void differentMemoryAreasNotMerged() {
        PlcReadRequest r = req(new String[][] {
            {"db", "%DB1.DBW0:INT"},
            {"flag", "%MW0:INT"},
        });
        List<S7ReadChunk> chunks = optimizer.splitReadRequest(r, context);
        assertEquals(1, chunks.size());
        // Two separate slots — different memory areas.
        assertEquals(2, chunks.get(0).slots().size());
    }
}
