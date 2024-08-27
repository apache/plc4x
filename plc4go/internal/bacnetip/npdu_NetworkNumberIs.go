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

type NetworkNumberIs struct {
	*_NPDU

	nniNet  uint16
	nniFlag bool
}

func NewNetworkNumberIs(opts ...func(*NetworkNumberIs)) (*NetworkNumberIs, error) {
	i := &NetworkNumberIs{}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMNetworkNumberIs(i.nniNet, i.nniFlag, 0), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	return i, nil
}

func WithNetworkNumberIsNET(net uint16) func(*NetworkNumberIs) {
	return func(n *NetworkNumberIs) {
		n.nniNet = net
	}
}

func WithNetworkNumberIsTerminationConfigured(configured bool) func(*NetworkNumberIs) {
	return func(n *NetworkNumberIs) {
		n.nniFlag = configured
	}
}

func (n *NetworkNumberIs) GetNniNet() uint16 {
	return n.nniNet
}

func (n *NetworkNumberIs) GetNniFlag() bool {
	return n.nniFlag
}

func (n *NetworkNumberIs) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		npdu.PutShort(n.nniNet)
		flag := uint8(0)
		if n.nniFlag {
			flag = 1
		}
		npdu.Put(flag)
		npdu.setNPDU(n.npdu)
		npdu.setNLM(n.nlm)
		npdu.setAPDU(n.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *NetworkNumberIs) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := n.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMNetworkNumberIsExactly:
				n.setNLM(nlm)
				n.nniNet = nlm.GetNetworkNumber()
				n.nniFlag = nlm.GetNetworkNumberConfigured()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *NetworkNumberIs) String() string {
	return fmt.Sprintf("NetworkNumberIs{%s, nniNet: %v, nniFlag: %v}", n._NPDU, n.nniNet, n.nniFlag)
}
