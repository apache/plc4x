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

type WhoIsRouterToNetwork struct {
	*_NPDU

	messageType uint8

	wirtnNetwork *uint16
}

func NewWhoIsRouterToNetwork(opts ...func(network *WhoIsRouterToNetwork)) (*WhoIsRouterToNetwork, error) {
	w := &WhoIsRouterToNetwork{
		messageType: 0x00,
	}
	for _, opt := range opts {
		opt(w)
	}
	npdu, err := NewNPDU(model.NewNLMWhoIsRouterToNetwork(w.wirtnNetwork, 0), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	w._NPDU = npdu.(*_NPDU)

	w.npduNetMessage = &w.messageType
	return w, nil
}

func WithWhoIsRouterToNetworkNet(net uint16) func(*WhoIsRouterToNetwork) {
	return func(n *WhoIsRouterToNetwork) {
		n.wirtnNetwork = &net
	}
}

func (n *WhoIsRouterToNetwork) GetWirtnNetwork() *uint16 {
	return n.wirtnNetwork
}

func (n *WhoIsRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		if n.wirtnNetwork != nil {
			npdu.PutShort(*n.wirtnNetwork)
		}
		npdu.setNLM(n.nlm)
		npdu.setAPDU(n.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *WhoIsRouterToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := n.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMWhoIsRouterToNetworkExactly:
				n.setNLM(nlm)
				n.wirtnNetwork = nlm.GetDestinationNetworkAddress()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *WhoIsRouterToNetwork) String() string {
	return fmt.Sprintf("WhoIsRouterToNetwork{%s, wirtnNetwork: %d}", n._NPDU, n.wirtnNetwork)
}
