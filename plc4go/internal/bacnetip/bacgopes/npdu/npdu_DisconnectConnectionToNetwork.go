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

type DisconnectConnectionToNetwork struct {
	*_NPDU

	messageType uint8

	dctnDNET uint16
}

func NewDisconnectConnectionToNetwork(args Args, kwArgs KWArgs, opts ...func(*DisconnectConnectionToNetwork)) (*DisconnectConnectionToNetwork, error) {
	i := &DisconnectConnectionToNetwork{
		messageType: 0x09,
	}
	for _, opt := range opts {
		opt(i)
	}
	if _, ok := kwArgs[KWCompNLM]; ok {
		kwArgs[KWCompNLM] = model.NewNLMDisconnectConnectionToNetwork(i.dctnDNET, 0)
	}
	npdu, err := NewNPDU(args, kwArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)

	i.npduNetMessage = &i.messageType
	return i, nil
}

func WithDisconnectConnectionToNetworkDNET(dnet uint16) func(*DisconnectConnectionToNetwork) {
	return func(n *DisconnectConnectionToNetwork) {
		n.dctnDNET = dnet
	}
}

func (n *DisconnectConnectionToNetwork) GetDctnDNET() uint16 {
	return n.dctnDNET
}

func (n *DisconnectConnectionToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(n); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.PutShort(n.dctnDNET)
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (n *DisconnectConnectionToNetwork) Decode(npdu Arg) error {
	if err := n._NPCI.Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMDisconnectConnectionToNetwork:
				n.dctnDNET = nlm.GetDestinationNetworkAddress()
				n.SetRootMessage(rm)
			}
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		n.SetPduData(npdu.GetPduData())
	}
	return nil
}

func (n *DisconnectConnectionToNetwork) String() string {
	return fmt.Sprintf("DisconnectConnectionToNetwork{%s, dctnDNET: %v}", n._NPDU, n.dctnDNET)
}
