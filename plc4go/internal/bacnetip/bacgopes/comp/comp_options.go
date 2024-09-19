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

// OptionalOption allows options to be applied that might be optional
func OptionalOption[V any, T any](value *V, opt func(V) GenericApplier[*T]) GenericApplier[*T] {
	if value != nil {
		return opt(*value)
	}
	return WrapGenericApplier(func(c *T) {})
}

// OptionalOption2 allows options to be applied that might be optional
func OptionalOption2[V1 any, V2 any, T any](value1 *V1, value2 *V2, opt func(V1, V2) GenericApplier[*T]) GenericApplier[*T] {
	v1Set := value1 != nil
	v2Set := value2 != nil
	if (v1Set && !v2Set) || (!v1Set && v2Set) {
		return WrapGenericApplier(func(c *T) {})
	}
	if v1Set {
		return opt(*value1, *value2)
	}
	return WrapGenericApplier(func(c *T) {})
}

// Option is a generic interface for transporting options which are meant to bubble up
type Option interface {
	isOption()
}

type defaultOption struct {
}

func (defaultOption) isOption() {}

// Combine Option s with different Option s
func Combine(options []Option, additionalOptions ...Option) []Option {
	return append(options, additionalOptions...)
}

type genericOption[T any] struct {
	defaultOption
	value T
}

// AsOption converts a value T as option
func AsOption[T any](value T) Option {
	return genericOption[T]{value: value}
}

// GenericApplier wraps a func(T) as option
type GenericApplier[T any] genericOption[func(T)]

// WrapGenericApplier is a factory function for GenericApplier
func WrapGenericApplier[T any](applier func(T)) GenericApplier[T] {
	return GenericApplier[T]{value: applier}
}

// ApplyAppliers applies all appliers to target if genericOption matches GenericApplier
func ApplyAppliers[T any](options []Option, target T) {
	for _, option := range options {
		switch option := option.(type) {
		case GenericApplier[T]:
			option.value(target)
		}
	}
}

// ApplyGenericOption applies to applier if genericOption matches T
func ApplyGenericOption[T any](options []Option, applier func(T)) {
	for _, option := range options {
		switch option := option.(type) {
		case genericOption[T]:
			applier(option.value)
		}
	}
}

// AddIfAbundant adds a generic Option for T if value is abundant
func AddIfAbundant[T any](options []Option, value T) []Option {
	for _, option := range options {
		switch option.(type) {
		case genericOption[T]:
			// We have a match so nothing to do
			return options
		}
	}
	return append(options, genericOption[T]{value: value})
}

// ExtractIfPresent extracts T if present
func ExtractIfPresent[T any](options []Option) (T, bool) {
	for _, option := range options {
		switch option := option.(type) {
		case genericOption[T]:
			return option.value, true
		}
	}
	return *new(T), false
}
