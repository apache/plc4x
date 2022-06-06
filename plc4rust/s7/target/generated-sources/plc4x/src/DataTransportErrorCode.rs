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
// package org.apache.plc4x.rust.s7.readwrite;
use std::io::{Error, ErrorKind, Read, Write};
use plc4rust::{Message, NoOption};
use plc4rust::read_buffer::ReadBuffer;
use plc4rust::write_buffer::WriteBuffer;


#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
pub enum DataTransportErrorCode {
    RESERVED, 
    OK, 
    ACCESS_DENIED, 
    INVALID_ADDRESS, 
    DATA_TYPE_NOT_SUPPORTED, 
    NOT_FOUND
}

impl TryFrom<u8> for DataTransportErrorCode {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x00 => Ok(DataTransportErrorCode::RESERVED),
            0xFF => Ok(DataTransportErrorCode::OK),
            0x03 => Ok(DataTransportErrorCode::ACCESS_DENIED),
            0x05 => Ok(DataTransportErrorCode::INVALID_ADDRESS),
            0x06 => Ok(DataTransportErrorCode::DATA_TYPE_NOT_SUPPORTED),
            0x0A => Ok(DataTransportErrorCode::NOT_FOUND),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for DataTransportErrorCode {
    fn into(self) -> u8 {
        match self {
            RESERVED => 0x00,
            OK => 0xFF,
            ACCESS_DENIED => 0x03,
            INVALID_ADDRESS => 0x05,
            DATA_TYPE_NOT_SUPPORTED => 0x06,
            NOT_FOUND => 0x0A
        }
    }
}

impl Message for DataTransportErrorCode {
    type M = DataTransportErrorCode;
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
        match DataTransportErrorCode::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

