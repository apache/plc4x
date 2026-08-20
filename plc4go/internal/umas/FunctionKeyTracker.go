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
)

// maxTrackedFunctionKeys bounds the tracker. A UMAS connection has one request in flight at a time,
// so a single entry is the normal case; the cap only matters for entries which are never consumed
// because their response never came, and it keeps such a stream of timeouts from growing the map
// without end.
const maxTrackedFunctionKeys = 256

// unknownFunctionKey is what the tracker answers with for a transaction it doesn't know. Passing it
// to the parser lands in the catch-all branch of the UmasPDUItem type switch (function key 0xFE
// alone parses as a UmasPDUSuccessResponse), which is the same thing plc4j's
// UmasFunctionKeyTracker.consumeFunctionKey does by returning 0.
const unknownFunctionKey = uint8(0)

// functionKeyTracker remembers which UMAS function key a request carried so the response can be
// parsed at all: every UMAS response uses the generic function key 0xFE and the mspec discriminates
// the concrete response type on the *request's* function key, which is nowhere on the wire of the
// response.
//
// Deliberate deviation from plc4j, whose UmasFunctionKeyTracker is a class with a static
// ConcurrentHashMap: state shared by every connection in the JVM, unbounded, and impossible to
// isolate in a test. Here the tracker belongs to the message codec, which is the one object that
// both sees every request go out and every response come in.
type functionKeyTracker struct {
	mutex sync.Mutex
	keys  map[uint16]uint8
	// order is the insertion order of the keys, used to evict the oldest entry once the map is at
	// its cap.
	order []uint16
}

func newFunctionKeyTracker() *functionKeyTracker {
	return &functionKeyTracker{
		keys: map[uint16]uint8{},
	}
}

// trackRequest records that the request with this transaction identifier carried this function key.
// Re-using a transaction identifier overwrites the previous entry rather than adding a second one,
// which is what makes the tracker survive the 16 bit wrap-around of the identifier.
func (t *functionKeyTracker) trackRequest(transactionIdentifier uint16, functionKey uint8) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if _, alreadyTracked := t.keys[transactionIdentifier]; !alreadyTracked {
		if len(t.order) >= maxTrackedFunctionKeys {
			oldest := t.order[0]
			t.order = t.order[1:]
			delete(t.keys, oldest)
		}
		t.order = append(t.order, transactionIdentifier)
	}
	t.keys[transactionIdentifier] = functionKey
}

// consumeFunctionKey hands out and forgets the function key of a transaction, answering
// unknownFunctionKey for a transaction which was never tracked.
func (t *functionKeyTracker) consumeFunctionKey(transactionIdentifier uint16) uint8 {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	functionKey, ok := t.keys[transactionIdentifier]
	if !ok {
		return unknownFunctionKey
	}
	delete(t.keys, transactionIdentifier)
	for i, tracked := range t.order {
		if tracked == transactionIdentifier {
			t.order = append(t.order[:i], t.order[i+1:]...)
			break
		}
	}
	return functionKey
}
