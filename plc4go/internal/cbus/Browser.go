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
	"runtime/debug"
	"sync"
	"time"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Browser struct {
	_default.DefaultBrowser
	connection      plc4go.PlcConnection
	sequenceCounter uint8

	log zerolog.Logger
}

func NewBrowser(connection plc4go.PlcConnection, _options ...options.WithOption) *Browser {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	browser := &Browser{
		connection:      connection,
		sequenceCounter: 0,

		log: customLogger,
	}
	browser.DefaultBrowser = _default.NewDefaultBrowser(browser, _options...)
	return browser
}

func (m *Browser) BrowseQuery(ctx context.Context, interceptor func(result apiModel.PlcBrowseItem) bool, queryName string, query apiModel.PlcQuery) (responseCode apiModel.PlcResponseCode, queryResults []apiModel.PlcBrowseItem) {
	switch query := query.(type) {
	case *unitInfoQuery:
		return m.browseUnitInfo(ctx, interceptor, queryName, query)
	default:
		m.log.Warn().Type("query", query).Msg("unsupported query type supplied %T")
		return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
	}
}

func (m *Browser) browseUnitInfo(ctx context.Context, interceptor func(result apiModel.PlcBrowseItem) bool, queryName string, query *unitInfoQuery) (responseCode apiModel.PlcResponseCode, queryResults []apiModel.PlcBrowseItem) {
	m.log.Trace().Msg("extract units")
	units, allUnits, err := m.extractUnits(ctx, query, m.getInstalledUnitAddressBytes)
	if err != nil {
		m.log.Error().Err(err).Msg("Error extracting units")
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	attributes, allAttributes := m.extractAttributes(query)

	if allUnits {
		m.log.Info().Msg("Querying all (available) units")
	} else {
		m.log.Debug().Interface("units", units).Msg("Querying units")
	}
unitLoop:
	for _, unit := range units {
		unitLog := m.log.With().Stringer("unit", unit).Logger()
		unitLog.Trace().Msg("checking unit")
		if err := ctx.Err(); err != nil {
			unitLog.Info().Err(err).Msg("Aborting scan at unit")
			return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
		}
		unitAddress := unit.GetAddress()
		if !allUnits && allAttributes {
			m.log.Info().
				Uint8("unitAddress", unitAddress).
				Msg("Querying all attributes of unit")
		}
		event := m.log.Info()
		if allUnits {
			event = m.log.Debug()
		}
		event.Uint8("unitAddress", unitAddress).Msg("Query unit")
		for _, attribute := range attributes {
			if err := ctx.Err(); err != nil {
				unitLog.Info().Err(err).Msg("Aborting scan at unit")
				return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
			}
			if !allUnits && !allAttributes {
				m.log.Info().
					Uint8("unitAddress", unitAddress).
					Stringer("attribute", attribute).
					Msg("Querying attribute of unit")
			} else {
				event.Uint8("unitAddress", unitAddress).
					Stringer("attribute", attribute).
					Msg("unit unitAddress: Query attribute")
			}
			m.log.Trace().Msg("Building request")
			readTagName := fmt.Sprintf("%s/%d/%s", queryName, unitAddress, attribute)
			readRequest, _ := m.connection.ReadRequestBuilder().
				AddTag(readTagName, NewCALIdentifyTag(unit, nil /*TODO: add bridge support*/, attribute, 1)).
				Build()
			timeout := 5 * time.Second
			timeoutCtx, timeoutCancel := context.WithTimeout(ctx, timeout)
			m.log.Trace().
				Stringer("readRequest", readRequest).
				Dur("timeout", timeout).
				Msg("Executing readRequest with timeout")
			requestResult := <-readRequest.ExecuteWithContext(timeoutCtx)
			m.log.Trace().Stringer("requestResult", requestResult).Msg("got a response")
			timeoutCancel()
			if err := requestResult.GetErr(); err != nil {
				if allUnits || allAttributes {
					event = m.log.Trace()
				}
				event.Err(err).
					Uint8("unitAddress", unitAddress).
					Stringer("attribute", attribute).
					Msg("unit unitAddress: Can't read attribute attribute")
				continue unitLoop
			}
			response := requestResult.GetResponse()
			if code := response.GetResponseCode(readTagName); code != apiModel.PlcResponseCode_OK {
				event.
					Uint8("unitAddress", unitAddress).
					Stringer("attribute", attribute).
					Stringer("code", code).
					Msg("unit unitAddress: error reading tag attribute. Code %s")
				continue unitLoop
			}
			queryResult := spiModel.NewDefaultPlcBrowseItem(
				NewCALIdentifyTag(unit, nil /*TODO: add bridge support*/, attribute, 1),
				queryName,
				"",
				true,
				false,
				false,
				nil,
				map[string]values.PlcValue{
					"CurrentValue": response.GetValue(readTagName),
				},
			)
			if interceptor != nil {
				m.log.Trace().Msg("forwarding query result to interceptor")
				interceptor(queryResult)
			}
			queryResults = append(queryResults, queryResult)
		}
	}
	return apiModel.PlcResponseCode_OK, queryResults
}

func (m *Browser) extractUnits(ctx context.Context, query *unitInfoQuery, getInstalledUnitAddressBytes func(ctx context.Context) (map[byte]any, error)) ([]readWriteModel.UnitAddress, bool, error) {
	if unitAddress := query.unitAddress; unitAddress != nil {
		return []readWriteModel.UnitAddress{unitAddress}, false, nil
	}

	// TODO: check if we still want the option to brute force all addresses
	installedUnitAddressBytes, err := getInstalledUnitAddressBytes(ctx)
	if err != nil {
		return nil, false, errors.New("Unable to get installed uints")
	}

	var units []readWriteModel.UnitAddress
	for i := 0; i <= 0xFF; i++ {
		unitAddressByte := byte(i)
		if _, ok := installedUnitAddressBytes[unitAddressByte]; ok {
			units = append(units, readWriteModel.NewUnitAddress(unitAddressByte))
		}
	}
	return units, true, nil
}

func (m *Browser) extractAttributes(query *unitInfoQuery) ([]readWriteModel.Attribute, bool) {
	if attribute := query.attribute; attribute != nil {
		return []readWriteModel.Attribute{*attribute}, false
	} else {
		var attributes []readWriteModel.Attribute
		for _, attribute := range readWriteModel.AttributeValues {
			attributes = append(attributes, attribute)
		}
		return attributes, true
	}
}

func (m *Browser) getInstalledUnitAddressBytes(ctx context.Context) (map[byte]any, error) {
	start := time.Now()
	defer func() {
		m.log.Debug().TimeDiff("duration", time.Now(), start).Msg("Ending unit address acquiring after duration")
	}()
	// We need to pre-subscribe to catch the 2 followup responses
	subscriptionRequest, err := m.connection.SubscriptionRequestBuilder().
		AddEventTagAddress("installationMMIMonitor", "mmimonitor/*/NETWORK_CONTROL").
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error subscribing to the installation MMI")
	}
	subCtx, subCtxCancel := context.WithTimeout(ctx, 2*time.Second)
	defer subCtxCancel()
	subscriptionResult := <-subscriptionRequest.ExecuteWithContext(subCtx)
	if err := subscriptionResult.GetErr(); err != nil {
		return nil, errors.Wrap(err, "Error subscribing to the mmi")
	}
	response := subscriptionResult.GetResponse()
	if responseCode := response.GetResponseCode("installationMMIMonitor"); responseCode != apiModel.PlcResponseCode_OK {
		return nil, errors.Errorf("Got %s", responseCode)
	}
	subscriptionHandle, err := response.GetSubscriptionHandle("installationMMIMonitor")
	if err != nil {
		return nil, errors.Wrap(err, "Error getting the subscription handle")
	}
	build, err := m.connection.UnsubscriptionRequestBuilder().AddHandles(subscriptionHandle).Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error building unsubscription request")
	}
	defer build.ExecuteWithContext(ctx)

	blockOffset0Received := false
	blockOffset0ReceivedChan := make(chan any, 100) // We only expect one, but we make it a bit bigger to no clog up
	blockOffset88Received := false
	blockOffset88ReceivedChan := make(chan any, 100) // We only expect one, but we make it a bit bigger to no clog up
	blockOffset176Received := false
	blockOffset176ReceivedChan := make(chan any, 100) // We only expect one, but we make it a bit bigger to no clog up
	result := make(map[byte]any)
	plcConsumerRegistration := subscriptionHandle.Register(func(event apiModel.PlcSubscriptionEvent) {
		m.log.Trace().Stringer("event", event).Msg("handling event")
		if responseCode := event.GetResponseCode("installationMMIMonitor"); responseCode != apiModel.PlcResponseCode_OK {
			m.log.Warn().Stringer("event", event).Msg("Ignoring")
			return
		}
		rootValue := event.GetValue("installationMMIMonitor")
		if !rootValue.IsStruct() {
			m.log.Warn().Stringer("rootValue", rootValue).Msg("Ignoring rootValue should be a struct")
			return
		}
		rootStruct := rootValue.GetStruct()
		if applicationValue := rootStruct["application"]; applicationValue == nil || !applicationValue.IsString() || applicationValue.GetString() != "NETWORK_CONTROL" {
			m.log.Warn().
				Interface("rootStruct", rootStruct).
				Msg("Ignoring rootStruct should contain a application tag of type string with value NETWORK_CONTROL")
			return
		}
		var blockStart int
		if blockStartValue := rootStruct["blockStart"]; blockStartValue == nil || !blockStartValue.IsByte() {
			m.log.Warn().
				Interface("rootStruct", rootStruct).
				Msg("Ignoring rootStruct should contain a blockStart tag of type byte")
			return
		} else {
			blockStart = int(blockStartValue.GetByte())
		}

		if plcListValue := rootStruct["values"]; plcListValue == nil || !plcListValue.IsList() {
			m.log.Warn().
				Interface("rootStruct", rootStruct).
				Msg("Ignoring rootStruct should contain a values tag of type list")
			return
		} else {
			for unitByteAddress, plcValue := range plcListValue.GetList() {
				unitByteAddress = blockStart + unitByteAddress
				if !plcValue.IsString() {
					m.log.Warn().
						Stringer("plcValue", plcValue).
						Int("unitByteAddress", unitByteAddress).
						Msg("Ignoring plcValue at unitByteAddress should be a string")
					return
				}
				switch plcValue.GetString() {
				case readWriteModel.GAVState_ON.PLC4XEnumName(), readWriteModel.GAVState_OFF.PLC4XEnumName():
					m.log.Debug().
						Int("unitByteAddress", unitByteAddress).
						Msg("unit does exists")
					result[byte(unitByteAddress)] = true
				case readWriteModel.GAVState_DOES_NOT_EXIST.PLC4XEnumName():
					m.log.Debug().
						Int("unitByteAddress", unitByteAddress).
						Msg("unit does not exists")
				case readWriteModel.GAVState_ERROR.PLC4XEnumName():
					m.log.Warn().
						Int("unitByteAddress", unitByteAddress).
						Msg("unit is in error state")
				}
			}
		}
		// We notify here so we don't exit to early
		switch blockStart {
		case 0:
			select {
			case blockOffset0ReceivedChan <- true:
				m.log.Trace().Msg("0 notified")
			default:
				m.log.Warn().Msg("0 blocked")
			}
		case 88:
			select {
			case blockOffset88ReceivedChan <- true:
				m.log.Trace().Msg("88 notified")
			default:
				m.log.Warn().Msg("88 blocked")
			}
		case 176:
			select {
			case blockOffset176ReceivedChan <- true:
				m.log.Trace().Msg("176 notified")
			default:
				m.log.Warn().Msg("176 blocked")
			}
		}
	})
	defer plcConsumerRegistration.Unregister()

	readRequest, err := m.connection.ReadRequestBuilder().
		AddTagAddress("installationMMI", "status/binary/0xFF").
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error building the installation MMI")
	}
	readCtx, readCtxCancel := context.WithTimeout(ctx, 2*time.Second)
	defer readCtxCancel()
	readWg := sync.WaitGroup{}
	readWg.Add(1)
	go func() {
		defer readWg.Done()
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		defer readCtxCancel()
		m.log.Debug().Stringer("readRequest", readRequest).Msg("sending read request")
		readRequestResult := <-readRequest.ExecuteWithContext(readCtx)
		if err := readRequestResult.GetErr(); err != nil {
			m.log.Warn().Err(err).Msg("Error reading the mmi")
			return
		}
		response := readRequestResult.GetResponse()
		if responseCode := response.GetResponseCode("installationMMI"); responseCode == apiModel.PlcResponseCode_OK {
			rootValue := response.GetValue("installationMMI")
			if !rootValue.IsStruct() {
				m.log.Warn().Err(err).Stringer("rootValue", rootValue).Msg("%v should be a struct")
				return
			}
			rootStruct := rootValue.GetStruct()
			if applicationValue := rootStruct["application"]; applicationValue == nil || !applicationValue.IsString() || applicationValue.GetString() != "NETWORK_CONTROL" {
				m.log.Warn().Err(err).
					Interface("rootStruct", rootStruct).
					Msg("%v should contain a application tag of type string with value NETWORK_CONTROL")
				return
			}
			var blockStart int
			if blockStartValue := rootStruct["blockStart"]; blockStartValue == nil || !blockStartValue.IsByte() || blockStartValue.GetByte() != 0 {
				m.log.Warn().Err(err).
					Interface("rootStruct", rootStruct).
					Msg("rootStruct should contain a blockStart tag of type byte with value 0")
				return
			} else {
				blockStart = int(blockStartValue.GetByte())
			}
			m.log.Debug().Int("blockStart", blockStart).Msg("Read MMI with block start")

			if plcListValue := rootStruct["values"]; plcListValue == nil || !plcListValue.IsList() {
				m.log.Warn().Err(err).
					Interface("rootStruct", rootStruct).
					Msg("rootStruct should contain a values tag of type list")
				return
			} else {
				for unitByteAddress, plcValue := range plcListValue.GetList() {
					unitByteAddress = blockStart + unitByteAddress
					if !plcValue.IsString() {
						m.log.Warn().Err(err).
							Stringer("plcValue", plcValue).
							Int("unitByteAddress", unitByteAddress).
							Msg("plcValue at unitByteAddress should be a string")
						return
					}
					switch plcValue.GetString() {
					case readWriteModel.GAVState_ON.PLC4XEnumName(), readWriteModel.GAVState_OFF.PLC4XEnumName():
						m.log.Debug().
							Int("unitByteAddress", unitByteAddress).
							Msg("unit does exists")
						result[byte(unitByteAddress)] = true
					case readWriteModel.GAVState_DOES_NOT_EXIST.PLC4XEnumName():
						m.log.Debug().
							Int("unitByteAddress", unitByteAddress).
							Msg("unit does not exists")
					case readWriteModel.GAVState_ERROR.PLC4XEnumName():
						m.log.Warn().
							Int("unitByteAddress", unitByteAddress).
							Msg("unit is in error state")
					}
				}
			}
			switch blockStart {
			case 0:
				blockOffset0Received = true
				m.log.Trace().Msg("block 0 read by read")
			case 88:
				blockOffset88Received = true
				m.log.Trace().Msg("block 88 read by read")
			case 176:
				blockOffset176Received = true
				m.log.Trace().Msg("block 176 read by read")
			}

		} else {
			m.log.Warn().
				Stringer("responseCode", responseCode).
				Msg("We got responseCode as response code for installation mmi so we rely on getting it via subscription")
		}
	}()

	syncCtx, syncCtxCancel := context.WithTimeout(ctx, 6*time.Second)
	defer syncCtxCancel()
	for !blockOffset0Received || !blockOffset88Received || !blockOffset176Received {
		select {
		case <-blockOffset0ReceivedChan:
			m.log.Trace().Msg("Offset 0 received")
			blockOffset0Received = true
		case <-blockOffset88ReceivedChan:
			m.log.Trace().Msg("Offset 88 received")
			blockOffset88Received = true
		case <-blockOffset176ReceivedChan:
			m.log.Trace().Msg("Offset 176 received")
			blockOffset176Received = true
		case <-syncCtx.Done():
			err = syncCtx.Err()
			m.log.Trace().Err(err).Msg("Ending prematurely")
			return nil, errors.Wrap(err, "error waiting for other offsets")
		}
	}
	readWg.Wait()
	return result, nil
}
