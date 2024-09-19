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

package apdu

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// ErrorSequenceContract provides a set of functions which can be overwritten by a sub struct
type ErrorSequenceContract interface {
	APCISequenceContractRequirement
	GetErrorChoice() *model.BACnetConfirmedServiceChoice
}

// ErrorSequenceContractRequirement is needed when one want to extend using SequenceContract
type ErrorSequenceContractRequirement interface {
	ErrorSequenceContract
	// SetErrorSequence callback is needed as we work in the constructor already with the finished object // TODO: maybe we need to return as init again as it might not be finished constructing....
	SetErrorSequence(es *ErrorSequence)
}

type ErrorSequence struct {
	*APCISequence
	*ErrorPDU

	_contract ErrorSequenceContract
}

func NewErrorSequence(args Args, kwArgs KWArgs, options ...Option) (*ErrorSequence, error) {
	e := &ErrorSequence{}
	ApplyAppliers(options, e)
	if e._contract == nil {
		e._contract = e
	} else {
		e._contract.(ErrorSequenceContractRequirement).SetErrorSequence(e)
	}
	options = AddLeafTypeIfAbundant(options, e)
	var err error
	kwArgs[KWConfirmedServiceChoice] = e._contract.GetErrorChoice()
	e.ErrorPDU, err = NewErrorPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating ErrorPDU")
	}
	e.APCISequence, err = NewAPCISequence(args, kwArgs, Combine(options, WithAPCISequenceExtension(e._contract))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APCISequence")
	}
	// We need to set the APCI to the same objects...
	e.APCISequence._APCI = e.ErrorPDU._APCI
	return e, nil
}

func WithErrorSequenceExtension(contract ErrorSequenceContractRequirement) GenericApplier[*ErrorSequence] {
	return WrapGenericApplier(func(a *ErrorSequence) { a._contract = contract })
}

func (u *ErrorSequence) SetAPCISequence(a *APCISequence) {
	u.APCISequence = a
}

func (*ErrorSequence) GetErrorChoice() *model.BACnetConfirmedServiceChoice {
	return nil
}

func (u *ErrorSequence) DeepCopy() any {
	panic("implement me")
}

func (u *ErrorSequence) String() string {
	return u.APCISequence.String()
}
