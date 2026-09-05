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
package org.apache.plc4x.java.ads;

import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.stream.IntStream;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link AdsTcpConnection#chunk(List, int)}, the helper that partitions the item
 * list of a sum command (multi-read/multi-write) into groups no larger than the
 * per-sum-command item cap.
 */
class AdsSumChunkTest {

    private static List<Integer> items(int count) {
        List<Integer> items = new ArrayList<>(count);
        IntStream.range(0, count).forEach(items::add);
        return items;
    }

    @Test
    void emptyListYieldsNoChunks() {
        assertTrue(AdsTcpConnection.chunk(Collections.emptyList(), 500).isEmpty());
    }

    @Test
    void singleItemYieldsOneChunk() {
        List<List<Integer>> chunks = AdsTcpConnection.chunk(items(1), 500);
        assertEquals(1, chunks.size());
        assertEquals(List.of(0), chunks.get(0));
    }

    @Test
    void listSmallerThanCapYieldsOneChunk() {
        List<List<Integer>> chunks = AdsTcpConnection.chunk(items(499), 500);
        assertEquals(1, chunks.size());
        assertEquals(499, chunks.get(0).size());
    }

    @Test
    void listEqualToCapYieldsOneChunk() {
        List<List<Integer>> chunks = AdsTcpConnection.chunk(items(500), 500);
        assertEquals(1, chunks.size());
        assertEquals(500, chunks.get(0).size());
    }

    @Test
    void listOneOverCapYieldsTwoChunks() {
        List<List<Integer>> chunks = AdsTcpConnection.chunk(items(501), 500);
        assertEquals(2, chunks.size());
        assertEquals(500, chunks.get(0).size());
        assertEquals(1, chunks.get(1).size());
    }

    @Test
    void largeListYieldsFullChunksPlusRemainder() {
        List<List<Integer>> chunks = AdsTcpConnection.chunk(items(1001), 500);
        assertEquals(3, chunks.size());
        assertEquals(500, chunks.get(0).size());
        assertEquals(500, chunks.get(1).size());
        assertEquals(1, chunks.get(2).size());
    }

    @Test
    void orderIsPreservedAcrossChunks() {
        List<Integer> input = items(1001);
        List<List<Integer>> chunks = AdsTcpConnection.chunk(input, 500);
        List<Integer> flattened = new ArrayList<>(input.size());
        chunks.forEach(flattened::addAll);
        assertEquals(input, flattened);
    }

    @Test
    void chunkSizeSmallerThanCapWorks() {
        List<List<Integer>> chunks = AdsTcpConnection.chunk(items(7), 3);
        assertEquals(3, chunks.size());
        assertEquals(List.of(0, 1, 2), chunks.get(0));
        assertEquals(List.of(3, 4, 5), chunks.get(1));
        assertEquals(List.of(6), chunks.get(2));
    }

    @Test
    void nonPositiveCapIsRejected() {
        assertThrows(IllegalArgumentException.class, () -> AdsTcpConnection.chunk(items(1), 0));
        assertThrows(IllegalArgumentException.class, () -> AdsTcpConnection.chunk(items(1), -1));
    }
}
