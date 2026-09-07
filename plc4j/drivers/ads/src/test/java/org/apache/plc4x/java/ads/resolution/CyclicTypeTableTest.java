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
package org.apache.plc4x.java.ads.resolution;

import org.apache.plc4x.java.ads.readwrite.AdsDataTypeTableEntry;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The data type table is uploaded from the device, and a field's type is a name looked up in it. So
 * a struct can have a field of its own type - and decoding that follows names round a flat map,
 * which the bound on how deeply the table may nest when parsed does not touch.
 */
class CyclicTypeTableTest {

    private final Map<String, AdsDataTypeTableEntry> table = new HashMap<>();

    /** A struct entry whose single child is of type {@code childTypeName}. */
    private AdsDataTypeTableEntry struct(String name, String childTypeName) {
        AdsDataTypeTableEntry child = Mockito.mock(AdsDataTypeTableEntry.class);
        Mockito.when(child.getMainName()).thenReturn("field");
        Mockito.when(child.getSecondaryName()).thenReturn(childTypeName);
        Mockito.when(child.getOffset()).thenReturn(0L);
        Mockito.when(child.getSize()).thenReturn(0L);
        Mockito.when(child.getChildren()).thenReturn(Collections.emptyList());
        Mockito.when(child.getArrayInfo()).thenReturn(Collections.emptyList());

        AdsDataTypeTableEntry entry = Mockito.mock(AdsDataTypeTableEntry.class);
        Mockito.when(entry.getMainName()).thenReturn(name);
        Mockito.when(entry.getSecondaryName()).thenReturn(name);
        Mockito.when(entry.getOffset()).thenReturn(0L);
        Mockito.when(entry.getSize()).thenReturn(0L);
        Mockito.when(entry.getChildren()).thenReturn(List.of(child));
        Mockito.when(entry.getArrayInfo()).thenReturn(Collections.emptyList());
        table.put(name, entry);
        return entry;
    }

    private static ReadBufferByteBased buffer() {
        return new ReadBufferByteBased(new byte[64]);
    }

    @Test
    void aStructWhoseFieldIsItsOwnTypeIsReportedNotFollowed() {
        // The shape the report named: offset and size zero, so neither of the existing guards
        // turns the child away, and its type is the struct it lives in.
        AdsDataTypeTableEntry selfReferential = struct("ST_Loop", "ST_Loop");
        ValueDecoder decoder = new ValueDecoder(table);

        BufferException thrown = assertThrows(BufferException.class,
            () -> decoder.decodeStructForTest(buffer(), selfReferential));
        assertTrue(thrown.getMessage().contains("ST_Loop"),
            "the failure should name the type that loops, but was: " + thrown.getMessage());
    }

    @Test
    void aLongerLoopIsAlsoReported() {
        AdsDataTypeTableEntry a = struct("ST_A", "ST_B");
        struct("ST_B", "ST_C");
        struct("ST_C", "ST_A");
        ValueDecoder decoder = new ValueDecoder(table);

        assertThrows(BufferException.class, () -> decoder.decodeStructForTest(buffer(), a));
    }
}
