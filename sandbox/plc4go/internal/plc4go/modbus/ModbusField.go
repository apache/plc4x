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
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go/model"
	"strconv"
)

type PlcField struct {
	fieldType uint8
	address   uint32
	quantity  uint32
	datatype  string
	model.PlcField
}

func NewModbusPlcField(fieldType uint8, address uint32, quantity uint32, datatype string) PlcField {
	return PlcField{
		fieldType: fieldType,
		address:   address,
		quantity:  quantity,
		datatype:  datatype,
	}
}

func NewModbusPlcFieldFromStrings(fieldType uint8, addressString string, quantityString string, datatype string) (model.PlcField, error) {
	address, err := strconv.Atoi(addressString)
	if err != nil {
		return nil, errors.New("Couldn't parse address string '" + addressString + "' into an int")
	}
	quantity, err := strconv.Atoi(quantityString)
	if err != nil {
		quantity = 1
	}
	return NewModbusPlcField(fieldType, uint32(address), uint32(quantity), datatype), nil
}
