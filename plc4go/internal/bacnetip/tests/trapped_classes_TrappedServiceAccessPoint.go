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

package tests

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

type TrappedServiceAccessPointRequirements interface {
	SapIndication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	SapConfirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

// TrappedServiceAccessPoint  Note that while this class inherits from ServiceAccessPoint, it doesn't
//
//		provide any stubbed behavior for SapIndication() or SapConfirmation(),
//		so if these functions are called it will still raise panic.
//
//		To provide these functions, write a ServiceAccessPoint derived class and
//		stuff it in the inheritance sequence:
//
//		    struct Snort{
//		     ServiceAccessPoint
//	      }
//		        func SapIndication(pdu):
//		            ...do something...
//		        func SapConfirmation(pdu):
//		            ...do something...
//
//		    struct TrappedSnort(TrappedServiceAccessPoint, Snort)
//
//		The Snort functions will be called after the PDU is trapped.
type TrappedServiceAccessPoint struct {
	bacnetip.ServiceAccessPointContract
	requirements TrappedServiceAccessPointRequirements

	sapRequestSent          bacnetip.PDU
	sapIndicationReceived   bacnetip.PDU
	sapResponseSent         bacnetip.PDU
	sapConfirmationReceived bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedServiceAccessPoint(localLog zerolog.Logger, requirements TrappedServiceAccessPointRequirements) (*TrappedServiceAccessPoint, error) {
	t := &TrappedServiceAccessPoint{
		requirements: requirements,
		log:          localLog,
	}
	var err error
	t.ServiceAccessPointContract, err = bacnetip.NewServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating SAP")
	}
	return t, nil
}

func (s *TrappedServiceAccessPoint) GetSapRequestSent() bacnetip.PDU {
	return s.sapRequestSent
}

func (s *TrappedServiceAccessPoint) GetSapIndicationReceived() bacnetip.PDU {
	return s.sapIndicationReceived
}

func (s *TrappedServiceAccessPoint) GetSapResponseSent() bacnetip.PDU {
	return s.sapResponseSent
}

func (s *TrappedServiceAccessPoint) GetSapConfirmationReceived() bacnetip.PDU {
	return s.sapConfirmationReceived
}

func (s *TrappedServiceAccessPoint) SapRequest(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapRequest")
	s.sapRequestSent = args.Get0PDU()
	return s.ServiceAccessPointContract.SapRequest(args, kwargs)
}

func (s *TrappedServiceAccessPoint) SapIndication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapIndication")
	s.sapIndicationReceived = args.Get0PDU()
	return s.requirements.SapIndication(args, kwargs)
}

func (s *TrappedServiceAccessPoint) SapResponse(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapResponse")
	s.sapResponseSent = args.Get0PDU()
	return s.ServiceAccessPointContract.SapResponse(args, kwargs)
}

func (s *TrappedServiceAccessPoint) SapConfirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapConfirmation")
	s.sapConfirmationReceived = args.Get0PDU()
	return s.requirements.SapConfirmation(args, kwargs)
}
