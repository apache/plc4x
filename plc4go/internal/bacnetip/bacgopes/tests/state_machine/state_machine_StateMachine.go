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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// StateMachine GA state machine consisting of states.  Every state machine has a start
//
//	state where the state machine begins when it is started.  It also has
//	an *unexpected receive* fail state where the state machine goes if
//	there is an unexpected (unmatched) tPDU received.
type StateMachine interface {
	StateMachineContract
	StateMachineRequirements
}

// StateMachineContract provides a set of functions which can be overwritten by a sub struct
type StateMachineContract interface {
	fmt.Stringer
	fmt.Formatter
	utils.Serializable
	NewState(string) State
	UnexpectedReceive(pdu PDU)
	BeforeSend(pdu PDU)
	AfterSend(pdu PDU)
	BeforeReceive(pdu PDU)
	Receive(args Args, kwArgs KWArgs) error
	AfterReceive(pdu PDU)
	EventSet(id string)
	Run() error
	Reset()
	GetTransactionLog() []TransactionLogEntry
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

// StateMachineRequirements provides a set of functions which must be overwritten by a sub struct
type StateMachineRequirements interface {
	StateMachineContract
	Send(args Args, kwArgs KWArgs) error
}

//go:generate go tool plc4xGenerator -type=stateMachine -prefix=state_machine_
type stateMachine struct {
	requirements StateMachineRequirements `ignore:"true"`

	interceptor    StateInterceptor `asPtr:"true"`
	stateDecorator func(state State) State

	states                 []State
	name                   string
	machineGroup           *StateMachineGroup `asPtr:"true"`
	stateSubStruct         any
	startState             State
	unexpectedReceiveState State
	transitionQueue        chan PDU
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
	transactionLog         []TransactionLogEntry

	_leafName string

	log zerolog.Logger
}

type TransactionLogEntry struct {
	Direction string
	Pdu       PDU
}

func (t TransactionLogEntry) String() string {
	return t.Pdu.String() + t.Pdu.String()
}

var _ StateMachineContract = (*stateMachine)(nil)

// NewStateMachine creates a new state machine. Make sure to call the init function (Init must be called after a new (can't be done in constructor as initialization is then not yet finished))
func NewStateMachine(localLog zerolog.Logger, requirements StateMachineRequirements, options ...Option) (sm StateMachineContract, init func()) {
	s := &stateMachine{
		requirements: requirements,

		log: localLog,
	}
	ApplyAppliers(options, s)
	s._leafName = ExtractLeafName(options, StructName())
	if _debug != nil {
		_debug("__init__(%s)", s.name)
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
			s.startState.setStateMachine(s.requirements)
		} else {
			s.startState = s.NewState("start")
			s.startState = s.stateDecorator(s.startState)
		}

		if s.unexpectedReceiveState != nil {
			if s.unexpectedReceiveState.getStateMachine() != nil {
				panic("start state already bound to a machine")
			}
			s.states = append(s.states, s.unexpectedReceiveState)
			s.unexpectedReceiveState.setStateMachine(s.requirements)
		} else {
			s.unexpectedReceiveState = s.NewState("unexpected receive").Fail("")
			s.unexpectedReceiveState = s.stateDecorator(s.unexpectedReceiveState)
		}

		s.transitionQueue = make(chan PDU, 100)

		s.stateTimeoutTask = NewTimeoutTask(s.StateTimeout, NoArgs, NoKWArgs())

		if s.timeout != 0 {
			s.timeoutState = s.NewState("state machine timeout").Fail("")
			s.timeoutTask = NewTimeoutTask(s.StateMachineTimeout, NoArgs, NoKWArgs(), WithTaskTime(*s.stateMachineTimeout))
		}
	}
}

func WithStateMachineName(name string) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.name = name })
}

func WithStateMachineStateInterceptor(interceptor StateInterceptor) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.interceptor = interceptor })
}

func WithStateMachineTimeout(timeout time.Duration) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.timeout = timeout })
}

func WithStateMachineStartState(startState State) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.startState = startState })
}

func WithStateMachineUnexpectedReceiveState(unexpectedReceiveState State) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.unexpectedReceiveState = unexpectedReceiveState })
}

func WithStateMachineMachineGroup(machineGroup *StateMachineGroup) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.machineGroup = machineGroup })
}

func WithStateMachineStateDecorator(stateDecorator func(state State) State) GenericApplier[*stateMachine] {
	return WrapGenericApplier(func(stateMachine *stateMachine) { stateMachine.stateDecorator = stateDecorator })
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

func (s *stateMachine) GetTransactionLog() []TransactionLogEntry {
	return s.transactionLog
}

func (s *stateMachine) GetCurrentState() State {
	return s.currentState
}

func (s *stateMachine) GetStartState() State {
	return s.startState
}

func (s *stateMachine) NewState(docString string) State {
	if _debug != nil {
		_debug("new_state(%s) %r %r", s.name, docString) // TODO: implement , state_subclass)
	}
	s.log.Trace().Str("docString", docString).Msg("NewState")
	_state := NewState(s.log, s.requirements, docString, WithStateStateInterceptor(s.interceptor))
	_state = s.stateDecorator(_state)
	if _debug != nil {
		_debug("    - state: %r", _state)
	}
	s.states = append(s.states, _state)
	return _state
}

func (s *stateMachine) Reset() {
	if _debug != nil {
		_debug("reset(%s)", s.name)
	}
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
	s.transactionLog = make([]TransactionLogEntry, 0)

	// we are not starting up
	s.startupFlag = false

	// give all the getStates a chance to reset
	for _, state := range s.states {
		state.Reset()
	}
}

func (s *stateMachine) Run() error {
	if _debug != nil {
		_debug("run(%s)", s.name)
	}
	s.log.Trace().Msg("Run")
	if s.running {
		panic("state machine is running")
	}
	if s.currentState != nil {
		panic("not running but has a current state")
	}

	if s.timeoutTask != nil {
		if _debug != nil {
			_debug("    - schedule runtime limit")
		}
		s.log.Debug().Msg("schedule runtime limit")
		s.timeoutTask.InstallTask(WithInstallTaskOptionsDelta(s.timeout))
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
		s.machineGroup.Started(s.requirements)

		// if it is stopped already, let the group know
		if !s.running {
			s.machineGroup.Stopped(s.requirements)
		}
	}
	return nil
}

// Called when the state machine should no longer be running.
func (s *stateMachine) halt() {
	if _debug != nil {
		_debug("halt(%s)", s.name)
	}
	s.log.Trace().Msg("Halt")
	// make sure we're running
	if !s.running {
		panic("state machine is not running")
	}

	// cancel the timeout task
	if s.timeoutTask != nil {
		if _debug != nil {
			_debug("    - cancel runtime limit")
		}
		s.log.Debug().Msg("cancel runtime limit")
		s.timeoutTask.SuspendTask()
	}

	close(s.transitionQueue)

	// no longer running
	s.running = false
}

// success Called when the state machine has successfully completed.
func (s *stateMachine) success() {
	if _debug != nil {
		_debug("success(%s)", s.name)
	}
	s.log.Trace().Msg("Success")
	isSuccessState := true
	s.isSuccessState = &isSuccessState
}

// success Called when the state machine has successfully completed.
func (s *stateMachine) fail() {
	if _debug != nil {
		_debug("fail(%s)", s.name)
	}
	s.log.Trace().Msg("Fail")
	isFailState := true
	s.isFailState = &isFailState
}

func (s *stateMachine) gotoState(state State) error {
	if _debug != nil {
		_debug("goto_state(%s) %r", s.name, state)
	}
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
	if _debug != nil {
		_debug("    - state entered")
	}
	s.log.Trace().Msg("state entered")

	if s.machineGroup != nil {
		for _, transition := range currentState.getSetEventTransitions() {
			if _debug != nil {
				_debug("    - setting event: %r", transition.eventId)
			}
			s.log.Debug().Str("eventId", transition.eventId).Msg("setting event")
			s.machineGroup.SetEvent(transition.eventId)
		}

		for _, transition := range currentState.getClearEventTransitions() {
			if _debug != nil {
				_debug("    - clearing event: %r", transition.eventId)
			}
			s.log.Debug().Str("eventId", transition.eventId).Msg("clearing event")
			s.machineGroup.ClearEvent(transition.eventId)
		}
	}

	if currentState.IsSuccessState() {
		if _debug != nil {
			_debug("    - success state")
		}
		s.log.Trace().Msg("Success state")
		s.stateTransitioning -= 1

		s.halt()
		s.success()

		if s.machineGroup != nil && s.startupFlag {
			s.machineGroup.Stopped(s.requirements)
		}

		return nil
	}

	if currentState.IsFailState() {
		if _debug != nil {
			_debug("    - fail state")
		}
		s.log.Trace().Msg("Fail state")
		s.stateTransitioning -= 1

		s.halt()
		s.fail()

		if s.machineGroup != nil && s.startupFlag {
			s.machineGroup.Stopped(s.requirements)
		}

		return nil
	}

	var nextState State

	if s.machineGroup != nil {
		didBreak := false
		for _, transition := range currentState.getWaitEventTransitions() {
			if _debug != nil {
				_debug("    - waiting event: %r", transition.eventId)
			}
			s.log.Debug().Str("eventID", transition.eventId).Msg("waiting event")
			if _, ok := s.machineGroup.events[transition.eventId]; ok {
				nextState = transition.nextState
				if _debug != nil {
					_debug("    - next_state: %r", nextState)
				}
				s.log.Debug().Stringer("nextState", nextState).Msg("nextState")
				if nextState != currentState {
					didBreak = true
					break
				}
			}
		}
		if !didBreak {
			if _debug != nil {
				_debug("    - no events already set")
			}
			s.log.Trace().Msg("no events already set")
		}
	} else {
		if _debug != nil {
			_debug("    - not part of a group")
		}
		s.log.Trace().Msg("not part of a group")
	}

	if callTransition := currentState.getCallTransition(); callTransition != nil {
		if _debug != nil {
			_debug("    - calling: %r", currentState.getCallTransition())
		}
		s.log.Debug().Interface("callTransition", callTransition).Msg("calling transition")
		f := callTransition.fnargs
		fn, args, kwArgs := f.fn, f.args, f.kwArgs
		if err := fn(args, kwArgs); err != nil {
			if _debug != nil {
				_debug("    - called, exception: %r", err)
			}
			var assertionError tests.AssertionError
			if !errors.As(err, &assertionError) {
				return err
			}
			s.log.Trace().Err(err).Msg("called exception")
			s.stateTransitioning -= 1

			s.halt()
			s.fail()

			if s.machineGroup != nil && !s.startupFlag {
				s.machineGroup.Stopped(s.requirements)
			}

			return nil
		}
		if _debug != nil {
			_debug("    - called, no exception")
		}

		nextState = callTransition.nextState
		if _debug != nil {
			_debug("    - next_state: %r", nextState)
		}
	} else {
		if _debug != nil {
			_debug("    - no calls")
		}
		s.log.Trace().Msg("no calls")
	}

	if nextState == nil {
		for _, transition := range currentState.getSendTransitions() {
			if _debug != nil {
				_debug("    - sending: %r", transition)
			}
			s.log.Debug().Stringer("transition", transition).Msg("sending transition")
			currentState.getInterceptor().BeforeSend(transition.pdu)
			if err := s.requirements.Send(NA(transition.pdu), NKW()); err != nil {
				return errors.Wrap(err, "failed to send")
			}
			currentState.getInterceptor().AfterSend(transition.pdu)

			nextState = transition.nextState
			if _debug != nil {
				_debug("    - next_state: %r", nextState)
			}
			s.log.Debug().Stringer("nextState", nextState).Msg("nextState")

			if nextState != currentState {
				break
			}
		}
	}

	if nextState == nil {
		if _debug != nil {
			_debug("    - nowhere to go")
		}
		s.log.Trace().Msg("nowhere to go")
	} else if nextState == s.currentState {
		if _debug != nil {
			_debug("    - going nowhere")
		}
		s.log.Trace().Msg("going nowhere")
	} else {
		if _debug != nil {
			_debug("    - going")
		}
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
				if _debug != nil {
					_debug("    - pdu: %r", pdu)
				}
				if err := s.Receive(NA(pdu), NKW()); err != nil {
					return errors.Wrap(err, "failed to receive")
				}
			default:
				break queueRead
			}
		}
	}
	return nil
}

func (s *stateMachine) BeforeSend(pdu PDU) {
	s.transactionLog = append(s.transactionLog, TransactionLogEntry{">>>", pdu})
}

func (s *stateMachine) AfterSend(pdu PDU) {
}

func (s *stateMachine) BeforeReceive(pdu PDU) {
	s.transactionLog = append(s.transactionLog, TransactionLogEntry{"<<< %s", pdu})
}

func (s *stateMachine) Receive(args Args, kwArgs KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Receive")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("receive(%s) %r", s.name, pdu)
	}

	// check to see if haven't started yet or we are transitioning
	if s.currentState == nil || s.stateTransitioning != 0 {
		if _debug != nil {
			_debug("    - queue for later")
		}
		s.log.Trace().Msg("queue for later")
		s.transitionQueue <- pdu
		return nil
	}

	// if this is not running it already completed
	if !s.running {
		if _debug != nil {
			_debug("    - already completed")
		}
		s.log.Trace().Msg("already completed")
		return nil
	}

	// reference the current state
	currentState := s.currentState
	if _debug != nil {
		_debug("    - current_state: %r", currentState)
	}
	s.log.Debug().Stringer("currentState", currentState).Msg("current_state")

	currentState.getInterceptor().BeforeReceive(pdu)

	var nextState State
	matchFound := false
	for _, transition := range currentState.getReceiveTransitions() {
		if s.MatchPDU(pdu, transition.criteria) {
			if _debug != nil {
				_debug("    - match found")
			}
			s.log.Trace().Msg("match found")
			matchFound = true

			currentState.getInterceptor().AfterReceive(pdu)

			nextState = transition.nextState
			if _debug != nil {
				_debug("    - next_state: %r", nextState)
			}
			s.log.Debug().Stringer("nextState", nextState).Msg("nextState")

			if nextState != currentState {
				break
			}
		} else {
			if _debug != nil {
				_debug("    - no matches")
			}
			s.log.Trace().Msg("no matches")
		}
	}

	if !matchFound {
		if _debug != nil {
			_debug("    - unexpected")
		}
		currentState.getInterceptor().UnexpectedReceive(pdu)
	} else if nextState != currentState {
		if _debug != nil {
			_debug("    - going")
		}
		if err := s.gotoState(nextState); err != nil {
			return errors.Wrap(err, "error going to state")
		}
	}
	return nil
}

func (s *stateMachine) AfterReceive(pdu PDU) {
	s.log.Trace().Stringer("pdu", pdu).Msg("AfterReceive")
}

func (s *stateMachine) UnexpectedReceive(pdu PDU) {
	if _debug != nil {
		_debug("unexpected_receive(%s) %r", s.name, pdu)
		_debug("    - current_state: %r", s.currentState)
	}
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
	if _debug != nil {
		_debug("event_set(%s) %r", s.name, eventId)
	}
	s.log.Debug().Str("eventId", eventId).Msg("EventSet")
	if !s.running {
		if _debug != nil {
			_debug("    - not running")
		}
		s.log.Trace().Msg("not running")
		return
	}

	if s.stateTransitioning == 1 {
		if _debug != nil {
			_debug("    - transitioning")
		}
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
			if _debug != nil {
				_debug("    - match found")
			}
			s.log.Trace().Msg("match found")
			matchFound = true

			currentState.EventSet(eventId)

			nextState = transition.nextState
			if _debug != nil {
				_debug("    - next_state: %r", nextState)
			}
			s.log.Debug().Stringer("nextState", nextState).Msg("nextState")

			if nextState != currentState {
				break
			}
		}
	}
	if len(currentState.getWaitEventTransitions()) == 0 {
		if _debug != nil {
			_debug("    - going nowhere")
		}
		s.log.Trace().Msg("going nowhere")
	}

	if matchFound && nextState != currentState {
		if _debug != nil {
			_debug("    - going")
		}
		s.log.Trace().Msg("going")
		if err := s.gotoState(nextState); err != nil {
			s.log.Error().Err(err).Msg("failed to go to next state")
		}
	}
}

func (s *stateMachine) StateTimeout(_ Args, _ KWArgs) error {
	if _debug != nil {
		_debug("state_timeout(%s)", s.name)
	}
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

func (s *stateMachine) StateMachineTimeout(_ Args, _ KWArgs) error {
	if _debug != nil {
		_debug("state_machine_timeout(%s)", s.name)
	}
	s.log.Trace().Msg("StateMachineTimeout")
	if !s.running {
		return errors.New("state machine is not running")
	}
	if err := s.gotoState(s.timeoutState); err != nil {
		return errors.Wrap(err, "failed to go to next state")
	}
	return nil
}

func (s *stateMachine) MatchPDU(pdu PDU, criteria criteria) bool {
	if _debug != nil {
		_debug("match_pdu(%s) %r %r", s.name, pdu, criteria)
	}
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

func (s *stateMachine) Format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		_, _ = state.Write([]byte(s.String()))
	case 'r':
		alternateString, _ := s.AlternateString()
		_, _ = state.Write([]byte(alternateString))
	}
}

func (s *stateMachine) AlternateString() (string, bool) {
	if IsDebuggingActive() {
		var stateText = ""
		if s.currentState == nil {
			stateText = "not started"
		} else if s.isSuccessState != nil && *s.isSuccessState {
			stateText = "success"
		} else if s.isFailState != nil && *s.isFailState {
			stateText = "fail"
		} else if !s.running {
			stateText = "idle"
		} else {
			stateText = "in"
		}
		if s.currentState != nil {
			stateText += " " + s.currentState.String()
		}

		return fmt.Sprintf("<%s(%s) %s at %p>", s._leafName, s.name, stateText, s), true
	}
	return "", false
}
