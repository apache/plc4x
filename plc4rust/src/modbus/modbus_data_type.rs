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
// package org.apache.plc4x.rust.modbus.readwrite;
use crate::plc4x_enum;
use std::io::{Error, ErrorKind, Read, Write};
use crate::{Message, NoOption, ReadBuffer, WriteBuffer};

pub struct ModbusDataTypeArguments {
    dataTypeSize: u8
}

plc4x_enum!
[enum u8 : ModbusDataType
    [1 => BOOL],
    [2 => BYTE],
    [3 => WORD],
    [4 => DWORD],
    [5 => LWORD],
    [6 => SINT],
    [7 => INT],
    [8 => DINT],
    [9 => LINT],
    [10 => USINT],
    [11 => UINT],
    [12 => UDINT],
    [13 => ULINT],
    [14 => REAL],
    [15 => LREAL],
    [16 => TIME],
    [17 => LTIME],
    [18 => DATE],
    [19 => LDATE],
    [20 => TIME_OF_DAY],
    [21 => LTIME_OF_DAY],
    [22 => DATE_AND_TIME],
    [23 => LDATE_AND_TIME],
    [24 => CHAR],
    [25 => WCHAR],
    [26 => STRING],
    [27 => WSTRING]
];

impl ModbusDataType {

    fn get_arguments(self) -> ModbusDataTypeArguments {
        match self {
            BOOL => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            BYTE => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            WORD => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            DWORD => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            },
            LWORD => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            SINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            INT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            DINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            },
            LINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            USINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            UINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            UDINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            },
            ULINT => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            REAL => {
                ModbusDataTypeArguments {
                    dataTypeSize: 4
                }
            },
            LREAL => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            TIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            LTIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            DATE => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            LDATE => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            TIME_OF_DAY => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            LTIME_OF_DAY => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            DATE_AND_TIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            LDATE_AND_TIME => {
                ModbusDataTypeArguments {
                    dataTypeSize: 8
                }
            },
            CHAR => {
                ModbusDataTypeArguments {
                    dataTypeSize: 1
                }
            },
            WCHAR => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            },
            STRING => {
                ModbusDataTypeArguments {
                    dataTypeSize: 1
                }
            },
            WSTRING => {
                ModbusDataTypeArguments {
                    dataTypeSize: 2
                }
            }
        }
    }
}


#[cfg(test)]
mod modbus_tests {
    use crate::{Endianess, Message, ReadBuffer, WriteBuffer};
    use crate::modbus::modbus_pdu::{ModbusPDUReadCoilsRequest, ModbusPDUReadCoilsResponse, ModbusPDUReadDiscreteInputsRequest, ModbusPDUReadDiscreteInputsResponse};
    use crate::modbus::{DriverType, ModbusADU, ModbusADUOptions, ModbusPDU};
    use crate::modbus::modbus_adu::ModbusTcpADU;
    use crate::modbus::modbus_data_type::{ModbusDataType, ModbusDataTypeArguments};

    #[test]
    fn read_write() {
        let data_type: ModbusDataTypeArguments = ModbusDataType::BOOL.get_arguments();
    }
}