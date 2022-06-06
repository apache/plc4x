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
pub enum SyntaxIdType {
    S7ANY, 
    PBC_ID, 
    ALARM_LOCKFREESET, 
    ALARM_INDSET, 
    ALARM_ACKSET, 
    ALARM_QUERYREQSET, 
    NOTIFY_INDSET, 
    NCK, 
    NCK_METRIC, 
    NCK_INCH, 
    DRIVEESANY, 
    SYM1200, 
    DBREAD
}

impl TryFrom<u8> for SyntaxIdType {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x01 => Ok(SyntaxIdType::S7ANY),
            0x13 => Ok(SyntaxIdType::PBC_ID),
            0x15 => Ok(SyntaxIdType::ALARM_LOCKFREESET),
            0x16 => Ok(SyntaxIdType::ALARM_INDSET),
            0x19 => Ok(SyntaxIdType::ALARM_ACKSET),
            0x1A => Ok(SyntaxIdType::ALARM_QUERYREQSET),
            0x1C => Ok(SyntaxIdType::NOTIFY_INDSET),
            0x82 => Ok(SyntaxIdType::NCK),
            0x83 => Ok(SyntaxIdType::NCK_METRIC),
            0x84 => Ok(SyntaxIdType::NCK_INCH),
            0xA2 => Ok(SyntaxIdType::DRIVEESANY),
            0xB2 => Ok(SyntaxIdType::SYM1200),
            0xB0 => Ok(SyntaxIdType::DBREAD),
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<u8> for SyntaxIdType {
    fn into(self) -> u8 {
        match self {
            S7ANY => 0x01,
            PBC_ID => 0x13,
            ALARM_LOCKFREESET => 0x15,
            ALARM_INDSET => 0x16,
            ALARM_ACKSET => 0x19,
            ALARM_QUERYREQSET => 0x1A,
            NOTIFY_INDSET => 0x1C,
            NCK => 0x82,
            NCK_METRIC => 0x83,
            NCK_INCH => 0x84,
            DRIVEESANY => 0xA2,
            SYM1200 => 0xB2,
            DBREAD => 0xB0
        }
    }
}

impl Message for SyntaxIdType {
    type M = SyntaxIdType;
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
        match SyntaxIdType::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

