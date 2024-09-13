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

package appservice

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate plc4xGenerator -type=ApplicationServiceAccessPoint -prefix=appservice_
type ApplicationServiceAccessPoint struct {
	ApplicationServiceElementContract
	ServiceAccessPointContract

	// pass through args
	argAseID        *int                       `ignore:"true"`
	argASEExtension *ApplicationServiceElement `ignore:"true"`
	argSapID        *int                       `ignore:"true"`
	argSap          *ServiceAccessPoint        `ignore:"true"`

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
	a.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, OptionalOption2(a.argAseID, a.argASEExtension, WithApplicationServiceElementAseID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}
	a.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, OptionalOption2(a.argSapID, a.argSap, WithServiceAccessPointSapID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating service access point")
	}
	return a, nil
}

func WithApplicationServiceAccessPointAseID(aseID int, argASEExtension ApplicationServiceElement) func(*ApplicationServiceAccessPoint) {
	return func(a *ApplicationServiceAccessPoint) {
		a.argAseID = &aseID
		a.argASEExtension = &argASEExtension
	}
}

func WithApplicationServiceAccessPointSapID(sapID int, sap ServiceAccessPoint) func(*ApplicationServiceAccessPoint) {
	return func(a *ApplicationServiceAccessPoint) {
		a.argSapID = &sapID
		a.argSap = &sap
	}
}

func (a *ApplicationServiceAccessPoint) Indication(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	apdu := GA[APDU](args, 0)

	switch _apdu := apdu.GetRootMessage().(type) {
	case readWriteModel.APDUConfirmedRequest:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetConfirmedServiceChoiceKnows(uint8(_apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		var apduService readWriteModel.BACnetConfirmedServiceChoice
		if sr := _apdu.GetServiceRequest(); sr != nil {
			apduService = sr.GetServiceChoice()
		}
		// Look up the struct associated with the service
		cr, ok := ConfirmedRequestTypes[apduService]
		if !ok {
			a.log.Debug().Stringer("apduService", apduService).Msg("unknown service type")
			errorFound = errors.New("unrecognized service")
		}

		var xpdu Decoder
		// no error so far, keep going
		if errorFound == nil {
			xpdu = cr()
			if err := xpdu.Decode(apdu); err != nil {
				// TODO: add advanced error check for  reject and abort
				panic("do it")
				errorFound = err
			}
		}

		// no error so far, keep going
		if errorFound == nil {
			a.log.Trace().Msg("no decoding error")
			if err := a.SapRequest(NA(xpdu), NoKWArgs); err != nil {
				panic("if no abort or reject bubble up")
				errorFound = err
			}
		}

		switch {
		case false: // TODO: check for Reject or Abort error
			panic("implement it")
			a.log.Debug().Err(errorFound).Msg("got error")

			// TODO: map it to a error... code temporary placeholder
			return a.Response(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, readWriteModel.NewAPDUReject(_apdu.GetInvokeId(), nil, 0)))), NoKWArgs)
		}
	case readWriteModel.APDUUnconfirmedRequest:
		var apduService readWriteModel.BACnetUnconfirmedServiceChoice
		if sr := _apdu.GetServiceRequest(); sr != nil {
			apduService = sr.GetServiceChoice()
		}
		// Look up the struct associated with the service
		ur, ok := UnconfirmedRequestTypes[apduService]
		if !ok {
			a.log.Debug().Stringer("apduService", apduService).Msg("unknown service type")
			return nil
		}

		xpdu := ur()
		if err := xpdu.Decode(apdu); err != nil {
			// TODO: add advanced error check for  reject and abort
			panic("do it")
		}

		// forward the decoded packet
		if err := a.SapRequest(NA(xpdu), NoKWArgs); err != nil {
			panic("if no abort or reject bubble up")
		}
	default:
		return errors.Errorf("unknown _PDU type %T", apdu)
	}
	return nil
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapIndication(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapIndication")

	apdu := GA[APDU](args, 0)

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
	err := a.Request(NA(xpdu), NoKWArgs)
	if err != nil {
		return errors.Wrap(err, "error forwarding the request ")
	}

	// if the upper layers of the application did not assign an invoke ID,
	// copy the one that was assigned on its way down the stack
	if isConfirmed && apdu.GetApduInvokeID() != nil {
		apdu.SetApduInvokeID(xpdu.GetApduInvokeID())
	}
	return err
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Confirmation(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")

	// TODO: check if we need to check apdu here

	// forward the decoded packet
	return a.SapResponse(args, kwArgs)
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapConfirmation(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapConfirmation")

	// TODO: check if we need to check apdu here

	return a.Response(args, kwArgs)
}
