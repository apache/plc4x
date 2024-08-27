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
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

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
		s.timeoutTask.InstallTask(bacnetip.WithInstallTaskOptionsDelta(s.timeout))
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

func (s *stateMachine) String() string {
	var fields []string
	if s.name != "" {
		fields = append(fields, "name=", s.name)
	}
	fields = append(fields, "running="+strconv.FormatBool(s.running))
	if s.isSuccessState != nil {
		fields = append(fields, "successState="+strconv.FormatBool(*s.isSuccessState))
	}
	if s.isFailState != nil {
		fields = append(fields, "failState="+strconv.FormatBool(*s.isFailState))
	}
	return fmt.Sprintf("StateMachine(%s)", strings.Join(fields, ", "))
}
