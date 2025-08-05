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

type NetworkNumberIs struct {
	*_NPDU

	messageType uint8

	nniNet  uint16
	nniFlag bool
}

func NewNetworkNumberIs(args Args, kwArgs KWArgs, options ...Option) (*NetworkNumberIs, error) {
	n := &NetworkNumberIs{
		messageType: 0x13,
	}
	ApplyAppliers(options, n)
	options = AddLeafTypeIfAbundant(options, n)
	options = AddNLMIfAbundant(options, model.NewNLMNetworkNumberIs(n.nniNet, n.nniFlag, 0))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	n._NPDU = npdu.(*_NPDU)
	n.AddDebugContents(n, "nniNet", "nniFlag")

	n.npduNetMessage = &n.messageType
	return n, nil
}

// TODO: check if this is rather a KWArgs
func WithNetworkNumberIsNET(net uint16) GenericApplier[*NetworkNumberIs] {
	return WrapGenericApplier(func(n *NetworkNumberIs) { n.nniNet = net })
}

// TODO: check if this is rather a KWArgs
func WithNetworkNumberIsTerminationConfigured(configured bool) GenericApplier[*NetworkNumberIs] {
	return WrapGenericApplier(func(n *NetworkNumberIs) { n.nniFlag = configured })
}

func (n *NetworkNumberIs) GetDebugAttr(attr string) any {
	switch attr {
	case "nniNet":
		return n.nniNet
	case "nniFlag":
		return n.nniFlag
	}
	return nil
}

func (n *NetworkNumberIs) GetNniNet() uint16 {
	return n.nniNet
}

func (n *NetworkNumberIs) GetNniFlag() bool {
	return n.nniFlag
}

func (n *NetworkNumberIs) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(n); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.PutShort(n.nniNet)
		flag := uint8(0)
		if n.nniFlag {
			flag = 1
		}
		npdu.Put(flag)
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (n *NetworkNumberIs) Decode(npdu Arg) error {
	if err := n.GetNPCI().Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		if err := n.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMNetworkNumberIs:
				n.nniNet = nlm.GetNetworkNumber()
				n.nniFlag = nlm.GetNetworkNumberConfigured()
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
