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

use std::io::{Error, ErrorKind, Read, Write};
use crate::{Message, NoOption, ReadBuffer, WriteBuffer};


// [enum DriverType
//     ['0x01' MODBUS_TCP  ]
//     ['0x02' MODBUS_RTU  ]
//     ['0x03' MODBUS_ASCII]
// ]
#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
pub enum DriverType {
    MODBUS_TCP,
    MODBUS_RTU,
    MODBUS_ASCII
}

impl TryFrom<u8> for DriverType {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x01 => {
                Ok(DriverType::MODBUS_TCP)
            },
            0x02 => {
                Ok(DriverType::MODBUS_RTU)
            },
            0x03 => {
                Ok(DriverType::MODBUS_ASCII)
            }
            _ => {
                Err(())
            }
        }
    }
}

impl Into<u8> for DriverType {
    fn into(self) -> u8 {
        match self {
            DriverType::MODBUS_TCP => {
                0x01
            }
            DriverType::MODBUS_RTU => {
                0x02
            }
            DriverType::MODBUS_ASCII => {
                0x03
            }
        }
    }
}

impl Message for DriverType {
    type M = DriverType;
    type P = NoOption;

    fn get_length_in_bits(&self) -> u32 {
        8
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_u8((*self).into())
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        assert!(parameter.is_none());
        let result = reader.read_u8()?;
        match DriverType::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                Err(Error::new(ErrorKind::InvalidInput, format!("Cannot parse {}", result)))
            }
        }
    }
}