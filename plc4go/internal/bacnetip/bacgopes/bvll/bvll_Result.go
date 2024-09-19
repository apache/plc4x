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

type Result struct {
	*_BVLPDU

	bvlciResultCode readWriteModel.BVLCResultCode
}

var _ BVLPDU = (*Result)(nil)

// TODO: check this arg desaster...
func NewResult(code *readWriteModel.BVLCResultCode, args Args, kwArgs KWArgs, options ...Option) (*Result, error) {
	r := &Result{}
	options = AddLeafTypeIfAbundant(options, r)
	r._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	if r.GetRootMessage() == nil {
		r.SetRootMessage(readWriteModel.NewBVLCResult(r.bvlciResultCode))
	}
	r.AddDebugContents(r, "bvlciResultCode")
	r.bvlciFunction = BVLCIResult
	r.bvlciLength = 6
	if code != nil {
		r.bvlciResultCode = *code
	}
	return r, nil
}

func (r *Result) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciResultCode":
		return r.bvlciResultCode
	default:
		return nil
	}
}

func (r *Result) GetBvlciResultCode() readWriteModel.BVLCResultCode {
	return r.bvlciResultCode
}

func (r *Result) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		bvlpdu.PutShort(uint16(r.bvlciResultCode))
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (r *Result) Decode(bvlpdu Arg) error {
	if err := r._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCResult:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCResult:
				r.bvlciResultCode = bvlc.GetCode()
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

func (r *Result) String() string {
	if r == nil {
		return "(*Result)(nil)"
	}
	return fmt.Sprintf("Result{%v, bvlciResultCode: %v}", r._BVLPDU, r.bvlciResultCode)
}
