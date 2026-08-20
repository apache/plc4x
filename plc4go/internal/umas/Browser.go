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

package umas

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// maxTypeNestingDepth bounds how far the browser follows a type into its members and array elements.
// The dictionary comes from the device and nothing stops it from describing a struct which contains
// itself, which would otherwise recurse until the stack gives out.
const maxTypeNestingDepth = 16

// Browse lists the symbols of the PLC project. Ported from plc4j's UmasConnection.onBrowse, which
// hands out the whole symbol table for every query; here a query filters it by name, see
// TagHandler.ParseQuery.
func (c *Connection) Browse(ctx context.Context, browseRequest apiModel.PlcBrowseRequest) <-chan apiModel.PlcBrowseRequestResult {
	return c.BrowseWithInterceptor(ctx, browseRequest, func(apiModel.PlcBrowseItem) bool { return true })
}

func (c *Connection) BrowseWithInterceptor(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseItem) bool) <-chan apiModel.PlcBrowseRequestResult {
	result := make(chan apiModel.PlcBrowseRequestResult, 1)
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcBrowseRequestResult(browseRequest, nil,
					errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		responseCodes := map[string]apiModel.PlcResponseCode{}
		results := map[string][]apiModel.PlcBrowseItem{}
		for _, queryName := range browseRequest.GetQueryNames() {
			code, items := c.browseQuery(ctx, interceptor, browseRequest.GetQuery(queryName))
			responseCodes[queryName] = code
			results[queryName] = items
		}
		utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcBrowseRequestResult(browseRequest,
			spiModel.NewDefaultPlcBrowseResponse(browseRequest, results, responseCodes), nil))
	})
	return result
}

func (c *Connection) browseQuery(ctx context.Context, interceptor func(result apiModel.PlcBrowseItem) bool, query apiModel.PlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	symbolQuery, ok := query.(SymbolQuery)
	if !ok {
		c.log.Warn().Type("query", query).Msg("Not a UMAS browse query")
		return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
	}
	// The handshake normally downloads the dictionary; a browse on a connection whose download
	// failed retries it, the way plc4j's browse does.
	if err := c.ensureDataDictionary(ctx); err != nil {
		c.log.Error().Err(err).Msg("Can't browse without the data dictionary")
		return apiModel.PlcResponseCode_REMOTE_ERROR, nil
	}
	var items []apiModel.PlcBrowseItem
	for _, symbol := range c.session.symbols() {
		if !symbolQuery.Matches(strings.ToLower(symbol.GetValue())) {
			continue
		}
		item := c.buildBrowseItem(symbol.GetValue(), symbol.GetDataType(), 0)
		if interceptor != nil && !interceptor(item) {
			continue
		}
		items = append(items, item)
	}
	return apiModel.PlcResponseCode_OK, items
}

// buildBrowseItem describes one symbol (or one member of one), following its type into the array
// element type and struct members the dictionary knows about. Ported from plc4j's buildBrowseItem.
func (c *Connection) buildBrowseItem(name string, dataTypeId uint16, depth int) apiModel.PlcBrowseItem {
	valueType := mapToPlcValueType(dataTypeId)
	var arrayInfo []apiModel.ArrayInfo
	children := map[string]apiModel.PlcBrowseItem{}

	if customType, known := c.session.customType(dataTypeId); known && depth < maxTypeNestingDepth {
		if customType.isArray {
			// An array reports the type of its elements plus its bounds, so a caller sees
			// "array of DINT" rather than "unknown".
			valueType = c.resolveValueType(customType.elementTypeId, depth+1)
			arrayInfo = buildArrayInfo(customType.dimensions)
			children = c.buildStructChildren(customType.elementTypeId, depth+1)
		} else {
			valueType = apiValues.Struct
			children = c.buildStructChildren(dataTypeId, depth+1)
		}
	}

	return spiModel.NewDefaultPlcBrowseItem(
		NewTagWithType(name, valueType, arrayInfo),
		name,
		c.dataTypeName(dataTypeId),
		true, // every symbol of the dictionary is readable
		true, // and writable - plc4j says the same, the PLC is what refuses a write it doesn't like
		true, // subscribable, because subscriptions are emulated by polling the read path
		children,
		map[string]apiValues.PlcValue{},
	)
}

// resolveValueType is what a symbol of this type looks like to plc4x, following an array of arrays
// down to its element type. Ported from plc4j's resolveValueType.
func (c *Connection) resolveValueType(dataTypeId uint16, depth int) apiValues.PlcValueType {
	if depth >= maxTypeNestingDepth {
		return apiValues.RAW_BYTE_ARRAY
	}
	customType, known := c.session.customType(dataTypeId)
	if !known {
		return mapToPlcValueType(dataTypeId)
	}
	if !customType.isArray {
		return apiValues.Struct
	}
	return c.resolveValueType(customType.elementTypeId, depth+1)
}

// buildStructChildren describes the members of a struct type. A type which isn't a struct has none.
func (c *Connection) buildStructChildren(dataTypeId uint16, depth int) map[string]apiModel.PlcBrowseItem {
	children := map[string]apiModel.PlcBrowseItem{}
	if depth >= maxTypeNestingDepth {
		return children
	}
	customType, known := c.session.customType(dataTypeId)
	if !known || customType.isArray {
		return children
	}
	for _, field := range customType.fields {
		children[field.GetValue()] = c.buildBrowseItem(field.GetValue(), field.GetDataType(), depth+1)
	}
	return children
}

// buildArrayInfo turns the dimensions of an array type into what plc4x calls array info.
//
// Note that the plc4go DefaultArrayInfo reports GetSize as upperBound - lowerBound, where plc4j's
// reports upperBound - lowerBound + 1 for the same bounds. The bounds themselves are carried over
// unchanged - they are what the dictionary says - so a consumer reading the bounds sees the same
// thing in both languages.
func buildArrayInfo(dimensions []readWriteModel.UmasArrayDimension) []apiModel.ArrayInfo {
	if len(dimensions) == 0 {
		return nil
	}
	arrayInfo := make([]apiModel.ArrayInfo, 0, len(dimensions))
	for _, dimension := range dimensions {
		arrayInfo = append(arrayInfo, &spiModel.DefaultArrayInfo{
			LowerBound: dimension.GetStartIndex(),
			UpperBound: dimension.GetUpperBound(),
		})
	}
	return arrayInfo
}

// dataTypeName is what to call this type in a browse result: the name the project gave it for a
// custom type, the UMAS type name for a primitive, and the raw id for a type id the dictionary has
// nothing to say about.
func (c *Connection) dataTypeName(dataTypeId uint16) string {
	if customType, known := c.session.customType(dataTypeId); known && customType.name != "" {
		return customType.name
	}
	if dataType, ok := resolveDataType(dataTypeId); ok {
		return dataType.String()
	}
	return fmt.Sprintf("UNKNOWN(%d)", dataTypeId)
}
