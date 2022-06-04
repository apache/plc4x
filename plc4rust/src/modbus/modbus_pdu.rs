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

use std::fs::read;
use std::io::{Error, Read, Write};
use std::io::ErrorKind::InvalidInput;
use crate::{Message, NoOption, plc4x_type, ReadBuffer, WriteBuffer};

// [discriminatedType ModbusPDU(bit response)
//     [discriminator bit         errorFlag]
//     [discriminator uint 7      functionFlag]
//     [typeSwitch errorFlag,functionFlag,response
//         ['true'                     ModbusPDUError
//             [simple ModbusErrorCode  exceptionCode]
//         ]
//
//         // Bit Access
//         ['false','0x02','false'     ModbusPDUReadDiscreteInputsRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
//         ['false','0x02','true'      ModbusPDUReadDiscreteInputsResponse
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//
//         ['false','0x01','false'     ModbusPDUReadCoilsRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
//         ['false','0x01','true'      ModbusPDUReadCoilsResponse
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//
//         ['false','0x05','false'     ModbusPDUWriteSingleCoilRequest
//             [simple     uint 16     address]
//             [simple     uint 16     value]
//         ]
//         ['false','0x05','true'      ModbusPDUWriteSingleCoilResponse
//             [simple     uint 16     address]
//             [simple     uint 16     value]
//         ]
//
//         ['false','0x0F','false'     ModbusPDUWriteMultipleCoilsRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//         ['false','0x0F','true'      ModbusPDUWriteMultipleCoilsResponse
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
//
//         // Uint 16 Access (short)
//         ['false','0x04','false'     ModbusPDUReadInputRegistersRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
//         ['false','0x04','true'      ModbusPDUReadInputRegistersResponse
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//
//         ['false','0x03','false'     ModbusPDUReadHoldingRegistersRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
//         ['false','0x03','true'      ModbusPDUReadHoldingRegistersResponse
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//
//         ['false','0x06','false'     ModbusPDUWriteSingleRegisterRequest
//             [simple     uint 16     address]
//             [simple     uint 16     value]
//         ]
//         ['false','0x06','true'      ModbusPDUWriteSingleRegisterResponse
//             [simple     uint 16     address]
//             [simple     uint 16     value]
//         ]
//
//         ['false','0x10','false'     ModbusPDUWriteMultipleHoldingRegistersRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//         ['false','0x10','true'      ModbusPDUWriteMultipleHoldingRegistersResponse
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
//
//         ['false','0x17','false'     ModbusPDUReadWriteMultipleHoldingRegistersRequest
//             [simple     uint 16     readStartingAddress]
//             [simple     uint 16     readQuantity]
//             [simple     uint 16     writeStartingAddress]
//             [simple     uint 16     writeQuantity]
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//         ['false','0x17','true'      ModbusPDUReadWriteMultipleHoldingRegistersResponse
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//
//         ['false','0x16','false'     ModbusPDUMaskWriteHoldingRegisterRequest
//             [simple     uint 16     referenceAddress]
//             [simple     uint 16     andMask]
//             [simple     uint 16     orMask]
//         ]
//         ['false','0x16','true'      ModbusPDUMaskWriteHoldingRegisterResponse
//             [simple     uint 16     referenceAddress]
//             [simple     uint 16     andMask]
//             [simple     uint 16     orMask]
//         ]
//
//         ['false','0x18','false'     ModbusPDUReadFifoQueueRequest
//             [simple     uint 16     fifoPointerAddress]
//         ]
//         ['false','0x18','true'      ModbusPDUReadFifoQueueResponse
//             [implicit   uint 16     byteCount     '(COUNT(fifoValue) * 2) + 2']
//             [implicit   uint 16     fifoCount     '(COUNT(fifoValue) * 2) / 2']
//             [array      uint 16     fifoValue     count   'fifoCount']
//         ]
//
//         // File Record Access
//         ['false','0x14','false'     ModbusPDUReadFileRecordRequest
//             [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
//             [array      ModbusPDUReadFileRecordRequestItem      items length 'byteCount']
//         ]
//         ['false','0x14','true'      ModbusPDUReadFileRecordResponse
//             [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
//             [array      ModbusPDUReadFileRecordResponseItem     items length 'byteCount']
//         ]
//
//         ['false','0x15','false'     ModbusPDUWriteFileRecordRequest
//             [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
//             [array      ModbusPDUWriteFileRecordRequestItem     items length 'byteCount']
//         ]
//         ['false','0x15','true'      ModbusPDUWriteFileRecordResponse
//             [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
//             [array      ModbusPDUWriteFileRecordResponseItem    items length 'byteCount']
//         ]
//
//         // Diagnostics (Serial Line Only)
//         ['false','0x07','false'     ModbusPDUReadExceptionStatusRequest
//         ]
//         ['false','0x07','true'      ModbusPDUReadExceptionStatusResponse
//             [simple     uint 8      value]
//         ]
//
//         ['false','0x08','false'     ModbusPDUDiagnosticRequest
//             [simple     uint 16     subFunction]
//             [simple     uint 16     data]
//         ]
//         ['false','0x08','true'      ModbusPDUDiagnosticResponse
//             [simple     uint 16     subFunction]
//             [simple     uint 16     data]
//         ]
//
//         ['false','0x0B','false'     ModbusPDUGetComEventCounterRequest
//         ]
//         ['false','0x0B','true'      ModbusPDUGetComEventCounterResponse
//             [simple     uint 16     status]
//             [simple     uint 16     eventCount]
//         ]
//
//         ['false','0x0C','false'     ModbusPDUGetComEventLogRequest
//         ]
//         ['false','0x0C','true'      ModbusPDUGetComEventLogResponse
//             [implicit   uint 8      byteCount    'COUNT(events) + 6']
//             [simple     uint 16     status]
//             [simple     uint 16     eventCount]
//             [simple     uint 16     messageCount]
//             [array      byte        events       count   'byteCount - 6']
//         ]
//
//         ['false','0x11','false'     ModbusPDUReportServerIdRequest
//         ]
//         ['false','0x11','true'      ModbusPDUReportServerIdResponse
//             // TODO: This is not specified very well in the spec ... investigate.
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
//
//         // Remark: Even if the Modbus spec states that supporting this type of request is mandatory
//         // I have not come across a single device that really supported it. Some devices just reacted
//         // with an error.
//         ['false','0x2B','false'     ModbusPDUReadDeviceIdentificationRequest
//             [const  uint 8                       meiType  0x0E]
//             [simple ModbusDeviceInformationLevel level        ]
//             [simple uint 8                       objectId     ]
//         ]
//         ['false','0x2B','true'      ModbusPDUReadDeviceIdentificationResponse
//             [const    uint 8                                 meiType          0x0E                              ]
//             [simple   ModbusDeviceInformationLevel           level                                              ]
//             [simple   bit                                    individualAccess                                   ]
//             [simple   ModbusDeviceInformationConformityLevel conformityLevel                                    ]
//             [simple   ModbusDeviceInformationMoreFollows     moreFollows                                        ]
//             [simple   uint 8                                 nextObjectId                                       ]
//             [implicit uint 8                                 numberOfObjects  'COUNT(objects)'                  ]
//             [array    ModbusDeviceInformationObject          objects          count            'numberOfObjects']
//         ]
//     ]
// ]
#[derive(PartialEq, Debug, Clone)]
pub enum ModbusPDU {
    ModbusPDUError(ModbusPDUError),
    ModbusPDUReadDiscreteInputsRequest(ModbusPDUReadDiscreteInputsRequest),
    ModbusPDUReadDiscreteInputsResponse(ModbusPDUReadDiscreteInputsResponse)
}

pub struct ModbusPDUOption {
    pub(crate) bit_response: bool
}

impl Message for ModbusPDU {
    type M = ModbusPDU;
    type P = ModbusPDUOption;

    fn get_length_in_bits(&self) -> u32 {
        match self {
            ModbusPDU::ModbusPDUError(msg) => {
                msg.get_length_in_bits()
            },
            ModbusPDU::ModbusPDUReadDiscreteInputsRequest(msg) => {
                msg.get_length_in_bits()
            }
            ModbusPDU::ModbusPDUReadDiscreteInputsResponse(msg) => {
                msg.get_length_in_bits()
            }
        }
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        // Write general stuff first
        match self {
            ModbusPDU::ModbusPDUError(_) => {
                todo!()
            }
            ModbusPDU::ModbusPDUReadDiscreteInputsRequest(msg) => {
                // Write discriminator
                writer.write_u_n(1, 0);
                writer.write_u_n(7, 0x02);
                msg.serialize(writer)
            }
            ModbusPDU::ModbusPDUReadDiscreteInputsResponse(msg) => {
                // Write discriminator
                writer.write_u_n(1, 0);
                writer.write_u_n(7, 0x02);
                msg.serialize(writer)
            }
        }
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        let response = parameter.expect("No option given!").bit_response;
        let error_flag = reader.read_bit()?;
        let function_flag = reader.read_u_n(7)? as u8;

        match (error_flag, function_flag, response) {
            (true, _, _) => {
                Ok(ModbusPDU::ModbusPDUError(ModbusPDUError {

                }))
            },
            (false, 0x02, false) => {
                Ok(ModbusPDU::ModbusPDUReadDiscreteInputsRequest(ModbusPDUReadDiscreteInputsRequest::parse(reader, None)?))
            }
            (false, 0x02, true) => {
                Ok(ModbusPDU::ModbusPDUReadDiscreteInputsResponse(ModbusPDUReadDiscreteInputsResponse::parse(reader, None)?))
            }
            (_, _, _) => {
                panic!("unnable to parse");
            }
        }
    }
}

#[derive(PartialEq, Debug, Clone)]
pub struct ModbusPDUError {

}

impl Message for ModbusPDUError {
    type M = ModbusPDUError;
    type P = NoOption;

    fn get_length_in_bits(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        todo!()
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        todo!()
    }
}

//         // Bit Access
//         ['false','0x02','false'     ModbusPDUReadDiscreteInputsRequest
//             [simple     uint 16     startingAddress]
//             [simple     uint 16     quantity]
//         ]
plc4x_type!
[type ModbusPDUReadDiscreteInputsRequest
    [simple u16 : startingAddress],
    [simple u16 : quantity]
];

//         ['false','0x02','true'      ModbusPDUReadDiscreteInputsResponse
//             [implicit   uint 8      byteCount     'COUNT(value)']
//             [array      byte        value         count   'byteCount']
//         ]
#[derive(PartialEq, Debug, Clone)]
pub struct ModbusPDUReadDiscreteInputsResponse {
    pub value: Vec<u8>
}

impl ModbusPDUReadDiscreteInputsResponse {
    pub fn byte_count(&self) -> u8 {
        self.value.len() as u8
    }
}

impl Message for ModbusPDUReadDiscreteInputsResponse {
    type M = ModbusPDUReadDiscreteInputsResponse;
    type P = NoOption;

    fn get_length_in_bits(&self) -> u32 {
        (self.value.len() * 8) as u32
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_u8(self.byte_count())?;
        writer.write_bytes(&self.value)
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        let byte_count = reader.read_u8()?;
        let value = reader.read_bytes(byte_count as usize)?;

        Ok(Self::M {
            value
        })
    }
}
