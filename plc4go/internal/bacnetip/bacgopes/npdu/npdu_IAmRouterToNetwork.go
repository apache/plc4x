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

type IAmRouterToNetwork struct {
	*_NPDU

	messageType uint8

	iartnNetworkList []uint16
}

func NewIAmRouterToNetwork(args Args, kwArgs KWArgs, options ...Option) (*IAmRouterToNetwork, error) {
	i := &IAmRouterToNetwork{
		messageType: 0x01,
	}
	ApplyAppliers(options, i)
	options = AddLeafTypeIfAbundant(options, i)
	options = AddNLMIfAbundant(options, model.NewNLMIAmRouterToNetwork(i.iartnNetworkList, 0))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	i.AddDebugContents(i, "iartnNetworkList")

	i.npduNetMessage = &i.messageType
	return i, nil
}

// TODO: check if this is rather a KWArgs
func WithIAmRouterToNetworkNetworkList(iartnNetworkList ...uint16) GenericApplier[*IAmRouterToNetwork] {
	return WrapGenericApplier(func(n *IAmRouterToNetwork) { n.iartnNetworkList = iartnNetworkList })
}

func (i *IAmRouterToNetwork) GetDebugAttr(attr string) any {
	switch attr {
	case "iartnNetworkList":
		return i.iartnNetworkList
	}
	return nil
}

func (i *IAmRouterToNetwork) GetIartnNetworkList() []uint16 {
	return i.iartnNetworkList
}

func (i *IAmRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(i); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		for _, net := range i.iartnNetworkList {
			npdu.PutShort(net)
		}
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (i *IAmRouterToNetwork) Decode(npdu Arg) error {
	if err := i._NPCI.Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMIAmRouterToNetwork:
				i.iartnNetworkList = nlm.GetDestinationNetworkAddresses()
				i.SetRootMessage(rm)
			}
		}
	}
	return nil
}
