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
            S7Tag.of("%DB1.DBB0:BYTE[8]"), new PlcRawByteArray(new byte[]{1, 2, 3, 4, 5, 6, 7, 8})));
        PlcWriteRequest req = new DefaultPlcWriteRequest(null, tags);
        List<S7WriteChunk> chunks = optimizer.splitWriteRequest(req, context);
        assertEquals(1, chunks.size());
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
        assertEquals(1, b.elementOffset());
        assertTrue(b.isSplitFragment());

        S7WriteChunk wc = new S7WriteChunk(java.util.List.of("x"));
        assertEquals(1, wc.tagNames().size());
    }
}
