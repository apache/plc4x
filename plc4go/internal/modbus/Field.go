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
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	model2 "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	model3 "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	AddressOffset = 1
)

type modbusField struct {
	model.PlcField

	FieldType FieldType
	Address   uint16
	Quantity  uint16
	Datatype  model2.ModbusDataType
}

func NewField(fieldType FieldType, address uint16, quantity uint16, datatype model2.ModbusDataType) modbusField {
	return modbusField{
		FieldType: fieldType,
		Address:   address - AddressOffset,
		Quantity:  quantity,
		Datatype:  datatype,
	}
}

func NewModbusPlcFieldFromStrings(fieldType FieldType, addressString string, quantityString string, datatype model2.ModbusDataType) (model.PlcField, error) {
	address, err := strconv.ParseUint(addressString, 10, 16)
	if err != nil {
		return nil, errors.Errorf("Couldn't parse address string '%s' into an int", addressString)
	}
	if quantityString == "" {
		log.Debug().Msg("No quantity supplied, assuming 1")
		quantityString = "1"
	}
	quantity, err := strconv.ParseUint(quantityString, 10, 16)
	if err != nil {
		log.Warn().Err(err).Msgf("Error during parsing for %s. Falling back to 1", quantityString)
		quantity = 1
	}
	return NewField(fieldType, uint16(address), uint16(quantity), datatype), nil
}

func (m modbusField) GetAddressString() string {
	return fmt.Sprintf("%dx%05d:%s[%d]", m.FieldType, m.Address, m.Datatype.String(), m.Quantity)
}

func (m modbusField) GetValueType() values.PlcValueType {
	if plcValueType, ok := values.PlcValueByName(m.Datatype.String()); !ok {
		return values.NULL
	} else {
		return plcValueType
	}
}

func (m modbusField) GetArrayInfo() []model.ArrayInfo {
	if m.Quantity != 1 {
		return []model.ArrayInfo{
			model3.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(m.Quantity),
			},
		}
	}
	return []model.ArrayInfo{}
}

func CastToModbusFieldFromPlcField(plcField model.PlcField) (modbusField, error) {
	if modbusField, ok := plcField.(modbusField); ok {
		return modbusField, nil
	}
	return modbusField{}, errors.New("couldn't cast to ModbusPlcField")
}

func (m modbusField) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m modbusField) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.FieldType.GetName()); err != nil {
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

	if err := writeBuffer.PopContext(m.FieldType.GetName()); err != nil {
		return err
	}
	return nil
}
