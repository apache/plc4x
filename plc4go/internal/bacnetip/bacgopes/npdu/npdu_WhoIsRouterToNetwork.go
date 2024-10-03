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

type WhoIsRouterToNetwork struct {
	*_NPDU

	messageType uint8

	wirtnNetwork *uint16
}

func NewWhoIsRouterToNetwork(args Args, kwArgs KWArgs, options ...Option) (*WhoIsRouterToNetwork, error) {
	w := &WhoIsRouterToNetwork{
		messageType: 0x00,
	}
	ApplyAppliers(options, w)
	options = AddLeafTypeIfAbundant(options, w)
	options = AddNLMIfAbundant(options, model.NewNLMWhoIsRouterToNetwork(w.wirtnNetwork, 0))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	w._NPDU = npdu.(*_NPDU)
	w.AddDebugContents(w, "wirtnNetwork")

	w.npduNetMessage = &w.messageType
	return w, nil
}

func WithWhoIsRouterToNetworkNet(net uint16) GenericApplier[*WhoIsRouterToNetwork] {
	return WrapGenericApplier(func(n *WhoIsRouterToNetwork) { n.wirtnNetwork = &net })
}

func (w *WhoIsRouterToNetwork) GetDebugAttr(attr string) any {
	switch attr {
	case "wirtnNetwork":
		if w.wirtnNetwork != nil {
			return *w.wirtnNetwork
		}
	}
	return nil
}

func (w *WhoIsRouterToNetwork) GetWirtnNetwork() *uint16 {
	return w.wirtnNetwork
}

func (w *WhoIsRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(w); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		if w.wirtnNetwork != nil {
			npdu.PutShort(*w.wirtnNetwork)
		}
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (w *WhoIsRouterToNetwork) Decode(npdu Arg) error {
	if err := w._NPCI.Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMWhoIsRouterToNetwork:
				w.wirtnNetwork = nlm.GetDestinationNetworkAddress()
				w.SetRootMessage(rm)
			}
		}
	}
	return nil
}
