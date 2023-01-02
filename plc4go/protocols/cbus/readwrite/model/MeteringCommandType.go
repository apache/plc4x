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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// MeteringCommandType is an enum
type MeteringCommandType uint8

type IMeteringCommandType interface {
	utils.Serializable
	NumberOfArguments() uint8
}

const (
	MeteringCommandType_EVENT MeteringCommandType = 0x00
)

var MeteringCommandTypeValues []MeteringCommandType

func init() {
	_ = errors.New
	MeteringCommandTypeValues = []MeteringCommandType{
		MeteringCommandType_EVENT,
	}
}

func (e MeteringCommandType) NumberOfArguments() uint8 {
	switch e {
	case 0x00:
		{ /* '0x00' */
			return 0xFF
		}
	default:
		{
			return 0
		}
	}
}

func MeteringCommandTypeFirstEnumForFieldNumberOfArguments(value uint8) (MeteringCommandType, error) {
	for _, sizeValue := range MeteringCommandTypeValues {
		if sizeValue.NumberOfArguments() == value {
			return sizeValue, nil
		}
	}
	return 0, errors.Errorf("enum for %v describing NumberOfArguments not found", value)
}
func MeteringCommandTypeByValue(value uint8) (enum MeteringCommandType, ok bool) {
	switch value {
	case 0x00:
		return MeteringCommandType_EVENT, true
	}
	return 0, false
}

func MeteringCommandTypeByName(value string) (enum MeteringCommandType, ok bool) {
	switch value {
	case "EVENT":
		return MeteringCommandType_EVENT, true
	}
	return 0, false
}

func MeteringCommandTypeKnows(value uint8) bool {
	for _, typeValue := range MeteringCommandTypeValues {
		if uint8(typeValue) == value {
			return true
		}
	}
	return false
}

func CastMeteringCommandType(structType interface{}) MeteringCommandType {
	castFunc := func(typ interface{}) MeteringCommandType {
		if sMeteringCommandType, ok := typ.(MeteringCommandType); ok {
			return sMeteringCommandType
		}
		return 0
	}
	return castFunc(structType)
}

func (m MeteringCommandType) GetLengthInBits() uint16 {
	return 4
}

func (m MeteringCommandType) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func MeteringCommandTypeParse(theBytes []byte) (MeteringCommandType, error) {
	return MeteringCommandTypeParseWithBuffer(utils.NewReadBufferByteBased(theBytes))
}

func MeteringCommandTypeParseWithBuffer(readBuffer utils.ReadBuffer) (MeteringCommandType, error) {
	val, err := readBuffer.ReadUint8("MeteringCommandType", 4)
	if err != nil {
		return 0, errors.Wrap(err, "error reading MeteringCommandType")
	}
	if enum, ok := MeteringCommandTypeByValue(val); !ok {
		Plc4xModelLog.Debug().Msgf("no value %x found for RequestType", val)
		return MeteringCommandType(val), nil
	} else {
		return enum, nil
	}
}

func (e MeteringCommandType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e MeteringCommandType) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint8("MeteringCommandType", 4, uint8(e), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e MeteringCommandType) PLC4XEnumName() string {
	switch e {
	case MeteringCommandType_EVENT:
		return "EVENT"
	}
	return ""
}

func (e MeteringCommandType) String() string {
	return e.PLC4XEnumName()
}
