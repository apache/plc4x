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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMIRawReq struct {
    CEMI
}

// The corresponding interface
type ICEMIRawReq interface {
    ICEMI
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIRawReq) MessageCode() uint8 {
    return 0x10
}

func (m CEMIRawReq) initialize() spi.Message {
    return m
}

func NewCEMIRawReq() CEMIInitializer {
    return &CEMIRawReq{}
}

func CastICEMIRawReq(structType interface{}) ICEMIRawReq {
    castFunc := func(typ interface{}) ICEMIRawReq {
        if iCEMIRawReq, ok := typ.(ICEMIRawReq); ok {
            return iCEMIRawReq
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIRawReq(structType interface{}) CEMIRawReq {
    castFunc := func(typ interface{}) CEMIRawReq {
        if sCEMIRawReq, ok := typ.(CEMIRawReq); ok {
            return sCEMIRawReq
        }
        if sCEMIRawReq, ok := typ.(*CEMIRawReq); ok {
            return *sCEMIRawReq
        }
        return CEMIRawReq{}
    }
    return castFunc(structType)
}

func (m CEMIRawReq) LengthInBits() uint16 {
    var lengthInBits uint16 = m.CEMI.LengthInBits()

    return lengthInBits
}

func (m CEMIRawReq) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIRawReqParse(io *utils.ReadBuffer) (CEMIInitializer, error) {

    // Create the instance
    return NewCEMIRawReq(), nil
}

func (m CEMIRawReq) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return CEMISerialize(io, m.CEMI, CastICEMI(m), ser)
}
