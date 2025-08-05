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

package primitivedata

import (
	"cmp"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ComparableAndOrdered interface {
	comparable
	cmp.Ordered
}

type IsAtomic interface {
	isAtomic() bool
}

// AtomicContract provides a set of functions which can be overwritten by a sub struct
type AtomicContract[T ComparableAndOrdered] interface {
	Compare(other any) int
	LowerThan(other any) bool
	Equals(other any) bool
	GetValue() T
}

// AtomicRequirements provides a set of functions which must be overwritten by a sub struct
type AtomicRequirements interface {
	IsValid(arg any) bool
}

// Atomic is an abstract struct
type Atomic[T ComparableAndOrdered] struct {
	AtomicContract[T]
	atomicRequirements AtomicRequirements

	_appTag readWriteModel.BACnetDataType

	value T
}

var _ IsAtomic = (*Atomic[int])(nil)
var _ AtomicContract[int] = (*Atomic[int])(nil)

func NewAtomic[T ComparableAndOrdered](subStruct interface {
	AtomicContract[T]
	AtomicRequirements
}) *Atomic[T] {
	return &Atomic[T]{
		AtomicContract:     subStruct,
		atomicRequirements: subStruct,
	}
}

func (a *Atomic[T]) isAtomic() bool {
	return true
}

func (a *Atomic[T]) GetAppTag() readWriteModel.BACnetDataType {
	return a._appTag
}

func (a *Atomic[T]) Compare(other any) int {
	otherValue := other.(AtomicContract[T]).GetValue()
	// now compare the values
	if a.value < otherValue {
		return -1
	} else if a.value > otherValue {
		return 1
	} else {
		return 0
	}
}

func (a *Atomic[T]) LowerThan(other any) bool {
	otherValue := other.(AtomicContract[T]).GetValue()
	// now compare the values
	return a.value < otherValue
}

func (a *Atomic[T]) Equals(other any) bool {
	otherValue := other.(AtomicContract[T]).GetValue()
	// now compare the values
	return a.value == otherValue
}

func (a *Atomic[T]) GetValue() T {
	return a.value
}

func (a *Atomic[T]) Coerce(arg any) T {
	return arg.(AtomicContract[T]).GetValue()
}
