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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
	"regexp"
	"strconv"
)

type FieldType uint8

const (
	FieldRandom FieldType = iota
	FieldState
	FieldStdOut
)

func (e FieldType) Name() string {
	switch e {
	case FieldRandom:
		return "RANDOM"
	case FieldState:
		return "STATE"
	case FieldStdOut:
		return "STDOUT"
	default:
		return "UNKNOWN"
	}
}

type FieldHandler struct {
	simulatedQuery *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		simulatedQuery: regexp.MustCompile(`^(?P<type>\w+)/(?P<name>[a-zA-Z0-9_\\.]+):(?P<dataType>[a-zA-Z0-9]+)(\[(?P<numElements>\d+)])?$`),
	}
}

func (m FieldHandler) ParseQuery(query string) (apiModel.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.simulatedQuery, query); match != nil {
		fieldTypeName, ok := match["type"]
		var fieldType FieldType
		if ok {
			switch fieldTypeName {
			case "RANDOM":
				fieldType = FieldRandom
				break
			case "STATE":
				fieldType = FieldState
				break
			case "STDOUT":
				fieldType = FieldStdOut
			default:
				return nil, errors.New("unknown field type '" + fieldTypeName + "'")
			}
		}
		fieldName, ok := match["name"]
		fieldDataTypeName, ok := match["dataType"]
		var fieldDataType model.SimulatedDataTypeSizes
		if ok {
			fieldDataType, _ = model.SimulatedDataTypeSizesByName(fieldDataTypeName)
			if fieldDataType == 0 {
				return nil, errors.New("unknown field data-type '" + fieldDataTypeName + "'")
			}
		}
		fieldNumElementsText, ok := match["numElements"]
		var fieldNumElements uint16
		if ok && len(fieldNumElementsText) > 0 {
			num, err := strconv.Atoi(fieldNumElementsText)
			if err != nil {
				return nil, errors.New("invalid size '" + fieldNumElementsText + "'")
			}
			fieldNumElements = uint16(num)
		} else {
			fieldNumElements = 1
		}
		return NewSimulatedField(fieldType, fieldName, fieldDataType, fieldNumElements), nil
	}
	return nil, errors.New("Invalid address format for address '" + query + "'")
}
