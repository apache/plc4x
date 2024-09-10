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
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type OriginalUnicastNPDU struct {
	*_BVLPDU

	// post construct function
	_postConstruct []func()
}

var _ BVLPDU = (*OriginalUnicastNPDU)(nil)

func NewOriginalUnicastNPDU(pdu PDU, opts ...func(*OriginalUnicastNPDU)) (*OriginalUnicastNPDU, error) {
	o := &OriginalUnicastNPDU{}
	for _, opt := range opts {
		opt(o)
	}
	switch npdu := pdu.(type) {
	case readWriteModel.NPDU:
		o._BVLPDU = NewBVLPDU(NoArgs, NewKWArgs(KWCompRootMessage, readWriteModel.NewBVLCOriginalUnicastNPDU(o.produceInnerNPDU(npdu)))).(*_BVLPDU)
	case nil:
		o._BVLPDU = NewBVLPDU(Nothing()).(*_BVLPDU)
	default:
		// TODO: re-encode seems expensive... check if there is a better option (e.g. only do it on the message bridge)
		data := pdu.GetPduData()
		parse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return nil, errors.Wrap(err, "error re-encoding")
		}
		o._BVLPDU = NewBVLPDU(NoArgs, NewKWArgs(KWCompRootMessage, readWriteModel.NewBVLCOriginalUnicastNPDU(o.produceInnerNPDU(parse)))).(*_BVLPDU)
	}
	// Do a post construct for a bit more easy initialization
	for _, f := range o._postConstruct {
		f()
	}
	o._postConstruct = nil
	return o, nil
}

func WithOriginalUnicastNPDUDestination(destination *Address) func(*OriginalUnicastNPDU) {
	return func(o *OriginalUnicastNPDU) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUDestination(destination)
		})
	}
}

func WithOriginalUnicastNPDUUserData(userData spi.Message) func(*OriginalUnicastNPDU) {
	return func(o *OriginalUnicastNPDU) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUUserData(userData)
		})
	}
}

func (o *OriginalUnicastNPDU) produceInnerNPDU(inNpdu readWriteModel.NPDU) (npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
	npdu = inNpdu
	return
}

func (o *OriginalUnicastNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(o); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		bvlpdu.PutData(o.GetPduData()...)
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (o *OriginalUnicastNPDU) Decode(bvlpdu Arg) error {
	if err := o._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCOriginalUnicastNPDU:
			o.SetRootMessage(rm)
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		o.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (o *OriginalUnicastNPDU) String() string {
	if o == nil {
		return "(*OriginalUnicastNPDU)(nil)"
	}
	return fmt.Sprintf("OriginalUnicastNPDU{%s}", o._BVLPDU)
}
