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

// ConfirmedRequestSequenceContract provides a set of functions which can be overwritten by a sub struct
type ConfirmedRequestSequenceContract interface {
	APCISequenceContractRequirement
	GetServiceChoice() *model.BACnetConfirmedServiceChoice
}

// ConfirmedRequestSequenceContractRequirement is needed when one want to extend using SequenceContract
type ConfirmedRequestSequenceContractRequirement interface {
	ConfirmedRequestSequenceContract
	// SetConfirmedRequestSequence callback is needed as we work in the constructor already with the finished object // TODO: maybe we need to return as init again as it might not be finished constructing....
	SetConfirmedRequestSequence(crs *ConfirmedRequestSequence)
}

type ConfirmedRequestSequence struct {
	*APCISequence
	*ConfirmedRequestPDU

	_contract ConfirmedRequestSequenceContract
}

func NewConfirmedRequestSequence(args Args, kwArgs KWArgs, options ...Option) (*ConfirmedRequestSequence, error) {
	u := &ConfirmedRequestSequence{}
	ApplyAppliers(options, u)
	if u._contract == nil {
		u._contract = u
	} else {
		u._contract.(ConfirmedRequestSequenceContractRequirement).SetConfirmedRequestSequence(u)
	}
	options = AddSharedSuperIfAbundant[_APCI](options)
	options = AddLeafTypeIfAbundant(options, u)
	kwArgs[KWConfirmedServiceChoice] = u._contract.GetServiceChoice()
	var err error
	u.ConfirmedRequestPDU, err = NewConfirmedRequestPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating ConfirmedRequestPDU")
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

func WithConfirmedRequestSequenceExtension(contract ConfirmedRequestSequenceContractRequirement) GenericApplier[*ConfirmedRequestSequence] {
	return WrapGenericApplier(func(a *ConfirmedRequestSequence) { a._contract = contract })
}

func (u *ConfirmedRequestSequence) SetAPCISequence(a *APCISequence) {
	u.APCISequence = a
}

func (*ConfirmedRequestSequence) GetServiceChoice() *model.BACnetConfirmedServiceChoice {
	return nil
}

func (u *ConfirmedRequestSequence) DeepCopy() any {
	panic("implement me")
}

func (u *ConfirmedRequestSequence) String() string {
	return u.APCISequence.String()
}
