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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type ForwardedNPDU struct {
	*_BVLPDU

	bvlciAddress *Address
}

var _ BVLPDU = (*ForwardedNPDU)(nil)

func NewForwardedNPDU(pdu PDU, opts ...func(*ForwardedNPDU)) (*ForwardedNPDU, error) {
	b := &ForwardedNPDU{}
	for _, opt := range opts {
		opt(b)
	}
	switch npdu := pdu.(type) {
	case model.NPDU:
		b._BVLPDU = NewBVLPDU(model.NewBVLCForwardedNPDU(b.produceInnerNPDU(npdu))).(*_BVLPDU)
	case nil:
		b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	default:
		data := pdu.GetPduData()
		parsedNPDU, err := model.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			b._BVLPDU = NewBVLPDU(model.NewBVLCForwardedNPDU(b.produceInnerNPDU(parsedNPDU))).(*_BVLPDU)
		} else {
			b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
			b._BVLPDU.SetPduData(data)
		}
	}
	b.bvlciFunction = 0x04
	b.bvlciLength = uint16(10 + len(b.GetPduData()))
	return b, nil
}

func WithForwardedNPDUAddress(addr *Address) func(*ForwardedNPDU) {
	return func(b *ForwardedNPDU) {
		b.bvlciAddress = addr
	}
}

func WithForwardedNPDUUserData(userData spi.Message) func(*ForwardedNPDU) {
	return func(b *ForwardedNPDU) {
		b.pduUserData = userData
	}
}

func (f *ForwardedNPDU) GetBvlciAddress() *Address {
	return f.bvlciAddress
}

func (f *ForwardedNPDU) produceInnerNPDU(inNpdu model.NPDU) (ip []uint8, port uint16, npdu model.NPDU, bvlcPayloadLength uint16) {
	ip = f.bvlciAddress.AddrAddress[:4]
	port = uint16(47808)
	if f.bvlciAddress.AddrPort != nil {
		port = *f.bvlciAddress.AddrPort
	}
	npdu = inNpdu
	return
}

func (f *ForwardedNPDU) Encode(bvlpdu Arg) error {
	// make sure the length is correct
	f.bvlciLength = uint16(10 + len(f.GetPduData()))

	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(f); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:

		// encode the addrress
		bvlpdu.PutData(f.bvlciAddress.AddrAddress...)

		// encode the rest of the data
		bvlpdu.PutData(f.GetPduData()...)
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (f *ForwardedNPDU) Decode(bvlpdu Arg) error {
	if err := f._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCForwardedNPDU:
			_, _ = bvlpdu.GetData(6) // TODO: do we really want to discard that?
			addr := rm.GetIp()
			port := rm.GetPort()
			var portArray = make([]byte, 2)
			binary.BigEndian.PutUint16(portArray, port)
			var err error
			address, err := NewAddress(zerolog.Nop(), append(addr, portArray...))
			if err != nil {
				return errors.Wrap(err, "error creating address")
			}
			f.bvlciAddress = address
			f.rootMessage = rm
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		f.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (f *ForwardedNPDU) String() string {
	return fmt.Sprintf("ForwardedNPDU{%v, bvlciAddress: %v}", f._BVLPDU, f.bvlciAddress)
}
