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
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "strconv"
)

// Constant values.
const BVLC_BACNETTYPE uint8 = 0x81

// The data-structure of this message
type BVLC struct {

}

// The corresponding interface
type IBVLC interface {
    spi.Message
    BvlcFunction() uint8
    Serialize(io spi.WriteBuffer) error
}

type BVLCInitializer interface {
    initialize() spi.Message
}

func BVLCBvlcFunction(m IBVLC) uint8 {
    return m.BvlcFunction()
}


func CastIBVLC(structType interface{}) IBVLC {
    castFunc := func(typ interface{}) IBVLC {
        if iBVLC, ok := typ.(IBVLC); ok {
            return iBVLC
        }
        return nil
    }
    return castFunc(structType)
}

func CastBVLC(structType interface{}) BVLC {
    castFunc := func(typ interface{}) BVLC {
        if sBVLC, ok := typ.(BVLC); ok {
            return sBVLC
        }
        if sBVLC, ok := typ.(*BVLC); ok {
            return *sBVLC
        }
        return BVLC{}
    }
    return castFunc(structType)
}

func (m BVLC) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Const Field (bacnetType)
    lengthInBits += 8

    // Discriminator Field (bvlcFunction)
    lengthInBits += 8

    // Implicit Field (bvlcLength)
    lengthInBits += 16

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits
}

func (m BVLC) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BVLCParse(io *spi.ReadBuffer) (spi.Message, error) {

    // Const Field (bacnetType)
    bacnetType, _bacnetTypeErr := io.ReadUint8(8)
    if _bacnetTypeErr != nil {
        return nil, errors.New("Error parsing 'bacnetType' field " + _bacnetTypeErr.Error())
    }
    if bacnetType != BVLC_BACNETTYPE {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BVLC_BACNETTYPE)) + " but got " + strconv.Itoa(int(bacnetType)))
    }

    // Discriminator Field (bvlcFunction) (Used as input to a switch field)
    bvlcFunction, _bvlcFunctionErr := io.ReadUint8(8)
    if _bvlcFunctionErr != nil {
        return nil, errors.New("Error parsing 'bvlcFunction' field " + _bvlcFunctionErr.Error())
    }

    // Implicit Field (bvlcLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    bvlcLength, _bvlcLengthErr := io.ReadUint16(16)
    if _bvlcLengthErr != nil {
        return nil, errors.New("Error parsing 'bvlcLength' field " + _bvlcLengthErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer BVLCInitializer
    var typeSwitchError error
    switch {
    case bvlcFunction == 0x00:
        initializer, typeSwitchError = BVLCResultParse(io)
    case bvlcFunction == 0x01:
        initializer, typeSwitchError = BVLCWideBroadcastDistributionTableParse(io)
    case bvlcFunction == 0x02:
        initializer, typeSwitchError = BVLCReadBroadcastDistributionTableParse(io)
    case bvlcFunction == 0x03:
        initializer, typeSwitchError = BVLCReadBroadcastDistributionTableAckParse(io)
    case bvlcFunction == 0x04:
        initializer, typeSwitchError = BVLCForwardedNPDUParse(io, bvlcLength)
    case bvlcFunction == 0x05:
        initializer, typeSwitchError = BVLCRegisterForeignDeviceParse(io)
    case bvlcFunction == 0x06:
        initializer, typeSwitchError = BVLCReadForeignDeviceTableParse(io)
    case bvlcFunction == 0x07:
        initializer, typeSwitchError = BVLCReadForeignDeviceTableAckParse(io)
    case bvlcFunction == 0x08:
        initializer, typeSwitchError = BVLCDeleteForeignDeviceTableEntryParse(io)
    case bvlcFunction == 0x09:
        initializer, typeSwitchError = BVLCDistributeBroadcastToNetworkParse(io)
    case bvlcFunction == 0x0A:
        initializer, typeSwitchError = BVLCOriginalUnicastNPDUParse(io, bvlcLength)
    case bvlcFunction == 0x0B:
        initializer, typeSwitchError = BVLCOriginalBroadcastNPDUParse(io, bvlcLength)
    case bvlcFunction == 0x0C:
        initializer, typeSwitchError = BVLCSecureBVLLParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Create the instance
    return initializer.initialize(), nil
}

func BVLCSerialize(io spi.WriteBuffer, m BVLC, i IBVLC, childSerialize func() error) error {

    // Const Field (bacnetType)
    _bacnetTypeErr := io.WriteUint8(8, 0x81)
    if _bacnetTypeErr != nil {
        return errors.New("Error serializing 'bacnetType' field " + _bacnetTypeErr.Error())
    }

    // Discriminator Field (bvlcFunction) (Used as input to a switch field)
    bvlcFunction := uint8(i.BvlcFunction())
    _bvlcFunctionErr := io.WriteUint8(8, (bvlcFunction))
    if _bvlcFunctionErr != nil {
        return errors.New("Error serializing 'bvlcFunction' field " + _bvlcFunctionErr.Error())
    }

    // Implicit Field (bvlcLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    bvlcLength := uint16(uint16(m.LengthInBytes()))
    _bvlcLengthErr := io.WriteUint16(16, (bvlcLength))
    if _bvlcLengthErr != nil {
        return errors.New("Error serializing 'bvlcLength' field " + _bvlcLengthErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}
