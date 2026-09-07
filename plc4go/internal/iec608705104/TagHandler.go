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

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
)

// TagHandler parses IEC 60870-5-104 tag addresses.
//
// The syntax is the one plc4j's Iec608705104Tag documents in a commented-out regex which never
// compiled (an unterminated group name, and a '^' anchor which bound to only the first alternative
// of the ASDU part) and which its tag handler never used - plc4j's parseTag ignores its argument and
// answers everything with ASDU 0 / IOA 0. Recovered and implemented here:
//
//	<asdu>/<ioa>
//	  asdu = N | N/N          the two octet common address of ASDU, low octet first
//	  ioa  = N | N.N.N        the three octet information object address, low octet first
//
// '*' is accepted in any position and matches whatever the station reports there.
//
// Both spellings of a part address the same number: 1/0 as a split ASDU is ASDU 1, and 2.0.0 as a
// split IOA is information object 2. The octet order is least significant first, which is both the
// order the octets sit in on the wire and the order the recovered regex documents for the ASDU part
// ("low/high"). The plc4j side is being fixed to the very same syntax.
type TagHandler struct {
	log zerolog.Logger
}

func NewTagHandler(_options ...options.WithOption) TagHandler {
	return TagHandler{
		log: options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

// ParseTag turns an address into a Tag. Anything which isn't one of the four shapes above is
// refused rather than quietly resolved to a default point.
func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	// Surrounding whitespace is trimmed, the way the plc4j tag does, so that the two languages
	// accept exactly the same set of addresses. Whitespace *inside* an address is still refused.
	parts := strings.Split(strings.TrimSpace(tagAddress), "/")
	var asduPart []string
	var ioaPart string
	switch len(parts) {
	case 2:
		// <asdu>/<ioa>
		asduPart = parts[:1]
		ioaPart = parts[1]
	case 3:
		// <asduLow>/<asduHigh>/<ioa>
		asduPart = parts[:2]
		ioaPart = parts[2]
	default:
		return nil, errors.Errorf("Unable to parse %s: an address is <asdu>/<ioa> where asdu is N or N/N", tagAddress)
	}

	asdu, err := parseComponents(asduPart, maxAsduAddress, maxOctet)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse the common address of ASDU of %s", tagAddress)
	}

	ioaComponents := strings.Split(ioaPart, ".")
	if len(ioaComponents) != 1 && len(ioaComponents) != 3 {
		return nil, errors.Errorf("Unable to parse %s: the information object address is N or N.N.N", tagAddress)
	}
	ioa, err := parseComponents(ioaComponents, maxInformationObjectAddress, maxOctet)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse the information object address of %s", tagAddress)
	}

	return Tag{asdu: asdu, ioa: ioa}, nil
}

// ParseQuery is not supported. A controlled station can be asked for everything it has with a
// general interrogation, but that answers with ASDUs rather than with a directory of addresses, and
// neither this driver nor plc4j's browses.
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("This driver doesn't support browsing")
}

// parseComponents turns the numbers of one address part into components. A part spelled as a single
// number is range-checked against wholeMax, a part spelled octet by octet against octetMax.
func parseComponents(parts []string, wholeMax uint64, octetMax uint64) ([]addressComponent, error) {
	limit := octetMax
	if len(parts) == 1 {
		limit = wholeMax
	}
	// plc4j's Iec608705104Tag caps the digits of a component at the width of its limit (\d{1,5} for a
	// whole ASDU, \d{1,8} for a whole IOA, \d{1,3} for an octet), so a zero padded "0000000001" is
	// refused there even though it is in range. The same cap here keeps an address portable in both
	// directions rather than only from Go to Java.
	maxDigits := len(strconv.FormatUint(limit, 10))
	components := make([]addressComponent, 0, len(parts))
	for _, part := range parts {
		if part == wildcard {
			components = append(components, addressComponent{any: true})
			continue
		}
		if len(part) > maxDigits {
			return nil, errors.Errorf("%q is longer than the %d digits a value here may have", part, maxDigits)
		}
		value, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return nil, errors.Errorf("%q is neither a number nor %s", part, wildcard)
		}
		if value > limit {
			return nil, errors.Errorf("%d is out of range, the largest value here is %d", value, limit)
		}
		components = append(components, addressComponent{value: uint32(value)})
	}
	return components, nil
}
