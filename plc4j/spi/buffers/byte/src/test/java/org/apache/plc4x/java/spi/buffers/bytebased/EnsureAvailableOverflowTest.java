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
package org.apache.plc4x.java.spi.buffers.bytebased;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Bounds tests for the availability check that guards every read.
 * <p>
 * {@code ensureAvailable} is the only thing standing between a field width taken off the wire and
 * the allocation {@code readBits} performs for it. Protocols carrying 32 bit length fields can hand
 * it a bit count close to {@link Integer#MAX_VALUE}; the check must stay correct there rather than
 * losing its meaning to integer overflow, and it must reject the value regardless of how large the
 * request is relative to the buffer.
 */
class EnsureAvailableOverflowTest {

    /**
     * A bit count large enough that {@code positionInBits + bitsNeeded} exceeds
     * {@link Integer#MAX_VALUE} and wraps to a negative number. Any check written as a sum rather
     * than as a comparison against the remaining bits silently accepts these values.
     */
    @Test
    @DisplayName("a bit count that overflows the position sum is still rejected")
    void rejectsBitCountThatOverflowsThePositionSum() throws Exception {
        // One byte consumed, so positionInBits is 8 and the sum below wraps negative.
        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[16]);
        buffer.readBits(8);

        int overflowing = Integer.MAX_VALUE - 4;
        assertThrows(BufferException.class, () -> buffer.readBits(overflowing),
            "readBits must refuse a request larger than the buffer even when the bounds "
                + "arithmetic would overflow");
    }

    /**
     * The whole window of bit counts for which the sum wraps. The defect is not a single value:
     * for a buffer positioned at bit 8, every count in [MAX_VALUE-7, MAX_VALUE] wraps.
     */
    @Test
    @DisplayName("every bit count in the wrapping window is rejected")
    void rejectsEveryBitCountInTheWrappingWindow() throws Exception {
        for (int offset = 0; offset < 8; offset++) {
            ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[16]);
            buffer.readBits(8);
            int numBits = Integer.MAX_VALUE - offset;
            assertThrows(BufferException.class, () -> buffer.readBits(numBits),
                "bit count " + numBits + " must be rejected");
        }
    }

    /**
     * The guard must not have been tightened into rejecting legitimate reads: a request for exactly
     * the remaining bits still succeeds, and one bit more does not.
     */
    @Test
    @DisplayName("reads up to the remaining bits still succeed")
    void acceptsExactlyTheRemainingBits() throws Exception {
        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[4]);
        buffer.readBits(8);
        assertEquals(24, buffer.getRemainingBits());

        byte[] rest = buffer.readBits(24);
        assertEquals(3, rest.length);
        assertEquals(0, buffer.getRemainingBits());

        ReadBufferByteBased other = new ReadBufferByteBased(new byte[4]);
        assertThrows(BufferException.class, () -> other.readBits(33));
    }
}
