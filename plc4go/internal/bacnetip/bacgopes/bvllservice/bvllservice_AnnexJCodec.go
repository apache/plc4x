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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

//go:generate plc4xGenerator -type=AnnexJCodec -prefix=bvllservice_
type AnnexJCodec struct {
	Client
	Server

	// pass through args
	argCid *int `ignore:"true"`
	argSid *int `ignore:"true"`

	log zerolog.Logger
}

func NewAnnexJCodec(localLog zerolog.Logger, opts ...func(*AnnexJCodec)) (*AnnexJCodec, error) {
	a := &AnnexJCodec{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}
	localLog.Debug().
		Interface("cid", a.argCid).
		Interface("sid", a.argSid).
		Msg("NewAnnexJCodec")
	client, err := NewClient(localLog, a, OptionalOption(a.argCid, WithClientCID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	a.Client = client
	server, err := NewServer(localLog, a, OptionalOption(a.argSid, WithServerSID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	a.Server = server
	return a, nil
}

func WithAnnexJCodecCid(cid int) func(*AnnexJCodec) {
	return func(a *AnnexJCodec) {
		a.argCid = &cid
	}
}

func WithAnnexJCodecSid(sid int) func(*AnnexJCodec) {
	return func(a *AnnexJCodec) {
		a.argSid = &sid
	}
}

func (b *AnnexJCodec) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")

	rpdu := Get[PDU](args, 0)

	// encode it as a generic BVLL PDU
	bvlpdu := NewBVLPDU(nil)
	if err := rpdu.(Encoder).Encode(bvlpdu); err != nil {
		return errors.Wrap(err, "error encoding PDU")
	}

	// encode it as a PDU
	pdu := NewPDU(nil)
	if err := bvlpdu.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding PDU")
	}

	// send it downstream
	return b.Request(NewArgs(pdu), NoKWArgs)
}

func (b *AnnexJCodec) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")

	pdu := Get[PDU](args, 0)

	// interpret as a BVLL PDU
	bvlpdu := NewBVLPDU(nil)
	if err := bvlpdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding pdu")
	}

	// get the class related to the function
	rpdu := BVLPDUTypes[bvlpdu.GetBvlcFunction()]()
	if err := rpdu.Decode(bvlpdu); err != nil {
		return errors.Wrap(err, "error decoding PDU")
	}

	// send it upstream
	return b.Response(NewArgs(rpdu), NoKWArgs)
}
