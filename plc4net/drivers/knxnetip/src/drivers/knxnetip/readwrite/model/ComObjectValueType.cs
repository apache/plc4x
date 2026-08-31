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

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public enum ComObjectValueType : byte
    {
        BIT1 = 0x00,
        BIT2 = 0x01,
        BIT3 = 0x02,
        BIT4 = 0x03,
        BIT5 = 0x04,
        BIT6 = 0x05,
        BIT7 = 0x06,
        BYTE1 = 0x07,
        BYTE2 = 0x08,
        BYTE3 = 0x09,
        BYTE4 = 0x0A,
        BYTE6 = 0x0B,
        BYTE8 = 0x0C,
        BYTE10 = 0x0D,
        BYTE14 = 0x0E,
    }

    public static class ComObjectValueTypeExtensions
    {
        public static byte GetSizeInBytes(this ComObjectValueType value) => value switch
        {
            ComObjectValueType.BIT1 => 1,
            ComObjectValueType.BIT2 => 1,
            ComObjectValueType.BIT3 => 1,
            ComObjectValueType.BIT4 => 1,
            ComObjectValueType.BIT5 => 1,
            ComObjectValueType.BIT6 => 1,
            ComObjectValueType.BIT7 => 1,
            ComObjectValueType.BYTE1 => 1,
            ComObjectValueType.BYTE2 => 2,
            ComObjectValueType.BYTE3 => 3,
            ComObjectValueType.BYTE4 => 4,
            ComObjectValueType.BYTE6 => 6,
            ComObjectValueType.BYTE8 => 8,
            ComObjectValueType.BYTE10 => 10,
            ComObjectValueType.BYTE14 => 14,
            _ => default,
        };
    }
}
