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

type RouterAvailableToNetwork struct {
	*_NPDU

	messageType uint8

	ratnNetworkList []uint16
}

func NewRouterAvailableToNetwork(args Args, kwArgs KWArgs, opts ...func(*RouterAvailableToNetwork)) (*RouterAvailableToNetwork, error) {
	i := &RouterAvailableToNetwork{
		messageType: 0x05,
	}
	for _, opt := range opts {
		opt(i)
	}
	if _, ok := kwArgs[KWCompNLM]; ok {
		kwArgs[KWCompNLM] = model.NewNLMRouterAvailableToNetwork(i.ratnNetworkList, 0)
	}
	npdu, err := NewNPDU(args, kwArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)

	i.npduNetMessage = &i.messageType
	return i, nil
}

func WithRouterAvailableToNetworkDnet(networkList []uint16) func(*RouterAvailableToNetwork) {
	return func(n *RouterAvailableToNetwork) {
		n.ratnNetworkList = networkList
	}
}

func (r *RouterAvailableToNetwork) GetRatnNetworkList() []uint16 {
	return r.ratnNetworkList
}

func (r *RouterAvailableToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		for _, net := range r.GetRatnNetworkList() {
			npdu.PutShort(net)
		}
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (r *RouterAvailableToNetwork) Decode(npdu Arg) error {
	if err := r._NPCI.Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMRouterAvailableToNetwork:
				r.ratnNetworkList = nlm.GetDestinationNetworkAddresses()
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

func (r *RouterAvailableToNetwork) String() string {
	return fmt.Sprintf("RouterAvailableToNetwork{%s, ratnNetworkList: %v}", r._NPDU, r.ratnNetworkList)
}
