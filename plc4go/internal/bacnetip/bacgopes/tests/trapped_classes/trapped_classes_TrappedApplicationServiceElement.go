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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type TrappedApplicationServiceElementRequirements interface {
	Indication(args Args, kwArgs KWArgs) error
	Confirmation(args Args, kwArgs KWArgs) error
}

// TrappedApplicationServiceElement  Note that while this class inherits from ApplicationServiceElement, it
//
//	doesn't provide any stubbed behavior for indication() or confirmation(),
//	so if these functions are called it will still raise NotImplementedError.
//
//	To provide these functions, write a ServiceAccessPoint derived class and
//	stuff it in the inheritance sequence:
//
//	    struct Snort(ApplicationServiceElement)
//	        func indication(pdu):
//	            ...do something...
//	        func confirmation(pdu)
//	            ...do something...
//
//	    struct TrappedSnort(TrappedApplicationServiceElement, Snort)
//
//	The Snort functions will be called after the PDU is trapped.
type TrappedApplicationServiceElement struct {
	ApplicationServiceElementContract
	requirements TrappedApplicationServiceElementRequirements

	requestSent          PDU
	indicationReceived   PDU
	responseSent         PDU
	confirmationReceived PDU

	log zerolog.Logger
}

var _ ApplicationServiceElement = (*TrappedApplicationServiceElement)(nil)

func NewTrappedApplicationServiceElement(localLog zerolog.Logger, requirements TrappedApplicationServiceElementRequirements, options ...Option) (*TrappedApplicationServiceElement, error) {
	t := &TrappedApplicationServiceElement{
		requirements: requirements,
		log:          localLog,
	}
	ApplyAppliers(options, t)
	optionsForParent := AddLeafTypeIfAbundant(options, t)
	var err error
	t.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating SAP")
	}
	if _debug != nil {
		_debug("__init__(%s)", t.GetElementId())
	}
	return t, nil
}

func (s *TrappedApplicationServiceElement) GetRequestSent() PDU {
	return s.requestSent
}

func (s *TrappedApplicationServiceElement) GetIndicationReceived() PDU {
	return s.indicationReceived
}

func (s *TrappedApplicationServiceElement) GetResponseSent() PDU {
	return s.responseSent
}

func (s *TrappedApplicationServiceElement) GetConfirmationReceived() PDU {
	return s.confirmationReceived
}

func (s *TrappedApplicationServiceElement) String() string {
	return fmt.Sprintf("TrappedApplicationServiceElement(TBD...)") // TODO: fill some info here
}

func (s *TrappedApplicationServiceElement) Request(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Request")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("request(%s) %r", s.GetElementId(), pdu)
	}
	s.requestSent = pdu
	return s.ApplicationServiceElementContract.Request(args, kwArgs)
}

func (s *TrappedApplicationServiceElement) Indication(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Indication")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("indication(%s) %r", s.GetElementId(), pdu)
	}
	s.indicationReceived = pdu
	return s.requirements.Indication(args, kwArgs)
}

func (s *TrappedApplicationServiceElement) Response(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Response")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("response(%s) %r", s.GetElementId(), pdu)
	}
	s.responseSent = pdu
	return s.ApplicationServiceElementContract.Response(args, kwArgs)
}

func (s *TrappedApplicationServiceElement) Confirmation(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation(%s) %r", s.GetElementId(), pdu)
	}
	s.confirmationReceived = pdu
	return s.requirements.Confirmation(args, kwArgs)
}
