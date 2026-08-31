//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    public enum ModbusDataType : byte
    {
        BOOL = 1,
        BYTE = 2,
        WORD = 3,
        DWORD = 4,
        LWORD = 5,
        SINT = 6,
        INT = 7,
        DINT = 8,
        LINT = 9,
        USINT = 10,
        UINT = 11,
        UDINT = 12,
        ULINT = 13,
        REAL = 14,
        LREAL = 15,
        TIME = 16,
        LTIME = 17,
        DATE = 18,
        LDATE = 19,
        TIME_OF_DAY = 20,
        LTIME_OF_DAY = 21,
        DATE_AND_TIME = 22,
        LDATE_AND_TIME = 23,
        CHAR = 24,
        WCHAR = 25,
        STRING = 26,
        WSTRING = 27,
    }

    public static class ModbusDataTypeExtensions
    {
        public static byte GetDataTypeSize(this ModbusDataType value) => value switch
        {
            ModbusDataType.BOOL => 2,
            ModbusDataType.BYTE => 2,
            ModbusDataType.WORD => 2,
            ModbusDataType.DWORD => 4,
            ModbusDataType.LWORD => 8,
            ModbusDataType.SINT => 2,
            ModbusDataType.INT => 2,
            ModbusDataType.DINT => 4,
            ModbusDataType.LINT => 8,
            ModbusDataType.USINT => 2,
            ModbusDataType.UINT => 2,
            ModbusDataType.UDINT => 4,
            ModbusDataType.ULINT => 8,
            ModbusDataType.REAL => 4,
            ModbusDataType.LREAL => 8,
            ModbusDataType.TIME => 8,
            ModbusDataType.LTIME => 8,
            ModbusDataType.DATE => 8,
            ModbusDataType.LDATE => 8,
            ModbusDataType.TIME_OF_DAY => 8,
            ModbusDataType.LTIME_OF_DAY => 8,
            ModbusDataType.DATE_AND_TIME => 8,
            ModbusDataType.LDATE_AND_TIME => 8,
            ModbusDataType.CHAR => 1,
            ModbusDataType.WCHAR => 2,
            ModbusDataType.STRING => 1,
            ModbusDataType.WSTRING => 2,
            _ => default,
        };
    }
}
