use std::io::{Error, ErrorKind, Read, Write};
use crate::{Message, ReadBuffer, WriteBuffer};

// [enum DriverType
//     ['0x01' MODBUS_TCP  ]
//     ['0x02' MODBUS_RTU  ]
//     ['0x03' MODBUS_ASCII]
// ]
#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
pub enum DriverType {
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
    type O = u8;

    fn get_length(&self) -> u32 {
        1
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_u8((*self).into())
    }

    fn _deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
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