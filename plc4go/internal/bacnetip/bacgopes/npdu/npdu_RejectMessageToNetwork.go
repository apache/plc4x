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

func NewRejectMessageToNetwork(args Args, kwArgs KWArgs, options ...Option) (*RejectMessageToNetwork, error) {
	r := &RejectMessageToNetwork{
		messageType: 0x03,
	}
	ApplyAppliers(options, r)
	options = AddLeafTypeIfAbundant(options, r)
	options = AddNLMIfAbundant(options, model.NewNLMRejectMessageToNetwork(r.rmtnRejectionReason, r.rmtnDNET, 0))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	r._NPDU = npdu.(*_NPDU)
	r.AddDebugContents(r, "rmtnRejectReason", "rmtnDNET")

	r.npduNetMessage = &r.messageType
	return r, nil
}

// TODO: check if this is rather a KWArgs
func WithRejectMessageToNetworkRejectionReason(reason model.NLMRejectMessageToNetworkRejectReason) GenericApplier[*RejectMessageToNetwork] {
	return WrapGenericApplier(func(n *RejectMessageToNetwork) { n.rmtnRejectionReason = reason })
}

// TODO: check if this is rather a KWArgs
func WithRejectMessageToNetworkDnet(dnet uint16) GenericApplier[*RejectMessageToNetwork] {
	return WrapGenericApplier(func(n *RejectMessageToNetwork) { n.rmtnDNET = dnet })
}

func (r *RejectMessageToNetwork) GetDebugAttr(attr string) any {
	switch attr {
	case "rmtnRejectReason":
		return r.rmtnRejectionReason
	case "rmtnDNET":
		return r.rmtnDNET
	}
	return nil
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
