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
	"github.com/pkg/errors"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// TODO: implement it...
type ComplexAckSequence struct {
	*APCISequence
	*ComplexAckPDU
}

// UnconfirmedRequestSequenceContract provides a set of functions which can be overwritten by a sub struct
type UnconfirmedRequestSequenceContract interface {
	APCISequenceContractRequirement
	GetServiceChoice() *readWriteModel.BACnetUnconfirmedServiceChoice
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

func NewUnconfirmedRequestSequence(opts ...func(*UnconfirmedRequestSequence)) (*UnconfirmedRequestSequence, error) {
	u := &UnconfirmedRequestSequence{}
	for _, opt := range opts {
		opt(u)
	}
	if u._contract == nil {
		u._contract = u
	} else {
		u._contract.(UnconfirmedRequestSequenceContractRequirement).SetUnconfirmedRequestSequence(u)
	}
	var err error
	u.APCISequence, err = NewAPCISequence(WithAPCISequenceExtension(u._contract))
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APCISequence")
	}
	u.UnconfirmedRequestPDU, err = NewUnconfirmedRequestPDU(func(pdu *UnconfirmedRequestPDU) {
		pdu.argChoice = u.GetServiceChoice()
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating UnconfirmedRequestPDU")
	}
	return u, nil
}

func WithUnconfirmedRequestSequenceExtension(contract UnconfirmedRequestSequenceContractRequirement) func(*UnconfirmedRequestSequence) {
	return func(a *UnconfirmedRequestSequence) {
		a._contract = contract
	}
}

func (u *UnconfirmedRequestSequence) GetServiceChoice() *readWriteModel.BACnetUnconfirmedServiceChoice {
	return nil
}

func (u *UnconfirmedRequestSequence) SetAPCISequence(a *APCISequence) {
	u.APCISequence = a
}
