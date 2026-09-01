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

package abeth

import (
	"fmt"
	"sync/atomic"
)

// maxTransactionCounter is the last transaction counter that is handed out. plc4j's
// AbEthConnection.nextTransactionCounter restarts at 1 once the counter reaches 0xFFFF, so 0xFFFE
// is the last value used and 0 never is.
const maxTransactionCounter = uint32(0xFFFE)

// session is the state a connected ab-eth connection carries: the session handle the PLC handed out
// during the connect handshake and the transaction counter that ties reads to their responses.
//
// The reader runs on whatever go routine executed the read request while the connection may be
// re-running its handshake, so both fields are atomic.
type session struct {
	sessionHandle      atomic.Uint32
	transactionCounter atomic.Uint32
}

func newSession() *session {
	return &session{}
}

// getSessionHandle is the handle every packet after the handshake has to carry.
func (s *session) getSessionHandle() uint32 {
	return s.sessionHandle.Load()
}

func (s *session) setSessionHandle(sessionHandle uint32) {
	s.sessionHandle.Store(sessionHandle)
}

// nextTransactionCounter hands out the counter of the next read. It cycles through 1..0xFFFE, which
// is what plc4j's AbEthConnection does; zero is never used, so a response carrying it can't be
// mistaken for an answer to a request of ours.
func (s *session) nextTransactionCounter() uint16 {
	for {
		current := s.transactionCounter.Load()
		next := current + 1
		if next > maxTransactionCounter {
			next = 1
		}
		if s.transactionCounter.CompareAndSwap(current, next) {
			return uint16(next)
		}
	}
}

func (s *session) String() string {
	return fmt.Sprintf("abeth.session{sessionHandle: %d, transactionCounter: %d}",
		s.sessionHandle.Load(), s.transactionCounter.Load())
}
