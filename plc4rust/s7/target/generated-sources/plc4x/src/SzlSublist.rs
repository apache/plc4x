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
pub enum SzlSublist {
    MODULE_IDENTIFICATION, 
    CPU_FEATURES, 
    USER_MEMORY_AREA, 
    SYSTEM_AREAS, 
    BLOCK_TYPES, 
    STATUS_MODULE_LEDS, 
    COMPONENT_IDENTIFICATION, 
    INTERRUPT_STATUS, 
    ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS, 
    COMMUNICATION_STATUS_DATA, 
    STATUS_SINGLE_MODULE_LED, 
    DP_MASTER_SYSTEM_INFORMATION, 
    MODULE_STATUS_INFORMATION, 
    RACK_OR_STATION_STATUS_INFORMATION, 
    RACK_OR_STATION_STATUS_INFORMATION_2, 
    ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION, 
    MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP, 
    DIAGNOSTIC_BUFFER, 
    MODULE_DIAGNOSTIC_DATA
}

impl TryFrom<u8> for SzlSublist {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x11 => Ok(SzlSublist::MODULE_IDENTIFICATION),
            0x12 => Ok(SzlSublist::CPU_FEATURES),
            0x13 => Ok(SzlSublist::USER_MEMORY_AREA),
            0x14 => Ok(SzlSublist::SYSTEM_AREAS),
            0x15 => Ok(SzlSublist::BLOCK_TYPES),
            0x19 => Ok(SzlSublist::STATUS_MODULE_LEDS),
            0x1C => Ok(SzlSublist::COMPONENT_IDENTIFICATION),
            0x22 => Ok(SzlSublist::INTERRUPT_STATUS),
            0x25 => Ok(SzlSublist::ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS),
            0x32 => Ok(SzlSublist::COMMUNICATION_STATUS_DATA),
            0x74 => Ok(SzlSublist::STATUS_SINGLE_MODULE_LED),
            0x90 => Ok(SzlSublist::DP_MASTER_SYSTEM_INFORMATION),
            0x91 => Ok(SzlSublist::MODULE_STATUS_INFORMATION),
            0x92 => Ok(SzlSublist::RACK_OR_STATION_STATUS_INFORMATION),
            0x94 => Ok(SzlSublist::RACK_OR_STATION_STATUS_INFORMATION_2),
            0x95 => Ok(SzlSublist::ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION),
            0x96 => Ok(SzlSublist::MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP),
            0xA0 => Ok(SzlSublist::DIAGNOSTIC_BUFFER),
            0xB1 => Ok(SzlSublist::MODULE_DIAGNOSTIC_DATA),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for SzlSublist {
    fn into(self) -> u8 {
        match self {
            MODULE_IDENTIFICATION => 0x11,
            CPU_FEATURES => 0x12,
            USER_MEMORY_AREA => 0x13,
            SYSTEM_AREAS => 0x14,
            BLOCK_TYPES => 0x15,
            STATUS_MODULE_LEDS => 0x19,
            COMPONENT_IDENTIFICATION => 0x1C,
            INTERRUPT_STATUS => 0x22,
            ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS => 0x25,
            COMMUNICATION_STATUS_DATA => 0x32,
            STATUS_SINGLE_MODULE_LED => 0x74,
            DP_MASTER_SYSTEM_INFORMATION => 0x90,
            MODULE_STATUS_INFORMATION => 0x91,
            RACK_OR_STATION_STATUS_INFORMATION => 0x92,
            RACK_OR_STATION_STATUS_INFORMATION_2 => 0x94,
            ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION => 0x95,
            MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP => 0x96,
            DIAGNOSTIC_BUFFER => 0xA0,
            MODULE_DIAGNOSTIC_DATA => 0xB1
        }
    }
}

impl Message for SzlSublist {
    type M = SzlSublist;
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
        match SzlSublist::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

