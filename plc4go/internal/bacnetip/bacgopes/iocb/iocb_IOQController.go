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

package iocb

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/core"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

type IOQControllerRequirements interface {
	IOControllerRequirements
	ProcessIO(iocb IOCBContract) error
}

//go:generate go tool plc4xGenerator -type=IOQController -prefix=iocb_
type IOQController struct {
	*IOController `stringer:"true"`

	state      IOQControllerStates
	ActiveIOCB IOCBContract `stringer:"true"`
	IoQueue    *IOQueue     `stringer:"true"`
	waitTime   time.Duration

	log      zerolog.Logger
	stateLog zerolog.Logger //TODO: add option so it can be logged to a different file
}

func NewIOQController(localLog zerolog.Logger, name string, requirements IOQControllerRequirements) (*IOQController, error) {
	i := &IOQController{
		log:      localLog,
		stateLog: localLog,
	}
	var err error
	i.IOController, err = NewIOController(localLog, name, requirements)
	if err != nil {
		return nil, errors.Wrap(err, "error creating IO controller")
	}

	// start idle
	i.state = IOQControllerStates_CTRL_IDLE
	i.log.Debug().
		Timestamp().
		Str("name", name).
		Stringer("state", i.state).
		Msg("creating")

	// no active iocb
	i.ActiveIOCB = nil

	// create an IOQueue for iocb's requested when not idle
	i.IoQueue = NewIOQueue(localLog, name+" queue")

	return i, nil
}

// Abort all pending requests
func (i *IOQController) Abort(err error) error {
	i.log.Debug().Err(err).Msg("Abort")

	if i.state == IOQControllerStates_CTRL_IDLE {
		i.log.Debug().Msg("idle")
		return nil
	}

	for {
		iocb, err := i.IoQueue.Get(false, nil)
		if err != nil {
			return errors.Wrap(err, "error getting something from queue")
		}
		if iocb == nil {
			break
		}
		i.log.Debug().Stringer("iocb", iocb).Msg("working with iocb")

		// change the state
		iocb.setIOState(IOCBState_ABORTED)
		iocb.setIOError(err)

		// notify the client
		iocb.Trigger()
	}

	if i.state != IOQControllerStates_CTRL_IDLE {
		i.log.Debug().Msg("busy after aborts")
	}

	return nil
}

// RequestIO Called by a client to start processing a request
func (i *IOQController) RequestIO(iocb IOCBContract) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("RequestIO")

	// bind the iocb to this controller
	iocb.setIOController(i)

	// if we're busy, queue it
	if i.state != IOQControllerStates_CTRL_IDLE {
		i.log.Debug().
			Stringer("activeIOCB", i.ActiveIOCB).
			Msg("busy, request queued, activeIOCB")

		iocb.setIOState(IOCBState_PENDING)

		if err := i.IoQueue.Put(iocb); err != nil {
			return errors.Wrap(err, "error putting iocb in the queue")
		}
		return nil
	}

	if err := i.requirements.ProcessIO(iocb); err != nil {
		i.log.Debug().Err(err).Msg("ProcessIO error")
		if err := i.Abort(err); err != nil {
			return errors.Wrap(err, "error sending abort")
		}
	}

	return nil
}

// ProcessIO Figure out how to respond to this request.  This must be provided by the derived class
func (i *IOQController) ProcessIO(IOCBContract) error {
	panic("IOController must implement ProcessIO()")
}

// ActiveIO Called by a handler to notify the controller that a request is being processed
func (i *IOQController) ActiveIO(iocb IOCBContract) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("ActiveIO")

	// base class work first, setting iocb state and timer data
	if err := i.IOController.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "error calling super active io")
	}

	// change our state
	i.state = IOQControllerStates_CTRL_ACTIVE
	i.stateLog.Debug().Timestamp().Str("name", i.name).Msg("active")

	// keep track of the iocb
	i.ActiveIOCB = iocb
	return nil
}

// CompleteIO Called by a handler to return data to the client
func (i *IOQController) CompleteIO(iocb IOCBContract, msg PDU) error {
	i.log.Debug().Stringer("iocb", iocb).Stringer("msg", msg).Msg("CompleteIO")

	// check to see if it is completing the active one
	if iocb != i.ActiveIOCB {
		return errors.New("not the current iocb")
	}

	// normal completion
	if err := i.IOController.CompleteIO(iocb, msg); err != nil {
		return errors.Wrap(err, "error completing io")
	}

	// no longer an active iocb
	i.ActiveIOCB = nil

	// check to see if we should wait a bit
	if i.waitTime != 0 {
		// change our state
		i.state = IOQControllerStates_CTRL_WAITING
		i.stateLog.Debug().Timestamp().Str("name", i.name).Msg("waiting")

		task := FunctionTask(i._waitTrigger, NoArgs, NoKWArgs())
		task.InstallTask(WithInstallTaskOptionsDelta(i.waitTime))
	} else {
		// change our state
		i.state = IOQControllerStates_CTRL_IDLE
		i.stateLog.Debug().Timestamp().Str("name", i.name).Msg("idle")

		// look for more to do
		Deferred(i._trigger, NoArgs, NoKWArgs())
	}

	return nil
}

// AbortIO Called by a handler or a client to abort a transaction
func (i *IOQController) AbortIO(iocb IOCBContract, err error) error {
	i.log.Debug().Err(err).Stringer("iocb", iocb).Msg("AbortIO")

	// Normal abort
	if err := i.IOController.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "ActiveIO failed")
	}

	// check to see if it is completing the active one
	if iocb != i.ActiveIOCB {
		i.log.Debug().Msg("not current iocb")
		return nil
	}

	// no longer an active iocb
	i.ActiveIOCB = nil

	// change our state
	i.state = IOQControllerStates_CTRL_IDLE
	i.stateLog.Debug().Timestamp().Str("name", i.name).Msg("idle")

	// look for more to do
	Deferred(i._trigger, NoArgs, NoKWArgs())
	return nil
}

// _trigger Called to launch the next request in the queue
func (i *IOQController) _trigger(_ Args, _ KWArgs) error {
	i.log.Debug().Msg("_trigger")

	// if we are busy, do nothing
	if i.state != IOQControllerStates_CTRL_IDLE {
		i.log.Debug().Msg("not idle")
		return nil
	}

	// if there is nothing to do, return
	if len(i.IoQueue.Queue) == 0 {
		i.log.Debug().Msg("empty queue")
		return nil
	}

	// get the next iocb
	iocb, err := i.IoQueue.Get(false, nil)
	if err != nil {
		panic("this should never happen")
	}

	if err := i.ProcessIO(iocb); err != nil {
		// if there was an error, abort the request
		if err := i.Abort(err); err != nil {
			i.log.Debug().Err(err).Msg("error aborting")
			return nil
		}
		return nil
	}

	// if we're idle, call again
	if i.state == IOQControllerStates_CTRL_IDLE {
		Deferred(i._trigger, NoArgs, NoKWArgs())
	}
	return nil
}

// _waitTrigger is called to launch the next request in the queue
func (i *IOQController) _waitTrigger(_ Args, _ KWArgs) error {
	i.log.Debug().Msg("_waitTrigger")

	// make sure we are waiting
	if i.state != IOQControllerStates_CTRL_WAITING {
		i.log.Debug().Msg("not waiting")
		return nil
	}

	// change our state
	i.state = IOQControllerStates_CTRL_IDLE
	i.stateLog.Debug().Timestamp().Str("name", i.name).Msg("idle")

	// look for more to do
	return i._trigger(Nothing())
}

func (i *IOQController) Close() error {
	if i.IoQueue != nil {
		return i.IoQueue.Close()
	}
	return nil
}
