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

// UnconfirmedRequestSequenceContract provides a set of functions which can be overwritten by a sub struct
type UnconfirmedRequestSequenceContract interface {
	APCISequenceContractRequirement
	GetServiceChoice() *model.BACnetUnconfirmedServiceChoice
}

// UnconfirmedRequestSequenceContractRequirement is needed when one want to extend using SequenceContract
type UnconfirmedRequestSequenceContractRequirement interface {
	UnconfirmedRequestSequenceContract
	// SetUnconfirmedRequestSequence callback is needed as we work in the constructor already with the finished object // TODO: maybe we need to return as init again as it might not be finished constructing....
	SetUnconfirmedRequestSequence(urs *UnconfirmedRequestSequence)
}

type UnconfirmedRequestSequence struct {
	*APCISequence
	*UnconfirmedRequestPDU

	_contract UnconfirmedRequestSequenceContract
}

func NewUnconfirmedRequestSequence(args Args, kwArgs KWArgs, options ...Option) (*UnconfirmedRequestSequence, error) {
	u := &UnconfirmedRequestSequence{}
	ApplyAppliers(options, u)
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	if u._contract == nil {
		u._contract = u
	} else {
		u._contract.(UnconfirmedRequestSequenceContractRequirement).SetUnconfirmedRequestSequence(u)
	}
	options = AddSharedSuperIfAbundant[_APCI](options)
	options = AddLeafTypeIfAbundant(options, u)
	var err error
	kwArgs[KWUnconfirmedServiceChoice] = u._contract.GetServiceChoice()
	u.UnconfirmedRequestPDU, err = NewUnconfirmedRequestPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating UnconfirmedRequestPDU")
	}
	//TODO: the sequence is usually init first but seems upstream does a init on the same level first before going deeper
	u.APCISequence, err = NewAPCISequence(args, kwArgs, Combine(options, WithAPCISequenceExtension(u._contract))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APCISequence")
	}
	if u.GetRootMessage() == nil {
		panic("this should be set by NewConfirmedRequestPDU")
		serviceRequest, _ := GAO[model.BACnetConfirmedServiceRequest](args, 0, nil)
		if serviceRequest != nil {
			apduConfirmedRequest := model.NewAPDUConfirmedRequest(false, false, false, model.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS, model.MaxApduLengthAccepted_NUM_OCTETS_1476, 0, nil, nil, serviceRequest, nil, nil, 0)
			u.SetRootMessage(apduConfirmedRequest)
		}
	}
	return u, nil
}

func WithUnconfirmedRequestSequenceExtension(contract UnconfirmedRequestSequenceContractRequirement) GenericApplier[*UnconfirmedRequestSequence] {
	return WrapGenericApplier(func(a *UnconfirmedRequestSequence) { a._contract = contract })
}

func (u *UnconfirmedRequestSequence) SetAPCISequence(a *APCISequence) {
	u.APCISequence = a
}

func (*UnconfirmedRequestSequence) GetServiceChoice() *model.BACnetUnconfirmedServiceChoice {
	return nil
}

func (u *UnconfirmedRequestSequence) DeepCopy() any {
	panic("implement me")
}

func (u *UnconfirmedRequestSequence) String() string {
	return u.APCISequence.String()
}
