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

package comp

import (
	"cmp"
	"fmt"
	"iter"
	"reflect"
	"sort"
)

// NillableKey is a key which can be used in maps
type NillableKey[T any] struct {
	Value T
	IsNil bool
}

func (n NillableKey[T]) Format(s fmt.State, v rune) {
	switch v {
	case 'v', 's', 'r':
		_, _ = fmt.Fprint(s, n.String())
	}
}

func (n NillableKey[T]) String() string {
	if n.IsNil {
		return "nil"
	}
	return fmt.Sprintf("%v", n.Value)
}

// NK creates a new NillableKey of type K
func NK[T any, K NillableKey[T]](value *T) K {
	var _nk NillableKey[T]
	if value == nil {
		_nk.IsNil = true
		return K(_nk)
	}
	_nk.Value = *value
	return K(_nk)
}

// DeepCopy copies things implementing Copyable
func DeepCopy[T Copyable](copyable Copyable) T {
	return copyable.DeepCopy().(T)
}

// CopyPtr copies things that are a pointer to something
func CopyPtr[T any](t *T) *T {
	if t == nil {
		return nil
	}
	tc := *t
	return &tc
}

// SortedMapIterator lets you iterate over an array in a deterministic way
func SortedMapIterator[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// OR returns a or b
func OR[T comparable](a T, b T) T {
	if reflect.ValueOf(a).IsNil() || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil()) { // TODO: check if there is another way than using reflect
		return b
	} else {
		return a
	}
}

// ToPtr gives a Ptr
func ToPtr[T any](value T) *T {
	return &value
}

// Try something and return panic as error
func Try(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return f()
}

// Try1 something and return panic as error
func Try1[T any](f func() (T, error)) (v T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return f()
}

// IsNil when nil checks aren't enough
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	valueOf := reflect.ValueOf(v)
	switch valueOf.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return valueOf.IsNil()
	default:
		return false
	}
}

func ToStringers[I any](in []I) []fmt.Stringer {
	return ConvertSlice[I, fmt.Stringer](in)
}

func ConvertSlice[I any, O any](in []I) (out []O) {
	out = make([]O, len(in))
	for i, v := range in {
		out[i] = any(v).(O)
	}
	return
}
