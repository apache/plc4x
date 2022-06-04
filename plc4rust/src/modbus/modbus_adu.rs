/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

use std::io::{Error, Read, Write};
use std::io::ErrorKind::InvalidInput;
use crate::{Message, NoOption, ReadBuffer, WriteBuffer};
use crate::modbus::{DriverType, ModbusPDU};
use crate::modbus::modbus_pdu::ModbusPDUOption;

// [discriminatedType ModbusADU(DriverType driverType, bit response) byteOrder='BIG_ENDIAN'
//     [typeSwitch driverType
//         ['MODBUS_TCP' ModbusTcpADU
//             // It is used for transaction pairing, the MODBUS server copies in the response the transaction
//             // identifier of the request.
//             [simple         uint 16     transactionIdentifier]
//
//             // It is used for intra-system multiplexing. The MODBUS protocol is identified by the value 0.
//             [const          uint 16     protocolIdentifier    0x0000]
//
//             // The length field is a byte count of the following fields, including the Unit Identifier and
//             // data fields.
//             [implicit       uint 16     length                'pdu.lengthInBytes + 1']
//
//             // This field is used for intra-system routing purpose. It is typically used to communicate to
//             // a MODBUS+ or a MODBUS serial line slave through a gateway between an Ethernet TCP-IP network
//             // and a MODBUS serial line. This field is set by the MODBUS Client in the request and must be
//             // returned with the same value in the response by the server.
//             [simple         uint 8      unitIdentifier]
//
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu]
//         ]
//         ['MODBUS_RTU' ModbusRtuADU
//             [simple         uint 8                  address]
//
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu    ]
//
//             [checksum       uint 16                 crc     'STATIC_CALL("rtuCrcCheck", address, pdu)']
//         ]
//         ['MODBUS_ASCII' ModbusAsciiADU
//             [simple         uint 8                  address]
//
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu    ]
//
//             [checksum       uint 8                  crc     'STATIC_CALL("asciiLrcCheck", address, pdu)']
//         ]
//     ]
// ]
#[derive(PartialEq, Debug, Clone)]
pub enum ModbusADU {
    ModbusTcpADU(ModbusTcpADU),
    ModbusRtuADU(ModbusRtuADU),
}

pub struct ModbusADUOptions {
    pub driver_type: DriverType,
    pub response: bool
}

impl Message for ModbusADU {
    type M = ModbusADU;
    type P = ModbusADUOptions;

    fn get_length_in_bits(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        match self {
            ModbusADU::ModbusTcpADU(msg) => {
                msg.serialize(writer)
            }
            ModbusADU::ModbusRtuADU(msg) => {
                panic!("Not implemented!")
            }
        }
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        let parameter = parameter.expect("No Options given!");
        match parameter.driver_type {
            DriverType::MODBUS_TCP => {
                Ok(ModbusADU::ModbusTcpADU(ModbusTcpADU::parse::<T>(reader, Some(ModbusTcpADUOptions { response: parameter.response }))?))
            }
            DriverType::MODBUS_RTU => {
                Ok(ModbusADU::ModbusRtuADU(ModbusRtuADU::parse::<T>(reader, None)?))
            }
            _ => {
                panic!("{}", format!("Unable to deserialize from {:?}, {:?}", parameter.driver_type, parameter.response));
            }
        }
    }
}

#[derive(PartialEq, Debug, Clone)]
pub struct ModbusTcpADU {
    pub transaction_identifier: u16,
    pub protocol_identifier: u16,
    pub unit_identifier: u8,
    pub pdu: ModbusPDU
}

pub struct ModbusTcpADUOptions {
    response: bool
}

impl ModbusTcpADU {
    fn length(&self) -> u16 {
        return (self.pdu.get_length_in_bytes() + 1) as u16
    }
}

impl Message for ModbusTcpADU {
    type M = ModbusTcpADU;
    type P = ModbusTcpADUOptions;

    fn get_length_in_bits(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_u16(self.transaction_identifier)?;
        writer.write_u16(self.protocol_identifier)?;
        writer.write_u8(self.unit_identifier)?;
        self.pdu.serialize(writer);
        Ok(0)
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        let transaction_identifier = reader.read_u16()?;
        let protocol_identifier = reader.read_u16()?;
        let unit_identifier = reader.read_u8()?;
        let pdu = ModbusPDU::parse(reader, Some(ModbusPDUOption {
            bit_response: parameter.unwrap().response
        }))?;
        Ok(Self::M {
            transaction_identifier,
            protocol_identifier,
            unit_identifier,
            pdu
        })
    }
}

#[derive(PartialEq, Debug, Clone)]
pub struct ModbusRtuADU {

}

pub struct ModbusRtuADUOptions {
    response: bool
}

impl Message for ModbusRtuADU {
    type M = ModbusRtuADU;
    type P = ModbusRtuADUOptions;

    fn get_length_in_bits(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        todo!()
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        todo!()
    }
}