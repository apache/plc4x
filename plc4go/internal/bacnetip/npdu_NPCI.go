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
	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type NPCI interface {
	PCI

	GetNPDUNetMessage() *uint8

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setNLM(model.NLM)
	getNLM() model.NLM
}

type _NPCI struct {
	*_PCI
	*DebugContents

	nlm model.NLM
}

var _ NPCI = (*_NPCI)(nil)

func NewNPCI(pduUserData spi.Message, nlm model.NLM) NPCI {
	n := &_NPCI{
		nlm: nlm,
	}
	n._PCI = newPCI(pduUserData, nil, nil, nil, false, model.NPDUNetworkPriority_NORMAL_MESSAGE)
	switch nlm := pduUserData.(type) {
	case model.NLMExactly:
		n.nlm = nlm
	}
	return n
}

func (n *_NPCI) GetNPDUNetMessage() *uint8 {
	if n.nlm == nil {
		return nil
	}
	messageType := n.nlm.GetMessageType()
	return &messageType
}

// Deprecated: check if needed as we do it in update
func (n *_NPCI) setNLM(nlm model.NLM) {
	n.nlm = nlm
}

func (n *_NPCI) getNLM() model.NLM {
	return n.nlm
}

func (n *_NPCI) Update(npci Arg) error {
	if err := n._PCI.Update(npci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch npci := npci.(type) {
	case NPCI:
		n.nlm = npci.getNLM()
		// TODO: update coordinates...
		return nil
	default:
		return errors.Errorf("invalid NPCI type %T", npci)
	}
}

func (n *_NPCI) Encode(pdu Arg) error {
	if err := pdu.(interface{ Update(Arg) error }).Update(n); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_NPCI) Decode(pdu Arg) error {
	if err := n._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_NPCI) deepCopy() *_NPCI {
	return &_NPCI{_PCI: n._PCI.deepCopy(), nlm: n.nlm}
}
