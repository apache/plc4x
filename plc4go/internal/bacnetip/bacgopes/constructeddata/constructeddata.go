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

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

var _debug = CreateDebugPrinter()

var _arrayOfMap = make(map[any]func(Args, KWArgs) (ElementKlass, error))
var _listOfMap = make(map[any]func(Args, KWArgs) (ElementKlass, error))
var _sequenceOfMap = make(map[any]func(Args, KWArgs) (ElementKlass, error))
var _sequenceOfClasses = make(map[any]struct{})
var _arrayOfClasses = make(map[any]struct{})
var _listOfClasses = make(map[any]struct{})

// V2E accepts a function which takes an Arg and maps it to a ElementKlass
func V2E[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	return func(args Args, kwArgs KWArgs) (ElementKlass, error) {
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
	return func(args Args, kwArgs KWArgs) (ElementKlass, error) {
		r, err := b(args)
		return any(r).(ElementKlass), err
	}
}

// TODO: big WIP
type _SequenceOf[T any] struct {
	name    string
	subType func(Args) (*T, error)
}

func (a _SequenceOf[T]) Encode(arg Arg) (O error) {
	panic("implement me")
}

func (a _SequenceOf[T]) GetAppTag() readWriteModel.BACnetDataType {
	panic("implement me")
}

func SequenceOfE[T any](klass func(arg Arg) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	return SequenceOfEs(func(args Args) (*T, error) {
		if len(args) == 0 {
			return klass(NoArgs)
		}
		if len(args) > 1 {
			panic(fmt.Sprintf("oh no %d", len(args)))
		}
		return klass(args[0])
	})
}

func SequenceOfEs[T any](klass func(args Args) (*T, error)) func(Args, KWArgs) (ElementKlass, error) {
	elementKlass, err := klass(NoArgs)
	if err != nil {
		return func(_ Args, _ KWArgs) (ElementKlass, error) {
			return nil, errors.Wrap(err, "can't get zero object")
		}
	}
	switch any(elementKlass).(type) {
	case IsAtomic:
	// TODO: add prototype check...
	default:
		// TODO: add check
	}

	// if this has already been built, return the cached one
	if v, ok := _sequenceOfMap[fmt.Sprintf("%T", elementKlass)]; ok {
		return v
	}

	// no SequenceOf(SequenceOf(...)) allowed
	if _, ok := _sequenceOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (ElementKlass, error) {
			return nil, errors.New("nested sequences disallowed")
		}
	}
	// no SequenceOf(Array(...)) allowed
	if _, ok := _arrayOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (ElementKlass, error) {
			return nil, errors.New("sequence of ArrayOf disallowed")
		}
	}

	// define a generic class for arrays
	_sequenceOf := struct {
		*_SequenceOf[T]
	}{&_SequenceOf[T]{}}

	// constrain it to a list of a specific type of item
	_sequenceOf.subType = klass

	// update the name
	_sequenceOf.name = "SequenceOf" + fmt.Sprintf("%T", elementKlass)

	// cache this type
	_sequenceOfMap[fmt.Sprintf("%T", elementKlass)] = func(_ Args, _ KWArgs) (ElementKlass, error) {
		return _sequenceOf, nil
	}
	_sequenceOfClasses[fmt.Sprintf("%T", _sequenceOf)] = struct{}{}

	return func(args Args, kwArgs KWArgs) (ElementKlass, error) {
		return _sequenceOf, nil
	}
}

// TODO: big WIP
type _ArrayOf[T any] struct {
	name        string
	subType     func(Args) (*T, error)
	fixedLength int
	prototype   any
}

var _ ElementKlass = (*_ArrayOf[any])(nil)

func (a *_ArrayOf[T]) GetAppTag() readWriteModel.BACnetDataType {
	//TODO implement me
	panic("implement me")
}

func (a *_ArrayOf[T]) Encode(arg Arg) (O error) {
	panic("implement me")
}

func ArrayOfE[T any](klass func(arg Arg) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (ElementKlass, error) {
	return ArrayOfEs(func(args Args) (*T, error) {
		if len(args) == 0 {
			return klass(NoArgs)
		}
		if len(args) > 1 {
			panic(fmt.Sprintf("oh no %d", len(args)))
		}
		return klass(args[0])
	}, fixedLength, prototype)
}

// TODO: finish // convert to kwArgs and check wtf we are doing here...
func ArrayOfEs[T any](klass func(args Args) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (ElementKlass, error) {
	elementKlass, err := klass(NoArgs)
	if err != nil {
		return func(_ Args, _ KWArgs) (ElementKlass, error) {
			return nil, errors.Wrap(err, "can't get zero object")
		}
	}
	switch any(elementKlass).(type) {
	case IsAtomic:
	// TODO: add prototype check...
	default:
		// TODO: add check
	}

	// build a signature of the parameters
	arraySignature := struct {
		string
		int
		any
	}{
		fmt.Sprintf("%T", elementKlass),
		fixedLength,
		prototype,
	}

	// if this has already been built, return the cached one
	if v, ok := _arrayOfMap[arraySignature]; ok {
		return v
	}

	// no ArrayOf(ArrayOf(...)) allowed
	if _, ok := _arrayOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (ElementKlass, error) {
			return nil, errors.New("nested arrays disallowed")
		}
	}
	// no ArrayOf(SequenceOf(...)) allowed
	if _, ok := _sequenceOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (ElementKlass, error) {
			return nil, errors.New("arrays of SequenceOf disallowed")
		}
	}

	// define a generic class for arrays
	_arrayOf := struct {
		*_ArrayOf[T]
	}{&_ArrayOf[T]{}}

	// constrain it to a list of a specific type of item
	_arrayOf.subType = klass
	_arrayOf.fixedLength = fixedLength
	_arrayOf.prototype = prototype

	// update the name
	_arrayOf.name = "ArrayOf" + fmt.Sprintf("%T", elementKlass)

	// cache this type
	_arrayOfMap[arraySignature] = func(_ Args, _ KWArgs) (ElementKlass, error) {
		return _arrayOf, nil
	}
	_arrayOfClasses[fmt.Sprintf("%T", _arrayOf)] = struct{}{}

	return func(args Args, kwArgs KWArgs) (ElementKlass, error) {
		return _arrayOf, nil
	}
}

// TODO: finish // convert to kwArgs and check wtf we are doing here...
func ListOfE[T any](b func(arg Arg) (*T, error)) func(Args) (*T, error) {
	panic("finish me")
}
