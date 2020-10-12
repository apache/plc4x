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
package model

import (
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type RoutingIndication struct {
	KNXNetIPMessage
}

// The corresponding interface
type IRoutingIndication interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m RoutingIndication) MsgType() uint16 {
	return 0x0530
}

func (m RoutingIndication) initialize() spi.Message {
	return m
}

func NewRoutingIndication() KNXNetIPMessageInitializer {
	return &RoutingIndication{}
}

func CastIRoutingIndication(structType interface{}) IRoutingIndication {
	castFunc := func(typ interface{}) IRoutingIndication {
		if iRoutingIndication, ok := typ.(IRoutingIndication); ok {
			return iRoutingIndication
		}
		return nil
	}
	return castFunc(structType)
}

func CastRoutingIndication(structType interface{}) RoutingIndication {
	castFunc := func(typ interface{}) RoutingIndication {
		if sRoutingIndication, ok := typ.(RoutingIndication); ok {
			return sRoutingIndication
		}
		return RoutingIndication{}
	}
	return castFunc(structType)
}

func (m RoutingIndication) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	return lengthInBits
}

func (m RoutingIndication) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func RoutingIndicationParse(io spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Create the instance
	return NewRoutingIndication(), nil
}

func (m RoutingIndication) Serialize(io spi.WriteBuffer) {

}
