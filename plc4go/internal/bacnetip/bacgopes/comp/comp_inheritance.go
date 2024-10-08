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
	"fmt"

	"github.com/pkg/errors"
)

var inheritanceDebug = func(_ string, _ ...interface{}) {}

type Initializer interface {
	Init(Args, KWArgs) error
}

// GenericConstructor is the generic signature for constructors
type GenericConstructor[T any] func(args Args, kwArgs KWArgs, options ...Option) (*T, error)

type LeafNameSupplier interface {
	GetLeafName() string
}

type leafType struct {
	defaultOption

	name string
}

// AddLeafTypeIfAbundant can be used to add a leaf type if not already set
func AddLeafTypeIfAbundant[T any](options []Option, leaf *T) []Option {
	for _, option := range options {
		switch option.(type) {
		case leafType:
			// We have a match so nothing to do
			return options
		}
	}
	return append(options, WithLeafType(leaf))
}

// WithLeafType can be used to add a leaf type
func WithLeafType[T any](leaf *T) Option {
	return leafType{name: "bacgopes." + fmt.Sprintf("%T", leaf)[1:]}
}

// ExtractLeafName or return default
func ExtractLeafName(options []Option, defaultName string) string {
	for _, option := range options {
		switch option := option.(type) {
		case leafType:
			return option.name
		}
	}
	return defaultName
}

// sharedHolder holds a common super
type sharedHolder[T any] struct {
	defaultOption
	commonSuper *T
}

// sharedHolders is a collection of different holders (useful if you don't want to define every super in the tree)
type sharedHolders struct {
	defaultOption
	holders map[any]struct{}
}

// AddMultiInheritanceSupportIfAbundant adds a WithMultiInheritanceSupport if abundant
func AddMultiInheritanceSupportIfAbundant(options []Option) []Option {
	for _, option := range options {
		switch option.(type) {
		case *sharedHolders:
			inheritanceDebug("*sharedHolders found on add. address: %p", option)
			// We have a match so nothing to do
			return options
		}
	}
	option := WithMultiInheritanceSupport()
	inheritanceDebug("add new *sharedHolders on add. address: %p", option)
	return append(options, option)
}

// WithMultiInheritanceSupport is a quick way to ensure that multi inheritance is respected from a root node.
func WithMultiInheritanceSupport() Option {
	return &sharedHolders{holders: map[any]struct{}{}}
}

// AddSharedSuperIfAbundant adds a WithSharedSuper if abundant
func AddSharedSuperIfAbundant[T any](options []Option) []Option {
	for _, option := range options {
		switch option := option.(type) {
		case *sharedHolder[T]:
			inheritanceDebug("*sharedHolder found on add. address: %p", option)
			// We have a match so nothing to do
			return options
		case *sharedHolders:
			for option := range option.holders {
				switch option := option.(type) {
				case *sharedHolder[T]:
					inheritanceDebug("*sharedHolder found in holders on add. address: %p", option)
					// We have a match so nothing to do
					return options
				}
			}
		}
	}
	option := WithSharedSuper[T]()
	inheritanceDebug("add new *sharedHolder on add. address: %p", option)
	return append(options, option)
}

// WithSharedSuper adds a shared super holder. Useful in multi-inheritance
func WithSharedSuper[T any]() Option {
	return &sharedHolder[T]{}
}

// CreateSharedSuperIfAbundant is the actual factory/retrieval function used when the actual object is about to be build
func CreateSharedSuperIfAbundant[T any](options []Option, constructor GenericConstructor[T], cArgs Args, ckwArgs KWArgs, cOptions ...Option) (*T, error) {
	var holders *sharedHolders
	// look for shared holders containers so we can add it when we find one
sharedHoldersScan:
	for _, option := range options {
		switch option := option.(type) {
		case *sharedHolders:
			holders = option
			inheritanceDebug("*sharedHolders found. address: %p", option)
			// check if we have an instance in there
			for holder := range option.holders {
				if option, ok := holder.(*sharedHolder[T]); ok {
					inheritanceDebug("found *sharedHolder[T]. address: %p", option)
					for _, existingOption := range options {
						if existingOption == option {
							inheritanceDebug("already there. address: %p", option)
							break sharedHoldersScan // we are good, no need to add
						}
					}
					inheritanceDebug("appending it. address: %p", option)
					options = append(options, option) // Append it to the options and let the code below do its magic
					break sharedHoldersScan           // we are good, no need to add
				}
			}
			// We haven't found on so we just add an empty one
			holder := new(sharedHolder[T])
			inheritanceDebug("adding a new holder. address: %p", holder)
			options = append(options, holder)
			break sharedHoldersScan // we are good, no need to add
		}
	}
	// look for shared holders
	for _, option := range options {
		switch option := option.(type) {
		case *sharedHolder[T]:
			inheritanceDebug("*sharedHolder found. address: %p", option)
			if holders != nil { // Add it to a holder instance if there are some
				inheritanceDebug("adding it to holders. address: %p", holders)
				holders.holders[option] = struct{}{}
			}
			commonSuper := option.commonSuper
			if commonSuper == nil {
				// Apparently it is not initialized so we do that now
				var err error
				commonSuper, err = constructor(cArgs, ckwArgs, cOptions...)
				if err != nil {
					return nil, errors.Wrap(err, "error creating object")
				}
				inheritanceDebug("new common super. address: %p (%[1]T)", commonSuper)
				option.commonSuper = commonSuper
			} else {
				inheritanceDebug("existing common super. address: %p (%[1]T)", commonSuper)
			}
			return commonSuper, nil
		}
	}
	inheritanceDebug("No holder, just fallback")
	// we don't share anything, so just create one
	commonSuper, err := constructor(cArgs, ckwArgs, cOptions...)
	inheritanceDebug("new common super fallback. address: %p (%[1]T)", commonSuper)
	return commonSuper, err
}
