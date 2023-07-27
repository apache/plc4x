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

package opcua

import (
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"regexp"
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

type TagHandler struct {
	tagAddress *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		tagAddress: regexp.MustCompile(`^ns=(?P<namespace>\d+);(?P<identifierType>[isgb])=(?P<identifier>[^;]+)?(;(?P<datatype>[a-zA-Z_]+))?`),
	}
}

type CommandAndArgumentsCount interface {
	fmt.Stringer
	PLC4XEnumName() string
	NumberOfArguments() uint8
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.tagAddress, tagAddress); match != nil {
		return m.handleTagAddress(match)
	} else {
		return nil, errors.Errorf("Unable to parse %s", tagAddress)
	}
}

func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("This driver doesn't support browsing")
}

func (m TagHandler) handleTagAddress(match map[string]string) (apiModel.PlcTag, error) {
	namespace, err := strconv.Atoi(match["namespace"])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing namespace")
	}
	identifier := match["identifier"]
	if identifier == "" {
		return nil, errors.New("empty identifier")
	}
	identifierType, ok := readWriteModel.OpcuaIdentifierTypeByValue(match["identifierType"])
	if !ok {
		return nil, errors.Errorf("No identifier found for " + match["identifierType"])
	}
	dataType := readWriteModel.OpcuaDataType_NULL
	if dataTypeMatch := match["dataType"]; dataTypeMatch != "" {
		dataType, ok = readWriteModel.OpcuaDataTypeByName(dataTypeMatch)
		if !ok {
			return nil, errors.Errorf("No identifier found for " + match["dataType"])
		}
	}
	return NewTag(namespace, identifier, identifierType, dataType), nil
}
