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
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type OriginalUnicastNPDU struct {
	*_BVLPDU

	// post construct function
	_postConstruct []func()
}

var _ BVLPDU = (*OriginalUnicastNPDU)(nil)

func NewOriginalUnicastNPDU(args Args, kwArgs KWArgs, options ...Option) (*OriginalUnicastNPDU, error) {
	o := &OriginalUnicastNPDU{}
	options = AddLeafTypeIfAbundant(options, o)
	o._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	switch npdu := o.GetRootMessage().(type) {
	case readWriteModel.NPDU:
		// Repackage
		o.SetRootMessage(readWriteModel.NewBVLCOriginalUnicastNPDU(o.produceInnerNPDU(npdu)))
	}
	o.bvlciFunction = BVLCIOriginalUnicastNPDU
	o.bvlciLength = uint16(4 + len(o.GetPduData()))
	return o, nil
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
