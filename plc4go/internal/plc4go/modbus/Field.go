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
	"fmt"
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

const (
	AddressOffset = 1
)

type PlcField struct {
	FieldType FieldType
	Address   uint16
	Quantity  uint16
	Datatype  model2.ModbusDataType
}

func NewField(fieldType FieldType, address uint16, quantity uint16, datatype model2.ModbusDataType) PlcField {
	return PlcField{
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

func (m PlcField) GetAddressString() string {
	return fmt.Sprintf("%dx%05d:%s[%d]", m.FieldType, m.Address, m.Datatype.String(), m.Quantity)
}

func (m PlcField) GetTypeName() string {
	return m.Datatype.String()
}

func (m PlcField) GetDataType() model2.ModbusDataType {
	return m.Datatype
}

func (m PlcField) GetQuantity() uint16 {
	return m.Quantity
}

func CastToModbusFieldFromPlcField(plcField model.PlcField) (PlcField, error) {
	if modbusField, ok := plcField.(PlcField); ok {
		return modbusField, nil
	}
	return PlcField{}, errors.New("couldn't cast to ModbusPlcField")
}

func (m PlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
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
