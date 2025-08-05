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

// BIPStateMachine is an application layer for BVLL messages that has no BVLL
//
//	processing like the 'simple', 'foreign', or 'bbmd' versions.  The client
//	state machine sits above and Annex-J codec so the send and receive PDUs are
//	BVLL PDUs.
type BIPStateMachine struct {
	*ClientStateMachine

	address *Address
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPStateMachine, error) {
	b := &BIPStateMachine{}
	if _debug != nil {
		_debug("__init__ %r %r", address, vlan)
	}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(b), WithLeafType(b))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexj")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)

	// bind the stack together
	err = Bind(localLog, b, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}
