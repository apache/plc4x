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
	"errors"
	"fmt"
	"regexp"
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
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

type TagHandler struct {
	simulatedQuery *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		simulatedQuery: regexp.MustCompile(`^(?P<type>\w+)/(?P<name>[a-zA-Z0-9_\\.]+):(?P<dataType>[a-zA-Z0-9]+)(\[(?P<numElements>\d+)])?$`),
	}
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
		tagNumElementsText, ok := match["numElements"]
		var tagNumElements uint16
		if ok && len(tagNumElementsText) > 0 {
			num, err := strconv.Atoi(tagNumElementsText)
			if err != nil {
				return nil, errors.New("invalid size '" + tagNumElementsText + "'")
			}
			tagNumElements = uint16(num)
		} else {
			tagNumElements = 1
		}
		return NewSimulatedTag(tagType, tagName, tagDataType, tagNumElements), nil
	}
	return nil, errors.New("Invalid address format for address '" + tagAddress + "'")
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
