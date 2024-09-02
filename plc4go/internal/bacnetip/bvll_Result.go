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

type Result struct {
	*_BVLPDU

	bvlciResultCode model.BVLCResultCode
}

var _ BVLPDU = (*Result)(nil)

func NewResult(opts ...func(result *Result)) (*Result, error) {
	b := &Result{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(model.NewBVLCResult(b.bvlciResultCode)).(*_BVLPDU)
	return b, nil
}

func WithResultBvlciResultCode(code model.BVLCResultCode) func(*Result) {
	return func(b *Result) {
		b.bvlciResultCode = code
	}
}

func (n *Result) GetBvlciResultCode() model.BVLCResultCode {
	return n.bvlciResultCode
}

func (n *Result) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.PutShort(uint16(n.bvlciResultCode))
		bvlpdu.setBVLC(n.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *Result) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := n.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCResult:
			switch bvlc := rm.(type) {
			case model.BVLCResult:
				n.setBVLC(bvlc)
				n.bvlciResultCode = bvlc.GetCode()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *Result) String() string {
	return fmt.Sprintf("Result{%v, bvlciResultCode: %v}", n._BVLPDU, n.bvlciResultCode)
}
