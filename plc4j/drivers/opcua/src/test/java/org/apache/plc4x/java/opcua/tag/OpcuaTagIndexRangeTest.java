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
package org.apache.plc4x.java.opcua.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

/**
 * Tests the array-index / IndexRange syntax on OPC UA tag addresses: {@code [n]}, {@code [lo..hi]}
 * (inclusive), an optional {@code ;base} lower bound (default 0), and the multi-dimensional form
 * {@code [..][..]}. The resolved value is an OPC UA IndexRange string (0-based, inclusive, comma
 * separated per dimension).
 */
class OpcuaTagIndexRangeTest {

    @Test
    void singleIndexIsZeroBased() {
        OpcuaTag tag = OpcuaTag.of("ns=3;s=Test/Array/Int[8]");
        // The index suffix is split off the identifier; the node address itself is unchanged.
        assertThat(tag.getIdentifier()).isEqualTo("Test/Array/Int");
        assertThat(tag.getIndexRange()).isEqualTo("8");
        // The index suffix survives in the address string (";NULL" is the untyped-tag type tail),
        // and re-parsing the address yields an equivalent tag.
        assertThat(tag.getAddressString()).startsWith("ns=3;s=Test/Array/Int[8]");
        OpcuaTag reparsed = OpcuaTag.of(tag.getAddressString());
        assertThat(reparsed.getIdentifier()).isEqualTo("Test/Array/Int");
        assertThat(reparsed.getIndexRange()).isEqualTo("8");
    }

    @Test
    void inclusiveRange() {
        // [3..8] -> indices 3..8 inclusive.
        assertThat(OpcuaTag.of("ns=3;s=Foo[3..8]").getIndexRange()).isEqualTo("3:8");
    }

    @Test
    void baseIsSubtracted() {
        // [3..8;1] -> base 1, so the 3rd..8th element -> 0-based OPC UA range 2:7.
        assertThat(OpcuaTag.of("ns=3;s=Foo[3..8;1]").getIndexRange()).isEqualTo("2:7");
        // Single index with a base.
        assertThat(OpcuaTag.of("ns=3;s=Foo[8;1]").getIndexRange()).isEqualTo("7");
    }

    @Test
    void multiDimensional() {
        assertThat(OpcuaTag.of("ns=3;s=Foo[1..2][0..1]").getIndexRange()).isEqualTo("1:2,0:1");
        assertThat(OpcuaTag.of("ns=3;s=Foo[1][2][3]").getIndexRange()).isEqualTo("1,2,3");
        // Per-dimension base.
        assertThat(OpcuaTag.of("ns=3;s=Foo[1..6;1][3]").getIndexRange()).isEqualTo("0:5,3");
    }

    @Test
    void noIndexLeavesTheTagUntouched() {
        OpcuaTag tag = OpcuaTag.of("ns=3;s=Test/Scalar/Bool");
        assertThat(tag.getIndexRange()).isNull();
        assertThat(tag.getIdentifier()).isEqualTo("Test/Scalar/Bool");
    }

    @Test
    void nonNumericBracketsAreNotAnIndex() {
        // A string identifier that legitimately contains '[...]' with non-numeric content is left
        // whole (the index grammar is strictly numeric), so such nodes remain addressable.
        OpcuaTag tag = OpcuaTag.of("ns=3;s=Foo[bar]");
        assertThat(tag.getIndexRange()).isNull();
        assertThat(tag.getIdentifier()).isEqualTo("Foo[bar]");
    }

    @Test
    void indexCombinesWithAttributeAndType() {
        OpcuaTag tag = OpcuaTag.of("ns=3;s=Foo[3..8];INT");
        assertThat(tag.getIdentifier()).isEqualTo("Foo");
        assertThat(tag.getIndexRange()).isEqualTo("3:8");
        assertThat(tag.getDataType().name()).isEqualTo("INT");
    }

    @Test
    void invalidRangesAreRejected() {
        // high < low
        assertThatThrownBy(() -> OpcuaTag.of("ns=3;s=Foo[8..3]"))
            .isInstanceOf(PlcInvalidTagException.class);
        // base pushes an index negative
        assertThatThrownBy(() -> OpcuaTag.of("ns=3;s=Foo[0..5;1]"))
            .isInstanceOf(PlcInvalidTagException.class);
    }
}
