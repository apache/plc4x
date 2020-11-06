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
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDU struct {
    Child IModbusPDUChild
    IModbusPDU
    IModbusPDUParent
}

// The corresponding interface
type IModbusPDU interface {
    ErrorFlag() bool
    FunctionFlag() uint8
    Response() bool
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IModbusPDUParent interface {
    SerializeParent(io utils.WriteBuffer, child IModbusPDU, serializeChildFunction func() error) error
}

type IModbusPDUChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *ModbusPDU)
    IModbusPDU
}

func NewModbusPDU() *ModbusPDU {
    return &ModbusPDU{}
}

func CastModbusPDU(structType interface{}) ModbusPDU {
    castFunc := func(typ interface{}) ModbusPDU {
        if casted, ok := typ.(ModbusPDU); ok {
            return casted
        }
        if casted, ok := typ.(*ModbusPDU); ok {
            return *casted
        }
        return ModbusPDU{}
    }
    return castFunc(structType)
}

func (m *ModbusPDU) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (errorFlag)
    lengthInBits += 1

    // Discriminator Field (functionFlag)
    lengthInBits += 7

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *ModbusPDU) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUParse(io *utils.ReadBuffer, response bool) (*ModbusPDU, error) {

    // Discriminator Field (errorFlag) (Used as input to a switch field)
    errorFlag, _errorFlagErr := io.ReadBit()
    if _errorFlagErr != nil {
        return nil, errors.New("Error parsing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Discriminator Field (functionFlag) (Used as input to a switch field)
    functionFlag, _functionFlagErr := io.ReadUint8(7)
    if _functionFlagErr != nil {
        return nil, errors.New("Error parsing 'functionFlag' field " + _functionFlagErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *ModbusPDU
    var typeSwitchError error
    switch {
    case errorFlag == true:
        _parent, typeSwitchError = ModbusPDUErrorParse(io)
    case errorFlag == false && functionFlag == 0x02 && response == false:
        _parent, typeSwitchError = ModbusPDUReadDiscreteInputsRequestParse(io)
    case errorFlag == false && functionFlag == 0x02 && response == true:
        _parent, typeSwitchError = ModbusPDUReadDiscreteInputsResponseParse(io)
    case errorFlag == false && functionFlag == 0x01 && response == false:
        _parent, typeSwitchError = ModbusPDUReadCoilsRequestParse(io)
    case errorFlag == false && functionFlag == 0x01 && response == true:
        _parent, typeSwitchError = ModbusPDUReadCoilsResponseParse(io)
    case errorFlag == false && functionFlag == 0x05 && response == false:
        _parent, typeSwitchError = ModbusPDUWriteSingleCoilRequestParse(io)
    case errorFlag == false && functionFlag == 0x05 && response == true:
        _parent, typeSwitchError = ModbusPDUWriteSingleCoilResponseParse(io)
    case errorFlag == false && functionFlag == 0x0F && response == false:
        _parent, typeSwitchError = ModbusPDUWriteMultipleCoilsRequestParse(io)
    case errorFlag == false && functionFlag == 0x0F && response == true:
        _parent, typeSwitchError = ModbusPDUWriteMultipleCoilsResponseParse(io)
    case errorFlag == false && functionFlag == 0x04 && response == false:
        _parent, typeSwitchError = ModbusPDUReadInputRegistersRequestParse(io)
    case errorFlag == false && functionFlag == 0x04 && response == true:
        _parent, typeSwitchError = ModbusPDUReadInputRegistersResponseParse(io)
    case errorFlag == false && functionFlag == 0x03 && response == false:
        _parent, typeSwitchError = ModbusPDUReadHoldingRegistersRequestParse(io)
    case errorFlag == false && functionFlag == 0x03 && response == true:
        _parent, typeSwitchError = ModbusPDUReadHoldingRegistersResponseParse(io)
    case errorFlag == false && functionFlag == 0x06 && response == false:
        _parent, typeSwitchError = ModbusPDUWriteSingleRegisterRequestParse(io)
    case errorFlag == false && functionFlag == 0x06 && response == true:
        _parent, typeSwitchError = ModbusPDUWriteSingleRegisterResponseParse(io)
    case errorFlag == false && functionFlag == 0x10 && response == false:
        _parent, typeSwitchError = ModbusPDUWriteMultipleHoldingRegistersRequestParse(io)
    case errorFlag == false && functionFlag == 0x10 && response == true:
        _parent, typeSwitchError = ModbusPDUWriteMultipleHoldingRegistersResponseParse(io)
    case errorFlag == false && functionFlag == 0x17 && response == false:
        _parent, typeSwitchError = ModbusPDUReadWriteMultipleHoldingRegistersRequestParse(io)
    case errorFlag == false && functionFlag == 0x17 && response == true:
        _parent, typeSwitchError = ModbusPDUReadWriteMultipleHoldingRegistersResponseParse(io)
    case errorFlag == false && functionFlag == 0x16 && response == false:
        _parent, typeSwitchError = ModbusPDUMaskWriteHoldingRegisterRequestParse(io)
    case errorFlag == false && functionFlag == 0x16 && response == true:
        _parent, typeSwitchError = ModbusPDUMaskWriteHoldingRegisterResponseParse(io)
    case errorFlag == false && functionFlag == 0x18 && response == false:
        _parent, typeSwitchError = ModbusPDUReadFifoQueueRequestParse(io)
    case errorFlag == false && functionFlag == 0x18 && response == true:
        _parent, typeSwitchError = ModbusPDUReadFifoQueueResponseParse(io)
    case errorFlag == false && functionFlag == 0x14 && response == false:
        _parent, typeSwitchError = ModbusPDUReadFileRecordRequestParse(io)
    case errorFlag == false && functionFlag == 0x14 && response == true:
        _parent, typeSwitchError = ModbusPDUReadFileRecordResponseParse(io)
    case errorFlag == false && functionFlag == 0x15 && response == false:
        _parent, typeSwitchError = ModbusPDUWriteFileRecordRequestParse(io)
    case errorFlag == false && functionFlag == 0x15 && response == true:
        _parent, typeSwitchError = ModbusPDUWriteFileRecordResponseParse(io)
    case errorFlag == false && functionFlag == 0x07 && response == false:
        _parent, typeSwitchError = ModbusPDUReadExceptionStatusRequestParse(io)
    case errorFlag == false && functionFlag == 0x07 && response == true:
        _parent, typeSwitchError = ModbusPDUReadExceptionStatusResponseParse(io)
    case errorFlag == false && functionFlag == 0x08 && response == false:
        _parent, typeSwitchError = ModbusPDUDiagnosticRequestParse(io)
    case errorFlag == false && functionFlag == 0x08 && response == true:
        _parent, typeSwitchError = ModbusPDUDiagnosticResponseParse(io)
    case errorFlag == false && functionFlag == 0x0B && response == false:
        _parent, typeSwitchError = ModbusPDUGetComEventCounterRequestParse(io)
    case errorFlag == false && functionFlag == 0x0B && response == true:
        _parent, typeSwitchError = ModbusPDUGetComEventCounterResponseParse(io)
    case errorFlag == false && functionFlag == 0x0C && response == false:
        _parent, typeSwitchError = ModbusPDUGetComEventLogRequestParse(io)
    case errorFlag == false && functionFlag == 0x0C && response == true:
        _parent, typeSwitchError = ModbusPDUGetComEventLogResponseParse(io)
    case errorFlag == false && functionFlag == 0x11 && response == false:
        _parent, typeSwitchError = ModbusPDUReportServerIdRequestParse(io)
    case errorFlag == false && functionFlag == 0x11 && response == true:
        _parent, typeSwitchError = ModbusPDUReportServerIdResponseParse(io)
    case errorFlag == false && functionFlag == 0x2B && response == false:
        _parent, typeSwitchError = ModbusPDUReadDeviceIdentificationRequestParse(io)
    case errorFlag == false && functionFlag == 0x2B && response == true:
        _parent, typeSwitchError = ModbusPDUReadDeviceIdentificationResponseParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *ModbusPDU) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *ModbusPDU) SerializeParent(io utils.WriteBuffer, child IModbusPDU, serializeChildFunction func() error) error {

    // Discriminator Field (errorFlag) (Used as input to a switch field)
    errorFlag := bool(child.ErrorFlag())
    _errorFlagErr := io.WriteBit((errorFlag))
    if _errorFlagErr != nil {
        return errors.New("Error serializing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Discriminator Field (functionFlag) (Used as input to a switch field)
    functionFlag := uint8(child.FunctionFlag())
    _functionFlagErr := io.WriteUint8(7, (functionFlag))
    if _functionFlagErr != nil {
        return errors.New("Error serializing 'functionFlag' field " + _functionFlagErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *ModbusPDU) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            }
        }
    }
}

func (m *ModbusPDU) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDU"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

