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
	"container/heap"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate go tool plc4xGenerator -type=IOQueue -prefix=iocb_
type IOQueue struct {
	name     string
	notEmpty sync.Cond
	Queue    PriorityQueue[int, IOCBContract]

	wg sync.WaitGroup

	log zerolog.Logger
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
func (i *IOQueue) Put(iocb IOCBContract) error {
	i.log.Debug().Stringer("iocb", iocb).Msg("Put")

	// requests should be pending before being queued
	if iocb.getIOState() != IOCBState_PENDING {
		return errors.New("invalid state transition")
	}

	// add the request to the end of the list of iocb's at same priority
	priority := iocb.getPriority()

	heap.Push(&i.Queue, PriorityItem[int, IOCBContract]{Value: iocb, Priority: priority})

	i.notEmpty.Broadcast()
	return nil
}

// Get a request from a queue, optionally block until a request is available.
func (i *IOQueue) Get(block bool, delay *time.Duration) (IOCBContract, error) {
	i.log.Debug().
		Bool("block", block).
		Interface("delay", delay).
		Msg("GA")

	// if the queue is empty, and we do not block return None
	if !block && len(i.Queue) == 0 {
		i.log.Debug().Msg("not blocking and empty")
		return nil, nil
	}

	// wait for something to be in the queue
	if len(i.Queue) == 0 {
		if delay != nil {
			gotSomething := make(chan struct{})
			i.wg.Add(1)
			go func() {
				defer i.wg.Done()
				i.notEmpty.Wait()
				close(gotSomething)
			}()
			timeout := time.NewTimer(*delay)
			select {
			case <-gotSomething:
			case <-timeout.C:
				return nil, nil
			}
		} else {
			i.notEmpty.Wait()
		}
	}

	if len(i.Queue) == 0 {
		return nil, nil
	}

	// extract the first element
	pi := heap.Pop(&i.Queue).(PriorityItem[int, IOCBContract])
	iocb := pi.Value
	iocb.clearQueue()

	// return the request
	return iocb, nil
}

// Remove a control block from the queue, called if the request
//
//	is canceled/aborted
func (i *IOQueue) Remove(iocb IOCBContract) error {
	for _, item := range i.Queue {
		if iocb == item.Value {
			heap.Remove(&i.Queue, item.Index)

			if len(i.Queue) == 0 {
				i.notEmpty.Broadcast()
			}
			return nil
		}
	}
	return nil
}

// Abort all the control blocks in the queue
func (i *IOQueue) Abort(err error) {
	for _, item := range i.Queue {
		item.Value.clearQueue()
		_ = item.Value.Abort(err)
	}

	//
	i.Queue = nil

	// the queue is now empty, clear the event
	i.notEmpty.Broadcast()
}

func (i *IOQueue) Close() error {
	defer utils.StopWarn(i.log)()
	i.log.Debug().Msg("IOQueue closing")
	defer func() {
		i.log.Debug().Msg("waiting for running tasks to finnish")
		i.wg.Wait()
		i.log.Debug().Msg("waiting done")
	}()
	return nil
}
