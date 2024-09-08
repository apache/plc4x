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
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type DeleteForeignDeviceTableEntry struct {
	*_BVLPDU

	bvlciAddress *Address
}

var _ BVLPDU = (*DeleteForeignDeviceTableEntry)(nil)

func NewDeleteForeignDeviceTableEntry(opts ...func(*DeleteForeignDeviceTableEntry)) (*DeleteForeignDeviceTableEntry, error) {
	d := &DeleteForeignDeviceTableEntry{}
	for _, opt := range opts {
		opt(d)
	}
	d._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCDeleteForeignDeviceTableEntry(d.buildIPArgs())).(*_BVLPDU)
	return d, nil
}

func WithDeleteForeignDeviceTableEntryAddress(address *Address) func(*DeleteForeignDeviceTableEntry) {
	return func(d *DeleteForeignDeviceTableEntry) {
		d.bvlciAddress = address
	}
}

func (d *DeleteForeignDeviceTableEntry) buildIPArgs() (ip []uint8, port uint16) {
	if d.bvlciAddress == nil {
		return
	}
	ip = d.bvlciAddress.AddrAddress[:4]
	port = *d.bvlciAddress.AddrPort
	return
}

func (d *DeleteForeignDeviceTableEntry) buildAddress(ip []uint8, port uint16) *Address {
	var portArray = make([]byte, 2)
	binary.BigEndian.PutUint16(portArray, port)
	address, _ := NewAddress(zerolog.Nop(), append(ip, portArray...))
	return address
}

func (d *DeleteForeignDeviceTableEntry) GetBvlciAddress() *Address {
	return d.bvlciAddress
}

func (d *DeleteForeignDeviceTableEntry) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(d); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		bvlpdu.PutData(d.bvlciAddress.AddrAddress...)
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (d *DeleteForeignDeviceTableEntry) Decode(bvlpdu Arg) error {
	if err := d._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCDeleteForeignDeviceTableEntry:
			d.bvlciAddress = d.buildAddress(rm.GetIp(), rm.GetPort())
			d.SetRootMessage(rm)
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		d.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (d *DeleteForeignDeviceTableEntry) String() string {
	if d == nil {
		return "(*DeleteForeignDeviceTableEntry)(nil)"
	}
	return fmt.Sprintf("DeleteForeignDeviceTableEntry{%v}", d._BVLPDU)
}
