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

import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;

import static org.junit.jupiter.api.Assertions.*;

class WritePacerTest {

    @Test
    void zeroDelayIsNoop() throws Exception {
        WritePacer pacer = new WritePacer(0);
        pacer.noteActivity();
        long start = System.nanoTime();
        pacer.awaitTurn();
        // Generously bounded: the claim is "did not wait", and no configured gap exists to wait
        // out, so the bound only has to be tighter than a wait - not tighter than a GC pause.
        assertTrue(System.nanoTime() - start < TimeUnit.MILLISECONDS.toNanos(500),
            "zero delay must not wait");
    }

    @Test
    void enforcesGapAfterActivity() throws Exception {
        WritePacer pacer = new WritePacer(60);
        // Sampled before noteActivity(), so it is never later than the pacer's own stamp: the
        // gap it then has to wait out is at least this long, whatever the machine does in
        // between. Measuring from after the call would let a stall eat into the margin.
        long activity = System.nanoTime();
        pacer.noteActivity();
        pacer.awaitTurn();
        assertTrue(System.nanoTime() - activity >= TimeUnit.MILLISECONDS.toNanos(60),
            "must wait out the gap");
    }

    @Test
    void noWaitWhenGapElapsed() throws Exception {
        WritePacer pacer = new WritePacer(200);
        pacer.noteActivity();
        Thread.sleep(250);
        long start = System.nanoTime();
        pacer.awaitTurn();
        // The gap is 200ms, so returning within 100ms proves it was not waited out. The wide
        // margin between the two is what keeps a scheduling stall from failing the test.
        assertTrue(System.nanoTime() - start < TimeUnit.MILLISECONDS.toNanos(100),
            "gap already elapsed");
    }

    @Test
    void activityDuringWaitExtendsGap() throws Exception {
        final long gapMs = 200;
        final long midWaitMs = 50;
        WritePacer pacer = new WritePacer(gapMs);
        AtomicLong midWaitActivity = new AtomicLong();

        long firstActivity = System.nanoTime();
        pacer.noteActivity();
        Thread traffic = new Thread(() -> {
            try {
                Thread.sleep(midWaitMs);
                // Stamped before the call, so it is never later than the pacer's own stamp.
                midWaitActivity.set(System.nanoTime());
                pacer.noteActivity(); // traffic arrives while a writer waits its turn
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });
        traffic.start();
        pacer.awaitTurn();
        long returned = System.nanoTime();
        traffic.join(1000);

        long activity = midWaitActivity.get();
        assertNotEquals(0, activity, "the traffic thread never recorded its activity");
        // The point of the test is activity arriving *during* the wait, so say so out loud: if a
        // stalled machine delivered it after the original gap had already expired, the premise is
        // gone and the run proves nothing - that is a broken test, not a broken pacer.
        assertTrue(activity - firstActivity < TimeUnit.MILLISECONDS.toNanos(gapMs),
            "test premise broken: the mid-wait activity arrived only after the original gap had expired");
        // The contract: the gap restarts from that activity. Both stamps come from the same clock
        // and bracket the pacer's own, so this holds no matter how the threads were scheduled.
        assertTrue(returned - activity >= TimeUnit.MILLISECONDS.toNanos(gapMs),
            "the gap must restart from the mid-wait activity");
    }
}
