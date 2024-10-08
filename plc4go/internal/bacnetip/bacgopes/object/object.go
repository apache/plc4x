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

var _arrayOfMap = make(map[any]func(Args, KWArgs) (PropertyKlass, error))
var _listOfMap = make(map[any]func(Args, KWArgs) (PropertyKlass, error))
var _arrayOfClasses = make(map[any]struct{})
var _sequenceOfClasses = make(map[any]struct{})
var _listOfClasses = make(map[any]struct{})

var RegisteredObjectTypes = make(map[any]struct{})

func RegisterObjectType(kwargs KWArgs) (any, error) {
	cls, aCls := KWO[any](kwargs, KWCls, nil)
	vendorId, _ := KWO[int](kwargs, KWVendorIdentifier, 0)

	if _debug != nil {
		_debug("register_object_type %r vendor_id=%s", cls, vendorId)
	}

	if !aCls {
		panic("implement me")
	}

	// make sure it's an Object derived class
	objClass, ok := cls.(Object)
	if !ok {
		return nil, errors.New("Object derived struct required")
	}

	// build a property dictionary by going through the class and all its parents
	_properties := map[string]Property{}
	for _, property := range objClass.GetProperties() {
		_properties[property.GetIdentifier()] = property
	}

	// if the object type hasn't been provided, make an immutable one
	if _, ok := _properties["objectType"]; !ok {
		_properties["objectType"] = NewReadableProperty("objectType", func(args Args, _ KWArgs) (PropertyKlass, error) { return NewObjectType(args) }, WithPropertyMutable(false))
	}

	// store it in the struct
	cls.(interface {
		Set_Properties(_properties map[string]Property)
	}).Set_Properties(_properties)

	RegisteredObjectTypes[cls] = struct{}{} // TODO: key needs to be objectType + vendor id

	// return the struct as decorator
	return cls, nil
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

var _ PropertyKlass = (*_ArrayOf[any])(nil)

func (a *_ArrayOf[T]) Encode(arg Arg) (O error) {
	panic("implement me")
}

func ArrayOfP[T any](klass func(arg Arg) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (PropertyKlass, error) {
	return ArrayOfPs(func(args Args) (*T, error) {
		if len(args) == 0 {
			return klass(NoArgs)
		}
		if len(args) > 1 {
			panic(fmt.Sprintf("oh no %d", len(args)))
		}
		return klass(args[0])
	}, fixedLength, prototype)
}

func ArrayOfPs[T any](klass func(args Args) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (PropertyKlass, error) {
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

// TODO: big WIP
type _ListOf[T any] struct {
	name    string
	subType func(Args) (*T, error)
}

func (a _ListOf[T]) Encode(arg Arg) (O error) {
	panic("implement me")
}

func ListOfP[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	return ListOfPs(func(args Args) (*T, error) {
		if len(args) == 0 {
			return b(NoArgs)
		}
		if len(args) > 1 {
			panic(fmt.Sprintf("oh no %d", len(args)))
		}
		return b(args[0])
	})
}

func ListOfPs[T any](klass func(args Args) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
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

	// if this has already been built, return the cached one
	if v, ok := _listOfMap[fmt.Sprintf("%T", elementKlass)]; ok {
		return v
	}

	// no ArrayOf(ArrayOf(...)) allowed
	if _, ok := _listOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (PropertyKlass, error) {
			return nil, errors.New("nested lists disallowed")
		}
	}
	// no ArrayOf(SequenceOf(...)) allowed
	if _, ok := _arrayOfClasses[fmt.Sprintf("%T", elementKlass)]; ok {
		return func(_ Args, _ KWArgs) (PropertyKlass, error) {
			return nil, errors.New("list of arrays disallowed")
		}
	}

	// define a generic class for lists
	_listOf := struct {
		*_ListOf[T]
	}{&_ListOf[T]{}}

	// constrain it to a list of a specific type of item
	_listOf.subType = klass

	// update the name
	_listOf.name = "ListOf" + fmt.Sprintf("%T", elementKlass)

	// cache this type
	_listOfMap[fmt.Sprintf("%T", elementKlass)] = func(_ Args, _ KWArgs) (PropertyKlass, error) {
		return _listOf, nil
	}
	_listOfClasses[fmt.Sprintf("%T", _listOf)] = struct{}{}

	return func(args Args, kwArgs KWArgs) (PropertyKlass, error) {
		return _listOf, nil
	}
}
