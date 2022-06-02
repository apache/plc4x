use std::io::{Error, Read, Write};
use std::io::ErrorKind::InvalidInput;
use crate::{Message, ReadBuffer, WriteBuffer};
use crate::modbus::{DriverType, ModbusPDUSubtypes};

pub struct ModbusADU {
    driver_type: DriverType,
    bit_response: bool
}

impl ModbusADU {

    fn deserialize_special<T: Read>(reader: &mut ReadBuffer<T>, driver_type: DriverType, bit_response: bool) -> Result<ModbusADUSubtype, Error> {
        match driver_type {
            DriverType::MODBUS_TCP => {
                Ok(ModbusADUSubtype::ModbusTcpADU(ModbusTcpADU::deserialize(reader)?))
            }
            DriverType::MODBUS_RTU => {
                Ok(ModbusADUSubtype::ModbusRtuADU(ModbusRtuADU::deserialize(reader)?))
            }
            _ => {
                Err(Error::new(InvalidInput, format!("Unable to deserialize from {:?}, {:?}", driver_type, bit_response)))
            }
        }
    }

}

enum ModbusADUSubtype {
    ModbusTcpADU(ModbusTcpADU),
    ModbusRtuADU(ModbusRtuADU),
}

impl Message for ModbusADU {
    type M = ModbusADU;

    fn get_length(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        todo!()
    }

    fn deserialize<T: Read>(reader: &mut ReadBuffer<T>) -> Result<Self::M, Error> {
        Err(Error::new(InvalidInput, "Cannot parse directly!"))
    }
}

struct ModbusTcpADU {
    transaction_identifier: u16,
    protocol_identifier: u16,
    unit_identifier: u8,
    pdu: ModbusPDUSubtypes
}

impl ModbusTcpADU {
    fn length(&self) -> u16 {
        return (self.pdu.get_length() + 1) as u16
    }
}

impl Message for ModbusTcpADU {
    type M = ModbusTcpADU;

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

struct ModbusRtuADU {

}

impl Message for ModbusRtuADU {
    type M = ModbusRtuADU;

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