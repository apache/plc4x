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

// TODO: big WIP
package object

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

var _debug = CreateDebugPrinter()

var _arrayOfMap map[any]func(Args, KWArgs) (PropertyKlass, error)
var _arrayOfClasses map[any]struct{}
var _sequenceOfClasses map[any]struct{}
var _listOfClasses map[any]struct{}

func init() {
	_arrayOfMap = make(map[any]func(Args, KWArgs) (PropertyKlass, error))
	_arrayOfClasses = make(map[any]struct{})
	_sequenceOfClasses = make(map[any]struct{})
	_listOfClasses = make(map[any]struct{})
}

var registeredObjectTypes map[any]struct{}

func init() {
	registeredObjectTypes = make(map[any]struct{})
}

type PropertyError struct {
	Key string
}

func (e PropertyError) Error() string {
	return fmt.Sprintf("Property '%s' is not a valid object", e.Key)
}

// V2P accepts a function which takes an Arg and maps it to a PropertyKlass
func V2P[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	return func(args Args, kwArgs KWArgs) (PropertyKlass, error) {
		var arg any
		if len(args) == 1 {
			arg = args[0]
		}
		r, err := b(arg)
		return any(r).(PropertyKlass), err
	}
}

// Vs2P accepts a function which takes an Args and maps it to a PropertyKlass
func Vs2P[T any](b func(args Args) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	return func(args Args, kwArgs KWArgs) (PropertyKlass, error) {
		r, err := b(args)
		return any(r).(PropertyKlass), err
	}
}

// TODO: finish
func SequenceOfP[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}

func SequenceOfsP[T any](b func(args Args) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}

// TODO: big WIP
type _ArrayOf[T any] struct {
	name        string
	subType     func(Args) (*T, error)
	fixedLength int
	prototype   any
}

func (a _ArrayOf[T]) Encode(arg Arg) (O error) {
	panic("implement me")
}

func ArrayOfP[T any](b func(arg Arg) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (PropertyKlass, error) {
	return ArrayOfsP(func(args Args) (*T, error) {
		if len(args) == 0 {
			return b(NoArgs)
		}
		if len(args) > 1 {
			panic(fmt.Sprintf("oh no %d", len(args)))
		}
		return b(args[0])
	}, fixedLength, prototype)
}

func ArrayOfsP[T any](klass func(args Args) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (PropertyKlass, error) {
	elementKlass, err := klass(NoArgs)
	if err != nil {
		return func(_ Args, _ KWArgs) (PropertyKlass, error) {
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
		return func(_ Args, _ KWArgs) (PropertyKlass, error) {
			return nil, errors.New("nested arrays disallowed")
		}
	}
	// no ArrayOf(SequenceOf(...)) allowed
	if _, ok := _sequenceOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (PropertyKlass, error) {
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
	_arrayOfMap[arraySignature] = func(_ Args, _ KWArgs) (PropertyKlass, error) {
		return _arrayOf, nil
	}
	_arrayOfClasses[fmt.Sprintf("%T", _arrayOf)] = struct{}{}

	return func(args Args, kwArgs KWArgs) (PropertyKlass, error) {
		return _arrayOf, nil
	}
}

// TODO: finish // convert to kwArgs and check wtf we are doing here...
func ListOfP[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}
