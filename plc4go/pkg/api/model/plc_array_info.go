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

import "fmt"

type ArrayInfo interface {
	fmt.Stringer
	// GetSize is the number of elements. Both bounds are inclusive, so [0..7] is eight.
	GetSize() uint32
	// GetLowerBound is the lower index as the address wrote it, not the resolved offset.
	GetLowerBound() uint32
	// GetUpperBound is the upper index as the address wrote it. Inclusive.
	GetUpperBound() uint32
	// GetBase is the array's declared lower bound, as in PLCs not every array starts at 0. An
	// address may state it - [4..7;1] selects elements 4 to 7 of an array declared from 1 - so
	// that the bounds above can be written the way the PLC program declares them. The offset of
	// an element from the start of the array is its index minus this value. Defaults to 0.
	GetBase() uint32
	// IsRange reports whether the address wrote this dimension as a range rather than a single
	// index. The two mean different things to a caller: a single index selects one element and
	// yields a scalar, while a range yields an array - even a range spanning one element. Equal
	// bounds alone cannot tell them apart, so the written form has to be remembered.
	IsRange() bool
}
