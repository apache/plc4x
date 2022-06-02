use std::io::{Error, Read, Write};
use crate::{Message, ReadBuffer, WriteBuffer};

pub enum ModbusPDUSubtypes {
    ModbusPDUError(ModbusPDUError)
}

impl Message for ModbusPDUSubtypes {
    type M = ModbusPDUSubtypes;

    fn get_length(&self) -> u32 {
        match self {
            ModbusPDUSubtypes::ModbusPDUError(m) => {
                m.get_length()
            }
        }
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        todo!()
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
        todo!()
    }
}

pub struct ModbusPDUError {

}

impl Message for ModbusPDUError {
    type M = ModbusPDUError;

    fn get_length(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        todo!()
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
        todo!()
    }
}