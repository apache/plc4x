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
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

// TrappedStateMachine This class is a simple wrapper around the stateMachine class that keeps the
//
//	latest copy of the pdu parameter in the BeforeSend(), AfterSend(), BeforeReceive(), AfterReceive() and UnexpectedReceive() calls.
//
//	It also provides a send() function, so when the machine runs it doesn't
//	throw an exception.
type TrappedStateMachine struct {
	*Trapper
	StateMachine

	sent bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedStateMachine(localLog zerolog.Logger) *TrappedStateMachine {
	t := &TrappedStateMachine{
		log: localLog,
	}
	var init func()
	t.StateMachine, init = NewStateMachine(localLog, t, WithStateMachineStateInterceptor(t), WithStateMachineStateDecorator(t.DecorateState))
	t.Trapper = NewTrapper(localLog, t.StateMachine)
	init() // bit later so everything is set up
	return t
}

func (t *TrappedStateMachine) GetSent() bacnetip.PDU {
	return t.sent
}

func (t *TrappedStateMachine) BeforeSend(pdu bacnetip.PDU) {
	t.StateMachine.BeforeSend(pdu)
}

func (t *TrappedStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	// keep a copy
	t.sent = args.Get0PDU()
	return nil
}

func (t *TrappedStateMachine) AfterSend(pdu bacnetip.PDU) {
	t.StateMachine.AfterSend(pdu)
}

func (t *TrappedStateMachine) BeforeReceive(pdu bacnetip.PDU) {
	t.StateMachine.BeforeReceive(pdu)
}

func (t *TrappedStateMachine) AfterReceive(pdu bacnetip.PDU) {
	t.StateMachine.AfterReceive(pdu)
}

func (t *TrappedStateMachine) UnexpectedReceive(pdu bacnetip.PDU) {
	t.StateMachine.UnexpectedReceive(pdu)
}

func (t *TrappedStateMachine) DecorateState(state State) State {
	return NewTrappedState(state, t.Trapper)
}
