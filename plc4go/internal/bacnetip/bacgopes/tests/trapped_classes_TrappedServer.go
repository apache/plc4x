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

package tests

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/utils"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

// TrappedServerContract provides a set of functions which can be overwritten by a sub struct
type TrappedServerContract interface {
	utils.Serializable
	Indication(args Args, kwargs KWArgs) error
	Response(Args, KWArgs) error
}

// TrappedServer An instance of this class sits at the bottom of a stack.
//
//go:generate plc4xGenerator -type=TrappedServer -prefix=trapped_classes_
type TrappedServer struct {
	TrappedServerContract `ignore:"true"`
	Server

	indicationReceived PDU
	responseSent       PDU

	log zerolog.Logger
}

func NewTrappedServer(localLog zerolog.Logger, opts ...func(*TrappedServer)) (*TrappedServer, error) {
	t := &TrappedServer{
		log: localLog,
	}
	t.TrappedServerContract = t
	for _, opt := range opts {
		opt(t)
	}
	var err error
	t.Server, err = NewServer(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "error building server")
	}
	return t, nil
}

func WithTrappedServerContract(trappedServerContract TrappedServerContract) func(*TrappedServer) {
	return func(t *TrappedServer) {
		t.TrappedServerContract = trappedServerContract
	}
}

func (t *TrappedServer) GetIndicationReceived() PDU {
	return t.indicationReceived
}

func (t *TrappedServer) GetResponseSent() PDU {
	return t.responseSent
}

func (t *TrappedServer) Indication(args Args, kwargs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	// a reference for checking
	t.indicationReceived = Get[PDU](args, 0)

	return nil
}

func (t *TrappedServer) Response(args Args, kwargs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Response")
	// a reference for checking
	t.responseSent = Get[PDU](args, 0)

	// continue with regular processing
	return t.Server.Response(args, kwargs)
}
