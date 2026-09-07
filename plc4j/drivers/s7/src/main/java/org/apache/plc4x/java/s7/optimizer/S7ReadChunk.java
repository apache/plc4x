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

import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItem;
import org.apache.plc4x.java.s7.tag.S7Tag;

import java.util.List;

/**
 * One PDU-sized worth of S7 read protocol items.
 * <p>Each {@code slots} entry is exactly one outgoing
 * {@link S7VarRequestParameterItem}; the connection sends one S7 message per chunk.
 * The slot's bindings tell the connection how to decode the corresponding response
 * payload back into one or more user-tag values.
 */
public record S7ReadChunk(List<Slot> slots) {

    /**
     * One protocol-level read item plus the metadata needed to decode and merge
     * its response into the final {@code PlcReadResponse}.
     */
    public record Slot(
        S7VarRequestParameterItem requestItem,
        // The tag describing the actual bytes being requested.
        // For a block-read, this is a synthetic "block" tag; for a split-array tag
        // it is a fragment with a reduced numberOfElements; otherwise it is the
        // user-supplied tag.
        S7Tag fragmentTag,
        List<Binding> bindings) {}

    /**
     * How a portion of one slot's response payload feeds into a final user tag.
     */
    public record Binding(
        String tagName,
        S7Tag originalTag,
        // Byte offset *within* the slot's response payload at which this user
        // tag's data starts. 0 for non-block-merged single-tag slots.
        int payloadByteOffset,
        // Element index inside {@code originalTag} where this slot's data starts.
        // Non-zero only for split-array tags.
        int elementOffset,
        // True if {@code originalTag} was split across multiple slots/chunks.
        boolean isSplitFragment) {}
}
