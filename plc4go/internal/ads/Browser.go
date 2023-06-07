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

package ads

import (
	"context"
	"github.com/pkg/errors"
	"runtime/debug"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

func (m *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return spiModel.NewDefaultPlcBrowseRequestBuilder(m.GetPlcTagHandler(), m)
}

func (m *Connection) Browse(ctx context.Context, browseRequest apiModel.PlcBrowseRequest) <-chan apiModel.PlcBrowseRequestResult {
	return m.BrowseWithInterceptor(ctx, browseRequest, func(result apiModel.PlcBrowseItem) bool {
		return true
	})
}

func (m *Connection) BrowseWithInterceptor(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseItem) bool) <-chan apiModel.PlcBrowseRequestResult {
	result := make(chan apiModel.PlcBrowseRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcBrowseRequestResult(browseRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		responseCodes := map[string]apiModel.PlcResponseCode{}
		results := map[string][]apiModel.PlcBrowseItem{}
		for _, queryName := range browseRequest.GetQueryNames() {
			query := browseRequest.GetQuery(queryName)
			responseCodes[queryName], results[queryName] = m.BrowseQuery(ctx, interceptor, queryName, query)
		}
		browseResponse := spiModel.NewDefaultPlcBrowseResponse(browseRequest, results, responseCodes)
		result <- spiModel.NewDefaultPlcBrowseRequestResult(browseRequest, browseResponse, nil)
	}()
	return result
}

func (m *Connection) BrowseQuery(_ context.Context, _ func(result apiModel.PlcBrowseItem) bool, _ string, query apiModel.PlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	switch query.(type) {
	case SymbolicPlcQuery:
		return m.executeSymbolicAddressQuery(query.(SymbolicPlcQuery))
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
}

func (m *Connection) executeSymbolicAddressQuery(query SymbolicPlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	// Process the data type and symbol tables to produce the response.
	tags := m.filterSymbols(query.GetSymbolicAddressPattern())
	return apiModel.PlcResponseCode_OK, tags
}

func (m *Connection) filterSymbols(filterExpression string) []apiModel.PlcBrowseItem {
	if len(filterExpression) == 0 {
		return nil
	}
	addressSegments := strings.Split(filterExpression, ".")

	// The symbol name consists of the first two segments of the address
	// Some addresses only have one segment, so in that case we'll simply use that.
	symbolName := addressSegments[0]
	remainingSegments := addressSegments[1:]
	if len(addressSegments) > 0 {
		symbolName = symbolName + "." + remainingSegments[0]
		remainingSegments = remainingSegments[1:]
	}

	if symbol, ok := m.driverContext.symbolTable[symbolName]; !ok {
		// Couldn't find the base symbol
		return nil
	} else if len(remainingSegments) == 0 {
		// TODO: Convert the symbol itself into a PlcBrowseTag
		return nil
	} else {
		symbolDataTypeName := symbol.GetDataTypeName()
		if symbolDataType, ok := m.driverContext.dataTypeTable[symbolDataTypeName]; ok {
			return m.filterDataTypes(symbolName, symbolDataType, symbolDataTypeName, remainingSegments)
		}
		// Couldn't find data type
		return nil
	}
}

func (m *Connection) filterDataTypes(parentName string, currentType driverModel.AdsDataTypeTableEntry, currentPath string, remainingAddressSegments []string) []apiModel.PlcBrowseItem {
	if len(remainingAddressSegments) == 0 {
		var arrayInfo []apiModel.ArrayInfo
		for _, ai := range currentType.GetArrayInfo() {
			arrayInfo = append(arrayInfo, &spiModel.DefaultArrayInfo{
				LowerBound: ai.GetLowerBound(),
				UpperBound: ai.GetUpperBound(),
			})
		}
		foundTag := spiModel.NewDefaultPlcBrowseItem(
			model.SymbolicPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: arrayInfo,
				},
				SymbolicAddress: parentName,
			},
			parentName,
			currentType.GetDataTypeName(),
			false,
			false,
			false,
			nil,
			nil,
		)
		return []apiModel.PlcBrowseItem{foundTag}
	}

	currentAddressSegment := remainingAddressSegments[0]
	remainingAddressSegments = remainingAddressSegments[1:]
	for _, child := range currentType.GetChildren() {
		if child.GetPropertyName() == currentAddressSegment {
			childTypeName := child.GetDataTypeName()
			if symbolDataType, ok := m.driverContext.dataTypeTable[childTypeName]; !ok {
				// TODO: Couldn't find data type with the name defined in the protperty.
				return nil
			} else {
				return m.filterDataTypes(parentName+"."+child.GetPropertyName(), symbolDataType,
					currentPath+"."+currentAddressSegment, remainingAddressSegments)
			}
		}
	}
	// TODO: Couldn't find property with the given name.
	return nil
}
