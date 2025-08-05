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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type IOQControllerStates int

const (
	IOQControllerStates_CTRL_IDLE IOQControllerStates = iota
	IOQControllerStates_CTRL_ACTIVE
	IOQControllerStates_CTRL_WAITING
)

func (i IOQControllerStates) String() string {
	switch i {
	case IOQControllerStates_CTRL_IDLE:
		return "IDLE"
	case IOQControllerStates_CTRL_ACTIVE:
		return "ACTIVE"
	case IOQControllerStates_CTRL_WAITING:
		return "WAITING"
	default:
		return "Unknown"
	}
}

type IOControllerRequirements interface {
	Abort(err error) error
	ProcessIO(iocb IOCBContract) error
	CompleteIO(iocb IOCBContract, pdu PDU) error
	AbortIO(iocb IOCBContract, err error) error
}

//go:generate go tool plc4xGenerator -type=IOController -prefix=iocb_
type IOController struct {
	name         string
	requirements IOControllerRequirements `ignore:"true"`

	_leafName string

	log zerolog.Logger
}

func NewIOController(localLog zerolog.Logger, name string, requirements IOControllerRequirements, options ...Option) (*IOController, error) {
	localLog.Debug().Str("name", name).Msg("NewIOController")
	return &IOController{
		// save the name
		name:         name,
		requirements: requirements,
		_leafName:    ExtractLeafName(options, StructName()),
		log:          localLog,
	}, nil
}

// Abort all requests, no default implementation.
func (i *IOController) Abort(err error) error {
	return nil
}

// RequestIO Called by a client to start processing a request.
func (i *IOController) RequestIO(iocb IOCBContract) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("RequestIO")

	// bind the iocb to this controller
	iocb.setIOController(i)

	// hopefully there won't be an error
	var err error

	// change the state
	iocb.setIOState(IOCBState_PENDING)

	// let derived class figure out how to process this
	err = i.requirements.ProcessIO(iocb)

	// if there was an error, abort the request
	if err != nil {
		return i.requirements.AbortIO(iocb, err)
	}
	return nil
}

// ProcessIO Figure out how to respond to this request.  This must be provided by the derived class.
func (i *IOController) ProcessIO(IOCBContract) error {
	panic("IOController must implement ProcessIO()")
}

// ActiveIO Called by a handler to notify the controller that a request is being processed
func (i *IOController) ActiveIO(iocb IOCBContract) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("ActiveIO")

	// requests should be idle or pending before coming active
	if iocb.getIOState() != IOCBState_IDLE && iocb.getIOState() != IOCBState_PENDING {
		return errors.Errorf("invalid state transition (currently %d)", iocb.getIOState())
	}

	// change the state
	iocb.setIOState(IOCBState_ACTIVE)
	return nil
}

// CompleteIO Called by a handler to return data to the client
func (i *IOController) CompleteIO(iocb IOCBContract, apdu PDU) error {
	i.log.Debug().
		Stringer("iocb", iocb).
		Stringer("apdu", apdu).
		Msg("ActiveIO")

	// if it completed, leave it alone
	if iocb.getIOState() == IOCBState_COMPLETED {
		return nil
	}

	// if it already aborted, leave it alone
	if iocb.getIOState() == IOCBState_ABORTED {
		return nil
	}

	// change the state
	iocb.setIOState(IOCBState_COMPLETED)
	iocb.setIOResponse(apdu)

	// notify the client
	iocb.Trigger()

	return nil
}

// AbortIO Called by a handler or a client to abort a transaction
func (i *IOController) AbortIO(iocb IOCBContract, err error) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("AbortIO")

	// if it completed, leave it alone
	if iocb.getIOState() == IOCBState_COMPLETED {
		return nil
	}

	// if it already aborted, leave it alone
	if iocb.getIOState() == IOCBState_ABORTED {
		return nil
	}

	// change the state
	iocb.setIOState(IOCBState_ABORTED)
	iocb.setIOError(err)

	// notify the client
	iocb.Trigger()

	return nil
}
