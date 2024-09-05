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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate plc4xGenerator -type=ApplicationServiceAccessPoint -prefix=appservice_
type ApplicationServiceAccessPoint struct {
	ApplicationServiceElementContract
	ServiceAccessPointContract

	// pass through args
	argAseID *int `ignore:"true"`
	argSapID *int `ignore:"true"`

	log zerolog.Logger
}

func NewApplicationServiceAccessPoint(localLog zerolog.Logger, opts ...func(*ApplicationServiceAccessPoint)) (*ApplicationServiceAccessPoint, error) {
	a := &ApplicationServiceAccessPoint{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}
	var err error
	a.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, func(ase *applicationServiceElement) {
		ase.elementID = a.argAseID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}
	a.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, func(sap *serviceAccessPoint) {
		sap.serviceID = a.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating service access point")
	}
	return a, nil
}

func WithApplicationServiceAccessPointAseID(aseID int) func(*ApplicationServiceAccessPoint) {
	return func(a *ApplicationServiceAccessPoint) {
		a.argAseID = &aseID
	}
}

func WithApplicationServiceAccessPointSapID(sapID int) func(*ApplicationServiceAccessPoint) {
	return func(a *ApplicationServiceAccessPoint) {
		a.argSapID = &sapID
	}
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Indication(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	apdu := args.Get0PDU()

	switch _apdu := apdu.GetRootMessage().(type) {
	case readWriteModel.APDUConfirmedRequest:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetConfirmedServiceChoiceKnows(uint8(_apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		if errorFound == nil {
			errorFound = a.SapRequest(NewArgs(apdu), NoKWArgs)
		}
		// TODO: the handling here gets a bit different now... need to wrap the head around how to do this (error handling etc)

		if errorFound == nil {
			if err := a.SapRequest(NewArgs(apdu), NoKWArgs); err != nil {
				return err
			}
		} else {
			a.log.Debug().Err(errorFound).Msg("got error")

			// TODO: map it to a error... code temporary placeholder
			return a.Response(NewArgs(NewPDU(readWriteModel.NewAPDUReject(_apdu.GetInvokeId(), nil, 0))), NoKWArgs)
		}
	case readWriteModel.APDUUnconfirmedRequest:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetUnconfirmedServiceChoiceKnows(uint8(_apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		if errorFound == nil {
			errorFound = a.SapRequest(NewArgs(apdu), NoKWArgs)
		}
		// TODO: the handling here gets a bit different now... need to wrap the head around how to do this (error handling etc)

		if errorFound == nil {
			if err := a.SapRequest(NewArgs(apdu), NoKWArgs); err != nil {
				return err
			}
		} else {
			a.log.Debug().Err(errorFound).Msg("got error")
		}

	default:
		return errors.Errorf("unknown _PDU type %T", apdu)
	}
	return nil
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapIndication(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapIndication")

	apdu := args.Get0APDU()

	isConfirmed := false
	var xpdu APDU
	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUConfirmedRequest:
		var err error
		xpdu, err = NewConfirmedRequestPDU(nil)
		if err != nil {
			return errors.Wrap(err, "error creating unconfirmed request")
		}
		if err := apdu.Encode(xpdu); err != nil {
			return errors.Wrap(err, "error encoding APDU")
		}
		isConfirmed = true
	case readWriteModel.APDUUnconfirmedRequest:
		var err error
		xpdu, err = NewUnconfirmedRequestPDU(nil)
		if err != nil {
			return errors.Wrap(err, "error creating unconfirmed request")
		}
		if err := apdu.Encode(xpdu); err != nil {
			return errors.Wrap(err, "error encoding APDU")
		}
	default:
		return errors.Errorf("unknown _PDU type %T", apdu)
	}

	// forward the encoded packet
	err := a.Request(NewArgs(xpdu), NoKWArgs)
	if err != nil {
		return errors.Wrap(err, "error forwarding the request ")
	}

	// if the upper layers of the application did not assign an invoke ID,
	// copy the one that was assigned on its way down the stack
	if isConfirmed && apdu.getApduInvokeID() != nil {
		apdu.setApduInvokeID(xpdu.getApduInvokeID())
	}
	return err
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Confirmation(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	// TODO: check if we need to check apdu here

	// forward the decoded packet
	return a.SapResponse(args, kwargs)
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapConfirmation(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")

	// TODO: check if we need to check apdu here

	return a.Response(args, kwargs)
}
