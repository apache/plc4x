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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

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
//
//go:generate go tool plc4xGenerator -type=StateMachineGroup -prefix=state_machine_
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
	s := &StateMachineGroup{
		events: map[string]struct{}{},
		log:    localLog,
	}
	if _debug != nil {
		_debug("__init__")
	}
	return s
}

// Append Add a state machine to the end of the list of state machines
func (s *StateMachineGroup) Append(machine StateMachine) {
	if _debug != nil {
		_debug("append %r", machine)
	}
	s.log.Debug().Stringer("stateMachine", machine).Msg("Append")
	if machine.getMachineGroup() != nil {
		panic("state machine group already contains this machine")
	}

	machine.setMachineGroup(s)

	s.stateMachines = append(s.stateMachines, machine)
}

// Remove a state machine from the list of state machines.
func (s *StateMachineGroup) Remove(machine StateMachine) {
	if _debug != nil {
		_debug("remove %r", machine)
	}
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
	if _debug != nil {
		_debug("reset")
	}
	s.log.Trace().Msg("Reset")
	for _, machine := range s.stateMachines {
		if _debug != nil {
			_debug("    - resetting: %r", machine)
		}
		s.log.Debug().Stringer("stateMachine", machine).Msg("Resetting")
		machine.Reset()
	}

	s.isSuccessState = false
	s.isFailState = false

	s.events = make(map[string]struct{})
}

// SetEvent save an event as 'set' and pass it to the state machines to see
//
//	if they are in a state that is waiting for the event.
func (s *StateMachineGroup) SetEvent(id string) {
	if _debug != nil {
		_debug("set_event %r", id)
	}
	s.log.Trace().Str("eventId", id).Msg("SetEvent")
	s.events[id] = struct{}{}
	if _debug != nil {
		_debug("    - event set")
	}

	for _, machine := range s.stateMachines {
		if _debug != nil {
			_debug("    - state_machine: %r", machine)
		}
		s.log.Debug().Stringer("stateMachine", machine).Msg("Setting")
		machine.EventSet(id)
	}
}

// ClearEvent Remove an event from the set of elements that are 'set'.
func (s *StateMachineGroup) ClearEvent(id string) {
	if _debug != nil {
		_debug("clear_event %r", id)
	}
	s.log.Trace().Str("eventId", id).Msg("ClearEvent")
	if _, ok := s.events[id]; ok {
		delete(s.events, id)
		if _debug != nil {
			_debug("    - event cleared")
		}
	} else {
		if _debug != nil {
			_debug("    - noop")
		}
	}
}

// Run Runs all the machines in the group.
func (s *StateMachineGroup) Run() error {
	if _debug != nil {
		_debug("run")
	}
	s.log.Trace().Msg("Run")
	s.startupFlag = true
	s.isRunning = true

	for _, machine := range s.stateMachines {
		if _debug != nil {
			_debug("    - starting: %r", machine)
		}
		s.log.Debug().Stringer("stateMachine", machine).Msg("starting")
		if err := machine.Run(); err != nil {
			return errors.Wrap(err, "failed to start machine")
		}
	}

	s.startupFlag = false
	if _debug != nil {
		_debug("    - all started")
	}
	s.log.Trace().Msg("all started")

	allSuccess, someFailed := s.CheckForSuccess()
	if allSuccess {
		s.Success()
	} else if someFailed {
		s.Fail()
	} else {
		if _debug != nil {
			_debug("    - some still running")
		}
		s.log.Trace().Msg("some still running")
	}
	return nil
}

// Started Called by a state machine in the group when it has completed its
//
//	transition into its starting state.
func (s *StateMachineGroup) Started(machine StateMachine) {
	if _debug != nil {
		_debug("started %r", machine)
	}
	s.log.Debug().Stringer("stateMachine", machine).Msg("started")
}

// Stopped Called by a state machine after it has halted and its Success()
//
//	or fail() method has been called.
func (s *StateMachineGroup) Stopped(machine StateMachine) {
	if _debug != nil {
		_debug("stopped %r", machine)
	}
	s.log.Debug().Stringer("stateMachine", machine).Msg("stopped")
	if s.startupFlag {
		if _debug != nil {
			_debug("    - still starting up")
		}
		s.log.Trace().Msg("still starting up")
		return
	}

	allSuccess, someFailed := s.CheckForSuccess()
	if allSuccess {
		s.Success()
	} else if someFailed {
		s.Fail()
	} else {
		if _debug != nil {
			_debug("    - some still running")
		}
		s.log.Trace().Msg("some still running")
	}
}

// CheckForSuccess Called after all of the machines have started, and each time a
//
//	machine has stopped, to see if the entire group should be considered
//	a Success or fail.
func (s *StateMachineGroup) CheckForSuccess() (allSuccess bool, someFailed bool) {
	if _debug != nil {
		_debug("check_for_success")
	}
	s.log.Trace().Msg("CheckForSuccess")
	allSuccess = true
	someFailed = false

	for _, machine := range s.stateMachines {
		if machine.IsRunning() {
			if _debug != nil {
				_debug("    - running: %r", machine)
			}
			s.log.Trace().Stringer("machine", machine).Msg("running")
			allSuccess = false
			someFailed = false
			break
		}

		if machine.getCurrentState() == nil {
			if _debug != nil {
				_debug("    - not started: %r", machine)
			}
			s.log.Trace().Stringer("machine", machine).Msg("not started")
			allSuccess = false
			someFailed = false
			break
		}

		allSuccess = allSuccess && machine.getCurrentState().IsSuccessState()
		someFailed = someFailed || machine.getCurrentState().IsFailState()
	}
	if _debug != nil {
		_debug("    - all_success: %r", allSuccess)
		_debug("    - some_failed: %r", someFailed)
	}
	s.log.Debug().Bool("allSuccess", allSuccess).Msg("allSuccess")
	s.log.Debug().Bool("someFailed", allSuccess).Msg("someFailed")
	return
}

// Halt halts all of the running machines in the group.
func (s *StateMachineGroup) Halt() {
	if _debug != nil {
		_debug("halt")
	}
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
	if _debug != nil {
		_debug("success")
	}
	s.log.Trace().Msg("Success")
	s.isRunning = false
	s.isSuccessState = true
}

// Fail Called when all of the machines in the group have halted and at
//
//	at least one of them is in a 'fail' final state.
func (s *StateMachineGroup) Fail() {
	if _debug != nil {
		_debug("fail")
	}
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
