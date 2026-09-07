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

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.junit.jupiter.api.Test;

import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A constructed data element may contain further constructed data, so the depth of the value tree
 * is decided by whoever sent the datagram. One level is a single opening tag, so a datagram that
 * fits comfortably inside a BACnet/IP frame describes a tree deeper than the parser has stack for.
 */
public class BACnetConstructedDataNestingTest {

    /**
     * Context tag 0 with a length/value/type of 6, which is how BACnet spells an opening tag. The
     * parser reads it as "constructed data follows" and descends a level.
     */
    private static final byte OPENING_TAG = 0x0E;

    private static ReadBuffer nestedOpeningTags(int levels) {
        byte[] data = new byte[levels];
        Arrays.fill(data, OPENING_TAG);
        return new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"));
    }

    private static void parse(ReadBuffer readBuffer) throws BufferException {
        BACnetConstructedDataElement.staticParse(readBuffer,
            BACnetObjectType.VENDOR_PROPRIETARY_VALUE,
            BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE,
            null);
    }

    @Test
    public void aDatagramFullOfOpeningTagsIsReportedAsAParseFailure() {
        // Enough levels to exhaust the parser's stack before the nesting bound was in place, and
        // still small enough to arrive as a single BACnet/IP frame.
        BufferException e = assertThrows(BufferException.class, () -> parse(nestedOpeningTags(1400)));
        assertTrue(e.getMessage().contains("nesting depth"),
            "the failure should name the bound that stopped it, but was: " + e.getMessage());
    }

    @Test
    public void aRealisticallyNestedTreeIsStillRead() {
        // The deepest message in the project's test suites nests 36 contexts. A tree at that scale
        // must run out of data, which is what a truncated tree should do, rather than run out of
        // budget.
        BufferException e = assertThrows(BufferException.class, () -> parse(nestedOpeningTags(64)));
        assertFalse(e.getMessage().contains("nesting depth"),
            "64 levels must not hit the depth bound, but did: " + e.getMessage());
    }

    @Test
    public void theParserNeverRunsOutOfStack() {
        // The bound exists so that this stays a BufferException. Assert the absence of the error
        // it replaced, rather than only the presence of the exception.
        for (int levels : new int[]{1000, 2000, 8000}) {
            assertThrows(BufferException.class, () -> parse(nestedOpeningTags(levels)),
                "a tree of " + levels + " levels must fail as a parse error");
        }
    }

}
