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

package bacnetip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

type TagHandler struct {
	addressPattern       *regexp.Regexp
	propertyFieldPattern *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		addressPattern:       regexp.MustCompile(`^(?P<objectType>[\d\w]+),(?P<objectInstance>\d+)/(?P<propertyIdentifiers>[\d\w]+(?:\[\d+])?(?:&[\d\w]+(?:\[\d+])?)*)`),
		propertyFieldPattern: regexp.MustCompile(`^(?P<propertyIdentifier>[\d\w]+)(?:\[(?P<arrayIndex>\d+)])?$`),
	}
}

const (
	OBJECT_TYPE          = "objectType"
	OBJECT_INSTANCE      = "objectInstance"
	PROPERTY_IDENTIFIERS = "propertyIdentifiers"
	PROPERTY_IDENTIFIER  = "propertyIdentifier"
	ARRAY_INDEX          = "arrayIndex"
)

func (m TagHandler) ParseTag(tagString string) (apiModel.PlcTag, error) {
	if addressMatch := utils.GetSubgroupMatches(m.addressPattern, tagString); addressMatch != nil {
		var result plcTag
		{
			objectTypeMatch := addressMatch[OBJECT_TYPE]
			if parsedObjectType, parseUintErr := strconv.ParseUint(objectTypeMatch, 10, 16); parseUintErr != nil {
				if objectType, ok := readWriteModel.BACnetObjectTypeByName(objectTypeMatch); !ok {
					return nil, errors.Errorf("Unknown object type %s", objectTypeMatch)
				} else {
					result.ObjectId.ObjectIdType = &objectType
				}
			} else {
				proprietaryType := uint16(parsedObjectType)
				result.ObjectId.ObjectIdTypeProprietary = &proprietaryType
			}
		}
		if parsedObjectInstance, err := strconv.ParseUint(addressMatch[OBJECT_INSTANCE], 10, 32); err != nil {
			return nil, errors.Wrap(err, "Error parsing object instance")
		} else {
			result.ObjectId.ObjectIdInstance = uint32(parsedObjectInstance)
		}

		for _, propertyString := range strings.Split(addressMatch[PROPERTY_IDENTIFIERS], "&") {
			var _property struct {
				PropertyIdentifier            *readWriteModel.BACnetPropertyIdentifier
				PropertyIdentifierProprietary *uint32
				ArrayIndex                    *uint
			}
			propertyMatch := utils.GetSubgroupMatches(m.propertyFieldPattern, propertyString)
			propertyIdentifierMatch := propertyMatch[PROPERTY_IDENTIFIER]
			if parsedPropertyType, parseUintErr := strconv.ParseUint(propertyIdentifierMatch, 10, 32); parseUintErr != nil {
				if propertyIdentifier, ok := readWriteModel.BACnetPropertyIdentifierByName(propertyIdentifierMatch); !ok {
					return nil, errors.Errorf("Unknown property type %s", propertyIdentifierMatch)
				} else {
					_property.PropertyIdentifier = &propertyIdentifier
				}
			} else {
				proprietaryType := uint32(parsedPropertyType)
				_property.PropertyIdentifierProprietary = &proprietaryType
			}
			if arrayIndexMatch := propertyMatch[ARRAY_INDEX]; arrayIndexMatch != "" {
				if parsedArrayIndex, err := strconv.ParseUint(arrayIndexMatch, 10, 32); err != nil {
					return nil, errors.Wrap(err, "Error parsing array index")
				} else {
					arrayIndex := uint(parsedArrayIndex)
					_property.ArrayIndex = &arrayIndex
				}
			}

			result.Properties = append(result.Properties, _property)
		}
		return result, nil
	}
	return nil, errors.Errorf("Unable to parse %s", tagString)
}

func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
