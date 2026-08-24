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
package org.apache.plc4x.java.opcua.readwrite;

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
 * A variant of type 24 holds further variants, so the depth of the value is decided by the server.
 * One level is a single encoding-mask byte, so any response large enough to be worth sending can
 * describe a value deeper than the parser has stack for.
 */
public class VariantNestingTest {

    /**
     * A variant encoding mask of 24 with neither the array-length nor the array-dimensions bit
     * set: one nested variant, one byte.
     */
    private static final byte NESTED_VARIANT = 24;

    private static ReadBuffer nestedVariants(int levels) {
        byte[] data = new byte[levels];
        Arrays.fill(data, NESTED_VARIANT);
        return new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"));
    }

    @Test
    public void aValueNestedPastTheBoundIsReportedAsAParseFailure() {
        BufferException e = assertThrows(BufferException.class,
            () -> Variant.staticParse(nestedVariants(3000)));
        assertTrue(e.getMessage().contains("nesting depth"),
            "the failure should name the bound that stopped it, but was: " + e.getMessage());
    }

    @Test
    public void theParserNeverRunsOutOfStack() {
        // Before the bound this was a StackOverflowError, which is neither a parse failure the
        // driver can report nor an error the receive path can contain.
        for (int levels : new int[]{2000, 4000, 16000}) {
            assertThrows(BufferException.class, () -> Variant.staticParse(nestedVariants(levels)),
                "a value nested " + levels + " deep must fail as a parse error");
        }
    }

    @Test
    public void aRealisticallyNestedValueIsStillRead() {
        // The deepest OPC UA message in the project's test suites nests 29 contexts. A value at
        // that scale must still be parsed as data.
        BufferException e = assertThrows(BufferException.class,
            () -> Variant.staticParse(nestedVariants(32)));
        assertFalse(e.getMessage().contains("nesting depth"),
            "32 levels must run out of data, not out of budget, but was: " + e.getMessage());
    }
}
