use std::io::{Error, ErrorKind, Read, Write};
use std::io::ErrorKind::InvalidInput;

use crate::Message;
use crate::read_buffer::ReadBuffer;
use crate::write_buffer::WriteBuffer;

// [type ModbusConstants
//     [const          uint 16     modbusTcpDefaultPort 502]
// ]
const MODBUS_TCP_DEFAULT_PORT: u16 = 502;

// [enum DriverType
//     ['0x01' MODBUS_TCP  ]
//     ['0x02' MODBUS_RTU  ]
//     ['0x03' MODBUS_ASCII]
// ]
#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
enum DriverType {
    MODBUS_TCP,
    MODBUS_RTU,
    MODBUS_ASCII
}

impl TryFrom<u8> for DriverType {
    type Error = ();

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            0x01 => {
                Ok(DriverType::MODBUS_TCP)
            },
            0x02 => {
                Ok(DriverType::MODBUS_RTU)
            },
            0x03 => {
                Ok(DriverType::MODBUS_ASCII)
            }
            _ => {
                Err(())
            }
        }
    }
}

impl Into<u8> for DriverType {
    fn into(self) -> u8 {
        match self {
            DriverType::MODBUS_TCP => {
                0x01
            }
            DriverType::MODBUS_RTU => {
                0x02
            }
            DriverType::MODBUS_ASCII => {
                0x03
            }
        }
    }
}

impl Message for DriverType {
    type M = DriverType;

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_u8((*self).into())
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
        let result = reader.read_u8()?;
        match DriverType::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                Err(Error::new(ErrorKind::InvalidInput, format!("Cannot parse {}", result)))
            }
        }
    }
}

// [discriminatedType ModbusADU(DriverType driverType, bit response) byteOrder='BIG_ENDIAN'
//     [typeSwitch driverType
//         ['MODBUS_TCP' ModbusTcpADU
//             // It is used for transaction pairing, the MODBUS server copies in the response the transaction
//             // identifier of the request.
//             [simple         uint 16     transactionIdentifier]
//
//             // It is used for intra-system multiplexing. The MODBUS protocol is identified by the value 0.
//             [const          uint 16     protocolIdentifier    0x0000]
//             // The length field is a byte count of the following fields, including the Unit Identifier and
//             // data fields.
//             [implicit       uint 16     length                'pdu.lengthInBytes + 1']
//
//             // This field is used for intra-system routing purpose. It is typically used to communicate to
//             // a MODBUS+ or a MODBUS serial line slave through a gateway between an Ethernet TCP-IP network
//             // and a MODBUS serial line. This field is set by the MODBUS Client in the request and must be
//             // returned with the same value in the response by the server.
//             [simple         uint 8      unitIdentifier]
//
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu]
//         ]
//         ['MODBUS_RTU' ModbusRtuADU
//             [simple         uint 8                  address]
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu    ]
//
//             [checksum       uint 16                 crc     'STATIC_CALL("rtuCrcCheck", address, pdu)']
//         ]
//         ['MODBUS_ASCII' ModbusAsciiADU
//             [simple         uint 8                  address]
//
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu    ]
//             [checksum       uint 8                  crc     'STATIC_CALL("asciiLrcCheck", address, pdu)']
//         ]
//     ]
// ]

//             // It is used for transaction pairing, the MODBUS server copies in the response the transaction
//             // identifier of the request.
//             [simple         uint 16     transactionIdentifier]
//
//             // It is used for intra-system multiplexing. The MODBUS protocol is identified by the value 0.
//             [const          uint 16     protocolIdentifier    0x0000]
//             // The length field is a byte count of the following fields, including the Unit Identifier and
//             // data fields.
//             [implicit       uint 16     length                'pdu.lengthInBytes + 1']
//
//             // This field is used for intra-system routing purpose. It is typically used to communicate to
//             // a MODBUS+ or a MODBUS serial line slave through a gateway between an Ethernet TCP-IP network
//             // and a MODBUS serial line. This field is set by the MODBUS Client in the request and must be
//             // returned with the same value in the response by the server.
//             [simple         uint 8      unitIdentifier]
//
//             // The actual modbus payload
//             [simple         ModbusPDU('response')   pdu]
struct ModbusTcpADU {
    transaction_identifier: u16,
    protocol_identifier: u16,

}

struct ModbusADU {
    driver_type: DriverType,
    bit_response: bool
}

impl Message for ModbusADU {
    type M = ModbusADU;

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        todo!()
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
        let driver_type = DriverType::deserialize(reader)?;
        let bit_response = reader.read_bit()?;

        // Now the switch begins
        match (driver_type, bit_response) {
            (_, _) => {
                Err(Error::new(InvalidInput, format!("Unable to deserialize from {:?}, {:?}", driver_type, bit_response)))
            }
        }
    }
}

// [type ModbusPDUReadFileRecordRequestItem
//     [simple     uint 8     referenceType]
//     [simple     uint 16    fileNumber   ]
//     [simple     uint 16    recordNumber ]
//     [simple     uint 16    recordLength ]
// ]
#[derive(PartialEq,Eq,Clone,Debug)]
struct ModbusPDUReadFileRecordRequestItem {
    reference_type: u8,
    file_number: u16,
    record_number: u16,
    record_length: u16,
}

impl Message for ModbusPDUReadFileRecordRequestItem {
    type M = ModbusPDUReadFileRecordRequestItem;

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error> {
        let mut size = writer.write_u8(self.reference_type)?;
        size += writer.write_u16(self.file_number)?;
        size += writer.write_u16(self.record_number)?;
        size += writer.write_u16(self.record_length)?;
        Ok(size)
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, std::io::Error> {
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

// [type ModbusPDUWriteFileRecordResponseItem
//     [simple     uint 8     referenceType]
//     [simple     uint 16    fileNumber]
//     [simple     uint 16    recordNumber]
//     [implicit   uint 16    recordLength   'COUNT(recordData) / 2']
//     [array      byte       recordData     length  'recordLength']
// ]
#[derive(PartialEq,Eq,Clone,Debug)]
struct ModbusPDUWriteFileRecordResponseItem {
    reference_type: u8,
    file_number: u16,
    record_number: u16,
    record_data: Vec<u8>,
}

impl ModbusPDUWriteFileRecordResponseItem {
    fn record_length(&self) -> u16 {
        (self.record_data.len() / 2) as u16
    }
}

impl Message for ModbusPDUWriteFileRecordResponseItem {
    type M = ModbusPDUWriteFileRecordResponseItem;

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error> {
        let mut size = writer.write_u8(self.reference_type)?;
        size += writer.write_u16(self.file_number)?;
        size += writer.write_u16(self.record_number)?;
        size += writer.write_u16(self.record_length())?;
        size += writer.write_bytes(&self.record_data)?;
        Ok(size)
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, std::io::Error> {
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

        if let Ok(msg) = ModbusPDUWriteFileRecordResponseItem::deserialize(&mut reader) {
            assert_eq!(message, msg);
        } else {
            assert!(false);
        }

    }
}
