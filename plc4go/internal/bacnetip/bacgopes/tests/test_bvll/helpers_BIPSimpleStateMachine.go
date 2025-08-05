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

type BIPSimpleStateMachine struct {
	*ClientStateMachine
	name string

	address *Address

	bip    *BIPSimple
	annexj *AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPSimpleStateMachine(localLog zerolog.Logger, netstring string, vlan *IPNetwork) (*BIPSimpleStateMachine, error) {
	b := &BIPSimpleStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(netstring), WithClientStateMachineExtension(b), WithLeafType(b))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(netstring))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip simple")
	}
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux")
	}

	// bind the stack together
	if err := Bind(localLog, b, b.bip, b.annexj, b.mux); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}
