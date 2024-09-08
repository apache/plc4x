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
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type RegisterForeignDevice struct {
	*_BVLPDU

	bvlciTimeToLive uint16
}

var _ BVLPDU = (*RegisterForeignDevice)(nil)

func NewRegisterForeignDevice(opts ...func(RegisterForeignDevice *RegisterForeignDevice)) (*RegisterForeignDevice, error) {
	b := &RegisterForeignDevice{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(model.NewBVLCRegisterForeignDevice(b.bvlciTimeToLive)).(*_BVLPDU)
	return b, nil
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
		case model.BVLCRegisterForeignDevice:
			switch bvlc := rm.(type) {
			case model.BVLCRegisterForeignDevice:
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

func (r *RegisterForeignDevice) String() string {
	if r == nil {
		return "(*RegisterForeignDevice)(nil)"
	}
	return fmt.Sprintf("RegisterForeignDevice{%v, bvlciTimeToLive: %v}", r._BVLPDU, r.bvlciTimeToLive)
}
