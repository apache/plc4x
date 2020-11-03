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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "reflect"
)

// The data-structure of this message
type ModbusSerialADU struct {
    TransactionId uint16
    Length uint16
    Address uint8
    Pdu IModbusPDU

}

// The corresponding interface
type IModbusSerialADU interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewModbusSerialADU(transactionId uint16, length uint16, address uint8, pdu IModbusPDU) spi.Message {
    return &ModbusSerialADU{TransactionId: transactionId, Length: length, Address: address, Pdu: pdu}
}

func CastIModbusSerialADU(structType interface{}) IModbusSerialADU {
    castFunc := func(typ interface{}) IModbusSerialADU {
        if iModbusSerialADU, ok := typ.(IModbusSerialADU); ok {
            return iModbusSerialADU
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusSerialADU(structType interface{}) ModbusSerialADU {
    castFunc := func(typ interface{}) ModbusSerialADU {
        if sModbusSerialADU, ok := typ.(ModbusSerialADU); ok {
            return sModbusSerialADU
        }
        if sModbusSerialADU, ok := typ.(*ModbusSerialADU); ok {
            return *sModbusSerialADU
        }
        return ModbusSerialADU{}
    }
    return castFunc(structType)
}

func (m ModbusSerialADU) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

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

func (m ModbusSerialADU) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusSerialADUParse(io *utils.ReadBuffer, response bool) (spi.Message, error) {

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
    _pduMessage, _err := ModbusPDUParse(io, response)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'pdu'. " + _err.Error())
    }
    var pdu IModbusPDU
    pdu, _pduOk := _pduMessage.(IModbusPDU)
    if !_pduOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_pduMessage).Name() + " to IModbusPDU")
    }

    // Create the instance
    return NewModbusSerialADU(transactionId, length, address, pdu), nil
}

func (m ModbusSerialADU) Serialize(io utils.WriteBuffer) error {

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
    pdu := CastIModbusPDU(m.Pdu)
    _pduErr := pdu.Serialize(io)
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
                        var dt *ModbusPDUError
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDiscreteInputsRequest":
                        var dt *ModbusPDUReadDiscreteInputsRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDiscreteInputsResponse":
                        var dt *ModbusPDUReadDiscreteInputsResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadCoilsRequest":
                        var dt *ModbusPDUReadCoilsRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadCoilsResponse":
                        var dt *ModbusPDUReadCoilsResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleCoilRequest":
                        var dt *ModbusPDUWriteSingleCoilRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleCoilResponse":
                        var dt *ModbusPDUWriteSingleCoilResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleCoilsRequest":
                        var dt *ModbusPDUWriteMultipleCoilsRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleCoilsResponse":
                        var dt *ModbusPDUWriteMultipleCoilsResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadInputRegistersRequest":
                        var dt *ModbusPDUReadInputRegistersRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadInputRegistersResponse":
                        var dt *ModbusPDUReadInputRegistersResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadHoldingRegistersRequest":
                        var dt *ModbusPDUReadHoldingRegistersRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadHoldingRegistersResponse":
                        var dt *ModbusPDUReadHoldingRegistersResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleRegisterRequest":
                        var dt *ModbusPDUWriteSingleRegisterRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleRegisterResponse":
                        var dt *ModbusPDUWriteSingleRegisterResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersRequest":
                        var dt *ModbusPDUWriteMultipleHoldingRegistersRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersResponse":
                        var dt *ModbusPDUWriteMultipleHoldingRegistersResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadWriteMultipleHoldingRegistersRequest":
                        var dt *ModbusPDUReadWriteMultipleHoldingRegistersRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadWriteMultipleHoldingRegistersResponse":
                        var dt *ModbusPDUReadWriteMultipleHoldingRegistersResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterRequest":
                        var dt *ModbusPDUMaskWriteHoldingRegisterRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterResponse":
                        var dt *ModbusPDUMaskWriteHoldingRegisterResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFifoQueueRequest":
                        var dt *ModbusPDUReadFifoQueueRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFifoQueueResponse":
                        var dt *ModbusPDUReadFifoQueueResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFileRecordRequest":
                        var dt *ModbusPDUReadFileRecordRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFileRecordResponse":
                        var dt *ModbusPDUReadFileRecordResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteFileRecordRequest":
                        var dt *ModbusPDUWriteFileRecordRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteFileRecordResponse":
                        var dt *ModbusPDUWriteFileRecordResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadExceptionStatusRequest":
                        var dt *ModbusPDUReadExceptionStatusRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadExceptionStatusResponse":
                        var dt *ModbusPDUReadExceptionStatusResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUDiagnosticRequest":
                        var dt *ModbusPDUDiagnosticRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUDiagnosticResponse":
                        var dt *ModbusPDUDiagnosticResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventCounterRequest":
                        var dt *ModbusPDUGetComEventCounterRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventCounterResponse":
                        var dt *ModbusPDUGetComEventCounterResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventLogRequest":
                        var dt *ModbusPDUGetComEventLogRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUGetComEventLogResponse":
                        var dt *ModbusPDUGetComEventLogResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReportServerIdRequest":
                        var dt *ModbusPDUReportServerIdRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReportServerIdResponse":
                        var dt *ModbusPDUReportServerIdResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDeviceIdentificationRequest":
                        var dt *ModbusPDUReadDeviceIdentificationRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    case "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDeviceIdentificationResponse":
                        var dt *ModbusPDUReadDeviceIdentificationResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Pdu = dt
                    }
            }
        }
    }
}

func (m ModbusSerialADU) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

