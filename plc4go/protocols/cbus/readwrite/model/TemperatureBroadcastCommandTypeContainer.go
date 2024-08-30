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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// TemperatureBroadcastCommandTypeContainer is an enum
type TemperatureBroadcastCommandTypeContainer uint8

type ITemperatureBroadcastCommandTypeContainer interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	NumBytes() uint8
	CommandType() TemperatureBroadcastCommandType
}

const (
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x02
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x0A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x12
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x1A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x22
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x2A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x32
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x3A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x42
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes  TemperatureBroadcastCommandTypeContainer = 0x4A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes TemperatureBroadcastCommandTypeContainer = 0x52
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes TemperatureBroadcastCommandTypeContainer = 0x5A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes TemperatureBroadcastCommandTypeContainer = 0x62
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes TemperatureBroadcastCommandTypeContainer = 0x6A
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes TemperatureBroadcastCommandTypeContainer = 0x72
	TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes TemperatureBroadcastCommandTypeContainer = 0x7A
)

var TemperatureBroadcastCommandTypeContainerValues []TemperatureBroadcastCommandTypeContainer

func init() {
	_ = errors.New
	TemperatureBroadcastCommandTypeContainerValues = []TemperatureBroadcastCommandTypeContainer{
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes,
		TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes,
	}
}

func (e TemperatureBroadcastCommandTypeContainer) NumBytes() uint8 {
	switch e {
	case 0x02:
		{ /* '0x02' */
			return 2
		}
	case 0x0A:
		{ /* '0x0A' */
			return 2
		}
	case 0x12:
		{ /* '0x12' */
			return 2
		}
	case 0x1A:
		{ /* '0x1A' */
			return 2
		}
	case 0x22:
		{ /* '0x22' */
			return 2
		}
	case 0x2A:
		{ /* '0x2A' */
			return 2
		}
	case 0x32:
		{ /* '0x32' */
			return 2
		}
	case 0x3A:
		{ /* '0x3A' */
			return 2
		}
	case 0x42:
		{ /* '0x42' */
			return 2
		}
	case 0x4A:
		{ /* '0x4A' */
			return 2
		}
	case 0x52:
		{ /* '0x52' */
			return 2
		}
	case 0x5A:
		{ /* '0x5A' */
			return 2
		}
	case 0x62:
		{ /* '0x62' */
			return 2
		}
	case 0x6A:
		{ /* '0x6A' */
			return 2
		}
	case 0x72:
		{ /* '0x72' */
			return 2
		}
	case 0x7A:
		{ /* '0x7A' */
			return 2
		}
	default:
		{
			return 0
		}
	}
}

func TemperatureBroadcastCommandTypeContainerFirstEnumForFieldNumBytes(value uint8) (enum TemperatureBroadcastCommandTypeContainer, ok bool) {
	for _, sizeValue := range TemperatureBroadcastCommandTypeContainerValues {
		if sizeValue.NumBytes() == value {
			return sizeValue, true
		}
	}
	return 0, false
}

func (e TemperatureBroadcastCommandTypeContainer) CommandType() TemperatureBroadcastCommandType {
	switch e {
	case 0x02:
		{ /* '0x02' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x0A:
		{ /* '0x0A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x12:
		{ /* '0x12' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x1A:
		{ /* '0x1A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x22:
		{ /* '0x22' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x2A:
		{ /* '0x2A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x32:
		{ /* '0x32' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x3A:
		{ /* '0x3A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x42:
		{ /* '0x42' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x4A:
		{ /* '0x4A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x52:
		{ /* '0x52' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x5A:
		{ /* '0x5A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x62:
		{ /* '0x62' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x6A:
		{ /* '0x6A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x72:
		{ /* '0x72' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	case 0x7A:
		{ /* '0x7A' */
			return TemperatureBroadcastCommandType_BROADCAST_EVENT
		}
	default:
		{
			return 0
		}
	}
}

func TemperatureBroadcastCommandTypeContainerFirstEnumForFieldCommandType(value TemperatureBroadcastCommandType) (enum TemperatureBroadcastCommandTypeContainer, ok bool) {
	for _, sizeValue := range TemperatureBroadcastCommandTypeContainerValues {
		if sizeValue.CommandType() == value {
			return sizeValue, true
		}
	}
	return 0, false
}
func TemperatureBroadcastCommandTypeContainerByValue(value uint8) (enum TemperatureBroadcastCommandTypeContainer, ok bool) {
	switch value {
	case 0x02:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes, true
	case 0x0A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes, true
	case 0x12:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes, true
	case 0x1A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes, true
	case 0x22:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes, true
	case 0x2A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes, true
	case 0x32:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes, true
	case 0x3A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes, true
	case 0x42:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes, true
	case 0x4A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes, true
	case 0x52:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes, true
	case 0x5A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes, true
	case 0x62:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes, true
	case 0x6A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes, true
	case 0x72:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes, true
	case 0x7A:
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes, true
	}
	return 0, false
}

func TemperatureBroadcastCommandTypeContainerByName(value string) (enum TemperatureBroadcastCommandTypeContainer, ok bool) {
	switch value {
	case "TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes, true
	case "TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes":
		return TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes, true
	}
	return 0, false
}

func TemperatureBroadcastCommandTypeContainerKnows(value uint8) bool {
	for _, typeValue := range TemperatureBroadcastCommandTypeContainerValues {
		if uint8(typeValue) == value {
			return true
		}
	}
	return false
}

func CastTemperatureBroadcastCommandTypeContainer(structType any) TemperatureBroadcastCommandTypeContainer {
	castFunc := func(typ any) TemperatureBroadcastCommandTypeContainer {
		if sTemperatureBroadcastCommandTypeContainer, ok := typ.(TemperatureBroadcastCommandTypeContainer); ok {
			return sTemperatureBroadcastCommandTypeContainer
		}
		return 0
	}
	return castFunc(structType)
}

func (m TemperatureBroadcastCommandTypeContainer) GetLengthInBits(ctx context.Context) uint16 {
	return 8
}

func (m TemperatureBroadcastCommandTypeContainer) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func TemperatureBroadcastCommandTypeContainerParse(ctx context.Context, theBytes []byte) (TemperatureBroadcastCommandTypeContainer, error) {
	return TemperatureBroadcastCommandTypeContainerParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func TemperatureBroadcastCommandTypeContainerParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (TemperatureBroadcastCommandTypeContainer, error) {
	log := zerolog.Ctx(ctx)
	_ = log
	val, err := /*TODO: migrate me*/ /*TODO: migrate me*/ readBuffer.ReadUint8("TemperatureBroadcastCommandTypeContainer", 8)
	if err != nil {
		return 0, errors.Wrap(err, "error reading TemperatureBroadcastCommandTypeContainer")
	}
	if enum, ok := TemperatureBroadcastCommandTypeContainerByValue(val); !ok {
		log.Debug().Interface("val", val).Msg("no value val found for TemperatureBroadcastCommandTypeContainer")
		return TemperatureBroadcastCommandTypeContainer(val), nil
	} else {
		return enum, nil
	}
}

func (e TemperatureBroadcastCommandTypeContainer) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e TemperatureBroadcastCommandTypeContainer) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	log := zerolog.Ctx(ctx)
	_ = log
	return /*TODO: migrate me*/ writeBuffer.WriteUint8("TemperatureBroadcastCommandTypeContainer", 8, uint8(uint8(e)), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e TemperatureBroadcastCommandTypeContainer) PLC4XEnumName() string {
	switch e {
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes"
	case TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes:
		return "TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes"
	}
	return fmt.Sprintf("Unknown(%v)", uint8(e))
}

func (e TemperatureBroadcastCommandTypeContainer) String() string {
	return e.PLC4XEnumName()
}
