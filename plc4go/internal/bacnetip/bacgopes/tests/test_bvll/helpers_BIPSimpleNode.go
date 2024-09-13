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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type BIPSimpleNode struct {
	name    string
	address *Address
	bip     *BIPSimple
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPSimpleNode(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPSimpleNode, error) {
	b := &BIPSimpleNode{}

	// save the name and address
	b.name = address
	var err error
	b.address, err = NewAddress(NA(address))
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

	// bind the stack together
	err = Bind(localLog, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}
