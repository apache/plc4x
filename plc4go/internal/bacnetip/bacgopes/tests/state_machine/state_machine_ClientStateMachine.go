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

package state_machine

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type ClientStateMachineContract interface {
	fmt.Stringer
	Client
	StateMachineContract
	Send(args Args, kwArgs KWArgs) error
	Confirmation(args Args, kwArgs KWArgs) error
}

// ClientStateMachine An instance of this class sits at the top of a stack.  tPDU's that the
//
//	state machine sends are sent down the stack and tPDU's coming up the
//	stack are fed as received tPDU's.
//
//go:generate go tool plc4xGenerator -type=ClientStateMachine -prefix=state_machine_
type ClientStateMachine struct {
	ClientContract
	StateMachineContract

	contract ClientStateMachineContract

	name string

	log zerolog.Logger
}

var _ ClientStateMachineContract = (*ClientStateMachine)(nil)

func NewClientStateMachine(localLog zerolog.Logger, options ...Option) (*ClientStateMachine, error) {
	c := &ClientStateMachine{
		log: localLog,
	}
	ApplyAppliers(options, c)
	optionsForParent := AddLeafTypeIfAbundant(options, c)
	if c.contract == nil {
		c.contract = c
	}
	if _debug != nil {
		_debug("__init__")
	}
	if c.name != "" {
		c.log = c.log.With().Str("name", c.name).Logger()
	}
	var err error
	c.ClientContract, err = NewClient(localLog, optionsForParent...) // TODO: do we need to pass cid?
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	var init func()
	c.StateMachineContract, init = NewStateMachine(localLog, c.contract, Combine(optionsForParent, WithStateMachineName(c.name))...)
	init()
	return c, nil
}

func WithClientStateMachineName(name string) GenericApplier[*ClientStateMachine] {
	return WrapGenericApplier(func(c *ClientStateMachine) { c.name = name })
}

func WithClientStateMachineExtension(contract ClientStateMachineContract) GenericApplier[*ClientStateMachine] {
	return WrapGenericApplier(func(c *ClientStateMachine) { c.contract = contract })
}

func (s *ClientStateMachine) Send(args Args, kwArgs KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Send")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("send(%s) %r", s.name, pdu)
	}
	return s.contract.Request(args, kwArgs)
}

func (s *ClientStateMachine) Confirmation(args Args, kwArgs KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation(%s) %r", s.name, pdu)
	}
	return s.contract.Receive(args, kwArgs)
}

func (s *ClientStateMachine) AlternateString() (string, bool) {
	if IsDebuggingActive() {
		return s.StateMachineContract.String(), true
	}
	return "", false
}
