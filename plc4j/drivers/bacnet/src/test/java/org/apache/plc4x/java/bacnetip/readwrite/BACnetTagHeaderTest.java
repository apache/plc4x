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

package org.apache.plc4x.java.bacnetip.readwrite;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * A tag header only carries the fields its own leading bits announce: the extended length exclusively
 * for a primitive, non-boolean tag, and the extended tag number exclusively for tag number 15.
 * Deriving {@code actualLength} and {@code actualTagNumber} therefore has to tolerate those fields
 * being absent, which the length-value-type-5 case below did not.
 */
class BACnetTagHeaderTest {

    private static final WithOption[] OPTS = {
        WithByteBasedOption.WithByteOrder("BIG_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    private static BACnetTagHeader parse(int... data) throws Exception {
        byte[] bytes = new byte[data.length];
        for (int i = 0; i < data.length; i++) {
            bytes[i] = (byte) data[i];
        }
        return BACnetTagHeader.staticParse(new ReadBufferByteBased(bytes, OPTS));
    }

    /**
     * {@code 0x15} is tag number 1 in the application class, so a boolean tag, with
     * length-value-type 5. Length-value-type 5 announces an extended length for a primitive tag, but
     * a boolean tag carries none, so {@code actualLength} has to fall back to the length-value-type.
     */
    @Test
    void booleanTagWithLengthValueTypeFiveHasNoExtendedLength() throws Exception {
        BACnetTagHeader header = parse(0x15);

        assertNull(header.getExtLength(), "a boolean tag carries no extended length");
        assertEquals(5, header.getActualLength());
    }

    /**
     * The extended length still has to be picked up for a tag that really does announce one:
     * {@code 0x25} is a primitive tag with length-value-type 5.
     */
    @Test
    void primitiveTagWithLengthValueTypeFiveReadsItsExtendedLength() throws Exception {
        BACnetTagHeader header = parse(0x25, 0x2A);

        assertEquals((short) 0x2A, header.getExtLength());
        assertEquals(0x2A, header.getActualLength());
    }

    /**
     * The extended tag number likewise: {@code 0xF0} announces tag number 15.
     */
    @Test
    void tagNumberFifteenReadsItsExtendedTagNumber() throws Exception {
        BACnetTagHeader header = parse(0xF0, 0x2A);

        assertEquals((short) 0x2A, header.getExtTagNumber());
        assertEquals(0x2A, header.getActualTagNumber());
    }

    /**
     * A header that announces a field the datagram is too short to hold is rejected, rather than
     * being completed with a missing field that the derived values would then trip over.
     */
    @Test
    void headerTruncatedBeforeAnAnnouncedFieldIsRejected() {
        assertThrows(BufferException.class, () -> parse(0x05), "extended length announced but absent");
        assertThrows(BufferException.class, () -> parse(0xF0), "extended tag number announced but absent");
    }

    /**
     * No tag header, however truncated or self-contradictory, may fail with anything other than a
     * {@link BufferException}. An unchecked exception escapes the driver's parse error handling and
     * takes the transport's receive thread down with it.
     */
    @Test
    void noTagHeaderFailsWithAnUncheckedException() {
        // Enough trailing bytes for the extended tag number plus the widest extended length.
        int[] trailer = {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF};
        for (int first = 0; first <= 0xFF; first++) {
            for (int trailerLength = 0; trailerLength <= trailer.length; trailerLength++) {
                int[] data = new int[1 + trailerLength];
                data[0] = first;
                System.arraycopy(trailer, 0, data, 1, trailerLength);

                String message = String.format("header 0x%02X with %d trailing bytes", first, trailerLength);
                try {
                    BACnetTagHeader header = parse(data);
                    // Touching the derived fields has to be safe too.
                    header.getActualLength();
                    header.getActualTagNumber();
                } catch (BufferException e) {
                    // A malformed header is allowed to be rejected.
                } catch (RuntimeException e) {
                    fail(message + " threw " + e, e);
                } catch (Exception e) {
                    fail(message + " threw " + e, e);
                }
            }
        }
    }
}
