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

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.function.Function;

public class MemoryChunkStorage implements ChunkStorage {

    private final List<byte[]> chunks = new ArrayList<>();
    private long size = 0;

    @Override
    public void append(byte[] frame) {
        chunks.add(frame);
        size += chunks.get(chunks.size() - 1).length;
    }

    public long size() {
        return size;
    }

    @Override
    public byte[] get() {
        Optional<byte[]> collect = chunks.stream().reduce((b1, b2) -> {
            byte[] combined = new byte[b1.length + b2.length];
            System.arraycopy(b1, 0, combined, 0, b1.length);
            System.arraycopy(b2, 0, combined, b1.length, b2.length);
            return combined;
        });
        return collect.orElse(new byte[0]);
    }

    @Override
    public void reset() {
        chunks.clear();
        size = 0;
    }


}
