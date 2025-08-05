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
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type IOCBState int

const (
	IOCBState_IDLE IOCBState = iota
	IOCBState_PENDING
	IOCBState_ACTIVE
	IOCBState_COMPLETED
	IOCBState_ABORTED
)

func (i IOCBState) String() string {
	switch i {
	case IOCBState_IDLE:
		return "IDLE"
	case IOCBState_PENDING:
		return "PENDING"
	case IOCBState_ACTIVE:
		return "ACTIVE"
	case IOCBState_COMPLETED:
		return "COMPLETED"
	case IOCBState_ABORTED:
		return "ABORTED"
	default:
		return "Unknown"
	}
}

// IOCBContract provides a set of functions which can be overwritten by a sub struct
type IOCBContract interface {
	fmt.Stringer
	setIOController(ioController IOControllerRequirements)
	setIOState(newState IOCBState)
	getIOState() IOCBState
	setIOResponse(msg PDU)
	Trigger()
	setIOError(err error)
	getRequest() PDU
	GetDestination() *Address
	getPriority() int
	clearQueue()
	Abort(err error) error
}

var _identNext = 1
var _identLock sync.Mutex

//go:generate go tool plc4xGenerator -type=IOCB -prefix=iocb_
type IOCB struct {
	ioID           int
	request        PDU       `stringer:"true"`
	destination    *Address  `stringer:"true"`
	ioState        IOCBState `stringer:"true"`
	ioResponse     PDU       `stringer:"true"`
	ioError        error
	ioController   IOControllerRequirements `ignore:"true"`
	ioComplete     sync.Cond
	ioCompleteDone bool
	ioCallback     []func() `ignore:"true"`
	ioQueue        []IOCBContract
	ioTimeout      *time.Timer
	ioTimoutCancel chan struct{}
	priority       int

	wg sync.WaitGroup

	log zerolog.Logger
}

var _ IOCBContract = (*IOCB)(nil)

func NewIOCB(localLog zerolog.Logger, request PDU, destination *Address) (*IOCB, error) {
	// lock the identity sequence number
	_identLock.Lock()

	// generate a unique identity for this block
	ioID := _identNext
	_identNext++

	// release the lock
	_identLock.Unlock()

	//  debugging postponed until ID acquired
	localLog.Debug().Int("ioID", ioID).Msg("NewIOCB(%d)")

	return &IOCB{
		// save the ID
		ioID: ioID,

		// save the request parameter
		request:     request,
		destination: destination,

		ioComplete: *sync.NewCond(&sync.Mutex{}),

		// start with an idle request
		ioState: IOCBState_IDLE,

		log: localLog,
	}, nil
}

func (i *IOCB) GetIOResponse() PDU {
	return i.ioResponse
}

func (i *IOCB) GetIOError() error {
	return i.ioError
}

// AddCallback Pass a function to be called when IO is complete.
func (i *IOCB) AddCallback(fn func()) {
	i.log.Debug().
		Int("ioID", i.ioID).
		Interface("callback", fn).
		Msg("AddCallback(ioID): callback")
	// store it
	i.ioCallback = append(i.ioCallback, fn)

	// already complete?
	if i.ioCompleteDone {
		i.Trigger()
	}
}

// Wait for the completion event to be set
func (i *IOCB) Wait() {
	i.log.Debug().
		Int("ioID", i.ioID).
		Msg("Wait")
	i.ioComplete.L.Lock()
	i.ioComplete.Wait()
	i.ioComplete.L.Unlock()
}

// Trigger Set the completion event and make the callback(s)
func (i *IOCB) Trigger() {
	i.ioComplete.L.Lock()
	i.log.Debug().
		Int("ioID", i.ioID).
		Msg("Trigger")

	// if it's queued, remove it from its queue
	myIndex := -1
	var meAsInterface IOCBContract = i
	for index, qe := range i.ioQueue {
		if qe == meAsInterface {
			myIndex = index
		}
	}
	if myIndex >= 0 {
		i.log.Debug().Msg("dequeue")
		i.ioQueue = append(i.ioQueue[:myIndex], i.ioQueue[myIndex+1:]...)
	}

	// if there's a timer, cancel it
	if i.ioTimeout != nil {
		i.log.Debug().Msg("cancel timeout")
		i.ioTimeout.Stop()
	}

	// set the completion event
	i.ioComplete.Broadcast()
	i.ioComplete.L.Unlock()
	i.log.Debug().Msg("complete event set")

	// make callback(s)
	for _, f := range i.ioCallback {
		f()
	}
}

// Complete Called to complete a transaction, usually when ProcessIO has shipped the IOCB off to some other thread or
//
//	function.
func (i *IOCB) Complete(apdu PDU) error {
	i.log.Debug().
		Int("ioID", i.ioID).
		Stringer("apdu", apdu).
		Msg("Complete")

	if i.ioController != nil {
		// pass to the controller
		return i.ioController.CompleteIO(i, apdu)
	} else {
		// just fill in the data
		i.ioState = IOCBState_COMPLETED
		i.ioResponse = apdu
		i.Trigger()
		return nil
	}
}

// Abort Called by a client to abort a transaction.
func (i *IOCB) Abort(err error) error {
	i.log.Debug().Err(err).
		Int("ioID", i.ioID).
		Msg("Abort")
	defer close(i.ioTimoutCancel)

	if i.ioController != nil {
		// pass to the controller
		return i.ioController.AbortIO(i, err)
	} else {
		// just fill in the data
		i.ioState = IOCBState_ABORTED
		i.ioError = err
		i.Trigger()
		return nil
	}
}

// SetTimeout Called to set a transaction timer.
func (i *IOCB) SetTimeout(delay time.Duration) {
	// if one has already been created, cancel it
	if i.ioTimeout != nil {
		i.ioTimeout.Reset(delay)
	} else {
		now := GetTaskManagerTime()
		i.ioTimeout = time.NewTimer(delay)
		i.ioTimoutCancel = make(chan struct{})
		i.wg.Add(1)
		go func() {
			defer i.wg.Done()
			select {
			case timeout := <-i.ioTimeout.C:
				_ = i.Abort(utils.NewTimeoutError(now.Sub(timeout)))
			case <-i.ioTimoutCancel:
			}
		}()
	}
}

func (i *IOCB) setIOController(ioController IOControllerRequirements) {
	i.ioController = ioController
}

func (i *IOCB) setIOState(newState IOCBState) {
	i.ioState = newState
}

func (i *IOCB) getIOState() IOCBState {
	return i.ioState
}

func (i *IOCB) setIOResponse(msg PDU) {
	i.ioResponse = msg
}

func (i *IOCB) setIOError(err error) {
	i.ioError = err
}

func (i *IOCB) getRequest() PDU {
	return i.request
}

func (i *IOCB) GetDestination() *Address {
	return i.destination
}

func (i *IOCB) getPriority() int {
	return i.priority
}

func (i *IOCB) clearQueue() {
	i.ioQueue = nil
}

func (i *IOCB) Close() error { // TODO: ensure this is getting called
	defer utils.StopWarn(i.log)()
	i.log.Debug().Msg("IOCB closing")
	defer func() {
		i.log.Debug().Msg("waiting for running tasks to finnish")
		i.wg.Wait()
		i.log.Debug().Msg("waiting done")
	}()
	return nil
}
