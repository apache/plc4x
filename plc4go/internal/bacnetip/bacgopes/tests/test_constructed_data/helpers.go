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

package test_constructed_data

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type SequenceEqualityRequirements interface {
	GetSequenceElements() []Element
}

type SequenceEquality struct {
	_requirements SequenceEqualityRequirements
}

func NewSequenceEquality(requirements SequenceEqualityRequirements) *SequenceEquality {
	if requirements == nil {
		panic("requirements cannot be nil")
	}
	return &SequenceEquality{_requirements: requirements}
}

func (s *SequenceEquality) Equals(other any) bool {
	for _, element := range s._requirements.GetSequenceElements() {
		if !element.IsOptional() && true {
			panic("what??")
		}
	}
	return true
}

type EmptySequence struct {
	*Sequence
	*SequenceEquality
}

func NewEmptySequence(kwArgs KWArgs) (*EmptySequence, error) {
	e := &EmptySequence{}
	var err error
	e.Sequence, err = NewSequence(NoArgs, kwArgs, WithSequenceExtension(e))
	if err != nil {
		return nil, errors.Wrap(err, "could not create sequence")
	}
	e.SequenceEquality = NewSequenceEquality(e)
	return e, nil
}

func (e *EmptySequence) SetSequence(sequence *Sequence) {
	e.Sequence = sequence
}

type SimpleSequence struct {
	*Sequence
	*SequenceEquality

	sequenceElements []Element
}

func NewSimpleSequence(kwArgs KWArgs) (*SimpleSequence, error) {
	s := &SimpleSequence{
		sequenceElements: []Element{
			NewElement("hydrogen", V2E(NewBoolean)),
		},
	}
	var err error
	s.Sequence, err = NewSequence(NoArgs, kwArgs, WithSequenceExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "could not create sequence")
	}
	s.SequenceEquality = NewSequenceEquality(s)
	return s, nil
}

func (e *SimpleSequence) SetSequence(sequence *Sequence) {
	e.Sequence = sequence
}

func (e *SimpleSequence) GetSequenceElements() []Element {
	return e.sequenceElements
}

type CompoundSequence1 struct {
	*Sequence
	*SequenceEquality

	sequenceElements []Element
}

func NewCompoundSequence1(kwArgs KWArgs) (*CompoundSequence1, error) {
	s := &CompoundSequence1{
		sequenceElements: []Element{
			NewElement("hydrogen", V2E(NewBoolean)),
			NewElement("helium", V2E(NewInteger)),
		},
	}
	var err error
	s.Sequence, err = NewSequence(NoArgs, kwArgs, WithSequenceExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "could not create sequence")
	}
	s.SequenceEquality = NewSequenceEquality(s)
	return s, nil
}

func (e *CompoundSequence1) SetSequence(sequence *Sequence) {
	e.Sequence = sequence
}

func (e *CompoundSequence1) GetSequenceElements() []Element {
	return e.sequenceElements
}

type CompoundSequence2 struct {
	*Sequence
	*SequenceEquality

	sequenceElements []Element
}

func NewCompoundSequence2(kwArgs KWArgs) (*CompoundSequence2, error) {
	s := &CompoundSequence2{
		sequenceElements: []Element{
			NewElement("lithium", V2E(NewBoolean), WithElementOptional(true)),
			NewElement("beryllium", V2E(NewInteger)),
		},
	}
	var err error
	s.Sequence, err = NewSequence(NoArgs, kwArgs, WithSequenceExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "could not create sequence")
	}
	s.SequenceEquality = NewSequenceEquality(s)
	return s, nil
}

func (e *CompoundSequence2) SetSequence(sequence *Sequence) {
	e.Sequence = sequence
}

func (e *CompoundSequence2) GetSequenceElements() []Element {
	return e.sequenceElements
}
