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
type KnxNetRemoteLogging struct {
    Version uint8
    ServiceId
}

// The corresponding interface
type IKnxNetRemoteLogging interface {
    IServiceId
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m KnxNetRemoteLogging) ServiceType() uint8 {
    return 0x06
}

func (m KnxNetRemoteLogging) initialize() spi.Message {
    return m
}

func NewKnxNetRemoteLogging(version uint8) ServiceIdInitializer {
    return &KnxNetRemoteLogging{Version: version}
}

func CastIKnxNetRemoteLogging(structType interface{}) IKnxNetRemoteLogging {
    castFunc := func(typ interface{}) IKnxNetRemoteLogging {
        if iKnxNetRemoteLogging, ok := typ.(IKnxNetRemoteLogging); ok {
            return iKnxNetRemoteLogging
        }
        return nil
    }
    return castFunc(structType)
}

func CastKnxNetRemoteLogging(structType interface{}) KnxNetRemoteLogging {
    castFunc := func(typ interface{}) KnxNetRemoteLogging {
        if sKnxNetRemoteLogging, ok := typ.(KnxNetRemoteLogging); ok {
            return sKnxNetRemoteLogging
        }
        if sKnxNetRemoteLogging, ok := typ.(*KnxNetRemoteLogging); ok {
            return *sKnxNetRemoteLogging
        }
        return KnxNetRemoteLogging{}
    }
    return castFunc(structType)
}

func (m KnxNetRemoteLogging) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ServiceId.LengthInBits()

    // Simple field (version)
    lengthInBits += 8

    return lengthInBits
}

func (m KnxNetRemoteLogging) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxNetRemoteLoggingParse(io *spi.ReadBuffer) (ServiceIdInitializer, error) {

    // Simple Field (version)
    version, _versionErr := io.ReadUint8(8)
    if _versionErr != nil {
        return nil, errors.New("Error parsing 'version' field " + _versionErr.Error())
    }

    // Create the instance
    return NewKnxNetRemoteLogging(version), nil
}

func (m KnxNetRemoteLogging) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Simple Field (version)
    version := uint8(m.Version)
    _versionErr := io.WriteUint8(8, (version))
    if _versionErr != nil {
        return errors.New("Error serializing 'version' field " + _versionErr.Error())
    }

        return nil
    }
    return ServiceIdSerialize(io, m.ServiceId, CastIServiceId(m), ser)
}
