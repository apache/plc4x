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
pub enum COTPTpduSize {
    SIZE_128, 
    SIZE_256, 
    SIZE_512, 
    SIZE_1024, 
    SIZE_2048, 
    SIZE_4096, 
    SIZE_8192
}

impl TryFrom<u8> for COTPTpduSize {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x07 => Ok(COTPTpduSize::SIZE_128),
            0x08 => Ok(COTPTpduSize::SIZE_256),
            0x09 => Ok(COTPTpduSize::SIZE_512),
            0x0a => Ok(COTPTpduSize::SIZE_1024),
            0x0b => Ok(COTPTpduSize::SIZE_2048),
            0x0c => Ok(COTPTpduSize::SIZE_4096),
            0x0d => Ok(COTPTpduSize::SIZE_8192),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for COTPTpduSize {
    fn into(self) -> u8 {
        match self {
            SIZE_128 => 0x07,
            SIZE_256 => 0x08,
            SIZE_512 => 0x09,
            SIZE_1024 => 0x0a,
            SIZE_2048 => 0x0b,
            SIZE_4096 => 0x0c,
            SIZE_8192 => 0x0d
        }
    }
}

impl Message for COTPTpduSize {
    type M = COTPTpduSize;
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
        match COTPTpduSize::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

pub struct COTPTpduSizeArguments {
    sizeInBytes: u16
}

impl COTPTpduSize {

    fn get_arguments(self) -> COTPTpduSizeArguments {
        match self {
            SIZE_128 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 128
                }
            }, 
            SIZE_256 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 256
                }
            }, 
            SIZE_512 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 512
                }
            }, 
            SIZE_1024 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 1024
                }
            }, 
            SIZE_2048 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 2048
                }
            }, 
            SIZE_4096 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 4096
                }
            }, 
            SIZE_8192 => {
                COTPTpduSizeArguments {
                    sizeInBytes: 8192
                }
            }
        }
    }
}
