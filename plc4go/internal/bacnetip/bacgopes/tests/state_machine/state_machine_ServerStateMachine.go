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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

//go:generate go tool plc4xGenerator -type=ServerStateMachine -prefix=state_machine_
type ServerStateMachine struct {
	ServerContract
	StateMachineContract

	name string

	log zerolog.Logger
}

func NewServerStateMachine(localLog zerolog.Logger, options ...Option) (*ServerStateMachine, error) {
	c := &ServerStateMachine{
		log: localLog,
	}
	ApplyAppliers(options, c)
	optionsForParent := AddLeafTypeIfAbundant(options, c)
	var err error
	c.ServerContract, err = NewServer(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Server")
	}
	var init func()
	c.StateMachineContract, init = NewStateMachine(localLog, c, Combine(options, WithStateMachineName(c.name))...)
	init()
	return c, nil
}

func WithServerStateMachineName(name string) GenericApplier[*ServerStateMachine] {
	return WrapGenericApplier(func(s *ServerStateMachine) { s.name = name })
}

func (s *ServerStateMachine) Send(args Args, kwArgs KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Send")
	return s.Response(args, kwArgs)
}

func (s *ServerStateMachine) Indication(args Args, kwArgs KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Indication")
	return s.Receive(args, kwArgs)
}
