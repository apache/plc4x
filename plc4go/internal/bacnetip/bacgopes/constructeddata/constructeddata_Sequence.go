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
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

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

type Sequence struct {
	_contract        SequenceContract
	sequenceElements []Element
	attr             map[string]any
}

// NewSequence Create a sequence element, optionally providing attribute/property values.
func NewSequence(args Args, kwArgs KWArgs, options ...Option) (*Sequence, error) {
	s := &Sequence{
		attr: make(map[string]any),
	}
	ApplyAppliers(options, s)
	if s._contract == nil {
		s._contract = s
	} else {
		s._contract.(SequenceContractRequirement).SetSequence(s)
	}
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}

	var myKWArgs = make(KWArgs)
	var otherKWArgs = make(KWArgs)
	for _, element := range s._contract.GetSequenceElements() {
		if a, ok := kwArgs[KnownKey(element.GetName())]; ok {
			myKWArgs[KnownKey(element.GetName())] = a
		}
	}
	for key, a := range kwArgs {
		if _, ok := myKWArgs[key]; !ok {
			otherKWArgs[key] = a
		}
	}

	if _debug != nil {
		_debug("    - my_kwargs: %r", myKWArgs)
	}
	if _debug != nil {
		_debug("    - other_kwargs: %r", otherKWArgs)
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

func WithSequenceExtension(contract SequenceContractRequirement) GenericApplier[*Sequence] {
	return WrapGenericApplier(func(s *Sequence) { s._contract = contract })
}

func (a *Sequence) GetSequenceElements() []Element {
	return a.sequenceElements
}

func (a *Sequence) GetAttr(key string) (any, bool) {
	v, ok := a.attr[key]
	return v, ok
}

func (a *Sequence) Encode(arg Arg) error {
	tagList, ok := arg.(*TagList)
	if !ok {
		return errors.New("arg is not a TagList")
	}
	if _debug != nil {
		_debug("encode %r", tagList)
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
			return errors.Wrap(err, "can't get zero object")
		}
		_, elementInSequenceOfClasses := _sequenceOfClasses[fmt.Sprintf("%T", elementKlass)]
		_, elementInListOfClasses := _listOfClasses[fmt.Sprintf("%T", elementKlass)]
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

			helper, err := element.GetKlass()(NA(value), NoKWArgs())
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
			helper, err := element.GetKlass()(NA(value), NoKWArgs())
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
	if _debug != nil {
		_debug("decode %r", tagList)
	}

	for _, element := range a._contract.GetSequenceElements() {
		tag := tagList.Peek()
		if _debug != nil {
			_debug("    - element, tag: %r, %r", element, tag)
		}

		elementKlass, err := element.GetKlass()(Nothing())
		if err != nil {
			return errors.Wrap(err, "can't get zero object")
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
			helper, err := element.GetKlass()(NA(tag), NoKWArgs())
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

func (a *Sequence) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	for _, element := range a.sequenceElements {
		value, ok := a.GetAttr(element.GetName())
		if element.IsOptional() && !ok {
			continue
		} else if !element.IsOptional() && !ok {
			_, _ = fmt.Fprintf(file, "%s%s is a missing required element of %s\n", strings.Repeat("    ", indent), element.GetName(), StructName())
			continue
		}

		elementKlass, err := element.GetKlass()(Nothing())
		if err != nil {
			if _debug != nil {
				_debug("can't get zero object. %w", err)
			}
			return
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
		isValueSubOfElementKlass := false // TODO: how to figure that out??
		if elementInSequenceOfClasses || elementInListOfClasses {
			_, _ = fmt.Fprintf(file, "%s%s\n", strings.Repeat("    ", indent), element.GetName())
			helper, err := element.GetKlass()(NA(value), NoKWArgs())
			if err != nil {
				if _debug != nil {
					_debug("    - helper class %s, err: %v", element.GetName(), err)
				}
				return
			}
			helper.(DebugContentPrinter).PrintDebugContents(indent+1, file, _ids)
		} else if isAtomic || isAnyAtomic {
			printVerb := VerbForType(value, 'r')
			_, _ = fmt.Fprintf(file, "%s%s = "+string(printVerb)+"\n", strings.Repeat("    ", indent), element.GetName(), value)
		} else if isValueSubOfElementKlass {
			value.(DebugContentPrinter).PrintDebugContents(indent+1, file, _ids)
		} else {
			_, _ = fmt.Fprintf(file, "%s%s must be a %T\n", strings.Repeat("    ", indent), element.GetName(), element.GetKlass())
		}
	}
}
