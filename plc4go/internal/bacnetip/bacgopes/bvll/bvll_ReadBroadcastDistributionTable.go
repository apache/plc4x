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

type ReadBroadcastDistributionTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadBroadcastDistributionTable)(nil)

func NewReadBroadcastDistributionTable(args Args, kwArgs KWArgs, options ...Option) (*ReadBroadcastDistributionTable, error) {
	r := &ReadBroadcastDistributionTable{}
	options = AddLeafTypeIfAbundant(options, r)
	r._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	if r.GetRootMessage() == nil {
		r.SetRootMessage(readWriteModel.NewBVLCReadBroadcastDistributionTable())
	}
	r.bvlciFunction = BVLCIReadBroadcastDistributionTable
	r.bvlciLength = 4
	return r, nil
}

func (r *ReadBroadcastDistributionTable) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	return nil
}

func (r *ReadBroadcastDistributionTable) Decode(bvlpdu Arg) error {
	if err := r._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCReadBroadcastDistributionTable:
			r.SetRootMessage(rm)
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		r.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (r *ReadBroadcastDistributionTable) String() string {
	if r == nil {
		return "(*ReadBroadcastDistributionTable)(nil)"
	}
	return fmt.Sprintf("ReadBroadcastDistributionTable{%v}", r._BVLPDU)
}
