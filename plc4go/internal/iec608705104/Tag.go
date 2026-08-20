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

package iec608705104

import (
	"strconv"
	"strings"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

const (
	// maxAsduAddress is the largest common address of ASDU the wire format can carry. The field is
	// two octets wide.
	maxAsduAddress = 0xFFFF
	// maxInformationObjectAddress is the largest information object address the wire format can
	// carry. The field is three octets wide.
	maxInformationObjectAddress = 0xFFFFFF
	// maxOctet is the largest value a single octet of a split address may take.
	maxOctet = 0xFF
	// wildcard is the address component which matches every value in its position.
	wildcard = "*"
)

// addressComponent is one number of a tag address. A component is either a concrete value or the
// wildcard, which matches whatever the station reports in that position.
type addressComponent struct {
	any   bool
	value uint32
}

// matches reports whether a concrete value satisfies this component.
func (c addressComponent) matches(value uint32) bool {
	return c.any || c.value == value
}

func (c addressComponent) String() string {
	if c.any {
		return wildcard
	}
	return strconv.FormatUint(uint64(c.value), 10)
}

// Tag is a parsed IEC 60870-5-104 tag address: a common address of ASDU plus an information object
// address, either of which may be spelled octet by octet and may use the wildcard in any position.
//
//	<asdu>/<ioa>
//	  asdu = N | N/N          the two octet common address, low octet first
//	  ioa  = N | N.N.N        the three octet information object address, low octet first
//
// Examples:
//
//	1/2          ASDU 1, information object 2
//	1/0/2        the same ASDU, spelled as its two octets (low 1, high 0)
//	1/2.0.0      the same information object, spelled as its three octets (low 2)
//	*/*          every point the station reports
//	1/*          every point of ASDU 1
//	*/5.*.0      every point whose IOA has 5 in its low octet and 0 in its high octet
//
// This is the syntax plc4j's Iec608705104Tag only ever carried as a commented-out (and never
// compiling) regex, recovered here and implemented for real: plc4j's Iec608705104TagHandler.parseTag
// ignores its argument entirely and answers every address with ASDU 0 / IOA 0.
//
// A Tag is also an apiModel.PlcSubscriptionTag, because subscribing is the only thing this protocol
// offers and the subscription request builder only accepts tags which are one.
type Tag struct {
	// asdu is one component when the common address was spelled as a single number and two
	// components (low, high) when it was spelled octet by octet.
	asdu []addressComponent
	// ioa is one component when the information object address was spelled as a single number and
	// three components (low, middle, high) when it was spelled octet by octet.
	ioa []addressComponent
}

var _ apiModel.PlcSubscriptionTag = Tag{}

// NewTag builds a tag matching exactly one point, which is what an address without a wildcard
// parses into and what an incoming information object is turned into for reporting.
func NewTag(asduAddress uint16, informationObjectAddress uint32) Tag {
	return Tag{
		asdu: []addressComponent{{value: uint32(asduAddress)}},
		ioa:  []addressComponent{{value: informationObjectAddress & maxInformationObjectAddress}},
	}
}

// Matches reports whether an information object the station reported belongs to this tag. Every
// component has to agree; a wildcard component agrees with anything.
//
// A split address is compared octet by octet, so 1/0/2 and 1/2 match the very same point - the
// spelling is a convenience, not a different address space.
func (t Tag) Matches(asduAddress uint16, informationObjectAddress uint32) bool {
	if !matchesComponents(t.asdu, uint32(asduAddress)) {
		return false
	}
	return matchesComponents(t.ioa, informationObjectAddress&maxInformationObjectAddress)
}

// matchesComponents compares a concrete address against its components, which are either a single
// whole-value component or a run of octets ordered from the least significant one up.
func matchesComponents(components []addressComponent, value uint32) bool {
	if len(components) == 1 {
		return components[0].matches(value)
	}
	for i, component := range components {
		if !component.matches((value >> (8 * i)) & maxOctet) {
			return false
		}
	}
	return true
}

// GetAsduAddress is the concrete common address this tag covers and whether there is one: a tag
// with a wildcard anywhere in its ASDU part covers many.
func (t Tag) GetAsduAddress() (uint16, bool) {
	value, ok := concreteValue(t.asdu)
	return uint16(value), ok
}

// GetInformationObjectAddress is the concrete information object address this tag covers and
// whether there is one.
func (t Tag) GetInformationObjectAddress() (uint32, bool) {
	return concreteValue(t.ioa)
}

// concreteValue reassembles the single address a run of components covers, reporting false as soon
// as one component is a wildcard.
func concreteValue(components []addressComponent) (uint32, bool) {
	if len(components) == 1 {
		if components[0].any {
			return 0, false
		}
		return components[0].value, true
	}
	value := uint32(0)
	for i, component := range components {
		if component.any {
			return 0, false
		}
		value |= component.value << (8 * i)
	}
	return value, true
}

// HasWildcard reports whether this tag covers more than a single point.
func (t Tag) HasWildcard() bool {
	for _, components := range [][]addressComponent{t.asdu, t.ioa} {
		for _, component := range components {
			if component.any {
				return true
			}
		}
	}
	return false
}

// GetAddressString spells the tag the way the tag handler parses it back, keeping the spelling the
// user chose so that re-parsing an address string yields the very same tag. plc4j's
// Iec608705104Tag.getAddressString returns null instead, which makes its tags unusable anywhere the
// address has to survive a round trip.
func (t Tag) GetAddressString() string {
	var sb strings.Builder
	for i, component := range t.asdu {
		if i > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(component.String())
	}
	sb.WriteByte('/')
	for i, component := range t.ioa {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(component.String())
	}
	return sb.String()
}

// GetValueType is what a subscription event for this tag carries. IEC 60870-5-104 addresses say
// nothing about the type of a point - the type identification travels with each ASDU, so the same
// address can report a boolean today and a scaled measurement after a configuration change. Every
// event therefore carries a struct: the value, the quality flags the wire came with, and the point's
// own timestamp. See AsduDecoder.go for the shape.
func (t Tag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

// GetArrayInfo is empty: a tag addresses one point, and a wildcard tag is a filter over points
// rather than an array of them - which of the covered points fires is not known in advance.
func (t Tag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

// GetPlcSubscriptionType is what a tag which wasn't added through one of the typed builder methods
// defaults to. A controlled station reports a point when it changes, so change-of-state it is.
func (t Tag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionChangeOfState
}

// GetDuration is not applicable: the reporting interval of a cyclically transmitted point is a
// property of the station's configuration, not something a subscription can ask for.
func (t Tag) GetDuration() time.Duration {
	return 0
}

func (t Tag) String() string {
	return "iec608705104.Tag{" + t.GetAddressString() + "}"
}
