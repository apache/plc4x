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
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetAccessPassbackMode is an enum
type BACnetAccessPassbackMode uint8

type IBACnetAccessPassbackMode interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
}

const (
	BACnetAccessPassbackMode_PASSBACK_OFF  BACnetAccessPassbackMode = 0
	BACnetAccessPassbackMode_HARD_PASSBACK BACnetAccessPassbackMode = 1
	BACnetAccessPassbackMode_SOFT_PASSBACK BACnetAccessPassbackMode = 2
)

var BACnetAccessPassbackModeValues []BACnetAccessPassbackMode

func init() {
	_ = errors.New
	BACnetAccessPassbackModeValues = []BACnetAccessPassbackMode{
		BACnetAccessPassbackMode_PASSBACK_OFF,
		BACnetAccessPassbackMode_HARD_PASSBACK,
		BACnetAccessPassbackMode_SOFT_PASSBACK,
	}
}

func BACnetAccessPassbackModeByValue(value uint8) (enum BACnetAccessPassbackMode, ok bool) {
	switch value {
	case 0:
		return BACnetAccessPassbackMode_PASSBACK_OFF, true
	case 1:
		return BACnetAccessPassbackMode_HARD_PASSBACK, true
	case 2:
		return BACnetAccessPassbackMode_SOFT_PASSBACK, true
	}
	return 0, false
}

func BACnetAccessPassbackModeByName(value string) (enum BACnetAccessPassbackMode, ok bool) {
	switch value {
	case "PASSBACK_OFF":
		return BACnetAccessPassbackMode_PASSBACK_OFF, true
	case "HARD_PASSBACK":
		return BACnetAccessPassbackMode_HARD_PASSBACK, true
	case "SOFT_PASSBACK":
		return BACnetAccessPassbackMode_SOFT_PASSBACK, true
	}
	return 0, false
}

func BACnetAccessPassbackModeKnows(value uint8) bool {
	for _, typeValue := range BACnetAccessPassbackModeValues {
		if uint8(typeValue) == value {
			return true
		}
	}
	return false
}

func CastBACnetAccessPassbackMode(structType any) BACnetAccessPassbackMode {
	castFunc := func(typ any) BACnetAccessPassbackMode {
		if sBACnetAccessPassbackMode, ok := typ.(BACnetAccessPassbackMode); ok {
			return sBACnetAccessPassbackMode
		}
		return 0
	}
	return castFunc(structType)
}

func (m BACnetAccessPassbackMode) GetLengthInBits(ctx context.Context) uint16 {
	return 8
}

func (m BACnetAccessPassbackMode) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetAccessPassbackModeParse(ctx context.Context, theBytes []byte) (BACnetAccessPassbackMode, error) {
	return BACnetAccessPassbackModeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetAccessPassbackModeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetAccessPassbackMode, error) {
	val, err := readBuffer.ReadUint8("BACnetAccessPassbackMode", 8)
	if err != nil {
		return 0, errors.Wrap(err, "error reading BACnetAccessPassbackMode")
	}
	if enum, ok := BACnetAccessPassbackModeByValue(val); !ok {
		Plc4xModelLog.Debug().Msgf("no value %x found for RequestType", val)
		return BACnetAccessPassbackMode(val), nil
	} else {
		return enum, nil
	}
}

func (e BACnetAccessPassbackMode) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e BACnetAccessPassbackMode) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint8("BACnetAccessPassbackMode", 8, uint8(e), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e BACnetAccessPassbackMode) PLC4XEnumName() string {
	switch e {
	case BACnetAccessPassbackMode_PASSBACK_OFF:
		return "PASSBACK_OFF"
	case BACnetAccessPassbackMode_HARD_PASSBACK:
		return "HARD_PASSBACK"
	case BACnetAccessPassbackMode_SOFT_PASSBACK:
		return "SOFT_PASSBACK"
	}
	return ""
}

func (e BACnetAccessPassbackMode) String() string {
	return e.PLC4XEnumName()
}
