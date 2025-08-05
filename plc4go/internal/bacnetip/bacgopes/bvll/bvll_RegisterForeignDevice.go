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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type RegisterForeignDevice struct {
	*_BVLPDU

	bvlciTimeToLive uint16
}

var _ BVLPDU = (*RegisterForeignDevice)(nil)

// TODO: check this args desaster
func NewRegisterForeignDevice(ttl *uint16, args Args, kwArgs KWArgs, options ...Option) (*RegisterForeignDevice, error) {
	r := &RegisterForeignDevice{}
	options = AddLeafTypeIfAbundant(options, r)
	r._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	if r.GetRootMessage() == nil {
		r.SetRootMessage(readWriteModel.NewBVLCRegisterForeignDevice(r.bvlciTimeToLive))
	}
	r.AddDebugContents(r, "bvlciTimeToLive")
	r.bvlciFunction = BVLCIRegisterForeignDevice
	r.bvlciLength = 6
	if ttl != nil {
		r.bvlciTimeToLive = *ttl
	}
	return r, nil
}

func (r *RegisterForeignDevice) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciTimeToLive":
		return r.bvlciTimeToLive
	default:
		return nil
	}
}

func WithRegisterForeignDeviceBvlciTimeToLive(ttl uint16) func(*RegisterForeignDevice) {
	return func(b *RegisterForeignDevice) {
		b.bvlciTimeToLive = ttl
	}
}

func (r *RegisterForeignDevice) GetBvlciTimeToLive() uint16 {
	return r.bvlciTimeToLive
}

func (r *RegisterForeignDevice) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		bvlpdu.PutShort(r.bvlciTimeToLive)
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (r *RegisterForeignDevice) Decode(bvlpdu Arg) error {
	if err := r._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCRegisterForeignDevice:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCRegisterForeignDevice:
				r.bvlciTimeToLive = bvlc.GetTtl()
				r.SetRootMessage(rm)
			}
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		r.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}
