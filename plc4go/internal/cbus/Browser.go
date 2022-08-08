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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/rs/zerolog/log"
)

type Browser struct {
	connection      *Connection
	messageCodec    spi.MessageCodec
	sequenceCounter uint8
}

func NewBrowser(connection *Connection, messageCodec spi.MessageCodec) *Browser {
	return &Browser{
		connection:      connection,
		messageCodec:    messageCodec,
		sequenceCounter: 0,
	}
}

func (m Browser) Browse(browseRequest apiModel.PlcBrowseRequest) <-chan apiModel.PlcBrowseRequestResult {
	return m.BrowseWithInterceptor(browseRequest, nil)
}

func (m Browser) BrowseWithInterceptor(browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseEvent) bool) <-chan apiModel.PlcBrowseRequestResult {
	result := make(chan apiModel.PlcBrowseRequestResult)

	go func() {
		responseCodes := map[string]apiModel.PlcResponseCode{}
		results := map[string][]apiModel.PlcBrowseFoundField{}
		for _, fieldName := range browseRequest.GetFieldNames() {
			field := browseRequest.GetField(fieldName)

			var queryResults []apiModel.PlcBrowseFoundField
			switch field := field.(type) {
			case *unitInfoField:
				var units []readWriteModel.UnitAddress
				var attributes []readWriteModel.Attribute
				if unitAddress := field.unitAddress; unitAddress != nil {
					units = append(units, *unitAddress)
				} else {
					for i := 0; i <= 0xFF; i++ {
						units = append(units, readWriteModel.NewUnitAddress(byte(i)))
					}
				}
				if attribute := field.attribute; attribute != nil {
					attributes = append(attributes, *attribute)
				} else {
					for _, attribute := range readWriteModel.AttributeValues {
						attributes = append(attributes, attribute)
					}
				}
			unitLoop:
				for _, unit := range units {
					for _, attribute := range attributes {
						unitAddress := unit.GetAddress()
						log.Info().Msgf("unit %d: Query %s", unitAddress, attribute)
						readFieldName := fmt.Sprintf("%s/%d/%s", fieldName, unitAddress, attribute)
						readRequest, _ := m.connection.ReadRequestBuilder().
							AddField(readFieldName, NewCALIdentifyField(unit, attribute, 1)).
							Build()
						requestResult := <-readRequest.Execute()
						if err := requestResult.GetErr(); err != nil {
							log.Info().Err(err).Msgf("unit %d: Can't read attribute %s", unitAddress, attribute)
							continue unitLoop
						}
						queryResults = append(queryResults, &model.DefaultPlcBrowseQueryResult{
							Field:        NewCALIdentifyField(unit, attribute, 1),
							Name:         fieldName,
							Readable:     true,
							Writable:     false,
							Subscribable: false,
							Attributes: map[string]values.PlcValue{
								"CurrentValue": requestResult.GetResponse().GetValue(readFieldName),
							},
						})
					}
				}
			default:
				responseCodes[fieldName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			}
		}
		result <- &model.DefaultPlcBrowseRequestResult{
			Request:  browseRequest,
			Response: model.NewDefaultPlcBrowseResponse(browseRequest, results, responseCodes),
			Err:      nil,
		}
	}()
	return result
}
