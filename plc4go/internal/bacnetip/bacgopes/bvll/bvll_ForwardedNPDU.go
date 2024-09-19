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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ForwardedNPDU struct {
	*_BVLPDU

	bvlciAddress *Address
}

var _ BVLPDU = (*ForwardedNPDU)(nil)

// TODO: check this args desaster...
func NewForwardedNPDU(addr *Address, args Args, kwArgs KWArgs, options ...Option) (*ForwardedNPDU, error) {
	f := &ForwardedNPDU{}
	options = AddLeafTypeIfAbundant(options, f)
	f._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	f.AddDebugContents(f, "bvlciAddress")
	switch npdu := f.GetRootMessage().(type) {
	case readWriteModel.NPDU:
		// Repackage
		f.SetRootMessage(readWriteModel.NewBVLCForwardedNPDU(f.produceInnerNPDU(npdu)))
	}
	f.bvlciFunction = BVLCIForwardedNPDU
	f.bvlciLength = uint16(10 + len(f.GetPduData()))
	f.bvlciAddress = addr
	return f, nil
}

func (f *ForwardedNPDU) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciAddress":
		if f.bvlciAddress != nil {
			return f.bvlciAddress
		}
	default:
		return nil
	}
	return nil
}

func (f *ForwardedNPDU) GetBvlciAddress() *Address {
	return f.bvlciAddress
}

func (f *ForwardedNPDU) produceInnerNPDU(inNpdu readWriteModel.NPDU) (ip []uint8, port uint16, npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
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
		case readWriteModel.BVLCForwardedNPDU:
			addr := rm.GetIp()
			port := rm.GetPort()
			var portArray = make([]byte, 2)
			binary.BigEndian.PutUint16(portArray, port)
			var err error
			address, err := NewAddress(NA(append(addr, portArray...)))
			if err != nil {
				return errors.Wrap(err, "error creating address")
			}
			f.bvlciAddress = address
			f.SetRootMessage(rm)
		}

		// get the address
		data, err := bvlpdu.GetData(6)
		if err != nil {
			return errors.Wrap(err, "error reading data")
		}
		f.bvlciAddress, err = NewAddress(NA(UnpackIpAddr(data)))
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		f.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (f *ForwardedNPDU) String() string {
	if f == nil {
		return "(*ForwardedNPDU)(nil)"
	}
	return fmt.Sprintf("ForwardedNPDU{%v, bvlciAddress: %v}", f._BVLPDU, f.bvlciAddress)
}
