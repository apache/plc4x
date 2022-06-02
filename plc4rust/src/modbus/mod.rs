// [type ModbusConstants
//     [const          uint 16     modbusTcpDefaultPort 502]
// ]
const MODBUS_TCP_DEFAULT_PORT: u16 = 502;

mod driver_type;
mod modbus_adu;
mod modbus_pdu;
mod modbus_pdu_read_file_record_request_item;
mod modbus_pdu_write_file_record_response_item;

pub use driver_type::DriverType;
pub use modbus_adu::ModbusADU;
pub use modbus_adu::ModbusADUOptions;
pub use modbus_pdu::ModbusPDU;
pub use modbus_pdu::ModbusPDUError;
pub use modbus_pdu_read_file_record_request_item::ModbusPDUReadFileRecordRequestItem;
pub use modbus_pdu_write_file_record_response_item::ModbusPDUWriteFileRecordResponseItem;

#[cfg(test)]
#[allow(unused_must_use)]
mod test {
    use crate::{Endianess, Message, ReadBuffer};
    use crate::modbus::ModbusPDUWriteFileRecordResponseItem;
    use crate::write_buffer::{WriteBuffer};

    #[test]
    fn ser_deser() {
        let message = ModbusPDUWriteFileRecordResponseItem {
            reference_type: 0,
            file_number: 0,
            record_number: 0,
            record_data: vec![1, 2, 3, 4],
        };

        let bytes: Vec<u8> = vec![];

        let mut writer = WriteBuffer::new(Endianess::BigEndian, bytes);

        message.serialize(&mut writer);

        let bytes = writer.writer.clone();

        assert_eq!(vec![0, 0, 0, 0, 0, 0, 2, 1, 2, 3, 4], bytes);

        let bytes = writer.writer.clone();
        let mut reader = ReadBuffer::new(Endianess::BigEndian, &*bytes);

        if let Ok(msg) = ModbusPDUWriteFileRecordResponseItem::deserialize_with_parameters(&mut reader, None) {
            assert_eq!(message, msg);
        } else {
            assert!(false);
        }

    }
}
