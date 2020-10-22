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
)

// The data-structure of this message
type KNXGroupAddress struct {

}

// The corresponding interface
type IKNXGroupAddress interface {
    spi.Message
    NumLevels() uint8
    Serialize(io spi.WriteBuffer) error
}

type KNXGroupAddressInitializer interface {
    initialize() spi.Message
}

func KNXGroupAddressNumLevels(m IKNXGroupAddress) uint8 {
    return m.NumLevels()
}


func CastIKNXGroupAddress(structType interface{}) IKNXGroupAddress {
    castFunc := func(typ interface{}) IKNXGroupAddress {
        if iKNXGroupAddress, ok := typ.(IKNXGroupAddress); ok {
            return iKNXGroupAddress
        }
        return nil
    }
    return castFunc(structType)
}

func CastKNXGroupAddress(structType interface{}) KNXGroupAddress {
    castFunc := func(typ interface{}) KNXGroupAddress {
        if sKNXGroupAddress, ok := typ.(KNXGroupAddress); ok {
            return sKNXGroupAddress
        }
        if sKNXGroupAddress, ok := typ.(*KNXGroupAddress); ok {
            return *sKNXGroupAddress
        }
        return KNXGroupAddress{}
    }
    return castFunc(structType)
}

func (m KNXGroupAddress) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits
}

func (m KNXGroupAddress) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KNXGroupAddressParse(io *spi.ReadBuffer, numLevels uint8) (spi.Message, error) {

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer KNXGroupAddressInitializer
    var typeSwitchError error
    switch {
    case numLevels == 1:
        initializer, typeSwitchError = KNXGroupAddressFreeLevelParse(io)
    case numLevels == 2:
        initializer, typeSwitchError = KNXGroupAddress2LevelParse(io)
    case numLevels == 3:
        initializer, typeSwitchError = KNXGroupAddress3LevelParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Create the instance
    return initializer.initialize(), nil
}

func KNXGroupAddressSerialize(io spi.WriteBuffer, m KNXGroupAddress, i IKNXGroupAddress, childSerialize func() error) error {

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}
