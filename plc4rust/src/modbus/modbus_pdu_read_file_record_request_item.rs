use std::io::{Read, Write};
use crate::{Message, ReadBuffer, WriteBuffer};

// [type ModbusPDUReadFileRecordRequestItem
//     [simple     uint 8     referenceType]
//     [simple     uint 16    fileNumber   ]
//     [simple     uint 16    recordNumber ]
//     [simple     uint 16    recordLength ]
// ]
#[derive(PartialEq,Eq,Clone,Debug)]
pub struct ModbusPDUReadFileRecordRequestItem {
    reference_type: u8,
    file_number: u16,
    record_number: u16,
    record_length: u16,
}

impl Message for ModbusPDUReadFileRecordRequestItem {
    type M = ModbusPDUReadFileRecordRequestItem;
    type O = u8;

    fn get_length(&self) -> u32 {
        56
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error> {
        let mut size = writer.write_u8(self.reference_type)?;
        size += writer.write_u16(self.file_number)?;
        size += writer.write_u16(self.record_number)?;
        size += writer.write_u16(self.record_length)?;
        Ok(size)
    }

    fn _deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, std::io::Error> {
        let reference_type = reader.read_u8()?;
        let file_number = reader.read_u16()?;
        let record_number = reader.read_u16()?;
        let record_length = reader.read_u16()?;

        Ok(Self::M {
            reference_type,
            file_number,
            record_number,
            record_length,
        })
    }
}