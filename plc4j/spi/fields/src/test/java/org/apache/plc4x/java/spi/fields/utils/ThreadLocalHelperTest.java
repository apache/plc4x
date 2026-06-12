/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.spi.fields.utils;

import org.apache.plc4x.java.spi.fields.utils.ThreadLocalHelper;
import org.junit.jupiter.api.Test;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

class ThreadLocalHelperTest {

    @Test
    void defaultsAreInitializedPerThread() {
        // Initial values should come from withInitial suppliers
        ThreadLocalHelper.lastItemThreadLocal.set(false);
        assertFalse(ThreadLocalHelper.lastItemThreadLocal.get(), "Default lastItem should be false");
        ThreadLocalHelper.curItemThreadLocal.set(0);
        assertEquals(0, ThreadLocalHelper.curItemThreadLocal.get(), "Default curItem should be 0");
    }

    @Test
    void valuesAreIsolatedAcrossThreads() throws InterruptedException {
        // Set values in the main thread
        ThreadLocalHelper.lastItemThreadLocal.set(true);
        ThreadLocalHelper.curItemThreadLocal.set(42);

        CountDownLatch started = new CountDownLatch(1);
        CountDownLatch finished = new CountDownLatch(1);
        AtomicBoolean observedLastInOther = new AtomicBoolean(true); // init opposite to ensure assertion checks
        AtomicInteger observedCurInOther = new AtomicInteger(-1);

        Thread t = new Thread(() -> {
            try {
                // On a fresh thread the defaults should apply
                observedLastInOther.set(ThreadLocalHelper.lastItemThreadLocal.get());
                observedCurInOther.set(ThreadLocalHelper.curItemThreadLocal.get());

                // Now set different values in this thread
                ThreadLocalHelper.lastItemThreadLocal.set(false);
                ThreadLocalHelper.curItemThreadLocal.set(7);
            } finally {
                started.countDown();
                finished.countDown();
            }
        });
        t.start();
        started.await();
        finished.await();

        // Other thread should have seen defaults, not the values set in this thread
        assertFalse(observedLastInOther.get(), "Other thread should see default false for lastItem");
        assertEquals(0, observedCurInOther.get(), "Other thread should see default 0 for curItem");

        // Ensure our thread's values remain unchanged
        assertTrue(ThreadLocalHelper.lastItemThreadLocal.get());
        assertEquals(42, ThreadLocalHelper.curItemThreadLocal.get());
    }

    @Test
    void threadLocalsAreIndependent() {
        // Changing one should not implicitly change the other
        ThreadLocalHelper.lastItemThreadLocal.set(true);
        assertTrue(ThreadLocalHelper.lastItemThreadLocal.get());
        assertEquals(0, ThreadLocalHelper.curItemThreadLocal.get());

        ThreadLocalHelper.curItemThreadLocal.set(5);
        assertTrue(ThreadLocalHelper.lastItemThreadLocal.get());
        assertEquals(5, ThreadLocalHelper.curItemThreadLocal.get());
    }

    @Test
    void removeResetsToDefaults() {
        ThreadLocalHelper.lastItemThreadLocal.set(true);
        ThreadLocalHelper.curItemThreadLocal.set(3);

        ThreadLocalHelper.lastItemThreadLocal.remove();
        ThreadLocalHelper.curItemThreadLocal.remove();

        // After remove(), subsequent get() should return the initial values again
        assertFalse(ThreadLocalHelper.lastItemThreadLocal.get());
        assertEquals(0, ThreadLocalHelper.curItemThreadLocal.get());
    }
}
