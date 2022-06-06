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

use crate::DataTransportSize::DataTransportSize;
use crate::DataTransportSize::DataTransportSize;

#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
pub enum TransportSize {
    BOOL, 
    BYTE, 
    WORD, 
    DWORD, 
    LWORD, 
    INT, 
    UINT, 
    SINT, 
    USINT, 
    DINT, 
    UDINT, 
    LINT, 
    ULINT, 
    REAL, 
    LREAL, 
    CHAR, 
    WCHAR, 
    STRING, 
    WSTRING, 
    TIME, 
    LTIME, 
    DATE, 
    TIME_OF_DAY, 
    TOD, 
    DATE_AND_TIME, 
    DT
}

impl TryFrom<u8> for TransportSize {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x01 => Ok(TransportSize::BOOL),
            0x02 => Ok(TransportSize::BYTE),
            0x03 => Ok(TransportSize::WORD),
            0x04 => Ok(TransportSize::DWORD),
            0x05 => Ok(TransportSize::LWORD),
            0x06 => Ok(TransportSize::INT),
            0x07 => Ok(TransportSize::UINT),
            0x08 => Ok(TransportSize::SINT),
            0x09 => Ok(TransportSize::USINT),
            0x0A => Ok(TransportSize::DINT),
            0x0B => Ok(TransportSize::UDINT),
            0x0C => Ok(TransportSize::LINT),
            0x0D => Ok(TransportSize::ULINT),
            0x0E => Ok(TransportSize::REAL),
            0x0F => Ok(TransportSize::LREAL),
            0x10 => Ok(TransportSize::CHAR),
            0x11 => Ok(TransportSize::WCHAR),
            0x12 => Ok(TransportSize::STRING),
            0x13 => Ok(TransportSize::WSTRING),
            0x14 => Ok(TransportSize::TIME),
            0x16 => Ok(TransportSize::LTIME),
            0x17 => Ok(TransportSize::DATE),
            0x18 => Ok(TransportSize::TIME_OF_DAY),
            0x19 => Ok(TransportSize::TOD),
            0x1A => Ok(TransportSize::DATE_AND_TIME),
            0x1B => Ok(TransportSize::DT),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for TransportSize {
    fn into(self) -> u8 {
        match self {
            BOOL => 0x01,
            BYTE => 0x02,
            WORD => 0x03,
            DWORD => 0x04,
            LWORD => 0x05,
            INT => 0x06,
            UINT => 0x07,
            SINT => 0x08,
            USINT => 0x09,
            DINT => 0x0A,
            UDINT => 0x0B,
            LINT => 0x0C,
            ULINT => 0x0D,
            REAL => 0x0E,
            LREAL => 0x0F,
            CHAR => 0x10,
            WCHAR => 0x11,
            STRING => 0x12,
            WSTRING => 0x13,
            TIME => 0x14,
            LTIME => 0x16,
            DATE => 0x17,
            TIME_OF_DAY => 0x18,
            TOD => 0x19,
            DATE_AND_TIME => 0x1A,
            DT => 0x1B
        }
    }
}

impl Message for TransportSize {
    type M = TransportSize;
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
        match TransportSize::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

pub struct TransportSizeArguments {
    code: u8, 
    shortName: u8, 
    sizeInBytes: u8, 
    baseType: Option<TransportSize>, 
    dataTransportSize: Option<DataTransportSize>, 
    dataProtocolId: String, 
    supported_S7_300: bool, 
    supported_S7_400: bool, 
    supported_S7_1200: bool, 
    supported_S7_1500: bool, 
    supported_LOGO: bool
}

impl TransportSize {

    fn get_arguments(self) -> TransportSizeArguments {
        match self {
            BOOL => {
                TransportSizeArguments {
                    code: 0x01,
                    shortName: b'X',
                    sizeInBytes: 1,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BIT),
                    dataProtocolId: String::from("IEC61131_BOOL"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            BYTE => {
                TransportSizeArguments {
                    code: 0x02,
                    shortName: b'B',
                    sizeInBytes: 1,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_BYTE"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            WORD => {
                TransportSizeArguments {
                    code: 0x04,
                    shortName: b'W',
                    sizeInBytes: 2,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_WORD"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            DWORD => {
                TransportSizeArguments {
                    code: 0x06,
                    shortName: b'D',
                    sizeInBytes: 4,
                    baseType: Some(TransportSize::WORD),
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_DWORD"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            LWORD => {
                TransportSizeArguments {
                    code: 0x00,
                    shortName: b'X',
                    sizeInBytes: 8,
                    baseType: None,
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_LWORD"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: false,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }, 
            INT => {
                TransportSizeArguments {
                    code: 0x05,
                    shortName: b'W',
                    sizeInBytes: 2,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::INTEGER),
                    dataProtocolId: String::from("IEC61131_INT"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            UINT => {
                TransportSizeArguments {
                    code: 0x05,
                    shortName: b'W',
                    sizeInBytes: 2,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: Some(DataTransportSize::INTEGER),
                    dataProtocolId: String::from("IEC61131_UINT"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            SINT => {
                TransportSizeArguments {
                    code: 0x02,
                    shortName: b'B',
                    sizeInBytes: 1,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_SINT"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            USINT => {
                TransportSizeArguments {
                    code: 0x02,
                    shortName: b'B',
                    sizeInBytes: 1,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_USINT"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            DINT => {
                TransportSizeArguments {
                    code: 0x07,
                    shortName: b'D',
                    sizeInBytes: 4,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: Some(DataTransportSize::INTEGER),
                    dataProtocolId: String::from("IEC61131_DINT"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            UDINT => {
                TransportSizeArguments {
                    code: 0x07,
                    shortName: b'D',
                    sizeInBytes: 4,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: Some(DataTransportSize::INTEGER),
                    dataProtocolId: String::from("IEC61131_UDINT"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            LINT => {
                TransportSizeArguments {
                    code: 0x00,
                    shortName: b'X',
                    sizeInBytes: 8,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_LINT"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: false,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }, 
            ULINT => {
                TransportSizeArguments {
                    code: 0x00,
                    shortName: b'X',
                    sizeInBytes: 16,
                    baseType: Some(TransportSize::INT),
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_ULINT"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: false,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }, 
            REAL => {
                TransportSizeArguments {
                    code: 0x08,
                    shortName: b'D',
                    sizeInBytes: 4,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::REAL),
                    dataProtocolId: String::from("IEC61131_REAL"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            LREAL => {
                TransportSizeArguments {
                    code: 0x30,
                    shortName: b'X',
                    sizeInBytes: 8,
                    baseType: Some(TransportSize::REAL),
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_LREAL"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }, 
            CHAR => {
                TransportSizeArguments {
                    code: 0x03,
                    shortName: b'B',
                    sizeInBytes: 1,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_CHAR"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            WCHAR => {
                TransportSizeArguments {
                    code: 0x13,
                    shortName: b'X',
                    sizeInBytes: 2,
                    baseType: None,
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_WCHAR"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            STRING => {
                TransportSizeArguments {
                    code: 0x03,
                    shortName: b'X',
                    sizeInBytes: 1,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_STRING"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            WSTRING => {
                TransportSizeArguments {
                    code: 0x00,
                    shortName: b'X',
                    sizeInBytes: 2,
                    baseType: None,
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_WSTRING"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            TIME => {
                TransportSizeArguments {
                    code: 0x0B,
                    shortName: b'X',
                    sizeInBytes: 4,
                    baseType: None,
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_TIME"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            LTIME => {
                TransportSizeArguments {
                    code: 0x00,
                    shortName: b'X',
                    sizeInBytes: 8,
                    baseType: Some(TransportSize::TIME),
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_LTIME"),
                    supported_S7_300: false,
                    supported_S7_400: false,
                    supported_S7_1200: false,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }, 
            DATE => {
                TransportSizeArguments {
                    code: 0x09,
                    shortName: b'X',
                    sizeInBytes: 2,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_DATE"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            TIME_OF_DAY => {
                TransportSizeArguments {
                    code: 0x06,
                    shortName: b'X',
                    sizeInBytes: 4,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_TIME_OF_DAY"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            TOD => {
                TransportSizeArguments {
                    code: 0x06,
                    shortName: b'X',
                    sizeInBytes: 4,
                    baseType: None,
                    dataTransportSize: Some(DataTransportSize::BYTE_WORD_DWORD),
                    dataProtocolId: String::from("IEC61131_TIME_OF_DAY"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: true,
                    supported_S7_1500: true,
                    supported_LOGO: true
                }
            }, 
            DATE_AND_TIME => {
                TransportSizeArguments {
                    code: 0x0F,
                    shortName: b'X',
                    sizeInBytes: 12,
                    baseType: None,
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_DATE_AND_TIME"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: false,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }, 
            DT => {
                TransportSizeArguments {
                    code: 0x0F,
                    shortName: b'X',
                    sizeInBytes: 12,
                    baseType: None,
                    dataTransportSize: None,
                    dataProtocolId: String::from("IEC61131_DATE_AND_TIME"),
                    supported_S7_300: true,
                    supported_S7_400: true,
                    supported_S7_1200: false,
                    supported_S7_1500: true,
                    supported_LOGO: false
                }
            }
        }
    }
}
