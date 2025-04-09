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

package utils

import (
	"reflect"
)

// InlineIf is basically an inline if like construct for golang
func InlineIf[T any](test bool, a func() T, b func() T) T {
	if test {
		return a()
	} else {
		return b()
	}
}

// ToPtr makes a pointer to T
func ToPtr[T any](v T) *T {
	return &v
}

// CopyPtr copies things that are a pointer to something
func CopyPtr[T any](t *T) *T {
	if t == nil {
		return nil
	}
	tc := *t
	return &tc
}

// DeepCopy copies things implementing Copyable
func DeepCopy[T Copyable](copyable Copyable) T {
	if IsNil(copyable) {
		var zero T
		return zero
	}
	return copyable.DeepCopy().(T)
}

// DeepCopySlice copies as slice into a new one. Note: if slice contains pointer use DeepCopySliceWithConverter in combination with ToPtr
func DeepCopySlice[I any, O any](in []I) (out []O) {
	out = make([]O, len(in))
	for i, v := range in {
		if copyable, ok := any(v).(Copyable); ok {
			out[i] = copyable.DeepCopy().(O)
		} else {
			out[i] = any(v).(O)
		}
	}
	return
}

// DeepCopySliceWithConverter copies as slice into a new one using converter for transformation
func DeepCopySliceWithConverter[I any, O any](in []I, converter func(I) O) (out []O) {
	out = make([]O, len(in))
	for i, v := range in {
		out[i] = converter(v)
	}
	return
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
