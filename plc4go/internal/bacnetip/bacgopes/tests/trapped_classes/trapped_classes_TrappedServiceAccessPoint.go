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

package trapped_classes

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type TrappedServiceAccessPointRequirements interface {
	SapIndication(args Args, kwArgs KWArgs) error
	SapConfirmation(args Args, kwArgs KWArgs) error
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
	ServiceAccessPointContract
	requirements TrappedServiceAccessPointRequirements

	sapRequestSent          PDU
	sapIndicationReceived   PDU
	sapResponseSent         PDU
	sapConfirmationReceived PDU

	log zerolog.Logger
}

func NewTrappedServiceAccessPoint(localLog zerolog.Logger, requirements TrappedServiceAccessPointRequirements, options ...Option) (*TrappedServiceAccessPoint, error) {
	t := &TrappedServiceAccessPoint{
		requirements: requirements,
		log:          localLog,
	}
	ApplyAppliers(options, t)
	optionsForParent := AddLeafTypeIfAbundant(options, t)
	var err error
	t.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating SAP")
	}
	if _debug != nil {
		_debug("__init__(%s)", t.GetServiceID())
	}
	return t, nil
}

func (s *TrappedServiceAccessPoint) GetSapRequestSent() PDU {
	return s.sapRequestSent
}

func (s *TrappedServiceAccessPoint) GetSapIndicationReceived() PDU {
	return s.sapIndicationReceived
}

func (s *TrappedServiceAccessPoint) GetSapResponseSent() PDU {
	return s.sapResponseSent
}

func (s *TrappedServiceAccessPoint) GetSapConfirmationReceived() PDU {
	return s.sapConfirmationReceived
}

func (s *TrappedServiceAccessPoint) SapRequest(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("SapRequest")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("sap_request(%s) %r", s.GetServiceID(), pdu)
	}
	s.sapRequestSent = pdu
	return s.ServiceAccessPointContract.SapRequest(args, kwArgs)
}

func (s *TrappedServiceAccessPoint) SapIndication(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("SapIndication")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("sap_indication(%s) %r", s.GetServiceID(), pdu)
	}
	s.sapIndicationReceived = pdu
	return s.requirements.SapIndication(args, kwArgs)
}

func (s *TrappedServiceAccessPoint) SapResponse(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("SapResponse")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("sap_response(%s) %r", s.GetServiceID(), pdu)
	}
	s.sapResponseSent = pdu
	return s.ServiceAccessPointContract.SapResponse(args, kwArgs)
}

func (s *TrappedServiceAccessPoint) SapConfirmation(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("SapConfirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("sap_confirmation(%s) %r", s.GetServiceID(), pdu)
	}
	s.sapConfirmationReceived = pdu
	return s.requirements.SapConfirmation(args, kwArgs)
}
