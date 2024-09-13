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
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type TrapperRequirements interface {
	BeforeSend(pdu PDU)
	AfterSend(pdu PDU)
	BeforeReceive(pdu PDU)
	AfterReceive(pdu PDU)
	UnexpectedReceive(pdu PDU)
}

// Trapper This class provides a set of utility functions that keeps the latest copy of the pdu parameter in the
// before_send(), after_send(), before_receive(), after_receive() and unexpected_receive() calls.
type Trapper struct {
	TrapperRequirements

	beforeSendPdu        PDU
	afterSendPdu         PDU
	beforeReceivePdu     PDU
	afterReceivePdu      PDU
	unexpectedReceivePdu PDU

	log zerolog.Logger
}

func NewTrapper(localLog zerolog.Logger, requirements TrapperRequirements) *Trapper {
	trapper := &Trapper{
		TrapperRequirements: requirements,
		log:                 localLog,
	}
	if _debug != nil {
		_debug("__init__ %r %r", nil, nil) //TODO: args and kwargs
	}
	// reset to initialize
	trapper.reset()
	return trapper
}

func (t *Trapper) reset() {
	if _debug != nil {
		_debug("reset")
	}
	t.log.Trace().Msg("Reset")
	// flush the copies
	t.beforeSendPdu = nil
	t.afterSendPdu = nil
	t.beforeReceivePdu = nil
	t.afterReceivePdu = nil
	t.unexpectedReceivePdu = nil
}

// BeforeSend is Called before each PDU about to be sent.
func (t *Trapper) BeforeSend(pdu PDU) {
	if _debug != nil {
		_debug("before_send %r", pdu)
	}
	t.log.Debug().Stringer("pdu", pdu).Msg("BeforeSend")
	//keep a copy
	t.beforeSendPdu = pdu

	// continue
	t.TrapperRequirements.BeforeSend(pdu)
}

func (t *Trapper) GetBeforeSendPdu() PDU {
	return t.beforeSendPdu
}

// AfterSend is Called after each PDU sent.
func (t *Trapper) AfterSend(pdu PDU) {
	if _debug != nil {
		_debug("after_send %r", pdu)
	}
	t.log.Debug().Stringer("pdu", pdu).Msg("AfterSend")
	//keep a copy
	t.afterSendPdu = pdu

	// continue
	t.TrapperRequirements.AfterSend(pdu)
}

func (t *Trapper) GetAfterSendPdu() PDU {
	return t.afterSendPdu
}

// BeforeReceive is Called with each PDU received before matching.
func (t *Trapper) BeforeReceive(pdu PDU) {
	if _debug != nil {
		_debug("before_receive %r", pdu)
	}
	t.log.Debug().Stringer("pdu", pdu).Msg("BeforeReceive")
	//keep a copy
	t.beforeReceivePdu = pdu

	// continue
	t.TrapperRequirements.BeforeReceive(pdu)
}

func (t *Trapper) GetBeforeReceivePdu() PDU {
	return t.beforeReceivePdu
}

// AfterReceive is Called with PDU received after match.
func (t *Trapper) AfterReceive(pdu PDU) {
	if _debug != nil {
		_debug("after_receive %r", pdu)
	}
	t.log.Debug().Stringer("pdu", pdu).Msg("AfterReceive")
	//keep a copy
	t.afterReceivePdu = pdu

	// continue
	t.TrapperRequirements.AfterReceive(pdu)
}

func (t *Trapper) GetAfterReceivePdu() PDU {
	return t.afterReceivePdu
}

// UnexpectedReceive is Called with PDU that did not match.  Unless this is trapped by the state, the default behaviour is to fail.
func (t *Trapper) UnexpectedReceive(pdu PDU) {
	if _debug != nil {
		_debug("unexpected_receive %r", pdu)
	}
	t.log.Debug().Stringer("pdu", pdu).Msg("UnexpectedReceive")
	//keep a copy
	t.unexpectedReceivePdu = pdu

	// continue
	t.TrapperRequirements.UnexpectedReceive(pdu)
}

func (t *Trapper) GetUnexpectedReceivePDU() PDU {
	return t.unexpectedReceivePdu
}
