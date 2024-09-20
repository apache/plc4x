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

package state_machine

import (
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

type StateInterceptor interface {
	BeforeSend(pdu PDU)
	AfterSend(pdu PDU)
	BeforeReceive(pdu PDU)
	AfterReceive(pdu PDU)
	UnexpectedReceive(pdu PDU)
}

type State interface {
	fmt.Stringer

	Send(pdu PDU, nextState State) State
	Receive(args Args, kwArgs KWArgs) State
	Reset()
	Fail(docstring string) State
	Success(docstring string) State
	ExitState()
	EnterState()
	EventSet(eventId string)
	Timeout(duration time.Duration, nextState State) State
	WaitEvent(eventId string, nextState State) State
	SetEvent(eventId string) State
	Doc(docstring string) State
	DocString() string
	Call(fn GenericFunction, args Args, kwArgs KWArgs) State

	getStateMachine() StateMachine
	setStateMachine(StateMachine)
	getDocString() string
	IsSuccessState() bool
	IsFailState() bool
	getSendTransitions() []SendTransition
	getReceiveTransitions() []ReceiveTransition
	getSetEventTransitions() []EventTransition
	getClearEventTransitions() []EventTransition
	getWaitEventTransitions() []EventTransition
	getTimeoutTransition() *TimeoutTransition
	getCallTransition() *CallTransition
	getInterceptor() StateInterceptor

	Equals(other State) bool
}

type state struct {
	interceptor StateInterceptor

	stateMachine          StateMachine
	docString             string
	isSuccessState        bool
	isFailState           bool
	sendTransitions       []SendTransition
	receiveTransitions    []ReceiveTransition
	setEventTransitions   []EventTransition
	clearEventTransitions []EventTransition
	waitEventTransitions  []EventTransition
	timeoutTransition     *TimeoutTransition
	callTransition        *CallTransition

	_leafName string

	log zerolog.Logger
}

func NewState(localLog zerolog.Logger, stateMachine StateMachine, docString string, options ...Option) State {
	s := &state{
		stateMachine: stateMachine,
		docString:    docString,

		_leafName: ExtractLeafName(options, StructName()),

		log: localLog,
	}
	ApplyAppliers(options, s)
	if _debug != nil {
		_debug("__init__ %r docString=%r", stateMachine, docString)
	}
	if s.interceptor == nil {
		s.interceptor = s
	}
	if docString != "" {
		s.log = s.log.With().Str("docString", docString).Logger()
	}
	return s
}

func WithStateStateInterceptor(interceptor StateInterceptor) GenericApplier[*state] {
	return WrapGenericApplier(func(state *state) { state.interceptor = interceptor })
}

func (s *state) getStateMachine() StateMachine {
	return s.stateMachine
}

func (s *state) setStateMachine(machine StateMachine) {
	s.stateMachine = machine
}

func (s *state) getDocString() string {
	return s.docString
}

func (s *state) getSendTransitions() []SendTransition {
	return s.sendTransitions
}

func (s *state) getReceiveTransitions() []ReceiveTransition {
	return s.receiveTransitions
}

func (s *state) getSetEventTransitions() []EventTransition {
	return s.setEventTransitions
}

func (s *state) getClearEventTransitions() []EventTransition {
	return s.clearEventTransitions
}

func (s *state) getWaitEventTransitions() []EventTransition {
	return s.waitEventTransitions
}

func (s *state) getTimeoutTransition() *TimeoutTransition {
	return s.timeoutTransition
}

func (s *state) getCallTransition() *CallTransition {
	return s.callTransition
}

func (s *state) getInterceptor() StateInterceptor {
	return s.interceptor
}

// Reset Override this method in a derived class if the state maintains counters or other information.  Called when the
//
//	associated state machine is Reset.
func (s *state) Reset() {
	if _debug != nil {
		_debug("reset")
	}
	s.log.Trace().Msg("Reset")
}

// Doc Change the documentation string (label) for the state.  The state
//
//	is returned for method chaining.
func (s *state) Doc(docString string) State {
	if _debug != nil {
		_debug("doc %r", docString)
	}
	s.log = s.log.With().Str("docString", docString).Logger() // TODO: would be better to use stringer or object to dynamically retrieve docstring and only build logger once
	s.log.Debug().Msg("Doc")
	s.docString = docString
	return s
}

func (s *state) DocString() string {
	return s.docString
}

// Success Mark a state as a successful final state.  The state is returned for method chaining. docString: an optional
//
//	label for the state
func (s *state) Success(docString string) State {
	if _debug != nil {
		_debug("success %r", docString)
	}
	s.log.Debug().Str("docString", docString).Msg("Success")
	if s.isSuccessState {
		panic("already a Success state")
	}
	if s.isFailState {
		panic("already a Fail state")
	}

	s.isSuccessState = true

	if docString != "" {
		s.docString = docString
	} else if s.docString == "" {
		s.docString = "Success"
	}
	return s
}

func (s *state) IsSuccessState() bool {
	return s.isSuccessState
}

// Fail Mark a state as a failure final state.  The state is returned for method chaining. docString: an optional
//
//	label for the state
func (s *state) Fail(docString string) State {
	if _debug != nil {
		_debug("fail %r", docString)
	}
	s.log.Debug().Str("docString", docString).Msg("Fail")
	if s.isSuccessState {
		panic("already a Success state")
	}
	if s.isFailState {
		panic("already a Fail state")
	}

	s.isFailState = true

	if docString != "" {
		s.docString = docString
	} else if s.docString == "" {
		s.docString = "Fail"
	}
	return s
}

func (s *state) IsFailState() bool {
	return s.isFailState
}

// EnterState Called when the state machine is entering the state.
func (s *state) EnterState() {
	if _debug != nil {
		_debug("enter_state(%s)", s.docString)
	}
	s.log.Debug().Msg("EnterState")
	if s.timeoutTransition != nil {
		if _debug != nil {
			_debug("    - waiting: %r", s.timeoutTransition.timeout)
		}
		s.log.Debug().Time("timeout", s.timeoutTransition.timeout).Msg("waiting")
		s.stateMachine.getStateTimeoutTask().InstallTask(WithInstallTaskOptionsWhen(s.timeoutTransition.timeout))
	} else {
		if _debug != nil {
			_debug("    - no timeout")
		}
		s.log.Trace().Msg("no timeout")
	}
}

// ExitState Called when the state machine is existing the state.
func (s *state) ExitState() {
	if _debug != nil {
		_debug("exit_state(%s)", s.docString)
	}
	s.log.Debug().Msg("ExitState")
	if s.timeoutTransition != nil {
		if _debug != nil {
			_debug("    - canceling timeout")
		}
		s.log.Trace().Msg("canceling timeout")
		s.stateMachine.getStateTimeoutTask().SuspendTask()
	}
}

// Send Create a SendTransition from this state to another, possibly new, state.  The next state is returned for method
//
//	chaining. pdu tPDU to send nextState state to transition to after sending
func (s *state) Send(pdu PDU, nextState State) State {
	if _debug != nil {
		_debug("send(%s) %r next_state=%r", s.docString, pdu, nextState)
	}
	s.log.Debug().Stringer("pdu", pdu).Msg("Send")
	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		s.log.Debug().Stringer("nextState", nextState).Msg("new nextState")
		if _debug != nil {
			_debug("    - new next_state: %r", nextState)
		}
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	s.sendTransitions = append(s.sendTransitions, SendTransition{
		Transition: Transition{nextState: nextState},
		pdu:        pdu,
	})
	return nextState
}

// BeforeSend Called before each tPDU about to be sent.
func (s *state) BeforeSend(pdu PDU) {
	s.stateMachine.BeforeSend(pdu)
}

// AfterSend Called after each tPDU about to be sent.
func (s *state) AfterSend(pdu PDU) {
	s.stateMachine.AfterSend(pdu)
}

// Receive Create a ReceiveTransition from this state to another, possibly new,
//
//	state.  The next state is returned for method chaining.
//
//	criteria tPDU to match
//	 next_state destination state after a successful match
func (s *state) Receive(args Args, kwArgs KWArgs) State {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Receive")
	pduType := args[0]
	pduAttrs := kwArgs
	if _debug != nil {
		_debug("receive(%s) %r %r", s.docString, pduType, pduAttrs)
	}
	var nextState State
	if _nextState, ok := pduAttrs["next_state"]; ok {
		nextState = _nextState.(State)
		if _debug != nil {
			_debug("    - next_state: %r", nextState)
		}
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	}

	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		if _debug != nil {
			_debug("    - new next_state: %r", nextState)
		}
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	// add this to the list of transitions
	criteria := criteria{
		pduType:  pduType,
		pduAttrs: pduAttrs,
	}
	if _debug != nil {
		_debug("    - criteria: %r", criteria)
	}
	s.log.Debug().Interface("criteria", criteria).Msg("criteria")
	s.receiveTransitions = append(s.receiveTransitions, ReceiveTransition{
		Transition: Transition{nextState: nextState},
		criteria:   criteria,
	})

	return nextState
}

// BeforeReceive Called with each tPDU received before matching.
func (s *state) BeforeReceive(pdu PDU) {
	s.stateMachine.BeforeReceive(pdu)
}

// AfterReceive Called with tPDU received after match.
func (s *state) AfterReceive(pdu PDU) {
	s.stateMachine.AfterReceive(pdu)
}

// Ignore Create a ReceiveTransition from this state to its, if match is successful the effect is to Ignore the tPDU.
//
//	criteria tPDU to match
func (s *state) Ignore(pduType any, pduAttrs map[KnownKey]any) State {
	if _debug != nil {
		_debug("ignore(%s) %r %r", s.docString, pduType, pduAttrs)
	}
	s.log.Debug().Interface("pduType", pduType).Interface("pduAttrs", pduAttrs).Msg("Ignore")
	criteria := criteria{
		pduType:  pduType,
		pduAttrs: pduAttrs,
	}
	if _debug != nil {
		_debug("    - criteria: %r", criteria)
	}
	s.receiveTransitions = append(s.receiveTransitions, ReceiveTransition{
		Transition: Transition{},
		criteria:   criteria,
	})

	return s
}

// UnexpectedReceive Called with PDU that did not match.
//
// Unless this is trapped by the state, the default behaviour is to fail.
func (s *state) UnexpectedReceive(pdu PDU) {
	if _debug != nil {
		_debug("unexpected_receive(%s) %r", s.docString, pdu)
	}
	s.log.Debug().Stringer("pdu", pdu).Msg("UnexpectedReceive")
	s.stateMachine.UnexpectedReceive(pdu)
}

// SetEvent Create an EventTransition for this state that sets an event.  The current state is returned for method
//
//	chaining. event_id event identifier
func (s *state) SetEvent(eventId string) State {
	if _debug != nil {
		_debug("set_event(%s) %r", s.docString, eventId)
	}
	s.log.Debug().Str("eventId", eventId).Msg("SetEvent")
	s.setEventTransitions = append(s.setEventTransitions, EventTransition{
		Transition: Transition{},
		eventId:    eventId,
	})
	return s
}

// EventSet Called with the event that was set.
func (s *state) EventSet(eventId string) {
	// Nothing
}

// ClearEvent Create an EventTransition for this state that clears an event.  The current state is returned for method
//
//	chaining. event_id event identifier
func (s *state) ClearEvent(eventId string) State {
	if _debug != nil {
		_debug("clear_event(%s) %r", s.docString, eventId)
	}
	s.log.Debug().Str("eventId", eventId).Msg("ClearEvent")
	s.clearEventTransitions = append(s.clearEventTransitions, EventTransition{
		Transition: Transition{},
		eventId:    eventId,
	})
	return s
}

// WaitEvent Create an EventTransition from this state to another, possibly new, state.  The next state is returned for
//
//	method chaining. pdu tPDU to send next_state state to transition to after sending
func (s *state) WaitEvent(eventId string, nextState State) State {
	if _debug != nil {
		_debug("wait_event(%s) %r next_state=%r", s.docString, eventId, nextState)
	}
	s.log.Debug().Str("eventId", eventId).Msg("WaitEvent")
	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		if _debug != nil {
			_debug("    - new next_state: %r", nextState)
		}
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}
	s.waitEventTransitions = append(s.waitEventTransitions, EventTransition{
		Transition: Transition{
			nextState: nextState,
		},
		eventId: eventId,
	})
	return nextState
}

// Timeout Create a TimeoutTransition from this state to another, possibly new,
//
//	state.  There can only be one timeout transition per state.  The next
//	state is returned for method chaining.
//
//	delay the amount of time to wait for a matching tPDU
//	next_state destination state after timeout
func (s *state) Timeout(delay time.Duration, nextState State) State {
	if _debug != nil {
		_debug("timeout(%s) %r next_state=%r", s.docString, delay, nextState)
	}
	s.log.Debug().Dur("delay", delay).Stringer("nextState", nextState).Msg("Timeout")
	if s.timeoutTransition != nil {
		panic("state already has a timeout")
	}

	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		if _debug != nil {
			_debug("    - new next_state: %r", nextState)
		}
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	now := GetTaskManagerTime()

	s.timeoutTransition = &TimeoutTransition{
		Transition: Transition{nextState: nextState},
		timeout:    now.Add(delay),
	}
	return nextState
}

// Call Create a CallTransition from this state to another, possibly new, state.  The next state is returned for method
//
//	chaining. criteria tPDU to match next_state destination state after a successful match
func (s *state) Call(fn GenericFunction, args Args, kwArgs KWArgs) State {
	if _debug != nil {
		_debug("call(%s) %r %r %r", s.docString, fn, args, kwArgs)
	}
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Call")
	if s.callTransition != nil {
		panic("state already has a 'Call' per state")
	}
	var nextState State
	if _nextState, ok := kwArgs["next_state"]; ok {
		nextState = _nextState.(State)
		if _debug != nil {
			_debug("    - next_state: %r", nextState)
		}
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
		delete(kwArgs, "next_state")
	}

	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		if _debug != nil {
			_debug("    - new next_state: %r", nextState)
		}
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	fnargs := fnargs{
		fn:     fn,
		args:   args,
		kwArgs: kwArgs,
	}
	if _debug != nil {
		_debug("    - fnargs: %r", fnargs)
	}
	s.callTransition = &CallTransition{
		Transition: Transition{nextState: nextState},
		fnargs:     fnargs,
	}
	return nextState
}

func (s *state) Equals(other State) bool {
	if s == other {
		return true
	}
	return false
}

func (s *state) Format(state fmt.State, v rune) {
	switch v {
	case 's', 'v', 'r':
		_, _ = fmt.Fprintf(state, "<%s(%s) at %p>", s._leafName, s.docString, s)
	}
}

func (s *state) String() string {
	if s == nil {
		return "<nil>(*state)"
	}
	if IsDebuggingActive() {
		return fmt.Sprintf("%s", s) // Delegate to format
	}
	return fmt.Sprintf("state(doc: %s, successState: %t, isFailState: %t)", s.docString, s.isSuccessState, s.isFailState)
}
