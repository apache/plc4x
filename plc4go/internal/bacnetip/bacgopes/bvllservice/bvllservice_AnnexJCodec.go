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

package bvllservice

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

//go:generate go tool plc4xGenerator -type=AnnexJCodec -prefix=bvllservice_
type AnnexJCodec struct {
	ClientContract
	ServerContract
	*DefaultRFormatter `ignore:"true"`

	log zerolog.Logger
}

func NewAnnexJCodec(localLog zerolog.Logger, options ...Option) (*AnnexJCodec, error) {
	a := &AnnexJCodec{
		DefaultRFormatter: NewDefaultRFormatter(),
		log:               localLog,
	}
	ApplyAppliers(options, a)
	optionsForParent := AddLeafTypeIfAbundant(options, a)
	var err error
	a.ClientContract, err = NewClient(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	a.ServerContract, err = NewServer(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	if _debug != nil {
		_debug("__init__ cid=%r sid=%r", a.GetClientID(), a.GetServerId())
	}
	localLog.Debug().
		Interface("cid", a.GetClientID()).
		Interface("sid", a.GetServerId()).
		Msg("NewAnnexJCodec")
	return a, nil
}

func (b *AnnexJCodec) Indication(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Indication")
	rpdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("indication %r", rpdu)
	}

	// encode it as a generic BVLL PDU
	bvlpdu := NewBVLPDU(Nothing())
	if err := rpdu.(Encoder).Encode(bvlpdu); err != nil {
		return errors.Wrap(err, "error encoding PDU")
	}

	// encode it as a PDU
	pdu := NewPDU(Nothing())
	if err := bvlpdu.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding PDU")
	}

	// send it downstream
	return b.Request(NA(pdu), NoKWArgs())
}

func (b *AnnexJCodec) Confirmation(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}

	// interpret as a BVLL PDU
	bvlpdu := NewBVLPDU(Nothing())
	if err := bvlpdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding pdu")
	}

	// get the class related to the function
	rpdu := BVLPDUTypes[bvlpdu.GetBvlcFunction()]()
	if err := rpdu.Decode(bvlpdu); err != nil {
		return errors.Wrap(err, "error decoding PDU")
	}

	// send it upstream
	return b.Response(NA(rpdu), NoKWArgs())
}
