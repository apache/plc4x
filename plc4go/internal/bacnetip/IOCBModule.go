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
	setIOController(ioController _IOController)
	setIOState(newState IOCBState)
	getIOState() IOCBState
	setIOResponse(msg _PDU)
	Trigger()
	setIOError(err error)
	getRequest() _PDU
	getDestination() *Address
	getPriority() int
	clearQueue()
	Abort(err error) error
}

var _identNext = 1
var _identLock sync.Mutex

type IOCB struct {
	ioID           int
	request        _PDU
	destination    *Address
	ioState        IOCBState
	ioResponse     _PDU
	ioError        error
	ioController   _IOController
	ioComplete     sync.Cond
	ioCompleteDone bool
	ioCallback     []func()
	ioQueue        []_IOCB
	ioTimeout      *time.Timer
	ioTimoutCancel chan any
	priority       int
}

func NewIOCB(request _PDU, destination *Address) (*IOCB, error) {
	// lock the identity sequence number
	_identLock.Lock()

	// generate a unique identity for this block
	ioID := _identNext
	_identNext++

	// release the lock
	_identLock.Unlock()

	//  debugging postponed until ID acquired
	log.Debug().Msgf("NewIOCB(%d)", ioID)

	return &IOCB{
		// save the ID
		ioID: ioID,

		// save the request parameter
		request:     request,
		destination: destination,

		ioComplete: *sync.NewCond(&sync.Mutex{}),

		// start with an idle request
		ioState: IOCBState_IDLE,
	}, nil
}

// AddCallback Pass a function to be called when IO is complete.
func (i *IOCB) AddCallback(fn func()) {
	log.Debug().Msgf("AddCallback(%d): %t", i.ioID, fn != nil)
	// store it
	i.ioCallback = append(i.ioCallback, fn)

	// already complete?
	if i.ioCompleteDone {
		i.Trigger()
	}
}

// Wait for the completion event to be set
func (i *IOCB) Wait() {
	log.Debug().Msgf("Wait(%d)", i.ioID)
	i.ioComplete.L.Lock()
	i.ioComplete.Wait()
	i.ioComplete.L.Unlock()
}

// Trigger Set the completion event and make the callback(s)
func (i *IOCB) Trigger() {
	i.ioComplete.L.Lock()
	log.Debug().Msgf("Trigger(%d)", i.ioID)

	// if it's queued, remove it from its queue
	myIndex := -1
	var meAsInterface _IOCB = i
	for index, qe := range i.ioQueue {
		if qe == meAsInterface {
			myIndex = index
		}
	}
	if myIndex >= 0 {
		log.Debug().Msg("dequeue")
		i.ioQueue = append(i.ioQueue[:myIndex], i.ioQueue[myIndex+1:]...)
	}

	// if there's a timer, cancel it
	if i.ioTimeout != nil {
		log.Debug().Msg("cancel timeout")
		i.ioTimeout.Stop()
	}

	// set the completion event
	i.ioComplete.Broadcast()
	i.ioComplete.L.Unlock()
	log.Debug().Msg("complete event set")

	// make callback(s)
	for _, f := range i.ioCallback {
		f()
	}
}

// Complete Called to complete a transaction, usually when ProcessIO has shipped the IOCB off to some other thread or
//
//	function.
func (i *IOCB) Complete(apdu _PDU) error {
	log.Debug().Msgf("Complete(%d)\n%s", i.ioID, apdu)

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
	log.Debug().Err(err).Msgf("Abort(%d)", i.ioID)
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

func (i *IOCB) setIOResponse(msg _PDU) {
	i.ioResponse = msg
}

func (i *IOCB) setIOError(err error) {
	i.ioError = err
}

func (i *IOCB) getRequest() _PDU {
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

type IOQueue struct {
	notEmpty sync.Cond
	queue    PriorityQueue
}

func NewIOQueue(name string) *IOQueue {
	log.Debug().Msgf("NewIOQueue %s", name)
	return &IOQueue{}
}

// Put an IOCB to a queue.  This is usually called by the function that filters requests and passes them out to the
//
//	correct processing thread.
func (i *IOQueue) Put(iocb _IOCB) error {
	log.Debug().Msgf("Put %s", iocb)

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
	log.Debug().Msgf("Get block=%t, delay=%s", block, delay)

	// if the queue is empty, and we do not block return None
	if !block && len(i.queue) == 0 {
		log.Debug().Msgf("not blocking and empty")
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
	CompleteIO(iocb _IOCB, pdu _PDU) error
	AbortIO(iocb _IOCB, err error) error
}

type IOController struct {
	name       string
	rootStruct _IOController
}

func NewIOController(name string, rootStruct _IOController) (*IOController, error) {
	log.Debug().Msgf("NewIOController name=%s", name)
	return &IOController{
		// save the name
		name:       name,
		rootStruct: rootStruct,
	}, nil
}

// Abort all requests, no default implementation.
func (i *IOController) Abort(err error) error {
	return nil
}

// RequestIO Called by a client to start processing a request.
func (i *IOController) RequestIO(iocb _IOCB) error {
	log.Debug().Msgf("RequestIO\n%s", iocb)

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
	log.Debug().Msgf("ActiveIO %s", iocb)

	// requests should be idle or pending before coming active
	if iocb.getIOState() != IOCBState_IDLE && iocb.getIOState() != IOCBState_PENDING {
		return errors.Errorf("invalid state transition (currently %d)", iocb.getIOState())
	}

	// change the state
	iocb.setIOState(IOCBState_ACTIVE)
	return nil
}

// CompleteIO Called by a handler to return data to the client
func (i *IOController) CompleteIO(iocb _IOCB, apdu _PDU) error {
	log.Debug().Msgf("CompleteIO %s\n%s", iocb, apdu)

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
	log.Debug().Err(err).Msgf("AbortIO %s", iocb)

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

type IOQController struct {
	*IOController
	state      IOQControllerStates
	activeIOCB _IOCB
	ioQueue    *IOQueue
	waitTime   time.Duration
	rootStruct _IOQController
}

func NewIOQController(name string, rootStruct _IOQController) (*IOQController, error) {
	i := &IOQController{
		rootStruct: rootStruct,
	}
	var err error
	i.IOController, err = NewIOController(name, i)
	if err != nil {
		return nil, errors.Wrap(err, "error creating IO controller")
	}

	// start idle
	i.state = IOQControllerStates_CTRL_IDLE
	log.Debug().Msgf("%s %s %s", time.Now(), name, i.state)

	// no active iocb
	i.activeIOCB = nil

	// create an IOQueue for iocb's requested when not idle
	i.ioQueue = NewIOQueue(name + " queue")

	return i, nil
}

// Abort all pending requests
func (i *IOQController) Abort(err error) error {
	log.Debug().Err(err).Msg("Abort")

	if i.state == IOQControllerStates_CTRL_IDLE {
		log.Debug().Msg("idle")
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
		log.Debug().Msgf("iocb %v", iocb)

		// change the state
		iocb.setIOState(IOCBState_ABORTED)
		iocb.setIOError(err)

		// notify the client
		iocb.Trigger()
	}

	if i.state != IOQControllerStates_CTRL_IDLE {
		log.Debug().Msg("busy after aborts")
	}

	return nil
}

// RequestIO Called by a client to start processing a request
func (i *IOQController) RequestIO(iocb _IOCB) error {
	log.Debug().Msgf("RequestIO %v", iocb)

	// bind the iocb to this controller
	iocb.setIOController(i)

	// if we're busy, queue it
	if i.state != IOQControllerStates_CTRL_IDLE {
		log.Debug().Msgf("busy, request queued, activeIOCB: %v", i.activeIOCB)

		iocb.setIOState(IOCBState_PENDING)

		if err := i.ioQueue.Put(iocb); err != nil {
			return errors.Wrap(err, "error putting iocb in the queue")
		}
		return nil
	}

	if err := i.rootStruct.ProcessIO(iocb); err != nil {
		log.Debug().Err(err).Msgf("ProcessIO error")
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
	log.Debug().Msgf("ActiveIO %v", iocb)

	// base class work first, setting iocb state and timer data
	if err := i.IOController.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "error calling super active io")
	}

	// change our state
	i.state = IOQControllerStates_CTRL_ACTIVE
	stateLog.Debug().Msgf("%s %s %s", time.Now(), i.name, "active")

	// keep track of the iocb
	i.activeIOCB = iocb
	return nil
}

// CompleteIO Called by a handler to return data to the client
func (i *IOQController) CompleteIO(iocb _IOCB, msg _PDU) error {
	log.Debug().Msgf("CompleteIO %v %v", iocb, msg)

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
		stateLog.Debug().Msgf("%s %s %s", time.Now(), i.name, "waiting")

		task := FunctionTask(i._waitTrigger)
		task.InstallTask(nil, &i.waitTime)
	} else {
		// change our state
		i.state = IOQControllerStates_CTRL_IDLE
		stateLog.Debug().Msgf("%s %s %s", time.Now(), i.name, "idle")

		// look for more to do
		Deferred(i._trigger)
	}

	return nil
}

// AbortIO Called by a handler or a client to abort a transaction
func (i *IOQController) AbortIO(iocb _IOCB, err error) error {
	log.Debug().Err(err).Msgf("AbortIO %v", iocb)

	// Normal abort
	if err := i.IOController.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "ActiveIO failed")
	}

	// check to see if it is completing the active one
	if iocb != i.activeIOCB {
		log.Debug().Msg("not current iocb")
		return nil
	}

	// no longer an active iocb
	i.activeIOCB = nil

	// change our state
	i.state = IOQControllerStates_CTRL_IDLE
	stateLog.Debug().Msgf("%s %s %s", time.Now(), i.name, "idle")

	// look for more to do
	Deferred(i._trigger)
	return nil
}

// _trigger Called to launch the next request in the queue
func (i *IOQController) _trigger() error {
	log.Debug().Msg("_trigger")

	// if we are busy, do nothing
	if i.state != IOQControllerStates_CTRL_IDLE {
		log.Debug().Msg("not idle")
		return nil
	}

	// if there is nothing to do, return
	if len(i.ioQueue.queue) == 0 {
		log.Debug().Msg("empty queue")
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
			log.Debug().Err(err).Msg("error aborting")
			return nil
		}
		return nil
	}

	// if we're idle, call again
	if i.state == IOQControllerStates_CTRL_IDLE {
		Deferred(i._trigger)
	}
	return nil
}

// _waitTrigger is called to launch the next request in the queue
func (i *IOQController) _waitTrigger() error {
	log.Debug().Msg("_waitTrigger")

	// make sure we are waiting
	if i.state != IOQControllerStates_CTRL_WAITING {
		log.Debug().Msg("not waiting")
		return nil
	}

	// change our state
	i.state = IOQControllerStates_CTRL_IDLE
	stateLog.Debug().Msgf("%s %s %s", time.Now(), i.name, "idle")

	// look for more to do
	return i._trigger()
}

type SieveQueue struct {
	*IOQController
	requestFn func(apdu _PDU)
	address   *Address
}

func NewSieveQueue(fn func(apdu _PDU), address *Address) (*SieveQueue, error) {
	s := &SieveQueue{}
	var err error
	s.IOQController, err = NewIOQController(address.String(), s)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a IOQController")
	}

	// Save a reference to the request function
	s.requestFn = fn
	s.address = address
	return s, nil
}

func (s *SieveQueue) ProcessIO(iocb _IOCB) error {
	log.Debug().Msgf("ProcessIO %s", iocb)

	// this is now an active request
	if err := s.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "error on active io")
	}

	// send the request
	s.requestFn(iocb.getRequest())
	return nil
}

func (s *SieveQueue) String() string {
	return fmt.Sprintf("%#q", s)
}
