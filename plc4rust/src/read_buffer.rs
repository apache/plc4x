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

use std::io::Read;
use crate::Endianess;

#[allow(dead_code)]
pub struct ReadBuffer<T: Read> {
    position: u64,
    endianness: Endianess,
    reader: T,
}

impl<T: Read> ReadBuffer<T> {
    pub(crate) fn new(endianess: Endianess, reader: T) -> ReadBuffer<T> {
        ReadBuffer {
            position: 0,
            endianness: endianess,
            reader
        }
    }
}

impl<T: Read> ReadBuffer<T> {

    pub(crate) fn read_bit(&self) -> Result<bool, std::io::Error> {
        todo!()
    }

    pub(crate) fn read_u8(&mut self) -> Result<u8, std::io::Error> {
        let mut byte = [0_u8; 1];
        self.reader.read(&mut byte)?;

        Ok(byte[0])
    }

    pub(crate) fn read_u16(&mut self) -> Result<u16, std::io::Error> {
        let mut bytes = [0_u8; 2];
        self.reader.read(&mut bytes)?;

        Ok(match self.endianness {
            Endianess::BigEndian => {
                u16::from_be_bytes(bytes)
            },
            Endianess::LittleEndian => {
                u16::from_le_bytes(bytes)
            },
        })
    }

    pub(crate) fn read_u_n(&self, number_of_bits: u8) -> Result<u64, std::io::Error> {
        todo!()
    }

    pub(crate) fn read_bytes(&mut self, length: usize) -> Result<Vec<u8>, std::io::Error> {
        let mut bytes = vec![0_8; length];
        self.reader.read(&mut bytes)?;

        Ok(bytes)
    }
}
