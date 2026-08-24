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

import org.apache.plc4x.java.spi.buffers.api.AbstractBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A message nests one context per type it descends into, so a type that contains itself lets the
 * sender decide how deep we recurse. The depth is bounded so that the nesting is reported as a
 * parse failure rather than as an exhausted stack.
 */
public class NestingDepthTest {

    private static ReadBufferByteBased buffer() {
        return new ReadBufferByteBased(new byte[0]);
    }

    private static void nest(ReadBufferByteBased readBuffer, int levels) throws BufferException {
        for (int i = 0; i < levels; i++) {
            readBuffer.pushContext(WithOption.WithName("level" + i));
        }
    }

    @Test
    public void nestingIsAllowedUpToTheBound() {
        // The buffer starts with the context it was constructed with, so the bound leaves room for
        // one less push than its own value.
        assertDoesNotThrow(() -> nest(buffer(), AbstractBuffer.getMaxNestingDepth() - 1));
    }

    @Test
    public void nestingBeyondTheBoundIsRejected() throws BufferException {
        ReadBufferByteBased readBuffer = buffer();
        nest(readBuffer, AbstractBuffer.getMaxNestingDepth() - 1);
        BufferException e = assertThrows(BufferException.class,
            () -> readBuffer.pushContext(WithOption.WithName("oneTooDeep")));
        assertTrue(e.getMessage().contains("nesting depth"),
            "the failure should name the bound that stopped it, but was: " + e.getMessage());
    }

    @Test
    public void aRejectedContextIsNotPushed() throws BufferException {
        ReadBufferByteBased readBuffer = buffer();
        nest(readBuffer, AbstractBuffer.getMaxNestingDepth() - 1);
        WithOption[] before = readBuffer.getContext();
        assertThrows(BufferException.class, () -> readBuffer.pushContext(WithOption.WithName("oneTooDeep")));
        assertEquals(before, readBuffer.getContext(),
            "a rejected push must leave the context stack where it was");
    }

    @Test
    public void leavingATypeMakesRoomAgain() throws BufferException {
        ReadBufferByteBased readBuffer = buffer();
        nest(readBuffer, AbstractBuffer.getMaxNestingDepth() - 1);
        readBuffer.popContext();
        assertDoesNotThrow(() -> readBuffer.pushContext(WithOption.WithName("backInBudget")));
    }

    @Test
    public void theDepthCanBeGivenBackToWhereItWas() throws BufferException {
        ReadBufferByteBased readBuffer = buffer();
        int before = readBuffer.getContextDepth();
        nest(readBuffer, 50);
        assertEquals(before + 50, readBuffer.getContextDepth());
        readBuffer.resetContextDepth(before);
        assertEquals(before, readBuffer.getContextDepth(),
            "giving the depth back must close every context opened since it was taken");
        assertDoesNotThrow(() -> nest(readBuffer, AbstractBuffer.getMaxNestingDepth() - 1 - before));
    }

    @Test
    public void givingBackADepthWeAreAlreadyAboveChangesNothing() throws BufferException {
        ReadBufferByteBased readBuffer = buffer();
        nest(readBuffer, 10);
        int depth = readBuffer.getContextDepth();
        readBuffer.resetContextDepth(depth + 5);
        assertEquals(depth, readBuffer.getContextDepth());
    }

    @Test
    public void theEnvironmentCanRaiseOrLowerTheBound() {
        assertEquals(64, AbstractBuffer.resolveMaxNestingDepth("64"));
        assertEquals(4096, AbstractBuffer.resolveMaxNestingDepth(" 4096 "));
    }

    @Test
    public void anUnusableEnvironmentValueLeavesTheDefaultInPlace() {
        for (String value : new String[]{null, "", "   ", "plenty", "0", "-1"}) {
            assertEquals(AbstractBuffer.DEFAULT_MAX_NESTING_DEPTH,
                AbstractBuffer.resolveMaxNestingDepth(value),
                "a value of '" + value + "' should leave the default in place");
        }
    }
}
