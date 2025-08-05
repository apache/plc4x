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
	"context"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type NPDU interface {
	Copyable
	readWriteModel.NPDU
	NPCI
	PDUData
}

type _NPDU struct {
	*_NPCI
	PDUData
}

var _ = (NPDU)(nil)

func NewNPDU(args Args, kwArgs KWArgs, options ...Option) (NPDU, error) {
	n := &_NPDU{}
	options = AddLeafTypeIfAbundant(options, n)
	n._NPCI = NewNPCI(args, kwArgs, options...).(*_NPCI)
	n.PDUData = NewPDUData(NoArgs, NoKWArgs(), options...)
	n.AddExtraPrinters(n.PDUData.(DebugContentPrinter))
	if n.GetRootMessage() != nil {
		data, _ := n.GetRootMessage().Serialize()
		n.SetPduData(data)
	}
	return n, nil
}

func (n *_NPDU) Encode(pdu Arg) error {
	if err := n._NPCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _NPCI")
	}
	switch pdu := pdu.(type) {
	case PDU:
		pdu.PutData(n.GetPduData()...)
	}
	return nil
}

func (n *_NPDU) Decode(pdu Arg) error {
	var rootMessage spi.Message
	switch pdu := pdu.(type) { // Save a root message as long as we have enough data
	case PDUData:
		data := pdu.GetPduData()
		rootMessage, _ = readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
	}
	switch pdu := pdu.(type) {
	case IPCI:
		if rootMessage != nil { // in this case we are good and want to parse from that
			pdu.SetRootMessage(rootMessage)
		}
	}
	if err := n._NPCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _NPCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		n.PutData(pdu.GetPduData()[n.bytesToDiscard:]...)
	}
	if rootMessage != nil {
		// Overwrite the root message again so we can use it for matching
		n.SetRootMessage(rootMessage)
	}
	return nil
}

func (n *_NPDU) getNPDUModel() (readWriteModel.NPDU, bool) {
	rm := n.GetRootMessage()
	npdu, ok := rm.(readWriteModel.NPDU)
	return npdu, ok
}

func (n *_NPDU) CreateNPDUBuilder() readWriteModel.NPDUBuilder {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return readWriteModel.NewNPDUBuilder()
	}
	return npdu.CreateNPDUBuilder()
}

func (n *_NPDU) GetProtocolVersionNumber() uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetProtocolVersionNumber()
	}
	return npdu.GetProtocolVersionNumber()
}

func (n *_NPDU) GetControl() readWriteModel.NPDUControl {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetControl()
	}
	return npdu.GetControl()
}

func (n *_NPDU) GetDestinationNetworkAddress() *uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetDestinationNetworkAddress()
	}
	return npdu.GetDestinationNetworkAddress()
}

func (n *_NPDU) GetDestinationLength() *uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetDestinationLength()
	}
	return npdu.GetDestinationLength()
}

func (n *_NPDU) GetDestinationAddress() []uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetDestinationAddress()
	}
	return npdu.GetDestinationAddress()
}

func (n *_NPDU) GetSourceNetworkAddress() *uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetSourceNetworkAddress()
	}
	return npdu.GetSourceNetworkAddress()
}

func (n *_NPDU) GetSourceLength() *uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetSourceLength()
	}
	return npdu.GetSourceLength()
}

func (n *_NPDU) GetSourceAddress() []uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetSourceAddress()
	}
	return npdu.GetSourceAddress()
}

func (n *_NPDU) GetHopCount() *uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return n._NPCI.GetHopCount()
	}
	return npdu.GetHopCount()
}

func (n *_NPDU) GetNlm() readWriteModel.NLM {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetNlm()
}

func (n *_NPDU) GetApdu() readWriteModel.APDU {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetApdu()
}

func (n *_NPDU) GetDestinationLengthAddon() uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0 // TODO: implement me
	}
	return npdu.GetDestinationLengthAddon()
}

func (n *_NPDU) GetSourceLengthAddon() uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0 // TODO: implement me
	}
	return npdu.GetSourceLengthAddon()
}

func (n *_NPDU) GetPayloadSubtraction() uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0 // TODO: implement me
	}
	return npdu.GetPayloadSubtraction()
}

func (n *_NPDU) IsNPDU() {
}

func (n *_NPDU) deepCopy() *_NPDU {
	return &_NPDU{_NPCI: n._NPCI.deepCopy(), PDUData: n.PDUData.DeepCopy().(PDUData)}
}

func (n *_NPDU) DeepCopy() any {
	return n.deepCopy()
}
