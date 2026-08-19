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

package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	AddressOffset = 1

	// maxWireAddress is the highest address that fits into the 16 bit address field of a modbus
	// PDU. Ported from REGISTER_MAX_ADDRESS in plc4j ModbusTagCoil/ModbusTagHoldingRegister.
	maxWireAddress = 65535
	// maxCoilQuantity is the largest number of coils/discrete inputs a single request may cover
	// (plc4j ModbusTagCoil.of).
	maxCoilQuantity = 2000
	// maxRegisterQuantity is the largest number of values a single request may cover in one of the
	// register areas (plc4j ModbusTagHoldingRegister.of, ModbusTagInputRegister.of and
	// ModbusTagExtendedRegister.of). Anything beyond that no longer fits into the 253 byte modbus
	// PDU.
	maxRegisterQuantity = 125
	// maxStringLength bounds the declared length of a single string so that the total payload size
	// still fits into the 16 bit arithmetic the codec does.
	maxStringLength = 65535
)

type modbusTag struct {
	apiModel.PlcTag

	TagType  TagType
	Address  uint16
	Quantity uint16
	Datatype readWriteModel.ModbusDataType
	// StringLength is the declared length of a single string. Nothing on the wire announces it, so
	// it is part of the address; for every data type that is not a string it is 1, which leaves the
	// size arithmetic unchanged (plc4j ModbusTag).
	StringLength uint16
	// UnitId, when set, is the unit identifier (slave id) this one tag is addressed at, overriding
	// the connection's default (plc4j ModbusTag.getUnitId).
	UnitId *uint8
	// ByteOrder, when set, is the payload byte order of this one tag, overriding the connection's
	// default-payload-byte-order (plc4j ModbusTag.getByteOrder).
	ByteOrder *ByteOrder
}

// logicalAddressOffset is how much higher the address the user writes is than the address that
// goes onto the wire. Every area but the extended registers is addressed one-based (plc4j
// ModbusTagCoil.getLogicalAddress and its siblings); the extended register area is addressed
// starting at zero, so its logical address is the wire address itself (plc4j
// ModbusTagExtendedRegister.getLogicalAddress).
func logicalAddressOffset(tagType TagType) uint16 {
	if tagType == ExtendedRegister {
		return 0
	}
	return AddressOffset
}

// NewTag builds a tag from the logical (user facing) address, which for every area but the
// extended registers is one higher than the address that goes onto the wire.
func NewTag(tagType TagType, address uint16, quantity uint16, datatype readWriteModel.ModbusDataType) apiModel.PlcTag {
	return newTagFromWireAddress(tagType, address-logicalAddressOffset(tagType), quantity, datatype, 1, tagConfig{})
}

func newTagFromWireAddress(tagType TagType, wireAddress uint16, quantity uint16, datatype readWriteModel.ModbusDataType, stringLength uint16, config tagConfig) modbusTag {
	return modbusTag{
		TagType:      tagType,
		Address:      wireAddress,
		Quantity:     quantity,
		Datatype:     datatype,
		StringLength: stringLength,
		UnitId:       config.unitId,
		ByteOrder:    config.byteOrder,
	}
}

// isStringType says whether the data type carries a length of its own.
func isStringType(datatype readWriteModel.ModbusDataType) bool {
	return datatype == readWriteModel.ModbusDataType_STRING || datatype == readWriteModel.ModbusDataType_WSTRING
}

// validateStringLength ports the rules plc4j's ModbusTag applies to the length in parentheses: a
// string has to have one, and nothing else may.
// The returned length is what the codec has to work with - 1 for every non-string type.
func validateStringLength(datatype readWriteModel.ModbusDataType, stringLengthString string) (uint16, error) {
	if !isStringType(datatype) {
		if stringLengthString != "" {
			return 0, errors.Errorf("a length in parentheses is only supported for STRING and WSTRING, not for %s", datatype.String())
		}
		return 1, nil
	}
	if stringLengthString == "" {
		return 0, errors.Errorf("%s requires the length of one string, for example '%s(20)'", datatype.String(), datatype.String())
	}
	stringLength, err := strconv.ParseUint(stringLengthString, 10, 32)
	if err != nil {
		return 0, errors.Errorf("Couldn't parse string length '%s' into an int", stringLengthString)
	}
	if stringLength == 0 {
		return 0, errors.Errorf("string length must be greater than zero. Was %d", stringLength)
	}
	if stringLength > maxStringLength {
		return 0, errors.Errorf("string length may not be larger than %d. Was %d", maxStringLength, stringLength)
	}
	return uint16(stringLength), nil
}

// lengthWords is how many registers the payload of this tag occupies (plc4j
// ModbusTag.getLengthWords). Rounded up, because a value narrower than a register still occupies a
// whole one, and never zero, because every request has to ask for something.
// The length is the one the register codec actually reads and writes, not the natural width of the
// data type times the number of values: several values narrower than a byte are packed, so a
// BOOL[3] is one register rather than the two that counting bytes would suggest.
// A payload that doesn't fit into the 16 bit quantity field of a request is an error rather than a
// silently truncated request.
func (m modbusTag) lengthWords() (uint16, error) {
	lengthBytes := lengthInBytes(m.Datatype, m.Quantity, m.StringLength)
	words := (lengthBytes + 1) / 2
	if words == 0 {
		return 1, nil
	}
	if words > maxWireAddress {
		return 0, errors.Errorf("the requested %d registers don't fit into a request, at most %d do", words, maxWireAddress)
	}
	return uint16(words), nil
}

// resolveUnitId is the unit identifier a request for this tag goes out with: the tag's own one if
// it declared one, the connection's default otherwise (plc4j ModbusTcpConnection.getUnitId).
func (m modbusTag) resolveUnitId(defaultUnitId uint8) uint8 {
	if m.UnitId != nil {
		return *m.UnitId
	}
	return defaultUnitId
}

// resolveByteOrder is the payload byte order this tag is read and written with: its own if it
// declared one, the connection's default otherwise (plc4j
// ModbusTcpConnection.getEffectiveByteOrder).
func (m modbusTag) resolveByteOrder(defaultByteOrder ByteOrder) ByteOrder {
	if m.ByteOrder != nil {
		return *m.ByteOrder
	}
	return defaultByteOrder
}

// validateAddressAndQuantity mirrors the range checks plc4j performs while constructing a tag:
// ModbusTag.java rejects a logical address of zero or less and a quantity of zero or less, and the
// per-area factories (ModbusTagCoil.of, ModbusTagHoldingRegister.of,
// ModbusTagExtendedRegister.of) reject an address beyond the 16 bit address space, a range running
// past the end of it and a quantity beyond what a single request can carry - 2000 for the bit
// areas, 125 for the register areas.
// Both arguments are the values as written by the user, i.e. address is the logical address.
func validateAddressAndQuantity(tagType TagType, address uint64, quantity uint64) error {
	// plc4j checks getLogicalAddress() <= 0 in the ModbusTag constructor, which covers the
	// extended register area too even though that one is addressed starting at zero on the wire.
	if address < 1 {
		return errors.Errorf("address must be greater than zero. Was %d", address)
	}
	offset := uint64(logicalAddressOffset(tagType))
	wireAddress := address - offset
	if wireAddress > maxWireAddress {
		return errors.Errorf("address must be less than or equal to %d. Was %d", maxWireAddress+offset, address)
	}
	if quantity == 0 {
		return errors.Errorf("quantity must be greater than zero. Was %d", quantity)
	}
	// plc4j rejects a range whose last address reaches maxWireAddress, and the strict bound keeps
	// the quantity well inside the 16 bit field it is written into further down.
	if wireAddress+quantity > maxWireAddress {
		return errors.Errorf("last requested address is out of range, should be between %d and %d. Was %d",
			offset, maxWireAddress, wireAddress+quantity)
	}
	switch tagType {
	case Coil, DiscreteInput:
		if quantity > maxCoilQuantity {
			return errors.Errorf("quantity may not be larger than %d. Was %d", maxCoilQuantity, quantity)
		}
	default:
		if quantity > maxRegisterQuantity {
			return errors.Errorf("quantity may not be larger than %d. Was %d", maxRegisterQuantity, quantity)
		}
	}
	return nil
}

func NewModbusPlcTagFromStrings(tagType TagType, addressString string, quantityString string, stringLengthString string, datatype readWriteModel.ModbusDataType, config tagConfig, _options ...options.WithOption) (apiModel.PlcTag, error) {
	// Parsed with 32 bits so that an out-of-range address is reported as such instead of as an
	// unparsable string.
	address, err := strconv.ParseUint(addressString, 10, 32)
	if err != nil {
		return nil, errors.Errorf("Couldn't parse address string '%s' into an int", addressString)
	}
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	if quantityString == "" {
		customLogger.Debug().Msg("No quantity supplied, assuming 1")
		quantityString = "1"
	}
	quantity, err := strconv.ParseUint(quantityString, 10, 32)
	if err != nil {
		// A quantity that was spelled out but doesn't parse is a broken address, not a request
		// for a single element.
		return nil, errors.Errorf("Couldn't parse quantity string '%s' into an int", quantityString)
	}
	if err := validateAddressAndQuantity(tagType, address, quantity); err != nil {
		return nil, err
	}
	stringLength, err := validateStringLength(datatype, stringLengthString)
	if err != nil {
		return nil, err
	}
	return newTagFromWireAddress(tagType, uint16(address-uint64(logicalAddressOffset(tagType))), uint16(quantity), datatype, stringLength, config), nil
}

func (m modbusTag) GetAddressString() string {
	dataType := m.Datatype.String()
	// A string's length is part of its address - without it the address wouldn't parse back.
	if isStringType(m.Datatype) {
		dataType = fmt.Sprintf("%s(%d)", dataType, m.StringLength)
	}
	// The logical address is what the user wrote and what the address has to parse back as (plc4j
	// ModbusTag.getAddressString uses getLogicalAddress for the same reason).
	address := fmt.Sprintf("%dx%05d:%s[%d]", m.TagType, uint32(m.Address)+uint32(logicalAddressOffset(m.TagType)), dataType, m.Quantity)
	// Same for the per-tag settings, which are written in curly braces behind the address.
	var config []string
	if m.UnitId != nil {
		config = append(config, fmt.Sprintf("unit-id: %d", *m.UnitId))
	}
	if m.ByteOrder != nil {
		config = append(config, fmt.Sprintf("byte-order: '%s'", m.ByteOrder.String()))
	}
	if len(config) > 0 {
		address += "{" + strings.Join(config, ", ") + "}"
	}
	return address
}

func (m modbusTag) GetValueType() apiValues.PlcValueType {
	if plcValueType, ok := apiValues.PlcValueTypeByName(m.Datatype.String()); !ok {
		return apiValues.NULL
	} else {
		return plcValueType
	}
}

// GetArrayInfo reports the number of elements as a half-open range [0, Quantity).
//
// Note that plc4j's ModbusTag returns an inclusive upper bound (quantity - 1). The Go SPI uses the
// opposite convention: spiModel.DefaultArrayInfo.GetSize is UpperBound - LowerBound, and every
// consumer (e.g. bacnetip's Reader, knxnetip's Subscriber) as well as every other Go driver (ads,
// s7, cbus, simulated) treats the upper bound as exclusive. Modbus follows the Go convention here,
// so a quantity of 5 yields [0, 5) and not [0, 4].
func (m modbusTag) GetArrayInfo() []apiModel.ArrayInfo {
	if m.Quantity != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(m.Quantity),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func castToModbusTagFromPlcTag(plcTag apiModel.PlcTag) (modbusTag, error) {
	if modbusTagVar, ok := plcTag.(modbusTag); ok {
		return modbusTagVar, nil
	}
	return modbusTag{}, errors.New("couldn't cast to ModbusPlcTag")
}

func (m modbusTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m modbusTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.TagType.GetName()); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint16("address", 16, m.Address); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("numberOfElements", 16, m.Quantity); err != nil {
		return err
	}
	dataType := m.Datatype.String()
	if err := writeBuffer.WriteString("dataType", uint32(len([]rune(dataType))*8), dataType, utils.WithEncoding("UTF8")); err != nil {
		return err
	}
	// Only written when the tag actually declared one, the same way plc4j's ModbusTag.serialize
	// leaves the field out for a tag that goes to the connection's default unit.
	if m.UnitId != nil {
		if err := writeBuffer.WriteUint8("unitId", 8, *m.UnitId); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext(m.TagType.GetName()); err != nil {
		return err
	}
	return nil
}
