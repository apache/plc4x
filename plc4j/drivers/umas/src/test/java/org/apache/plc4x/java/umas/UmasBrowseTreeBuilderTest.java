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
package org.apache.plc4x.java.umas;

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.umas.UmasConnection.BrowseTreeBuilder;
import org.apache.plc4x.java.umas.readwrite.UmasArrayDimension;
import org.apache.plc4x.java.umas.readwrite.UmasUDTDefinition;
import org.junit.jupiter.api.Test;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The data dictionary is the device's to write, so it can describe a type that contains itself.
 * Walking that without noticing does not come back.
 */
class UmasBrowseTreeBuilderTest {

    private final Map<Integer, List<UmasUDTDefinition>> fields = new HashMap<>();
    private final Map<Integer, Integer> elementTypeIds = new HashMap<>();
    private final Map<Integer, List<UmasArrayDimension>> dimensions = new HashMap<>();

    private static UmasUDTDefinition field(String name, int dataType) {
        return new UmasUDTDefinition(dataType, 0, 0, 0, name);
    }

    private BrowseTreeBuilder builder() {
        return new BrowseTreeBuilder(fields, elementTypeIds, dimensions);
    }

    @Test
    void aStructThatNamesItselfIsListedWithoutItsFields() {
        // Type 100 has a field of type 100.
        fields.put(100, List.of(field("self", 100)));

        PlcBrowseItem item = assertDoesNotThrow(() -> builder().buildBrowseItem("root", 100));
        assertNotNull(item);
        assertEquals(1, item.getChildren().size(), "the field itself is still listed");
        assertTrue(item.getChildren().get("self").getChildren().isEmpty(),
            "but it is not expanded again, because that is where the walk would not end");
    }

    /**
     * The case a depth-1 check would miss: the loop closes two types later.
     */
    @Test
    void aLongerLoopIsAlsoNoticed() {
        fields.put(200, List.of(field("toB", 201)));
        fields.put(201, List.of(field("toC", 202)));
        fields.put(202, List.of(field("backToA", 200)));

        PlcBrowseItem item = assertDoesNotThrow(() -> builder().buildBrowseItem("root", 200));
        PlcBrowseItem b = item.getChildren().get("toB");
        PlcBrowseItem c = b.getChildren().get("toC");
        PlcBrowseItem backToA = c.getChildren().get("backToA");
        assertTrue(backToA.getChildren().isEmpty(), "the type that closes the loop stops the walk");
    }

    /**
     * The reason the check is scoped to the path and not the whole browse: one type used by two
     * different fields is ordinary, and both should be expanded.
     */
    @Test
    void aTypeUsedTwiceInDifferentBranchesIsExpandedInBoth() {
        fields.put(300, List.of(field("left", 301), field("right", 301)));
        fields.put(301, List.of(field("leaf", 0x0A)));

        PlcBrowseItem item = builder().buildBrowseItem("root", 300);
        assertEquals(1, item.getChildren().get("left").getChildren().size(),
            "the left branch expands the shared type");
        assertEquals(1, item.getChildren().get("right").getChildren().size(),
            "and so does the right, because neither is inside the other");
    }

    @Test
    void anArrayWhoseElementTypeLoopsBackIsReportedAsRawBytes() {
        // Type 400 is an array of 400.
        elementTypeIds.put(400, 400);

        PlcBrowseItem item = assertDoesNotThrow(() -> builder().buildBrowseItem("arr", 400));
        assertNotNull(item);
    }

    @Test
    void anOrdinaryNestedStructIsStillWalkedInFull() {
        fields.put(500, List.of(field("inner", 501)));
        fields.put(501, List.of(field("value", 0x0A)));

        PlcBrowseItem item = builder().buildBrowseItem("root", 500);
        assertEquals(1, item.getChildren().size());
        assertEquals(1, item.getChildren().get("inner").getChildren().size());
    }
}
