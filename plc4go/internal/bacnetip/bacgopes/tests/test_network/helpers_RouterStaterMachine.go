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

package test_network

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
)

type RouterStateMachine struct {
	*RouterNode
	StateMachineContract
}

func NewRouterStateMachine(localLog zerolog.Logger, options ...Option) (*RouterStateMachine, error) {
	r := &RouterStateMachine{}
	ApplyAppliers(options, r)
	optionsForParent := AddLeafTypeIfAbundant(options, r)
	var err error
	r.RouterNode, err = NewRouterNode(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating router node")
	}
	var initFunc func()
	r.StateMachineContract, initFunc = NewStateMachine(localLog, r, optionsForParent...)
	initFunc()
	if !LogTestNetwork {
		r.log = zerolog.Nop()
	}
	return r, nil
}

func (r *RouterStateMachine) Send(args Args, kwArgs KWArgs) error {
	panic("not available")
}

func (r *RouterStateMachine) String() string {
	return "RouterStateMachine"
}
