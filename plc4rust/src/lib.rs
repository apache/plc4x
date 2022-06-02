use std::io::{Read, Write};

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

trait Message {
    type M;

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error>;
    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, std::io::Error>;

}

#[cfg(test)]
#[allow(unused_must_use)]
mod tests {

}
