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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ApplicationIOController struct {
	*IOController
	*Application

	queueByAddress map[string]*SieveQueue

	// pass through args
	argDeviceInfoCache *DeviceInfoCache
	argAseID           *int

	log zerolog.Logger
}

func NewApplicationIOController(localLog zerolog.Logger, localDevice *LocalDeviceObject, opts ...func(controller *ApplicationIOController)) (*ApplicationIOController, error) {
	a := &ApplicationIOController{
		// queues for each address
		queueByAddress: make(map[string]*SieveQueue),
		log:            localLog,
	}
	for _, opt := range opts {
		opt(a)
	}
	var err error
	a.IOController, err = NewIOController(localLog, "", a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	a.Application, err = NewApplication(localLog, localDevice, func(application *Application) {
		application.deviceInfoCache = a.argDeviceInfoCache
		application.argAseID = a.argAseID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating application")
	}
	return a, nil
}

func WithApplicationIOControllerDeviceInfoCache(deviceInfoCache *DeviceInfoCache) func(*ApplicationIOController) {
	return func(a *ApplicationIOController) {
		a.argDeviceInfoCache = deviceInfoCache
	}
}

func WithApplicationIOControllerAseID(aseID *int) func(*ApplicationIOController) {
	return func(a *ApplicationIOController) {
		a.argAseID = aseID
	}
}

func (a *ApplicationIOController) ProcessIO(iocb _IOCB) error {
	a.log.Debug().Stringer("iocb", iocb).Msg("ProcessIO")

	// get the destination address from the pdu
	destinationAddress := iocb.getDestination()
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
	if queue.activeIOCB == nil {
		a.log.Debug().Stringer("address", address).Msg("no active request for")
		return nil
	}

	// this request is complete
	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUSimpleAck, readWriteModel.APDUComplexAck:
		if err := queue.CompleteIO(queue.activeIOCB, apdu); err != nil {
			return err
		}
	case readWriteModel.APDUError, readWriteModel.APDUReject, readWriteModel.APDUAbort:
		// TODO: extract error
		if err := queue.AbortIO(queue.activeIOCB, errors.Errorf("%s", apdu)); err != nil {
			return err
		}
	default:
		return errors.New("unrecognized APDU type")
	}
	a.log.Debug().Msg("controller finished")
	// if the queue is empty and idle, forget about the controller
	if len(queue.ioQueue.queue) == 0 && queue.activeIOCB == nil {
		delete(a.queueByAddress, address.String())
	}
	return nil
}

func (a *ApplicationIOController) _AppRequest(apdu PDU) {
	a.log.Debug().Stringer("apdu", apdu).Msg("_AppRequest")

	// send it downstream, bypass the guard
	if err := a.Application.Request(NewArgs(apdu), NoKWArgs); err != nil {
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

func (a *ApplicationIOController) Request(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")
	apdu := args.Get0PDU()

	// if this is not unconfirmed request, tell the application to use the IOCB interface
	if _, ok := apdu.(readWriteModel.APDUUnconfirmedRequest); !ok {
		return errors.New("use IOCB for confirmed requests")
	}

	// send it downstream
	return a.Application.Request(args, kwargs)
}

func (a *ApplicationIOController) Confirmation(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	apdu := args.Get0PDU()

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
