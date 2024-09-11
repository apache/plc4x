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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// APCISequenceContract provides a set of functions which can be overwritten by a sub struct
type APCISequenceContract interface {
	SequenceContractRequirement
}

// APCISequenceContractRequirement is needed when one want to extend using SequenceContract
type APCISequenceContractRequirement interface {
	APCISequenceContract
	// SetAPCISequence callback is needed as we work in the constructor already with the finished object // TODO: maybe we need to return as init again as it might not be finished constructing....
	SetAPCISequence(a *APCISequence)
}

type APCISequence struct {
	*_APCI
	*Sequence

	_contract APCISequenceContract

	tagList *TagList
}

func NewAPCISequence(args Args, kwargs KWArgs, opts ...func(*APCISequence)) (*APCISequence, error) {
	a := &APCISequence{}
	for _, opt := range opts {
		opt(a)
	}
	if a._contract == nil {
		a._contract = a
	} else {
		a._contract.(APCISequenceContractRequirement).SetAPCISequence(a)
	}
	apdu := GA[readWriteModel.APDU](args, 0) // TODO: might break but should be fine as we are a GA*
	a._APCI = NewAPCI(apdu).(*_APCI)
	var err error
	a.Sequence, err = NewSequence(args, kwargs, WithSequenceExtension(a._contract))
	if err != nil {
		return nil, errors.Wrap(err, "error creating sequence")
	}

	// start with an empty tag list
	a.tagList = NewTagList(nil)
	return a, nil
}

func WithAPCISequenceExtension(contract APCISequenceContractRequirement) func(*APCISequence) {
	return func(a *APCISequence) {
		a._contract = contract
	}
}

func (a *APCISequence) SetSequence(sequence *Sequence) {
	a.Sequence = sequence
}

func (a *APCISequence) Encode(apdu Arg) error {
	switch apdu := apdu.(type) {
	case Updater:
		if err := apdu.Update(a); err != nil {
			return errors.Wrap(err, "error updating APDU")
		}
	}

	// create a tag list
	a.tagList = NewTagList(nil)

	if err := a.Sequence.Encode(a.tagList); err != nil {
		return errors.Wrap(err, "error encoding TagList")
	}

	switch apdu := apdu.(type) {
	case PDUData:
		// encode the tag list
		a.tagList.Encode(apdu)
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

func (a *APCISequence) Decode(apdu Arg) error {
	// copy the header fields
	if err := a.Update(apdu); err != nil {
		return errors.Wrap(err, "error updating APDU")
	}

	// create a tag list and decode the rest of the data
	a.tagList = NewTagList(nil)
	switch apdu := apdu.(type) {
	case PDUData:
		if err := a.tagList.Decode(apdu); err != nil {
			return errors.Wrap(err, "error decoding TagList")
		}
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	// pass the taglist to the Sequence for additional decoding
	if err := a.Sequence.Decode(a.tagList); err != nil {
		return errors.Wrap(err, "error encoding TagList")
	}

	if len(a.tagList.GetTagList()) > 0 {
		return errors.New("trailing unmatched tags")
	}
	return nil
}
