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

package bacnetip

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var stateLog = log.Logger

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

type _IOCB interface {
	fmt.Stringer
	setIOController(ioController _IOController)
	setIOState(newState IOCBState)
	getIOState() IOCBState
	setIOResponse(msg PDU)
	Trigger()
	setIOError(err error)
	getRequest() PDU
	getDestination() *Address
	getPriority() int
	clearQueue()
	Abort(err error) error
}

var _identNext = 1
var _identLock sync.Mutex

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=IOCB
type IOCB struct {
	ioID           int
	request        PDU
	destination    *Address
	ioState        IOCBState
	ioResponse     PDU
	ioError        error
	ioController   _IOController
	ioComplete     sync.Cond
	ioCompleteDone bool
	ioCallback     []func() `ignore:"true"`
	ioQueue        []_IOCB
	ioTimeout      *time.Timer `ignore:"true"`
	ioTimoutCancel chan any
	priority       int

	log zerolog.Logger `ignore:"true"`
}

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
	var meAsInterface _IOCB = i
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
		now := time.Now()
		i.ioTimeout = time.NewTimer(delay)
		i.ioTimoutCancel = make(chan any)
		go func() {
			select {
			case timeout := <-i.ioTimeout.C:
				_ = i.Abort(utils.NewTimeoutError(now.Sub(timeout)))
			case <-i.ioTimoutCancel:
			}
		}()
	}
}

func (i *IOCB) setIOController(ioController _IOController) {
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

func (i *IOCB) getDestination() *Address {
	return i.destination
}

func (i *IOCB) getPriority() int {
	return i.priority
}

func (i *IOCB) clearQueue() {
	i.ioQueue = nil
}

// An PriorityItem is something we manage in a priority queue.
type PriorityItem struct {
	value    _IOCB // The value of the item; arbitrary.
	priority int   // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*PriorityItem

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return (*pq)[i].priority > (*pq)[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *PriorityItem, value _IOCB, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=IOQueue
type IOQueue struct {
	name     string
	notEmpty sync.Cond
	queue    PriorityQueue

	log zerolog.Logger `ignore:"true"`
}

func NewIOQueue(localLog zerolog.Logger, name string) *IOQueue {
	localLog.Debug().Str("name", name).Msg("NewIOQueue")
	return &IOQueue{
		name: name,
		log:  localLog,
	}
}

// Put an IOCB to a queue.  This is usually called by the function that filters requests and passes them out to the
//
//	correct processing thread.
func (i *IOQueue) Put(iocb _IOCB) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("Put")

	// requests should be pending before being queued
	if iocb.getIOState() != IOCBState_PENDING {
		return errors.New("invalid state transition")
	}

	// add the request to the end of the list of iocb's at same priority
	priority := iocb.getPriority()

	heap.Push(&i.queue, PriorityItem{iocb, priority, 0})

	i.notEmpty.Broadcast()
	return nil
}

// Get a request from a queue, optionally block until a request is available.
func (i *IOQueue) Get(block bool, delay *time.Duration) (_IOCB, error) {
	i.log.Debug().
		Bool("block", block).
		Interface("delay", delay).
		Msg("Get")

	// if the queue is empty, and we do not block return None
	if !block && len(i.queue) == 0 {
		i.log.Debug().Msg("not blocking and empty")
		return nil, nil
	}

	// wait for something to be in the queue
	if len(i.queue) == 0 {
		if delay != nil {
			gotSomething := make(chan any)
			go func() {
				i.notEmpty.Wait()
				close(gotSomething)
			}()
			timeout := time.NewTimer(*delay)
			defer utils.CleanupTimer(timeout)
			select {
			case <-gotSomething:
			case <-timeout.C:
				return nil, nil
			}
		} else {
			i.notEmpty.Wait()
		}
	}

	if len(i.queue) == 0 {
		return nil, nil
	}

	// extract the first element
	pi := heap.Pop(&i.queue).(PriorityItem)
	iocb := pi.value
	iocb.clearQueue()

	// return the request
	return iocb, nil
}

// Remove a control block from the queue, called if the request
//
//	is canceled/aborted
func (i *IOQueue) Remove(iocb _IOCB) error {
	for _, item := range i.queue {
		if iocb == item.value {
			heap.Remove(&i.queue, item.index)

			if len(i.queue) == 0 {
				i.notEmpty.Broadcast()
			}
			return nil
		}
	}
	return nil
}

// Abort all the control blocks in the queue
func (i *IOQueue) Abort(err error) {
	for _, item := range i.queue {
		item.value.clearQueue()
		_ = item.value.Abort(err)
	}

	//
	i.queue = nil

	// the queue is now empty, clear the event
	i.notEmpty.Broadcast()
}

type _IOController interface {
	Abort(err error) error
	ProcessIO(iocb _IOCB) error
	CompleteIO(iocb _IOCB, pdu PDU) error
	AbortIO(iocb _IOCB, err error) error
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=IOController
type IOController struct {
	name       string
	rootStruct _IOController

	log zerolog.Logger `ignore:"true"`
}

func NewIOController(localLog zerolog.Logger, name string, rootStruct _IOController) (*IOController, error) {
	localLog.Debug().Str("name", name).Msg("NewIOController")
	return &IOController{
		// save the name
		name:       name,
		rootStruct: rootStruct,
		log:        localLog,
	}, nil
}

// Abort all requests, no default implementation.
func (i *IOController) Abort(err error) error {
	return nil
}

// RequestIO Called by a client to start processing a request.
func (i *IOController) RequestIO(iocb _IOCB) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("RequestIO")

	// bind the iocb to this controller
	iocb.setIOController(i)

	// hopefully there won't be an error
	var err error

	// change the state
	iocb.setIOState(IOCBState_PENDING)

	// let derived class figure out how to process this
	err = i.rootStruct.ProcessIO(iocb)

	// if there was an error, abort the request
	if err != nil {
		return i.rootStruct.AbortIO(iocb, err)
	}
	return nil
}

// ProcessIO Figure out how to respond to this request.  This must be provided by the derived class.
func (i *IOController) ProcessIO(_IOCB) error {
	panic("IOController must implement ProcessIO()")
}

// ActiveIO Called by a handler to notify the controller that a request is being processed
func (i *IOController) ActiveIO(iocb _IOCB) error {
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
func (i *IOController) CompleteIO(iocb _IOCB, apdu PDU) error {
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
func (i *IOController) AbortIO(iocb _IOCB, err error) error {
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

type _IOQController interface {
	ProcessIO(iocb _IOCB) error
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=IOQController
type IOQController struct {
	*IOController
	state      IOQControllerStates
	activeIOCB _IOCB
	ioQueue    *IOQueue
	waitTime   time.Duration `stringer:"true"`
	rootStruct _IOQController

	log zerolog.Logger `ignore:"true"`
}

func NewIOQController(localLog zerolog.Logger, name string, rootStruct _IOQController) (*IOQController, error) {
	i := &IOQController{
		rootStruct: rootStruct,
		log:        localLog,
	}
	var err error
	i.IOController, err = NewIOController(localLog, name, i)
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
	i.activeIOCB = nil

	// create an IOQueue for iocb's requested when not idle
	i.ioQueue = NewIOQueue(localLog, name+" queue")

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
		iocb, err := i.ioQueue.Get(false, nil)
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
func (i *IOQController) RequestIO(iocb _IOCB) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("RequestIO")

	// bind the iocb to this controller
	iocb.setIOController(i)

	// if we're busy, queue it
	if i.state != IOQControllerStates_CTRL_IDLE {
		i.log.Debug().
			Stringer("activeIOCB", i.activeIOCB).
			Msg("busy, request queued, activeIOCB")

		iocb.setIOState(IOCBState_PENDING)

		if err := i.ioQueue.Put(iocb); err != nil {
			return errors.Wrap(err, "error putting iocb in the queue")
		}
		return nil
	}

	if err := i.rootStruct.ProcessIO(iocb); err != nil {
		i.log.Debug().Err(err).Msg("ProcessIO error")
		if err := i.Abort(err); err != nil {
			return errors.Wrap(err, "error sending abort")
		}
	}

	return nil
}

// ProcessIO Figure out how to respond to this request.  This must be provided by the derived class
func (i *IOQController) ProcessIO(_IOCB) error {
	panic("IOController must implement ProcessIO()")
}

// ActiveIO Called by a handler to notify the controller that a request is being processed
func (i *IOQController) ActiveIO(iocb _IOCB) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("ActiveIO")

	// base class work first, setting iocb state and timer data
	if err := i.IOController.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "error calling super active io")
	}

	// change our state
	i.state = IOQControllerStates_CTRL_ACTIVE
	stateLog.Debug().Timestamp().Str("name", i.name).Msg("active")

	// keep track of the iocb
	i.activeIOCB = iocb
	return nil
}

// CompleteIO Called by a handler to return data to the client
func (i *IOQController) CompleteIO(iocb _IOCB, msg PDU) error {
	i.log.Debug().Stringer("iocb", iocb).Stringer("msg", msg).Msg("CompleteIO")

	// check to see if it is completing the active one
	if iocb != i.activeIOCB {
		return errors.New("not the current iocb")
	}

	// normal completion
	if err := i.IOController.CompleteIO(iocb, msg); err != nil {
		return errors.Wrap(err, "error completing io")
	}

	// no longer an active iocb
	i.activeIOCB = nil

	// check to see if we should wait a bit
	if i.waitTime != 0 {
		// change our state
		i.state = IOQControllerStates_CTRL_WAITING
		stateLog.Debug().Timestamp().Str("name", i.name).Msg("waiting")

		task := FunctionTask(i._waitTrigger, NoArgs, NoKWArgs)
		task.InstallTask(InstallTaskOptions{Delta: &i.waitTime})
	} else {
		// change our state
		i.state = IOQControllerStates_CTRL_IDLE
		stateLog.Debug().Timestamp().Str("name", i.name).Msg("idle")

		// look for more to do
		Deferred(i._trigger, NoArgs, NoKWArgs)
	}

	return nil
}

// AbortIO Called by a handler or a client to abort a transaction
func (i *IOQController) AbortIO(iocb _IOCB, err error) error {
	i.log.Debug().Err(err).Stringer("iocb", iocb).Msg("AbortIO")

	// Normal abort
	if err := i.IOController.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "ActiveIO failed")
	}

	// check to see if it is completing the active one
	if iocb != i.activeIOCB {
		i.log.Debug().Msg("not current iocb")
		return nil
	}

	// no longer an active iocb
	i.activeIOCB = nil

	// change our state
	i.state = IOQControllerStates_CTRL_IDLE
	stateLog.Debug().Timestamp().Str("name", i.name).Msg("idle")

	// look for more to do
	Deferred(i._trigger, NoArgs, NoKWArgs)
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
	if len(i.ioQueue.queue) == 0 {
		i.log.Debug().Msg("empty queue")
		return nil
	}

	// get the next iocb
	iocb, err := i.ioQueue.Get(false, nil)
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
		Deferred(i._trigger, NoArgs, NoKWArgs)
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
	stateLog.Debug().Timestamp().Str("name", i.name).Msg("idle")

	// look for more to do
	return i._trigger(NoArgs, NoKWArgs)
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=SieveQueue
type SieveQueue struct {
	*IOQController
	requestFn func(apdu PDU)
	address   *Address

	log zerolog.Logger `ignore:"true"`
}

func NewSieveQueue(localLog zerolog.Logger, fn func(apdu PDU), address *Address) (*SieveQueue, error) {
	s := &SieveQueue{}
	var err error
	s.IOQController, err = NewIOQController(localLog, address.String(), s)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a IOQController")
	}

	// Save a reference to the request function
	s.requestFn = fn
	s.address = address
	return s, nil
}

func (s *SieveQueue) ProcessIO(iocb _IOCB) error {
	s.log.Debug().Stringer("iocb", iocb).Msg("ProcessIO")

	// this is now an active request
	if err := s.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "error on active io")
	}

	// send the request
	s.requestFn(iocb.getRequest())
	return nil
}
