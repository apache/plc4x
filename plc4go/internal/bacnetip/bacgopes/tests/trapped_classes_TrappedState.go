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

package tests

import (
	"fmt"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
)

// TrappedState This class is a simple wrapper around the state class that keeps the latest copy of the pdu parameter in
// the BeforeSend(), AfterSend(), BeforeReceive(), AfterReceive() and UnexpectedReceive() calls.
type TrappedState struct {
	*Trapper
	State
}

func NewTrappedState(state State, trapper *Trapper) *TrappedState {
	t := &TrappedState{
		State:   state,
		Trapper: trapper,
	}
	return t
}

func (t *TrappedState) Equals(other State) bool {
	if t.State.Equals(other) { //TODO: we always want to match the inner
		return true
	}
	if otherTs, ok := other.(*TrappedState); ok {
		return t.State.Equals(otherTs.State)
	}
	return false
}

func (t *TrappedState) String() string {
	return fmt.Sprintf("TrappedState(%v)", t.State)
}

func (t *TrappedState) BeforeSend(pdu bacgopes.PDU) {
	t.Trapper.BeforeSend(pdu)
}

func (t *TrappedState) AfterSend(pdu bacgopes.PDU) {
	t.Trapper.AfterSend(pdu)
}

func (t *TrappedState) BeforeReceive(pdu bacgopes.PDU) {
	t.Trapper.BeforeReceive(pdu)
}

func (t *TrappedState) AfterReceive(pdu bacgopes.PDU) {
	t.Trapper.AfterReceive(pdu)
}

func (t *TrappedState) UnexpectedReceive(pdu bacgopes.PDU) {
	t.Trapper.UnexpectedReceive(pdu)
}

func (t *TrappedState) getInterceptor() StateInterceptor {
	return t
}
