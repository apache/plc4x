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

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class SynchronizerTest {

    @Test
    void testResponseSent() {
        Synchronizer sync = new Synchronizer();

        // Start a thread that signals after a delay
        new Thread(() -> {
            Delay.delay(50);
            sync.responseSent();
        }).start();

        // Wait for response
        boolean received = sync.waitForResponseSent(200);
        assertTrue(received, "Should receive response signal");
    }

    @Test
    void testTimeout() {
        Synchronizer sync = new Synchronizer();

        // Don't send any signal
        boolean received = sync.waitForResponseSent(50);
        assertFalse(received, "Should timeout waiting for response");
    }

    @Test
    void testReset() {
        Synchronizer sync = new Synchronizer();

        sync.responseSent();
        boolean received1 = sync.waitForResponseSent(10);
        assertTrue(received1, "Should receive first signal");

        sync.reset();
        boolean received2 = sync.waitForResponseSent(10);
        assertFalse(received2, "Should timeout after reset");
    }
}
