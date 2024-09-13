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

func NewConfirmedRequestSequence(serviceRequest /*TODO: breaks a bit the consistency, maybe we just convert it to args to be flexible*/ model.BACnetConfirmedServiceRequest, kwArgs KWArgs, opts ...func(*ConfirmedRequestSequence)) (*ConfirmedRequestSequence, error) {
	u := &ConfirmedRequestSequence{}
	for _, opt := range opts {
		opt(u)
	}
	if u._contract == nil {
		u._contract = u
	} else {
		u._contract.(ConfirmedRequestSequenceContractRequirement).SetConfirmedRequestSequence(u)
	}
	var err error
	u.APCISequence, err = NewAPCISequence(NA(model.NewAPDUConfirmedRequest(false, false, false, model.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS, model.MaxApduLengthAccepted_NUM_OCTETS_1476, 0, nil, nil, serviceRequest, nil, nil, 0)), kwArgs, WithAPCISequenceExtension(u._contract))
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APCISequence")
	}
	u.ConfirmedRequestPDU, err = NewConfirmedRequestPDU(serviceRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error creating ConfirmedRequestPDU")
	}
	// We need to set the APCI to the same objects...
	u.APCISequence._APCI = u.ConfirmedRequestPDU._APCI
	return u, nil
}

func WithConfirmedRequestSequenceExtension(contract ConfirmedRequestSequenceContractRequirement) func(*ConfirmedRequestSequence) {
	return func(a *ConfirmedRequestSequence) {
		a._contract = contract
	}
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
