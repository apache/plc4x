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

// OpcuaNodeIdServicesVariableSystem is an enum
type OpcuaNodeIdServicesVariableSystem int32

type IOpcuaNodeIdServicesVariableSystem interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
}

const (
	OpcuaNodeIdServicesVariableSystem_SystemStatusChangeEventType_SystemState OpcuaNodeIdServicesVariableSystem = 11696
)

var OpcuaNodeIdServicesVariableSystemValues []OpcuaNodeIdServicesVariableSystem

func init() {
	_ = errors.New
	OpcuaNodeIdServicesVariableSystemValues = []OpcuaNodeIdServicesVariableSystem{
		OpcuaNodeIdServicesVariableSystem_SystemStatusChangeEventType_SystemState,
	}
}

func OpcuaNodeIdServicesVariableSystemByValue(value int32) (enum OpcuaNodeIdServicesVariableSystem, ok bool) {
	switch value {
	case 11696:
		return OpcuaNodeIdServicesVariableSystem_SystemStatusChangeEventType_SystemState, true
	}
	return 0, false
}

func OpcuaNodeIdServicesVariableSystemByName(value string) (enum OpcuaNodeIdServicesVariableSystem, ok bool) {
	switch value {
	case "SystemStatusChangeEventType_SystemState":
		return OpcuaNodeIdServicesVariableSystem_SystemStatusChangeEventType_SystemState, true
	}
	return 0, false
}

func OpcuaNodeIdServicesVariableSystemKnows(value int32) bool {
	for _, typeValue := range OpcuaNodeIdServicesVariableSystemValues {
		if int32(typeValue) == value {
			return true
		}
	}
	return false
}

func CastOpcuaNodeIdServicesVariableSystem(structType any) OpcuaNodeIdServicesVariableSystem {
	castFunc := func(typ any) OpcuaNodeIdServicesVariableSystem {
		if sOpcuaNodeIdServicesVariableSystem, ok := typ.(OpcuaNodeIdServicesVariableSystem); ok {
			return sOpcuaNodeIdServicesVariableSystem
		}
		return 0
	}
	return castFunc(structType)
}

func (m OpcuaNodeIdServicesVariableSystem) GetLengthInBits(ctx context.Context) uint16 {
	return 32
}

func (m OpcuaNodeIdServicesVariableSystem) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func OpcuaNodeIdServicesVariableSystemParse(ctx context.Context, theBytes []byte) (OpcuaNodeIdServicesVariableSystem, error) {
	return OpcuaNodeIdServicesVariableSystemParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func OpcuaNodeIdServicesVariableSystemParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (OpcuaNodeIdServicesVariableSystem, error) {
	log := zerolog.Ctx(ctx)
	_ = log
	val, err := readBuffer.ReadInt32("OpcuaNodeIdServicesVariableSystem", 32)
	if err != nil {
		return 0, errors.Wrap(err, "error reading OpcuaNodeIdServicesVariableSystem")
	}
	if enum, ok := OpcuaNodeIdServicesVariableSystemByValue(val); !ok {
		log.Debug().Interface("val", val).Msg("no value val found for OpcuaNodeIdServicesVariableSystem")
		return OpcuaNodeIdServicesVariableSystem(val), nil
	} else {
		return enum, nil
	}
}

func (e OpcuaNodeIdServicesVariableSystem) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e OpcuaNodeIdServicesVariableSystem) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	log := zerolog.Ctx(ctx)
	_ = log
	return writeBuffer.WriteInt32("OpcuaNodeIdServicesVariableSystem", 32, int32(int32(e)), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e OpcuaNodeIdServicesVariableSystem) PLC4XEnumName() string {
	switch e {
	case OpcuaNodeIdServicesVariableSystem_SystemStatusChangeEventType_SystemState:
		return "SystemStatusChangeEventType_SystemState"
	}
	return fmt.Sprintf("Unknown(%v)", int32(e))
}

func (e OpcuaNodeIdServicesVariableSystem) String() string {
	return e.PLC4XEnumName()
}
