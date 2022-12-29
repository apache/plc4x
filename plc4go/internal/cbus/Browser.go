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
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/model"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Browser struct {
	_default.DefaultBrowser
	connection      *Connection
	messageCodec    spi.MessageCodec
	sequenceCounter uint8
}

func NewBrowser(connection *Connection, messageCodec spi.MessageCodec) *Browser {
	browser := Browser{
		connection:      connection,
		messageCodec:    messageCodec,
		sequenceCounter: 0,
	}
	browser.DefaultBrowser = _default.NewDefaultBrowser(browser)
	return &browser
}

func (m Browser) BrowseQuery(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseItem) bool, queryName string, query apiModel.PlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	var queryResults []apiModel.PlcBrowseItem
	switch query := query.(type) {
	case *unitInfoQuery:
		allUnits := false
		var units []readWriteModel.UnitAddress
		allAttributes := false
		var attributes []readWriteModel.Attribute
		if unitAddress := query.unitAddress; unitAddress != nil {
			units = append(units, *unitAddress)
		} else {
			// TODO: check if we still want the option to brute force all addresses
			installedUnitAddressBytes, err := m.getInstalledUnitAddressBytes(ctx)
			if err != nil {
				log.Warn().Err(err).Msg("Unable to get installed uints")
				return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
			}

			allUnits = true
			for i := 0; i <= 0xFF; i++ {
				unitAddressByte := byte(i)
				if _, ok := installedUnitAddressBytes[unitAddressByte]; ok {
					units = append(units, readWriteModel.NewUnitAddress(unitAddressByte))
				}
			}
		}
		if attribute := query.attribute; attribute != nil {
			attributes = append(attributes, *attribute)
		} else {
			allAttributes = true
			for _, attribute := range readWriteModel.AttributeValues {
				attributes = append(attributes, attribute)
			}
		}

		if allUnits {
			log.Info().Msg("Querying all (available) units")
		}
	unitLoop:
		for _, unit := range units {
			if err := ctx.Err(); err != nil {
				log.Info().Err(err).Msgf("Aborting scan at unit %s", unit)
				return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
			}
			unitAddress := unit.GetAddress()
			if !allUnits && allAttributes {
				log.Info().Msgf("Querying all attributes of unit %d", unitAddress)
			}
			event := log.Info()
			if allUnits {
				event = log.Debug()
			}
			event.Msgf("Query unit  %d", unitAddress)
			for _, attribute := range attributes {
				if err := ctx.Err(); err != nil {
					log.Info().Err(err).Msgf("Aborting scan at unit %s", unit)
					return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
				}
				if !allUnits && !allAttributes {
					log.Info().Msgf("Querying attribute %s of unit %d", attribute, unitAddress)
				} else {
					event.Msgf("unit %d: Query %s", unitAddress, attribute)
				}
				readTagName := fmt.Sprintf("%s/%d/%s", queryName, unitAddress, attribute)
				readRequest, _ := m.connection.ReadRequestBuilder().
					AddTag(readTagName, NewCALIdentifyTag(unit, attribute, 1)).
					Build()
				timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*2)
				requestResult := <-readRequest.ExecuteWithContext(timeoutCtx)
				timeoutCancel()
				if err := requestResult.GetErr(); err != nil {
					if !allUnits && !allAttributes {
						event.Err(err).Msgf("unit %d: Can't read attribute %s", unitAddress, attribute)
					}
					continue unitLoop
				}
				response := requestResult.GetResponse()
				if code := response.GetResponseCode(readTagName); code != apiModel.PlcResponseCode_OK {
					event.Msgf("unit %d: error reading tag %s. Code %s", unitAddress, attribute, code)
					continue unitLoop
				}
				queryResult := &model.DefaultPlcBrowseItem{
					Tag:          NewCALIdentifyTag(unit, attribute, 1),
					Name:         queryName,
					Readable:     true,
					Writable:     false,
					Subscribable: false,
					Options: map[string]values.PlcValue{
						"CurrentValue": response.GetValue(readTagName),
					},
				}
				if interceptor != nil {
					interceptor(queryResult)
				}
				queryResults = append(queryResults, queryResult)
			}
		}
	default:
		return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
	}
	return apiModel.PlcResponseCode_OK, queryResults
}

func (m Browser) getInstalledUnitAddressBytes(ctx context.Context) (map[byte]any, error) {
	// We need to presubscribe to catch the 2 followup responses
	subscriptionRequest, err := m.connection.SubscriptionRequestBuilder().
		AddEventTagAddress("installationMMIMonitor", "mmimonitor/*/NETWORK_CONTROL").
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error subscribing to the installation MMI")
	}
	subCtx, subCtxCancel := context.WithTimeout(ctx, time.Second*2)
	subscriptionResult := <-subscriptionRequest.ExecuteWithContext(subCtx)
	subCtxCancel()
	if err := subscriptionResult.GetErr(); err != nil {
		return nil, errors.Wrap(err, "Error subscribing to the mmi")
	}
	if responseCode := subscriptionResult.GetResponse().GetResponseCode("installationMMIMonitor"); responseCode != apiModel.PlcResponseCode_OK {
		return nil, errors.Errorf("Got %s", responseCode)
	}
	subscriptionHandle, err := subscriptionResult.GetResponse().GetSubscriptionHandle("installationMMIMonitor")
	if err != nil {
		return nil, errors.Wrap(err, "Error getting the subscription handle")
	}

	blockOffset0Received := false
	blockOffset0ReceivedChan := make(chan any, 100) // We only expect one, but we make it a bit bigger to no clog up
	blockOffset88Received := false
	blockOffset88ReceivedChan := make(chan any, 100) // We only expect one, but we make it a bit bigger to no clog up
	blockOffset176Received := false
	blockOffset176ReceivedChan := make(chan any, 100) // We only expect one, but we make it a bit bigger to no clog up
	result := make(map[byte]any)
	plcConsumerRegistration := subscriptionHandle.Register(func(event apiModel.PlcSubscriptionEvent) {
		if responseCode := event.GetResponseCode("installationMMIMonitor"); responseCode != apiModel.PlcResponseCode_OK {
			log.Warn().Msgf("Ignoring %v", event)
			return
		}
		rootValue := event.GetValue("installationMMIMonitor")
		if !rootValue.IsStruct() {
			log.Warn().Msgf("Ignoring %v should be a struct", rootValue)
			return
		}
		rootStruct := rootValue.GetStruct()
		if applicationValue := rootStruct["application"]; applicationValue == nil || !applicationValue.IsString() || applicationValue.GetString() != "NETWORK_CONTROL" {
			log.Warn().Msgf("Ignoring %v should contain a application tag of type string with value NETWORK_CONTROL", rootStruct)
			return
		}
		var blockStart int
		if blockStartValue := rootStruct["blockStart"]; blockStartValue == nil || !blockStartValue.IsByte() {
			log.Warn().Msgf("Ignoring %v should contain a blockStart tag of type byte", rootStruct)
			return
		} else {
			blockStart = int(blockStartValue.GetByte())
		}
		switch blockStart {
		case 88:
			select {
			case blockOffset88ReceivedChan <- true:
			default:
			}
		case 176:
			select {
			case blockOffset176ReceivedChan <- true:
			default:
			}
		case 0:
			select {
			case blockOffset0ReceivedChan <- true:
			default:
			}
		}

		if plcListValue := rootStruct["values"]; plcListValue == nil || !plcListValue.IsList() {
			log.Warn().Msgf("Ignoring %v should contain a values tag of type list", rootStruct)
			return
		} else {
			for unitByteAddress, plcValue := range plcListValue.GetList() {
				unitByteAddress = blockStart + unitByteAddress
				if !plcValue.IsString() {
					log.Warn().Msgf("Ignoring %v at %d should be a string", plcValue, unitByteAddress)
					return
				}
				switch plcValue.GetString() {
				case readWriteModel.GAVState_ON.PLC4XEnumName(), readWriteModel.GAVState_OFF.PLC4XEnumName():
					log.Debug().Msgf("unit %d does exists", unitByteAddress)
					result[byte(unitByteAddress)] = true
				case readWriteModel.GAVState_DOES_NOT_EXIST.PLC4XEnumName():
					log.Debug().Msgf("unit %d does not exists", unitByteAddress)
				case readWriteModel.GAVState_ERROR.PLC4XEnumName():
					log.Warn().Msgf("unit %d is in error state", unitByteAddress)
				}
			}
		}
	})
	defer plcConsumerRegistration.Unregister()

	readRequest, err := m.connection.ReadRequestBuilder().
		AddTagAddress("installationMMI", "status/binary/0xFF").
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting the installation MMI")
	}
	readCtx, readCtxCancel := context.WithTimeout(ctx, time.Second*2)
	readRequestResult := <-readRequest.ExecuteWithContext(readCtx)
	readCtxCancel()
	if err := readRequestResult.GetErr(); err != nil {
		return nil, errors.Wrap(err, "Error reading the mmi")
	}
	if responseCode := readRequestResult.GetResponse().GetResponseCode("installationMMI"); responseCode == apiModel.PlcResponseCode_OK {
		rootValue := readRequestResult.GetResponse().GetValue("installationMMI")
		if !rootValue.IsStruct() {
			return nil, errors.Errorf("%v should be a struct", rootValue)
		}
		rootStruct := rootValue.GetStruct()
		if applicationValue := rootStruct["application"]; applicationValue == nil || !applicationValue.IsString() || applicationValue.GetString() != "NETWORK_CONTROL" {
			return nil, errors.Errorf("%v should contain a application tag of type string with value NETWORK_CONTROL", rootStruct)
		}
		var blockStart int
		if blockStartValue := rootStruct["blockStart"]; blockStartValue == nil || !blockStartValue.IsByte() || blockStartValue.GetByte() != 0 {
			return nil, errors.Errorf("%v should contain a blockStart tag of type byte with value 0", rootStruct)
		} else {
			blockStart = int(blockStartValue.GetByte())
		}

		if plcListValue := rootStruct["values"]; plcListValue == nil || !plcListValue.IsList() {
			return nil, errors.Errorf("%v should contain a values tag of type list", rootStruct)
		} else {
			for unitByteAddress, plcValue := range plcListValue.GetList() {
				unitByteAddress = blockStart + unitByteAddress
				if !plcValue.IsString() {
					return nil, errors.Errorf("%v at %d should be a string", plcValue, unitByteAddress)
				}
				switch plcValue.GetString() {
				case readWriteModel.GAVState_ON.PLC4XEnumName(), readWriteModel.GAVState_OFF.PLC4XEnumName():
					log.Debug().Msgf("unit %d does exists", unitByteAddress)
					result[byte(unitByteAddress)] = true
				case readWriteModel.GAVState_DOES_NOT_EXIST.PLC4XEnumName():
					log.Debug().Msgf("unit %d does not exists", unitByteAddress)
				case readWriteModel.GAVState_ERROR.PLC4XEnumName():
					log.Warn().Msgf("unit %d is in error state", unitByteAddress)
				}
			}
		}
		blockOffset0Received = true
	} else {
		log.Warn().Msgf("We got %s as response code for installation mmi so we rely on getting it via subscription", responseCode)
	}

	syncCtx, syncCtxCancel := context.WithTimeout(ctx, time.Second*2)
	defer syncCtxCancel()
	for !blockOffset0Received || !blockOffset88Received || !blockOffset176Received {
		select {
		case <-blockOffset0ReceivedChan:
			log.Trace().Msg("Offset 0 received")
			blockOffset0Received = true
		case <-blockOffset88ReceivedChan:
			log.Trace().Msg("Offset 88 received")
			blockOffset88Received = true
		case <-blockOffset176ReceivedChan:
			log.Trace().Msg("Offset 176 received")
			blockOffset176Received = true
		case <-syncCtx.Done():
			return nil, errors.Wrap(err, "error waiting for other offsets")
		}
	}
	return result, nil
}
