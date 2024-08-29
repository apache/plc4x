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
	"fmt"

	"github.com/pkg/errors"

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

func (n *RegisterForeignDevice) GetBvlciTimeToLive() uint16 {
	return n.bvlciTimeToLive
}

func (n *RegisterForeignDevice) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.PutShort(n.bvlciTimeToLive)
		bvlpdu.setBVLC(n.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *RegisterForeignDevice) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := n.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCRegisterForeignDeviceExactly:
			switch bvlc := rm.(type) {
			case model.BVLCRegisterForeignDevice:
				n.setBVLC(bvlc)
				n.bvlciTimeToLive = bvlc.GetTtl()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *RegisterForeignDevice) String() string {
	return fmt.Sprintf("RegisterForeignDevice{%v, bvlciTimeToLive: %v}", n._BVLPDU, n.bvlciTimeToLive)
}
