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

use std::io::Write;
use std::marker::PhantomData;

use crate::Endianess;

pub struct WriteBuffer<T: Write> {
    pub(crate) position: u64,
    pub(crate) endianness: Endianess,
    pub(crate) bit_writer: BitWriter<T>,
    pub(crate) writer: T,
}

pub struct BitWriter<T: Write> {
    pub(crate) position: u8,
    pub(crate) value: u8,
    pub(crate) phantom_data: PhantomData<T>
}

impl<T: Write> BitWriter<T> {

    fn new() -> BitWriter<T> {
        BitWriter {
            position: 0,
            value: 0,
            phantom_data: PhantomData::default()
        }
    }

    // Writes the given value as the given number of bits to the Bitwriter
    // If it "overflows" the "full" byte is returned
    fn write(&mut self, value: u64, bits: u8, writer: &mut dyn Write)  -> std::io::Result<usize> {
        let mut results: usize = 0;
        // Write until the byte is full or bits are over
        let mut bit_index: u8 = 0;
        let bit_offset = self.position;
        loop {
            if self.position == 8 {
                // Flush and then go to 0 again
                results += self.flush(writer)?;
            }
            if bit_index == bits {
                break;
            }
            let mask = (((value >> bit_index) & (0x01)) << self.position) as u8;
            self.value = self.value | mask;

            bit_index += 1;
            self.position += 1;
        }
        Ok(results)
    }

    fn flush(&mut self, writer: &mut dyn Write) -> std::io::Result<usize> {
        let result = writer.write(&[self.value]);
        self.position = 0;
        self.value = 0;
        result
    }

}

#[macro_export]
macro_rules! write_int {
    ($func:ident, $type:ty) => {
        pub fn $func(&mut self, x: $type) -> std::io::Result<usize> {
        let bytes = match self.endianness {
            Endianess::LittleEndian => {
                x.to_le_bytes()
            }
            Endianess::BigEndian => {
                x.to_be_bytes()
            }
        };
        self.write(&bytes)
    }
    };
}

#[allow(dead_code)]
impl<T: Write> WriteBuffer<T> {

    pub fn new(endianess: Endianess, writer: T) -> WriteBuffer<T> {
        WriteBuffer {
            position: 0,
            endianness: endianess,
            bit_writer: BitWriter::new(),
            writer: writer
        }
    }

    fn write(&mut self, bytes: &[u8]) -> std::io::Result<usize> {
        assert!(self.bit_writer.position == 0);
        let bytes_written = self.writer.write(bytes)?;
        self.position = self.position + bytes_written as u64;
        Ok(bytes_written)
    }

    pub fn write_bit(&mut self, value: bool) -> std::io::Result<usize> {
        if value {
            self.write_u_n(1, 1)
        } else {
            self.write_u_n(1, 0)
        }
    }
    pub fn write_u_n(&mut self, num_bits: u8, value: u64) -> std::io::Result<usize> {
        self.bit_writer.write(value, num_bits, &mut self.writer)
    }

    pub fn write_u8(&mut self, x: u8) -> std::io::Result<usize> {
        self.write(&[x])
    }

    write_int!(write_u16, u16);
    write_int!(write_u32, u32);

    write_int!(write_u64, u64);
    write_int!(write_u128, u128);

    write_int!(write_i8, i8);
    write_int!(write_i16, i16);
    write_int!(write_i32, i32);
    write_int!(write_i64, i64);
    write_int!(write_i128, i128);

    write_int!(write_f32, f32);
    write_int!(write_f64, f64);

    pub fn write_bytes(&mut self, bytes: &[u8]) -> std::io::Result<usize> {
        self.write(bytes)
    }
}

#[cfg(test)]
#[allow(unused_must_use)]
mod test {
    use std::io::Write;
    use std::marker::PhantomData;

    use crate::Endianess;
    use crate::write_buffer::{BitWriter, WriteBuffer};

    #[test]
    fn test_it() {
        let mut target: u8 = 0x1;

        let value: u8 = 0x03;
        let mut position: u8 = 1;
        let num_bits = 2;

        for bit_index in 0..num_bits {
            let mask = ((value >> bit_index) & (0x01)) << position;
            target = target | mask;
            position += 1;
        }

        assert_eq!(target, 0x07);
    }

    #[test]
    fn test_write() {
        let mut writer: BitWriter<Vec<u8>> = BitWriter {
            position: 0,
            value: 0,
            phantom_data: PhantomData::default()
        };

        let buffer: Vec<u8> = vec![];

        let mut noop_writer: Box<dyn Write> = Box::new(buffer);
        writer.write(0x01, 1, &mut noop_writer);
        assert_eq!(writer.value, 0x01);
        assert_eq!(writer.position, 1);

        writer.write(0x01, 1, &mut noop_writer);
        assert_eq!(writer.value, 0x03);
        assert_eq!(writer.position, 2);

        writer.write(0x01, 1, &mut noop_writer);
        assert_eq!(writer.value, 0x07);
        assert_eq!(writer.position, 3);

        writer.write(0x03, 2, &mut noop_writer);
        assert_eq!(writer.value, 31);
        assert_eq!(writer.position, 5);

        // Now overflow
        writer.write(0x00, 3, &mut noop_writer);
        assert_eq!(writer.value, 0);
        assert_eq!(writer.position, 0);
    }

    #[test]
    fn test_write_byte() {
        let mut writer: BitWriter<Vec<u8>> = BitWriter {
            position: 0,
            value: 0,
            phantom_data: PhantomData::default()
        };

        let mut bytes: Vec<u8> = vec![];

        // Now overflow
        writer.write(0xFF, 8, &mut bytes);
        assert_eq!(writer.value, 0);
        assert_eq!(writer.position, 0);
        assert_eq!(*bytes.get(0).unwrap(), 0xFF);
    }

    #[test]
    fn write_bit_via_writer() {
        let bytes: Vec<u8> = vec![];

        let mut writer = WriteBuffer::new(Endianess::BigEndian, bytes);

        &writer.write_u_n(9, 0xFFFF);
        assert_eq!(writer.bit_writer.position, 1);
        assert_eq!(writer.bit_writer.value, 0x01);

        let bytes = writer.writer;

        assert_eq!(*bytes.get(0).unwrap(), 0xFF);
        assert_eq!(bytes.get(1), None);
    }
}
