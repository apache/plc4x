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

package basetypes

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/errors"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type NameValue struct {
	*Sequence

	sequenceElements []Element
	name             string
	value            any
}

func NewNameValue(args Args) (*NameValue, error) {
	s := &NameValue{
		sequenceElements: []Element{
			NewElement("name", V2E(NewCharacterString)),
			NewElement("value", Vs2E(NewAnyAtomic), WithElementOptional(true)),
		},
	}
	// default to no value
	s.name, _ = GAO[string](args, 0, "")
	s.value, _ = GAO[any](args, 1, nil)

	if s.value == nil {
		return s, nil
	}
	switch value := s.value.(type) {
	case IsAtomic:
		s.value = value
	case DateTime:
		s.value = value
	case Tag:
		var err error
		s.value, err = value.AppToObject()
		if err != nil {
			return nil, errors.Wrap(err, "error converting tag")
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", value)
	}
	return s, nil
}

func (s *NameValue) Encode(arg Arg) error {
	tagList := arg.(*TagList)
	// build a tag and encode the name into it
	tag, err := NewTag(NoArgs)
	if err != nil {
		return errors.Wrap(err, "error converting tag")
	}
	characterString, err := NewCharacterString(s.name)
	if err != nil {
		return errors.Wrap(err, "error creating character string")
	}
	if err := characterString.Encode(tag); err != nil {
		return errors.Wrap(err, "error encoding")
	}
	context, err := tag.AppToContext(0)
	if err != nil {
		return errors.Wrap(err, "error converting tag")
	}
	tagList.Append(context)

	// the value is optional
	if s.value != nil {
		if v, ok := s.value.(*DateTime); ok {
			// has its own encoder
			if err := v.Encode(tagList); err != nil {
				return errors.Wrap(err, "error converting tag")
			}
		} else if e, ok := s.value.(Encoder); ok {
			// atomic values encode into a tag
			tag, err = NewTag(NoArgs)
			if err != nil {
				return errors.Wrap(err, "error creating tag")
			}
			if err := e.Encode(tag); err != nil {
				return errors.Wrap(err, "error converting tag")
			}
			tagList.Append(tag)
		}
	}
	return nil
}

func (s *NameValue) Decode(arg Arg) error {
	tagList := arg.(*TagList)

	// no contents yet
	s.name = ""
	s.value = ""

	// look for the context encoded character string
	tag := tagList.Peek()
	if tag == nil || (tag.GetTagClass() != readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS) || (tag.GetTagNumber() != 0) {
		return MissingRequiredParameter{RejectException: RejectException{Exception: Exception{Message: fmt.Sprintf("%s is a missing required element of %p", s.name, s)}}}
	}

	// pop it off and save the value
	tagList.Pop()
	tag, err := tag.ContextToApp(uint(readWriteModel.BACnetDataType_CHARACTER_STRING))
	if err != nil {
		return errors.Wrap(err, "error converting tag")
	}
	characterString, err := NewCharacterString(tag)
	if err != nil {
		return errors.Wrap(err, "error converting tag to string")
	}
	s.name = characterString.GetValue()

	// look for the optional application encoded value
	tag = tagList.Peek()
	if tag != nil && (tag.GetTagClass() == readWriteModel.TagClass_APPLICATION_TAGS) {
		// if it is a date check the next one for a time
		if tag.GetTagNumber() == uint(readWriteModel.BACnetDataType_DATE) && (len(tagList.GetTagList()) >= 2) {
			nextTag := tagList.GetTagList()[1]

			if nextTag.GetTagClass() == readWriteModel.TagClass_APPLICATION_TAGS && (nextTag.GetTagNumber() == uint(readWriteModel.BACnetDataType_TIME)) {
				s.value, err = NewDateTime(NoArgs)
				if err != nil {
					return errors.Wrap(err, "error creating date time")
				}
				if err := s.value.(Decoder).Decode(tagList); err != nil {
					return errors.Wrap(err, "error decoding taglist")
				}
			}

			// just a primitive value
			if s.value == nil {
				tagList.Pop()
				s.value, err = tag.AppToObject()
				if err != nil {
					return errors.Wrap(err, "error converting tag")
				}
			}
		}
	}
	return nil
}
