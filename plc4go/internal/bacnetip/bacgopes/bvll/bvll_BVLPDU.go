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

package bvll

import (
	"context"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type BVLPDU interface {
	readWriteModel.BVLC
	BVLCI
	PDUData
}

type _BVLPDU struct {
	*_BVLCI
	PDUData
}

var _ BVLPDU = (*_BVLPDU)(nil)

func NewBVLPDU(args Args, kwArgs KWArgs, options ...Option) BVLPDU {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	b := &_BVLPDU{}
	options = AddLeafTypeIfAbundant(options, b)
	b._BVLCI = NewBVLCI(b, args, kwArgs, options...).(*_BVLCI)
	b.PDUData = NewPDUData(args, kwArgs, options...)
	b.AddExtraPrinters(b.PDUData.(DebugContentPrinter))
	if b.GetRootMessage() != nil {
		data, _ := b.GetRootMessage().Serialize()
		b.SetPduData(data[4:])
	}
	return b
}

func (b *_BVLPDU) Encode(pdu Arg) error {
	if _debug != nil {
		_debug("encode %s", pdu)
	}
	if err := b._BVLCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _BVLCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		pdu.PutData(b.GetPduData()...)
	}
	return nil
}

func (b *_BVLPDU) Decode(pdu Arg) error {
	if _debug != nil {
		_debug("decode %s", pdu)
	}
	var rootMessage spi.Message
	switch pdu := pdu.(type) { // Save a root message as long as we have enough data
	case PDUData:
		rootMessage, _ = Try1(func() (readWriteModel.BVLC, error) {
			return readWriteModel.BVLCParse[readWriteModel.BVLC](context.Background(), pdu.GetPduData())
		})
	}
	switch pdu := pdu.(type) {
	case IPCI:
		if rootMessage != nil { // in this case we are good and want to parse from that
			pdu.SetRootMessage(rootMessage)
		}
	}
	if err := b._BVLCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _BVLCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		b.PutData(pdu.GetPduData()[b.bytesToDiscard:]...)
	}
	if rootMessage != nil {
		// Overwrite the root message again so we can use it for matching
		b.SetRootMessage(rootMessage)
	}
	return nil
}

func (b *_BVLPDU) CreateBVLCBuilder() readWriteModel.BVLCBuilder {
	switch rm := b.GetRootMessage().(type) {
	case readWriteModel.BVLC:
		return rm.CreateBVLCBuilder()
	default:
		return readWriteModel.NewBVLCBuilder()
	}
}

func (b *_BVLPDU) GetBvlcFunction() uint8 {
	switch rm := b.GetRootMessage().(type) {
	case readWriteModel.BVLC:
		return rm.GetBvlcFunction()
	default:
		return b.bvlciFunction
	}
}

func (b *_BVLPDU) GetBvlcPayloadLength() uint16 {
	switch rm := b.GetRootMessage().(type) {
	case readWriteModel.BVLC:
		return rm.GetBvlcPayloadLength()
	default:
		return b.bvlciLength
	}
}

func (b *_BVLPDU) IsBVLC() {
}

func (b *_BVLPDU) deepCopy() *_BVLPDU {
	return &_BVLPDU{_BVLCI: b._BVLCI.deepCopy(), PDUData: b.PDUData.DeepCopy().(PDUData)}
}

func (b *_BVLPDU) DeepCopy() any {
	return b.deepCopy()
}
