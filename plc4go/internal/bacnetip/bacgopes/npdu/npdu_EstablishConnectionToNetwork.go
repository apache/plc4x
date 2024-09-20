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

type EstablishConnectionToNetwork struct {
	*_NPDU

	messageType uint8

	ectnDNET            uint16
	ectnTerminationTime uint8
}

func NewEstablishConnectionToNetwork(args Args, kwArgs KWArgs, options ...Option) (*EstablishConnectionToNetwork, error) {
	e := &EstablishConnectionToNetwork{
		messageType: 0x08,
	}
	ApplyAppliers(options, e)
	options = AddLeafTypeIfAbundant(options, e)
	options = AddNLMIfAbundant(options, model.NewNLMEstablishConnectionToNetwork(e.ectnDNET, e.ectnTerminationTime, 0))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	e._NPDU = npdu.(*_NPDU)
	e.AddDebugContents(e, "ectnDNET", "ectnTerminationTime")

	e.npduNetMessage = &e.messageType
	return e, nil
}

// TODO: check if this is rather a KWArgs
func WithEstablishConnectionToNetworkDNET(dnet uint16) GenericApplier[*EstablishConnectionToNetwork] {
	return WrapGenericApplier(func(n *EstablishConnectionToNetwork) { n.ectnDNET = dnet })
}

// TODO: check if this is rather a KWArgs
func WithEstablishConnectionToNetworkTerminationTime(terminationTime uint8) GenericApplier[*EstablishConnectionToNetwork] {
	return WrapGenericApplier(func(n *EstablishConnectionToNetwork) { n.ectnTerminationTime = terminationTime })
}

func (e *EstablishConnectionToNetwork) GetDebugAttr(attr string) any {
	switch attr {
	case "ectnDNET":
		return e.ectnDNET
	case "ectnTerminationTime":
		return e.ectnTerminationTime
	}
	return nil
}

func (e *EstablishConnectionToNetwork) GetEctnDNET() uint16 {
	return e.ectnDNET
}

func (e *EstablishConnectionToNetwork) GetEctnTerminationTime() uint8 {
	return e.ectnTerminationTime
}

func (e *EstablishConnectionToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(e); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.PutShort(e.ectnDNET)
		npdu.Put(e.ectnTerminationTime)
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (e *EstablishConnectionToNetwork) Decode(npdu Arg) error {
	if err := e._NPCI.Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMEstablishConnectionToNetwork:
				e.ectnDNET = nlm.GetDestinationNetworkAddress()
				e.ectnTerminationTime = nlm.GetTerminationTime()
				e.SetRootMessage(rm)
			}
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		e.SetPduData(npdu.GetPduData())
	}
	return nil
}
