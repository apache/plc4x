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
	return m.BrowseWithInterceptor(ctx, browseRequest, func(result apiModel.PlcBrowseEvent) bool {
		return true
	})
}

func (m *Connection) BrowseWithInterceptor(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseEvent) bool) <-chan apiModel.PlcBrowseRequestResult {
	result := make(chan apiModel.PlcBrowseRequestResult)
	go func() {
		responseCodes := map[string]apiModel.PlcResponseCode{}
		results := map[string][]apiModel.PlcBrowseFoundField{}
		for _, fieldName := range browseRequest.GetFieldNames() {
			field := browseRequest.GetField(fieldName)
			responseCodes[fieldName], results[fieldName] = m.BrowseField(ctx, browseRequest, interceptor, fieldName, field)
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

func (m *Connection) BrowseField(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseEvent) bool, fieldName string, field apiModel.PlcField) (apiModel.PlcResponseCode, []apiModel.PlcBrowseFoundField) {
	switch field.(type) {
	case SymbolicPlcField:
		return m.executeSymbolicAddressQuery(ctx, field.(SymbolicPlcField))
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
}

func (m *Connection) executeSymbolicAddressQuery(ctx context.Context, field SymbolicPlcField) (apiModel.PlcResponseCode, []apiModel.PlcBrowseFoundField) {
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
	fields := m.filterSymbols(field.SymbolicAddress)
	return apiModel.PlcResponseCode_OK, fields
}

func (m *Connection) filterSymbols(filterExpression string) []apiModel.PlcBrowseFoundField {
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
		// TODO: Convert the symbol itself into a PlcBrowseField
		return nil
	} else {
		symbolDataTypeName := symbol.GetDataTypeName()
		if symbolDataType, ok := m.dataTypeTable[symbolDataTypeName]; !ok {
			// Couldn't find data type
			return nil
		} else {
			return m.filterDataTypes(symbolName, symbolDataType, symbolDataTypeName, remainingSegments)
		}
	}
}

/*
func LALALA(){
	for (AdsSymbolTableEntry symbol : symbolTable.values()) {
		// Get the datatype of this entry.
		AdsDataTypeTableEntry dataType = dataTypeTable.get(symbol.getDataTypeName());
		if (dataType == null) {
			System.out.printf("couldn't find datatype: %s%n", symbol.getDataTypeName());
			continue;
		}
		String itemName = (symbol.getComment() == null || symbol.getComment().isEmpty()) ? symbol.getName() : symbol.getComment();
		// Convert the plc value type from the ADS specific one to the PLC4X global one.
		org.apache.plc4x.java.api.types.PlcValueType plc4xPlcValueType = org.apache.plc4x.java.api.types.PlcValueType.valueOf(getPlcValueTypeForAdsDataType(dataType).toString());

		// If this type has children, add entries for its children.
		List<PlcBrowseItem> children = getBrowseItems(symbol.getName(), symbol.getGroup(), symbol.getOffset(), !symbol.getFlagReadOnly(), dataType);

		// Populate a map of protocol-dependent options.
		Map<String, PlcValue> options = new HashMap<>();
		options.put("comment", new PlcSTRING(symbol.getComment()));
		options.put("group-id", new PlcUDINT(symbol.getGroup()));
		options.put("offset", new PlcUDINT(symbol.getOffset()));
		options.put("size-in-bytes", new PlcUDINT(symbol.getSize()));

		if(plc4xPlcValueType == org.apache.plc4x.java.api.types.PlcValueType.List) {
			List<PlcBrowseItemArrayInfo> arrayInfo = new ArrayList<>();
			for (AdsDataTypeArrayInfo adsDataTypeArrayInfo : dataType.getArrayInfo()) {
				arrayInfo.add(new DefaultBrowseItemArrayInfo(
					adsDataTypeArrayInfo.getLowerBound(), adsDataTypeArrayInfo.getUpperBound()));
			}
			// Add the type itself.
			values.add(new DefaultListPlcBrowseItem(symbol.getName(), itemName, plc4xPlcValueType, arrayInfo,
				true, !symbol.getFlagReadOnly(), true, children, options));
		} else {
			// Add the type itself.
			values.add(new DefaultPlcBrowseItem(symbol.getName(), itemName, plc4xPlcValueType, true,
				!symbol.getFlagReadOnly(), true, children, options));
		}
	}
	DefaultPlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, PlcResponseCode.OK, values);

}
*/

func (m *Connection) filterDataTypes(parentName string, currentType model.AdsDataTypeTableEntry, currentPath string, remainingAddressSegments []string) []apiModel.PlcBrowseFoundField {
	if len(remainingAddressSegments) == 0 {
		var numElements int32
		var startElement int32
		var endElement int32
		if len(currentType.GetArrayInfo()) > 0 {
			numElements = int32(currentType.GetArrayInfo()[0].GetNumElements())
			startElement = int32(currentType.GetArrayInfo()[0].GetLowerBound())
			endElement = int32(currentType.GetArrayInfo()[0].GetUpperBound())
		}
		foundField := &model2.DefaultPlcBrowseQueryResult{
			Field: SymbolicPlcField{
				PlcField: PlcField{
					NumElements:  numElements,
					StartElement: startElement,
					EndElement:   endElement,
				},
				SymbolicAddress: parentName,
			},
			Name:              parentName,
			Readable:          false,
			Writable:          false,
			Subscribable:      false,
			PossibleDataTypes: nil,
			Attributes:        nil,
		}
		return []apiModel.PlcBrowseFoundField{foundField}
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
