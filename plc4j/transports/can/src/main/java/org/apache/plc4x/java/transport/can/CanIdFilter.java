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

import java.util.*;

/**
 * Configurable filter for CAN frame identifiers.
 * <p>
 * Determines which CAN frame identifiers a transport instance accepts.
 * A frame matches if its identifier is in the explicit ID set OR falls within
 * any configured range. If no IDs or ranges are configured, all frames are accepted.
 * <p>
 * Instances are immutable and thread-safe.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Accept only specific IDs
 * CanIdFilter filter = CanIdFilter.acceptIds(0x100, 0x200, 0x300);
 *
 * // Accept a range
 * CanIdFilter rangeFilter = CanIdFilter.acceptRange(0x100, 0x1FF);
 *
 * // Combine IDs and ranges
 * CanIdFilter combined = CanIdFilter.builder()
 *     .addId(0x000)
 *     .addRange(0x100, 0x1FF)
 *     .build();
 * }</pre>
 */
public class CanIdFilter {

    private static final CanIdFilter ACCEPT_ALL = new CanIdFilter(Collections.emptySet(), Collections.emptyList());

    private final Set<Integer> acceptedIds;
    private final List<IdRange> acceptedRanges;

    private CanIdFilter(Set<Integer> acceptedIds, List<IdRange> acceptedRanges) {
        this.acceptedIds = Collections.unmodifiableSet(new HashSet<>(acceptedIds));
        this.acceptedRanges = Collections.unmodifiableList(new ArrayList<>(acceptedRanges));
    }

    /**
     * Returns a filter that accepts all CAN frame identifiers (no filtering).
     *
     * @return an accept-all filter
     */
    public static CanIdFilter acceptAll() {
        return ACCEPT_ALL;
    }

    /**
     * Returns a filter that accepts only the specified CAN identifiers.
     *
     * @param ids the accepted CAN identifiers
     * @return a filter accepting only the given IDs
     * @throws IllegalArgumentException if any ID is negative
     */
    public static CanIdFilter acceptIds(int... ids) {
        Set<Integer> idSet = new HashSet<>();
        for (int id : ids) {
            validateId(id);
            idSet.add(id);
        }
        return new CanIdFilter(idSet, Collections.emptyList());
    }

    /**
     * Returns a filter that accepts CAN identifiers within the given range (inclusive).
     *
     * @param start the start of the range (inclusive)
     * @param end   the end of the range (inclusive)
     * @return a filter accepting IDs in [start, end]
     * @throws IllegalArgumentException if start or end is negative, or start > end
     */
    public static CanIdFilter acceptRange(int start, int end) {
        return new CanIdFilter(Collections.emptySet(), List.of(new IdRange(start, end)));
    }

    /**
     * Creates a new builder for constructing filters with combined IDs and ranges.
     *
     * @return a new builder
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Tests whether the given CAN identifier matches this filter.
     * If the filter is empty (no IDs or ranges configured), all identifiers match.
     *
     * @param identifier the CAN identifier to test
     * @return true if the identifier matches the filter or the filter is empty
     */
    public boolean matches(int identifier) {
        if (isEmpty()) {
            return true;
        }
        if (acceptedIds.contains(identifier)) {
            return true;
        }
        for (IdRange range : acceptedRanges) {
            if (identifier >= range.start() && identifier <= range.end()) {
                return true;
            }
        }
        return false;
    }

    /**
     * Returns whether this filter is empty (accepts all frames).
     *
     * @return true if no IDs or ranges are configured
     */
    public boolean isEmpty() {
        return acceptedIds.isEmpty() && acceptedRanges.isEmpty();
    }

    /**
     * Returns the set of explicitly accepted CAN identifiers.
     *
     * @return unmodifiable set of accepted IDs
     */
    public Set<Integer> getAcceptedIds() {
        return acceptedIds;
    }

    /**
     * Returns the list of accepted CAN identifier ranges.
     *
     * @return unmodifiable list of accepted ranges
     */
    public List<IdRange> getAcceptedRanges() {
        return acceptedRanges;
    }

    private static void validateId(int id) {
        if (id < 0) {
            throw new IllegalArgumentException("CAN identifier must not be negative, got: " + id);
        }
        // We don't enforce max here because the filter might be used before we know
        // whether the frame will be standard or extended
        if (id > CanFrame.MAX_EXTENDED_ID) {
            throw new IllegalArgumentException(String.format(
                    "CAN identifier must not exceed 0x%X, got: 0x%X",
                    CanFrame.MAX_EXTENDED_ID, id));
        }
    }

    @Override
    public String toString() {
        if (isEmpty()) {
            return "CanIdFilter{acceptAll}";
        }
        return "CanIdFilter{ids=" + acceptedIds + ", ranges=" + acceptedRanges + "}";
    }

    /**
     * Represents an inclusive range of CAN identifiers [start, end].
     *
     * @param start the inclusive start of the range
     * @param end   the inclusive end of the range
     */
    public record IdRange(int start, int end) {
        /**
         * Constructs a validated ID range.
         *
         * @throws IllegalArgumentException if start or end is negative, end exceeds max, or start > end
         */
        public IdRange {
            if (start < 0) {
                throw new IllegalArgumentException("Range start must not be negative, got: " + start);
            }
            if (end < 0) {
                throw new IllegalArgumentException("Range end must not be negative, got: " + end);
            }
            if (end > CanFrame.MAX_EXTENDED_ID) {
                throw new IllegalArgumentException(String.format(
                        "Range end must not exceed 0x%X, got: 0x%X", CanFrame.MAX_EXTENDED_ID, end));
            }
            if (start > end) {
                throw new IllegalArgumentException(String.format(
                        "Range start (%d) must not exceed end (%d)", start, end));
            }
        }

        @Override
        public String toString() {
            return String.format("[0x%X–0x%X]", start, end);
        }
    }

    /**
     * Builder for constructing {@link CanIdFilter} instances with combined IDs and ranges.
     */
    public static class Builder {
        private final Set<Integer> ids = new HashSet<>();
        private final List<IdRange> ranges = new ArrayList<>();

        private Builder() {
        }

        /**
         * Adds a single CAN identifier to accept.
         *
         * @param id the CAN identifier
         * @return this builder
         * @throws IllegalArgumentException if id is negative or exceeds max
         */
        public Builder addId(int id) {
            validateId(id);
            ids.add(id);
            return this;
        }

        /**
         * Adds an inclusive range of CAN identifiers to accept.
         *
         * @param start the inclusive start of the range
         * @param end   the inclusive end of the range
         * @return this builder
         * @throws IllegalArgumentException if range is invalid
         */
        public Builder addRange(int start, int end) {
            ranges.add(new IdRange(start, end));
            return this;
        }

        /**
         * Builds the filter.
         *
         * @return the constructed filter
         */
        public CanIdFilter build() {
            return new CanIdFilter(ids, ranges);
        }
    }
}
