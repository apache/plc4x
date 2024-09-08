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

package npdu

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type RejectMessageToNetwork struct {
	*_NPDU

	messageType uint8

	rmtnRejectionReason model.NLMRejectMessageToNetworkRejectReason
	rmtnDNET            uint16
}

func NewRejectMessageToNetwork(opts ...func(*RejectMessageToNetwork)) (*RejectMessageToNetwork, error) {
	i := &RejectMessageToNetwork{
		messageType: 0x03,
	}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMRejectMessageToNetwork(i.rmtnRejectionReason, i.rmtnDNET, 0), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)

	i.npduNetMessage = &i.messageType
	return i, nil
}

func WithRejectMessageToNetworkRejectionReason(reason model.NLMRejectMessageToNetworkRejectReason) func(*RejectMessageToNetwork) {
	return func(n *RejectMessageToNetwork) {
		n.rmtnRejectionReason = reason
	}
}

func WithRejectMessageToNetworkDnet(dnet uint16) func(*RejectMessageToNetwork) {
	return func(n *RejectMessageToNetwork) {
		n.rmtnDNET = dnet
	}
}

func (r *RejectMessageToNetwork) GetRmtnRejectionReason() model.NLMRejectMessageToNetworkRejectReason {
	return r.rmtnRejectionReason
}

func (r *RejectMessageToNetwork) GetRmtnDNET() uint16 {
	return r.rmtnDNET
}

func (r *RejectMessageToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.Put(byte(r.rmtnRejectionReason))
		npdu.PutShort(r.rmtnDNET)
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (r *RejectMessageToNetwork) Decode(npdu Arg) error {
	if err := r._NPCI.Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		if err := r.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMRejectMessageToNetwork:
				r.rmtnRejectionReason = nlm.GetRejectReason()
				r.rmtnDNET = nlm.GetDestinationNetworkAddress()
				r.SetRootMessage(rm)
			}
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		r.SetPduData(npdu.GetPduData())
	}
	return nil
}

func (r *RejectMessageToNetwork) String() string {
	return fmt.Sprintf("RejectMessageToNetwork{%s, rmtnRejectionReason: %s, rmtnDNET: %v}", r._NPDU, r.rmtnRejectionReason, r.rmtnDNET)
}
