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

package model

import (
	"context"
	"fmt"

	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Code generated by code-generation. DO NOT EDIT.

// OpcuaNodeIdServicesVariableY is an enum
type OpcuaNodeIdServicesVariableY int32

type IOpcuaNodeIdServicesVariableY interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
}

const (
	OpcuaNodeIdServicesVariableY_YArrayItemType_Definition       OpcuaNodeIdServicesVariableY = 12030
	OpcuaNodeIdServicesVariableY_YArrayItemType_ValuePrecision   OpcuaNodeIdServicesVariableY = 12031
	OpcuaNodeIdServicesVariableY_YArrayItemType_InstrumentRange  OpcuaNodeIdServicesVariableY = 12032
	OpcuaNodeIdServicesVariableY_YArrayItemType_EURange          OpcuaNodeIdServicesVariableY = 12033
	OpcuaNodeIdServicesVariableY_YArrayItemType_EngineeringUnits OpcuaNodeIdServicesVariableY = 12034
	OpcuaNodeIdServicesVariableY_YArrayItemType_Title            OpcuaNodeIdServicesVariableY = 12035
	OpcuaNodeIdServicesVariableY_YArrayItemType_AxisScaleType    OpcuaNodeIdServicesVariableY = 12036
	OpcuaNodeIdServicesVariableY_YArrayItemType_XAxisDefinition  OpcuaNodeIdServicesVariableY = 12037
)

var OpcuaNodeIdServicesVariableYValues []OpcuaNodeIdServicesVariableY

func init() {
	_ = errors.New
	OpcuaNodeIdServicesVariableYValues = []OpcuaNodeIdServicesVariableY{
		OpcuaNodeIdServicesVariableY_YArrayItemType_Definition,
		OpcuaNodeIdServicesVariableY_YArrayItemType_ValuePrecision,
		OpcuaNodeIdServicesVariableY_YArrayItemType_InstrumentRange,
		OpcuaNodeIdServicesVariableY_YArrayItemType_EURange,
		OpcuaNodeIdServicesVariableY_YArrayItemType_EngineeringUnits,
		OpcuaNodeIdServicesVariableY_YArrayItemType_Title,
		OpcuaNodeIdServicesVariableY_YArrayItemType_AxisScaleType,
		OpcuaNodeIdServicesVariableY_YArrayItemType_XAxisDefinition,
	}
}

func OpcuaNodeIdServicesVariableYByValue(value int32) (enum OpcuaNodeIdServicesVariableY, ok bool) {
	switch value {
	case 12030:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_Definition, true
	case 12031:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_ValuePrecision, true
	case 12032:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_InstrumentRange, true
	case 12033:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_EURange, true
	case 12034:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_EngineeringUnits, true
	case 12035:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_Title, true
	case 12036:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_AxisScaleType, true
	case 12037:
		return OpcuaNodeIdServicesVariableY_YArrayItemType_XAxisDefinition, true
	}
	return 0, false
}

func OpcuaNodeIdServicesVariableYByName(value string) (enum OpcuaNodeIdServicesVariableY, ok bool) {
	switch value {
	case "YArrayItemType_Definition":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_Definition, true
	case "YArrayItemType_ValuePrecision":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_ValuePrecision, true
	case "YArrayItemType_InstrumentRange":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_InstrumentRange, true
	case "YArrayItemType_EURange":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_EURange, true
	case "YArrayItemType_EngineeringUnits":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_EngineeringUnits, true
	case "YArrayItemType_Title":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_Title, true
	case "YArrayItemType_AxisScaleType":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_AxisScaleType, true
	case "YArrayItemType_XAxisDefinition":
		return OpcuaNodeIdServicesVariableY_YArrayItemType_XAxisDefinition, true
	}
	return 0, false
}

func OpcuaNodeIdServicesVariableYKnows(value int32) bool {
	for _, typeValue := range OpcuaNodeIdServicesVariableYValues {
		if int32(typeValue) == value {
			return true
		}
	}
	return false
}

func CastOpcuaNodeIdServicesVariableY(structType any) OpcuaNodeIdServicesVariableY {
	castFunc := func(typ any) OpcuaNodeIdServicesVariableY {
		if sOpcuaNodeIdServicesVariableY, ok := typ.(OpcuaNodeIdServicesVariableY); ok {
			return sOpcuaNodeIdServicesVariableY
		}
		return 0
	}
	return castFunc(structType)
}

func (m OpcuaNodeIdServicesVariableY) GetLengthInBits(ctx context.Context) uint16 {
	return 32
}

func (m OpcuaNodeIdServicesVariableY) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func OpcuaNodeIdServicesVariableYParse(ctx context.Context, theBytes []byte) (OpcuaNodeIdServicesVariableY, error) {
	return OpcuaNodeIdServicesVariableYParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func OpcuaNodeIdServicesVariableYParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (OpcuaNodeIdServicesVariableY, error) {
	log := zerolog.Ctx(ctx)
	_ = log
	val, err := readBuffer.ReadInt32("OpcuaNodeIdServicesVariableY", 32)
	if err != nil {
		return 0, errors.Wrap(err, "error reading OpcuaNodeIdServicesVariableY")
	}
	if enum, ok := OpcuaNodeIdServicesVariableYByValue(val); !ok {
		log.Debug().Interface("val", val).Msg("no value val found for OpcuaNodeIdServicesVariableY")
		return OpcuaNodeIdServicesVariableY(val), nil
	} else {
		return enum, nil
	}
}

func (e OpcuaNodeIdServicesVariableY) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e OpcuaNodeIdServicesVariableY) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	log := zerolog.Ctx(ctx)
	_ = log
	return writeBuffer.WriteInt32("OpcuaNodeIdServicesVariableY", 32, int32(e), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e OpcuaNodeIdServicesVariableY) PLC4XEnumName() string {
	switch e {
	case OpcuaNodeIdServicesVariableY_YArrayItemType_Definition:
		return "YArrayItemType_Definition"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_ValuePrecision:
		return "YArrayItemType_ValuePrecision"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_InstrumentRange:
		return "YArrayItemType_InstrumentRange"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_EURange:
		return "YArrayItemType_EURange"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_EngineeringUnits:
		return "YArrayItemType_EngineeringUnits"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_Title:
		return "YArrayItemType_Title"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_AxisScaleType:
		return "YArrayItemType_AxisScaleType"
	case OpcuaNodeIdServicesVariableY_YArrayItemType_XAxisDefinition:
		return "YArrayItemType_XAxisDefinition"
	}
	return fmt.Sprintf("Unknown(%v)", int32(e))
}

func (e OpcuaNodeIdServicesVariableY) String() string {
	return e.PLC4XEnumName()
}