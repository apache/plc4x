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

type NetworkLayerStateMachine struct {
	*ClientStateMachine

	address *Address

	log   zerolog.Logger
	codec *NPDUCodec
	node  *Node
}

func NewNetworkLayerStateMachine(localLog zerolog.Logger, address string, vlan *Network) (*NetworkLayerStateMachine, error) {
	n := &NetworkLayerStateMachine{
		log: localLog,
	}
	var err error
	n.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(n), WithLeafType(n))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	n.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creaing address")
	}

	// create a network layer encoder/decoder
	n.codec, err = NewNPDUCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating codec")
	}
	if LogTestNetwork {
		n.log.Debug().Stringer("codec", n.codec).Msg("codec")
	}

	// create a node, added to the network
	n.node, err = NewNode(localLog, n.address, WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	if LogTestNetwork {
		n.log.Debug().Stringer("node", n.node).Msg("node")
	}

	// bind this to the node
	if err := Bind(localLog, n, n.codec, n.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	if !LogTestNetwork {
		n.log = zerolog.Nop()
	}
	return n, nil
}
