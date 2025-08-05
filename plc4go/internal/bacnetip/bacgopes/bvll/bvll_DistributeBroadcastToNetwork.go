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

type DistributeBroadcastToNetwork struct {
	*_BVLPDU
}

var _ BVLPDU = (*DistributeBroadcastToNetwork)(nil)

func NewDistributeBroadcastToNetwork(args Args, kwArgs KWArgs, options ...Option) (*DistributeBroadcastToNetwork, error) {
	d := &DistributeBroadcastToNetwork{}
	options = AddLeafTypeIfAbundant(options, d)
	d._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	switch npdu := d.GetRootMessage().(type) {
	case readWriteModel.NPDU:
		// Repackage
		d.SetRootMessage(readWriteModel.NewBVLCDistributeBroadcastToNetwork(d.produceInnerNPDU(npdu)))
	}
	d.bvlciFunction = BVLCIDistributeBroadcastToNetwork
	d.bvlciLength = uint16(4 + len(d.GetPduData()))
	return d, nil
}

func (d *DistributeBroadcastToNetwork) produceInnerNPDU(inNpdu readWriteModel.NPDU) (npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
	npdu = inNpdu
	return
}

func (d *DistributeBroadcastToNetwork) Encode(bvlpdu Arg) error {
	d.bvlciLength = uint16(4 + len(d.GetPduData()))
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(d); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}

	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		bvlpdu.PutData(d.GetPduData()...)
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (d *DistributeBroadcastToNetwork) Decode(bvlpdu Arg) error {
	if err := d._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCDistributeBroadcastToNetwork:
			d.SetRootMessage(rm)
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		data, err := bvlpdu.GetData(len(bvlpdu.GetPduData()))
		if err != nil {
			return errors.Wrap(err, "error getting data")
		}
		d.SetPduData(data)
	}
	return nil
}

func (d *DistributeBroadcastToNetwork) String() string {
	if d == nil {
		return "(*DistributeBroadcastToNetwork)(nil)"
	}
	return fmt.Sprintf("DistributeBroadcastToNetwork{%s}", d._BVLPDU)
}
