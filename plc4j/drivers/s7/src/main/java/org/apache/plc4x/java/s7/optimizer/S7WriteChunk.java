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

import org.apache.plc4x.java.s7.tag.S7Tag;

import java.util.List;

/**
 * One PDU-sized worth of write protocol items. Writes are not block-merged; a user tag becomes
 * one item, except a BOOL array that shares a byte with something else, which is written in
 * several (see {@link S7Optimizer#writeItemsFor}). Each item's response status is mapped back to
 * its tag by name.
 */
public record S7WriteChunk(List<Item> items) {

    /**
     * One write protocol item.
     *
     * @param tagName      the user tag the item writes to
     * @param tag          the part of that tag the item writes - the tag itself unless it was split
     * @param firstElement the index of the tag's first element the item carries
     */
    public record Item(String tagName, S7Tag tag, int firstElement) {}

    /** The user tags this chunk writes to, each once, in order. */
    public List<String> tagNames() {
        return items.stream().map(Item::tagName).distinct().toList();
    }
}
