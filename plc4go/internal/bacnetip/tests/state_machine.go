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
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Transition Instances of this class are transitions betweeen getStates of a state machine.
type Transition struct {
	nextState State
}

func (t Transition) String() string {
	return fmt.Sprintf("Transition{nextState: %s}", t.nextState)
}

type SendTransition struct {
	Transition
	pdu bacnetip.PDU
}

func (t SendTransition) String() string {
	return fmt.Sprintf("SendTransition{Transition: %s, pdu: %s}", t.Transition, t.pdu)
}

type criteria struct {
	pduType  any
	pduAttrs map[bacnetip.KnownKey]any
}

func (c criteria) String() string {
	return fmt.Sprintf("criteria{%s, %v}", c.pduType, c.pduAttrs)
}

type ReceiveTransition struct {
	Transition
	criteria criteria
}

func (t ReceiveTransition) String() string {
	return fmt.Sprintf("ReceiveTransition{Transition: %s, criteria: %s}", t.Transition, t.criteria)
}

type EventTransition struct {
	Transition
	eventId string
}

func (t EventTransition) String() string {
	return fmt.Sprintf("EventTransition{Transition: %s, eventId: %s}", t.Transition, t.eventId)
}

type TimeoutTransition struct {
	Transition
	timeout time.Time
}

func (t TimeoutTransition) String() string {
	return fmt.Sprintf("TimeoutTransition{Transition: %s, timeout: %s}", t.Transition, t.timeout)
}

type fnargs struct {
	fn     func(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	args   bacnetip.Args
	kwargs bacnetip.KWArgs
}

func (f fnargs) String() string {
	return fmt.Sprintf("fnargs{fn: %t, args: %s, kwargs: %s}", f.fn == nil, f.args, f.kwargs)
}

type CallTransition struct {
	Transition
	fnargs fnargs
}

func (t CallTransition) String() string {
	return fmt.Sprintf("CallTransition{Transition: %s, fnargs: %s}", t.Transition, t.fnargs)
}

func MatchPdu(localLog zerolog.Logger, pdu bacnetip.PDU, pduType any, pduAttrs map[bacnetip.KnownKey]any) bool {
	// check the type
	if pduType != nil && fmt.Sprintf("%T", pdu) != fmt.Sprintf("%T", pduType) {
		localLog.Debug().Type("got", pdu).Type("want", pduType).Msg("failed match, wrong type")
		return false
	}
	for attrName, attrValue := range pduAttrs {
		switch attrName {
		case bacnetip.KWPPDUSource:
			if !pdu.GetPDUSource().Equals(attrValue) {
				localLog.Debug().Msg("source doesn't match")
				return false
			}
		case bacnetip.KWPDUDestination:
			if !pdu.GetPDUDestination().Equals(attrValue) {
				localLog.Debug().Msg("destination doesn't match")
				return false
			}
		case "x": // only used in test cases
			return bytes.Equal(pdu.(interface{ X() []byte }).X(), attrValue.([]byte))
		case "y": // only used in test cases
			return false
		case "a": // only used in test cases
			a := pdu.(interface{ A() int }).A()
			if a == 0 {
				return false
			}
			return a == attrValue.(int)
		case "b": // only used in test cases
			b := pdu.(interface{ B() int }).B()
			if b == 0 {
				return false
			}
			return b == attrValue.(int)
		case bacnetip.KWPDUData:
			got := pdu.GetPduData()
			want := attrValue
			equal := reflect.DeepEqual(got, want)
			if !equal {
				switch want := want.(type) {
				case []byte:
					localLog.Debug().Bytes("got", got).Bytes("want", want).Msg("mismatch")
				default:
					localLog.Debug().Bytes("got", got).Interface("want", want).Msg("mismatch")
				}
			}
			return equal
		case bacnetip.KWWirtnNetwork:
			wirtn, ok := pdu.(*bacnetip.WhoIsRouterToNetwork)
			if !ok {
				return false
			}
			net := wirtn.GetWirtnNetwork()
			if net == nil {
				return false
			}
			return *net == attrValue
		case bacnetip.KWIartnNetworkList:
			iamrtn, ok := pdu.(*bacnetip.IAmRouterToNetwork)
			if !ok {
				return false
			}
			net := iamrtn.GetIartnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				return false
			}
			return slices.Equal(net, uint16s)
		case bacnetip.KWIcbrtnNetwork:
			iamrtn, ok := pdu.(*bacnetip.ICouldBeRouterToNetwork)
			if !ok {
				return false
			}
			return iamrtn.GetDestinationNetworkAddress() == attrValue
		case bacnetip.KWIcbrtnPerformanceIndex:
			iamrtn, ok := pdu.(*bacnetip.ICouldBeRouterToNetwork)
			if !ok {
				return false
			}
			return iamrtn.GetPerformanceIndex() == attrValue
		default:
			panic("implement " + attrName)
		}
	}
	// TODO: implement
	/*
	  #
	    # check for matching attribute values
	    for attr_name, attr_value in pdu_attrs.items():
	        if not hasattr(pdu, attr_name):
	            if _debug: matchbacnetip.PDU._debug("    - failed match, missing attr: %r", attr_name)
	            return False
	        if getattr(pdu, attr_name) != attr_value:
	            if _debug: stateMachine._debug("    - failed match, attr value: %r, %r", attr_name, attr_value)
	            return False
	*/
	localLog.Trace().Msg("successful match")
	return true
}

type TimeoutTask struct {
	*bacnetip.OneShotTask

	fn     func(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	args   bacnetip.Args
	kwargs bacnetip.KWArgs
}

func NewTimeoutTask(fn func(args bacnetip.Args, kwargs bacnetip.KWArgs) error, args bacnetip.Args, kwargs bacnetip.KWArgs, when *time.Time) *TimeoutTask {
	task := &TimeoutTask{
		fn:     fn,
		args:   args,
		kwargs: kwargs,
	}
	task.OneShotTask = bacnetip.NewOneShotTask(task, when)
	return task
}

func (t *TimeoutTask) ProcessTask() error {
	return t.fn(t.args, t.kwargs)
}

func (t *TimeoutTask) String() string {
	return fmt.Sprintf("TimeoutTask(%v, fn: %t, args: %s, kwargs: %s)", t.Task, t.fn != nil, t.args, t.kwargs)
}

type StateInterceptor interface {
	BeforeSend(pdu bacnetip.PDU)
	AfterSend(pdu bacnetip.PDU)
	BeforeReceive(pdu bacnetip.PDU)
	AfterReceive(pdu bacnetip.PDU)
	UnexpectedReceive(pdu bacnetip.PDU)
}

type State interface {
	fmt.Stringer

	Send(pdu bacnetip.PDU, nextState State) State
	Receive(args bacnetip.Args, kwargs bacnetip.KWArgs) State
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
	Call(fn func(args bacnetip.Args, kwargs bacnetip.KWArgs) error, args bacnetip.Args, kwargs bacnetip.KWArgs) State

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

	log zerolog.Logger
}

func (s *state) Equals(other State) bool {
	if s == other {
		return true
	}
	return false
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

func NewState(localLog zerolog.Logger, stateMachine StateMachine, docString string, opts ...func(state *state)) State {
	s := &state{
		stateMachine: stateMachine,
		docString:    docString,

		log: localLog.With().Str("docString", docString).Logger(),
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.interceptor == nil {
		s.interceptor = s
	}
	return s
}

func WithStateStateInterceptor(interceptor StateInterceptor) func(state *state) {
	return func(state *state) {
		state.interceptor = interceptor
	}
}

func (s *state) String() string {
	if s == nil {
		return "<nil>(*state)"
	}
	return fmt.Sprintf("state(doc: %s, successState: %t, isFailState: %t)", s.docString, s.isSuccessState, s.isFailState)
}

// Reset Override this method in a derived class if the state maintains counters or other information.  Called when the
//
//	associated state machine is Reset.
func (s *state) Reset() {
	s.log.Trace().Msg("Reset")
}

// Doc Change the documentation string (label) for the state.  The state
//
//	is returned for method chaining.
func (s *state) Doc(docString string) State {
	s.log.Debug().Str("docString", docString).Msg("Doc")
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
	s.log.Debug().Msg("EnterState")
	if s.timeoutTransition != nil {
		s.log.Debug().Time("timeout", s.timeoutTransition.timeout).Msg("waiting")
		s.stateMachine.getStateTimeoutTask().InstallTask(bacnetip.InstallTaskOptions{When: &s.timeoutTransition.timeout})
	} else {
		s.log.Trace().Msg("no timeout")
	}
}

// ExitState Called when the state machine is existing the state.
func (s *state) ExitState() {
	s.log.Debug().Msg("ExitState")
	if s.timeoutTransition != nil {
		s.log.Trace().Msg("canceling timeout")
		s.stateMachine.getStateTimeoutTask().SuspendTask()
	}
}

// Send Create a SendTransition from this state to another, possibly new, state.  The next state is returned for method
//
//	chaining. pdu tPDU to send nextState state to transition to after sending
func (s *state) Send(pdu bacnetip.PDU, nextState State) State {
	s.log.Debug().Stringer("pdu", pdu).Msg("Send")
	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		s.log.Debug().Stringer("nextState", nextState).Msg("new nextState")
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
func (s *state) BeforeSend(pdu bacnetip.PDU) {
	s.stateMachine.BeforeSend(pdu)
}

// AfterSend Called after each tPDU about to be sent.
func (s *state) AfterSend(pdu bacnetip.PDU) {
	s.stateMachine.AfterSend(pdu)
}

// Receive Create a ReceiveTransition from this state to another, possibly new,
//
//	state.  The next state is returned for method chaining.
//
//	criteria tPDU to match
//	 next_state destination state after a successful match
func (s *state) Receive(args bacnetip.Args, kwargs bacnetip.KWArgs) State {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Receive")
	pduType := args[0]
	pduAttrs := kwargs
	var nextState State
	if _nextState, ok := pduAttrs["next_state"]; ok {
		nextState = _nextState.(State)
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	}

	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	// add this to the list of transitions
	criteria := criteria{
		pduType:  pduType,
		pduAttrs: pduAttrs,
	}
	s.log.Debug().Interface("criteria", criteria).Msg("criteria")
	s.receiveTransitions = append(s.receiveTransitions, ReceiveTransition{
		Transition: Transition{nextState: nextState},
		criteria:   criteria,
	})

	return nextState
}

// BeforeReceive Called with each tPDU received before matching.
func (s *state) BeforeReceive(pdu bacnetip.PDU) {
	s.stateMachine.BeforeReceive(pdu)
}

// AfterReceive Called with tPDU received after match.
func (s *state) AfterReceive(pdu bacnetip.PDU) {
	s.stateMachine.AfterReceive(pdu)
}

// Ignore Create a ReceiveTransition from this state to itself, if match is successful the effect is to Ignore the tPDU.
//
//	criteria tPDU to match
func (s *state) Ignore(pduType any, pduAttrs map[bacnetip.KnownKey]any) State {
	s.log.Debug().Interface("pduType", pduType).Interface("pduAttrs", pduAttrs).Msg("Ignore")
	s.receiveTransitions = append(s.receiveTransitions, ReceiveTransition{
		Transition: Transition{},
		criteria: criteria{
			pduType:  pduType,
			pduAttrs: pduAttrs,
		},
	})

	return s
}

// UnexpectedReceive Called with tPDU that did not match.  Unless this is trapped by the state, the default behaviour is
//
//	to fail.
func (s *state) UnexpectedReceive(pdu bacnetip.PDU) {
	s.log.Debug().Stringer("pdu", pdu).Msg("UnexpectedReceive")
	s.stateMachine.UnexpectedReceive(pdu)
}

// SetEvent Create an EventTransition for this state that sets an event.  The current state is returned for method
//
//	chaining. event_id event identifier
func (s *state) SetEvent(eventId string) State {
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
	s.log.Debug().Str("eventId", eventId).Msg("WaitEvent")
	if nextState == nil {
		nextState = s.stateMachine.NewState("")
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
	s.log.Debug().Dur("delay", delay).Stringer("nextState", nextState).Msg("Timeout")
	if s.timeoutTransition != nil {
		panic("state already has a timeout")
	}

	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	now := bacnetip.GetTaskManagerTime()

	s.timeoutTransition = &TimeoutTransition{
		Transition: Transition{nextState: nextState},
		timeout:    now.Add(delay),
	}
	return nextState
}

// Call Create a CallTransition from this state to another, possibly new, state.  The next state is returned for method
//
//	chaining. criteria tPDU to match next_state destination state after a successful match
func (s *state) Call(fn func(args bacnetip.Args, kwargs bacnetip.KWArgs) error, args bacnetip.Args, kwargs bacnetip.KWArgs) State {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Call")
	if s.callTransition != nil {
		panic("state already has a 'Call' per state")
	}
	var nextState State
	if _nextState, ok := kwargs["next_state"]; ok {
		nextState = _nextState.(State)
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
		delete(kwargs, "next_state")
	}

	if nextState == nil {
		nextState = s.stateMachine.NewState("")
		s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
	} else if !slices.ContainsFunc(s.stateMachine.getStates(), nextState.Equals) {
		panic("off the rails")
	}

	s.callTransition = &CallTransition{
		Transition: Transition{nextState: nextState},
		fnargs: fnargs{
			fn:     fn,
			args:   args,
			kwargs: kwargs,
		},
	}
	return nextState
}

type StateMachineRequirements interface {
	Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

// StateMachine A state machine consisting of states.  Every state machine has a start
//
//	state where the state machine begins when it is started.  It also has
//	an *unexpected receive* fail state where the state machine goes if
//	there is an unexpected (unmatched) tPDU received.
type StateMachine interface {
	fmt.Stringer

	NewState(string) State
	UnexpectedReceive(pdu bacnetip.PDU)
	BeforeSend(pdu bacnetip.PDU)
	Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	AfterSend(pdu bacnetip.PDU)
	BeforeReceive(pdu bacnetip.PDU)
	Receive(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	AfterReceive(pdu bacnetip.PDU)
	EventSet(id string)
	Run() error
	Reset()
	GetTransactionLog() []string
	GetCurrentState() State
	GetUnexpectedReceiveState() State
	GetStartState() State
	IsRunning() bool
	IsSuccessState() bool
	IsFailState() bool

	getStateTimeoutTask() *TimeoutTask
	getStates() []State
	getCurrentState() State
	setMachineGroup(machineGroup *StateMachineGroup)
	getMachineGroup() *StateMachineGroup
	halt()
}

type stateMachine struct {
	StateMachineRequirements

	interceptor    StateInterceptor
	stateDecorator func(state State) State

	states                 []State
	name                   string
	machineGroup           *StateMachineGroup
	stateSubStruct         any
	startState             State
	unexpectedReceiveState State
	transitionQueue        chan bacnetip.PDU
	stateTimeoutTask       *TimeoutTask
	timeout                time.Duration
	timeoutState           State
	stateMachineTimeout    *time.Time
	timeoutTask            *TimeoutTask
	running                bool
	startupFlag            bool
	isSuccessState         *bool
	isFailState            *bool
	stateTransitioning     int
	currentState           State
	transactionLog         []string

	log zerolog.Logger
}

// NewStateMachine creates a new state machine. Make sure to call the init function (Init must be called after a new (can't be done in constructor as initialization is then not yet finished))
func NewStateMachine(localLog zerolog.Logger, stateMachineRequirements StateMachineRequirements, opts ...func(machine *stateMachine)) (sm StateMachine, init func()) {
	s := &stateMachine{
		StateMachineRequirements: stateMachineRequirements,

		log: localLog,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.name != "" {
		s.log = s.log.With().Str("name", s.name).Logger()
	}
	if s.stateDecorator == nil {
		s.stateDecorator = func(state State) State {
			return state
		}
	}
	return s, func() {
		s.Reset()

		if s.startState != nil {
			if s.startState.getStateMachine() != nil {
				panic("start state already bound to a machine")
			}
			s.states = append(s.states, s.startState)
			s.startState.setStateMachine(s)
		} else {
			s.startState = s.NewState("start")
			s.startState = s.stateDecorator(s.startState)
		}

		if s.unexpectedReceiveState != nil {
			if s.unexpectedReceiveState.getStateMachine() != nil {
				panic("start state already bound to a machine")
			}
			s.states = append(s.states, s.unexpectedReceiveState)
			s.unexpectedReceiveState.setStateMachine(s)
		} else {
			s.unexpectedReceiveState = s.NewState("unexpected receive").Fail("")
			s.unexpectedReceiveState = s.stateDecorator(s.unexpectedReceiveState)
		}

		s.transitionQueue = make(chan bacnetip.PDU, 100)

		s.stateTimeoutTask = NewTimeoutTask(s.StateTimeout, bacnetip.NoArgs, bacnetip.NoKWArgs, nil)

		if s.timeout != 0 {
			s.timeoutState = s.NewState("state machine timeout").Fail("")
			s.timeoutTask = NewTimeoutTask(s.StateMachineTimeout, bacnetip.NoArgs, bacnetip.NoKWArgs, s.stateMachineTimeout)
		}
	}
}

func WithStateMachineName(name string) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.name = name
	}
}

func WithStateMachineStateInterceptor(interceptor StateInterceptor) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.interceptor = interceptor
	}
}

func WithStateMachineTimeout(timeout time.Duration) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.timeout = timeout
	}
}

func WithStateMachineStartState(startState State) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.startState = startState
	}
}

func WithStateMachineUnexpectedReceiveState(unexpectedReceiveState State) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.unexpectedReceiveState = unexpectedReceiveState
	}
}

func WithStateMachineMachineGroup(machineGroup *StateMachineGroup) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.machineGroup = machineGroup
	}
}

func WithStateMachineStateDecorator(stateDecorator func(state State) State) func(stateMachine *stateMachine) {
	return func(stateMachine *stateMachine) {
		stateMachine.stateDecorator = stateDecorator
	}
}

func (s *stateMachine) getStateTimeoutTask() *TimeoutTask {
	return s.stateTimeoutTask
}

func (s *stateMachine) getStates() []State {
	return s.states
}

func (s *stateMachine) getMachineGroup() *StateMachineGroup {
	return s.machineGroup
}

func (s *stateMachine) getCurrentState() State {
	return s.currentState
}

func (s *stateMachine) setMachineGroup(machineGroup *StateMachineGroup) {
	s.machineGroup = machineGroup
}

func (s *stateMachine) GetTransactionLog() []string {
	return s.transactionLog
}

func (s *stateMachine) GetCurrentState() State {
	return s.currentState
}

func (s *stateMachine) GetStartState() State {
	return s.startState
}

func (s *stateMachine) String() string {
	return fmt.Sprintf("stateMachine(name=%s)", s.name)
}

func (s *stateMachine) NewState(docString string) State {
	s.log.Trace().Str("docString", docString).Msg("NewState")
	_state := NewState(s.log, s, docString, WithStateStateInterceptor(s.interceptor))
	_state = s.stateDecorator(_state)
	s.states = append(s.states, _state)
	return _state
}

func (s *stateMachine) Reset() {
	s.log.Trace().Msg("Reset")
	// make sure we're not running
	if s.running {
		panic("state machine is running")
	}

	// flags for remembering Success or fail
	s.isSuccessState = nil
	s.isFailState = nil

	// no current state, empty transaction log
	s.currentState = nil
	s.transactionLog = make([]string, 0)

	// we are not starting up
	s.startupFlag = false

	// give all the getStates a chance to reset
	for _, state := range s.states {
		state.Reset()
	}
}

func (s *stateMachine) Run() error {
	s.log.Trace().Msg("Run")
	if s.running {
		panic("state machine is running")
	}
	if s.currentState != nil {
		panic("not running but has a current state")
	}

	if s.timeoutTask != nil {
		s.log.Debug().Msg("schedule runtime limit")
		s.timeoutTask.InstallTask(bacnetip.InstallTaskOptions{Delta: &s.timeout})
	}

	// we are starting up
	s.startupFlag = true

	// go to the start state
	if err := s.gotoState(s.startState); err != nil {
		return errors.Wrap(err, "error going to start state")
	}

	// startup complete
	s.startupFlag = false

	// if it is part of a group, let the group know
	if s.machineGroup != nil {
		s.machineGroup.Started(s)

		// if it is stopped already, let the group know
		if !s.running {
			s.machineGroup.Stopped(s)
		}
	}
	return nil
}

// Called when the state machine should no longer be running.
func (s *stateMachine) halt() {
	s.log.Trace().Msg("Halt")
	// make sure we're running
	if !s.running {
		panic("state machine is not running")
	}

	// cancel the timeout task
	if s.timeoutTask != nil {
		s.log.Debug().Msg("cancel runtime limit")
		s.timeoutTask.SuspendTask()
	}

	close(s.transitionQueue)

	// no longer running
	s.running = false
}

// success Called when the state machine has successfully completed.
func (s *stateMachine) success() {
	s.log.Trace().Msg("Success")
	isSuccessState := true
	s.isSuccessState = &isSuccessState
}

// success Called when the state machine has successfully completed.
func (s *stateMachine) fail() {
	s.log.Trace().Msg("Fail")
	isFailState := true
	s.isFailState = &isFailState
}

func (s *stateMachine) gotoState(state State) error {
	s.log.Debug().Stringer("state", state).Msg("gotoState")
	//where do you think you're going?
	if !slices.ContainsFunc(s.states, state.Equals) {
		return errors.New("off the rails")
	}

	s.stateTransitioning += 1

	if s.currentState != nil {
		s.currentState.ExitState()
	} else if state == s.startState {
		// starting  up
		s.running = true
	} else {
		return errors.New("start at the start state")
	}

	s.currentState = state
	currentState := state

	currentState.EnterState()
	s.log.Trace().Msg("state entered")

	if s.machineGroup != nil {
		for _, transition := range currentState.getSetEventTransitions() {
			s.log.Debug().Str("eventId", transition.eventId).Msg("setting event")
			s.machineGroup.SetEvent(transition.eventId)
		}

		for _, transition := range currentState.getClearEventTransitions() {
			s.log.Debug().Str("eventId", transition.eventId).Msg("clearing event")
			s.machineGroup.ClearEvent(transition.eventId)
		}
	}

	if currentState.IsSuccessState() {
		s.log.Trace().Msg("Success state")
		s.stateTransitioning -= 1

		s.halt()
		s.success()

		if s.machineGroup != nil && s.startupFlag {
			s.machineGroup.Stopped(s)
		}

		return nil
	}

	if currentState.IsFailState() {
		s.log.Trace().Msg("Fail state")
		s.stateTransitioning -= 1

		s.halt()
		s.fail()

		if s.machineGroup != nil && s.startupFlag {
			s.machineGroup.Stopped(s)
		}

		return nil
	}

	var nextState State

	if s.machineGroup != nil {
		didBreak := false
		for _, transition := range currentState.getWaitEventTransitions() {
			s.log.Debug().Str("eventID", transition.eventId).Msg("waiting event")
			if _, ok := s.machineGroup.events[transition.eventId]; ok {
				nextState = transition.nextState
				s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
				if nextState != currentState {
					didBreak = true
					break
				}
			}
		}
		if !didBreak {
			s.log.Trace().Msg("no events already set")
		}
	} else {
		s.log.Trace().Msg("not part of a group")
	}

	if callTransition := currentState.getCallTransition(); callTransition != nil {
		s.log.Debug().Interface("callTransition", callTransition).Msg("calling transition")
		f := callTransition.fnargs
		fn, args, kwargs := f.fn, f.args, f.kwargs
		if err := fn(args, kwargs); err != nil {
			var assertionError AssertionError
			if !errors.As(err, &assertionError) {
				return err
			}
			s.log.Trace().Err(err).Msg("called exception")
			s.stateTransitioning -= 1

			s.halt()
			s.fail()

			if s.machineGroup != nil && !s.startupFlag {
				s.machineGroup.Stopped(s)
			}

			return nil
		} else {
			s.log.Trace().Msg("called, exception")
		}

		nextState = callTransition.nextState
	} else {
		s.log.Trace().Msg("no calls")
	}

	if nextState == nil {
		for _, transition := range currentState.getSendTransitions() {
			s.log.Debug().Stringer("transition", transition).Msg("sending transition")
			currentState.getInterceptor().BeforeSend(transition.pdu)
			if err := s.StateMachineRequirements.Send(bacnetip.NewArgs(transition.pdu), bacnetip.NewKWArgs()); err != nil {
				return errors.Wrap(err, "failed to send")
			}
			currentState.getInterceptor().AfterSend(transition.pdu)

			nextState = transition.nextState
			s.log.Debug().Stringer("nextState", nextState).Msg("nextState")

			if nextState != currentState {
				break
			}
		}
	}

	if nextState == nil {
		s.log.Trace().Msg("nowhere to go")
	} else if nextState == s.currentState {
		s.log.Trace().Msg("going nowhere")
	} else {
		s.log.Trace().Msg("going")
		if err := s.gotoState(nextState); err != nil {
			return errors.Wrap(err, "error in recursion")
		}
	}

	s.stateTransitioning -= 1

	if s.stateTransitioning == 0 {
	queueRead:
		for s.running {
			select {
			case pdu := <-s.transitionQueue:
				if err := s.Receive(bacnetip.NewArgs(pdu), bacnetip.NewKWArgs()); err != nil {
					return errors.Wrap(err, "failed to receive")
				}
			default:
				break queueRead
			}
		}
	}
	return nil
}

func (s *stateMachine) BeforeSend(pdu bacnetip.PDU) {
	s.transactionLog = append(s.transactionLog, fmt.Sprintf("<<<%v", pdu))
}

func (s *stateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	panic("not implemented")
}

func (s *stateMachine) AfterSend(pdu bacnetip.PDU) {
}

func (s *stateMachine) BeforeReceive(pdu bacnetip.PDU) {
	s.transactionLog = append(s.transactionLog, fmt.Sprintf(">>>%v", pdu))
}

func (s *stateMachine) Receive(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Receive")
	pdu := args.Get0PDU()
	if s.currentState == nil || s.stateTransitioning != 0 {
		s.log.Trace().Msg("queue for later")
		s.transitionQueue <- pdu
		return nil
	}

	if !s.running {
		s.log.Trace().Msg("already completed")
		return nil
	}

	currentState := s.currentState
	s.log.Debug().Stringer("currentState", currentState).Msg("current_state")

	currentState.getInterceptor().BeforeReceive(pdu)

	var nextState State
	matchFound := false
	for _, transition := range currentState.getReceiveTransitions() {
		if s.MatchPDU(pdu, transition.criteria) {
			s.log.Trace().Msg("match found")
			matchFound = true

			currentState.getInterceptor().AfterReceive(pdu)

			nextState = transition.nextState
			s.log.Debug().Stringer("nextState", nextState).Msg("nextState")

			if nextState != currentState {
				break
			}
		} else {
			s.log.Trace().Msg("no matches")
		}
	}

	if !matchFound {
		currentState.getInterceptor().UnexpectedReceive(pdu)
	} else if nextState != currentState {
		if err := s.gotoState(nextState); err != nil {
			return errors.Wrap(err, "error going to state")
		}
	}
	return nil
}

func (s *stateMachine) AfterReceive(pdu bacnetip.PDU) {
	s.log.Trace().Stringer("pdu", pdu).Msg("AfterReceive")
}

func (s *stateMachine) UnexpectedReceive(pdu bacnetip.PDU) {
	s.log.Trace().Stringer("pdu", pdu).Msg("UnexpectedReceive")
	s.log.Trace().Stringer("currentState", s.currentState).Msg("currentState")
	if err := s.gotoState(s.unexpectedReceiveState); err != nil {
		s.log.Error().Err(err).Msg("error going to unexpected state")
	}
}

func (s *stateMachine) GetUnexpectedReceiveState() State {
	return s.unexpectedReceiveState
}

func (s *stateMachine) EventSet(eventId string) {
	s.log.Debug().Str("eventId", eventId).Msg("EventSet")
	if !s.running {
		s.log.Trace().Msg("not running")
		return
	}

	if s.stateTransitioning == 1 {
		s.log.Trace().Msg("transitioning")
		return
	}
	if s.currentState == nil {
		panic("current state is nil")
	}
	currentState := s.currentState

	var nextState State
	matchFound := false
	for _, transition := range currentState.getWaitEventTransitions() {
		if transition.eventId == eventId {
			s.log.Trace().Msg("match found")
			matchFound = true

			currentState.EventSet(eventId)

			nextState = transition.nextState
			s.log.Debug().Stringer("nextState", nextState).Msg("nextState")

			if nextState != currentState {
				break
			}
		}
	}
	if len(currentState.getWaitEventTransitions()) == 0 {
		s.log.Trace().Msg("going nowhere")
	}

	if matchFound && nextState != currentState {
		s.log.Trace().Msg("going")
		if err := s.gotoState(nextState); err != nil {
			s.log.Error().Err(err).Msg("failed to go to next state")
		}
	}
}

func (s *stateMachine) StateTimeout(_ bacnetip.Args, _ bacnetip.KWArgs) error {
	s.log.Trace().Msg("StateTimeout")
	if !s.running {
		return errors.New("state machine is not running")
	}
	if s.currentState.getTimeoutTransition() == nil {
		return errors.New("state timeout, but timeout transition is nil")
	}
	if err := s.gotoState(s.currentState.getTimeoutTransition().nextState); err != nil {
		return errors.Wrap(err, "failed to go to next state")
	}
	return nil
}

func (s *stateMachine) StateMachineTimeout(_ bacnetip.Args, _ bacnetip.KWArgs) error {
	s.log.Trace().Msg("StateMachineTimeout")
	if !s.running {
		return errors.New("state machine is not running")
	}
	if err := s.gotoState(s.timeoutState); err != nil {
		return errors.Wrap(err, "failed to go to next state")
	}
	return nil
}

func (s *stateMachine) MatchPDU(pdu bacnetip.PDU, criteria criteria) bool {
	s.log.Debug().Stringer("pdu", pdu).Stringer("criteria", criteria).Msg("MatchPDU")
	return MatchPdu(s.log, pdu, criteria.pduType, criteria.pduAttrs)
}

func (s *stateMachine) IsRunning() bool {
	return s.running
}

func (s *stateMachine) IsSuccessState() bool {
	if s.isSuccessState == nil {
		return false
	}
	return *s.isSuccessState
}

func (s *stateMachine) IsFailState() bool {
	if s.isFailState == nil {
		return false
	}
	return *s.isFailState
}

// StateMachineGroup  A state machine group is a collection of state machines that are all
//
//	started and stopped together.  There are methods available to derived
//	classes that are called when all of the machines in the group have
//	completed, either all successfully or at least one has failed.
//
//	.. note:: When creating a group of state machines, add the ones that
//	    are expecting to receive one or more tPDU's first before the ones
//	    that send tPDU's.  They will be started first, and be ready for the
//	    tPDU that might be sent.
type StateMachineGroup struct {
	stateMachines  []StateMachine
	isSuccessState bool
	isFailState    bool
	events         map[string]struct{}
	startupFlag    bool
	isRunning      bool

	log zerolog.Logger
}

func NewStateMachineGroup(localLog zerolog.Logger) *StateMachineGroup {
	return &StateMachineGroup{
		events: map[string]struct{}{},
		log:    localLog,
	}
}

// Append Add a state machine to the end of the list of state machines
func (s *StateMachineGroup) Append(machine StateMachine) {
	s.log.Debug().Stringer("stateMachine", machine).Msg("Append")
	if machine.getMachineGroup() != nil {
		panic("state machine group already contains this machine")
	}

	machine.setMachineGroup(s)

	s.stateMachines = append(s.stateMachines, machine)
}

// Remove a state machine from the list of state machines.
func (s *StateMachineGroup) Remove(machine StateMachine) {
	s.log.Debug().Stringer("stateMachine", machine).Msg("Remove")
	if machine.getMachineGroup() != s {
		panic("state machine is not a member of this group")
	}

	machine.setMachineGroup(nil)
	for i, stateMachine := range s.stateMachines {
		if stateMachine == machine {
			s.stateMachines = append(s.stateMachines[:i], s.stateMachines[i+1:]...)
			break
		}
	}
}

// Reset resets all the machines in the group.
func (s *StateMachineGroup) Reset() {
	s.log.Trace().Msg("Reset")
	for _, stateMachine := range s.stateMachines {
		s.log.Debug().Stringer("stateMachine", stateMachine).Msg("Resetting")
		stateMachine.Reset()
	}

	s.isSuccessState = false
	s.isFailState = false

	s.events = make(map[string]struct{})
}

// SetEvent save an event as 'set' and pass it to the state machines to see
//
//	if they are in a state that is waiting for the event.
func (s *StateMachineGroup) SetEvent(id string) {
	s.log.Trace().Str("eventId", id).Msg("SetEvent")
	s.events[id] = struct{}{}

	for _, machine := range s.stateMachines {
		s.log.Debug().Stringer("stateMachine", machine).Msg("Setting")
		machine.EventSet(id)
	}
}

// ClearEvent Remove an event from the set of elements that are 'set'.
func (s *StateMachineGroup) ClearEvent(id string) {
	s.log.Trace().Str("eventId", id).Msg("ClearEvent")
	delete(s.events, id)
}

// Run Runs all the machines in the group.
func (s *StateMachineGroup) Run() error {
	s.log.Trace().Msg("Run")
	s.startupFlag = true
	s.isRunning = true

	for _, machine := range s.stateMachines {
		s.log.Debug().Stringer("stateMachine", machine).Msg("starting")
		if err := machine.Run(); err != nil {
			return errors.Wrap(err, "failed to start machine")
		}
	}

	s.startupFlag = false
	s.log.Trace().Msg("all started")

	allSuccess, someFailed := s.CheckForSuccess()
	if allSuccess {
		s.Success()
	} else if someFailed {
		s.Fail()
	} else {
		s.log.Trace().Msg("some still running")
	}
	return nil
}

// Started Called by a state machine in the group when it has completed its
//
//	transition into its starting state.
func (s *StateMachineGroup) Started(machine *stateMachine) {
	s.log.Debug().Stringer("stateMachine", machine).Msg("started")
}

// Stopped Called by a state machine after it has halted and its Success()
//
//	or fail() method has been called.
func (s *StateMachineGroup) Stopped(machine *stateMachine) {
	s.log.Debug().Stringer("stateMachine", machine).Msg("stopped")
	if s.startupFlag {
		s.log.Trace().Msg("still starting up")
		return
	}

	allSuccess, someFailed := s.CheckForSuccess()
	if allSuccess {
		s.Success()
	} else if someFailed {
		s.Fail()
	} else {
		s.log.Trace().Msg("some still running")
	}
}

// CheckForSuccess Called after all of the machines have started, and each time a
//
//	machine has stopped, to see if the entire group should be considered
//	a Success or fail.
func (s *StateMachineGroup) CheckForSuccess() (allSuccess bool, someFailed bool) {
	s.log.Trace().Msg("CheckForSuccess")
	allSuccess = true
	someFailed = false

	for _, machine := range s.stateMachines {
		if machine.IsRunning() {
			s.log.Trace().Stringer("machine", machine).Msg("running")
			allSuccess = false
			someFailed = false
			break
		}

		if machine.getCurrentState() == nil {
			s.log.Trace().Stringer("machine", machine).Msg("not started")
			allSuccess = false
			someFailed = false
			break
		}

		allSuccess = allSuccess && machine.getCurrentState().IsSuccessState()
		someFailed = someFailed || machine.getCurrentState().IsFailState()
	}
	s.log.Debug().Bool("allSuccess", allSuccess).Msg("allSuccess")
	s.log.Debug().Bool("someFailed", allSuccess).Msg("someFailed")
	return
}

// Halt halts all of the running machines in the group.
func (s *StateMachineGroup) Halt() {
	s.log.Trace().Msg("Halt")
	for _, machine := range s.stateMachines {
		if machine.IsRunning() {
			machine.halt()
		}
	}
}

// Success Called when all of the machines in the group have halted and they
//
//	are all in a 'Success' final state.
func (s *StateMachineGroup) Success() {
	s.log.Trace().Msg("Success")
	s.isRunning = false
	s.isSuccessState = true
}

// Fail Called when all of the machines in the group have halted and at
//
//	at least one of them is in a 'fail' final state.
func (s *StateMachineGroup) Fail() {
	s.log.Trace().Msg("Fail")
	s.isRunning = false
	s.isFailState = true
}

func (s *StateMachineGroup) GetStateMachines() []StateMachine {
	return s.stateMachines
}

func (s *StateMachineGroup) IsRunning() bool {
	return s.isRunning
}

func (s *StateMachineGroup) IsSuccessState() bool {
	return s.isSuccessState
}

func (s *StateMachineGroup) IsFailState() bool {
	return s.isFailState
}

// ClientStateMachine An instance of this class sits at the top of a stack.  tPDU's that the
//
//	state machine sends are sent down the stack and tPDU's coming up the
//	stack are fed as received tPDU's.
type ClientStateMachine struct {
	*bacnetip.Client
	StateMachine

	name string

	log zerolog.Logger
}

func NewClientStateMachine(localLog zerolog.Logger, opts ...func(*ClientStateMachine)) (*ClientStateMachine, error) {
	c := &ClientStateMachine{
		log: localLog,
	}
	for _, opt := range opts {
		opt(c)
	}
	var err error
	c.Client, err = bacnetip.NewClient(localLog, c)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	var init func()
	c.StateMachine, init = NewStateMachine(localLog, c, WithStateMachineName(c.name))
	init()
	return c, nil
}

func WithClientStateMachineName(name string) func(*ClientStateMachine) {
	return func(c *ClientStateMachine) {
		c.name = name
	}
}

func (s *ClientStateMachine) String() string {
	return fmt.Sprintf("ClientStateMachine{Client: %v, StateMachine: %v}", s.Client, s.StateMachine)
}

func (s *ClientStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return s.Request(args, kwargs)
}

func (s *ClientStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	return s.Receive(args, kwargs)
}

type ServerStateMachine struct {
	*bacnetip.Server
	StateMachine

	name string

	log zerolog.Logger
}

func NewServerStateMachine(localLog zerolog.Logger, opts ...func(*ServerStateMachine)) (*ServerStateMachine, error) {
	c := &ServerStateMachine{
		log: localLog,
	}
	for _, opt := range opts {
		opt(c)
	}
	var err error
	c.Server, err = bacnetip.NewServer(localLog, c)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Server")
	}
	var init func()
	c.StateMachine, init = NewStateMachine(localLog, c, WithStateMachineName(c.name))
	init()
	return c, nil
}

func WithServerStateMachineName(name string) func(*ServerStateMachine) {
	return func(s *ServerStateMachine) {
		s.name = name
	}
}

func (s *ServerStateMachine) String() string {
	return fmt.Sprintf("ServerStateMachine(TBD...)") // TODO: fill some info here
}

func (s *ServerStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return s.Response(args, kwargs)
}

func (s *ServerStateMachine) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	return s.Receive(args, kwargs)
}

type TrafficLog struct {
	traffic []struct {
		time.Time
		bacnetip.Args
	}
}

// Call Capture the current time and the arguments.
func (t *TrafficLog) Call(args bacnetip.Args) {
	t.traffic = append(t.traffic, struct {
		time.Time
		bacnetip.Args
	}{Time: time.Now(), Args: args})
}

// Dump the traffic, pass the correct handler like SomeClass._debug
func (t *TrafficLog) Dump(handlerFn func(format string, args bacnetip.Args)) {
	if t == nil {
		return
	}
	for _, args := range t.traffic {
		argFormat := "   %6.3f:"
		for _, arg := range args.Args[1:] {
			_ = arg
			argFormat += " %v"
		}
		handlerFn(argFormat, args.Args)
	}
}
