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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

type TrapperRequirements interface {
	BeforeSend(pdu bacnetip.PDU)
	AfterSend(pdu bacnetip.PDU)
	BeforeReceive(pdu bacnetip.PDU)
	AfterReceive(pdu bacnetip.PDU)
	UnexpectedReceive(pdu bacnetip.PDU)
}

// Trapper This class provides a set of utility functions that keeps the latest copy of the pdu parameter in the
// before_send(), after_send(), before_receive(), after_receive() and unexpected_receive() calls.
type Trapper struct {
	TrapperRequirements

	beforeSendPdu        bacnetip.PDU
	afterSendPdu         bacnetip.PDU
	beforeReceivePdu     bacnetip.PDU
	afterReceivePdu      bacnetip.PDU
	unexpectedReceivePdu bacnetip.PDU

	log zerolog.Logger
}

func NewTrapper(localLog zerolog.Logger, requirements TrapperRequirements) *Trapper {
	trapper := &Trapper{
		TrapperRequirements: requirements,
		log:                 localLog,
	}
	// reset to initialize
	trapper.reset()
	return trapper
}

func (t *Trapper) reset() {
	t.log.Trace().Msg("Reset")
	// flush the copies
	t.beforeSendPdu = nil
	t.afterSendPdu = nil
	t.beforeReceivePdu = nil
	t.afterReceivePdu = nil
	t.unexpectedReceivePdu = nil
}

// BeforeSend is Called before each PDU about to be sent.
func (t *Trapper) BeforeSend(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("BeforeSend")
	//keep a copy
	t.beforeSendPdu = pdu

	// continue
	t.TrapperRequirements.BeforeSend(pdu)
}

func (t *Trapper) GetBeforeSendPdu() bacnetip.PDU {
	return t.beforeSendPdu
}

// AfterSend is Called after each PDU sent.
func (t *Trapper) AfterSend(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("AfterSend")
	//keep a copy
	t.afterSendPdu = pdu

	// continue
	t.TrapperRequirements.AfterSend(pdu)
}

func (t *Trapper) GetAfterSendPdu() bacnetip.PDU {
	return t.afterSendPdu
}

// BeforeReceive is Called with each PDU received before matching.
func (t *Trapper) BeforeReceive(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("BeforeReceive")
	//keep a copy
	t.beforeReceivePdu = pdu

	// continue
	t.TrapperRequirements.BeforeReceive(pdu)
}

func (t *Trapper) GetBeforeReceivePdu() bacnetip.PDU {
	return t.beforeReceivePdu
}

// AfterReceive is Called with PDU received after match.
func (t *Trapper) AfterReceive(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("AfterReceive")
	//keep a copy
	t.afterReceivePdu = pdu

	// continue
	t.TrapperRequirements.AfterReceive(pdu)
}

func (t *Trapper) GetAfterReceivePdu() bacnetip.PDU {
	return t.afterReceivePdu
}

// UnexpectedReceive is Called with PDU that did not match.  Unless this is trapped by the state, the default behaviour is to fail.
func (t *Trapper) UnexpectedReceive(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("UnexpectedReceive")
	//keep a copy
	t.unexpectedReceivePdu = pdu

	// continue
	t.TrapperRequirements.UnexpectedReceive(pdu)
}

func (t *Trapper) GetUnexpectedReceivePDU() bacnetip.PDU {
	return t.unexpectedReceivePdu
}

// TrappedState This class is a simple wrapper around the state class that keeps the latest copy of the pdu parameter in
// the BeforeSend(), AfterSend(), BeforeReceive(), AfterReceive() and UnexpectedReceive() calls.
type TrappedState struct {
	*Trapper
	State
}

func NewTrappedState(state State, trapper *Trapper) *TrappedState {
	t := &TrappedState{
		State:   state,
		Trapper: trapper,
	}
	return t
}

func (t *TrappedState) Equals(other State) bool {
	if t.State.Equals(other) { //TODO: we always want to match the inner
		return true
	}
	if otherTs, ok := other.(*TrappedState); ok {
		return t.State.Equals(otherTs.State)
	}
	return false
}

func (t *TrappedState) String() string {
	return fmt.Sprintf("TrappedState(%v)", t.State)
}

func (t *TrappedState) BeforeSend(pdu bacnetip.PDU) {
	t.Trapper.BeforeSend(pdu)
}

func (t *TrappedState) AfterSend(pdu bacnetip.PDU) {
	t.Trapper.AfterSend(pdu)
}

func (t *TrappedState) BeforeReceive(pdu bacnetip.PDU) {
	t.Trapper.BeforeReceive(pdu)
}

func (t *TrappedState) AfterReceive(pdu bacnetip.PDU) {
	t.Trapper.AfterReceive(pdu)
}

func (t *TrappedState) UnexpectedReceive(pdu bacnetip.PDU) {
	t.Trapper.UnexpectedReceive(pdu)
}

func (t *TrappedState) getInterceptor() StateInterceptor {
	return t
}

// TrappedStateMachine This class is a simple wrapper around the stateMachine class that keeps the
//
//	latest copy of the pdu parameter in the BeforeSend(), AfterSend(), BeforeReceive(), AfterReceive() and UnexpectedReceive() calls.
//
//	It also provides a send() function, so when the machine runs it doesn't
//	throw an exception.
type TrappedStateMachine struct {
	*Trapper
	StateMachine

	sent bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedStateMachine(localLog zerolog.Logger) *TrappedStateMachine {
	t := &TrappedStateMachine{
		log: localLog,
	}
	var init func()
	t.StateMachine, init = NewStateMachine(localLog, t, WithStateMachineStateInterceptor(t), WithStateMachineStateDecorator(t.DecorateState))
	t.Trapper = NewTrapper(localLog, t.StateMachine)
	init() // bit later so everything is set up
	return t
}

func (t *TrappedStateMachine) GetSent() bacnetip.PDU {
	return t.sent
}

func (t *TrappedStateMachine) BeforeSend(pdu bacnetip.PDU) {
	t.StateMachine.BeforeSend(pdu)
}

func (t *TrappedStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	// keep a copy
	t.sent = args.Get0PDU()
	return nil
}

func (t *TrappedStateMachine) AfterSend(pdu bacnetip.PDU) {
	t.StateMachine.AfterSend(pdu)
}

func (t *TrappedStateMachine) BeforeReceive(pdu bacnetip.PDU) {
	t.StateMachine.BeforeReceive(pdu)
}

func (t *TrappedStateMachine) AfterReceive(pdu bacnetip.PDU) {
	t.StateMachine.AfterReceive(pdu)
}

func (t *TrappedStateMachine) UnexpectedReceive(pdu bacnetip.PDU) {
	t.StateMachine.UnexpectedReceive(pdu)
}

func (t *TrappedStateMachine) DecorateState(state State) State {
	return NewTrappedState(state, t.Trapper)
}

// TrappedClientContract provides a set of functions which can be overwritten by a sub struct
type TrappedClientContract interface {
	Request(bacnetip.Args, bacnetip.KWArgs) error
	Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

// TrappedClient  An instance of this class sits at the top of a stack.
type TrappedClient struct {
	TrappedClientContract
	*bacnetip.Client

	requestSent          bacnetip.PDU
	confirmationReceived bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedClient(localLog zerolog.Logger, opts ...func(*TrappedClient)) (*TrappedClient, error) {
	t := &TrappedClient{
		log: localLog,
	}
	t.TrappedClientContract = t
	for _, opt := range opts {
		opt(t)
	}
	var err error
	t.Client, err = bacnetip.NewClient(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "error building client")
	}
	return t, nil
}

func WithTrappedClientContract(trappedClientContract TrappedClientContract) func(*TrappedClient) {
	return func(t *TrappedClient) {
		t.TrappedClientContract = trappedClientContract
	}
}

func (t *TrappedClient) GetRequestSent() bacnetip.PDU {
	return t.requestSent
}

func (t *TrappedClient) GetConfirmationReceived() bacnetip.PDU {
	return t.confirmationReceived
}

func (t *TrappedClient) Request(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Request")
	// a reference for checking
	t.requestSent = args.Get0PDU()

	// continue with regular processing
	return t.Client.Request(args, kwargs)
}

func (t *TrappedClient) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	// a reference for checking
	t.confirmationReceived = args.Get0PDU()
	return nil
}

func (t *TrappedClient) String() string {
	return fmt.Sprintf("TrappedClient{%s, requestSent: %v, confirmationReceived: %v}", t.Client, t.requestSent, t.confirmationReceived)
}

// TrappedServerContract provides a set of functions which can be overwritten by a sub struct
type TrappedServerContract interface {
	Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Response(bacnetip.Args, bacnetip.KWArgs) error
}

// TrappedServer An instance of this class sits at the bottom of a stack.
type TrappedServer struct {
	TrappedServerContract
	*bacnetip.Server

	indicationReceived bacnetip.PDU
	responseSent       bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedServer(localLog zerolog.Logger, opts ...func(*TrappedServer)) (*TrappedServer, error) {
	t := &TrappedServer{
		log: localLog,
	}
	t.TrappedServerContract = t
	for _, opt := range opts {
		opt(t)
	}
	var err error
	t.Server, err = bacnetip.NewServer(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "error building server")
	}
	return t, nil
}

func WithTrappedServerContract(trappedServerContract TrappedServerContract) func(*TrappedServer) {
	return func(t *TrappedServer) {
		t.TrappedServerContract = trappedServerContract
	}
}

func (t *TrappedServer) GetIndicationReceived() bacnetip.PDU {
	return t.indicationReceived
}

func (t *TrappedServer) GetResponseSent() bacnetip.PDU {
	return t.responseSent
}

func (t *TrappedServer) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	// a reference for checking
	t.indicationReceived = args.Get0PDU()

	return nil
}

func (t *TrappedServer) Response(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Response")
	// a reference for checking
	t.responseSent = args.Get0PDU()

	// continue with regular processing
	return t.Server.Response(args, kwargs)
}

func (t *TrappedServer) String() string {
	return fmt.Sprintf("TrappedServer{%s, indicationReceived: %v, responseSent: %v}", t.Server, t.indicationReceived, t.responseSent)
}

type TrappedServerStateMachine struct {
	*TrappedServer
	*TrappedStateMachine

	log zerolog.Logger
}

func NewTrappedServerStateMachine(localLog zerolog.Logger) (*TrappedServerStateMachine, error) {
	t := &TrappedServerStateMachine{log: localLog}
	var err error
	t.TrappedServer, err = NewTrappedServer(localLog, WithTrappedServerContract(t))
	if err != nil {
		return nil, errors.Wrap(err, "error building trapped server")
	}
	t.TrappedStateMachine = NewTrappedStateMachine(localLog)
	return t, nil
}

func (t *TrappedServerStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return t.Response(args, kwargs)
}

func (t *TrappedServerStateMachine) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	return t.Receive(args, kwargs)
}

type TrappedServiceAccessPointRequirements interface {
	SapRequest(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	SapIndication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	SapResponse(args bacnetip.Args, kwargs bacnetip.KWArgs) error
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
	TrappedServiceAccessPointRequirements
	*bacnetip.ServiceAccessPoint

	sapRequestSent          bacnetip.PDU
	sapIndicationReceived   bacnetip.PDU
	sapResponseSent         bacnetip.PDU
	sapConfirmationReceived bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedServiceAccessPoint(localLog zerolog.Logger, requirements TrappedServiceAccessPointRequirements) (*TrappedServiceAccessPoint, error) {
	t := &TrappedServiceAccessPoint{
		TrappedServiceAccessPointRequirements: requirements,
		log:                                   localLog,
	}
	var err error
	t.ServiceAccessPoint, err = bacnetip.NewServiceAccessPoint(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "error building service access point")
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
	return s.TrappedServiceAccessPointRequirements.SapRequest(args, kwargs)
}

func (s *TrappedServiceAccessPoint) SapIndication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapIndication")
	s.sapIndicationReceived = args.Get0PDU()
	return s.TrappedServiceAccessPointRequirements.SapIndication(args, kwargs)
}

func (s *TrappedServiceAccessPoint) SapResponse(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapResponse")
	s.sapResponseSent = args.Get0PDU()
	return s.TrappedServiceAccessPointRequirements.SapResponse(args, kwargs)
}

func (s *TrappedServiceAccessPoint) SapConfirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapConfirmation")
	s.sapConfirmationReceived = args.Get0PDU()
	return s.TrappedServiceAccessPointRequirements.SapConfirmation(args, kwargs)
}

type TrappedApplicationServiceElementRequirements interface {
	Request(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Response(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error
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
	TrappedApplicationServiceElementRequirements
	*bacnetip.ApplicationServiceElement

	requestSent          bacnetip.PDU
	indicationReceived   bacnetip.PDU
	responseSent         bacnetip.PDU
	confirmationReceived bacnetip.PDU

	log zerolog.Logger
}

func NewTrappedApplicationServiceElement(localLog zerolog.Logger, requirements TrappedApplicationServiceElementRequirements) (*TrappedApplicationServiceElement, error) {
	t := &TrappedApplicationServiceElement{
		TrappedApplicationServiceElementRequirements: requirements,
		log: localLog,
	}
	var err error
	t.ApplicationServiceElement, err = bacnetip.NewApplicationServiceElement(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service element")
	}
	return t, nil
}

func (s *TrappedApplicationServiceElement) GetRequestSent() bacnetip.PDU {
	return s.requestSent
}

func (s *TrappedApplicationServiceElement) GetIndicationReceived() bacnetip.PDU {
	return s.indicationReceived
}

func (s *TrappedApplicationServiceElement) GetResponseSent() bacnetip.PDU {
	return s.responseSent
}

func (s *TrappedApplicationServiceElement) GetConfirmationReceived() bacnetip.PDU {
	return s.confirmationReceived
}

func (s *TrappedApplicationServiceElement) String() string {
	return fmt.Sprintf("TrappedApplicationServiceElement(TBD...)") // TODO: fill some info here
}

func (s *TrappedApplicationServiceElement) Request(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Request")
	s.requestSent = args.Get0PDU()
	return s.TrappedApplicationServiceElementRequirements.Request(args, kwargs)
}

func (s *TrappedApplicationServiceElement) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	s.indicationReceived = args.Get0PDU()
	return s.TrappedApplicationServiceElementRequirements.Indication(args, kwargs)
}

func (s *TrappedApplicationServiceElement) Response(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Response")
	s.responseSent = args.Get0PDU()
	return s.TrappedApplicationServiceElementRequirements.Response(args, kwargs)
}

func (s *TrappedApplicationServiceElement) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	s.confirmationReceived = args.Get0PDU()
	return s.TrappedApplicationServiceElementRequirements.Confirmation(args, kwargs)
}
