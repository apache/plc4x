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

package constructeddata

import (
	"reflect"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

var _sequenceOfClasses map[any]struct{}
var _listOfClasses map[any]struct{}

func init() {
	_sequenceOfClasses = make(map[any]struct{})
	_listOfClasses = make(map[any]struct{})
}

type Element interface {
	GetName() string
	GetKlass() func(Args, KWArgs) (ElementKlass, error)
	GetContext() *int
	IsOptional() bool
	Encode(tagList Arg) error
}

type ElementKlass interface {
	Encode(Arg) error
	GetAppTag() readWriteModel.BACnetDataType
}

// V2E accepts a function which takes an Arg and maps it to a ElementKlass
func V2E[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	return func(args Args, kwargs KWArgs) (ElementKlass, error) {
		var arg any
		if len(args) == 1 {
			arg = args[0]
		}
		r, err := b(arg)
		return any(r).(ElementKlass), err
	}
}

// Vs2E accepts a function which takes an Args and maps it to a ElementKlass
func Vs2E[T any](b func(args Args) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	return func(args Args, kwargs KWArgs) (ElementKlass, error) {
		r, err := b(args)
		return any(r).(ElementKlass), err
	}
}

// TODO: finish
type _Element struct {
	Name     string
	Klass    func(Args, KWArgs) (ElementKlass, error)
	Context  *int
	Optional bool
}

func NewElement(name string, klass func(Args, KWArgs) (ElementKlass, error), opts ...func(*_Element)) Element {
	e := &_Element{
		Name:  name,
		Klass: klass,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

var _ Element = (*_Element)(nil)

func WithElementOptional(optional bool) func(*_Element) {
	return func(e *_Element) {
		e.Optional = optional
	}
}

func WithElementContext(context int) func(*_Element) {
	return func(e *_Element) {
		e.Context = &context
	}
}

func (e *_Element) GetName() string {
	return e.Name
}

func (e *_Element) GetKlass() func(Args, KWArgs) (ElementKlass, error) {
	return e.Klass
}

func (e *_Element) GetContext() *int {
	return e.Context
}

func (e *_Element) IsOptional() bool {
	return e.Optional
}

func (e *_Element) Encode(tagList Arg) error {
	//TODO implement me
	panic("implement me")
}

// SequenceContract provides a set of functions which can be overwritten by a sub struct
type SequenceContract interface {
	GetSequenceElements() []Element
}

// SequenceContractRequirement is needed when one want to extend using SequenceContract
type SequenceContractRequirement interface {
	SequenceContract
	// SetSequence callback is needed as we work in the constructor already with the finished object // TODO: maybe we need to return as init again as it might not be finished constructing....
	SetSequence(s *Sequence)
}

// TODO: finish
type Sequence struct {
	_contract        SequenceContract
	sequenceElements []Element
	attr             map[string]any
}

// NewSequence Create a sequence element, optionally providing attribute/property values.
func NewSequence(args Args, kwargs KWArgs, opts ...func(*Sequence)) (*Sequence, error) {
	s := &Sequence{
		attr: make(map[string]any),
	}
	for _, opt := range opts {
		opt(s)
	}
	if s._contract == nil {
		s._contract = s
	} else {
		s._contract.(SequenceContractRequirement).SetSequence(s)
	}

	var myKWArgs = make(KWArgs)
	var otherKWArgs = make(KWArgs)
	for _, element := range s._contract.GetSequenceElements() {
		if a, ok := kwargs[KnownKey(element.GetName())]; ok {
			myKWArgs[KnownKey(element.GetName())] = a
		}
	}
	for key, a := range kwargs {
		if _, ok := myKWArgs[key]; !ok {
			otherKWArgs[key] = a
		}
	}

	if len(otherKWArgs) > 0 {
		return nil, errors.Errorf("invalid arguments %v", otherKWArgs)
	}

	// set the attribute/property values for the ones provided
	for _, element := range s._contract.GetSequenceElements() {
		a, ok := myKWArgs[KnownKey(element.GetName())]
		if ok {
			s.attr[element.GetName()] = a
		}
	}
	return s, nil
}

func WithSequenceExtension(contract SequenceContractRequirement) func(*Sequence) {
	return func(s *Sequence) {
		s._contract = contract
	}
}

func (a *Sequence) GetSequenceElements() []Element {
	return a.sequenceElements
}

func (a *Sequence) Encode(arg Arg) error {
	tagList, ok := arg.(*TagList)
	if !ok {
		return errors.New("arg is not a TagList")
	}

	for _, element := range a._contract.GetSequenceElements() {
		value, ok := a.attr[element.GetName()]
		if element.IsOptional() && !ok {
			continue
		}
		if !element.IsOptional() && !ok {
			return errors.Errorf("%s is a missing required element of %T", element.GetName(), a)
		}
		elementKlass, err := element.GetKlass()(Nothing())
		if err != nil {
			return errors.New("can't get zero object")
		}
		_, elementInSequenceOfClasses := _sequenceOfClasses[elementKlass]
		_, elementInListOfClasses := _listOfClasses[elementKlass]
		isAtomic := false
		switch elementKlass.(type) {
		case IsAtomic, IsAnyAtomic:
			isAtomic = true
		}
		isValue := reflect.TypeOf(value) == reflect.TypeOf(elementKlass)

		if elementInSequenceOfClasses || elementInListOfClasses {
			// might need to encode an opening tag
			if element.GetContext() != nil {
				openingTag, err := NewOpeningTag(*element.GetContext())
				if err != nil {
					return errors.Wrap(err, "error creating opening tag")
				}
				tagList.Append(openingTag)
			}

			helper, err := element.GetKlass()(NewArgs(value), NoKWArgs)
			if err != nil {
				return errors.Wrap(err, "error klass element")
			}

			// encode the value
			if err := helper.Encode(tagList); err != nil {
				return errors.Wrap(err, "error encoding tag list")
			}

			// might need to encode a closing tag
			if element.GetContext() != nil {
				closingTag, err := NewClosingTag(*element.GetContext())
				if err != nil {
					return errors.Wrap(err, "error creating closing tag")
				}
				tagList.Append(closingTag)
			}
		} else if isAtomic {
			helper, err := element.GetKlass()(NewArgs(value), NoKWArgs)
			if err != nil {
				return errors.Wrap(err, "error klass element")
			}

			// build a tag and encode the data into it
			tag, err := NewTag(nil)
			if err != nil {
				return errors.Wrap(err, "error creating tag")
			}
			// encode the value
			if err := helper.Encode(tag); err != nil {
				return errors.Wrap(err, "error encoding tag list")
			}

			// convert it to context encoding if necessary
			if element.GetContext() != nil {
				tag, err = tag.AppToContext(uint(*element.GetContext()))
				if err != nil {
					return errors.Wrap(err, "error converting tag to context")
				}
			}
			tagList.Append(tag)
		} else if isValue {
			// might need to encode an opening tag
			if element.GetContext() != nil {
				openingTag, err := NewOpeningTag(*element.GetContext())
				if err != nil {
					return errors.Wrap(err, "error creating opening tag")
				}
				tagList.Append(openingTag)
			}

			// encode the tag
			if err := value.(interface{ Encode(Arg) error }).Encode(tagList); err != nil { // TODO: ugly case, need a encode interface soon
				return errors.Wrap(err, "error encoding tag list")
			}

			// might need to encode a closing tag
			if element.GetContext() != nil {
				closingTag, err := NewClosingTag(*element.GetContext())
				if err != nil {
					return errors.Wrap(err, "error creating closing tag")
				}
				tagList.Append(closingTag)
			}
		}
	}
	return nil
}

func (a *Sequence) Decode(arg Arg) error {
	tagList, ok := arg.(*TagList)
	if !ok {
		return errors.New("arg is not a TagList")
	}

	for _, element := range a._contract.GetSequenceElements() {
		tag := tagList.Peek()

		elementKlass, err := element.GetKlass()(Nothing())
		if err != nil {
			return errors.New("can't get zero object")
		}
		_, elementInSequenceOfClasses := _sequenceOfClasses[elementKlass]
		_, elementInListOfClasses := _listOfClasses[elementKlass]
		isAtomic := false
		isAnyAtomic := false
		switch elementKlass.(type) {
		case IsAtomic:
			isAtomic = true
		case IsAnyAtomic:
			isAnyAtomic = true
		}
		// no more elements
		if tag == nil {
			if element.IsOptional() {
				// ommited optional element
				a.attr[element.GetName()] = nil
			} else if elementInSequenceOfClasses || elementInListOfClasses {
				// empty list
				//a.attr[element.GetName()] = nil // TODO: what to do???
			} else {
				return errors.Errorf("%s is a missing required element of %T", element.GetName(), a)
			}
		} else if tag.GetTagClass() == TagClosingTagClass {
			if !element.IsOptional() {
				return errors.Errorf("%s is a missing required element of %T", element.GetName(), a)
			}

			// ommited optional element
			// a.attr[element.GetName()] = nil // TODO: don't set it for now as we use _,ok:=
		} else if elementInSequenceOfClasses {
			// check for context encoding
			panic("finish me") // TODO: finish me
		} else if isAnyAtomic {
			// convert it to application encoding
			panic("finish me") // TODO: finish me
		} else if isAtomic {
			// convert it to application encoding
			if context := element.GetContext(); context != nil {
				if tag.GetTagClass() != readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS && tag.GetTagNumber() != uint(*context) {
					if !element.IsOptional() {
						return errors.Errorf("%s expected context tag %d", element.GetName(), *context)
					} else {
						// TODO: we don't do this
						//a.attr[element.GetName()] = nil
						continue
					}
				}
				atomicTag := tag.(interface {
					GetAppTag() readWriteModel.BACnetDataType
				})
				tag, err = tag.ContextToApp(uint(atomicTag.GetAppTag()))
				if err != nil {
					return errors.Wrap(err, "error converting tag")
				}
			} else {
				if tag.GetTagClass() != readWriteModel.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(elementKlass.GetAppTag()) {
					if !element.IsOptional() {
						return errors.Errorf("%s expected context tag %d", element.GetName(), context)
					} else {
						// TODO: we don't do this
						//a.attr[element.GetName()] = nil
						continue
					}
				}
			}

			// consume the tag
			tagList.Pop()

			// a helper cooperates between the atomic value and the tag
			helper, err := element.GetKlass()(NewArgs(tag), NoKWArgs)
			if err != nil {
				return errors.Wrap(err, "error creating helper")
			}
			a.attr[element.GetName()] = helper
		} else if isAnyAtomic { // TODO: what is upstream doing here??? how???
			// convert it to application encoding
			panic("finish me") // TODO: finish me
		} else {
			panic("finish me") // TODO: finish me
		}
	}
	return nil
}

// TODO: finish
func SequenceOf[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	panic("finish me")
}

func SequenceOfs[T any](b func(args Args) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	panic("finish me")
}

// TODO: finish // convert to kwargs and check wtf we are doing here...
func ArrayOf[T any](b func(arg Arg) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (ElementKlass, error) {
	panic("finish me")
}

// TODO: finish
type List struct {
}

// TODO: finish
type Array struct {
}

// TODO: finish
type Choice struct {
}

type AnyContract interface {
	Encode(taglist TagList) error
	Decode(taglist TagList) error
	castIn(arg Arg) error
	castOut(arg Arg) error
}

type Any struct {
	_contract AnyContract // TODO: finish

	tagList TagList
}

func NewAny(args Args) (*Any, error) {
	a := &Any{}

	// cast the args
	for _, arg := range args {
		if err := a.castIn(arg); err != nil {
			return nil, errors.Wrapf(err, "error casting arg %v", arg)
		}
	}
	return a, nil
}

func WithAnyContract(contract AnyContract) func(*Any) {
	return func(a *Any) {
		a._contract = contract
	}
}

func (a *Any) Encode(tagList TagList) error {
	a.tagList.Extend(tagList.GetTagList()...)
	return nil
}

func (a *Any) Decode(tagList TagList) error {
	lvl := 0
	for len(tagList.GetTagList()) != 0 {
		tag := tagList.Peek()
		if tag.GetTagClass() == TagOpeningTagClass {
			lvl++
		} else if tag.GetTagClass() == TagClosingTagClass {
			lvl--
			if lvl < 0 {
				break
			}
		}
		a.tagList.Append(tagList.Pop())
	}

	// make sure everything balances
	if lvl > 0 {
		return errors.New("mismatched open/close tags")
	}
	return nil
}

func (a *Any) castIn(element Arg) error {
	t := NewTagList(nil)
	switch element.(type) {
	case IsAtomic:
		tag, err := NewTag(nil)
		if err != nil {
			return errors.New("error creating empty tag")
		}
		if err := element.(interface{ Encode(arg Arg) error }).Encode(tag); err != nil {
			return errors.New("error encoding element")
		}
		t.Append(tag)
	case IsAnyAtomic:
		tag, err := NewTag(nil)
		if err != nil {
			return errors.New("error creating empty tag")
		}
		if err := element.(interface{ GetValue() any }).GetValue().(interface{ Encode(Tag) error }).Encode(tag); err != nil {
			return errors.New("error encoding element")
		}
		t.Append(tag)
	default:
		if err := element.(interface{ Encode(arg Arg) error }).Encode(t); err != nil {
			return errors.New("error encoding element")
		}
	}
	a.tagList.Extend(t.GetTagList()...)
	return nil
}

func (a *Any) castOut(element Arg) error {
	panic("implement me")
	return nil
}

type IsAnyAtomic interface {
	isAnyAtomic() bool
}

// TODO: finish me
type AnyAtomic struct {
	value any
}

var _ IsAnyAtomic = (*AnyAtomic)(nil)

func NewAnyAtomic(args Args) (*AnyAtomic, error) {
	a := &AnyAtomic{}
	return a, nil
}

func (a *AnyAtomic) isAnyAtomic() bool {
	return true
}

func (a *AnyAtomic) GetValue() any {
	return a.value
}

func (a *AnyAtomic) Encode(arg Arg) {}

func (a *AnyAtomic) Decode(arg Arg) error {
	return nil
}

// TODO: finish
type SequenceOfAny struct {
}
