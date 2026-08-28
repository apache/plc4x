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
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Connection struct {
	_default.DefaultConnection

	messageCodec       spi.MessageCodec
	requestInterceptor interceptors.RequestInterceptor
	configuration      model.Configuration
	driverContext      *DriverContext
	tracer             tracer.Tracer

	subscriptions map[uint32]apiModel.PlcSubscriptionHandle

	wg sync.WaitGroup // use to track spawned go routines

	passLogToModel bool
	log            zerolog.Logger
	_options       []options.WithOption // Used to pass them downstream

	// Nesting budget for parsing the data-type table uploaded from the device. An entry may
	// contain further entries, so without a budget the depth of the tree is the device's choice.
	// Mirrors the Java driver's "max-data-type-table-depth" connection option.
	maxDataTypeTableDepth uint16
}

// defaultMaxDataTypeTableDepth matches the default of the Java driver's
// "max-data-type-table-depth" option. Real type hierarchies are only a handful of levels deep.
const defaultMaxDataTypeTableDepth uint16 = 20

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(messageCodec spi.MessageCodec, configuration model.Configuration, connectionOptions map[string][]string, _options ...options.WithOption) (*Connection, error) {
	driverContext, err := NewDriverContext(configuration)
	if err != nil {
		return nil, err
	}
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		messageCodec:   messageCodec,
		configuration:  configuration,
		driverContext:  driverContext,
		subscriptions:  map[uint32]apiModel.PlcSubscriptionHandle{},
		passLogToModel: passLoggerToModel,
		log:            customLogger,
		_options:       _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			// TODO: Connection Id is probably "" all the time.
			connection.tracer = tracer.NewTracer(driverContext.connectionId, _options...)
		}
	}
	connection.maxDataTypeTableDepth = defaultMaxDataTypeTableDepth
	if depthOption, ok := connectionOptions["max-data-type-table-depth"]; ok && len(depthOption) == 1 {
		depth, err := strconv.ParseUint(depthOption[0], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid max-data-type-table-depth %q: %v", depthOption[0], err)
		}
		connection.maxDataTypeTableDepth = uint16(depth)
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

func (m *Connection) Connect(ctx context.Context) error {
	m.log.Trace().Msg("Connecting")

	// Reset the driver context (Actually this should not be required, but just to be on the safe side)
	m.driverContext.clear()
	if err := m.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting to message codec")
	}

	// For testing purposes we can skip the waiting for a complete connection
	if !m.driverContext.awaitSetupComplete {
		m.wg.Go(func() {
			if err := m.setupConnection(ctx); err != nil {
				m.log.Error().Err(err).Msg("Error during connection setup")
			}
		})
		m.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		m.SetConnected(true)
		return nil
	}

	if err := m.setupConnection(ctx); err != nil {
		return errors.Wrap(err, "error setting up connection")
	}
	return nil
}

func (m *Connection) setupConnection(ctx context.Context) error {
	// First read the device info (Including TwinCat version and PLC name)
	deviceInfoResponse, err := m.ExecuteAdsReadDeviceInfoRequest(ctx)
	if err != nil {
		return errors.Wrap(err, "error reading device info")
	}
	m.driverContext.adsVersion = fmt.Sprintf("%d.%d.%d", deviceInfoResponse.GetMajorVersion(), deviceInfoResponse.GetMinorVersion(), deviceInfoResponse.GetVersion())
	m.driverContext.deviceName = GetZeroTerminatedString(deviceInfoResponse.GetDevice())

	// Read the online-version
	// (The order online- before symbol-version matches the Java driver and the shared driver testsuite.)
	onlineVersionResponse, err := m.ExecuteAdsReadWriteRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_VALBYNAME), 0, 4, nil, []byte("TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt"))
	if err != nil {
		return errors.Wrap(err, "error reading online version")
	}
	rb := utils.NewReadBufferByteBased(onlineVersionResponse.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	m.driverContext.onlineVersion, err = rb.ReadUint32("", 32)
	if err != nil {
		return errors.Wrap(err, "error reading online version")
	}

	// Read the symbol-version (offline changes)
	symbolVersionResponse, err := m.ExecuteAdsReadRequest(ctx, uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_VERSION), 0, 1)
	if err != nil {
		return errors.Wrap(err, "error reading symbol version")
	}
	// The length of the response data is wire-controlled, so it must be checked
	// before indexing to avoid a panic on an empty response.
	if len(symbolVersionResponse.GetData()) < 1 {
		return errors.New("error reading symbol version: empty response data")
	}
	m.driverContext.symbolVersion = symbolVersionResponse.GetData()[0]

	// Read the data type and symbol table
	err = m.readSymbolTableAndDatatypeTable(ctx)
	if err != nil {
		return errors.Wrap(err, "error reading symbol table and data type table")
	}

	// Start the worker for handling incoming messages
	// (Messages that are not responses to outgoing messages)
	defaultIncomingMessageChannel := m.messageCodec.GetDefaultIncomingMessageChannel()
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
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
					m.log.Warn().Interface("message", message).Msg("Got unexpected type of incoming ADS message")
				}
			default:
				m.log.Warn().Interface("message", message).Msg("Got unexpected type of incoming ADS message")
			}
		}
		m.log.Info().Msg("Done waiting for messages ...")
	})

	// Subscribe for changes to the symbol or the offline-versions
	// (Cyclic with a 1s check interval: the wire request is an ON_CHANGE device notification
	// whose cycle time is the interval, matching the Java driver's Duration.ofMillis(1000).
	// The online- before offline-version order matches the Java driver and the driver testsuite.)
	versionChangeRequest, err := m.SubscriptionRequestBuilder().
		AddCyclicTagAddress("onlineVersion", "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt", time.Second).
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
		AddCyclicTagAddress("offlineVersion", "0xF008/0x0000:USINT", time.Second).
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
		Build()
	if err != nil {
		return errors.Wrap(err, "error building subscription request")
	}
	subscriptionResultChan := versionChangeRequest.Execute(ctx)
	subscriptionRequestResult := <-subscriptionResultChan
	if err := subscriptionRequestResult.GetErr(); err != nil {
		return errors.Wrap(err, "error subscribing to version change")
	}

	return nil
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
	for range numDataTypes {
		dataType, err := readWriteModel.AdsDataTypeTableEntryParseWithBuffer(ctx, readBuffer, m.maxDataTypeTableDepth)
		if err != nil {
			return nil, fmt.Errorf("error parsing table: %v", err)
		}
		// Key by the main name: that's the name symbols reference via their dataTypeName
		// (and what the Java driver keys by); the secondary name is the aliased/simple type.
		dataTypes[dataType.GetMainName()] = dataType
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
	for range numSymbols {
		symbol, err := readWriteModel.AdsSymbolTableEntryParseWithBuffer(ctx, readBuffer)
		if err != nil {
			return nil, fmt.Errorf("error parsing table")
		}
		symbols[symbol.GetName()] = symbol
	}
	return symbols, nil
}

// directTagFor returns a request tag as a fully usable direct tag: symbolic tags are
// resolved against the tables loaded during connection setup, and direct tags created
// straight from an address string get their data type table entry filled in.
func (m *Connection) directTagFor(ctx context.Context, tag apiModel.PlcTag) (*model.DirectPlcTag, error) {
	var directTag model.DirectPlcTag
	switch typedTag := tag.(type) {
	case model.SymbolicPlcTag:
		resolvedTag, err := m.resolveSymbolicTag(ctx, typedTag)
		if err != nil {
			return nil, errors.Wrap(err, "error resolving symbolic tag")
		}
		return resolvedTag, nil
	case *model.SymbolicPlcTag:
		resolvedTag, err := m.resolveSymbolicTag(ctx, *typedTag)
		if err != nil {
			return nil, errors.Wrap(err, "error resolving symbolic tag")
		}
		return resolvedTag, nil
	case model.DirectPlcTag:
		directTag = typedTag
	case *model.DirectPlcTag:
		directTag = *typedTag
	default:
		return nil, errors.Errorf("invalid tag type %T", tag)
	}
	if directTag.DataType == nil {
		dataType, ok := m.driverContext.dataTypeTable[directTag.ValueType.String()]
		if !ok {
			return nil, errors.Errorf("no entry for data type %s in the data type table", directTag.ValueType)
		}
		directTag.DataType = dataType
	}
	// A direct address names a location whose type is the element type, so its selection is the
	// whole of its shape. The offset was already moved to the first selected element while the
	// address was parsed (see TagHandler.applySelectionOffset); what is added here is the shape
	// to decode, without which the extra elements were transferred and then dropped.
	if len(directTag.SelectedArrayInfo) == 0 && len(directTag.ArrayInfo) > 0 {
		elements := uint32(1)
		for _, dimension := range directTag.ArrayInfo {
			elements *= dimension.GetSize()
		}
		directTag.SelectedArrayInfo = []readWriteModel.AdsDataTypeArrayInfo{
			readWriteModel.NewAdsDataTypeArrayInfo(uint32(directTag.ArrayInfo[0].GetLowerBound()), elements),
		}
		directTag.SelectedSizeInBytes = directTag.DataType.GetSize() * elements
	}
	return &directTag, nil
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
	resolved, err := m.resolveSymbolicAddress(ctx, addressParts, dataType, symbol.GetGroup(), symbol.GetOffset())
	if err != nil {
		return nil, err
	}
	// The trailing selection is not part of the symbolic path - it says which elements of the
	// resolved location to read. Applying it here is what makes MAIN.arr[1..4] read those four
	// elements; without it the resolved location was the whole array at its original offset, and
	// the selection was parsed, rendered and silently discarded.
	return m.applySymbolicSelection(resolved, symbolicTag.ArrayInfo, symbolicAddress)
}

// applySymbolicSelection narrows a resolved location to the elements the address selected.
//
// A selection is refused rather than approximated when it cannot be applied exactly: reading the
// wrong elements is indistinguishable from reading the right ones once the values come back.
func (m *Connection) applySymbolicSelection(tag *model.DirectPlcTag, selection []apiModel.ArrayInfo, address string) (*model.DirectPlcTag, error) {
	if len(selection) == 0 {
		return tag, nil
	}
	declared := tag.ArrayInfo
	if len(declared) == 0 {
		return nil, fmt.Errorf("address %s selects elements, but the PLC does not declare %s as an array",
			address, tag.DataType.GetMainName())
	}
	if len(selection) != len(declared) {
		return nil, fmt.Errorf("address %s selects %d dimension(s) of %s, which the PLC declares "+
			"with %d - name every dimension or none of them",
			address, len(selection), tag.DataType.GetMainName(), len(declared))
	}
	itemType, err := m.arrayItemTypeFor(tag.DataType)
	if err != nil {
		return nil, err
	}
	itemSize := itemType.GetSize()
	if itemSize == 0 {
		return nil, fmt.Errorf("address %s selects elements of %s, whose size the data type table "+
			"gives as zero", address, tag.DataType.GetMainName())
	}
	for dimension := range selection {
		if err := checkSelectedDimension(selection, declared, dimension, address); err != nil {
			return nil, err
		}
	}

	// Row-major strides: a step along a dimension skips every element of the dimensions inside
	// it. The innermost stride is one element, and each dimension outwards multiplies by the
	// number of elements the dimension inside it declares.
	stride := itemSize
	elements := uint32(1)
	var shape []readWriteModel.AdsDataTypeArrayInfo
	for dimension := len(selection) - 1; dimension >= 0; dimension-- {
		tag.IndexOffset += (selection[dimension].GetLowerBound() - declared[dimension].GetLowerBound()) * stride
		stride *= declared[dimension].GetSize()
		elements *= selection[dimension].GetSize()
		if selection[dimension].IsRange() {
			// A range is an array even when it spans one element, so the shape follows what the
			// address wrote rather than the count. A bare index contributes no level - it moves
			// the start and collapses, which is what makes it a scalar.
			shape = append([]readWriteModel.AdsDataTypeArrayInfo{
				readWriteModel.NewAdsDataTypeArrayInfo(
					selection[dimension].GetLowerBound(), selection[dimension].GetSize()),
			}, shape...)
		}
	}

	tag.DataType = itemType
	tag.ValueType, tag.StringLength = m.getPlcValueForAdsDataTypeTableEntry(itemType)
	tag.SelectedSizeInBytes = itemSize * elements
	tag.SelectedArrayInfo = shape
	if len(shape) > 0 {
		tag.ArrayInfo = selection
	} else {
		tag.ArrayInfo = nil
	}
	return tag, nil
}

// checkSelectedDimension holds one dimension of a selection to what the PLC declares, and to what
// a single read can express.
//
// One read covers one contiguous run of memory. Scanning outwards from the innermost dimension,
// that means every dimension inside a dimension selecting more than one element must be selected
// whole: on an ARRAY [0..9,0..4], "[0..9,1..3]" names ten separate three-element runs, and the
// contiguous block of thirty elements starting at [0,1] that a single read would return is not
// what was asked for. Refusing says so; returning that block would not.
func checkSelectedDimension(selection, declared []apiModel.ArrayInfo, dimension int, address string) error {
	selected, available := selection[dimension], declared[dimension]
	if selected.GetLowerBound() < available.GetLowerBound() || selected.GetUpperBound() > available.GetUpperBound() {
		return fmt.Errorf("address %s selects [%d..%d] of a dimension the PLC declares as [%d..%d]",
			address, selected.GetLowerBound(), selected.GetUpperBound(),
			available.GetLowerBound(), available.GetUpperBound())
	}
	if dimension == 0 || selected.GetSize() == available.GetSize() {
		return nil
	}
	for outer := 0; outer < dimension; outer++ {
		if selection[outer].GetSize() > 1 {
			return fmt.Errorf("address %s selects part of dimension %d while dimension %d spans "+
				"%d elements, which is not one contiguous read - select the whole of the inner "+
				"dimension, or one element of the outer one",
				address, dimension, outer, selection[outer].GetSize())
		}
	}
	return nil
}

func (m *Connection) resolveSymbolicAddress(ctx context.Context, addressParts []string, curDataType readWriteModel.AdsDataTypeTableEntry, indexGroup uint32, indexOffset uint32) (*model.DirectPlcTag, error) {
	// If we've reached then end of the resolution, return the final entry.
	if len(addressParts) == 0 {
		// The dimensions the symbol table declares. They are arrays by definition - the device
		// says so - which is what Range records: without it the shape rule would read them as a
		// bare index and report the whole array as a scalar. The declared lower bound is also
		// the base, so an address written with the PLC's own indices lines up with it.
		var arrayInfo []apiModel.ArrayInfo
		for _, adsArrayInfo := range curDataType.GetArrayInfo() {
			arrayInfo = append(arrayInfo, &spiModel.DefaultArrayInfo{
				LowerBound: adsArrayInfo.GetLowerBound(),
				UpperBound: adsArrayInfo.GetUpperBound(),
				Base:       adsArrayInfo.GetLowerBound(),
				Range:      true,
			})
		}
		plcValueType, stringLength := m.getPlcValueForAdsDataTypeTableEntry(curDataType)
		return &model.DirectPlcTag{
			ArrayInfo:    arrayInfo,
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
		if child.GetMainName() == curAddressPart {
			childDataTypeName := child.GetSecondaryName()
			childDataType, ok := m.driverContext.dataTypeTable[childDataTypeName]
			if !ok {
				return nil, fmt.Errorf("couldn't find data type %s for property %s of data type %s",
					childDataTypeName, curAddressPart, curDataType.GetSecondaryName())
			}
			return m.resolveSymbolicAddress(ctx, restAddressParts, childDataType, indexGroup, indexOffset+child.GetOffset())
		}
	}
	return nil, fmt.Errorf("couldn't find property named %s for data type %s",
		curAddressPart, curDataType.GetSecondaryName())
}

func (m *Connection) getPlcValueForAdsDataTypeTableEntry(entry readWriteModel.AdsDataTypeTableEntry) (apiValues.PlcValueType, int32) {
	stringLength := -1
	// The main name carries the type's name ("BYTE", "STRING(80)", ...); the secondary
	// name is only set for aliased types (matching the Java driver's use of getMainName()).
	dataTypeName := entry.GetMainName()
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
	plcValueType, ok := apiValues.PlcValueTypeByName(dataTypeName)
	if !ok {
		return apiValues.NULL, -1
	}
	return plcValueType, int32(stringLength)
}

func (m *Connection) GetMessageCodec() spi.MessageCodec {
	return m.messageCodec
}

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
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
