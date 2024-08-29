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
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

type TrapperRequirements interface {
	BeforeSend(pdu bacnetip.PDU)
	AfterSend(pdu bacnetip.PDU)
	BeforeReceive(pdu bacnetip.PDU)
	AfterReceive(pdu bacnetip.PDU)
	UnexpectedReceive(pdu bacnetip.PDU)
}

// Trapper This class provides a set of utility functions that keeps the latest copy of the pdu parameter in the
// before_send(), after_send(), before_receive(), after_receive() and unexpected_receive() calls.
type Trapper struct {
	TrapperRequirements

	beforeSendPdu        bacnetip.PDU
	afterSendPdu         bacnetip.PDU
	beforeReceivePdu     bacnetip.PDU
	afterReceivePdu      bacnetip.PDU
	unexpectedReceivePdu bacnetip.PDU

	log zerolog.Logger
}

func NewTrapper(localLog zerolog.Logger, requirements TrapperRequirements) *Trapper {
	trapper := &Trapper{
		TrapperRequirements: requirements,
		log:                 localLog,
	}
	// reset to initialize
	trapper.reset()
	return trapper
}

func (t *Trapper) reset() {
	t.log.Trace().Msg("Reset")
	// flush the copies
	t.beforeSendPdu = nil
	t.afterSendPdu = nil
	t.beforeReceivePdu = nil
	t.afterReceivePdu = nil
	t.unexpectedReceivePdu = nil
}

// BeforeSend is Called before each PDU about to be sent.
func (t *Trapper) BeforeSend(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("BeforeSend")
	//keep a copy
	t.beforeSendPdu = pdu

	// continue
	t.TrapperRequirements.BeforeSend(pdu)
}

func (t *Trapper) GetBeforeSendPdu() bacnetip.PDU {
	return t.beforeSendPdu
}

// AfterSend is Called after each PDU sent.
func (t *Trapper) AfterSend(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("AfterSend")
	//keep a copy
	t.afterSendPdu = pdu

	// continue
	t.TrapperRequirements.AfterSend(pdu)
}

func (t *Trapper) GetAfterSendPdu() bacnetip.PDU {
	return t.afterSendPdu
}

// BeforeReceive is Called with each PDU received before matching.
func (t *Trapper) BeforeReceive(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("BeforeReceive")
	//keep a copy
	t.beforeReceivePdu = pdu

	// continue
	t.TrapperRequirements.BeforeReceive(pdu)
}

func (t *Trapper) GetBeforeReceivePdu() bacnetip.PDU {
	return t.beforeReceivePdu
}

// AfterReceive is Called with PDU received after match.
func (t *Trapper) AfterReceive(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("AfterReceive")
	//keep a copy
	t.afterReceivePdu = pdu

	// continue
	t.TrapperRequirements.AfterReceive(pdu)
}

func (t *Trapper) GetAfterReceivePdu() bacnetip.PDU {
	return t.afterReceivePdu
}

// UnexpectedReceive is Called with PDU that did not match.  Unless this is trapped by the state, the default behaviour is to fail.
func (t *Trapper) UnexpectedReceive(pdu bacnetip.PDU) {
	t.log.Debug().Stringer("pdu", pdu).Msg("UnexpectedReceive")
	//keep a copy
	t.unexpectedReceivePdu = pdu

	// continue
	t.TrapperRequirements.UnexpectedReceive(pdu)
}

func (t *Trapper) GetUnexpectedReceivePDU() bacnetip.PDU {
	return t.unexpectedReceivePdu
}
