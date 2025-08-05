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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

//go:generate go tool plc4xGenerator -type=NPDUCodec -prefix=helpers_
type NPDUCodec struct {
	ClientContract
	ServerContract

	log zerolog.Logger
}

func NewNPDUCodec(localLog zerolog.Logger) (*NPDUCodec, error) {
	n := &NPDUCodec{
		log: localLog,
	}
	var err error
	n.ClientContract, err = NewClient(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	n.ServerContract, err = NewServer(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	if !LogTestNetwork {
		n.log = zerolog.Nop()
	}
	return n, nil
}

func (n *NPDUCodec) Indication(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")

	npdu := GA[NPDU](args, 0)

	// first a generic _NPDU
	xpdu, err := NewNPDU(Nothing())
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding xpdu")
	}

	// Now as a vanilla PDU
	ypdu := NewPDU(Nothing())
	if err := xpdu.Encode(ypdu); err != nil {
		return errors.Wrap(err, "error decoding xpdu")
	}
	n.log.Debug().Stringer("ypdu", ypdu).Msg("encoded")

	// send it downstream
	return n.Request(NA(ypdu), NoKWArgs())
}

func (n *NPDUCodec) Confirmation(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")

	pdu := GA[PDU](args, 0)

	// decode as generic _NPDU
	xpdu, err := NewNPDU(Nothing())
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err := xpdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding xpdu")
	}

	// drop application layer message
	if xpdu.GetNPDUNetMessage() == nil {
		n.log.Trace().Msg("drop message")
		return nil
	}

	// do a deeper decode of the _NPDU
	ypdu, err := NPDUTypes[*xpdu.GetNPDUNetMessage()](Nothing())
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err := ypdu.Decode(xpdu); err != nil {
		return errors.Wrap(err, "error decoding ypdu")
	}

	return n.Response(NA(ypdu), NoKWArgs())
}
