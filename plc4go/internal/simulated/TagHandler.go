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

package simulated

import (
	"fmt"
	"regexp"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TagType uint8

const (
	TagRandom TagType = iota
	TagState
	TagStdOut
)

func (e TagType) Name() string {
	switch e {
	case TagRandom:
		return "RANDOM"
	case TagState:
		return "STATE"
	case TagStdOut:
		return "STDOUT"
	default:
		return "UNKNOWN"
	}
}

func (e TagType) String() string {
	return e.Name()
}

type TagHandler struct {
	simulatedQuery *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		// The selection sits between the name and the type, as it does in plc4j's
		// SimulatedTag.ADDRESS_PATTERN: "RANDOM/foo[0..3]:INT".
		simulatedQuery: regexp.MustCompile(`^(?P<type>\w+)/(?P<name>[a-zA-Z0-9_\\.]+)` + spiModel.ArrayGroupPattern + `:(?P<dataType>[a-zA-Z0-9]+)$`),
	}
}

// elementsOf resolves an address's array expression to a number of elements. This driver
// addresses a named variable rather than a numeric offset, so a selection that does not start at
// the first element has nothing to apply to and is reported rather than quietly ignored.
//
// The count is carried as a uint16 from here on, so it is bounded as one: an unchecked cast
// turned a count above 65535 into whatever the low two bytes happened to be, handing back a tag
// of 4464 elements for a request of 70000 and an empty one for 65536.
func elementsOf(expression string, tagAddress string) (uint16, bool, error) {
	if expression == "" {
		return 1, false, nil
	}
	dimensions, err := spiModel.ParseArrayExpression(expression, tagAddress, spiModel.SingleDimension)
	if err != nil {
		return 0, false, err
	}
	dimension := dimensions[0]
	if dimension.GetLowerBound() != dimension.GetBase() {
		return 0, false, errors.Errorf("Array selection '%s' in tag '%s' must start at the first element: "+
			"this driver addresses a named variable, so there is no offset to start from",
			expression, tagAddress)
	}
	elements := dimension.GetSize()
	if elements < 1 || elements > 0xFFFF {
		return 0, false, errors.Errorf("A tag of %d elements in '%s' is more than a simulated tag may hold",
			elements, tagAddress)
	}
	// A range is an array even when it spans one element, which the count cannot say.
	return uint16(elements), dimension.IsRange(), nil
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.simulatedQuery, tagAddress); match != nil {
		tagTypeName, ok := match["type"]
		var tagType TagType
		if ok {
			switch tagTypeName {
			case "RANDOM":
				tagType = TagRandom
				break
			case "STATE":
				tagType = TagState
				break
			case "STDOUT":
				tagType = TagStdOut
			default:
				return nil, errors.New("unknown tag type '" + tagTypeName + "'")
			}
		}
		tagName, ok := match["name"]
		tagDataTypeName, ok := match["dataType"]
		var tagDataType model.SimulatedDataTypeSizes
		if ok {
			tagDataType, _ = model.SimulatedDataTypeSizesByName(tagDataTypeName)
			if tagDataType == 0 {
				return nil, errors.New("unknown tag data-type '" + tagDataTypeName + "'")
			}
		}
		tagNumElements, explicitRange, err := elementsOf(match["array"], tagAddress)
		if err != nil {
			return nil, err
		}
		return NewSimulatedTagWithShape(tagType, tagName, tagDataType, tagNumElements, explicitRange), nil
	}
	// "RANDOM/foo:INT[4]" - the count after the type - no longer parses, so name the form to
	// write rather than reporting only that nothing matched.
	return nil, spiModel.InvalidAddressError(tagAddress,
		"{type}/{name}[selection]:{TYPE} - for example RANDOM/foo[0..3]:INT")
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
