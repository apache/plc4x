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
		// TODO: re-encode seems expensive... check if there is a better option (e.g. only do it on the message bridge)
		data := pdu.GetPduData()
		parse, err := model.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return nil, errors.Wrap(err, "error re-encoding")
		}
		b._BVLPDU = NewBVLPDU(model.NewBVLCForwardedNPDU(b.produceInnerNPDU(parse))).(*_BVLPDU)
	}
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

func (w *ForwardedNPDU) GetBvlciAddress() *Address {
	return w.bvlciAddress
}

func (w *ForwardedNPDU) produceInnerNPDU(inNpdu model.NPDU) (ip []uint8, port uint16, npdu model.NPDU, bvlcPayloadLength uint16) {
	ip = w.bvlciAddress.AddrAddress[:4]
	port = uint16(47808)
	if w.bvlciAddress.AddrPort != nil {
		port = *w.bvlciAddress.AddrPort
	}
	npdu = inNpdu
	return
}

func (w *ForwardedNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}

		// encode the addrress
		bvlpdu.PutData(w.bvlciAddress.AddrAddress...)

		// encode the rest of the data
		bvlpdu.PutData(w.GetPduData()...)

		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ForwardedNPDU) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCForwardedNPDU:
			addr := rm.GetIp()
			port := rm.GetPort()
			var portArray = make([]byte, 2)
			binary.BigEndian.PutUint16(portArray, port)
			var err error
			address, err := NewAddress(zerolog.Nop(), append(addr, portArray...))
			if err != nil {
				return errors.Wrap(err, "error creating address")
			}
			w.bvlciAddress = address
			npdu := rm.GetNpdu()
			pduData, err := npdu.Serialize()
			if err != nil {
				return errors.Wrap(err, "error serializing NPDU")
			}
			w.SetPduData(pduData)

			w.setBVLC(rm)
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ForwardedNPDU) String() string {
	return fmt.Sprintf("ForwardedNPDU{%v, bvlciAddress: %v}", w._BVLPDU, w.bvlciAddress)
}
