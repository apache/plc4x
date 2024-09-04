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

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
)

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

func (t *TrappedServerStateMachine) Send(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")
	return t.Response(args, kwargs)
}

func (t *TrappedServerStateMachine) Indication(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	return t.Receive(args, kwargs)
}
