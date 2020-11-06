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
    log "github.com/sirupsen/logrus"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusSerialADU struct {
    TransactionId uint16
    Length uint16
    Address uint8
    Pdu *ModbusPDU
    IModbusSerialADU
}

// The corresponding interface
type IModbusSerialADU interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

func NewModbusSerialADU(transactionId uint16, length uint16, address uint8, pdu *ModbusPDU) *ModbusSerialADU {
    return &ModbusSerialADU{TransactionId: transactionId, Length: length, Address: address, Pdu: pdu}
}

func CastModbusSerialADU(structType interface{}) ModbusSerialADU {
    castFunc := func(typ interface{}) ModbusSerialADU {
        if casted, ok := typ.(ModbusSerialADU); ok {
            return casted
        }
        if casted, ok := typ.(*ModbusSerialADU); ok {
            return *casted
        }
        return ModbusSerialADU{}
    }
    return castFunc(structType)
}

func (m *ModbusSerialADU) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (transactionId)
    lengthInBits += 16

    // Reserved Field (reserved)
    lengthInBits += 16

    // Simple field (length)
    lengthInBits += 16

    // Simple field (address)
    lengthInBits += 8

    // Simple field (pdu)
    lengthInBits += m.Pdu.LengthInBits()

    return lengthInBits
}

func (m *ModbusSerialADU) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusSerialADUParse(io *utils.ReadBuffer, response bool) (*ModbusSerialADU, error) {

    // Simple Field (transactionId)
    transactionId, _transactionIdErr := io.ReadUint16(16)
    if _transactionIdErr != nil {
        return nil, errors.New("Error parsing 'transactionId' field " + _transactionIdErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint16(16)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint16(0x0000) {
            log.WithFields(log.Fields{
                "expected value": uint16(0x0000),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (length)
    length, _lengthErr := io.ReadUint16(16)
    if _lengthErr != nil {
        return nil, errors.New("Error parsing 'length' field " + _lengthErr.Error())
    }

    // Simple Field (address)
    address, _addressErr := io.ReadUint8(8)
    if _addressErr != nil {
        return nil, errors.New("Error parsing 'address' field " + _addressErr.Error())
    }

    // Simple Field (pdu)
    pdu, _pduErr := ModbusPDUParse(io, response)
    if _pduErr != nil {
        return nil, errors.New("Error parsing 'pdu' field " + _pduErr.Error())
    }

    // Create the instance
    return NewModbusSerialADU(transactionId, length, address, pdu), nil
}

func (m *ModbusSerialADU) Serialize(io utils.WriteBuffer) error {

    // Simple Field (transactionId)
    transactionId := uint16(m.TransactionId)
    _transactionIdErr := io.WriteUint16(16, (transactionId))
    if _transactionIdErr != nil {
        return errors.New("Error serializing 'transactionId' field " + _transactionIdErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint16(16, uint16(0x0000))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (length)
    length := uint16(m.Length)
    _lengthErr := io.WriteUint16(16, (length))
    if _lengthErr != nil {
        return errors.New("Error serializing 'length' field " + _lengthErr.Error())
    }

    // Simple Field (address)
    address := uint8(m.Address)
    _addressErr := io.WriteUint8(8, (address))
    if _addressErr != nil {
        return errors.New("Error serializing 'address' field " + _addressErr.Error())
    }

    // Simple Field (pdu)
    _pduErr := m.Pdu.Serialize(io)
    if _pduErr != nil {
        return errors.New("Error serializing 'pdu' field " + _pduErr.Error())
    }

    return nil
}

func (m *ModbusSerialADU) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "transactionId":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TransactionId = data
            case "length":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Length = data
            case "address":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Address = data
            case "pdu":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUError":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDiscreteInputsRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDiscreteInputsResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadCoilsRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadCoilsResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleCoilRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleCoilResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleCoilsRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleCoilsResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadInputRegistersRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadInputRegistersResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadHoldingRegistersRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadHoldingRegistersResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleRegisterRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleRegisterResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadWriteMultipleHoldingRegistersRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadWriteMultipleHoldingRegistersResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFifoQueueRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFifoQueueResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFileRecordRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFileRecordResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteFileRecordRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteFileRecordResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadExceptionStatusRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadExceptionStatusResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUDiagnosticRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUDiagnosticResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventCounterRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventCounterResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventLogRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventLogResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReportServerIdRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReportServerIdResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDeviceIdentificationRequest":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDeviceIdentificationResponse":
                        var dt *ModbusPDU
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    }
            }
        }
    }
}

func (m *ModbusSerialADU) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusSerialADU"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.TransactionId, xml.StartElement{Name: xml.Name{Local: "transactionId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Length, xml.StartElement{Name: xml.Name{Local: "length"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Address, xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Pdu, xml.StartElement{Name: xml.Name{Local: "pdu"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

