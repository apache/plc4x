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

// SimulatedDataTypeSizes is an enum
type SimulatedDataTypeSizes uint8

type ISimulatedDataTypeSizes interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	DataTypeSize() uint8
}

const (
	SimulatedDataTypeSizes_BOOL           SimulatedDataTypeSizes = 1
	SimulatedDataTypeSizes_BYTE           SimulatedDataTypeSizes = 2
	SimulatedDataTypeSizes_WORD           SimulatedDataTypeSizes = 3
	SimulatedDataTypeSizes_DWORD          SimulatedDataTypeSizes = 4
	SimulatedDataTypeSizes_LWORD          SimulatedDataTypeSizes = 5
	SimulatedDataTypeSizes_SINT           SimulatedDataTypeSizes = 6
	SimulatedDataTypeSizes_INT            SimulatedDataTypeSizes = 7
	SimulatedDataTypeSizes_DINT           SimulatedDataTypeSizes = 8
	SimulatedDataTypeSizes_LINT           SimulatedDataTypeSizes = 9
	SimulatedDataTypeSizes_USINT          SimulatedDataTypeSizes = 10
	SimulatedDataTypeSizes_UINT           SimulatedDataTypeSizes = 11
	SimulatedDataTypeSizes_UDINT          SimulatedDataTypeSizes = 12
	SimulatedDataTypeSizes_ULINT          SimulatedDataTypeSizes = 13
	SimulatedDataTypeSizes_REAL           SimulatedDataTypeSizes = 14
	SimulatedDataTypeSizes_LREAL          SimulatedDataTypeSizes = 15
	SimulatedDataTypeSizes_TIME           SimulatedDataTypeSizes = 16
	SimulatedDataTypeSizes_LTIME          SimulatedDataTypeSizes = 17
	SimulatedDataTypeSizes_DATE           SimulatedDataTypeSizes = 18
	SimulatedDataTypeSizes_LDATE          SimulatedDataTypeSizes = 19
	SimulatedDataTypeSizes_TIME_OF_DAY    SimulatedDataTypeSizes = 20
	SimulatedDataTypeSizes_LTIME_OF_DAY   SimulatedDataTypeSizes = 21
	SimulatedDataTypeSizes_DATE_AND_TIME  SimulatedDataTypeSizes = 22
	SimulatedDataTypeSizes_LDATE_AND_TIME SimulatedDataTypeSizes = 23
	SimulatedDataTypeSizes_CHAR           SimulatedDataTypeSizes = 24
	SimulatedDataTypeSizes_WCHAR          SimulatedDataTypeSizes = 25
	SimulatedDataTypeSizes_STRING         SimulatedDataTypeSizes = 26
	SimulatedDataTypeSizes_WSTRING        SimulatedDataTypeSizes = 27
)

var SimulatedDataTypeSizesValues []SimulatedDataTypeSizes

func init() {
	_ = errors.New
	SimulatedDataTypeSizesValues = []SimulatedDataTypeSizes{
		SimulatedDataTypeSizes_BOOL,
		SimulatedDataTypeSizes_BYTE,
		SimulatedDataTypeSizes_WORD,
		SimulatedDataTypeSizes_DWORD,
		SimulatedDataTypeSizes_LWORD,
		SimulatedDataTypeSizes_SINT,
		SimulatedDataTypeSizes_INT,
		SimulatedDataTypeSizes_DINT,
		SimulatedDataTypeSizes_LINT,
		SimulatedDataTypeSizes_USINT,
		SimulatedDataTypeSizes_UINT,
		SimulatedDataTypeSizes_UDINT,
		SimulatedDataTypeSizes_ULINT,
		SimulatedDataTypeSizes_REAL,
		SimulatedDataTypeSizes_LREAL,
		SimulatedDataTypeSizes_TIME,
		SimulatedDataTypeSizes_LTIME,
		SimulatedDataTypeSizes_DATE,
		SimulatedDataTypeSizes_LDATE,
		SimulatedDataTypeSizes_TIME_OF_DAY,
		SimulatedDataTypeSizes_LTIME_OF_DAY,
		SimulatedDataTypeSizes_DATE_AND_TIME,
		SimulatedDataTypeSizes_LDATE_AND_TIME,
		SimulatedDataTypeSizes_CHAR,
		SimulatedDataTypeSizes_WCHAR,
		SimulatedDataTypeSizes_STRING,
		SimulatedDataTypeSizes_WSTRING,
	}
}

func (e SimulatedDataTypeSizes) DataTypeSize() uint8 {
	switch e {
	case 1:
		{ /* '1' */
			return 1
		}
	case 10:
		{ /* '10' */
			return 1
		}
	case 11:
		{ /* '11' */
			return 2
		}
	case 12:
		{ /* '12' */
			return 4
		}
	case 13:
		{ /* '13' */
			return 8
		}
	case 14:
		{ /* '14' */
			return 4
		}
	case 15:
		{ /* '15' */
			return 8
		}
	case 16:
		{ /* '16' */
			return 8
		}
	case 17:
		{ /* '17' */
			return 8
		}
	case 18:
		{ /* '18' */
			return 8
		}
	case 19:
		{ /* '19' */
			return 8
		}
	case 2:
		{ /* '2' */
			return 1
		}
	case 20:
		{ /* '20' */
			return 8
		}
	case 21:
		{ /* '21' */
			return 8
		}
	case 22:
		{ /* '22' */
			return 8
		}
	case 23:
		{ /* '23' */
			return 8
		}
	case 24:
		{ /* '24' */
			return 1
		}
	case 25:
		{ /* '25' */
			return 2
		}
	case 26:
		{ /* '26' */
			return 255
		}
	case 27:
		{ /* '27' */
			return 127
		}
	case 3:
		{ /* '3' */
			return 2
		}
	case 4:
		{ /* '4' */
			return 4
		}
	case 5:
		{ /* '5' */
			return 8
		}
	case 6:
		{ /* '6' */
			return 1
		}
	case 7:
		{ /* '7' */
			return 2
		}
	case 8:
		{ /* '8' */
			return 4
		}
	case 9:
		{ /* '9' */
			return 8
		}
	default:
		{
			return 0
		}
	}
}

func SimulatedDataTypeSizesFirstEnumForFieldDataTypeSize(value uint8) (enum SimulatedDataTypeSizes, ok bool) {
	for _, sizeValue := range SimulatedDataTypeSizesValues {
		if sizeValue.DataTypeSize() == value {
			return sizeValue, true
		}
	}
	return 0, false
}
func SimulatedDataTypeSizesByValue(value uint8) (enum SimulatedDataTypeSizes, ok bool) {
	switch value {
	case 1:
		return SimulatedDataTypeSizes_BOOL, true
	case 10:
		return SimulatedDataTypeSizes_USINT, true
	case 11:
		return SimulatedDataTypeSizes_UINT, true
	case 12:
		return SimulatedDataTypeSizes_UDINT, true
	case 13:
		return SimulatedDataTypeSizes_ULINT, true
	case 14:
		return SimulatedDataTypeSizes_REAL, true
	case 15:
		return SimulatedDataTypeSizes_LREAL, true
	case 16:
		return SimulatedDataTypeSizes_TIME, true
	case 17:
		return SimulatedDataTypeSizes_LTIME, true
	case 18:
		return SimulatedDataTypeSizes_DATE, true
	case 19:
		return SimulatedDataTypeSizes_LDATE, true
	case 2:
		return SimulatedDataTypeSizes_BYTE, true
	case 20:
		return SimulatedDataTypeSizes_TIME_OF_DAY, true
	case 21:
		return SimulatedDataTypeSizes_LTIME_OF_DAY, true
	case 22:
		return SimulatedDataTypeSizes_DATE_AND_TIME, true
	case 23:
		return SimulatedDataTypeSizes_LDATE_AND_TIME, true
	case 24:
		return SimulatedDataTypeSizes_CHAR, true
	case 25:
		return SimulatedDataTypeSizes_WCHAR, true
	case 26:
		return SimulatedDataTypeSizes_STRING, true
	case 27:
		return SimulatedDataTypeSizes_WSTRING, true
	case 3:
		return SimulatedDataTypeSizes_WORD, true
	case 4:
		return SimulatedDataTypeSizes_DWORD, true
	case 5:
		return SimulatedDataTypeSizes_LWORD, true
	case 6:
		return SimulatedDataTypeSizes_SINT, true
	case 7:
		return SimulatedDataTypeSizes_INT, true
	case 8:
		return SimulatedDataTypeSizes_DINT, true
	case 9:
		return SimulatedDataTypeSizes_LINT, true
	}
	return 0, false
}

func SimulatedDataTypeSizesByName(value string) (enum SimulatedDataTypeSizes, ok bool) {
	switch value {
	case "BOOL":
		return SimulatedDataTypeSizes_BOOL, true
	case "USINT":
		return SimulatedDataTypeSizes_USINT, true
	case "UINT":
		return SimulatedDataTypeSizes_UINT, true
	case "UDINT":
		return SimulatedDataTypeSizes_UDINT, true
	case "ULINT":
		return SimulatedDataTypeSizes_ULINT, true
	case "REAL":
		return SimulatedDataTypeSizes_REAL, true
	case "LREAL":
		return SimulatedDataTypeSizes_LREAL, true
	case "TIME":
		return SimulatedDataTypeSizes_TIME, true
	case "LTIME":
		return SimulatedDataTypeSizes_LTIME, true
	case "DATE":
		return SimulatedDataTypeSizes_DATE, true
	case "LDATE":
		return SimulatedDataTypeSizes_LDATE, true
	case "BYTE":
		return SimulatedDataTypeSizes_BYTE, true
	case "TIME_OF_DAY":
		return SimulatedDataTypeSizes_TIME_OF_DAY, true
	case "LTIME_OF_DAY":
		return SimulatedDataTypeSizes_LTIME_OF_DAY, true
	case "DATE_AND_TIME":
		return SimulatedDataTypeSizes_DATE_AND_TIME, true
	case "LDATE_AND_TIME":
		return SimulatedDataTypeSizes_LDATE_AND_TIME, true
	case "CHAR":
		return SimulatedDataTypeSizes_CHAR, true
	case "WCHAR":
		return SimulatedDataTypeSizes_WCHAR, true
	case "STRING":
		return SimulatedDataTypeSizes_STRING, true
	case "WSTRING":
		return SimulatedDataTypeSizes_WSTRING, true
	case "WORD":
		return SimulatedDataTypeSizes_WORD, true
	case "DWORD":
		return SimulatedDataTypeSizes_DWORD, true
	case "LWORD":
		return SimulatedDataTypeSizes_LWORD, true
	case "SINT":
		return SimulatedDataTypeSizes_SINT, true
	case "INT":
		return SimulatedDataTypeSizes_INT, true
	case "DINT":
		return SimulatedDataTypeSizes_DINT, true
	case "LINT":
		return SimulatedDataTypeSizes_LINT, true
	}
	return 0, false
}

func SimulatedDataTypeSizesKnows(value uint8) bool {
	for _, typeValue := range SimulatedDataTypeSizesValues {
		if uint8(typeValue) == value {
			return true
		}
	}
	return false
}

func CastSimulatedDataTypeSizes(structType any) SimulatedDataTypeSizes {
	castFunc := func(typ any) SimulatedDataTypeSizes {
		if sSimulatedDataTypeSizes, ok := typ.(SimulatedDataTypeSizes); ok {
			return sSimulatedDataTypeSizes
		}
		return 0
	}
	return castFunc(structType)
}

func (m SimulatedDataTypeSizes) GetLengthInBits(ctx context.Context) uint16 {
	return 8
}

func (m SimulatedDataTypeSizes) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func SimulatedDataTypeSizesParse(ctx context.Context, theBytes []byte) (SimulatedDataTypeSizes, error) {
	return SimulatedDataTypeSizesParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func SimulatedDataTypeSizesParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (SimulatedDataTypeSizes, error) {
	log := zerolog.Ctx(ctx)
	_ = log
	val, err := /*TODO: migrate me*/ /*TODO: migrate me*/ readBuffer.ReadUint8("SimulatedDataTypeSizes", 8)
	if err != nil {
		return 0, errors.Wrap(err, "error reading SimulatedDataTypeSizes")
	}
	if enum, ok := SimulatedDataTypeSizesByValue(val); !ok {
		log.Debug().Interface("val", val).Msg("no value val found for SimulatedDataTypeSizes")
		return SimulatedDataTypeSizes(val), nil
	} else {
		return enum, nil
	}
}

func (e SimulatedDataTypeSizes) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e SimulatedDataTypeSizes) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	log := zerolog.Ctx(ctx)
	_ = log
	return /*TODO: migrate me*/ writeBuffer.WriteUint8("SimulatedDataTypeSizes", 8, uint8(uint8(e)), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e SimulatedDataTypeSizes) PLC4XEnumName() string {
	switch e {
	case SimulatedDataTypeSizes_BOOL:
		return "BOOL"
	case SimulatedDataTypeSizes_USINT:
		return "USINT"
	case SimulatedDataTypeSizes_UINT:
		return "UINT"
	case SimulatedDataTypeSizes_UDINT:
		return "UDINT"
	case SimulatedDataTypeSizes_ULINT:
		return "ULINT"
	case SimulatedDataTypeSizes_REAL:
		return "REAL"
	case SimulatedDataTypeSizes_LREAL:
		return "LREAL"
	case SimulatedDataTypeSizes_TIME:
		return "TIME"
	case SimulatedDataTypeSizes_LTIME:
		return "LTIME"
	case SimulatedDataTypeSizes_DATE:
		return "DATE"
	case SimulatedDataTypeSizes_LDATE:
		return "LDATE"
	case SimulatedDataTypeSizes_BYTE:
		return "BYTE"
	case SimulatedDataTypeSizes_TIME_OF_DAY:
		return "TIME_OF_DAY"
	case SimulatedDataTypeSizes_LTIME_OF_DAY:
		return "LTIME_OF_DAY"
	case SimulatedDataTypeSizes_DATE_AND_TIME:
		return "DATE_AND_TIME"
	case SimulatedDataTypeSizes_LDATE_AND_TIME:
		return "LDATE_AND_TIME"
	case SimulatedDataTypeSizes_CHAR:
		return "CHAR"
	case SimulatedDataTypeSizes_WCHAR:
		return "WCHAR"
	case SimulatedDataTypeSizes_STRING:
		return "STRING"
	case SimulatedDataTypeSizes_WSTRING:
		return "WSTRING"
	case SimulatedDataTypeSizes_WORD:
		return "WORD"
	case SimulatedDataTypeSizes_DWORD:
		return "DWORD"
	case SimulatedDataTypeSizes_LWORD:
		return "LWORD"
	case SimulatedDataTypeSizes_SINT:
		return "SINT"
	case SimulatedDataTypeSizes_INT:
		return "INT"
	case SimulatedDataTypeSizes_DINT:
		return "DINT"
	case SimulatedDataTypeSizes_LINT:
		return "LINT"
	}
	return fmt.Sprintf("Unknown(%v)", uint8(e))
}

func (e SimulatedDataTypeSizes) String() string {
	return e.PLC4XEnumName()
}
