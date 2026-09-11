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

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public enum DataTransportSize : byte
    {
        NULL = 0x00,
        BIT = 0x03,
        BYTE_WORD_DWORD = 0x04,
        INTEGER = 0x05,
        DINTEGER = 0x06,
        REAL = 0x07,
        OCTET_STRING = 0x09,
    }

    public static class DataTransportSizeExtensions
    {
        public static bool GetSizeInBits(this DataTransportSize value) => value switch
        {
            DataTransportSize.NULL => false,
            DataTransportSize.BIT => true,
            DataTransportSize.BYTE_WORD_DWORD => true,
            DataTransportSize.INTEGER => true,
            DataTransportSize.DINTEGER => false,
            DataTransportSize.REAL => false,
            DataTransportSize.OCTET_STRING => false,
            _ => default,
        };
    }
}
