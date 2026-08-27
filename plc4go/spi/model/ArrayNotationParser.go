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

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

// The array notation shared by every PLC4X tag address, in both language bindings.
//
//	array-expression = group , { group } ;
//	group            = "[" , dimension , { "," , dimension } , "]" ;
//	dimension        = bounds , [ ";" , base ] ;
//	bounds           = index , [ ".." , index ] ;
//
// A range is inclusive of both bounds, so [0..7] is eight elements. A bare index is one element,
// so [4] is the fifth. ";base" states the array's declared lower bound - [4..7;1] selects
// elements 4 to 7 of an array declared from 1, which sit at offsets 3 to 6 - and defaults to 0.
// Each bracket group is one dimension, in written order, and several dimensions may be written
// comma-separated inside one bracket, which is the spelling Allen-Bradley and others use.
//
// This file is the single definition of that grammar in plc4go. The specification is shared with
// plc4j - see specs/002-unified-array-notation - but the code is not: any difference in what is
// accepted, or in what it means, is a defect rather than a Go dialect.
//
// Patterns are compiled once: address parsing runs per tag per request.

const (
	// dimensionPattern is one dimension: an index or an inclusive range, with an optional
	// declared lower bound.
	dimensionPattern = `\d+(?:\.\.\d+)?(?:;\d+)?`

	// ArrayExpressionPattern is the array expression as a regex fragment with no capturing
	// groups. A driver embeds this in its own address pattern so that the grammar has one
	// definition rather than a copy per driver - see ArrayGroupPattern.
	ArrayExpressionPattern = `(?:\[` + dimensionPattern + `(?:,` + dimensionPattern + `)*])+`

	// ArrayGroupPattern is the array expression as an optional named group, ready to splice into
	// a driver's address pattern between the address and the type. The group is named "array".
	ArrayGroupPattern = `(?P<array>` + ArrayExpressionPattern + `)?`
)

var (
	// expressionPattern matches a trailing run of strictly numeric bracket groups.
	expressionPattern = regexp.MustCompile(ArrayExpressionPattern + `$`)

	// wholeExpressionPattern matches an expression in its entirety.
	wholeExpressionPattern = regexp.MustCompile(`^` + ArrayExpressionPattern + `$`)

	// groupPattern matches one bracket group, whose content is one or more comma-separated
	// dimensions.
	groupPattern = regexp.MustCompile(`\[([^\]]*)]`)

	// singleDimensionPattern matches one dimension of a group.
	singleDimensionPattern = regexp.MustCompile(`^(\d+)(?:\.\.(\d+))?(?:;(\d+))?$`)

	// legacyAfterType matches an address in the pre-migration shape "address:TYPE[n]", where the
	// selection came after the type and meant a count.
	legacyAfterType = regexp.MustCompile(`^(.+?):([A-Za-z_][A-Za-z_0-9]*(?:\(\d+\))?)\[(\d+)]$`)

	// legacyCountSuffix matches the pre-migration shape "address:TYPE:n", where a count trailed.
	legacyCountSuffix = regexp.MustCompile(`^(.+?):([A-Za-z_][A-Za-z_0-9]*):(\d+)$`)
)

// AddressPart is the part of an address before any trailing array expression. An address with no
// such expression is returned unchanged - including one whose brackets are not numeric, such as
// an OPC UA string identifier that happens to contain them.
func AddressPart(address string) string {
	if loc := expressionPattern.FindStringIndex(address); loc != nil {
		return address[:loc[0]]
	}
	return address
}

// ExpressionPart is the trailing array expression of an address, or the empty string if it has
// none.
func ExpressionPart(address string) string {
	if match := expressionPattern.FindString(address); match != "" {
		return match
	}
	return ""
}

// SelectsSingleElement reports whether an expression selects a single element rather than an
// array, which decides what a tag reports from GetArrayInfo: a bare index yields a scalar, a
// range yields an array even when it spans one element. [1] is a scalar and [1..1] is an array
// of one, so equal bounds alone cannot tell them apart - only the written form can.
//
// An expression is a single element when every one of its dimensions is a bare index.
func SelectsSingleElement(expression string) bool {
	return expression != "" && !strings.Contains(expression, "..")
}

// ParseArrayExpression parses an array expression into one ArrayInfo per dimension, in written
// order. An empty expression selects nothing and yields no dimensions.
//
// The address is quoted back in any error so the caller can find it.
func ParseArrayExpression(expression string, address string, constraints AddressConstraints) ([]apiModel.ArrayInfo, error) {
	if expression == "" {
		return nil, nil
	}
	if !wholeExpressionPattern.MatchString(expression) {
		return nil, fmt.Errorf("invalid array expression '%s' in tag '%s': expected [index], "+
			"[lo..hi] or either with a ';base', repeated once per dimension", expression, address)
	}

	var dimensions []apiModel.ArrayInfo
	for _, group := range groupPattern.FindAllStringSubmatch(expression, -1) {
		// A group may hold several dimensions separated by commas - the spelling Allen-Bradley
		// and others use. "[1..2,3..4]" and "[1..2][3..4]" are the same selection.
		for _, part := range strings.Split(group[1], ",") {
			dimension, err := parseDimension(part, address, constraints)
			if err != nil {
				return nil, err
			}
			dimensions = append(dimensions, dimension)
		}
	}

	if len(dimensions) > constraints.MaxDimensions {
		return nil, fmt.Errorf("array expression '%s' in tag '%s' has %d dimensions, but this "+
			"protocol carries at most %d", expression, address, len(dimensions), constraints.MaxDimensions)
	}
	if constraints.OnlyTrailingDimensionMayBeRange {
		for i := 0; i < len(dimensions)-1; i++ {
			if dimensions[i].GetSize() > 1 {
				return nil, fmt.Errorf("array expression '%s' in tag '%s' spans %d elements in "+
					"dimension %d, but this protocol carries one element count for the whole "+
					"address, so only the last dimension may be a range",
					expression, address, dimensions[i].GetSize(), i+1)
			}
		}
	}
	return dimensions, nil
}

func parseDimension(part string, address string, constraints AddressConstraints) (apiModel.ArrayInfo, error) {
	match := singleDimensionPattern.FindStringSubmatch(part)
	if match == nil {
		return nil, fmt.Errorf("invalid array dimension '%s' in tag '%s'", part, address)
	}

	lowerBound, err := parseIndex(match[1], part, address)
	if err != nil {
		return nil, err
	}
	isRange := match[2] != ""
	upperBound := lowerBound
	if isRange {
		if upperBound, err = parseIndex(match[2], part, address); err != nil {
			return nil, err
		}
	}
	var base uint32
	if match[3] != "" {
		if base, err = parseIndex(match[3], part, address); err != nil {
			return nil, err
		}
	}

	if upperBound < lowerBound {
		return nil, fmt.Errorf("invalid array range '%s' in tag '%s': the upper bound %d is "+
			"below the lower bound %d", part, address, upperBound, lowerBound)
	}
	if lowerBound < base {
		return nil, fmt.Errorf("invalid array range '%s' in tag '%s': index %d lies below the "+
			"declared lower bound %d", part, address, lowerBound, base)
	}
	// The bound applies to the offset the protocol actually encodes - the start of the selection
	// - not to its last element. A CIP request carries a start index and an element count, so
	// [0..300] is encodable where [300] is not.
	if lowerBound-base > constraints.MaxIndex {
		return nil, fmt.Errorf("invalid array range '%s' in tag '%s': index %d is out of range "+
			"0 to %d for this protocol", part, address, lowerBound-base, constraints.MaxIndex)
	}

	return &DefaultArrayInfo{
		LowerBound: lowerBound,
		UpperBound: upperBound,
		Base:       base,
		Range:      isRange,
	}, nil
}

func parseIndex(value string, part string, address string) (uint32, error) {
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid array range '%s' in tag '%s': '%s' is not a number this "+
			"protocol can address", part, address, value)
	}
	return uint32(parsed), nil
}

// RenderArrayExpression renders dimensions back to their canonical form: one bracket per
// dimension, omitting what is defaulted - a base of 0 is dropped, and a bare index stays bare.
// The comma-separated spelling is accepted on input but never produced, so [1..2,3..4] renders as
// [1..2][3..4]. A one-element range still renders as a range: [8..8] is an array of one and [8]
// is a scalar, so collapsing it would change what the address means.
func RenderArrayExpression(dimensions []apiModel.ArrayInfo) string {
	if len(dimensions) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, dimension := range dimensions {
		sb.WriteString("[")
		sb.WriteString(strconv.FormatUint(uint64(dimension.GetLowerBound()), 10))
		if dimension.IsRange() {
			sb.WriteString("..")
			sb.WriteString(strconv.FormatUint(uint64(dimension.GetUpperBound()), 10))
		}
		if dimension.GetBase() != 0 {
			sb.WriteString(";")
			sb.WriteString(strconv.FormatUint(uint64(dimension.GetBase()), 10))
		}
		sb.WriteString("]")
	}
	return sb.String()
}

// CurrentFormOf returns how to rewrite an address written before the array notation was unified,
// and whether one could be worked out.
//
// The brackets moved from after the type to before it, and a count became a range. An upgrading
// user who sees only "does not match pattern" has to work that out from a regex; this hands them
// the address they meant.
func CurrentFormOf(address string) (string, bool) {
	if match := legacyAfterType.FindStringSubmatch(address); match != nil {
		return match[1] + rangeFor(match[3]) + ":" + match[2], true
	}
	if match := legacyCountSuffix.FindStringSubmatch(address); match != nil {
		return match[1] + rangeFor(match[3]) + ":" + match[2], true
	}
	return "", false
}

func rangeFor(count string) string {
	elements, err := strconv.Atoi(count)
	if err != nil || elements <= 1 {
		return "[0]"
	}
	return "[0.." + strconv.Itoa(elements-1) + "]"
}

// InvalidAddressError reports an address the driver could not parse, naming the form it expected
// and - when the address looks like one written before the notation was unified - the address to
// write instead.
func InvalidAddressError(address string, expectedForm string) error {
	message := fmt.Sprintf("invalid address '%s': expected %s", address, expectedForm)
	if current, ok := CurrentFormOf(address); ok {
		message += fmt.Sprintf(". The array notation moved before the type and a count became a "+
			"range, so this address is now written '%s'", current)
	}
	return fmt.Errorf("%s", message)
}
