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
)

type BVLCIRequirements interface {
	GetPduData() []uint8
}

type BVLCI interface {
	PCI

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setBvlciType(bvlciType uint8)
	getBvlciType() uint8
	setBvlciFunction(bvlciFunction uint8)
	getBvlciFunction() uint8
	setBvlciLength(bvlciLength uint16)
	getBvlciLength() uint16

	getBVLCI() BVLCI
}

type _BVLCI struct {
	PCI
	*DebugContents

	_requirements BVLCIRequirements

	bvlciType     uint8
	bvlciFunction uint8
	bvlciLength   uint16

	// Deprecated: hacky workaround
	bytesToDiscard int
}

var _ BVLCI = (*_BVLCI)(nil)

func NewBVLCI(requirements BVLCIRequirements, args Args, kwArgs KWArgs, options ...Option) BVLCI {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	b := &_BVLCI{
		_requirements: requirements,

		bvlciType: 0x81,
	}
	options = AddLeafTypeIfAbundant(options, b)
	b.DebugContents = NewDebugContents(b, "bvlciType", "bvlciFunction", "bvlciLength")
	b.PCI = NewPCI(args, kwArgs)
	b.AddExtraPrinters(b.PCI.(DebugContentPrinter))
	if bvlc, ok := ExtractRootMessage(options).(readWriteModel.BVLC); ok {
		b.bvlciFunction = bvlc.GetBvlcFunction()
		b.bvlciLength = bvlc.GetLengthInBytes(context.Background())
	}
	return b
}

func (b *_BVLCI) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciType":
		return b.bvlciType
	case "bvlciFunction":
		return b.bvlciFunction
	case "bvlciLength":
		return b.bvlciLength
	}
	return nil
}

func (b *_BVLCI) Update(bvlci Arg) error {
	if err := b.PCI.Update(bvlci); err != nil {
		return errors.Wrap(err, "Update BVLCI")
	}
	switch bvlci := bvlci.(type) {
	case BVLCI:
		b.bvlciType = bvlci.getBvlciType()
		b.bvlciFunction = bvlci.getBvlciFunction()
		b.bvlciLength = bvlci.getBvlciLength()
		return nil
	default:
		return errors.Errorf("invalid BVLCI type %T", bvlci)
	}
}

func (b *_BVLCI) Encode(pdu Arg) error {
	if _debug != nil {
		_debug("encode %s", pdu)
	}
	switch pdu := pdu.(type) {
	case PCI:
		if err := pdu.GetPCI().Update(b); err != nil {
			return errors.Wrap(err, "error updating pdu")
		}
	}
	switch pdu := pdu.(type) {
	case PDUData:
		pdu.Put(b.bvlciType)
		pdu.Put(b.bvlciFunction)

		if int(b.bvlciLength) != len(b._requirements.GetPduData())+4 {
			return errors.Errorf("invalid BVLCI length %d != %d", b.bvlciLength, len(b._requirements.GetPduData())+4)
		}

		pdu.PutShort(b.bvlciLength)
	}
	return nil
}

func (b *_BVLCI) Decode(pdu Arg) error {
	if _debug != nil {
		_debug("decode %s", pdu)
	}
	if err := b.PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	readBytes := 0 // TODO: as long as we use the read like this we should be good (not using plc4x to parse that away)
	switch pdu := pdu.(type) {
	case PDUData:
		var err error
		b.bvlciType, err = pdu.Get()
		if err != nil {
			return errors.Wrap(err, "error reading bvlci type")
		}
		if b.bvlciType != 0x81 {
			return errors.New("invalid BVLCI type")
		}

		b.bvlciFunction, err = pdu.Get()
		if err != nil {
			return errors.Wrap(err, "error reading bvlci function")
		}
		bvlciLength, err := pdu.GetShort()
		if err != nil {
			return errors.Wrap(err, "error reading bvlci length")
		}
		b.bvlciLength = uint16(bvlciLength)
		if int(b.bvlciLength) != len(pdu.GetPduData())+4 {
			return errors.Errorf("invalid BVLCI length %d != %d", b.bvlciLength, len(pdu.GetPduData())+4)
		}
	}
	b.bytesToDiscard = readBytes
	return nil
}

func (b *_BVLCI) setBvlciType(bvlciType uint8) {
	b.bvlciType = bvlciType
}

func (b *_BVLCI) getBvlciType() uint8 {
	return b.bvlciType
}

func (b *_BVLCI) setBvlciFunction(bvlciFunction uint8) {
	b.bvlciFunction = bvlciFunction
}

func (b *_BVLCI) getBvlciFunction() uint8 {
	return b.bvlciFunction
}

func (b *_BVLCI) setBvlciLength(bvlciLength uint16) {
	b.bvlciLength = bvlciLength
}

func (b *_BVLCI) getBvlciLength() uint16 {
	return b.bvlciLength
}

func (b *_BVLCI) getBVLCI() BVLCI {
	return b
}

func (b *_BVLCI) deepCopy() *_BVLCI {
	newB := &_BVLCI{
		b.PCI.DeepCopy().(PCI),
		nil,
		nil, // TODO: what to do with that? Clone will be useless...
		b.bvlciType,
		b.bvlciFunction,
		b.bvlciLength,
		b.bytesToDiscard,
	}
	newB.DebugContents = NewDebugContents(newB, "bvlciType", "bvlciFunction", "bvlciLength") // TODO: bit ugly to repeat that here again but what are the options...
	newB.AddExtraPrinters(newB.PCI.(DebugContentPrinter))
	return newB
}
