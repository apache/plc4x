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
	"github.com/apache/plc4x/plc4go/internal/bacnetip/globals"
)

type ClientStateMachineContract interface {
	fmt.Stringer
	bacnetip.Client
	StateMachineContract
	Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

// ClientStateMachine An instance of this class sits at the top of a stack.  tPDU's that the
//
//	state machine sends are sent down the stack and tPDU's coming up the
//	stack are fed as received tPDU's.
type ClientStateMachine struct {
	bacnetip.Client
	StateMachineContract

	contract ClientStateMachineContract

	name string

	log zerolog.Logger
}

var _ ClientStateMachineContract = (*ClientStateMachine)(nil)

func NewClientStateMachine(localLog zerolog.Logger, opts ...func(*ClientStateMachine)) (*ClientStateMachine, error) {
	c := &ClientStateMachine{
		log: localLog,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.contract == nil {
		c.contract = c
	}
	if c.name != "" {
		c.log = c.log.With().Str("name", c.name).Logger()
	}
	var err error
	c.Client, err = bacnetip.NewClient(localLog, c.contract)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	var init func()
	c.StateMachineContract, init = NewStateMachine(localLog, c.contract, WithStateMachineName(c.name))
	init()
	return c, nil
}

func WithClientStateMachineName(name string) func(*ClientStateMachine) {
	return func(c *ClientStateMachine) {
		c.name = name
	}
}

func WithClientStateMachineExtension(contract ClientStateMachineContract) func(*ClientStateMachine) {
	return func(c *ClientStateMachine) {
		c.contract = contract
	}
}

func (s *ClientStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return s.contract.Request(args, kwargs)
}

func (s *ClientStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	return s.contract.Receive(args, kwargs)
}

func (s *ClientStateMachine) String() string {
	if globals.ExtendedGeneralOutput {
		return fmt.Sprintf("ClientStateMachine{%v, %v, name=%s}", s.Client, s.StateMachineContract, s.name)
	} else {
		return s.StateMachineContract.String()
	}
}
