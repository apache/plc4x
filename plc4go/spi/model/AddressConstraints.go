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

package model

import "math"

// AddressConstraints is what a protocol can actually encode of an array selection. The notation
// is the same for every driver, but the wire formats are not: an EtherNet/IP array index travels
// in a CIP MemberID whose instance field is a uint 8, and a driver addressing linear memory has
// no second dimension to express. A driver states its limits here, and a selection that exceeds
// them is reported when the address is parsed rather than truncated when it is serialized.
//
// This mirrors the plc4j type of the same name; the two bindings share a specification, not code.
type AddressConstraints struct {
	// MaxIndex is the largest start offset the protocol can encode. It bounds where a selection
	// begins, not where it ends: a protocol carrying a start index and an element count can read
	// past this bound, it just cannot start past it.
	MaxIndex uint32
	// MaxDimensions is how many dimensions the wire format carries.
	MaxDimensions int
	// OnlyTrailingDimensionMayBeRange reports whether every dimension but the last must be a
	// single element - true where one request carries a single element count for the whole
	// address, so that a range anywhere else would not be one contiguous read.
	OnlyTrailingDimensionMayBeRange bool
}

// Unconstrained imposes no limits beyond the grammar itself.
var Unconstrained = AddressConstraints{
	MaxIndex:      math.MaxUint32,
	MaxDimensions: math.MaxInt,
}

// SingleDimension is a protocol addressing linear memory: any index, but only one dimension.
var SingleDimension = AddressConstraints{
	MaxIndex:      math.MaxUint32,
	MaxDimensions: 1,
}

// WithMaxIndex returns a copy bounded to the given largest start offset.
func (c AddressConstraints) WithMaxIndex(maxIndex uint32) AddressConstraints {
	c.MaxIndex = maxIndex
	return c
}

// WithMaxDimensions returns a copy carrying at most the given number of dimensions.
func (c AddressConstraints) WithMaxDimensions(maxDimensions int) AddressConstraints {
	c.MaxDimensions = maxDimensions
	return c
}

// WithOnlyTrailingDimensionMayBeRange returns a copy in which only the last dimension may span
// more than one element.
func (c AddressConstraints) WithOnlyTrailingDimensionMayBeRange(onlyTrailing bool) AddressConstraints {
	c.OnlyTrailingDimensionMayBeRange = onlyTrailing
	return c
}
