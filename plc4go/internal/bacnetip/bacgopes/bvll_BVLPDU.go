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

package bacgopes

import (
	"context"

	"github.com/pkg/errors"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type BVLPDU interface {
	readWriteModel.BVLC
	BVLCI
	PDUData
}

//go:generate plc4xGenerator -type=_BVLPDU -prefix=bvll
type _BVLPDU struct {
	*_BVLCI
	*_PDUData
}

var _ BVLPDU = (*_BVLPDU)(nil)

func NewBVLPDU(bvlc readWriteModel.BVLC) BVLPDU {
	b := &_BVLPDU{}
	b._BVLCI = NewBVLCI(b, bvlc).(*_BVLCI)
	b._PDUData = NewPDUData(NoArgs).(*_PDUData)
	if b.rootMessage != nil {
		b.data, _ = b.rootMessage.Serialize()
		b.data = b.data[4:]
	}
	return b
}

func (b *_BVLPDU) Encode(pdu Arg) error {
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
	var rootMessage spi.Message
	switch pdu := pdu.(type) { // Save a root message as long as we have enough data
	case PDUData:
		rootMessage, _ = readWriteModel.BVLCParse[readWriteModel.BVLC](context.Background(), pdu.GetPduData())
	}
	if err := b._BVLCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _BVLCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		b.PutData(pdu.GetPduData()...)
	}
	if rootMessage != nil {
		// Overwrite the root message again so we can use it for matching
		b.rootMessage = rootMessage
	}
	return nil
}

func (b *_BVLPDU) GetBvlcFunction() uint8 {
	switch rm := b.rootMessage.(type) {
	case readWriteModel.BVLC:
		return rm.GetBvlcFunction()
	default:
		return 0
	}
}

func (b *_BVLPDU) GetBvlcPayloadLength() uint16 {
	switch rm := b.rootMessage.(type) {
	case readWriteModel.BVLC:
		return rm.GetBvlcPayloadLength()
	default:
		return 0
	}
}

func (b *_BVLPDU) IsBVLC() {
}

func (b *_BVLPDU) deepCopy() *_BVLPDU {
	return &_BVLPDU{_BVLCI: b._BVLCI.deepCopy(), _PDUData: b._PDUData.deepCopy()}
}

func (b *_BVLPDU) DeepCopy() any {
	return b.deepCopy()
}
