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

package trapped_classes

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TrappedServerContract provides a set of functions which can be overwritten by a sub struct
type TrappedServerContract interface {
	utils.Serializable
	Indication(args Args, kwArgs KWArgs) error
	Response(Args, KWArgs) error
}

// TrappedServer An instance of this class sits at the bottom of a stack.
//
//go:generate go tool plc4xGenerator -type=TrappedServer -prefix=trapped_classes_
type TrappedServer struct {
	TrappedServerContract `ignore:"true"`
	ServerContract

	indicationReceived PDU
	responseSent       PDU

	log zerolog.Logger
}

func NewTrappedServer(localLog zerolog.Logger, options ...Option) (*TrappedServer, error) {
	t := &TrappedServer{
		log: localLog,
	}
	t.TrappedServerContract = t
	ApplyAppliers(options, t)
	optionsForParent := AddLeafTypeIfAbundant(options, t)
	if _debug != nil {
		_debug("__init__")
	}
	var err error
	t.ServerContract, err = NewServer(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error building server")
	}
	return t, nil
}

func WithTrappedServerContract(trappedServerContract TrappedServerContract) GenericApplier[*TrappedServer] {
	return WrapGenericApplier(func(t *TrappedServer) { t.TrappedServerContract = trappedServerContract })
}

func (t *TrappedServer) GetIndicationReceived() PDU {
	return t.indicationReceived
}

func (t *TrappedServer) GetResponseSent() PDU {
	return t.responseSent
}

func (t *TrappedServer) Indication(args Args, kwArgs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Indication")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("indication %r", pdu)
	}

	// a reference for checking
	t.indicationReceived = pdu
	return nil
}

func (t *TrappedServer) Response(args Args, kwArgs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Response")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("response %r", pdu)
	}

	// a reference for checking
	t.responseSent = pdu

	// continue with regular processing
	return t.ServerContract.Response(args, kwArgs)
}
