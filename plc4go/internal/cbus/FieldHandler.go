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

package cbus

import (
	"encoding/hex"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type FieldType uint8

//go:generate stringer -type FieldType
const (
	STATUS FieldType = iota
	CAL
)

func (i FieldType) GetName() string {
	return i.String()
}

type FieldHandler struct {
	statusRequestPattern *regexp.Regexp
	calPattern           *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		statusRequestPattern: regexp.MustCompile(`^status/(?P<statusRequestType>(?P<binary>binary)|level=0x(?P<level>00|20|40|60|80|A0|C0|E0))/(?P<applicationId>.*)`),
		calPattern:           regexp.MustCompile(`^cal/(?P<calType>recall=\[(?P<recallParamNo>[\w\d]+),(?P<recallCount>\d+)]|identify=\[(?P<identifyAttribute>[\w\d]+)]|getstatus=\[(?P<getstatusParamNo>[\w\d]+),(?P<getstatusCount>\d+)])`),
	}
}

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.statusRequestPattern, query); match != nil {
		var level *byte
		var statusRequestType StatusRequestType
		statusRequestArgument := match["statusRequestType"]
		if statusRequestArgument != "" {
			if match["binary"] != "" {
				statusRequestType = StatusRequestTypeBinaryState
			} else if levelArgument := match["level"]; levelArgument != "" {
				statusRequestType = StatusRequestTypeLevel
				decodedHex, _ := hex.DecodeString(match["level"])
				if len(decodedHex) != 1 {
					panic("invalid state. Should have exactly 1")
				}
				level = &decodedHex[0]
			} else {
				return nil, errors.Errorf("Unknown statusRequestType%s", statusRequestArgument)
			}
		}
		var applicationId readWriteModel.ApplicationId
		applicationIdArgument := match["applicationId"]
		if strings.HasPrefix(applicationIdArgument, "0x") {
			decodedHex, err := hex.DecodeString(applicationIdArgument[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Not a valid hex")
			}
			if len(decodedHex) != 1 {
				return nil, errors.Errorf("Hex must be exatly one byte")
			}
			applicationId = readWriteModel.ApplicationId(decodedHex[0])
		} else {
			atoi, err := strconv.ParseUint(applicationIdArgument, 10, 8)
			if err != nil {
				applicationId = readWriteModel.ApplicationId(atoi)
			} else {
				applicationIdByName, ok := readWriteModel.ApplicationIdByName(applicationIdArgument)
				if !ok {
					return nil, errors.Errorf("Unknown applicationId%s", applicationIdArgument)
				}
				applicationId = applicationIdByName
			}
		}
		return NewStatusField(statusRequestType, level, applicationId, 1), nil
	} else if match := utils.GetSubgroupMatches(m.calPattern, query); match != nil {
		calTypeArgument := match["calType"]
		switch {
		case strings.HasPrefix(calTypeArgument, "recall="):
			panic("Implement me ")
		case strings.HasPrefix(calTypeArgument, "identify="):
			panic("Implement me ")
		case strings.HasPrefix(calTypeArgument, "getstatus="):
			panic("Implement me ")
		default:
			return nil, errors.Errorf("Invalid cal type %s", calTypeArgument)
		}
	} else {
		return nil, errors.Errorf("Unable to parse %s", query)
	}
}
