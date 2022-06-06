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
pub enum ModbusDataType {
    BOOL, 
    BYTE, 
    WORD, 
    DWORD, 
    LWORD, 
    SINT, 
    INT, 
    DINT, 
    LINT, 
    USINT, 
    UINT, 
    UDINT, 
    ULINT, 
    REAL, 
    LREAL, 
    TIME, 
    LTIME, 
    DATE, 
    LDATE, 
    TIME_OF_DAY, 
    LTIME_OF_DAY, 
    DATE_AND_TIME, 
    LDATE_AND_TIME, 
    CHAR, 
    WCHAR, 
    STRING, 
    WSTRING
}

impl TryFrom<u8> for ModbusDataType {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            1 => Ok(ModbusDataType::BOOL),
            2 => Ok(ModbusDataType::BYTE),
            3 => Ok(ModbusDataType::WORD),
            4 => Ok(ModbusDataType::DWORD),
            5 => Ok(ModbusDataType::LWORD),
            6 => Ok(ModbusDataType::SINT),
            7 => Ok(ModbusDataType::INT),
            8 => Ok(ModbusDataType::DINT),
            9 => Ok(ModbusDataType::LINT),
            10 => Ok(ModbusDataType::USINT),
            11 => Ok(ModbusDataType::UINT),
            12 => Ok(ModbusDataType::UDINT),
            13 => Ok(ModbusDataType::ULINT),
            14 => Ok(ModbusDataType::REAL),
            15 => Ok(ModbusDataType::LREAL),
            16 => Ok(ModbusDataType::TIME),
            17 => Ok(ModbusDataType::LTIME),
            18 => Ok(ModbusDataType::DATE),
            19 => Ok(ModbusDataType::LDATE),
            20 => Ok(ModbusDataType::TIME_OF_DAY),
            21 => Ok(ModbusDataType::LTIME_OF_DAY),
            22 => Ok(ModbusDataType::DATE_AND_TIME),
            23 => Ok(ModbusDataType::LDATE_AND_TIME),
            24 => Ok(ModbusDataType::CHAR),
            25 => Ok(ModbusDataType::WCHAR),
            26 => Ok(ModbusDataType::STRING),
            27 => Ok(ModbusDataType::WSTRING),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for ModbusDataType {
    fn into(self) -> u8 {
        match self {
            BOOL => 1,
            BYTE => 2,
            WORD => 3,
            DWORD => 4,
            LWORD => 5,
            SINT => 6,
            INT => 7,
            DINT => 8,
            LINT => 9,
            USINT => 10,
            UINT => 11,
            UDINT => 12,
            ULINT => 13,
            REAL => 14,
            LREAL => 15,
            TIME => 16,
            LTIME => 17,
            DATE => 18,
            LDATE => 19,
            TIME_OF_DAY => 20,
            LTIME_OF_DAY => 21,
            DATE_AND_TIME => 22,
            LDATE_AND_TIME => 23,
            CHAR => 24,
            WCHAR => 25,
            STRING => 26,
            WSTRING => 27
        }
    }
}

impl Message for ModbusDataType {
    type M = ModbusDataType;
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
        match ModbusDataType::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

pub struct ModbusDataTypeArguments {
    dataTypeSize: u8
}

impl ModbusDataType {

    fn get_arguments(self) -> ModbusDataTypeArguments {
        match self {
            BOOL => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            BYTE => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            WORD => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            DWORD => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            }, 
            LWORD => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            SINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            INT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            DINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            }, 
            LINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            USINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            UINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            UDINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            }, 
            ULINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            REAL => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            }, 
            LREAL => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            TIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            LTIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            DATE => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            LDATE => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            TIME_OF_DAY => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            LTIME_OF_DAY => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            DATE_AND_TIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            LDATE_AND_TIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            }, 
            CHAR => {
                ModbusDataTypeArguments {
                    dataTypeSize: 1
                }
            }, 
            WCHAR => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }, 
            STRING => {
                ModbusDataTypeArguments {
                    dataTypeSize: 1
                }
            }, 
            WSTRING => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }
        }
    }
}
