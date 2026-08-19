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

package s7

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

type Browser struct {
	_default.DefaultBrowser

	connection *Connection

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewBrowser(connection *Connection, _options ...options.WithOption) *Browser {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	browser := &Browser{
		connection: connection,
		log:        customLogger,
		_options:   _options,
	}
	browser.DefaultBrowser = _default.NewDefaultBrowser(browser, _options...)
	return browser
}

// BrowseQuery returns the static memory areas plus one item per data block reported by the
// PLC's "list blocks of type DB" block function. Browse rides the UserData services, so
// devices without them get UNSUPPORTED. A failing block-list request degrades to the static
// areas only - some PLCs don't honor the block functions.
func (m *Browser) BrowseQuery(ctx context.Context, interceptor func(result apiModel.PlcBrowseItem) bool, queryName string, query apiModel.PlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	if !m.connection.driverContext.UserDataServicesSupported {
		return apiModel.PlcResponseCode_UNSUPPORTED, nil
	}
	subscribable := true
	items := staticAreaBrowseItems(subscribable)

	tpduId := m.connection.tpduGenerator.getAndIncrement()
	response, err := m.connection.sendUserData(ctx, tpduId, buildListBlocksOfTypeRequest(tpduId, blockTypeDataBlock), "browse_list_blocks")
	if err != nil {
		m.log.Debug().Err(err).Msg("Block-list browse failed; returning static areas only")
	} else if blockNumbers, err := parseListBlocksOfTypeResponse(response); err != nil {
		m.log.Debug().Err(err).Msg("Block-list browse rejected; returning static areas only")
	} else {
		m.log.Debug().Int("count", len(blockNumbers)).Msg("Block-list browse found data blocks")
		for _, blockNumber := range blockNumbers {
			items = append(items, dataBlockBrowseItem(blockNumber, subscribable))
		}
	}

	results := make([]apiModel.PlcBrowseItem, 0, len(items))
	for _, item := range items {
		if interceptor == nil || interceptor(item) {
			results = append(results, item)
		}
	}
	return apiModel.PlcResponseCode_OK, results
}

func staticAreaBrowseItems(subscribable bool) []apiModel.PlcBrowseItem {
	return []apiModel.PlcBrowseItem{
		staticAreaBrowseItem("%M", readWriteModel.MemoryArea_FLAGS_MARKERS, subscribable),
		staticAreaBrowseItem("%I", readWriteModel.MemoryArea_INPUTS, subscribable),
		staticAreaBrowseItem("%Q", readWriteModel.MemoryArea_OUTPUTS, subscribable),
	}
}

func staticAreaBrowseItem(name string, area readWriteModel.MemoryArea, subscribable bool) apiModel.PlcBrowseItem {
	// Placeholder tag for the area root, the caller composes the real offset / data type.
	return spiModel.NewDefaultPlcBrowseItem(
		NewTag(area, 0, 0, 0, 1, readWriteModel.TransportSize_BYTE),
		name,
		"",
		true,
		true,
		subscribable,
		nil,
		nil,
	)
}

func dataBlockBrowseItem(blockNumber uint16, subscribable bool) apiModel.PlcBrowseItem {
	return spiModel.NewDefaultPlcBrowseItem(
		NewTag(readWriteModel.MemoryArea_DATA_BLOCKS, blockNumber, 0, 0, 1, readWriteModel.TransportSize_BYTE),
		fmt.Sprintf("DB%d", blockNumber),
		"",
		true,
		true,
		subscribable,
		nil,
		nil,
	)
}
