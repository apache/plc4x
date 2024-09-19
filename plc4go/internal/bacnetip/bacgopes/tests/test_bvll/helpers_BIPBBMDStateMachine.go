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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type BIPBBMDStateMachine struct {
	*ClientStateMachine

	address *Address
	bip     *BIPBBMD
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPBBMDStateMachine, error) {
	b := &BIPBBMDStateMachine{
		log: localLog,
	}
	if _debug != nil {
		_debug("__init__ %r %r", address, vlan)
	}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(b), WithLeafType(b))
	if err != nil {
		return nil, errors.New("error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPBBMD(localLog, b.address)
	b.annexj, err = NewAnnexJCodec(localLog)

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	if _debug != nil {
		_debug("    - bdt_address: %r", bdtAddress)
	}
	b.log.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	if err := b.bip.AddPeer(quick.Address(bdtAddress)); err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// bind the stack together
	err = Bind(b.log, b, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}
