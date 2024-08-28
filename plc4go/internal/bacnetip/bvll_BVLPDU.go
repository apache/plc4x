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
	"context"
	"fmt"

	"github.com/pkg/errors"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type BVLPDU interface {
	readWriteModel.BVLC
	BVLCI
	PDUData

	setBVLC(readWriteModel.BVLC)
	getBVLC() readWriteModel.BVLC
}

type _BVLPDU struct {
	*_BVLCI
	*_PDUData

	bvlc readWriteModel.BVLC
}

var _ BVLPDU = (*_BVLPDU)(nil)

func NewBVLPDU(bvlc readWriteModel.BVLC) BVLPDU {
	b := &_BVLPDU{
		bvlc: bvlc,
	}
	b._BVLCI = NewBVLCI(bvlc).(*_BVLCI)
	b._PDUData = NewPDUData(NoArgs).(*_PDUData)
	return b
}

// Deprecated: check if needed as we do it in update
func (b *_BVLPDU) setBVLC(bvlc readWriteModel.BVLC) {
	b.bvlc = bvlc
}

func (b *_BVLPDU) getBVLC() readWriteModel.BVLC {
	return b.bvlc
}

func (b *_BVLPDU) Update(bvlci Arg) error {
	if err := b._BVLCI.Update(bvlci); err != nil {
		return errors.Wrap(err, "Update BVLCI")
	}
	switch bvlci := bvlci.(type) {
	case BVLCI:
		b.bvlc = b.getBVLC()
		// TODO: update coordinates....
		return nil
	default:
		return errors.Errorf("invalid BVLCI type %T", bvlci)
	}
}

func (b *_BVLPDU) Encode(pdu Arg) error {
	if err := b._BVLCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _BVLCI")
	}
	serialize, err := b.bvlc.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing BVLC")
	}
	pdu.(interface{ PutData(n ...byte) }).PutData(serialize...) // TODO: ugly cast...
	return nil
}

func (b *_BVLPDU) Decode(pdu Arg) error {
	if err := b._BVLCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _BVLCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		data := pdu.GetPduData()
		b.PutData(data...)
		var err error
		b.bvlc, err = readWriteModel.BVLCParse(context.Background(), data)
		if err != nil {
			return errors.Wrap(err, "error parsing NPDU")
		}
		b.rootMessage = b.bvlc
	}
	return nil
}

func (b *_BVLPDU) GetBvlcFunction() uint8 {
	if b.bvlc == nil {
		return 0
	}
	return b.bvlc.GetBvlcFunction()
}

func (b *_BVLPDU) GetBvlcPayloadLength() uint16 {
	if b.bvlc == nil {
		return 0
	}
	return b.bvlc.GetBvlcPayloadLength()
}

func (b *_BVLPDU) deepCopy() *_BVLPDU {
	return &_BVLPDU{_BVLCI: b._BVLCI.deepCopy(), _PDUData: b._PDUData.deepCopy(), bvlc: b.bvlc}
}

func (b *_BVLPDU) DeepCopy() any {
	return b.deepCopy()
}

func (b *_BVLPDU) String() string {
	return fmt.Sprintf("_BVLPDU{%s, PDUData: %s}", b._BVLCI, b._PDUData)
}
