/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#[cfg(test)]
mod modbus_tests {
    use crate::{Endianess, Message, ReadBuffer, WriteBuffer};
    use crate::modbus::modbus_pdu::{ModbusPDUReadDiscreteInputsRequest, ModbusPDUReadDiscreteInputsResponse};
    use crate::modbus::{DriverType, ModbusADU, ModbusADUOptions, ModbusPDU};
    use crate::modbus::modbus_adu::ModbusTcpADU;

    #[test]
    fn read_write() {
        let pdu = ModbusPDU::ModbusPDUReadDiscreteInputsRequest(
            ModbusPDUReadDiscreteInputsRequest {
                startingAddress: 0x03,
                quantity: 1
            }
        );

        let adu = ModbusADU::ModbusTcpADU(ModbusTcpADU {
            transaction_identifier: 1,
            protocol_identifier: 0x0000,
            unit_identifier: 1,
            pdu
        });


        // Send this over a wire
        let mut bytes: Vec<u8> = vec![];
        let mut write_buffer = WriteBuffer::new(Endianess::BigEndian, bytes);
        let result = adu.serialize(&mut write_buffer);

        assert_eq!(vec![0, 1, 0, 0, 1, 4, 0, 3, 0, 1], write_buffer.writer);

        let bytes = write_buffer.writer;

        let mut read_buffer = ReadBuffer::new(Endianess::BigEndian, &*bytes);

        let deserialized = ModbusADU::parse(&mut read_buffer, Some(ModbusADUOptions {
            driver_type: DriverType::MODBUS_TCP,
            response: false
        }));

        assert!(deserialized.is_ok());
        assert_eq!(adu, deserialized.unwrap());
    }

    #[test]
    fn read_write_response() {
        let pdu = ModbusPDU::ModbusPDUReadDiscreteInputsResponse(
            ModbusPDUReadDiscreteInputsResponse {
                value: vec![1]
            }
        );

        let adu = ModbusADU::ModbusTcpADU(ModbusTcpADU {
            transaction_identifier: 1,
            protocol_identifier: 0x0000,
            unit_identifier: 1,
            pdu
        });


        // Send this over a wire
        let mut bytes: Vec<u8> = vec![];
        let mut write_buffer = WriteBuffer::new(Endianess::BigEndian, bytes);
        let result = adu.serialize(&mut write_buffer);

        assert_eq!(vec![0, 1, 0, 0, 1, 4, 1, 1], write_buffer.writer);

        let bytes = write_buffer.writer;

        let mut read_buffer = ReadBuffer::new(Endianess::BigEndian, &*bytes);

        let deserialized = ModbusADU::parse(&mut read_buffer, Some(ModbusADUOptions {
            driver_type: DriverType::MODBUS_TCP,
            response: true
        }));

        assert!(deserialized.is_ok());
        assert_eq!(adu, deserialized.unwrap());
    }

}