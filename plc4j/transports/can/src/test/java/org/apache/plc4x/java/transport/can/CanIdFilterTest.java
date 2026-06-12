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
package org.apache.plc4x.java.transport.can;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class CanIdFilterTest {

    @Test
    void acceptAllMatchesEverything() {
        CanIdFilter filter = CanIdFilter.acceptAll();

        assertTrue(filter.matches(0x000));
        assertTrue(filter.matches(0x123));
        assertTrue(filter.matches(0x7FF));
        assertTrue(filter.matches(0x1FFFFFFF));
        assertTrue(filter.isEmpty());
    }

    @Test
    void acceptIdsMatchesListedIdsOnly() {
        CanIdFilter filter = CanIdFilter.acceptIds(0x100, 0x200, 0x300);

        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x200));
        assertTrue(filter.matches(0x300));
        assertFalse(filter.matches(0x150));
        assertFalse(filter.matches(0x000));
        assertFalse(filter.isEmpty());
    }

    @Test
    void acceptRangeMatchesInclusiveRange() {
        CanIdFilter filter = CanIdFilter.acceptRange(0x100, 0x1FF);

        assertTrue(filter.matches(0x100));  // start boundary
        assertTrue(filter.matches(0x150));  // middle
        assertTrue(filter.matches(0x1FF));  // end boundary
        assertFalse(filter.matches(0x0FF)); // below range
        assertFalse(filter.matches(0x200)); // above range
        assertFalse(filter.isEmpty());
    }

    @Test
    void combinedIdsAndRanges() {
        CanIdFilter filter = CanIdFilter.builder()
                .addId(0x000)
                .addId(0x7FF)
                .addRange(0x100, 0x1FF)
                .build();

        assertTrue(filter.matches(0x000));  // explicit ID
        assertTrue(filter.matches(0x7FF));  // explicit ID
        assertTrue(filter.matches(0x150));  // in range
        assertFalse(filter.matches(0x050)); // not in IDs or range
        assertFalse(filter.matches(0x200)); // above range, not in IDs
        assertFalse(filter.isEmpty());
    }

    @Test
    void emptyFilterMatchesAll() {
        CanIdFilter filter = CanIdFilter.builder().build();

        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x000));
        assertTrue(filter.isEmpty());
    }

    @Test
    void outOfRangeIdRejectedOnAcceptIds() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.acceptIds(-1));
    }

    @Test
    void outOfRangeIdRejectedOnBuilder() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.builder().addId(-1));
    }

    @Test
    void idExceedingMaxRejected() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.acceptIds(0x20000000));
    }

    @Test
    void startGreaterThanEndRejected() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.acceptRange(0x200, 0x100));
    }

    @Test
    void negativeRangeStartRejected() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.acceptRange(-1, 0x100));
    }

    @Test
    void negativeRangeEndRejected() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.acceptRange(0, -1));
    }

    @Test
    void rangeEndExceedingMaxRejected() {
        assertThrows(IllegalArgumentException.class, () ->
                CanIdFilter.acceptRange(0, 0x20000000));
    }

    @Test
    void isEmptyBehavior() {
        assertTrue(CanIdFilter.acceptAll().isEmpty());
        assertTrue(CanIdFilter.builder().build().isEmpty());

        assertFalse(CanIdFilter.acceptIds(0x100).isEmpty());
        assertFalse(CanIdFilter.acceptRange(0x100, 0x200).isEmpty());
    }

    @Test
    void getAcceptedIds() {
        CanIdFilter filter = CanIdFilter.acceptIds(0x100, 0x200);

        assertEquals(2, filter.getAcceptedIds().size());
        assertTrue(filter.getAcceptedIds().contains(0x100));
        assertTrue(filter.getAcceptedIds().contains(0x200));
    }

    @Test
    void getAcceptedRanges() {
        CanIdFilter filter = CanIdFilter.acceptRange(0x100, 0x1FF);

        assertEquals(1, filter.getAcceptedRanges().size());
        assertEquals(0x100, filter.getAcceptedRanges().getFirst().start());
        assertEquals(0x1FF, filter.getAcceptedRanges().getFirst().end());
    }

    @Test
    void acceptedIdsUnmodifiable() {
        CanIdFilter filter = CanIdFilter.acceptIds(0x100);

        assertThrows(UnsupportedOperationException.class, () ->
                filter.getAcceptedIds().add(0x200));
    }

    @Test
    void acceptedRangesUnmodifiable() {
        CanIdFilter filter = CanIdFilter.acceptRange(0x100, 0x200);

        assertThrows(UnsupportedOperationException.class, () ->
                filter.getAcceptedRanges().add(new CanIdFilter.IdRange(0x300, 0x400)));
    }

    @Test
    void singlePointRange() {
        CanIdFilter filter = CanIdFilter.acceptRange(0x100, 0x100);

        assertTrue(filter.matches(0x100));
        assertFalse(filter.matches(0x0FF));
        assertFalse(filter.matches(0x101));
    }

    @Test
    void multipleRanges() {
        CanIdFilter filter = CanIdFilter.builder()
                .addRange(0x100, 0x1FF)
                .addRange(0x300, 0x3FF)
                .build();

        assertTrue(filter.matches(0x150));
        assertTrue(filter.matches(0x350));
        assertFalse(filter.matches(0x250));
    }

    @Test
    void toStringAcceptAll() {
        String str = CanIdFilter.acceptAll().toString();
        assertTrue(str.contains("acceptAll"));
    }

    @Test
    void toStringWithFilters() {
        String str = CanIdFilter.acceptIds(0x100).toString();
        assertFalse(str.contains("acceptAll"));
        assertTrue(str.contains("ids="));
    }

    @Test
    void idRangeToString() {
        CanIdFilter.IdRange range = new CanIdFilter.IdRange(0x100, 0x1FF);
        String str = range.toString();
        assertTrue(str.contains("100"));
        assertTrue(str.contains("1FF"));
    }

    @Test
    void zeroIdAccepted() {
        CanIdFilter filter = CanIdFilter.acceptIds(0);
        assertTrue(filter.matches(0));
        assertFalse(filter.matches(1));
    }
}
