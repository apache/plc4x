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

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// TODO: implement it...
type APCISequence struct {
	*_APCI
	*Sequence

	tagList *TagList
}

func NewAPCISequence() (*APCISequence, error) {
	a := &APCISequence{}
	a._APCI = NewAPCI(nil).(*_APCI) // TODO: what to pass up?
	var err error
	a.Sequence, err = NewSequence(NoKWArgs, WithSequenceContract(a))
	if err != nil {
		return nil, errors.Wrap(err, "error creating sequence")
	}

	// start with an empty tag list
	a.tagList = NewTagList(nil)
	return a, nil
}

func (a *APCISequence) SetSequence(sequence *Sequence) {
	a.Sequence = sequence
}

func (a *APCISequence) Encode(apdu Arg) error {
	switch apdu := apdu.(type) {
	case APDU:
		if err := apdu.Update(a); err != nil {
			return errors.Wrap(err, "error updating APDU")
		}

		// create a tag list
		a.tagList = NewTagList(nil)
		if err := a.Sequence.Encode(a.tagList); err != nil {
			return errors.Wrap(err, "error encoding TagList")
		}

		// encode the tag list
		a.tagList.Encode(apdu)
		return nil
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
}

func (a *APCISequence) Decode(apdu Arg) error {
	switch apdu := apdu.(type) {
	case APDU:
		if err := a.Update(apdu); err != nil {
			return errors.Wrap(err, "error updating APDU")
		}
		switch pduUserData := apdu.GetRootMessage().(type) {
		case model.APDUExactly:
			a.tagList = NewTagList(nil)
			if err := a.tagList.Decode(apdu); err != nil {
				return errors.Wrap(err, "error decoding TagList")
			}
			// pass the taglist to the Sequence for additional decoding
			if err := a.Sequence.Decode(a.tagList); err != nil {
				return errors.Wrap(err, "error encoding TagList")
			}

			_ = pduUserData
		}
		return nil
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
}
