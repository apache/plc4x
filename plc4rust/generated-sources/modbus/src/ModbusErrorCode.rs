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
// package org.apache.plc4x.rust.modbus.readwrite;
use std::io::{Error, ErrorKind, Read, Write};
use plc4rust::{Message, NoOption};
use plc4rust::read_buffer::ReadBuffer;
use plc4rust::write_buffer::WriteBuffer;


#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
pub enum ModbusErrorCode {
    ILLEGAL_FUNCTION, 
    ILLEGAL_DATA_ADDRESS, 
    ILLEGAL_DATA_VALUE, 
    SLAVE_DEVICE_FAILURE, 
    ACKNOWLEDGE, 
    SLAVE_DEVICE_BUSY, 
    NEGATIVE_ACKNOWLEDGE, 
    MEMORY_PARITY_ERROR, 
    GATEWAY_PATH_UNAVAILABLE, 
    GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND
}

impl TryFrom<u8> for ModbusErrorCode {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            1 => Ok(ModbusErrorCode::ILLEGAL_FUNCTION),
            2 => Ok(ModbusErrorCode::ILLEGAL_DATA_ADDRESS),
            3 => Ok(ModbusErrorCode::ILLEGAL_DATA_VALUE),
            4 => Ok(ModbusErrorCode::SLAVE_DEVICE_FAILURE),
            5 => Ok(ModbusErrorCode::ACKNOWLEDGE),
            6 => Ok(ModbusErrorCode::SLAVE_DEVICE_BUSY),
            7 => Ok(ModbusErrorCode::NEGATIVE_ACKNOWLEDGE),
            8 => Ok(ModbusErrorCode::MEMORY_PARITY_ERROR),
            10 => Ok(ModbusErrorCode::GATEWAY_PATH_UNAVAILABLE),
            11 => Ok(ModbusErrorCode::GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for ModbusErrorCode {
    fn into(self) -> u8 {
        match self {
            ILLEGAL_FUNCTION => 1,
            ILLEGAL_DATA_ADDRESS => 2,
            ILLEGAL_DATA_VALUE => 3,
            SLAVE_DEVICE_FAILURE => 4,
            ACKNOWLEDGE => 5,
            SLAVE_DEVICE_BUSY => 6,
            NEGATIVE_ACKNOWLEDGE => 7,
            MEMORY_PARITY_ERROR => 8,
            GATEWAY_PATH_UNAVAILABLE => 10,
            GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND => 11
        }
    }
}

impl Message for ModbusErrorCode {
    type M = ModbusErrorCode;
    type P = NoOption;

    fn get_length_in_bits(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_u8((*self).into())
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        assert!(parameter.is_none());
        let result = reader.read_u8()?;
        match ModbusErrorCode::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

