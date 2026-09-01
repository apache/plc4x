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

import apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"

var _ apiModel.ArrayInfo = &DefaultArrayInfo{}

//go:generate go tool plc4xGenerator -type=DefaultArrayInfo
type DefaultArrayInfo struct {
	LowerBound uint32
	UpperBound uint32
	// Base is the array's declared lower bound; 0 for an array that does not declare one.
	Base uint32
	// Range records whether the address wrote this dimension as a range. A one-element range is
	// still a range, so this cannot be derived from the bounds - see apiModel.ArrayInfo.IsRange.
	Range bool
}

// GetSize is the number of elements. Both bounds are inclusive, so {0, 7} is eight elements.
// This used to return UpperBound-LowerBound, treating the upper bound as exclusive, which
// disagreed with plc4j about the same address and with the drivers that build from an inclusive
// range.
func (t *DefaultArrayInfo) GetSize() uint32 {
	return t.UpperBound - t.LowerBound + 1
}

func (t *DefaultArrayInfo) GetLowerBound() uint32 {
	return t.LowerBound
}

func (t *DefaultArrayInfo) GetUpperBound() uint32 {
	return t.UpperBound
}

func (t *DefaultArrayInfo) GetBase() uint32 {
	return t.Base
}

func (t *DefaultArrayInfo) IsRange() bool {
	return t.Range
}
