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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type SnifferStateMachine struct {
	*ClientStateMachine

	address *Address
	node    *Node

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *Network) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	var err error
	s.ClientStateMachine, err = NewClientStateMachine(s.log, WithClientStateMachineName(address), WithClientStateMachineExtension(s), WithLeafType(s))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	s.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// create a promiscuous node, added to the network
	s.node, err = NewNode(s.log, s.address, WithNodePromiscuous(true), WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	if LogTestNetwork {
		s.log.Debug().Stringer("node", s.node).Msg("node")
	}

	// bind the stack together
	if err := Bind(s.log, s, s.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	if !LogTestNetwork {
		s.log = zerolog.Nop()
	}
	return s, nil
}
