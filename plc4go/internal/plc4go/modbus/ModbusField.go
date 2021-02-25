//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package modbus

import (
	"encoding/xml"
	"errors"
	"fmt"
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"strconv"
)

const (
	MODBUS_PROTOCOL_ADDRESS_OFFSET = 1
)

type ModbusPlcField struct {
	FieldType ModbusFieldType
	Address   uint16
	Quantity  uint16
	Datatype  model2.ModbusDataType
}

func NewModbusPlcField(fieldType ModbusFieldType, address uint16, quantity uint16, datatype model2.ModbusDataType) ModbusPlcField {
	return ModbusPlcField{
		FieldType: fieldType,
		Address:   address - MODBUS_PROTOCOL_ADDRESS_OFFSET,
		Quantity:  quantity,
		Datatype:  datatype,
	}
}

func NewModbusPlcFieldFromStrings(fieldType ModbusFieldType, addressString string, quantityString string, datatype model2.ModbusDataType) (model.PlcField, error) {
	address, err := strconv.Atoi(addressString)
	if err != nil {
		return nil, errors.New("Couldn't parse address string '" + addressString + "' into an int")
	}
	quantity, err := strconv.Atoi(quantityString)
	if err != nil {
		quantity = 1
	}
	return NewModbusPlcField(fieldType, uint16(address), uint16(quantity), datatype), nil
}

func (m ModbusPlcField) GetAddressString() string {
	switch m.FieldType {
	case MODBUS_FIELD_COIL:
		return fmt.Sprintf("0x%05d:%s[%d]", m.Address, m.Datatype.String(), m.Quantity)
	case MODBUS_FIELD_DISCRETE_INPUT:
		return fmt.Sprintf("1x%05d:%s[%d]", m.Address, m.Datatype.String(), m.Quantity)
	case MODBUS_FIELD_INPUT_REGISTER:
		return fmt.Sprintf("3x%05d:%s[%d]", m.Address, m.Datatype.String(), m.Quantity)
	case MODBUS_FIELD_HOLDING_REGISTER:
		return fmt.Sprintf("4x%05d:%s[%d]", m.Address, m.Datatype.String(), m.Quantity)
	case MODBUS_FIELD_EXTENDED_REGISTER:
		return fmt.Sprintf("6x%05d:%s[%d]", m.Address, m.Datatype.String(), m.Quantity)
	}
	return ""
}

func (m ModbusPlcField) GetTypeName() string {
	return m.Datatype.String()
}

func (m ModbusPlcField) GetDataType() model2.ModbusDataType {
	return m.Datatype
}

func (m ModbusPlcField) GetQuantity() uint16 {
	return m.Quantity
}

func CastToModbusFieldFromPlcField(plcField model.PlcField) (ModbusPlcField, error) {
	if modbusField, ok := plcField.(ModbusPlcField); ok {
		return modbusField, nil
	}
	return ModbusPlcField{}, errors.New("couldn't cast to ModbusPlcField")
}

func (m ModbusPlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.Address, xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Quantity, xml.StartElement{Name: xml.Name{Local: "numberOfElements"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Datatype.String(), xml.StartElement{Name: xml.Name{Local: "dataType"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}
