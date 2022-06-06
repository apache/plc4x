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
pub enum MemoryArea {
    COUNTERS, 
    TIMERS, 
    DIRECT_PERIPHERAL_ACCESS, 
    INPUTS, 
    OUTPUTS, 
    FLAGS_MARKERS, 
    DATA_BLOCKS, 
    INSTANCE_DATA_BLOCKS, 
    LOCAL_DATA
}

impl TryFrom<u8> for MemoryArea {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x1C => Ok(MemoryArea::COUNTERS),
            0x1D => Ok(MemoryArea::TIMERS),
            0x80 => Ok(MemoryArea::DIRECT_PERIPHERAL_ACCESS),
            0x81 => Ok(MemoryArea::INPUTS),
            0x82 => Ok(MemoryArea::OUTPUTS),
            0x83 => Ok(MemoryArea::FLAGS_MARKERS),
            0x84 => Ok(MemoryArea::DATA_BLOCKS),
            0x85 => Ok(MemoryArea::INSTANCE_DATA_BLOCKS),
            0x86 => Ok(MemoryArea::LOCAL_DATA),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for MemoryArea {
    fn into(self) -> u8 {
        match self {
            COUNTERS => 0x1C,
            TIMERS => 0x1D,
            DIRECT_PERIPHERAL_ACCESS => 0x80,
            INPUTS => 0x81,
            OUTPUTS => 0x82,
            FLAGS_MARKERS => 0x83,
            DATA_BLOCKS => 0x84,
            INSTANCE_DATA_BLOCKS => 0x85,
            LOCAL_DATA => 0x86
        }
    }
}

impl Message for MemoryArea {
    type M = MemoryArea;
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
        match MemoryArea::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

pub struct MemoryAreaArguments {
    shortName: String
}

impl MemoryArea {

    fn get_arguments(self) -> MemoryAreaArguments {
        match self {
            COUNTERS => {
                MemoryAreaArguments {
                    shortName: String::from("C")
                }
            }, 
            TIMERS => {
                MemoryAreaArguments {
                    shortName: String::from("T")
                }
            }, 
            DIRECT_PERIPHERAL_ACCESS => {
                MemoryAreaArguments {
                    shortName: String::from("D")
                }
            }, 
            INPUTS => {
                MemoryAreaArguments {
                    shortName: String::from("I")
                }
            }, 
            OUTPUTS => {
                MemoryAreaArguments {
                    shortName: String::from("Q")
                }
            }, 
            FLAGS_MARKERS => {
                MemoryAreaArguments {
                    shortName: String::from("M")
                }
            }, 
            DATA_BLOCKS => {
                MemoryAreaArguments {
                    shortName: String::from("DB")
                }
            }, 
            INSTANCE_DATA_BLOCKS => {
                MemoryAreaArguments {
                    shortName: String::from("DBI")
                }
            }, 
            LOCAL_DATA => {
                MemoryAreaArguments {
                    shortName: String::from("LD")
                }
            }
        }
    }
}
