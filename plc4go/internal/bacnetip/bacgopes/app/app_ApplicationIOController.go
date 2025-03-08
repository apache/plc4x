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

package app

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/iocb"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate go tool plc4xGenerator -type=ApplicationIOController -prefix=app_
type ApplicationIOController struct {
	*IOController
	*Application

	queueByAddress map[string]*SieveQueue

	log zerolog.Logger
}

func NewApplicationIOController(localLog zerolog.Logger, options ...Option) (*ApplicationIOController, error) {
	a := &ApplicationIOController{
		// queues for each address
		queueByAddress: make(map[string]*SieveQueue),
		log:            localLog,
	}
	ApplyAppliers(options, a)
	optionsForParent := AddLeafTypeIfAbundant(options, a)
	var err error
	a.IOController, err = NewIOController(localLog, "", a, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	a.Application, err = NewApplication(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application")
	}
	return a, nil
}

func (a *ApplicationIOController) ProcessIO(iocb IOCBContract) error {
	a.log.Debug().Stringer("iocb", iocb).Msg("ProcessIO")

	// get the destination address from the pdu
	destinationAddress := iocb.GetDestination()
	a.log.Debug().Stringer("destinationAddress", destinationAddress).Msg("working with destinationAddress")

	// look up the queue
	queue, ok := a.queueByAddress[destinationAddress.String()]
	if !ok {
		newQueue, _ := NewSieveQueue(a.log, a._AppRequest, destinationAddress)
		queue = newQueue
		a.queueByAddress[destinationAddress.String()] = queue
	}
	a.log.Debug().Stringer("queue", queue).Msg("working with queue")

	// ask the queue to process the request
	return queue.RequestIO(iocb)
}

func (a *ApplicationIOController) _AppComplete(address *Address, apdu PDU) error {
	a.log.Debug().
		Stringer("address", address).
		Stringer("apdu", apdu).
		Msg("_AppComplete")

	// look up the queue
	queue, ok := a.queueByAddress[address.String()]
	if !ok {
		a.log.Debug().Stringer("address", address).Msg("no queue for")
		return nil
	}
	a.log.Debug().Stringer("queue", queue).Msg("working with queue")

	// make sure it has an active iocb
	if queue.ActiveIOCB == nil {
		a.log.Debug().Stringer("address", address).Msg("no active request for")
		return nil
	}

	// this request is complete
	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUSimpleAck, readWriteModel.APDUComplexAck:
		if err := queue.CompleteIO(queue.ActiveIOCB, apdu); err != nil {
			return err
		}
	case readWriteModel.APDUError, readWriteModel.APDUReject, readWriteModel.APDUAbort:
		// TODO: extract error
		if err := queue.AbortIO(queue.ActiveIOCB, errors.Errorf("%s", apdu)); err != nil {
			return err
		}
	default:
		return errors.New("unrecognized APDU type")
	}
	a.log.Debug().Msg("controller finished")
	// if the queue is empty and idle, forget about the controller
	if len(queue.IoQueue.Queue) == 0 && queue.ActiveIOCB == nil {
		delete(a.queueByAddress, address.String())
	}
	return nil
}

func (a *ApplicationIOController) _AppRequest(apdu PDU) {
	a.log.Debug().Stringer("apdu", apdu).Msg("_AppRequest")

	// send it downstream, bypass the guard
	if err := a.Application.Request(NA(apdu), NoKWArgs()); err != nil {
		a.log.Error().Stack().Err(err).Msg("Uh oh")
		return
	}

	// if this was an unconfirmed request, it's complete, no message
	if _, ok := apdu.(readWriteModel.APDUUnconfirmedRequest); ok {
		if err := a._AppComplete(apdu.GetPDUDestination(), apdu); err != nil {
			a.log.Error().Err(err).Msg("AppRequest failed")
			return
		}
	}
}

func (a *ApplicationIOController) Request(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Request")
	apdu := GA[PDU](args, 0)

	// if this is not unconfirmed request, tell the application to use the IOCB interface
	if _, ok := apdu.(readWriteModel.APDUUnconfirmedRequest); !ok {
		return errors.New("use IOCB for confirmed requests")
	}

	// send it downstream
	return a.Application.Request(args, kwArgs)
}

func (a *ApplicationIOController) Confirmation(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")
	apdu := GA[PDU](args, 0)

	// this is an ack, error, reject or abort
	return a._AppComplete(apdu.GetPDUSource(), apdu)
}

func (a *ApplicationIOController) Close() error {
	for addr, queue := range a.queueByAddress {
		a.log.Debug().Str("addr", addr).Msg("Closing")
		if err := queue.Close(); err != nil {
			a.log.Warn().Str("addr", addr).Err(err).Stringer("queue", queue).Msg("error closing")
		}
	}
	return nil
}
