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
package org.apache.plc4x.java.opcua.protocol.chunk;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class MemoryChunkStorageTest {

    @Test
    void emptyStorageReturnsEmptyBuffer() {
        MemoryChunkStorage storage = new MemoryChunkStorage();
        assertThat(storage.size()).isZero();
        assertThat(storage.get()).isEmpty();
    }

    @Test
    void appendConcatenatesChunksInOrder() {
        MemoryChunkStorage storage = new MemoryChunkStorage();
        storage.append(new byte[]{1, 2, 3});
        storage.append(new byte[]{4, 5});
        storage.append(new byte[]{6});

        assertThat(storage.size()).isEqualTo(6);
        assertThat(storage.get()).containsExactly(1, 2, 3, 4, 5, 6);
    }

    @Test
    void singleChunkIsReturnedAsIs() {
        MemoryChunkStorage storage = new MemoryChunkStorage();
        storage.append(new byte[]{9, 9, 9});
        assertThat(storage.size()).isEqualTo(3);
        assertThat(storage.get()).containsExactly(9, 9, 9);
    }

    @Test
    void resetClearsContent() {
        MemoryChunkStorage storage = new MemoryChunkStorage();
        storage.append(new byte[]{1, 2, 3});
        storage.reset();
        assertThat(storage.size()).isZero();
        assertThat(storage.get()).isEmpty();
    }
}
