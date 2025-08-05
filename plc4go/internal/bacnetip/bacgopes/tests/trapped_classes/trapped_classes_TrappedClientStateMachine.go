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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type TrappedClientStateMachine struct {
	*TrappedServer
	*TrappedStateMachine

	log zerolog.Logger
}

func NewTrappedClientStateMachine(localLog zerolog.Logger, options ...Option) (*TrappedClientStateMachine, error) {
	t := &TrappedClientStateMachine{log: localLog}
	ApplyAppliers(options, t)
	optionsForParent := AddLeafTypeIfAbundant(options, t)
	if _debug != nil {
		_debug("__init__")
	}
	var err error
	t.TrappedServer, err = NewTrappedServer(localLog, Combine(optionsForParent, WithTrappedServerContract(t))...)
	if err != nil {
		return nil, errors.Wrap(err, "error building trapped server")
	}
	t.TrappedStateMachine = NewTrappedStateMachine(localLog, optionsForParent...)
	return t, nil
}

func (t *TrappedClientStateMachine) Send(args Args, kwArgs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Send")

	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("send %r", pdu)
	}

	return t.Response(args, kwArgs)
}

func (t *TrappedClientStateMachine) Confirmation(args Args, kwArgs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Confirmation")

	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}

	return t.Receive(args, kwArgs)
}
