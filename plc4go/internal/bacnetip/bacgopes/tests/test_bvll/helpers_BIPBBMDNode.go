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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type BIPBBMDNode struct {
	name    string
	address *Address

	bip    *BIPBBMD
	annexj *AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDNode(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPBBMDNode, error) {
	b := &BIPBBMDNode{
		log: localLog,
	}

	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	var err error
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	b.log.Debug().Str("address", address).Msg("address")

	// BACnet/IP interpreter
	b.bip, err = NewBIPBBMD(b.log, b.address)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}
	b.annexj, err = NewAnnexJCodec(b.log)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	b.log.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	bbmdAddress, err := NewAddress(NA(bdtAddress))
	if err != nil {
		return nil, errors.Wrap(err, "error creating bbmd address")
	}
	err = b.bip.AddPeer(bbmdAddress)
	if err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(b.log, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// bind the stack together
	err = Bind(b.log, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}
