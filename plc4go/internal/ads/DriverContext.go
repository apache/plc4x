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
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync/atomic"

	model3 "github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	model2 "github.com/apache/plc4x/plc4go/spi/model"
)

type DriverContext struct {
	connectionId            string
	invokeId                uint32
	adsVersion              string
	deviceName              string
	symbolVersion           uint8
	onlineVersion           uint32
	dataTypeTable           map[string]driverModel.AdsDataTypeTableEntry
	symbolTable             map[string]driverModel.AdsSymbolTableEntry
	awaitSetupComplete      bool
	awaitDisconnectComplete bool
}

func NewDriverContext(configuration model3.Configuration) (*DriverContext, error) {
	return &DriverContext{
		invokeId: 0,
	}, nil
}

func (m *DriverContext) clear() {
	m.connectionId = ""
	m.invokeId = 0
	m.adsVersion = ""
	m.deviceName = ""
	m.symbolVersion = 0
	m.onlineVersion = 0
	m.dataTypeTable = map[string]driverModel.AdsDataTypeTableEntry{}
	m.symbolTable = map[string]driverModel.AdsSymbolTableEntry{}
	m.awaitSetupComplete = false
	m.awaitDisconnectComplete = false
}

func (m *DriverContext) getDirectTagForSymbolTag(symbolicPlcTag model3.SymbolicPlcTag) (*model3.DirectPlcTag, error) {
	address := symbolicPlcTag.SymbolicAddress
	addressSegments := strings.Split(address, ".")
	var symbolName string
	var remainingSegments []string
	if len(addressSegments) == 1 {
		symbolName = addressSegments[0]
		remainingSegments = []string{}
	} else if len(addressSegments) > 1 {
		symbolName = addressSegments[0] + "." + addressSegments[1]
		remainingSegments = addressSegments[2:]
	} else {
		return nil, fmt.Errorf("can't resolve empty address")
	}

	symbolEntry, ok := m.symbolTable[symbolName]
	if !ok {
		return nil, fmt.Errorf("couldn't find symbol with name %s", symbolName)
	}
	dataTypeEntry, ok := m.dataTypeTable[symbolEntry.GetDataTypeName()]
	if !ok {
		return nil, fmt.Errorf("couldn't find data type with name %s for symbol %s", symbolEntry.GetDataTypeName(), symbolName)
	}
	return m.resolveDirectTag(remainingSegments, dataTypeEntry, symbolEntry.GetGroup(), symbolEntry.GetOffset())
}

func (m *DriverContext) resolveDirectTag(remainingSegments []string, currentDatatype driverModel.AdsDataTypeTableEntry, indexGroup uint32, indexOffset uint32) (*model3.DirectPlcTag, error) {
	if len(remainingSegments) == 0 {
		return &model3.DirectPlcTag{
			IndexGroup:   indexGroup,
			IndexOffset:  indexOffset,
			ValueType:    m.getDataTypeForDataTypeTableEntry(currentDatatype),
			StringLength: m.getStringLengthForDataTypeTableEntry(currentDatatype),
			DataType:     currentDatatype,
			PlcTag: model3.PlcTag{
				ArrayInfo: m.getArrayInfoForDataTypeTableEntry(currentDatatype),
			},
		}, nil
	}

	currentSegment := remainingSegments[0]
	remainingSegments = remainingSegments[1:]
	for _, child := range currentDatatype.GetChildren() {
		if child.GetPropertyName() == currentSegment {
			childDataType, ok := m.dataTypeTable[child.GetDataTypeName()]
			if !ok {
				return nil, fmt.Errorf("couldn't find data type with name %s", child.GetDataTypeName())
			}
			return m.resolveDirectTag(remainingSegments, childDataType, indexGroup, indexOffset+child.GetOffset())
		}
	}
	return nil, fmt.Errorf("couldn't find child with name %s in type %s", currentSegment, currentDatatype.GetDataTypeName())
}

func (m *DriverContext) getDataTypeForDataTypeTableEntry(entry driverModel.AdsDataTypeTableEntry) values.PlcValueType {
	if entry.GetArrayInfo() != nil && len(entry.GetArrayInfo()) > 0 {
		return values.List
	}
	if entry.GetNumChildren() > 0 {
		return values.Struct
	}
	dataTypeName := entry.GetDataTypeName()
	if strings.HasPrefix(dataTypeName, "STRING(") {
		dataTypeName = "STRING"
	} else if strings.HasPrefix(dataTypeName, "WSTRING(") {
		dataTypeName = "WSTRING"
	}
	plcValueType, _ := values.PlcValueByName(dataTypeName)
	return plcValueType
}

func (m *DriverContext) getStringLengthForDataTypeTableEntry(entry driverModel.AdsDataTypeTableEntry) int32 {
	dataTypeName := entry.GetDataTypeName()
	if strings.HasPrefix(dataTypeName, "STRING(") {
		lenStr := dataTypeName[7 : len(dataTypeName)-1]
		lenVal, err := strconv.Atoi(lenStr)
		if err != nil {
			return -1
		}
		return int32(lenVal)
	} else if strings.HasPrefix(dataTypeName, "WSTRING(") {
		lenStr := dataTypeName[8 : len(dataTypeName)-1]
		lenVal, err := strconv.Atoi(lenStr)
		if err != nil {
			return -1
		}
		return int32(lenVal)
	}
	return 0
}

func (m *DriverContext) getArrayInfoForDataTypeTableEntry(entry driverModel.AdsDataTypeTableEntry) []model.ArrayInfo {
	var arrayInfos []model.ArrayInfo
	for _, adsArrayInfo := range entry.GetArrayInfo() {
		arrayInfo := model2.DefaultArrayInfo{
			LowerBound: adsArrayInfo.GetLowerBound(),
			UpperBound: adsArrayInfo.GetUpperBound(),
		}
		arrayInfos = append(arrayInfos, arrayInfo)
	}
	return arrayInfos
}

func (m *DriverContext) getInvokeId() uint32 {
	// Calculate a new transaction identifier
	transactionIdentifier := atomic.AddUint32(&m.invokeId, 1)
	if transactionIdentifier > math.MaxUint8 {
		transactionIdentifier = 1
		atomic.StoreUint32(&m.invokeId, 1)
	}
	return transactionIdentifier
}
