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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.junit.jupiter.api.Test;

import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.*;

class TestContextTest {

    @Test
    void testInitialState() {
        TestContext context = new TestContext();

        assertFalse(context.hasPendingResponse());
        assertNull(context.getPendingResponse());
    }

    @Test
    void testSetAndGetPendingResponse() {
        TestContext context = new TestContext();
        CompletableFuture<String> future = CompletableFuture.completedFuture("result");

        context.setPendingResponse(future);
        assertTrue(context.hasPendingResponse());

        CompletableFuture<?> retrieved = context.getPendingResponse();
        assertEquals(future, retrieved);

        // After getting, it should be cleared
        assertFalse(context.hasPendingResponse());
        assertNull(context.getPendingResponse());
    }

    @Test
    void testOverwritePendingResponse() {
        TestContext context = new TestContext();
        CompletableFuture<String> future1 = CompletableFuture.completedFuture("first");
        CompletableFuture<String> future2 = CompletableFuture.completedFuture("second");

        context.setPendingResponse(future1);
        context.setPendingResponse(future2);

        CompletableFuture<?> retrieved = context.getPendingResponse();
        assertEquals(future2, retrieved);
    }

    @Test
    void testClearAfterGet() {
        TestContext context = new TestContext();
        CompletableFuture<Integer> future = CompletableFuture.completedFuture(42);

        context.setPendingResponse(future);
        context.getPendingResponse(); // This should clear it

        assertFalse(context.hasPendingResponse());
    }
}
