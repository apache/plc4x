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
	"math"
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
	// ExplicitRange records whether the address wrote the selection as a range. A one-element
	// range is still a range - [4] is a scalar and [4..4] a list of one - which no count can say.
	ExplicitRange bool
	Datatype      readWriteModel.ModbusDataType
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
	return newTagFromWireAddress(tagType, address-logicalAddressOffset(tagType), quantity, datatype, 1, tagConfig{}, quantity > 1)
}

func newTagFromWireAddress(tagType TagType, wireAddress uint16, quantity uint16, datatype readWriteModel.ModbusDataType, stringLength uint16, config tagConfig, explicitRange bool) modbusTag {
	return modbusTag{
		ExplicitRange: explicitRange,
		TagType:       tagType,
		Address:       wireAddress,
		Quantity:      quantity,
		Datatype:      datatype,
		StringLength:  stringLength,
		UnitId:        config.unitId,
		ByteOrder:     config.byteOrder,
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
// address is the logical address as written by the user and quantity the number of elements the
// user selected; registers is what those elements occupy on the wire, which is what the address
// space and the per-request ceilings are measured in.
func validateAddressAndQuantity(tagType TagType, address uint64, quantity uint64, registers uint64) error {
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
	if wireAddress+registers > maxWireAddress {
		return errors.Errorf("last requested address is out of range, should be between %d and %d. Was %d",
			offset, maxWireAddress, wireAddress+registers)
	}
	switch tagType {
	case Coil, DiscreteInput:
		if registers > maxCoilQuantity {
			return errors.Errorf("quantity may not be larger than %d. Was %d", maxCoilQuantity, registers)
		}
	default:
		if registers > maxRegisterQuantity {
			return errors.Errorf("quantity may not be larger than %d registers. Was %d", maxRegisterQuantity, registers)
		}
	}
	return nil
}

// modbusConstraints: a Modbus read covers one contiguous run of registers or bits, so an address
// selects from a single dimension.
var modbusConstraints = spiModel.SingleDimension

// selectionOf reads the array expression an address carries and returns how far past the written
// address the selection starts, in registers, and how many it spans. An address with no
// expression reads one element where it says.
//
// The offset is consumed into the address here, the way plc4j's per-area of() does: a Modbus
// address is a register number, so "holding-register:1[4..7]" is the same read as
// "holding-register:5[0..3]".
func selectionOf(expression string, address string) (uint64, uint64, bool, error) {
	if expression == "" {
		return 0, 1, false, nil
	}
	dimensions, err := spiModel.ParseArrayExpression(expression, address, modbusConstraints)
	if err != nil {
		return 0, 0, false, err
	}
	dimension := dimensions[0]
	// The third value is not derivable from the others: [4] and [4..4] both select one element,
	// and only the range is an array.
	return uint64(dimension.GetLowerBound() - dimension.GetBase()), uint64(dimension.GetSize()),
		dimension.IsRange(), nil
}

func NewModbusPlcTagFromStrings(tagType TagType, addressString string, arrayExpression string, stringLengthString string, datatype readWriteModel.ModbusDataType, config tagConfig, _options ...options.WithOption) (apiModel.PlcTag, error) {
	// Parsed with 32 bits so that an out-of-range address is reported as such instead of as an
	// unparsable string.
	address, err := strconv.ParseUint(addressString, 10, 32)
	if err != nil {
		return nil, errors.Errorf("Couldn't parse address string '%s' into an int", addressString)
	}
	offset, quantity, explicitRange, err := selectionOf(arrayExpression, addressString+arrayExpression)
	if err != nil {
		return nil, err
	}
	stringLength, err := validateStringLength(datatype, stringLengthString)
	if err != nil {
		return nil, err
	}
	// The offset counts elements; a register address counts registers. They are the same number
	// only for a one-register type, so "holding-register:1[4]:DINT" would otherwise land four
	// registers short of the fifth DINT.
	registerOffset, err := registerOffsetOf(tagType, offset, datatype, stringLength)
	if err != nil {
		return nil, err
	}
	address += registerOffset
	// The wire carries registers, not elements. The address-space and per-request limits are
	// about what a request can hold, so they have to be checked against the register count: 63
	// DINTs are 126 registers and do not fit into the 125 a read carries, however few elements
	// that is.
	registers := registerCountOf(tagType, datatype, quantity, stringLength)
	if err := validateAddressAndQuantity(tagType, address, quantity, registers); err != nil {
		return nil, err
	}
	return newTagFromWireAddress(tagType, uint16(address-uint64(logicalAddressOffset(tagType))), uint16(quantity), datatype, stringLength, config, explicitRange), nil
}

// registerOffsetOf is how far into the area a selection starts, in addresses.
//
// A bit area addresses individual bits, so one element is one address there and nothing is
// scaled. A register area addresses registers, and the conversion goes through the total bit
// offset rather than rounding each element up on its own: the register codec packs elements
// narrower than a register - widthBits reports 8 for a CHAR and 1 for a BOOL - so rounding per
// element would place the start where nothing was written. An offset that does not land on a
// register boundary cannot be addressed by a Modbus read at all, and is reported rather than
// quietly moved to the register before or after it.
func registerOffsetOf(tagType TagType, elementOffset uint64, dataType readWriteModel.ModbusDataType, stringLength uint16) (uint64, error) {
	switch tagType {
	case Coil, DiscreteInput:
		return elementOffset, nil
	}
	bits := elementOffset * widthBits(dataType, stringLength)
	if bits%16 != 0 {
		return 0, errors.Errorf("selection starts %d bits into the address, which is not a register "+
			"boundary - a read starts at a register, so this offset cannot be addressed", bits)
	}
	return bits / 16, nil
}

// registerCountOf is how many addresses the selection occupies on the wire: bits in a bit area,
// where one element is one address, and registers everywhere else. This is the count the request
// carries, and so the one the address-space and per-request limits apply to.
func registerCountOf(tagType TagType, dataType readWriteModel.ModbusDataType, quantity uint64, stringLength uint16) uint64 {
	switch tagType {
	case Coil, DiscreteInput:
		return quantity
	}
	if quantity > math.MaxUint16 {
		// Too large for lengthInBytes to be asked, and far beyond any request. Reported as-is so
		// the quantity limits below reject it rather than a truncated conversion passing.
		return quantity
	}
	registers := (lengthInBytes(dataType, uint16(quantity), stringLength) + 1) / 2
	if registers < 1 {
		// A value narrower than a register still occupies a whole one, and every request has to
		// ask for something.
		return 1
	}
	return registers
}

func (m modbusTag) GetAddressString() string {
	dataType := m.Datatype.String()
	// A string's length is part of its address - without it the address wouldn't parse back.
	if isStringType(m.Datatype) {
		dataType = fmt.Sprintf("%s(%d)", dataType, m.StringLength)
	}
	// The logical address is what the user wrote and what the address has to parse back as (plc4j
	// ModbusTag.getAddressString uses getLogicalAddress for the same reason).
	address := fmt.Sprintf("%dx%05d%s:%s", m.TagType, uint32(m.Address)+uint32(logicalAddressOffset(m.TagType)),
		spiModel.RenderArrayExpression(m.GetArrayInfo()), dataType)
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

// GetArrayInfo reports the shape of the value the caller receives, as an inclusive range: a
// quantity of 5 yields [0..4], the same as plc4j.
//
// The bounds used to be exclusive here, documented as a deliberate divergence from plc4j. It was
// not one worth keeping: once the range is what the user wrote, the bounds are the indices the
// address stated, and an exclusive upper bound reports a number that appears nowhere in it.
//
// The indices are relative to the value the caller receives, not to what was written: a Modbus
// address is a register number, so the driver consumes the start of the selection into the
// address when it resolves it.
func (m modbusTag) GetArrayInfo() []apiModel.ArrayInfo {
	// The flag decides the shape; the count only sizes it.
	if m.ExplicitRange {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(m.Quantity) - 1,
				Range:      true,
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
