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
package org.apache.plc4x.java.spi.drivers.model;

/**
 * What a protocol can actually encode of an array selection. The notation is the same for every
 * driver, but the wire formats are not: an EtherNet/IP array index travels in a CIP MemberID
 * whose instance field is a uint 8, and a driver addressing linear memory has no second dimension
 * to express. A driver states its limits here, and a selection that exceeds them is reported when
 * the address is parsed rather than truncated when it is serialized.
 *
 * @param maxIndex the largest start offset the protocol can encode. It bounds where a
 *        selection begins, not where it ends: a protocol carrying a start index and an
 *        element count can read past this bound, it just cannot start past it
 * @param maxDimensions how many dimensions the wire format carries
 * @param onlyTrailingDimensionMayBeRange whether every dimension but the last must be a single
 *        element - true where one request carries a single element count for the whole address
 */
public record AddressConstraints(int maxIndex, int maxDimensions, boolean onlyTrailingDimensionMayBeRange) {

    /** No limits beyond the grammar itself. */
    public static final AddressConstraints UNCONSTRAINED =
        new AddressConstraints(Integer.MAX_VALUE, Integer.MAX_VALUE, false);

    /** A protocol addressing linear memory: any index, but only one dimension. */
    public static final AddressConstraints SINGLE_DIMENSION =
        new AddressConstraints(Integer.MAX_VALUE, 1, false);

    public AddressConstraints {
        if (maxIndex < 0) {
            throw new IllegalArgumentException("maxIndex must not be negative");
        }
        if (maxDimensions < 1) {
            throw new IllegalArgumentException("maxDimensions must be at least 1");
        }
    }

    public AddressConstraints withMaxIndex(int maxIndex) {
        return new AddressConstraints(maxIndex, maxDimensions, onlyTrailingDimensionMayBeRange);
    }

    public AddressConstraints withMaxDimensions(int maxDimensions) {
        return new AddressConstraints(maxIndex, maxDimensions, onlyTrailingDimensionMayBeRange);
    }

    public AddressConstraints withOnlyTrailingDimensionMayBeRange(boolean onlyTrailing) {
        return new AddressConstraints(maxIndex, maxDimensions, onlyTrailing);
    }
}
