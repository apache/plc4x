use std::io::{Read, Write};
use crate::{Message, NoOption, ReadBuffer, WriteBuffer};

// [type ModbusPDUWriteFileRecordResponseItem
//     [simple     uint 8     referenceType]
//     [simple     uint 16    fileNumber]
//     [simple     uint 16    recordNumber]
//     [implicit   uint 16    recordLength   'COUNT(recordData) / 2']
//     [array      byte       recordData     length  'recordLength']
// ]
#[derive(PartialEq,Eq,Clone,Debug)]
pub struct ModbusPDUWriteFileRecordResponseItem {
    pub(crate) reference_type: u8,
    pub(crate) file_number: u16,
    pub(crate) record_number: u16,
    pub(crate) record_data: Vec<u8>,
}

impl ModbusPDUWriteFileRecordResponseItem {
    fn record_length(&self) -> u16 {
        (self.record_data.len() / 2) as u16
    }
}

impl Message for ModbusPDUWriteFileRecordResponseItem {
    type M = ModbusPDUWriteFileRecordResponseItem;
    type P = NoOption;

    fn get_length_in_bits(&self) -> u32 {
        8 * (56 + self.record_data.len()) as u32
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error> {
        let mut size = writer.write_u8(self.reference_type)?;
        size += writer.write_u16(self.file_number)?;
        size += writer.write_u16(self.record_number)?;
        size += writer.write_u16(self.record_length())?;
        size += writer.write_bytes(&self.record_data)?;
        Ok(size)
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, std::io::Error> {
        let reference_type = reader.read_u8()?;
        let file_number = reader.read_u16()?;
        let record_number = reader.read_u16()?;
        let record_length = reader.read_u16()?;
        let record_data = reader.read_bytes(2 * record_length as usize)?;

        Ok(Self::M {
            reference_type,
            file_number,
            record_number,
            record_data,
        })
    }
}