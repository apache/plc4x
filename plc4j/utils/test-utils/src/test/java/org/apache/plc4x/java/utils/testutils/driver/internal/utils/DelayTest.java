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
package org.apache.plc4x.java.utils.testutils.driver.internal.utils;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

class DelayTest {

    @Test
    void testShortDelay() {
        long start = System.currentTimeMillis();
        Delay.shortDelay();
        long elapsed = System.currentTimeMillis() - start;

        // Should be at least 10ms, but allow some tolerance
        assertTrue(elapsed >= 5, "Short delay should be at least 5ms, was: " + elapsed);
    }

    @Test
    void testCustomDelay() {
        long start = System.currentTimeMillis();
        Delay.delay(20);
        long elapsed = System.currentTimeMillis() - start;

        assertTrue(elapsed >= 15, "Custom delay should be at least 15ms, was: " + elapsed);
    }

    @Test
    void testMediumDelay() {
        long start = System.currentTimeMillis();
        Delay.mediumDelay();
        long elapsed = System.currentTimeMillis() - start;

        assertTrue(elapsed >= 40, "Medium delay should be at least 40ms, was: " + elapsed);
    }

    @Test
    void testLongDelay() {
        long start = System.currentTimeMillis();
        Delay.longDelay();
        long elapsed = System.currentTimeMillis() - start;

        assertTrue(elapsed >= 90, "Long delay should be at least 90ms, was: " + elapsed);
    }

    @Test
    void testZeroDelay() {
        long start = System.currentTimeMillis();
        Delay.delay(0);
        long elapsed = System.currentTimeMillis() - start;

        assertTrue(elapsed < 100, "Zero delay should complete quickly, was: " + elapsed);
    }
}
