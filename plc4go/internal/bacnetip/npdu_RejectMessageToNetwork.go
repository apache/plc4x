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

package bacnetip

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type RejectMessageToNetwork struct {
	*_NPDU

	rmtnRejectionReason model.NLMRejectMessageToNetworkRejectReason
	rmtnDNET            uint16
}

func NewRejectMessageToNetwork(opts ...func(*RejectMessageToNetwork)) (*RejectMessageToNetwork, error) {
	i := &RejectMessageToNetwork{}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMRejectMessageToNetwork(i.rmtnRejectionReason, i.rmtnDNET, 0), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
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

func (n *RejectMessageToNetwork) GetRmtnRejectionReason() model.NLMRejectMessageToNetworkRejectReason {
	return n.rmtnRejectionReason
}

func (n *RejectMessageToNetwork) GetRmtnDNET() uint16 {
	return n.rmtnDNET
}

func (n *RejectMessageToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		npdu.Put(byte(n.rmtnRejectionReason))
		npdu.PutShort(n.rmtnDNET)
		npdu.setNLM(n.nlm)
		npdu.setAPDU(n.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *RejectMessageToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := n.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMRejectMessageToNetworkExactly:
				n.setNLM(nlm)
				n.rmtnRejectionReason = nlm.GetRejectReason()
				n.rmtnDNET = nlm.GetDestinationNetworkAddress()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *RejectMessageToNetwork) String() string {
	return fmt.Sprintf("RejectMessageToNetwork{%s, rmtnRejectionReason: %s, rmtnDNET: %v}", n._NPDU, n.rmtnRejectionReason, n.rmtnDNET)
}
