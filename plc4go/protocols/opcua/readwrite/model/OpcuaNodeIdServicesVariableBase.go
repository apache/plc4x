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

// OpcuaNodeIdServicesVariableBase is an enum
type OpcuaNodeIdServicesVariableBase int32

type IOpcuaNodeIdServicesVariableBase interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
}

const (
	OpcuaNodeIdServicesVariableBase_BaseEventType_EventId               OpcuaNodeIdServicesVariableBase = 2042
	OpcuaNodeIdServicesVariableBase_BaseEventType_EventType             OpcuaNodeIdServicesVariableBase = 2043
	OpcuaNodeIdServicesVariableBase_BaseEventType_SourceNode            OpcuaNodeIdServicesVariableBase = 2044
	OpcuaNodeIdServicesVariableBase_BaseEventType_SourceName            OpcuaNodeIdServicesVariableBase = 2045
	OpcuaNodeIdServicesVariableBase_BaseEventType_Time                  OpcuaNodeIdServicesVariableBase = 2046
	OpcuaNodeIdServicesVariableBase_BaseEventType_ReceiveTime           OpcuaNodeIdServicesVariableBase = 2047
	OpcuaNodeIdServicesVariableBase_BaseEventType_Message               OpcuaNodeIdServicesVariableBase = 2050
	OpcuaNodeIdServicesVariableBase_BaseEventType_Severity              OpcuaNodeIdServicesVariableBase = 2051
	OpcuaNodeIdServicesVariableBase_BaseEventType_LocalTime             OpcuaNodeIdServicesVariableBase = 3190
	OpcuaNodeIdServicesVariableBase_BaseAnalogType_InstrumentRange      OpcuaNodeIdServicesVariableBase = 17567
	OpcuaNodeIdServicesVariableBase_BaseAnalogType_EURange              OpcuaNodeIdServicesVariableBase = 17568
	OpcuaNodeIdServicesVariableBase_BaseAnalogType_EngineeringUnits     OpcuaNodeIdServicesVariableBase = 17569
	OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassId      OpcuaNodeIdServicesVariableBase = 31771
	OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassName    OpcuaNodeIdServicesVariableBase = 31772
	OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassId   OpcuaNodeIdServicesVariableBase = 31773
	OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassName OpcuaNodeIdServicesVariableBase = 31774
)

var OpcuaNodeIdServicesVariableBaseValues []OpcuaNodeIdServicesVariableBase

func init() {
	_ = errors.New
	OpcuaNodeIdServicesVariableBaseValues = []OpcuaNodeIdServicesVariableBase{
		OpcuaNodeIdServicesVariableBase_BaseEventType_EventId,
		OpcuaNodeIdServicesVariableBase_BaseEventType_EventType,
		OpcuaNodeIdServicesVariableBase_BaseEventType_SourceNode,
		OpcuaNodeIdServicesVariableBase_BaseEventType_SourceName,
		OpcuaNodeIdServicesVariableBase_BaseEventType_Time,
		OpcuaNodeIdServicesVariableBase_BaseEventType_ReceiveTime,
		OpcuaNodeIdServicesVariableBase_BaseEventType_Message,
		OpcuaNodeIdServicesVariableBase_BaseEventType_Severity,
		OpcuaNodeIdServicesVariableBase_BaseEventType_LocalTime,
		OpcuaNodeIdServicesVariableBase_BaseAnalogType_InstrumentRange,
		OpcuaNodeIdServicesVariableBase_BaseAnalogType_EURange,
		OpcuaNodeIdServicesVariableBase_BaseAnalogType_EngineeringUnits,
		OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassId,
		OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassName,
		OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassId,
		OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassName,
	}
}

func OpcuaNodeIdServicesVariableBaseByValue(value int32) (enum OpcuaNodeIdServicesVariableBase, ok bool) {
	switch value {
	case 17567:
		return OpcuaNodeIdServicesVariableBase_BaseAnalogType_InstrumentRange, true
	case 17568:
		return OpcuaNodeIdServicesVariableBase_BaseAnalogType_EURange, true
	case 17569:
		return OpcuaNodeIdServicesVariableBase_BaseAnalogType_EngineeringUnits, true
	case 2042:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_EventId, true
	case 2043:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_EventType, true
	case 2044:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_SourceNode, true
	case 2045:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_SourceName, true
	case 2046:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_Time, true
	case 2047:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ReceiveTime, true
	case 2050:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_Message, true
	case 2051:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_Severity, true
	case 31771:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassId, true
	case 31772:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassName, true
	case 31773:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassId, true
	case 31774:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassName, true
	case 3190:
		return OpcuaNodeIdServicesVariableBase_BaseEventType_LocalTime, true
	}
	return 0, false
}

func OpcuaNodeIdServicesVariableBaseByName(value string) (enum OpcuaNodeIdServicesVariableBase, ok bool) {
	switch value {
	case "BaseAnalogType_InstrumentRange":
		return OpcuaNodeIdServicesVariableBase_BaseAnalogType_InstrumentRange, true
	case "BaseAnalogType_EURange":
		return OpcuaNodeIdServicesVariableBase_BaseAnalogType_EURange, true
	case "BaseAnalogType_EngineeringUnits":
		return OpcuaNodeIdServicesVariableBase_BaseAnalogType_EngineeringUnits, true
	case "BaseEventType_EventId":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_EventId, true
	case "BaseEventType_EventType":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_EventType, true
	case "BaseEventType_SourceNode":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_SourceNode, true
	case "BaseEventType_SourceName":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_SourceName, true
	case "BaseEventType_Time":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_Time, true
	case "BaseEventType_ReceiveTime":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ReceiveTime, true
	case "BaseEventType_Message":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_Message, true
	case "BaseEventType_Severity":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_Severity, true
	case "BaseEventType_ConditionClassId":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassId, true
	case "BaseEventType_ConditionClassName":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassName, true
	case "BaseEventType_ConditionSubClassId":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassId, true
	case "BaseEventType_ConditionSubClassName":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassName, true
	case "BaseEventType_LocalTime":
		return OpcuaNodeIdServicesVariableBase_BaseEventType_LocalTime, true
	}
	return 0, false
}

func OpcuaNodeIdServicesVariableBaseKnows(value int32) bool {
	for _, typeValue := range OpcuaNodeIdServicesVariableBaseValues {
		if int32(typeValue) == value {
			return true
		}
	}
	return false
}

func CastOpcuaNodeIdServicesVariableBase(structType any) OpcuaNodeIdServicesVariableBase {
	castFunc := func(typ any) OpcuaNodeIdServicesVariableBase {
		if sOpcuaNodeIdServicesVariableBase, ok := typ.(OpcuaNodeIdServicesVariableBase); ok {
			return sOpcuaNodeIdServicesVariableBase
		}
		return 0
	}
	return castFunc(structType)
}

func (m OpcuaNodeIdServicesVariableBase) GetLengthInBits(ctx context.Context) uint16 {
	return 32
}

func (m OpcuaNodeIdServicesVariableBase) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func OpcuaNodeIdServicesVariableBaseParse(ctx context.Context, theBytes []byte) (OpcuaNodeIdServicesVariableBase, error) {
	return OpcuaNodeIdServicesVariableBaseParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func OpcuaNodeIdServicesVariableBaseParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (OpcuaNodeIdServicesVariableBase, error) {
	log := zerolog.Ctx(ctx)
	_ = log
	val, err := readBuffer.ReadInt32("OpcuaNodeIdServicesVariableBase", 32)
	if err != nil {
		return 0, errors.Wrap(err, "error reading OpcuaNodeIdServicesVariableBase")
	}
	if enum, ok := OpcuaNodeIdServicesVariableBaseByValue(val); !ok {
		log.Debug().Interface("val", val).Msg("no value val found for OpcuaNodeIdServicesVariableBase")
		return OpcuaNodeIdServicesVariableBase(val), nil
	} else {
		return enum, nil
	}
}

func (e OpcuaNodeIdServicesVariableBase) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e OpcuaNodeIdServicesVariableBase) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	log := zerolog.Ctx(ctx)
	_ = log
	return writeBuffer.WriteInt32("OpcuaNodeIdServicesVariableBase", 32, int32(int32(e)), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e OpcuaNodeIdServicesVariableBase) PLC4XEnumName() string {
	switch e {
	case OpcuaNodeIdServicesVariableBase_BaseAnalogType_InstrumentRange:
		return "BaseAnalogType_InstrumentRange"
	case OpcuaNodeIdServicesVariableBase_BaseAnalogType_EURange:
		return "BaseAnalogType_EURange"
	case OpcuaNodeIdServicesVariableBase_BaseAnalogType_EngineeringUnits:
		return "BaseAnalogType_EngineeringUnits"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_EventId:
		return "BaseEventType_EventId"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_EventType:
		return "BaseEventType_EventType"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_SourceNode:
		return "BaseEventType_SourceNode"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_SourceName:
		return "BaseEventType_SourceName"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_Time:
		return "BaseEventType_Time"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_ReceiveTime:
		return "BaseEventType_ReceiveTime"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_Message:
		return "BaseEventType_Message"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_Severity:
		return "BaseEventType_Severity"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassId:
		return "BaseEventType_ConditionClassId"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionClassName:
		return "BaseEventType_ConditionClassName"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassId:
		return "BaseEventType_ConditionSubClassId"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_ConditionSubClassName:
		return "BaseEventType_ConditionSubClassName"
	case OpcuaNodeIdServicesVariableBase_BaseEventType_LocalTime:
		return "BaseEventType_LocalTime"
	}
	return fmt.Sprintf("Unknown(%v)", int32(e))
}

func (e OpcuaNodeIdServicesVariableBase) String() string {
	return e.PLC4XEnumName()
}
