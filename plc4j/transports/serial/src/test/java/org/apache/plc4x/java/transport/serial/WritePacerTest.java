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
package org.apache.plc4x.java.transport.serial;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class WritePacerTest {

    @Test
    void zeroDelayIsNoop() throws Exception {
        WritePacer pacer = new WritePacer(0);
        pacer.noteActivity();
        long start = System.currentTimeMillis();
        pacer.awaitTurn();
        assertTrue(System.currentTimeMillis() - start < 20, "zero delay must not wait");
    }

    @Test
    void enforcesGapAfterActivity() throws Exception {
        WritePacer pacer = new WritePacer(60);
        pacer.noteActivity();
        long start = System.currentTimeMillis();
        pacer.awaitTurn();
        assertTrue(System.currentTimeMillis() - start >= 50, "must wait out the gap");
    }

    @Test
    void noWaitWhenGapElapsed() throws Exception {
        WritePacer pacer = new WritePacer(30);
        pacer.noteActivity();
        Thread.sleep(40);
        long start = System.currentTimeMillis();
        pacer.awaitTurn();
        assertTrue(System.currentTimeMillis() - start < 20, "gap already elapsed");
    }

    @Test
    void activityDuringWaitExtendsGap() throws Exception {
        WritePacer pacer = new WritePacer(80);
        pacer.noteActivity();
        Thread traffic = new Thread(() -> {
            try {
                Thread.sleep(40);
                pacer.noteActivity(); // traffic arrives while a writer waits its turn
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });
        traffic.start();
        long start = System.currentTimeMillis();
        pacer.awaitTurn();
        traffic.join(1000);
        assertTrue(System.currentTimeMillis() - start >= 110,
            "the gap must restart from the mid-wait activity (40ms + 80ms)");
    }
}
