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

type ServerStateMachine struct {
	*bacnetip.Server
	StateMachine

	name string

	log zerolog.Logger
}

func NewServerStateMachine(localLog zerolog.Logger, opts ...func(*ServerStateMachine)) (*ServerStateMachine, error) {
	c := &ServerStateMachine{
		log: localLog,
	}
	for _, opt := range opts {
		opt(c)
	}
	var err error
	c.Server, err = bacnetip.NewServer(localLog, c)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Server")
	}
	var init func()
	c.StateMachine, init = NewStateMachine(localLog, c, WithStateMachineName(c.name))
	init()
	return c, nil
}

func WithServerStateMachineName(name string) func(*ServerStateMachine) {
	return func(s *ServerStateMachine) {
		s.name = name
	}
}

func (s *ServerStateMachine) String() string {
	return fmt.Sprintf("ServerStateMachine(TBD...)") // TODO: fill some info here
}

func (s *ServerStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return s.Response(args, kwargs)
}

func (s *ServerStateMachine) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Trace().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	return s.Receive(args, kwargs)
}
