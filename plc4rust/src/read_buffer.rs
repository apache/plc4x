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

use std::fs::read;
use std::io::{Error, ErrorKind, Read};
use std::marker::PhantomData;
use crate::Endianess;

#[allow(dead_code)]
pub struct ReadBuffer<T: Read> {
    position: u64,
    endianness: Endianess,
    bit_reader: BitReader<T>,
    reader: T,
}

impl<T: Read> ReadBuffer<T> {
    pub(crate) fn new(endianess: Endianess, reader: T) -> ReadBuffer<T> {
        ReadBuffer {
            position: 0,
            endianness: endianess,
            bit_reader: BitReader::new(),
            reader
        }
    }
}

impl<T: Read> ReadBuffer<T> {

    pub(crate) fn read_bit(&mut self) -> Result<bool, std::io::Error> {
        Ok(self.read_u_n(1)? > 0)
    }

    pub(crate) fn read_u8(&mut self) -> Result<u8, std::io::Error> {
        let mut byte = [0_u8; 1];
        self.reader.read(&mut byte)?;
        // println!("read_u8: {}", byte[0]);

        Ok(byte[0])
    }

    pub(crate) fn read_u16(&mut self) -> Result<u16, std::io::Error> {
        let mut bytes = [0_u8; 2];
        self.reader.read(&mut bytes)?;

        // println!("read_u16: {}, {} -> {}", bytes[0], bytes[1], u16::from_le_bytes(bytes));

        Ok(match self.endianness {
            Endianess::BigEndian => {
                u16::from_be_bytes(bytes)
            },
            Endianess::LittleEndian => {
                u16::from_le_bytes(bytes)
            },
        })
    }

    pub(crate) fn read_u_n(&mut self, number_of_bits: u8) -> Result<u64, std::io::Error> {
        match self.bit_reader.read(number_of_bits, &mut self.reader) {
            Ok(value) => {
                Ok(value as u64)
            }
            Err(_) => {
                Err(Error::new(ErrorKind::InvalidInput, "Something went wrong"))
            }
        }
    }

    pub(crate) fn read_bytes(&mut self, length: usize) -> Result<Vec<u8>, std::io::Error> {
        let mut bytes = vec![0_8; length];
        self.reader.read(&mut bytes)?;

        Ok(bytes)
    }
}

pub struct BitReader<T: Read> {
    pub(crate) position: u8,
    pub(crate) value: Option<u8>,
    pub(crate) phantom_data: PhantomData<T>
}

impl<T: Read> BitReader<T> {

    fn new() -> BitReader<T> {
        BitReader {
            position: 0,
            value: None,
            phantom_data: PhantomData::default()
        }
    }

    // Writes the given value as the given number of bits to the Bitwriter
    // If it "overflows" the "full" byte is returned
    fn read(&mut self, bits: u8, reader: &mut dyn Read)  -> std::io::Result<u8> {
        assert!(bits <= 8);
        let mut results: u8 = 0;
        let mut bit_count: u8 = 0;
        let bit_offset = self.position;
        loop {
            if bit_count == bits {
                break;
            }
            if self.position == 8 || self.value.is_none(){
                self.fetch(reader);
            }
            let mask = (0x01 << self.position) as u8;
            let bit = self.value.unwrap() & mask;
            // Add the bit to our result
            results = results | (bit >> (self.position - bit_count));

            bit_count += 1;
            self.position += 1;
        }
        Ok(results)
    }

    fn fetch(&mut self, reader: &mut dyn Read) -> std::io::Result<usize> {
        let mut buffer = [0x00_u8; 1];
        let result = reader.read(&mut buffer);

        // println!("Bit Reader Consumed {}", buffer[0]);

        self.position = 0;
        self.value = Some(buffer[0]);
        result
    }

}