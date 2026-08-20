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

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Nesting tests for the ADS data type table.
 * <p>
 * A table entry may carry an array of further table entries, so the depth of the tree is decided
 * entirely by the device: the driver uploads the whole table at connect time and every byte of it
 * comes off the wire. One entry header costs {@value #ENTRY_HEADER_BYTES} bytes, so a table of a
 * few hundred kilobytes describes a tree deep enough to exhaust the parser's stack. Depth must
 * therefore be bounded by the parser rather than by the device's goodwill.
 */
class AdsDataTypeTableEntryNestingTest {

    /**
     * Fixed cost of one entry with empty names, no array dimensions, no guid and no optional
     * blocks: six 32 bit numbers, the 32 bit type id, four flag bytes, five 16 bit lengths and the
     * three string terminators.
     */
    private static final int ENTRY_HEADER_BYTES = 45;

    /**
     * Deep enough that an unbounded recursive-descent parser cannot survive it, while staying a
     * plausible upload size: 20000 * 45 bytes is under a megabyte, and the ADS length field is 32
     * bits wide so a real device could send far more.
     */
    private static final int OVERFLOW_DEPTH = 20_000;

    /** The budget the driver hands the parser for a real upload. */
    private static final int DRIVER_BUDGET = 255;

    @Test
    @DisplayName("a deeply nested table is rejected without exhausting the stack")
    void deeplyNestedTableIsRejectedNotFatal() {
        byte[] encoded = nestedEntries(OVERFLOW_DEPTH);

        assertThrows(BufferException.class,
            () -> AdsDataTypeTableEntry.staticParse(new ReadBufferByteBased(encoded), DRIVER_BUDGET),
            "a table nested " + OVERFLOW_DEPTH + " deep must fail as a parse error, not as a "
                + "StackOverflowError escaping into the caller");
    }

    @Test
    @DisplayName("nesting a real device would produce still parses")
    void modestNestingStillParses() throws Exception {
        AdsDataTypeTableEntry entry = AdsDataTypeTableEntry.staticParse(
            new ReadBufferByteBased(nestedEntries(8)), DRIVER_BUDGET);

        assertNotNull(entry);
        assertEquals(7, depthOf(entry), "all eight levels must have been parsed");
    }

    /**
     * Pins the boundary rather than just the extremes: the budget counts levels, so a table
     * nested exactly to the budget is accepted and one level more is not.
     */
    @Test
    @DisplayName("the budget admits exactly its own number of levels")
    void budgetAdmitsExactlyItsOwnNumberOfLevels() throws Exception {
        int budget = 8;

        AdsDataTypeTableEntry atLimit = AdsDataTypeTableEntry.staticParse(
            new ReadBufferByteBased(nestedEntries(budget)), budget);
        assertEquals(budget - 1, depthOf(atLimit),
            "a table nested exactly to the budget must parse in full");

        assertThrows(BufferException.class,
            () -> AdsDataTypeTableEntry.staticParse(
                new ReadBufferByteBased(nestedEntries(budget + 1)), budget),
            "one level past the budget must be rejected");
    }

    private static int depthOf(AdsDataTypeTableEntry entry) {
        int depth = 0;
        AdsDataTypeTableEntry cursor = entry;
        while (!cursor.getChildren().isEmpty()) {
            cursor = cursor.getChildren().get(0);
            depth++;
        }
        return depth;
    }

    /**
     * Builds {@code depth} entries nested one inside the next, innermost first. Every entry
     * declares an {@code entryLength} covering itself and everything below it, so the trailing
     * "rest" array consumes nothing.
     */
    private static byte[] nestedEntries(int depth) {
        byte[] current = entry(new byte[0], 0);
        for (int level = 1; level < depth; level++) {
            current = entry(current, 1);
        }
        return current;
    }

    private static byte[] entry(byte[] child, int numChildren) {
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        writeUInt32(out, ENTRY_HEADER_BYTES + child.length); // entryLength, covers the subtree
        writeUInt32(out, 0);                                 // version
        writeUInt32(out, 0);                                 // hashValue
        writeUInt32(out, 0);                                 // typeHashValue
        writeUInt32(out, 0);                                 // size
        writeUInt32(out, 0);                                 // offset
        writeUInt32(out, 0);                                 // dataType: ADST_VOID
        // Four flag bytes, all clear: no guid, no method infos, no attributes, no extended infos,
        // so none of the optional trailing blocks is present.
        out.write(0);
        out.write(0);
        out.write(0);
        out.write(0);
        writeUInt16(out, 0);           // mainNameLength
        writeUInt16(out, 0);           // secondaryNameLength
        writeUInt16(out, 0);           // commentLength
        writeUInt16(out, 0);           // arrayDimensions
        writeUInt16(out, numChildren); // numChildren
        out.write(0);                  // mainNameTerminator
        out.write(0);                  // secondaryNameTerminator
        out.write(0);                  // commentTerminator
        out.write(child, 0, child.length);
        return out.toByteArray();
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
