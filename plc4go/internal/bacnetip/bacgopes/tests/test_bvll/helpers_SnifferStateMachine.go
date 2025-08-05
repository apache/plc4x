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

package test_bvll

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type SnifferStateMachine struct {
	*ClientStateMachine

	address *Address
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	if _debug != nil {
		_debug("__init__ %r %r", address, vlan)
	}
	machine, err := NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(s), WithLeafType(s))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}
	s.ClientStateMachine = machine

	// save the name and address
	s.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	s.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexj")
	}

	// fake multiplexer has a VLAN node in it
	s.mux, err = NewFauxMultiplexer(localLog, s.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// might receive all packets and allow spoofing
	s.mux.node.SetPromiscuous(true)
	s.mux.node.SetSpoofing(true)

	// bind the stack together
	if err := Bind(localLog, s, s.annexj, s.mux); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return s, nil
}
