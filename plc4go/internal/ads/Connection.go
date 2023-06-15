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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/rs/zerolog"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

type Connection struct {
	_default.DefaultConnection

	messageCodec       spi.MessageCodec
	requestInterceptor interceptors.RequestInterceptor
	configuration      model.Configuration
	driverContext      *DriverContext
	tracer             tracer.Tracer

	subscriptions map[uint32]apiModel.PlcSubscriptionHandle

	passLogToModel bool
	log            zerolog.Logger
}

func NewConnection(messageCodec spi.MessageCodec, configuration model.Configuration, connectionOptions map[string][]string, _options ...options.WithOption) (*Connection, error) {
	driverContext, err := NewDriverContext(configuration)
	if err != nil {
		return nil, err
	}
	connection := &Connection{
		messageCodec:   messageCodec,
		configuration:  configuration,
		driverContext:  driverContext,
		subscriptions:  map[uint32]apiModel.PlcSubscriptionHandle{},
		passLogToModel: options.ExtractPassLoggerToModel(_options...),
		log:            options.ExtractCustomLogger(_options...),
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			// TODO: Connection Id is probably "" all the time.
			connection.tracer = tracer.NewTracer(driverContext.connectionId, _options...)
		}
	}
	tagHandler := NewTagHandlerWithDriverContext(driverContext)
	valueHandler := NewValueHandlerWithDriverContext(driverContext, tagHandler, _options...)
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
			_default.WithPlcValueHandler(valueHandler),
		)...,
	)
	return connection, nil
}

func (m *Connection) GetConnectionId() string {
	return m.driverContext.connectionId
}

func (m *Connection) IsTraceEnabled() bool {
	return m.tracer != nil
}

func (m *Connection) GetTracer() tracer.Tracer {
	return m.tracer
}

func (m *Connection) GetConnection() plc4go.PlcConnection {
	return m
}

func (m *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	m.log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)

	// Reset the driver context (Actually this should not be required, but just to be on the safe side)
	m.driverContext.clear()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		err := m.messageCodec.ConnectWithContext(ctx)
		if err != nil {
			ch <- _default.NewDefaultPlcConnectionConnectResult(m, err)
		}

		m.setupConnection(ctx, ch)
	}()
	return ch
}

func (m *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	// First read the device info (Including TwinCat version and PLC name)
	deviceInfoResponse, err := m.ExecuteAdsReadDeviceInfoRequest(ctx)
	if err != nil {
		ch <- _default.NewDefaultPlcConnectionCloseResult(nil, err)
		return
	}
	m.driverContext.adsVersion = fmt.Sprintf("%d.%d.%d", deviceInfoResponse.GetMajorVersion(), deviceInfoResponse.GetMinorVersion(), deviceInfoResponse.GetVersion())
	m.driverContext.deviceName = GetZeroTerminatedString(deviceInfoResponse.GetDevice())

	// Read the symbol-version (offline changes)
	symbolVersionResponse, err := m.ExecuteAdsReadRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_VERSION), 0, 1)
	if err != nil {
		ch <- _default.NewDefaultPlcConnectionCloseResult(nil, err)
		return
	}
	m.driverContext.symbolVersion = symbolVersionResponse.GetData()[0]

	// Read the online-version
	onlineVersionResponse, err := m.ExecuteAdsReadWriteRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_VALBYNAME), 0, 4, nil, []byte("TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt"))
	if err != nil {
		ch <- _default.NewDefaultPlcConnectionCloseResult(nil, err)
		return
	}
	rb := utils.NewReadBufferByteBased(onlineVersionResponse.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	m.driverContext.onlineVersion, err = rb.ReadUint32("", 32)
	if err != nil {
		ch <- _default.NewDefaultPlcConnectionCloseResult(nil, err)
		return
	}

	// Read the data type and symbol table
	err = m.readSymbolTableAndDatatypeTable(ctx)
	if err != nil {
		ch <- _default.NewDefaultPlcConnectionCloseResult(nil, err)
		return
	}

	// Start the worker for handling incoming messages
	// (Messages that are not responses to outgoing messages)
	defaultIncomingMessageChannel := m.messageCodec.GetDefaultIncomingMessageChannel()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			}
		}()
		for message := range defaultIncomingMessageChannel {
			switch message.(type) {
			case readWriteModel.AmsTCPPacket:
				amsTCPPacket := message.(readWriteModel.AmsTCPPacket)
				switch amsTCPPacket.GetUserdata().(type) {
				// Forward all device notification requests to the subscriber component.
				case readWriteModel.AdsDeviceNotificationRequest:
					m.handleIncomingDeviceNotificationRequest(
						amsTCPPacket.GetUserdata().(readWriteModel.AdsDeviceNotificationRequest))
				default:
					m.log.Warn().Msgf("Got unexpected type of incoming ADS message %v", message)
				}
			default:
				m.log.Warn().Msgf("Got unexpected type of incoming ADS message %v", message)
			}
		}
		m.log.Info().Msg("Done waiting for messages ...")
	}()

	// Subscribe for changes to the symbol or the offline-versions
	versionChangeRequest, err := m.SubscriptionRequestBuilder().
		AddChangeOfStateTagAddress("offlineVersion", "0xF008/0x0000:USINT").
		AddPreRegisteredConsumer("offlineVersion", func(event apiModel.PlcSubscriptionEvent) {
			if event.GetResponseCode("offlineVersion") == apiModel.PlcResponseCode_OK {
				newVersion := event.GetValue("offlineVersion").GetUint8()
				if newVersion != m.driverContext.symbolVersion {
					m.log.Info().Msg("detected offline version change: reloading symbol- and data-type-table.")
					err := m.readSymbolTableAndDatatypeTable(ctx)
					if err != nil {
						m.log.Error().Err(err).Msg("error updating data-type and symbol tables")
					}
				}
			}
		}).
		AddChangeOfStateTagAddress("onlineVersion", "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt").
		AddPreRegisteredConsumer("onlineVersion", func(event apiModel.PlcSubscriptionEvent) {
			if event.GetResponseCode("onlineVersion") == apiModel.PlcResponseCode_OK {
				newVersion := event.GetValue("onlineVersion").GetUint32()
				if newVersion != m.driverContext.onlineVersion {
					m.log.Info().Msg("detected online version change: reloading symbol- and data-type-table.")
					err := m.readSymbolTableAndDatatypeTable(ctx)
					if err != nil {
						m.log.Error().Err(err).Msg("error updating data-type and symbol tables")
					}
				}
			}
		}).
		Build()
	subscriptionResultChan := versionChangeRequest.Execute()
	subscriptionRequestResult := <-subscriptionResultChan
	if subscriptionRequestResult.GetErr() != nil {
		ch <- _default.NewDefaultPlcConnectionCloseResult(nil, subscriptionRequestResult.GetErr())
		return
	}

	// Return the finished connection
	ch <- _default.NewDefaultPlcConnectionConnectResult(m, nil)
}

func (m *Connection) readSymbolTableAndDatatypeTable(ctx context.Context) error {
	// First read the sizes of the data type and symbol table, if needed.
	tableSizes, err := m.readDataTypeTableAndSymbolTableSizes(ctx)
	if err != nil {
		return err
	}

	// Then read the data type table, if needed.
	m.driverContext.dataTypeTable, err = m.readDataTypeTable(ctx, tableSizes.GetDataTypeLength(), tableSizes.GetDataTypeCount())
	if err != nil {
		return err
	}

	// Then read the symbol table, if needed.
	m.driverContext.symbolTable, err = m.readSymbolTable(ctx, tableSizes.GetSymbolLength(), tableSizes.GetSymbolCount())
	if err != nil {
		return err
	}
	return nil
}

func (m *Connection) readDataTypeTableAndSymbolTableSizes(ctx context.Context) (readWriteModel.AdsTableSizes, error) {
	response, err := m.ExecuteAdsReadRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYMBOL_AND_DATA_TYPE_SIZES), 0x00000000, 24)
	if err != nil {
		return nil, fmt.Errorf("error reading table: %v", err)
	}

	// Parse and process the response
	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	tableSizes, err := readWriteModel.AdsTableSizesParse(ctxForModel, response.GetData())
	if err != nil {
		return nil, fmt.Errorf("error parsing table: %v", err)
	}
	return tableSizes, nil
}

func (m *Connection) readDataTypeTable(ctx context.Context, dataTableSize uint32, numDataTypes uint32) (map[string]readWriteModel.AdsDataTypeTableEntry, error) {
	response, err := m.ExecuteAdsReadRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_DATA_TYPE_TABLE_UPLOAD), 0x00000000, dataTableSize)
	if err != nil {
		return nil, fmt.Errorf("error reading data-type table: %v", err)
	}

	// Parse and process the response
	readBuffer := utils.NewReadBufferByteBased(response.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	dataTypes := map[string]readWriteModel.AdsDataTypeTableEntry{}
	for i := uint32(0); i < numDataTypes; i++ {
		dataType, err := readWriteModel.AdsDataTypeTableEntryParseWithBuffer(context.Background(), readBuffer)
		if err != nil {
			return nil, fmt.Errorf("error parsing table: %v", err)
		}
		dataTypes[dataType.GetDataTypeName()] = dataType
	}
	return dataTypes, nil
}

func (m *Connection) readSymbolTable(ctx context.Context, symbolTableSize uint32, numSymbols uint32) (map[string]readWriteModel.AdsSymbolTableEntry, error) {
	response, err := m.ExecuteAdsReadRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_UPLOAD), 0x00000000, symbolTableSize)
	if err != nil {
		return nil, fmt.Errorf("error reading data-type table: %v", err)
	}

	// Parse and process the response
	readBuffer := utils.NewReadBufferByteBased(response.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	symbols := map[string]readWriteModel.AdsSymbolTableEntry{}
	for i := uint32(0); i < numSymbols; i++ {
		symbol, err := readWriteModel.AdsSymbolTableEntryParseWithBuffer(context.Background(), readBuffer)
		if err != nil {
			return nil, fmt.Errorf("error parsing table")
		}
		symbols[symbol.GetName()] = symbol
	}
	return symbols, nil
}

func (m *Connection) resolveSymbolicTag(ctx context.Context, symbolicTag model.SymbolicPlcTag) (*model.DirectPlcTag, error) {
	// Find the initial datatype, based on the first to segments.
	symbolicAddress := symbolicTag.SymbolicAddress
	addressParts := strings.Split(symbolicAddress, ".")
	symbolName := ""
	if len(addressParts) == 1 {
		symbolName = addressParts[0]
		addressParts = addressParts[1:]
	} else if len(addressParts) > 1 {
		symbolName = addressParts[0] + "." + addressParts[1]
		addressParts = addressParts[2:]
	} else {
		return nil, errors.New("invalid address")
	}
	symbol, ok := m.driverContext.symbolTable[symbolName]
	if !ok {
		return nil, fmt.Errorf("couldn't find tag with address %s", symbolName)
	}
	dataTypeName := symbol.GetDataTypeName()
	dataType, ok := m.driverContext.dataTypeTable[dataTypeName]
	if !ok {
		return nil, fmt.Errorf("couldn't find data type with name %s for tag with address %s", dataTypeName, symbolName)
	}
	// Start resolving the address.
	return m.resolveSymbolicAddress(ctx, addressParts, dataType, symbol.GetGroup(), symbol.GetOffset())
}

func (m *Connection) resolveSymbolicAddress(ctx context.Context, addressParts []string, curDataType readWriteModel.AdsDataTypeTableEntry, indexGroup uint32, indexOffset uint32) (*model.DirectPlcTag, error) {
	// If we've reached then end of the resolution, return the final entry.
	if len(addressParts) == 0 {
		var arrayInfo []apiModel.ArrayInfo
		for _, adsArrayInfo := range curDataType.GetArrayInfo() {
			arrayInfo = append(arrayInfo, &spiModel.DefaultArrayInfo{
				LowerBound: adsArrayInfo.GetLowerBound(),
				UpperBound: adsArrayInfo.GetUpperBound(),
			})
		}
		plcValueType, stringLength := m.getPlcValueForAdsDataTypeTableEntry(curDataType)
		return &model.DirectPlcTag{
			PlcTag: model.PlcTag{
				ArrayInfo: arrayInfo,
			},
			IndexGroup:   indexGroup,
			IndexOffset:  indexOffset,
			ValueType:    plcValueType,
			StringLength: stringLength,
			DataType:     curDataType,
		}, nil
	}

	// Resolve the next level of the address.
	curAddressPart := addressParts[0]
	restAddressParts := addressParts[1:]
	for _, child := range curDataType.GetChildren() {
		if child.GetPropertyName() == curAddressPart {
			childDataTypeName := child.GetDataTypeName()
			childDataType, ok := m.driverContext.dataTypeTable[childDataTypeName]
			if !ok {
				return nil, fmt.Errorf("couldn't find data type %s for property %s of data type %s",
					childDataTypeName, curAddressPart, curDataType.GetDataTypeName())
			}
			return m.resolveSymbolicAddress(ctx, restAddressParts, childDataType, indexGroup, indexOffset+child.GetOffset())
		}
	}
	return nil, fmt.Errorf("couldn't find property named %s for data type %s",
		curAddressPart, curDataType.GetDataTypeName())
}

func (m *Connection) getPlcValueForAdsDataTypeTableEntry(entry readWriteModel.AdsDataTypeTableEntry) (apiValues.PlcValueType, int32) {
	stringLength := -1
	dataTypeName := entry.GetDataTypeName()
	if strings.HasPrefix(dataTypeName, "STRING(") {
		var err error
		stringLength, err = strconv.Atoi(dataTypeName[7 : len(dataTypeName)-1])
		if err != nil {
			return apiValues.NULL, -1
		}
		dataTypeName = "STRING"
	} else if strings.HasPrefix(dataTypeName, "WSTRING(") {
		var err error
		stringLength, err = strconv.Atoi(dataTypeName[8 : len(dataTypeName)-1])
		if err != nil {
			return apiValues.NULL, -1
		}
		dataTypeName = "WSTRING"
	}
	plcValueType, ok := apiValues.PlcValueByName(dataTypeName)
	if !ok {
		return apiValues.NULL, -1
	}
	return plcValueType, int32(stringLength)
}

func (m *Connection) GetMessageCodec() spi.MessageCodec {
	return m.messageCodec
}

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    true,
	}
}

func (m *Connection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m *Connection) String() string {
	return fmt.Sprintf("ads.Connection{}")
}

func GetZeroTerminatedString(data []byte) string {
	for i := range data {
		if data[i] == 0x00 {
			return string(data[0:i])
		}
	}
	return ""
}
