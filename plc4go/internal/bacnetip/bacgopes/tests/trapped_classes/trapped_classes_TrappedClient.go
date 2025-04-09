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

// TrappedClientContract provides a set of functions which can be overwritten by a sub struct
type TrappedClientContract interface {
	utils.Serializable
	Request(Args, KWArgs) error
	Confirmation(args Args, kwArgs KWArgs) error
}

// TrappedClient  An instance of this class sits at the top of a stack.
//
//go:generate go tool plc4xGenerator -type=TrappedClient -prefix=trapped_classes_
type TrappedClient struct {
	TrappedClientContract `ignore:"true"`
	ClientContract

	requestSent          PDU
	confirmationReceived PDU

	log zerolog.Logger
}

func NewTrappedClient(localLog zerolog.Logger, options ...Option) (*TrappedClient, error) {
	t := &TrappedClient{
		log: localLog,
	}
	t.TrappedClientContract = t
	ApplyAppliers(options, t)
	optionsForParent := AddLeafTypeIfAbundant(options, t)
	if _debug != nil {
		_debug("__init__")
	}
	var err error
	t.ClientContract, err = NewClient(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error building client")
	}
	return t, nil
}

func WithTrappedClientContract(trappedClientContract TrappedClientContract) func(*TrappedClient) {
	return func(t *TrappedClient) {
		t.TrappedClientContract = trappedClientContract
	}
}

func (t *TrappedClient) GetRequestSent() PDU {
	return t.requestSent
}

func (t *TrappedClient) GetConfirmationReceived() PDU {
	return t.confirmationReceived
}

func (t *TrappedClient) Request(args Args, kwArgs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Request")
	// a reference for checking
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("request %r", pdu)
	}
	t.requestSent = pdu

	// continue with regular processing
	return t.ClientContract.Request(args, kwArgs)
}

func (t *TrappedClient) Confirmation(args Args, kwArgs KWArgs) error {
	t.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Confirmation")
	// a reference for checking
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}
	t.confirmationReceived = pdu
	return nil
}
