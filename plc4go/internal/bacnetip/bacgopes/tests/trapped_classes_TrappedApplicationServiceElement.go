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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type TrappedApplicationServiceElementRequirements interface {
	Indication(args Args, kwargs KWArgs) error
	Confirmation(args Args, kwargs KWArgs) error
}

// TrappedApplicationServiceElement  Note that while this class inherits from ApplicationServiceElement, it
//
//	doesn't provide any stubbed behavior for indication() or confirmation(),
//	so if these functions are called it will still raise NotImplementedError.
//
//	To provide these functions, write a ServiceAccessPoint derived class and
//	stuff it in the inheritance sequence:
//
//	    class Snort(ApplicationServiceElement):
//	        def indication(self, pdu):
//	            ...do something...
//	        def confirmation(self, pdu):
//	            ...do something...
//
//	    class TrappedSnort(TrappedApplicationServiceElement, Snort): pass
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

func NewTrappedApplicationServiceElement(localLog zerolog.Logger, requirements TrappedApplicationServiceElementRequirements) (*TrappedApplicationServiceElement, error) {
	t := &TrappedApplicationServiceElement{
		requirements: requirements,
		log:          localLog,
	}
	var err error
	t.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating SAP")
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

func (s *TrappedApplicationServiceElement) Request(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Request")
	s.requestSent = GA[PDU](args, 0)
	return s.ApplicationServiceElementContract.Request(args, kwargs)
}

func (s *TrappedApplicationServiceElement) Indication(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	s.indicationReceived = GA[PDU](args, 0)
	return s.requirements.Indication(args, kwargs)
}

func (s *TrappedApplicationServiceElement) Response(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Response")
	s.responseSent = GA[PDU](args, 0)
	return s.ApplicationServiceElementContract.Response(args, kwargs)
}

func (s *TrappedApplicationServiceElement) Confirmation(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	s.confirmationReceived = GA[PDU](args, 0)
	return s.requirements.Confirmation(args, kwargs)
}
