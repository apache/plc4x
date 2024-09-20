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

type ICouldBeRouterToNetwork struct {
	*_NPDU

	messageType uint8

	icbrtnNetwork          uint16
	icbrtnPerformanceIndex uint8
}

func NewICouldBeRouterToNetwork(args Args, kwArgs KWArgs, options ...Option) (*ICouldBeRouterToNetwork, error) {
	i := &ICouldBeRouterToNetwork{
		messageType: 0x02,
	}
	ApplyAppliers(options, i)
	options = AddLeafTypeIfAbundant(options, i)
	options = AddNLMIfAbundant(options, model.NewNLMICouldBeRouterToNetwork(i.icbrtnNetwork, i.icbrtnPerformanceIndex, 0))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	i.AddDebugContents(i, "icbrtnNetwork", "icbrtnPerformanceIndex")

	i.npduNetMessage = &i.messageType
	return i, nil
}

// TODO: check if this is rather a KWArgs
func WithICouldBeRouterToNetworkNetwork(icbrtnNetwork uint16) GenericApplier[*ICouldBeRouterToNetwork] {
	return WrapGenericApplier(func(n *ICouldBeRouterToNetwork) { n.icbrtnNetwork = icbrtnNetwork })
}

// TODO: check if this is rather a KWArgs
func WithICouldBeRouterToNetworkPerformanceIndex(icbrtnPerformanceIndex uint8) GenericApplier[*ICouldBeRouterToNetwork] {
	return WrapGenericApplier(func(n *ICouldBeRouterToNetwork) { n.icbrtnPerformanceIndex = icbrtnPerformanceIndex })
}

func (i *ICouldBeRouterToNetwork) GetDebugAttr(attr string) any {
	switch attr {
	case "icbrtnNetwork":
		return i.icbrtnNetwork
	case "icbrtnPerformanceIndex":
		return i.icbrtnPerformanceIndex
	}
	return nil
}

func (i *ICouldBeRouterToNetwork) GetIcbrtnNetwork() uint16 {
	return i.icbrtnNetwork
}

func (i *ICouldBeRouterToNetwork) GetIcbrtnPerformanceIndex() uint8 {
	return i.icbrtnPerformanceIndex
}

func (i *ICouldBeRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(i); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.PutShort(i.icbrtnNetwork)
		npdu.Put(i.icbrtnPerformanceIndex)
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (i *ICouldBeRouterToNetwork) Decode(npdu Arg) error {
	if err := i.GetNPCI().Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMICouldBeRouterToNetwork:
				i.icbrtnNetwork = nlm.GetDestinationNetworkAddress()
				i.icbrtnPerformanceIndex = nlm.GetPerformanceIndex()
				i.SetRootMessage(rm)
			}
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		i.SetPduData(npdu.GetPduData())
	}
	return nil
}
