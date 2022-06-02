use std::any::Any;
use std::collections::HashMap;
use std::io::{Error, Read, Write};
use std::io::ErrorKind::InvalidInput;

use crate::read_buffer::ReadBuffer;
use crate::write_buffer::WriteBuffer;

mod write_buffer;
mod modbus;
mod read_buffer;

#[allow(dead_code)]
pub enum Endianess {
    LittleEndian,
    BigEndian
}

struct NoOption {}

trait Message {
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

        let mut bytes: Vec<u8> = vec![];
        let mut read_buffer = ReadBuffer::new(Endianess::BigEndian, &*bytes);

        let result = ModbusADU::parse(&mut read_buffer, Some(options));
    }

}
