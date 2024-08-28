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

type IAmRouterToNetwork struct {
	*_NPDU

	iartnNetworkList []uint16
}

func NewIAmRouterToNetwork(opts ...func(*IAmRouterToNetwork)) (*IAmRouterToNetwork, error) {
	i := &IAmRouterToNetwork{}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMIAmRouterToNetwork(i.iartnNetworkList, 0), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	return i, nil
}

func WithIAmRouterToNetworkNetworkList(iartnNetworkList ...uint16) func(*IAmRouterToNetwork) {
	return func(n *IAmRouterToNetwork) {
		n.iartnNetworkList = iartnNetworkList
	}
}

func (i *IAmRouterToNetwork) GetIartnNetworkList() []uint16 {
	return i.iartnNetworkList
}

func (i *IAmRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(i); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		for _, net := range i.iartnNetworkList {
			npdu.PutShort(net)
		}
		npdu.setNLM(i.nlm)
		npdu.setAPDU(i.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (i *IAmRouterToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := i.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMIAmRouterToNetworkExactly:
				i.setNLM(nlm)
				i.iartnNetworkList = nlm.GetDestinationNetworkAddresses()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (i *IAmRouterToNetwork) String() string {
	return fmt.Sprintf("IAmRouterToNetwork{%s, iartnNetworkList: %v}", i._NPDU, i.iartnNetworkList)
}
