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

type ICouldBeRouterToNetwork struct {
	*_NPDU

	icbrtnNetwork          uint16
	icbrtnPerformanceIndex uint8
}

func NewICouldBeRouterToNetwork(opts ...func(*ICouldBeRouterToNetwork)) (*ICouldBeRouterToNetwork, error) {
	i := &ICouldBeRouterToNetwork{}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMICouldBeRouterToNetwork(i.icbrtnNetwork, i.icbrtnPerformanceIndex, 0), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	return i, nil
}

func WithICouldBeRouterToNetworkNetwork(icbrtnNetwork uint16) func(*ICouldBeRouterToNetwork) {
	return func(n *ICouldBeRouterToNetwork) {
		n.icbrtnNetwork = icbrtnNetwork
	}
}

func WithICouldBeRouterToNetworkPerformanceIndex(icbrtnPerformanceIndex uint8) func(*ICouldBeRouterToNetwork) {
	return func(n *ICouldBeRouterToNetwork) {
		n.icbrtnPerformanceIndex = icbrtnPerformanceIndex
	}
}

func (i *ICouldBeRouterToNetwork) GetIcbrtnNetwork() uint16 {
	return i.icbrtnNetwork
}

func (i *ICouldBeRouterToNetwork) GetIcbrtnPerformanceIndex() uint8 {
	return i.icbrtnPerformanceIndex
}

func (i *ICouldBeRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(i); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		npdu.PutShort(i.icbrtnNetwork)
		npdu.Put(i.icbrtnPerformanceIndex)
		npdu.setNPDU(i.npdu)
		npdu.setNLM(i.nlm)
		npdu.setAPDU(i.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (i *ICouldBeRouterToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := i.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMICouldBeRouterToNetworkExactly:
				i.setNLM(nlm)
				i.icbrtnNetwork = nlm.GetDestinationNetworkAddress()
				i.icbrtnPerformanceIndex = nlm.GetPerformanceIndex()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (i *ICouldBeRouterToNetwork) String() string {
	return fmt.Sprintf("ICouldBeRouterToNetwork{%s, icbrtnNetwork: %v, icbrtnPerformanceIndex: %v}", i._NPDU, i.icbrtnNetwork, i.icbrtnPerformanceIndex)
}
