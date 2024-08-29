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

type ReadForeignDeviceTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadForeignDeviceTable)(nil)

func NewReadForeignDeviceTable(opts ...func(*ReadForeignDeviceTable)) (*ReadForeignDeviceTable, error) {
	b := &ReadForeignDeviceTable{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(model.NewBVLCReadForeignDeviceTable()).(*_BVLPDU)
	return b, nil
}

func (w *ReadForeignDeviceTable) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadForeignDeviceTable) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCReadForeignDeviceTableExactly:
			switch bvlc := rm.(type) {
			case model.BVLCReadForeignDeviceTable:
				w.setBVLC(bvlc)
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadForeignDeviceTable) String() string {
	return fmt.Sprintf("ReadForeignDeviceTable{%v}", w._BVLPDU)
}
