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
	"plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/model"
	"strconv"
)

const (
	MODBUS_PROTOCOL_ADDRESS_OFFSET = 1
)

type ModbusPlcField struct {
	FieldType ModbusFieldType
	Address   uint16
	Quantity  uint16
	Datatype  string
}

func NewModbusPlcField(fieldType ModbusFieldType, address uint16, quantity uint16, datatype string) ModbusPlcField {
	return ModbusPlcField{
		FieldType: fieldType,
		Address:   address - MODBUS_PROTOCOL_ADDRESS_OFFSET,
		Quantity:  quantity,
		Datatype:  datatype,
	}
}

func NewModbusPlcFieldFromStrings(fieldType ModbusFieldType, addressString string, quantityString string, datatype string) (model.PlcField, error) {
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

func (m ModbusPlcField) GetTypeName() string {
	return m.Datatype
}

func CastFromPlcField(plcField model.PlcField) (ModbusPlcField, error) {
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
	if err := e.EncodeElement(m.Datatype, xml.StartElement{Name: xml.Name{Local: "dataType"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}
