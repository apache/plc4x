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

// BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter is an enum
type BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter uint8

type IBACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter interface {
	utils.Serializable
}

const (
	BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ALL       BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter = 0
	BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ACKED     BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter = 1
	BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_NOT_ACKED BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter = 2
)

var BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterValues []BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter

func init() {
	_ = errors.New
	BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterValues = []BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter{
		BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ALL,
		BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ACKED,
		BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_NOT_ACKED,
	}
}

func BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterByValue(value uint8) (enum BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter, ok bool) {
	switch value {
	case 0:
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ALL, true
	case 1:
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ACKED, true
	case 2:
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_NOT_ACKED, true
	}
	return 0, false
}

func BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterByName(value string) (enum BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter, ok bool) {
	switch value {
	case "ALL":
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ALL, true
	case "ACKED":
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ACKED, true
	case "NOT_ACKED":
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_NOT_ACKED, true
	}
	return 0, false
}

func BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterKnows(value uint8) bool {
	for _, typeValue := range BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterValues {
		if uint8(typeValue) == value {
			return true
		}
	}
	return false
}

func CastBACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter(structType interface{}) BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter {
		if sBACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter, ok := typ.(BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter); ok {
			return sBACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter
		}
		return 0
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) GetLengthInBits() uint16 {
	return 8
}

func (m BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterParse(theBytes []byte) (BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter, error) {
	return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterParseWithBuffer(utils.NewReadBufferByteBased(theBytes))
}

func BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterParseWithBuffer(readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter, error) {
	val, err := readBuffer.ReadUint8("BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter", 8)
	if err != nil {
		return 0, errors.Wrap(err, "error reading BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter")
	}
	if enum, ok := BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterByValue(val); !ok {
		Plc4xModelLog.Debug().Msgf("no value %x found for RequestType", val)
		return BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter(val), nil
	} else {
		return enum, nil
	}
}

func (e BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := e.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (e BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint8("BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter", 8, uint8(e), utils.WithAdditionalStringRepresentation(e.PLC4XEnumName()))
}

// PLC4XEnumName returns the name that is used in code to identify this enum
func (e BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) PLC4XEnumName() string {
	switch e {
	case BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ALL:
		return "ALL"
	case BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_ACKED:
		return "ACKED"
	case BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter_NOT_ACKED:
		return "NOT_ACKED"
	}
	return ""
}

func (e BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter) String() string {
	return e.PLC4XEnumName()
}
