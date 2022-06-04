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

extern crate core;

use std::io::{Read, Write};

use crate::read_buffer::ReadBuffer;
use crate::write_buffer::WriteBuffer;

pub mod write_buffer;
pub mod modbus;
pub mod read_buffer;
mod r#enum;
mod types;
mod s7;

#[allow(dead_code)]
#[derive(Debug)]
pub enum Endianess {
    LittleEndian,
    BigEndian
}

pub struct NoOption {}

pub trait Message {
    type M;
    type P;

    fn get_length_in_bits(&self) -> u32;

    fn get_length_in_bytes(&self) -> u32 {
        self.get_length_in_bits()/8
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error>;

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, std::io::Error>;

}

#[cfg(test)]
#[allow(unused_must_use)]
mod tests {
    use crate::{Endianess, Message, ReadBuffer};
    use crate::modbus::{DriverType, ModbusADU, ModbusADUOptions};

    #[test]
    fn deserialize_adu() {
        let options = ModbusADUOptions {
            driver_type: DriverType::MODBUS_TCP,
            response: false
        };

        let bytes: Vec<u8> = vec![];
        let mut read_buffer = ReadBuffer::new(Endianess::BigEndian, &*bytes);

        let _ = ModbusADU::parse(&mut read_buffer, Some(options));
    }

}
