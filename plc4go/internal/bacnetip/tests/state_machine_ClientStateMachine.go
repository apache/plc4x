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

// ClientStateMachine An instance of this class sits at the top of a stack.  tPDU's that the
//
//	state machine sends are sent down the stack and tPDU's coming up the
//	stack are fed as received tPDU's.
type ClientStateMachine struct {
	*bacnetip.Client
	StateMachine

	name string

	log zerolog.Logger
}

func NewClientStateMachine(localLog zerolog.Logger, opts ...func(*ClientStateMachine)) (*ClientStateMachine, error) {
	c := &ClientStateMachine{
		log: localLog,
	}
	for _, opt := range opts {
		opt(c)
	}
	var err error
	c.Client, err = bacnetip.NewClient(localLog, c)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	var init func()
	c.StateMachine, init = NewStateMachine(localLog, c, WithStateMachineName(c.name))
	init()
	return c, nil
}

func WithClientStateMachineName(name string) func(*ClientStateMachine) {
	return func(c *ClientStateMachine) {
		c.name = name
	}
}

func (s *ClientStateMachine) String() string {
	return fmt.Sprintf("ClientStateMachine{%v, %v}", s.Client, s.StateMachine)
}

func (s *ClientStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return s.Request(args, kwargs)
}

func (s *ClientStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	return s.Receive(args, kwargs)
}
