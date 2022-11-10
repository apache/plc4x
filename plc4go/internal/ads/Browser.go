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
	"encoding/binary"
	"fmt"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	model2 "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func (m *Connection) Browse(ctx context.Context, browseRequest apiModel.PlcBrowseRequest) <-chan apiModel.PlcBrowseRequestResult {
	return m.BrowseWithInterceptor(ctx, browseRequest, func(result apiModel.PlcBrowseItem) bool {
		return true
	})
}

func (m *Connection) BrowseWithInterceptor(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseItem) bool) <-chan apiModel.PlcBrowseRequestResult {
	result := make(chan apiModel.PlcBrowseRequestResult)
	go func() {
		responseCodes := map[string]apiModel.PlcResponseCode{}
		results := map[string][]apiModel.PlcBrowseItem{}
		for _, queryName := range browseRequest.GetQueryNames() {
			query := browseRequest.GetQuery(queryName)
			responseCodes[queryName], results[queryName] = m.BrowseQuery(ctx, browseRequest, interceptor, queryName, query)
		}
		browseResponse := model2.NewDefaultPlcBrowseResponse(browseRequest, results, responseCodes)
		result <- &model2.DefaultPlcBrowseRequestResult{
			Request:  browseRequest,
			Response: &browseResponse,
			Err:      nil,
		}
	}()
	return result
}

func (m *Connection) BrowseQuery(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseItem) bool, queryName string, query apiModel.PlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	switch query.(type) {
	case SymbolicPlcQuery:
		return m.executeSymbolicAddressQuery(ctx, query.(SymbolicPlcQuery))
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
}

func (m *Connection) executeSymbolicAddressQuery(ctx context.Context, query SymbolicPlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	var err error

	// First read the sizes of the data type and symbol table, if needed.
	var tableSizes model.AdsTableSizes
	if m.dataTypeTable == nil || m.symbolTable == nil {
		tableSizes, err = m.readDataTypeTableAndSymbolTableSizes(ctx)
		if err != nil {
			return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
		}
	}

	// Then read the data type table, if needed.
	if m.dataTypeTable == nil {
		m.dataTypeTable, err = m.readDataTypeTable(ctx, tableSizes.GetDataTypeLength(), tableSizes.GetDataTypeCount())
		if err != nil {
			return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
		}
	}

	// Then read the symbol table, if needed.
	if m.symbolTable == nil {
		m.symbolTable, err = m.readSymbolTable(ctx, tableSizes.GetSymbolLength(), tableSizes.GetSymbolCount())
		if err != nil {
			return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
		}
	}

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

	if symbol, ok := m.symbolTable[symbolName]; !ok {
		// Couldn't find the base symbol
		return nil
	} else if len(remainingSegments) == 0 {
		// TODO: Convert the symbol itself into a PlcBrowseTag
		return nil
	} else {
		symbolDataTypeName := symbol.GetDataTypeName()
		if symbolDataType, ok := m.dataTypeTable[symbolDataTypeName]; ok {
			return m.filterDataTypes(symbolName, symbolDataType, symbolDataTypeName, remainingSegments)
		}
		// Couldn't find data type
		return nil
	}
}

func (m *Connection) filterDataTypes(parentName string, currentType model.AdsDataTypeTableEntry, currentPath string, remainingAddressSegments []string) []apiModel.PlcBrowseItem {
	if len(remainingAddressSegments) == 0 {
		arrayInfo := []apiModel.ArrayInfo{}
		for _, ai := range currentType.GetArrayInfo() {
			arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
				LowerBound: ai.GetLowerBound(),
				UpperBound: ai.GetUpperBound(),
			})
		}
		foundTag := &model2.DefaultPlcBrowseItem{
			Tag: SymbolicPlcTag{
				PlcTag: PlcTag{
					arrayInfo: arrayInfo,
				},
				SymbolicAddress: parentName,
			},
			Name:         parentName,
			DataTypeName: currentType.GetDataTypeName(),
			Readable:     false,
			Writable:     false,
			Subscribable: false,
			Options:      nil,
		}
		return []apiModel.PlcBrowseItem{foundTag}
	}

	currentAddressSegment := remainingAddressSegments[0]
	remainingAddressSegments = remainingAddressSegments[1:]
	for _, child := range currentType.GetChildren() {
		if child.GetPropertyName() == currentAddressSegment {
			childTypeName := child.GetDataTypeName()
			if symbolDataType, ok := m.dataTypeTable[childTypeName]; !ok {
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

func (m *Connection) readDataTypeTableAndSymbolTableSizes(ctx context.Context) (model.AdsTableSizes, error) {
	response, err := m.ExecuteAdsReadRequest(ctx, uint32(model.ReservedIndexGroups_ADSIGRP_SYMBOL_AND_DATA_TYPE_SIZES), 0x00000000, 24)
	if err != nil {
		return nil, fmt.Errorf("error reading table: %v", err)
	}

	// Parse and process the response
	tableSizes, err := model.AdsTableSizesParse(response.GetData())
	if err != nil {
		return nil, fmt.Errorf("error parsing table: %v", err)
	}
	return tableSizes, nil
}

func (m *Connection) readDataTypeTable(ctx context.Context, dataTableSize uint32, numDataTypes uint32) (map[string]model.AdsDataTypeTableEntry, error) {
	response, err := m.ExecuteAdsReadRequest(ctx, uint32(model.ReservedIndexGroups_ADSIGRP_DATA_TYPE_TABLE_UPLOAD), 0x00000000, dataTableSize)
	if err != nil {
		return nil, fmt.Errorf("error reading data-type table: %v", err)
	}

	// Parse and process the response
	readBuffer := utils.NewReadBufferByteBased(response.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	dataTypes := map[string]model.AdsDataTypeTableEntry{}
	for i := uint32(0); i < numDataTypes; i++ {
		dataType, err := model.AdsDataTypeTableEntryParseWithBuffer(readBuffer)
		if err != nil {
			return nil, fmt.Errorf("error parsing table: %v", err)
		}
		dataTypes[dataType.GetDataTypeName()] = dataType
	}
	return dataTypes, nil
}

func (m *Connection) readSymbolTable(ctx context.Context, symbolTableSize uint32, numSymbols uint32) (map[string]model.AdsSymbolTableEntry, error) {
	response, err := m.ExecuteAdsReadRequest(ctx, uint32(model.ReservedIndexGroups_ADSIGRP_SYM_UPLOAD), 0x00000000, symbolTableSize)
	if err != nil {
		return nil, fmt.Errorf("error reading data-type table: %v", err)
	}

	// Parse and process the response
	readBuffer := utils.NewReadBufferByteBased(response.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	symbols := map[string]model.AdsSymbolTableEntry{}
	for i := uint32(0); i < numSymbols; i++ {
		symbol, err := model.AdsSymbolTableEntryParseWithBuffer(readBuffer)
		if err != nil {
			return nil, fmt.Errorf("error parsing table")
		}
		symbols[symbol.GetName()] = symbol
	}
	return symbols, nil
}
