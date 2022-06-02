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

trait Message {
    type M;
    type O;

    fn get_length(&self) -> u32;

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error>;

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
        Err(Error::new(InvalidInput, "Cannot parse directly!"))
    }

    fn deserialize_with_parameters<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::O>) -> Result<Self::M, std::io::Error> {
        match parameter {
            None => {
                Self::deserialize(reader)
            }
            Some(_) => {
                Err(Error::new(InvalidInput, "not implemented!"))
            }
        }
    }

}

trait ParametrizedMessage: Message {

}


#[cfg(test)]
#[allow(unused_must_use)]
mod tests {

}
