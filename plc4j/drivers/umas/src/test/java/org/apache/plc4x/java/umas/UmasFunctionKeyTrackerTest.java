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
package org.apache.plc4x.java.umas;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class UmasFunctionKeyTrackerTest {

    @Test
    void trackedFunctionKeyIsConsumedOnce() {
        UmasFunctionKeyTracker tracker = new UmasFunctionKeyTracker();
        tracker.trackRequest(42, (short) 0x59);
        assertThat(tracker.consumeFunctionKey(42)).isEqualTo((short) 0x59);
        // Second consume returns 0 — the entry is removed on first consume.
        assertThat(tracker.consumeFunctionKey(42)).isZero();
    }

    @Test
    void unsolicitedMessagesReturnZero() {
        assertThat(new UmasFunctionKeyTracker().consumeFunctionKey(99999)).isZero();
    }

    /**
     * A transaction id is a counter belonging to one connection, so two connections reach the same
     * id in the ordinary course of things. While one map was shared across the process, the second
     * to arrive overwrote the first and a response was then parsed under the other conversation's
     * function key - which is what decides the response subtype.
     */
    @Test
    void twoConnectionsUsingTheSameTransactionIdDoNotOverwriteEachOther() {
        UmasFunctionKeyTracker first = new UmasFunctionKeyTracker();
        UmasFunctionKeyTracker second = new UmasFunctionKeyTracker();

        first.trackRequest(7, (short) 0x59);
        second.trackRequest(7, (short) 0x21);

        assertThat(first.consumeFunctionKey(7)).isEqualTo((short) 0x59);
        assertThat(second.consumeFunctionKey(7)).isEqualTo((short) 0x21);
    }

    @Test
    void aRequestThatWillNotBeAnsweredIsForgotten() {
        UmasFunctionKeyTracker tracker = new UmasFunctionKeyTracker();
        tracker.trackRequest(11, (short) 0x59);
        tracker.forget(11);
        assertThat(tracker.trackedCount()).isZero();
        // And a later request reaching the same id does not find the old key.
        assertThat(tracker.consumeFunctionKey(11)).isZero();
    }

    @Test
    void closingForgetsEverythingStillTracked() {
        UmasFunctionKeyTracker tracker = new UmasFunctionKeyTracker();
        tracker.trackRequest(1, (short) 0x59);
        tracker.trackRequest(2, (short) 0x21);
        tracker.clear();
        assertThat(tracker.trackedCount()).isZero();
    }
}
