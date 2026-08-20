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

package umas

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The two cases plc4j's UmasFunctionKeyTrackerTest asserts.
func TestFunctionKeyTracker_TrackedKeyIsConsumedOnce(t *testing.T) {
	tracker := newFunctionKeyTracker()
	tracker.trackRequest(42, 0x22)
	assert.Equal(t, uint8(0x22), tracker.consumeFunctionKey(42))
	// The entry is gone after the first consume.
	assert.Equal(t, unknownFunctionKey, tracker.consumeFunctionKey(42))
}

func TestFunctionKeyTracker_UnknownTransactionsAnswerZero(t *testing.T) {
	tracker := newFunctionKeyTracker()
	assert.Equal(t, unknownFunctionKey, tracker.consumeFunctionKey(9999))
}

func TestFunctionKeyTracker_KeepsSeveralTransactionsApart(t *testing.T) {
	tracker := newFunctionKeyTracker()
	tracker.trackRequest(1, 0x02)
	tracker.trackRequest(2, 0x22)
	tracker.trackRequest(3, 0x23)
	assert.Equal(t, uint8(0x22), tracker.consumeFunctionKey(2))
	assert.Equal(t, uint8(0x23), tracker.consumeFunctionKey(3))
	assert.Equal(t, uint8(0x02), tracker.consumeFunctionKey(1))
}

// A transaction identifier comes back around every 65535 requests, and the second use has to win:
// the response still to come belongs to the newer request.
func TestFunctionKeyTracker_ReusingATransactionOverwritesIt(t *testing.T) {
	tracker := newFunctionKeyTracker()
	tracker.trackRequest(7, 0x22)
	tracker.trackRequest(7, 0x23)
	assert.Equal(t, uint8(0x23), tracker.consumeFunctionKey(7))
	assert.Equal(t, unknownFunctionKey, tracker.consumeFunctionKey(7))
}

// Entries whose response never arrives would otherwise pile up for the life of the connection.
func TestFunctionKeyTracker_EvictsTheOldestOnceItIsFull(t *testing.T) {
	tracker := newFunctionKeyTracker()
	for i := 0; i < maxTrackedFunctionKeys; i++ {
		tracker.trackRequest(uint16(i), 0x22)
	}
	// One more entry pushes the oldest one out and nothing else.
	tracker.trackRequest(uint16(maxTrackedFunctionKeys), 0x23)
	assert.Equal(t, unknownFunctionKey, tracker.consumeFunctionKey(0), "the oldest entry should have been evicted")
	assert.Equal(t, uint8(0x22), tracker.consumeFunctionKey(1))
	assert.Equal(t, uint8(0x23), tracker.consumeFunctionKey(uint16(maxTrackedFunctionKeys)))

	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	assert.LessOrEqual(t, len(tracker.keys), maxTrackedFunctionKeys)
	assert.Equal(t, len(tracker.keys), len(tracker.order), "the eviction order has to stay in step with the map")
}

// The codec tracks from Send and consumes from Receive, which run on different go routines.
func TestFunctionKeyTracker_IsSafeForConcurrentUse(t *testing.T) {
	tracker := newFunctionKeyTracker()
	var waitGroup sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		waitGroup.Add(1)
		go func(worker int) {
			defer waitGroup.Done()
			for i := 0; i < 100; i++ {
				identifier := uint16(worker*100 + i)
				tracker.trackRequest(identifier, 0x22)
				tracker.consumeFunctionKey(identifier)
			}
		}(worker)
	}
	waitGroup.Wait()
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	assert.Empty(t, tracker.keys)
	assert.Empty(t, tracker.order)
}
