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

package bacgopes

import (
	"container/heap"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate plc4xGenerator -type=IOQueue -prefix=iocb_
type IOQueue struct {
	name     string
	notEmpty sync.Cond
	queue    PriorityQueue[int, _IOCB]

	wg sync.WaitGroup

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

	heap.Push(&i.queue, PriorityItem[int, _IOCB]{iocb, priority, 0})

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
			i.wg.Add(1)
			go func() {
				defer i.wg.Done()
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
	pi := heap.Pop(&i.queue).(PriorityItem[int, _IOCB])
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

func (i *IOQueue) Close() error {
	i.log.Debug().Msg("IOQueue closing")
	defer func() {
		i.log.Debug().Msg("waiting for running tasks to finnish")
		i.wg.Wait()
		i.log.Debug().Msg("waiting done")
	}()
	return nil
}
