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
	"github.com/apache/plc4x/plc4go/spi/options"
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

const (
	AddressOffset = 1
)

type modbusTag struct {
	apiModel.PlcTag

	TagType  TagType
	Address  uint16
	Quantity uint16
	Datatype readWriteModel.ModbusDataType
}

func NewTag(tagType TagType, address uint16, quantity uint16, datatype readWriteModel.ModbusDataType) modbusTag {
	return modbusTag{
		TagType:  tagType,
		Address:  address - AddressOffset,
		Quantity: quantity,
		Datatype: datatype,
	}
}

func NewModbusPlcTagFromStrings(tagType TagType, addressString string, quantityString string, datatype readWriteModel.ModbusDataType, _options ...options.WithOption) (apiModel.PlcTag, error) {
	address, err := strconv.ParseUint(addressString, 10, 16)
	if err != nil {
		return nil, errors.Errorf("Couldn't parse address string '%s' into an int", addressString)
	}
	localLogger := options.ExtractCustomLogger(_options...)
	if quantityString == "" {
		localLogger.Debug().Msg("No quantity supplied, assuming 1")
		quantityString = "1"
	}
	quantity, err := strconv.ParseUint(quantityString, 10, 16)
	if err != nil {
		localLogger.Warn().Err(err).Msgf("Error during parsing for %s. Falling back to 1", quantityString)
		quantity = 1
	}
	return NewTag(tagType, uint16(address), uint16(quantity), datatype), nil
}

func (m modbusTag) GetAddressString() string {
	return fmt.Sprintf("%dx%05d:%s[%d]", m.TagType, m.Address, m.Datatype.String(), m.Quantity)
}

func (m modbusTag) GetValueType() apiValues.PlcValueType {
	if plcValueType, ok := apiValues.PlcValueByName(m.Datatype.String()); !ok {
		return apiValues.NULL
	} else {
		return plcValueType
	}
}

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

func CastToModbusTagFromPlcTag(plcTag apiModel.PlcTag) (modbusTag, error) {
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
	if err := writeBuffer.WriteString("dataType", uint32(len([]rune(dataType))*8), "UTF-8", dataType); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.TagType.GetName()); err != nil {
		return err
	}
	return nil
}
